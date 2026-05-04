package ledger

// PackageRef identifies a DAML package by its registered name and version.
type PackageRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ── FetchParticipants ────────────────────────────────────────────────────────

// FetchParticipantsInput is the input to [FetchParticipantsOp].
type FetchParticipantsInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
	SynchronizerID       string `json:"synchronizer_id"`
}

// FetchParticipantsOutput is the output of [FetchParticipantsOp].
type FetchParticipantsOutput struct {
	Participants          []string `json:"participants"`
	PartySigningKeysB64   []string `json:"party_signing_keys_b64,omitempty"`
	PartySigningThreshold uint32   `json:"party_signing_threshold,omitempty"`
}

// ── GrantPartyRights ─────────────────────────────────────────────────────────

// GrantPartyRightsInput is the input to [GrantPartyRightsOp].
type GrantPartyRightsInput struct {
	// ParticipantID makes the cache key unique per participant so each actor
	// grants rights on their own Ledger API connection.
	ParticipantID        string `json:"participant_id"`
	DecentralizedPartyID string `json:"decentralized_party_id"`
}

// GrantPartyRightsOutput is the output of [GrantPartyRightsOp].
type GrantPartyRightsOutput struct {
	ParticipantID string `json:"participant_id"`
	UserID        string `json:"user_id"`
	// Granted is true when rights were actively granted; false when skipped
	// because UserID was empty (no-auth environment).
	Granted bool `json:"granted"`
}

// ── UploadDars ───────────────────────────────────────────────────────────────

// UploadDarsInput is the input to [UploadDarsOp].
type UploadDarsInput struct {
	ParticipantID string       `json:"participant_id"`
	Packages      []PackageRef `json:"packages"`
}

// UploadDarsOutput is the output of [UploadDarsOp].
type UploadDarsOutput struct {
	ParticipantID string   `json:"participant_id"`
	PackageIDs    []string `json:"package_ids"`
}

// ── VerifyParty ──────────────────────────────────────────────────────────────

// VerifyPartyInput is the input to [VerifyPartyOp].
type VerifyPartyInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
}

// VerifyPartyOutput is the output of [VerifyPartyOp].
type VerifyPartyOutput struct {
	PartyID  string `json:"party_id"`
	Verified bool   `json:"verified"`
}

// ── PrepareSubmission ────────────────────────────────────────────────────────

// PrepareSubmissionInput is the input to [PrepareSubmissionOp].
type PrepareSubmissionInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
	SynchronizerID       string `json:"synchronizer_id"`
	PackageID            string `json:"package_id"`
	TemplateModule       string `json:"template_module"`
	TemplateEntity       string `json:"template_entity"`
	ContractArgs         string `json:"contract_args"`
}

// PrepareSubmissionOutput is the output of [PrepareSubmissionOp].
type PrepareSubmissionOutput struct {
	PreparedTransactionHash string `json:"prepared_transaction_hash"`
	PreparedTxB64           string `json:"prepared_tx_b64"`
	HashingSchemeVersion    int32  `json:"hashing_scheme_version"`
}

// ── SignSubmission ───────────────────────────────────────────────────────────

// SignSubmissionInput is the input to [SignSubmissionOp].
type SignSubmissionInput struct {
	ParticipantID           string   `json:"participant_id"`
	PreparedTransactionHash string   `json:"prepared_transaction_hash"`
	PreparedTxB64           string   `json:"prepared_tx_b64"`
	KnownSigningKeysB64     []string `json:"known_signing_keys_b64,omitempty"`
}

// SignSubmissionOutput is the output of [SignSubmissionOp].
type SignSubmissionOutput struct {
	ParticipantID  string `json:"participant_id"`
	SignatureB64   string `json:"signature_b64"`
	KeyFingerprint string `json:"key_fingerprint"`
}

// ── ExecuteSubmission ────────────────────────────────────────────────────────

// ExecuteSubmissionInput is the input to [ExecuteSubmissionOp].
type ExecuteSubmissionInput struct {
	DecentralizedPartyID string   `json:"decentralized_party_id"`
	PreparedTxB64        string   `json:"prepared_tx_b64"`
	SignaturesB64        []string `json:"signatures_b64"`
	HashingSchemeVersion int32    `json:"hashing_scheme_version"`
}

// ExecuteSubmissionOutput is the output of [ExecuteSubmissionOp].
type ExecuteSubmissionOutput struct {
	ContractID string `json:"contract_id"`
}

// ── VerifyContract ───────────────────────────────────────────────────────────

// VerifyContractInput is the input to [VerifyContractOp].
type VerifyContractInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
	ContractID           string `json:"contract_id"`
	PackageID            string `json:"package_id"`
	TemplateModule       string `json:"template_module"`
	TemplateEntity       string `json:"template_entity"`
}

// VerifyContractOutput is the output of [VerifyContractOp].
type VerifyContractOutput struct {
	Verified   bool   `json:"verified"`
	ContractID string `json:"contract_id"`
}
