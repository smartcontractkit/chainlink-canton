package tests

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

type ContractDeployFlowTestSuite struct {
	OnboardingFlowTestSuite
}

func (s *ContractDeployFlowTestSuite) SetupSuite() {
	t := s.T()

	// Run the full onboarding ceremony so s.PartyID is populated.
	// CeremonyTestSuite.SetupSuite (called inside) loads the chain and stores it
	// in s.chain, so we can reuse the same endpoints below.
	s.OnboardingFlowTestSuite.SetupSuite()
	s.performOnboarding(t, operations.NewMemoryReporter())
}

// TestOnboardingFlow overrides the inherited method to skip re-running the
// full onboarding ceremony (it was already done in SetupSuite).
func (s *ContractDeployFlowTestSuite) TestOnboardingFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Onboarding already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// TestContractDeployFlow validates the contract deployment ceremony against a
// real CTF Canton environment with 3 KMS-backed participants.
//
// Flow:
//   - Run 1 (p1): uploads DARs (1/3) → ErrThresholdNotMet
//   - Run 2 (p2): uploads DARs (2/3) → ErrThresholdNotMet
//   - Run 3 (p3): uploads DARs (3/3) → verifies party → prepares submission →
//     signs (1/3) → ErrThresholdNotMet
//   - Run 4 (p1): signs (2/3) → ErrThresholdNotMet
//   - Run 5 (p2): signs (3/3) → executes → verifies contract in ACS
func (s *ContractDeployFlowTestSuite) TestContractDeployFlow() {
	t := s.T()

	kmsCfgs := []client.KMSConfig{
		s.kmsConfigFor(0, "onboarding"),
		s.kmsConfigFor(1, "onboarding"),
		s.kmsConfigFor(2, "onboarding"),
	}
	sr, recorders := s.runContractDeployFlow(t, s.PartyID, "contract-deploy", kmsCfgs)
	assertRecordersUsed(t, recorders)
	t.Logf("Package IDs: %v", sr.Output.PackageIDs)
	t.Logf("Prepared TX hash: %s", sr.Output.PreparedTransactionHash)
	t.Logf("Contract ID: %s", sr.Output.ContractID)
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
