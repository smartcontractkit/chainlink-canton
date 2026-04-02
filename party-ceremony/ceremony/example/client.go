package example

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/chainlink/canton-party-ceremony/internal/client"
)

// ── Mock Canton client ────────────────────────────────────────────────────────

// MockCantonClient is a deterministic, in-memory implementation of
// [CantonClient].  It derives all outputs from a sha-256 of the input
// parameters so that results are stable across runs without real cryptography
// or network access.
type MockCantonClient struct {
	ParticipantID string
}

func (m *MockCantonClient) GenerateSigningKey(usage string) (KeyMaterial, error) {
	raw := sha256.Sum256([]byte(m.ParticipantID + ":" + usage))
	fp := fmt.Sprintf("1220%x", raw[:8])

	return KeyMaterial{
		Format:      "DER",
		KeyDataB64:  base64.StdEncoding.EncodeToString(raw[:]),
		KeySpec:     "EC_CURVE25519",
		Fingerprint: fp,
	}, nil
}

func (m *MockCantonClient) GetParticipantID() (string, error) {
	return m.ParticipantID, nil
}

func (m *MockCantonClient) GetParticipantUID() (string, error) {
	raw := sha256.Sum256([]byte("uid:" + m.ParticipantID))
	return fmt.Sprintf("%s::1220%x", m.ParticipantID, raw[:8]), nil
}

func (m *MockCantonClient) AuthorizeProposal(req AuthorizeRequest) (string, error) {
	payload := fmt.Sprintf("%s:%s:%s:%d", req.Mapping, req.SynchronizerID, req.PartyID, req.Serial)
	raw := sha256.Sum256([]byte(payload))

	return base64.StdEncoding.EncodeToString(raw[:]), nil
}

func (m *MockCantonClient) SignTransactions(req SignTransactionsRequest) (SignaturePair, error) {
	dnsSig := sha256.Sum256([]byte("dns:" + m.ParticipantID + ":" + req.DNSTxB64))
	p2pSig := sha256.Sum256([]byte("p2p:" + m.ParticipantID + ":" + req.P2PTxB64))

	return SignaturePair{
		DNSSignatureB64: base64.StdEncoding.EncodeToString(dnsSig[:]),
		P2PSignatureB64: base64.StdEncoding.EncodeToString(p2pSig[:]),
		SignedBy:        m.ParticipantID,
	}, nil
}

func (m *MockCantonClient) AddTransactions(_, _ string) error { return nil }

func (m *MockCantonClient) PollUntilConfirmed(_, _ string) error { return nil }

// NewMockCantonClientFromConfig returns a MockCantonClient.
//
// In production this would use cfg.AdminHost, cfg.AdminPort, and cfg.AdminJWT
// to dial the Canton admin gRPC endpoint and return a real client.
func NewMockCantonClientFromConfig(cfg client.ClientConfig) CantonClient {
	return &MockCantonClient{ParticipantID: cfg.ParticipantID}
}
