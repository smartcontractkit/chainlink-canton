package changesets

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/stretchr/testify/assert"
)

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
