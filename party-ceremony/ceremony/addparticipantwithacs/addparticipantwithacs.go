// Package addparticipantwithacs implements the combined add-participant +
// ACS replication ceremony for Canton decentralized parties.
//
// This ceremony extends the standard add-participant flow with ACS (Active
// Contract Set) replication: the new participant receives all pre-existing
// contracts via an offline snapshot rather than relying on Canton's built-in
// party replication. The onboarding flag on the P2P mapping ensures Canton
// defers automatic replication until the manual import is complete.
//
// # Steps
//
//  1. GenerateNewMemberKeyOp  – new participant generates namespace + DAML keys.
//  2. ProposeNewNSDOp         – new participant publishes namespace delegation.
//  3. ReadCurrentStateOp      – read current DNS and P2P topology state.
//  4. CreateAddDNSProposalOp  – coordinator creates updated DNS proposal.
//  5. SignAddDNSProposalOp    – each existing participant signs the proposal.
//  6. SubmitAddDNSOp          – merge signatures, submit updated DNS.
//     1b. RecordTargetLedgerOffsetOp – target records offset before topology changes.
//  7. RecordLedgerOffsetOp    – source records offset BEFORE P2P authorization.
//  8. ProposeAddP2PWithOnboardingOp – P2P mapping with onboarding flag on new participant.
//  9. DisconnectSynchronizerOp – target disconnects after onboarding consent.
//  10. ExportAcsOp            – source exports party ACS at recorded offset.
//  11. ImportAcsOp            – target imports ACS while disconnected.
//  12. ReconnectSynchronizerOp – target reconnects after successful import.
//  13. ClearOnboardingFlagOp   – target clears the onboarding flag.
//
// Steps 1b and 7–13 are participant-scoped: only the designated actor executes
// the handler, but any participant may resume the sequence and pick up cached
// operation reports from the shared reporter.
package addparticipantwithacs

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
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/replication"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [AddParticipantWithAcsSequence] when not
// enough actors have contributed keys, signatures, or proposals.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must resume")

// replicationStepErr classifies an operation result for async multi-actor resume.
// A nil error means the step succeeded (including via a cached report from the
// designated actor). When the current runner is not the designated actor and the
// operation failed, ErrThresholdNotMet is returned so the caller can signal
// that another participant must act. When the runner is the designated actor,
// the underlying error is returned unchanged.
func replicationStepErr(runnerUID, actorUID string, err error) error {
	if err == nil {
		return nil
	}
	if runnerUID != actorUID {
		return ErrThresholdNotMet
	}

	return err
}

