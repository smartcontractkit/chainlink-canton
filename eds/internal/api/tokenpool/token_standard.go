package tokenpool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

type transferFactory func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error)

func getTransferFactory(ctx context.Context, poolOwner types.PARTY, acs store.ActiveContractStoreInterface, cfg config.TransferFactory) (transferFactory, error) {
	switch cfg.Type {
	case config.FactoryTypeDisabled:
		return nil, nil //nolint:nilnil
	case config.FactoryTypeAddress:
		factoryAddress := *cfg.InstanceAddress

		templateId, err := contracts.TemplateIDFromString(*cfg.TemplateId)
		if err != nil {
			return nil, fmt.Errorf("invalid TemplateId for TransferPreapproval: %w", err)
		}
		acs.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: templateId,
			PartyID:    *cfg.Party,
		})

		return func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
			activeTransferFactory, ok := acs.Get(factoryAddress)
			if !ok {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", factoryAddress)
			}

			return activeTransferFactory.GetCreatedEvent().GetContractId(), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
		}, nil
	case config.FactoryTypeURL:
		// If authentication has been configured, add an interceptor that adds the Authorization header
		var options []transferInstructionV1.ClientOption
		if cfg.TokenStandardAuthConfig != nil {
			authProvider, err := cfg.TokenStandardAuthConfig.NewProvider(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create auth provider: %w", err)
			}
			interceptor := func(ctx context.Context, req *http.Request) error {
				token, err := authProvider.TokenSource().Token()
				if err != nil {
					return fmt.Errorf("failed to retrieve token: %w", err)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

				return nil
			}
			options = append(options, transferInstructionV1.WithRequestEditorFn(interceptor))
		}
		transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(
			*cfg.TokenStandardURL,
			options...,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create TransferInstructionV1 client with URL %q: %w", *cfg.TokenStandardURL, err)
		}

		return func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
			resp, err := transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
				ChoiceArguments: map[string]any{
					"expectedAdmin": instrumentId.Admin,
					"transfer": map[string]any{
						"sender":   poolOwner,
						"receiver": poolOwner,
						"amount":   "1.0",
						"instrumentId": map[string]any{
							"admin": instrumentId.Admin,
							"id":    instrumentId.Id,
						},
						"lock":             nil,
						"requestedAt":      time.Now().Add(time.Second * -10).Format(time.RFC3339),
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
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to call GetTransferFactory: %w", err)
			}
			if resp.StatusCode() != http.StatusOK {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
			}

			disclosedContracts := make([]oapiCommon.DisclosedContract, len(resp.JSON200.ChoiceContext.DisclosedContracts))
			for i, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
				disclosedContracts[i] = oapiCommon.DisclosedContract{
					TemplateId:       contract.TemplateId,
					ContractId:       contract.ContractId,
					CreatedEventBlob: contract.CreatedEventBlob,
					SynchronizerId:   contract.SynchronizerId,
				}
			}

			choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
			if err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to convert choice context: %w", err)
			}

			return resp.JSON200.FactoryId, choiceContext, disclosedContracts, nil
		}, nil
	}

	return nil, nil //nolint:nilnil
}

type burnMintFactory func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error)

