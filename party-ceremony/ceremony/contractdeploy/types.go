package contractdeploy

import "github.com/chainlink/canton-party-ceremony/ceremony/ops/ledger"

// PackageRef identifies a DAML package by its registered name and version.
// The name corresponds to a [contracts.Package] constant (e.g. "mcms") and the
// version to a registered version string (e.g. "current" or "0.0.1").
type PackageRef = ledger.PackageRef

// ContractDeployInput is the top-level input to [ContractDeploySequence].
type ContractDeployInput struct {
	// DecentralizedPartyID is the full party ID (e.g. "prefix::namespace")
	// from a prior onboarding ceremony.
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// SynchronizerID is the Canton synchronizer to target.
	SynchronizerID string `json:"synchronizer_id"`

	// Packages lists the DAML packages to upload by name and version.
	// Each entry resolves to a DAR via the DARLoader in [ContractDeployDeps].
	Packages []PackageRef `json:"packages"`

	// TemplateModule is the fully-qualified DAML module name (e.g. "MCMS.Main").
	TemplateModule string `json:"template_module"`

	// TemplateEntity is the DAML template entity name (e.g. "MCMS").
	TemplateEntity string `json:"template_entity"`

	// ContractArgs is a JSON-encoded string of the contract creation arguments.
	// The format depends on the specific template being deployed.
	ContractArgs string `json:"contract_args"`
}

// Phase represents the current execution phase of the contract deploy ceremony.
type Phase string

const (
	PhaseVerifyParty    Phase = "verify-party"
	PhaseFetchMembers   Phase = "fetch-members"
	PhaseDARUpload      Phase = "dar-upload"
	PhasePrepare        Phase = "prepare"
	PhaseSigning        Phase = "signing"
	PhaseExecute        Phase = "execute"
	PhaseVerifyContract Phase = "verify-contract"
	PhaseCompleted      Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every ContractDeployOutput.
// It is built progressively as the sequence advances, so it is always present
// — even when the sequence returns an error (e.g. ErrThresholdNotMet).
type CeremonyState struct {
	Phase          Phase    `json:"phase"`
	Participants   []string `json:"participants,omitempty"`
	DARsUploaded   []string `json:"dars_uploaded"`
	DARsRequired   int      `json:"dars_required"`
	Signed         []string `json:"signed"`
	SignRequired   int      `json:"sign_required"`
	PreparedTxHash string   `json:"prepared_tx_hash,omitempty"`
}

// ContractDeployOutput is the final result of a completed [ContractDeploySequence].
// State is always populated — even when ExecuteSequence returns an error —
// making it the primary way to inspect ceremony progress.
type ContractDeployOutput struct {
	PackageIDs              []string      `json:"package_ids"`
	PreparedTransactionHash string        `json:"prepared_transaction_hash"`
	ContractID              string        `json:"contract_id"`
	State                   CeremonyState `json:"state"`
}

// FetchParticipantsInput is the input to [FetchParticipantsOp].
type FetchParticipantsInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
	SynchronizerID       string `json:"synchronizer_id"`
}

// FetchParticipantsOutput is the output of [FetchParticipantsOp].
type FetchParticipantsOutput struct {
	// Participants is the list of participant UIDs hosting the decentralized party
	// (e.g. "PAR::name::fingerprint").
	Participants []string `json:"participants"`
}

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

// VerifyPartyInput is the input to [VerifyPartyOp].
type VerifyPartyInput struct {
	DecentralizedPartyID string `json:"decentralized_party_id"`
}

// VerifyPartyOutput is the output of [VerifyPartyOp].
type VerifyPartyOutput struct {
	PartyID  string `json:"party_id"`
	Verified bool   `json:"verified"`
}

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
	// PreparedTxB64 is the base64-encoded serialized PreparedTransaction proto.
	// Needed by the signing and execution steps.
	PreparedTxB64 string `json:"prepared_tx_b64"`
	// HashingSchemeVersion is forwarded from PrepareSubmissionResponse and must
	// be passed as-is to ExecuteSubmissionRequest.
	HashingSchemeVersion int32 `json:"hashing_scheme_version"`
}

// SignSubmissionInput is the input to [SignSubmissionOp].
type SignSubmissionInput struct {
	ParticipantID           string `json:"participant_id"`
	PreparedTransactionHash string `json:"prepared_transaction_hash"`
	PreparedTxB64           string `json:"prepared_tx_b64"`
}

// SignSubmissionOutput is the output of [SignSubmissionOp].
type SignSubmissionOutput struct {
	ParticipantID  string `json:"participant_id"`
	SignatureB64   string `json:"signature_b64"`   // base64-encoded serialised v2.Signature proto
	KeyFingerprint string `json:"key_fingerprint"` // Canton key fingerprint used to sign
}

// ExecuteSubmissionInput is the input to [ExecuteSubmissionOp].
type ExecuteSubmissionInput struct {
	DecentralizedPartyID string   `json:"decentralized_party_id"`
	PreparedTxB64        string   `json:"prepared_tx_b64"`
	SignaturesB64        []string `json:"signatures_b64"` // one per participant
	HashingSchemeVersion int32    `json:"hashing_scheme_version"`
}

// ExecuteSubmissionOutput is the output of [ExecuteSubmissionOp].
type ExecuteSubmissionOutput struct {
	ContractID string `json:"contract_id"`
}

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
