package executor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/core"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
)

const RequestSizeLimit = 100_000

func makeClient(t *testing.T, server *Server) oapiExecutor.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	router.Use(middleware.RequestSizeLimiterMiddleware(RequestSizeLimit))
	oapiExecutor.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiExecutor.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func TestServer_PostExecutorSend(t *testing.T) {
	t.Parallel()

	t.Run("With Allowlist Enabled", func(t *testing.T) {
		t.Parallel()

		executorAddress := contracts.NewRawInstanceAddress("executor1", "owner")
		ccv1 := contracts.NewRawInstanceAddress("ccv1", "owner")
		ccv2 := contracts.NewRawInstanceAddress("ccv2", "owner")
		ccv3 := contracts.NewRawInstanceAddress("ccv3", "owner")
		ccv4 := contracts.NewRawInstanceAddress("ccv4", "owner")
		executor := executorBinding.Executor{
			InstanceId:    types.TEXT(executorAddress.InstanceID()),
			Owner:         types.PARTY(executorAddress.Owner()),
			MaxCCVsPerMsg: 2,
			DynamicConfig: executorBinding.DynamicConfig{
				AllowedFinalityConfig: core.FinalityConfig{WaitForFinality: new(types.UNIT)},
				CcvAllowlistEnabled:   true,
			},
			AllowedCCVs: []chainlinkapi.RawInstanceAddress{
				ccv1.Binding(),
				ccv2.Binding(),
				ccv3.Binding(),
			},
			RemoteChainConfigs: nil,
		}
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{
			Executors: []config.Executor{
				{
					ContractIdentifier: config.ContractIdentifier{
						PartyID:         executorAddress.Owner(),
						InstanceAddress: executorAddress.InstanceAddress(),
					},
				},
			},
		})
		require.NoError(t, err)
		client := makeClient(t, server)

		mockActiveContractStore.EXPECT().Get(executorAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contract1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package1",
				ModuleName: "module1",
				EntityName: "entity1",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(executor),
		}}, true)

		t.Run("With InstanceAddress", func(t *testing.T) {
			// All requested addresses are InstanceAddresses
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.InstanceAddress().Hex(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv1.InstanceAddress().Hex()),
					testhelpers.MakeRawAddress(ccv2.InstanceAddress().Hex()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeHashedAddress(executorAddress.InstanceAddress().Hex())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Equal(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "contract1",
					TemplateId: "package1:module1:entity1",
				},
			}, resp.JSON200.DisclosedContracts)
			require.Equal(t, "contract1", resp.JSON200.ContractId)
			require.Equal(t, executorAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
			require.Equal(t, executorAddress.String(), resp.JSON200.RawInstanceAddress)
			require.NotNil(t, resp.JSON200.ContextData)
		})
		t.Run("With RawInstanceAddress", func(t *testing.T) {
			// All requested addresses are RawInstanceAddresses
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress(ccv3.String()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeRawAddress(executorAddress.String())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Equal(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "contract1",
					TemplateId: "package1:module1:entity1",
				},
			}, resp.JSON200.DisclosedContracts)
		})
		t.Run("Too many CCVs", func(t *testing.T) {
			// Request too many CCVs
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv1.String()),
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress(ccv3.String()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeRawAddress(executorAddress.String())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "too many CCVs")
		})
		t.Run("Unallowed CCV", func(t *testing.T) {
			// Request a CCV that isn't allowed
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress(ccv4.String()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeRawAddress(executorAddress.String())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "not in the allowed CCV list")
		})
		t.Run("Invalid CCV", func(t *testing.T) {
			// Request a CCV with an invalid address
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress("notavalidaddress"),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeRawAddress(executorAddress.String())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid CCV")
		})
		t.Run("Address parameter mismatch", func(t *testing.T) {
			// Request a message with executor=WithAddress, but path parameter doesn't match the address in the message itself
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress("notavalidaddress"),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeHashedAddress("0x42")), // address that doesn't match the path parameter
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "executor in message doesn't match requested executor")
		})
		t.Run("NoExecutor in message", func(t *testing.T) {
			// Request a message with executor=NoExecutor, should fail
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress("notavalidaddress"),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
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
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no_executor")
		})
		t.Run("Empty/Default executor in message", func(t *testing.T) {
			// Request a message with executor=Empty, should succeed
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.InstanceAddress().Hex(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv1.InstanceAddress().Hex()),
					testhelpers.MakeRawAddress(ccv2.InstanceAddress().Hex()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type: oapiCommon.Empty,
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Equal(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "contract1",
					TemplateId: "package1:module1:entity1",
				},
			}, resp.JSON200.DisclosedContracts)
		})
		t.Run("Unknown executor", func(t *testing.T) {
			// Request an executor that in unknown/unconfigured
			t.Parallel()

			executorAddress := contracts.HexToInstanceAddress("0x7777")
			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.Hex(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv1.InstanceAddress().Hex()),
					testhelpers.MakeRawAddress(ccv2.InstanceAddress().Hex()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeHashedAddress(executorAddress.Hex())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "not found")
		})
	})
	t.Run("With Allowlist Disabled", func(t *testing.T) {
		t.Parallel()

		executorAddress := contracts.NewRawInstanceAddress("executor1", "owner")
		ccv1 := contracts.NewRawInstanceAddress("ccv1", "owner")
		ccv2 := contracts.NewRawInstanceAddress("ccv2", "owner")
		ccv3 := contracts.NewRawInstanceAddress("ccv3", "owner")
		ccv4 := contracts.NewRawInstanceAddress("ccv4", "owner")
		executor := executorBinding.Executor{
			InstanceId:    types.TEXT(executorAddress.InstanceID()),
			Owner:         types.PARTY(executorAddress.Owner()),
			MaxCCVsPerMsg: 3,
			DynamicConfig: executorBinding.DynamicConfig{
				AllowedFinalityConfig: core.FinalityConfig{WaitForFinality: new(types.UNIT)},
				CcvAllowlistEnabled:   false, // allowlist disabled, CCVs in the request shouldn't be checked
			},
			AllowedCCVs: []chainlinkapi.RawInstanceAddress{
				ccv1.Binding(),
				ccv2.Binding(),
				ccv3.Binding(),
			},
			RemoteChainConfigs: nil,
		}
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{
			Executors: []config.Executor{
				{
					ContractIdentifier: config.ContractIdentifier{
						PartyID:         "owner",
						InstanceAddress: executorAddress.InstanceAddress(),
					},
				},
			},
		})
		require.NoError(t, err)
		client := makeClient(t, server)

		mockActiveContractStore.EXPECT().Get(executorAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contract1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package1",
				ModuleName: "module1",
				EntityName: "entity1",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(executor),
		}}, true)

		t.Run("With InstanceAddress", func(t *testing.T) {
			// All requested addresses are InstanceAddresses
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.InstanceAddress().Hex(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv3.InstanceAddress().Hex()), // allowed CCV
					testhelpers.MakeRawAddress(ccv4.InstanceAddress().Hex()), // unallowed CCV - should not be checked
					testhelpers.MakeRawAddress("invalidaddress"),             // invalid address - should not be checked
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeHashedAddress(executorAddress.String())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Equal(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "contract1",
					TemplateId: "package1:module1:entity1",
				},
			}, resp.JSON200.DisclosedContracts)
		})
		t.Run("Too many CCVs", func(t *testing.T) {
			// Request too many CCVs
			t.Parallel()

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.String(), oapiExecutor.ExecutorSendRequest{
				Ccvs: []oapiCommon.RawOrHashedAddress{
					testhelpers.MakeRawAddress(ccv1.String()),
					testhelpers.MakeRawAddress(ccv2.String()),
					testhelpers.MakeRawAddress(ccv3.String()),
					testhelpers.MakeRawAddress(ccv4.String()),
				},
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(testhelpers.MakeHashedAddress(executorAddress.String())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "too many CCVs")
		})
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()

		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()

			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{})
			require.NoError(t, err)
			client := makeClient(t, server)

			resp, err := client.PostExecutorSendWithBodyWithResponse(t.Context(), "0x1", "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()

			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{})
			require.NoError(t, err)
			client := makeClient(t, server)

			resp, err := client.PostExecutorSendWithBodyWithResponse(t.Context(), "0x1", "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid message", func(t *testing.T) {
			// Request with a message that fails message validation check
			t.Parallel()

			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{})
			require.NoError(t, err)
			client := makeClient(t, server)

			resp, err := client.PostExecutorSendWithResponse(t.Context(), "0x1", oapiExecutor.ExecutorSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: "0", // invalid selector
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type: oapiCommon.Empty,
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid message")
		})
		t.Run("Invalid executor path parameter", func(t *testing.T) {
			// Request with a message that fails message validation check
			t.Parallel()

			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{})
			require.NoError(t, err)
			client := makeClient(t, server)

			resp, err := client.PostExecutorSendWithResponse(t.Context(), "invalidaddress", oapiExecutor.ExecutorSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type: oapiCommon.Empty,
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid RawOrHashedAddress")
		})
		t.Run("Contract not found in store", func(t *testing.T) {
			t.Parallel()

			executorAddress := contracts.NewRawInstanceAddress("executor1", "owner")
			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{
				Executors: []config.Executor{
					{
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: executorAddress.InstanceAddress(),
						},
					},
				},
			})
			require.NoError(t, err)
			client := makeClient(t, server)

			mockActiveContractStore.EXPECT().Get(executorAddress.InstanceAddress()).Return(nil, false)

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.InstanceAddress().Hex(), oapiExecutor.ExecutorSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type: oapiCommon.Empty,
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		})
		t.Run("No contract returned from store", func(t *testing.T) {
			t.Parallel()

			executorAddress := contracts.NewRawInstanceAddress("executor1", "owner")
			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{
				Executors: []config.Executor{
					{
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: executorAddress.InstanceAddress(),
						},
					},
				},
			})
			require.NoError(t, err)
			client := makeClient(t, server)

			mockActiveContractStore.EXPECT().Get(executorAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contract1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: nil, // No CreateArguments returned
			}}, true)

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.InstanceAddress().Hex(), oapiExecutor.ExecutorSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type: oapiCommon.Empty,
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		})
		t.Run("Invalid contract returned from store", func(t *testing.T) {
			t.Parallel()

			executorAddress := contracts.NewRawInstanceAddress("executor1", "owner")
			mockActiveContractStore := mocks.NewMockActiveContractStore(t)
			mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
			server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.ExecutorAPIConfig{
				Executors: []config.Executor{
					{
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: executorAddress.InstanceAddress(),
						},
					},
				},
			})
			require.NoError(t, err)
			client := makeClient(t, server)

			mockActiveContractStore.EXPECT().Get(executorAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contract1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(executorBinding.Executor{
					DynamicConfig: executorBinding.DynamicConfig{
						AllowedFinalityConfig: core.FinalityConfig{WaitForFinality: &types.UNIT{}},
					},
					AllowedCCVs: []chainlinkapi.RawInstanceAddress{
						chainlinkapi.RawInstanceAddress{Unpack: "invalidaddress"},
					},
				}),
			}}, true)

			resp, err := client.PostExecutorSendWithResponse(t.Context(), executorAddress.InstanceAddress().Hex(), oapiExecutor.ExecutorSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: "123",
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type: oapiCommon.Empty,
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
				},
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		})
	})
}

