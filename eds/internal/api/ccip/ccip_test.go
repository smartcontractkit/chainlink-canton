package ccip

import (
	"encoding/hex"
	"math/big"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/testhelpers"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

const RequestSizeLimit = 100_000

func makeClient(t *testing.T, server *Server) oapiCCIP.ClientWithResponsesInterface {
	t.Helper()

	router := gin.Default()
	router.Use(middleware.RequestSizeLimiterMiddleware(RequestSizeLimit))
	oapiCCIP.RegisterHandlers(router, server)
	s := httptest.NewServer(router)
	client, err := oapiCCIP.NewClientWithResponses(s.URL)
	require.NoError(t, err)

	return client
}

func TestServer_PostPerPartyRouterFactory(t *testing.T) {
	t.Parallel()

	perPartyRouterFactoryAddress := contracts.NewRawInstanceAddress("factory", "owner")
	cfg := config.CCIPAPIConfig{
		PerPartyRouterFactory: config.ContractIdentifier{
			InstanceAddress: perPartyRouterFactoryAddress.InstanceAddress(),
		},
	}

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiCCIP.ClientWithResponsesInterface) {
		t.Helper()
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockActiveContractStore, client
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		perPartyRouterFactory := ccipruntime.PerPartyRouterFactory{
			InstanceId: types.TEXT(perPartyRouterFactoryAddress.InstanceID()),
			CcipOwner:  types.PARTY(perPartyRouterFactoryAddress.Owner()),
		}
		mockActiveContractStore.EXPECT().Get(perPartyRouterFactoryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contractId1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package1",
				ModuleName: "module1",
				EntityName: "entity1",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(perPartyRouterFactory),
		}}, true)

		resp, err := client.PostPerPartyRouterFactoryWithResponse(t.Context(), oapiCCIP.CCIPPerPartyRouterFactoryRequest{
			PartyID: "userParty",
		})
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, "contractId1", resp.JSON200.ContractId)
		require.Equal(t, perPartyRouterFactoryAddress.InstanceAddress().Hex(), resp.JSON200.InstanceAddress)
		require.Equal(t, perPartyRouterFactoryAddress.String(), resp.JSON200.RawInstanceAddress)
		require.Equal(t, []oapiCommon.DisclosedContract{
			{
				ContractId: "contractId1",
				TemplateId: "package1:module1:entity1",
			},
		}, resp.JSON200.DisclosedContracts)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostPerPartyRouterFactoryWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostPerPartyRouterFactoryWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Contract not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(perPartyRouterFactoryAddress.InstanceAddress()).Return(nil, false)
			resp, err := client.PostPerPartyRouterFactoryWithResponse(t.Context(), oapiCCIP.CCIPPerPartyRouterFactoryRequest{})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("No contract returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(perPartyRouterFactoryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: nil, // No CreateArguments returned
			}}, true)
			resp, err := client.PostPerPartyRouterFactoryWithResponse(t.Context(), oapiCCIP.CCIPPerPartyRouterFactoryRequest{})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
	})
}

