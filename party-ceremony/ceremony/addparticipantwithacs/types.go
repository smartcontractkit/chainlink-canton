package addparticipantwithacs

// Phase represents the current execution phase of the add-participant-with-acs ceremony.
type Phase string

const (
	PhaseKeyGen             Phase = "key-gen"
	PhaseRecordTargetOffset Phase = "record-target-offset"
	PhaseNSD                Phase = "nsd"
	PhaseReadState          Phase = "read-state"
	PhaseDNSProposal        Phase = "dns-proposal"
	PhaseDNSSigning         Phase = "dns-signing"
	PhaseDNSSubmit          Phase = "dns-submit"
	PhaseRecordOffset       Phase = "record-offset"
	PhaseP2POnboarding      Phase = "p2p-onboarding"
	PhaseAcsExport          Phase = "acs-export"
	PhaseAcsImport          Phase = "acs-import"
	PhaseClearOnboarding    Phase = "clear-onboarding"
	PhaseCompleted          Phase = "completed"
)

// AddParticipantWithAcsInput is the top-level input to [AddParticipantWithAcsSequence].
type AddParticipantWithAcsInput struct {
	// DecentralizedPartyID is the full party identifier in the format
	// "<prefix>::<namespace>", e.g. "cbtc-network::1220abcdef...".
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// NewParticipantID is the Canton UID of the participant being added
	// (e.g. "PAR::newnode::fingerprint").
	NewParticipantID string `json:"new_participant_id"`

	// NamespaceName is the human-readable label used when generating the
	// new participant's namespace and DAML signing keys.
	NamespaceName string `json:"namespace_name"`

	// SynchronizerID is the Canton synchronizer to update.
	SynchronizerID string `json:"synchronizer_id"`

	// SynchronizerAlias is the human-readable alias for the synchronizer,
	// used for disconnect/reconnect during ACS import.
	SynchronizerAlias string `json:"synchronizer_alias"`

	// SourceParticipantID is the Canton UID of the existing participant
	// that will export the ACS for the new participant.
	SourceParticipantID string `json:"source_participant_id"`

	// NewThreshold overrides the threshold after the addition.
	// When 0 or omitted, the sequence keeps the current threshold.
	NewThreshold int `json:"new_threshold,omitempty"`
}

// CeremonyState is the live snapshot embedded in every AddParticipantWithAcsOutput.
type CeremonyState struct {
	Phase                   Phase    `json:"phase"`
	NewMemberKeyReady       bool     `json:"new_member_key_ready"`
	NSDProposed             bool     `json:"nsd_proposed"`
	CollectedSigners        []string `json:"collected_signers"`
	RequiredSigners         []string `json:"required_signers,omitempty"`
	PendingSigners          []string `json:"pending_signers"`
	DNSThreshold            int      `json:"dns_threshold"`
	ProposalHash            string   `json:"proposal_hash,omitempty"`
	P2PExistingProposed     int      `json:"p2p_existing_proposed"`
	P2PExistingRequired     int      `json:"p2p_existing_required"`
	NewParticipantConsented bool     `json:"new_participant_consented"`
	AllOwners               []string `json:"all_owners,omitempty"`
	TargetLedgerOffset      int64    `json:"target_ledger_offset"`
	LedgerOffsetRecorded    bool     `json:"ledger_offset_recorded"`
	AcsExported             bool     `json:"acs_exported"`
	AcsImported             bool     `json:"acs_imported"`
	OnboardingFlagCleared   bool     `json:"onboarding_flag_cleared"`
}

// AddParticipantWithAcsOutput is the final result of a completed
// [AddParticipantWithAcsSequence].
type AddParticipantWithAcsOutput struct {
	State        CeremonyState `json:"state"`
	AllOwners    []string      `json:"all_owners"`
	NewThreshold int           `json:"new_threshold"`
	DNSUpdated   bool          `json:"dns_updated"`
	P2PUpdated   bool          `json:"p2p_updated"`
	AcsImported  bool          `json:"acs_imported"`
}
