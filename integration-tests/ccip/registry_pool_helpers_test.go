package tests

import (
	"context"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/model"
	damlledger "github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

// submitCreateAndExercise submits a single CreateAndExercise command: payload is
// created and choice is exercised on the resulting contract in one atomic
// transaction, authorized by actAs. This is the real usage pattern Initialize is
// designed for - a caller never needs to see the pool's contractId, since it's
// created and exercised on in the same command.
//
// disclosedContracts must include any contract the choice body touches that the
// submitting participant doesn't already have in its own ACS (e.g. a cross-participant
// TAR) - authorization (actAs) and visibility are separate concerns: a party's
// signature can satisfy a choice's controller check without the submitting
// participant being able to see the contract at all.
func submitCreateAndExercise(
	ctx context.Context,
	participant canton.Participant,
	actAs string,
	payload interface {
		model.CreateCommander
		GetTemplateID() string
	},
	choice string,
	choiceArg any,
	disclosedContracts ...*apiv2.DisclosedContract,
) (*apiv2.SubmitAndWaitForTransactionResponse, error) {
	createCmd := payload.CreateCommand()
	return participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			ActAs:     []string{actAs},
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_CreateAndExercise{
					CreateAndExercise: &apiv2.CreateAndExerciseCommand{
						TemplateId:      contracts.IdentifierFromBinding(payload),
						CreateArguments: damlledger.MapToRecord(createCmd.Arguments),
						Choice:          choice,
						ChoiceArgument:  damlledger.MapToValue(choiceArg),
					},
				},
			}},
			DisclosedContracts: disclosedContracts,
		},
	})
}
