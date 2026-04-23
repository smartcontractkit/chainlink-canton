package global

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
)

type Server struct {
	logger              zerolog.Logger
	activeContractStore *store.ActiveContractStore

	maxBatchLimit int
}

var _ oapiGlobal.ServerInterface = &Server{}

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	config config.GlobalAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:              logger.With().Str("component", "GlobalAPI").Logger(),
		activeContractStore: activeContractStore,
		maxBatchLimit:       config.MaxBatchSize,
	}

	if config.MaxBatchSize <= 0 {
		return nil, fmt.Errorf("MaxBatchSize must be greater than zero")
	}

	return s, nil
}

// PostGetExplicitDisclosureBatch (POST /ccip/v1/global/disclosure/batch)
func (s Server) PostGetExplicitDisclosureBatch(c *gin.Context) {
	var req oapiGlobal.GetExplicitDisclosureBatchRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	if len(req.Addresses) > s.maxBatchLimit {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "max batch limit exceeded"})
		return
	}

	resp := oapiGlobal.GetExplicitDisclosureBatchResponse{
		Disclosures: make([]oapiCommon.DisclosedContract, len(req.Addresses)),
	}
	for i, address := range req.Addresses {
		instanceAddress, err := converters.ResolveRawOrHashedAddress(address)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
			return
		}

		activeContract, ok := s.activeContractStore.Get(instanceAddress)
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "active contract not found"})
			return
		}

		resp.Disclosures[i] = converters.ActiveContractToDisclosedContract(activeContract)
	}

	c.JSON(http.StatusOK, resp)
}
