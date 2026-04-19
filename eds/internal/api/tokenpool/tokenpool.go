package tokenpool

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-ccv/protocol"

	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"

	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
)

type ContractConfig struct {
	Type config.TokenPoolType
}

type Server struct {
	logger                 zerolog.Logger
	activeContractStore    store.ActiveContractStore
	instrumentHoldingStore store.InstrumentHoldingStore

	contractConfigs map[contracts.InstanceAddress]ContractConfig
}

var _ oapiTokenPool.ServerInterface = &Server{}

func NewServer(
	logger zerolog.Logger,
	activeContractStore store.ActiveContractStore,
	instrumentHoldingStore store.InstrumentHoldingStore,
	config config.TokenPoolAPIConfig,

) *Server {
	s := &Server{
		logger:                 logger,
		activeContractStore:    activeContractStore,
		instrumentHoldingStore: instrumentHoldingStore,
		contractConfigs:        make(map[contracts.InstanceAddress]ContractConfig),
	}
	for _, tokenPool := range config.TokenPools {
		s.contractConfigs[tokenPool.InstanceAddress] = ContractConfig{
			Type: tokenPool.Type,
		}
	}

	return s
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
	}

	var resp *oapiTokenPool.TokenPoolSendResponse
	switch cfg.Type {
	case config.TokenPoolTypeLockRelease:
		lockReleaseTokenPool, err := ParseLockReleaseTokenPool(activeTokenPoolContract.CreatedEvent)
		if err != nil {
			s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse lock release token pool contract")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

			return
		}

		// TODO validate that the TokenTransfer in the message is actually for this token pool

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

		resp = &oapiTokenPool.TokenPoolSendResponse{
			ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
			InstanceAddress:    lockReleaseTokenPool.Address.InstanceAddress().Hex(),
			RawInstanceAddress: lockReleaseTokenPool.Address.String(),
			RequiredCCVs:       requiredCCVs,
			ContextData:        nil,                        // TODO
			TokenInput:         oapiTokenPool.TokenInput{}, // TODO
			DisclosedContracts: []oapiCommon.DisclosedContract{
				converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
				converters.ActiveContractToDisclosedContract(rateLimiter),
			},
		}
	case config.TokenPoolTypeBurnMint:
		fallthrough
	default:
		s.logger.Error().Stringer("address", instanceAddress).Msgf("unknown token pool type: %s", cfg.Type)
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
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
	}

	var resp *oapiTokenPool.TokenPoolExecuteResponse
	switch cfg.Type {
	case config.TokenPoolTypeLockRelease:
		lockReleaseTokenPool, err := ParseLockReleaseTokenPool(activeTokenPoolContract.CreatedEvent)
		if err != nil {
			s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse lock release token pool contract")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

			return
		}

		// TODO validate that the TokenTransfer in the message is actually for this token pool

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

		requiredCCVs := make([]oapiCommon.RawOrHashedAddress, len(remoteChainConfig.OutboundCCVs))
		for i, v := range remoteChainConfig.OutboundCCVs {
			requiredCCVs[i] = converters.RawInstanceAddressAsRawOrHashedAddress(v)
		}

		resp = &oapiTokenPool.TokenPoolExecuteResponse{
			ContractId:         activeTokenPoolContract.GetCreatedEvent().GetContractId(),
			InstanceAddress:    lockReleaseTokenPool.Address.InstanceAddress().Hex(),
			RawInstanceAddress: lockReleaseTokenPool.Address.String(),
			RequiredCCVs:       requiredCCVs,
			ContextData:        nil,                        // TODO
			TokenInput:         oapiTokenPool.TokenInput{}, // TODO
			DisclosedContracts: []oapiCommon.DisclosedContract{
				converters.ActiveContractToDisclosedContract(activeTokenPoolContract),
				converters.ActiveContractToDisclosedContract(rateLimiter),
			},
		}
	case config.TokenPoolTypeBurnMint:
		fallthrough
	default:
		s.logger.Error().Stringer("address", instanceAddress).Msgf("unknown token pool type: %s", cfg.Type)
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	c.JSON(http.StatusOK, resp)
}
