package tests

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chainlink/canton-party-ceremony/ceremony/contractdeploy"
	"github.com/chainlink/canton-party-ceremony/internal/client"
)

type ContractDeployFlowTestSuite struct {
	OnboardingFlowTestSuite

	LedgerClients []client.LedgerClient
}

func (s *ContractDeployFlowTestSuite) SetupSuite() {
	t := s.T()

	// Run the full onboarding ceremony so s.PartyID is populated.
	// CeremonyTestSuite.SetupSuite (called inside) loads the chain and stores it
	// in s.chain, so we can reuse the same endpoints below.
	s.OnboardingFlowTestSuite.SetupSuite()
	s.performOnboarding(t, operations.NewMemoryReporter())

	chain := s.chain

	ledgerClients := make([]client.LedgerClient, len(chain.Participants))
	for i, p := range chain.Participants {
		lc, conn := s.NewLedgerClient(p)
		t.Cleanup(func() { _ = conn.Close() })
		ledgerClients[i] = lc
	}
	s.LedgerClients = ledgerClients
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
//     ErrSigningNotImplemented
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
		{AdminClient: s.Actors[0].client, LedgerClient: s.LedgerClients[0], DARLoader: darLoader, Logger: logger.Test(t)},
		{AdminClient: s.Actors[1].client, LedgerClient: s.LedgerClients[1], DARLoader: darLoader, Logger: logger.Test(t)},
		{AdminClient: s.Actors[2].client, LedgerClient: s.LedgerClients[2], DARLoader: darLoader, Logger: logger.Test(t)},
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

	// Run 3 (p3): uploads DARs (3/3) → verifies party → prepares → ErrSigningNotImplemented
	t.Log("Run 3: p3 uploads DARs (3/3) + verify + prepare → signing stub")
	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[2], input)
	require.ErrorContains(t, err, contractdeploy.ErrSigningNotImplemented.Error(),
		"run 3: expected signing-not-implemented")

	// Verify partial output: package IDs and prepared transaction hash.
	assert.NotEmpty(t, sr.Output.PackageIDs, "should have at least one package ID")
	assert.NotEmpty(t, sr.Output.PreparedTransactionHash, "should have a prepared transaction hash")
	t.Logf("Package IDs: %v", sr.Output.PackageIDs)
	t.Logf("Prepared TX hash: %s", sr.Output.PreparedTransactionHash)

	// Verify idempotency: re-run should produce the same result from cache.
	t.Log("Run 4: p1 idempotency check")
	srCached, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrSigningNotImplemented.Error(),
		"run 4: idempotent re-run should reach signing stub")
	assert.Equal(t, sr.Output.PackageIDs, srCached.Output.PackageIDs,
		"cached package IDs should match")
	assert.Equal(t, sr.Output.PreparedTransactionHash, srCached.Output.PreparedTransactionHash,
		"cached prepared TX hash should match")
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
