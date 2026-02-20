package testhelpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

func ChoiceContextFromData(choiceContextData map[string]any) (*apiv2.Value, error) {
	values, ok := choiceContextData["values"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no values found in choice context")
	}

	// ref: https://docs.digitalasset.com/build/3.5/reference/json-api/lf-value-specification.html
	// AnyValue is a variant
	var fields []*apiv2.TextMap_Entry
	for k, v := range values {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value *apiv2.Value
		switch tag {
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Text{Text: valueString}}
		case "AV_Int":
			// JSON numbers come as float64
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(valueFloat)}}
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: valueString}}
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return nil, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: valueBool}}
		case "AV_Date":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Date{Date: int32(t.Unix() / 86400)}} // days since epoch
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: t.UnixMicro()}}
		case "AV_RelTime":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(valueFloat)}}},
			}}}}
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: valueString}}
		default:
			// Add lists and maps
			return nil, fmt.Errorf("unimplemented tag: %v", tag)
		}

		fields = append(fields, &apiv2.TextMap_Entry{
			Key: k,
			Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
				Constructor: tag,
				Value:       value,
			}}},
		})
	}

	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "values",
			Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: fields}}},
		},
	}}}}, nil
}

func GetRegistryAdmin(ctx context.Context, metadataClient tokenMetadataV1.ClientWithResponsesInterface) (string, error) {
	registryInfoResponse, err := metadataClient.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting registry info: %w", err)
	}
	if registryInfoResponse.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d: %v", registryInfoResponse.StatusCode(), registryInfoResponse.Body)
	}

	return registryInfoResponse.JSON200.AdminId, nil
}

func GetTransferFactory(ctx context.Context, transferInstructionClient transferInstructionV1.ClientWithResponsesInterface, registryAdmin, sender, receiver string) (string, []*apiv2.DisclosedContract, *apiv2.Value, error) {
	transferFactoryResponse, err := transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": registryAdmin,
			"transfer": map[string]any{
				"sender":   sender,
				"receiver": receiver,
				"amount":   "100.00",
				"instrumentId": map[string]any{
					"admin": registryAdmin,
					"id":    "Amulet",
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(time.Hour * -1).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				"inputHoldingCids": []string{},
				"meta": map[string]any{
					"values": map[string]any{},
				},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{
					"values": map[string]any{},
				},
				"meta": map[string]any{
					"values": map[string]any{},
				},
			},
		},
		ExcludeDebugFields: nil,
	})
	if err != nil {
		return "", nil, nil, fmt.Errorf("error getting transferFactory response: %w", err)
	}
	if transferFactoryResponse.StatusCode() != http.StatusOK {
		return "", nil, nil, fmt.Errorf("unexpected status code: %d: %v", transferFactoryResponse.StatusCode(), transferFactoryResponse.Body)
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		if err != nil {
			return "", nil, nil, fmt.Errorf("failed to parse template id: %w", err)
		}
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if err != nil {
			return "", nil, nil, fmt.Errorf("failed to decode created event blob: %w", err)
		}
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}
	choiceContext, err := ChoiceContextFromData(transferFactoryResponse.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	factoryCid := transferFactoryResponse.JSON200.FactoryId

	return factoryCid, disclosedContracts, choiceContext, nil
}

func MintAMT(
	ctx context.Context,
	participant canton.Participant,
	metadataClient tokenMetadataV1.ClientWithResponsesInterface,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	scanProxyClient scanProxy.ClientWithResponsesInterface,
	toParty string,
	amount string,
) (string, error) {
	// Get Instrument Admin
	registryAdmin, err := GetRegistryAdmin(ctx, metadataClient)
	if err != nil {
		return "", fmt.Errorf("failed to get registry admin: %w", err)
	}

	// Get AmuletRules Contract
	amuletRulesCid, amuletRulesId, err := GetAmuletRulesContract(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get amulet rules contract: %w", err)
	}

	// Create Transfer Factory
	_, disclosedContracts, _, err := GetTransferFactory(ctx, transferInstructionClient, registryAdmin, registryAdmin, toParty)
	if err != nil {
		return "", fmt.Errorf("failed to get transfer factory: %w", err)
	}

	// Get open mining round
	openMiningRoundCid, err := GetFirstOpenMiningRound(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get open mining round: %w", err)
	}

	// Mint AMT
	response, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: amuletRulesId,
							ContractId: amuletRulesCid,
							Choice:     "AmuletRules_DevNet_Tap",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: toParty}},
								}, {
									Label: "amount",
									Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: amount}},
								}, {
									Label: "openRound",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: openMiningRoundCid}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{toParty},
			DisclosedContracts: disclosedContracts,
		},
		TransactionFormat: nil,
	})
	if err != nil {
		return "", fmt.Errorf("failed to mint AMT: %w", err)
	}

	tokenHoldingCid := ""
	for _, event := range response.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenHoldingCid = e.Created.ContractId
		}
	}

	return tokenHoldingCid, nil
}
