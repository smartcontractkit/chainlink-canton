package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
)

func TestFactoryAddressRef(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 123
	raw := contracts.InstanceID("core-factory").RawInstanceAddress("owner-party")

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
		Qualifier:     QualifierCCIP,
		Address:       raw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(raw.String()),
	}))
	require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
		Qualifier:     QualifierPools,
		Address:       "other",
		Labels:        datastore.NewLabelSet("pools-factory@owner-party"),
	}))

	got, err := FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCIP)
	require.NoError(t, err)
	require.Equal(t, QualifierCCIP, got.Qualifier)
	require.Equal(t, raw.InstanceAddress().String(), got.Address)

	_, err = FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCV)
	require.Error(t, err)
}

func TestFactoryAddressRefFromRefs(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 99
	refs := []datastore.AddressRef{{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
		Qualifier:     QualifierPools,
		Address:       "addr",
		Labels:        datastore.NewLabelSet("pools-factory@owner"),
	}}

	got, err := FactoryAddressRefFromRefs(chainSelector, QualifierPools, refs)
	require.NoError(t, err)
	require.Equal(t, QualifierPools, got.Qualifier)
}
