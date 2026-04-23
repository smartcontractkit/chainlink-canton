package example

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// ── Dependency types ─────────────────────────────────────────────────────────

// KeyMaterial holds a single Canton signing key entry.
// All fields are JSON-serializable strings so they can be persisted in
// member-<id>.json.
type KeyMaterial struct {
	Format      string `json:"format"`
	KeyDataB64  string `json:"key_data_b64"`
	KeySpec     string `json:"key_spec"`
	Fingerprint string `json:"fingerprint"`
}

// AuthorizeRequest mirrors the fields we pass to
// TopologyManagerWriteService.Authorize.
type AuthorizeRequest struct {
	Mapping        string   `json:"mapping"`         // e.g. "DecentralizedNamespaceDefinition"
	SynchronizerID string   `json:"synchronizer_id"` // empty → Authorized store
	Owners         []string `json:"owners,omitempty"`
	PartyID        string   `json:"party_id,omitempty"`
	Serial         int      `json:"serial"`
}

// SignTransactionsRequest mirrors SignTransactionsRequest (topology write API).
type SignTransactionsRequest struct {
	DNSTxB64       string `json:"dns_tx_b64"`
	P2PTxB64       string `json:"p2p_tx_b64"`
	SynchronizerID string `json:"synchronizer_id"`
}

// SignaturePair holds the two signatures (DNS and P2P) produced by one signer.
type SignaturePair struct {
	DNSSignatureB64 string `json:"dns_sig_b64"`
	P2PSignatureB64 string `json:"p2p_sig_b64"`
	SignedBy        string `json:"signed_by"`
}

// CantonClient is the interface that each Canton admin operation requires.
// In production this would wrap gRPC stubs for VaultService,
// TopologyManagerWriteService and TopologyManagerReadService.
type CantonClient interface {
	// GenerateSigningKey produces key material for a participant.
	GenerateSigningKey(usage string) (KeyMaterial, error)

	// GetParticipantID returns the participant's ID
	GetParticipantID() (string, error)

	// GetParticipantUID returns the participant's unique identifier (UID) as
	// reported by IdentityInitializationService.GetId.
	GetParticipantUID() (string, error)

	// AuthorizeProposal calls TopologyManagerWriteService.Authorize and returns
	// the raw base-64-encoded signed topology transaction bytes.
	AuthorizeProposal(req AuthorizeRequest) (string, error)

	// SignTransactions calls TopologyManagerWriteService.SignTransactions and
	// returns a pair of signatures (one for DNS, one for P2P).
	SignTransactions(req SignTransactionsRequest) (SignaturePair, error)

	// AddTransactions calls TopologyManagerWriteService.AddTransactions.
	AddTransactions(txB64, synchronizerID string) error

	// PollUntilConfirmed polls the read service until the topology change is
	// visible at head state.  This is the "wait for DNS / P2P" loop that wraps
	// ListDecentralizedNamespaceDefinition / ListPartyToParticipant.
	PollUntilConfirmed(filter, synchronizerID string) error
}

// CantonDeps is the dependency container passed to every operation handler.
type CantonDeps struct {
	Client    CantonClient
	Logger    logger.Logger
	Confirmer ceremony.Confirmer // nil means no confirmation prompt
}
