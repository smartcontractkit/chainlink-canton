package keyrotation

import "time"

// KeyRotationInput is the top-level input to [KeyRotationSequence].
type KeyRotationInput struct {
	// DecentralizedPartyID is the full party identifier in the format
	// "<prefix>::<namespace>", e.g. "cbtc-network::1220abcdef...".
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// TargetParticipantID is the Canton UID of the participant whose key
	// is being rotated (e.g. "PAR::name::fingerprint").
	TargetParticipantID string `json:"target_participant_id"`

	// TargetNamespaceFingerprint is the current namespace fingerprint of the
	// target participant. Used to identify which DNS owner entry to replace
	// during namespace key rotation.
	TargetNamespaceFingerprint string `json:"target_namespace_fingerprint"`

	// SynchronizerID is the Canton synchronizer to update.
	SynchronizerID string `json:"synchronizer_id"`

	// RotateNamespaceKey controls whether the namespace signing key is rotated.
	// When true, the DNS owner list is updated with the new namespace fingerprint.
	RotateNamespaceKey bool `json:"rotate_namespace_key"`

	// RotateDamlKey controls whether the DAML (protocol) signing key is rotated.
	// When true, the P2P signing keys are updated with the new DAML key.
	RotateDamlKey bool `json:"rotate_daml_key"`

	// NewThreshold overrides the threshold after the rotation.
	// When 0 or omitted, the sequence keeps the current threshold.
	NewThreshold int `json:"new_threshold,omitempty"`
}

// Phase represents the current execution phase of the key rotation ceremony.
type Phase string

