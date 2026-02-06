package contract

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/ledger"
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
			contractID, err := FindContractIDByInstanceAddress(b.GetContext(), b.Logger, deps.StateServiceClient, deps.Party, params.Template.GetTemplateID(), input.InstanceAddress)
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to find contract by InstanceAddress %s: %w", input.InstanceAddress.Hex(), err)
			}

			// Get template ID and choice name from the method
			exerciseCommand := params.Method(contractID, input.Args)

			// Parse template ID to get package ID, module name, and entity name
			packageID, moduleName, entityName, err := parseTemplateIDFromString(exerciseCommand.TemplateID)
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to parse template ID %s: %w", exerciseCommand.TemplateID, err)
			}

			// Convert args struct to ledger.MapToValue for ChoiceArgument
			choiceArgument := ledger.MapToValue(input.Args)

			submitResp, err := deps.CommandServiceClient.SubmitAndWaitForTransaction(b.GetContext(), &apiv2.SubmitAndWaitForTransactionRequest{
				Commands: &apiv2.Commands{
					CommandId: uuid.Must(uuid.NewUUID()).String(),
					ActAs:     input.ActAs,
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
				},
			})
			if err != nil {
				return ExerciseOutput{}, fmt.Errorf("failed to submit exercise command: %w", err)
			}

			// Note: apiv2.SubmitAndWaitForTransactionResponse doesn't expose UpdateId directly
			// The transaction was successfully submitted, which is what matters
			return ExerciseOutput{
				ChainSelector: input.ChainSelector,
				ExecInfo: &ExecInfo{
					UpdateID: submitResp.GetTransaction().GetUpdateId(),
				},
			}, nil
		},
	)
}

// FindContractIDByInstanceAddress finds a contract ID by its instance address
func FindContractIDByInstanceAddress(ctx context.Context, logger logger.Logger, stateService apiv2.StateServiceClient, party, templateId string, instanceAddress contracts.InstanceAddress) (string, error) {
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get ledger end: %w", err)
	}

	// Parse template ID to get package ID, module name, and entity name
	packageID, moduleName, entityName, err := parseTemplateIDFromString(templateId)
	if err != nil {
		return "", fmt.Errorf("failed to parse template ID: %w", err)
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
									IncludeCreatedEventBlob: false,
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
		return "", fmt.Errorf("failed to get active contracts: %w", err)
	}
	defer activeContractsResp.CloseSend()

	var contractID string
	for {
		activeContract, err := activeContractsResp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", fmt.Errorf("failed to receive active contracts: %w", err)
		}

		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			createArguments := c.ActiveContract.GetCreatedEvent().GetCreateArguments()
			if createArguments == nil {
				logger.Debugw("Skipping contract with nil create arguments", "contractID", c.ActiveContract.GetCreatedEvent().ContractId)
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
				logger.Debugw("Skipping contract with missing instanceId field", "contractID", c.ActiveContract.GetCreatedEvent().ContractId)
				continue
			}

			instanceID := contracts.InstanceID(contractInstanceId)
			if !instanceID.Valid() {
				logger.Debugw("Skipping contract with invalid instanceId field", "contractID", c.ActiveContract.GetCreatedEvent().ContractId, "instanceId", contractInstanceId)
				continue
			}

			if instanceAddress != instanceID.InstanceAddress() {
				logger.Debugw("Skipping contract with different instanceId field", "contractID", c.ActiveContract.GetCreatedEvent().ContractId, "instanceId", contractInstanceId)
				continue
			}

			if contractID != "" {
				return "", fmt.Errorf("multiple active contracts found for InstanceID %s", instanceID.String())
			}
			contractID = c.ActiveContract.GetCreatedEvent().ContractId
		}
	}

	if contractID == "" {
		return "", fmt.Errorf("no active contract found for InstanceAddress %s", instanceAddress.String())
	}

	return contractID, nil
}

// parseTemplateIDFromString parses a template ID string like "#package:Module:Entity" into its components
func parseTemplateIDFromString(templateID string) (packageID, moduleName, entityName string, err error) {
	if !strings.HasPrefix(templateID, "#") {
		return "", "", "", fmt.Errorf("template ID must start with #")
	}
	parts := strings.Split(templateID, ":")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("template ID must have format #package:module:entity, got: %s", templateID)
	}

	return parts[0], parts[1], parts[2], nil
}
