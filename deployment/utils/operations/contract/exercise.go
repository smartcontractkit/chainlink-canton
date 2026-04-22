package contract

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type ExecInfo struct {
	// The ID of the update that was confirmed as the result of this execution.
	UpdateID string `json:"update_id"`
}

type ExerciseOutput struct {
	ChainSelector uint64                 `json:"chainSelector"`
	Tx            mcms_types.Transaction `json:"tx"`
	ExecInfo      *ExecInfo              `json:"execInfo,omitempty"`
}

func (o ExerciseOutput) Executed() bool {
	return o.ExecInfo != nil
}

// DisclosedContract copies the regenerated protobuf message.
// The generated protobufs cannot be used as an operation's input since they contain unexported fields.
type DisclosedContract struct {
	TemplateId       contracts.TemplateID `json:"template_id"`
	ContractId       string               `json:"contract_id"`
	CreatedEventBlob []byte               `json:"created_event_blob"`
	SynchronizerId   string               `json:"synchronizer_id"`
}

func DisclosedContractsFromProto(dcs []*apiv2.DisclosedContract) []DisclosedContract {
	disclosedContracts := make([]DisclosedContract, len(dcs))
	for i, dc := range dcs {
		disclosedContracts[i] = DisclosedContract{
			TemplateId: contracts.TemplateID{
				PackageID:  dc.TemplateId.PackageId,
				ModuleName: dc.TemplateId.ModuleName,
				EntityName: dc.TemplateId.EntityName,
			},
			ContractId:       dc.ContractId,
			CreatedEventBlob: dc.CreatedEventBlob,
			SynchronizerId:   dc.SynchronizerId,
		}
	}

	return disclosedContracts
}

type ChoiceInput[ARGS any] struct {
	// The InstanceAddress this operation is targeting. Will be resolved to an active contract.
	InstanceAddress contracts.InstanceAddress `json:"instanceAddress"`
	// RawInstanceAddress is the "instanceId@partyId" format required by the Canton MCMS SDK
	// for AdditionalFields.TargetInstanceAddress. Must be set when MCMSEnabled is true.
	RawInstanceAddress string `json:"rawInstanceAddress,omitempty"`
	Args               ARGS   `json:"args"`
	MCMSEnabled        bool   `json:"mcmsEnabled,omitempty"`
	DisclosedContracts []DisclosedContract
}

type ExerciseParams[ARGS any] struct {
	// Name is the name of the operation.
	Name string
	// Version is the version of the operation.
	Version *semver.Version
	// Description is a brief description of the operation.
	Description string
	// ContractType is the type of the target contract.
	ContractType deployment.ContractType
	// Modifier can be used to modify the argument for this exercise.
	Modifier func(chain canton.Chain, input ARGS) (ARGS, error)
	// Validate is an optional function to validate the input arguments.
	Validate func(input ARGS) error

	// Template is the binding struct of the target contract.
	Template common.Template
	// Method is the bindings method to call the choice.
	Method func(contractID string, args ARGS) *model.ExerciseCommand

	// EncodeMethod encodes the choice args to hex for MCMS proposals.
	// Uses the binding's Encoder (e.g., globalConfig.Encoder().ApplyDestChainConfigUpdates).
	// When nil, MCMS encoding is not available for this operation.
	EncodeMethod func(args ARGS) (*bind.EncodedChoice, error)
}

