package ccv

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

const RequestSizeLimit = 100_000

func makeClient(t *testing.T, server *Server) oapiCCV.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	router.Use(middleware.RequestSizeLimiterMiddleware(RequestSizeLimit))
	oapiCCV.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiCCV.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func TestServer_PostCCVSend(t *testing.T) {
	t.Parallel()

	ccv1 := contracts.NewRawInstanceAddress("ccv1", "owner")
	ccv2 := contracts.NewRawInstanceAddress("ccv2", "owner")
	cfg := config.CCVAPIConfig{CCVs: []config.CCV{
		{ContractIdentifier: config.ContractIdentifier{InstanceAddress: ccv1.InstanceAddress()}},
	}}

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiCCV.ClientWithResponsesInterface) {
		t.Helper()
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockActiveContractStore, client
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
	}

	t.Run("Success - InstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "ccv1ContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "CommitteeVerifier",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(committeeverifier.CommitteeVerifier{
				InstanceId: types.TEXT(ccv1.InstanceID()),
				Owner:      types.PARTY(ccv1.Owner()),
				CcipOwner:  "owner",
			}),
		}}, true)

		resp, err := client.PostCCVSendWithResponse(t.Context(), ccv1.InstanceAddress().Hex(), oapiCCV.CCVSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "ccv1ContractId", resp.JSON200.ContractId)
		require.Equal(t, ccv1.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, ccv1.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "ccv1ContractId",
				TemplateId: "package:module:CommitteeVerifier",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Success - RawInstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "ccv1ContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "CommitteeVerifier",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(committeeverifier.CommitteeVerifier{
				InstanceId: types.TEXT(ccv1.InstanceID()),
				Owner:      types.PARTY(ccv1.Owner()),
				CcipOwner:  "owner",
			}),
		}}, true)

		resp, err := client.PostCCVSendWithResponse(t.Context(), ccv1.String(), oapiCCV.CCVSendRequest{
			Message: validMessage,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "ccv1ContractId", resp.JSON200.ContractId)
		require.Equal(t, ccv1.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, ccv1.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "ccv1ContractId",
				TemplateId: "package:module:CommitteeVerifier",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVSendWithBodyWithResponse(t.Context(), ccv1.InstanceAddress().String(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVSendWithBodyWithResponse(t.Context(), ccv1.InstanceAddress().String(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid address", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVSendWithResponse(t.Context(), "invalidAddress", oapiCCV.CCVSendRequest{Message: validMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Unknown address", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVSendWithResponse(t.Context(), ccv2.InstanceAddress().Hex(), oapiCCV.CCVSendRequest{Message: validMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "CCV address not found")
		})
		t.Run("CCV not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(nil, false)

			resp, err := client.PostCCVSendWithResponse(t.Context(), ccv1.InstanceAddress().Hex(), oapiCCV.CCVSendRequest{Message: validMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("CCV not returned by store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(nil, true) // Nil ActiveContract - parsing will fail

			resp, err := client.PostCCVSendWithResponse(t.Context(), ccv1.InstanceAddress().Hex(), oapiCCV.CCVSendRequest{Message: validMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_PostCCVExecute(t *testing.T) {
	t.Parallel()

	ccv1 := contracts.NewRawInstanceAddress("ccv1", "owner")
	ccv2 := contracts.NewRawInstanceAddress("ccv2", "owner")
	cfg := config.CCVAPIConfig{CCVs: []config.CCV{
		{ContractIdentifier: config.ContractIdentifier{InstanceAddress: ccv1.InstanceAddress()}},
	}}
	message := protocol.Message{
		SourceChainSelector: protocol.ChainSelector(123),
		DestChainSelector:   protocol.ChainSelector(456),
	}
	encodedMessage, err := message.Encode()
	validEncodedMessage := hex.EncodeToString(encodedMessage)
	require.NoError(t, err)

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiCCV.ClientWithResponsesInterface) {
		t.Helper()
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockActiveContractStore, client
	}

	t.Run("Success - InstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "ccv1ContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "CommitteeVerifier",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(committeeverifier.CommitteeVerifier{
				InstanceId: types.TEXT(ccv1.InstanceID()),
				Owner:      types.PARTY(ccv1.Owner()),
				CcipOwner:  "owner",
			}),
		}}, true)

		resp, err := client.PostCCVExecuteWithResponse(t.Context(), ccv1.InstanceAddress().Hex(), oapiCCV.CCVExecuteRequest{EncodedMessage: validEncodedMessage})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "ccv1ContractId", resp.JSON200.ContractId)
		require.Equal(t, ccv1.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, ccv1.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "ccv1ContractId",
				TemplateId: "package:module:CommitteeVerifier",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Success - RawInstanceAddress", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "ccv1ContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "CommitteeVerifier",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(committeeverifier.CommitteeVerifier{
				InstanceId: types.TEXT(ccv1.InstanceID()),
				Owner:      types.PARTY(ccv1.Owner()),
				CcipOwner:  "owner",
			}),
		}}, true)

		resp, err := client.PostCCVExecuteWithResponse(t.Context(), ccv1.String(), oapiCCV.CCVExecuteRequest{EncodedMessage: validEncodedMessage})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "ccv1ContractId", resp.JSON200.ContractId)
		require.Equal(t, ccv1.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, ccv1.String(), resp.JSON200.RawInstanceAddress)
		require.NotNil(t, resp.JSON200.ContextData)
		require.ElementsMatch(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "ccv1ContractId",
				TemplateId: "package:module:CommitteeVerifier",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVExecuteWithBodyWithResponse(t.Context(), ccv1.InstanceAddress().String(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVExecuteWithBodyWithResponse(t.Context(), ccv1.InstanceAddress().String(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid address", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVExecuteWithResponse(t.Context(), "invalidAddress", oapiCCV.CCVExecuteRequest{EncodedMessage: validEncodedMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Unknown address", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCVExecuteWithResponse(t.Context(), ccv2.InstanceAddress().Hex(), oapiCCV.CCVExecuteRequest{EncodedMessage: validEncodedMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "CCV address not found")
		})
		t.Run("CCV not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(nil, false)

			resp, err := client.PostCCVExecuteWithResponse(t.Context(), ccv1.InstanceAddress().Hex(), oapiCCV.CCVExecuteRequest{EncodedMessage: validEncodedMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("CCV not returned by store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(ccv1.InstanceAddress()).Return(nil, true) // Nil ActiveContract - parsing will fail

			resp, err := client.PostCCVExecuteWithResponse(t.Context(), ccv1.InstanceAddress().Hex(), oapiCCV.CCVExecuteRequest{EncodedMessage: validEncodedMessage})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_FilterContracts(t *testing.T) {
	t.Parallel()

	cfg := config.CCVAPIConfig{
		CCVs: []config.CCV{
			{ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x1")}},
			{ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x2")}},
			{ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x3")}},
		},
	}
	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything).Times(len(cfg.CCVs))
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
	require.NoError(t, err)

	t.Run("All filters", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.CCVs[0].InstanceAddress,
				cfg.CCVs[1].InstanceAddress,
				cfg.CCVs[2].InstanceAddress,
			},
		)
		assert.Len(t, res, 3)
		assert.Contains(t, res, cfg.CCVs[0].InstanceAddress)
		assert.Contains(t, res, cfg.CCVs[1].InstanceAddress)
		assert.Contains(t, res, cfg.CCVs[2].InstanceAddress)
	})
	t.Run("Filters unknown addresses", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.CCVs[0].InstanceAddress,
				contracts.HexToInstanceAddress("0x999"),
				cfg.CCVs[2].InstanceAddress,
				contracts.HexToInstanceAddress("0x998"),
				contracts.HexToInstanceAddress("0x997"),
				cfg.CCVs[1].InstanceAddress,
				contracts.HexToInstanceAddress("0x996"),
			},
		)
		assert.Len(t, res, 3)
		assert.Contains(t, res, cfg.CCVs[0].InstanceAddress)
		assert.Contains(t, res, cfg.CCVs[1].InstanceAddress)
		assert.Contains(t, res, cfg.CCVs[2].InstanceAddress)
		assert.NotContains(t, res, contracts.HexToInstanceAddress("0x999"))
		assert.NotContains(t, res, contracts.HexToInstanceAddress("0x998"))
		assert.NotContains(t, res, contracts.HexToInstanceAddress("0x997"))
		assert.NotContains(t, res, contracts.HexToInstanceAddress("0x996"))
	})
	t.Run("Filters empty list", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts([]contracts.InstanceAddress{})
		assert.Empty(t, res)
	})
	t.Run("Filters nil list", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(nil)
		assert.Empty(t, res)
	})
}
