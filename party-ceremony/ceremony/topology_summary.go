package ceremony

import (
	"encoding/base64"
	"fmt"

	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"google.golang.org/protobuf/proto"
)

// SummarizeTopologyTx decodes a base64-encoded SignedTopologyTransaction and
// extracts a human-readable [TopologySignDetail]. The proposalHash and
// signerIdentity are passed through from the caller (operation input) since
// they are not embedded in the proto itself.
func SummarizeTopologyTx(txB64, proposalHash, signerIdentity string) (TopologySignDetail, error) {
	txBytes, err := base64.StdEncoding.DecodeString(txB64)
	if err != nil {
		return TopologySignDetail{}, fmt.Errorf("decoding base64: %w", err)
	}

	var stx protov30.SignedTopologyTransaction
	if err := proto.Unmarshal(txBytes, &stx); err != nil {
		return TopologySignDetail{}, fmt.Errorf("unmarshalling SignedTopologyTransaction: %w", err)
	}

	detail := TopologySignDetail{
		ExistingSignatures: len(stx.GetSignatures()),
		ProposalHash:       proposalHash,
		SignerIdentity:     signerIdentity,
	}

	// The Transaction field is a serialised TopologyTransaction.
	var ttx protov30.TopologyTransaction
	if err := proto.Unmarshal(stx.GetTransaction(), &ttx); err != nil {
		// If inner decoding fails, return what we have with a generic mapping type.
		detail.MappingType = "unknown (inner decode failed)"
		return detail, nil //nolint:nilerr // partial result is acceptable
	}

	detail.Operation = ttx.GetOperation().String()
	detail.Serial = ttx.GetSerial()

	mapping := ttx.GetMapping()
	if mapping == nil {
		detail.MappingType = "unknown (no mapping)"
		return detail, nil
	}

	switch m := mapping.GetMapping().(type) {
	case *protov30.TopologyMapping_DecentralizedNamespaceDefinition:
		dns := m.DecentralizedNamespaceDefinition
		detail.MappingType = "DecentralizedNamespaceDefinition"
		detail.DNSNamespace = dns.GetDecentralizedNamespace()
		detail.DNSThreshold = dns.GetThreshold()
		detail.DNSOwners = dns.GetOwners()

	case *protov30.TopologyMapping_PartyToParticipant:
		p2p := m.PartyToParticipant
		detail.MappingType = "PartyToParticipant"
		detail.P2PParty = p2p.GetParty()
		detail.P2PThreshold = p2p.GetThreshold()
		for _, hp := range p2p.GetParticipants() {
			detail.P2PParticipants = append(detail.P2PParticipants,
				fmt.Sprintf("%s (%s)", hp.GetParticipantUid(), hp.GetPermission()))
		}

	case *protov30.TopologyMapping_NamespaceDelegation:
		nsd := m.NamespaceDelegation
		detail.MappingType = "NamespaceDelegation"
		detail.NSDNamespace = nsd.GetNamespace()

	default:
		detail.MappingType = fmt.Sprintf("%T", mapping.GetMapping())
		detail.RawMappingType = detail.MappingType
	}

	return detail, nil
}
