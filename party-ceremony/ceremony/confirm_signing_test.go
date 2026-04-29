package ceremony_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// ── Minimal mock Canton client for signing tests ────────────────────────────

type signTestClient struct {
	participantID string
	signCalled    bool
}

func (m *signTestClient) GetParticipantUID(context.Context) (string, error) {
	return m.participantID, nil
}

func (m *signTestClient) GetParticipantID(context.Context) (string, error) {
	return m.participantID, nil
}

func (m *signTestClient) SignTransactions(_ context.Context, txs []*protov30.SignedTopologyTransaction, _ string) ([]*protov30.SignedTopologyTransaction, error) {
	m.signCalled = true
	result := make([]*protov30.SignedTopologyTransaction, len(txs))
	for i, tx := range txs {
		result[i] = &protov30.SignedTopologyTransaction{
			Transaction: tx.GetTransaction(),
			Signatures: append(tx.GetSignatures(), &cryptov30.Signature{
				SignedBy:  m.participantID,
				Signature: []byte("fake-sig"),
			}),
		}
	}

	return result, nil
}

// Unused interface methods — required to satisfy client.CantonClient.
var errNotImplemented = errors.New("not implemented")

func (m *signTestClient) GenerateSigningKey(context.Context, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errNotImplemented
}
func (m *signTestClient) RegisterKmsSigningKey(context.Context, string, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errNotImplemented
}
func (m *signTestClient) GetNamespaceFingerprint(context.Context, string, string, []string) (string, error) {
	return "", nil
}
func (m *signTestClient) GetNamespaceKeyName(context.Context, string, []string) (string, error) {
	return "", nil
}
func (m *signTestClient) GetProtocolKeyFingerprint(context.Context, []string) (string, string, error) {
	return "", "", nil
}
func (m *signTestClient) Authorize(context.Context, uint32, *protov30.TopologyMapping, string, bool, ...string) (*protov30.SignedTopologyTransaction, error) {
	return nil, errNotImplemented
}
func (m *signTestClient) AddTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) error {
	return nil
}
func (m *signTestClient) DNSExists(context.Context, string, string) (bool, error) { return false, nil }
func (m *signTestClient) NSDExists(context.Context, string, string) (bool, error) { return false, nil }
func (m *signTestClient) P2PExists(context.Context, string, string) (bool, error) { return false, nil }
func (m *signTestClient) GetDNS(context.Context, string, string) (*client.DNSState, error) {
	return nil, errNotImplemented
}
func (m *signTestClient) GetP2P(context.Context, string, string) (*client.P2PState, error) {
	return nil, errNotImplemented
}
func (m *signTestClient) ListDecentralizedNamespaces(context.Context, string) ([]*client.DNSState, error) {
	return nil, nil
}
func (m *signTestClient) UploadDar(context.Context, []byte) (string, error) { return "", nil }

// ── Test helpers ─────────────────────────────────────────────────────────────

func buildTestDNSTx(t *testing.T) string {
	t.Helper()
	mapping := &protov30.TopologyMapping{
		Mapping: &protov30.TopologyMapping_DecentralizedNamespaceDefinition{
			DecentralizedNamespaceDefinition: &protov30.DecentralizedNamespaceDefinition{
				DecentralizedNamespace: "test-ns",
				Threshold:              2,
				Owners:                 []string{"fp-a", "fp-b"},
			},
		},
	}
	innerTx := &protov30.TopologyTransaction{
		Operation: protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE,
		Serial:    1,
		Mapping:   mapping,
	}
	innerBytes, _ := proto.Marshal(innerTx)
	stx := &protov30.SignedTopologyTransaction{
		Transaction: innerBytes,
		Proposal:    true,
	}
	stxBytes, _ := proto.Marshal(stx)
	txB64 := base64.StdEncoding.EncodeToString(stxBytes)

	return txB64
}

func txHash(txB64 string) string {
	raw, _ := base64.StdEncoding.DecodeString(txB64)
	h := sha256.Sum256(raw)

	return base64.StdEncoding.EncodeToString(h[:])
}

// ── Tests ────────────────────────────────────────────────────────────────────

func TestSignDNSProposalOp_ConfirmRejected(t *testing.T) {
	t.Parallel()

	txB64 := buildTestDNSTx(t)
	mockClient := &signTestClient{participantID: "p1"}
	deps := ceremony.CantonDeps{
		Client:    mockClient,
		Logger:    logger.Nop(),
		Confirmer: ceremony.AlwaysRejectConfirmer{},
	}
	bundle := operations.NewBundle(t.Context, logger.Nop(), operations.NewMemoryReporter())

	input := topology.SignDNSProposalInput{
		ParticipantID:      "p1",
		ProposalHashSHA256: txHash(txB64),
		DNSTxB64:           txB64,
		SynchronizerID:     "global",
	}

	_, err := operations.ExecuteOperation(bundle, topology.SignDNSProposalOp, deps, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), ceremony.ErrUserRejected.Error())
	assert.False(t, mockClient.signCalled, "SignTransactions must NOT be called when user rejects")
}

func TestSignDNSProposalOp_ConfirmAccepted(t *testing.T) {
	t.Parallel()

	txB64 := buildTestDNSTx(t)
	mockClient := &signTestClient{participantID: "p1"}
	deps := ceremony.CantonDeps{
		Client:    mockClient,
		Logger:    logger.Nop(),
		Confirmer: ceremony.NoOpConfirmer{},
	}
	bundle := operations.NewBundle(t.Context, logger.Nop(), operations.NewMemoryReporter())

	input := topology.SignDNSProposalInput{
		ParticipantID:      "p1",
		ProposalHashSHA256: txHash(txB64),
		DNSTxB64:           txB64,
		SynchronizerID:     "global",
	}

	out, err := operations.ExecuteOperation(bundle, topology.SignDNSProposalOp, deps, input)
	require.NoError(t, err)
	assert.True(t, mockClient.signCalled, "SignTransactions must be called when user accepts")
	assert.Equal(t, "p1", out.Output.ParticipantID)
}

func TestSignDNSProposalOp_NilConfirmerSkipsPrompt(t *testing.T) {
	t.Parallel()

	txB64 := buildTestDNSTx(t)
	mockClient := &signTestClient{participantID: "p1"}
	deps := ceremony.CantonDeps{
		Client: mockClient,
		Logger: logger.Nop(),
		// Confirmer is nil — should skip confirmation
	}
	bundle := operations.NewBundle(t.Context, logger.Nop(), operations.NewMemoryReporter())

	input := topology.SignDNSProposalInput{
		ParticipantID:      "p1",
		ProposalHashSHA256: txHash(txB64),
		DNSTxB64:           txB64,
		SynchronizerID:     "global",
	}

	_, err := operations.ExecuteOperation(bundle, topology.SignDNSProposalOp, deps, input)
	require.NoError(t, err)
	assert.True(t, mockClient.signCalled)
}
