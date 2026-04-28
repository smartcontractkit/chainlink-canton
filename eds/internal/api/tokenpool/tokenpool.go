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
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
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
	FactoryResolver factoryResolver
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
		switch tokenPool.Type {
		case config.TokenPoolTypeLockRelease:
			s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}),
				PartyID:    tokenPool.PartyID,
			})
		case config.TokenPoolTypeBurnMint:
			s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{}),
				PartyID:    tokenPool.PartyID,
			})
		default:
			return nil, fmt.Errorf("unsupported token pool type: %s", tokenPool.Type)
		}

		if resolver := newStaticFactoryResolver(tokenPool.TransferFactoryID, tokenPool.BurnMintFactoryID); resolver != nil {
			contractConfig.FactoryResolver = resolver
		} else if tokenPool.TokenStandardURL != nil {
			resolver, err := getTokenStandardFactoryResolver(ctx, *tokenPool.TokenStandardURL, tokenPool.TokenStandardAuthConfig, contractConfig.Owner)
			if err != nil {
				return nil, fmt.Errorf("failed to get token standard resolver for token pool with address %s: %w", tokenPool.InstanceAddress, err)
			}
			contractConfig.FactoryResolver = resolver
		}

		s.contractConfigs[tokenPool.InstanceAddress] = contractConfig
	}

	return s, nil
}

