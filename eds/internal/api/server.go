package api

import (
	"encoding/base64"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

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
	}

	disclosures, err := s.disclosureSvc.GetCCIPExecuteDisclosures(c.Request.Context(), disclosure.CCIPExecuteRequest{})
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
		convertDisclosedContract(disclosures.DefaultCCV),
	}

	resp := edsv1.CCIPExecuteResponse{
		ChoiceContext: &edsv1.ChoiceContext{
			ChoiceContextData:  choiceContextData,
			DisclosedContracts: disclosedContracts,
		},
	}

	c.JSON(200, resp)
}

// (POST /ccip/v1/message/send)
func (s Server) CcipSend(c *gin.Context) {
	var req edsv1.CCIPSendRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, edsv1.ErrorResponse{Error: err.Error()})
	}
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
