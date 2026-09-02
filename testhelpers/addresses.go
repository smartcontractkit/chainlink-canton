package testhelpers

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
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
	labels := addressRef.Labels.List()
	if len(labels) == 0 {
		return datastore.AddressRef{}, "", fmt.Errorf(
			"address ref %s/%s has no labels (raw instance address missing)",
			contractType, qualifier,
		)
	}
	rawInstanceAddress, err := contracts.RawInstanceAddressFromString(labels[0])
	if err != nil {
		return datastore.AddressRef{}, "", err
	}

	return addressRef, rawInstanceAddress, nil
}
