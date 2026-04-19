package ccv

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"

	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

type ContractConfig struct{}

type Server struct {
	logger              zerolog.Logger
	activeContractStore store.ActiveContractStore

	contractConfigs map[contracts.InstanceAddress]ContractConfig
}

var _ oapiCCV.ServerInterface = &Server{}

func NewServer(
	logger zerolog.Logger,
	activeContractStore store.ActiveContractStore,
	config config.CCVAPIConfig,
) *Server {
	s := &Server{
		logger:              logger,
		activeContractStore: activeContractStore,
		contractConfigs:     make(map[contracts.InstanceAddress]ContractConfig),
	}
	for _, ccv := range config.CCVs {
		s.contractConfigs[ccv.InstanceAddress] = ContractConfig{}
	}

	return s
}

func (s Server) PostCCVSend(c *gin.Context, address string) {
	var req oapiCCV.CCVSendRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	instanceAddress, err := converters.ResolveAddress(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	_, ok := s.contractConfigs[instanceAddress]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "CCV address not found"})
		return
	}

	activeCCVContract, ok := s.activeContractStore.Get(instanceAddress)
	if !ok {
		s.logger.Error().Stringer("address", instanceAddress).Msg("active CCV contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	parsedCommitteeVerifierContract, err := ParseCommitteeVerifier(activeCCVContract.GetCreatedEvent())
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse committee verifier contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	// TODO check if destination chain is configured on the CCV

	resp := oapiCCV.CCVSendResponse{
		ContractId:         activeCCVContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedCommitteeVerifierContract.Address.InstanceAddress().Hex(),
		RawInstanceAddress: parsedCommitteeVerifierContract.Address.String(),
		ContextData: map[string]any{
			"values": map[string]struct {
				Tag   string `json:"tag"`
				Value string `json:"value"`
			}{},
		},
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeCCVContract),
		},
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) PostCCVExecute(c *gin.Context, address string) {
	s.logger.Info().Msg("Got request for CCVExecute")
	var req oapiCCV.CCVExecuteRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	instanceAddress, err := converters.ResolveAddress(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	_, ok := s.contractConfigs[instanceAddress]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "CCV address not found"})
		return
	}

	activeCCVContract, ok := s.activeContractStore.Get(instanceAddress)
	if !ok {
		s.logger.Error().Stringer("address", instanceAddress).Msg("active CCV contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	parsedCommitteeVerifierContract, err := ParseCommitteeVerifier(activeCCVContract.GetCreatedEvent())
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse committee verifier contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	resp := oapiCCV.CCVExecuteResponse{
		ContractId:         activeCCVContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedCommitteeVerifierContract.Address.InstanceAddress().Hex(),
		RawInstanceAddress: parsedCommitteeVerifierContract.Address.String(),
		ContextData: map[string]any{
			"values": map[string]struct {
				Tag   string `json:"tag"`
				Value string `json:"value"`
			}{},
		},
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeCCVContract),
		},
	}

	c.JSON(http.StatusOK, resp)
}
