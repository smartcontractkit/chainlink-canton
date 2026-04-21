package executor

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"

	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
)

type ContractConfig struct{}

type Server struct {
	logger              zerolog.Logger
	activeContractStore *store.ActiveContractStore

	contractConfigs map[contracts.InstanceAddress]ContractConfig
}

var _ oapiExecutor.ServerInterface = &Server{}

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	config config.ExecutorAPIConfig,
) *Server {
	s := &Server{
		logger:              logger,
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

	return s
}

func (s *Server) PostExecutorSend(c *gin.Context, address string) {
	var req oapiExecutor.ExecutorSendRequest
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

	// TODO check the provided CCVs against CCV allow list
	// TODO Validate that the executor specified in the message matches this executor (and the user didn't specify no-execution)

	resp := oapiExecutor.ExecutorSendResponse{
		ContractId:         activeExecutorContract.GetCreatedEvent().GetContractId(),
		InstanceAddress:    parsedExecutor.Address.InstanceAddress().Hex(),
		RawInstanceAddress: parsedExecutor.Address.String(),
		ContextData: map[string]any{
			"values": map[string]struct {
				Tag   string `json:"tag"`
				Value string `json:"value"`
			}{},
		},
		DisclosedContracts: []oapiCommon.DisclosedContract{
			converters.ActiveContractToDisclosedContract(activeExecutorContract),
		},
	}

	c.JSON(http.StatusOK, resp)
}
