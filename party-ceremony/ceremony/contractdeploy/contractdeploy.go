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
//  4. SignSubmissionOp     – (TODO) each participant signs the prepared hash.
//  5. ExecuteSubmissionOp  – (TODO) coordinator aggregates signatures and executes.
//  6. VerifyContractOp     – (TODO) verify contract exists in Active Contract Set.
//
// Steps 4-6 are not yet implemented because Canton does not provide a
// vault-based transaction signing API. The signing mechanism is deferred to
// future work.
//
// # Async / Resume Pattern
//
// Same pattern as the onboarding and kick ceremonies: all operations are
// idempotent via the Operations framework cache. [ErrThresholdNotMet] is
// returned until all participants have uploaded their DARs.
// [ErrSigningNotImplemented] is returned when the sequence reaches the
// signing step.
package contractdeploy

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [ContractDeploySequence] when not all
// participants have completed their DAR upload. Callers should treat this as
// a "come back later" signal.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must upload DARs")

// ErrSigningNotImplemented is returned by [ContractDeploySequence] when the
// sequence reaches the signing step, which is not yet implemented.
var ErrSigningNotImplemented = errors.New("signing not yet implemented: DAML transaction signing requires external key handling")

// ContractDeploySequence orchestrates the contract deployment ceremony.
// It is designed to be called multiple times by different actors:
//
//   - On each call the sequence re-enters the state machine. Operations whose
//     (definition, input) hash already has a successful report are skipped
//     instantly — they reflect work done by previous actors.
//   - If not all participants have uploaded DARs, [ErrThresholdNotMet] is
//     returned so the caller knows to retry after more actors run.
//   - Once all DARs are uploaded an the party is verified, the coordinator
//     prepares the submission. The sequence then returns
//     [ErrSigningNotImplemented] at the signing step.
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

		// ── Step 4: Sign submission (TODO) ───────────────────────────────────
		// Each participant would sign the prepared transaction hash.
		for _, pid := range in.Participants {
			_, sigErr := operations.ExecuteOperation(b, SignSubmissionOp, deps, SignSubmissionInput{
				ParticipantID:           pid,
				PreparedTransactionHash: prep.PreparedTransactionHash,
				PreparedTxB64:           prep.PreparedTxB64,
			})
			if sigErr != nil {
				// ErrSigningNotImplemented is expected — the sequence stops here.
				return ContractDeployOutput{
					PackageIDs:              packageIDs,
					PreparedTransactionHash: prep.PreparedTransactionHash,
				}, fmt.Errorf("sign-submission %s: %w", pid, sigErr)
			}
		}

		// Steps 5 and 6 would follow here once signing is implemented.

		return ContractDeployOutput{
			PackageIDs:              packageIDs,
			PreparedTransactionHash: prep.PreparedTransactionHash,
		}, nil
	},
)
