package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

// TODO: move this helper into ccv.Lib.
func GetChain(t *testing.T, chainType string, cfg *ccv.Cfg, harness tcapi.TestHarness) cciptestinterfaces.CCIP17 {
	var chain *blockchain.Input
	for _, bc := range cfg.Blockchains {
		if bc.Type == chainType {
			chain = bc
			break
		}
	}
	require.NotNil(t, chain, "need at least one chain for this test")

	var family string
	switch chainType {
	case blockchain.TypeCanton:
		family = chainsel.FamilyCanton
	case blockchain.TypeAnvil:
		family = chainsel.FamilyEVM
	default:
		t.Fatalf("unsupported chain type %q", chainType)
	}

	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(chain.ChainID, family)
	require.NoError(t, err)

	chainMap, err := harness.Lib.ChainsMap(t.Context())
	require.NoError(t, err)

	return chainMap[chainDetails.ChainSelector]
}

func GetFeeTransferFactoryInput(
	ctx context.Context,
	participant canton.Participant,
	sender string,
	receiver string,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ExtraArgs, []*ledgerv2.DisclosedContract, error) {
	auth := func(ctx context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("get participant token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}

	scanProxyBaseURL := fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL)
	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(scanProxyBaseURL, tokenMetadataV1.WithRequestEditorFn(auth))
	if err != nil {
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("create token metadata client: %w", err)
	}
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(scanProxyBaseURL, transferInstructionV1.WithRequestEditorFn(auth))
	if err != nil {
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("create transfer instruction client: %w", err)
	}

	registryInfoResponse, err := tokenMetadataClient.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("get registry info: %w", err)
	}
	if registryInfoResponse.StatusCode() != http.StatusOK || registryInfoResponse.JSON200 == nil {
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("get registry info unexpected status: %d", registryInfoResponse.StatusCode())
	}
	registryAdmin := registryInfoResponse.JSON200.AdminId

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
				"requestedAt":      time.Now().Add(-time.Hour).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
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
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("get transfer factory: %w", err)
	}
	if transferFactoryResponse.StatusCode() != http.StatusOK || transferFactoryResponse.JSON200 == nil {
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("get transfer factory unexpected status: %d", transferFactoryResponse.StatusCode())
	}

	disclosedContracts := make([]*ledgerv2.DisclosedContract, 0, len(transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts))
	for _, dc := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
		id, err := templateIdFromString(dc.TemplateId)
		if err != nil {
			return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("parse transfer factory template id: %w", err)
		}
		createdEventBlob, err := base64.StdEncoding.DecodeString(dc.CreatedEventBlob)
		if err != nil {
			return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("decode transfer factory created event blob: %w", err)
		}
		disclosedContracts = append(disclosedContracts, &ledgerv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       dc.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   dc.SynchronizerId,
		})
	}

	choiceContext, err := ChoiceContextFromData(transferFactoryResponse.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", splice_api_token_metadata_v1.ExtraArgs{}, nil, fmt.Errorf("convert transfer factory choice context: %w", err)
	}
	transferFactoryContextValues := make(types.TEXTMAP)
	if choiceContextRecord := choiceContext.GetRecord(); choiceContextRecord != nil && len(choiceContextRecord.Fields) > 0 {
		valuesField := choiceContextRecord.Fields[0]
		if valuesField.GetLabel() == "values" && valuesField.GetValue().GetTextMap() != nil {
			for _, entry := range valuesField.GetValue().GetTextMap().GetEntries() {
				if v := entry.GetValue().GetVariant(); v != nil {
					cid := types.CONTRACT_ID(v.GetValue().GetContractId())
					transferFactoryContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVContractId: &cid}
				} else if entry.GetValue().GetText() != "" {
					txt := types.TEXT(entry.GetValue().GetText())
					transferFactoryContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVText: &txt}
				}
			}
		}
	}

	return types.CONTRACT_ID(transferFactoryResponse.JSON200.FactoryId), splice_api_token_metadata_v1.ExtraArgs{
		Context: splice_api_token_metadata_v1.ChoiceContext{Values: transferFactoryContextValues},
		Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
	}, disclosedContracts, nil
}

func templateIdFromString(s string) (*ledgerv2.Identifier, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid template id format: %s", s)
	}

	return &ledgerv2.Identifier{
		PackageId:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

func ChoiceContextFromData(choiceContextData map[string]any) (*ledgerv2.Value, error) {
	values, ok := choiceContextData["values"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no values found in choice context")
	}

	fields := make([]*ledgerv2.TextMap_Entry, 0, len(values))
	for k, v := range values {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value *ledgerv2.Value
		switch tag {
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: valueString}}
		case "AV_Int":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: int64(valueFloat)}}
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: valueString}}
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return nil, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: valueBool}}
		case "AV_Date":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Date{Date: int32(t.Unix() / 86400)}} //nolint:gosec
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Time value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Time value is not a RFC3339 time: %s", valueString)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Timestamp{Timestamp: t.UnixMicro()}}
		case "AV_RelTime":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_RelTime value is not a number: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
				{Label: "microseconds", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: int64(valueFloat)}}},
			}}}}
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_ContractId{ContractId: valueString}}
		default:
			return nil, fmt.Errorf("unimplemented tag: %v", tag)
		}

		fields = append(fields, &ledgerv2.TextMap_Entry{
			Key: k,
			Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Variant{Variant: &ledgerv2.Variant{
				Constructor: tag,
				Value:       value,
			}}},
		})
	}

	return &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
		{
			Label: "values",
			Value: &ledgerv2.Value{Sum: &ledgerv2.Value_TextMap{TextMap: &ledgerv2.TextMap{Entries: fields}}},
		},
	}}}}, nil
}
