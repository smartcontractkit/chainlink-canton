package contract

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

// ParticipantAt returns the participant at index on chain.
// Zero value index selects the first participant.
func ParticipantAt(chain canton.Chain, index int) (canton.Participant, error) {
	if index < 0 || index >= len(chain.Participants) {
		return canton.Participant{}, fmt.Errorf(
			"participant index %d out of range for chain with %d participants",
			index, len(chain.Participants),
		)
	}

	return chain.Participants[index], nil
}
