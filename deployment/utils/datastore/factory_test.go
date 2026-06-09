package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
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
		Qualifier:     QualifierCCV,
		Address:       "other",
		Labels:        datastore.NewLabelSet("ccv-factory@owner-party"),
	}))

	got, err := FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCIP)
	require.NoError(t, err)
	require.Equal(t, QualifierCCIP, got.Qualifier)
	require.Equal(t, raw.InstanceAddress().String(), got.Address)

	got, err = FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCV)
	require.NoError(t, err)
	require.Equal(t, QualifierCCV, got.Qualifier)

	_, err = FactoryAddressRef(ds.Seal(), chainSelector, "missing")
	require.Error(t, err)
}

func TestFactoryAddressRefFromRefs(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 99
	refs := []datastore.AddressRef{{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
		Qualifier:     QualifierCCV,
		Address:       "addr",
		Labels:        datastore.NewLabelSet("ccv-factory@owner"),
	}}

	got, err := FactoryAddressRefFromRefs(chainSelector, QualifierCCV, refs)
	require.NoError(t, err)
	require.Equal(t, QualifierCCV, got.Qualifier)
}

func TestFactoryAddressRefFallsBackToLegacyVersion(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 456
	raw := contracts.InstanceID("legacy-factory").RawInstanceAddress("owner-party")

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       semver.MustParse("0.1.0"),
		Qualifier:     QualifierCCIP,
		Address:       raw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(raw.String()),
	}))

	got, err := FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCIP)
	require.NoError(t, err)
	require.Equal(t, QualifierCCIP, got.Qualifier)
	require.Equal(t, raw.InstanceAddress().String(), got.Address)
}
