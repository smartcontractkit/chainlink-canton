// package example contains the dummy party-onboarding ceremony workflow.
//
// # Overview
//
// This file is an annotated reference implementation of an *async, multi-participant*
// party-onboarding ceremony built on top of the Operations framework.
//
//  1. Member init     – every participant generates mock signing-key material.
//  2. Propose         – the coordinator assembles a topology proposal from member
//     data and "submits" it to Canton (mocked here).
//  3. Sign            – each required signer independently signs the proposal.
//  4. Submit          – once the signature threshold is met the coordinator
//     aggregates the signatures and finalises the ceremony.
//
// # Async / Resume Pattern
//
// The ceremony is *asynchronous*.  Different operators run the tool at
// different times.  Each `ExecuteOperation` call is **idempotent**: if a
// successful report with the same (definition, input) hash already exists in
// the Reporter the result is returned from the cache and the side-effect is
// *not* repeated.  This means the sequence can be called many times (by
// different actors) and will only advance to the next action that is still
// pending for the current actor.
//
// When the signature threshold has not yet been reached the sequence returns a
// sentinel error (`ErrThresholdNotMet`).  The caller is expected to re-run the
// sequence later (after more signers have completed their signing step).
package example

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [OnboardingSequence] when the number of
// collected signatures have not yet reached the required threshold.  Callers
// should treat this as a transient "come back later" signal.
var ErrThresholdNotMet = errors.New("signature threshold not met: more signers must resume the sequence")

// OnboardingInput is the top-level input to [OnboardingSequence].
type OnboardingInput struct {
	NamespaceName  string   `json:"namespace_name"`
	PartyName      string   `json:"party_name"`
	Participants   []string `json:"participants"`
	SynchronizerID string   `json:"synchronizer_id"`
	// Threshold overrides strict-majority when > 0.
	Threshold int `json:"threshold"`
}

// Phase represents the current execution phase of the onboarding ceremony.
type Phase string

