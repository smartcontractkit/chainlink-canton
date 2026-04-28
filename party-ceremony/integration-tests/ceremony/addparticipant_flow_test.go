package tests

import (
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
)

// AddParticipantFlowTestSuite validates adding a participant back to a
// decentralized party that was previously reduced by a kick.
//
// SetupSuite runs:
//  1. Onboard a 3-member party (threshold=3) via inherited KickFlowTestSuite.
//  2. Kick participant 3 → 2-member party (threshold=2) via performKick.
//
// TestAddParticipantFlow then adds participant 3 back.
type AddParticipantFlowTestSuite struct {
	KickFlowTestSuite
}

func (s *AddParticipantFlowTestSuite) SetupSuite() {
	// Runs CeremonyTestSuite.SetupSuite (chain init) + performOnboarding.
	s.KickFlowTestSuite.SetupSuite()
	// Run kick so TestAddParticipantFlow finds a 2-member party.
	s.performKick()
}

// performKick executes the full kick ceremony to remove participant 3 (actors[2]).
// Called in SetupSuite so TestAddParticipantFlow (alphabetically before TestKickFlow)
// finds the expected 2-member post-kick state.
func (s *AddParticipantFlowTestSuite) performKick() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID

	decNS := strings.SplitN(s.PartyID, "::", 2)[1]
	dnsState, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after onboarding")

	kickedNSFP, err := actors[2].deps.Client.GetNamespaceFingerprint(t.Context(), onboardingNamespaceName, synchronizerID, dnsState.Owners)
	require.NoError(t, err, "GetNamespaceFingerprint for kicked participant")

	kickInput := kick.KickInput{
		DecentralizedPartyID:       s.PartyID,
		KickedParticipantID:        actors[2].uid,
		KickedNamespaceFingerprint: kickedNSFP,
		NewThreshold:               2,
		RemainingParticipants:      []string{actors[0].uid, actors[1].uid},
		SynchronizerID:             synchronizerID,
	}

	sharedKickReporter := operations.NewMemoryReporter()
	newKickBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedKickReporter)
	}

	kickDeps := [3]ceremony.CantonDeps{
		s.OnboardingDeps(0),
		s.OnboardingDeps(1),
		s.OnboardingDeps(2),
	}

	// Run 1 (p1): reads state, creates DNS proposal, signs (1/3) → ErrThresholdNotMet.
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[0], kickInput)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())

	// Run 2 (p3 — kicked, still a current DNS owner): signs DNS (2/3) → ErrThresholdNotMet.
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[2], kickInput)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())

	// Run 3 (p2): signs DNS (3/3) → DNS submitted → P2P (1/2) → ErrThresholdNotMet.
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[1], kickInput)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())

	// Run 4 (p1): P2P (2/2) → confirmed → SUCCESS.
	_, err = operations.ExecuteSequence(newKickBundle(), kick.KickSequence, kickDeps[0], kickInput)
	require.NoError(t, err, "kick ceremony should complete successfully")
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
//   - Run 3 (p2): signs DNS (2/2) → submits → P2P from p2 (1/2 existing) → ErrThresholdNotMet.
//   - Run 4 (p1): P2P from p1 (2/2 existing) → new participant consent pending → ErrThresholdNotMet.
//   - Run 5 (p3 — new): consents to P2P hosting → P2P confirmed → SUCCESS.
func (s *AddParticipantFlowTestSuite) TestAddParticipantFlow() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID

	// ── Phase 1: Verify post-kick state (2-member party) ──────────────────
	t.Log("Phase 1: verifying post-kick state")
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]

	dnsState, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after kick")
	require.Len(t, dnsState.Owners, 2, "should have 2 owners after kick")

	p2pState, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after kick")
	require.Len(t, p2pState.Participants, 2, "should have 2 P2P participants after kick")

	newParticipantUID := actors[2].uid

	t.Logf("Adding participant back: %s", newParticipantUID)

	// ── Phase 2: Add participant 3 ────────────────────────────────────────
	t.Log("Phase 2: add-participant ceremony")

	addInput := addparticipant.AddParticipantInput{
		DecentralizedPartyID: s.PartyID,
		NewParticipantID:     newParticipantUID,
		NamespaceName:        "inttest-add-participant",
		SynchronizerID:       synchronizerID,
		NewThreshold:         2,
	}

	sharedAddReporter := operations.NewMemoryReporter()
	newAddBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedAddReporter)
	}

	addDeps := [3]ceremony.CantonDeps{
		s.OnboardingDeps(0),
		s.OnboardingDeps(1),
		s.OnboardingDeps(2),
	}

	// Run 1 (p3 — new): generates keys, NSD, reads state → 0/2 DNS sigs → ErrThresholdNotMet.
	t.Log("Add run 1: p3 generates keys + NSD + reads state (0/2 DNS sigs)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[2], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 1: DNS signing (0/2)")

	// Run 2 (p1): signs DNS (1/2) → ErrThresholdNotMet.
	t.Log("Add run 2: p1 signs DNS (1/2)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[0], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 2: DNS signing (1/2)")

	// Run 3 (p2): signs DNS (2/2) → submits → P2P from p2 (1/2 existing) → ErrThresholdNotMet.
	t.Log("Add run 3: p2 signs DNS (2/2) + P2P proposal (1/2 existing)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[1], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 3: P2P proposals (1/2 existing)")

	// Run 4 (p1): P2P from p1 (2/2 existing) → new participant consent pending → ErrThresholdNotMet.
	t.Log("Add run 4: p1 proposes P2P (2/2 existing) — new participant consent pending")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[0], addInput)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "add run 4: new participant consent pending")

	// Run 5 (p3 — new): consents to P2P hosting → confirmed → SUCCESS.
	t.Log("Add run 5: p3 consents to P2P hosting — completing add-participant ceremony")
	addResult, err := operations.ExecuteSequence(newAddBundle(), addparticipant.AddParticipantSequence, addDeps[2], addInput)
	require.NoError(t, err, "add-participant ceremony should complete successfully")

	// ── Verify output ─────────────────────────────────────────────────────
	assert.True(t, addResult.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, addResult.Output.P2PUpdated, "P2PUpdated should be true")
	assert.Len(t, addResult.Output.AllOwners, 3, "should have 3 owners after add")

	// ── Verify via topology read API ──────────────────────────────────────
	updatedDNS, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after add")
	assert.Len(t, updatedDNS.Owners, 3, "should have 3 owners after add")

	updatedP2P, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
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
