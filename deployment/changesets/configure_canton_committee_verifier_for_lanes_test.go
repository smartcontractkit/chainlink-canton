package changesets

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCantonOnlyConfigureLanesConfig(t *testing.T) {
	t.Parallel()

	cfg := v2cs.ConfigureChainsForLanesFromTopologyConfig{
		Chains: []v2cs.PartialChainConfig{
			{ChainSelector: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector},
			{ChainSelector: chainsel.CANTON_TESTNET.Selector},
		},
	}

	filtered := cantonOnlyConfigureLanesConfig(cfg)
	require.Len(t, filtered.Chains, 1)
	assert.Equal(t, chainsel.CANTON_TESTNET.Selector, filtered.Chains[0].ChainSelector)
}

func TestFamiliesInConfigureLanesConfig(t *testing.T) {
	t.Parallel()

	cfg := v2cs.ConfigureChainsForLanesFromTopologyConfig{
		Chains: []v2cs.PartialChainConfig{
			{
				ChainSelector: chainsel.CANTON_TESTNET.Selector,
				CommitteeVerifiers: []v2cs.CommitteeVerifierInputConfig{
					{
						RemoteChains: map[uint64]v2cs.CommitteeVerifierRemoteChainConfig{
							chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {},
						},
					},
				},
			},
		},
	}

	families := familiesInConfigureLanesConfig(cfg)
	assert.Contains(t, families, chainsel.FamilyCanton)
	assert.Contains(t, families, chainsel.FamilyEVM)
}

func TestHasCantonChainsInConfig(t *testing.T) {
	t.Parallel()

	cfg := v2cs.ConfigureChainsForLanesFromTopologyConfig{
		Chains: []v2cs.PartialChainConfig{
			{ChainSelector: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector},
		},
	}
	assert.False(t, hasCantonChainsInConfig(cfg))

	cfg.Chains = append(cfg.Chains, v2cs.PartialChainConfig{
		ChainSelector: chainsel.CANTON_TESTNET.Selector,
	})
	assert.True(t, hasCantonChainsInConfig(cfg))
}
