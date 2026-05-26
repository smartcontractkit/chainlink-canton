package sequences

import (
	"testing"

	"github.com/stretchr/testify/require"

	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func TestResolveTokenPoolParties_devenvFallback(t *testing.T) {
	t.Parallel()

	ccip, pool, err := resolveTokenPoolParties(tokenadapters.DeployTokenPoolInput{
		TokenRef: &datastore.AddressRef{Labels: datastore.NewLabelSet()},
	}, cldfcanton.Participant{PartyID: "solo-party"})
	require.NoError(t, err)
	require.Equal(t, "solo-party", ccip)
	require.Equal(t, "solo-party", pool)
}

func TestResolveTokenPoolParties_requiresLabelWithReadAs(t *testing.T) {
	t.Parallel()

	_, _, err := resolveTokenPoolParties(tokenadapters.DeployTokenPoolInput{
		TokenRef: &datastore.AddressRef{Labels: datastore.NewLabelSet()},
	}, cldfcanton.Participant{
		PartyID:        "operator",
		ReadAsPartyIDs: []string{"ccip-owner::abc"},
	})
	require.Error(t, err)
}

func TestResolveTokenPoolParties_fromCcipOwnerLabel(t *testing.T) {
	t.Parallel()

	ccip, pool, err := resolveTokenPoolParties(tokenadapters.DeployTokenPoolInput{
		TokenRef: &datastore.AddressRef{
			Labels: datastore.NewLabelSet("ccip-owner:ccip-party"),
		},
	}, cldfcanton.Participant{
		PartyID:        "operator",
		ReadAsPartyIDs: []string{"ccip-owner::abc"},
	})
	require.NoError(t, err)
	require.Equal(t, "ccip-party", ccip)
	require.Equal(t, "ccip-party", pool)
}

func TestResolveTokenPoolParties_ignoresPoolOwnerLabel(t *testing.T) {
	t.Parallel()

	ccip, pool, err := resolveTokenPoolParties(tokenadapters.DeployTokenPoolInput{
		TokenRef: &datastore.AddressRef{
			Labels: datastore.NewLabelSet("ccip-owner:ccip-party", "pool-owner:other-party"),
		},
	}, cldfcanton.Participant{
		PartyID:        "operator",
		ReadAsPartyIDs: []string{"ccip-owner::abc"},
	})
	require.NoError(t, err)
	require.Equal(t, "ccip-party", ccip)
	require.Equal(t, "ccip-party", pool)
}
