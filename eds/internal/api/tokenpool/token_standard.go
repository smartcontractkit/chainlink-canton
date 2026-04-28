package tokenpool

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

type transferFactory func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, common.CCIPContext, []oapiCommon.DisclosedContract, error)

func getTransferFactory(ctx context.Context, poolOwner types.PARTY, acs *store.ActiveContractStore, cfg config.TransferFactory) (transferFactory, error) {
	switch cfg.Type {
	case config.FactoryTypeAddress:
		factoryAddress := *cfg.Address

		templateId, err := contracts.TemplateIDFromString(*cfg.TemplateId)
		if err != nil {
			return nil, fmt.Errorf("invalid TemplateId for TransferPreapproval: %w", err)
		}
		acs.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: templateId,
			PartyID:    *cfg.Party,
		})

		return func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, common.CCIPContext, []oapiCommon.DisclosedContract, error) {
			activeTransferFactory, ok := acs.Get(factoryAddress)
			if !ok {
				return "", common.CCIPContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", factoryAddress)
			}

			return activeTransferFactory.GetCreatedEvent().GetContractId(), common.CCIPContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
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

		return func(ctx context.Context, instrumentId splice_api_token_holding_v1.InstrumentId) (string, common.CCIPContext, []oapiCommon.DisclosedContract, error) {
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
				return "", common.CCIPContext{}, nil, fmt.Errorf("failed to call GetTransferFactory: %w", err)
			}
			if resp.StatusCode() != http.StatusOK {
				return "", common.CCIPContext{}, nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
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

			ccipContext, err := contracts.CCIPContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
			if err != nil {
				return "", common.CCIPContext{}, nil, fmt.Errorf("failed to convert choice context: %w", err)
			}

			return resp.JSON200.FactoryId, ccipContext, disclosedContracts, nil
		}, nil
	}

	return nil, nil
}

type burnMintFactory func(ctx context.Context) (string, []oapiCommon.DisclosedContract, error)

func getBurnMintFactory(acs *store.ActiveContractStore, cfg config.BurnMintFactory) (burnMintFactory, error) {
	switch cfg.Type {
	case config.FactoryTypeAddress:
		factoryAddress := *cfg.Address

		templateId, err := contracts.TemplateIDFromString(*cfg.TemplateId)
		if err != nil {
			return nil, fmt.Errorf("invalid TemplateId for TransferPreapproval: %w", err)
		}
		acs.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: templateId,
			PartyID:    *cfg.Party,
		})

		return func(ctx context.Context) (string, []oapiCommon.DisclosedContract, error) {
			activeTransferFactory, ok := acs.Get(factoryAddress)
			if !ok {
				return "", nil, fmt.Errorf("no active contract found for transfer factory at address %s", factoryAddress)
			}

			return activeTransferFactory.GetCreatedEvent().GetContractId(), []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
		}, nil
	}

	return nil, nil
}
