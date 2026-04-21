// Package addparticipant implements the decentralized party add-participant
// ceremony using Canton Admin gRPC APIs via the Operations framework.
//
// # Overview
//
// The ceremony adds a new participant to an existing decentralized party
// (reverse of the kick workflow):
//
//  1. GenerateNewMemberKeyOp  – new participant generates namespace + DAML keys.
//  2. ProposeNewNSDOp         – new participant publishes namespace delegation.
//  3. ReadCurrentStateOp      – read current DNS and P2P topology state.
//  4. CreateAddDNSProposalOp  – coordinator creates updated DNS proposal
//     (new owner added, threshold adjusted).
//  5. SignAddDNSProposalOp    – each existing participant signs the proposal.
//  6. SubmitAddDNSOp          – merge signatures, submit updated DNS.
//  7. ProposeAddP2POp         – each existing participant proposes the updated
//     P2P mapping (new participant added with Confirmation permission).
//
// # Authorization Rules
//
// For a serial > 1 DecentralizedNamespaceDefinition update, Canton requires
// threshold-of-current-owners signatures. The new participant is NOT yet an
// owner and cannot sign. Only existing members sign the DNS update.
//
// # Async / Resume Pattern
//
// Same pattern as onboarding and kick: all operations are idempotent via the
// Operations framework cache. ErrThresholdNotMet is returned until enough
// actors have contributed their keys, signatures, or proposals.
package addparticipant

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"
	"github.com/chainlink/canton-party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [AddParticipantSequence] when not enough
// actors have contributed keys, signatures, or proposals.
// Callers should treat this as a "come back later" signal.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must resume")

