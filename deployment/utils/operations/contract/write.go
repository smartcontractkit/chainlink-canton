package contract

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonOps "github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
)

type ExecInfo struct {
	UpdateID string `json:"update_id"`
}

type WriteOutput struct {
	ChainSelector uint64    `json:"chainSelector"`
	ExecInfo      *ExecInfo `json:"exec_info"`
}

func (o WriteOutput) Executed() bool {
	return o.ExecInfo != nil
}

type ChoiceInput[ARGS any] struct {
	ChainSelector uint64               `json:"chainSelector"`
	InstanceID    contracts.InstanceID `json:"instanceId"`
	ActAs         []string             `json:"act_as"`
	Args          ARGS                 `json:"args"`
}

type WriteParams[ARGS any] struct {
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

	Template common.Template
	// Method is the bindings method to call the choice.
	Method func(contractID string, args ARGS) *model.ExerciseCommand
}

func NewWrite[ARGS any](params WriteParams[ARGS]) *operations.Operation[ChoiceInput[ARGS], WriteOutput, cantonOps.CantonDeps] {
	return operations.NewOperation(
		params.Name,
		params.Version,
		params.Description,
		func(b operations.Bundle, deps cantonOps.CantonDeps, input ChoiceInput[ARGS]) (WriteOutput, error) {
			if params.Validate != nil {
				if err := params.Validate(input.Args); err != nil {
					return WriteOutput{}, fmt.Errorf("validate input: %w", err)
				}
			}

			// Find contract by InstanceID
			contractID, err := findContractByInstanceID(b.GetContext(), deps.BindingClient, deps.Party, params.Template.GetTemplateID(), input.InstanceID)
			if err != nil {
				return WriteOutput{}, fmt.Errorf("failed to find contract by InstanceID: %w", err)
			}
			exerciseCommand := params.Method(contractID, input.Args)

			submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(b.GetContext(), &model.SubmitAndWaitRequest{
				Commands: &model.Commands{
					CommandID: uuid.Must(uuid.NewUUID()).String(),
					ActAs:     input.ActAs,
					Commands:  []*model.Command{{Command: exerciseCommand}},
				}},
			)
			if err != nil {
				return WriteOutput{}, fmt.Errorf("failed to submit exercise command: %w", err)
			}

			return WriteOutput{
				ChainSelector: input.ChainSelector,
				ExecInfo: &ExecInfo{
					UpdateID: submitResp.UpdateID,
				},
			}, nil
		},
	)
}

func findContractByInstanceID(ctx context.Context, bindingClient *client.DamlBindingClient, party, templateId string, instanceID contracts.InstanceID) (contractId string, err error) {
	currentOffset, err := bindingClient.StateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get ledger end: %w", err)
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
			Verbose: false,
		},
	})

	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				// Channel closed, stream ended
				if contractId == "" {
					return "", fmt.Errorf("no active contract found for InstanceID %s", instanceID.String())
				}

				return contractId, nil
			}
			if resp != nil && resp.ContractEntry != nil {
				if entry, ok := resp.ContractEntry.(*model.ActiveContractEntry); ok {
					if contractId != "" {
						// ContractID was already found. This is an error, InstanceID must be unique
						return "", fmt.Errorf("multiple active contracts found for InstanceID %s", instanceID.String())
					}
					contractId = entry.ActiveContract.CreatedEvent.ContractID
				}
			}
		case err := <-errChan:
			if err != nil {
				return "", fmt.Errorf("failed to get active contracts: %w", err)
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
