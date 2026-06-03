package tests

import (
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipantwithacs"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// AddParticipantWithAcsFlowTestSuite validates adding a participant back to a
// decentralized party with ACS replication. The new participant receives
// pre-existing contracts via an offline snapshot rather than built-in party
// replication, which is gated by the onboarding flag on the P2P mapping.
//
// SetupSuite runs:
//  1. Onboard a 3-member party (threshold=3) via inherited KickFlowTestSuite.
//  2. Deploy a contract to all 3 participants — stores ContractID + PackageIDs.
//  3. Kick participant 3 → 2-member party (threshold=2).
//
// TestAddParticipantWithAcsFlow then adds participant 3 back with ACS import.
type AddParticipantWithAcsFlowTestSuite struct {
	AddParticipantFlowTestSuite

	SynchronizerAlias string
	ContractID        string
	PackageIDs        []string
}

func (s *AddParticipantWithAcsFlowTestSuite) SetupSuite() {
	t := s.T()

	// Step 1: Run onboarding (sets s.PartyID via KickFlowTestSuite.SetupSuite).
	// We deliberately skip AddParticipantFlowTestSuite.SetupSuite because it
	// calls performKick before we have a chance to deploy the contract.
	s.KickFlowTestSuite.SetupSuite()

	// Step 2: Discover the synchronizer alias. This is needed by ImportAcsOp
	// for the disconnect/reconnect of the target participant.
	syncInfos, err := s.Actors[0].deps.Client.ListConnectedSynchronizers(t.Context())
	require.NoError(t, err, "ListConnectedSynchronizers to discover synchronizer alias")
	for _, info := range syncInfos {
		if info.SynchronizerID == s.SynchronizerID {
			s.SynchronizerAlias = info.Alias
			break
		}
	}
	require.NotEmpty(t, s.SynchronizerAlias, "could not resolve alias for synchronizer %s", s.SynchronizerID)
	t.Logf("Resolved synchronizer alias: %s", s.SynchronizerAlias)

	// Step 3: Deploy a contract to the 3-member party BEFORE kicking p3.
	// After the kick, p3 will not be able to see this contract; after the
	// AddParticipantWithAcs ceremony it should become visible via ACS import.
	t.Log("SetupSuite: deploying contract to 3-member party")
	kmsCfgs := []client.KMSConfig{
		s.kmsConfigFor(0, "onboarding"),
		s.kmsConfigFor(1, "onboarding"),
		s.kmsConfigFor(2, "onboarding"),
	}
	sr, recorders := s.runContractDeployFlow(t, s.PartyID, "add-with-acs-pre-deploy", kmsCfgs)
	assertRecordersMatchKMSConfig(t, kmsCfgs, recorders)
	s.ContractID = sr.Output.ContractID
	s.PackageIDs = sr.Output.PackageIDs
	require.NotEmpty(t, s.ContractID, "pre-deploy contract ID must be set")
	require.NotEmpty(t, s.PackageIDs, "pre-deploy package IDs must be set")
	t.Logf("Pre-deploy contract ID: %s", s.ContractID)

	// Step 4: Kick participant 3 → 2-member party.
	t.Log("SetupSuite: kicking participant 3")
	s.performKick()
}

// TestOnboardingFlow overrides the inherited method — onboarding was already
// done in SetupSuite.
func (s *AddParticipantWithAcsFlowTestSuite) TestOnboardingFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Onboarding already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// TestKickFlow overrides the inherited method — kick was already done in
// SetupSuite.
func (s *AddParticipantWithAcsFlowTestSuite) TestKickFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Kick already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// TestAddParticipantFlow overrides the inherited method — this suite adds a
// participant with ACS replication instead of topology-only.
func (s *AddParticipantWithAcsFlowTestSuite) TestAddParticipantFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Add-participant-with-ACS suite skips topology-only add; PartyID=%s", s.PartyID)
}

// ── Test ─────────────────────────────────────────────────────────────────────

