package factory

import (
	"context"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/daRegistry"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

func NewBurnMintFactory(ctx context.Context, poolOwner types.PARTY, acs store.ActiveContractStoreInterface, cfg config.Factory) (DisclosureFactory, error) {
	switch cfg.Type {
	case config.FactoryTypeDisabled:
		return nil, nil //nolint:nilnil
	case config.FactoryTypeAddress:
		factoryAddress := *cfg.InstanceAddress

		templateId, err := contracts.TemplateIDFromString(*cfg.TemplateId)
		if err != nil {
			return nil, fmt.Errorf("invalid TemplateId for BurnMintFactory: %w", err)
		}
		acs.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: templateId,
			PartyID:    *cfg.Party,
		})

		return AddressBurnMintFactory{
			factoryAddress: factoryAddress,
			acs:            acs,
		}, nil
	case config.FactoryTypeURL, config.FactoryTypeURLRequests:
		// If authentication has been configured, add an interceptor that adds the Authorization header
		var options []daRegistry.ClientOption
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
			options = append(options, daRegistry.WithRequestEditorFn(interceptor))
		}
		daRegistryClient, err := daRegistry.NewClientWithResponses(
			*cfg.TokenStandardURL,
			options...,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create DARegistry client with URL %q: %w", *cfg.TokenStandardURL, err)
		}

		//nolint:exhaustive // these are the only two possible types
		switch cfg.Type {
		case config.FactoryTypeURL:
			return URLBurnMintFactory{
				poolOwner:        poolOwner,
				daRegistryClient: daRegistryClient,
			}, nil
		case config.FactoryTypeURLRequests:
			return RequestBurnMintFactory{
				poolOwner:        poolOwner,
				daRegistryClient: daRegistryClient,
			}, nil
		}
	}

	return nil, nil //nolint:nilnil
}

// AddressBurnMintFactory is a DisclosureFactory implementation that looks up a BurnMintFactory
// from an InstanceAddress directly.
type AddressBurnMintFactory struct {
	factoryAddress contracts.InstanceAddress
	acs            store.ActiveContractStoreInterface
}

func (f AddressBurnMintFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	activeBurnMintFactory, ok := f.acs.Get(f.factoryAddress)
	if !ok {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", f.factoryAddress)
	}

	return types.CONTRACT_ID(activeBurnMintFactory.GetCreatedEvent().GetContractId()), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeBurnMintFactory)}, nil
}

func (f AddressBurnMintFactory) GetExecuteDisclosures(ctx context.Context,
	message *protocol.Message,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	inputHoldingCids []types.CONTRACT_ID,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	activeBurnMintFactory, ok := f.acs.Get(f.factoryAddress)
	if !ok {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", f.factoryAddress)
	}

	return types.CONTRACT_ID(activeBurnMintFactory.GetCreatedEvent().GetContractId()), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeBurnMintFactory)}, nil
}

// URLBurnMintFactory is a DisclosureFactory implementation that requests a BurnMintFactory from another URL.
// It calls DA's Registry to retrieve both the factory and ChoiceContext using the getBurnMintFactory endpoint for
// both send & execute directions.
type URLBurnMintFactory struct {
	poolOwner        types.PARTY
	daRegistryClient daRegistry.ClientWithResponsesInterface
}

func (f URLBurnMintFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}
	tokenTransfer := *message.TokenTransfer
	instrumentId := tokenTransfer.Token
	// API expects empty arrays to be `[]`, not `null`
	inputHoldingCids := []string{}
	if tokenTransfer.HoldingContractIds != nil {
		inputHoldingCids = append(inputHoldingCids, *tokenTransfer.HoldingContractIds...)
	}

	resp, err := f.daRegistryClient.GetBurnMintFactoryWithResponse(ctx, daRegistry.GetBurnMintFactoryRequest{
		InstrumentId: daRegistry.InstrumentId{
			Admin: instrumentId.Admin,
			Id:    instrumentId.Id,
		},
		InputHoldingCids: inputHoldingCids,
		Outputs:          nil, // TODO: technically, this should contain the created change
	})
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to call GetBurnMintFactory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	disclosedContracts := make([]oapiCommon.DisclosedContract, len(resp.JSON200.ChoiceContext.DisclosedContracts))
	for i, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		synchronizerId := ""
		if contract.SynchronizerId != nil {
			synchronizerId = *contract.SynchronizerId
		}
		disclosedContracts[i] = oapiCommon.DisclosedContract{
			TemplateId:       contract.TemplateId,
			ContractId:       contract.ContractId,
			CreatedEventBlob: contract.CreatedEventBlob,
			SynchronizerId:   synchronizerId,
		}
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	return types.CONTRACT_ID(resp.JSON200.FactoryId), choiceContext, disclosedContracts, nil
}

