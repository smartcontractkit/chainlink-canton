package tests

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
)

type OnboardingFlowTestSuite struct {
	CeremonyTestSuite

	PartyID string
}

func (s *OnboardingFlowTestSuite) SetupSuite() {
	s.CeremonyTestSuite.SetupSuite()
}

func (s *OnboardingFlowTestSuite) onboardingNamespaceName() string {
	return s.uniqueName("onboarding-ns")
}

func (s *OnboardingFlowTestSuite) onboardingPartyPrefix() string {
	return s.uniqueName("party")
}

// performOnboarding executes the full 7-step onboarding ceremony using the
// provided reporter (so the caller may reuse it for idempotency checks) and
// sets s.PartyID to the resulting party identifier.
func (s *OnboardingFlowTestSuite) performOnboarding(t *testing.T, reporter operations.Reporter) operations.SequenceReport[onboarding.OnboardingInput, onboarding.OnboardingOutput] {
	t.Helper()

	input := onboarding.OnboardingInput{
		NamespaceName:  s.onboardingNamespaceName(),
		PartyPrefix:    s.onboardingPartyPrefix(),
		Participants:   s.ParticipantIDs,
		SynchronizerID: s.SynchronizerID,
		Threshold:      3, // all must sign
	}

	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), reporter)
	}

	kmsCfg1 := s.kmsConfigFor(0, "onboarding")
	kmsCfg2 := s.kmsConfigFor(1, "onboarding")
	kmsCfg3 := s.kmsConfigFor(2, "onboarding")
	deps1 := s.depsFor(0, kmsCfg1)
	deps2 := s.depsFor(1, kmsCfg2)
	deps3 := s.depsFor(2, kmsCfg3)

	// Run 1: Actor 1 generates key (1/3 keys)
	t.Log("Actor 1: run 1 — key gen (1/3)")
	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps1, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "actor 1 run 1: expected threshold-not-met")

	// Run 2: Actor 2 generates key (2/3 keys)
	t.Log("Actor 2: run 2 — key gen (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps2, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "actor 2 run 2: expected threshold-not-met")

	// Run 3: Actor 3 generates key (3/3), creates DNS proposal, signs (1/3 DNS)
	t.Log("Actor 3: run 3 — key gen (3/3) + DNS sig (1/3)")
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps3, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "actor 3 run 3: expected threshold-not-met (1/3 DNS sigs)")

	// Run 4: Actor 1 signs DNS proposal (2/3 DNS)
	t.Log("Actor 1: run 4 — DNS sig (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps1, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "actor 1 run 4: expected threshold-not-met (2/3 DNS sigs)")

	// Run 5: Actor 2 signs DNS (3/3) → DNS submitted; proposes P2P (1/3)
	t.Log("Actor 2: run 5 — DNS sig (3/3) + submit + P2P proposal (1/3)")
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps2, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "actor 2 run 5: expected threshold-not-met (1/3 P2P)")

	// Run 6: Actor 3 proposes P2P (2/3)
	t.Log("Actor 3: run 6 — P2P proposal (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps3, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "actor 3 run 6: expected threshold-not-met (2/3 P2P)")

	// Run 7: Actor 1 proposes P2P (3/3) → P2P confirmed → ceremony done
	t.Log("Actor 1: run 7 — P2P proposal (3/3) — completing ceremony")
	sr, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps1, input)
	require.NoError(t, err, "actor 1 run 7: ceremony should complete successfully")

	s.PartyID = sr.Output.PartyID
	s.assertKMSKeysRegistered(0, kmsCfg1)
	s.assertKMSKeysRegistered(1, kmsCfg2)
	s.assertKMSKeysRegistered(2, kmsCfg3)

	return sr
}

// ── Test ─────────────────────────────────────────────────────────────────────

