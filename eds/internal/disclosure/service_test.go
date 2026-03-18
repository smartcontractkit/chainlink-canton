package disclosure

import (
	"context"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

// testInstanceAddress returns a deterministic InstanceAddress for testing.
func testInstanceAddress(instanceID string, party string) contracts.InstanceAddress {
	return contracts.InstanceID(instanceID).RawInstanceAddress(types.PARTY(party)).InstanceAddress()
}

func defaultTestConfig(contractStore *mocks.MockContractStore, instrumentStore *mocks.MockInstrumentHoldingStore) DisclosureServiceConfig {
	return DisclosureServiceConfig{
		ContractStore:          contractStore,
		InstrumentHoldingStore: instrumentStore,
		PerPartyRouterFactory:  testInstanceAddress("router-factory", "party"),
		OnRamp:                 testInstanceAddress("onramp", "party"),
		OffRamp:                testInstanceAddress("offramp", "party"),
		GlobalConfig:           testInstanceAddress("global-config", "party"),
		TokenAdminRegistry:     testInstanceAddress("token-admin-registry", "party"),
		RMNRemote:              testInstanceAddress("rmn-remote", "party"),
		FeeQuoter:              testInstanceAddress("fee-quoter", "party"),
		CCVs:                   []contracts.InstanceAddress{testInstanceAddress("ccv1", "party")},
		TokenPools:             []contracts.InstanceAddress{testInstanceAddress("token-pool", "party")},
	}
}

// makeActiveContract returns a minimal ActiveContract for disclosure tests.
func makeActiveContract(contractID string, templateID *apiv2.Identifier) *apiv2.ActiveContract {
	return &apiv2.ActiveContract{
		CreatedEvent: &apiv2.CreatedEvent{
			ContractId: contractID,
			TemplateId: templateID,
		},
		SynchronizerId: "sync-1",
	}
}

func TestNewDisclosureService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	contractStore := mocks.NewMockContractStore(t)
	instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
	config := defaultTestConfig(contractStore, instrumentStore)

	svc := NewDisclosureService(ctx, config)
	require.NotNil(t, svc)
}

func TestDisclosureService_GetDisclosure(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns error when address is not configured", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		unknownAddr := testInstanceAddress("unknown-contract", "party")
		disclosure, err := svc.GetDisclosure(ctx, unknownAddr)
		require.Error(t, err)
		require.Nil(t, disclosure)
		require.Contains(t, err.Error(), "not found")
		require.Contains(t, err.Error(), unknownAddr.Hex())
	})

	t.Run("returns error when contract store returns nil", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.OnRamp).Return(nil)

		disclosure, err := svc.GetDisclosure(ctx, config.OnRamp)
		require.Error(t, err)
		require.Nil(t, disclosure)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("returns disclosure when contract is found", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		templateID := &apiv2.Identifier{
			PackageId:  "#test-pkg",
			ModuleName: "Test",
			EntityName: "OnRamp",
		}
		activeContract := makeActiveContract("contract-123", templateID)
		contractStore.EXPECT().GetContract(ctx, config.OnRamp).Return(activeContract)

		disclosure, err := svc.GetDisclosure(ctx, config.OnRamp)
		require.NoError(t, err)
		require.NotNil(t, disclosure)
		require.Equal(t, "contract-123", disclosure.ContractId)
		require.Equal(t, templateID, disclosure.TemplateId)
		require.Equal(t, "sync-1", disclosure.SynchronizerId)
	})
}

func TestDisclosureService_GetCCIPSendDisclosures(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns all disclosures when all lookups succeed", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.OnRamp).Return(makeActiveContract("onramp-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.GlobalConfig).Return(makeActiveContract("global-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.TokenAdminRegistry).Return(makeActiveContract("tar-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.RMNRemote).Return(makeActiveContract("rmn-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.FeeQuoter).Return(makeActiveContract("fee-1", nil))
		for _, ccvAddr := range config.CCVs {
			contractStore.EXPECT().GetContract(ctx, ccvAddr).Return(makeActiveContract("ccv-1", nil))
		}

		result, err := svc.GetCCIPSendDisclosures(ctx, CCIPSendRequest{CCVs: config.CCVs})
		require.NoError(t, err)
		require.NotNil(t, result.OnRamp)
		require.Equal(t, "onramp-1", result.OnRamp.ContractId)
		require.NotNil(t, result.GlobalConfig)
		require.NotNil(t, result.TokenAdminRegistry)
		require.NotNil(t, result.RMNRemote)
		require.NotNil(t, result.FeeQuoter)
		require.Len(t, result.CCVs, len(config.CCVs))
	})

	t.Run("returns error when onRamp disclosure fails", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.OnRamp).Return(nil)

		result, err := svc.GetCCIPSendDisclosures(ctx, CCIPSendRequest{CCVs: config.CCVs})
		require.Error(t, err)
		require.Contains(t, err.Error(), "onRamp")
		require.Contains(t, err.Error(), "not found")
		require.Nil(t, result.OnRamp)
	})
}

func TestDisclosureService_GetCCIPExecuteDisclosures(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns disclosures with nil token pool and holding when message has no token transfer", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.OffRamp).Return(makeActiveContract("offramp-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.GlobalConfig).Return(makeActiveContract("global-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.TokenAdminRegistry).Return(makeActiveContract("tar-1", nil))
		contractStore.EXPECT().GetContract(ctx, config.RMNRemote).Return(makeActiveContract("rmn-1", nil))
		for _, ccvAddr := range config.CCVs {
			contractStore.EXPECT().GetContract(ctx, ccvAddr).Return(makeActiveContract("ccv-1", nil))
		}

		result, err := svc.GetCCIPExecuteDisclosures(ctx, CCIPExecuteRequest{
			Message: protocol.Message{}, // no token transfer
			CCVs:    config.CCVs,
		})
		require.NoError(t, err)
		require.NotNil(t, result.OffRamp)
		require.Nil(t, result.TokenPool)
		require.Nil(t, result.TokenPoolHolding)
	})

	t.Run("returns error when offRamp disclosure fails", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.OffRamp).Return(nil)

		result, err := svc.GetCCIPExecuteDisclosures(ctx, CCIPExecuteRequest{
			Message: protocol.Message{},
			CCVs:    config.CCVs,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "offRamp")
		require.Nil(t, result.OffRamp)
	})
}

func TestDisclosureService_GetPerPartyRouterFactory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("returns disclosure when contract is found", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.PerPartyRouterFactory).Return(makeActiveContract("router-factory-1", nil))

		result, err := svc.GetPerPartyRouterFactory(ctx, PerPartyRouterFactoryRequest{})
		require.NoError(t, err)
		require.NotNil(t, result.PerPartyRouterFactory)
		require.Equal(t, "router-factory-1", result.PerPartyRouterFactory.ContractId)
	})

	t.Run("returns error when contract store returns nil", func(t *testing.T) {
		t.Parallel()
		contractStore := mocks.NewMockContractStore(t)
		instrumentStore := mocks.NewMockInstrumentHoldingStore(t)
		config := defaultTestConfig(contractStore, instrumentStore)
		svc := NewDisclosureService(ctx, config)

		contractStore.EXPECT().GetContract(ctx, config.PerPartyRouterFactory).Return(nil)

		result, err := svc.GetPerPartyRouterFactory(ctx, PerPartyRouterFactoryRequest{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "perPartyRouterFactory")
		require.Contains(t, err.Error(), "not found")
		require.Nil(t, result.PerPartyRouterFactory)
	})
}
