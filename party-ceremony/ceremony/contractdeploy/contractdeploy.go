// Package contractdeploy implements the contract deployment ceremony for a
// decentralized party using the Canton Ledger API's InteractiveSubmissionService.
//
// # Overview
//
// The ceremony deploys a DAML contract owned by a decentralized party. It follows
// the contract deploy workflow from the Decentralized Party Management spec:
//
//  1. UploadDarsOp         – each participant uploads DARs via Admin PackageService.
//  2. VerifyPartyOp        – verify the decentralized party is visible on the Ledger API.
//  3. PrepareSubmissionOp  – coordinator prepares the contract creation via
//     InteractiveSubmissionService.PrepareSubmission.
//  4. SignSubmissionOp     – each participant signs the prepared hash with their
//     signing key (via [ContractDeployDeps.Signer]).
//  5. ExecuteSubmissionOp  – coordinator aggregates signatures and executes.
//  6. VerifyContractOp     – verify contract exists in Active Contract Set.
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
	retry "github.com/avast/retry-go/v4"

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
//   - If not all participants have uploaded DARs or signed, [ErrThresholdNotMet]
//     is returned so the caller knows to retry after more actors run.
//   - Once all steps succeed, the sequence returns the deployed contract ID.
var ContractDeploySequence = operations.NewSequence(
	"contract-deploy/canton-ceremony/deploy-contract",
	semver.MustParse("1.0.0"),
	"Deploy a DAML contract through a decentralized party (DAR upload → verify → prepare → sign → execute → verify)",
	func(b operations.Bundle, deps ContractDeployDeps, in ContractDeployInput) (ContractDeployOutput, error) {
		// ── Step 1: DAR upload ───────────────────────────────────────────────
		// Each participant uploads all configured DARs to their node.
		uploads := make([]UploadDarsOutput, 0)
		for _, pid := range in.Participants {
			r, err := operations.ExecuteOperation(b, UploadDarsOp, deps, UploadDarsInput{
				ParticipantID: pid,
				Packages:      in.Packages,
			})
			if err != nil {
				if !retry.IsRecoverable(err) {
					return ContractDeployOutput{}, fmt.Errorf("upload-dars %s: %w", pid, err)
				}
				deps.Logger.Infow("DAR upload pending", "participant", pid, "err", err)

				continue
			}
			uploads = append(uploads, r.Output)
		}

		// Gate: all participants must have uploaded DARs before proceeding.
		if len(uploads) < len(in.Participants) {
			deps.Logger.Warnw("Not all participants have uploaded DARs",
				"uploaded", len(uploads), "required", len(in.Participants))

			return ContractDeployOutput{}, fmt.Errorf("%w: %d/%d participants have uploaded DARs",
				ErrThresholdNotMet, len(uploads), len(in.Participants))
		}

		// Collect the first participant's package IDs as the canonical set
		// (all participants upload the same DARs, so package IDs are identical).
		packageIDs := uploads[0].PackageIDs

		deps.Logger.Infow("All DARs uploaded",
			"participants", len(uploads),
			"packages", len(packageIDs),
		)

		// ── Step 2: Verify party ─────────────────────────────────────────────
		_, err := operations.ExecuteOperation(b, VerifyPartyOp, deps, VerifyPartyInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
		})
		if err != nil {
			return ContractDeployOutput{}, fmt.Errorf("verify-party: %w", err)
		}

		// ── Step 3: Prepare submission ───────────────────────────────────────
		// Use the first package ID as the contract's package.
		pkgID := ""
		if len(packageIDs) > 0 {
			pkgID = packageIDs[0]
		}

		prepReport, err := operations.ExecuteOperation(b, PrepareSubmissionOp, deps, PrepareSubmissionInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
			PackageID:            pkgID,
			TemplateModule:       in.TemplateModule,
			TemplateEntity:       in.TemplateEntity,
			ContractArgs:         in.ContractArgs,
		})
		if err != nil {
			return ContractDeployOutput{}, fmt.Errorf("prepare-submission: %w", err)
		}

		prep := prepReport.Output

		deps.Logger.Infow("Submission prepared",
			"hash", prep.PreparedTransactionHash,
		)

		// ── Step 4: Sign submission ──────────────────────────────────────────
		// Each participant signs the prepared transaction hash with their signing key.
		signs := make([]SignSubmissionOutput, 0, len(in.Participants))
		for _, pid := range in.Participants {
			r, signErr := operations.ExecuteOperation(b, SignSubmissionOp, deps, SignSubmissionInput{
				ParticipantID:           pid,
				PreparedTransactionHash: prep.PreparedTransactionHash,
				PreparedTxB64:           prep.PreparedTxB64,
			})
			if signErr != nil {
				if !retry.IsRecoverable(signErr) {
					return ContractDeployOutput{
						PackageIDs:              packageIDs,
						PreparedTransactionHash: prep.PreparedTransactionHash,
					}, fmt.Errorf("sign-submission %s: %w", pid, signErr)
				}
				deps.Logger.Infow("Signing pending", "participant", pid, "err", signErr)

				continue
			}
			signs = append(signs, r.Output)
		}

		// Gate: all participants must have signed before executing.
		if len(signs) < len(in.Participants) {
			deps.Logger.Warnw("Not all participants have signed",
				"signed", len(signs), "required", len(in.Participants))

			return ContractDeployOutput{
					PackageIDs:              packageIDs,
					PreparedTransactionHash: prep.PreparedTransactionHash,
				}, fmt.Errorf("%w: %d/%d participants have signed",
					ErrThresholdNotMet, len(signs), len(in.Participants))
		}

		// Collect all base64-encoded signatures.
		sigsB64 := make([]string, len(signs))
		for i, s := range signs {
			sigsB64[i] = s.SignatureB64
		}

		deps.Logger.Infow("All participants signed",
			"participants", len(signs),
		)

		// ── Step 5: Execute submission ───────────────────────────────────────
		executeReport, err := operations.ExecuteOperation(b, ExecuteSubmissionOp, deps, ExecuteSubmissionInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			PreparedTxB64:        prep.PreparedTxB64,
			SignaturesB64:        sigsB64,
			HashingSchemeVersion: prep.HashingSchemeVersion,
		})
		if err != nil {
			return ContractDeployOutput{
				PackageIDs:              packageIDs,
				PreparedTransactionHash: prep.PreparedTransactionHash,
			}, fmt.Errorf("execute-submission: %w", err)
		}

		// ── Step 6: Verify contract ───────────────────────────────────────────
		verifyReport, err := operations.ExecuteOperation(b, VerifyContractOp, deps, VerifyContractInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			PackageID:            pkgID,
			TemplateModule:       in.TemplateModule,
			TemplateEntity:       in.TemplateEntity,
			ContractID:           executeReport.Output.ContractID,
		})
		if err != nil {
			return ContractDeployOutput{
				PackageIDs:              packageIDs,
				PreparedTransactionHash: prep.PreparedTransactionHash,
			}, fmt.Errorf("verify-contract: %w", err)
		}

		contractID := verifyReport.Output.ContractID

		deps.Logger.Infow("Contract deployed successfully",
			"contract_id", contractID,
			"packages", len(packageIDs),
		)

		return ContractDeployOutput{
			PackageIDs:              packageIDs,
			PreparedTransactionHash: prep.PreparedTransactionHash,
			ContractID:              contractID,
		}, nil
	},
)
