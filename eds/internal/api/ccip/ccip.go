package ccip

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccv/protocol"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

type Server struct {
	logger              zerolog.Logger
	activeContractStore *store.ActiveContractStore

	perPartyRouterFactory contracts.InstanceAddress
	onRamp                contracts.InstanceAddress
	offRamp               contracts.InstanceAddress
	globalConfig          contracts.InstanceAddress
	tokenAdminRegistry    contracts.InstanceAddress
	rmnRemote             contracts.InstanceAddress
	feeQuoter             contracts.InstanceAddress
}

var _ oapiCCIP.ServerInterface = &Server{}

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	config config.CCIPAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:                logger.With().Str("component", "CCIPAPI").Logger(),
		activeContractStore:   activeContractStore,
		perPartyRouterFactory: config.PerPartyRouterFactory.InstanceAddress,
		onRamp:                config.OnRamp.InstanceAddress,
		offRamp:               config.OffRamp.InstanceAddress,
		globalConfig:          config.GlobalConfig.InstanceAddress,
		tokenAdminRegistry:    config.TokenAdminRegistry.InstanceAddress,
		rmnRemote:             config.RMNRemote.InstanceAddress,
		feeQuoter:             config.FeeQuoter.InstanceAddress,
	}

	s.activeContractStore.RegisterTemplates(
		[]store.RegisteredTemplate{
			{
				TemplateID: contracts.TemplateIDFromBinding(perpartyrouter.PerPartyRouterFactory{}),
				PartyID:    config.OnRamp.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(onramp.OnRamp{}),
				PartyID:    config.OnRamp.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(offramp.OffRamp{}),
				PartyID:    config.OffRamp.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(common.GlobalConfig{}),
				PartyID:    config.GlobalConfig.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(tokenadminregistry.TokenAdminRegistry{}),
				PartyID:    config.TokenAdminRegistry.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(rmn.RMNRemote{}),
				PartyID:    config.RMNRemote.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(feequoter.FeeQuoter{}),
				PartyID:    config.FeeQuoter.PartyID,
			},
		}...,
	)

	return s, nil
}

func (s Server) PostPerPartyRouterFactory(c *gin.Context) {
	var req oapiCCIP.CCIPPerPartyRouterFactoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	activePerPartyRouterFactoryContract, ok := s.activeContractStore.Get(s.perPartyRouterFactory)
	if !ok {
		s.logger.Error().Stringer("address", s.perPartyRouterFactory).Msg("active per-party router factory contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	parsedPerPartyRouterFactory, err := ParsePerPartyRouterFactory(activePerPartyRouterFactoryContract.GetCreatedEvent())
	if err != nil {
		s.logger.Err(err).Stringer("address", s.perPartyRouterFactory).Msg("failed to parse per-party router factory contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := &oapiCCIP.CCIPPerPartyRouterFactoryResponse{
		ContractId:         activePerPartyRouterFactoryContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedPerPartyRouterFactory.Address.InstanceID(),
		RawInstanceAddress: parsedPerPartyRouterFactory.Address.String(),
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activePerPartyRouterFactoryContract),
		},
	}
	c.JSON(http.StatusOK, resp)
}

func (s Server) GetTokenAdminRegistryToken(c *gin.Context, instrumentId oapiCommon.HashedInstrumentId) {
	encodedInstrumentId := contracts.HexToEncodedInstrumentID(instrumentId)

	if (encodedInstrumentId == contracts.EncodedInstrumentID{}) {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "invalid instrumentId"})
		return
	}

	activeTokenAdminRegistryContract, ok := s.activeContractStore.Get(s.tokenAdminRegistry)
	if !ok {
		s.logger.Error().Stringer("address", s.tokenAdminRegistry).Msg("active tokenAdminRegistry contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	parsedTokenAdminRegistry, err := ParseTokenAdminRegistry(activeTokenAdminRegistryContract.GetCreatedEvent())
	if err != nil {
		s.logger.Err(err).Stringer("address", s.tokenAdminRegistry).Msg("failed to parse token admin registry contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	tokenConfig, ok := parsedTokenAdminRegistry.TokenConfigs[encodedInstrumentId]
	if !ok || tokenConfig.TokenPool == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("no token pool registered for token: %s", encodedInstrumentId.Hex())})
		return
	}
	tokenPoolRawInstanceAddress := oapiCommon.RawInstanceAddress(contracts.InstanceID(tokenConfig.TokenPool.PoolInstanceId).RawInstanceAddress(tokenConfig.TokenPool.PoolOwner))

	resp := &oapiCCIP.LookupTokenPoolResponse{
		RawInstanceAddress: tokenPoolRawInstanceAddress,
	}
	c.JSON(http.StatusOK, resp)
}

