package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
)

func TestTokenDirectionsForEnv(t *testing.T) {
	t.Parallel()

	cfg, err := decodeTokenConfigTOML(defaultTokenConfigPath())
	require.NoError(t, err)

	devenvDirs, err := tokenDirectionsForEnv(cfg, EnvDevenv)
	require.NoError(t, err)
	require.Equal(t, "BurnMintTokenPool", devenvDirs.EVMToCanton.PoolType)
	require.Contains(t, devenvDirs.EVMToCanton.PoolQualifier, "BurnMintTokenPool 2.0.0 [default]")
	require.Equal(t, "LockReleaseTokenPool", devenvDirs.CantonToEVM.PoolType)

	prodDirs, err := tokenDirectionsForEnv(cfg, EnvProdTestnet)
	require.NoError(t, err)
	require.Equal(t, "TEST", prodDirs.EVMToCanton.PoolQualifier)
	require.Equal(t, "LINK", prodDirs.EVMToCanton.RemotePoolQualifier)
	require.Equal(t, "LINK", prodDirs.CantonToEVM.PoolQualifier)
}

func TestLoadTokenDirection_envSelection(t *testing.T) {
	t.Setenv(envTokenTestConfig, defaultTokenConfigPath())

	devenvDir := loadTokenDirection(t, EnvDevenv, directionEVMToCanton)
	require.Equal(t, "BurnMintTokenPool", string(devenvDir.PoolRef.Type))
	require.Contains(t, devenvDir.PoolRef.Qualifier, "BurnMintTokenPool 2.0.0 [default]")
	require.Nil(t, devenvDir.RemotePoolRef)

	prodDir := loadTokenDirection(t, EnvProdTestnet, directionEVMToCanton)
	require.Equal(t, "TEST", prodDir.PoolRef.Qualifier)
	require.NotNil(t, prodDir.RemotePoolRef)
	require.Equal(t, "LINK", prodDir.RemotePoolRef.Qualifier)
}

func TestLoadTokenDirection_missingEnvSection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token_transfer_config.toml")
	require.NoError(t, os.WriteFile(path, []byte(`
[devenv.evm_to_canton]
pool_type = "BurnMintTokenPool"
pool_version = "2.0.0"
pool_qualifier = "TEST"
transfer_amount = "1"
execution_gas_limit = 1
finality_config = 0
`), 0o600))

	t.Setenv(envTokenTestConfig, path)

	_, err := tokenDirectionsForEnv(mustDecodeTokenConfig(t, path), EnvProdTestnet)
	require.Error(t, err)
	require.Contains(t, err.Error(), "prod-testnet")
}

func mustDecodeTokenConfig(t *testing.T, path string) tokenConfigTOML {
	t.Helper()

	cfg, err := decodeTokenConfigTOML(path)
	require.NoError(t, err)

	return cfg
}

func decodeTokenConfigTOML(path string) (tokenConfigTOML, error) {
	var cfg tokenConfigTOML
	_, err := toml.DecodeFile(path, &cfg)

	return cfg, err
}
