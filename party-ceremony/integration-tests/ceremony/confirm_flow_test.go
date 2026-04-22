package tests

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/ceremony/onboarding"
)

// ConfirmFlowTestSuite validates that the Confirmer integration works
// end-to-end against a real Canton environment.
type ConfirmFlowTestSuite struct {
	CeremonyTestSuite
}

func (s *ConfirmFlowTestSuite) SetupSuite() {
	s.CeremonyTestSuite.SetupSuite()
}

// TestOnboardingWithNoOpConfirmer runs the full onboarding ceremony with
// NoOpConfirmer set on all deps. The ceremony should complete identically
// to the base case — the confirmer is in the path but auto-approves.
func (s *ConfirmFlowTestSuite) TestOnboardingWithNoOpConfirmer() {
	t := s.T()

	input := onboarding.OnboardingInput{
		NamespaceName:  "inttest-confirm-noop",
		PartyPrefix:    "confirm-party",
		Participants:   s.ParticipantIDs,
		SynchronizerID: s.SynchronizerID,
		Threshold:      3,
	}

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	// Build deps with NoOpConfirmer for all actors.
	deps := make([]ceremony.CantonDeps, 3)
	for i := range 3 {
		deps[i] = ceremony.CantonDeps{
			Client:    s.Actors[i].client,
			Logger:    logger.Test(t),
			Confirmer: ceremony.NoOpConfirmer{},
		}
	}

	// Run the standard 7-run flow with confirmer in the path.
	for i := range 3 {
		_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[i], input)
		require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	}

	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[0], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())

	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[1], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())

	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[2], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())

	sr, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[0], input)
	require.NoError(t, err, "ceremony should complete with NoOpConfirmer")
	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
}

// TestOnboardingRejectsDNSSigning runs the ceremony through key-gen and
// proposal phases, then uses AlwaysRejectConfirmer to verify that the
// DNS signing step is blocked and no DNS is submitted.
func (s *ConfirmFlowTestSuite) TestOnboardingRejectsDNSSigning() {
	t := s.T()

	input := onboarding.OnboardingInput{
		NamespaceName:  "inttest-confirm-reject",
		PartyPrefix:    "reject-party",
		Participants:   s.ParticipantIDs,
		SynchronizerID: s.SynchronizerID,
		Threshold:      3,
	}

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	// Phase 1: Key gen with no confirmer (nil) — 3 runs, each generating keys.
	for i := range 3 {
		deps := ceremony.CantonDeps{
			Client: s.Actors[i].client,
			Logger: logger.Test(t),
			// No confirmer for key gen phase
		}
		_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps, input)
		require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	}

	// Phase 2: Actor 3 (run 3) created the DNS proposal and would try to sign.
	// Now switch to AlwaysRejectConfirmer for the signing step.
	rejectDeps := ceremony.CantonDeps{
		Client:    s.Actors[0].client,
		Logger:    logger.Test(t),
		Confirmer: ceremony.AlwaysRejectConfirmer{},
	}

	// This run should fail because the confirmer rejects signing.
	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, rejectDeps, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), ceremony.ErrUserRejected.Error(),
		"expected user-rejected error when confirmer rejects signing")

	// Verify DNS was NOT submitted (namespace should not exist yet).
	state, ok := onboarding.LatestSequenceState(sharedReporter, input)
	if ok {
		assert.NotEqual(t, "completed", string(state.Phase),
			"ceremony should NOT have reached completion")
	}
}
