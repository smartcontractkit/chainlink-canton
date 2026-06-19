package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveConfigPath(t *testing.T) {
	t.Setenv(envConfigFile, "")

	require.Equal(t, "env-canton-evm-out.toml", ResolveConfigPath(EnvDevenv))
	require.Equal(t, "env-prod-testnet-out.toml", ResolveConfigPath(EnvProdTestnet))

	t.Setenv(envConfigFile, "env-prod-testnet.ci.toml")
	require.Equal(t, "env-prod-testnet.ci.toml", ResolveConfigPath(EnvProdTestnet))

	t.Setenv(envConfigFile, "ccip/devenv/custom.toml")
	require.Equal(t, "custom.toml", ResolveConfigPath(EnvDevenv))
}

func TestParseCCIPEnv(t *testing.T) {
	env, err := ParseCCIPEnv("devenv")
	require.NoError(t, err)
	require.Equal(t, EnvDevenv, env)

	env, err = ParseCCIPEnv("prod-testnet")
	require.NoError(t, err)
	require.Equal(t, EnvProdTestnet, env)

	_, err = ParseCCIPEnv("unknown")
	require.Error(t, err)
}
