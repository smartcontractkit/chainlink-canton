package sequences

import (
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/core"
)

func TestRmnRemoteDeployTemplate_passesCustomObservers(t *testing.T) {
	t.Parallel()

	observers := []types.PARTY{"ccvOwner::abc", "other::def"}
	cursed := []types.TEXT{"deadbeef"}

	tmpl := rmnRemoteDeployTemplate(
		types.PARTY("rmn-owner::xyz"),
		types.PARTY("ccip-owner::xyz"),
		RMNRemoteParams{Template: core.RMNRemote{
			CustomObservers: observers,
			CursedSubjects:  cursed,
		}},
	)

	require.Equal(t, types.PARTY("rmn-owner::xyz"), tmpl.RmnOwner)
	require.Equal(t, types.PARTY("ccip-owner::xyz"), tmpl.CcipOwner)
	require.Equal(t, observers, tmpl.CustomObservers)
	require.Equal(t, cursed, tmpl.CursedSubjects)
}

func TestRequireRMNOwnerParty(t *testing.T) {
	t.Parallel()

	party, err := requireRMNOwnerParty(DeployChainContractsParams{RMNOwnerParty: "rmn-owner::abc"})
	require.NoError(t, err)
	require.Equal(t, "rmn-owner::abc", string(party))

	_, err = requireRMNOwnerParty(DeployChainContractsParams{})
	require.ErrorContains(t, err, "RMNOwnerParty is required")
}
