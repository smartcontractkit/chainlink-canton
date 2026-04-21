package ceremony

import (
	"context"
	"errors"
)

// ErrUserRejected is returned when a user declines to sign a transaction
// during an interactive confirmation prompt.
var ErrUserRejected = errors.New("user rejected transaction signing")

// TopologySignDetail contains the decoded fields of a topology transaction
// that are presented to the user for review before signing.
type TopologySignDetail struct {
	// MappingType is the human-readable name of the topology mapping
	// (e.g. "DecentralizedNamespaceDefinition", "PartyToParticipant").
	MappingType string
	// Operation is the topology change operation (e.g. "ADD_REPLACE", "REMOVE").
	Operation string
	// Serial is the topology transaction serial number.
	Serial uint32
	// ExistingSignatures is the number of signatures already present on the transaction.
	ExistingSignatures int
	// ProposalHash is the SHA256 hex string of the proposal (passed through from the operation input).
	ProposalHash string
	// SignerIdentity is the participant UID that will sign.
	SignerIdentity string

	// ── DNS-specific fields (populated when mapping is DecentralizedNamespaceDefinition) ──
	DNSNamespace string
	DNSThreshold int32
	DNSOwners    []string

	// ── P2P-specific fields (populated when mapping is PartyToParticipant) ──
	P2PParty        string
	P2PThreshold    uint32
	P2PParticipants []string

	// ── NSD-specific fields (populated when mapping is NamespaceDelegation) ──
	NSDNamespace string

	// RawMappingType is set when the mapping type is not one of the known types above.
	RawMappingType string
}

// DAMLSignDetail contains transaction metadata presented to the user
// for review before signing a DAML (Ledger API) transaction.
type DAMLSignDetail struct {
	// TransactionHash is the hex-encoded prepared transaction hash.
	TransactionHash string
	// SignerIdentity is the participant UID that will sign.
	SignerIdentity string
}

// Confirmer is called before signing operations to give the user an
// opportunity to review and approve (or reject) the transaction.
//
// Implementations must return nil to proceed with signing, or
// [ErrUserRejected] to abort. Any other error is treated as a
// transient failure.
type Confirmer interface {
	ConfirmTopologySign(ctx context.Context, detail TopologySignDetail) error
	ConfirmDAMLSign(ctx context.Context, detail DAMLSignDetail) error
}

// NoOpConfirmer always approves — used when --confirm is not set.
type NoOpConfirmer struct{}

func (NoOpConfirmer) ConfirmTopologySign(context.Context, TopologySignDetail) error { return nil }
func (NoOpConfirmer) ConfirmDAMLSign(context.Context, DAMLSignDetail) error         { return nil }

// AlwaysRejectConfirmer always rejects — used in negative-path tests.
type AlwaysRejectConfirmer struct{}

func (AlwaysRejectConfirmer) ConfirmTopologySign(context.Context, TopologySignDetail) error {
	return ErrUserRejected
}

func (AlwaysRejectConfirmer) ConfirmDAMLSign(context.Context, DAMLSignDetail) error {
	return ErrUserRejected
}