// TestOnboardingFlow validates the full 7-step decentralized party onboarding
// ceremony against a real CTF Canton environment. The suite is run once with
// KMS-backed ceremony deps and once with generated local keys.
//
// The ceremony is async: each actor runs the OnboardingSequence independently,
// sharing a MemoryReporter so cached operation results are visible across runs.
// The sequence returns ErrThresholdNotMet until enough actors have contributed:
//
//   - Runs 1–2 (actors 1–2): key-gen threshold not yet met (1/3, 2/3 keys).
//   - Run  3  (actor  3):    all keys present; DNS proposal created; actor 3
//     signs (1/3 DNS sigs) → ErrThresholdNotMet.
//   - Run  4  (actor  1):    actor 1 signs (2/3 DNS sigs) → ErrThresholdNotMet.
//   - Run  5  (actor  2):    actor 2 signs (3/3) → DNS submitted and confirmed;
//     actor 2 proposes P2P (1/3 P2P) → ErrThresholdNotMet.
//   - Run  6  (actor  3):    actor 3 proposes P2P (2/3 P2P) → ErrThresholdNotMet.
//   - Run  7  (actor  1):    actor 1 proposes P2P (3/3) → confirmed → success.
func (s *OnboardingFlowTestSuite) TestOnboardingFlow() {
	t := s.T()

	// Shared reporter: cached operation results are visible to all actors,
	// including the idempotency re-run (run 8) below.
	sharedReporter := operations.NewMemoryReporter()
	sr := s.performOnboarding(t, sharedReporter)

	// ── Verify output fields ─────────────────────────────────────────────
	assert.True(t, strings.HasPrefix(sr.Output.PartyID, s.onboardingPartyPrefix()+"::"),
		"PartyID should start with the suite party prefix, got: %s", sr.Output.PartyID)
	assert.True(t, sr.Output.DNSConfirmed, "DNSConfirmed should be true")
	assert.True(t, sr.Output.P2PConfirmed, "P2PConfirmed should be true")

	// ── Verify topology via admin API ────────────────────────────────────
	parts := strings.SplitN(sr.Output.PartyID, "::", 2)
	require.Len(t, parts, 2, "PartyID %q should be <name>::<namespace>", sr.Output.PartyID)
	decNS := parts[1]

	dnsOK, err := s.Actors[0].client.DNSExists(t.Context(), decNS, s.SynchronizerID)
	require.NoError(t, err, "DNSExists query failed")
	assert.True(t, dnsOK, "DecentralizedNamespaceDefinition should be active in topology")

	p2pOK, err := s.Actors[0].client.P2PExists(t.Context(), sr.Output.PartyID, s.SynchronizerID)
	require.NoError(t, err, "P2PExists query failed")
	assert.True(t, p2pOK, "PartyToParticipant mapping should be active in topology")
	p2pState, err := s.Actors[0].client.GetP2P(t.Context(), sr.Output.PartyID, s.SynchronizerID)
	require.NoError(t, err, "GetP2P after onboarding")
	require.NotNil(t, p2pState.PartySigningKeys, "P2P signing keys should be active")
	assert.Len(t, p2pState.PartySigningKeys.Keys, 3, "should have one protocol signing key per participant")
	assert.Equal(t, uint32(3), p2pState.PartySigningKeys.Threshold, "P2P signing-key threshold should match onboarding threshold")

	// ── Verify idempotency: re-run produces the same output from cache ───
	t.Log("Actor 1: run 8 — idempotency check")
	deps1 := s.depsFor(0, s.kmsConfigFor(0, "onboarding"))
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}
	input := onboarding.OnboardingInput{
		NamespaceName:  s.onboardingNamespaceName(),
		PartyPrefix:    s.onboardingPartyPrefix(),
		Participants:   s.ParticipantIDs,
		SynchronizerID: s.SynchronizerID,
		Threshold:      3,
	}
	srCached, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps1, input)
	require.NoError(t, err, "idempotent re-run should succeed")
	assert.Equal(t, sr.Output.PartyID, srCached.Output.PartyID, "cached PartyID should match")
	assert.True(t, srCached.Output.DNSConfirmed, "cached DNSConfirmed should be true")
	assert.True(t, srCached.Output.P2PConfirmed, "cached P2PConfirmed should be true")
	s.assertReportsDoNotContainKMS(sharedReporter)
}
