package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnvironments(t *testing.T) {
	t.Run("loads valid environments config", func(t *testing.T) {
		// Create temp file
		content := `
environments:
  mainnet-v1:
    party: "ccip-owner::12345"
    description: "Production mainnet"
    contracts:
      router: "mainnet-v1"
      onRamp: "mainnet-v1"
      feeQuoter: "mainnet-v1"
      offRamp: "mainnet-v1"
      ccv: "mainnet-v1"
      tokenAdminRegistry: "mainnet-v1"
      tokenPool: "mainnet-v1"
  testnet:
    party: "ccip-test::67890"
    description: "Test network"
    contracts:
      router: "testnet"
      onRamp: "testnet"
      feeQuoter: "testnet"
      offRamp: "testnet"
      ccv: "testnet"
      tokenAdminRegistry: "testnet"
      tokenPool: "testnet"
`
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "environments.yaml")
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)

		config, err := LoadEnvironments(tmpFile)
		require.NoError(t, err)

		assert.Len(t, config.Environments, 2)

		mainnet, ok := config.Environments["mainnet-v1"]
		assert.True(t, ok)
		assert.Equal(t, "ccip-owner::12345", mainnet.Party)
		assert.Equal(t, "Production mainnet", mainnet.Description)
		assert.Equal(t, "mainnet-v1", mainnet.Contracts.Router)
		assert.Equal(t, "mainnet-v1", mainnet.Contracts.OnRamp)

		testnet, ok := config.Environments["testnet"]
		assert.True(t, ok)
		assert.Equal(t, "ccip-test::67890", testnet.Party)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := LoadEnvironments("/non/existent/path.yaml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read")
	})

	t.Run("returns error for invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "invalid.yaml")
		err := os.WriteFile(tmpFile, []byte("not: valid: yaml: content:"), 0644)
		require.NoError(t, err)

		_, err = LoadEnvironments(tmpFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse")
	})

	t.Run("returns error for empty environments", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "empty.yaml")
		err := os.WriteFile(tmpFile, []byte("environments: {}"), 0644)
		require.NoError(t, err)

		_, err = LoadEnvironments(tmpFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no environments defined")
	})
}

func TestEnvironmentsConfig_EnvironmentNames(t *testing.T) {
	config := &EnvironmentsConfig{
		Environments: map[string]EnvironmentConfig{
			"mainnet": {},
			"testnet": {},
			"staging": {},
		},
	}

	names := config.EnvironmentNames()
	assert.Len(t, names, 3)
	assert.Contains(t, names, "mainnet")
	assert.Contains(t, names, "testnet")
	assert.Contains(t, names, "staging")
}

func TestEnvironmentsConfig_GetEnvironment(t *testing.T) {
	config := &EnvironmentsConfig{
		Environments: map[string]EnvironmentConfig{
			"mainnet": {
				Party:       "party-123",
				Description: "Main network",
			},
		},
	}

	t.Run("returns environment when exists", func(t *testing.T) {
		env, ok := config.GetEnvironment("mainnet")
		assert.True(t, ok)
		assert.Equal(t, "party-123", env.Party)
		assert.Equal(t, "Main network", env.Description)
	})

	t.Run("returns false when environment does not exist", func(t *testing.T) {
		_, ok := config.GetEnvironment("nonexistent")
		assert.False(t, ok)
	})
}
