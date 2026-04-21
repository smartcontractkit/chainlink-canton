// Package keyrotation implements the decentralized party key rotation ceremony
// using Canton Admin gRPC APIs via the Operations framework.
//
// # Overview
//
// The ceremony replaces a participant's namespace key and/or DAML signing key
// in the topology (see the Key Rotation section in docs/Decentralized Party
// Management):
//
//  1. ReadCurrentStateOp           – read current DNS and P2P topology state
//     (reused from kick).
//  2. GenerateRotatedKeyOp         – target participant generates new key(s)
//     and discovers old DAML key.
//  3. ProposeRotatedNSDOp          – target publishes NSD for new namespace key
//     (reused from addparticipant).
//  4. CreateRotationDNSProposalOp  – coordinator creates updated DNS proposal
//     (old fingerprint replaced with new).
//  5. SignRotationDNSProposalOp    – each member signs the DNS proposal
//     (reused from kick).
//  6. SubmitRotationDNSOp          – merge signatures, submit updated DNS
//     (reused from kick).
//  7. ProposeRotationP2POp         – each member proposes updated P2P mapping
//     with new DAML key.
//
// Steps 3-6 only execute when RotateNamespaceKey is true.
// Step 7 only executes when RotateDamlKey is true.
// When both are true, namespace rotation completes first (DNS authority must
// be updated before P2P proposals).
//
// # Authorization Rules
//
// For namespace key rotation, Canton requires threshold-of-current-owners
// signatures for the serial > 1 DNS update. All current members are eligible
// signers.
//
// For DAML key rotation, each member proposes the updated P2P mapping
// independently. Canton accumulates proposals and activates when threshold
// is reached.
//
// # Async / Resume Pattern
//
// Same pattern as onboarding, kick, and add-participant: all operations are
// idempotent via the Operations framework cache. ErrThresholdNotMet is
// returned until enough actors have contributed.
package keyrotation

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"
	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/ceremony/addparticipant"
	"github.com/chainlink/canton-party-ceremony/ceremony/kick"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [KeyRotationSequence] when not enough
// participants have contributed keys, signatures, or proposals.
// Callers should treat this as a "come back later" signal.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must resume")

