package datastore

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
)

// FeeQuoterExerciseAddrs holds Canton fee quoter addresses for on-ledger exercises.
type FeeQuoterExerciseAddrs struct {
	InstanceAddress    contracts.InstanceAddress
	RawInstanceAddress string
}

// ResolveFeeQuoterExerciseAddrs resolves fee quoter addresses from ChainDefinition.FeeQuoter
// bytes (via GetFQAddress in topology). MCMS paths also need the raw instance address label,
// looked up from ds by matching the instance-address hash.
func ResolveFeeQuoterExerciseAddrs(
	ds datastore.DataStore,
	chainSelector uint64,
	feeQuoterBytes []byte,
	mcmsEnabled bool,
) (FeeQuoterExerciseAddrs, error) {
	if len(feeQuoterBytes) == 0 {
		return FeeQuoterExerciseAddrs{}, fmt.Errorf("fee quoter address bytes are required")
	}

	addrs := FeeQuoterExerciseAddrs{
		InstanceAddress: contracts.BytesToInstanceAddress(feeQuoterBytes),
	}
	if !mcmsEnabled {
		return addrs, nil
	}
	if ds == nil {
		return FeeQuoterExerciseAddrs{}, fmt.Errorf("datastore is required to resolve fee quoter raw instance address for MCMS")
	}

	raw, err := GetRawInstanceAddressFromInstanceAddressBytes(
		ds,
		chainSelector,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		feeQuoterBytes,
	)
	if err != nil {
		return FeeQuoterExerciseAddrs{}, err
	}
	addrs.RawInstanceAddress = raw.String()

	return addrs, nil
}

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
