package token_standard

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiTokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

const RequestSizeLimit = 100_000

func makeMetadataClient(t *testing.T, server *Server) oapiTokenMetadataV1.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	router.Use(middleware.RequestSizeLimiterMiddleware(RequestSizeLimit))
	oapiTokenMetadataV1.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiTokenMetadataV1.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func makeTransferInstructionClient(t *testing.T, server *Server) oapiTransferInstruction.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	router.Use(middleware.RequestSizeLimiterMiddleware(RequestSizeLimit))
	oapiTransferInstruction.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiTransferInstruction.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func makeRegistries(size int) map[string]config.Registry {
	registries := make(map[string]config.Registry, size)
	for i := range size {
		id := fmt.Sprintf("Token%v", i+1)
		name := fmt.Sprintf("Token Name %v", i+1)
		symbol := fmt.Sprintf("TKN%v", i+1)
		registries[id] = config.Registry{
			ContractIdentifier: config.ContractIdentifier{
				PartyID:         "owner",
				InstanceAddress: contracts.HexToInstanceAddress(fmt.Sprintf("0x%v", i+1)),
			},
			TokenType:   config.TokenTypeLINK,
			TokenId:     id,
			TokenName:   name,
			TokenSymbol: symbol,
		}
	}

	return registries
}

func TestServer_ListInstruments(t *testing.T) {
	t.Parallel()

	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
	mockActiveContractStore.EXPECT().GetByTemplateId(types.PARTY("admin"), contracts.TemplateIDFromBinding(link.LockedLinkHolding{})).
		Return([]*apiv2.ActiveContract{
			{
				CreatedEvent: &apiv2.CreatedEvent{
					CreateArguments: bindings.MarshalTemplateToRecord(link.LockedLinkHolding{
						LockedInstrumentId: splice_api_token_holding_v1.InstrumentId{
							Admin: "admin",
							Id:    "LINK", // Should be ignored, as it doesn't match any of the configured tokens
						},
						LockedAmount: "1.0",
					}),
				},
			},
		}, true)
	mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)
	mockInstrumentHoldingStore.EXPECT().RegisterParty([]string{"admin"})
	mockInstrumentHoldingStore.EXPECT().ListHoldings(mock.Anything).Return([]*apiv2.ActiveContract{
		{
			CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(link.LinkHolding{HoldingAmount: "100"}),
			},
		},
	}, nil)
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, mockInstrumentHoldingStore, config.TokenStandardAPIConfig{
		Enabled:    true,
		Admin:      "admin",
		Registries: makeRegistries(30),
	})
	require.NoError(t, err)
	client := makeMetadataClient(t, server)

	t.Run("Success - Default page size", func(t *testing.T) {
		t.Parallel()
		// 1st call
		resp, err := client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
			PageSize:  nil,
			PageToken: nil,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Len(t, resp.JSON200.Instruments, int(DefaultPageSize))
		require.NotNil(t, resp.JSON200.NextPageToken)

		// 2nd call
		resp, err = client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
			PageSize:  nil,
			PageToken: resp.JSON200.NextPageToken,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Len(t, resp.JSON200.Instruments, int(30-DefaultPageSize))
		require.Nil(t, resp.JSON200.NextPageToken)
	})
	t.Run("Success - Custom page size", func(t *testing.T) {
		t.Parallel()
		resp, err := client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
			PageSize:  new(int32(50)),
			PageToken: nil,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Len(t, resp.JSON200.Instruments, 30)
		require.Nil(t, resp.JSON200.NextPageToken)
		for _, instrument := range resp.JSON200.Instruments {
			// Validate that TotalSupply is calculated for each instrument (mocked to be 100, the LockedHolding must be ignored since it's for another instrument ID)
			require.NotNil(t, instrument.TotalSupply)
			require.Equal(t, "100.0000000000", *instrument.TotalSupply)
		}
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Zero page size", func(t *testing.T) {
			t.Parallel()
			resp, err := client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
				PageSize:  new(int32(0)),
				PageToken: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "pageSize must be a positive integer")
		})
		t.Run("Maximum page size", func(t *testing.T) {
			t.Parallel()
			resp, err := client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
				PageSize:  new(int32(1025)),
				PageToken: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "pageSize cannot be greater than")
		})
		t.Run("Invalid page token", func(t *testing.T) {
			t.Parallel()
			resp, err := client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
				PageToken: new("invalid"),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid pageToken")
		})
		t.Run("Invalid page token", func(t *testing.T) {
			t.Parallel()
			resp, err := client.ListInstrumentsWithResponse(t.Context(), &oapiTokenMetadataV1.ListInstrumentsParams{
				PageToken: new("50"),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid pageToken")
		})
	})
}

