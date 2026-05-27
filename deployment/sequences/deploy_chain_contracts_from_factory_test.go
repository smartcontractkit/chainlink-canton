package sequences

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequireRMNOwnerParty(t *testing.T) {
	t.Parallel()

	party, err := requireRMNOwnerParty(DeployChainContractsParams{RMNOwnerParty: "rmn-owner::abc"})
	require.NoError(t, err)
	require.Equal(t, "rmn-owner::abc", string(party))

	_, err = requireRMNOwnerParty(DeployChainContractsParams{})
	require.ErrorContains(t, err, "RMNOwnerParty is required")
}
