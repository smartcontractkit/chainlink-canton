package tests

import (
	"os"
	"sort"
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

// AddParticipantWithAcsResyncFlowTestSuite covers two correctness properties
// of the add-participant-with-acs ceremony that the happy-path suite does not
// exercise:
//
//  1. **Resync on ACS divergence.** A contract present in the source's exported
//     snapshot is archived (by p1+p2) AFTER the snapshot is captured but
//     BEFORE the target reconnects. The synchronizer must deliver the archive
//     event to the target during reconnect; ClearPartyOnboardingFlag must
//     block until that catch-up is past the decision deadline.
//
//  2. **New participant actually confirms.** After the ceremony, the source
//     (p2) is disconnected from the synchronizer and a transaction is
//     submitted with signatures from p1 and p3 only. With threshold=2 and p2
//     down, the submission only commits if p3's participant confirms.
//
// SetupSuite mirrors AddParticipantWithAcsFlowTestSuite but deploys TWO
// contracts pre-kick so the ceremony has a contract to archive mid-flight.
type AddParticipantWithAcsResyncFlowTestSuite struct {
	AddParticipantWithAcsFlowTestSuite

	ContractIDB string // second pre-deployed contract (the one we archive mid-ceremony)

	UpdatedContractIDB string // replacement for ContractIDB created during the resync window
}

func (s *AddParticipantWithAcsResyncFlowTestSuite) SetupSuite() {
	t := s.T()

	// Step 1: Inherited onboarding (3-member party).
	s.KickFlowTestSuite.SetupSuite()

	// Step 2: Discover synchronizer alias (needed by ImportAcsOp).
	syncInfos, err := s.Actors[0].deps.Client.ListConnectedSynchronizers(t.Context())
	require.NoError(t, err, "ListConnectedSynchronizers")
	for _, info := range syncInfos {
		if info.SynchronizerID == s.SynchronizerID {
			s.SynchronizerAlias = info.Alias
			break
		}
	}
	require.NotEmpty(t, s.SynchronizerAlias, "could not resolve alias for synchronizer %s", s.SynchronizerID)

	// Step 3: Deploy TWO contracts to the 3-member party.
	kmsCfgs := []client.KMSConfig{
		s.kmsConfigFor(0, "onboarding"),
		s.kmsConfigFor(1, "onboarding"),
		s.kmsConfigFor(2, "onboarding"),
	}

	t.Log("SetupSuite: deploying contract A (preserved through ceremony)")
	srA, _ := s.runContractDeployFlow(t, s.PartyID, "resync-pre-deploy-a", kmsCfgs)
	s.ContractID = srA.Output.ContractID
	s.PackageIDs = srA.Output.PackageIDs
	require.NotEmpty(t, s.ContractID, "contract A ID must be set")

	t.Log("SetupSuite: deploying contract B (archived mid-ceremony)")
	srB, _ := s.runContractDeployFlow(t, s.PartyID, "resync-pre-deploy-b", kmsCfgs)
	s.ContractIDB = srB.Output.ContractID
	require.NotEmpty(t, s.ContractIDB, "contract B ID must be set")
	t.Logf("Pre-deploy contract IDs: A=%s B=%s", s.ContractID, s.ContractIDB)

	// Step 4: Kick participant 3 → 2-member party (threshold=2).
	t.Log("SetupSuite: kicking participant 3")
	s.performKick()
}

// TestAddParticipantWithAcsFlow is inherited but we don't want to run it on
// this suite — the happy-path is covered by the parent suite. Override to
// no-op.
func (s *AddParticipantWithAcsResyncFlowTestSuite) TestAddParticipantWithAcsFlow() {
	s.T().Log("Skipped — covered by AddParticipantWithAcsFlowTestSuite")
}

// TestResyncOnArchive runs the add-participant-with-acs ceremony and archives
// contract B between Run 6 (ACS export) and Run 7 (target import + reconnect).
// After Run 7, the target's ACS must contain only contract A.
func (s *AddParticipantWithAcsResyncFlowTestSuite) TestResyncOnArchive() {
	t := s.T()

	// Keep containers alive on failure for inspection.
	defer func() {
		if t.Failed() && os.Getenv("KEEP_CONTAINERS") != "" {
			t.Logf("KEEP_CONTAINERS set and test failed — pausing before teardown.")
			select {}
		}
	}()

	actors := s.Actors
	synchronizerID := s.SynchronizerID
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]

	// Verify pre-state.
	p2pState, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after kick")
	require.Len(t, p2pState.Participants, 2, "should have 2 participants after kick")

	// ── Ceremony setup ────────────────────────────────────────────────────
	sourceUID := actors[1].uid
	newParticipantUID := actors[2].uid
	addInput := addparticipantwithacs.AddParticipantWithAcsInput{
		DecentralizedPartyID: s.PartyID,
		NewParticipantID:     newParticipantUID,
		NamespaceName:        s.uniqueName("add-with-acs-resync-ns"),
		SynchronizerID:       synchronizerID,
		SynchronizerAlias:    s.SynchronizerAlias,
		SourceParticipantID:  sourceUID,
		NewThreshold:         2,
	}

	acsKMS := s.kmsConfigFor(2, "add-participant-with-acs-resync")
	addDeps := [3]ceremony.CantonDeps{
		s.OnboardingDeps(0),
		s.OnboardingDeps(1),
		s.depsFor(2, acsKMS),
	}

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	// Runs 1–6: same as the happy-path test.
	for i, deps := range []ceremony.CantonDeps{addDeps[2], addDeps[0], addDeps[1], addDeps[0], addDeps[2], addDeps[1]} {
		t.Logf("Resync ceremony run %d", i+1)
		_, runErr := operations.ExecuteSequence(newBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, deps, addInput)
		require.ErrorContains(t, runErr, addparticipantwithacs.ErrThresholdNotMet.Error(),
			"run %d: expected threshold-not-met", i+1)
	}

	// ── Inject: update contract B with p1+p2 (between Run 6 and Run 7) ────
	// At this point the source has the ACS snapshot containing B (contract ID
	// B-old). We exercise the consuming `UpdateValue` choice on B, which
	// archives B-old and creates B-new with a fresh contract ID. The
	// synchronizer must deliver both events (archive of B-old + create of
	// B-new) to the target during reconnect — i.e. target's ACS must converge
	// to {A, B-new}, NOT {A, B-old}.
	t.Log("Updating contract B on source-side (p1+p2) before target reconnects")
	updateSigners := []signerActor{
		{actorIndex: 0, kmsCfg: s.kmsConfigFor(0, "onboarding")},
		{actorIndex: 1, kmsCfg: s.kmsConfigFor(1, "onboarding")},
	}
	templateID := &apiv2.Identifier{
		PackageId:  s.PackageIDs[0],
		ModuleName: "Main",
		EntityName: "DisclosedTarget",
	}
	updateArg := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "newValue", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 99}}},
		},
	}}}
	newContractB := s.exerciseChoiceAsParty(t, 1 /*submitter=p2*/, updateSigners,
		s.PartyID, synchronizerID,
		templateID, s.ContractIDB, "UpdateValue", updateArg,
	)
	require.NotEmpty(t, newContractB, "UpdateValue must return the new contract ID")
	require.NotEqual(t, s.ContractIDB, newContractB, "new contract ID must differ from the consumed one")
	s.UpdatedContractIDB = newContractB
	t.Logf("Contract B updated: old=%s new=%s", s.ContractIDB, newContractB)

	// ── Run 7: target imports snapshot + reconnects + clears flag ──────────
	t.Log("Resync ceremony run 7: target imports + reconnects + clears flag")
	addResult, err := operations.ExecuteSequence(newBundle(), addparticipantwithacs.AddParticipantWithAcsSequence, addDeps[2], addInput)
	require.NoError(t, err, "run 7 should succeed")
	assert.True(t, addResult.Output.State.OnboardingFlagCleared, "onboarding flag should be cleared")

	// ── Assert: target sees {A, newContractB}, NOT the snapshot's old B ───
	p3Ledger, p3Conn := s.NewLedgerClient(s.chain.Participants[2])
	t.Cleanup(func() { _ = p3Conn.Close() })
	if s.chain.Participants[2].UserID != "" {
		err = p3Ledger.GrantPartyRights(t.Context(), s.chain.Participants[2].UserID, s.PartyID)
		require.NoError(t, err, "GrantPartyRights to p3 user post-ceremony")
	}

	var p3Contracts []*apiv2.CreatedEvent
	var p3IDs []string
	targetOK := assert.Eventually(t, func() bool {
		var qErr error
		p3Contracts, qErr = p3Ledger.GetActiveContractsByTemplateForParty(
			t.Context(), s.PartyID, disclosedTargetPackageName, "Main", "DisclosedTarget",
		)
		if qErr != nil {
			t.Logf("p3 query not yet ready: %v", qErr)
			return false
		}

		p3IDs = activeContractIDs(p3Contracts)
		seen := contractIDSet(p3IDs)

		return seen[s.ContractID] && seen[newContractB] && !seen[s.ContractIDB]
	}, 30*time.Second, 1*time.Second, "p3 should see A and newContractB, but not old B, after resync")
	require.Truef(t, targetOK,
		"p3 active contracts=%v; expected A=%s and newContractB=%s active, oldB=%s absent",
		p3IDs, s.ContractID, newContractB, s.ContractIDB)
	seenOnTarget := contractIDSet(p3IDs)
	assert.True(t, seenOnTarget[s.ContractID], "p3 should see contract A (unchanged through ceremony)")
	assert.True(t, seenOnTarget[newContractB], "p3 should see newContractB (created during resync window)")
	assert.False(t, seenOnTarget[s.ContractIDB], "p3 must NOT see the archived snapshot version of B")

	// Sanity: decNS still resolves on p1 (proves topology didn't regress).
	_, err = actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after resync ceremony")
}

