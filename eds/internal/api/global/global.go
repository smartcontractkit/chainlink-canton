package global

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"

	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
)

type Server struct {
	logger              zerolog.Logger //nolint:unused
	activeContractStore *store.ActiveContractStore

	maxBatchLimit int
}

var _ oapiGlobal.ServerInterface = &Server{}

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