func TestServer_GetTokenAdminRegistryToken(t *testing.T) {
	t.Parallel()

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenAdmin",
		Id:    "ChainLink",
	}
	encodedInstrumentId := contracts.EncodeInstrumentID(instrumentId)
	tokenPoolAddress := contracts.NewRawInstanceAddress("poolId", "poolOwner")
	tokenAdminRegistryAddress := contracts.NewRawInstanceAddress("factory", "owner")
	tokenConfigAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(hex.EncodeToString(encodedInstrumentId.Bytes())), types.PARTY(tokenAdminRegistryAddress.Owner()))
	cfg := config.CCIPAPIConfig{
		TokenAdminRegistry: config.ContractIdentifier{
			InstanceAddress: tokenAdminRegistryAddress.InstanceAddress(),
		},
	}

	tokenAdminRegistry := core.TokenAdminRegistry{
		InstanceId: types.TEXT(tokenAdminRegistryAddress.InstanceID()),
		CcipOwner:  types.PARTY(tokenAdminRegistryAddress.Owner()),
		EntryCount: 2,
	}
	tokenConfig := core.TokenConfig{
		InstanceId:         types.TEXT(hex.EncodeToString(encodedInstrumentId.Bytes())),
		RegistryInstanceId: types.TEXT(tokenAdminRegistryAddress.InstanceID()),
		RegistryOwner:      types.PARTY(tokenAdminRegistryAddress.Owner()),
		Index:              1,
		IsCCIPManaged:      false,
		InstrumentId:       instrumentId,
		TokenPool: &core.PoolRegistration2{
			PoolOwner:      types.PARTY(tokenPoolAddress.Owner()),
			PoolInstanceId: types.TEXT(tokenPoolAddress.InstanceID()),
		},
	}

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiCCIP.ClientWithResponsesInterface) {
		t.Helper()
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockActiveContractStore, client
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(tokenAdminRegistryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contractId1",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package1",
				ModuleName: "module1",
				EntityName: "entity1",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(tokenAdminRegistry),
		}}, true)
		mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "contractId2",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package2",
				ModuleName: "module2",
				EntityName: "entity2",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(tokenConfig),
		}}, true).Once()

		resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), encodedInstrumentId.String())
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		require.Equal(t, tokenPoolAddress.String(), resp.JSON200.RawInstanceAddress)
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		t.Run("Invalid InstrumentId", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), "invalidInstrumentId")
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid instrumentId")
		})
		t.Run("TokenAdminRegistry not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(tokenAdminRegistryAddress.InstanceAddress()).Return(nil, false)

			resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), encodedInstrumentId.String())
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenAdminRegistry not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(tokenAdminRegistryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: nil, // No CreateArguments returned
			}}, true)

			resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), encodedInstrumentId.String())
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenConfig not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(tokenAdminRegistryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(tokenAdminRegistry),
			}}, true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(nil, false).Once()

			resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), encodedInstrumentId.String())
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no token config registered for token")
		})
		t.Run("TokenConfig not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(tokenAdminRegistryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(tokenAdminRegistry),
			}}, true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId2",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package2",
					ModuleName: "module2",
					EntityName: "entity2",
				},
				CreateArguments: nil, // No CreateArguments returned
			}}, true).Once()

			resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), encodedInstrumentId.String())
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("No TokenPool set in TokenConfig", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(tokenAdminRegistryAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId1",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package1",
					ModuleName: "module1",
					EntityName: "entity1",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(tokenAdminRegistry),
			}}, true)
			tokenConfig := tokenConfig
			tokenConfig.TokenPool = nil
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: "contractId2",
				TemplateId: &apiv2.Identifier{
					PackageId:  "package2",
					ModuleName: "module2",
					EntityName: "entity2",
				},
				CreateArguments: bindings.MarshalTemplateToRecord(tokenConfig),
			}}, true).Once()

			resp, err := client.GetTokenAdminRegistryTokenWithResponse(t.Context(), encodedInstrumentId.String())
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no token pool registered for token")
		})
	})
}

