package contractdeploy_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/chainlink/canton-party-ceremony/ceremony/contractdeploy"
	"github.com/chainlink/canton-party-ceremony/internal/client"
)

// ── Mock Canton client (Admin API) ──────────────────────────────────────────

type mockAdminClient struct {
	participantID string
}

func newMockAdminClient(pid string) *mockAdminClient {
	return &mockAdminClient{participantID: pid}
}

func (m *mockAdminClient) GetParticipantUID(_ context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockAdminClient) GetParticipantID(_ context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockAdminClient) GenerateSigningKey(_ context.Context, _ string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return &cryptov30.SigningPublicKey{}, nil
}

func (m *mockAdminClient) GetNamespaceFingerprint(_ context.Context, _ string, _ string, _ []string) (string, error) {
	return "mock-ns-fp", nil
}

func (m *mockAdminClient) Authorize(_ context.Context, _ uint32, _ *protov30.TopologyMapping, _ string, _ bool, _ ...string) (*protov30.SignedTopologyTransaction, error) {
	return &protov30.SignedTopologyTransaction{}, nil
}

func (m *mockAdminClient) SignTransactions(_ context.Context, txs []*protov30.SignedTopologyTransaction, _ string) ([]*protov30.SignedTopologyTransaction, error) {
	return txs, nil
}

func (m *mockAdminClient) AddTransactions(_ context.Context, _ []*protov30.SignedTopologyTransaction, _ string) error {
	return nil
}

func (m *mockAdminClient) DNSExists(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func (m *mockAdminClient) NSDExists(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func (m *mockAdminClient) P2PExists(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func (m *mockAdminClient) GetDNS(_ context.Context, _ string, _ string) (*client.DNSState, error) {
	return &client.DNSState{}, nil
}

func (m *mockAdminClient) GetP2P(_ context.Context, _ string, _ string) (*client.P2PState, error) {
	return &client.P2PState{}, nil
}

func (m *mockAdminClient) ListDecentralizedNamespaces(_ context.Context, _ string) ([]*client.DNSState, error) {
	return []*client.DNSState{}, nil
}

func (m *mockAdminClient) UploadDar(_ context.Context, darBytes []byte) (string, error) {
	id := hex.EncodeToString(darBytes[:min(8, len(darBytes))])
	return "pkg-" + id, nil
}

// ── Mock Ledger client ──────────────────────────────────────────────────────

type mockLedgerClient struct {
	partyExists bool
}

func (m *mockLedgerClient) PartyExists(_ context.Context, _ string) (bool, error) {
	return m.partyExists, nil
}

func (m *mockLedgerClient) PrepareSubmission(
	_ context.Context,
	_ []*apiv2.Command,
	_ []string,
	_ []string,
	_ string,
) (*interactive.PrepareSubmissionResponse, error) {
	return &interactive.PrepareSubmissionResponse{
		PreparedTransactionHash: []byte{0xab, 0xcd, 0xef, 0x12},
	}, nil
}

func (m *mockLedgerClient) GetActiveContractsByTemplate(
	_ context.Context, _ string, _ string, _ string, _ string,
) ([]*apiv2.CreatedEvent, error) {
	return nil, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func fakeDARLoader() contractdeploy.DARLoader {
	return func(name, version string) ([]byte, error) {
		return []byte(fmt.Sprintf("fake-dar-%s-%s", name, version)), nil
	}
}

func newDeps(participantID string, partyExists bool) contractdeploy.ContractDeployDeps {
	return contractdeploy.ContractDeployDeps{
		AdminClient:  newMockAdminClient(participantID),
		LedgerClient: &mockLedgerClient{partyExists: partyExists},
		DARLoader:    fakeDARLoader(),
		Logger:       logger.Nop(),
	}
}

func baseInput(pkgs []contractdeploy.PackageRef) contractdeploy.ContractDeployInput {
	return contractdeploy.ContractDeployInput{
		DecentralizedPartyID: "test-party::1220abcdef",
		SynchronizerID:       "global",
		Participants:         []string{"p1", "p2"},
		Packages:             pkgs,
		TemplateModule:       "MCMS.Main",
		TemplateEntity:       "MCMS",
		ContractArgs:         `{}`,
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

func TestContractDeploySequence_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// Run 1: p1 uploads DARs -> threshold not met (1/2)
	_, err := operations.ExecuteSequence(
		newBundle(),
		contractdeploy.ContractDeploySequence,
		newDeps("p1", true),
		input,
	)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())
}

func TestContractDeploySequence_SigningNotImplemented(t *testing.T) {
	t.Parallel()

	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// Run 1: p1 uploads DARs -> threshold not met (1/2)
	_, err := operations.ExecuteSequence(
		newBundle(),
		contractdeploy.ContractDeploySequence,
		newDeps("p1", true),
		input,
	)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	// Run 2: p2 uploads DARs -> all uploaded, verifies party, prepares,
	// then hits signing not implemented
	_, err = operations.ExecuteSequence(
		newBundle(),
		contractdeploy.ContractDeploySequence,
		newDeps("p2", true),
		input,
	)
	require.ErrorContains(t, err, contractdeploy.ErrSigningNotImplemented.Error())
}

func TestContractDeploySequence_Idempotent(t *testing.T) {
	t.Parallel()

	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// Complete through preparation (will stop at signing)
	_, err1 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err1, contractdeploy.ErrThresholdNotMet.Error())

	_, err2 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p2", true), input)
	require.ErrorContains(t, err2, contractdeploy.ErrSigningNotImplemented.Error())

	// Run 3: idempotent same reporter, returns same cached result
	sr3, err3 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err3, contractdeploy.ErrSigningNotImplemented.Error())

	assert.NotEmpty(t, sr3.Output.PackageIDs, "package IDs should be populated")
	assert.NotEmpty(t, sr3.Output.PreparedTransactionHash, "prepared transaction hash should be populated")
}

func TestContractDeploySequence_PartyNotFound(t *testing.T) {
	t.Parallel()

	input := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: "test-party::1220abcdef",
		SynchronizerID:       "global",
		Participants:         []string{"p1"},
		Packages:             []contractdeploy.PackageRef{{Name: "mcms", Version: "current"}},
		TemplateModule:       "MCMS.Main",
		TemplateEntity:       "MCMS",
		ContractArgs:         `{}`,
	}

	b := optest.NewBundle(t)
	_, err := operations.ExecuteSequence(
		b,
		contractdeploy.ContractDeploySequence,
		newDeps("p1", false),
		input,
	)
	require.ErrorContains(t, err, "party")
	require.ErrorContains(t, err, "not found")
}

func TestContractDeploySequence_MultipleDARs(t *testing.T) {
	t.Parallel()

	input := baseInput([]contractdeploy.PackageRef{
		{Name: "mcms", Version: "current"},
		{Name: "globalconfig", Version: "1.0.0"},
	})

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// Both participants upload both DARs
	_, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p2", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrSigningNotImplemented.Error())

	assert.Len(t, sr.Output.PackageIDs, 2, "should have uploaded 2 DARs")
}

// ── Interface compliance ────────────────────────────────────────────────────

func TestGRPCLedgerClientImplementsInterface(t *testing.T) {
	t.Parallel()
	var _ client.LedgerClient = (*client.GRPCLedgerClient)(nil)
}
