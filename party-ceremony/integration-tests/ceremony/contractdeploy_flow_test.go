package tests

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chainlink/canton-party-ceremony/ceremony/contractdeploy"
	"github.com/chainlink/canton-party-ceremony/ceremony/onboarding"
	"github.com/chainlink/canton-party-ceremony/internal/client"
)

type ContractDeployFlowTestSuite struct {
	OnboardingFlowTestSuite

	LedgerClients []client.LedgerClient
	Signers       []client.TransactionSigner
}

func (s *ContractDeployFlowTestSuite) SetupSuite() {
	t := s.T()

	// Run the full onboarding ceremony so s.PartyID is populated.
	// CeremonyTestSuite.SetupSuite (called inside) loads the chain and stores it
	// in s.chain, so we can reuse the same endpoints below.
	s.OnboardingFlowTestSuite.SetupSuite()
	// Capture the reporter so we can extract DAML key fingerprints below.
	onboardingReporter := operations.NewMemoryReporter()
	s.performOnboarding(t, onboardingReporter)

	chain := s.chain

	// In JWT-auth environments (CTF), decentralized parties created via topology
	// are not automatically granted to any user. Grant actAs+readAs rights so that
	// PrepareSubmission and related Ledger API calls are authorized.
	for i, p := range chain.Participants {
		if p.UserID == "" {
			continue // no-auth env; rights enforced
		}
		_, err := p.LedgerServices.Admin.UserManagement.GrantUserRights(
			t.Context(),
			&adminv2.GrantUserRightsRequest{
				UserId: p.UserID,
				Rights: []*adminv2.Right{
					{Kind: &adminv2.Right_CanActAs_{CanActAs: &adminv2.Right_CanActAs{Party: s.PartyID}}},
					{Kind: &adminv2.Right_CanReadAs_{CanReadAs: &adminv2.Right_CanReadAs{Party: s.PartyID}}},
				},
			},
		)
		require.NoError(t, err, "failed to grant party rights for participant %d (user=%s, party=%s)", i+1, p.UserID, s.PartyID)
		t.Logf("Participant %d (%s): granted actAs+readAs for party %s", i+1, p.UserID, s.PartyID)
	}

	ledgerClients := make([]client.LedgerClient, len(chain.Participants))
	for i, p := range chain.Participants {
		lc, conn := s.NewLedgerClient(p)
		t.Cleanup(func() { _ = conn.Close() })
		ledgerClients[i] = lc
	}
	s.LedgerClients = ledgerClients

	// Extract the DAML (PROTOCOL) key fingerprint for each participant from the
	// onboarding ceremony reports. Using the reporter avoids relying on
	// ListMyKeys ordering, which is non-deterministic in persistent vaults.
	damlFingerprints := map[string]string{}
	allReports, err := onboardingReporter.GetReports()
	require.NoError(t, err, "getting onboarding reports")
	for _, r := range allReports {
		if r.Def.ID != "onboarding/canton-ceremony/create-member-key" {
			continue
		}
		if out, ok := r.Output.(onboarding.CreateMemberKeyOutput); ok && out.DamlKeyFingerprint != "" {
			damlFingerprints[out.ParticipantID] = out.DamlKeyFingerprint
		}
	}

	// Create a VaultSigner for each participant using the DAML signing key
	// fingerprint from the ceremony (the key registered in PartyToParticipant.PartySigningKeys).
	signers := make([]client.TransactionSigner, len(chain.Participants))
	for i, p := range chain.Participants {
		adminConn, err := grpc.NewClient(p.Endpoints.AdminAPIURL,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "dial admin API for vault (participant %d)", i+1)
		t.Cleanup(func() { _ = adminConn.Close() })

		fp, hasFP := damlFingerprints[s.ParticipantIDs[i]]
		require.True(t, hasFP, "DAML key fingerprint not found for participant %d (%s)", i+1, s.ParticipantIDs[i])

		vault := cryptoadminv30.NewVaultServiceClient(adminConn)
		signer, err := client.NewVaultSigner(context.Background(), vault, fp)
		require.NoError(t, err, "create vault signer (participant %d)", i+1)
		signers[i] = signer
		t.Logf("Participant %d VaultSigner ready, fingerprint=%s", i+1, fp)
	}
	s.Signers = signers
}

