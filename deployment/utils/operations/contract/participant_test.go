package contract

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/stretchr/testify/require"
)

func TestParticipantAt(t *testing.T) {
	t.Parallel()

	chain := canton.Chain{
		Participants: []canton.Participant{
			{PartyID: "party-0"},
			{PartyID: "party-1"},
		},
	}

	p0, err := ParticipantAt(chain, 0)
	require.NoError(t, err)
	require.Equal(t, "party-0", p0.PartyID)

	p1, err := ParticipantAt(chain, 1)
	require.NoError(t, err)
	require.Equal(t, "party-1", p1.PartyID)

	_, err = ParticipantAt(chain, 2)
	require.Error(t, err)

	_, err = ParticipantAt(chain, -1)
	require.Error(t, err)
}
