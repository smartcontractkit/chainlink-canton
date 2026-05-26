package adapters

import (
	"testing"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/require"
)

func TestDataStoreFromConfigureChainForLanesInput_FamilyExtras(t *testing.T) {
	t.Parallel()

	ds := datastore.NewMemoryDataStore().Seal()
	input := ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: 42,
		FamilyExtras: map[string]any{
			ConfigureLanesDataStoreFamilyExtraKey: ds,
		},
	}

	got, err := dataStoreFromConfigureChainForLanesInput(input)
	require.NoError(t, err)
	require.Equal(t, ds, got)
}

func TestDataStoreFromConfigureChainForLanesInput_Cache(t *testing.T) {
	t.Parallel()

	ds := datastore.NewMemoryDataStore().Seal()
	cacheConfigureLanesDataStore(8706591216959472610, ds)

	got, err := dataStoreFromConfigureChainForLanesInput(ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: 8706591216959472610,
	})
	require.NoError(t, err)
	require.Equal(t, ds, got)
}

func TestDataStoreFromConfigureChainForLanesInput_Missing(t *testing.T) {
	t.Parallel()

	_, err := dataStoreFromConfigureChainForLanesInput(ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: 1,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), ConfigureLanesDataStoreFamilyExtraKey)
}
