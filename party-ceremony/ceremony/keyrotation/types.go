package keyrotation

// KeyRotationInput is the top-level input to [KeyRotationSequence].
type KeyRotationInput struct {
	// DecentralizedPartyID is the full party identifier in the format
	// "<prefix>::<namespace>", e.g. "cbtc-network::1220abcdef...".
	DecentralizedPartyID string `json:"decentralized_party_id"`

	// TargetParticipantID is the Canton UID of the participant whose key
	// is being rotated (e.g. "PAR::name::fingerprint").
	TargetParticipantID string `json:"target_participant_id"`

	// TargetNamespaceFingerprint is the current namespace fingerprint of the
	// target participant. Used to identify which DNS owner entry to replace
	// during namespace key rotation.
	TargetNamespaceFingerprint string `json:"target_namespace_fingerprint"`

	// SynchronizerID is the Canton synchronizer to update.
	SynchronizerID string `json:"synchronizer_id"`

	// RotateNamespaceKey controls whether the namespace signing key is rotated.
	// When true, the DNS owner list is updated with the new namespace fingerprint.
	RotateNamespaceKey bool `json:"rotate_namespace_key"`

	// RotateDamlKey controls whether the DAML (protocol) signing key is rotated.
	// When true, the P2P signing keys are updated with the new DAML key.
	RotateDamlKey bool `json:"rotate_daml_key"`

	// NewThreshold overrides the threshold after the rotation.
	// When 0 or omitted, the sequence keeps the current threshold.
	NewThreshold int `json:"new_threshold,omitempty"`
}

// Phase represents the current execution phase of the key rotation ceremony.
type Phase string

const (
	PhaseReadState   Phase = "read-state"
	PhaseKeyGen      Phase = "key-gen"
	PhaseNSD         Phase = "nsd"
	PhaseDNSProposal Phase = "dns-proposal"
	PhaseDNSSigning  Phase = "dns-signing"
	PhaseDNSSubmit   Phase = "dns-submit"
	PhaseP2P         Phase = "p2p"
	PhaseCompleted   Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every KeyRotationOutput.
// It is built progressively as the sequence advances, so it is always present
// — even when the sequence returns an error (e.g. ErrThresholdNotMet).
type CeremonyState struct {
	Phase             Phase    `json:"phase"`
	TargetKeyGenReady bool     `json:"target_key_gen_ready"`
	NSDProposed       bool     `json:"nsd_proposed"`
	RequiredSigners   []string `json:"required_signers,omitempty"`
	CollectedSigners  []string `json:"collected_signers"`
	PendingSigners    []string `json:"pending_signers"`
	DNSThreshold      int      `json:"dns_threshold"`
	ProposalHash      string   `json:"proposal_hash,omitempty"`
	P2PProposedCount  int      `json:"p2p_proposed_count"`
	P2PRequired       int      `json:"p2p_required"`
	RotateNamespace   bool     `json:"rotate_namespace"`
	RotateDaml        bool     `json:"rotate_daml"`
}

// KeyRotationOutput is the final result of a completed [KeyRotationSequence].
// State is always populated — even when ExecuteSequence returns an error —
// making it the primary way to inspect ceremony progress.
type KeyRotationOutput struct {
	NamespaceKeyRotated     bool          `json:"namespace_key_rotated"`
	DamlKeyRotated          bool          `json:"daml_key_rotated"`
	NewNamespaceFingerprint string        `json:"new_namespace_fingerprint,omitempty"`
	NewDamlKeyFingerprint   string        `json:"new_daml_key_fingerprint,omitempty"`
	DNSUpdated              bool          `json:"dns_updated"`
	P2PUpdated              bool          `json:"p2p_updated"`
	State                   CeremonyState `json:"state"`
}
