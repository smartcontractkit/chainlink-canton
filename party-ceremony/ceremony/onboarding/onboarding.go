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

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"

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
		out := OnboardingOutput{
			State: CeremonyState{
				Phase:     PhaseKeyGen,
				Threshold: in.Threshold,
			},
		}

		// ── Step 1: Member key generation ────────────────────────────────────
		members := make([]keys.CreateMemberKeyOutput, 0)
		for _, pid := range in.Participants {
			deps.Logger.Info("creating key for participant", pid)
			r, err := operations.ExecuteOperation(b, keys.CreateMemberKeyOp, deps, keys.CreateMemberKeyInput{
				NamespaceName:     in.NamespaceName,
				ParticipantID:     pid,
				KmsNamespaceKeyID: in.KmsNamespaceKeyID,
				KmsProtocolKeyID:  in.KmsProtocolKeyID,
			})
			if err != nil {
				if !retry.IsRecoverable(err) {
					return out, fmt.Errorf("create-member-key %s: %w", pid, err)
				}
				deps.Logger.Infow("Member key pending", "participant", pid, "err", err)

				continue
			}
			members = append(members, r.Output)
			out.State.KeysGenerated = append(out.State.KeysGenerated, pid)
		}

		for _, m := range members {
			deps.Logger.Infow("Member key generated",
				"participant", m.ParticipantID,
				"uid", m.ParticipantUID,
				"namespace_fingerprint", m.NamespaceFingerprint)
		}

		pid, err := deps.Client.GetParticipantUID(b.GetContext())
		if err != nil {
			return out, fmt.Errorf("fetching client participant ID: %w", err)
		}

		// ── Step 2: Namespace delegation ─────────────────────────────────────
		out.State.Phase = PhaseNSD
		for _, m := range members {
			if m.ParticipantID != pid {
				deps.Logger.Infow("Namespace delegation pending",
					"participant", m.ParticipantID,
					"reason", "current participant does not own the key",
				)

				continue
			}
			_, err := operations.ExecuteOperation(b, topology.ProposeNamespaceDelegationOp, deps, topology.ProposeNSDInput{
				ParticipantID:  m.ParticipantID,
				SigningKeyB64:  m.SigningKeyB64,
				Namespace:      m.NamespaceFingerprint,
				SynchronizerID: in.SynchronizerID,
			})
			if err != nil {
				return out, fmt.Errorf("propose-nsd %s: %w", m.ParticipantID, err)
			}
			out.State.NSDsProposed = append(out.State.NSDsProposed, m.ParticipantID)
		}

		// Gate: all participants must have generated their keys before the DNS
		// proposal can be created.
		if len(members) < len(in.Participants) {
			deps.Logger.Warnw("Not all participants have generated their member keys",
				"collected", len(members), "required", len(in.Participants))

			return out, fmt.Errorf("%w: %d/%d participants have generated their member keys",
				ErrThresholdNotMet, len(members), len(in.Participants))
		}

		// Gate: wait for all namespace delegations to be visible.
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
			return out, fmt.Errorf("waiting for namespace delegations: %w", err)
		}

		deps.Logger.Infow("All namespace delegations confirmed", "count", len(members))

		// ── Step 3: DNS proposal creation ────────────────────────────────────
		out.State.Phase = PhaseDNSProposal
		proposalReport, err := operations.ExecuteOperation(b, topology.CreateDNSProposalOp, deps, topology.CreateDNSProposalInput{
			NamespaceName:  in.NamespaceName,
			Members:        members,
			SynchronizerID: in.SynchronizerID,
			Threshold:      in.Threshold,
		})
		if err != nil {
			return out, fmt.Errorf("create-dns-proposal: %w", err)
		}
		proposal := proposalReport.Output
		out.State.ProposalHash = proposal.ProposalHashSHA256
		out.State.RequiredSigners = proposal.RequiredSigners
		out.State.Threshold = proposal.Threshold

		// ── Step 4: DNS signature collection ─────────────────────────────────
		out.State.Phase = PhaseDNSSigning
		var allSignedTxsB64 []string
		for _, signerID := range proposal.RequiredSigners {
			sigReport, sigErr := operations.ExecuteOperation(b, topology.SignDNSProposalOp, deps, topology.SignDNSProposalInput{
				ParticipantID:      signerID,
				ProposalHashSHA256: proposal.ProposalHashSHA256,
				DNSTxB64:           proposal.DNSTxB64,
				SynchronizerID:     in.SynchronizerID,
			})
			if sigErr != nil {
				deps.Logger.Infow("Signature pending", "signer", signerID, "err", sigErr)
				out.State.PendingSigners = append(out.State.PendingSigners, signerID)

				continue
			}
			allSignedTxsB64 = append(allSignedTxsB64, sigReport.Output.SignedDNSTxB64)
			out.State.CollectedSigners = append(out.State.CollectedSigners, signerID)
		}

		deps.Logger.Infow("Collected DNS signatures",
			"count", len(out.State.CollectedSigners), "required", proposal.Threshold)

		// Gate: require threshold signatures for initial DNS creation.
		if len(out.State.CollectedSigners) < proposal.Threshold {
			deps.Logger.Warnw("Threshold not met",
				"collected", len(out.State.CollectedSigners), "required", proposal.Threshold)

			return out, fmt.Errorf("%w: %d/%d",
				ErrThresholdNotMet, len(out.State.CollectedSigners), proposal.Threshold)
		}

		// ── Step 5: Submit DNS ───────────────────────────────────────────────
		out.State.Phase = PhaseDNSSubmit
		_, err = operations.ExecuteOperation(
			b, topology.SubmitDNSOp, deps,
			topology.SubmitDNSInput{
				SignedDNSTxsB64: allSignedTxsB64,
				SynchronizerID:  in.SynchronizerID,
				FilterNamespace: proposal.DecentralizedNS,
			},
			operations.WithRetry[topology.SubmitDNSInput, ceremony.CantonDeps](),
		)
		if err != nil {
			return out, fmt.Errorf("submit-dns: %w", err)
		}

		// Check that the DNS is confirmed in the topology state.
		ctx := b.GetContext()
		err = retry.Do(
			func() error {
				deps.Logger.Infow("Checking DNS confirmation", "namespace", proposal.DecentralizedNS)
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
			return out, fmt.Errorf("waiting for DNS confirmation: %w", err)
		}

		// ── Step 6: PartyToParticipant mapping ───────────────────────────────
		out.State.Phase = PhaseP2P
		partyID := fmt.Sprintf("%s::%s", in.PartyPrefix, proposal.DecentralizedNS)
		out.State.P2PRequired = in.Threshold

		for _, m := range members {
			_, p2pErr := operations.ExecuteOperation(b, topology.ProposeP2POp, deps, topology.ProposeP2PInput{
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
			out.State.P2PProposedCount++
		}

		deps.Logger.Infow("Collected P2P proposals",
			"count", out.State.P2PProposedCount, "required", in.Threshold)

		if out.State.P2PProposedCount < in.Threshold {
			deps.Logger.Warnw("P2P threshold not met",
				"collected", out.State.P2PProposedCount, "required", in.Threshold)

			return out, fmt.Errorf("%w: %d/%d P2P proposals collected",
				ErrThresholdNotMet, out.State.P2PProposedCount, in.Threshold)
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
			return out, fmt.Errorf("waiting for P2P confirmation: %w", p2pErr)
		}

		deps.Logger.Infow("P2P mapping confirmed", "party_id", partyID)

		out.State.Phase = PhaseCompleted
		out.State.PendingSigners = nil
		out.PartyID = partyID
		out.DNSConfirmed = true
		out.P2PConfirmed = true

		return out, nil
	},
)
