package global

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
)

type Server struct {
	logger              zerolog.Logger
	activeContractStore *store.ActiveContractStore
	addressFilters      []InstanceAddressFilter

	maxBatchLimit int
}

var _ oapiGlobal.ServerInterface = &Server{}

type InstanceAddressFilter interface {
	// FilterContracts returns the sub-set of contracts that are tracked by the InstanceAddressFilter
	FilterContracts(addresses []contracts.InstanceAddress) []contracts.InstanceAddress
}

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	config config.GlobalAPIConfig,
	addressFilters ...InstanceAddressFilter,
) (*Server, error) {
	s := &Server{
		logger:              logger.With().Str("component", "GlobalAPI").Logger(),
		activeContractStore: activeContractStore,
		addressFilters:      addressFilters,
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
		return
	}

	if len(req.Addresses) > s.maxBatchLimit {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: "max batch limit exceeded"})
		return
	}

	// Parse all requested addresses to InstanceAddresses
	requestedAddresses := make([]contracts.InstanceAddress, len(req.Addresses))
	for i, address := range req.Addresses {
		instanceAddress, err := converters.ResolveRawOrHashedAddress(address)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiCommon.ErrorResponse{Error: err.Error()})
			return
		}
		requestedAddresses[i] = instanceAddress
	}

	// Query all filters to determine which of the requested addresses should be returned.
	// This ensures that if an address is requested which is not known by any of the filters, an error will be returned.
	filteredAddresses := make(map[contracts.InstanceAddress]bool, len(requestedAddresses))
	for _, filter := range s.addressFilters {
		for _, address := range filter.FilterContracts(requestedAddresses) {
			filteredAddresses[address] = true
		}
	}

	resp := oapiGlobal.GetExplicitDisclosureBatchResponse{
		Disclosures: make([]oapiCommon.DisclosedContract, len(req.Addresses)),
	}
	for i, instanceAddress := range requestedAddresses {
		// Validate that the address should be returned, if not return an error
		if ok := filteredAddresses[instanceAddress]; !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: fmt.Sprintf("contract not found: %s", instanceAddress.Hex())})
			return
		}

		// Get the active contract for this address.
		// If the active contract cannot be found, something is misconfigured.
		activeContract, ok := s.activeContractStore.Get(instanceAddress)
		if !ok {
			s.logger.Error().Stringer("address", instanceAddress).Msg("active contract not found for address that passed filters")
			c.AbortWithStatusJSON(http.StatusNotFound, oapiCommon.ErrorResponse{Error: "active contract not found"})
			return
		}

		resp.Disclosures[i] = converters.ActiveContractToDisclosedContract(activeContract)
	}

	c.JSON(http.StatusOK, resp)
}
