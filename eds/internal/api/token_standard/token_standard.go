package token_standard

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/global"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

const (
	DefaultPageSize = int32(25)
	MaximumPageSize = int32(1024)
)

type TokenConfig struct {
	Type            config.TokenType
	RegistryAddress contracts.InstanceAddress
	TokenID         string
	TokenName       string
	TokenSymbol     string
}

type tokenSupply struct {
	UpdatedAt   time.Time
	TotalSupply string
}

type Server struct {
	logger                 zerolog.Logger
	activeContractStore    store.ActiveContractStoreInterface
	instrumentHoldingStore store.InstrumentHoldingStoreInterface

	admin types.PARTY
	// maps Token ID -> TokenConfig
	tokens map[types.TEXT]TokenConfig

	// Caching a token's totalSupply for the configured duration.
	// Calculating the total supply is somewhat expensive, as it requires fetching all holdings for the instrument of which there could be many.
	tokenSupplyCacheTimeout time.Duration
	// Keyed by Token ID, this map stores the last known total supply for each token and the time it was last updated.
	tokenSupplies  map[string]tokenSupply
	tokenSupplyMux *sync.Mutex
}

var (
	_ oapiTokenMetadataV1.ServerInterface     = &Server{}
	_ oapiTransferInstruction.ServerInterface = &Server{}
)

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore store.ActiveContractStoreInterface,
	instrumentHoldingStore store.InstrumentHoldingStoreInterface,
	cfg config.TokenStandardAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:                  logger.With().Str("component", "TokenStandardAPI").Logger(),
		activeContractStore:     activeContractStore,
		instrumentHoldingStore:  instrumentHoldingStore,
		admin:                   types.PARTY(cfg.Admin),
		tokens:                  make(map[types.TEXT]TokenConfig),
		tokenSupplyCacheTimeout: cfg.SupplyCacheTTL,
		tokenSupplies:           make(map[string]tokenSupply),
		tokenSupplyMux:          new(sync.Mutex),
	}

	for _, registry := range cfg.Registries {
		// Validate config
		if len(registry.TokenId) == 0 {
			return nil, fmt.Errorf("empty TokenId in config")
		}
		if _, ok := s.tokens[types.TEXT(registry.TokenId)]; ok {
			return nil, fmt.Errorf("duplicate TokenId in config: %s", registry.TokenId)
		}

		// Default to the token's ID if name/symbol are not set
		tokenName := registry.TokenId
		if registry.TokenName != "" {
			tokenName = registry.TokenName
		}
		tokenSymbol := registry.TokenId
		if registry.TokenSymbol != "" {
			tokenSymbol = registry.TokenSymbol
		}

		s.tokens[types.TEXT(registry.TokenId)] = TokenConfig{
			Type:            registry.TokenType,
			RegistryAddress: registry.InstanceAddress,
			TokenID:         registry.TokenId,
			TokenName:       tokenName,
			TokenSymbol:     tokenSymbol,
		}
		switch registry.TokenType {
		case config.TokenTypeLINK:
			s.activeContractStore.RegisterTemplates(
				store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(link.LinkRegistry{}),
					PartyID:    registry.PartyID,
				},
				store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(link.LinkTransferInstruction{}),
					PartyID:    registry.PartyID,
				},
				// Register LockLinkHoldings, needed to accurately calculate the total supply
				store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(link.LockedLinkHolding{}),
					PartyID:    registry.PartyID,
				},
				// Register LinkTransferPreapprovals, needed to return the correct preapproval from GetTransferFactory
				store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(link.LinkTransferPreapproval{}),
					PartyID:    registry.PartyID,
				},
			)
		default:
			return nil, fmt.Errorf("unsupported token type: %s", registry.TokenType)
		}
	}

	s.instrumentHoldingStore.RegisterParty(cfg.Admin)

	return s, nil
}

func (s Server) GetRegistryInfo(c *gin.Context) {
	c.JSON(http.StatusOK, oapiTokenMetadataV1.GetRegistryInfoResponse{
		AdminId: string(s.admin),
		SupportedApis: map[string]int32{
			"splice-api-token-metadata-v1": 1,
		},
	})
}

// Token Metadata V1

