package sequences

import (
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestRemoteChainConfigured(t *testing.T) {
	t.Parallel()

	require.True(t, remoteChainConfigured(map[types.NUMERIC]any{"123.": nil}, "123"))
	require.True(t, remoteChainConfigured(map[types.NUMERIC]any{"123.0000000000": nil}, "123"))
	require.True(t, remoteChainConfigured(map[types.NUMERIC]any{"123": nil}, "123"))
	require.False(t, remoteChainConfigured(map[types.NUMERIC]any{"1234.": nil}, "123"))
	require.False(t, remoteChainConfigured(map[types.NUMERIC]any{}, "123"))
}
