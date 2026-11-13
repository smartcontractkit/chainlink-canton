package exercisechoice

import "github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"

// ExerciseChoiceInput is the top-level input to [ExerciseChoiceSequence].
type ExerciseChoiceInput struct {
	// DecentralizedPartyID is the full party ID (e.g. "prefix::namespace")
	// from a prior onboarding ceremony.
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// SynchronizerID is the Canton synchronizer to target.
	SynchronizerID string `json:"synchronizer_id"`

	// PackageID is a package ID or a "#<package-name>" reference. The package
	// must already be vetted on every hosting participant.
	PackageID string `json:"package_id"`

	// TemplateModule is the fully-qualified DAML module name (e.g. "Utility.Registry.App.V0.Service.Provider").
	TemplateModule string `json:"template_module"`

	// TemplateEntity is the DAML template entity name (e.g. "ProviderService").
	TemplateEntity string `json:"template_entity"`

	// ContractID is the contract to exercise the choice on.
	ContractID string `json:"contract_id"`

	// Choice is the choice name (e.g. "ProviderService_CreateProviderConfiguration").
	Choice string `json:"choice"`

	// ChoiceArgument is the choice argument as a Ledger API Value in proto JSON.
	ChoiceArgument string `json:"choice_argument"`
}

// Phase represents the current execution phase of the exercise choice ceremony.
type Phase string

const (
	PhaseVerifyParty  Phase = "verify-party"
	PhaseFetchMembers Phase = "fetch-members"
	PhaseGrantRights  Phase = "grant-rights"
	PhasePrepare      Phase = "prepare"
	PhaseSigning      Phase = "signing"
	PhaseExecute      Phase = "execute"
	PhaseCompleted    Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every ExerciseChoiceOutput.
// It is populated even when the sequence returns an error.
type CeremonyState struct {
	Phase          Phase    `json:"phase"`
	Participants   []string `json:"participants,omitempty"`
	Granted        []string `json:"granted"`
	Signed         []string `json:"signed"`
	SignRequired   int      `json:"sign_required"`
	PreparedTxHash string   `json:"prepared_tx_hash,omitempty"`
}

// ExerciseChoiceOutput is the result of [ExerciseChoiceSequence].
type ExerciseChoiceOutput struct {
	PreparedTransactionHash string                   `json:"prepared_transaction_hash"`
	UpdateID                string                   `json:"update_id"`
	CreatedContracts        []client.CreatedContract `json:"created_contracts"`
	State                   CeremonyState            `json:"state"`
}