func (s Server) PostTokenPoolSend(c *gin.Context, address string) {
	var req oapiTokenPool.TokenPoolSendRequest
	if err := c.ShouldBind(&req); err != nil {
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

	switch cfg.Type {
	case config.TokenPoolTypeLockRelease:
		s.lockReleaseTokenPoolSend(c, cfg, instanceAddress, activeTokenPoolContract, destinationChainSelector)
	case config.TokenPoolTypeBurnMint:
		s.burnMintTokenPoolSend(c, cfg, instanceAddress, activeTokenPoolContract, destinationChainSelector)
	default:
		s.logger.Error().Stringer("address", instanceAddress).Msgf("unknown token pool type: %s", cfg.Type)
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
	}
}

func (s Server) lockReleaseTokenPoolSend(
	c *gin.Context,
	cfg ContractConfig,
	instanceAddress contracts.InstanceAddress,
	activeTokenPoolContract *apiv2.ActiveContract,
	destinationChainSelector uint64,
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
	contextValues := map[string]struct {
		Tag   string `json:"tag"`
		Value string `json:"value"`
	}{
		"rate-limiter": {Tag: "AV_ContractId", Value: rateLimiter.GetCreatedEvent().GetContractId()},
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.OutboundCCVs))
	for i, v := range remoteChainConfig.OutboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	factories, err := s.resolveFactories(c, cfg, lockReleaseTokenPool.InstrumentId)
	if err != nil {
		s.logger.Error().Err(err).Msg("token factory resolver returned an error")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiTokenPool.TokenPoolSendResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    lockReleaseTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: lockReleaseTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData: map[string]any{
			"values": contextValues,
		},
		TokenInput: oapiTokenPool.TokenInput{
			ExtraArgs: struct {
				Context  map[string]any `json:"context"`
				Metadata map[string]any `json:"metadata"`
			}{
				Context:  factories.ChoiceContext,
				Metadata: map[string]any{},
			},
			TokenPoolHoldings: nil,
			TransferFactory:   oapiCommon.ContractId(factories.TransferFactory),
		},
		DisclosedContracts: append(
			factories.DisclosedContracts,
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
) {
	burnMintTokenPool, err := ParseBurnMintTokenPool(activeTokenPoolContract.CreatedEvent)
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse burn/mint token pool contract")
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
	contextValues := map[string]struct {
		Tag   string `json:"tag"`
		Value string `json:"value"`
	}{
		"rate-limiter": {Tag: "AV_ContractId", Value: rateLimiter.GetCreatedEvent().GetContractId()},
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.OutboundCCVs))
	for i, v := range remoteChainConfig.OutboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	factories, err := s.resolveFactories(c, cfg, burnMintTokenPool.InstrumentId)
	if err != nil {
		s.logger.Error().Err(err).Msg("token factory resolver returned an error")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}
	if factories.BurnMintFactory == nil || *factories.BurnMintFactory == "" {
		s.logger.Error().Stringer("address", instanceAddress).Msg("burn/mint token pool is missing burn/mint factory configuration")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "burn/mint factory not configured"})
		return
	}

	resp := &oapiTokenPool.TokenPoolSendResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    burnMintTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: burnMintTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData: map[string]any{
			"values": contextValues,
		},
		TokenInput: oapiTokenPool.TokenInput{
			ExtraArgs: struct {
				Context  map[string]any `json:"context"`
				Metadata map[string]any `json:"metadata"`
			}{
				Context:  factories.ChoiceContext,
				Metadata: map[string]any{},
			},
			TokenPoolHoldings: nil,
			TransferFactory:   oapiCommon.ContractId(factories.TransferFactory),
			BurnMintFactory:   toOAPIContractID(factories.BurnMintFactory),
		},
		DisclosedContracts: append(
			factories.DisclosedContracts,
			converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
			converters.ActiveContractToDisclosedContract(rateLimiter),
		),
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) PostTokenPoolExecute(c *gin.Context, address string) {
	var req oapiTokenPool.TokenPoolExecuteRequest
	if err := c.ShouldBind(&req); err != nil {
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

	switch cfg.Type {
	case config.TokenPoolTypeLockRelease:
		s.lockReleaseTokenPoolExecute(c, cfg, instanceAddress, activeTokenPoolContract, sourceChainSelector, message)
	case config.TokenPoolTypeBurnMint:
		s.burnMintTokenPoolExecute(c, cfg, instanceAddress, activeTokenPoolContract, sourceChainSelector, message)
	default:
		s.logger.Error().Stringer("address", instanceAddress).Msgf("unknown token pool type: %s", cfg.Type)
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
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

	contextValues, rateLimiter, ok := s.getInboundRateLimiterContext(c, instanceAddress, sourceChainSelector, remoteChainConfig, message)
	if !ok {
		return
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
	tokenPoolHoldings := make([]oapiCommon.ContractId, len(holdings))
	disclosedHoldings := make([]oapiCommon.DisclosedContract, len(holdings))
	for i, holding := range holdings {
		tokenPoolHoldings[i] = holding.GetCreatedEvent().GetContractId()
		disclosedHoldings[i] = converters.ActiveContractToDisclosedContract(holding)
	}

	factories, err := s.resolveFactories(c, cfg, lockReleaseTokenPool.InstrumentId)
	if err != nil {
		s.logger.Error().Err(err).Msg("token factory resolver returned an error")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiTokenPool.TokenPoolExecuteResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    lockReleaseTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: lockReleaseTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData: map[string]any{
			"values": contextValues,
		},
		TokenInput: oapiTokenPool.TokenInput{
			TokenPoolHoldings: tokenPoolHoldings,
			ExtraArgs: struct {
				Context  map[string]any `json:"context"`
				Metadata map[string]any `json:"metadata"`
			}{
				Context:  factories.ChoiceContext,
				Metadata: map[string]any{},
			},
			TransferFactory: oapiCommon.ContractId(factories.TransferFactory),
		},
		DisclosedContracts: slices.Concat(
			disclosedHoldings,
			factories.DisclosedContracts,
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
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse burn/mint token pool contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	remoteChainConfig, ok := burnMintTokenPool.RemoteChainConfigs[sourceChainSelector]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("unsupported source chain selector: %v", sourceChainSelector)})
		return
	}

	contextValues, rateLimiter, ok := s.getInboundRateLimiterContext(c, instanceAddress, sourceChainSelector, remoteChainConfig, message)
	if !ok {
		return
	}

	requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.InboundCCVs))
	for i, v := range remoteChainConfig.InboundCCVs {
		requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
	}

	factories, err := s.resolveFactories(c, cfg, burnMintTokenPool.InstrumentId)
	if err != nil {
		s.logger.Error().Err(err).Msg("token factory resolver returned an error")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}
	if factories.BurnMintFactory == nil || *factories.BurnMintFactory == "" {
		s.logger.Error().Stringer("address", instanceAddress).Msg("burn/mint token pool is missing burn/mint factory configuration")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "burn/mint factory not configured"})
		return
	}

	resp := &oapiTokenPool.TokenPoolExecuteResponse{
		ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    burnMintTokenPool.Address.InstanceAddress().Hex(),
		RawInstanceAddress: burnMintTokenPool.Address.String(),
		RequiredCCVs:       requiredCCVs,
		ContextData: map[string]any{
			"values": contextValues,
		},
		TokenInput: oapiTokenPool.TokenInput{
			TokenPoolHoldings: nil,
			ExtraArgs: struct {
				Context  map[string]any `json:"context"`
				Metadata map[string]any `json:"metadata"`
			}{
				Context:  factories.ChoiceContext,
				Metadata: map[string]any{},
			},
			TransferFactory: oapiCommon.ContractId(factories.TransferFactory),
			BurnMintFactory: toOAPIContractID(factories.BurnMintFactory),
		},
		DisclosedContracts: append(
			factories.DisclosedContracts,
			converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
			converters.ActiveContractToDisclosedContract(rateLimiter),
		),
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) resolveFactories(
	ctx context.Context,
	cfg ContractConfig,
	instrumentId splice_api_token_holding_v1.InstrumentId,
) (tokenFactories, error) {
	if cfg.FactoryResolver == nil {
		return tokenFactories{
			ChoiceContext:      map[string]any{"values": map[string]any{}},
			DisclosedContracts: nil,
		}, nil
	}

	factories, err := cfg.FactoryResolver(ctx, instrumentId)
	if err != nil {
		return tokenFactories{}, err
	}
	if factories.ChoiceContext == nil {
		factories.ChoiceContext = map[string]any{"values": map[string]any{}}
	}

	return factories, nil
}

