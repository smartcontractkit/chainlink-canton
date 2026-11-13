// Package exercisechoice exercises a choice on an existing contract as a
// decentralized party. It is the contract deploy ceremony without the DAR upload.
package exercisechoice

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned by [ExerciseChoiceSequence] when not all
// participants have completed their step.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must complete their step")

// ExerciseChoiceSequence orchestrates the exercise choice ceremony.
var ExerciseChoiceSequence = operations.NewSequence(
	"exercise-choice/canton-ceremony/exercise-choice",
	semver.MustParse("1.0.0"),
	"Exercise a choice through a decentralized party (verify → fetch participants → grant rights → prepare → sign → execute)",
	func(b operations.Bundle, deps ledger.ContractDeployDeps, in ExerciseChoiceInput) (ExerciseChoiceOutput, error) {
		out := ExerciseChoiceOutput{
			State: CeremonyState{Phase: PhaseVerifyParty},
		}

		_, err := operations.ExecuteOperation(b, ledger.VerifyPartyOp, deps, ledger.VerifyPartyInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
		})
		if err != nil {
			return out, fmt.Errorf("verify-party: %w", err)
		}

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
		out.State.SignRequired = len(participants)

		out.State.Phase = PhaseGrantRights
		for _, pid := range participants {
			if _, grantErr := operations.ExecuteOperation(b, ledger.GrantPartyRightsOp, deps, ledger.GrantPartyRightsInput{
				ParticipantID:        pid,
				DecentralizedPartyID: in.DecentralizedPartyID,
			}); grantErr != nil {
				deps.Logger.Infow("Party rights grant pending", "participant", pid, "err", grantErr)

				continue
			}
			out.State.Granted = append(out.State.Granted, pid)
		}
		if len(out.State.Granted) < len(participants) {
			return out, fmt.Errorf("%w: %d/%d participants have granted party rights",
				ErrThresholdNotMet, len(out.State.Granted), len(participants))
		}

		out.State.Phase = PhasePrepare
		prepReport, err := operations.ExecuteOperation(b, ledger.PrepareExerciseOp, deps, ledger.PrepareExerciseInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
			PackageID:            in.PackageID,
			TemplateModule:       in.TemplateModule,
			TemplateEntity:       in.TemplateEntity,
			ContractID:           in.ContractID,
			Choice:               in.Choice,
			ChoiceArgument:       in.ChoiceArgument,
		})
		if err != nil {
			return out, fmt.Errorf("prepare-exercise: %w", err)
		}
		prep := prepReport.Output
		out.PreparedTransactionHash = prep.PreparedTransactionHash
		out.State.PreparedTxHash = prep.PreparedTransactionHash

		out.State.Phase = PhaseSigning
		sigsB64 := make([]string, 0, len(participants))
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
			sigsB64 = append(sigsB64, r.Output.SignatureB64)
			out.State.Signed = append(out.State.Signed, pid)
		}
		if len(sigsB64) < len(participants) {
			return out, fmt.Errorf("%w: %d/%d participants have signed",
				ErrThresholdNotMet, len(sigsB64), len(participants))
		}

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

		out.UpdateID = executeReport.Output.UpdateID
		out.CreatedContracts = executeReport.Output.CreatedContracts
		out.State.Phase = PhaseCompleted

		deps.Logger.Infow("Choice exercised",
			"choice", in.Choice,
			"update_id", out.UpdateID,
			"created", len(out.CreatedContracts),
		)

		return out, nil
	},
)
