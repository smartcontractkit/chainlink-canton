package api

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-ccv/protocol"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/disclosure"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
)

type Server struct {
	logger        zerolog.Logger
	disclosureSvc *disclosure.DisclosureService
}

var _ edsv1.ServerInterface = &Server{}

func NewServer(logger zerolog.Logger, disclosureSvc *disclosure.DisclosureService) *Server {
	return &Server{
		logger:        logger,
		disclosureSvc: disclosureSvc,
	}
}

// (POST /ccip/v1/message/execute)
func (s Server) CcipExecute(c *gin.Context) {
	var req edsv1.CCIPExecuteRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, edsv1.ErrorResponse{Error: err.Error()})
		return
	}

	// Parse requested CCVs
	ccvs := make([]contracts.InstanceAddress, len(req.Ccvs))
	for i, ccv := range req.Ccvs {
		instanceAddress := contracts.HexToInstanceAddress(ccv)
		if (instanceAddress == contracts.InstanceAddress{}) {
			c.JSON(400, edsv1.ErrorResponse{Error: fmt.Sprintf("invalid CCV address: %s", ccv)})
			return
		}
		ccvs[i] = instanceAddress
	}

	// Decode the CCIP encoded message
	var message *protocol.Message
	if req.EncodedMessage != "" {
		messageBytes, err := hex.DecodeString(req.EncodedMessage)
		if err != nil {
			c.JSON(http.StatusBadRequest, edsv1.ErrorResponse{Error: fmt.Sprintf("invalid encoded message: %s", err.Error())})
			return
		}
		message, err = protocol.DecodeMessage(messageBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, edsv1.ErrorResponse{Error: fmt.Sprintf("invalid encoded message: %s", err.Error())})
			return
		}
	}

	disclosures, err := s.disclosureSvc.GetCCIPExecuteDisclosures(c.Request.Context(), disclosure.CCIPExecuteRequest{
		Message: message,
		CCVs:    ccvs,
	})
	if err != nil {
		s.logger.Err(err).Msg("failed to get disclosures for CCIP execute")
		c.JSON(500, edsv1.ErrorResponse{Error: "failed to get disclosures"})

		return
	}

	choiceContextData := map[string]any{
		"values": map[string]struct {
			Tag   string `json:"tag"`
			Value string `json:"value"`
		}{
			"off-ramp": {
				Tag:   "AV_ContractId",
				Value: disclosures.OffRamp.GetContractId(),
			},
			"global-config": {
				Tag:   "AV_ContractId",
				Value: disclosures.GlobalConfig.GetContractId(),
			},
			"token-admin-registry": {
				Tag:   "AV_ContractId",
				Value: disclosures.TokenAdminRegistry.GetContractId(),
			},
			"rmn-remote": {
				Tag:   "AV_ContractId",
				Value: disclosures.RMNRemote.GetContractId(),
			},
		},
	}
	disclosedContracts := []edsv1.DisclosedContract{
		convertDisclosedContract(disclosures.OffRamp),
		convertDisclosedContract(disclosures.GlobalConfig),
		convertDisclosedContract(disclosures.TokenAdminRegistry),
		convertDisclosedContract(disclosures.RMNRemote),
	}

	var (
		poolExtraContext  = make(map[string]any)
		tokenPool         *edsv1.DisclosedContract
		tokenPoolHoldings *edsv1.DisclosedContract
	)
	if disclosures.TokenPool != nil {
		tokenPoolDisclosure := convertDisclosedContract(disclosures.TokenPool)
		tokenPoolHoldingsDisclosure := convertDisclosedContract(disclosures.TokenPoolHolding)
		disclosedContracts = append(disclosedContracts,
			tokenPoolDisclosure,
			tokenPoolHoldingsDisclosure,
			convertDisclosedContract(disclosures.InboundRateLimiter),
			convertDisclosedContract(disclosures.OutboundRateLimiter),
		)

		// update the poolExtraContext with the relevant data.
		values := map[string]any{
			"rate-limiter": map[string]any{
				"tag":   "AV_ContractId",
				"value": disclosures.InboundRateLimiter.GetContractId(),
			},
			// TODO: outbound rate limiter not needed?
		}
		disclosedContracts = append(disclosedContracts, convertDisclosedContract(disclosures.InboundCustomBlockConfirmationsRateLimiter))
		values["inbound-custom-block-confirmations-rate-limiter"] = map[string]any{
			"tag":   "AV_ContractId",
			"value": disclosures.InboundCustomBlockConfirmationsRateLimiter.GetContractId(),
		}
		poolExtraContext = map[string]any{
			"values": map[string]any{
				"rate-limiter": values["rate-limiter"],
				"inbound-custom-block-confirmations-rate-limiter": values["inbound-custom-block-confirmations-rate-limiter"],
			},
		}
		tokenPool = &tokenPoolDisclosure
		tokenPoolHoldings = &tokenPoolHoldingsDisclosure
	}

	resp := edsv1.CCIPExecuteResponse{
		ChoiceContext: edsv1.ChoiceContext{
			ChoiceContextData:  choiceContextData,
			DisclosedContracts: disclosedContracts,
		},
		Ccvs:              make(map[string]edsv1.OptionalDisclosure, len(disclosures.CCVs)),
		PoolExtraContext:  poolExtraContext,
		TokenPool:         tokenPool,
		TokenPoolHoldings: tokenPoolHoldings,
	}

	for address, contract := range disclosures.CCVs {
		var disclosedContract *edsv1.DisclosedContract
		if contract != nil {
			d := convertDisclosedContract(contract)
			disclosedContract = &d
		}

		resp.Ccvs[address.String()] = edsv1.OptionalDisclosure{
			DisclosedContract:  disclosedContract,
			RegisteredContract: nil,
		}
	}

	c.JSON(200, resp)
}

