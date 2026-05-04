package topology

import (
	"time"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
)

// ── ProposeNamespaceDelegation (NSD) ─────────────────────────────────────────

// ProposeNSDInput is the input to [ProposeNamespaceDelegationOp].
type ProposeNSDInput struct {
	ParticipantID  string `json:"participant_id"`
	SigningKeyB64  string `json:"signing_key_b64"`
	Namespace      string `json:"namespace"`
	SynchronizerID string `json:"synchronizer_id"`
}

// ProposeNSDOutput is the output of [ProposeNamespaceDelegationOp].
type ProposeNSDOutput struct {
	ParticipantID      string `json:"participant_id"`
	DelegationProposed bool   `json:"delegation_proposed"`
}

// ── ReadCurrentState ─────────────────────────────────────────────────────────

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
	// Party signing keys (optional, populated when available in topology)
	PartySigningKeysB64   []string `json:"party_signing_keys_b64,omitempty"`
	PartySigningThreshold uint32   `json:"party_signing_threshold,omitempty"`
}

// ── MemberKeyOutput ──────────────────────────────────────────────────────────

// MemberKeyOutput is an alias for [keys.CreateMemberKeyOutput].
// Topology operations accept member lists in this form, allowing callers
// to pass the output of CreateMemberKeyOp directly without any conversion.
type MemberKeyOutput = keys.CreateMemberKeyOutput

// ── CreateDNSProposal (onboarding – initial DNS creation) ────────────────────

// CreateDNSProposalInput is the input to [CreateDNSProposalOp].
type CreateDNSProposalInput struct {
	NamespaceName  string                       `json:"namespace_name"`
	Members        []keys.CreateMemberKeyOutput `json:"members"`
	SynchronizerID string                       `json:"synchronizer_id"`
	Threshold      int                          `json:"threshold"`
}

// CreateDNSProposalOutput is the output of [CreateDNSProposalOp].
type CreateDNSProposalOutput struct {
	DecentralizedNS    string   `json:"decentralized_namespace"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	DNSTxB64           string   `json:"dns_transaction_b64"`
	RequiredSigners    []string `json:"required_signers"`
	Threshold          int      `json:"threshold"`
}

// ── CreateKickDNSProposal ────────────────────────────────────────────────────

// CreateKickDNSProposalInput is the input to [CreateKickDNSProposalOp].
type CreateKickDNSProposalInput struct {
	DecentralizedNamespace     string   `json:"decentralized_namespace"`
	CurrentOwners              []string `json:"current_owners"`
	KickedNamespaceFingerprint string   `json:"kicked_namespace_fingerprint"`
	NewThreshold               int      `json:"new_threshold"`
	CurrentSerial              int      `json:"current_serial"`
	RemainingParticipants      []string `json:"remaining_participants"`
	KickedParticipantID        string   `json:"kicked_participant_id"`
	SynchronizerID             string   `json:"synchronizer_id"`
}

// CreateKickDNSProposalOutput is the output of [CreateKickDNSProposalOp].
type CreateKickDNSProposalOutput struct {
	DNSTxB64           string   `json:"dns_tx_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	NewOwners          []string `json:"new_owners"`
	NewThreshold       int      `json:"new_threshold"`
	RequiredSigners    []string `json:"required_signers"`
}

// ── CreateAddDNSProposal ─────────────────────────────────────────────────────

// CreateAddDNSProposalInput is the input to [CreateAddDNSProposalOp].
type CreateAddDNSProposalInput struct {
	DecentralizedNamespace  string   `json:"decentralized_namespace"`
	CurrentOwners           []string `json:"current_owners"`
	NewOwnerFingerprint     string   `json:"new_owner_fingerprint"`
	NewThreshold            int      `json:"new_threshold"`
	CurrentSerial           int      `json:"current_serial"`
	ExistingParticipantUIDs []string `json:"existing_participant_uids"`
	SynchronizerID          string   `json:"synchronizer_id"`
}

// CreateAddDNSProposalOutput is the output of [CreateAddDNSProposalOp].
type CreateAddDNSProposalOutput struct {
	DNSTxB64           string   `json:"dns_tx_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	NewOwners          []string `json:"new_owners"`
	NewThreshold       int      `json:"new_threshold"`
	RequiredSigners    []string `json:"required_signers"`
}

// ── CreateRotationDNSProposal ────────────────────────────────────────────────

// CreateRotationDNSProposalInput is the input to [CreateRotationDNSProposalOp].
type CreateRotationDNSProposalInput struct {
	DecentralizedNamespace  string   `json:"decentralized_namespace"`
	CurrentOwners           []string `json:"current_owners"`
	OldNamespaceFingerprint string   `json:"old_namespace_fingerprint"`
	NewNamespaceFingerprint string   `json:"new_namespace_fingerprint"`
	CurrentThreshold        int      `json:"current_threshold"`
	CurrentSerial           int      `json:"current_serial"`
	AllParticipantIDs       []string `json:"all_participant_ids"`
	SynchronizerID          string   `json:"synchronizer_id"`
}

// CreateRotationDNSProposalOutput is the output of [CreateRotationDNSProposalOp].
type CreateRotationDNSProposalOutput struct {
	DNSTxB64           string   `json:"dns_tx_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	NewOwners          []string `json:"new_owners"`
	RequiredSigners    []string `json:"required_signers"`
}

