package onboarding

import "time"

// OnboardingInput is the top-level input to [OnboardingSequence].
type OnboardingInput struct {
	NamespaceName  string   `json:"namespace_name"`
	PartyPrefix    string   `json:"party_prefix"`
	Participants   []string `json:"participants"`
	SynchronizerID string   `json:"synchronizer_id"`
	// Threshold sets the signing/confirmation threshold for the decentralized
	// namespace and the party-to-participant mapping. When 0 the sequence
	// defaults to strict-majority: floor(n/2)+1.
	Threshold int `json:"threshold"`
}

// OnboardingOutput is the final result of a completed [OnboardingSequence].
type OnboardingOutput struct {
	PartyID      string `json:"party_id"`
	DNSConfirmed bool   `json:"dns_confirmed"`
	P2PConfirmed bool   `json:"p2p_confirmed"`
}

// ── Per-operation I/O types ──────────────────────────────────────────────────

// CreateMemberKeyInput is the input to [CreateMemberKeyOp].
type CreateMemberKeyInput struct {
	NamespaceName string `json:"namespace_name"`
	ParticipantID string `json:"participant_id"`
}

// CreateMemberKeyOutput is the output of [CreateMemberKeyOp].
type CreateMemberKeyOutput struct {
	ParticipantID        string `json:"participant_id"`
	ParticipantUID       string `json:"participant_uid"`
	NamespaceFingerprint string `json:"namespace_fingerprint"`
	// SigningKeyB64 is the base64-encoded proto-marshalled SigningPublicKey
	// returned by GenerateSigningKey. It preserves the exact Format, Scheme,
	// KeySpec, and Usage fields that Canton assigned, avoiding reconstruction errors.
	SigningKeyB64 string `json:"signing_key_b64"`
}

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

// CreateDNSProposalInput is the input to [CreateDNSProposalOp].
type CreateDNSProposalInput struct {
	NamespaceName  string                  `json:"namespace_name"`
	Members        []CreateMemberKeyOutput `json:"members"`
	SynchronizerID string                  `json:"synchronizer_id"`
	Threshold      int                     `json:"threshold"`
}

// CreateDNSProposalOutput is the output of [CreateDNSProposalOp].
type CreateDNSProposalOutput struct {
	DecentralizedNS    string   `json:"decentralized_namespace"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	DNSTxB64           string   `json:"dns_transaction_b64"`
	RequiredSigners    []string `json:"required_signers"`
	Threshold          int      `json:"threshold"`
}

// SignDNSProposalInput is the input to [SignDNSProposalOp].
type SignDNSProposalInput struct {
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

// SubmitDNSInput is the input to [SubmitDNSOp].
type SubmitDNSInput struct {
	// SignedDNSTxsB64 holds one base64-encoded proto-marshalled
	// SignedTopologyTransaction per signer. SubmitDNSOp merges the Signature
	// lists from all entries before calling AddTransactions, so each signer can
	// independently add only their own signature to the original proposal.
	SignedDNSTxsB64 []string `json:"signed_dns_txs_b64"`
	SynchronizerID  string   `json:"synchronizer_id"`
	FilterNamespace string   `json:"filter_namespace"`
}

// SubmitDNSOutput is the output of [SubmitDNSOp].
type SubmitDNSOutput struct {
	DNSSubmitted bool `json:"dns_submitted"`
}

// ProposeP2PInput is the input to [ProposeP2POp].
// ParticipantID is included so the idempotency hash is unique per participant,
// mirroring the SignDNSProposalInput pattern.
type ProposeP2PInput struct {
	ParticipantID  string                  `json:"participant_id"`
	PartyID        string                  `json:"party_id"`
	Members        []CreateMemberKeyOutput `json:"members"`
	SynchronizerID string                  `json:"synchronizer_id"`
	Threshold      int                     `json:"threshold"`
}

// ProposeP2POutput is the output of [ProposeP2POp].
type ProposeP2POutput struct {
	ParticipantID string `json:"participant_id"`
	Proposed      bool   `json:"proposed"`
}
