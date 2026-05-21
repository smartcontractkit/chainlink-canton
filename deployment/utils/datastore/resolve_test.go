package datastore

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
)

func TestGetRawInstanceAddressFromInstanceAddressBytes(t *testing.T) {
	t.Parallel()

	raw := contracts.NewRawInstanceAddress(contracts.InstanceID("feequoter-test"), types.PARTY("ccip-owner"))
	ref := datastore.AddressRef{
		Address:       raw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(raw.String()),
		ChainSelector: 42,
		Type:          datastore.ContractType(feequoterop.ContractType),
		Version:       feequoterop.Version,
	}
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(ref))

	got, err := GetRawInstanceAddressFromInstanceAddressBytes(
		ds.Seal(),
		42,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		raw.InstanceAddress().Bytes(),
	)
	require.NoError(t, err)
	require.Equal(t, raw, got)
}

func TestResolveFeeQuoterExerciseAddrsMCMS(t *testing.T) {
	t.Parallel()

	raw := contracts.NewRawInstanceAddress(contracts.InstanceID("feequoter-test"), types.PARTY("ccip-owner"))
	ref := datastore.AddressRef{
		Address:       raw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(raw.String()),
		ChainSelector: 42,
		Type:          datastore.ContractType(feequoterop.ContractType),
		Version:       feequoterop.Version,
	}
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(ref))

	addrs, err := ResolveFeeQuoterExerciseAddrs(ds.Seal(), 42, raw.InstanceAddress().Bytes(), true)
	require.NoError(t, err)
	require.Equal(t, raw.InstanceAddress(), addrs.InstanceAddress)
	require.Equal(t, raw.String(), addrs.RawInstanceAddress)
}

func TestResolveFeeQuoterExerciseAddrsDirect(t *testing.T) {
	t.Parallel()

	raw := contracts.NewRawInstanceAddress(contracts.InstanceID("feequoter-test"), types.PARTY("ccip-owner"))

	addrs, err := ResolveFeeQuoterExerciseAddrs(nil, 42, raw.InstanceAddress().Bytes(), false)
	require.NoError(t, err)
	require.Equal(t, raw.InstanceAddress(), addrs.InstanceAddress)
	require.Empty(t, addrs.RawInstanceAddress)
}
