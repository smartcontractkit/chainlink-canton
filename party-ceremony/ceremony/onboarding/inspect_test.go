package onboarding_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/chainlink/canton-party-ceremony/ceremony/onboarding"
)

func bundleWith(t *testing.T, r operations.Reporter) operations.Bundle {
	t.Helper()
	return operations.NewBundle(t.Context, logger.Nop(), r)
}

// TestState_ThresholdNotMet verifies that state is populated on the error path
// when the key-gen gate fires.
func TestState_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := multiActorInput()

	sr, err := operations.ExecuteSequence(bundleWith(t, reporter), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())

	s := sr.Output.State
	// p1 generates its own key, but p2's key is missing - stuck at key-gen/NSD.
	assert.Contains(t, []onboarding.Phase{onboarding.PhaseKeyGen, onboarding.PhaseNSD}, s.Phase)
	assert.Equal(t, 2, s.Threshold)
	assert.Contains(t, s.KeysGenerated, "p1")
	assert.Empty(t, sr.Output.PartyID)
}

// TestState_ResumePartialSigning exercises multi-actor resume and verifies
// state advances at each step.
func TestState_ResumePartialSigning(t *testing.T) {
	t.Parallel()

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, sharedReporter) }
	input := multiActorInput()

	// Run 1: p1 generates key + NSD - gate fires (1/2 keys).
	sr1, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	assert.Contains(t, sr1.Output.State.KeysGenerated, "p1")

	// Run 2: p2 generates key + NSD - DNS proposal - p2 signs (1/2) - threshold not met.
	sr2, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p2"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	assert.Equal(t, onboarding.PhaseDNSSigning, sr2.Output.State.Phase)
	assert.NotEmpty(t, sr2.Output.State.ProposalHash)
	assert.Equal(t, []string{"p2"}, sr2.Output.State.CollectedSigners)
	assert.Equal(t, []string{"p1"}, sr2.Output.State.PendingSigners)

	// Run 3: p1 signs (2/2) - DNS submitted - p1 proposes P2P (1/2) - threshold not met.
	sr3, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	assert.Equal(t, onboarding.PhaseP2P, sr3.Output.State.Phase)
	assert.ElementsMatch(t, []string{"p1", "p2"}, sr3.Output.State.CollectedSigners)

	// Run 4: p2 proposes P2P (2/2) - confirmed - completed.
	sr4, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p2"), input)
	require.NoError(t, err)
	assert.Equal(t, onboarding.PhaseCompleted, sr4.Output.State.Phase)
	assert.True(t, sr4.Output.DNSConfirmed)
	assert.True(t, sr4.Output.P2PConfirmed)
	assert.NotEmpty(t, sr4.Output.PartyID)
}

// TestState_Completed verifies the happy-path state snapshot.
func TestState_Completed(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := baseInput()

	sr, err := operations.ExecuteSequence(bundleWith(t, reporter), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.NoError(t, err)

	s := sr.Output.State
	assert.Equal(t, onboarding.PhaseCompleted, s.Phase)
	assert.Contains(t, s.KeysGenerated, "p1")
	assert.NotEmpty(t, s.ProposalHash)
	assert.Contains(t, s.CollectedSigners, "p1")
	assert.Empty(t, s.PendingSigners)
}

// TestState_LatestSequenceState verifies observer reads via LatestSequenceState.
func TestState_LatestSequenceState(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseInput()

	// No runs yet.
	_, found := onboarding.LatestSequenceState(reporter, input)
	assert.False(t, found)

	// Run once.
	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.NoError(t, err)

	s, found := onboarding.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.Equal(t, onboarding.PhaseCompleted, s.Phase)
}

// TestState_Scoping verifies two ceremonies don't bleed into each other.
func TestState_Scoping(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	inputA := onboarding.OnboardingInput{
		PartyPrefix:    "party-a",
		NamespaceName:  "scope-a",
		Threshold:      1,
		SynchronizerID: "global",
		Participants:   []string{"p1"},
	}
	inputB := onboarding.OnboardingInput{
		PartyPrefix:    "party-b",
		NamespaceName:  "scope-b",
		Threshold:      1,
		SynchronizerID: "global",
		Participants:   []string{"p2"},
	}

	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), inputA)
	require.NoError(t, err)

	sA, foundA := onboarding.LatestSequenceState(reporter, inputA)
	assert.True(t, foundA)
	assert.Equal(t, onboarding.PhaseCompleted, sA.Phase)

	_, foundB := onboarding.LatestSequenceState(reporter, inputB)
	assert.False(t, foundB)
}