// TestOnboardingFlow overrides the inherited method to skip re-running the
// full onboarding ceremony (it was already done in SetupSuite).
func (s *ContractDeployFlowTestSuite) TestOnboardingFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Onboarding already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// TestContractDeployFlow validates the contract deployment ceremony against a
// real Canton environment with 3 participants.
//
// Flow:
//   - Run 1 (p1): uploads DARs (1/3) → ErrThresholdNotMet
//   - Run 2 (p2): uploads DARs (2/3) → ErrThresholdNotMet
//   - Run 3 (p3): uploads DARs (3/3) → verifies party → prepares submission →
//     signs (all participants) → executes → verifies contract in ACS
func (s *ContractDeployFlowTestSuite) TestContractDeployFlow() {
	t := s.T()

	// Build contract args: DisclosedTarget { owner: Party, value: Int }
	// The decentralized party is the signatory.
	contractArgs := buildContractArgs(t, s.PartyID)

	input := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: s.PartyID,
		SynchronizerID:       s.SynchronizerID,
		Participants:         s.ParticipantIDs,
		Packages:             []contractdeploy.PackageRef{{Name: "test-test", Version: "0.0.1"}},
		TemplateModule:       "Main",
		TemplateEntity:       "DisclosedTarget",
		ContractArgs:         contractArgs,
	}

	darDir := filepath.Join("..", "..", "..", "contracts", "dars")
	darLoader := contractdeploy.FileDARLoader(darDir)

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	deps := [3]contractdeploy.ContractDeployDeps{
		{AdminClient: s.Actors[0].client, LedgerClient: s.LedgerClients[0], DARLoader: darLoader, Signer: s.Signers[0], Logger: logger.Test(t)},
		{AdminClient: s.Actors[1].client, LedgerClient: s.LedgerClients[1], DARLoader: darLoader, Signer: s.Signers[1], Logger: logger.Test(t)},
		{AdminClient: s.Actors[2].client, LedgerClient: s.LedgerClients[2], DARLoader: darLoader, Signer: s.Signers[2], Logger: logger.Test(t)},
	}

	// Run 1 (p1): uploads DARs (1/3) → threshold not met
	t.Log("Run 1: p1 uploads DARs (1/3)")
	_, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(),
		"run 1: expected threshold-not-met (1/3)")

	// Run 2 (p2): uploads DARs (2/3) → threshold not met
	t.Log("Run 2: p2 uploads DARs (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[1], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(),
		"run 2: expected threshold-not-met (2/3)")

	// Run 3 (p3): uploads DARs (3/3) → verifies → prepares → p3 signs (1/3) → threshold not met for signing
	t.Log("Run 3: p3 uploads DARs (3/3) + prepares + signs (1/3)")
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[2], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(),
		"run 3: expected threshold-not-met for signing (1/3)")

	// Run 4 (p1): p1 signs (2/3) → threshold not met for signing
	t.Log("Run 4: p1 signs (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(),
		"run 4: expected threshold-not-met for signing (2/3)")

	// Run 5 (p2): p2 signs (3/3) → executes → verifies contract
	t.Log("Run 5: p2 signs (3/3) + executes + verifies")
	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[1], input)
	require.NoError(t, err, "run 5: expected full success")

	// Verify output.
	assert.NotEmpty(t, sr.Output.PackageIDs, "should have at least one package ID")
	assert.NotEmpty(t, sr.Output.PreparedTransactionHash, "should have a prepared transaction hash")
	assert.NotEmpty(t, sr.Output.ContractID, "should have a deployed contract ID")
	t.Logf("Package IDs: %v", sr.Output.PackageIDs)
	t.Logf("Prepared TX hash: %s", sr.Output.PreparedTransactionHash)
	t.Logf("Contract ID: %s", sr.Output.ContractID)

	// Verify idempotency: re-run should produce the same cached result.
	t.Log("Run 6: p1 idempotency check")
	srCached, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.NoError(t, err, "run 6: idempotent re-run should succeed")
	assert.Equal(t, sr.Output.PackageIDs, srCached.Output.PackageIDs,
		"cached package IDs should match")
	assert.Equal(t, sr.Output.ContractID, srCached.Output.ContractID,
		"cached contract ID should match")
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// buildContractArgs constructs the protobuf JSON for a DisclosedTarget contract:
//
//	template DisclosedTarget with
//	    owner : Party
//	    value : Int
//
// The JSON is in proto JSON format for protojson.Unmarshal.
func buildContractArgs(t *testing.T, partyID string) string {
	t.Helper()
	args := map[string]any{
		"fields": []map[string]any{
			{"label": "owner", "value": map[string]any{"party": partyID}},
			{"label": "value", "value": map[string]any{"int64": "42"}},
		},
	}
	b, err := json.Marshal(args)
	require.NoError(t, err, "marshalling contract args")

	return string(b)
}
