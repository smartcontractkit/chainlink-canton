package changesets

import (
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain"

	"github.com/stretchr/testify/require"
)

const (
	cantonProdSelector  uint64 = 9268731218649498074
	sepoliaProdSelector uint64 = 16015286601757825753
)

func TestWithCantonProductionMinNOPCheckBypassed_inflatesOnlyForValidation(t *testing.T) {
	t.Parallel()

	topology := prodLikeTopologyWithFourCantonNOPs()
	cantonKey := selectorKey(cantonProdSelector)
	originalCommittee := append([]string(nil), topology.NOPTopology.Committees["default"].ChainConfigs[cantonKey].NOPAliases...)
	originalExecutorPool := append([]string(nil), topology.ExecutorPools["default"].ChainConfigs[cantonKey].NOPAliases...)

	err := WithCantonProductionMinNOPCheckBypassed(topology, func() error {
		inflatedCommittee := topology.NOPTopology.Committees["default"].ChainConfigs[cantonKey].NOPAliases
		require.GreaterOrEqual(t, uniqueTopologyNOPCount(inflatedCommittee), minProductionCantonChainNOPs)

		inflatedExecutorPool := topology.ExecutorPools["default"].ChainConfigs[cantonKey].NOPAliases
		require.GreaterOrEqual(t, uniqueTopologyNOPCount(inflatedExecutorPool), minProductionChainNOPs)
		require.NoError(t, topology.ValidateForEnvironment("prod_testnet"))

		return nil
	})
	require.NoError(t, err)

	restoredCommittee := topology.NOPTopology.Committees["default"].ChainConfigs[cantonKey].NOPAliases
	require.Equal(t, originalCommittee, restoredCommittee)
	restoredExecutorPool := topology.ExecutorPools["default"].ChainConfigs[cantonKey].NOPAliases
	require.Equal(t, originalExecutorPool, restoredExecutorPool)
	require.ErrorContains(t, topology.ValidateForEnvironment("prod_testnet"), "requires at least 9 unique NOPs")
}

func TestIsCantonCommitteeChainSelector(t *testing.T) {
	t.Parallel()

	require.True(t, isCantonCommitteeChainSelector(selectorKey(cantonProdSelector)))
	require.False(t, isCantonCommitteeChainSelector(selectorKey(sepoliaProdSelector)))
	require.False(t, isCantonCommitteeChainSelector("not-a-selector"))
}

func prodLikeTopologyWithFourCantonNOPs() *offchain.EnvironmentTopology {
	cantonNOPs := []string{
		"ccv-prod-testnet-0",
		"ccv-prod-testnet-1",
		"ccv-prod-testnet-2",
		"ccv-prod-testnet-3",
	}
	sepoliaNOPs := []string{
		"ccv-prod-testnet-0", "ccv-prod-testnet-1", "ccv-prod-testnet-2", "ccv-prod-testnet-3",
		"ccv-prod-testnet-4", "ccv-prod-testnet-5", "ccv-prod-testnet-6", "ccv-prod-testnet-7",
		"ccv-prod-testnet-8", "ccv-prod-testnet-9", "ccv-prod-testnet-10", "ccv-prod-testnet-11",
		"ccv-prod-testnet-12", "ccv-prod-testnet-13", "ccv-prod-testnet-14", "ccv-prod-testnet-15",
	}

	nops := make([]offchain.NOPConfig, 0, len(sepoliaNOPs))
	for _, alias := range sepoliaNOPs {
		nops = append(nops, offchain.NOPConfig{Alias: alias, Name: alias})
	}

	return &offchain.EnvironmentTopology{
		IndexerAddress: []string{"https://indexer.example"},
		NOPTopology: &offchain.NOPTopology{
			NOPs: nops,
			Committees: map[string]offchain.CommitteeConfig{
				"default": {
					Qualifier: "default",
					Aggregators: []offchain.AggregatorConfig{
						{Name: "agg", Address: "http://aggregator"},
					},
					ChainConfigs: map[string]offchain.ChainCommitteeConfig{
						selectorKey(cantonProdSelector): {
							NOPAliases: cantonNOPs,
							Threshold:  3,
						},
						selectorKey(sepoliaProdSelector): {
							NOPAliases: sepoliaNOPs,
							Threshold:  9,
						},
					},
				},
			},
		},
		ExecutorPools: map[string]offchain.ExecutorPoolConfig{
			"default": {
				ChainConfigs: map[string]offchain.ChainExecutorPoolConfig{
					selectorKey(cantonProdSelector): {
						NOPAliases: []string{"ccv-prod-testnet-0"},
					},
					selectorKey(sepoliaProdSelector): {
						NOPAliases: sepoliaNOPs,
					},
				},
			},
		},
	}
}

func selectorKey(sel uint64) string {
	return strconv.FormatUint(sel, 10)
}
