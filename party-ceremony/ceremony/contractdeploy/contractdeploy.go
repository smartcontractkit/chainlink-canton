// Package contractdeploy implements the contract deployment ceremony for a
// decentralized party using the Canton Ledger API's InteractiveSubmissionService.
//
// # Overview
//
// The ceremony deploys a DAML contract owned by a decentralized party. It follows
// the contract deploy workflow from the Decentralized Party Management spec:
//
//  1. VerifyPartyOp        – verify the decentralized party is visible on the Ledger API.
//  2. FetchParticipantsOp  – fetch the hosting participant UIDs from the Canton topology
//     store (PartyToParticipant mapping). The list drives steps 3 and 4.
//  3. UploadDarsOp         – each participant uploads DARs via Admin PackageService.
//  4. PrepareSubmissionOp  – coordinator prepares the contract creation via
//     InteractiveSubmissionService.PrepareSubmission.
//  5. SignSubmissionOp     – each participant signs the prepared hash with their
//     signing key (via [ContractDeployDeps.Signer]).
//  6. ExecuteSubmissionOp  – coordinator aggregates signatures and executes.
//  7. VerifyContractOp     – verify contract exists in Active Contract Set.
//
// # Async / Resume Pattern
//
// Same pattern as the onboarding and kick ceremonies: all operations are
// idempotent via the Operations framework cache. [ErrThresholdNotMet] is
// returned until all participants have uploaded their DARs and signed the
// transaction. The sequence completes when all steps succeed.
package contractdeploy

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [ContractDeploySequence] when not all
// participants have completed their step. Callers should treat this as
// a "come back later" signal.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must complete their step")

