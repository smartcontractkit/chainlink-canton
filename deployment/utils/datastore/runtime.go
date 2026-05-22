package datastore

import (
	"sync"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

var (
	runtimeDataStoreMu sync.RWMutex
	runtimeDataStore   datastore.DataStore
)

// SetRuntimeDataStore sets the datastore used during Canton lane configure sequences
// to resolve raw instance addresses from contract address bytes.
func SetRuntimeDataStore(ds datastore.DataStore) {
	runtimeDataStoreMu.Lock()
	defer runtimeDataStoreMu.Unlock()
	runtimeDataStore = ds
}

// RuntimeDataStore returns the datastore set by SetRuntimeDataStore.
func RuntimeDataStore() datastore.DataStore {
	runtimeDataStoreMu.RLock()
	defer runtimeDataStoreMu.RUnlock()

	return runtimeDataStore
}
