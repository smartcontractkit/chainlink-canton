package tokenpool

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
)

const RequestSizeLimit = 100_000

func makeClient(t *testing.T, server *Server) oapiTokenPool.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	router.Use(middleware.RequestSizeLimiterMiddleware(RequestSizeLimit))
	oapiTokenPool.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiTokenPool.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func TestServer_PostTokenPoolSend_LockRelease(t *testing.T) {
	t.Parallel()

	tokenPoolAddress := contracts.NewRawInstanceAddress("lockReleasePool1", "poolOwner")
	outboundRL := contracts.NewRawInstanceAddress("outboundRL", "owner")
	inboundRL := contracts.NewRawInstanceAddress("inboundRL", "owner")
	inboundCustomRL := contracts.NewRawInstanceAddress("inboundCustomRL", "owner")
	outboundCCV1 := contracts.NewRawInstanceAddress("outboundCCV1", "owner")
	outboundCCV2 := contracts.NewRawInstanceAddress("outboundCCV2", "owner")

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenAdmin",
		Id:    "TestToken",
	}

	lockReleaseTokenPool := lockreleasetokenpool.LockReleaseTokenPool{
		InstanceId:   types.TEXT(tokenPoolAddress.InstanceID()),
		PoolOwner:    types.PARTY(tokenPoolAddress.Owner()),
		CcipOwner:    "ccipOwner",
		InstrumentId: instrumentId,
		Decimals:     18,
		RemoteChainConfigs: map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig{
			"456": {
				FinalityConfig:     ccipcodec.FinalityConfig{WaitForFinality: new(types.UNIT)},
				InboundRateLimiter: inboundRL.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: inboundCustomRL.Binding(),
				OutboundRateLimiter:                        outboundRL.Binding(),
				OutboundCCVs:                               []chainlinkapi.RawInstanceAddress{outboundCCV1.Binding(), outboundCCV2.Binding()},
			},
		},
		TokenTransferFeeConfigs: nil,
		PoolReceiveContext:      splice_api_token_metadata_v1.ChoiceContext{},
		TransferTimeout:         lockreleasetokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
		Deps:                    lockreleasetokenpool.LockReleaseTokenPoolDeps{},
	}

	tokenPoolActiveContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "lockReleaseContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package1",
			ModuleName: "CCIP.LockReleaseTokenPoolV2",
			EntityName: "LockReleaseTokenPool",
		},
		CreateArguments: bindings.MarshalTemplateToRecord(lockReleaseTokenPool),
	}}

	rateLimiterContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "rateLimiterContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package2",
			ModuleName: "CCIP.RateLimiterV2",
			EntityName: "RateLimiter",
		},
	}}

	holdingContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "holdingContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package3",
			ModuleName: "Splice.Token",
			EntityName: "Holding",
		},
	}}

	cfg := config.TokenPoolAPIConfig{
		TokenPools: map[string]config.TokenPool{
			"LRPool": {
				Type:      config.TokenPoolTypeLockRelease,
				PoolOwner: tokenPoolAddress.Owner(),
				ContractIdentifier: config.ContractIdentifier{
					InstanceAddress: tokenPoolAddress.InstanceAddress(),
				},
			},
		},
	}

	validMessage := oapiCommon.Message{
		DestinationChainSelector: "456",
		Executor: struct {
			Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiCommon.MessageExecutorType `json:"type"`
		}{
			Type: oapiCommon.NoExecutor,
		},
		FeeToken: oapiCommon.InstrumentId{
			Admin: "feeAdmin",
			Id:    "LINK",
		},
		Payload:  "0xdeadbeef",
		Receiver: "0x1234567890",
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: "100",
			Token: oapiCommon.InstrumentId{
				Admin: "tokenAdmin",
				Id:    "TestToken",
			},
		},
	}

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, *mocks.MockInstrumentHoldingStore, oapiTokenPool.ClientWithResponsesInterface) {
		t.Helper()
		mockACS := mocks.NewMockActiveContractStore(t)
		mockACS.EXPECT().RegisterTemplates(mock.Anything).Maybe()
		mockIHS := mocks.NewMockInstrumentHoldingStore(t)
		mockIHS.EXPECT().RegisterParty(mock.Anything).Maybe()
		server, err := NewServer(t.Context(), zerolog.Nop(), mockACS, mockIHS, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockACS, mockIHS, client
	}

	t.Run("Success - with holdings", func(t *testing.T) {
		t.Parallel()
		mockACS, mockIHS, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(rateLimiterContract, true)
		mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return([]*apiv2.ActiveContract{holdingContract}, true)

		resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "lockReleaseContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.Len(t, resp.JSON200.RequiredCCVs, 2)
		// Disclosed contracts: existing pool holdings + token pool + rate limiter.
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "holdingContractId",
				TemplateId: "package3:Splice.Token:Holding",
			},
			{
				ContractId: "lockReleaseContractId",
				TemplateId: "package1:CCIP.LockReleaseTokenPoolV2:LockReleaseTokenPool",
			},
			{
				ContractId: "rateLimiterContractId",
				TemplateId: "package2:CCIP.RateLimiterV2:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Success - without holdings", func(t *testing.T) {
		t.Parallel()
		mockACS, mockIHS, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(rateLimiterContract, true)
		mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return(nil, false)

		resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "lockReleaseContractId", resp.JSON200.ContractId)
		// Disclosed contracts: token pool + rate limiter (no holdings)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "lockReleaseContractId",
				TemplateId: "package1:CCIP.LockReleaseTokenPoolV2:LockReleaseTokenPool",
			},
			{
				ContractId: "rateLimiterContractId",
				TemplateId: "package2:CCIP.RateLimiterV2:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Success - RawInstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockACS, mockIHS, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(rateLimiterContract, true)
		mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return(nil, false)

		resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.String(), oapiTokenPool.TokenPoolSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "lockReleaseContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "lockReleaseContractId",
				TemplateId: "package1:CCIP.LockReleaseTokenPoolV2:LockReleaseTokenPool",
			},
			{
				ContractId: "rateLimiterContractId",
				TemplateId: "package2:CCIP.RateLimiterV2:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			resp, err := client.PostTokenPoolSendWithBodyWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			resp, err := client.PostTokenPoolSendWithBodyWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid address", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), "invalidAddress", oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Unknown address", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			unknownAddress := contracts.NewRawInstanceAddress("unknown", "owner")
			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), unknownAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "token pool address not found")
		})
		t.Run("Token pool not found in store", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Invalid destination chain selector", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			invalidMessage := validMessage
			invalidMessage.DestinationChainSelector = "notanumber"

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: invalidMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid destination chain selector")
		})
		t.Run("No token transfer in message", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			noTokenMessage := validMessage
			noTokenMessage.TokenTransfer = nil

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: noTokenMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "message does not contain a token transfer")
		})
		t.Run("Wrong token for pool", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			unsupportedTokenMessage := validMessage
			unsupportedTokenMessage.TokenTransfer = &oapiCommon.TokenTransfer{
				Amount: "1",
				Token: oapiCommon.InstrumentId{
					Admin: "wrongAdmin",
					Id:    "WrongID",
				},
			}

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: unsupportedTokenMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "wrong pool for Message.TokenTransfer")
		})
		t.Run("Unsupported destination chain selector", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			unsupportedChainMessage := validMessage
			unsupportedChainMessage.DestinationChainSelector = "999"

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: unsupportedChainMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "unsupported destination chain selector")
		})
		t.Run("Outbound rate limiter not found", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Token pool not returned from store (nil created event)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId:      "lockReleaseContractId",
				CreateArguments: nil, // nil create arguments will cause parse failure
			}}, true)

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_PostTokenPoolSend_BurnMint(t *testing.T) {
	t.Parallel()

	tokenPoolAddress := contracts.NewRawInstanceAddress("burnMintPool1", "poolOwner")
	outboundRL := contracts.NewRawInstanceAddress("outboundRL", "owner")
	inboundRL := contracts.NewRawInstanceAddress("inboundRL", "owner")
	inboundCustomRL := contracts.NewRawInstanceAddress("inboundCustomRL", "owner")
	outboundCCV1 := contracts.NewRawInstanceAddress("outboundCCV1", "owner")

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenAdmin",
		Id:    "TestToken",
	}

	burnMintPool := burnminttokenpool.BurnMintTokenPool{
		InstanceId:   types.TEXT(tokenPoolAddress.InstanceID()),
		PoolOwner:    types.PARTY(tokenPoolAddress.Owner()),
		CcipOwner:    "ccipOwner",
		InstrumentId: instrumentId,
		Decimals:     18,
		RemoteChainConfigs: map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
			"456": {
				FinalityConfig:     ccipcodec.FinalityConfig{WaitForFinality: new(types.UNIT)},
				InboundRateLimiter: inboundRL.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: inboundCustomRL.Binding(),
				OutboundRateLimiter:                        outboundRL.Binding(),
				OutboundCCVs:                               []chainlinkapi.RawInstanceAddress{outboundCCV1.Binding()},
			},
		},
		TransferTimeout: burnminttokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
	}

	tokenPoolActiveContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "burnMintContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package1",
			ModuleName: "CCIP.BurnMintTokenPoolV2",
			EntityName: "BurnMintTokenPool",
		},
		CreateArguments: bindings.MarshalTemplateToRecord(burnMintPool),
	}}

	rateLimiterContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "rateLimiterContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package2",
			ModuleName: "CCIP.RateLimiterV2",
			EntityName: "RateLimiter",
		},
	}}

	cfg := config.TokenPoolAPIConfig{
		TokenPools: map[string]config.TokenPool{
			"BMPool": {
				Type:      config.TokenPoolTypeBurnMint,
				PoolOwner: tokenPoolAddress.Owner(),
				ContractIdentifier: config.ContractIdentifier{
					InstanceAddress: tokenPoolAddress.InstanceAddress(),
				},
			},
		},
	}

	validMessage := oapiCommon.Message{
		DestinationChainSelector: "456",
		Executor: struct {
			Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiCommon.MessageExecutorType `json:"type"`
		}{
			Type: oapiCommon.NoExecutor,
		},
		FeeToken: oapiCommon.InstrumentId{
			Admin: "feeAdmin",
			Id:    "LINK",
		},
		Payload:  "0xdeadbeef",
		Receiver: "0x1234567890",
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: "100",
			Token: oapiCommon.InstrumentId{
				Admin: "tokenAdmin",
				Id:    "TestToken",
			},
		},
	}

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, *mocks.MockInstrumentHoldingStore, oapiTokenPool.ClientWithResponsesInterface) {
		t.Helper()
		mockACS := mocks.NewMockActiveContractStore(t)
		mockACS.EXPECT().RegisterTemplates(mock.Anything).Maybe()
		mockIHS := mocks.NewMockInstrumentHoldingStore(t)
		mockIHS.EXPECT().RegisterParty(mock.Anything).Maybe()
		server, err := NewServer(t.Context(), zerolog.Nop(), mockACS, mockIHS, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockACS, mockIHS, client
	}
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockACS, _, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(rateLimiterContract, true)

		resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "burnMintContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.Len(t, resp.JSON200.RequiredCCVs, 1)
		// Disclosed contracts: token pool + rate limiter
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "burnMintContractId",
				TemplateId: "package1:CCIP.BurnMintTokenPoolV2:BurnMintTokenPool",
			},
			{
				ContractId: "rateLimiterContractId",
				TemplateId: "package2:CCIP.RateLimiterV2:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Success - RawInstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockACS, _, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(rateLimiterContract, true)

		resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.String(), oapiTokenPool.TokenPoolSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "burnMintContractId", resp.JSON200.ContractId)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "burnMintContractId",
				TemplateId: "package1:CCIP.BurnMintTokenPoolV2:BurnMintTokenPool",
			},
			{
				ContractId: "rateLimiterContractId",
				TemplateId: "package2:CCIP.RateLimiterV2:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Wrong token for pool", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			unsupportedTokenMessage := validMessage
			unsupportedTokenMessage.TokenTransfer = &oapiCommon.TokenTransfer{
				Amount: "1",
				Token: oapiCommon.InstrumentId{
					Admin: "wrongAdmin",
					Id:    "WrongID",
				},
			}

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: unsupportedTokenMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "wrong pool for Message.TokenTransfer")
		})
		t.Run("Unsupported destination chain selector", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			unsupportedChainMessage := validMessage
			unsupportedChainMessage.DestinationChainSelector = "999"

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: unsupportedChainMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "unsupported destination chain selector")
		})
		t.Run("Outbound rate limiter not found", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(outboundRL.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Token pool not returned from store (nil create arguments)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId:      "burnMintContractId",
				CreateArguments: nil,
			}}, true)

			resp, err := client.PostTokenPoolSendWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}