func (s Server) getInboundRateLimiterContext(
	c *gin.Context,
	instanceAddress contracts.InstanceAddress,
	sourceChainSelector uint64,
	remoteChainConfig RemoteChainConfig,
	message *protocol.Message,
) (map[string]struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}, *apiv2.ActiveContract, bool) {
	contextValues := make(map[string]struct {
		Tag   string `json:"tag"`
		Value string `json:"value"`
	})
	var rateLimiter *apiv2.ActiveContract
	var ok bool
	if message.Finality == protocol.FinalityWaitForFinality {
		rateLimiter, ok = s.activeContractStore.Get(remoteChainConfig.InboundRateLimiter)
		if !ok {
			s.logger.Error().Uint64("sourceChainSelector", sourceChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("inbound active rate limiter not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return nil, nil, false
		}
		contextValues["inbound-rate-limiter"] = struct {
			Tag   string `json:"tag"`
			Value string `json:"value"`
		}{Tag: "AV_ContractId", Value: rateLimiter.GetCreatedEvent().GetContractId()}
	} else {
		rateLimiter, ok = s.activeContractStore.Get(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter)
		if !ok {
			s.logger.Error().Uint64("sourceChainSelector", sourceChainSelector).Stringer("poolAddress", instanceAddress).Stringer("rateLimiterAddress", remoteChainConfig.OutboundRateLimiter).Msg("custom inbound active rate limiter not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
			return nil, nil, false
		}
		contextValues["inbound-custom-block-confirmations-rate-limiter"] = struct {
			Tag   string `json:"tag"`
			Value string `json:"value"`
		}{Tag: "AV_ContractId", Value: rateLimiter.GetCreatedEvent().GetContractId()}
	}

	return contextValues, rateLimiter, true
}

func toOAPIContractID(contractID *string) *oapiCommon.ContractId {
	if contractID == nil || *contractID == "" {
		return nil
	}

	value := oapiCommon.ContractId(*contractID)
	return &value
}
