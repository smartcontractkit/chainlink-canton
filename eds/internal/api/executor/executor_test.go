package executor

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