// KeyRotationSequence orchestrates the key rotation ceremony. It is designed
// to be called multiple times by different actors:
//
//   - Each call re-enters the state machine. Operations already cached in the
//     reporter are skipped instantly.
//   - If the DNS signature threshold or P2P proposal threshold is not yet met,
//     ErrThresholdNotMet is returned so the caller knows to retry after more
//     members act.
//   - Participants are dynamically fetched from the P2P topology state (not
//     supplied as input).
//   - The old DAML key fingerprint is auto-discovered by cross-referencing
//     the target's vault with the party's signing keys.
var KeyRotationSequence = operations.NewSequence(
	"keyrotation/canton-ceremony/decentralized-party",
	semver.MustParse("1.0.0"),
	"Async decentralized party key rotation (read state → generate key → NSD → propose DNS → sign → submit → P2P update)",
	func(b operations.Bundle, deps ceremony.CantonDeps, in KeyRotationInput) (KeyRotationOutput, error) {
		ctx := b.GetContext()

		out := KeyRotationOutput{
			State: CeremonyState{
				Phase:           PhaseReadState,
				RotateNamespace: in.RotateNamespaceKey,
				RotateDaml:      in.RotateDamlKey,
			},
		}

		// ── Input validation ─────────────────────────────────────────────────
		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("invalid decentralized_party_id %q: expected format <prefix>::<namespace>",
					in.DecentralizedPartyID),
			)
		}
		decNS := parts[1]

		if !in.RotateNamespaceKey && !in.RotateDamlKey {
			return out, operations.NewUnrecoverableError(
				errors.New("at least one of rotate_namespace_key or rotate_daml_key must be true"),
			)
		}

		if in.TargetParticipantID == "" {
			return out, operations.NewUnrecoverableError(
				errors.New("target_participant_id is required"),
			)
		}

		// ── Step 1: Read current topology state ──────────────────────────────
		stateReport, err := operations.ExecuteOperation(b, kick.ReadCurrentStateOp, deps, kick.ReadCurrentStateInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("read-current-state: %w", err)
		}
		currentState := stateReport.Output

		// Dynamically derive the participant list from topology.
		allParticipantUIDs := currentState.P2PParticipantUIDs

		// ── Sequence-level validation ────────────────────────────────────────
		if in.RotateNamespaceKey {
			if in.TargetNamespaceFingerprint == "" {
				return out, operations.NewUnrecoverableError(
					errors.New("target_namespace_fingerprint is required when rotate_namespace_key is true"),
				)
			}
			if !slices.Contains(currentState.DNSOwners, in.TargetNamespaceFingerprint) {
				return out, operations.NewUnrecoverableError(
					fmt.Errorf("target namespace fingerprint %q not found in DNS owners %v",
						in.TargetNamespaceFingerprint, currentState.DNSOwners),
				)
			}
		}

		if in.RotateDamlKey && len(currentState.PartySigningKeysB64) == 0 {
			return out, operations.NewUnrecoverableError(
				errors.New("DAML key rotation requested but no party signing keys found in P2P topology"),
			)
		}

		// Compute effective threshold (keep current unless overridden).
		newThreshold := in.NewThreshold
		if newThreshold <= 0 {
			newThreshold = int(currentState.DNSThreshold)
		}

		out.State.DNSThreshold = int(currentState.DNSThreshold)
		out.State.Phase = PhaseKeyGen

		// ── Step 2: Generate rotated keys (target only) ──────────────────────
		keyReport, err := operations.ExecuteOperation(b, GenerateRotatedKeyOp, deps, GenerateRotatedKeyInput{
			ParticipantID:       in.TargetParticipantID,
			SynchronizerID:      in.SynchronizerID,
			DNSOwners:           currentState.DNSOwners,
			RotateNamespaceKey:  in.RotateNamespaceKey,
			RotateDamlKey:       in.RotateDamlKey,
			KnownSigningKeysB64: currentState.PartySigningKeysB64,
		})
		if err != nil {
			deps.Logger.Infow("Rotated key generation pending",
				"target", in.TargetParticipantID, "err", err)

			return out, fmt.Errorf("%w: target participant has not generated rotated keys yet",
				ErrThresholdNotMet)
		}
		rotatedKeys := keyReport.Output
		out.State.TargetKeyGenReady = true

		// ── Steps 3-6: Namespace key rotation ────────────────────────────────
		dnsUpdated := false
		if in.RotateNamespaceKey {
			out.State.Phase = PhaseNSD

			// Step 3: Propose NSD for new namespace key (target only).
			_, err = operations.ExecuteOperation(b, addparticipant.ProposeNewNSDOp, deps, addparticipant.ProposeNewNSDInput{
				ParticipantID:  in.TargetParticipantID,
				SigningKeyB64:  rotatedKeys.NewNamespaceKeyB64,
				Namespace:      rotatedKeys.NewNamespaceFingerprint,
				SynchronizerID: in.SynchronizerID,
			})
			if err != nil {
				deps.Logger.Infow("Rotated NSD proposal pending",
					"target", in.TargetParticipantID, "err", err)

				return out, fmt.Errorf("%w: target participant has not proposed rotated NSD yet",
					ErrThresholdNotMet)
			}
			out.State.NSDProposed = true

			// Poll until NSD is visible on the synchronizer.
			err = retry.Do(
				func() error {
					exists, qErr := deps.Client.NSDExists(ctx, rotatedKeys.NewNamespaceFingerprint, in.SynchronizerID)
					if qErr != nil {
						return retry.Unrecoverable(fmt.Errorf("checking NSD for rotated key: %w", qErr))
					}
					if !exists {
						return fmt.Errorf("rotated namespace delegation not yet visible for %s",
							rotatedKeys.NewNamespaceFingerprint)
					}

					return nil
				},
				retry.Context(ctx),
				retry.Attempts(30),
				retry.Delay(500*time.Millisecond),
			)
			if err != nil {
				return out, fmt.Errorf("waiting for rotated NSD: %w", err)
			}

			deps.Logger.Infow("Rotated NSD confirmed",
				"namespace", rotatedKeys.NewNamespaceFingerprint)

			out.State.Phase = PhaseDNSProposal

			// Step 4: Create DNS proposal (swap old→new fingerprint in owners).
			proposalReport, propErr := operations.ExecuteOperation(b, CreateRotationDNSProposalOp, deps, CreateRotationDNSProposalInput{
				DecentralizedNamespace:  decNS,
				CurrentOwners:           currentState.DNSOwners,
				OldNamespaceFingerprint: in.TargetNamespaceFingerprint,
				NewNamespaceFingerprint: rotatedKeys.NewNamespaceFingerprint,
				CurrentThreshold:        newThreshold,
				CurrentSerial:           int(currentState.DNSSerial),
				AllParticipantIDs:       allParticipantUIDs,
				SynchronizerID:          in.SynchronizerID,
			})
			if propErr != nil {
				return out, fmt.Errorf("create-rotation-dns-proposal: %w", propErr)
			}
			proposal := proposalReport.Output
			out.State.RequiredSigners = proposal.RequiredSigners
			out.State.ProposalHash = proposal.ProposalHashSHA256
			out.State.Phase = PhaseDNSSigning

			// Step 5: Collect DNS signatures from all members.
			var allSignedTxsB64 []string
			var dnsSigCount int
			for _, signerUID := range proposal.RequiredSigners {
				sigReport, sigErr := operations.ExecuteOperation(b, kick.SignKickDNSProposalOp, deps, kick.SignKickDNSProposalInput{
					ParticipantID:      signerUID,
					ProposalHashSHA256: proposal.ProposalHashSHA256,
					DNSTxB64:           proposal.DNSTxB64,
					SynchronizerID:     in.SynchronizerID,
				})
				if sigErr != nil {
					deps.Logger.Infow("DNS rotation signature pending", "signer", signerUID, "err", sigErr)
					out.State.PendingSigners = append(out.State.PendingSigners, signerUID)

					continue
				}
				allSignedTxsB64 = append(allSignedTxsB64, sigReport.Output.SignedDNSTxB64)
				dnsSigCount++
				out.State.CollectedSigners = append(out.State.CollectedSigners, signerUID)
			}

			deps.Logger.Infow("Collected rotation DNS signatures",
				"collected", dnsSigCount, "required", currentState.DNSThreshold,
			)

			// Gate: Canton requires currentThreshold signatures for serial > 1.
			if dnsSigCount < int(currentState.DNSThreshold) {
				return out, fmt.Errorf("%w: %d/%d DNS signatures collected",
					ErrThresholdNotMet, dnsSigCount, currentState.DNSThreshold,
				)
			}

			out.State.Phase = PhaseDNSSubmit

			// Step 6: Submit the rotation DNS update.
			_, err = operations.ExecuteOperation(
				b, kick.SubmitKickDNSOp, deps,
				kick.SubmitKickDNSInput{
					SignedDNSTxsB64: allSignedTxsB64,
					SynchronizerID:  in.SynchronizerID,
					FilterNamespace: decNS,
				},
				operations.WithRetry[kick.SubmitKickDNSInput, ceremony.CantonDeps](),
			)
			if err != nil {
				return out, fmt.Errorf("submit-rotation-dns: %w", err)
			}

			// Poll until the updated DNS is visible (new fingerprint in owners, old gone).
			err = retry.Do(
				func() error {
					deps.Logger.Infow("Polling rotation DNS confirmation", "namespace", decNS)
					dnsState, qErr := deps.Client.GetDNS(ctx, decNS, in.SynchronizerID)
					if qErr != nil {
						return fmt.Errorf("polling DNS state: %w", qErr)
					}
					if slices.Contains(dnsState.Owners, in.TargetNamespaceFingerprint) {
						return fmt.Errorf("old namespace fingerprint %q still present in DNS owners",
							in.TargetNamespaceFingerprint)
					}
					if !slices.Contains(dnsState.Owners, rotatedKeys.NewNamespaceFingerprint) {
						return fmt.Errorf("new namespace fingerprint %q not yet present in DNS owners",
							rotatedKeys.NewNamespaceFingerprint)
					}
					deps.Logger.Infow("Rotation DNS confirmed", "namespace", decNS)

					return nil
				},
				retry.Context(ctx),
				retry.Attempts(30),
				retry.Delay(500*time.Millisecond),
			)
			if err != nil {
				return out, fmt.Errorf("waiting for rotation DNS confirmation: %w", err)
			}

			dnsUpdated = true
		}

		// ── Step 7: DAML key rotation P2P update ─────────────────────────────
		p2pUpdated := false
		if in.RotateDamlKey {
			out.State.Phase = PhaseP2P
			out.State.P2PRequired = newThreshold

			var p2pProposedCount int
			for _, uid := range allParticipantUIDs {
				_, p2pErr := operations.ExecuteOperation(b, ProposeRotationP2POp, deps, ProposeRotationP2PInput{
					ParticipantID:         uid,
					PartyID:               in.DecentralizedPartyID,
					AllParticipantUIDs:    allParticipantUIDs,
					NewP2PThreshold:       newThreshold,
					CurrentP2PSerial:      int(currentState.P2PSerial),
					SynchronizerID:        in.SynchronizerID,
					CurrentSigningKeysB64: currentState.PartySigningKeysB64,
					OldDamlKeyB64:         rotatedKeys.OldDamlKeyB64,
					NewDamlKeyB64:         rotatedKeys.NewDamlKeyB64,
					SigningKeysThreshold:  currentState.PartySigningThreshold,
				})
				if p2pErr != nil {
					deps.Logger.Infow("P2P rotation proposal pending", "participant", uid, "err", p2pErr)
					continue
				}
				p2pProposedCount++
			}
			out.State.P2PProposedCount = p2pProposedCount

			deps.Logger.Infow("Collected rotation P2P proposals",
				"collected", p2pProposedCount, "required", newThreshold,
			)

			if p2pProposedCount < newThreshold {
				return out, fmt.Errorf("%w: %d/%d P2P proposals collected",
					ErrThresholdNotMet, p2pProposedCount, newThreshold,
				)
			}

			// Poll until the updated P2P signing keys are confirmed.
			err = retry.Do(
				func() error {
					deps.Logger.Infow("Checking rotation P2P confirmation", "party", in.DecentralizedPartyID)
					p2pState, qErr := deps.Client.GetP2P(ctx, in.DecentralizedPartyID, in.SynchronizerID)
					if qErr != nil {
						return retry.Unrecoverable(fmt.Errorf("polling P2P state: %w", qErr))
					}
					if p2pState.PartySigningKeys == nil {
						return fmt.Errorf("P2P signing keys not yet available")
					}
					// Verify the new key is present by checking if the old key is gone.
					if slices.Contains(p2pState.PartySigningKeys.Keys, rotatedKeys.OldDamlKeyB64) {
						return fmt.Errorf("old DAML key still present in P2P signing keys")
					}
					deps.Logger.Infow("Rotation P2P confirmed", "party", in.DecentralizedPartyID)

					return nil
				},
				retry.Context(ctx),
				retry.Attempts(20),
				retry.Delay(1*time.Second),
			)
			if err != nil {
				return out, fmt.Errorf("waiting for rotation P2P confirmation: %w", err)
			}

			p2pUpdated = true
		}

		out.NamespaceKeyRotated = in.RotateNamespaceKey
		out.DamlKeyRotated = in.RotateDamlKey
		out.NewNamespaceFingerprint = rotatedKeys.NewNamespaceFingerprint
		out.NewDamlKeyFingerprint = rotatedKeys.NewDamlKeyFingerprint
		out.DNSUpdated = dnsUpdated
		out.P2PUpdated = p2pUpdated
		out.State.Phase = PhaseCompleted
		out.State.PendingSigners = nil

		return out, nil
	},
)
