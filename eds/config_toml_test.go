package eds

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
)

//go:embed config.toml
var exampleConfigToml string

//go:embed config_secrets.toml
var exampleSecretsConfigToml string

func TestExampleConfig(t *testing.T) {
	t.Parallel()

	// Read and parse configs separately
	exampleConfig, err := config.Read(strings.NewReader(exampleConfigToml))
	require.NoError(t, err)

	secretsConfig, err := config.Read(strings.NewReader(exampleSecretsConfigToml))
	require.NoError(t, err)

	// Read and merge multiple config files
	cfg, err := config.ReadAndMerge(strings.NewReader(exampleConfigToml), strings.NewReader(exampleSecretsConfigToml))
	require.NoError(t, err)
	err = cfg.Validate()
	require.NoError(t, err)

	// Check that the merge was successful
	exampleConfig.Node.AuthConfig.JWT = secretsConfig.Node.AuthConfig.JWT
	assert.Equal(t, exampleConfig, cfg)
}
