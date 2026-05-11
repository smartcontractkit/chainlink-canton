package ccv

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/global"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

type ContractConfig struct{}

type Server struct {
	logger              zerolog.Logger
	activeContractStore *store.ActiveContractStore

	contractConfigs map[contracts.InstanceAddress]ContractConfig
}

var _ oapiCCV.ServerInterface = &Server{}

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	config config.CCVAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:              logger.With().Str("component", "CCVAPI").Logger(),
		activeContractStore: activeContractStore,
		contractConfigs:     make(map[contracts.InstanceAddress]ContractConfig),
	}
	for _, ccv := range config.CCVs {
		s.contractConfigs[ccv.InstanceAddress] = ContractConfig{}
		s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(ccvs.CommitteeVerifier{}),
			PartyID:    ccv.PartyID,
		})
	}

	return s, nil
}

func (s Server) PostCCVSend(c *gin.Context, address string) {
	var req oapiCCV.CCVSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	contextData, err := converters.SerializeChoiceContext(splice_api_token_metadata_v1.ChoiceContext{
		Values: map[string]splice_api_token_metadata_v1.AnyValue{
			// Empty for now
		},
	})
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := oapiCCV.CCVSendResponse{
		ContractId:         activeCCVContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedCommitteeVerifierContract.Address.InstanceAddress().Hex(),
		RawInstanceAddress: parsedCommitteeVerifierContract.Address.String(),
		ContextData:        contextData,
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeCCVContract),
		},
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) PostCCVExecute(c *gin.Context, address string) {
	var req oapiCCV.CCVExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	contextData, err := converters.SerializeChoiceContext(splice_api_token_metadata_v1.ChoiceContext{
		Values: map[string]splice_api_token_metadata_v1.AnyValue{
			// Empty for now
		},
	})
	if err != nil {
		s.logger.Err(err).Msg("failed to serialize CCIP context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := oapiCCV.CCVExecuteResponse{
		ContractId:         activeCCVContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedCommitteeVerifierContract.Address.InstanceAddress().Hex(),
		RawInstanceAddress: parsedCommitteeVerifierContract.Address.String(),
		ContextData:        contextData,
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeCCVContract),
		},
	}

	c.JSON(http.StatusOK, resp)
}

var _ global.InstanceAddressFilter = &Server{}

// FilterContracts returns the sub-set of contracts that are tracked by the CCV API Server
func (s Server) FilterContracts(addresses []contracts.InstanceAddress) []contracts.InstanceAddress {
	var out []contracts.InstanceAddress
	for _, address := range addresses {
		if _, ok := s.contractConfigs[address]; ok {
			out = append(out, address)
		}
	}

	return out
}
