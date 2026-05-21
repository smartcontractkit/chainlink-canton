package tests

import (
	"context"
	"testing"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	utilstests "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
)

// TODO: move this helper into ccv.Lib.
func GetChain(t *testing.T, chainType string, cfg *ccv.Cfg, lib ccv.Lib) cciptestinterfaces.CCIP17 {
	chainMap, err := lib.ChainsMap(t.Context())
	require.NoError(t, err)
	return GetChainFromMap(t, chainType, cfg, lib, chainMap)
}

// GetChainFromMap returns a chain from an existing ChainsMap result. lib.ChainsMap
// constructs new impls on every call, so tests must reuse the same map they wired
// (via WireLibIntoChains) rather than calling ChainsMap again through GetChain.
func GetChainFromMap(
	t *testing.T,
	chainType string,
	cfg *ccv.Cfg,
	lib ccv.Lib,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
) cciptestinterfaces.CCIP17 {
	t.Helper()

	selector := chainSelectorForType(t, chainType, cfg)

	c, ok := chainMap[selector]
	require.True(t, ok, "chain selector %d not in ChainsMap", selector)
	if la, ok := c.(libAware); ok {
		la.SetLib(lib)
	}

	return c
}

func chainSelectorForType(t *testing.T, chainType string, cfg *ccv.Cfg) uint64 {
	t.Helper()

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

	return chainDetails.ChainSelector
}

// libAware is implemented by chain impls that need a back-reference to ccv.Lib
// (e.g. to fetch verifier results inside ConfirmExecOnDest). EVM does not
// implement this; Canton does.
type libAware interface {
	SetLib(ccv.Lib)
}

// WireLibIntoChains injects lib into every chain that implements libAware. Test
// runners must call this once after lib.ChainsMap and before invoking
// chain methods that depend on the lib (today: Canton's ConfirmExecOnDest).
func WireLibIntoChains(lib ccv.Lib, chains map[uint64]cciptestinterfaces.CCIP17) {
	for _, c := range chains {
		if la, ok := c.(libAware); ok {
			la.SetLib(lib)
		}
	}
}

func AssertSingleVerifierResult(
	t *testing.T,
	ctx context.Context,
	lib ccv.Lib,
	messageID [32]byte,
) tcapi.AssertionResult {
	t.Helper()

	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)

	aggregatorClients, err := lib.AllAggregators()
	require.NoError(t, err)
	aggregatorClient := aggregatorClients[common.DefaultCommitteeVerifierQualifier]
	indexerMonitor, err := lib.IndexerMonitor()
	require.NoError(t, err)

	testCtx, cleanupFn := tcapi.NewTestingContext(
		ctx,
		chainMap,
		aggregatorClient,
		indexerMonitor,
	)
	defer cleanupFn()

	result, err := testCtx.AssertMessage(messageID, tcapi.AssertMessageOptions{
		TickInterval:            time.Second,
		Timeout:                 utilstests.WaitTimeout(t),
		ExpectedVerifierResults: 1,
		AssertVerifierLogs:      false,
		AssertExecutorLogs:      false,
	})
	require.NoError(t, err)
	require.NotNil(t, result.AggregatedResult)
	require.Len(t, result.IndexedVerifications.Results, 1)

	return result
}
