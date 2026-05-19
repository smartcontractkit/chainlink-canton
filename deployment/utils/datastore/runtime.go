package datastore

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

var runtimeDataStoreGetter func() datastore.DataStore

// RegisterRuntimeDataStoreGetter wires the runtime datastore owned by deployment/adapters.
// Canton sequences call RuntimeDataStore; adapters.SetRuntimeDataStore registers the getter
// so sequences do not import adapters (that would create an import cycle).
func RegisterRuntimeDataStoreGetter(fn func() datastore.DataStore) {
	runtimeDataStoreGetter = fn
}

// RuntimeDataStore returns the datastore registered by adapters.SetRuntimeDataStore.
func RuntimeDataStore() datastore.DataStore {
	if runtimeDataStoreGetter == nil {
		return nil
	}

	return runtimeDataStoreGetter()
}
