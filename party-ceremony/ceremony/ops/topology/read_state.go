package topology

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/chainlink/canton-party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ReadCurrentStateOp queries the Canton topology store for the current
// DecentralizedNamespaceDefinition and PartyToParticipant state. It provides
// the serial numbers and existing owner/participant lists required by all
// subsequent topology operations.
//
// This operation is read-only and naturally idempotent.
var ReadCurrentStateOp = operations.NewOperation(
	"canton-ceremony/topology/read-current-state",
	semver.MustParse("1.0.0"),
	"Query current DNS and P2P topology state from the synchronizer",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ReadCurrentStateInput) (ReadCurrentStateOutput, error) {
		if in.DecentralizedPartyID == "" {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				errors.New("read-current-state: decentralized_party_id is required"),
			)
		}

		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("read-current-state: invalid decentralized_party_id format %q: expected <prefix>::<namespace>",
					in.DecentralizedPartyID),
			)
		}
		decNS := parts[1]

		ctx := b.GetContext()

		dnsState, err := deps.Client.GetDNS(ctx, decNS, in.SynchronizerID)
		if err != nil {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("reading DNS state for namespace %q: %w", decNS, err),
			)
		}

		p2pState, err := deps.Client.GetP2P(ctx, in.DecentralizedPartyID, in.SynchronizerID)
		if err != nil {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("reading P2P state for party %q: %w", in.DecentralizedPartyID, err),
			)
		}

		participantUIDs := make([]string, len(p2pState.Participants))
		for i, p := range p2pState.Participants {
			participantUIDs[i] = p.ParticipantUID
		}

		deps.Logger.Infow("Current state read",
			"namespace", decNS,
			"dns_owners_count", len(dnsState.Owners),
			"dns_threshold", dnsState.Threshold,
			"dns_serial", dnsState.Serial,
			"p2p_participants_count", len(participantUIDs),
			"p2p_threshold", p2pState.Threshold,
			"p2p_serial", p2pState.Serial,
		)

		out := ReadCurrentStateOutput{
			DecentralizedNamespace: dnsState.DecentralizedNamespace,
			DNSOwners:              dnsState.Owners,
			DNSThreshold:           dnsState.Threshold,
			DNSSerial:              dnsState.Serial,
			P2PParticipantUIDs:     participantUIDs,
			P2PThreshold:           p2pState.Threshold,
			P2PSerial:              p2pState.Serial,
		}
		if p2pState.PartySigningKeys != nil {
			out.PartySigningKeysB64 = p2pState.PartySigningKeys.Keys
			out.PartySigningThreshold = p2pState.PartySigningKeys.Threshold
		}

		return out, nil
	},
)
