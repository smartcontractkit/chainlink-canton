package onboarding

// OnboardingInput is the top-level input to [OnboardingSequence].
type OnboardingInput struct {
	NamespaceName string `json:"namespace_name"`
	// KmsVaultName optionally overrides NamespaceName for KMS vault registration.
	// Use when reusing KMS keys already registered under another ceremony name.
	KmsVaultName   string   `json:"kms_vault_name,omitempty"`
	PartyPrefix    string   `json:"party_prefix"`
	Participants   []string `json:"participants"`
	SynchronizerID string   `json:"synchronizer_id"`
	// Threshold sets the signing/confirmation threshold for the decentralized
	// namespace and the party-to-participant mapping. When 0 the sequence
	// defaults to strict-majority: floor(n/2)+1.
	Threshold int `json:"threshold"`
}

// Phase represents the current execution phase of the onboarding ceremony.
type Phase string

const (
	PhaseKeyGen      Phase = "key-gen"
	PhaseNSD         Phase = "nsd"
	PhaseDNSProposal Phase = "dns-proposal"
	PhaseDNSSigning  Phase = "dns-signing"
	PhaseDNSSubmit   Phase = "dns-submit"
	PhaseP2P         Phase = "p2p"
	PhaseCompleted   Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every OnboardingOutput.
// It is built progressively as the sequence advances, so it is always present
// — even when the sequence returns an error (e.g. ErrThresholdNotMet).
type CeremonyState struct {
	Phase            Phase    `json:"phase"`
	KeysGenerated    []string `json:"keys_generated"`
	NSDsProposed     []string `json:"nsds_proposed"`
	RequiredSigners  []string `json:"required_signers,omitempty"`
	CollectedSigners []string `json:"collected_signers"`
	PendingSigners   []string `json:"pending_signers"`
	Threshold        int      `json:"threshold"`
	ProposalHash     string   `json:"proposal_hash,omitempty"`
	P2PProposedCount int      `json:"p2p_proposed_count"`
	P2PRequired      int      `json:"p2p_required"`
}

// OnboardingOutput is the final result of a completed [OnboardingSequence].
// State is always populated — even when ExecuteSequence returns an error —
// making it the primary way to inspect ceremony progress.
type OnboardingOutput struct {
	PartyID      string        `json:"party_id"`
	DNSConfirmed bool          `json:"dns_confirmed"`
	P2PConfirmed bool          `json:"p2p_confirmed"`
	State        CeremonyState `json:"state"`
}