const (
	PhaseReadState   Phase = "read-state"
	PhaseKeyGen      Phase = "key-gen"
	PhaseNSD         Phase = "nsd"
	PhaseDNSProposal Phase = "dns-proposal"
	PhaseDNSSigning  Phase = "dns-signing"
	PhaseDNSSubmit   Phase = "dns-submit"
	PhaseP2P         Phase = "p2p"
	PhaseCompleted   Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every KeyRotationOutput.
// It is built progressively as the sequence advances, so it is always present
// — even when the sequence returns an error (e.g. ErrThresholdNotMet).
type CeremonyState struct {
	Phase             Phase    `json:"phase"`
	TargetKeyGenReady bool     `json:"target_key_gen_ready"`
	NSDProposed       bool     `json:"nsd_proposed"`
	RequiredSigners   []string `json:"required_signers,omitempty"`
	CollectedSigners  []string `json:"collected_signers"`
	PendingSigners    []string `json:"pending_signers"`
	DNSThreshold      int      `json:"dns_threshold"`
	ProposalHash      string   `json:"proposal_hash,omitempty"`
	P2PProposedCount  int      `json:"p2p_proposed_count"`
	P2PRequired       int      `json:"p2p_required"`
	RotateNamespace   bool     `json:"rotate_namespace"`
	RotateDaml        bool     `json:"rotate_daml"`
}

// KeyRotationOutput is the final result of a completed [KeyRotationSequence].
// State is always populated — even when ExecuteSequence returns an error —
// making it the primary way to inspect ceremony progress.
type KeyRotationOutput struct {
	NamespaceKeyRotated     bool          `json:"namespace_key_rotated"`
	DamlKeyRotated          bool          `json:"daml_key_rotated"`
	NewNamespaceFingerprint string        `json:"new_namespace_fingerprint,omitempty"`
	NewDamlKeyFingerprint   string        `json:"new_daml_key_fingerprint,omitempty"`
	DNSUpdated              bool          `json:"dns_updated"`
	P2PUpdated              bool          `json:"p2p_updated"`
	State                   CeremonyState `json:"state"`
}

// ── Per-operation I/O types ──────────────────────────────────────────────────

// GenerateRotatedKeyInput is the input to [GenerateRotatedKeyOp].
type GenerateRotatedKeyInput struct {
	// ParticipantID identifies the target participant. Only this participant's
	// node can execute this operation (UID check enforced).
	ParticipantID string `json:"participant_id"`

	// SynchronizerID is the Canton synchronizer to query when discovering the
	// existing key name from the vault.
	SynchronizerID string `json:"synchronizer_id"`

	// DNSOwners is the list of current namespace fingerprints from the
	// DecentralizedNamespaceDefinition. Used to look up the existing
	// namespace key name from the vault via GetNamespaceKeyName.
	DNSOwners []string `json:"dns_owners"`

	// RotateNamespaceKey controls whether a new namespace key is generated.
	RotateNamespaceKey bool `json:"rotate_namespace_key"`

	// RotateDamlKey controls whether a new DAML key is generated.
	RotateDamlKey bool `json:"rotate_daml_key"`

	// KnownSigningKeysB64 is the list of current party signing keys
	// (base64-encoded proto). Used to discover the target's old DAML key
	// via vault cross-reference when RotateDamlKey is true.
	KnownSigningKeysB64 []string `json:"known_signing_keys_b64,omitempty"`
}

// GenerateRotatedKeyOutput is the output of [GenerateRotatedKeyOp].
type GenerateRotatedKeyOutput struct {
	ParticipantID  string `json:"participant_id"`
	ParticipantUID string `json:"participant_uid"`

	// Namespace key rotation fields (populated when RotateNamespaceKey is true).
	NewNamespaceKeyB64      string `json:"new_namespace_key_b64,omitempty"`
	NewNamespaceFingerprint string `json:"new_namespace_fingerprint,omitempty"`

	// DAML key rotation fields (populated when RotateDamlKey is true).
	NewDamlKeyB64         string `json:"new_daml_key_b64,omitempty"`
	NewDamlKeyFingerprint string `json:"new_daml_key_fingerprint,omitempty"`

	// Old DAML key info (auto-discovered from vault, populated when RotateDamlKey is true).
	OldDamlKeyFingerprint string `json:"old_daml_key_fingerprint,omitempty"`
	OldDamlKeyB64         string `json:"old_daml_key_b64,omitempty"`
}

// CreateRotationDNSProposalInput is the input to [CreateRotationDNSProposalOp].
type CreateRotationDNSProposalInput struct {
	DecentralizedNamespace  string   `json:"decentralized_namespace"`
	CurrentOwners           []string `json:"current_owners"`
	OldNamespaceFingerprint string   `json:"old_namespace_fingerprint"`
	NewNamespaceFingerprint string   `json:"new_namespace_fingerprint"`
	CurrentThreshold        int      `json:"current_threshold"`
	CurrentSerial           int      `json:"current_serial"`
	// AllParticipantIDs are the Canton UIDs of all current members.
	// All members are eligible signers for the DNS update.
	AllParticipantIDs []string `json:"all_participant_ids"`
	SynchronizerID    string   `json:"synchronizer_id"`
}

// CreateRotationDNSProposalOutput is the output of [CreateRotationDNSProposalOp].
type CreateRotationDNSProposalOutput struct {
	DNSTxB64           string   `json:"dns_tx_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	NewOwners          []string `json:"new_owners"`
	// RequiredSigners is the list of all current member UIDs. Canton requires
	// threshold-of-current-owners for serial > 1 DNS updates.
	RequiredSigners []string `json:"required_signers"`
}

// ProposeRotationP2PInput is the input to [ProposeRotationP2POp].
type ProposeRotationP2PInput struct {
	// ParticipantID uniquely identifies this proposer, giving each actor a
	// distinct idempotency hash.
	ParticipantID string `json:"participant_id"`

	// PartyID is the full party identifier (e.g. "prefix::namespace").
	PartyID string `json:"party_id"`

	// AllParticipantUIDs includes all current participant UIDs. Key rotation
	// does not change membership, only the signing keys.
	AllParticipantUIDs []string `json:"all_participant_uids"`

	// NewP2PThreshold is the P2P threshold to set.
	NewP2PThreshold int `json:"new_p2p_threshold"`

	// CurrentP2PSerial is the current P2P serial number.
	CurrentP2PSerial int `json:"current_p2p_serial"`

	SynchronizerID string `json:"synchronizer_id"`

	// CurrentSigningKeysB64 is the full list of current DAML signing keys
	// (base64-encoded proto). One entry will be replaced.
	CurrentSigningKeysB64 []string `json:"current_signing_keys_b64"`

	// OldDamlKeyB64 is the target's current DAML key to replace.
	OldDamlKeyB64 string `json:"old_daml_key_b64"`

	// NewDamlKeyB64 is the replacement DAML key.
	NewDamlKeyB64 string `json:"new_daml_key_b64"`

	// SigningKeysThreshold is the threshold for the signing keys.
	SigningKeysThreshold uint32 `json:"signing_keys_threshold"`
}

// ProposeRotationP2POutput is the output of [ProposeRotationP2POp].
type ProposeRotationP2POutput struct {
	ParticipantID string    `json:"participant_id"`
	Proposed      bool      `json:"proposed"`
	ProposedAt    time.Time `json:"proposed_at"`
}