func (s Server) ListInstruments(c *gin.Context, params oapiTokenMetadataV1.ListInstrumentsParams) {
	pageSize := DefaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	if pageSize <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: "pageSize must be a positive integer"})
		return
	}
	if pageSize > MaximumPageSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: fmt.Sprintf("pageSize cannot be greater than %d", MaximumPageSize)})
		return
	}
	// Using the index of the alpha-sorted list of configured tokens as the pageToken
	// If not specified, it will start at 0
	// If there are more elements to be returned, the next index will be returned as the response's pageToken
	pageToken := 0
	if params.PageToken != nil {
		parsedIndex, err := strconv.Atoi(*params.PageToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: "invalid pageToken"})
			return
		}
		pageToken = parsedIndex
	}

	// Ensure the provided pageToken is valid
	if pageToken < 0 || pageToken >= len(s.tokens) {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: "invalid pageToken"})
		return
	}

	// Sort all configured tokens by ID
	sortedTokenConfigs := maps.Values(s.tokens)
	slices.SortFunc(sortedTokenConfigs, func(a, b TokenConfig) int {
		return strings.Compare(a.TokenID, b.TokenID)
	})

	endIndex := min(pageToken+int(pageSize), len(sortedTokenConfigs))

	instrumentIds := make([]TokenConfig, endIndex-pageToken)
	copy(instrumentIds, sortedTokenConfigs[pageToken:endIndex])

	instruments := make([]oapiTokenMetadataV1.Instrument, len(instrumentIds))
	for i, tokenConfig := range instrumentIds {
		totalSupply, err := s.getTotalSupplyForInstrument(tokenConfig.TokenID)
		if err != nil {
			s.logger.Error().Err(err).Str("tokenId", tokenConfig.TokenID).Msg("Failed to get total supply for instrument")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTokenMetadataV1.ErrorResponse{Error: "internal server error"})
			return
		}

		instruments[i] = oapiTokenMetadataV1.Instrument{
			Decimals: 10,
			Id:       tokenConfig.TokenID,
			Name:     tokenConfig.TokenName,
			SupportedApis: map[string]int32{
				"splice-api-token-transfer-instruction-v1": 1,
			},
			Symbol:          tokenConfig.TokenSymbol,
			TotalSupply:     &totalSupply.TotalSupply,
			TotalSupplyAsOf: &totalSupply.UpdatedAt,
		}
	}

	var nextPageToken *string
	if endIndex < len(sortedTokenConfigs) {
		nextPageToken = new(strconv.Itoa(endIndex))
	}

	c.JSON(http.StatusOK, oapiTokenMetadataV1.ListInstrumentsResponse{
		Instruments:   instruments,
		NextPageToken: nextPageToken,
	})
}

func (s Server) GetInstrument(c *gin.Context, instrumentId string) {
	tokenConfig, ok := s.tokens[types.TEXT(instrumentId)]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTokenMetadataV1.ErrorResponse{Error: fmt.Sprintf("No instrument with id %s found", instrumentId)})
		return
	}

	totalSupply, err := s.getTotalSupplyForInstrument(tokenConfig.TokenID)
	if err != nil {
		s.logger.Error().Err(err).Str("tokenId", tokenConfig.TokenID).Msg("Failed to get total supply for instrument")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTokenMetadataV1.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, oapiTokenMetadataV1.Instrument{
		Decimals: 10,
		Id:       tokenConfig.TokenID,
		Name:     tokenConfig.TokenName,
		Symbol:   tokenConfig.TokenSymbol,
		SupportedApis: map[string]int32{
			"splice-api-token-transfer-instruction-v1": 1,
		},
		TotalSupply:     &totalSupply.TotalSupply,
		TotalSupplyAsOf: &totalSupply.UpdatedAt,
	})
}

