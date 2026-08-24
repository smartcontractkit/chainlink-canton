package executor

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	executorBinding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/global"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/parse"
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
	// Parse and validate request
	var req oapiExecutor.ExecutorSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, oapiCommon.ErrorResponse{Error: "request body too large"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})

		return
	}
	if err := parse.ValidateMessage(req.Message); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid message: %s", err.Error())})
		return
	}

	instanceAddress, err := converters.ResolveAddress(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}
	// Validate that the address path parameter matches the executor address specified in the message
	switch req.Message.Executor.Type {
	case oapiCommon.NoExecutor:
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "executor type is no_executor, can't request executor disclosures"})
		return
	case oapiCommon.Empty:
		// If the message contains an empty executor, that means the default should apply.
		// Cannot check against this executor since we don't know if we're the default.
	case oapiCommon.WithAddress:
		// parse.ValidateMessage will have already verified that the address is present and valid
		requestInstanceAddress, err := converters.ResolveRawOrHashedAddress(*req.Message.Executor.Address)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid executor address in message: %s", err.Error())})
			return
		}
		if instanceAddress != requestInstanceAddress {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("executor in message doesn't match requested executor: %s!=%s", requestInstanceAddress, instanceAddress)})
			return
		}
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

	// Validate request's CCVs against the contract settings
	if int64(len(req.Ccvs)) > parsedExecutor.MaxCCVsPerMessage {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("too many CCVs: %d provided, but max allowed is %d", len(req.Ccvs), parsedExecutor.MaxCCVsPerMessage)})
		return
	}
	// If the allowlist is enabled, validate that the provided CCVs are all allowed
	if parsedExecutor.CCVAllowlistEnabled {
		for i, ccv := range req.Ccvs {
			instanceAddress, err := converters.ResolveRawOrHashedAddress(ccv)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("invalid CCV at index %d: %s", i, err.Error())})
				return
			}
			allowed := false
			for _, allowedCCV := range parsedExecutor.AllowedCCVs {
				if instanceAddress == allowedCCV.InstanceAddress() {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: fmt.Sprintf("CCV at index %d is not in the allowed CCV list", i)})
				return
			}
		}
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
