package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	utilstests "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
)

func GetContractAddress(
	t *testing.T,
	ds datastore.DataStore,
	chainSelector uint64,
	contractType datastore.ContractType,
	version, qualifier, contractName string,
) protocol.UnknownAddress {
	t.Helper()

	ref, err := ds.Addresses().Get(
		datastore.NewAddressRefKey(chainSelector, contractType, semver.MustParse(version), qualifier),
	)
	require.NoErrorf(t, err, "failed to get %s address for chain selector %d, ContractType: %s, ContractVersion: %s",
		contractName, chainSelector, contractType, version)

	addr, err := protocol.NewUnknownAddressFromHex(ref.Address)
	require.NoError(t, err)

	return addr
}

func GetChain(t *testing.T, chainType string, cfg *ccv.Cfg, lib ccv.Lib) cciptestinterfaces.CCIP17 {
	chainMap, err := lib.ChainsMap(t.Context())
	require.NoError(t, err)
	return GetChainFromMap(t, chainType, cfg, chainMap)
}

// GetChainFromMap returns a chain from an existing ChainsMap result. lib.ChainsMap
// constructs new impls on every call, so tests must reuse the same map they wired
// (via WireVerifierObservationFromLib) rather than calling ChainsMap again through GetChain.
func GetChainFromMap(
	t *testing.T,
	chainType string,
	cfg *ccv.Cfg,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
) cciptestinterfaces.CCIP17 {
	t.Helper()

	selector := chainSelectorForType(t, chainType, cfg)

	c, ok := chainMap[selector]
	require.True(t, ok, "chain selector %d not in ChainsMap", selector)

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

// verifierObservationAware is implemented by chain impls that need off-chain verifier
// clients inside ConfirmExecOnDest. EVM does not implement this; Canton does.
type verifierObservationAware interface {
	SetVerifierObservation(cantondevenv.VerifierObservation)
}

// WireVerifierObservationFromLib builds aggregator/indexer clients from lib and
// injects them into every chain that implements verifierObservationAware. Call once
// after lib.ChainsMap. Requires NewLibFromCCVEnv (not CLDF-only Lib).
func WireVerifierObservationFromLib(lib ccv.Lib, chains map[uint64]cciptestinterfaces.CCIP17) error {
	obs, err := cantondevenv.VerifierObservationFromLib(lib)
	if err != nil {
		return err
	}
	for _, c := range chains {
		if vo, ok := c.(verifierObservationAware); ok {
			vo.SetVerifierObservation(obs)
		}
	}

	return nil
}

func AssertSingleVerifierResult(
	t *testing.T,
	ctx context.Context,
	lib ccv.Lib,
	messageID [32]byte,
) tcapi.AssertionResult {
	t.Helper()

	obs, err := cantondevenv.VerifierObservationFromLib(lib)
	require.NoError(t, err)

	result, err := cantondevenv.AssertMessageWithVerifierObservation(ctx, obs, messageID, tcapi.AssertMessageOptions{
		TickInterval:            time.Second,
		Timeout:                 utilstests.WaitTimeout(t),
		ExpectedVerifierResults: 1,
		AssertVerifierLogs:      false,
		AssertExecutorLogs:      false,
	})
	require.NoError(t, err)
	require.Len(t, result.IndexedVerifications.Results, 1)

	return result
}
