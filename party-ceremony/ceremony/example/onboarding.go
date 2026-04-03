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

// OnboardingOutput is the final result of a completed [OnboardingSequence].
type OnboardingOutput struct {
	PartyID      string `json:"party_id"`
	DNSConfirmed bool   `json:"dns_confirmed"`
	P2PConfirmed bool   `json:"p2p_confirmed"`
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
				return OnboardingOutput{}, fmt.Errorf("init-member %s: %w", pid, err)
			}
			members = append(members, r.Output)
		}

		for _, m := range members {
			if m.ParticipantUID == "" || m.NamespaceFingerprint == "" {
				return OnboardingOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("member %q has not completed initialisation: uid or namespace fingerprint is missing", m.ParticipantID),
				)
			}
		}

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
			return OnboardingOutput{}, fmt.Errorf("create-proposal: %w", err)
		}
		proposal := proposalReport.Output

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
				// This signer has not yet run their signing step.
				deps.Logger.Infow("Signature pending", "signer", signerID, "err", sigErr)
				continue
			}
			collectedSigs = append(collectedSigs, sigReport.Output)
		}

		deps.Logger.Infow("Collected signatures", "count", len(collectedSigs), "required", proposal.Threshold)

		// Gate: require at least `threshold` signatures before submitting.
		if len(collectedSigs) < proposal.Threshold {
			deps.Logger.Warnw("Threshold not met", "collected", len(collectedSigs), "required", proposal.Threshold)
			return OnboardingOutput{}, fmt.Errorf("%w: %d/%d",
				ErrThresholdNotMet, len(collectedSigs), proposal.Threshold)
		}

		// ── Step 4: Submit ───────────────────────────────────────────────────
		// WithRetry enables the default retry policy (10 attempts, exponential
		// backoff) for transient network / polling errors.
		submitReport, err := operations.ExecuteOperation(
			b, SubmitProposalOp, deps,
			SubmitProposalInput{Proposal: proposal, Signatures: collectedSigs},
			operations.WithRetry[SubmitProposalInput, CantonDeps](),
		)
		if err != nil {
			return OnboardingOutput{}, fmt.Errorf("submit-proposal: %w", err)
		}

		return OnboardingOutput{
			PartyID:      submitReport.Output.PartyID,
			DNSConfirmed: submitReport.Output.DNSConfirmed,
			P2PConfirmed: submitReport.Output.P2PConfirmed,
		}, nil
	},
)