func NewExercise[ARGS any](params ExerciseParams[ARGS]) *operations.Operation[ChoiceInput[ARGS], ExerciseOutput, canton.Chain] {
	return operations.NewOperation(
		params.Name,
		params.Version,
		params.Description,
		func(b operations.Bundle, deps canton.Chain, input ChoiceInput[ARGS]) (ExerciseOutput, error) {
			if params.Validate != nil {
				if err := params.Validate(input.Args); err != nil {
					return ExerciseOutput{}, fmt.Errorf("validate input: %w", err)
				}
			}
			if params.Modifier != nil {
				var err error
				input.Args, err = params.Modifier(deps, input.Args)
				if err != nil {
					return ExerciseOutput{}, fmt.Errorf("failed to modify input: %w", err)
				}
			}

			// If MCMS enabled, encode and return without executing on-chain
			if input.MCMSEnabled {
				if params.EncodeMethod == nil {
					return ExerciseOutput{}, fmt.Errorf("MCMSEnabled is true but no EncodeMethod is defined for operation %s", params.Name)
				}
				if input.RawInstanceAddress == "" {
					return ExerciseOutput{}, fmt.Errorf("MCMSEnabled is true but RawInstanceAddress is empty for operation %s", params.Name)
				}
				encodedChoice, err := params.EncodeMethod(input.Args)
				if err != nil {
					return ExerciseOutput{}, fmt.Errorf("failed to encode choice args for MCMS: %w", err)
				}
				mcmsTx, err := NewCantonTransaction(input.RawInstanceAddress, input.InstanceAddress, encodedChoice, params.ContractType, params.Template.GetTemplateID())
				if err != nil {
					return ExerciseOutput{}, fmt.Errorf("failed to build MCMS transaction: %w", err)
				}

				return ExerciseOutput{
					ChainSelector: deps.ChainSelector(),
					Tx:            mcmsTx,
				}, nil
			}

			// Direct execution path
			participant := deps.Participants[0]

			contractID, err := FindActiveContractIDByInstanceAddress(b.GetContext(), participant.LedgerServices.State, participant.PartyID, params.Template.GetTemplateID(), input.InstanceAddress)
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to find contract by InstanceAddress %s: %w", input.InstanceAddress.Hex(), err)
			}

			exerciseCommand := params.Method(contractID, input.Args)

			packageID, moduleName, entityName, err := contracts.ParseTemplateIDFromString(exerciseCommand.TemplateID)
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to parse template ID %s: %w", exerciseCommand.TemplateID, err)
			}

			choiceArgument := ledger.MapToValue(exerciseCommand.Arguments)

			disclosedContracts := make([]*apiv2.DisclosedContract, len(input.DisclosedContracts))
			for i, contract := range input.DisclosedContracts {
				disclosedContracts[i] = &apiv2.DisclosedContract{
					TemplateId:       contract.TemplateId.ToLedgerIdentifier(),
					ContractId:       contract.ContractId,
					CreatedEventBlob: contract.CreatedEventBlob,
					SynchronizerId:   contract.SynchronizerId,
				}
			}

			submitResp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(b.GetContext(), &apiv2.SubmitAndWaitForTransactionRequest{
				Commands: &apiv2.Commands{
					CommandId: uuid.Must(uuid.NewUUID()).String(),
					ActAs:     []string{participant.PartyID},
					Commands: []*apiv2.Command{{
						Command: &apiv2.Command_Exercise{
							Exercise: &apiv2.ExerciseCommand{
								TemplateId: &apiv2.Identifier{
									PackageId:  packageID,
									ModuleName: moduleName,
									EntityName: entityName,
								},
								ContractId:     contractID,
								Choice:         exerciseCommand.Choice,
								ChoiceArgument: choiceArgument,
							},
						},
					}},
					DisclosedContracts: disclosedContracts,
				},
			})
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to submit exercise command: %w", err)
			}

			return ExerciseOutput{
				ChainSelector: deps.ChainSelector(),
				ExecInfo: &ExecInfo{
					UpdateID: submitResp.GetTransaction().GetUpdateId(),
				},
			}, nil
		},
	)
}

// FindActiveContractByInstanceAddress finds an active contract by its instance address. It returns an error if there are multiple or zero active contracts matching the instance address.
// The returned ActiveContract includes the CreatedEventBlob required for explicit disclosures.
func FindActiveContractByInstanceAddress(ctx context.Context, stateService apiv2.StateServiceClient, party, templateId string, instanceAddress contracts.InstanceAddress) (*apiv2.ActiveContract, error) {
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger end: %w", err)
	}

	// Parse template ID to get package ID, module name, and entity name
	packageID, moduleName, entityName, err := contracts.ParseTemplateIDFromString(templateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template ID: %w", err)
	}

	activeContractsResp, err := stateService.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &apiv2.TemplateFilter{
									TemplateId: &apiv2.Identifier{
										PackageId:  packageID,
										ModuleName: moduleName,
										EntityName: entityName,
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts: %w", err)
	}
	defer activeContractsResp.CloseSend()

	var activeContract *apiv2.ActiveContract
	for {
		activeContractResp, err := activeContractsResp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}

		if c, ok := activeContractResp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			createArguments := c.ActiveContract.GetCreatedEvent().GetCreateArguments()
			if createArguments == nil {
				continue
			}

			var contractInstanceId string
			for _, field := range createArguments.GetFields() {
				if field.GetLabel() == "instanceId" {
					contractInstanceId = field.GetValue().GetText()
					break
				}
			}
			if contractInstanceId == "" {
				continue
			}

			// Get signatory of contract and compute instance address, then compare with the provided instance address
			instanceID := contracts.InstanceID(contractInstanceId)
			signatories := c.ActiveContract.GetCreatedEvent().GetSignatories()
			if len(signatories) != 1 {
				continue
			}
			gotAddress := instanceID.RawInstanceAddress(types.PARTY(signatories[0])).InstanceAddress()

			if instanceAddress != gotAddress {
				continue
			}

			if activeContract != nil {
				return nil, fmt.Errorf("multiple active contracts found for InstanceAddress %s", instanceAddress.String())
			}
			activeContract = c.ActiveContract
		}
	}

	if activeContract == nil {
		return nil, fmt.Errorf("no active contract found for InstanceAddress %s", instanceAddress.String())
	}

	return activeContract, nil
}

// FindActiveContractIDByInstanceAddress finds an active contract ID by its instance address. It returns an error if there are multiple or zero active contracts matching the instance address.
func FindActiveContractIDByInstanceAddress(ctx context.Context, stateService apiv2.StateServiceClient, party, templateId string, instanceAddress contracts.InstanceAddress) (string, error) {
	activeContract, err := FindActiveContractByInstanceAddress(ctx, stateService, party, templateId, instanceAddress)
	if err != nil {
		return "", err
	}

	return activeContract.GetCreatedEvent().GetContractId(), nil
}
