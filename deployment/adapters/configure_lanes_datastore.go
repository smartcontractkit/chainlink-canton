package adapters

import (
	"sync"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

// configureLanesDataStoreCache holds e.DataStore keyed by chain selector while the generic
// ConfigureChainsForLanesFromTopology changeset resolves contracts. ccv/devenv and other callers
// that use the bare v2 changeset do not inject FamilyExtras["canton.dataStore"]; caching the
// datastore from adapter lookup methods lets configureChainForLanes recover it in the same pass.
var configureLanesDataStoreCache sync.Map

func cacheConfigureLanesDataStore(chainSelector uint64, ds datastore.DataStore) {
	if ds != nil {
		configureLanesDataStoreCache.Store(chainSelector, ds)
	}
}

func cachedConfigureLanesDataStore(chainSelector uint64) (datastore.DataStore, bool) {
	value, ok := configureLanesDataStoreCache.Load(chainSelector)
	if !ok {
		return nil, false
	}
	ds, ok := value.(datastore.DataStore)

	return ds, ok && ds != nil
}

// cantonChainFamilyWithDataStoreCache wraps CantonChainFamilyAdapter and records the datastore
// passed to contract-resolution methods before lane configure dispatch runs.
type cantonChainFamilyWithDataStoreCache struct {
	CantonChainFamilyAdapter
}

func (a *cantonChainFamilyWithDataStoreCache) GetRouterAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	cacheConfigureLanesDataStore(chainSelector, ds)

	return a.CantonChainFamilyAdapter.GetRouterAddress(ds, chainSelector)
}

func (a *cantonChainFamilyWithDataStoreCache) GetTestRouter(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	cacheConfigureLanesDataStore(chainSelector, ds)

	return a.CantonChainFamilyAdapter.GetTestRouter(ds, chainSelector)
}

func (a *cantonChainFamilyWithDataStoreCache) GetOnRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	cacheConfigureLanesDataStore(chainSelector, ds)

	return a.CantonChainFamilyAdapter.GetOnRampAddress(ds, chainSelector)
}

func (a *cantonChainFamilyWithDataStoreCache) GetOffRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	cacheConfigureLanesDataStore(chainSelector, ds)

	return a.CantonChainFamilyAdapter.GetOffRampAddress(ds, chainSelector)
}

func (a *cantonChainFamilyWithDataStoreCache) GetFQAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	cacheConfigureLanesDataStore(chainSelector, ds)

	return a.CantonChainFamilyAdapter.GetFQAddress(ds, chainSelector)
}

func (a *cantonChainFamilyWithDataStoreCache) ResolveExecutor(ds datastore.DataStore, chainSelector uint64, qualifier string) (string, error) {
	cacheConfigureLanesDataStore(chainSelector, ds)

	return a.CantonChainFamilyAdapter.ResolveExecutor(ds, chainSelector, qualifier)
}