func TestServer_GetInstrument(t *testing.T) {
	t.Parallel()

	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
	mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)
	mockInstrumentHoldingStore.EXPECT().RegisterParty([]string{"admin"})
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, mockInstrumentHoldingStore, config.TokenStandardAPIConfig{
		Enabled: true,
		Admin:   "admin",
		Registries: map[string]config.Registry{
			"LINK": {
				ContractIdentifier: config.ContractIdentifier{
					PartyID:         "admin",
					InstanceAddress: contracts.HexToInstanceAddress("0x1"),
				},
				TokenType: config.TokenTypeLINK,
				TokenId:   "LINK",
			},
			"link2": {
				ContractIdentifier: config.ContractIdentifier{
					PartyID:         "admin",
					InstanceAddress: contracts.HexToInstanceAddress("0x2"),
				},
				TokenType:   config.TokenTypeLINK,
				TokenId:     "link2",
				TokenName:   "Link Token 2",
				TokenSymbol: "LINK2",
			},
		},
		SupplyCacheTTL: time.Second,
	})
	require.NoError(t, err)
	client := makeMetadataClient(t, server)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockInstrumentHoldingStore.EXPECT().ListHoldings(splice_api_token_holding_v1.InstrumentId{
			Admin: "admin",
			Id:    "LINK",
		}).Return([]*apiv2.ActiveContract{
			{
				CreatedEvent: &apiv2.CreatedEvent{
					CreateArguments: bindings.MarshalTemplateToRecord(link.LinkHolding{HoldingAmount: "100"}),
				},
			}, {
				CreatedEvent: &apiv2.CreatedEvent{
					CreateArguments: bindings.MarshalTemplateToRecord(link.LinkHolding{HoldingAmount: "23.456"}),
				},
			},
		}, nil).Once()
		mockActiveContractStore.EXPECT().GetByTemplateId(types.PARTY("admin"), contracts.TemplateIDFromBinding(link.LockedLinkHolding{})).
			Return([]*apiv2.ActiveContract{
				{
					CreatedEvent: &apiv2.CreatedEvent{
						CreateArguments: bindings.MarshalTemplateToRecord(link.LockedLinkHolding{
							LockedInstrumentId: splice_api_token_holding_v1.InstrumentId{
								Admin: "admin",
								Id:    "LINK",
							},
							LockedAmount: "0.000789",
						}),
					},
				},
			}, true).Once()
		resp, err := client.GetInstrumentWithResponse(t.Context(), "LINK")
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "LINK", resp.JSON200.Id)
		require.Equal(t, "LINK", resp.JSON200.Name)
		require.Equal(t, "LINK", resp.JSON200.Symbol)
		require.Equal(t, int8(10), resp.JSON200.Decimals)
		// Validate that the TotalSupply is calculated correctly
		require.NotNil(t, resp.JSON200.TotalSupply)
		require.Equal(t, "123.4567890000", *resp.JSON200.TotalSupply)
		require.NotNil(t, resp.JSON200.TotalSupplyAsOf)
		require.WithinDuration(t, time.Now(), *resp.JSON200.TotalSupplyAsOf, time.Second*5)

		// Validate that calling the same endpoint again doesn't cause the totalSupply to be re-calculated
		_, err = client.GetInstrumentWithResponse(t.Context(), "LINK")
		require.NoError(t, err)

		mockInstrumentHoldingStore.EXPECT().ListHoldings(splice_api_token_holding_v1.InstrumentId{
			Admin: "admin",
			Id:    "link2",
		}).Return([]*apiv2.ActiveContract{}, nil)
		mockActiveContractStore.EXPECT().GetByTemplateId(types.PARTY("admin"), contracts.TemplateIDFromBinding(link.LockedLinkHolding{})).
			Return([]*apiv2.ActiveContract{}, true)
		resp, err = client.GetInstrumentWithResponse(t.Context(), "link2")
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "link2", resp.JSON200.Id)
		require.Equal(t, "Link Token 2", resp.JSON200.Name)
		require.Equal(t, "LINK2", resp.JSON200.Symbol)
		require.Equal(t, int8(10), resp.JSON200.Decimals)
	})
	t.Run("Not found", func(t *testing.T) {
		t.Parallel()
		resp, err := client.GetInstrumentWithResponse(t.Context(), "USDC")
		require.NoError(t, err)
		require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
	})
}

