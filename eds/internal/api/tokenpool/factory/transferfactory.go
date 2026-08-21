package factory

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
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

func NewTransferFactory(ctx context.Context, poolOwner types.PARTY, acs store.ActiveContractStoreInterface, cfg config.Factory) (DisclosureFactory, error) {
	switch cfg.Type {
	case config.FactoryTypeDisabled:
		return nil, nil //nolint:nilnil
	case config.FactoryTypeAddress:
		factoryAddress := *cfg.InstanceAddress

		templateId, err := contracts.TemplateIDFromString(*cfg.TemplateId)
		if err != nil {
			return nil, fmt.Errorf("invalid TemplateId for TransferFactory: %w", err)
		}
		acs.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: templateId,
			PartyID:    *cfg.Party,
		})

		return AddressTransferFactory{
			factoryAddress: factoryAddress,
			acs:            acs,
		}, nil
	case config.FactoryTypeURL:
		// If authentication has been configured, add an interceptor that adds the Authorization header
		var options []transferInstructionV1.ClientOption
		if cfg.TokenStandardAuthConfig != nil {
			authProvider, err := cfg.TokenStandardAuthConfig.NewProvider(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create auth provider: %w", err)
			}
			// Try to get a token to validate the auth works
			_, err = authProvider.TokenSource().Token()
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve token: %w", err)
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

		return URLTransferFactory{
			poolOwner:                 poolOwner,
			transferInstructionClient: transferInstructionClient,
		}, nil
	case config.FactoryTypeURLRequests:
		return nil, fmt.Errorf("invalid factory type %q: URLRequests is not supported for TransferFactory", cfg.Type)
	}

	return nil, nil //nolint:nilnil
}

type AddressTransferFactory struct {
	factoryAddress contracts.InstanceAddress
	acs            store.ActiveContractStoreInterface
}

func (f AddressTransferFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	activeTransferFactory, ok := f.acs.Get(f.factoryAddress)
	if !ok {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", f.factoryAddress)
	}

	return types.CONTRACT_ID(activeTransferFactory.GetCreatedEvent().GetContractId()), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
}

func (f AddressTransferFactory) GetExecuteDisclosures(
	ctx context.Context,
	message *protocol.Message,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	inputHoldingCids []types.CONTRACT_ID,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	activeTransferFactory, ok := f.acs.Get(f.factoryAddress)
	if !ok {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", f.factoryAddress)
	}

	return types.CONTRACT_ID(activeTransferFactory.GetCreatedEvent().GetContractId()), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
}

type URLTransferFactory struct {
	poolOwner                 types.PARTY
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface
}

func (f URLTransferFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}
	tokenTransfer := message.TokenTransfer
	instrumentId := tokenTransfer.Token
	// API expects empty arrays to be `[]`, not `null`
	inputHoldingCids := []string{}
	if tokenTransfer.HoldingContractIds != nil {
		inputHoldingCids = append(inputHoldingCids, *tokenTransfer.HoldingContractIds...)
	}

	// TODO: for backward-compatibility, use poolOwner as sender if receiver is not specified
	sender := f.poolOwner
	if message.Sender != "" {
		sender = types.PARTY(message.Sender)
	}

	resp, err := f.transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": instrumentId.Admin,
			"transfer": map[string]any{
				"sender":   sender,
				"receiver": f.poolOwner,
				// TODO: this isn't currently used. If we'd wanted to take the amount from message.TokenTransfer it would have to be properly scaled by the TP's decimals
				"amount": "1.0",
				"instrumentId": map[string]any{
					"admin": instrumentId.Admin,
					"id":    instrumentId.Id,
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(time.Second * -10).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				"inputHoldingCids": inputHoldingCids,
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

	return types.CONTRACT_ID(resp.JSON200.FactoryId), choiceContext, disclosedContracts, nil
}

func (f URLTransferFactory) GetExecuteDisclosures(
	ctx context.Context,
	message *protocol.Message,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	inputHoldingCids []types.CONTRACT_ID,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}
	// API expects empty arrays to be `[]`, not `null`
	if inputHoldingCids == nil {
		inputHoldingCids = []types.CONTRACT_ID{}
	}

	// TODO: for backward-compatibility, use poolOwner as receiver if not specified
	if receiver == "" {
		receiver = f.poolOwner
	}

	resp, err := f.transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": instrumentId.Admin,
			"transfer": map[string]any{
				"sender":   f.poolOwner,
				"receiver": receiver,
				// TODO: this isn't currently used. If we'd wanted to take the amount from message.TokenTransfer it would have to be properly scaled by the TP's decimals
				"amount": "1.0",
				"instrumentId": map[string]any{
					"admin": instrumentId.Admin,
					"id":    instrumentId.Id,
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(time.Second * -10).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				"inputHoldingCids": inputHoldingCids,
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

	return types.CONTRACT_ID(resp.JSON200.FactoryId), choiceContext, disclosedContracts, nil
}
