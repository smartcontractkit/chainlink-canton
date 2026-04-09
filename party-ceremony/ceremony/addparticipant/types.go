package addparticipant

import "time"

// AddParticipantInput is the top-level input to [AddParticipantSequence].
type AddParticipantInput struct {
	// DecentralizedPartyID is the full party identifier in the format
	// "<prefix>::<namespace>", e.g. "cbtc-network::1220abcdef...".
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// NewParticipantID is the Canton UID of the participant being added
	// (e.g. "PAR::newnode::fingerprint").
	NewParticipantID string `json:"new_participant_id"`

	// ExistingParticipants is the ordered list of Canton UIDs for all
	// participants that are already members of the decentralized party.
	// Only these actors sign the DNS update and propose P2P.
	ExistingParticipants []string `json:"existing_participants"`

	// NamespaceName is the human-readable label used when generating the
	// new participant's namespace and DAML signing keys.
	NamespaceName string `json:"namespace_name"`

	// SynchronizerID is the Canton synchronizer to update.
	SynchronizerID string `json:"synchronizer_id"`

	// NewThreshold overrides the threshold after the addition.
	// When 0 or omitted, the sequence keeps the current threshold.
	NewThreshold int `json:"new_threshold,omitempty"`
}

// AddParticipantOutput is the final result of a completed [AddParticipantSequence].
type AddParticipantOutput struct {
	DNSUpdated   bool     `json:"dns_updated"`
	P2PUpdated   bool     `json:"p2p_updated"`
	NewThreshold int      `json:"new_threshold"`
	AllOwners    []string `json:"all_owners"`
}

// ── Per-operation I/O types ──────────────────────────────────────────────────

// GenerateNewMemberKeyInput is the input to [GenerateNewMemberKeyOp].
type GenerateNewMemberKeyInput struct {
	NamespaceName string `json:"namespace_name"`
	ParticipantID string `json:"participant_id"`
}

// GenerateNewMemberKeyOutput is the output of [GenerateNewMemberKeyOp].
type GenerateNewMemberKeyOutput struct {
	ParticipantID        string `json:"participant_id"`
	ParticipantUID       string `json:"participant_uid"`
	NamespaceFingerprint string `json:"namespace_fingerprint"`
	// SigningKeyB64 is the base64-encoded proto-marshalled SigningPublicKey
	// returned by GenerateSigningKey (NamespaceOnly usage).
	SigningKeyB64 string `json:"signing_key_b64"`
	// DamlKeyB64 is the base64-encoded proto-marshalled PROTOCOL SigningPublicKey.
	DamlKeyB64         string `json:"daml_key_b64"`
	DamlKeyFingerprint string `json:"daml_key_fingerprint"`
}

// ProposeNewNSDInput is the input to [ProposeNewNSDOp].
type ProposeNewNSDInput struct {
	ParticipantID  string `json:"participant_id"`
	SigningKeyB64  string `json:"signing_key_b64"`
	Namespace      string `json:"namespace"`
	SynchronizerID string `json:"synchronizer_id"`
}

// ProposeNewNSDOutput is the output of [ProposeNewNSDOp].
type ProposeNewNSDOutput struct {
	ParticipantID      string `json:"participant_id"`
	DelegationProposed bool   `json:"delegation_proposed"`
}

// ReadCurrentStateInput is the input to [ReadCurrentStateOp].
type ReadCurrentStateInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
	SynchronizerID       string `json:"synchronizer_id"`
}

// ReadCurrentStateOutput is the output of [ReadCurrentStateOp].
type ReadCurrentStateOutput struct {
	// DNS state
	DecentralizedNamespace string   `json:"decentralized_namespace"`
	DNSOwners              []string `json:"dns_owners"`
	DNSThreshold           int32    `json:"dns_threshold"`
	DNSSerial              int32    `json:"dns_serial"`
	// P2P state
	P2PParticipantUIDs []string `json:"p2p_participant_uids"`
	P2PThreshold       uint32   `json:"p2p_threshold"`
	P2PSerial          int32    `json:"p2p_serial"`
}

