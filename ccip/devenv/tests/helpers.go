package tests

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
)

// TODO: move this helper into ccv.Lib.
func GetChain(t *testing.T, chainType string, cfg *ccv.Cfg, harness tcapi.TestHarness) cciptestinterfaces.CCIP17 {
	var chain *blockchain.Input
	for _, bc := range cfg.Blockchains {
		if bc.Type == chainType {
			chain = bc
			break
		}
	}
	require.NotNil(t, chain, "need at least one chain for this test")

	var family string
	switch chainType {
	case blockchain.TypeCanton:
		family = chainsel.FamilyCanton
	case blockchain.TypeAnvil:
		family = chainsel.FamilyEVM
	default:
		t.Fatalf("unsupported chain type %q", chainType)
	}

	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(chain.ChainID, family)
	require.NoError(t, err)

	chainMap, err := harness.Lib.ChainsMap(t.Context())
	require.NoError(t, err)

	return chainMap[chainDetails.ChainSelector]
}
