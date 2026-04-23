package testhelpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
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
			// Int64s are encoded as JSON numbers or strings, depending on the encoder settings
			switch val := rawValue.(type) {
			case string:
				valueInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: valueInt}}
			case json.Number:
				valueInt, err := strconv.ParseInt(val.String(), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: valueInt}}
			case float64:
				// Some encoders may encode int64s as JSON numbers, which are float64s in Go. This can cause precision loss for large int64s, but we can still parse them if they fit within uint64.
				if val < 0 || val > float64(^uint64(0)) {
					return nil, fmt.Errorf("AV_Int value is out of range for uint64: %f", val)
				}
				value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(val)}}
			default:
				return nil, fmt.Errorf("AV_Int value is not a string or number: %T", rawValue)
			}
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
			t, err := time.Parse(time.DateOnly, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a DateOnly time: %s", valueString)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Date{Date: int32(t.Unix() / 86400)}} //nolint:gosec // days since epoch
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Time value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: t.UnixMicro()}}
		case "AV_RelTime":
			// Int64s are encoded as JSON numbers or strings, depending on the encoder settings
			switch val := rawValue.(type) {
			case string:
				valueInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: valueInt}}},
				}}}}
			case json.Number:
				valueInt, err := strconv.ParseInt(val.String(), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: valueInt}}},
				}}}}
			case float64:
				// Some encoders may encode int64s as JSON numbers, which are float64s in Go. This can cause precision loss for large int64s, but we can still parse them if they fit within uint64.
				if val < 0 || val > float64(^uint64(0)) {
					return nil, fmt.Errorf("AV_Int value is out of range for uint64: %f", val)
				}
				value = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(val)}}},
				}}}}
			default:
				return nil, fmt.Errorf("AV_Int value is not a string or number: %T", rawValue)
			}
		case "AV_Party":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Party is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Party{Party: valueString}}
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: valueString}}
		default:
			// TODO Add lists and maps
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

// ExtractChoiceContextValues converts a Splice ChoiceContext proto value into
// the typed map expected by TokenInput.ExtraArgs.Context.Values.
func ExtractChoiceContextValues(choiceContext *apiv2.Value) map[string]splice_api_token_metadata_v1.AnyValue {
	contextValues := make(map[string]splice_api_token_metadata_v1.AnyValue)
	if rec := choiceContext.GetRecord(); rec != nil && len(rec.Fields) > 0 {
		valuesField := rec.Fields[0]
		if valuesField.GetLabel() == "values" && valuesField.GetValue().GetTextMap() != nil {
			for _, entry := range valuesField.GetValue().GetTextMap().GetEntries() {
				if v := entry.GetValue().GetVariant(); v != nil {
					contextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVContractId: new(types.CONTRACT_ID(v.GetValue().GetContractId()))}
				} else if entry.GetValue().GetText() != "" {
					contextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVText: new(types.TEXT(entry.GetValue().GetText()))}
				}
			}
		}
	}

	return contextValues
}

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}

// GetRegistryAdmin reads the already-provisioned token registry admin (AdminId)
// from scan-proxy registry info. It does not register or mutate admin state.
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

func GetTransferFactory(ctx context.Context, transferInstructionClient transferInstructionV1.ClientWithResponsesInterface, registryAdmin, sender, receiver string) (string, []*apiv2.DisclosedContract, map[string]any, error) {
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

	return transferFactoryResponse.JSON200.FactoryId, disclosedContracts, transferFactoryResponse.JSON200.ChoiceContext.ChoiceContextData, nil
}

