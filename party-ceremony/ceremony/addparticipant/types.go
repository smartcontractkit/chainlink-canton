package addparticipant

// AddParticipantInput is the top-level input to [AddParticipantSequence].
type AddParticipantInput struct {
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

	// NewThreshold overrides the threshold after the addition.
	// When 0 or omitted, the sequence keeps the current threshold.
	NewThreshold int `json:"new_threshold,omitempty"`
}

// AddParticipantOutput is the final result of a completed [AddParticipantSequence].
// Phase represents the current execution phase of the add-participant ceremony.
type Phase string

const (
	PhaseKeyGen      Phase = "key-gen"
	PhaseNSD         Phase = "nsd"
	PhaseReadState   Phase = "read-state"
	PhaseDNSProposal Phase = "dns-proposal"
	PhaseDNSSigning  Phase = "dns-signing"
	PhaseDNSSubmit   Phase = "dns-submit"
	PhaseP2P         Phase = "p2p"
	PhaseCompleted   Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every AddParticipantOutput.
// It is built progressively as the sequence advances, so it is always present
type CeremonyState struct {
	Phase                   Phase    `json:"phase"`
	NewMemberKeyReady       bool     `json:"new_member_key_ready"`
	NSDProposed             bool     `json:"nsd_proposed"`
	RequiredSigners         []string `json:"required_signers,omitempty"`
	CollectedSigners        []string `json:"collected_signers"`
	PendingSigners          []string `json:"pending_signers"`
	DNSThreshold            int      `json:"dns_threshold"`
	ProposalHash            string   `json:"proposal_hash,omitempty"`
	P2PExistingProposed     int      `json:"p2p_existing_proposed"`
	P2PExistingRequired     int      `json:"p2p_existing_required"`
	NewParticipantConsented bool     `json:"new_participant_consented"`
	AllOwners               []string `json:"all_owners,omitempty"`
}

type AddParticipantOutput struct {
	DNSUpdated   bool          `json:"dns_updated"`
	P2PUpdated   bool          `json:"p2p_updated"`
	NewThreshold int           `json:"new_threshold"`
	AllOwners    []string      `json:"all_owners"`
	State        CeremonyState `json:"state"`
}
