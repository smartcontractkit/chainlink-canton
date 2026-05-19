package executor

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/global"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
)

type ContractConfig struct{}

type Server struct {
	logger              zerolog.Logger
	activeContractStore store.ActiveContractStoreInterface

	contractConfigs map[contracts.InstanceAddress]ContractConfig
}

var _ oapiExecutor.ServerInterface = &Server{}

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore store.ActiveContractStoreInterface,
	config config.ExecutorAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:              logger.With().Str("component", "ExecutorAPI").Logger(),
		activeContractStore: activeContractStore,
		contractConfigs:     make(map[contracts.InstanceAddress]ContractConfig),
	}
	for _, executor := range config.Executors {
		s.contractConfigs[executor.InstanceAddress] = ContractConfig{}
		s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(executorBinding.Executor{}),
			PartyID:    executor.PartyID,
		})
	}

	return s, nil
}

func (s *Server) PostExecutorSend(c *gin.Context, address string) {
	var req oapiExecutor.ExecutorSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, oapiCommon.ErrorResponse{Error: "request body too large"})
			return
		}
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
		c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "executor address not found"})
		return
	}

	activeExecutorContract, ok := s.activeContractStore.Get(instanceAddress)
	if !ok {
		s.logger.Error().Stringer("address", instanceAddress).Msg("active executor contract not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiCommon.ErrorResponse{Error: "internal server error"})

		return
	}

	parsedExecutor, err := ParseExecutor(activeExecutorContract.GetCreatedEvent())
	if err != nil {
		s.logger.Err(err).Stringer("address", instanceAddress).Msg("failed to parse executor contract")
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

	// TODO check the provided CCVs against CCV allow list
	// TODO Validate that the executor specified in the message matches this executor (and the user didn't specify no-execution)

	resp := oapiExecutor.ExecutorSendResponse{
		ContractId:         activeExecutorContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedExecutor.Address.InstanceAddress().Hex(),
		RawInstanceAddress: parsedExecutor.Address.String(),
		ContextData:        contextData,
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeExecutorContract),
		},
	}

	c.JSON(http.StatusOK, resp)
}

var _ global.InstanceAddressFilter = &Server{}

// FilterContracts returns the sub-set of contracts that are tracked by the Executor API Server
func (s *Server) FilterContracts(addresses []contracts.InstanceAddress) []contracts.InstanceAddress {
	var out []contracts.InstanceAddress
	for _, address := range addresses {
		if _, ok := s.contractConfigs[address]; ok {
			out = append(out, address)
		}
	}

	return out
}
