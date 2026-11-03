package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/eds/api/common"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/openapi/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
)

func NewTransferFactory(ctx context.Context, poolOwner types.PARTY, acs store.ActiveContractStoreInterface, cfg config.Factory) (TransferDisclosureFactory, error) {
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

		// DA's bridge API lives on the same host as the token-standard API.
		bridgeURL := ""
		if u, err := url.Parse(*cfg.TokenStandardURL); err == nil && u.Host != "" {
			bridgeURL = u.Scheme + "://" + u.Host + "/api/bridge/v0/ccip/holdings"
		}

		return URLTransferFactory{
			poolOwner:                 poolOwner,
			transferInstructionClient: transferInstructionClient,
			bridgeHoldingsURL:         bridgeURL,
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

func (f AddressTransferFactory) GetFeeTokenDisclosures(
	ctx context.Context,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	sender types.PARTY,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	// The address-based transfer factory is instrument-independent: it is simply the
	// factory contract itself.
	activeTransferFactory, ok := f.acs.Get(f.factoryAddress)
	if !ok {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no active contract found for transfer factory at address %s", f.factoryAddress)
	}

	return types.CONTRACT_ID(activeTransferFactory.GetCreatedEvent().GetContractId()), splice_api_token_metadata_v1.ChoiceContext{}, []oapiCommon.DisclosedContract{converters.ActiveContractToDisclosedContract(activeTransferFactory)}, nil
}

func (f AddressTransferFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message, _ map[string]types.TEXT) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
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
	bridgeHoldingsURL         string
}

func (f URLTransferFactory) GetSendDisclosures(ctx context.Context, message oapiCommon.Message, poolReceiveMeta map[string]types.TEXT) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	if message.TokenTransfer == nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("no TokenTransfer in message")
	}
	tokenTransfer := message.TokenTransfer
	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(tokenTransfer.Token.Admin),
		Id:    types.TEXT(tokenTransfer.Token.Id),
	}
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

	factoryID, choiceContext, disclosedContracts, err := f.getTransferFactory(ctx, instrumentId, inputHoldingCids, sender, f.poolOwner, poolReceiveMeta)
	if err != nil {
		return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, err
	}

	// A named receiver label switches the backend to two-step ("offer") mode,
	// which omits the direct-mode context entries and disclosures the offer path
	// still requires — merge the direct-mode response in.
	if len(poolReceiveMeta) > 0 {
		_, directContext, directDisclosures, err := f.getTransferFactory(ctx, instrumentId, inputHoldingCids, sender, f.poolOwner, nil)
		if err != nil {
			return "", splice_api_token_metadata_v1.ChoiceContext{}, nil, fmt.Errorf("failed to fetch direct-mode context for merge: %w", err)
		}
		for k, v := range directContext.Values {
			if _, ok := choiceContext.Values[k]; !ok {
				choiceContext.Values[k] = v
			}
		}
		seen := make(map[string]bool, len(disclosedContracts))
		for _, d := range disclosedContracts {
			seen[d.ContractId] = true
		}
		for _, d := range directDisclosures {
			if !seen[d.ContractId] {
				disclosedContracts = append(disclosedContracts, d)
			}
		}
	}

	return factoryID, choiceContext, disclosedContracts, nil
}

func (f URLTransferFactory) GetFeeTokenDisclosures(
	ctx context.Context,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	sender types.PARTY,
	receiver types.PARTY,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	return f.getTransferFactory(ctx, instrumentId, []string{}, sender, receiver, nil)
}

func (f URLTransferFactory) getTransferFactory(
	ctx context.Context,
	instrumentId splice_api_token_holding_v1.InstrumentId,
	inputHoldingCids []string,
	sender types.PARTY,
	receiver types.PARTY,
	meta map[string]types.TEXT,
) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error) {
	// The receiver label in the transfer meta selects the context mode: a named
	// label (e.g. cip-112/receiver.id) yields the two-step context, "" the
	// direct/preapproved one.
	metaValues := make(map[string]any, len(meta))
	for k, v := range meta {
		metaValues[k] = string(v)
	}
	resp, err := f.transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": instrumentId.Admin,
			"transfer": map[string]any{
				"sender":   sender,
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
					"values": metaValues,
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
	inputHoldingCidStrings := make([]string, len(inputHoldingCids))
	for i, cid := range inputHoldingCids {
		inputHoldingCidStrings[i] = string(cid)
	}

	return f.getTransferFactory(ctx, instrumentId, inputHoldingCidStrings, f.poolOwner, receiver, nil)
}

// bridgeHoldingsResponse matches DA's GET /api/bridge/v0/ccip/holdings response.
type bridgeHoldingsResponse struct {
	DisclosedHoldings []struct {
		ContractId       string `json:"contractId"`
		TemplateId       string `json:"templateId"`
		SynchronizerId   string `json:"synchronizerId"`
		CreatedEventBlob string `json:"createdEventBlob"`
	} `json:"disclosedHoldings"`
}

// GetBridgeHoldings fetches the owner's publicly disclosed holdings from DA's
// bridge API, when the local holdings store cannot see them. instrumentId is
// the bare id ("TESTLINK6LR"), not the "id@admin" form.
func (f URLTransferFactory) GetBridgeHoldings(
	ctx context.Context,
	owner types.PARTY,
	instrumentId splice_api_token_holding_v1.InstrumentId,
) ([]oapiCommon.DisclosedContract, error) {
	if f.bridgeHoldingsURL == "" {
		return nil, fmt.Errorf("bridge holdings URL unavailable (token-standard URL not set)")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.bridgeHoldingsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build bridge holdings request: %w", err)
	}
	q := req.URL.Query()
	q.Set("adminPartyId", string(owner))
	q.Set("instrumentId", string(instrumentId.Id))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bridge holdings request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("bridge holdings returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed bridgeHoldingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("bridge holdings response: %w", err)
	}

	out := make([]oapiCommon.DisclosedContract, len(parsed.DisclosedHoldings))
	for i, h := range parsed.DisclosedHoldings {
		out[i] = oapiCommon.DisclosedContract{
			ContractId:       h.ContractId,
			TemplateId:       h.TemplateId,
			CreatedEventBlob: h.CreatedEventBlob,
			SynchronizerId:   h.SynchronizerId,
		}
	}

	return out, nil
}
