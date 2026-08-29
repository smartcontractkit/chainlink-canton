package ledger

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

// CTFClient wraps a CLDF LocalNet participant.
type CTFClient struct {
	participant canton.Participant
}

func NewCTFClient(participant canton.Participant) *CTFClient {
	return &CTFClient{participant: participant}
}

func (c *CTFClient) Participant() canton.Participant {
	return c.participant
}

func (c *CTFClient) ForParty(party string) canton.Participant {
	p := c.participant
	p.PartyID = party
	return p
}

func (c *CTFClient) SubmitCreate(ctx context.Context, actAs string, payload any) (*apiv2.SubmitAndWaitForTransactionResponse, error) {
	tmpl, ok := payload.(interface{ GetTemplateID() string })
	if !ok {
		return nil, fmt.Errorf("payload %T must implement GetTemplateID", payload)
	}

	var createArgs *apiv2.Record
	if creatable, ok := payload.(interface{ CreateCommand() *model.CreateCommand }); ok {
		createArgs = ledger.MapToRecord(creatable.CreateCommand().Arguments)
	} else {
		createArgs = ledger.ConvertToRecord(payload)
	}

	return c.participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId:      contracts.IdentifierFromBinding(tmpl),
						CreateArguments: createArgs,
					},
				},
			}},
			ActAs: []string{actAs},
		},
	})
}

func (c *CTFClient) SubmitExercise(ctx context.Context, actAs string, template interface{ GetTemplateID() string }, contractID, choice string, choiceArg any) (*apiv2.SubmitAndWaitForTransactionResponse, error) {
	return c.SubmitExerciseDisclosed(ctx, actAs, template, contractID, choice, choiceArg, nil)
}

func (c *CTFClient) SubmitExerciseDisclosed(
	ctx context.Context,
	actAs string,
	template interface{ GetTemplateID() string },
	contractID, choice string,
	choiceArg any,
	disclosed *apiv2.DisclosedContract,
) (*apiv2.SubmitAndWaitForTransactionResponse, error) {
	var disclosedContracts []*apiv2.DisclosedContract
	if disclosed != nil {
		disclosedContracts = []*apiv2.DisclosedContract{disclosed}
	}

	return c.participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(template),
						ContractId:     contractID,
						Choice:         choice,
						ChoiceArgument: ledger.MapToValue(choiceArg),
					},
				},
			}},
			ActAs:              []string{actAs},
			DisclosedContracts: disclosedContracts,
		},
	})
}

func (c *CTFClient) SubmitExerciseMulti(
	ctx context.Context,
	actAs []string,
	template interface{ GetTemplateID() string },
	contractID, choice string,
	choiceArg any,
	disclosed []*apiv2.DisclosedContract,
) (*apiv2.SubmitAndWaitForTransactionResponse, error) {
	return c.participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(template),
						ContractId:     contractID,
						Choice:         choice,
						ChoiceArgument: ledger.MapToValue(choiceArg),
					},
				},
			}},
			ActAs:              actAs,
			DisclosedContracts: disclosed,
		},
	})
}
