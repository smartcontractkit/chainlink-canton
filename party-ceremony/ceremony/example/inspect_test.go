package example_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	ceremony "github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/example"
)

func bundleWith(t *testing.T, r operations.Reporter) operations.Bundle {
	t.Helper()
	return operations.NewBundle(t.Context, logger.Nop(), r)
}

// TestState_ThresholdNotMet: the direct caller reads sr.Output.State after
// a partial-signing run, exactly like reading queryResult in Temporal.
func TestState_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := ceremony.OnboardingInput{
		NamespaceName:  "state-threshold",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "threshold-party",
		SynchronizerID: "global",
		Threshold:      2,
	}

	// p1 runs; their signing succeeds but p2/p3 don't match → ErrThresholdNotMet.
	sr, err := operations.ExecuteSequence(bundleWith(t, reporter), ceremony.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error())

	// State is fully populated even on the error path.
	s := sr.Output.State
	assert.Equal(t, ceremony.PhaseSigning, s.Phase)
	assert.Equal(t, 2, s.Threshold)
	assert.ElementsMatch(t, []string{"p1", "p2", "p3"}, s.InitializedMembers)
	assert.Equal(t, []string{"p1"}, s.CollectedSigners)
	assert.ElementsMatch(t, []string{"p2", "p3"}, s.PendingSigners)
	assert.NotEmpty(t, s.ProposalHash)
	assert.Empty(t, sr.Output.PartyID) // not yet completed
}

// TestState_ResumePartialSigning: after each actor's run, state advances.
// Read directly from the SequenceReport — no reporter query needed.
func TestState_ResumePartialSigning(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }

	input := ceremony.OnboardingInput{
		NamespaceName:  "state-partial",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "partial-party",
		SynchronizerID: "global",
		Threshold:      3,
	}

	// Run 1: only p1 signs.
	sr1, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error())
	assert.Equal(t, ceremony.PhaseSigning, sr1.Output.State.Phase)
	assert.Equal(t, []string{"p1"}, sr1.Output.State.CollectedSigners)
	assert.ElementsMatch(t, []string{"p2", "p3"}, sr1.Output.State.PendingSigners)

	// Run 2: p2 signs.
	sr2, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p2"), input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error())
	assert.Equal(t, ceremony.PhaseSigning, sr2.Output.State.Phase)
	assert.ElementsMatch(t, []string{"p1", "p2"}, sr2.Output.State.CollectedSigners)
	assert.Equal(t, []string{"p3"}, sr2.Output.State.PendingSigners)

	// Run 3: p3 signs → threshold met → ceremony completes.
	sr3, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p3"), input)
	require.NoError(t, err)
	assert.Equal(t, ceremony.PhaseCompleted, sr3.Output.State.Phase)
	assert.ElementsMatch(t, []string{"p1", "p2", "p3"}, sr3.Output.State.CollectedSigners)
	assert.Empty(t, sr3.Output.State.PendingSigners)
	assert.NotEmpty(t, sr3.Output.PartyID)
	assert.True(t, sr3.Output.DNSConfirmed)
}

// TestState_Completed: happy-path state is PhaseCompleted with all fields set.
func TestState_Completed(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := ceremony.OnboardingInput{
		NamespaceName:  "state-completed",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "completed-party",
		SynchronizerID: "global",
		Threshold:      1,
	}

	sr, err := operations.ExecuteSequence(bundleWith(t, reporter), ceremony.OnboardingSequence, newDeps("p1"), input)
	require.NoError(t, err)

	s := sr.Output.State
	assert.Equal(t, ceremony.PhaseCompleted, s.Phase)
	assert.ElementsMatch(t, []string{"p1", "p2", "p3"}, s.InitializedMembers)
	assert.NotEmpty(t, s.ProposalHash)
	assert.Contains(t, s.CollectedSigners, "p1")
	assert.Empty(t, s.PendingSigners)
}

// TestState_LatestSequenceState: an observer (not the direct caller of
// ExecuteSequence) reads the same state via LatestSequenceState.
func TestState_LatestSequenceState(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }

	input := ceremony.OnboardingInput{
		NamespaceName:  "state-observer",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "observer-party",
		SynchronizerID: "global",
		Threshold:      3,
	}

	// No runs yet — observer sees nothing.
	_, found := ceremony.LatestSequenceState(reporter, input)
	assert.False(t, found)

	// p1 runs.
	_, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error())

	// Observer reads without having called ExecuteSequence itself.
	s, found := ceremony.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.Equal(t, ceremony.PhaseSigning, s.Phase)
	assert.Equal(t, []string{"p1"}, s.CollectedSigners)

	// p2 runs; observer re-reads and sees the updated state.
	_, err = operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p2"), input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error())

	s, found = ceremony.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.ElementsMatch(t, []string{"p1", "p2"}, s.CollectedSigners)
	assert.Equal(t, []string{"p3"}, s.PendingSigners)
}

// TestState_Scoping: two ceremonies on the same reporter don't bleed into
// each other's LatestSequenceState results.
func TestState_Scoping(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }

	inputA := ceremony.OnboardingInput{
		NamespaceName: "scope-a", Participants: []string{"p1"}, PartyName: "party-a",
		SynchronizerID: "global", Threshold: 1,
	}
	inputB := ceremony.OnboardingInput{
		NamespaceName: "scope-b", Participants: []string{"p2"}, PartyName: "party-b",
		SynchronizerID: "global", Threshold: 1,
	}

	// Only ceremony A runs.
	_, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p1"), inputA)
	require.NoError(t, err)

	sA, foundA := ceremony.LatestSequenceState(reporter, inputA)
	assert.True(t, foundA)
	assert.Equal(t, ceremony.PhaseCompleted, sA.Phase)

	_, foundB := ceremony.LatestSequenceState(reporter, inputB)
	assert.False(t, foundB, "ceremony B should have no state")
}
