package token_standard

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/global"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	oapiTokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

const (
	DefaultPageSize = int32(25)
	MaximumPageSize = int32(1024)
)

type TokenConfig struct {
	Type            config.TokenType
	RegistryAddress contracts.InstanceAddress
}

type Server struct {
	logger              zerolog.Logger
	activeContractStore store.ActiveContractStoreInterface

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
	activeContractStore store.ActiveContractStoreInterface,
	cfg config.TokenStandardAPIConfig,
) (*Server, error) {
	s := &Server{
		logger:              logger.With().Str("component", "TokenStandardAPI").Logger(),
		activeContractStore: activeContractStore,
		admin:               types.PARTY(cfg.Admin),
		tokens:              make(map[types.TEXT]TokenConfig),
	}

	for _, registry := range cfg.Registries {
		// Validate config
		if len(registry.TokenId) == 0 {
			return nil, fmt.Errorf("empty TokenId in config")
		}
		if _, ok := s.tokens[types.TEXT(registry.TokenId)]; ok {
			return nil, fmt.Errorf("duplicate TokenId in config: %s", registry.TokenId)
		}

		s.tokens[types.TEXT(registry.TokenId)] = TokenConfig{
			Type:            registry.TokenType,
			RegistryAddress: registry.InstanceAddress,
		}
		switch registry.TokenType {
		case config.TokenTypeLINK:
			s.activeContractStore.RegisterTemplates(
				store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(link.LinkRegistry{}),
					PartyID:    registry.PartyID,
				},
				store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(link.LinkTransferInstruction{}),
					PartyID:    registry.PartyID,
				},
			)
		default:
			return nil, fmt.Errorf("unsupported token type: %s", registry.TokenType)
		}
	}

	return s, nil
}

func (s Server) GetRegistryInfo(c *gin.Context) {
	c.JSON(http.StatusOK, oapiTokenMetadataV1.GetRegistryInfoResponse{
		AdminId: string(s.admin),
		SupportedApis: map[string]int32{
			"splice-api-token-metadata-v1": 1,
		},
	})
}

// Token Metadata V1

func (s Server) ListInstruments(c *gin.Context, params oapiTokenMetadataV1.ListInstrumentsParams) {
	pageSize := DefaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	if pageSize <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: "pageSize must be a positive integer"})
		return
	}
	if pageSize > MaximumPageSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: fmt.Sprintf("pageSize cannot be greater than %d", MaximumPageSize)})
		return
	}
	// Using the index of the alpha-sorted list of configured tokens as the pageToken
	// If not specified, it will start at 0
	// If there are more elements to be returned, the next index will be returned as the response's pageToken
	pageToken := 0
	if params.PageToken != nil {
		parsedIndex, err := strconv.Atoi(*params.PageToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: "invalid pageToken"})
			return
		}
		pageToken = parsedIndex
	}

	sortedTokens := maps.Keys(s.tokens)
	slices.Sort(sortedTokens)

	if pageToken < 0 || pageToken >= len(sortedTokens) {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTokenMetadataV1.ErrorResponse{Error: "invalid pageToken"})
		return
	}

	endIndex := min(pageToken+int(pageSize), len(sortedTokens))

	instrumentIds := make([]string, endIndex-pageToken)
	for i, tokenId := range sortedTokens[pageToken:endIndex] {
		instrumentIds[i] = string(tokenId)
	}

	instruments := make([]oapiTokenMetadataV1.Instrument, len(instrumentIds))
	for i, id := range instrumentIds {
		// TODO add Name/Symbol to config?
		// TODO add total supply information to response
		instruments[i] = oapiTokenMetadataV1.Instrument{
			Decimals: 10,
			Id:       id,
			Name:     id,
			SupportedApis: map[string]int32{
				"splice-api-token-transfer-instruction-v1": 1,
			},
			Symbol:          id,
			TotalSupply:     nil,
			TotalSupplyAsOf: nil,
		}
	}

	var nextPageToken *string
	if endIndex < len(sortedTokens) {
		nextPageToken = new(strconv.Itoa(endIndex))
	}

	c.JSON(http.StatusOK, oapiTokenMetadataV1.ListInstrumentsResponse{
		Instruments:   instruments,
		NextPageToken: nextPageToken,
	})
}

