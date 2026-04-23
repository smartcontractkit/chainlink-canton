package addparticipant_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
)

func bundleWith(t *testing.T, r operations.Reporter) operations.Bundle {
	t.Helper()
	return operations.NewBundle(t.Context, logger.Nop(), r)
}

// TestState_ThresholdNotMet verifies state at the DNS signing gate.
func TestState_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	state := &mockState{}
	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseInput()

	// p3 (new) generates keys, NSD, reads state, creates DNS proposal, 0/2 sigs.
	sr, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error())

	s := sr.Output.State
	assert.Equal(t, addparticipant.PhaseDNSSigning, s.Phase)
	assert.True(t, s.NewMemberKeyReady)
	assert.True(t, s.NSDProposed)
	assert.NotEmpty(t, s.ProposalHash)
	assert.Empty(t, s.CollectedSigners)
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

	// Complete the ceremony.
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)

	sr, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.NoError(t, err)

	s := sr.Output.State
	assert.Equal(t, addparticipant.PhaseCompleted, s.Phase)
	assert.Empty(t, s.PendingSigners)
	assert.True(t, sr.Output.DNSUpdated)
	assert.True(t, sr.Output.P2PUpdated)
	assert.True(t, s.NewParticipantConsented)
}

// TestState_LatestSequenceState verifies observer reads.
func TestState_LatestSequenceState(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	state := &mockState{}
	input := baseInput()

	// No runs yet.
	_, found := addparticipant.LatestSequenceState(reporter, input)
	assert.False(t, found)

	// Run once (partial).
	_, err := operations.ExecuteSequence(bundleWith(t, reporter), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error())

	s, found := addparticipant.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.Equal(t, addparticipant.PhaseDNSSigning, s.Phase)
}

// TestState_Scoping verifies two add-participant ceremonies don't bleed.
func TestState_Scoping(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	state := &mockState{}
	inputA := baseInput()
	inputB := addparticipant.AddParticipantInput{
		DecentralizedPartyID: "other-party::1220112233",
		NewParticipantID:     "p4",
		NamespaceName:        "other-ns",
		SynchronizerID:       testSyncID,
	}

	_, err := operations.ExecuteSequence(bundleWith(t, reporter), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), inputA)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error())

	_, foundA := addparticipant.LatestSequenceState(reporter, inputA)
	assert.True(t, foundA)

	_, foundB := addparticipant.LatestSequenceState(reporter, inputB)
	assert.False(t, foundB)
}