// ContractDeploySequence orchestrates the contract deployment ceremony.
// It is designed to be called multiple times by different actors:
//
//   - On each call the sequence re-enters the state machine. Operations whose
//     (definition, input) hash already has a successful report are skipped
//     instantly — they reflect work done by previous actors.
//   - Participants are resolved dynamically from the Canton topology store via
//     [FetchParticipantsOp];
//   - If not all participants have uploaded DARs or signed, [ErrThresholdNotMet]
//     is returned so the caller knows to retry after more actors run.
//   - Once all steps succeed, the sequence returns the deployed contract ID.
var ContractDeploySequence = operations.NewSequence(
	"contract-deploy/canton-ceremony/deploy-contract",
	semver.MustParse("1.0.0"),
	"Deploy a DAML contract through a decentralized party (verify → fetch participants → DAR upload → prepare → sign → execute → verify)",
	func(b operations.Bundle, deps ledger.ContractDeployDeps, in ContractDeployInput) (ContractDeployOutput, error) {
		out := ContractDeployOutput{
			State: CeremonyState{Phase: PhaseVerifyParty},
		}

		// ── Step 1: Verify party ─────────────────────────────────────────────
		_, err := operations.ExecuteOperation(b, ledger.VerifyPartyOp, deps, ledger.VerifyPartyInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
		})
		if err != nil {
			return out, fmt.Errorf("verify-party: %w", err)
		}

		// ── Step 2: Fetch participants ────────────────────────────────────────
		out.State.Phase = PhaseFetchMembers
		fetchReport, err := operations.ExecuteOperation(b, ledger.FetchParticipantsOp, deps, ledger.FetchParticipantsInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("fetch-participants: %w", err)
		}

		participants := fetchReport.Output.Participants
		out.State.Participants = participants
		out.State.DARsRequired = len(participants)
		out.State.SignRequired = len(participants)

		deps.Logger.Infow("Participants resolved from topology",
			"party", in.DecentralizedPartyID,
			"count", len(participants),
		)

		// ── Step 3: Grant party rights + DAR upload ─────────────────────────
		out.State.Phase = PhaseDARUpload
		uploads := make([]ledger.UploadDarsOutput, 0)
		grants := 0
		for _, pid := range participants {
			// Grant actAs/readAs rights for this participant's user. In no-auth
			// environments (deps.UserID=="") this is a no-op that always succeeds.
			if _, grantErr := operations.ExecuteOperation(b, ledger.GrantPartyRightsOp, deps, ledger.GrantPartyRightsInput{
				ParticipantID:        pid,
				DecentralizedPartyID: in.DecentralizedPartyID,
			}); grantErr != nil {
				deps.Logger.Infow("Party rights grant pending", "participant", pid, "err", grantErr)
			} else {
				grants++
			}

			r, uploadErr := operations.ExecuteOperation(b, ledger.UploadDarsOp, deps, ledger.UploadDarsInput{
				ParticipantID: pid,
				Packages:      in.Packages,
			})
			if uploadErr != nil {
				deps.Logger.Infow("DAR upload pending", "participant", pid, "err", uploadErr)

				continue
			}
			uploads = append(uploads, r.Output)
			out.State.DARsUploaded = append(out.State.DARsUploaded, pid)
		}

		if len(uploads) < len(participants) {
			deps.Logger.Warnw("Not all participants have uploaded DARs",
				"uploaded", len(uploads), "required", len(participants))

			return out, fmt.Errorf("%w: %d/%d participants have uploaded DARs",
				ErrThresholdNotMet, len(uploads), len(participants))
		}

		if grants < len(participants) {
			deps.Logger.Warnw("Not all participants have granted party rights",
				"granted", grants, "required", len(participants))

			return out, fmt.Errorf("%w: %d/%d participants have granted party rights",
				ErrThresholdNotMet, grants, len(participants))
		}

		packageIDs := uploads[0].PackageIDs
		out.PackageIDs = packageIDs

		deps.Logger.Infow("All DARs uploaded",
			"participants", len(uploads),
			"packages", len(packageIDs),
		)

		// ── Step 4: Prepare submission ───────────────────────────────────────
		out.State.Phase = PhasePrepare
		pkgID := ""
		if len(packageIDs) > 0 {
			pkgID = packageIDs[0]
		}

		prepReport, err := operations.ExecuteOperation(b, ledger.PrepareSubmissionOp, deps, ledger.PrepareSubmissionInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
			PackageID:            pkgID,
			TemplateModule:       in.TemplateModule,
			TemplateEntity:       in.TemplateEntity,
			ContractArgs:         in.ContractArgs,
		})
		if err != nil {
			return out, fmt.Errorf("prepare-submission: %w", err)
		}

		prep := prepReport.Output
		out.PreparedTransactionHash = prep.PreparedTransactionHash
		out.State.PreparedTxHash = prep.PreparedTransactionHash

		deps.Logger.Infow("Submission prepared",
			"hash", prep.PreparedTransactionHash,
		)

		// ── Step 5: Sign submission ──────────────────────────────────────────
		out.State.Phase = PhaseSigning
		signs := make([]ledger.SignSubmissionOutput, 0, len(participants))
		for _, pid := range participants {
			r, signErr := operations.ExecuteOperation(b, ledger.SignSubmissionOp, deps, ledger.SignSubmissionInput{
				ParticipantID:           pid,
				PreparedTransactionHash: prep.PreparedTransactionHash,
				PreparedTxB64:           prep.PreparedTxB64,
				KnownSigningKeysB64:     fetchReport.Output.PartySigningKeysB64,
			})
			if signErr != nil {
				deps.Logger.Infow("Signing pending", "participant", pid, "err", signErr)

				continue
			}
			signs = append(signs, r.Output)
			out.State.Signed = append(out.State.Signed, pid)
		}

		if len(signs) < len(participants) {
			deps.Logger.Warnw("Not all participants have signed",
				"signed", len(signs), "required", len(participants))

			return out, fmt.Errorf("%w: %d/%d participants have signed",
				ErrThresholdNotMet, len(signs), len(participants))
		}

		sigsB64 := make([]string, len(signs))
		for i, s := range signs {
			sigsB64[i] = s.SignatureB64
		}

		deps.Logger.Infow("All participants signed",
			"participants", len(signs),
		)

		// ── Step 6: Execute submission ───────────────────────────────────────
		out.State.Phase = PhaseExecute
		executeReport, err := operations.ExecuteOperation(b, ledger.ExecuteSubmissionOp, deps, ledger.ExecuteSubmissionInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			PreparedTxB64:        prep.PreparedTxB64,
			SignaturesB64:        sigsB64,
			HashingSchemeVersion: prep.HashingSchemeVersion,
		})
		if err != nil {
			return out, fmt.Errorf("execute-submission: %w", err)
		}

		// ── Step 7: Verify contract ───────────────────────────────────────────
		out.State.Phase = PhaseVerifyContract
		verifyReport, err := operations.ExecuteOperation(b, ledger.VerifyContractOp, deps, ledger.VerifyContractInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			PackageID:            pkgID,
			TemplateModule:       in.TemplateModule,
			TemplateEntity:       in.TemplateEntity,
			ContractID:           executeReport.Output.ContractID,
		})
		if err != nil {
			return out, fmt.Errorf("verify-contract: %w", err)
		}

		contractID := verifyReport.Output.ContractID

		deps.Logger.Infow("Contract deployed successfully",
			"contract_id", contractID,
			"packages", len(packageIDs),
		)

		out.State.Phase = PhaseCompleted
		out.ContractID = contractID

		return out, nil
	},
)
