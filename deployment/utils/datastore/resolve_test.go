package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestGetRawInstanceAddressFromInstanceAddressBytes(t *testing.T) {
	t.Parallel()

	raw := contracts.NewRawInstanceAddress(contracts.InstanceID("feequoter-test"), types.PARTY("ccip-owner"))
	ref := datastore.AddressRef{
		Address:       raw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(raw.String()),
		ChainSelector: 42,
		Type:          datastore.ContractType("FeeQuoter"),
		Version:       semver.MustParse("2.0.0"),
	}
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(ref))

	got, err := GetRawInstanceAddressFromInstanceAddressBytes(
		ds.Seal(),
		42,
		datastore.ContractType("FeeQuoter"),
		semver.MustParse("2.0.0"),
		raw.InstanceAddress().Bytes(),
	)
	require.NoError(t, err)
	require.Equal(t, raw, got)
}
