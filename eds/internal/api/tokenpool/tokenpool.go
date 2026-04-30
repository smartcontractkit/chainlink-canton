package tokenpool

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
)

type ContractConfig struct {
	Type            config.TokenPoolType
	Owner           types.PARTY
	transferFactory transferFactory
	burnMintFactory burnMintFactory
	preapproval     preapprovalFactory
}

type Server struct {
	logger                 zerolog.Logger
	activeContractStore    *store.ActiveContractStore
	instrumentHoldingStore *store.InstrumentHoldingStore

	contractConfigs map[contracts.InstanceAddress]ContractConfig
}

var _ oapiTokenPool.ServerInterface = &Server{}

func NewServer(
	ctx context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	instrumentHoldingStore *store.InstrumentHoldingStore,
	cfg config.TokenPoolAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:                 logger.With().Str("component", "TokenPoolAPI").Logger(),
		activeContractStore:    activeContractStore,
		instrumentHoldingStore: instrumentHoldingStore,
		contractConfigs:        make(map[contracts.InstanceAddress]ContractConfig),
	}
	for _, tokenPool := range cfg.TokenPools {
		contractConfig := ContractConfig{
			Type:  tokenPool.Type,
			Owner: types.PARTY(tokenPool.PoolOwner),
		}
		s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(common.RateLimiter{}),
			PartyID:    tokenPool.PartyID,
		})
		s.instrumentHoldingStore.RegisterParty(tokenPool.PoolOwner)
		s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(tokenadminregistry.TokenConfig{}),
			PartyID:    tokenPool.PartyID,
		})
		switch tokenPool.Type {
		case config.TokenPoolTypeLockRelease:
			s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}),
				PartyID:    tokenPool.PartyID,
			})
			// If the TransferFactory is configured, this will make this API automatically retrieve the necessary
			// ContractIds, Context, and disclosures from the instrument's TransferFactory.
			// If not enabled, users will have to get these information from the TransferFactory API themselves.
			if tokenPool.TransferFactory != nil {
				getFactoryFunc, err := getTransferFactory(ctx, types.PARTY(tokenPool.PoolOwner), activeContractStore, *tokenPool.TransferFactory)
				if err != nil {
					return nil, fmt.Errorf("failed to get transfer factory for token pool with address %s: %w", tokenPool.InstanceAddress, err)
				}
				contractConfig.transferFactory = getFactoryFunc
			}
		case config.TokenPoolTypeBurnMint:
			s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{}),
				PartyID:    tokenPool.PartyID,
			})
			if tokenPool.BurnMintFactory != nil {
				getFactoryFunc, err := getBurnMintFactory(activeContractStore, *tokenPool.BurnMintFactory)
				if err != nil {
					return nil, fmt.Errorf("failed to get burn mint factory for token pool with address %s: %w", tokenPool.InstanceAddress, err)
				}
				contractConfig.burnMintFactory = getFactoryFunc
			}
		default:
			return nil, fmt.Errorf("unsupported token pool type: %s", tokenPool.Type)
		}

		// If the TransferPreapproval is configured, this will make the API keep track of and return the given pre-approval
		if tokenPool.TransferPreapproval != nil {
			preapprovalFactoryFunc, err := getPreapprovalFactory(activeContractStore, tokenPool.TransferPreapproval.ContextKey, contractConfig.Owner, *tokenPool.TransferPreapproval)
			if err != nil {
				return nil, fmt.Errorf("failed to get preapproval factory for token pool with address %s: %w", tokenPool.InstanceAddress, err)
			}
			contractConfig.preapproval = preapprovalFactoryFunc
		}

		s.contractConfigs[tokenPool.InstanceAddress] = contractConfig
	}

	return s, nil
}

