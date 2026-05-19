package token_standard

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

	cfg := config.TokenStandardAPIConfig{
		Registries: map[string]config.Registry{
			"Token1": {TokenId: "1", TokenType: config.TokenTypeLINK, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x1")}},
			"Token2": {TokenId: "2", TokenType: config.TokenTypeLINK, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x2")}},
			"Token3": {TokenId: "3", TokenType: config.TokenTypeLINK, ContractIdentifier: config.ContractIdentifier{InstanceAddress: contracts.HexToInstanceAddress("0x3")}},
		},
	}
	mockActiveContractStore := mocks.NewMockActiveContractStore(t)
	mockActiveContractStore.EXPECT().RegisterTemplates(mock.Anything).Times(len(cfg.Registries))
	server, err := NewServer(t.Context(), zerolog.Nop(), mockActiveContractStore, cfg)
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
