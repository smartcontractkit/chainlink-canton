package eds

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
)

//go:embed testdata/config.toml
var exampleConfigToml string

//go:embed testdata/config_secrets.toml
var exampleSecretsConfigToml string

//go:embed testdata/config_merged_golden.toml
var exampleMergedGoldenToml string

//go:embed testdata/eds-token-pools.toml
var poolDiscoveryConfigToml string

//go:embed testdata/eds-token-pools_secrets.toml
var poolDiscoverySecretsConfigToml string

func TestExampleConfig(t *testing.T) {
	t.Parallel()

	want, err := config.Read(strings.NewReader(exampleMergedGoldenToml))
	require.NoError(t, err)

	got, err := config.ReadAndMerge(strings.NewReader(exampleConfigToml), strings.NewReader(exampleSecretsConfigToml))
	require.NoError(t, err)
	require.NoError(t, got.Validate())

	assert.Equal(t, want, got)
}

// TestPoolDiscoveryConfig proves the standalone registry-discovery instance config (no
// hand-configured pools, token_pool_api populated entirely by discovery) parses and validates.
func TestPoolDiscoveryConfig(t *testing.T) {
	t.Parallel()

	got, err := config.ReadAndMerge(strings.NewReader(poolDiscoveryConfigToml), strings.NewReader(poolDiscoverySecretsConfigToml))
	require.NoError(t, err)
	require.NoError(t, got.Validate())

	assert.True(t, got.TokenPoolAPIConfig.Enabled)
	assert.Empty(t, got.TokenPoolAPIConfig.TokenPools)
	assert.True(t, got.RegistryAPIConfig.Enabled)
	assert.NotEmpty(t, got.RegistryAPIConfig.PartyID)
}