func (s Server) PostTokenPoolSend(c *gin.Context, address string) {
	var req oapiTokenPool.TokenPoolSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}
	instanceAddress, err := converters.ResolveAddress(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}
	cfg, ok := s.contractConfigs[instanceAddress]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "token pool address not found"})
		return
	}
	activeTokenPoolContract, ok := s.activeContractStore.Get(instanceAddress)
	if !ok {
		s.logger.Error().Stringer("address", instanceAddress).Msg("active token pool contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	destinationChainSelector, err := strconv.ParseUint(req.Message.DestinationChainSelector, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "invalid destination chain selector"})
		return
	}
	tokenTransfer := req.Message.TokenTransfer
	if tokenTransfer == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "message does not contain a token transfer"})
		return
	}

	// TODO check that the token transfer is for this pool

	switch cfg.Type {
	case config.TokenPoolTypeLockRelease:
		s.lockReleaseTokenPoolSend(c, cfg, instanceAddress, activeTokenPoolContract, destinationChainSelector, req.Message)
		return
	case config.TokenPoolTypeBurnMint:
		s.burnMintTokenPoolSend(c, cfg, instanceAddress, activeTokenPoolContract, destinationChainSelector, req.Message)
		return
	default:
		s.logger.Error().Stringer("address", instanceAddress).Msgf("unknown token pool type: %s", cfg.Type)
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}
}

