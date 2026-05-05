package kick

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

// Phase represents the current execution phase of the kick ceremony.
type Phase string

const (
	PhaseReadState   Phase = "read-state"
	PhaseDNSProposal Phase = "dns-proposal"
	PhaseDNSSigning  Phase = "dns-signing"
	PhaseDNSSubmit   Phase = "dns-submit"
	PhaseP2P         Phase = "p2p"
	PhaseCompleted   Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every KickOutput.
// It is built progressively as the sequence advances, so it is always present
// — even when the sequence returns an error (e.g. ErrThresholdNotMet).
type CeremonyState struct {
	Phase             Phase    `json:"phase"`
	RequiredSigners   []string `json:"required_signers,omitempty"`
	CollectedSigners  []string `json:"collected_signers"`
	PendingSigners    []string `json:"pending_signers"`
	DNSThreshold      int      `json:"dns_threshold"`
	ProposalHash      string   `json:"proposal_hash,omitempty"`
	P2PProposedCount  int      `json:"p2p_proposed_count"`
	P2PRequired       int      `json:"p2p_required"`
	KickedParticipant string   `json:"kicked_participant"`
	RemainingOwners   []string `json:"remaining_owners,omitempty"`
}

// KickOutput is the final result of a completed [KickSequence].
// State is always populated — even when ExecuteSequence returns an error —
// making it the primary way to inspect ceremony progress.
type KickOutput struct {
	DNSUpdated      bool          `json:"dns_updated"`
	P2PUpdated      bool          `json:"p2p_updated"`
	NewThreshold    int           `json:"new_threshold"`
	RemainingOwners []string      `json:"remaining_owners"`
	State           CeremonyState `json:"state"`
}