// ── SignDNSProposal ──────────────────────────────────────────────────────────

// SignDNSProposalInput is the input to [SignDNSProposalOp].
type SignDNSProposalInput struct {
	// ParticipantID uniquely identifies this signer, giving each actor a
	// distinct idempotency hash so the framework cache does not block
	// subsequent actors from executing their own signing step.
	ParticipantID      string `json:"participant_id"`
	ProposalHashSHA256 string `json:"proposal_hash_sha256"`
	DNSTxB64           string `json:"dns_tx_b64"`
	SynchronizerID     string `json:"synchronizer_id"`
}

// SignDNSProposalOutput is the output of [SignDNSProposalOp].
type SignDNSProposalOutput struct {
	ParticipantID  string    `json:"participant_id"`
	SignedDNSTxB64 string    `json:"signed_dns_tx_b64"`
	SignedBy       string    `json:"signed_by"`
	SignedAt       time.Time `json:"signed_at"`
}

// ── SubmitDNS ────────────────────────────────────────────────────────────────

// SubmitDNSInput is the input to [SubmitDNSOp].
type SubmitDNSInput struct {
	// SignedDNSTxsB64 holds one base64-encoded proto-marshalled
	// SignedTopologyTransaction per signer. [SubmitDNSOp] merges the
	// Signature lists before calling AddTransactions.
	SignedDNSTxsB64 []string `json:"signed_dns_txs_b64"`
	SynchronizerID  string   `json:"synchronizer_id"`
	FilterNamespace string   `json:"filter_namespace"`
}

// SubmitDNSOutput is the output of [SubmitDNSOp].
type SubmitDNSOutput struct {
	DNSSubmitted bool `json:"dns_submitted"`
}

// ── ProposeP2P (onboarding – initial P2P creation) ───────────────────────────

// ProposeP2PInput is the input to [ProposeP2POp].
type ProposeP2PInput struct {
	ParticipantID  string                       `json:"participant_id"`
	PartyID        string                       `json:"party_id"`
	Members        []keys.CreateMemberKeyOutput `json:"members"`
	SynchronizerID string                       `json:"synchronizer_id"`
	Threshold      int                          `json:"threshold"`
}

// ProposeP2POutput is the output of [ProposeP2POp].
type ProposeP2POutput struct {
	ParticipantID string `json:"participant_id"`
	Proposed      bool   `json:"proposed"`
}

// ── ProposeKickP2P ───────────────────────────────────────────────────────────

// ProposeKickP2PInput is the input to [ProposeKickP2POp].
type ProposeKickP2PInput struct {
	ParticipantID         string   `json:"participant_id"`
	PartyID               string   `json:"party_id"`
	RemainingParticipants []string `json:"remaining_participants"`
	NewP2PThreshold       int      `json:"new_p2p_threshold"`
	CurrentP2PSerial      int      `json:"current_p2p_serial"`
	SynchronizerID        string   `json:"synchronizer_id"`
	PartySigningKeysB64   []string `json:"party_signing_keys_b64,omitempty"`
}

// ProposeKickP2POutput is the output of [ProposeKickP2POp].
type ProposeKickP2POutput struct {
	ParticipantID string `json:"participant_id"`
	Proposed      bool   `json:"proposed"`
}

// ── ProposeAddP2P ────────────────────────────────────────────────────────────

// ProposeAddP2PInput is the input to [ProposeAddP2POp].
type ProposeAddP2PInput struct {
	ParticipantID       string   `json:"participant_id"`
	PartyID             string   `json:"party_id"`
	AllParticipantUIDs  []string `json:"all_participant_uids"`
	NewP2PThreshold     int      `json:"new_p2p_threshold"`
	CurrentP2PSerial    int      `json:"current_p2p_serial"`
	SynchronizerID      string   `json:"synchronizer_id"`
	PartySigningKeysB64 []string `json:"party_signing_keys_b64,omitempty"`
}

// ProposeAddP2POutput is the output of [ProposeAddP2POp].
type ProposeAddP2POutput struct {
	ParticipantID string `json:"participant_id"`
	Proposed      bool   `json:"proposed"`
}

// ── ProposeRotationP2P ───────────────────────────────────────────────────────

// ProposeRotationP2PInput is the input to [ProposeRotationP2POp].
type ProposeRotationP2PInput struct {
	ParticipantID         string   `json:"participant_id"`
	PartyID               string   `json:"party_id"`
	AllParticipantUIDs    []string `json:"all_participant_uids"`
	NewP2PThreshold       int      `json:"new_p2p_threshold"`
	CurrentP2PSerial      int      `json:"current_p2p_serial"`
	SynchronizerID        string   `json:"synchronizer_id"`
	CurrentSigningKeysB64 []string `json:"current_signing_keys_b64"`
	OldDamlKeyB64         string   `json:"old_daml_key_b64"`
	NewDamlKeyB64         string   `json:"new_daml_key_b64"`
	// PartySigningKeysThreshold preserves the current topology's signing-key
	// threshold while replacing only the rotated protocol key.
	PartySigningKeysThreshold uint32 `json:"party_signing_keys_threshold"`
}

// ProposeRotationP2POutput is the output of [ProposeRotationP2POp].
type ProposeRotationP2POutput struct {
	ParticipantID string    `json:"participant_id"`
	Proposed      bool      `json:"proposed"`
	ProposedAt    time.Time `json:"proposed_at"`
}
