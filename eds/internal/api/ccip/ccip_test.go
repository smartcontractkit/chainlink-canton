package ccip

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/mocks"
)

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
