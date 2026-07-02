package devenv

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

const (
	// OwnerParticipantIndex is Participants[0]: deploy, EDS, verifiers, fee aggregator.
	OwnerParticipantIndex = 0
	// ClientParticipantIndex is Participants[1]: send, receive, execute, token holdings.
	ClientParticipantIndex = 1
)

func (c *Chain) ClientParticipant() (canton.Participant, int, error) {
	if len(c.chain.Participants) == 0 {
		return canton.Participant{}, 0, fmt.Errorf("canton chain has no participants")
	}
	idx := c.clientParticipantIndex()

	return c.chain.Participants[idx], idx, nil
}

func (c *Chain) OwnerParticipant() (canton.Participant, error) {
	return ownerParticipantFromBlockchain(c.chain)
}

// ownerParticipantFromBlockchain returns the owner participant (index 0), or an error if none exist.
func ownerParticipantFromBlockchain(chain canton.Chain) (canton.Participant, error) {
	if len(chain.Participants) == 0 {
		return canton.Participant{}, fmt.Errorf("canton chain has no participants")
	}

	return chain.Participants[OwnerParticipantIndex], nil
}

func (c *Chain) clientParticipantIndex() int {
	if len(c.chain.Participants) > ClientParticipantIndex {
		return ClientParticipantIndex
	}

	return OwnerParticipantIndex
}