// AddParticipantSequence orchestrates the full seven-step add-participant
// ceremony. It is designed to be called multiple times by different actors:
//
//   - Each call re-enters the state machine. Operations already cached in the
//     reporter are skipped instantly.
//   - If the key generation, DNS signature threshold, or P2P proposal threshold
//     is not yet met, ErrThresholdNotMet is returned.
//
// The new participant generates keys and publishes NSD. Existing members sign
// the DNS update and propose the P2P mapping.
var AddParticipantSequence = operations.NewSequence(
	"add-participant/canton-ceremony/decentralized-party",
	semver.MustParse("1.0.0"),
	"Async decentralized party add-participant (key-gen → NSD → read state → propose DNS update → sign → submit → P2P update)",
	func(b operations.Bundle, deps ceremony.CantonDeps, in AddParticipantInput) (AddParticipantOutput, error) {
		ctx := b.GetContext()
		out := AddParticipantOutput{
			State: CeremonyState{Phase: PhaseKeyGen},
		}

		// Validate party ID format up front.
		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("invalid decentralized_party_id %q: expected format <prefix>::<namespace>",
					in.DecentralizedPartyID),
			)
		}
		decNS := parts[1]

		// ── Step 1: New participant generates keys ───────────────────────────
		keyReport, err := operations.ExecuteOperation(b, GenerateNewMemberKeyOp, deps, GenerateNewMemberKeyInput{
			NamespaceName: in.NamespaceName,
			ParticipantID: in.NewParticipantID,
		})
		if err != nil {
			deps.Logger.Infow("New member key generation pending",
				"new_participant", in.NewParticipantID, "err", err)

			return out, fmt.Errorf("%w: new participant has not generated keys yet",
				ErrThresholdNotMet)
		}
		newMember := keyReport.Output
		out.State.NewMemberKeyReady = true

		// ── Step 2: New participant publishes NSD ────────────────────────────
		out.State.Phase = PhaseNSD
		_, err = operations.ExecuteOperation(b, ProposeNewNSDOp, deps, ProposeNewNSDInput{
			ParticipantID:  in.NewParticipantID,
			SigningKeyB64:  newMember.SigningKeyB64,
			Namespace:      newMember.NamespaceFingerprint,
			SynchronizerID: in.SynchronizerID,
		})
		if err != nil {
			deps.Logger.Infow("NSD proposal pending",
				"new_participant", in.NewParticipantID, "err", err)

			return out, fmt.Errorf("%w: new participant has not proposed NSD yet",
				ErrThresholdNotMet)
		}
		out.State.NSDProposed = true

		// Poll until NSD is visible on the synchronizer.
		err = retry.Do(
			func() error {
				exists, qErr := deps.Client.NSDExists(ctx, newMember.NamespaceFingerprint, in.SynchronizerID)
				if qErr != nil {
					return retry.Unrecoverable(fmt.Errorf("checking NSD for new participant: %w", qErr))
				}
				if !exists {
					return fmt.Errorf("namespace delegation not yet visible for %s (namespace %s)",
						in.NewParticipantID, newMember.NamespaceFingerprint)
				}

				return nil
			},
			retry.Context(ctx),
			retry.Attempts(30),
			retry.Delay(500*time.Millisecond),
		)
		if err != nil {
			return out, fmt.Errorf("waiting for new participant NSD: %w", err)
		}

		deps.Logger.Infow("New participant NSD confirmed",
			"namespace", newMember.NamespaceFingerprint)

		// ── Step 3: Read current topology state ──────────────────────────────
		out.State.Phase = PhaseReadState
		stateReport, err := operations.ExecuteOperation(b, ReadCurrentStateOp, deps, ReadCurrentStateInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("read-current-state: %w", err)
		}
		currentState := stateReport.Output

		// ── Sequence-level validation ────────────────────────────────────────
		if slices.Contains(currentState.DNSOwners, newMember.NamespaceFingerprint) {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("new participant namespace fingerprint %q already exists in DNS owners %v",
					newMember.NamespaceFingerprint, currentState.DNSOwners),
			)
		}

		if slices.Contains(currentState.P2PParticipantUIDs, in.NewParticipantID) {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("new participant %q already exists in P2P mapping %v",
					in.NewParticipantID, currentState.P2PParticipantUIDs),
			)
		}

		newThreshold := in.NewThreshold
		if newThreshold <= 0 {
			newThreshold = int(currentState.DNSThreshold)
		}

		if len(currentState.P2PParticipantUIDs) < int(currentState.DNSThreshold) {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf(
					"add is impossible: %d existing participants cannot reach current DNS threshold of %d",
					len(currentState.P2PParticipantUIDs), currentState.DNSThreshold,
				),
			)
		}

		out.State.DNSThreshold = int(currentState.DNSThreshold)

		// ── Step 4: Create add DNS proposal ──────────────────────────────────
		out.State.Phase = PhaseDNSProposal
		proposalReport, err := operations.ExecuteOperation(b, CreateAddDNSProposalOp, deps, CreateAddDNSProposalInput{
			DecentralizedNamespace:  decNS,
			CurrentOwners:           currentState.DNSOwners,
			NewOwnerFingerprint:     newMember.NamespaceFingerprint,
			NewThreshold:            newThreshold,
			CurrentSerial:           int(currentState.DNSSerial),
			ExistingParticipantUIDs: currentState.P2PParticipantUIDs,
			SynchronizerID:          in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("create-add-dns-proposal: %w", err)
		}
		proposal := proposalReport.Output
		out.State.ProposalHash = proposal.ProposalHashSHA256
		out.State.RequiredSigners = proposal.RequiredSigners
		out.State.AllOwners = proposal.NewOwners

		// ── Step 5: Collect DNS signatures from existing participants ─────────
		out.State.Phase = PhaseDNSSigning
		var allSignedTxsB64 []string
		for _, signerUID := range proposal.RequiredSigners {
			sigReport, sigErr := operations.ExecuteOperation(b, SignAddDNSProposalOp, deps, SignAddDNSProposalInput{
				ParticipantID:      signerUID,
				ProposalHashSHA256: proposal.ProposalHashSHA256,
				DNSTxB64:           proposal.DNSTxB64,
				SynchronizerID:     in.SynchronizerID,
			})
			if sigErr != nil {
				deps.Logger.Infow("DNS add signature pending", "signer", signerUID, "err", sigErr)
				out.State.PendingSigners = append(out.State.PendingSigners, signerUID)

				continue
			}
			allSignedTxsB64 = append(allSignedTxsB64, sigReport.Output.SignedDNSTxB64)
			out.State.CollectedSigners = append(out.State.CollectedSigners, signerUID)
		}

		deps.Logger.Infow("Collected add DNS signatures",
			"collected", len(out.State.CollectedSigners), "required", currentState.DNSThreshold,
		)

		if len(out.State.CollectedSigners) < int(currentState.DNSThreshold) {
			return out, fmt.Errorf("%w: %d/%d DNS signatures collected",
				ErrThresholdNotMet, len(out.State.CollectedSigners), currentState.DNSThreshold,
			)
		}

		// ── Step 6: Submit the add DNS update ─────────────────────────────────
		out.State.Phase = PhaseDNSSubmit
		_, err = operations.ExecuteOperation(
			b, SubmitAddDNSOp, deps,
			SubmitAddDNSInput{
				SignedDNSTxsB64: allSignedTxsB64,
				SynchronizerID:  in.SynchronizerID,
				FilterNamespace: decNS,
			},
			operations.WithRetry[SubmitAddDNSInput, ceremony.CantonDeps](),
		)
		if err != nil {
			return out, fmt.Errorf("submit-add-dns: %w", err)
		}

		expectedOwnerCount := len(currentState.DNSOwners) + 1
		err = retry.Do(
			func() error {
				deps.Logger.Infow("Polling add DNS confirmation", "namespace", decNS)
				dnsState, qErr := deps.Client.GetDNS(ctx, decNS, in.SynchronizerID)
				if qErr != nil {
					return fmt.Errorf("polling DNS state: %w", qErr)
				}
				if len(dnsState.Owners) != expectedOwnerCount {
					return fmt.Errorf("DNS update not yet visible: have %d owners, want %d",
						len(dnsState.Owners), expectedOwnerCount)
				}
				deps.Logger.Infow("Add DNS confirmed", "namespace", decNS, "owners", expectedOwnerCount)

				return nil
			},
			retry.Context(ctx),
			retry.Attempts(30),
			retry.Delay(500*time.Millisecond),
		)
		if err != nil {
			return out, fmt.Errorf("waiting for add DNS confirmation: %w", err)
		}

		// ── Step 7: P2P proposals from existing participants + new participant consent ──
		out.State.Phase = PhaseP2P
		out.State.P2PExistingRequired = int(currentState.DNSThreshold)

		allParticipantUIDs := append(append([]string{}, currentState.P2PParticipantUIDs...), newMember.ParticipantUID)

		for _, uid := range currentState.P2PParticipantUIDs {
			_, p2pErr := operations.ExecuteOperation(b, ProposeAddP2POp, deps, ProposeAddP2PInput{
				ParticipantID:      uid,
				PartyID:            in.DecentralizedPartyID,
				AllParticipantUIDs: allParticipantUIDs,
				NewP2PThreshold:    newThreshold,
				CurrentP2PSerial:   int(currentState.P2PSerial),
				SynchronizerID:     in.SynchronizerID,
			})
			if p2pErr != nil {
				deps.Logger.Infow("P2P add proposal pending", "participant", uid, "err", p2pErr)
				continue
			}
			out.State.P2PExistingProposed++
		}

		_, newConsentErr := operations.ExecuteOperation(b, ProposeAddP2POp, deps, ProposeAddP2PInput{
			ParticipantID:      in.NewParticipantID,
			PartyID:            in.DecentralizedPartyID,
			AllParticipantUIDs: allParticipantUIDs,
			NewP2PThreshold:    newThreshold,
			CurrentP2PSerial:   int(currentState.P2PSerial),
			SynchronizerID:     in.SynchronizerID,
		})
		out.State.NewParticipantConsented = newConsentErr == nil
		if newConsentErr != nil {
			deps.Logger.Infow("New participant P2P consent pending",
				"participant", in.NewParticipantID, "err", newConsentErr)
		}

		deps.Logger.Infow("Collected add P2P proposals",
			"existing", out.State.P2PExistingProposed, "existing_required", currentState.DNSThreshold,
			"new_participant_proposed", out.State.NewParticipantConsented,
		)

		if out.State.P2PExistingProposed < int(currentState.DNSThreshold) {
			return out, fmt.Errorf("%w: %d/%d P2P proposals collected from existing participants",
				ErrThresholdNotMet, out.State.P2PExistingProposed, currentState.DNSThreshold,
			)
		}

		if !out.State.NewParticipantConsented {
			return out, fmt.Errorf("%w: new participant has not yet consented to P2P hosting",
				ErrThresholdNotMet)
		}

		// Poll until the updated P2P is confirmed (new participant present).
		err = retry.Do(
			func() error {
				deps.Logger.Infow("Checking add P2P confirmation", "party", in.DecentralizedPartyID)
				p2pState, qErr := deps.Client.GetP2P(ctx, in.DecentralizedPartyID, in.SynchronizerID)
				if qErr != nil {
					return retry.Unrecoverable(fmt.Errorf("polling P2P state: %w", qErr))
				}
				found := false
				for _, p := range p2pState.Participants {
					if p.ParticipantUID == newMember.ParticipantUID {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("new participant %q not yet present in P2P mapping", newMember.ParticipantUID)
				}
				deps.Logger.Infow("Add P2P confirmed", "party", in.DecentralizedPartyID)

				return nil
			},
			retry.Context(ctx),
			retry.Attempts(20),
			retry.Delay(1*time.Second),
		)
		if err != nil {
			return out, fmt.Errorf("waiting for add P2P confirmation: %w", err)
		}

		out.State.Phase = PhaseCompleted
		out.State.PendingSigners = nil
		out.DNSUpdated = true
		out.P2PUpdated = true
		out.NewThreshold = newThreshold
		out.AllOwners = proposal.NewOwners

		return out, nil
	},
)
