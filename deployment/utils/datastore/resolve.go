package datastore

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

// AddressRefFromInstanceAddressBytes finds the datastore ref whose Address matches the
// given 32-byte Canton instance address hash.
func AddressRefFromInstanceAddressBytes(
	ds datastore.DataStore,
	chainSelector uint64,
	contractType datastore.ContractType,
	version *semver.Version,
	instanceAddressBytes []byte,
) (datastore.AddressRef, error) {
	if ds == nil {
		return datastore.AddressRef{}, fmt.Errorf("datastore is nil")
	}
	if len(instanceAddressBytes) == 0 {
		return datastore.AddressRef{}, fmt.Errorf("instance address bytes are empty")
	}

	hexAddr := contracts.BytesToInstanceAddress(instanceAddressBytes).Hex()
	refs, err := ds.Addresses().Fetch()
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("fetch addresses: %w", err)
	}
	for _, ref := range refs {
		if ref.ChainSelector == chainSelector &&
			ref.Type == contractType &&
			ref.Version.Equal(version) &&
			ref.Address == hexAddr {
			return ref, nil
		}
	}

	return datastore.AddressRef{}, fmt.Errorf(
		"no %s address ref found for instance address %s on chain %d",
		contractType,
		hexAddr,
		chainSelector,
	)
}

// GetRawInstanceAddressFromInstanceAddressBytes resolves instanceId@party from the
// datastore label for a contract identified by its instance-address bytes.
//
// Used when lane input carries only ChainDefinition.FeeQuoter ([]byte hash) but on-ledger
// or MCMS paths need RawInstanceAddress; the datastore ref is looked up by matching Address.
func GetRawInstanceAddressFromInstanceAddressBytes(
	ds datastore.DataStore,
	chainSelector uint64,
	contractType datastore.ContractType,
	version *semver.Version,
	instanceAddressBytes []byte,
) (contracts.RawInstanceAddress, error) {
	ref, err := AddressRefFromInstanceAddressBytes(ds, chainSelector, contractType, version, instanceAddressBytes)
	if err != nil {
		return "", err
	}

	return GetRawInstanceAddressFromAddressRef(ref)
}
