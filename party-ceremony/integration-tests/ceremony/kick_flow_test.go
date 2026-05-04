package tests

import (
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
)

type KickFlowTestSuite struct {
	OnboardingFlowTestSuite
}

func (s *KickFlowTestSuite) SetupSuite() {
	// SetupSuite runs the full onboarding ceremony so s.PartyID is populated
	// before any test methods run (avoiding a race with TestKickFlow).
	s.OnboardingFlowTestSuite.SetupSuite()
	s.performOnboarding(s.T(), operations.NewMemoryReporter())
}

// TestOnboardingFlow overrides the inherited method to be a lightweight
// assertion. The full onboarding was already run inside SetupSuite; re-running
// with a fresh reporter here would generate new keys (different fingerprints)
// and create a separate party, which would overwrite s.PartyID unexpectedly.
func (s *KickFlowTestSuite) TestOnboardingFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Onboarding already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// ── Test ─────────────────────────────────────────────────────────────────────

// TestKickFlow validates the full end-to-end kick ceremony:
//
//  1. Run the onboarding ceremony to create a 3-member decentralized party.
//  2. Run the kick ceremony to remove participant 3 (p3).
//
// Kick ceremony progression (actors p1 and p2 only for P2P, all 3 for DNS signing,
// currentDNSThreshold=3 since serial=1 requires all owners, newThreshold=2 after kick):
//
//   - Run 1 (p1): reads state → creates DNS proposal → signs (1/3) → ErrThresholdNotMet.
//   - Run 2 (p3 — kicked, still a current DNS owner): signs DNS (2/3) → ErrThresholdNotMet.
//   - Run 3 (p2): signs DNS (3/3) → DNS submitted → confirms → proposes P2P (1/2) → ErrThresholdNotMet.
//
// NOTE: For the DNS signing gate the threshold is the CURRENT DNS threshold
// (3-of-3 for a freshly onboarded party with all-must-sign). The kicked
// participant only participates in DNS signing (Canton requires it because
// threshold=3), but NOT in P2P proposals.
//
//   - Run 4 (p1): P2P proposal (2/2) → P2P confirmed → SUCCESS.
func (s *KickFlowTestSuite) TestKickFlow() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID

	// ── Phase 1: Onboard a 3-member party (threshold=3, all must sign) ────────
	t.Log("Phase 1: onboarding a 3-member decentralized party")

	// Determine the kicked participant's namespace fingerprint. We need the DNS
	// state to identify which owner fingerprint belongs to participant 3.
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]
	dnsState, err := actors[0].client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after onboarding")
	require.Len(t, dnsState.Owners, 3, "should have 3 owners after onboarding")

	p2pState, err := actors[0].client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after onboarding")
	require.Len(t, p2pState.Participants, 3, "should have 3 P2P participants after onboarding")
	require.NotNil(t, p2pState.PartySigningKeys, "P2P signing keys should be present after onboarding")
	_, kickedProtocolKeyB64, err := actors[2].client.GetProtocolKeyFingerprint(t.Context(), p2pState.PartySigningKeys.Keys)
	require.NoError(t, err, "GetProtocolKeyFingerprint for kicked participant")

	// Identify the kicked participant's namespace fingerprint. The fingerprint
	// corresponding to actor[2] is discovered by asking their node for its own
	// namespace fingerprint via GetNamespaceFingerprint.
	kickedNSFP, err := actors[2].client.GetNamespaceFingerprint(t.Context(), s.onboardingNamespaceName(), synchronizerID, dnsState.Owners)
	require.NoError(t, err, "GetNamespaceFingerprint for kicked participant")

	kickedUID := actors[2].uid
	remainingUIDs := []string{actors[0].uid, actors[1].uid}

	t.Logf("Kicking participant: %s (ns fingerprint: %s)", kickedUID, kickedNSFP)

	// ── Phase 2: Kick participant 3 ───────────────────────────────────────────
	t.Log("Phase 2: kick ceremony")

	kickInput := kick.KickInput{
		DecentralizedPartyID:       s.PartyID,
		KickedParticipantID:        kickedUID,
		KickedNamespaceFingerprint: kickedNSFP,
		// NOTE: currentDNSThreshold=3 for the freshly onboarded party, so all
		// 3 actors must sign the DNS update (including the kicked participant).
		// The remaining participants list only drives P2P proposals.
		NewThreshold:          2,
		RemainingParticipants: remainingUIDs,
		SynchronizerID:        synchronizerID,
	}

	sharedKickReporter := operations.NewMemoryReporter()
	newKickBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedKickReporter)
	}

	kickDeps := [3]ceremony.CantonDeps{
		s.kmsDepsFor(0),
		s.kmsDepsFor(1),
		s.kmsDepsFor(2),
	}

	// Run 1 (p1): reads state, creates DNS proposal, signs (1/3) → ErrThresholdNotMet.
	t.Log("Kick run 1: p1 signs DNS (1/3)")
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[0], kickInput)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "kick run 1: DNS signing (1/3)")

	// Run 2 (p3 — kicked, but still a current DNS owner): signs DNS (2/3) → ErrThresholdNotMet.
	t.Log("Kick run 2: p3 signs DNS (2/3)")
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[2], kickInput)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "kick run 2: DNS signing (2/3)")

	// Run 3 (p2): signs DNS (3/3) → DNS submitted → confirms → proposes P2P (1/2) → ErrThresholdNotMet.
	t.Log("Kick run 3: p2 signs DNS (3/3) + P2P proposal (1/2)")
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[1], kickInput)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "kick run 3: P2P proposals (1/2)")

	// Run 4 (p1): P2P proposal (2/2) → P2P confirmed → SUCCESS.
	t.Log("Kick run 4: p1 proposes P2P (2/2) — completing kick ceremony")
	kickResult, err := operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[0], kickInput)
	require.NoError(t, err, "kick ceremony should complete successfully")

	// ── Verify output ─────────────────────────────────────────────────────────
	assert.True(t, kickResult.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, kickResult.Output.P2PUpdated, "P2PUpdated should be true")
	assert.Len(t, kickResult.Output.RemainingOwners, 2, "should have 2 remaining owners")
	assert.NotContains(t, kickResult.Output.RemainingOwners, kickedNSFP, "kicked FP must be removed")

	// ── Verify via topology read API ──────────────────────────────────────────
	updatedDNS, err := actors[0].client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after kick")
	assert.Len(t, updatedDNS.Owners, 2, "should have 2 owners after kick")
	assert.NotContains(t, updatedDNS.Owners, kickedNSFP, "kicked fingerprint must be gone from DNS")

	updatedP2P, err := actors[0].client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after kick")
	assert.Len(t, updatedP2P.Participants, 2, "should have 2 P2P participants after kick")
	require.NotNil(t, updatedP2P.PartySigningKeys, "P2P signing keys should remain active after kick")
	assert.Len(t, updatedP2P.PartySigningKeys.Keys, 2, "kicked participant protocol signing key should be removed")
	assert.NotContains(t, updatedP2P.PartySigningKeys.Keys, kickedProtocolKeyB64, "kicked participant protocol signing key must be removed")
	assert.Equal(t, uint32(2), updatedP2P.PartySigningKeys.Threshold, "P2P signing-key threshold should be reduced")
	for _, p := range updatedP2P.Participants {
		assert.NotEqual(t, kickedUID, p.ParticipantUID, "kicked participant must be gone from P2P")
	}

	// ── Verify idempotency ────────────────────────────────────────────────────
	t.Log("Kick idempotency check")
	kickResultCached, err := operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[0], kickInput)
	require.NoError(t, err, "idempotent re-run should succeed")
	assert.Equal(t, kickResult.Output, kickResultCached.Output, "cached result should match")
	s.assertReportsDoNotContainKMS(sharedKickReporter)
}