func (f URLBurnMintFactory) GetExecuteDisclosures(ctx context.Context,
	message *protocol.Message,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	_ []types.CONTRACT_ID,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}

	// TODO: for backward-compatibility, use poolOwner as receiver if not specified
	if receiver == "" {
		receiver = f.poolOwner
	}

	resp, err := f.daRegistryClient.GetBurnMintFactoryWithResponse(ctx, daRegistry.GetBurnMintFactoryRequest{
		InstrumentId: daRegistry.InstrumentId{
			Admin: string(instrumentId.Admin),
			Id:    string(instrumentId.Id),
		},
		InputHoldingCids: []string{},
		Outputs: []daRegistry.MintOutput{
			{
				// TODO: Amount should be taken from message.TokenTransfer, but would have to be properly scaled by the TP's decimals
				Amount: "1.0",
				Owner:  string(receiver),
			},
		},
	})
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to call GetBurnMintFactory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	disclosedContracts := make([]oapiCommon.DisclosedContract, len(resp.JSON200.ChoiceContext.DisclosedContracts))
	for i, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		synchronizerId := ""
		if contract.SynchronizerId != nil {
			synchronizerId = *contract.SynchronizerId
		}
		disclosedContracts[i] = oapiCommon.DisclosedContract{
			TemplateId:       contract.TemplateId,
			ContractId:       contract.ContractId,
			CreatedEventBlob: contract.CreatedEventBlob,
			SynchronizerId:   synchronizerId,
		}
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	return types.CONTRACT_ID(resp.JSON200.FactoryId), choiceContext, disclosedContracts, nil
}

// RequestBurnMintFactory is a DisclosureFactory implementation that requests a BurnMintFactory from another URL.
// It calls DA's Registry to retrieve both the factory and ChoiceContext.
// Contrary to URLBurnMintFactory, it uses separate endpoints for send/execute:
// - getMintRequestCreateContext for execution
// - getBurnRequestCreateContext for sending
type RequestBurnMintFactory struct {
	poolOwner        types.PARTY
	daRegistryClient daRegistry.ClientWithResponsesInterface
}

func (f RequestBurnMintFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}
	tokenTransfer := *message.TokenTransfer
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

	resp, err := f.daRegistryClient.GetBurnRequestCreateContextWithResponse(ctx, daRegistry.RequestBurnRequest{
		InstrumentId: daRegistry.InstrumentId{
			Admin: instrumentId.Admin,
			Id:    instrumentId.Id,
		},
		Holder:             string(sender),
		HoldingContractIds: inputHoldingCids,
	})
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to call GetBurnRequestCreateContext: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	disclosedContracts := make([]oapiCommon.DisclosedContract, len(resp.JSON200.ChoiceContext.DisclosedContracts))
	for i, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		synchronizerId := ""
		if contract.SynchronizerId != nil {
			synchronizerId = *contract.SynchronizerId
		}
		disclosedContracts[i] = oapiCommon.DisclosedContract{
			TemplateId:       contract.TemplateId,
			ContractId:       contract.ContractId,
			CreatedEventBlob: contract.CreatedEventBlob,
			SynchronizerId:   synchronizerId,
		}
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	return types.CONTRACT_ID(resp.JSON200.FactoryId), choiceContext, disclosedContracts, nil
}

func (f RequestBurnMintFactory) GetExecuteDisclosures(ctx context.Context,
	message *protocol.Message,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	_ []types.CONTRACT_ID,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}

	// TODO: for backward-compatibility, use poolOwner as receiver if not specified
	if receiver == "" {
		receiver = f.poolOwner
	}

	resp, err := f.daRegistryClient.GetMintRequestCreateContextWithResponse(ctx, daRegistry.RequestMintRequest{
		InstrumentId: daRegistry.InstrumentId{
			Admin: string(instrumentId.Admin),
			Id:    string(instrumentId.Id),
		},
		Holder: string(receiver), // The holder of the minted holdings will be the receiver
	})
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to call GetMintRequestCreateContext: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	disclosedContracts := make([]oapiCommon.DisclosedContract, len(resp.JSON200.ChoiceContext.DisclosedContracts))
	for i, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		synchronizerId := ""
		if contract.SynchronizerId != nil {
			synchronizerId = *contract.SynchronizerId
		}
		disclosedContracts[i] = oapiCommon.DisclosedContract{
			TemplateId:       contract.TemplateId,
			ContractId:       contract.ContractId,
			CreatedEventBlob: contract.CreatedEventBlob,
			SynchronizerId:   synchronizerId,
		}
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	return types.CONTRACT_ID(resp.JSON200.FactoryId), choiceContext, disclosedContracts, nil
}