const (
	PhaseInit      Phase = "init"
	PhaseProposal  Phase = "proposal"
	PhaseSigning   Phase = "signing"
	PhaseSubmit    Phase = "submit"
	PhaseCompleted Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every OnboardingOutput.
// It is built progressively as the sequence advances, so it is always present
type CeremonyState struct {
	Phase              Phase    `json:"phase"`
	InitializedMembers []string `json:"initialized_members"`
	RequiredSigners    []string `json:"required_signers,omitempty"`
	CollectedSigners   []string `json:"collected_signers"`
	PendingSigners     []string `json:"pending_signers"`
	Threshold          int      `json:"threshold"`
	ProposalHash       string   `json:"proposal_hash,omitempty"`
}

// OnboardingOutput is the result of [OnboardingSequence].
// State is always populated — even when ExecuteSequence returns an error —
// making it the primary way to inspect ceremony progress without any
// post-hoc report analysis.
type OnboardingOutput struct {
	PartyID      string        `json:"party_id"`
	DNSConfirmed bool          `json:"dns_confirmed"`
	P2PConfirmed bool          `json:"p2p_confirmed"`
	State        CeremonyState `json:"state"`
}

// OnboardingSequence orchestrates the full four-step ceremony.  It is designed
// to be called multiple times by different actors:
//
//   - On each call the sequence re-enters the state machine.  Operations whose
//     (definition, input) hash already has a successful report are skipped
//     instantly — they reflect work done by previous actors.
//   - If the signature threshold is not yet met [ErrThresholdNotMet] is
//     returned so the caller knows to retry after more signers have acted.
//   - Once the threshold is met the coordinator's SubmitProposalOp runs (or is
//     retrieved from cache if already completed).
var OnboardingSequence = operations.NewSequence(
	"example/canton-ceremony/onboarding",
	semver.MustParse("1.0.0"),
	"Full async decentralized party-onboarding ceremony (init → propose → sign → submit)",
	func(b operations.Bundle, deps CantonDeps, in OnboardingInput) (OnboardingOutput, error) {
		out := OnboardingOutput{
			State: CeremonyState{
				Phase:     PhaseInit,
				Threshold: in.Threshold,
			},
		}

		// ── Step 1: Member initialisation ────────────────────────────────────
		// All participants must have run member-init before the proposal can be
		// created.  Each call is idempotent: if a participant has already init'd
		// (i.e. their report exists in the reporter) the operation returns their
		// cached key material immediately.
		members := make([]InitMemberOutput, 0, len(in.Participants))
		for _, pid := range in.Participants {
			r, err := operations.ExecuteOperation(b, InitMemberOp, deps, InitMemberInput{
				NamespaceName: in.NamespaceName,
				ParticipantID: pid,
			})
			if err != nil {
				return out, fmt.Errorf("init-member %s: %w", pid, err)
			}
			members = append(members, r.Output)
			out.State.InitializedMembers = append(out.State.InitializedMembers, pid)
		}

		for _, m := range members {
			if m.ParticipantUID == "" || m.NamespaceFingerprint == "" {
				return out, operations.NewUnrecoverableError(
					fmt.Errorf("member %q has not completed initialisation: uid or namespace fingerprint is missing", m.ParticipantID),
				)
			}
		}

		out.State.Phase = PhaseProposal

		// ── Step 2: Proposal creation ────────────────────────────────────────
		// The coordinator (or anyone) creates the DNS + P2P proposals.  Idempotent: repeated
		// calls with the same member set return the cached proposal.
		proposalReport, err := operations.ExecuteOperation(b, CreateProposalOp, deps, CreateProposalInput{
			NamespaceName:  in.NamespaceName,
			PartyName:      in.PartyName,
			Members:        members,
			SynchronizerID: in.SynchronizerID,
			Threshold:      in.Threshold,
		})
		if err != nil {
			return out, fmt.Errorf("create-proposal: %w", err)
		}
		proposal := proposalReport.Output
		out.State.ProposalHash = proposal.ProposalHashSHA256
		out.State.RequiredSigners = proposal.RequiredSigners
		out.State.Threshold = proposal.Threshold
		out.State.Phase = PhaseSigning

		// ── Step 3: Signature collection ─────────────────────────────────────
		// Each required signer runs SignProposalOp with their own participant ID.
		// Signers who have already signed get their result from the cache.
		// Signers who have not yet acted will fail here — we skip those failures
		// and count only the successful (cached or freshly executed) reports.
		//
		// The threshold check below decides whether the ceremony can proceed.
		var collectedSigs []SignProposalOutput
		for _, signerID := range proposal.RequiredSigners {
			sigReport, sigErr := operations.ExecuteOperation(b, SignProposalOp, deps, SignProposalInput{
				NamespaceName:      in.NamespaceName,
				ParticipantID:      signerID,
				ProposalHashSHA256: proposal.ProposalHashSHA256,
				DNSTxB64:           proposal.DNSTxB64,
				P2PTxB64:           proposal.P2PTxB64,
				SynchronizerID:     in.SynchronizerID,
			})
			if sigErr != nil {
				// This signer has not yet acted; record them as pending.
				deps.Logger.Infow("Signature pending", "signer", signerID, "err", sigErr)
				out.State.PendingSigners = append(out.State.PendingSigners, signerID)

				continue
			}
			collectedSigs = append(collectedSigs, sigReport.Output)
			out.State.CollectedSigners = append(out.State.CollectedSigners, signerID)
		}

		deps.Logger.Infow("Collected signatures", "count", len(collectedSigs), "required", proposal.Threshold)

		// Gate: require at least `threshold` signatures before submitting.
		// out is returned alongside the error so the framework captures the
		// current state snapshot in the sequence report.
		if len(collectedSigs) < proposal.Threshold {
			deps.Logger.Warnw("Threshold not met", "collected", len(collectedSigs), "required", proposal.Threshold)
			return out, fmt.Errorf("%w: %d/%d",
				ErrThresholdNotMet, len(collectedSigs), proposal.Threshold)
		}

		out.State.Phase = PhaseSubmit

		// ── Step 4: Submit ───────────────────────────────────────────────────
		// WithRetry enables the default retry policy (10 attempts, exponential
		// backoff) for transient network / polling errors.
		submitReport, err := operations.ExecuteOperation(
			b, SubmitProposalOp, deps,
			SubmitProposalInput{Proposal: proposal, Signatures: collectedSigs},
			operations.WithRetry[SubmitProposalInput, CantonDeps](),
		)
		if err != nil {
			return out, fmt.Errorf("submit-proposal: %w", err)
		}

		out.State.Phase = PhaseCompleted
		out.State.PendingSigners = nil // all outstanding signers are moot once submitted
		out.PartyID = submitReport.Output.PartyID
		out.DNSConfirmed = submitReport.Output.DNSConfirmed
		out.P2PConfirmed = submitReport.Output.P2PConfirmed

		return out, nil
	},
)
