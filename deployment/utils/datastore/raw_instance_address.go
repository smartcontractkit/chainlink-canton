package datastore

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

func GetRawInstanceAddressFromAddressRef(addressRef datastore.AddressRef) (contracts.RawInstanceAddress, error) {
	labels := addressRef.Labels.List()
	if len(labels) == 0 {
		return "", fmt.Errorf("getting raw instance address from address ref: no labels found for ref: %s", addressRef.Address)
	}

	return contracts.RawInstanceAddressFromString(labels[0])
}
