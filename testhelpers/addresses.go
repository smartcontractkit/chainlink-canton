package testhelpers

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func ResolveAddressFromDatastore(
	ds datastore.DataStore,
	chainsel uint64,
	contractType deployment.ContractType,
	version *semver.Version,
	qualifier string,
) (
	datastore.AddressRef,
	contracts.RawInstanceAddress,
	error,
) {
	addressRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainsel, datastore.ContractType(contractType), version, qualifier))
	if err != nil {
		return datastore.AddressRef{}, "", err
	}
	rawInstanceAddress, err := contracts.RawInstanceAddressFromString(addressRef.Labels.List()[0])
	if err != nil {
		return datastore.AddressRef{}, "", err
	}

	return addressRef, rawInstanceAddress, nil
}