func (s Server) getTotalSupplyForInstrument(instrumentId string) (tokenSupply, error) {
	s.tokenSupplyMux.Lock()
	defer s.tokenSupplyMux.Unlock()

	// Check cache first
	cachedSupply, ok := s.tokenSupplies[instrumentId]
	if ok && time.Since(cachedSupply.UpdatedAt) < s.tokenSupplyCacheTimeout {
		return cachedSupply, nil
	}

	// If not in cache or cache is stale, calculate total supply
	tokenConfig := s.tokens[types.TEXT(instrumentId)]
	snapshotTime := time.Now()
	holdings, err := s.instrumentHoldingStore.ListHoldings(splice_api_token_holding_v1.InstrumentId{
		Admin: s.admin,
		Id:    types.TEXT(instrumentId),
	})
	if err != nil {
		return tokenSupply{}, fmt.Errorf("failed to list holdings for instrument %s: %w", instrumentId, err)
	}

	totalSupply := new(big.Rat)
	switch tokenConfig.Type {
	case config.TokenTypeLINK:
		// Iterate over all LINK holdings and sum them up
		for i, holding := range holdings {
			parsedHolding, err := bindings.UnmarshalCreatedEvent[link.LinkHolding](holding.GetCreatedEvent())
			if err != nil {
				return tokenSupply{}, fmt.Errorf("failed to unmarshal LinkHolding at index %d: %w", i, err)
			}
			parsedAmount, ok := new(big.Rat).SetString(string(parsedHolding.HoldingAmount))
			if !ok {
				return tokenSupply{}, fmt.Errorf("failed to parse HoldingAmount as big.Rat at index %d: %s", i, parsedHolding.HoldingAmount)
			}
			totalSupply.Add(totalSupply, parsedAmount)
		}

		// List LockLinkHoldings to also sum up
		lockedHoldings, _ := s.activeContractStore.GetByTemplateId(s.admin, contracts.TemplateIDFromBinding(link.LockedLinkHolding{}))
		for _, activeContract := range lockedHoldings {
			parsedLockedHolding, err := bindings.UnmarshalCreatedEvent[link.LockedLinkHolding](activeContract.GetCreatedEvent())
			if err != nil {
				return tokenSupply{}, fmt.Errorf("failed to unmarshal LockedLinkHolding: %w", err)
			}
			// Skipped holdings that might be for another instrument
			if parsedLockedHolding.LockedInstrumentId.Admin != s.admin ||
				parsedLockedHolding.LockedInstrumentId.Id != types.TEXT(instrumentId) {
				continue
			}
			parsedAmount, ok := new(big.Rat).SetString(string(parsedLockedHolding.LockedAmount))
			if !ok {
				return tokenSupply{}, fmt.Errorf("failed to parse LockedAmount as big.Rat: %s", parsedLockedHolding.LockedAmount)
			}
			totalSupply.Add(totalSupply, parsedAmount)
		}
	}

	// Cache totalSupply
	tokenSupply := tokenSupply{
		UpdatedAt:   snapshotTime,
		TotalSupply: totalSupply.FloatString(10), // 10 decimal places
	}
	s.tokenSupplies[instrumentId] = tokenSupply

	return tokenSupply, nil
}

func (s Server) GetTransferFactory(c *gin.Context) {
	var req oapiTransferInstruction.GetFactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, oapiTransferInstruction.ErrorResponse{Error: "request body too large"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: err.Error()})

		return
	}

	var choiceArguments splice_api_token_transfer_instruction_v1.TransferFactoryTransfer
	if err := ledger.MapToStruct(req.ChoiceArguments, &choiceArguments); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("invalid `choiceArguments` format: %s", err.Error())})
		return
	}

	if choiceArguments.ExpectedAdmin != s.admin {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: "invalid `expectedAdmin`, must match registry admin"})
		return
	}
	instrumentId := choiceArguments.Transfer.InstrumentId
	if instrumentId.Admin != s.admin {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: "invalid `transfer.instrumentId.admin`, must match registry admin"})
		return
	}

	requestedId := instrumentId.Id

	cfg, ok := s.tokens[requestedId]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No instrument with id %s found", requestedId)})
		return
	}

	var disclosedContracts []oapiTransferInstruction.DisclosedContract
	activeRegistryContract, ok := s.activeContractStore.Get(cfg.RegistryAddress)
	if !ok {
		s.logger.Error().Str("registryAddress", cfg.RegistryAddress.String()).Msg("Failed to retrieve active registry contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTransferInstruction.ErrorResponse{Error: "internal server error"})
		return
	}
	disclosedRegistry := ActiveContractToDisclosedContract(activeRegistryContract, req.ExcludeDebugFields != nil && !*req.ExcludeDebugFields)
	disclosedContracts = append(disclosedContracts, disclosedRegistry)

	// Look up the (optional) transfer preapproval for the receiver
	transferPreapproval, err := s.getLinkTransferPreApprovalForParty(s.admin, choiceArguments.Transfer.Receiver)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to retrieve transfer preapproval")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTransferInstruction.ErrorResponse{Error: "internal server error"})
		return
	}
	choiceContext := splice_api_token_metadata_v1.ChoiceContext{Values: make(map[string]splice_api_token_metadata_v1.AnyValue)}
	if transferPreapproval != nil {
		choiceContext.Values[string(link.TransferPreapprovalContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(transferPreapproval.GetCreatedEvent().GetContractId())),
		}
		disclosedPreapproval := ActiveContractToDisclosedContract(transferPreapproval, req.ExcludeDebugFields != nil && !*req.ExcludeDebugFields)
		disclosedContracts = append(disclosedContracts, disclosedPreapproval)
	}

	contextData, err := converters.SerializeChoiceContext(choiceContext)
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := oapiTransferInstruction.TransferFactoryWithChoiceContext{
		ChoiceContext: oapiTransferInstruction.ChoiceContext{
			ChoiceContextData:  contextData,
			DisclosedContracts: disclosedContracts,
		},
		FactoryId:    activeRegistryContract.GetCreatedEvent().GetContractId(),
		TransferKind: oapiTransferInstruction.Offer,
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) getLinkTransferPreApprovalForParty(admin, receiver types.PARTY) (*apiv2.ActiveContract, error) {
	preapprovals, _ := s.activeContractStore.GetByTemplateId(admin, contracts.TemplateIDFromBinding(link.LinkTransferPreapproval{}))
	for _, ac := range preapprovals {
		preapproval, err := bindings.UnmarshalCreatedEvent[link.LinkTransferPreapproval](ac.GetCreatedEvent())
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal LinkTransferPreapproval: %w", err)
		}
		if preapproval.PreapprovalReceiver == receiver {
			return ac, nil
		}
	}

	return nil, nil //nolint:nilnil // not finding a preapproval is expected
}

