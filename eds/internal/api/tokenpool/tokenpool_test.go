package tokenpool

import (
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
)

func TestServer_FilterContracts(t *testing.T) {
	t.Parallel()

	// Need to mock the two possible token pool types
	// In order to reconstruct all possible contracts, the Server will query each configured token pool contract,
	// parse it and inspect its rate limiters.
	rl1 := contracts.NewRawInstanceAddress("rateLimiter1", "owner")
	rl2 := contracts.NewRawInstanceAddress("rateLimiter2", "owner")
	rl3 := contracts.NewRawInstanceAddress("rateLimiter3", "owner")
	tokenPool1 := burnminttokenpool.BurnMintTokenPool{
		RemoteChainConfigs: map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
			types.NUMERIC("123"): {
				FinalityConfig:     core.FinalityConfig{WaitForFinality: new(types.UNIT)}, // Unit types must be set in order for the unmarshal to work
				InboundRateLimiter: rl1.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: rl2.Binding(),
				OutboundRateLimiter:                        rl3.Binding(),
			},
		},
		TransferTimeout: burnminttokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
	}
	rl4 := contracts.NewRawInstanceAddress("rateLimiter4", "owner")
	rl5 := contracts.NewRawInstanceAddress("rateLimiter5", "owner")
	rl6 := contracts.NewRawInstanceAddress("rateLimiter6", "owner")
	tokenPool2 := lockreleasetokenpool.LockReleaseTokenPool{
		RemoteChainConfigs: map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig{
			"456": {
				FinalityConfig:     core.FinalityConfig{WaitForFinality: new(types.UNIT)},
				InboundRateLimiter: rl4.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: rl5.Binding(),
				OutboundRateLimiter:                        rl6.Binding(),
			},
		},
		TokenTransferFeeConfigs: nil,
		PoolReceiveContext:      splice_api_token_metadata_v1.ChoiceContext{},
		TransferTimeout:         lockreleasetokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
		Deps:                    lockreleasetokenpool.LockReleaseTokenPoolDeps{},
	}
	_ = tokenPool1
	_ = tokenPool2

	cfg := config.TokenPoolAPIConfig{
		TokenPools: map[string]config.TokenPool{
			"TokenPool1": {Type: config.TokenPoolTypeBurnMint, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x1")}},
			"TokenPool2": {Type: config.TokenPoolTypeLockRelease, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x2")}},
		},
	}
	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything)
	mockInstrumentHoldingStore := mocks.NewMockInstrumentHoldingStore(t)
	mockInstrumentHoldingStore.EXPECT().RegisterParty(mock.Anything)
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, mockInstrumentHoldingStore, cfg)
	require.NoError(t, err)

	mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x1")).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{CreateArguments: bindings.MarshalTemplateToRecord(tokenPool1)}}, true)
	mockActiveContractStore.EXPECT().Get(contracts.HexToInstanceAddress("0x2")).Return(&apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{CreateArguments: bindings.MarshalTemplateToRecord(tokenPool2)}}, true)

	t.Run("All filters", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts([]contracts.InstanceAddress{
			contracts.HexToInstanceAddress("0x1"),
			contracts.HexToInstanceAddress("0x2"),
			rl1.InstanceAddress(),
			rl2.InstanceAddress(),
			rl3.InstanceAddress(),
			rl4.InstanceAddress(),
			rl5.InstanceAddress(),
			rl6.InstanceAddress(),
		})
		assert.Len(t, res, 8)
		assert.Contains(t, res, contracts.HexToInstanceAddress("0x1"))
		assert.Contains(t, res, contracts.HexToInstanceAddress("0x2"))
		assert.Contains(t, res, rl1.InstanceAddress())
		assert.Contains(t, res, rl2.InstanceAddress())
		assert.Contains(t, res, rl3.InstanceAddress())
		assert.Contains(t, res, rl4.InstanceAddress())
		assert.Contains(t, res, rl5.InstanceAddress())
		assert.Contains(t, res, rl6.InstanceAddress())
	})
	t.Run("Filters unknown adddresses", func(t *testing.T) {
		t.Parallel()

		res := server.FilterContracts([]contracts.InstanceAddress{
			contracts.HexToInstanceAddress("0x1"),
			rl5.InstanceAddress(),
		})
		assert.Len(t, res, 2)
		assert.Contains(t, res, contracts.HexToInstanceAddress("0x1"))
		assert.Contains(t, res, rl5.InstanceAddress())
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
