package replication

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// DisconnectSynchronizerOp disconnects the target participant at the documented
// offline party replication point: immediately after the target authorizes
// hosting with the onboarding flag set.
var DisconnectSynchronizerOp = operations.NewOperation(
	"canton-ceremony/replication/disconnect-synchronizer",
	semver.MustParse("1.0.0"),
	"Disconnect target participant from synchronizer before ACS export/import",
	func(b operations.Bundle, deps ceremony.CantonDeps, in DisconnectSynchronizerInput) (DisconnectSynchronizerOutput, error) {
		if in.ParticipantID == "" || in.SynchronizerAlias == "" {
			return DisconnectSynchronizerOutput{}, operations.NewUnrecoverableError(
				errors.New("disconnect-synchronizer: participant_id and synchronizer_alias are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return DisconnectSynchronizerOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return DisconnectSynchronizerOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		deps.Logger.Infow("Disconnecting target from synchronizer before ACS replication",
			"participant", in.ParticipantID,
			"alias", in.SynchronizerAlias,
		)
		if err := deps.Client.DisconnectSynchronizer(ctx, in.SynchronizerAlias); err != nil {
			return DisconnectSynchronizerOutput{}, fmt.Errorf("disconnecting synchronizer: %w", err)
		}

		return DisconnectSynchronizerOutput{
			ParticipantID: in.ParticipantID,
			Disconnected:  true,
		}, nil
	},
)
