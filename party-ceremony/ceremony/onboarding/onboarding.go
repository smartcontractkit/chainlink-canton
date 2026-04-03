// Package onboarding implements the real decentralized party onboarding
// ceremony using Canton Admin gRPC APIs via the Operations framework.
//
// # Overview
//
// The ceremony follows the official Canton docs for setting up a decentralized
// party (see decentralized_parties.rst):
//
//  1. CreateMemberKeyOp     – each participant generates a namespace signing key.
//  2. ProposeNamespaceDelegationOp – each participant publishes their NSD.
//  3. CreateDNSProposalOp   – coordinator creates DecentralizedNamespaceDefinition proposal.
//  4. SignDNSProposalOp      – each signer adds their signature.
//  5. SubmitDNSOp           – coordinator submits the fully-signed DNS.
//  6. CreateAndSubmitP2POp  – each participant authorizes PartyToParticipant mapping.
//
// # Async / Resume Pattern
//
// The ceremony is asynchronous. Different operators run the tool at different
// times. Each ExecuteOperation call is idempotent: if a successful report with
// the same (definition, input) hash already exists in the Reporter the result
// is returned from the cache and the side-effect is not repeated.
//
// When the signature threshold has not yet been reached the sequence returns
// [ErrThresholdNotMet]. The caller is expected to re-run the sequence later
// (after more signers have completed their signing step).
package onboarding

