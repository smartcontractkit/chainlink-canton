package tests

import (
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/ceremony/addparticipant"
)

// AddParticipantFlowTestSuite validates adding a participant back to a
// decentralized party that was previously reduced by a kick.
//
// It embeds KickFlowTestSuite, which embeds OnboardingFlowTestSuite. The
// embedded SetupSuite runs:
//  1. Onboard a 3-member party (threshold=3).
//  2. Kick participant 3 → 2-member party (threshold=2).
//
// TestAddParticipantFlow then adds participant 3 back.
type AddParticipantFlowTestSuite struct {
	KickFlowTestSuite
}

func (s *AddParticipantFlowTestSuite) SetupSuite() {
	s.KickFlowTestSuite.SetupSuite()
}

// TestOnboardingFlow overrides the inherited method — onboarding was already
// done in SetupSuite via the inherited KickFlowTestSuite.
func (s *AddParticipantFlowTestSuite) TestOnboardingFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Onboarding already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// TestKickFlow overrides the inherited method — kick was already done in
// SetupSuite via the inherited KickFlowTestSuite.
func (s *AddParticipantFlowTestSuite) TestKickFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Kick already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// ── Test ─────────────────────────────────────────────────────────────────────

// TestAddParticipantFlow validates the full end-to-end add-participant ceremony:
//
//  1. Onboard 3-member party (done in SetupSuite).
//  2. Kick participant 3 → 2-member party (done in SetupSuite).
//  3. Add participant 3 back → 3-member party.
//
// Add-participant ceremony progression:
//
//   - Run 1 (p3 — new): generates keys, proposes NSD, reads state, creates DNS
//     proposal → 0/2 DNS sigs → ErrThresholdNotMet.
//   - Run 2 (p1): signs DNS (1/2) → ErrThresholdNotMet.
//   - Run 3 (p2): signs DNS (2/2) → submits → P2P (1/2) → ErrThresholdNotMet.
//   - Run 4 (p1): P2P (2/2) → P2P confirmed → SUCCESS.
func (s *AddParticipantFlowTestSuite) TestAddParticipantFlow() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID

	// ── Phase 1: Verify post-kick state (2-member party) ──────────────────
	t.Log("Phase 1: verifying post-kick state")
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]

	dnsState, err := actors[0].client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after kick")
	require.Len(t, dnsState.Owners, 2, "should have 2 owners after kick")

	p2pState, err := actors[0].client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after kick")
	require.Len(t, p2pState.Participants, 2, "should have 2 P2P participants after kick")

	// Remaining participant UIDs (p1 and p2).
	existingUIDs := []string{actors[0].uid, actors[1].uid}
	newParticipantUID := actors[2].uid

	t.Logf("Adding participant back: %s", newParticipantUID)
	t.Logf("Existing participants: %v", existingUIDs)

	// ── Phase 2: Add participant 3 ────────────────────────────────────────
	t.Log("Phase 2: add-participant ceremony")

	addInput := addparticipant.AddParticipantInput{
		DecentralizedPartyID: s.PartyID,
		NewParticipantID:     newParticipantUID,
		ExistingParticipants: existingUIDs,
		NamespaceName:        "inttest-add-participant",
		SynchronizerID:       synchronizerID,
		NewThreshold:         2,
	}

	sharedAddReporter := operations.NewMemoryReporter()
	newAddBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedAddReporter)
	}

	addDeps := [3]ceremony.CantonDeps{
		{Client: actors[0].client, Logger: logger.Test(t)},
		{Client: actors[1].client, Logger: logger.Test(t)},
		{Client: actors[2].client, Logger: logger.Test(t)},
	}

	// Run 1 (p3 — new): generates keys, NSD, reads state → 0/2 DNS sigs → ErrThresholdNotMet.
	t.Log("Add run 1: p3 generates keys + NSD + reads state (0/2 DNS sigs)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[2], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 1: DNS signing (0/2)")

	// Run 2 (p1): signs DNS (1/2) → ErrThresholdNotMet.
	t.Log("Add run 2: p1 signs DNS (1/2)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[0], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 2: DNS signing (1/2)")

	// Run 3 (p2): signs DNS (2/2) → submits → P2P (1/2) → ErrThresholdNotMet.
	t.Log("Add run 3: p2 signs DNS (2/2) + P2P proposal (1/2)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[1], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 3: P2P proposals (1/2)")

	// Run 4 (p1): P2P (2/2) → confirmed → SUCCESS.
	t.Log("Add run 4: p1 proposes P2P (2/2) — completing add-participant ceremony")
	addResult, err := operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[0], addInput)
	require.NoError(t, err, "add-participant ceremony should complete successfully")

	// ── Verify output ─────────────────────────────────────────────────────
	assert.True(t, addResult.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, addResult.Output.P2PUpdated, "P2PUpdated should be true")
	assert.Len(t, addResult.Output.AllOwners, 3, "should have 3 owners after add")

	// ── Verify via topology read API ──────────────────────────────────────
	updatedDNS, err := actors[0].client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after add")
	assert.Len(t, updatedDNS.Owners, 3, "should have 3 owners after add")

	updatedP2P, err := actors[0].client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after add")
	assert.Len(t, updatedP2P.Participants, 3, "should have 3 P2P participants after add")

	// Verify new participant is in P2P.
	var foundNew bool
	for _, p := range updatedP2P.Participants {
		if p.ParticipantUID == newParticipantUID {
			foundNew = true
			break
		}
	}
	assert.True(t, foundNew, "new participant should be present in P2P mapping")

	// ── Verify idempotency ────────────────────────────────────────────────
	t.Log("Add-participant idempotency check")
	addResultCached, err := operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[0], addInput)
	require.NoError(t, err, "idempotent re-run should succeed")
	assert.Equal(t, addResult.Output, addResultCached.Output, "cached result should match")
}