// AddParticipantWithAcsSequence orchestrates the full eleven-step combined
// add-participant + ACS replication ceremony. It is designed to be called
// multiple times by different actors in an async workflow.
var AddParticipantWithAcsSequence = operations.NewSequence(
	"add-participant-with-acs/canton-ceremony/decentralized-party",
	semver.MustParse("1.0.0"),
	"Async decentralized party add-participant with ACS replication (topology → export → import → clear flag)",
	func(b operations.Bundle, deps ceremony.CantonDeps, in AddParticipantWithAcsInput) (AddParticipantWithAcsOutput, error) {
		ctx := b.GetContext()
		out := AddParticipantWithAcsOutput{
			State: CeremonyState{Phase: PhaseKeyGen},
		}

		// Validate party ID format.
		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return out, operations.NewUnrecoverableError(
				fmt.Errorf("invalid decentralized_party_id %q: expected format <prefix>::<namespace>",
					in.DecentralizedPartyID),
			)
		}
		decNS := parts[1]

		runnerUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return out, fmt.Errorf("fetching runner participant UID: %w", err)
		}

		// ── Step 1: New participant generates keys ───────────────────────────
		keyReport, err := operations.ExecuteOperation(b, keys.CreateMemberKeyOp, deps, keys.CreateMemberKeyInput{
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

		// ── Step 1b: Record target ledger offset ─────────────────────────────
		// Captured before any new P2P activation reaches the target's ledger,
		// and after the previous activation/kick has already been processed.
		// This is the BeginOffsetExclusive used by ClearPartyOnboardingFlag in
		// Step 11 — forward search from here skips any stale activation
		// (without the onboarding flag) and lands on the new one.
		out.State.Phase = PhaseRecordTargetOffset
		targetOffsetReport, err := operations.ExecuteOperation(b, replication.RecordTargetLedgerOffsetOp, deps, replication.RecordTargetLedgerOffsetInput{
			ParticipantID:  in.NewParticipantID,
			SynchronizerID: in.SynchronizerID,
		})
		if stepErr := replicationStepErr(runnerUID, in.NewParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("Target ledger offset recording pending",
					"runner", runnerUID, "target", in.NewParticipantID)
			} else {
				deps.Logger.Infow("Target ledger offset recording pending",
					"target", in.NewParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: target participant has not recorded ledger offset yet",
				ErrThresholdNotMet)
		}
		out.State.TargetLedgerOffset = targetOffsetReport.Output.LedgerOffset

		// ── Step 2: New participant publishes NSD ────────────────────────────
		out.State.Phase = PhaseNSD
		_, err = operations.ExecuteOperation(b, topology.ProposeNamespaceDelegationOp, deps, topology.ProposeNSDInput{
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
		stateReport, err := operations.ExecuteOperation(b, topology.ReadCurrentStateOp, deps, topology.ReadCurrentStateInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("read-current-state: %w", err)
		}
		currentState := stateReport.Output

		// Validate the add is possible.
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
		proposalReport, err := operations.ExecuteOperation(b, topology.CreateAddDNSProposalOp, deps, topology.CreateAddDNSProposalInput{
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
			sigReport, sigErr := operations.ExecuteOperation(b, topology.SignDNSProposalOp, deps, topology.SignDNSProposalInput{
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
			b, topology.SubmitDNSOp, deps,
			topology.SubmitDNSInput{
				SignedDNSTxsB64: allSignedTxsB64,
				SynchronizerID:  in.SynchronizerID,
				FilterNamespace: decNS,
			},
			operations.WithRetry[topology.SubmitDNSInput, ceremony.CantonDeps](),
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

		// ── Step 7: Record ledger offset on source BEFORE P2P auth ──────────
		out.State.Phase = PhaseRecordOffset
		offsetReport, err := operations.ExecuteOperation(b, replication.RecordLedgerOffsetOp, deps, replication.RecordLedgerOffsetInput{
			ParticipantID:  in.SourceParticipantID,
			SynchronizerID: in.SynchronizerID,
		})
		if stepErr := replicationStepErr(runnerUID, in.SourceParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("Ledger offset recording pending",
					"runner", runnerUID, "source", in.SourceParticipantID)
			} else {
				deps.Logger.Infow("Ledger offset recording pending",
					"source", in.SourceParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: source participant has not recorded ledger offset yet",
				ErrThresholdNotMet)
		}
		out.State.LedgerOffsetRecorded = true
		ledgerOffset := offsetReport.Output.LedgerOffset
		ledgerTimestamp := offsetReport.Output.TimestampSeconds

		// ── Step 8: P2P proposals with onboarding flag ──────────────────────
		out.State.Phase = PhaseP2POnboarding
		out.State.P2PExistingRequired = int(currentState.DNSThreshold)

		allParticipantUIDs := append(append([]string{}, currentState.P2PParticipantUIDs...), newMember.ParticipantUID)
		var partySigningKeysB64 []string
		if len(currentState.PartySigningKeysB64) > 0 {
			partySigningKeysB64 = append([]string{}, currentState.PartySigningKeysB64...)
			partySigningKeysB64 = append(partySigningKeysB64, newMember.DamlKeyB64)
		}

		for _, uid := range currentState.P2PParticipantUIDs {
			_, p2pErr := operations.ExecuteOperation(b, topology.ProposeAddP2PWithOnboardingOp, deps, topology.ProposeAddP2PWithOnboardingInput{
				ParticipantID:       uid,
				PartyID:             in.DecentralizedPartyID,
				AllParticipantUIDs:  allParticipantUIDs,
				NewParticipantUID:   in.NewParticipantID,
				NewP2PThreshold:     newThreshold,
				CurrentP2PSerial:    int(currentState.P2PSerial),
				SynchronizerID:      in.SynchronizerID,
				PartySigningKeysB64: partySigningKeysB64,
			})
			if p2pErr != nil {
				deps.Logger.Infow("P2P add proposal with onboarding pending", "participant", uid, "err", p2pErr)
				continue
			}
			out.State.P2PExistingProposed++
		}

		// New participant must also consent.
		_, newConsentErr := operations.ExecuteOperation(b, topology.ProposeAddP2PWithOnboardingOp, deps, topology.ProposeAddP2PWithOnboardingInput{
			ParticipantID:       in.NewParticipantID,
			PartyID:             in.DecentralizedPartyID,
			AllParticipantUIDs:  allParticipantUIDs,
			NewParticipantUID:   in.NewParticipantID,
			NewP2PThreshold:     newThreshold,
			CurrentP2PSerial:    int(currentState.P2PSerial),
			SynchronizerID:      in.SynchronizerID,
			PartySigningKeysB64: partySigningKeysB64,
		})
		out.State.NewParticipantConsented = newConsentErr == nil
		if newConsentErr != nil {
			deps.Logger.Infow("New participant P2P consent pending",
				"participant", in.NewParticipantID, "err", newConsentErr)
		}

		deps.Logger.Infow("Collected add P2P proposals with onboarding",
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

		// Poll until the updated P2P is confirmed.
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

		out.DNSUpdated = true
		out.P2PUpdated = true

		// ── Step 9: Disconnect target before ACS export/import ───────────────
		// Per Canton's offline party replication flow, the target must go
		// offline immediately after authorizing hosting with the onboarding
		// flag set. That prevents it from processing post-snapshot events until
		// the source ACS has been imported.
		out.State.Phase = PhaseTargetDisconnect
		_, err = operations.ExecuteOperation(b, replication.DisconnectSynchronizerOp, deps, replication.DisconnectSynchronizerInput{
			ParticipantID:     in.NewParticipantID,
			SynchronizerAlias: in.SynchronizerAlias,
		})
		if stepErr := replicationStepErr(runnerUID, in.NewParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("Target disconnect pending",
					"runner", runnerUID, "target", in.NewParticipantID)
			} else {
				deps.Logger.Infow("Target disconnect pending",
					"target", in.NewParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: target participant has not disconnected for ACS replication yet",
				ErrThresholdNotMet)
		}
		out.State.TargetDisconnected = true

		// ── Step 10: Export ACS from source ──────────────────────────────────
		out.State.Phase = PhaseAcsExport
		exportReport, err := operations.ExecuteOperation(b, replication.ExportAcsOp, deps, replication.ExportAcsInput{
			ParticipantID:    in.SourceParticipantID,
			PartyIDs:         []string{in.DecentralizedPartyID},
			SynchronizerID:   in.SynchronizerID,
			LedgerOffset:     ledgerOffset,
			TimestampSeconds: ledgerTimestamp,
		})
		if stepErr := replicationStepErr(runnerUID, in.SourceParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("ACS export pending",
					"runner", runnerUID, "source", in.SourceParticipantID)
			} else {
				deps.Logger.Infow("ACS export pending",
					"source", in.SourceParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: source participant has not exported ACS yet",
				ErrThresholdNotMet)
		}
		out.State.AcsExported = true

		// ── Step 11: Import ACS into target ──────────────────────────────────
		out.State.Phase = PhaseAcsImport
		_, err = operations.ExecuteOperation(b, replication.ImportAcsOp, deps, replication.ImportAcsInput{
			ParticipantID:     in.NewParticipantID,
			SynchronizerID:    in.SynchronizerID,
			SynchronizerAlias: in.SynchronizerAlias,
			AcsSnapshotB64:    exportReport.Output.AcsSnapshotB64,
			ExpectedSHA256:    exportReport.Output.SHA256,
		})
		if stepErr := replicationStepErr(runnerUID, in.NewParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("ACS import pending",
					"runner", runnerUID, "target", in.NewParticipantID)
			} else {
				deps.Logger.Infow("ACS import pending",
					"target", in.NewParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: target participant has not imported ACS yet",
				ErrThresholdNotMet)
		}
		out.State.AcsImported = true

		// ── Step 12: Reconnect target ────────────────────────────────────────
		out.State.Phase = PhaseTargetReconnect
		_, err = operations.ExecuteOperation(b, replication.ReconnectSynchronizerOp, deps, replication.ReconnectSynchronizerInput{
			ParticipantID:     in.NewParticipantID,
			SynchronizerAlias: in.SynchronizerAlias,
		})
		if stepErr := replicationStepErr(runnerUID, in.NewParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("Target reconnect pending",
					"runner", runnerUID, "target", in.NewParticipantID)
			} else {
				deps.Logger.Infow("Target reconnect pending",
					"target", in.NewParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: target participant has not reconnected after ACS import yet",
				ErrThresholdNotMet)
		}
		out.State.TargetReconnected = true

		// ── Step 13: Clear onboarding flag ───────────────────────────────────
		// Use the target ledger offset recorded in Step 1b. Forward search
		// from here skips past any prior activation of this party on the
		// target (e.g. before an earlier kick) and lands on the new one,
		// which carries the onboarding flag.
		out.State.Phase = PhaseClearOnboarding
		_, err = operations.ExecuteOperation(b, replication.ClearOnboardingFlagOp, deps, replication.ClearOnboardingFlagInput{
			ParticipantID:        in.NewParticipantID,
			PartyID:              in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
			BeginOffsetExclusive: out.State.TargetLedgerOffset,
		})
		if stepErr := replicationStepErr(runnerUID, in.NewParticipantID, err); stepErr != nil {
			if errors.Is(stepErr, ErrThresholdNotMet) {
				deps.Logger.Infow("Clear onboarding flag pending",
					"runner", runnerUID, "target", in.NewParticipantID)
			} else {
				deps.Logger.Infow("Clear onboarding flag pending",
					"target", in.NewParticipantID, "err", stepErr)
			}

			return out, fmt.Errorf("%w: target participant has not cleared onboarding flag yet",
				ErrThresholdNotMet)
		}
		out.State.OnboardingFlagCleared = true

		out.State.Phase = PhaseCompleted
		out.State.PendingSigners = nil
		out.AcsImported = true
		out.NewThreshold = newThreshold
		out.AllOwners = proposal.NewOwners

		return out, nil
	},
)