func TestServer_PostCCIPSend(t *testing.T) {
	t.Parallel()

	sourceSelector := "123"
	destSelector := "456"
	perPartyRouterFactoryAddress := contracts.NewRawInstanceAddress("perPartyRouterFactory", "owner")
	onRampAddress := contracts.NewRawInstanceAddress("onRamp", "owner")
	offRampAddress := contracts.NewRawInstanceAddress("offRamp", "owner")
	globalConfigAddress := contracts.NewRawInstanceAddress("globalConfig", "owner")
	tokenAdminRegistryAddress := contracts.NewRawInstanceAddress("tokenAdminRegistry", "owner")
	rmnRemoteAddress := contracts.NewRawInstanceAddress("rmnRemote", "owner")
	feeQuoterAddress := contracts.NewRawInstanceAddress("feeQuoter", "owner")
	tokenPoolAddress := contracts.NewRawInstanceAddress("tokenPool", "owner")
	cfg := config.CCIPAPIConfig{
		PerPartyRouterFactory: config.ContractIdentifier{InstanceAddress: perPartyRouterFactoryAddress.InstanceAddress()},
		OnRamp:                config.ContractIdentifier{InstanceAddress: onRampAddress.InstanceAddress()},
		OffRamp:               config.ContractIdentifier{InstanceAddress: offRampAddress.InstanceAddress()},
		GlobalConfig:          config.ContractIdentifier{InstanceAddress: globalConfigAddress.InstanceAddress()},
		TokenAdminRegistry:    config.ContractIdentifier{InstanceAddress: tokenAdminRegistryAddress.InstanceAddress()},
		RMNRemote:             config.ContractIdentifier{InstanceAddress: rmnRemoteAddress.InstanceAddress()},
		FeeQuoter:             config.ContractIdentifier{InstanceAddress: feeQuoterAddress.InstanceAddress()},
	}
	ccv1 := contracts.NewRawInstanceAddress("ccv1", "owner") // lane-mandated
	ccv2 := contracts.NewRawInstanceAddress("ccv2", "owner")
	ccv3 := contracts.NewRawInstanceAddress("ccv3", "owner")
	ccv4 := contracts.NewRawInstanceAddress("ccv4", "owner")
	ccv5 := contracts.NewRawInstanceAddress("ccv5", "owner") // default
	ccv6 := contracts.NewRawInstanceAddress("ccv6", "owner") // lane-mandated + default
	defaultExecutor := contracts.NewRawInstanceAddress("executor", "owner")
	// Token Transfers only
	tokenInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenAdmin",
		Id:    "LinkToken",
	}
	encodedTokenInstrumentId := contracts.EncodeInstrumentID(tokenInstrumentId)
	tokenConfigAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(hex.EncodeToString(encodedTokenInstrumentId.Bytes())), "owner")
	feeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "feeAdmin",
		Id:    "LINK",
	}
	encodedFeeInstrumentId := contracts.EncodeInstrumentID(feeInstrumentId)
	feeTokenConfigAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(hex.EncodeToString(encodedFeeInstrumentId.Bytes())), "owner")

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiCCIP.ClientWithResponsesInterface) {
		t.Helper()
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockActiveContractStore, client
	}
	tokenConfigContract := func(contractId string, instrumentId splice_api_token_holding_v1.InstrumentId, tokenPool *core.PoolRegistration2) *apiv2.ActiveContract {
		encodedInstrumentId := contracts.EncodeInstrumentID(instrumentId)
		return &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: contractId,
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "TokenConfig",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(core.TokenConfig{
				InstanceId:         types.TEXT(hex.EncodeToString(encodedInstrumentId.Bytes())),
				RegistryInstanceId: types.TEXT(tokenAdminRegistryAddress.InstanceID()),
				RegistryOwner:      types.PARTY(tokenAdminRegistryAddress.Owner()),
				Index:              0,
				InstrumentId:       instrumentId,
				TokenPool:          tokenPool,
			}),
		}}
	}
	tokenPoolRegistration := &core.PoolRegistration2{
		PoolOwner:      types.PARTY(tokenPoolAddress.Owner()),
		PoolInstanceId: types.TEXT(tokenPoolAddress.InstanceID()),
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "onRampContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "OnRamp",
			},
		}}, true)
		mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "globalConfigContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "GlobalConfig",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
				InstanceId:    "globalconfig",
				CcipOwner:     "owner",
				ChainSelector: types.NUMERIC(sourceSelector),
				DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
					types.NUMERIC(destSelector): {
						IsEnabled:       true,
						DefaultExecutor: new(defaultExecutor.Binding()),
						LaneMandatedCCVs: []chainlinkapi.RawInstanceAddress{
							ccv1.Binding(),
							ccv6.Binding(),
						},
						DefaultCCVs: []chainlinkapi.RawInstanceAddress{
							ccv5.Binding(),
							ccv6.Binding(),
						},
						MessageNetworkFeeUSDCents: "0",
						TokenNetworkFeeUSDCents:   "0",
					},
				},
				SourceChainConfigs: nil,
			}),
		}}, true)
		mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "tokenAdminRegistryContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "TokenAdminRegistry",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
				InstanceId: "tokenadminregistry",
				CcipOwner:  "owner",
				EntryCount: 1,
			}),
		}}, true)
		mockActiveContractStore.EXPECT().Get(feeTokenConfigAddress.InstanceAddress()).Return(tokenConfigContract("feeTokenConfigContractId", feeInstrumentId, nil), true)
		mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(tokenConfigContract("tokenConfigContractId", tokenInstrumentId, tokenPoolRegistration), true)
		mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "rmnRemoteContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "RMNRemote",
			},
		}}, true)
		mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "feeQuoterContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "FeeQuoter",
			},
		}}, true)

		t.Run("Message only", func(t *testing.T) {
			t.Parallel()
			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: destSelector,
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
				SenderRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv1.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv4),
				}),
				TokenPoolRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),
				}),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.RawOrHashedAddress{
				converters.InstanceAddressAsRawOrHashedAddress(ccv1.InstanceAddress()), // SenderRequiredCCVs + LaneMandatedCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv4),                // SenderRequiredCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv6),                // LaneMandatedCCVs
				// No DefaultCCVs / TokenPoolRequiredCCVs
			}, resp.JSON200.Ccvs)
			require.ElementsMatch(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "onRampContractId",
					TemplateId: "package:module:OnRamp",
				}, {
					ContractId: "globalConfigContractId",
					TemplateId: "package:module:GlobalConfig",
				}, {
					ContractId: "tokenAdminRegistryContractId",
					TemplateId: "package:module:TokenAdminRegistry",
				}, {
					ContractId: "rmnRemoteContractId",
					TemplateId: "package:module:RMNRemote",
				}, {
					ContractId: "feeQuoterContractId",
					TemplateId: "package:module:FeeQuoter",
				}, {
					ContractId: "feeTokenConfigContractId",
					TemplateId: "package:module:TokenConfig",
				},
			}, resp.JSON200.DisclosedContracts)
			require.Equal(t, new(converters.RawInstanceAddressAsRawOrHashedAddress(defaultExecutor)), resp.JSON200.Executor)
			require.Equal(t, "feeTokenConfigContractId", resp.JSON200.FeeTokenConfigCid)
		})
		t.Run("Token Transfer", func(t *testing.T) {
			t.Parallel()
			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: destSelector,
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
					TokenTransfer: &oapiCommon.TokenTransfer{
						Amount: "42",
						Token: oapiCommon.InstrumentId{
							Admin: oapiCommon.PartyId(tokenInstrumentId.Admin),
							Id:    string(tokenInstrumentId.Id),
						},
					},
				},
				SenderRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv1.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv4),
				}),
				TokenPoolRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),
				}),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.RawOrHashedAddress{
				converters.InstanceAddressAsRawOrHashedAddress(ccv1.InstanceAddress()), // SenderRequiredCCVs + LaneMandatedCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv4),                // SenderRequiredCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv6),                // LaneMandatedCCVs
				converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()), // TokenPoolRequiredCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),                // TokenPoolRequiredCCVs
			}, resp.JSON200.Ccvs)
			require.ElementsMatch(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "onRampContractId",
					TemplateId: "package:module:OnRamp",
				}, {
					ContractId: "globalConfigContractId",
					TemplateId: "package:module:GlobalConfig",
				}, {
					ContractId: "tokenAdminRegistryContractId",
					TemplateId: "package:module:TokenAdminRegistry",
				}, {
					ContractId: "rmnRemoteContractId",
					TemplateId: "package:module:RMNRemote",
				}, {
					ContractId: "feeQuoterContractId",
					TemplateId: "package:module:FeeQuoter",
				}, {
					ContractId: "feeTokenConfigContractId",
					TemplateId: "package:module:TokenConfig",
				}, {
					ContractId: "tokenConfigContractId",
					TemplateId: "package:module:TokenConfig",
				},
			}, resp.JSON200.DisclosedContracts)
			require.Equal(t, new(converters.RawInstanceAddressAsRawOrHashedAddress(defaultExecutor)), resp.JSON200.Executor)
			require.Equal(t, "feeTokenConfigContractId", resp.JSON200.FeeTokenConfigCid)
		})
		t.Run("Pure Token Transfer", func(t *testing.T) {
			t.Parallel()
			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: destSelector,
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
					Payload:  "", // No payload
					GasLimit: 0,  // No gas limit -> should lead to sender-required CCVs to not be included
					Receiver: "0x1234567890",
					TokenTransfer: &oapiCommon.TokenTransfer{
						Amount: "42",
						Token: oapiCommon.InstrumentId{
							Admin: oapiCommon.PartyId(tokenInstrumentId.Admin),
							Id:    string(tokenInstrumentId.Id),
						},
					},
				},
				SenderRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv1.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv4),
				}),
				TokenPoolRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),
				}),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.RawOrHashedAddress{
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv1),                // LaneMandatedCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv6),                // LaneMandatedCCVs
				converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()), // TokenPoolRequiredCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),                // TokenPoolRequiredCCVs
				// SenderRequiredCCVs must not be included as this is a pure token transfer
			}, resp.JSON200.Ccvs)
			require.ElementsMatch(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "onRampContractId",
					TemplateId: "package:module:OnRamp",
				}, {
					ContractId: "globalConfigContractId",
					TemplateId: "package:module:GlobalConfig",
				}, {
					ContractId: "tokenAdminRegistryContractId",
					TemplateId: "package:module:TokenAdminRegistry",
				}, {
					ContractId: "rmnRemoteContractId",
					TemplateId: "package:module:RMNRemote",
				}, {
					ContractId: "feeQuoterContractId",
					TemplateId: "package:module:FeeQuoter",
				}, {
					ContractId: "feeTokenConfigContractId",
					TemplateId: "package:module:TokenConfig",
				}, {
					ContractId: "tokenConfigContractId",
					TemplateId: "package:module:TokenConfig",
				},
			}, resp.JSON200.DisclosedContracts)
			require.Equal(t, new(converters.RawInstanceAddressAsRawOrHashedAddress(defaultExecutor)), resp.JSON200.Executor)
			require.Equal(t, "feeTokenConfigContractId", resp.JSON200.FeeTokenConfigCid)
		})
		t.Run("Token Transfer - with executor + default CCVs", func(t *testing.T) {
			t.Parallel()
			customExecutorAddress := contracts.NewRawInstanceAddress("customExecutor", "owner")
			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: destSelector,
					Executor: struct {
						Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
						Type    oapiCommon.MessageExecutorType `json:"type"`
					}{
						Type:    oapiCommon.WithAddress,
						Address: new(converters.InstanceAddressAsRawOrHashedAddress(customExecutorAddress.InstanceAddress())),
					},
					FeeToken: oapiCommon.InstrumentId{
						Admin: "feeAdmin",
						Id:    "LINK",
					},
					Payload:  "0xdeadbeef",
					Receiver: "0x1234567890",
					TokenTransfer: &oapiCommon.TokenTransfer{
						Amount: "42",
						Token: oapiCommon.InstrumentId{
							Admin: oapiCommon.PartyId(tokenInstrumentId.Admin),
							Id:    string(tokenInstrumentId.Id),
						},
					},
				},
				SenderRequiredCCVs:    nil, // No SenderRequiredCCVs - API should return default CCVs
				TokenPoolRequiredCCVs: nil, // No TokenPoolRequiredCCVs - API should return default CCVs
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.RawOrHashedAddress{
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv1), // LaneMandatedCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv5), // LaneMandatedCCVs + DefaultCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv6), // DefaultCCVs
				// No TokenPoolRequiredCCVs
			}, resp.JSON200.Ccvs)
			require.Equal(t, new(converters.InstanceAddressAsRawOrHashedAddress(customExecutorAddress.InstanceAddress())), resp.JSON200.Executor)
		})
		t.Run("Token Transfer - no executor + custom + default CCVs", func(t *testing.T) {
			t.Parallel()
			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: oapiCommon.Message{
					DestinationChainSelector: destSelector,
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
						Amount: "42",
						Token: oapiCommon.InstrumentId{
							Admin: oapiCommon.PartyId(tokenInstrumentId.Admin),
							Id:    string(tokenInstrumentId.Id),
						},
					},
				},
				SenderRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()),
					converters.RawInstanceAddressAsRawOrHashedAddress(contracts.RawInstanceAddress("default-ccvs")), // sentinel that should lead to the DefaultCCVs to be included in the response
					converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),
				}),
				TokenPoolRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{
					converters.RawInstanceAddressAsRawOrHashedAddress(contracts.RawInstanceAddress("default-ccvs")), // sentinel that should lead to the DefaultCCVs to be included in the response
					converters.InstanceAddressAsRawOrHashedAddress(ccv4.InstanceAddress()),
				}),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.RawOrHashedAddress{
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv1),                // LaneMandatedCCVs
				converters.InstanceAddressAsRawOrHashedAddress(ccv2.InstanceAddress()), // SenderRequiredCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv5),                // DefaultCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv6),                // DefaultCCVs
				converters.RawInstanceAddressAsRawOrHashedAddress(ccv3),                // SenderRequiredCCVs
				converters.InstanceAddressAsRawOrHashedAddress(ccv4.InstanceAddress()), // TokenPoolRequiredCCVs
				// No TokenPoolRequiredCCVs
			}, resp.JSON200.Ccvs)
			require.Nil(t, resp.JSON200.Executor)
		})
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		validMessage := oapiCommon.Message{
			DestinationChainSelector: destSelector,
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
				Amount: "42",
				Token: oapiCommon.InstrumentId{
					Admin: oapiCommon.PartyId(tokenInstrumentId.Admin),
					Id:    string(tokenInstrumentId.Id),
				},
			},
		}
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPSendWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPSendWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid message", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			message := validMessage
			message.DestinationChainSelector = "0" // invalid selector
			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: message,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid message")
		})
		t.Run("Too many SenderRequiredCCVs", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message:            validMessage,
				SenderRequiredCCVs: new(slices.Repeat([]oapiCommon.RawOrHashedAddress{converters.InstanceAddressAsRawOrHashedAddress(contracts.HexToInstanceAddress("0x1"))}, MaxNumCCVs+1)),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "sender required CCVs exceeds maximum value")
		})
		t.Run("Too many TokenPoolRequiredCCVs", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message:               validMessage,
				TokenPoolRequiredCCVs: new(slices.Repeat([]oapiCommon.RawOrHashedAddress{converters.InstanceAddressAsRawOrHashedAddress(contracts.HexToInstanceAddress("0x1"))}, MaxNumCCVs+1)),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "token pool required CCVs exceeds maximum value")
		})
		t.Run("OnRamp not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("GlobalConfig not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenAdminRegistry not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("RMNRemote not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("FeeQuoter not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("GlobalConfig not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true) // No ActiveContract - parsing will fail
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("GlobalConfig - no DestChainConfig", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:       "globalconfig",
					CcipOwner:        "owner",
					ChainSelector:    types.NUMERIC(sourceSelector),
					DestChainConfigs: nil, // No DestChainConfigs
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "unsupported destination chain selector")
		})
		t.Run("GlobalConfig - destination chain disabled", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 false,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "destination chain is disabled")
		})
		t.Run("Invalid SenderRequiredCCV", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message:            validMessage,
				SenderRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{converters.RawInstanceAddressAsRawOrHashedAddress(contracts.RawInstanceAddress("invalid-address"))}),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid sender required CCV address")
		})
		t.Run("Invalid TokenPoolRequiredCCV", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message:               validMessage,
				TokenPoolRequiredCCVs: new([]oapiCommon.RawOrHashedAddress{converters.RawInstanceAddressAsRawOrHashedAddress(contracts.RawInstanceAddress("invalid-address"))}),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid token pool required CCV address")
		})
		t.Run("Empty resolvedCCVs", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
							// This should lead to no resolved CCVs as there are no default CCVs and no provided CCVs
							// This shouldn't happen, since either lane-mandated or default CCVs always need to be set, still breaking out early in the API
							DefaultCCVs:      nil, // No default CCVs
							LaneMandatedCCVs: nil, // No lane-mandated CCVs
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true) // No ActiveContract - parsing will fail
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message:               validMessage,
				TokenPoolRequiredCCVs: nil, // No TokenPoolRequiredCCVs
				SenderRequiredCCVs:    nil, //  No SenderRequiredCCVs
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "every message must be validated by at least one CCV")
		})
		t.Run("TokenAdminRegistry not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
							DefaultCCVs:               []chainlinkapi.RawInstanceAddress{contracts.NewRawInstanceAddress("ccv1", "owner").Binding()},
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true) // No ActiveContract - parsing will fail
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenConfig not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
							DefaultCCVs:               []chainlinkapi.RawInstanceAddress{contracts.NewRawInstanceAddress("ccv1", "owner").Binding()},
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
					InstanceId: "tokenadminregistry",
					CcipOwner:  "owner",
					EntryCount: 1,
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(feeTokenConfigAddress.InstanceAddress()).Return(tokenConfigContract("feeTokenConfigContractId", feeInstrumentId, nil), true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(nil, false)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no token config registered for token")
		})
		t.Run("TokenConfig not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
							DefaultCCVs:               []chainlinkapi.RawInstanceAddress{contracts.NewRawInstanceAddress("ccv1", "owner").Binding()},
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
					InstanceId: "tokenadminregistry",
					CcipOwner:  "owner",
					EntryCount: 1,
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(feeTokenConfigAddress.InstanceAddress()).Return(tokenConfigContract("feeTokenConfigContractId", feeInstrumentId, nil), true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(nil, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("No TokenPool set in TokenConfig", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OnRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "owner",
					ChainSelector: types.NUMERIC(sourceSelector),
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC(destSelector): {
							IsEnabled:                 true,
							MessageNetworkFeeUSDCents: "0",
							TokenNetworkFeeUSDCents:   "0",
							DefaultCCVs:               []chainlinkapi.RawInstanceAddress{contracts.NewRawInstanceAddress("ccv1", "owner").Binding()},
						},
					},
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
					InstanceId: "tokenadminregistry",
					CcipOwner:  "owner",
					EntryCount: 1,
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.FeeQuoter.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(feeTokenConfigAddress.InstanceAddress()).Return(tokenConfigContract("feeTokenConfigContractId", feeInstrumentId, nil), true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenConfig{
					InstanceId:         types.TEXT(hex.EncodeToString(encodedTokenInstrumentId.Bytes())),
					RegistryInstanceId: types.TEXT(tokenAdminRegistryAddress.InstanceID()),
					RegistryOwner:      types.PARTY(tokenAdminRegistryAddress.Owner()),
					Index:              0,
					InstrumentId:       tokenInstrumentId,
					TokenPool:          nil,
				}),
			}}, true)

			resp, err := client.PostCCIPSendWithResponse(t.Context(), oapiCCIP.CCIPSendRequest{
				Message: validMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no token pool registered for token")
		})
	})
}

