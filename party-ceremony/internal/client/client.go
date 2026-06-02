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

	// RegisterKmsSigningKey registers a pre-existing KMS signing key in the
	// participant's vault. kmsKeyID is the external KMS identifier (e.g. an
	// AWS KMS ARN); name and usage mirror [GenerateSigningKey].
	// Returns the public key proto for the registered key.
	RegisterKmsSigningKey(ctx context.Context, kmsKeyID string, name string, usage []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error)

	// GetNamespaceFingerprint returns the namespace fingerprint for a named
	// key by cross-referencing the vault (ListMyKeys) with namespace
	// delegations (ListNamespaceDelegation).
	// knownOwners optionally restricts the search to a set of candidate
	// namespace fingerprints (e.g. the current party's DNS owners), which
	// avoids returning stale fingerprints from previous runs when the vault
	// holds multiple keys with the same name.
	GetNamespaceFingerprint(ctx context.Context, keyName string, synchronizerID string, knownOwners []string) (string, error)

	// GetNamespaceKeyName returns the human-readable name of this participant's
	// NAMESPACE signing key that belongs to the decentralized namespace identified
	// by knownOwners. It cross-references the vault (ListMyKeys with NAMESPACE
	// usage) with namespace delegations (ListNamespaceDelegation) to find the
	// matching key's label. Used by key rotation to auto-discover the original
	// onboarding key name without requiring user input.
	GetNamespaceKeyName(ctx context.Context, synchronizerID string, knownOwners []string) (string, error)

	// GetProtocolKeyFingerprint returns the fingerprint and base64-encoded
	// proto of this participant's PROTOCOL signing key that is currently
	// active in the party's signing keys. It cross-references vault keys
	// (ListMyKeys with PROTOCOL usage) against the provided knownSigningKeys
	// (base64-encoded proto-marshalled keys from P2P PartySigningKeys).
	GetProtocolKeyFingerprint(ctx context.Context, knownSigningKeys []string) (fingerprint string, keyB64 string, err error)

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

	// ── ACS Replication (ParticipantRepairService) ───────────────────────

	// ExportAcs exports the Active Contract Set for the given parties and
	// synchronizer at the specified ledger offset. The returned bytes are the
	// concatenated ACS snapshot chunks.
	ExportAcs(ctx context.Context, partyIDs []string, synchronizerID string, ledgerOffset int64) ([]byte, error)

	// ImportAcs imports an ACS snapshot into the participant for the given
	// synchronizer. The participant must be disconnected from the synchronizer
	// before calling this.
	ImportAcs(ctx context.Context, acsSnapshot []byte, synchronizerID string) error

	// ── Synchronizer Connectivity ────────────────────────────────────────

	// DisconnectSynchronizer disconnects the participant from the given
	// synchronizer alias.
	DisconnectSynchronizer(ctx context.Context, synchronizerAlias string) error

	// ReconnectSynchronizer reconnects the participant to the given
	// synchronizer alias.
	ReconnectSynchronizer(ctx context.Context, synchronizerAlias string) error

	// ListConnectedSynchronizers returns the aliases and IDs of all currently
	// connected synchronizers.
	ListConnectedSynchronizers(ctx context.Context) ([]SynchronizerInfo, error)

	// ── Party Management ─────────────────────────────────────────────────

	// ClearPartyOnboardingFlag clears the onboarding flag for a party on the
	// given synchronizer. Returns true if the flag was successfully cleared.
	ClearPartyOnboardingFlag(ctx context.Context, partyID string, synchronizerID string, beginOffsetExclusive int64) (bool, error)

	// ── Inspection ───────────────────────────────────────────────────────

	// LookupOffsetByTime returns the ledger offset at the given timestamp.
	// Used to record a reference point before topology changes.
	LookupOffsetByTime(ctx context.Context, timestamp int64) (int64, error)
}
