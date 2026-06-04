package devenv

import (
	"context"
	"fmt"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

// VerifierObservation holds off-chain clients used to wait for CCIP verifier
// results (aggregator + indexer). ConfirmExecOnDest needs these; it does not
// need ChainsMap or other Lib methods.
//
// Build from a CCV env Lib via [VerifierObservationFromLib] (requires
// [ccv.NewLibFromCCVEnv] — CLDF-only Lib backends cannot provide aggregator/indexer).
type VerifierObservation struct {
	AggregatorClient *ccv.AggregatorClient
	IndexerMonitor   *ccv.IndexerMonitor
}

func (o VerifierObservation) wired() bool {
	return o.AggregatorClient != nil && o.IndexerMonitor != nil
}

// VerifierObservationFromLib extracts aggregator and indexer clients from lib.
func VerifierObservationFromLib(lib ccv.Lib) (VerifierObservation, error) {
	if lib == nil {
		return VerifierObservation{}, fmt.Errorf("VerifierObservationFromLib: lib is nil")
	}

	aggregatorClients, err := lib.AllAggregators()
	if err != nil {
		return VerifierObservation{}, fmt.Errorf("all aggregators: %w", err)
	}
	aggregatorClient, ok := aggregatorClients[devenvcommon.DefaultCommitteeVerifierQualifier]
	if !ok || aggregatorClient == nil {
		return VerifierObservation{}, fmt.Errorf("no aggregator client for qualifier %q", devenvcommon.DefaultCommitteeVerifierQualifier)
	}

	indexerMonitor, err := lib.IndexerMonitor()
	if err != nil {
		return VerifierObservation{}, fmt.Errorf("indexer monitor: %w", err)
	}

	return VerifierObservation{
		AggregatorClient: aggregatorClient,
		IndexerMonitor:   indexerMonitor,
	}, nil
}

// AssertMessageWithVerifierObservation waits for verifier results for messageID
// using aggregator and indexer only (no chain map).
func AssertMessageWithVerifierObservation(
	ctx context.Context,
	obs VerifierObservation,
	messageID protocol.Bytes32,
	opts tcapi.AssertMessageOptions,
) (tcapi.AssertionResult, error) {
	if !obs.wired() {
		return tcapi.AssertionResult{}, fmt.Errorf("verifier observation not wired (aggregator and indexer required)")
	}

	testCtx, cleanupFn := tcapi.NewTestingContext(ctx, nil, obs.AggregatorClient, obs.IndexerMonitor)
	defer cleanupFn()

	return testCtx.AssertMessage([32]byte(messageID), opts)
}

// SetVerifierObservation injects off-chain verifier clients for ConfirmExecOnDest.
// Test runners typically call [WireVerifierObservationFromLib] once after lib.ChainsMap.
func (c *Chain) SetVerifierObservation(obs VerifierObservation) {
	c.verifierObs = obs
}
