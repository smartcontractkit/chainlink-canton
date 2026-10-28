package adapters

import (
	"context"
	"testing"

	ccvadapters "github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

func TestCantonCommitteeVerifierOnchainUnsupportedWrites(t *testing.T) {
	t.Parallel()

	adapter := &CantonCommitteeVerifierOnchain{}
	env := deployment.Environment{}

	t.Run("ApplySignatureConfigs", func(t *testing.T) {
		ops, err := adapter.ApplySignatureConfigs(context.Background(), env, 0, "", ccvadapters.SignatureConfigChange{})
		require.Nil(t, ops)
		require.EqualError(t, err, "apply signature configs is not supported on Canton")
	})

	t.Run("SetAllowedFinalityConfig", func(t *testing.T) {
		ops, err := adapter.SetAllowedFinalityConfig(context.Background(), env, 0, "", false, false, 0)
		require.Nil(t, ops)
		require.EqualError(t, err, "set allowed finality config is not supported on Canton")
	})

	t.Run("ApplyAllowlistUpdates", func(t *testing.T) {
		ops, err := adapter.ApplyAllowlistUpdates(context.Background(), env, 0, "", 0, false, nil, nil)
		require.Nil(t, ops)
		require.EqualError(t, err, "apply allowlist updates is not supported on Canton")
	})
}