// TestNewParticipantConfirms verifies that p3 (the newly-added participant) is
// a real confirmer in the post-ceremony party. After the resync ceremony, p2
// is disconnected from the synchronizer and a transaction is submitted with
// signatures from p1 and p3 only. With threshold=2 and p2 unreachable, the
// only way the submission commits is if p3 contributes a confirmation.
//
// Depends on TestResyncOnArchive having already completed (p3 added back).
func (s *AddParticipantWithAcsResyncFlowTestSuite) TestNewParticipantConfirms() {
	t := s.T()

	defer func() {
		if t.Failed() && os.Getenv("KEEP_CONTAINERS") != "" {
			t.Logf("KEEP_CONTAINERS set and test failed — pausing before teardown.")
			select {}
		}
	}()

	// Sanity-check: party should now be 3-member (TestResyncOnArchive
	// completed). If running this test standalone (e.g. via -run), the
	// dependency was not satisfied — skip with a clear message.
	p2pState, err := s.Actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, s.SynchronizerID)
	require.NoError(t, err, "GetP2P precondition")
	if len(p2pState.Participants) != 3 {
		t.Skipf("TestNewParticipantConfirms requires the resync ceremony to have completed first (party has %d members)", len(p2pState.Participants))
	}

	// Disconnect p2 from the synchronizer.
	t.Log("Disconnecting p2 from synchronizer")
	require.NoError(t, s.Actors[1].deps.Client.DisconnectSynchronizer(t.Context(), s.SynchronizerAlias),
		"DisconnectSynchronizer on p2")
	t.Cleanup(func() {
		// Best-effort reconnect on cleanup.
		_ = s.Actors[1].deps.Client.ReconnectSynchronizer(s.T().Context(), s.SynchronizerAlias)
	})

	require.NotEmpty(t, s.UpdatedContractIDB,
		"TestNewParticipantConfirms requires TestResyncOnArchive to record B's replacement contract ID")

	// Archive contract A using only p1 + p3 signatures. Submit via p1 (since
	// p2 is disconnected and p3 may still be settling its rights).
	signers := []signerActor{
		{actorIndex: 0, kmsCfg: s.kmsConfigFor(0, "onboarding")},
		{actorIndex: 2, kmsCfg: s.kmsConfigFor(2, "add-participant-with-acs-resync")},
	}
	templateID := &apiv2.Identifier{
		PackageId:  s.PackageIDs[0],
		ModuleName: "Main",
		EntityName: "DisclosedTarget",
	}

	t.Log("Submitting archive of contract A with p1+p3 signatures (p2 disconnected)")
	_ = s.exerciseChoiceAsParty(t, 0 /*submitter=p1*/, signers,
		s.PartyID, s.SynchronizerID,
		templateID, s.ContractID, "Archive", nil,
	)

	// Verify on p1 and p3: contract A is gone while B's replacement remains.
	p1Ledger, p1Conn := s.NewLedgerClient(s.chain.Participants[0])
	t.Cleanup(func() { _ = p1Conn.Close() })
	p3Ledger, p3Conn := s.NewLedgerClient(s.chain.Participants[2])
	t.Cleanup(func() { _ = p3Conn.Close() })

	for idx, lc := range []client.LedgerClient{p1Ledger, p3Ledger} {
		participantLabel := []string{"p1", "p3"}[idx]
		var ids []string
		ok := assert.Eventually(t, func() bool {
			cs, qErr := lc.GetActiveContractsByTemplateForParty(t.Context(), s.PartyID, disclosedTargetPackageName, "Main", "DisclosedTarget")
			if qErr != nil {
				return false
			}

			ids = activeContractIDs(cs)
			seen := contractIDSet(ids)

			return !seen[s.ContractID] && seen[s.UpdatedContractIDB] && !seen[s.ContractIDB]
		}, 15*time.Second, 500*time.Millisecond,
			"%s should see contract A archived while B's replacement remains active", participantLabel)
		require.Truef(t, ok,
			"%s active contracts=%v; expected A=%s absent, updatedB=%s active, oldB=%s absent",
			participantLabel, ids, s.ContractID, s.UpdatedContractIDB, s.ContractIDB)
	}

	// Reconnect p2 and verify it converges to the same ACS.
	t.Log("Reconnecting p2; verifying it converges to the post-archive state")
	require.NoError(t, s.Actors[1].deps.Client.ReconnectSynchronizer(t.Context(), s.SynchronizerAlias),
		"ReconnectSynchronizer on p2")

	p2Ledger, p2Conn := s.NewLedgerClient(s.chain.Participants[1])
	t.Cleanup(func() { _ = p2Conn.Close() })
	var p2IDs []string
	p2OK := assert.Eventually(t, func() bool {
		cs, qErr := p2Ledger.GetActiveContractsByTemplateForParty(t.Context(), s.PartyID, disclosedTargetPackageName, "Main", "DisclosedTarget")
		if qErr != nil {
			return false
		}

		p2IDs = activeContractIDs(cs)
		seen := contractIDSet(p2IDs)

		return !seen[s.ContractID] && seen[s.UpdatedContractIDB] && !seen[s.ContractIDB]
	}, 30*time.Second, 1*time.Second, "p2 should converge to post-archive state after reconnecting")
	require.Truef(t, p2OK,
		"p2 active contracts=%v; expected A=%s absent, updatedB=%s active, oldB=%s absent",
		p2IDs, s.ContractID, s.UpdatedContractIDB, s.ContractIDB)
}

func activeContractIDs(contracts []*apiv2.CreatedEvent) []string {
	ids := make([]string, 0, len(contracts))
	for _, c := range contracts {
		if id := c.GetContractId(); id != "" {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)

	return ids
}

func contractIDSet(ids []string) map[string]bool {
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		seen[id] = true
	}

	return seen
}
