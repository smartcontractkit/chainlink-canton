package keyrotation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/chainlink/canton-party-ceremony/ceremony/keyrotation"
)

func bundleWith(t *testing.T, r operations.Reporter) operations.Bundle {
	t.Helper()
	return operations.NewBundle(t.Context, logger.Nop(), r)
}

// TestState_ThresholdNotMet_DNS verifies state at the DNS signing gate
// during namespace key rotation.
func TestState_ThresholdNotMet_DNS(t *testing.T) {
	t.Parallel()

	state := &mockState{}
	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseNamespaceInput()

	// p1 (target): generates key, NSD, DNS proposal, signs (1/2) → threshold not met.
	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())

	s := sr.Output.State
	assert.Equal(t, keyrotation.PhaseDNSSigning, s.Phase)
	assert.True(t, s.TargetKeyGenReady)
	assert.True(t, s.NSDProposed)
	assert.True(t, s.RotateNamespace)
	assert.False(t, s.RotateDaml)
	assert.NotEmpty(t, s.ProposalHash)
	assert.Contains(t, s.CollectedSigners, testTargetUID)
	assert.NotEmpty(t, s.PendingSigners)
	assert.False(t, sr.Output.DNSUpdated)
}

// TestState_ThresholdNotMet_P2P verifies state at the P2P gate
// during DAML key rotation.
func TestState_ThresholdNotMet_P2P(t *testing.T) {
	t.Parallel()

	state := &mockState{}
	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseDamlInput()

	// p1 (target): generates key, proposes P2P (1/2) → threshold not met.
	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())

	s := sr.Output.State
	assert.Equal(t, keyrotation.PhaseP2P, s.Phase)
	assert.True(t, s.TargetKeyGenReady)
	assert.False(t, s.RotateNamespace)
	assert.True(t, s.RotateDaml)
	assert.Equal(t, 1, s.P2PProposedCount)
	assert.Equal(t, 2, s.P2PRequired) // threshold=2 from mock
	assert.False(t, sr.Output.P2PUpdated)
}

// TestState_Completed verifies the happy-path state snapshot for namespace-only.
func TestState_Completed(t *testing.T) {
	t.Parallel()

	state := &mockState{}
	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, sharedReporter) }
	input := baseNamespaceInput()

	// Run 1 (p1 target): key gen, NSD, DNS proposal, signs (1/2) → threshold not met.
	_, _ = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)

	// Run 2 (p2): signs DNS (2/2) → submit → confirmed → SUCCESS.
	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState("p2", state), input)
	require.NoError(t, err)

	s := sr.Output.State
	assert.Equal(t, keyrotation.PhaseCompleted, s.Phase)
	assert.Empty(t, s.PendingSigners)
	assert.True(t, s.TargetKeyGenReady)
	assert.True(t, s.NSDProposed)
	assert.True(t, sr.Output.DNSUpdated)
	assert.False(t, sr.Output.P2PUpdated)
}

// TestState_LatestSequenceState verifies observer reads via LatestSequenceState.
func TestState_LatestSequenceState(t *testing.T) {
	t.Parallel()

	state := &mockState{}
	reporter := operations.NewMemoryReporter()
	input := baseNamespaceInput()

	// No runs yet.
	_, found := keyrotation.LatestSequenceState(reporter, input)
	assert.False(t, found)

	// Run once (partial).
	_, err := operations.ExecuteSequence(bundleWith(t, reporter), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())

	s, found := keyrotation.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.Equal(t, keyrotation.PhaseDNSSigning, s.Phase)
}

// TestState_Scoping verifies two key rotation ceremonies don't bleed.
func TestState_Scoping(t *testing.T) {
	t.Parallel()

	state := &mockState{}
	reporter := operations.NewMemoryReporter()
	inputA := baseNamespaceInput()
	inputB := keyrotation.KeyRotationInput{
		DecentralizedPartyID:       "other-party::1220112233",
		TargetParticipantID:        "p99",
		TargetNamespaceFingerprint: "fp-other",
		SynchronizerID:             testSyncID,
		RotateNamespaceKey:         true,
	}

	_, err := operations.ExecuteSequence(bundleWith(t, reporter), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), inputA)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())

	_, foundA := keyrotation.LatestSequenceState(reporter, inputA)
	assert.True(t, foundA)

	_, foundB := keyrotation.LatestSequenceState(reporter, inputB)
	assert.False(t, foundB)
}
