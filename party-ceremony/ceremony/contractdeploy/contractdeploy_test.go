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

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
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

func (m *mockAdminClient) RegisterKmsSigningKey(_ context.Context, _ string, _ string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return &cryptov30.SigningPublicKey{}, nil
}

func (m *mockAdminClient) GetNamespaceFingerprint(_ context.Context, _ string, _ string, _ []string) (string, error) {
	return "mock-ns-fp", nil
}

func (m *mockAdminClient) GetNamespaceKeyName(_ context.Context, _ string, _ []string) (string, error) {
	return "mock-ns-key", nil
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
	return &client.P2PState{
		Participants: []client.P2PParticipantInfo{
			{ParticipantUID: "p1"},
			{ParticipantUID: "p2"},
		},
		PartySigningKeys: &client.P2PSigningKeysInfo{
			Keys:      []string{"p1-protocol-key-b64", "p2-protocol-key-b64"},
			Threshold: 2,
		},
	}, nil
}

func (m *mockAdminClient) ListDecentralizedNamespaces(_ context.Context, _ string) ([]*client.DNSState, error) {
	return []*client.DNSState{}, nil
}

func (m *mockAdminClient) GetProtocolKeyFingerprint(_ context.Context, _ []string) (string, string, error) {
	return "mock-protocol-fp", "mock-protocol-key-b64", nil
}