// TestAddParticipantWithAcsFlow validates the full end-to-end combined
// add-participant + ACS replication ceremony:
//
//  1. Onboard 3-member party + deploy contract + kick p3 (done in SetupSuite).
//  2. Assert p3 cannot see the pre-deployed contract (not hosting the party).
//  3. Run the AddParticipantWithAcs ceremony in 7 resume rounds.
//  4. Assert p3 CAN see the contract after ACS import.
//  5. Deploy a second contract — all 3 participants see it via normal replication.
//  6. Idempotent re-run returns cached output.
//
// Ceremony progression (p2 is source, p3 is new participant, threshold=2):
//
//   - Run 1 (p3/new):    key gen + NSD + reads state → DNS proposal → 0/2 sigs → ErrThresholdNotMet.
//   - Run 2 (p1):        signs DNS (1/2) → ErrThresholdNotMet.
//   - Run 3 (p2/source): signs DNS (2/2) → submits DNS → records ledger offset →
//     P2P w/flag (1/2 existing) → ErrThresholdNotMet.
//   - Run 4 (p1):        P2P w/flag (2/2 existing) → new participant consent pending → ErrThresholdNotMet.
//   - Run 5 (p3/new):    P2P consent w/flag → P2P confirmed → ACS export pending → ErrThresholdNotMet.
//   - Run 6 (p2/source): ACS export → snapshot stored in reporter → ACS import pending → ErrThresholdNotMet.
//   - Run 7 (p3/target): disconnect → ACS import → reconnect → clear onboarding flag → SUCCESS.
func (s *AddParticipantWithAcsFlowTestSuite) TestAddParticipantWithAcsFlow() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]

	// ── Phase 1: Verify post-kick state (2-member party) ──────────────────
	t.Log("Phase 1: verifying post-kick state")

	dnsState, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after kick")
	require.Len(t, dnsState.Owners, 2, "should have 2 owners after kick")

	p2pState, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after kick")
	require.Len(t, p2pState.Participants, 2, "should have 2 P2P participants after kick")

	// ── Phase 2: Verify p3 cannot see the pre-deployed contract ──────────
	// After the kick p3's participant no longer hosts the party, so its ACS is
	// empty for contracts owned by that party. Canton may also return
	// PermissionDenied when querying a non-hosted party — both outcomes
	// confirm p3 has no visibility.
	t.Log("Phase 2: verifying p3 cannot see pre-deployed contract")

	p3LedgerClient, p3LedgerConn := s.NewLedgerClient(s.chain.Participants[2])
	t.Cleanup(func() { _ = p3LedgerConn.Close() })

	preImportContracts, err := p3LedgerClient.GetActiveContractsByTemplate(
		t.Context(), s.PartyID, s.PackageIDs[0], "Main", "DisclosedTarget",
	)
	if err != nil {
		// PermissionDenied (or similar) is expected when p3 no longer hosts
		// the party — the error itself confirms p3 cannot see the contract.
		t.Logf("GetActiveContractsByTemplate on p3 returned expected error (not hosting): %v", err)
	} else {
		assert.Empty(t, preImportContracts, "p3 should not see the contract before ACS import (not hosting)")
	}

	// ── Phase 3: Run AddParticipantWithAcs ceremony ───────────────────────
	t.Log("Phase 3: add-participant-with-acs ceremony")

	sourceUID := actors[1].uid         // p2 is the ACS source
	newParticipantUID := actors[2].uid // p3 is the new participant (target)

	addInput := addparticipantwithacs.AddParticipantWithAcsInput{
		DecentralizedPartyID: s.PartyID,
		NewParticipantID:     newParticipantUID,
		NamespaceName:        s.uniqueName("add-participant-with-acs-ns"),
		SynchronizerID:       synchronizerID,
		SynchronizerAlias:    s.SynchronizerAlias,
		SourceParticipantID:  sourceUID,
		NewThreshold:         2,
	}

	acsKMS := s.kmsConfigFor(2, "add-participant-with-acs")
	addDeps := [3]ceremony.CantonDeps{
		s.OnboardingDeps(0),
		s.OnboardingDeps(1),
		s.depsFor(2, acsKMS),
	}

	sharedAddReporter := operations.NewMemoryReporter()
	newAddBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedAddReporter)
	}

	// Run 1 (p3/new): generates keys, NSD, reads state → DNS proposal → 0/2 sigs → ErrThresholdNotMet.
	t.Log("ACS run 1: p3 generates keys + NSD + reads state (0/2 DNS sigs)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[2], addInput)
	require.ErrorContains(t, err, addparticipantwithacs.ErrThresholdNotMet.Error(), "run 1: DNS signing (0/2)")

	// Run 2 (p1): signs DNS (1/2) → ErrThresholdNotMet.
	t.Log("ACS run 2: p1 signs DNS (1/2)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[0], addInput)
	require.ErrorContains(t, err, addparticipantwithacs.ErrThresholdNotMet.Error(), "run 2: DNS signing (1/2)")

	// Run 3 (p2/source): signs DNS (2/2) → submits DNS → records ledger offset →
	// P2P w/flag from p2 (1/2 existing) → ErrThresholdNotMet.
	t.Log("ACS run 3: p2 signs DNS (2/2) + submits + records offset + P2P (1/2 existing)")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[1], addInput)
	require.ErrorContains(t, err, addparticipantwithacs.ErrThresholdNotMet.Error(), "run 3: P2P proposals (1/2 existing)")

	// Run 4 (p1): P2P w/flag from p1 (2/2 existing) → new participant consent pending → ErrThresholdNotMet.
	t.Log("ACS run 4: p1 proposes P2P w/flag (2/2 existing) — new participant consent pending")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[0], addInput)
	require.ErrorContains(t, err, addparticipantwithacs.ErrThresholdNotMet.Error(), "run 4: new participant consent pending")

	// Run 5 (p3/new): consents to P2P hosting → P2P confirmed → ACS export pending → ErrThresholdNotMet.
	t.Log("ACS run 5: p3 consents to P2P w/flag → P2P confirmed → ACS export pending")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[2], addInput)
	require.ErrorContains(t, err, addparticipantwithacs.ErrThresholdNotMet.Error(), "run 5: ACS export pending")

	// Run 6 (p2/source): exports ACS snapshot → stored in reporter → ACS import pending → ErrThresholdNotMet.
	t.Log("ACS run 6: p2 exports ACS snapshot")
	_, err = operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[1], addInput)
	require.ErrorContains(t, err, addparticipantwithacs.ErrThresholdNotMet.Error(), "run 6: ACS import pending")

	// Run 7 (p3/target): disconnect → import ACS → reconnect → clear onboarding flag → SUCCESS.
	t.Log("ACS run 7: p3 imports ACS + clears onboarding flag — completing ceremony")
	addResult, err := operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[2], addInput)
	require.NoError(t, err, "add-participant-with-acs ceremony should complete successfully")

	// ── Verify output fields ──────────────────────────────────────────────
	assert.True(t, addResult.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, addResult.Output.P2PUpdated, "P2PUpdated should be true")
	assert.True(t, addResult.Output.AcsImported, "AcsImported should be true")
	assert.Len(t, addResult.Output.AllOwners, 3, "should have 3 owners after add")
	assert.Equal(t, 2, addResult.Output.NewThreshold, "new threshold should be 2")
	assert.True(t, addResult.Output.State.OnboardingFlagCleared, "onboarding flag should be cleared")

	// ── Phase 4: Verify topology via admin API ────────────────────────────
	t.Log("Phase 4: verifying topology after ceremony")

	updatedDNS, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after add-participant-with-acs")
	assert.Len(t, updatedDNS.Owners, 3, "should have 3 owners after ceremony")

	updatedP2P, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after add-participant-with-acs")
	assert.Len(t, updatedP2P.Participants, 3, "should have 3 P2P participants after ceremony")
	require.NotNil(t, updatedP2P.PartySigningKeys, "P2P signing keys should be present after add")
	assert.Len(t, updatedP2P.PartySigningKeys.Keys, 3, "should have 3 protocol signing keys after add")

	var foundNew bool
	for _, p := range updatedP2P.Participants {
		if p.ParticipantUID == newParticipantUID {
			foundNew = true
			break
		}
	}
	assert.True(t, foundNew, "new participant should be present in P2P mapping after ceremony")

	// ── Phase 5: Verify p3 CAN see pre-deployed contract via ACS import ───
	// The ACS snapshot exported during the ceremony contained the pre-deployed
	// contract, so p3's ledger should now have it visible.
	t.Log("Phase 5: verifying p3 can see pre-deployed contract after ACS import")

	// Grant party ledger rights to p3's user now that the party is hosted again.
	if s.chain.Participants[2].UserID != "" {
		err = p3LedgerClient.GrantPartyRights(t.Context(), s.chain.Participants[2].UserID, s.PartyID)
		require.NoError(t, err, "GrantPartyRights to p3 user after ceremony")
	}

	// p3's Ledger API may briefly return PermissionDenied (or empty results)
	// immediately after onboarding-flag clearance while the participant
	// reconciles the imported ACS with the just-activated party. Poll until
	// the pre-deployed contract becomes visible, up to 30s.
	var postImportContracts []*apiv2.CreatedEvent
	require.Eventually(t, func() bool {
		var qErr error
		// Use the by-party filter: p3's user has CanReadAs for s.PartyID via
		// the GrantPartyRights call above, but does NOT have the
		// participant/IDP admin claims that FiltersForAnyParty requires.
		postImportContracts, qErr = p3LedgerClient.GetActiveContractsByTemplateForParty(
			t.Context(), s.PartyID, s.PackageIDs[0], "Main", "DisclosedTarget",
		)
		if qErr != nil {
			t.Logf("GetActiveContractsByTemplate on p3 not yet ready: %v", qErr)
			return false
		}

		return len(postImportContracts) >= 1
	}, 30*time.Second, 1*time.Second, "p3 should see the pre-deployed contract after ACS import")
	require.Len(t, postImportContracts, 1, "p3 should see exactly the pre-deployed contract after ACS import")
	assert.Equal(t, s.ContractID, postImportContracts[0].GetContractId(),
		"p3's ACS should contain the contract created before the kick")

	// ── Phase 6: Deploy a second contract — all 3 see it naturally ────────
	// Unlike the first contract (replicated via ACS import), the second contract
	// is created after p3 rejoins the party and should appear on all participants
	// via Canton's normal party replication.
	t.Log("Phase 6: deploying second contract to verify normal replication to p3")

	kmsCfgsForSecondDeploy := []client.KMSConfig{
		s.kmsConfigFor(0, "onboarding"),
		s.kmsConfigFor(1, "onboarding"),
		acsKMS, // p3 uses the ACS ceremony KMS config
	}
	sr2, recorders2 := s.runContractDeployFlow(t, s.PartyID, "add-with-acs-post-deploy", kmsCfgsForSecondDeploy)
	assertRecordersMatchKMSConfig(t, kmsCfgsForSecondDeploy, recorders2)
	require.NotEmpty(t, sr2.Output.ContractID, "second contract deploy should produce a contract ID")
	t.Logf("Second contract ID: %s", sr2.Output.ContractID)

	// All three participants should now see two contracts: original + second.
	for i, p := range s.chain.Participants {
		lc, conn := s.NewLedgerClient(p)
		t.Cleanup(func() { _ = conn.Close() })

		if p.UserID != "" {
			err = lc.GrantPartyRights(t.Context(), p.UserID, s.PartyID)
			require.NoError(t, err, "GrantPartyRights to participant %d for post-deploy check", i+1)
		}

		// Use FiltersByParty (CanReadAs is enough); FiltersForAnyParty would
		// fail on participants without super-reader-wildcard admin claims.
		allContracts, qErr := lc.GetActiveContractsByTemplateForParty(
			t.Context(), s.PartyID, sr2.Output.PackageIDs[0], "Main", "DisclosedTarget",
		)
		require.NoError(t, qErr, "GetActiveContractsByTemplateForParty on participant %d after second deploy", i+1)
		assert.Len(t, allContracts, 2, "participant %d should see 2 contracts after second deploy", i+1)
	}

	// ── Phase 7: Idempotency ──────────────────────────────────────────────
	t.Log("Phase 7: idempotency check")
	addResultCached, err := operations.ExecuteSequence(newAddBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[0], addInput)
	require.NoError(t, err, "idempotent re-run should succeed")
	assert.Equal(t, addResult.Output, addResultCached.Output, "cached result should match original")
	s.assertKMSKeysRegistered(2, acsKMS)
	s.assertReportsDoNotContainKMS(sharedAddReporter)
}
