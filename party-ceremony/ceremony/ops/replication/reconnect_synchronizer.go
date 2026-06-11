package replication

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// ReconnectSynchronizerOp reconnects the target participant after its party ACS
// snapshot has been imported successfully.
var ReconnectSynchronizerOp = operations.NewOperation(
	"canton-ceremony/replication/reconnect-synchronizer",
	semver.MustParse("1.0.0"),
	"Reconnect target participant to synchronizer after ACS import",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ReconnectSynchronizerInput) (ReconnectSynchronizerOutput, error) {
		if in.ParticipantID == "" || in.SynchronizerAlias == "" {
			return ReconnectSynchronizerOutput{}, operations.NewUnrecoverableError(
				errors.New("reconnect-synchronizer: participant_id and synchronizer_alias are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ReconnectSynchronizerOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ReconnectSynchronizerOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		deps.Logger.Infow("Reconnecting target to synchronizer after ACS import",
			"participant", in.ParticipantID,
			"alias", in.SynchronizerAlias,
		)
		if err := deps.Client.ReconnectSynchronizer(ctx, in.SynchronizerAlias); err != nil {
			return ReconnectSynchronizerOutput{}, fmt.Errorf("reconnecting synchronizer: %w", err)
		}

		return ReconnectSynchronizerOutput{
			ParticipantID: in.ParticipantID,
			Reconnected:   true,
		}, nil
	},
)
