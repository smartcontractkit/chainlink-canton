package client

// DNSState holds the current state of a DecentralizedNamespaceDefinition
// topology mapping as read from Canton's topology store.
type DNSState struct {
	// DecentralizedNamespace is the 68-character hex multihash identifying the
	// decentralized namespace.
	DecentralizedNamespace string

	// Owners is the list of owner namespace fingerprints.
	Owners []string

	// Threshold is the number of owner signatures required for topology updates
	// within the decentralized namespace.
	Threshold int32

	// Serial is the topology serial number of this mapping. Updates must use
	// serial = Serial + 1.
	Serial int32
}

// P2PParticipantInfo holds the participant UID and permission level for a
// single hosting participant in a PartyToParticipant topology mapping.
type P2PParticipantInfo struct {
	// ParticipantUID is the Canton-assigned participant UID
	// (e.g. "PAR::name::fingerprint").
	ParticipantUID string

	// Permission is the human-readable permission level
	// (e.g. "PARTICIPANT_PERMISSION_CONFIRMATION").
	Permission string
}

// P2PState holds the current state of a PartyToParticipant topology mapping
// as read from Canton's topology store.
type P2PState struct {
	// Party is the full party ID (e.g. "prefix::namespace").
	Party string

	// Participants is the list of hosting participants with their permissions.
	Participants []P2PParticipantInfo

	// Threshold is the confirmation threshold for the party.
	Threshold uint32

	// Serial is the topology serial number of this mapping.
	Serial int32
}