func TestServer_PostCCIPExecute(t *testing.T) {
	t.Parallel()

	sourceSelector := uint64(123)
	destSelector := uint64(456)
	perPartyRouterFactoryAddress := contracts.NewRawInstanceAddress("perPartyRouterFactory", "owner")
	onRampAddress := contracts.NewRawInstanceAddress("onRamp", "owner")
	offRampAddress := contracts.NewRawInstanceAddress("offRamp", "owner")
	globalConfigAddress := contracts.NewRawInstanceAddress("globalConfig", "owner")
	tokenAdminRegistryAddress := contracts.NewRawInstanceAddress("tokenAdminRegistry", "owner")
	rmnRemoteAddress := contracts.NewRawInstanceAddress("rmnRemote", "owner")
	feeQuoterAddress := contracts.NewRawInstanceAddress("feeQuoter", "owner")
	tokenPoolAddress := contracts.NewRawInstanceAddress("tokenPool", "owner")
	cfg := config.CCIPAPIConfig{
		PerPartyRouterFactory: config.ContractIdentifier{InstanceAddress: perPartyRouterFactoryAddress.InstanceAddress()},
		OnRamp:                config.ContractIdentifier{InstanceAddress: onRampAddress.InstanceAddress()},
		OffRamp:               config.ContractIdentifier{InstanceAddress: offRampAddress.InstanceAddress()},
		GlobalConfig:          config.ContractIdentifier{InstanceAddress: globalConfigAddress.InstanceAddress()},
		TokenAdminRegistry:    config.ContractIdentifier{InstanceAddress: tokenAdminRegistryAddress.InstanceAddress()},
		RMNRemote:             config.ContractIdentifier{InstanceAddress: rmnRemoteAddress.InstanceAddress()},
		FeeQuoter:             config.ContractIdentifier{InstanceAddress: feeQuoterAddress.InstanceAddress()},
	}
	tokenInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: "tokenOwner",
		Id:    "LinkToken",
	}
	encodedTokenInstrumentId := contracts.EncodeInstrumentID(tokenInstrumentId)
	tokenConfigAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(hex.EncodeToString(encodedTokenInstrumentId.Bytes())), "owner")

	setup := func(t *testing.T) (*mocks.MockActiveContractStore, oapiCCIP.ClientWithResponsesInterface) {
		t.Helper()
		mockActiveContractStore := mocks.NewMockActiveContractStore(t)
		mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
		server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
		require.NoError(t, err)
		client := makeClient(t, server)

		return mockActiveContractStore, client
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockActiveContractStore, client := setup(t)

		mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "offRampContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "OffRamp",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(ccipruntime.OffRamp{
				InstanceId: "offramp",
				CcipOwner:  "owner",
			}),
		}}, true)
		mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "globalConfigContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "GlobalConfig",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
				InstanceId:         "globalconfig",
				CcipOwner:          "owner",
				ChainSelector:      types.NUMERIC(strconv.FormatUint(destSelector, 10)),
				DestChainConfigs:   nil,
				SourceChainConfigs: nil,
			}),
		}}, true)
		mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "tokenAdminRegistryContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "TokenAdminRegistry",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
				InstanceId: "tokenadminregistry",
				CcipOwner:  "owner",
				EntryCount: 1,
			}),
		}}, true)
		mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "rmnRemoteContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "RMNRemote",
			},
		}}, true)
		mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
			ContractId: "tokenConfigContractId",
			TemplateId: &apiv2.Identifier{
				PackageId:  "package",
				ModuleName: "module",
				EntityName: "TokenConfig",
			},
			CreateArguments: bindings.MarshalTemplateToRecord(core.TokenConfig{
				InstanceId:         types.TEXT(hex.EncodeToString(encodedTokenInstrumentId.Bytes())),
				RegistryInstanceId: types.TEXT(tokenAdminRegistryAddress.InstanceID()),
				RegistryOwner:      types.PARTY(tokenAdminRegistryAddress.Owner()),
				Index:              0,
				InstrumentId:       tokenInstrumentId,
				TokenPool: &core.PoolRegistration2{
					PoolOwner:      types.PARTY(tokenPoolAddress.Owner()),
					PoolInstanceId: types.TEXT(tokenPoolAddress.InstanceID()),
				},
			}),
		}}, true)

		t.Run("Message only", func(t *testing.T) {
			t.Parallel()
			message := protocol.Message{
				SourceChainSelector: protocol.ChainSelector(sourceSelector),
				DestChainSelector:   protocol.ChainSelector(destSelector),
			}
			encodedMessage, err := message.Encode()
			require.NoError(t, err)
			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: hex.EncodeToString(encodedMessage),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "offRampContractId",
					TemplateId: "package:module:OffRamp",
				}, {
					ContractId: "globalConfigContractId",
					TemplateId: "package:module:GlobalConfig",
				}, {
					ContractId: "tokenAdminRegistryContractId",
					TemplateId: "package:module:TokenAdminRegistry",
				}, {
					ContractId: "rmnRemoteContractId",
					TemplateId: "package:module:RMNRemote",
				},
			}, resp.JSON200.DisclosedContracts)
			require.NotNil(t, resp.JSON200.ContextData)
			require.Nil(t, resp.JSON200.TokenPool)
		})
		t.Run("Token Transfer", func(t *testing.T) {
			t.Parallel()
			message := protocol.Message{
				SourceChainSelector: protocol.ChainSelector(sourceSelector),
				DestChainSelector:   protocol.ChainSelector(destSelector),
				TokenTransfer: &protocol.TokenTransfer{
					Amount:                 big.NewInt(42),
					DestTokenAddress:       encodedTokenInstrumentId[:],
					DestTokenAddressLength: 32,
				},
			}
			encodedMessage, err := message.Encode()
			require.NoError(t, err)
			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: hex.EncodeToString(encodedMessage),
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusOK, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.ElementsMatch(t, []oapiCommon.DisclosedContract{
				{
					ContractId: "offRampContractId",
					TemplateId: "package:module:OffRamp",
				}, {
					ContractId: "globalConfigContractId",
					TemplateId: "package:module:GlobalConfig",
				}, {
					ContractId: "tokenAdminRegistryContractId",
					TemplateId: "package:module:TokenAdminRegistry",
				}, {
					ContractId: "rmnRemoteContractId",
					TemplateId: "package:module:RMNRemote",
				}, {
					ContractId: "tokenConfigContractId",
					TemplateId: "package:module:TokenConfig",
				},
			}, resp.JSON200.DisclosedContracts)
			require.NotNil(t, resp.JSON200.ContextData)
			require.Equal(t, new(tokenPoolAddress.String()), resp.JSON200.TokenPool)
		})
	})
	t.Run("Failure cases", func(t *testing.T) {
		t.Parallel()
		message := protocol.Message{
			SourceChainSelector: protocol.ChainSelector(sourceSelector),
			DestChainSelector:   protocol.ChainSelector(destSelector),
			TokenTransfer: &protocol.TokenTransfer{
				Amount:                 big.NewInt(42),
				DestTokenAddress:       encodedTokenInstrumentId[:],
				DestTokenAddressLength: 32,
			},
		}
		encodedMessage, err := message.Encode()
		validEncodedMessage := hex.EncodeToString(encodedMessage)
		require.NoError(t, err)
		t.Run("Invalid request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPExecuteWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit))
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("Oversized request", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPExecuteWithBodyWithResponse(t.Context(), "application/json", testhelpers.MakeOversizedRequest(RequestSizeLimit+1))
			require.NoError(t, err)
			require.Equalf(t, http.StatusRequestEntityTooLarge, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "request body too large")
		})
		t.Run("Invalid message hex", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: "invalidhex",
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid encoded message")
		})
		t.Run("Invalid message", func(t *testing.T) {
			t.Parallel()
			_, client := setup(t)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: "0x1234",
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "invalid encoded message")
		})
		t.Run("OffRamp not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("GlobalConfig not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenAdminRegistry not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("RMNRemote not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, false)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenAdminRegistry not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(nil, true) // No ActiveContract - parsing will fail
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("TokenConfig not found in store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
					InstanceId: "tokenadminregistry",
					CcipOwner:  "owner",
					EntryCount: 1,
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(nil, false)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no token config registered for token")
		})
		t.Run("TokenConfig not returned from store", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
					InstanceId: "tokenadminregistry",
					CcipOwner:  "owner",
					EntryCount: 1,
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(nil, true)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusInternalServerError, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
		})
		t.Run("No TokenPool set in TokenConfig", func(t *testing.T) {
			t.Parallel()
			mockActiveContractStore, client := setup(t)

			mockActiveContractStore.EXPECT().Get(cfg.OffRamp.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.GlobalConfig.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(cfg.TokenAdminRegistry.InstanceAddress).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenAdminRegistry{
					InstanceId: "tokenadminregistry",
					CcipOwner:  "owner",
					EntryCount: 1,
				}),
			}}, true)
			mockActiveContractStore.EXPECT().Get(cfg.RMNRemote.InstanceAddress).Return(nil, true)
			mockActiveContractStore.EXPECT().Get(tokenConfigAddress.InstanceAddress()).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.TokenConfig{
					InstanceId:         types.TEXT(hex.EncodeToString(encodedTokenInstrumentId.Bytes())),
					RegistryInstanceId: types.TEXT(tokenAdminRegistryAddress.InstanceID()),
					RegistryOwner:      types.PARTY(tokenAdminRegistryAddress.Owner()),
					Index:              0,
					InstrumentId:       tokenInstrumentId,
					TokenPool:          nil,
				}),
			}}, true)

			resp, err := client.PostCCIPExecuteWithResponse(t.Context(), oapiCCIP.CCIPExecuteRequest{
				EncodedMessage: validEncodedMessage,
			})
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, resp.StatusCode(), "unexpected response code, response: %s", string(resp.Body))
			require.Contains(t, string(resp.Body), "no token pool registered for token")
		})
	})
}

