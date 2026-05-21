package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
)

func TestFactoryAddressRef(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 8706591216959472610

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		Address:       "0xabc",
		ChainSelector: chainSelector,
		Labels:        datastore.NewLabelSet("factory-ccip@party::1220"),
		Qualifier:     QualifierCCIP,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
	}))

	got, err := FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCIP)
	require.NoError(t, err)
	require.Equal(t, QualifierCCIP, got.Qualifier)

	_, err = FactoryAddressRef(ds.Seal(), chainSelector, QualifierCCV)
	require.Error(t, err)
}

func TestFactoryAddressRefFromRefs(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 8706591216959472610

	refs := []datastore.AddressRef{{
		Address:       "0xdef",
		ChainSelector: chainSelector,
		Labels:        datastore.NewLabelSet("factory-ccv@party::1220"),
		Qualifier:     QualifierCCV,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
	}}

	got, err := FactoryAddressRefFromRefs(chainSelector, QualifierCCV, refs)
	require.NoError(t, err)
	require.Equal(t, QualifierCCV, got.Qualifier)
}
