// Package kick implements the decentralized party kick ceremony using Canton
// Admin gRPC APIs via the Operations framework.
//
// # Overview
//
// The ceremony removes a participant from an existing decentralized party
// (see the Kick workflow in docs/Decentralized Party Management):
//
//  1. ReadCurrentStateOp      – read current DNS and P2P topology state.
//  2. CreateKickDNSProposalOp – coordinator creates updated DNS proposal
//     (kicked owner removed, threshold reduced).
//  3. SignKickDNSProposalOp   – each remaining participant signs the proposal.
//  4. SubmitKickDNSOp         – merge signatures, submit updated DNS.
//  5. ProposeKickP2POp        – each remaining participant proposes the updated
//     P2P mapping (kicked participant removed).
//
// # Authorization Rules
//
// For a serial > 1 DecentralizedNamespaceDefinition update, Canton requires
// threshold-of-current-owners signatures. Since the kicked participant is
// excluded, this means len(RemainingParticipants) >= currentDNSThreshold.
//
// # Async / Resume Pattern
//
// Same pattern as the onboarding ceremony: all operations are idempotent
// via the Operations framework cache. ErrThresholdNotMet is returned until
// enough remaining actors have signed or proposed.
package kick

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [KickSequence] when not enough remaining
// participants have signed the DNS proposal or submitted P2P proposals.
// Callers should treat this as a "come back later" signal.
var ErrThresholdNotMet = errors.New("signature threshold not met: more remaining participants must resume")

