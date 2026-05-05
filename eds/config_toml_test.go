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

func TestExampleConfig(t *testing.T) {
	t.Parallel()

	want, err := config.Read(strings.NewReader(exampleMergedGoldenToml))
	require.NoError(t, err)

	got, err := config.ReadAndMerge(strings.NewReader(exampleConfigToml), strings.NewReader(exampleSecretsConfigToml))
	require.NoError(t, err)
	require.NoError(t, got.Validate())

	assert.Equal(t, want, got)
}