func (s Server) PostCCIPSend(c *gin.Context) {
	var req oapiCCIP.CCIPSendRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	destinationChainSelector, err := strconv.ParseUint(req.Message.DestinationChainSelector, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "invalid destination chain selector"})
		return
	}
	payload, err := hex.DecodeString(strings.TrimPrefix(req.Message.Payload, "0x"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "invalid payload"})
		return
	}

	activeOnRampContract, ok := s.activeContractStore.Get(s.onRamp)
	if !ok {
		s.logger.Error().Stringer("address", s.onRamp).Msg("active onRamp contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeGlobalConfigContract, ok := s.activeContractStore.Get(s.globalConfig)
	if !ok {
		s.logger.Error().Stringer("address", s.globalConfig).Msg("active globalConfig contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeTokenAdminRegistryContract, ok := s.activeContractStore.Get(s.tokenAdminRegistry)
	if !ok {
		s.logger.Error().Stringer("address", s.tokenAdminRegistry).Msg("active tokenAdminRegistry contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeRMNRemoteContract, ok := s.activeContractStore.Get(s.rmnRemote)
	if !ok {
		s.logger.Error().Stringer("address", s.rmnRemote).Msg("active rmnRemote contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeFeeQuoterContract, ok := s.activeContractStore.Get(s.feeQuoter)
	if !ok {
		s.logger.Error().Stringer("address", s.feeQuoter).Msg("active feeQuoter contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	// Determine CCVs
	parsedGlobalConfig, err := ParseGlobalConfig(activeGlobalConfigContract.GetCreatedEvent())
	if err != nil {
		s.logger.Err(err).Msg("failed to parse global config contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	destChainConfig, ok := parsedGlobalConfig.DestChainConfigs[destinationChainSelector]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "unsupported destination chain selector"})
		return
	}
	// Lane mandated CCVs are always required
	resolvedCCVs := make(map[contracts.InstanceAddress]oapiCommon.RawOrHashedAddress, len(destChainConfig.LaneMandatedCCVs))
	for _, laneMandatedCCV := range destChainConfig.LaneMandatedCCVs {
		resolvedCCVs[laneMandatedCCV.InstanceAddress()] = converters.RawInstanceAddressAsRawOrHashedAddress(laneMandatedCCV)
	}
	addDefaults := false
	// If the message contains a payload, add the sender-required CCVs
	if len(payload) > 0 {
		if req.SenderRequiredCCVs == nil || len(*req.SenderRequiredCCVs) == 0 {
			// If no CCVs are specified, use defaults
			addDefaults = true
		} else {
			for _, address := range *req.SenderRequiredCCVs {
				if converters.RawOrHashedAddressAsString(address) == "default-ccvs" {
					addDefaults = true
					continue
				}
				instanceAddress, err := converters.ResolveRawOrHashedAddress(address)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid sender required CCV address: %v", err)})
					return
				}
				resolvedCCVs[instanceAddress] = address
			}
		}
	}
	// If the message contains a token transfer, add the pool-required CCVs
	if req.Message.TokenTransfer != nil {
		if req.TokenPoolRequiredCCVs == nil || len(*req.TokenPoolRequiredCCVs) == 0 {
			// If no CCVs are specified, use defaults
			addDefaults = true
		} else {
			for _, address := range *req.TokenPoolRequiredCCVs {
				if converters.RawOrHashedAddressAsString(address) == "default-ccvs" {
					addDefaults = true
					continue
				}
				instanceAddress, err := converters.ResolveRawOrHashedAddress(address)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid token pool required CCV address: %v", err)})
					return
				}
				resolvedCCVs[instanceAddress] = address
			}
		}
	}
	if addDefaults {
		for _, defaultCCV := range destChainConfig.DefaultCCVs {
			resolvedCCVs[defaultCCV.InstanceAddress()] = converters.RawInstanceAddressAsRawOrHashedAddress(defaultCCV)
		}
	}

	// Determine Executor
	var executor *oapiCommon.RawOrHashedAddress
	switch req.Message.Executor.Type {
	case oapiCommon.NoExecutor:
		// No Executor
	case oapiCommon.Empty:
		// Default Executor
		// If set to none, the default executor is the no-exec executor
		if destChainConfig.DefaultExecutor != nil {
			defaultExecutor := converters.RawInstanceAddressAsRawOrHashedAddress(*destChainConfig.DefaultExecutor)
			executor = &defaultExecutor
		}
	case oapiCommon.WithAddress:
		// Use the specified executor
		executor = req.Message.Executor.Address
	}

	resp := &oapiCCIP.CCIPSendResponse{
		Ccvs: maps.Values(resolvedCCVs),
		ContextData: map[string]any{
			"values": map[string]struct {
				Tag   string `json:"tag"`
				Value string `json:"value"`
			}{
				"on-ramp": {
					Tag:   "AV_ContractId",
					Value: activeOnRampContract.GetCreatedEvent().GetContractId(),
				},
				"global-config": {
					Tag:   "AV_ContractId",
					Value: activeGlobalConfigContract.GetCreatedEvent().GetContractId(),
				},
				"token-admin-registry": {
					Tag:   "AV_ContractId",
					Value: activeTokenAdminRegistryContract.GetCreatedEvent().GetContractId(),
				},
				"rmn-remote": {
					Tag:   "AV_ContractId",
					Value: activeRMNRemoteContract.GetCreatedEvent().GetContractId(),
				},
				"fee-quoter": {
					Tag:   "AV_ContractId",
					Value: activeFeeQuoterContract.GetCreatedEvent().GetContractId(),
				},
			},
		},
		Executor: executor,
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeOnRampContract),
			converters.ActiveContractToDisclosedContract(activeGlobalConfigContract),
			converters.ActiveContractToDisclosedContract(activeTokenAdminRegistryContract),
			converters.ActiveContractToDisclosedContract(activeRMNRemoteContract),
			converters.ActiveContractToDisclosedContract(activeFeeQuoterContract),
		},
	}
	c.JSON(http.StatusOK, resp)
}

func (s Server) PostCCIPExecute(c *gin.Context) {
	var req oapiCCIP.CCIPExecuteRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
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

	// Get contracts
	activeOffRampContract, ok := s.activeContractStore.Get(s.offRamp)
	if !ok {
		s.logger.Error().Stringer("address", s.offRamp).Msg("active offRamp contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeGlobalConfigContract, ok := s.activeContractStore.Get(s.globalConfig)
	if !ok {
		s.logger.Error().Stringer("address", s.globalConfig).Msg("active globalConfig contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeTokenAdminRegistryContract, ok := s.activeContractStore.Get(s.tokenAdminRegistry)
	if !ok {
		s.logger.Error().Stringer("address", s.tokenAdminRegistry).Msg("active tokenAdminRegistry contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}
	activeRMNRemoteContract, ok := s.activeContractStore.Get(s.rmnRemote)
	if !ok {
		s.logger.Error().Stringer("address", s.rmnRemote).Msg("active rmnRemote contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	var tokenPool *oapiCommon.RawInstanceAddress
	if message.TokenTransfer != nil {
		destTokenAddress := contracts.BytesToEncodedInstrumentID(message.TokenTransfer.DestTokenAddress)

		parsedTokenAdminRegistry, err := ParseTokenAdminRegistry(activeTokenAdminRegistryContract.GetCreatedEvent())
		if err != nil {
			s.logger.Err(err).Msg("failed to parse token admin registry contract")
			c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

			return
		}

		tokenConfig, ok := parsedTokenAdminRegistry.TokenConfigs[destTokenAddress]
		if !ok || tokenConfig.TokenPool == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("no token pool registered for token: %s", destTokenAddress.Hex())})
			return
		}
		tokenPoolRawInstanceAddress := oapiCommon.RawInstanceAddress(contracts.InstanceID(tokenConfig.TokenPool.PoolInstanceId).RawInstanceAddress(tokenConfig.TokenPool.PoolOwner))
		tokenPool = &tokenPoolRawInstanceAddress
	}

	resp := &oapiCCIP.CCIPExecuteResponse{
		ContextData: map[string]any{
			"values": map[string]struct {
				Tag   string `json:"tag"`
				Value string `json:"value"`
			}{
				"off-ramp": {
					Tag:   "AV_ContractId",
					Value: activeOffRampContract.GetCreatedEvent().GetContractId(),
				},
				"global-config": {
					Tag:   "AV_ContractId",
					Value: activeGlobalConfigContract.GetCreatedEvent().GetContractId(),
				},
				"token-admin-registry": {
					Tag:   "AV_ContractId",
					Value: activeTokenAdminRegistryContract.GetCreatedEvent().GetContractId(),
				},
				"rmn-remote": {
					Tag:   "AV_ContractId",
					Value: activeRMNRemoteContract.GetCreatedEvent().GetContractId(),
				},
			},
		},
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeOffRampContract),
			converters.ActiveContractToDisclosedContract(activeGlobalConfigContract),
			converters.ActiveContractToDisclosedContract(activeTokenAdminRegistryContract),
			converters.ActiveContractToDisclosedContract(activeRMNRemoteContract),
		},
		TokenPool: tokenPool,
	}
	c.JSON(http.StatusOK, resp)
}