func getBurnMintFactory(ctx context.Context, acs store.ActiveContractStoreInterface, cfg config.BurnMintFactory) (burnMintFactory, error) {
	switch cfg.Type {
	case config.FactoryTypeDisabled:
		return nil, nil //nolint:nilnil
	case config.FactoryTypeAddress:
		factoryAddress := *cfg.InstanceAddress

		templateId, err := contracts.TemplateIDFromString(*cfg.TemplateId)
		if err != nil {
			return nil, fmt.Errorf("invalid TemplateId for TransferPreapproval: %w", err)
		}
		acs.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: templateId,
			PartyID:    *cfg.Party,
		})

		return func(ctx context.Context, _ splice_api_token_holding_v1.InstrumentId) (string, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
			activeTransferFactory, ok := acs.Get(factoryAddress)
			if !ok {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", factoryAddress)
			}

			return activeTransferFactory.GetCreatedEvent().GetContractId(), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
		}, nil
	case config.FactoryTypeURL:
		// DA utility-registry factory resolution endpoint (e.g. /mint/v0/request):
		//   POST {TokenStandardURL}  (URL is fully configured, including path)
		//   Request:  { holder, instrumentId }  (mint endpoint shape)
		//   Response: FactoryWithChoiceContext { factoryId, choiceContext.{choiceContextData, disclosedContracts} }
		//
		// We use the /mint/v0/request endpoint because DA's /burn-mint-factory endpoint
		// currently has a bug where issuer-credentials is returned as [] while
		// /mint/v0/request correctly populates them. Both endpoints return the same
		// factoryId (AllocationFactory) and the same context keys that
		// AllocationFactory_InternalBurnMint reads (instrument-configuration +
		// issuer-credentials). The mint endpoint is read-only — it does not actually
		// mint; it returns the factory + context needed to later exercise the choice.
		// Revisit this when DA fixes the /burn-mint-factory endpoint to populate
		// issuer-credentials directly.
		//
		// Auth is optional via TokenStandardAuthConfig (DA's public endpoints have security: []).
		httpClient := &http.Client{}

		var requestEditor transferInstructionV1.RequestEditorFn
		if cfg.TokenStandardAuthConfig != nil {
			authProvider, err := cfg.TokenStandardAuthConfig.NewProvider(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create auth provider for BurnMintFactory: %w", err)
			}
			requestEditor = func(ctx context.Context, req *http.Request) error {
				token, err := authProvider.TokenSource().Token()
				if err != nil {
					return fmt.Errorf("failed to retrieve token: %w", err)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

				return nil
			}
		}

		factoryURL := *cfg.TokenStandardURL

		return func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
			// DA's mint request body shape: { holder, instrumentId }.
			// holder is required by the schema but the backend does not use it for
			// factory/context resolution — the choiceContext is keyed by instrumentId.
			// Send instrumentId.admin as a dummy holder (same party, harmless).
			requestBody := map[string]any{
				"holder": instrumentId.Admin,
				"instrumentId": map[string]any{
					"admin": instrumentId.Admin,
					"id":    instrumentId.Id,
				},
			}
			bodyBytes, err := json.Marshal(requestBody)
			if err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to marshal BurnMintFactory request: %w", err)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, factoryURL, bytes.NewReader(bodyBytes))
			if err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to create BurnMintFactory request: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")
			if requestEditor != nil {
				if err := requestEditor(ctx, req); err != nil {
					return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to apply auth to BurnMintFactory request: %w", err)
				}
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to call BurnMintFactory endpoint %q: %w", factoryURL, err)
			}
			defer resp.Body.Close()

			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to read BurnMintFactory response: %w", err)
			}
			if resp.StatusCode != http.StatusOK {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("BurnMintFactory endpoint %q returned status %d: %s", factoryURL, resp.StatusCode, string(respBytes))
			}

			var factoryResp struct {
				FactoryId     string `json:"factoryId"`
				ChoiceContext struct {
					ChoiceContextData  map[string]any `json:"choiceContextData"`
					DisclosedContracts []struct {
						TemplateId       string `json:"templateId"`
						ContractId       string `json:"contractId"`
						CreatedEventBlob string `json:"createdEventBlob"`
						SynchronizerId   string `json:"synchronizerId"`
					} `json:"disclosedContracts"`
				} `json:"choiceContext"`
			}
			if err := json.Unmarshal(respBytes, &factoryResp); err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to parse BurnMintFactory response: %w", err)
			}

			disclosedContracts := make([]oapiCommon.DisclosedContract, len(factoryResp.ChoiceContext.DisclosedContracts))
			for i, c := range factoryResp.ChoiceContext.DisclosedContracts {
				disclosedContracts[i] = oapiCommon.DisclosedContract{
					TemplateId:       c.TemplateId,
					ContractId:       c.ContractId,
					CreatedEventBlob: c.CreatedEventBlob,
					SynchronizerId:   c.SynchronizerId,
				}
			}

			choiceContext, err := contracts.ChoiceContextFromData(factoryResp.ChoiceContext.ChoiceContextData)
			if err != nil {
				return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to convert BurnMintFactory choice context: %w", err)
			}

			return factoryResp.FactoryId, choiceContext, disclosedContracts, nil
		}, nil
	}

	return nil, nil //nolint:nilnil
}