import (
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"
	"github.com/chainlink/canton-party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [OnboardingSequence] when the number of
// collected signatures have not yet reached the required threshold. Callers
// should treat this as a transient "come back later" signal.
var ErrThresholdNotMet = errors.New("signature threshold not met: more signers must resume the sequence")

// OnboardingSequence orchestrates the full six-step decentralized party
// onboarding ceremony. It is designed to be called multiple times by different
// actors:
//
//   - On each call the sequence re-enters the state machine. Operations whose
//     (definition, input) hash already has a successful report are skipped
//     instantly — they reflect work done by previous actors.
//   - If the DNS signature threshold is not yet met [ErrThresholdNotMet] is
//     returned so the caller knows to retry after more signers have acted.
//   - Once the threshold is met the coordinator's SubmitDNSOp runs, followed
//     by each participant independently authorizing the P2P mapping.
var OnboardingSequence = operations.NewSequence(
	"onboarding/canton-ceremony/decentralized-party",
	semver.MustParse("1.0.0"),
	"Full async decentralized party onboarding (key-gen → NSD → DNS → sign → submit → P2P)",
	func(b operations.Bundle, deps ceremony.CantonDeps, in OnboardingInput) (OnboardingOutput, error) {
		// ── Step 1: Member key generation ────────────────────────────────────
		members := make([]CreateMemberKeyOutput, 0)
		for _, pid := range in.Participants {
			deps.Logger.Info("creating key for participant", pid)
			r, err := operations.ExecuteOperation(b, CreateMemberKeyOp, deps, CreateMemberKeyInput{
				NamespaceName: in.NamespaceName,
				ParticipantID: pid,
			})
			if err != nil {
				if !retry.IsRecoverable(err) {
					return OnboardingOutput{}, fmt.Errorf("create-member-key %s: %w", pid, err)
				}
				deps.Logger.Infow("Member key pending", "participant", pid, "err", err)

				continue
			}
			members = append(members, r.Output)
		}

		for _, m := range members {
			deps.Logger.Infow("Member key generated",
				"participant", m.ParticipantID,
				"uid", m.ParticipantUID,
				"namespace_fingerprint", m.NamespaceFingerprint)
		}

		pid, err := deps.Client.GetParticipantUID(b.GetContext())
		if err != nil {
			return OnboardingOutput{}, fmt.Errorf("fetching client participant ID: %w", err)
		}

		// ── Step 2: Namespace delegation ─────────────────────────────────────
		// Each participant publishes their namespace delegation to the synchronizer.
		for _, m := range members {
			// Can only be executed from the participant that owns the key
			if m.ParticipantID != pid {
				deps.Logger.Infow("Namespace delegation pending",
					"participant", m.ParticipantID,
					"reason", "current participant does not own the key",
				)

				continue
			}
			_, err := operations.ExecuteOperation(b, ProposeNamespaceDelegationOp, deps, ProposeNSDInput{
				ParticipantID:  m.ParticipantID,
				SigningKeyB64:  m.SigningKeyB64,
				Namespace:      m.NamespaceFingerprint,
				SynchronizerID: in.SynchronizerID,
			})
			if err != nil {
				return OnboardingOutput{}, fmt.Errorf("propose-nsd %s: %w", m.ParticipantID, err)
			}
		}

		// Gate: all participants must have generated their keys before the DNS
		// proposal can be created. A partial member set would produce a wrong
		// decentralized namespace; actors whose keys are not yet cached return
		// ErrThresholdNotMet so callers know to retry after more actors run.
		if len(members) < len(in.Participants) {
			deps.Logger.Warnw("Not all participants have generated their member keys",
				"collected", len(members), "required", len(in.Participants))

			return OnboardingOutput{}, fmt.Errorf("%w: %d/%d participants have generated their member keys",
				ErrThresholdNotMet, len(members), len(in.Participants))
		}

		// Gate: wait for all namespace delegations to be visible in the
		// synchronizer's topology state.
		err = retry.Do(
			func() error {
				for _, m := range members {
					exists, qErr := deps.Client.NSDExists(b.GetContext(), m.NamespaceFingerprint, in.SynchronizerID)
					if qErr != nil {
						return retry.Unrecoverable(fmt.Errorf("checking NSD for %s: %w", m.ParticipantID, qErr))
					}
					if !exists {
						return fmt.Errorf("namespace delegation not yet visible for %s (namespace %s)",
							m.ParticipantID, m.NamespaceFingerprint)
					}
				}

				return nil
			},
			retry.Context(b.GetContext()),
			retry.Attempts(30),
			retry.Delay(500*time.Millisecond),
		)
		if err != nil {
			return OnboardingOutput{}, fmt.Errorf("waiting for namespace delegations: %w", err)
		}

		deps.Logger.Infow("All namespace delegations confirmed", "count", len(members))

		// ── Step 3: DNS proposal creation ────────────────────────────────────
		proposalReport, err := operations.ExecuteOperation(b, CreateDNSProposalOp, deps, CreateDNSProposalInput{
			NamespaceName:  in.NamespaceName,
			Members:        members,
			SynchronizerID: in.SynchronizerID,
			Threshold:      in.Threshold,
		})
		if err != nil {
			return OnboardingOutput{}, fmt.Errorf("create-dns-proposal: %w", err)
		}
		proposal := proposalReport.Output

		// ── Step 4: DNS signature collection ─────────────────────────────────
		// Each required signer signs the DNS proposal independently.
		// Signers who have not yet acted will fail here — we skip those
		// failures and count only the successful reports.
		// We collect ALL signed transactions so SubmitDNSOp can merge their
		// Signature lists into one fully-authorized transaction.
		var allSignedTxsB64 []string
		var collectedCount int
		for _, signerID := range proposal.RequiredSigners {
			sigReport, sigErr := operations.ExecuteOperation(b, SignDNSProposalOp, deps, SignDNSProposalInput{
				ParticipantID:      signerID,
				ProposalHashSHA256: proposal.ProposalHashSHA256,
				DNSTxB64:           proposal.DNSTxB64,
				SynchronizerID:     in.SynchronizerID,
			})
			if sigErr != nil {
				deps.Logger.Infow("Signature pending", "signer", signerID, "err", sigErr)
				continue
			}
			allSignedTxsB64 = append(allSignedTxsB64, sigReport.Output.SignedDNSTxB64)
			collectedCount++
		}

		deps.Logger.Infow("Collected DNS signatures",
			"count", collectedCount, "required", proposal.Threshold)

		// Gate: require all signatures for initial DNS creation.
		if collectedCount < proposal.Threshold {
			deps.Logger.Warnw("Threshold not met",
				"collected", collectedCount, "required", proposal.Threshold)

			return OnboardingOutput{}, fmt.Errorf("%w: %d/%d",
				ErrThresholdNotMet, collectedCount, proposal.Threshold)
		}

		// ── Step 5: Submit DNS ───────────────────────────────────────────────
		// WithRetry enables the default retry policy for transient network errors.
		_, err = operations.ExecuteOperation(
			b, SubmitDNSOp, deps,
			SubmitDNSInput{
				SignedDNSTxsB64: allSignedTxsB64,
				SynchronizerID:  in.SynchronizerID,
				FilterNamespace: proposal.DecentralizedNS,
			},
			operations.WithRetry[SubmitDNSInput, ceremony.CantonDeps](),
		)
		if err != nil {
			return OnboardingOutput{}, fmt.Errorf("submit-dns: %w", err)
		}

		// Check that the DNS is confirmed in the topology state.
		ctx := b.GetContext()
		err = retry.Do(
			func() error {
				deps.Logger.Infow("Checking DNS confirmation", "namespace", proposal.DecentralizedNS)
				// Poll until DNS is visible at head state.
				exists, err := deps.Client.DNSExists(ctx, proposal.DecentralizedNS, in.SynchronizerID)
				if err != nil {
					return fmt.Errorf("checking DNS confirmation: %w", err)
				}
				if !exists {
					return fmt.Errorf("DNS not yet confirmed for namespace %s", proposal.DecentralizedNS)
				}

				deps.Logger.Infow("DNS submitted and confirmed", "namespace", proposal.DecentralizedNS)

				return nil
			},
			retry.Context(b.GetContext()),
			retry.Attempts(5),
			retry.Delay(5*time.Second),
		)
		if err != nil {
			return OnboardingOutput{}, fmt.Errorf("waiting for DNS confirmation: %w", err)
		}

		// ── Step 6: PartyToParticipant mapping ───────────────────────────────
		// Each participant independently proposes the same P2P mapping using
		// their own client.
		// Canton accumulates proposals and activates the mapping once the
		// threshold is reached.
		partyID := fmt.Sprintf("%s::%s", in.PartyPrefix, proposal.DecentralizedNS)

		var p2pProposedCount int
		for _, m := range members {
			_, p2pErr := operations.ExecuteOperation(b, ProposeP2POp, deps, ProposeP2PInput{
				ParticipantID:  m.ParticipantID,
				PartyID:        partyID,
				Members:        members,
				SynchronizerID: in.SynchronizerID,
				Threshold:      in.Threshold,
			})
			if p2pErr != nil {
				deps.Logger.Infow("P2P proposal pending", "participant", m.ParticipantID, "err", p2pErr)
				continue
			}
			p2pProposedCount++
		}

		deps.Logger.Infow("Collected P2P proposals",
			"count", p2pProposedCount, "required", in.Threshold)

		if p2pProposedCount < in.Threshold {
			deps.Logger.Warnw("P2P threshold not met",
				"collected", p2pProposedCount, "required", in.Threshold)

			return OnboardingOutput{}, fmt.Errorf("%w: %d/%d P2P proposals collected",
				ErrThresholdNotMet, p2pProposedCount, in.Threshold)
		}

		// Poll until Canton activates the mapping.
		p2pCtx := b.GetContext()
		p2pErr := retry.Do(
			func() error {
				exists, qErr := deps.Client.P2PExists(p2pCtx, partyID, in.SynchronizerID)
				if qErr != nil {
					return retry.Unrecoverable(fmt.Errorf("checking P2P confirmation: %w", qErr))
				}
				if !exists {
					return fmt.Errorf("P2P not yet confirmed for party %s", partyID)
				}

				return nil
			},
			retry.Context(b.GetContext()),
			retry.Attempts(30),
			retry.Delay(500*time.Millisecond),
		)
		if p2pErr != nil {
			return OnboardingOutput{}, fmt.Errorf("waiting for P2P confirmation: %w", p2pErr)
		}

		deps.Logger.Infow("P2P mapping confirmed", "party_id", partyID)

		return OnboardingOutput{
			PartyID:      partyID,
			DNSConfirmed: true,
			P2PConfirmed: true,
		}, nil
	},
)