// CreateAddDNSProposalInput is the input to [CreateAddDNSProposalOp].
type CreateAddDNSProposalInput struct {
	DecentralizedNamespace string   `json:"decentralized_namespace"`
	CurrentOwners          []string `json:"current_owners"`
	NewOwnerFingerprint    string   `json:"new_owner_fingerprint"`
	NewThreshold           int      `json:"new_threshold"`
	CurrentSerial          int      `json:"current_serial"`
	// ExistingParticipants are the Canton UIDs that are already members.
	// Only existing members sign the DNS update (the new participant is not
	// yet an owner and cannot sign).
	ExistingParticipants []string `json:"existing_participants"`
	SynchronizerID       string   `json:"synchronizer_id"`
}

// CreateAddDNSProposalOutput is the output of [CreateAddDNSProposalOp].
type CreateAddDNSProposalOutput struct {
	DNSTxB64           string   `json:"dns_tx_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	NewOwners          []string `json:"new_owners"`
	NewThreshold       int      `json:"new_threshold"`
	// RequiredSigners is the existing participant UIDs. Only current DNS
	// owners can sign a serial > 1 DNS update.
	RequiredSigners []string `json:"required_signers"`
}

// SignAddDNSProposalInput is the input to [SignAddDNSProposalOp].
type SignAddDNSProposalInput struct {
	// ParticipantID uniquely identifies this signer, giving each actor a
	// distinct idempotency hash.
	ParticipantID      string `json:"participant_id"`
	ProposalHashSHA256 string `json:"proposal_hash_sha256"`
	DNSTxB64           string `json:"dns_tx_b64"`
	SynchronizerID     string `json:"synchronizer_id"`
}

// SignAddDNSProposalOutput is the output of [SignAddDNSProposalOp].
type SignAddDNSProposalOutput struct {
	ParticipantID  string    `json:"participant_id"`
	SignedDNSTxB64 string    `json:"signed_dns_tx_b64"`
	SignedBy       string    `json:"signed_by"`
	SignedAt       time.Time `json:"signed_at"`
}

// SubmitAddDNSInput is the input to [SubmitAddDNSOp].
type SubmitAddDNSInput struct {
	SignedDNSTxsB64 []string `json:"signed_dns_txs_b64"`
	SynchronizerID  string   `json:"synchronizer_id"`
	FilterNamespace string   `json:"filter_namespace"`
}

// SubmitAddDNSOutput is the output of [SubmitAddDNSOp].
type SubmitAddDNSOutput struct {
	DNSSubmitted bool `json:"dns_submitted"`
}

// ProposeAddP2PInput is the input to [ProposeAddP2POp].
type ProposeAddP2PInput struct {
	// ParticipantID uniquely identifies this proposer, giving each existing
	// actor a distinct idempotency hash.
	ParticipantID string `json:"participant_id"`
	PartyID       string `json:"party_id"`
	// AllParticipantUIDs includes both existing and new participant UIDs.
	// All receive Confirmation permission.
	AllParticipantUIDs []string `json:"all_participant_uids"`
	// DamlKeys holds the base64-encoded DAML signing keys for ALL participants
	// (existing + new). These are registered in PartyToParticipant.SigningKeysWithThreshold.
	DamlKeys         []string `json:"daml_keys"`
	NewP2PThreshold  int      `json:"new_p2p_threshold"`
	CurrentP2PSerial int      `json:"current_p2p_serial"`
	SynchronizerID   string   `json:"synchronizer_id"`
}

// ProposeAddP2POutput is the output of [ProposeAddP2POp].
type ProposeAddP2POutput struct {
	ParticipantID string `json:"participant_id"`
	Proposed      bool   `json:"proposed"`
}
