package kick_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
)

func bundleWith(t *testing.T, r operations.Reporter) operations.Bundle {
	t.Helper()
	return operations.NewBundle(t.Context, logger.Nop(), r)
}

// TestState_ThresholdNotMet verifies state is populated on the DNS signing gate.
func TestState_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := baseInput()

	sr, err := operations.ExecuteSequence(bundleWith(t, reporter), kick.KickSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())

	s := sr.Output.State
	assert.Equal(t, kick.PhaseDNSSigning, s.Phase)
	assert.Equal(t, testKickedUID, s.KickedParticipant)
	assert.NotEmpty(t, s.ProposalHash)
	assert.Contains(t, s.CollectedSigners, "p1")
	assert.NotEmpty(t, s.PendingSigners)
	assert.False(t, sr.Output.DNSUpdated)
}

// TestState_Completed verifies the happy-path state snapshot.
func TestState_Completed(t *testing.T) {
	t.Parallel()

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, sharedReporter) }
	state := &mockState{}
	input := baseInput()

	// Run through the full ceremony (p1 sign, p2 sign+P2P, p1 P2P).
	_, _ = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p2", state), input)

	sr, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.NoError(t, err)

	s := sr.Output.State
	assert.Equal(t, kick.PhaseCompleted, s.Phase)
	assert.Empty(t, s.PendingSigners)
	assert.True(t, sr.Output.DNSUpdated)
	assert.True(t, sr.Output.P2PUpdated)
}

// TestState_LatestSequenceState verifies observer reads via LatestSequenceState.
func TestState_LatestSequenceState(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := baseInput()

	// No runs yet.
	_, found := kick.LatestSequenceState(reporter, input)
	assert.False(t, found)

	// Run once (partial).
	_, err := operations.ExecuteSequence(bundleWith(t, reporter), kick.KickSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())

	s, found := kick.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.Equal(t, kick.PhaseDNSSigning, s.Phase)
}

// TestState_Scoping verifies two kick ceremonies don't bleed.
func TestState_Scoping(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	inputA := baseInput()
	inputB := kick.KickInput{
		DecentralizedPartyID:       "other-party::1220112233",
		KickedParticipantID:        testKickedUID,
		KickedNamespaceFingerprint: testKickedFP,
		RemainingParticipants:      []string{"p1", "p2"},
		SynchronizerID:             testSyncID,
	}

	_, err := operations.ExecuteSequence(bundleWith(t, reporter), kick.KickSequence, newDeps("p1"), inputA)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())

	_, foundA := kick.LatestSequenceState(reporter, inputA)
	assert.True(t, foundA)

	_, foundB := kick.LatestSequenceState(reporter, inputB)
	assert.False(t, foundB)
}
