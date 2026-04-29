package token_standard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiTokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

type TokenConfig struct {
	Type            config.TokenType
	RegistryAddress contracts.InstanceAddress
}

type Server struct {
	logger              zerolog.Logger
	activeContractStore *store.ActiveContractStore

	admin  types.PARTY
	tokens map[types.TEXT]TokenConfig
}

var (
	_ oapiTokenMetadataV1.ServerInterface     = &Server{}
	_ oapiTransferInstruction.ServerInterface = &Server{}
)

func NewServer(
	_ context.Context,
	logger zerolog.Logger,
	activeContractStore *store.ActiveContractStore,
	cfg config.TokenStandardAPIConfig,
) (*Server, error) {
	s := &Server{
		logger: logger.With().Str("component", "TokenStandardAPI").Logger(),
		admin:  types.PARTY(cfg.Admin),
		tokens: make(map[types.TEXT]TokenConfig),
	}

	for registryAddress, registry := range cfg.Registries {
		s.tokens[types.TEXT(registry.TokenId)] = TokenConfig{
			Type:            registry.TokenType,
			RegistryAddress: contracts.HexToInstanceAddress(registryAddress),
		}
		switch registry.TokenType {
		case config.TokenTypeLINK:
			s.activeContractStore.RegisterTemplates(store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(link.LinkRegistry{}),
				PartyID:    registry.PartyID,
			})
		}
	}

	return s, nil
}

// Token Metadata V1

func (s Server) GetRegistryInfo(c *gin.Context) {
	c.JSON(http.StatusOK, oapiTokenMetadataV1.GetRegistryInfoResponse{
		AdminId: string(s.admin),
		SupportedApis: map[string]int32{
			"splice-api-token-metadata-v1": 1,
		},
	})
}

func (s Server) ListInstruments(c *gin.Context, params oapiTokenMetadataV1.ListInstrumentsParams) {
	// TODO implement me
	panic("implement me")
}

func (s Server) GetInstrument(c *gin.Context, instrumentId string) {
	// TODO implement me
	panic("implement me")
}

// Transfer Instruction V1

func (s Server) GetTransferFactory(c *gin.Context) {
	var req oapiTransferInstruction.GetFactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: err.Error()})
		return
	}

	var choiceArguments splice_api_token_transfer_instruction_v1.TransferFactoryTransfer
	if err := ledger.RecordToStruct(req.ChoiceArguments, &choiceArguments); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("invalid `choiceArguments` format: %s", err.Error())})
		return
	}

	if choiceArguments.ExpectedAdmin != s.admin {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: "invalid `expectedAdmin`, must match registry admin"})
		return
	}
	instrumentId := choiceArguments.Transfer.InstrumentId
	if instrumentId.Admin != s.admin {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: "invalid `transfer.instrumentId.admin`, must match registry admin"})
		return
	}

	requestedId := instrumentId.Id

	cfg, ok := s.tokens[requestedId]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No instrument with id %s found", requestedId)})
		return
	}

	activeRegistryContract, ok := s.activeContractStore.Get(cfg.RegistryAddress)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTransferInstruction.ErrorResponse{Error: "internal server error"})
	}

	resp := oapiTransferInstruction.TransferFactoryWithChoiceContext{
		ChoiceContext: oapiTransferInstruction.ChoiceContext{},
		FactoryId:     activeRegistryContract.GetCreatedEvent().GetContractId(),
		TransferKind:  oapiTransferInstruction.Offer,
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) GetTransferInstructionAcceptContext(c *gin.Context, transferInstructionId string) {
	// TODO implement me
	panic("implement me")
}

func (s Server) GetTransferInstructionRejectContext(c *gin.Context, transferInstructionId string) {
	// TODO implement me
	panic("implement me")
}

func (s Server) GetTransferInstructionWithdrawContext(c *gin.Context, transferInstructionId string) {
	// TODO implement me
	panic("implement me")
}