// KickSequence orchestrates the full five-step decentralized party kick
// ceremony. It is designed to be called multiple times by different actors:
//
//   - Each call re-enters the state machine. Operations already cached in the
//     reporter are skipped instantly.
//   - If the DNS or P2P signature threshold is not yet met, ErrThresholdNotMet
//     is returned so the caller knows to retry after more remaining members act.
//
// The kicked participant may participate in DNS signing (Canton requires
// threshold-of-current-owners for serial > 1 updates) but is explicitly
// excluded from P2P proposals — only remaining participants propose the
// updated mapping.
var KickSequence = operations.NewSequence(
	"kick/canton-ceremony/decentralized-party",
	semver.MustParse("1.0.0"),
	"Async decentralized party kick (read state → propose DNS update → sign → submit → P2P update)",
	func(b operations.Bundle, deps ceremony.CantonDeps, in KickInput) (KickOutput, error) {
		ctx := b.GetContext()
		out := KickOutput{
			State: CeremonyState{
				Phase:             PhaseReadState,
				KickedParticipant: in.KickedParticipantID,
			},
		}

		// Validate party ID format up front — this is not operation-specific.
		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("invalid decentralized_party_id %q: expected format <prefix>::<namespace>",
					in.DecentralizedPartyID),
			)
		}
		decNS := parts[1]

		// ── Step 1: Read current topology state ──────────────────────────────
		stateReport, err := operations.ExecuteOperation(b, topology.ReadCurrentStateOp, deps, topology.ReadCurrentStateInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("read-current-state: %w", err)
		}
		currentState := stateReport.Output
		runnerUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return out, fmt.Errorf("fetching runner participant UID: %w", err)
		}
		if len(currentState.PartySigningKeysB64) > 0 && slices.Contains(currentState.P2PParticipantUIDs, runnerUID) {
			_, err = operations.ExecuteOperation(b, keys.ResolveProtocolSigningKeyOp, deps, keys.ResolveProtocolSigningKeyInput{
				ParticipantID:       runnerUID,
				KnownSigningKeysB64: currentState.PartySigningKeysB64,
			})
			if err != nil {
				return out, fmt.Errorf("resolve-local-protocol-signing-key: %w", err)
			}
		}

		// ── Sequence-level validation ─────────────────────────────────────────
		newOwners := removeOwner(currentState.DNSOwners, in.KickedNamespaceFingerprint)
		if len(newOwners) == len(currentState.DNSOwners) {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("kicked namespace fingerprint %q not found in DNS owners %v",
					in.KickedNamespaceFingerprint, currentState.DNSOwners),
			)
		}
		if len(newOwners) < 2 {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("kick would reduce owners to %d (minimum 2)", len(newOwners)),
			)
		}

		if int(currentState.DNSThreshold) > len(in.RemainingParticipants)+1 {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf(
					"kick is impossible: DNS signing threshold (%d) exceeds available signers (%d remaining + 1 kicked)",
					currentState.DNSThreshold, len(in.RemainingParticipants),
				),
			)
		}

		kickedInP2P := slices.Contains(currentState.P2PParticipantUIDs, in.KickedParticipantID)
		if !kickedInP2P {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("kicked participant %q not found in P2P mapping participants %v",
					in.KickedParticipantID, currentState.P2PParticipantUIDs),
			)
		}

		newThreshold := in.NewThreshold
		if newThreshold <= 0 {
			newThreshold = int(currentState.DNSThreshold)
		}

		if len(in.RemainingParticipants) < newThreshold {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf(
					"kick is impossible: %d remaining participants cannot reach P2P threshold of %d",
					len(in.RemainingParticipants), newThreshold,
				),
			)
		}

		out.State.DNSThreshold = int(currentState.DNSThreshold)
		out.State.RemainingOwners = newOwners

		// ── Step 2: Create kick DNS proposal ─────────────────────────────────
		out.State.Phase = PhaseDNSProposal
		proposalReport, err := operations.ExecuteOperation(b, topology.CreateKickDNSProposalOp, deps, topology.CreateKickDNSProposalInput{
			DecentralizedNamespace:     decNS,
			CurrentOwners:              currentState.DNSOwners,
			KickedNamespaceFingerprint: in.KickedNamespaceFingerprint,
			NewThreshold:               newThreshold,
			CurrentSerial:              int(currentState.DNSSerial),
			RemainingParticipants:      in.RemainingParticipants,
			KickedParticipantID:        in.KickedParticipantID,
			SynchronizerID:             in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("create-kick-dns-proposal: %w", err)
		}
		proposal := proposalReport.Output
		out.State.ProposalHash = proposal.ProposalHashSHA256
		out.State.RequiredSigners = proposal.RequiredSigners

		// ── Step 3: Collect DNS signatures ────────────────────────────────────
		out.State.Phase = PhaseDNSSigning
		var allSignedTxsB64 []string
		for _, signerUID := range proposal.RequiredSigners {
			sigReport, sigErr := operations.ExecuteOperation(b, topology.SignDNSProposalOp, deps, topology.SignDNSProposalInput{
				ParticipantID:      signerUID,
				ProposalHashSHA256: proposal.ProposalHashSHA256,
				DNSTxB64:           proposal.DNSTxB64,
				SynchronizerID:     in.SynchronizerID,
			})
			if sigErr != nil {
				deps.Logger.Infow("DNS kick signature pending", "signer", signerUID, "err", sigErr)
				out.State.PendingSigners = append(out.State.PendingSigners, signerUID)

				continue
			}
			allSignedTxsB64 = append(allSignedTxsB64, sigReport.Output.SignedDNSTxB64)
			out.State.CollectedSigners = append(out.State.CollectedSigners, signerUID)
		}

		deps.Logger.Infow("Collected kick DNS signatures",
			"collected", len(out.State.CollectedSigners), "required", currentState.DNSThreshold,
		)

		if len(out.State.CollectedSigners) < int(currentState.DNSThreshold) {
			return out, fmt.Errorf("%w: %d/%d DNS signatures collected",
				ErrThresholdNotMet, len(out.State.CollectedSigners), currentState.DNSThreshold,
			)
		}

		// ── Step 4: Submit the kicked DNS update ──────────────────────────────
		out.State.Phase = PhaseDNSSubmit
		_, err = operations.ExecuteOperation(
			b, topology.SubmitDNSOp, deps,
			topology.SubmitDNSInput{
				SignedDNSTxsB64: allSignedTxsB64,
				SynchronizerID:  in.SynchronizerID,
				FilterNamespace: decNS,
			},
			operations.WithRetry[topology.SubmitDNSInput, ceremony.CantonDeps](),
		)
		if err != nil {
			return out, fmt.Errorf("submit-kick-dns: %w", err)
		}

		expectedOwnerCount := len(newOwners)
		err = retry.Do(
			func() error {
				deps.Logger.Infow("Polling kick DNS confirmation", "namespace", decNS)
				dnsState, qErr := deps.Client.GetDNS(ctx, decNS, in.SynchronizerID)
				if qErr != nil {
					return fmt.Errorf("polling DNS state: %w", qErr)
				}
				if len(dnsState.Owners) != expectedOwnerCount {
					return fmt.Errorf("DNS update not yet visible: have %d owners, want %d",
						len(dnsState.Owners), expectedOwnerCount)
				}
				deps.Logger.Infow("Kick DNS confirmed", "namespace", decNS, "owners", expectedOwnerCount)

				return nil
			},
			retry.Context(ctx),
			retry.Attempts(30),
			retry.Delay(500*time.Millisecond),
		)
		if err != nil {
			return out, fmt.Errorf("waiting for kick DNS confirmation: %w", err)
		}

		// ── Step 5: Each remaining participant proposes the updated P2P mapping ──
		out.State.Phase = PhaseP2P
		out.State.P2PRequired = newThreshold

		if runnerUID == in.KickedParticipantID {
			return out, fmt.Errorf("%w: kicked participant may not propose P2P update",
				ErrThresholdNotMet)
		}
		remainingSigningKeysB64, keyErr := protocolSigningKeysForParticipants(
			b,
			deps,
			in.RemainingParticipants,
			currentState.PartySigningKeysB64,
		)
		if keyErr != nil {
			return out, keyErr
		}
		allP2PProposers := in.RemainingParticipants
		for _, uid := range allP2PProposers {
			_, p2pErr := operations.ExecuteOperation(b, topology.ProposeKickP2POp, deps, topology.ProposeKickP2PInput{
				ParticipantID:         uid,
				PartyID:               in.DecentralizedPartyID,
				RemainingParticipants: in.RemainingParticipants,
				NewP2PThreshold:       newThreshold,
				CurrentP2PSerial:      int(currentState.P2PSerial),
				SynchronizerID:        in.SynchronizerID,
				PartySigningKeysB64:   remainingSigningKeysB64,
			})
			if p2pErr != nil {
				deps.Logger.Infow("P2P kick proposal pending", "participant", uid, "err", p2pErr)
				continue
			}
			out.State.P2PProposedCount++
		}

		deps.Logger.Infow("Collected kick P2P proposals",
			"collected", out.State.P2PProposedCount, "required", newThreshold,
		)

		if out.State.P2PProposedCount < newThreshold {
			return out, fmt.Errorf("%w: %d/%d P2P proposals collected",
				ErrThresholdNotMet, out.State.P2PProposedCount, newThreshold,
			)
		}

		// Poll until the updated P2P mapping is confirmed (kicked participant absent).
		err = retry.Do(
			func() error {
				deps.Logger.Infow("Checking kick P2P confirmation", "party", in.DecentralizedPartyID)
				p2pState, qErr := deps.Client.GetP2P(ctx, in.DecentralizedPartyID, in.SynchronizerID)
				if qErr != nil {
					return retry.Unrecoverable(fmt.Errorf("polling P2P state: %w", qErr))
				}
				for _, p := range p2pState.Participants {
					if p.ParticipantUID == in.KickedParticipantID {
						return fmt.Errorf("kicked participant %q still present in P2P mapping", in.KickedParticipantID)
					}
				}
				if len(remainingSigningKeysB64) > 0 {
					if p2pState.PartySigningKeys == nil {
						return fmt.Errorf("P2P signing keys not yet available")
					}
					if len(p2pState.PartySigningKeys.Keys) != len(remainingSigningKeysB64) {
						return fmt.Errorf("P2P signing keys not yet updated: have %d, want %d",
							len(p2pState.PartySigningKeys.Keys), len(remainingSigningKeysB64))
					}
					for _, keyB64 := range remainingSigningKeysB64 {
						if !slices.Contains(p2pState.PartySigningKeys.Keys, keyB64) {
							return fmt.Errorf("remaining party signing key not yet present in P2P mapping")
						}
					}
				}
				deps.Logger.Infow("Kick P2P confirmed", "party", in.DecentralizedPartyID)

				return nil
			},
			retry.Context(ctx),
			retry.Attempts(20),
			retry.Delay(1*time.Second),
		)
		if err != nil {
			return out, fmt.Errorf("waiting for kick P2P confirmation: %w", err)
		}

		out.State.Phase = PhaseCompleted
		out.State.PendingSigners = nil
		out.DNSUpdated = true
		out.P2PUpdated = true
		out.NewThreshold = newThreshold
		out.RemainingOwners = newOwners

		return out, nil
	},
)

// removeOwner returns a copy of owners with the given fingerprint removed.
// Returns the original slice unchanged if the fingerprint is not found.
func removeOwner(owners []string, fingerprint string) []string {
	result := make([]string, 0, len(owners))
	for _, o := range owners {
		if o != fingerprint {
			result = append(result, o)
		}
	}

	return result
}

func protocolSigningKeysForParticipants(
	b operations.Bundle,
	deps ceremony.CantonDeps,
	participantIDs []string,
	knownSigningKeysB64 []string,
) ([]string, error) {
	if len(knownSigningKeysB64) == 0 {
		return nil, nil
	}

	out := make([]string, 0, len(participantIDs))
	for _, uid := range participantIDs {
		report, err := operations.ExecuteOperation(b, keys.ResolveProtocolSigningKeyOp, deps, keys.ResolveProtocolSigningKeyInput{
			ParticipantID:       uid,
			KnownSigningKeysB64: knownSigningKeysB64,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: protocol signing key for remaining participant %q pending: %w", ErrThresholdNotMet, uid, err)
		}
		out = append(out, report.Output.KeyB64)
	}

	return out, nil
}