func TestServer_GetTransferFactory(t *testing.T) {
	t.Parallel()

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiTransferInstruction.ClientWithResponsesInterface) {
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)
		mockInstrumentHoldingStore.EXPECT().RegisterParty([]string{"admin"})
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, mockInstrumentHoldingStore, config.TokenStandardAPIConfig{
			Enabled: true,
			Admin:   "admin",
			Registries: map[string]config.Registry{
				"LINK": {
					ContractIdentifier: config.ContractIdentifier{
						PartyID:         "admin",
						InstanceAddress: contracts.HexToInstanceAddress("0x1"),
					},
					TokenType: config.TokenTypeLINK,
					TokenId:   "LINK",
				},
			},
		})
		require.NoError(t, err)
		client := makeTransferInstructionClient(t, server)

		return mockActiveContractStore, client
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)
		mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x1")).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contractId1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "LinkToken",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(link.LinkRegistry{
				RegistryAdmin: "admin",
				RegistryInstrumentId: splice_api_token_holding_v1.InstrumentId{
					Admin: "admin",
					Id:    "LINK",
				},
				InstanceId: "linktoken",
			}),
		}}, true)
		mockActiveContractStore.EXPECT().GetByTemplateId(types.PARTY("admin"), contracts.TemplateIDFromBinding(link.LinkTransferPreapproval{})).Return([]*apiv2.ActiveContract{}, false)

		resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
			ChoiceArguments: map[string]any{
				"expectedAdmin": "admin",
				"transfer": map[string]any{
					"sender":   "sender",
					"receiver": "receiver",
					"amount":   "100",
					"instrumentId": map[string]any{
						"admin": "admin",
						"id":    "LINK",
					},
					"requestedAt":      time.Now().Add(time.Hour * -1).Format(time.RFC3339),
					"executeBefore":    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
					"inputHoldingCids": []string{""},
					"meta": map[string]any{
						"values": map[string]any{},
					},
				},
				"extraArgs": map[string]any{
					"context": map[string]any{
						"values": map[string]any{},
					},
					"meta": map[string]any{
						"values": map[string]any{},
					},
				},
			},
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, []oapiTransferInstruction.DisclosedContract{
			{
				ContractId: "contractId1",
				TemplateId: "package:module:LinkToken",
			},
		}, resp.JSON200.ChoiceContext.DisclosedContracts)
		expectedChoiceContext, err := converters.SerializeChoiceContext(splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}})
		require.NoError(t, err)
		require.Equal(t, expectedChoiceContext, resp.JSON200.ChoiceContext.ChoiceContextData)
		require.Equal(t, "contractId1", resp.JSON200.FactoryId)
	})
	t.Run("Success - with receiver TransferPreapproval", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)
		mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x1")).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contractId1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "LinkToken",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(link.LinkRegistry{
				RegistryAdmin: "admin",
				RegistryInstrumentId: splice_api_token_holding_v1.InstrumentId{
					Admin: "admin",
					Id:    "LINK",
				},
				InstanceId: "linktoken",
			}),
		}}, true)
		mockActiveContractStore.EXPECT().GetByTemplateId(types.PARTY("admin"), contracts.TemplateIDFromBinding(link.LinkTransferPreapproval{})).Return([]*apiv2.ActiveContract{
			{
				CreatedEvent: &apiv2.CreatedEvent{
					ContractId: "senderPreapprovalContractId",
					CreateArguments: bindings.MarshalTemplateToRecord(link.LinkTransferPreapproval{
						PreapprovalAdmin:    "admin",
						PreapprovalReceiver: "sender",
					}),
				},
			}, {
				CreatedEvent: &apiv2.CreatedEvent{
					ContractId: "receiverPreapprovalContractId",
					TemplateId: &apiv2.Identifier{
						PackageId:  "package",
						ModuleName: "module",
						EntityName: "LinkTransferPreapproval",
					},
					CreateArguments: bindings.MarshalTemplateToRecord(link.LinkTransferPreapproval{
						PreapprovalAdmin:    "admin",
						PreapprovalReceiver: "receiver",
					}),
				},
			},
		}, false)

		resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
			ChoiceArguments: map[string]any{
				"expectedAdmin": "admin",
				"transfer": map[string]any{
					"sender":   "sender",
					"receiver": "receiver",
					"amount":   "100",
					"instrumentId": map[string]any{
						"admin": "admin",
						"id":    "LINK",
					},
					"requestedAt":      time.Now().Add(time.Hour * -1).Format(time.RFC3339),
					"executeBefore":    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
					"inputHoldingCids": []string{""},
					"meta": map[string]any{
						"values": map[string]any{},
					},
				},
				"extraArgs": map[string]any{
					"context": map[string]any{
						"values": map[string]any{},
					},
					"meta": map[string]any{
						"values": map[string]any{},
					},
				},
			},
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, []oapiTransferInstruction.DisclosedContract{
			{
				ContractId: "contractId1",
				TemplateId: "package:module:LinkToken",
			}, {
				ContractId: "receiverPreapprovalContractId",
				TemplateId: "package:module:LinkTransferPreapproval",
			},
		}, resp.JSON200.ChoiceContext.DisclosedContracts)
		// Check that the preapproval is picked up correctly
		expectedChoiceContext, err := converters.SerializeChoiceContext(splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{
			string(link.TransferPreapprovalContextKey): {
				AVContractId: new(types.CONTRACT_ID("receiverPreapprovalContractId")),
			},
		}})
		require.NoError(t, err)
		require.Equal(t, expectedChoiceContext, resp.JSON200.ChoiceContext.ChoiceContextData)
		require.Equal(t, "contractId1", resp.JSON200.FactoryId)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferFactoryWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferFactoryWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid choice arguments", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
				ChoiceArguments: map[string]any{
					"transfer": "invalidvalue",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "choiceArguments")
		})
		t.Run("Mismatched expectedAdmin", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
				ChoiceArguments: map[string]any{
					"expectedAdmin": "wrongParty",
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "expectedAdmin")
		})
		t.Run("Mismatched expectedAdmin", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
				ChoiceArguments: map[string]any{
					"expectedAdmin": "admin",
					"transfer": map[string]any{
						"instrumentId": map[string]any{
							"admin": "wrongparty",
							"id":    "LINK",
						},
					},
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "transfer.instrumentId.admin")
		})
		t.Run("Unknown instrument", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
				ChoiceArguments: map[string]any{
					"expectedAdmin": "admin",
					"transfer": map[string]any{
						"instrumentId": map[string]any{
							"admin": "admin",
							"id":    "USDC",
						},
					},
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "No instrument with id USDC found")
		})
		t.Run("Registry not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)
			mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x1")).Return(nil, false)

			resp, err := client.GetTransferFactoryWithResponse(t.Context(), oapiTransferInstruction.GetFactoryRequest{
				ChoiceArguments: map[string]any{
					"expectedAdmin": "admin",
					"transfer": map[string]any{
						"instrumentId": map[string]any{
							"admin": "admin",
							"id":    "LINK",
						},
					},
				},
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_GetTransferInstructionAcceptContext(t *testing.T) {
	t.Parallel()

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiTransferInstruction.ClientWithResponsesInterface) {
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)
		mockInstrumentHoldingStore.EXPECT().RegisterParty([]string{"admin"})
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, mockInstrumentHoldingStore, config.TokenStandardAPIConfig{
			Enabled: true,
			Admin:   "admin",
			Registries: map[string]config.Registry{
				"LINK": {
					ContractIdentifier: config.ContractIdentifier{
						PartyID:         "admin",
						InstanceAddress: contracts.HexToInstanceAddress("0x1"),
					},
					TokenType: config.TokenTypeLINK,
					TokenId:   "LINK",
				},
			},
		})
		require.NoError(t, err)
		client := makeTransferInstructionClient(t, server)

		return mockActiveContractStore, client
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)
		cid := types.CONTRACT_ID("transferInstructionContractId")
		mockActiveContractStore.EXPECT().GetByContractId(cid).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: string(cid),
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "LinkTransferInstruction",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(link.LinkTransferInstruction{
				InstructionAdmin: "admin",
				InstructionTransfer: splice_api_token_transfer_instruction_v1.Transfer{
					InstrumentId: splice_api_token_holding_v1.InstrumentId{
						Admin: "admin",
						Id:    "LINK",
					},
					Amount: "1",
				},
			}),
		}}, true)
		mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x1")).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contractId1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "LinkToken",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(link.LinkRegistry{
				RegistryAdmin: "admin",
				RegistryInstrumentId: splice_api_token_holding_v1.InstrumentId{
					Admin: "admin",
					Id:    "LINK",
				},
				InstanceId: "linktoken",
			}),
		}}, true)

		resp, err := client.GetTransferInstructionAcceptContextWithResponse(t.Context(), string(cid), oapiTransferInstruction.GetChoiceContextRequest{
			Meta: nil,
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, []oapiTransferInstruction.DisclosedContract{
			{
				ContractId: "contractId1",
				TemplateId: "package:module:LinkToken",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferInstructionAcceptContextWithBodyWithResponse(t.Context(), "abcdefg", "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTransferInstructionAcceptContextWithBodyWithResponse(t.Context(), "abcdefg", "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("TransferInstruction not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)
			cid := types.CONTRACT_ID("abcdefg")
			mockActiveContractStore.EXPECT().GetByContractId(cid).Return(nil, false)

			resp, err := client.GetTransferInstructionAcceptContextWithResponse(t.Context(), string(cid), oapiTransferInstruction.GetChoiceContextRequest{
				Meta: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "No transfer instruction with id abcdefg found")
		})
		t.Run("TransferInstruction not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)
			cid := types.CONTRACT_ID("abcdefg")
			mockActiveContractStore.EXPECT().GetByContractId(cid).Return(nil, true)

			resp, err := client.GetTransferInstructionAcceptContextWithResponse(t.Context(), string(cid), oapiTransferInstruction.GetChoiceContextRequest{
				Meta: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "No transfer instruction with id abcdefg found")
		})
		t.Run("Instrument admin mismatch", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)
			cid := types.CONTRACT_ID("abcdefg")
			mockActiveContractStore.EXPECT().GetByContractId(cid).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: string(cid),
				TemplateId: &apiv2.Identifier{
					PackageId:  "package",
					ModuleName: "module",
					EntityName: "LinkTransferInstruction",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(link.LinkTransferInstruction{
					InstructionAdmin: "anotherAdmin",
					InstructionTransfer: splice_api_token_transfer_instruction_v1.Transfer{
						InstrumentId: splice_api_token_holding_v1.InstrumentId{
							Admin: "anotherAdmin",
							Id:    "LINK",
						},
						Amount: "1",
					},
				}),
			}}, true)

			resp, err := client.GetTransferInstructionAcceptContextWithResponse(t.Context(), string(cid), oapiTransferInstruction.GetChoiceContextRequest{
				Meta: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "No transfer instruction with id abcdefg found")
		})
		t.Run("Unknown instrument", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)
			cid := types.CONTRACT_ID("abcdefg")
			mockActiveContractStore.EXPECT().GetByContractId(cid).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: string(cid),
				TemplateId: &apiv2.Identifier{
					PackageId:  "package",
					ModuleName: "module",
					EntityName: "LinkTransferInstruction",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(link.LinkTransferInstruction{
					InstructionAdmin: "admin",
					InstructionTransfer: splice_api_token_transfer_instruction_v1.Transfer{
						InstrumentId: splice_api_token_holding_v1.InstrumentId{
							Admin: "admin",
							Id:    "USDC",
						},
						Amount: "1",
					},
				}),
			}}, true)

			resp, err := client.GetTransferInstructionAcceptContextWithResponse(t.Context(), string(cid), oapiTransferInstruction.GetChoiceContextRequest{
				Meta: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusNotFound, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "No transfer instruction with id abcdefg found")
		})
		t.Run("Registry not found", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)
			cid := types.CONTRACT_ID("transferInstructionContractId")
			mockActiveContractStore.EXPECT().GetByContractId(cid).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: string(cid),
				TemplateId: &apiv2.Identifier{
					PackageId:  "package",
					ModuleName: "module",
					EntityName: "LinkTransferInstruction",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(link.LinkTransferInstruction{
					InstructionAdmin: "admin",
					InstructionTransfer: splice_api_token_transfer_instruction_v1.Transfer{
						InstrumentId: splice_api_token_holding_v1.InstrumentId{
							Admin: "admin",
							Id:    "LINK",
						},
						Amount: "1",
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x1")).Return(nil, false)

			resp, err := client.GetTransferInstructionAcceptContextWithResponse(t.Context(), string(cid), oapiTransferInstruction.GetChoiceContextRequest{
				Meta: nil,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_FilterContracts(t *testing.T) {
	t.Parallel()

	cfg := config.TokenStandardAPIConfig{
		Admin: "admin",
		Registries: map[string]config.Registry{
			"Token1": {TokenId: "1", TokenType: config.TokenTypeLINK, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x1")}},
			"Token2": {TokenId: "2", TokenType: config.TokenTypeLINK, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x2")}},
			"Token3": {TokenId: "3", TokenType: config.TokenTypeLINK, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x3")}},
		},
	}
	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything).Times(len(cfg.Registries))
	mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)
	mockInstrumentHoldingStore.EXPECT().RegisterParty([]string{cfg.Admin})
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, mockInstrumentHoldingStore, cfg)
	require.NoError(t, err)

	t.Run("All filters", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.Registries["Token1"].InstanceAddress,
				cfg.Registries["Token2"].InstanceAddress,
				cfg.Registries["Token3"].InstanceAddress,
			},
		)
		assert.Len(t, res, 3)
		assert.Contains(t, res, cfg.Registries["Token1"].InstanceAddress)
		assert.Contains(t, res, cfg.Registries["Token2"].InstanceAddress)
		assert.Contains(t, res, cfg.Registries["Token3"].InstanceAddress)
	})
	t.Run("Filters unknown addresses", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.Registries["Token1"].InstanceAddress,
				contracts.HexToInstanceAddress("0x999"),
				cfg.Registries["Token3"].InstanceAddress,
				contracts.HexToInstanceAddress("0x998"),
				contracts.HexToInstanceAddress("0x997"),
				cfg.Registries["Token2"].InstanceAddress,
				contracts.HexToInstanceAddress("0x996"),
			},
		)
		assert.Len(t, res, 3)
		assert.Contains(t, res, cfg.Registries["Token1"].InstanceAddress)
		assert.Contains(t, res, cfg.Registries["Token2"].InstanceAddress)
		assert.Contains(t, res, cfg.Registries["Token3"].InstanceAddress)
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