// (POST /ccip/v1/message/send)
func (s Server) CcipSend(c *gin.Context) {
	var req edsv1.CCIPSendRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, edsv1.ErrorResponse{Error: err.Error()})
		return
	}

	// Parse requested CCVs
	ccvs := make([]contracts.InstanceAddress, len(req.Ccvs))
	for i, ccv := range req.Ccvs {
		instanceAddress := contracts.HexToInstanceAddress(ccv)
		if (instanceAddress == contracts.InstanceAddress{}) {
			c.JSON(400, edsv1.ErrorResponse{Error: fmt.Sprintf("invalid CCV address: %s", ccv)})
			return
		}
		ccvs[i] = instanceAddress
	}

	disclosures, err := s.disclosureSvc.GetCCIPSendDisclosures(c.Request.Context(), disclosure.CCIPSendRequest{
		CCVs: ccvs,
	})
	if err != nil {
		s.logger.Err(err).Msg("failed to get disclosures for CCIP send")
		c.JSON(500, edsv1.ErrorResponse{Error: "failed to get disclosures"})

		return
	}

	choiceContextData := map[string]any{
		"values": map[string]struct {
			Tag   string `json:"tag"`
			Value string `json:"value"`
		}{
			"on-ramp": {
				Tag:   "AV_ContractId",
				Value: disclosures.OnRamp.GetContractId(),
			},
			"global-config": {
				Tag:   "AV_ContractId",
				Value: disclosures.GlobalConfig.GetContractId(),
			},
			"token-admin-registry": {
				Tag:   "AV_ContractId",
				Value: disclosures.TokenAdminRegistry.GetContractId(),
			},
			"rmn-remote": {
				Tag:   "AV_ContractId",
				Value: disclosures.RMNRemote.GetContractId(),
			},
			"fee-quoter": {
				Tag:   "AV_ContractId",
				Value: disclosures.FeeQuoter.GetContractId(),
			},
		},
	}
	disclosedContracts := []edsv1.DisclosedContract{
		convertDisclosedContract(disclosures.OnRamp),
		convertDisclosedContract(disclosures.GlobalConfig),
		convertDisclosedContract(disclosures.TokenAdminRegistry),
		convertDisclosedContract(disclosures.RMNRemote),
		convertDisclosedContract(disclosures.FeeQuoter),
		convertDisclosedContract(disclosures.DefaultExecutor),
	}

	resp := edsv1.CCIPSendResponse{
		ChoiceContext: edsv1.ChoiceContext{
			ChoiceContextData:  choiceContextData,
			DisclosedContracts: disclosedContracts,
		},
		Ccvs:            make(map[string]edsv1.OptionalDisclosure, len(disclosures.CCVs)),
		DefaultExecutor: convertDisclosedContract(disclosures.DefaultExecutor),
	}

	for address, contract := range disclosures.CCVs {
		var disclosedContract *edsv1.DisclosedContract
		if contract != nil {
			d := convertDisclosedContract(contract)
			disclosedContract = &d
		}

		resp.Ccvs[address.String()] = edsv1.OptionalDisclosure{
			DisclosedContract:  disclosedContract,
			RegisteredContract: nil,
		}
	}

	c.JSON(200, resp)
}

// (GET /ccip/v1/disclosure/{instanceAddress})
func (s Server) GetDisclosure(c *gin.Context, requestedInstanceAddress string) {
	instanceAddress := contracts.HexToInstanceAddress(requestedInstanceAddress)
	if (instanceAddress == contracts.InstanceAddress{}) {
		c.JSON(400, edsv1.ErrorResponse{Error: fmt.Sprintf("invalid address: %s", requestedInstanceAddress)})
	}

	disclosedContract, err := s.disclosureSvc.GetDisclosure(c, instanceAddress)
	if err != nil {
		c.JSON(http.StatusNotFound, edsv1.ErrorResponse{Error: fmt.Sprintf("disclosure not found for address: %s", requestedInstanceAddress)})
		return
	}

	c.JSON(http.StatusOK, convertDisclosedContract(disclosedContract))
}

// (POST /ccip/v1/perPartyRouter/factory)
func (s Server) PerPartyRouterFactory(c *gin.Context) {
	var req edsv1.CCIPPerPartyRouterFactoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, edsv1.ErrorResponse{Error: err.Error()})
	}

	disclosures, err := s.disclosureSvc.GetPerPartyRouterFactory(c.Request.Context(), disclosure.PerPartyRouterFactoryRequest{})
	if err != nil {
		s.logger.Err(err).Msg("failed to get disclosures for per party router factory")
		c.JSON(500, edsv1.ErrorResponse{Error: "failed to get disclosures"})

		return
	}

	resp := edsv1.CCIPPerPartyRouterFactoryResponse{
		DisclosedContracts: []edsv1.DisclosedContract{
			convertDisclosedContract(disclosures.PerPartyRouterFactory),
		},
		PerPartyRouterFactoryId: disclosures.PerPartyRouterFactory.GetContractId(),
	}

	c.JSON(200, resp)
}

func convertDisclosedContract(contract *apiv2.DisclosedContract) edsv1.DisclosedContract {
	return edsv1.DisclosedContract{
		ContractId:       contract.ContractId,
		CreatedEventBlob: base64.StdEncoding.EncodeToString(contract.CreatedEventBlob),
		SynchronizerId:   contract.SynchronizerId,
		TemplateId:       fmt.Sprintf("%s:%s:%s", contract.GetTemplateId().GetPackageId(), contract.GetTemplateId().GetModuleName(), contract.GetTemplateId().GetEntityName()),
	}
}
