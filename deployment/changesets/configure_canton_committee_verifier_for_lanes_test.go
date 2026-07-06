package changesets

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain"
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

func TestTopologyForCantonCommitteeVerifierConfigure_DropsEVMOnlySecondary(t *testing.T) {
	t.Parallel()

	cantonSel := chainsel.CANTON_TESTNET.Selector
	sepoliaSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	cantonKey := "10109143320554840099"
	sepoliaKey := "16015286601757825753"

	cfg := v2cs.ConfigureChainsForLanesFromTopologyConfig{
		Topology: &offchain.EnvironmentTopology{
			NOPTopology: &offchain.NOPTopology{
				Committees: map[string]offchain.CommitteeConfig{
					"default": {
						Qualifier: "default",
						ChainConfigs: map[string]offchain.ChainCommitteeConfig{
							cantonKey:  {NOPAliases: []string{"nop-1"}, Threshold: 1},
							sepoliaKey: {NOPAliases: []string{"nop-1", "nop-2"}, Threshold: 2},
						},
					},
					"secondary": {
						Qualifier: "secondary",
						ChainConfigs: map[string]offchain.ChainCommitteeConfig{
							sepoliaKey: {NOPAliases: []string{"nop-3"}, Threshold: 1},
						},
					},
				},
			},
		},
		BuildLanesCrossFamilyConfig: v2cs.BuildLanesCrossFamilyConfig{
			Lanes: []v2cs.CrossFamilyLanePair{
				{ChainA: cantonSel, ChainB: sepoliaSel},
			},
		},
	}

	filtered := topologyForCantonCommitteeVerifierConfigure(cfg)
	require.NotNil(t, filtered.Topology)
	require.NotNil(t, filtered.Topology.NOPTopology)

	_, hasDefault := filtered.Topology.NOPTopology.Committees["default"]
	_, hasSecondary := filtered.Topology.NOPTopology.Committees["secondary"]
	assert.True(t, hasDefault)
	assert.False(t, hasSecondary)
}

func TestKeepCommitteeForCantonCommitteeVerifierConfigure(t *testing.T) {
	t.Parallel()

	cantonSel := chainsel.CANTON_TESTNET.Selector
	sepoliaSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	cantonKey := "10109143320554840099"
	sepoliaKey := "16015286601757825753"
	remotes := remotesByCantonSelector(
		[]v2cs.CrossFamilyLanePair{{ChainA: cantonSel, ChainB: sepoliaSel}},
		[]uint64{cantonSel},
	)

	assert.True(t, keepCommitteeForCantonCommitteeVerifierConfigure(
		offchain.CommitteeConfig{
			Qualifier: "default",
			ChainConfigs: map[string]offchain.ChainCommitteeConfig{
				cantonKey:  {},
				sepoliaKey: {},
			},
		},
		[]uint64{cantonSel},
		remotes,
	))

	assert.False(t, keepCommitteeForCantonCommitteeVerifierConfigure(
		offchain.CommitteeConfig{
			Qualifier: "secondary",
			ChainConfigs: map[string]offchain.ChainCommitteeConfig{
				sepoliaKey: {},
			},
		},
		[]uint64{cantonSel},
		remotes,
	))
}