func AcceptPendingTransferInstruction(
	ctx context.Context,
	participant canton.Participant,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	executingParty string,
	pendingTransferInstructionCID string,
) error {
	acceptContextResp, err := transferInstructionClient.GetTransferInstructionAcceptContextWithResponse(ctx, pendingTransferInstructionCID, transferInstructionV1.GetChoiceContextRequest{})
	if err != nil {
		return fmt.Errorf("get transfer instruction accept context: %w", err)
	}
	if acceptContextResp.StatusCode() != http.StatusOK || acceptContextResp.JSON200 == nil {
		return fmt.Errorf("unexpected transfer instruction accept context status=%d", acceptContextResp.StatusCode())
	}

	acceptDisclosures := make([]*apiv2.DisclosedContract, 0, len(acceptContextResp.JSON200.DisclosedContracts))
	for _, contract := range acceptContextResp.JSON200.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		if err != nil {
			return fmt.Errorf("parse accept-context template id: %w", err)
		}
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if err != nil {
			return fmt.Errorf("decode accept-context created event blob: %w", err)
		}
		acceptDisclosures = append(acceptDisclosures, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}

	acceptContext, err := ChoiceContextFromData(acceptContextResp.JSON200.ChoiceContextData)
	if err != nil {
		return fmt.Errorf("convert transfer instruction accept context: %w", err)
	}

	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
					ContractId: pendingTransferInstructionCID,
					Choice:     "TransferInstruction_Accept",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "context", Value: acceptContext},
							{Label: "meta", Value: emptyMetadata},
						}}}}},
					}}}},
				}},
			}},
			ActAs:              []string{executingParty},
			DisclosedContracts: acceptDisclosures,
		},
	})
	if err != nil {
		return fmt.Errorf("accept pending transfer instruction: %w", err)
	}

	return nil
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
	_, amuletRulesContract, err := GetAmuletRulesContract(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get amulet rules contract: %w", err)
	}

	// Create Transfer Factory
	_, disclosedContracts, _, err := GetTransferFactory(ctx, transferInstructionClient, registryAdmin, registryAdmin, toParty)
	if err != nil {
		return "", fmt.Errorf("failed to get transfer factory: %w", err)
	}

	// Get open mining round
	openMiningRoundContract, err := GetFirstOpenMiningRound(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get open mining round: %w", err)
	}

	// Mint AMT
	response, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: amuletRulesContract.TemplateId,
							ContractId: amuletRulesContract.ContractId,
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
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: openMiningRoundContract.ContractId}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{toParty},
			DisclosedContracts: DeduplicateDisclosedContracts(append(disclosedContracts, amuletRulesContract, openMiningRoundContract)...),
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

// CreateTransferPreapproval creates an AmuletRules TransferPreapproval for the specified party as a receiver.
// It returns the ContractId of the Preapproval contract at the time of creation. Due to Amulets implementation
// of preapprovals, the contract ID is not long-lived and should not be used directly. When interacting/using
// the preapproval, always look it up using the ACS.
func CreateTransferPreapproval(
	ctx context.Context,
	participant canton.Participant,
	scanProxyClient scanProxy.ClientWithResponsesInterface,
	party string,
	holdingCid string,
) (types.CONTRACT_ID, error) {
	dsoPartyId, amuletRulesContract, err := GetAmuletRulesContract(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get DSO info: %w", err)
	}
	// Get open mining round
	openMiningRoundContract, err := GetFirstOpenMiningRound(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get open mining round: %w", err)
	}

	response, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: amuletRulesContract.TemplateId,
							ContractId: amuletRulesContract.ContractId,
							Choice:     "AmuletRules_CreateTransferPreapproval",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "context",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "amuletRules",
											Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: amuletRulesContract.ContractId}},
										}, {
											Label: "context",
											Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
												{Label: "openMiningRound", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: openMiningRoundContract.ContractId}}},
												{Label: "issuingMiningRounds", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
												{Label: "validatorRights", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
												{Label: "featuredAppRight", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}},
											}}}},
										},
									}}}},
								}, {
									Label: "inputs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
											Constructor: "InputAmulet",
											Value:       &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: holdingCid}},
										}}},
									}}}},
								}, {
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: party}},
								}, {
									Label: "provider",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: party}},
								}, {
									Label: "expiresAt",
									Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().Add(time.Hour * 24).UnixMicro()}},
								}, {
									Label: "expectedDso",
									Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: dsoPartyId}}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{party},
			DisclosedContracts: []*apiv2.DisclosedContract{amuletRulesContract, openMiningRoundContract},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create TransferPreapproval: %w", err)
	}

	preapprovalCid := ""
	for _, event := range response.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.TemplateId.EntityName == "TransferPreapproval" {
				preapprovalCid = e.Created.ContractId
			}
		}
	}
	if preapprovalCid == "" {
		return "", fmt.Errorf("failed to find preapproval")
	}

	return types.CONTRACT_ID(preapprovalCid), nil
}

func DeduplicateDisclosedContracts(disclosedContracts ...*apiv2.DisclosedContract) []*apiv2.DisclosedContract {
	m := make(map[string]*apiv2.DisclosedContract)
	for _, contract := range disclosedContracts {
		m[contract.ContractId] = contract
	}

	return maps.Values(m)
}