func TestServer_FilterContracts(t *testing.T) {
	t.Parallel()

	cfg := config.ExecutorAPIConfig{
		Executors: []config.Executor{
			{ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x1")}},
			{ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x2")}},
			{ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x3")}},
		},
	}
	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything).Times(len(cfg.Executors))
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
	require.NoError(t, err)

	t.Run("All filters", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.Executors[0].InstanceAddress,
				cfg.Executors[1].InstanceAddress,
				cfg.Executors[2].InstanceAddress,
			},
		)
		assert.Len(t, res, 3)
		assert.Contains(t, res, cfg.Executors[0].InstanceAddress)
		assert.Contains(t, res, cfg.Executors[1].InstanceAddress)
		assert.Contains(t, res, cfg.Executors[2].InstanceAddress)
	})
	t.Run("Filters unknown addresses", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.Executors[0].InstanceAddress,
				contracts.HexToInstanceAddress("0x999"),
				cfg.Executors[2].InstanceAddress,
				contracts.HexToInstanceAddress("0x998"),
				contracts.HexToInstanceAddress("0x997"),
				cfg.Executors[1].InstanceAddress,
				contracts.HexToInstanceAddress("0x996"),
			},
		)
		assert.Len(t, res, 3)
		assert.Contains(t, res, cfg.Executors[0].InstanceAddress)
		assert.Contains(t, res, cfg.Executors[1].InstanceAddress)
		assert.Contains(t, res, cfg.Executors[2].InstanceAddress)
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
