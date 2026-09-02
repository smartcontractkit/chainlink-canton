package datastore

import (
	"fmt"

	datastore2 "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

func ToInstanceAddress(ref datastore.AddressRef) (contracts.InstanceAddress, error) {
	if ref.Address == "" {
		return contracts.InstanceAddress{}, fmt.Errorf("address is empty in ref: %s", datastore2.SprintRef(ref))
	}

	return contracts.HexToInstanceAddress(ref.Address), nil
}

func ToInstanceAddressBytes(ref datastore.AddressRef) ([]byte, error) {
	addr, err := ToInstanceAddress(ref)
	if err != nil {
		return nil, err
	}

	return addr.Bytes(), nil
}