func TestServer_FilterContracts(t *testing.T) {
	t.Parallel()

	cfg := config.CCIPAPIConfig{
		PerPartyRouterFactory: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x1"),
		},
		OnRamp: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x2"),
		},
		OffRamp: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x3"),
		},
		GlobalConfig: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x4"),
		},
		TokenAdminRegistry: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x5"),
		},
		RMNRemote: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x6"),
		},
		FeeQuoter: config.ContractIdentifier{
			InstanceAddress: contracts.HexToInstanceAddress("0x7"),
		},
	}
	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
	require.NoError(t, err)

	t.Run("All filters", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.PerPartyRouterFactory.InstanceAddress,
				cfg.OnRamp.InstanceAddress,
				cfg.OffRamp.InstanceAddress,
				cfg.GlobalConfig.InstanceAddress,
				cfg.TokenAdminRegistry.InstanceAddress,
				cfg.RMNRemote.InstanceAddress,
				cfg.FeeQuoter.InstanceAddress,
			},
		)
		assert.Len(t, res, 7)
		assert.Contains(t, res, cfg.PerPartyRouterFactory.InstanceAddress)
		assert.Contains(t, res, cfg.OnRamp.InstanceAddress)
		assert.Contains(t, res, cfg.OffRamp.InstanceAddress)
		assert.Contains(t, res, cfg.GlobalConfig.InstanceAddress)
		assert.Contains(t, res, cfg.TokenAdminRegistry.InstanceAddress)
		assert.Contains(t, res, cfg.RMNRemote.InstanceAddress)
		assert.Contains(t, res, cfg.FeeQuoter.InstanceAddress)
	})
	t.Run("Filters unknown address", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts(
			[]contracts.InstanceAddress{
				cfg.PerPartyRouterFactory.InstanceAddress,
				contracts.HexToInstanceAddress("0x123456789"),
				cfg.FeeQuoter.InstanceAddress,
			},
		)
		assert.Len(t, res, 2)
		assert.Contains(t, res, cfg.PerPartyRouterFactory.InstanceAddress)
		assert.Contains(t, res, cfg.FeeQuoter.InstanceAddress)
		assert.NotContainsf(t, res, contracts.HexToInstanceAddress("0x123456789"), "unknown address should be filtered out")
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