func (s Server) GetInstrument(c *gin.Context, instrumentId string) {
	_, ok := s.tokens[types.TEXT(instrumentId)]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTokenMetadataV1.ErrorResponse{Error: fmt.Sprintf("No instrument with id %s found", instrumentId)})
		return
	}

	// TODO add Name/Symbol to config?
	// TODO add total supply information to response
	c.JSON(http.StatusOK, oapiTokenMetadataV1.Instrument{
		Decimals: 10,
		Id:       instrumentId,
		Name:     instrumentId,
		SupportedApis: map[string]int32{
			"splice-api-token-transfer-instruction-v1": 1,
		},
		Symbol:          instrumentId,
		TotalSupply:     nil,
		TotalSupplyAsOf: nil,
	})
}

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
		s.logger.Error().Str("registryAddress", cfg.RegistryAddress.String()).Msg("Failed to retrieve active registry contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTransferInstruction.ErrorResponse{Error: "internal server error"})
		return
	}
	disclosedRegistry := ActiveContractToDisclosedContract(activeRegistryContract, req.ExcludeDebugFields != nil && !*req.ExcludeDebugFields)

	resp := oapiTransferInstruction.TransferFactoryWithChoiceContext{
		ChoiceContext: oapiTransferInstruction.ChoiceContext{
			ChoiceContextData:  nil,
			DisclosedContracts: []oapiTransferInstruction.DisclosedContract{disclosedRegistry},
		},
		FactoryId:    activeRegistryContract.GetCreatedEvent().GetContractId(),
		TransferKind: oapiTransferInstruction.Offer,
	}

	c.JSON(http.StatusOK, resp)
}

// Transfer Instruction V1

func (s Server) GetTransferInstructionAcceptContext(c *gin.Context, transferInstructionId string) {
	var req oapiTransferInstruction.GetChoiceContextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, oapiTransferInstruction.ErrorResponse{Error: err.Error()})
		return
	}

	transferInstruction, ok := s.activeContractStore.GetByContractId(types.CONTRACT_ID(transferInstructionId))
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}
	parsedTransferInstruction, err := bindings.UnmarshalCreatedEvent[link.LinkTransferInstruction](transferInstruction.GetCreatedEvent())
	if err != nil {
		// If the unmarshal fails, we assume that the user provided an incorrect ContractId/a ContractId of a differently-types contract
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}
	instrumentId := parsedTransferInstruction.InstructionTransfer.InstrumentId
	if instrumentId.Admin != s.admin {
		// If the TransferInstruction is for a different instrument admin, we treat it as not-found
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}
	cfg, ok := s.tokens[instrumentId.Id]
	if !ok {
		// If the TransferInstruction is for a non-configured instrument, we treat it as not-found
		c.AbortWithStatusJSON(http.StatusNotFound, oapiTransferInstruction.ErrorResponse{Error: fmt.Sprintf("No transfer instruction with id %s found", transferInstructionId)})
		return
	}

	// Request is validated

	activeRegistryContract, ok := s.activeContractStore.Get(cfg.RegistryAddress)
	if !ok {
		s.logger.Error().Str("registryAddress", cfg.RegistryAddress.String()).Msg("Failed to retrieve active registry contract")
		c.AbortWithStatusJSON(http.StatusInternalServerError, oapiTransferInstruction.ErrorResponse{Error: "internal server error"})
		return
	}
	disclosedRegistry := ActiveContractToDisclosedContract(activeRegistryContract, false)

	resp := oapiTransferInstruction.ChoiceContext{
		ChoiceContextData: nil,
		// Don't return a disclosure for the TransferInstruction itself - the receiver will already be an observer on that contract
		DisclosedContracts: []oapiTransferInstruction.DisclosedContract{disclosedRegistry},
	}

	c.JSON(http.StatusOK, resp)
}

func (s Server) GetTransferInstructionRejectContext(c *gin.Context, transferInstructionId string) {
	s.GetTransferInstructionAcceptContext(c, transferInstructionId)
}

func (s Server) GetTransferInstructionWithdrawContext(c *gin.Context, transferInstructionId string) {
	s.GetTransferInstructionAcceptContext(c, transferInstructionId)
}

var _ global.InstanceAddressFilter = &Server{}

func (s Server) FilterContracts(addresses []contracts.InstanceAddress) []contracts.InstanceAddress {
	var out []contracts.InstanceAddress
	for _, address := range addresses {
		for _, tokenConfig := range s.tokens {
			if address == tokenConfig.RegistryAddress {
				out = append(out, address)
			}
		}
	}

	return out
}
