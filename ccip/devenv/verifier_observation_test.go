package devenv

import (
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/stretchr/testify/require"
)

func TestVerifierObservation_wired_indexerOnly(t *testing.T) {
	t.Parallel()

	obs := VerifierObservation{
		IndexerMonitor: &ccv.IndexerMonitor{},
	}
	require.True(t, obs.wired())
}

func TestVerifierObservation_wired_indexerNil(t *testing.T) {
	t.Parallel()

	obs := VerifierObservation{
		AggregatorClient: &ccv.AggregatorClient{},
	}
	require.False(t, obs.wired())
}
