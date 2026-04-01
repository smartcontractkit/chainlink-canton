package kick

import "time"

// KickInput is the top-level input to [KickSequence].
type KickInput struct {
	// DecentralizedPartyID is the full party identifier in the format
	// "<prefix>::<namespace>", e.g. "cbtc-network::1220abcdef...".
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// KickedParticipantID is the Canton UID of the participant being removed
	// from the PartyToParticipant mapping (e.g. "PAR::name::fingerprint").
	KickedParticipantID string `json:"kicked_participant_id"`

	// KickedNamespaceFingerprint is the namespace fingerprint of the participant
	// being kicked. Used to remove their owner entry from the
	// DecentralizedNamespaceDefinition.
	KickedNamespaceFingerprint string `json:"kicked_namespace_fingerprint"`

	// RemainingParticipants is the ordered list of Canton UIDs for all
	// participants that will remain after the kick. Only these actors
	// participate in the ceremony.
	RemainingParticipants []string `json:"remaining_participants"`

	// SynchronizerID is the Canton synchronizer to update.
	SynchronizerID string `json:"synchronizer_id"`

	// NewThreshold overrides the automatic threshold after the kick.
	// When 0 or omitted, the sequence defaults to previous threshold.
	NewThreshold int `json:"new_threshold,omitempty"`
}

// KickOutput is the final result of a completed [KickSequence].
type KickOutput struct {
	DNSUpdated      bool     `json:"dns_updated"`
	P2PUpdated      bool     `json:"p2p_updated"`
	NewThreshold    int      `json:"new_threshold"`
	RemainingOwners []string `json:"remaining_owners"`
}

// ── Per-operation I/O types ──────────────────────────────────────────────────

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

// CreateKickDNSProposalInput is the input to [CreateKickDNSProposalOp].
type CreateKickDNSProposalInput struct {
	DecentralizedNamespace     string   `json:"decentralized_namespace"`
	CurrentOwners              []string `json:"current_owners"`
	KickedNamespaceFingerprint string   `json:"kicked_namespace_fingerprint"`
	NewThreshold               int      `json:"new_threshold"`
	CurrentSerial              int      `json:"current_serial"`
	// RemainingParticipants are the Canton UIDs that stay in the party after the kick.
	RemainingParticipants []string `json:"remaining_participants"`
	// KickedParticipantID is included so it can be appended to RequiredSigners:
	// Canton requires threshold-of-current-owners for serial > 1 DNS updates,
	// so the kicked participant must also sign (they are still a current owner).
	KickedParticipantID string `json:"kicked_participant_id"`
	SynchronizerID      string `json:"synchronizer_id"`
}

// CreateKickDNSProposalOutput is the output of [CreateKickDNSProposalOp].
type CreateKickDNSProposalOutput struct {
	DNSTxB64           string   `json:"dns_tx_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	NewOwners          []string `json:"new_owners"`
	NewThreshold       int      `json:"new_threshold"`
	// RequiredSigners is remaining + kicked participant UIDs. Canton needs
	// threshold-of-current-owners for serial > 1 DNS updates, so the kicked
	// participant must contribute their signature too.
	RequiredSigners []string `json:"required_signers"`
}

// SignKickDNSProposalInput is the input to [SignKickDNSProposalOp].
type SignKickDNSProposalInput struct {
	// ParticipantID uniquely identifies this signer, giving each actor a
	// distinct idempotency hash so the framework cache does not block
	// subsequent actors from executing their own signing step.
	ParticipantID      string `json:"participant_id"`
	ProposalHashSHA256 string `json:"proposal_hash_sha256"`
	DNSTxB64           string `json:"dns_tx_b64"`
	SynchronizerID     string `json:"synchronizer_id"`
}

// SignKickDNSProposalOutput is the output of [SignKickDNSProposalOp].
type SignKickDNSProposalOutput struct {
	ParticipantID  string    `json:"participant_id"`
	SignedDNSTxB64 string    `json:"signed_dns_tx_b64"`
	SignedBy       string    `json:"signed_by"`
	SignedAt       time.Time `json:"signed_at"`
}

// SubmitKickDNSInput is the input to [SubmitKickDNSOp].
type SubmitKickDNSInput struct {
	// SignedDNSTxsB64 holds one base64-encoded proto-marshalled
	// SignedTopologyTransaction per signer. [SubmitKickDNSOp] merges the
	// Signature lists before calling AddTransactions.
	SignedDNSTxsB64 []string `json:"signed_dns_txs_b64"`
	SynchronizerID  string   `json:"synchronizer_id"`
	FilterNamespace string   `json:"filter_namespace"`
}

// SubmitKickDNSOutput is the output of [SubmitKickDNSOp].
type SubmitKickDNSOutput struct {
	DNSSubmitted bool `json:"dns_submitted"`
}

// ProposeKickP2PInput is the input to [ProposeKickP2POp].
type ProposeKickP2PInput struct {
	// ParticipantID uniquely identifies this proposer, giving each remaining
	// actor a distinct idempotency hash.
	ParticipantID string `json:"participant_id"`
	PartyID       string `json:"party_id"`
	// RemainingParticipants are the Canton UIDs of participants that will remain
	// after the kick. All receive Confirmation permission.
	RemainingParticipants []string `json:"remaining_participants"`
	NewP2PThreshold       int      `json:"new_p2p_threshold"`
	CurrentP2PSerial      int      `json:"current_p2p_serial"`
	SynchronizerID        string   `json:"synchronizer_id"`
}

// ProposeKickP2POutput is the output of [ProposeKickP2POp].
type ProposeKickP2POutput struct {
	ParticipantID string `json:"participant_id"`
	Proposed      bool   `json:"proposed"`
}
