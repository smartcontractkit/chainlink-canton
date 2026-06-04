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
		BuildLanesCrossFamilyConfig: v2cs.BuildLanesCrossFamilyConfig{
			Lanes: []v2cs.CrossFamilyLanePair{
				{
					ChainA: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
					ChainB: chainsel.CANTON_TESTNET.Selector,
				},
				{
					ChainA: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
					ChainB: chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
				},
			},
		},
	}

	filtered := cantonOnlyConfigureLanesConfig(cfg)
	require.Len(t, filtered.Lanes, 1)
	assert.Equal(t, chainsel.CANTON_TESTNET.Selector, filtered.Lanes[0].ChainB)
}

func TestFamiliesInConfigureLanesConfig(t *testing.T) {
	t.Parallel()

	cfg := v2cs.ConfigureChainsForLanesFromTopologyConfig{
		BuildLanesCrossFamilyConfig: v2cs.BuildLanesCrossFamilyConfig{
			Lanes: []v2cs.CrossFamilyLanePair{
				{
					ChainA: chainsel.CANTON_TESTNET.Selector,
					ChainB: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
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
		BuildLanesCrossFamilyConfig: v2cs.BuildLanesCrossFamilyConfig{
			Lanes: []v2cs.CrossFamilyLanePair{
				{
					ChainA: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
					ChainB: chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
				},
			},
		},
	}
	assert.False(t, hasCantonChainsInConfig(cfg))

	cfg.Lanes = append(cfg.Lanes, v2cs.CrossFamilyLanePair{
		ChainA: chainsel.CANTON_TESTNET.Selector,
		ChainB: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
	})
	assert.True(t, hasCantonChainsInConfig(cfg))
}
