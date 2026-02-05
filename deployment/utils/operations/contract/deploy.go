package contract

import (
	"fmt"
	"reflect"

	"github.com/aws/smithy-go/ptr"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
)

type DeployInput[TT common.Template] struct {
	ChainSelector uint64      `json:"chainSelector"`
	Qualifier     *string     `json:"qualifier,omitempty"`
	ActAs         []string    `json:"act_as"`
	Template      TT          `json:"createCommand"`
	OwnerParty    types.PARTY `json:"ownerParty"`
}

type DeployParams[ARGS any] struct {
	// Name is the name of the operation.
	Name string
	// TypeAndVersion is the type and version of the contract that's being deployed.
	TypeAndVersion deployment.TypeAndVersion
	// Description is a brief description of the operation.
	Description string
	// Validate is an optional function to validate the constructor arguments before deployment.
	Validate func(template ARGS) error
	// TODO remove and make part of template itself
	PackageName string
	// Prefix is the prefix used to create this contract's InstanceID.
	Prefix string `json:"prefix"`
}

// NewDeploy creates a new deploy operation for the given contract template.
//
// InstanceId:
// The template *must* have an InstanceId field of type types.TEXT.
// The instance ID is generated using the provided prefix and the owner party, it should not be set by the caller.
func NewDeploy[TT common.Template](params DeployParams[TT]) *operations.Operation[DeployInput[TT], datastore.AddressRef, dependencies.CantonDeps] {
	return operations.NewOperation(
		params.Name,
		&params.TypeAndVersion.Version,
		params.Description,
		func(b operations.Bundle, deps dependencies.CantonDeps, input DeployInput[TT]) (datastore.AddressRef, error) {
			// Generate InstanceID
			instanceID, err := contracts.NewInstanceID(params.Prefix, string(input.OwnerParty))
			if err != nil {
				return datastore.AddressRef{}, fmt.Errorf("failed to create instance ID: %w", err)
			}
			// Set InstanceID in the template
			templWithID, err := setInstanceID(input.Template, instanceID)
			if err != nil {
				return datastore.AddressRef{}, fmt.Errorf("failed to set InstanceID in template: %w", err)
			}

			// Validate
			if params.Validate != nil {
				if err := params.Validate(input.Template); err != nil {
					return datastore.AddressRef{}, fmt.Errorf("validate input: %w", err)
				}
			}
			if input.ChainSelector != deps.Chain.Selector {
				return datastore.AddressRef{}, fmt.Errorf("input deps selector %d does not match operation chain selector %d", input.ChainSelector, deps.Chain.Selector)
			}

			createCommand := templWithID.(TT).CreateCommand()

			// Convert model.CreateCommand to apiv2.CreateCommand
			apiv2CreateCmd, err := convertCreateCommandToAPIV2(createCommand)
			if err != nil {
				return datastore.AddressRef{}, fmt.Errorf("failed to convert create command: %w", err)
			}

			submitResp, err := deps.CommandServiceClient.SubmitAndWaitForTransaction(b.GetContext(), &apiv2.SubmitAndWaitForTransactionRequest{
				Commands: &apiv2.Commands{
					CommandId: uuid.Must(uuid.NewUUID()).String(),
					ActAs:     input.ActAs,
					Commands: []*apiv2.Command{{
						Command: &apiv2.Command_Create{Create: apiv2CreateCmd},
					}},
				},
			})
			if err != nil {
				return datastore.AddressRef{}, fmt.Errorf("failed to submit create command: %w", err)
			}
			contractId, err := getDeployedContractIDFromEvents(submitResp.GetTransaction().GetEvents(), input.Template, params.PackageName)
			if err != nil {
				return datastore.AddressRef{}, fmt.Errorf("failed to get deployed contract ID: %w", err)
			}
			b.Logger.Debugw(fmt.Sprintf("Deployed %s to %s", params.TypeAndVersion, deps.Chain), "contractID", contractId, "instanceID", instanceID.String(), "instanceAddress", instanceID.InstanceAddress().Hex())

			return datastore.AddressRef{
				Address:       instanceID.InstanceAddress().Hex(),
				ChainSelector: input.ChainSelector,
				Type:          datastore.ContractType(params.TypeAndVersion.Type),
				Version:       &params.TypeAndVersion.Version,
				Qualifier:     ptr.ToString(input.Qualifier),
			}, nil
		},
	)
}

// setInstanceID sets the InstanceId field of the given template to the provided instanceID.
func setInstanceID(template common.Template, instanceID contracts.InstanceID) (common.Template, error) {
	t := reflect.TypeOf(template)
	v := reflect.ValueOf(template)

	deref := false
	if t.Kind() != reflect.Pointer {
		t = reflect.PointerTo(t)
		v = reflect.New(t.Elem())
		v.Elem().Set(reflect.ValueOf(template))
		template = v.Interface().(common.Template)
		deref = true
	}

	if _, ok := t.Elem().FieldByName("InstanceId"); !ok {
		return nil, fmt.Errorf("no InstanceId field found in template")
	}

	field := v.Elem().FieldByName("InstanceId")
	instanceIDText := types.TEXT(instanceID.String())

	if !field.Type().AssignableTo(reflect.TypeFor[types.TEXT]()) {
		return nil, fmt.Errorf("cannot assign InstanceId field of type %v", field.Type())
	}

	if !field.CanSet() {
		return nil, fmt.Errorf("cannot assign InstanceId field of type %v", field.Type())
	}

	field.Set(reflect.ValueOf(instanceIDText))

	if deref {
		return v.Elem().Interface().(common.Template), nil
	}

	return template, nil
}

// convertCreateCommandToAPIV2 converts a model.CreateCommand to apiv2.CreateCommand
func convertCreateCommandToAPIV2(cmd *model.CreateCommand) (*apiv2.CreateCommand, error) {
	packageID, moduleName, entityName, err := parseTemplateID(cmd.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template ID %s: %w", cmd.TemplateID, err)
	}

	// Convert arguments map to apiv2.Value (Record)
	createArgumentsValue, err := convertMapToAPIV2Value(cmd.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to convert arguments: %w", err)
	}

	// Extract Record from Value
	recordValue, ok := createArgumentsValue.GetSum().(*apiv2.Value_Record)
	if !ok {
		return nil, fmt.Errorf("failed to extract Record from Value")
	}

	return &apiv2.CreateCommand{
		TemplateId: &apiv2.Identifier{
			PackageId:  packageID,
			ModuleName: moduleName,
			EntityName: entityName,
		},
		CreateArguments: recordValue.Record,
	}, nil
}

// TODO: packageName add package name to bindings instead
func getDeployedContractIDFromEvents(events []*apiv2.Event, template common.Template, packageName string) (string, error) {
	for _, event := range events {
		created := event.GetCreated()
		if created == nil {
			continue
		}

		templateID := created.GetTemplateId()
		if templateID == nil {
			continue
		}

		// Compare template IDs
		expectedTemplateID := contracts.ReplacePackageIdWithNameInTemplateID(template.GetTemplateID(), packageName)
		eventTemplateID := fmt.Sprintf("#%s:%s:%s", templateID.GetPackageId(), templateID.GetModuleName(), templateID.GetEntityName())
		eventTemplateIDWithPackageName := contracts.ReplacePackageIdWithNameInTemplateID(eventTemplateID, packageName)

		if expectedTemplateID == eventTemplateIDWithPackageName {
			return created.GetContractId(), nil
		}
	}

	return "", fmt.Errorf("failed to find contract in transaction events for template ID %s", template.GetTemplateID())
}