// Transfer Instruction V1

func (s Server) GetTransferInstructionAcceptContext(c *gin.Context, transferInstructionId string) {
	var req oapiTransferInstruction.GetChoiceContextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, oapiTransferInstruction.ErrorResponse{Error: "request body too large"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: err.Error()})

		return
	}

	transferInstruction, ok := s.activeContractStore.GetByContractId(types.CONTRACT_ID(transferInstructionId))
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}
	parsedTransferInstruction, err := bindings.UnmarshalCreatedEvent[link.LinkTransferInstruction](transferInstruction.GetCreatedEvent())
	if err != nil {
		// If the unmarshal fails, we assume that the user provided an incorrect ContractId/a ContractId of a differently-types contract
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}
	instrumentId := parsedTransferInstruction.InstructionTransfer.InstrumentId
	if instrumentId.Admin != s.admin {
		// If the TransferInstruction is for a different instrument admin, we treat it as not-found
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}
	cfg, ok := s.tokens[instrumentId.Id]
	if !ok {
		// If the TransferInstruction is for a non-configured instrument, we treat it as not-found
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}

	// Request is validated

	activeRegistryContract, ok := s.activeContractStore.Get(cfg.RegistryAddress)
	if !ok {
		s.logger.Error().Str("registryAddress", cfg.RegistryAddress.String()).Msg("Failed to retrieve active registry contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTransferInstruction.ErrorResponse{Error: "internal server error"})
		return
	}
	disclosedRegistry := ActiveContractToDisclosedContract(activeRegistryContract, false)

	resp := oapiTransferInstruction.ChoiceContext{
		ChoiceContextData: nil,
		// Don't return a disclosure for the TransferInstruction itself - the receiver will already be an observer on that contract
		DisclosedContracts: []oapiTransferInstruction.DisclosedContract{disclosedRegistry},
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) GetTransferInstructionRejectContext(c *gin.Context, transferInstructionId string) {
	s.GetTransferInstructionAcceptContext(c, transferInstructionId)
}

func (s Server) GetTransferInstructionWithdrawContext(c *gin.Context, transferInstructionId string) {
	s.GetTransferInstructionAcceptContext(c, transferInstructionId)
}

var _ global.InstanceAddressFilter = &Server{}

func (s Server) FilterContracts(addresses []contracts.InstanceAddress) []contracts.InstanceAddress {
	var out []contracts.InstanceAddress
	for _, address := range addresses {
		for _, tokenConfig := range s.tokens {
			if address == tokenConfig.RegistryAddress {
				out = append(out, address)
			}
		}
	}

	return out
}
