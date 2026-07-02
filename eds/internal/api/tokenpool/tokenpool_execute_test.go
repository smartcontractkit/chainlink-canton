package tokenpool

import (
	"encoding/hex"
	"math/big"
	"net/http"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
)

func TestServer_PostTokenPoolExecute_LockRelease(t *testing.T) {
	t.Parallel()

	tokenPoolAddress := contracts.NewRawInstanceAddress("lockReleasePool1", "poolOwner")
	outboundRL := contracts.NewRawInstanceAddress("outboundRL", "owner")
	inboundRL := contracts.NewRawInstanceAddress("inboundRL", "owner")
	inboundCustomRL := contracts.NewRawInstanceAddress("inboundCustomRL", "owner")
	inboundCCV1 := contracts.NewRawInstanceAddress("inboundCCV1", "owner")

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenAdmin",
		Id:    "TestToken",
	}
	encodedInstrumentId := contracts.EncodeInstrumentID(instrumentId)

	lockReleaseTokenPool := lockreleasetokenpool.LockReleaseTokenPool{
		InstanceId:   types.TEXT(tokenPoolAddress.InstanceID()),
		PoolOwner:    types.PARTY(tokenPoolAddress.Owner()),
		CcipOwner:    "ccipOwner",
		InstrumentId: instrumentId,
		Decimals:     18,
		RemoteChainConfigs: map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig{
			"123": {
				FinalityConfig:     ccipcodec.FinalityConfig{WaitForFinality: new(types.UNIT)},
				InboundRateLimiter: inboundRL.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: inboundCustomRL.Binding(),
				OutboundRateLimiter:                        outboundRL.Binding(),
				InboundCCVs:                                []chainlinkapi.RawInstanceAddress{inboundCCV1.Binding()},
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
			ModuleName: "CCIP.LockReleaseTokenPool",
			EntityName: "LockReleaseTokenPool",
		},
		CreateArguments: bindings.MarshalTemplateToRecord(lockReleaseTokenPool),
	}}

	inboundRateLimiterContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "inboundRateLimiterContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package2",
			ModuleName: "CCIP.Common",
			EntityName: "RateLimiter",
		},
	}}

	inboundCustomRateLimiterContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "inboundCustomRateLimiterContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package2",
			ModuleName: "CCIP.Common",
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

	// Create a valid encoded message with TokenTransfer and WaitForFinality
	validMessage := protocol.Message{
		SourceChainSelector: protocol.ChainSelector(123),
		DestChainSelector:   protocol.ChainSelector(456),
		Finality:            protocol.FinalityWaitForFinality,
		TokenTransfer: &protocol.TokenTransfer{
			Amount:                 big.NewInt(100),
			DestTokenAddress:       encodedInstrumentId[:],
			DestTokenAddressLength: 32,
		},
	}
	validEncodedMessageBytes, err := validMessage.Encode()
	require.NoError(t, err)
	validEncodedMessage := hex.EncodeToString(validEncodedMessageBytes)

	// Create a valid encoded message with CustomBlockConfirmations finality
	customFinalityMessage := protocol.Message{
		SourceChainSelector: protocol.ChainSelector(123),
		DestChainSelector:   protocol.ChainSelector(456),
		Finality:            5, // custom block confirmations
		TokenTransfer: &protocol.TokenTransfer{
			Amount:                 big.NewInt(100),
			DestTokenAddress:       encodedInstrumentId[:],
			DestTokenAddressLength: 32,
		},
	}
	customFinalityEncodedBytes, err := customFinalityMessage.Encode()
	require.NoError(t, err)
	customFinalityEncodedMessage := hex.EncodeToString(customFinalityEncodedBytes)

	// Create a message without token transfer
	noTokenMessage := protocol.Message{
		SourceChainSelector: protocol.ChainSelector(123),
		DestChainSelector:   protocol.ChainSelector(456),
	}
	noTokenEncodedBytes, err := noTokenMessage.Encode()
	require.NoError(t, err)
	noTokenEncodedMessage := hex.EncodeToString(noTokenEncodedBytes)

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

	t.Run("Success - WaitForFinality", func(t *testing.T) {
		t.Parallel()
		mockACS, mockIHS, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(inboundRateLimiterContract, true)
		mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return([]*apiv2.ActiveContract{holdingContract}, true)

		resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
			EncodedMessage: validEncodedMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "lockReleaseContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.Len(t, resp.JSON200.RequiredCCVs, 1)
		// Disclosed contracts: holding + token pool + rate limiter
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "holdingContractId",
				TemplateId: "package3:Splice.Token:Holding",
			},
			{
				ContractId: "lockReleaseContractId",
				TemplateId: "package1:CCIP.LockReleaseTokenPool:LockReleaseTokenPool",
			},
			{
				ContractId: "inboundRateLimiterContractId",
				TemplateId: "package2:CCIP.Common:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})

	t.Run("Success - CustomBlockConfirmations", func(t *testing.T) {
		t.Parallel()
		mockACS, mockIHS, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(inboundCustomRL.InstanceAddress()).Return(inboundCustomRateLimiterContract, true)
		mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return([]*apiv2.ActiveContract{holdingContract}, true)

		resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
			EncodedMessage: customFinalityEncodedMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "lockReleaseContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		// Disclosed contracts: holding + token pool + custom rate limiter
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "holdingContractId",
				TemplateId: "package3:Splice.Token:Holding",
			},
			{
				ContractId: "lockReleaseContractId",
				TemplateId: "package1:CCIP.LockReleaseTokenPool:LockReleaseTokenPool",
			},
			{
				ContractId: "inboundCustomRateLimiterContractId",
				TemplateId: "package2:CCIP.Common:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})

	t.Run("Success - RawInstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockACS, mockIHS, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(inboundRateLimiterContract, true)
		mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return([]*apiv2.ActiveContract{holdingContract}, true)

		resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.String(), oapiTokenPool.TokenPoolExecuteRequest{
			EncodedMessage: validEncodedMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "lockReleaseContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "holdingContractId",
				TemplateId: "package3:Splice.Token:Holding",
			},
			{
				ContractId: "lockReleaseContractId",
				TemplateId: "package1:CCIP.LockReleaseTokenPool:LockReleaseTokenPool",
			},
			{
				ContractId: "inboundRateLimiterContractId",
				TemplateId: "package2:CCIP.Common:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})

	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()

		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			resp, err := client.PostTokenPoolExecuteWithBodyWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			resp, err := client.PostTokenPoolExecuteWithBodyWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})

		t.Run("Invalid address", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), "invalidAddress", oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Unknown address", func(t *testing.T) {
			t.Parallel()
			_, _, client := setup(t)

			unknownAddress := contracts.NewRawInstanceAddress("unknown", "owner")
			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), unknownAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "token pool address not found")
		})

		t.Run("Token pool not found in store", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Invalid message hex", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: "invalidhex",
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid encoded message")
		})

		t.Run("Invalid message encoding", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: "0x1234",
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid encoded message")
		})

		t.Run("No token transfer in message", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: noTokenEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "message does not contain a token transfer")
		})
		t.Run("Wrong token in message", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			instrumentId := contracts.EncodeInstrumentID(splice_api_token_holding_v1.InstrumentId{
				Admin: "wrongAdmin",
				Id:    "WrongID",
			})
			unsupportedMessage := protocol.Message{
				SourceChainSelector: protocol.ChainSelector(999),
				DestChainSelector:   protocol.ChainSelector(456),
				TokenTransfer: &protocol.TokenTransfer{
					Amount:                 big.NewInt(100),
					DestTokenAddress:       instrumentId[:],
					DestTokenAddressLength: 32,
				},
			}
			unsupportedBytes, err := unsupportedMessage.Encode()
			require.NoError(t, err)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: hex.EncodeToString(unsupportedBytes),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "wrong pool for Message.TokenTransfer")
		})
		t.Run("Unsupported source chain selector", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			// Source chain selector 999 is not configured
			unsupportedMessage := protocol.Message{
				SourceChainSelector: protocol.ChainSelector(999),
				DestChainSelector:   protocol.ChainSelector(456),
				TokenTransfer: &protocol.TokenTransfer{
					Amount:                 big.NewInt(100),
					DestTokenAddress:       encodedInstrumentId[:],
					DestTokenAddressLength: 32,
				},
			}
			unsupportedBytes, err := unsupportedMessage.Encode()
			require.NoError(t, err)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: hex.EncodeToString(unsupportedBytes),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "unsupported source chain selector")
		})

		t.Run("Inbound rate limiter not found (WaitForFinality)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Inbound custom rate limiter not found (CustomBlockConfirmations)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(inboundCustomRL.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: customFinalityEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Holdings not found", func(t *testing.T) {
			t.Parallel()
			mockACS, mockIHS, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(inboundRateLimiterContract, true)
			mockIHS.EXPECT().GetHolding(types.PARTY(tokenPoolAddress.Owner()), instrumentId).Return(nil, false)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Token pool not returned from store (nil create arguments)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId:      "lockReleaseContractId",
				CreateArguments: nil,
			}}, true)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_PostTokenPoolExecute_BurnMint(t *testing.T) {
	t.Parallel()

	tokenPoolAddress := contracts.NewRawInstanceAddress("burnMintPool1", "poolOwner")
	outboundRL := contracts.NewRawInstanceAddress("outboundRL", "owner")
	inboundRL := contracts.NewRawInstanceAddress("inboundRL", "owner")
	inboundCustomRL := contracts.NewRawInstanceAddress("inboundCustomRL", "owner")
	inboundCCV1 := contracts.NewRawInstanceAddress("inboundCCV1", "owner")

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenAdmin",
		Id:    "TestToken",
	}
	encodedInstrumentId := contracts.EncodeInstrumentID(instrumentId)

	burnMintPool := burnminttokenpool.BurnMintTokenPool{
		InstanceId:   types.TEXT(tokenPoolAddress.InstanceID()),
		PoolOwner:    types.PARTY(tokenPoolAddress.Owner()),
		CcipOwner:    "ccipOwner",
		InstrumentId: instrumentId,
		Decimals:     18,
		RemoteChainConfigs: map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
			"123": {
				FinalityConfig:     ccipcodec.FinalityConfig{WaitForFinality: new(types.UNIT)},
				InboundRateLimiter: inboundRL.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: inboundCustomRL.Binding(),
				OutboundRateLimiter:                        outboundRL.Binding(),
				InboundCCVs:                                []chainlinkapi.RawInstanceAddress{inboundCCV1.Binding()},
			},
		},
		TransferTimeout: burnminttokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
	}

	tokenPoolActiveContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "burnMintContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package1",
			ModuleName: "CCIP.BurnMintTokenPool",
			EntityName: "BurnMintTokenPool",
		},
		CreateArguments: bindings.MarshalTemplateToRecord(burnMintPool),
	}}

	inboundRateLimiterContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "inboundRateLimiterContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package2",
			ModuleName: "CCIP.Common",
			EntityName: "RateLimiter",
		},
	}}

	inboundCustomRateLimiterContract := &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
		ContractId: "inboundCustomRateLimiterContractId",
		TemplateId: &apiv2.Identifier{
			PackageId:  "package2",
			ModuleName: "CCIP.Common",
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

	// Create a valid encoded message with TokenTransfer and WaitForFinality
	validMessage := protocol.Message{
		SourceChainSelector: protocol.ChainSelector(123),
		DestChainSelector:   protocol.ChainSelector(456),
		Finality:            protocol.FinalityWaitForFinality,
		TokenTransfer: &protocol.TokenTransfer{
			Amount:                 big.NewInt(100),
			DestTokenAddress:       encodedInstrumentId[:],
			DestTokenAddressLength: 32,
		},
	}
	validEncodedMessageBytes, err := validMessage.Encode()
	require.NoError(t, err)
	validEncodedMessage := hex.EncodeToString(validEncodedMessageBytes)

	// Custom block confirmations finality
	customFinalityMessage := protocol.Message{
		SourceChainSelector: protocol.ChainSelector(123),
		DestChainSelector:   protocol.ChainSelector(456),
		Finality:            5,
		TokenTransfer: &protocol.TokenTransfer{
			Amount:                 big.NewInt(100),
			DestTokenAddress:       encodedInstrumentId[:],
			DestTokenAddressLength: 32,
		},
	}
	customFinalityEncodedBytes, err := customFinalityMessage.Encode()
	require.NoError(t, err)
	customFinalityEncodedMessage := hex.EncodeToString(customFinalityEncodedBytes)

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

	t.Run("Success - WaitForFinality", func(t *testing.T) {
		t.Parallel()
		mockACS, _, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(inboundRateLimiterContract, true)

		resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
			EncodedMessage: validEncodedMessage,
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
				TemplateId: "package1:CCIP.BurnMintTokenPool:BurnMintTokenPool",
			},
			{
				ContractId: "inboundRateLimiterContractId",
				TemplateId: "package2:CCIP.Common:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})

	t.Run("Success - CustomBlockConfirmations", func(t *testing.T) {
		t.Parallel()
		mockACS, _, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(inboundCustomRL.InstanceAddress()).Return(inboundCustomRateLimiterContract, true)

		resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
			EncodedMessage: customFinalityEncodedMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "burnMintContractId", resp.JSON200.ContractId)
		require.Equal(t, tokenPoolAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "burnMintContractId",
				TemplateId: "package1:CCIP.BurnMintTokenPool:BurnMintTokenPool",
			},
			{
				ContractId: "inboundCustomRateLimiterContractId",
				TemplateId: "package2:CCIP.Common:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})

	t.Run("Success - RawInstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockACS, _, client := setup(t)

		mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
		mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(inboundRateLimiterContract, true)

		resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.String(), oapiTokenPool.TokenPoolExecuteRequest{
			EncodedMessage: validEncodedMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "burnMintContractId", resp.JSON200.ContractId)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "burnMintContractId",
				TemplateId: "package1:CCIP.BurnMintTokenPool:BurnMintTokenPool",
			},
			{
				ContractId: "inboundRateLimiterContractId",
				TemplateId: "package2:CCIP.Common:RateLimiter",
			},
		}, resp.JSON200.DisclosedContracts)
	})

	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()

		t.Run("Wrong token in message", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			instrumentId := contracts.EncodeInstrumentID(splice_api_token_holding_v1.InstrumentId{
				Admin: "wrongAdmin",
				Id:    "WrongID",
			})
			unsupportedMessage := protocol.Message{
				SourceChainSelector: protocol.ChainSelector(999),
				DestChainSelector:   protocol.ChainSelector(456),
				TokenTransfer: &protocol.TokenTransfer{
					Amount:                 big.NewInt(100),
					DestTokenAddress:       instrumentId[:],
					DestTokenAddressLength: 32,
				},
			}
			unsupportedBytes, err := unsupportedMessage.Encode()
			require.NoError(t, err)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: hex.EncodeToString(unsupportedBytes),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "wrong pool for Message.TokenTransfer")
		})
		t.Run("Unsupported source chain selector", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)

			unsupportedMessage := protocol.Message{
				SourceChainSelector: protocol.ChainSelector(999),
				DestChainSelector:   protocol.ChainSelector(456),
				TokenTransfer: &protocol.TokenTransfer{
					Amount:                 big.NewInt(100),
					DestTokenAddress:       encodedInstrumentId[:],
					DestTokenAddressLength: 32,
				},
			}
			unsupportedBytes, err := unsupportedMessage.Encode()
			require.NoError(t, err)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: hex.EncodeToString(unsupportedBytes),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "unsupported source chain selector")
		})

		t.Run("Inbound rate limiter not found (WaitForFinality)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(inboundRL.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})

		t.Run("Inbound custom rate limiter not found (CustomBlockConfirmations)", func(t *testing.T) {
			t.Parallel()
			mockACS, _, client := setup(t)

			mockACS.EXPECT().Get(tokenPoolAddress.InstanceAddress()).Return(tokenPoolActiveContract, true)
			mockACS.EXPECT().Get(inboundCustomRL.InstanceAddress()).Return(nil, false)

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: customFinalityEncodedMessage,
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

			resp, err := client.PostTokenPoolExecuteWithResponse(t.Context(), tokenPoolAddress.InstanceAddress().Hex(), oapiTokenPool.TokenPoolExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}