func (m *mockAdminClient) UploadDar(_ context.Context, darBytes []byte) (string, error) {
	id := hex.EncodeToString(darBytes[:min(8, len(darBytes))])
	return "pkg-" + id, nil
}
func (m *mockAdminClient) ExportAcs(_ context.Context, _ []string, _ string, _ int64) ([]byte, error) {
	return nil, nil
}
func (m *mockAdminClient) ImportAcs(_ context.Context, _ []byte, _ string) error    { return nil }
func (m *mockAdminClient) DisconnectSynchronizer(_ context.Context, _ string) error { return nil }
func (m *mockAdminClient) ReconnectSynchronizer(_ context.Context, _ string) error  { return nil }
func (m *mockAdminClient) ListConnectedSynchronizers(_ context.Context) ([]client.SynchronizerInfo, error) {
	return nil, nil
}
func (m *mockAdminClient) ClearPartyOnboardingFlag(_ context.Context, _ string, _ string, _ int64) (bool, error) {
	return true, nil
}
func (m *mockAdminClient) LookupOffsetByTime(_ context.Context, _ int64) (int64, error) {
	return 0, nil
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

func (m *mockLedgerClient) ExecuteSubmission(
	_ context.Context,
	_ *interactive.PreparedTransaction,
	_ *interactive.PartySignatures,
	_ interactive.HashingSchemeVersion,
) (string, error) {
	return "fake-contract-0xdeadbeef", nil
}

func (m *mockLedgerClient) GetActiveContractsByTemplate(
	_ context.Context, _ string, _ string, _ string, _ string,
) ([]*apiv2.CreatedEvent, error) {
	return []*apiv2.CreatedEvent{{ContractId: "fake-contract-0xdeadbeef"}}, nil
}
func (m *mockLedgerClient) GetActiveContractsByTemplateForParty(
	_ context.Context, _ string, _ string, _ string, _ string,
) ([]*apiv2.CreatedEvent, error) {
	return []*apiv2.CreatedEvent{{ContractId: "fake-contract-0xdeadbeef"}}, nil
}
func (m *mockLedgerClient) GrantPartyRights(_ context.Context, _, _ string) error {
	return nil
}

// ── Mock Signer ──────────────────────────────────────────────────────────────

type mockSigner struct{}

func (m *mockSigner) Sign(_ context.Context, hash []byte) (*apiv2.Signature, error) {
	return &apiv2.Signature{
		Format:               apiv2.SignatureFormat_SIGNATURE_FORMAT_RAW,
		Signature:            hash, // echo hash as deterministic fake signature
		SignedBy:             "mock-key-fp",
		SigningAlgorithmSpec: apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_ED25519,
	}, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func fakeDARLoader() ledger.DARLoader {
	return func(name, version string) ([]byte, error) {
		return fmt.Appendf(nil, "fake-dar-%s-%s", name, version), nil
	}
}

func newDeps(participantID string, partyExists bool) ledger.ContractDeployDeps {
	return ledger.ContractDeployDeps{
		AdminClient:  newMockAdminClient(participantID),
		LedgerClient: &mockLedgerClient{partyExists: partyExists},
		DARLoader:    fakeDARLoader(),
		Signer:       &mockSigner{},
		Logger:       logger.Nop(),
	}
}

func newFactoryDeps(
	participantID string,
	partyExists bool,
	factory client.TransactionSignerFactory,
) ledger.ContractDeployDeps {
	return ledger.ContractDeployDeps{
		AdminClient:   newMockAdminClient(participantID),
		LedgerClient:  &mockLedgerClient{partyExists: partyExists},
		DARLoader:     fakeDARLoader(),
		SignerFactory: factory,
		Logger:        logger.Nop(),
	}
}

func baseInput(pkgs []contractdeploy.PackageRef) contractdeploy.ContractDeployInput {
	return contractdeploy.ContractDeployInput{
		DecentralizedPartyID: "test-party::1220abcdef",
		SynchronizerID:       "global",
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

func TestContractDeploySequence_FullFlow(t *testing.T) {
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

	// Run 2: p2 uploads DARs — all DARs done; p2 signs (1/2) → threshold not met for signing
	_, err = operations.ExecuteSequence(
		newBundle(),
		contractdeploy.ContractDeploySequence,
		newDeps("p2", true),
		input,
	)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(), "run 2: p2 signed, p1 pending")

	// Run 3: p1 signs (2/2) → execute → verify
	sr, err := operations.ExecuteSequence(
		newBundle(),
		contractdeploy.ContractDeploySequence,
		newDeps("p1", true),
		input,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, sr.Output.ContractID, "contract ID should be set on full completion")
	assert.NotEmpty(t, sr.Output.PackageIDs)
}

func TestContractDeploySequence_Idempotent(t *testing.T) {
	t.Parallel()

	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// Run 1: p1 uploads DARs only — threshold not met
	_, err1 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err1, contractdeploy.ErrThresholdNotMet.Error())

	// Run 2: p2 uploads DARs — all DARs done; p2 signs (1/2) → threshold not met for signing
	_, err2 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p2", true), input)
	require.ErrorContains(t, err2, contractdeploy.ErrThresholdNotMet.Error(), "run 2: p2 signed, p1 pending")

	// Run 3: p1 signs (2/2) → execute → verify → success
	sr2, err3 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.NoError(t, err3)
	assert.NotEmpty(t, sr2.Output.ContractID)

	// Run 4: idempotent — same shared reporter, returns cached complete result
	sr3, err4 := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.NoError(t, err4)
	assert.Equal(t, sr2.Output.ContractID, sr3.Output.ContractID, "cached contract ID should match")
	assert.NotEmpty(t, sr3.Output.PackageIDs, "package IDs should be populated")
	assert.NotEmpty(t, sr3.Output.PreparedTransactionHash, "prepared transaction hash should be populated")
}

func TestContractDeploySequence_PartyNotFound(t *testing.T) {
	t.Parallel()

	input := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: "test-party::1220abcdef",
		SynchronizerID:       "global",
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

	// Run 1: p1 uploads both DARs → 1/2 → threshold not met
	_, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	// Run 2: p2 uploads both DARs — all done; p2 signs (1/2) → threshold not met for signing
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p2", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	// Run 3: p1 signs (2/2) → execute → verify
	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.NoError(t, err)

	assert.Len(t, sr.Output.PackageIDs, 2, "should have uploaded 2 DARs")
	assert.NotEmpty(t, sr.Output.ContractID)
}

func TestContractDeploySequence_UsesSignerFactoryWithTopologyKeys(t *testing.T) {
	t.Parallel()

	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	type factoryCall struct {
		participantID string
		keys          []string
	}
	var calls []factoryCall
	factory := func(_ context.Context, participantID string, knownSigningKeysB64 []string) (client.TransactionSigner, error) {
		calls = append(calls, factoryCall{
			participantID: participantID,
			keys:          append([]string{}, knownSigningKeysB64...),
		})

		return &mockSigner{}, nil
	}

	_, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newFactoryDeps("p1", true, factory), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newFactoryDeps("p2", true, factory), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newFactoryDeps("p1", true, factory), input)
	require.NoError(t, err)

	require.Len(t, calls, 2)
	assert.Equal(t, "p2", calls[0].participantID)
	assert.Equal(t, []string{"p1-protocol-key-b64", "p2-protocol-key-b64"}, calls[0].keys)
	assert.Equal(t, "p1", calls[1].participantID)
	assert.Equal(t, []string{"p1-protocol-key-b64", "p2-protocol-key-b64"}, calls[1].keys)
}

// ── Interface compliance ────────────────────────────────────────────────────

func TestGRPCLedgerClientImplementsInterface(t *testing.T) {
	t.Parallel()
	var _ client.LedgerClient = (*client.GRPCLedgerClient)(nil)
}
