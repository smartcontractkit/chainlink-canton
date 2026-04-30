package contractdeploy

import "github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"

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
