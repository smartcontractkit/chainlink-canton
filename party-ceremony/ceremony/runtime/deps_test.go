package runtime

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/stretchr/testify/require"
)

func TestFromCantonParticipant(t *testing.T) {
	participant, err := FromCantonParticipant(canton.Participant{
		Name: "participantA",
		Endpoints: canton.ParticipantEndpoints{
			AdminAPIURL:      "localhost:1201",
			GRPCLedgerAPIURL: "localhost:1301",
		},
		UserID: "user-participant1",
	})
	require.NoError(t, err)

	require.Equal(t, "participantA", participant.Name)
	require.Equal(t, "localhost:1201", participant.AdminAPIURL)
	require.Equal(t, "localhost:1301", participant.GRPCLedgerAPIURL)
	require.Equal(t, "user-participant1", participant.UserID)
	require.Empty(t, participant.AdminJWT)
	require.Empty(t, participant.LedgerJWT)
}
