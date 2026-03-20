package devenv

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

// TODO: dedupe with integration-tests GetTransferFactory helper once shared helper location is established.
func getTokenRegistryAdmin(ctx context.Context, participant canton.Participant) (types.PARTY, error) {
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
		return "", fmt.Errorf("create token metadata client: %w", err)
	}

	registryInfoResponse, err := tokenMetadataClient.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("get registry info: %w", err)
	}
	if registryInfoResponse.StatusCode() != http.StatusOK || registryInfoResponse.JSON200 == nil {
		return "", fmt.Errorf("get registry info unexpected status: %d", registryInfoResponse.StatusCode())
	}

	return types.PARTY(registryInfoResponse.JSON200.AdminId), nil
}

// TODO: dedupe with integration-tests GetTransferFactory helper once shared helper location is established.
func getFeeTransferFactoryInput(
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
		id, err := TemplateIdFromString(dc.TemplateId)
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
