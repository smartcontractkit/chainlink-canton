package global

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
)

func makeClient(t *testing.T, server *Server) oapiGlobal.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	oapiGlobal.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiGlobal.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("Fails with invalid batch size", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		_, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 0})
		require.Error(t, err)
	})
	t.Run("Succeeds with filters", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockFilter1 := mocks.NewMockInstanceAddressFilter(t)
		mockFilter2 := mocks.NewMockInstanceAddressFilter(t)
		mockFilter3 := mocks.NewMockInstanceAddressFilter(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 10}, mockFilter1, mockFilter2, mockFilter3)
		require.NoError(t, err)
		client := makeClient(t, server)

		address1 := contracts.HexToInstanceAddress("0x1")
		address2 := contracts.NewRawInstanceAddress("instanceId", "party")

		addresses := []contracts.InstanceAddress{
			address1,
			address2.InstanceAddress(),
			address1,
		}
		// address1 is returned by multiple filters
		mockFilter1.EXPECT().FilterContracts(addresses).Return(addresses).Once()
		mockFilter2.EXPECT().FilterContracts(addresses).Return([]contracts.InstanceAddress{address2.InstanceAddress()}).Once()
		mockFilter3.EXPECT().FilterContracts(addresses).Return(nil).Once()

		mockActiveContractStore.EXPECT().Get(address1).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contract1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package1",
				ModuleName: "module1",
				EntityName: "entity1",
			},
		}}, true)
		mockActiveContractStore.EXPECT().Get(address2.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contract2",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package2",
				ModuleName: "module2",
				EntityName: "entity2",
			},
		}}, true).Once()

		resp, err := client.PostGetExplicitDisclosureBatchWithResponse(t.Context(), oapiGlobal.GetExplicitDisclosureBatchRequest{
			Addresses: []oapiCommon.RawOrHashedAddress{
				testhelpers.MakeHashedAddress(address1.Hex()),
				testhelpers.MakeRawAddress(address2.String()),
				// Duplicate address, shouldn't fail
				testhelpers.MakeHashedAddress(address1.Hex()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())
		require.Equal(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "contract1",
				TemplateId: "package1:module1:entity1",
			},
			{
				ContractId: "contract2",
				TemplateId: "package2:module2:entity2",
			},
			{
				ContractId: "contract1",
				TemplateId: "package1:module1:entity1",
			},
		}, resp.JSON200.Disclosures)
	})
	t.Run("Succeeds on zero addresses", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockFilter1 := mocks.NewMockInstanceAddressFilter(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 10}, mockFilter1)
		require.NoError(t, err)
		client := makeClient(t, server)

		resp, err := client.PostGetExplicitDisclosureBatchWithResponse(t.Context(), oapiGlobal.GetExplicitDisclosureBatchRequest{
			Addresses: []oapiCommon.RawOrHashedAddress{},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())
		require.Empty(t, resp.JSON200.Disclosures)
	})
	t.Run("Fails on filters", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockFilter1 := mocks.NewMockInstanceAddressFilter(t)
		mockFilter2 := mocks.NewMockInstanceAddressFilter(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 10}, mockFilter1, mockFilter2)
		require.NoError(t, err)
		client := makeClient(t, server)

		address1 := contracts.HexToInstanceAddress("0x1")
		address2 := contracts.NewRawInstanceAddress("instanceId", "party")
		address3 := contracts.HexToInstanceAddress("0x3")

		addresses := []contracts.InstanceAddress{
			address1,
			address2.InstanceAddress(),
			address3,
		}
		// address1 & address3 do not match any filter
		mockFilter1.EXPECT().FilterContracts(addresses).Return(nil).Once()
		mockFilter2.EXPECT().FilterContracts(addresses).Return([]contracts.InstanceAddress{address2.InstanceAddress()}).Once()

		resp, err := client.PostGetExplicitDisclosureBatchWithResponse(t.Context(), oapiGlobal.GetExplicitDisclosureBatchRequest{
			Addresses: []oapiCommon.RawOrHashedAddress{
				testhelpers.MakeHashedAddress(address1.Hex()),
				testhelpers.MakeRawAddress(address2.String()),
				testhelpers.MakeHashedAddress(address3.Hex()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode())
		require.Contains(t, resp.JSON404.Error, address1.Hex())
		require.Contains(t, resp.JSON404.Error, address3.Hex())
	})
	t.Run("Fails on invalid addresses in request", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 10})
		require.NoError(t, err)
		client := makeClient(t, server)

		address1 := contracts.HexToInstanceAddress("this is invalid hex")
		address2 := contracts.NewRawInstanceAddress("valid", "address")

		resp, err := client.PostGetExplicitDisclosureBatchWithResponse(t.Context(), oapiGlobal.GetExplicitDisclosureBatchRequest{
			Addresses: []oapiCommon.RawOrHashedAddress{
				testhelpers.MakeHashedAddress(address1.Hex()),
				testhelpers.MakeRawAddress(address2.String()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
	t.Run("Fails on invalid request", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 10})
		require.NoError(t, err)
		client := makeClient(t, server)

		resp, err := client.PostGetExplicitDisclosureBatchWithBodyWithResponse(t.Context(), "application/json", strings.NewReader("invalid request"))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
	t.Run("Fails when batch size exceeded", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 1})
		require.NoError(t, err)
		client := makeClient(t, server)

		address1 := contracts.HexToInstanceAddress("0x1")
		address2 := contracts.NewRawInstanceAddress("instanceId", "party")

		resp, err := client.PostGetExplicitDisclosureBatchWithResponse(t.Context(), oapiGlobal.GetExplicitDisclosureBatchRequest{
			Addresses: []oapiCommon.RawOrHashedAddress{
				testhelpers.MakeHashedAddress(address1.Hex()),
				testhelpers.MakeRawAddress(address2.String()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
	t.Run("Fails on ActiveContractStore", func(t *testing.T) {
		t.Parallel()

		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockFilter1 := mocks.NewMockInstanceAddressFilter(t)

		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, config.GlobalAPIConfig{MaxBatchSize: 10}, mockFilter1)
		require.NoError(t, err)
		client := makeClient(t, server)

		address1 := contracts.HexToInstanceAddress("0x1")
		address2 := contracts.HexToInstanceAddress("0x2")

		addresses := []contracts.InstanceAddress{
			address1,
			address2,
		}
		// address1 is returned by multiple filters
		mockFilter1.EXPECT().FilterContracts(addresses).Return(addresses).Once()

		mockActiveContractStore.EXPECT().Get(address1).Return(nil, false).Once()

		resp, err := client.PostGetExplicitDisclosureBatchWithResponse(t.Context(), oapiGlobal.GetExplicitDisclosureBatchRequest{
			Addresses: []oapiCommon.RawOrHashedAddress{
				testhelpers.MakeHashedAddress(address1.Hex()),
				testhelpers.MakeHashedAddress(address2.Hex()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})
}