func (s Server) lockReleaseTokenPoolSend(
	c *gin.Context,
	cfg ContractConfig,
	instanceAddress contracts.InstanceAddress,
	activeTokenPoolContract *apiv2.ActiveContract,
	destinationChainSelector uint64,
	_ oapiCommon.Message,
) {
	lockReleaseTokenPool, err := ParseLockReleaseTokenPool(activeTokenPoolContract.CreatedEvent)
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse lock release token pool contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	remoteChainConfig, ok := lockReleaseTokenPool.RemoteChainConfigs[destinationChainSelector]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("unsupported destination chain selector: %v", destinationChainSelector)})
		return
	}

	rateLimiter, ok := s.activeContractStore.Get(remoteChainConfig.OutboundRateLimiter)
	if !ok {
		s.logger.Error().Uint64("destinationChainSelector", destinationChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("outbound active rate limiter not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.OutboundCCVs))
	for i, v := range remoteChainConfig.OutboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	// Add Token Pool holdings to choiceContext - not currently used by the Token Pool, but may be required in the future
	holdings, ok := s.instrumentHoldingStore.GetHolding(cfg.Owner, lockReleaseTokenPool.InstrumentId)
	if !ok {
		s.logger.Error().Str("owner", string(cfg.Owner)).Any("instrumentId", lockReleaseTokenPool.InstrumentId).Msg("no holdings found for lock release token pool")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}
	tokenPoolHoldings := make([]splice_api_token_metadata_v1.AnyValue, len(holdings))
	disclosedHoldings := make([]oapiCommon.DisclosedContract, len(holdings))
	for i, holding := range holdings {
		tokenPoolHoldings[i] = splice_api_token_metadata_v1.AnyValue{AVContractId: new(types.CONTRACT_ID(holding.GetCreatedEvent().GetContractId()))}
		disclosedHoldings[i] = converters.ActiveContractToDisclosedContract(holding)
	}

	// The ChoiceContext that will be passed to the Token Pool
	choiceContext := splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{
		string(common.RateLimiterKey):                            {AVContractId: new(types.CONTRACT_ID(rateLimiter.GetCreatedEvent().GetContractId()))},
		string(lockreleasetokenpool.TokenPoolHoldingsContextKey): {AVList: new(tokenPoolHoldings)},
	}}

	// The ChoiceContext that will be passed to the TransferFactory by the Token Pool
	transferFactoryContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: make(map[string]splice_api_token_metadata_v1.AnyValue),
	}
	var factoryDisclosures []oapiCommon.DisclosedContract
	// Get ExtraArgs and TransferFactory from Token Standard API (if enabled)
	if cfg.transferFactory != nil {
		transferFactory, transferContext, disclosedFactoryContracts, err := cfg.transferFactory(c, lockReleaseTokenPool.InstrumentId)
		if err != nil {
			s.logger.Error().Err(err).Msg("transfer factory returned an error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
		choiceContext.Values[string(lockreleasetokenpool.TransferFactoryContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(transferFactory)),
		}
		transferFactoryContext.Values = transferContext.Values
		factoryDisclosures = append(factoryDisclosures, disclosedFactoryContracts...)
	}

	// Get TransferPreapproval (if enabled)
	if cfg.preapproval != nil {
		contextKey, activeTransferPreapproval, err := cfg.preapproval(c)
		if err != nil {
			s.logger.Err(err).Msg("failed to get transfer preapproval")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}

		transferFactoryContext.Values[contextKey] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(activeTransferPreapproval.GetCreatedEvent().GetContractId())),
		}
		factoryDisclosures = append(factoryDisclosures, converters.ActiveContractToDisclosedContract(activeTransferPreapproval))
	}

	// If the TransferFactory context contains any values, set it as part of the choiceContext
	if len(transferFactoryContext.Values) > 0 {
		choiceContext.Values[string(lockreleasetokenpool.TransferFactoryExtraArgsContextValuesContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVMap: new(transferFactoryContext.Values),
		}
	}

	contextData, err := converters.SerializeChoiceContext(choiceContext)
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiTokenPool.TokenPoolSendResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    lockReleaseTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: lockReleaseTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData:        contextData,
		DisclosedContracts: append(
			factoryDisclosures,
			converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
			converters.ActiveContractToDisclosedContract(rateLimiter),
		),
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) burnMintTokenPoolSend(
	c *gin.Context,
	cfg ContractConfig,
	instanceAddress contracts.InstanceAddress,
	activeTokenPoolContract *apiv2.ActiveContract,
	destinationChainSelector uint64,
	_ oapiCommon.Message,
) {
	burnMintTokenPool, err := ParseBurnMintTokenPool(activeTokenPoolContract.CreatedEvent)
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse lock release token pool contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	remoteChainConfig, ok := burnMintTokenPool.RemoteChainConfigs[destinationChainSelector]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("unsupported destination chain selector: %v", destinationChainSelector)})
		return
	}

	rateLimiter, ok := s.activeContractStore.Get(remoteChainConfig.OutboundRateLimiter)
	if !ok {
		s.logger.Error().Uint64("destinationChainSelector", destinationChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("outbound active rate limiter not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.OutboundCCVs))
	for i, v := range remoteChainConfig.OutboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	// The ChoiceContext that will be passed to the Token Pool
	choiceContext := splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{
		string(common.RateLimiterKey): {AVContractId: new(types.CONTRACT_ID(rateLimiter.GetCreatedEvent().GetContractId()))},
	}}

	// The ChoiceContext that will be passed to the BurnMintFactory by the Token Pool
	transferFactoryContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: make(map[string]splice_api_token_metadata_v1.AnyValue),
	}

	// Get BurnMintFactory (if enabled)
	var factoryDisclosures []oapiCommon.DisclosedContract
	if cfg.burnMintFactory != nil {
		burnMintFactory, disclosedFactoryContracts, err := cfg.burnMintFactory(c)
		if err != nil {
			s.logger.Error().Err(err).Msg("transfer factory returned an error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
		choiceContext.Values[string(burnminttokenpool.BurnMintFactoryContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(burnMintFactory)),
		}
		factoryDisclosures = append(factoryDisclosures, disclosedFactoryContracts...)
	}

	// Get TransferPreapproval (if enabled)
	if cfg.preapproval != nil {
		contextKey, activeTransferPreapproval, err := cfg.preapproval(c)
		if err != nil {
			s.logger.Err(err).Msg("failed to get transfer preapproval")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}

		transferFactoryContext.Values[contextKey] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(activeTransferPreapproval.GetCreatedEvent().GetContractId())),
		}
		factoryDisclosures = append(factoryDisclosures, converters.ActiveContractToDisclosedContract(activeTransferPreapproval))
	}

	// If the TransferFactory context contains any values, set it as part of the choiceContext
	if len(transferFactoryContext.Values) > 0 {
		choiceContext.Values[string(burnminttokenpool.BurnMintFactoryExtraArgsContextValuesContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVMap: new(transferFactoryContext.Values),
		}
	}

	contextData, err := converters.SerializeChoiceContext(choiceContext)
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiTokenPool.TokenPoolSendResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    burnMintTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: burnMintTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData:        contextData,
		DisclosedContracts: append(
			factoryDisclosures,
			converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
			converters.ActiveContractToDisclosedContract(rateLimiter),
		),
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) PostTokenPoolExecute(c *gin.Context, address string) {
	var req oapiTokenPool.TokenPoolExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}
	instanceAddress, err := converters.ResolveAddress(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}
	cfg, ok := s.contractConfigs[instanceAddress]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "token pool address not found"})
		return
	}
	activeTokenPoolContract, ok := s.activeContractStore.Get(instanceAddress)
	if !ok {
		s.logger.Error().Stringer("address", instanceAddress).Msg("active token pool contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	// Parse the encoded message
	messageBytes, err := hex.DecodeString(strings.TrimPrefix(req.EncodedMessage, "0x"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid encoded message: %s", err.Error())})
		return
	}
	message, err := protocol.DecodeMessage(messageBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid encoded message: %s", err.Error())})
		return
	}
	sourceChainSelector := uint64(message.SourceChainSelector)
	tokenTransfer := message.TokenTransfer
	if tokenTransfer == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "message does not contain a token transfer"})
		return
	}

	// TODO check that the Token Transfer is for this pool

	switch cfg.Type {
	case config.TokenPoolTypeLockRelease:
		s.lockReleaseTokenPoolExecute(c, cfg, instanceAddress, activeTokenPoolContract, sourceChainSelector, message)
		return
	case config.TokenPoolTypeBurnMint:
		s.burnMintTokenPoolExecute(c, cfg, instanceAddress, activeTokenPoolContract, sourceChainSelector, message)
		return
	default:
		s.logger.Error().Stringer("address", instanceAddress).Msgf("unknown token pool type: %s", cfg.Type)
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}
}

