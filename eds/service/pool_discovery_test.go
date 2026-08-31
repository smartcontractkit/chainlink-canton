package service

import (
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
)

func TestPoolDiscoveryService_CheckForNewPools(t *testing.T) {
	t.Parallel()

	const observerParty = "chainlink-registry-observer::fingerprint"
	const poolOwner = "thirdPartyIssuer::fingerprint"
	const instanceID = "pool1"

	pool := burnminttokenpool.BurnMintTokenPool{
		InstanceId:      types.TEXT(instanceID),
		PoolOwner:       types.PARTY(poolOwner),
		CcipOwner:       types.PARTY(poolOwner),
		TransferTimeout: burnminttokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
	}
	// The address a discovered pool is registered under must match what any other caller
	// (e.g. FilterContracts, PostTokenPoolSend) resolves from instanceId+poolOwner - not from
	// the ledger's raw ContractId.
	wantAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(instanceID), types.PARTY(poolOwner)).InstanceAddress()

	activeContract := &apiv2.ActiveContract{
		CreatedEvent: &apiv2.CreatedEvent{
			ContractId:      "someContractId",
			CreateArguments: bindings.MarshalTemplateToRecord(pool),
		},
	}

	mockStore := mocks.NewMockActiveContractStore(t)
	mockStore.EXPECT().RegisterTemplates(mock.Anything)
	mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)

	server, err := tokenpool.NewServer(t.Context(), zerolog.Nop(), mockStore, mockInstrumentHoldingStore, config.TokenPoolAPIConfig{})
	require.NoError(t, err)

	svc := NewPoolDiscoveryService(zerolog.Nop(), mockStore, server, observerParty)

	mockStore.EXPECT().
		GetByTemplateId(types.PARTY(observerParty), contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{})).
		Return([]*apiv2.ActiveContract{activeContract}, true)
	mockStore.EXPECT().
		GetByTemplateId(types.PARTY(observerParty), contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{})).
		Return(nil, false)

	svc.checkForNewPools(t.Context())

	_, discovered := svc.discoveredPools[wantAddress]
	assert.True(t, discovered, "pool should be discovered under the instanceId+poolOwner derived address")

	// Registering must have used the pool's own owner, not the observer party.
	mockStore.EXPECT().Get(wantAddress).Return(activeContract, true)
	res := server.FilterContracts([]contracts.InstanceAddress{wantAddress})
	assert.Contains(t, res, wantAddress)

	// A second scan of the same active contract must not re-register it (the earlier
	// GetByTemplateId expectation has no call-count limit, so it still matches here).
	svc.checkForNewPools(t.Context())
}

func TestPoolDiscoveryService_CheckForNewPools_LockRelease(t *testing.T) {
	t.Parallel()

	const observerParty = "chainlink-registry-observer::fingerprint"
	const poolOwner = "thirdPartyIssuer::fingerprint"
	const instanceID = "pool2"

	pool := lockreleasetokenpool.LockReleaseTokenPool{
		InstanceId:      types.TEXT(instanceID),
		PoolOwner:       types.PARTY(poolOwner),
		CcipOwner:       types.PARTY(poolOwner),
		TransferTimeout: lockreleasetokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
	}
	wantAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(instanceID), types.PARTY(poolOwner)).InstanceAddress()

	activeContract := &apiv2.ActiveContract{
		CreatedEvent: &apiv2.CreatedEvent{
			ContractId:      "someOtherContractId",
			CreateArguments: bindings.MarshalTemplateToRecord(pool),
		},
	}

	mockStore := mocks.NewMockActiveContractStore(t)
	mockStore.EXPECT().RegisterTemplates(mock.Anything)
	mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)

	server, err := tokenpool.NewServer(t.Context(), zerolog.Nop(), mockStore, mockInstrumentHoldingStore, config.TokenPoolAPIConfig{})
	require.NoError(t, err)

	svc := NewPoolDiscoveryService(zerolog.Nop(), mockStore, server, observerParty)

	mockStore.EXPECT().
		GetByTemplateId(types.PARTY(observerParty), contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{})).
		Return(nil, false)
	mockStore.EXPECT().
		GetByTemplateId(types.PARTY(observerParty), contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{})).
		Return([]*apiv2.ActiveContract{activeContract}, true)

	svc.checkForNewPools(t.Context())

	_, discovered := svc.discoveredPools[wantAddress]
	assert.True(t, discovered, "lock release pool should be discovered under the instanceId+poolOwner derived address")

	mockStore.EXPECT().Get(wantAddress).Return(activeContract, true)
	res := server.FilterContracts([]contracts.InstanceAddress{wantAddress})
	assert.Contains(t, res, wantAddress)
}
