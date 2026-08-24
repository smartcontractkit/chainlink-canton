package adapters

import (
	"math/big"
	"testing"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
)

func TestResolveTokenPricesForRemoteDest_DefaultLinkToken(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 9268731218649498074
	ccipOwner := "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(feequoterop.ContractType),
		Version:       feequoterop.Version,
		Qualifier:     "",
		Address:       "0xabc",
		Labels:        datastore.NewLabelSet("feequoter-scxln@" + ccipOwner),
	}))

	prices, err := ResolveTokenPricesForRemoteDest(ds.Seal(), ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: chainSelector,
	}, 16015286601757825753, nil)
	require.NoError(t, err)
	require.Equal(t, map[string]*big.Int{
		ccipOwner + ":link-token": usdPerTokenToScaled(defaultLinkUsdPerTokenDollars),
	}, prices)
}

func TestResolveTokenPricesForRemoteDest_IncludesNativeToken(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 9268731218649498074
	ccipOwner := "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"
	dso := "DSO::1220be58c29e65de40bf273be1dc2b266d43a9a002ea5b18955aeef7aac881bb471a"

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(feequoterop.ContractType),
		Version:       feequoterop.Version,
		Qualifier:     "",
		Address:       "0xabc",
		Labels:        datastore.NewLabelSet("feequoter-scxln@" + ccipOwner),
	}))

	native := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(dso),
		Id:    types.TEXT("Amulet"),
	}
	prices, err := ResolveTokenPricesForRemoteDest(ds.Seal(), ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: chainSelector,
	}, 16015286601757825753, &native)
	require.NoError(t, err)
	nativePrice, err := parseUsdPerTokenPriceString(defaultNativeUsdPerToken)
	require.NoError(t, err)
	require.Equal(t, map[string]*big.Int{
		ccipOwner + ":link-token": usdPerTokenToScaled(defaultLinkUsdPerTokenDollars),
		dso + ":Amulet":           nativePrice,
	}, prices)
}

func TestTokenPricesFromFamilyExtras_Override(t *testing.T) {
	t.Parallel()

	prices, err := tokenPricesFromFamilyExtras(map[string]any{
		CantonRemoteTokenPricesFamilyExtraKey: map[string]any{
			"16015286601757825753": map[string]any{
				"ccipOwner::1220:link-token": "15",
			},
		},
	}, 16015286601757825753)
	require.NoError(t, err)
	require.Equal(t, usdPerTokenToScaled(15), prices["ccipOwner::1220:link-token"])
}

func TestPartyFromDeployLabels(t *testing.T) {
	t.Parallel()

	party, ok := partyFromDeployLabels(datastore.NewLabelSet(
		"feequoter-scxln@ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551",
	))
	require.True(t, ok)
	require.Equal(t, "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551", party)
}
