package contract

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonOps "github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
)

type ExecInfo struct {
	// The ID of the update that was confirmed as the result of this execution.
	UpdateID string `json:"update_id"`
}

type ExerciseOutput struct {
	ChainSelector uint64    `json:"chainSelector"`
	ExecInfo      *ExecInfo `json:"exec_info"`
}

func (o ExerciseOutput) Executed() bool {
	return o.ExecInfo != nil
}

type ChoiceInput[ARGS any] struct {
	ChainSelector uint64 `json:"chainSelector"`
	// The InstanceAddress this operation is targeting. Will be resolved to an active contract.
	InstanceAddress contracts.InstanceAddress `json:"instanceAddress"`
	ActAs           []string                  `json:"act_as"`
	Args            ARGS                      `json:"args"`
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
	// Validate is an optional function to validate the input arguments.
	Validate func(input ARGS) error

	// Template is the binding struct of the target contract.
	Template common.Template
	// Method is the bindings method to call the choice.
	Method func(contractID string, args ARGS) *model.ExerciseCommand
}

func NewExercise[ARGS any](params ExerciseParams[ARGS]) *operations.Operation[ChoiceInput[ARGS], ExerciseOutput, cantonOps.CantonDeps] {
	return operations.NewOperation(
		params.Name,
		params.Version,
		params.Description,
		func(b operations.Bundle, deps cantonOps.CantonDeps, input ChoiceInput[ARGS]) (ExerciseOutput, error) {
			if params.Validate != nil {
				if err := params.Validate(input.Args); err != nil {
					return ExerciseOutput{}, fmt.Errorf("validate input: %w", err)
				}
			}

			// Find contract by InstanceAddress
			createEvent, err := FindCreateEventByInstanceAddress(b.GetContext(), b.Logger, deps.BindingClient, deps.Party, params.Template.GetTemplateID(), input.InstanceAddress)
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to find contract by InstanceAddress %s: %w", input.InstanceAddress.Hex(), err)
			}
			exerciseCommand := params.Method(createEvent.ContractID, input.Args)

			submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(b.GetContext(), &model.SubmitAndWaitRequest{
				Commands: &model.Commands{
					CommandID: uuid.Must(uuid.NewUUID()).String(),
					ActAs:     input.ActAs,
					Commands:  []*model.Command{{Command: exerciseCommand}},
				}},
			)
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to submit exercise command: %w", err)
			}

			return ExerciseOutput{
				ChainSelector: input.ChainSelector,
				ExecInfo: &ExecInfo{
					UpdateID: submitResp.UpdateID,
				},
			}, nil
		},
	)
}

func FindCreateEventByInstanceAddress(ctx context.Context, logger logger.Logger, bindingClient *client.DamlBindingClient, party, templateId string, instanceAddress contracts.InstanceAddress) (*model.CreatedEvent, error) {
	currentOffset, err := bindingClient.StateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger end: %w", err)
	}
	responseChan, errChan := bindingClient.StateService.GetActiveContracts(ctx, &model.GetActiveContractsRequest{
		ActiveAtOffset: currentOffset.Offset,
		EventFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				party: {
					Inclusive: &model.InclusiveFilters{
						TemplateFilters: []*model.TemplateFilter{
							{
								TemplateID:              templateId,
								IncludeCreatedEventBlob: false,
							},
						},
					},
				},
			},
			// Verbose is needed for the record labels to be returned
			Verbose: true,
		},
	})

	var createEvent *model.CreatedEvent

	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				// Channel closed, stream ended
				if createEvent == nil {
					return nil, fmt.Errorf("no active contract found for InstanceAddress %s", instanceAddress.String())
				}

				return createEvent, nil
			}
			if resp != nil && resp.ContractEntry != nil {
				if entry, ok := resp.ContractEntry.(*model.ActiveContractEntry); ok {
					// Compare the instanceId field

					createArguments, ok := entry.ActiveContract.CreatedEvent.CreateArguments.(*v2.Record)
					if !ok {
						logger.Debugw("Skipping contract with unexpected create arguments type", "contractID", entry.ActiveContract.CreatedEvent.ContractID)
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
						logger.Debugw("Skipping contract with missing instanceId field", "contractID", entry.ActiveContract.CreatedEvent.ContractID)
						continue
					}
					instanceID := contracts.InstanceID(contractInstanceId)
					if !instanceID.Valid() {
						logger.Debugw("Skipping contract with invalid instanceId field", "contractID", entry.ActiveContract.CreatedEvent.ContractID, "instanceId", contractInstanceId)
						continue
					}
					if instanceAddress != instanceID.InstanceAddress() {
						// Not the contract we're looking for
						logger.Debugw("Skipping contract with different instanceId field", "contractID", entry.ActiveContract.CreatedEvent.ContractID, "instanceId", contractInstanceId)
						continue
					}

					if createEvent != nil {
						// contract was already found. This is an error, InstanceID must be unique
						return nil, fmt.Errorf("multiple active contracts found for InstanceID %s", instanceID.String())
					}
					createEvent = entry.ActiveContract.CreatedEvent
				}
			}
		case err := <-errChan:
			if err != nil {
				return nil, fmt.Errorf("failed to get active contracts: %w", err)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
