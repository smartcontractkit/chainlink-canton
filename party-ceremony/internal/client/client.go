package client

import (
	"context"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
)

// CantonClient abstracts the Canton Admin gRPC APIs needed by the onboarding
// operations. In production this wraps real gRPC stubs; for testing it can be
// backed by a mock.
type CantonClient interface {
	// ── Identity ─────────────────────────────────────────────────────────

	// GetParticipantUID returns the participant's unique identifier as
	// reported by IdentityInitializationService.GetId (e.g. "PAR::name::fingerprint").
	GetParticipantUID(ctx context.Context) (string, error)

	// GetParticipantID returns the human-readable participant identifier
	// (the segment before "::" in the UID, e.g. "participant1").
	GetParticipantID(ctx context.Context) (string, error)

	// ── Key Management (VaultService) ────────────────────────────────────

	// GenerateSigningKey creates a new signing key in the participant's vault.
	// name is the human-readable label; usage specifies the key purpose.
	// Returns the generated public key proto.
	GenerateSigningKey(ctx context.Context, name string, usage []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error)

	// GetNamespaceFingerprint returns the namespace fingerprint for a named
	// key by cross-referencing the vault (ListMyKeys) with namespace
	// delegations (ListNamespaceDelegation).
	// knownOwners optionally restricts the search to a set of candidate
	// namespace fingerprints (e.g. the current party's DNS owners), which
	// avoids returning stale fingerprints from previous runs when the vault
	// holds multiple keys with the same name.
	GetNamespaceFingerprint(ctx context.Context, keyName string, synchronizerID string, knownOwners []string) (string, error)

	// ── Topology Write (TopologyManagerWriteService) ─────────────────────

	// Authorize calls TopologyManagerWriteService.Authorize to create a
	// topology proposal. When mustFullyAuthorize is false, Canton adds the
	// proposer's signature and stores the transaction as a proposal.
	// signedBy lists the key fingerprints Canton should sign with; nil/empty
	// lets Canton auto-select based on existing topology delegations.
	// Returns the signed topology transaction proto.
	Authorize(ctx context.Context, serial uint32, mapping *protov30.TopologyMapping, synchronizerID string, mustFullyAuthorize bool, signedBy ...string) (*protov30.SignedTopologyTransaction, error)

	// SignTransactions signs topology transactions by calling
	// TopologyManagerWriteService.SignTransactions. Canton auto-selects the
	// appropriate namespace keys when signedBy is empty.
	SignTransactions(ctx context.Context, txs []*protov30.SignedTopologyTransaction, synchronizerID string) ([]*protov30.SignedTopologyTransaction, error)

	// AddTransactions submits fully-signed topology transactions via
	// TopologyManagerWriteService.AddTransactions.
	AddTransactions(ctx context.Context, txs []*protov30.SignedTopologyTransaction, synchronizerID string) error

	// ── Topology Read (TopologyManagerReadService) ───────────────────────

	// DNSExists returns true if a DecentralizedNamespaceDefinition for the
	// given namespace is active (ADD_REPLACE) in the specified store.
	// Pass a synchronizerID to query the synchronizer store, or "" for the
	// authorized store.
	DNSExists(ctx context.Context, namespace string, synchronizerID string) (bool, error)

	// P2PExists returns true if a PartyToParticipant mapping for the given
	// party UID is active (ADD_REPLACE) in the specified store.
	// Pass a synchronizerID to query the synchronizer store, or "" for the
	// authorized store.
	P2PExists(ctx context.Context, partyUID string, synchronizerID string) (bool, error)

	// NSDExists returns true if a NamespaceDelegation for the given namespace
	// is active (ADD_REPLACE) in the specified store. Pass a synchronizerID to
	// query the synchronizer store, or "" for the authorized store.
	NSDExists(ctx context.Context, namespace string, synchronizerID string) (bool, error)

	// GetDNS returns the current DecentralizedNamespaceDefinition state for the
	// given namespace from the specified store. Returns an error if the mapping
	// does not exist or is not in ADD_REPLACE state.
	GetDNS(ctx context.Context, namespace string, synchronizerID string) (*DNSState, error)

	// GetP2P returns the current PartyToParticipant state for the given partyUID
	// from the specified store. Returns an error if the mapping does not exist
	// or is not in ADD_REPLACE state.
	GetP2P(ctx context.Context, partyUID string, synchronizerID string) (*P2PState, error)

	// ListDecentralizedNamespaces returns all active DecentralizedNamespaceDefinition
	// mappings visible in the specified store. Used by the query-parties command.
	ListDecentralizedNamespaces(ctx context.Context, synchronizerID string) ([]*DNSState, error)

	// ── Package Management (PackageService) ─────────────────────────────

	// UploadDar uploads a DAR file to the participant via PackageService.UploadDar.
	// The DAR is vetted and vetting is synchronised before returning.
	// Returns the main package ID from the uploaded DAR.
	UploadDar(ctx context.Context, darBytes []byte) (string, error)
}
