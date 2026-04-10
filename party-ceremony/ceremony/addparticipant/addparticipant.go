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

		// Validate party ID format up front.
		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return AddParticipantOutput{}, operations.NewUnrecoverableError(
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

			return AddParticipantOutput{}, fmt.Errorf("%w: new participant has not generated keys yet",
				ErrThresholdNotMet)
		}
		newMember := keyReport.Output

		// ── Step 2: New participant publishes NSD ────────────────────────────
		_, err = operations.ExecuteOperation(b, ProposeNewNSDOp, deps, ProposeNewNSDInput{
			ParticipantID:  in.NewParticipantID,
			SigningKeyB64:  newMember.SigningKeyB64,
			Namespace:      newMember.NamespaceFingerprint,
			SynchronizerID: in.SynchronizerID,
		})
		if err != nil {
			deps.Logger.Infow("NSD proposal pending",
				"new_participant", in.NewParticipantID, "err", err)

			return AddParticipantOutput{}, fmt.Errorf("%w: new participant has not proposed NSD yet",
				ErrThresholdNotMet)
		}

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
			return AddParticipantOutput{}, fmt.Errorf("waiting for new participant NSD: %w", err)
		}

		deps.Logger.Infow("New participant NSD confirmed",
			"namespace", newMember.NamespaceFingerprint)

		// ── Step 3: Read current topology state ──────────────────────────────
		stateReport, err := operations.ExecuteOperation(b, ReadCurrentStateOp, deps, ReadCurrentStateInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return AddParticipantOutput{}, fmt.Errorf("read-current-state: %w", err)
		}
		currentState := stateReport.Output

		// ── Sequence-level validation ────────────────────────────────────────
		// Verify the new participant is NOT already in DNS owners.
		if slices.Contains(currentState.DNSOwners, newMember.NamespaceFingerprint) {
			return AddParticipantOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("new participant namespace fingerprint %q already exists in DNS owners %v",
					newMember.NamespaceFingerprint, currentState.DNSOwners),
			)
		}

		// Verify the new participant is NOT already in P2P.
		if slices.Contains(currentState.P2PParticipantUIDs, in.NewParticipantID) {
			return AddParticipantOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("new participant %q already exists in P2P mapping %v",
					in.NewParticipantID, currentState.P2PParticipantUIDs),
			)
		}

		// Compute post-add threshold: keep current unless overridden.
		newThreshold := in.NewThreshold
		if newThreshold <= 0 {
			newThreshold = int(currentState.DNSThreshold)
		}

		// Verify existing participants can reach the current DNS threshold.
		if len(currentState.P2PParticipantUIDs) < int(currentState.DNSThreshold) {
			return AddParticipantOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf(
					"add is impossible: %d existing participants cannot reach current DNS threshold of %d",
					len(currentState.P2PParticipantUIDs), currentState.DNSThreshold,
				),
			)
		}

		// ── Step 4: Create add DNS proposal ──────────────────────────────────
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
			return AddParticipantOutput{}, fmt.Errorf("create-add-dns-proposal: %w", err)
		}
		proposal := proposalReport.Output

		// ── Step 5: Collect DNS signatures from existing participants ─────────
		var allSignedTxsB64 []string
		var dnsSigCount int
		for _, signerUID := range proposal.RequiredSigners {
			sigReport, sigErr := operations.ExecuteOperation(b, SignAddDNSProposalOp, deps, SignAddDNSProposalInput{
				ParticipantID:      signerUID,
				ProposalHashSHA256: proposal.ProposalHashSHA256,
				DNSTxB64:           proposal.DNSTxB64,
				SynchronizerID:     in.SynchronizerID,
			})
			if sigErr != nil {
				deps.Logger.Infow("DNS add signature pending", "signer", signerUID, "err", sigErr)
				continue
			}
			allSignedTxsB64 = append(allSignedTxsB64, sigReport.Output.SignedDNSTxB64)
			dnsSigCount++
		}

		deps.Logger.Infow("Collected add DNS signatures",
			"collected", dnsSigCount, "required", currentState.DNSThreshold,
		)

		// Gate: Canton requires currentThreshold signatures for serial > 1.
		if dnsSigCount < int(currentState.DNSThreshold) {
			return AddParticipantOutput{}, fmt.Errorf("%w: %d/%d DNS signatures collected",
				ErrThresholdNotMet, dnsSigCount, currentState.DNSThreshold,
			)
		}

		// ── Step 6: Submit the add DNS update ─────────────────────────────────
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
			return AddParticipantOutput{}, fmt.Errorf("submit-add-dns: %w", err)
		}

		// Poll until the updated DNS is visible (owner count increased).
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
			return AddParticipantOutput{}, fmt.Errorf("waiting for add DNS confirmation: %w", err)
		}

		// ── Step 7: P2P proposals from existing participants + new participant consent ──
		// Existing members propose the updated P2P mapping for namespace authority
		// (Canton requires threshold-of-current-owners signatures).
		// The new participant must ALSO call Authorize to consent to hosting the party
		// on their node — Canton requires participant consent independently of the
		// namespace threshold, same as in the onboarding ceremony.

		// Build full participant UID list (existing + new).
		allParticipantUIDs := append(append([]string{}, currentState.P2PParticipantUIDs...), newMember.ParticipantUID)

		var existingProposedCount int
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
			existingProposedCount++
		}

		// New participant must consent to hosting by proposing the same P2P mapping.
		_, newConsentErr := operations.ExecuteOperation(b, ProposeAddP2POp, deps, ProposeAddP2PInput{
			ParticipantID:      in.NewParticipantID,
			PartyID:            in.DecentralizedPartyID,
			AllParticipantUIDs: allParticipantUIDs,
			NewP2PThreshold:    newThreshold,
			CurrentP2PSerial:   int(currentState.P2PSerial),
			SynchronizerID:     in.SynchronizerID,
		})
		newParticipantProposed := newConsentErr == nil
		if newConsentErr != nil {
			deps.Logger.Infow("New participant P2P consent pending",
				"participant", in.NewParticipantID, "err", newConsentErr)
		}

		deps.Logger.Infow("Collected add P2P proposals",
			"existing", existingProposedCount, "existing_required", currentState.DNSThreshold,
			"new_participant_proposed", newParticipantProposed,
		)

		if existingProposedCount < int(currentState.DNSThreshold) {
			return AddParticipantOutput{}, fmt.Errorf("%w: %d/%d P2P proposals collected from existing participants",
				ErrThresholdNotMet, existingProposedCount, currentState.DNSThreshold,
			)
		}

		if !newParticipantProposed {
			return AddParticipantOutput{}, fmt.Errorf("%w: new participant has not yet consented to P2P hosting",
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
			return AddParticipantOutput{}, fmt.Errorf("waiting for add P2P confirmation: %w", err)
		}

		return AddParticipantOutput{
			DNSUpdated:   true,
			P2PUpdated:   true,
			NewThreshold: newThreshold,
			AllOwners:    proposal.NewOwners,
		}, nil
	},
)