func (s Server) lockReleaseTokenPoolExecute(
	c *gin.Context,
	cfg ContractConfig,
	instanceAddress contracts.InstanceAddress,
	activeTokenPoolContract *apiv2.ActiveContract,
	sourceChainSelector uint64,
	message *protocol.Message,
) {
	lockReleaseTokenPool, err := ParseLockReleaseTokenPool(activeTokenPoolContract.CreatedEvent)
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse lock release token pool contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	remoteChainConfig, ok := lockReleaseTokenPool.RemoteChainConfigs[sourceChainSelector]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("unsupported source chain selector: %v", sourceChainSelector)})
		return
	}

	var rateLimiter *apiv2.ActiveContract
	if message.Finality == protocol.FinalityWaitForFinality {
		rateLimiter, ok = s.activeContractStore.Get(remoteChainConfig.InboundRateLimiter)
		if !ok {
			s.logger.Error().Uint64("sourceChainSelector", sourceChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("inbound active rate limiter not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
	} else {
		rateLimiter, ok = s.activeContractStore.Get(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter)
		if !ok {
			s.logger.Error().Uint64("sourceChainSelector", sourceChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("custom inbound active rate limiter not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.InboundCCVs))
	for i, v := range remoteChainConfig.InboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	holdings, ok := s.instrumentHoldingStore.GetHolding(cfg.Owner, lockReleaseTokenPool.InstrumentId)
	if !ok {
		s.logger.Error().Str("owner", string(cfg.Owner)).Any("instrumentId", lockReleaseTokenPool.InstrumentId).Msg("no holdings found for lock release token pool")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}
	tokenPoolHoldings := make([]splice_api_token_metadata_v1.AnyValue, len(holdings))
	disclosedHoldings := make([]oapiCommon.DisclosedContract, len(holdings))
	for i, holding := range holdings {
		tokenPoolHoldings[i] = splice_api_token_metadata_v1.AnyValue{AVContractId: new(types.CONTRACT_ID(holding.GetCreatedEvent().GetContractId()))}
		disclosedHoldings[i] = converters.ActiveContractToDisclosedContract(holding)
	}

	// The ChoiceContext that will be passed to the Token Pool
	choiceContext := splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{
		string(common.RateLimiterKey):                            {AVContractId: new(types.CONTRACT_ID(rateLimiter.GetCreatedEvent().GetContractId()))},
		string(lockreleasetokenpool.TokenPoolHoldingsContextKey): {AVList: new(tokenPoolHoldings)},
	}}

	// The ChoiceContext that will be passed to the TransferFactory by the Token Pool
	transferFactoryContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: make(map[string]splice_api_token_metadata_v1.AnyValue),
	}
	var factoryDisclosures []oapiCommon.DisclosedContract
	// Get ExtraArgs and TransferFactory from Token Standard API (if enabled)
	if cfg.transferFactory != nil {
		transferFactory, transferContext, disclosedFactoryContracts, err := cfg.transferFactory(c, lockReleaseTokenPool.InstrumentId)
		if err != nil {
			s.logger.Error().Err(err).Msg("transfer factory returned an error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
		choiceContext.Values[string(lockreleasetokenpool.TransferFactoryContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(transferFactory)),
		}
		transferFactoryContext.Values = transferContext.Values
		factoryDisclosures = append(factoryDisclosures, disclosedFactoryContracts...)
	}

	// If the TransferFactory context contains any values, set it as part of the choiceContext
	if len(transferFactoryContext.Values) > 0 {
		choiceContext.Values[string(lockreleasetokenpool.TransferFactoryExtraArgsContextValuesContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVMap: new(transferFactoryContext.Values),
		}
	}

	contextData, err := converters.SerializeChoiceContext(choiceContext)
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiTokenPool.TokenPoolExecuteResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    lockReleaseTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: lockReleaseTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData:        contextData,
		DisclosedContracts: slices.Concat(
			disclosedHoldings,
			factoryDisclosures,
			[]oapiCommon.DisclosedContract{
				converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
				converters.ActiveContractToDisclosedContract(rateLimiter),
			},
		),
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) burnMintTokenPoolExecute(
	c *gin.Context,
	cfg ContractConfig,
	instanceAddress contracts.InstanceAddress,
	activeTokenPoolContract *apiv2.ActiveContract,
	sourceChainSelector uint64,
	message *protocol.Message,
) {
	burnMintTokenPool, err := ParseBurnMintTokenPool(activeTokenPoolContract.CreatedEvent)
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse lock release token pool contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	remoteChainConfig, ok := burnMintTokenPool.RemoteChainConfigs[sourceChainSelector]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("unsupported source chain selector: %v", sourceChainSelector)})
		return
	}

	ccipContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: map[string]splice_api_token_metadata_v1.AnyValue{},
	}
	var rateLimiter *apiv2.ActiveContract
	if message.Finality == protocol.FinalityWaitForFinality {
		rateLimiter, ok = s.activeContractStore.Get(remoteChainConfig.InboundRateLimiter)
		if !ok {
			s.logger.Error().Uint64("sourceChainSelector", sourceChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("inbound active rate limiter not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
		ccipContext.Values[string(common.RateLimiterKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(rateLimiter.GetCreatedEvent().GetContractId())),
		}
	} else {
		rateLimiter, ok = s.activeContractStore.Get(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter)
		if !ok {
			s.logger.Error().Uint64("sourceChainSelector", sourceChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("custom inbound active rate limiter not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
		ccipContext.Values[string(common.RateLimiterKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(rateLimiter.GetCreatedEvent().GetContractId())),
		}
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.InboundCCVs))
	for i, v := range remoteChainConfig.InboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	// The ChoiceContext that will be passed to the Token Pool
	choiceContext := splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{
		string(common.RateLimiterKey): {AVContractId: new(types.CONTRACT_ID(rateLimiter.GetCreatedEvent().GetContractId()))},
	}}

	// The ChoiceContext that will be passed to the BurnMintFactory by the Token Pool
	transferFactoryContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: make(map[string]splice_api_token_metadata_v1.AnyValue),
	}

	// Get BurnMintFactory (if enabled)
	var factoryDisclosures []oapiCommon.DisclosedContract
	if cfg.burnMintFactory != nil {
		burnMintFactory, disclosedFactoryContracts, err := cfg.burnMintFactory(c)
		if err != nil {
			s.logger.Error().Err(err).Msg("transfer factory returned an error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return
		}
		choiceContext.Values[string(burnminttokenpool.BurnMintFactoryContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVContractId: new(types.CONTRACT_ID(burnMintFactory)),
		}
		factoryDisclosures = append(factoryDisclosures, disclosedFactoryContracts...)
	}

	// If the TransferFactory context contains any values, set it as part of the choiceContext
	if len(transferFactoryContext.Values) > 0 {
		choiceContext.Values[string(burnminttokenpool.BurnMintFactoryExtraArgsContextValuesContextKey)] = splice_api_token_metadata_v1.AnyValue{
			AVMap: new(transferFactoryContext.Values),
		}
	}

	contextData, err := converters.SerializeChoiceContext(choiceContext)
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiTokenPool.TokenPoolExecuteResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    burnMintTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: burnMintTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData:        contextData,
		DisclosedContracts: slices.Concat(
			factoryDisclosures,
			[]oapiCommon.DisclosedContract{
				converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
				converters.ActiveContractToDisclosedContract(rateLimiter),
			},
		),
	}

	c.JSON(http.StatusOK, resp)
}
