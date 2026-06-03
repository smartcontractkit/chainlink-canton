package replication

import (
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// RecordLedgerOffsetOp records the current ledger offset on the source
// participant. This offset is used as the snapshot point for ACS export.
var RecordLedgerOffsetOp = operations.NewOperation(
	"canton-ceremony/replication/record-ledger-offset",
	semver.MustParse("1.0.0"),
	"Record current ledger offset on source participant for ACS export",
	func(b operations.Bundle, deps ceremony.CantonDeps, in RecordLedgerOffsetInput) (RecordLedgerOffsetOutput, error) {
		if in.ParticipantID == "" || in.SynchronizerID == "" {
			return RecordLedgerOffsetOutput{}, operations.NewUnrecoverableError(
				errors.New("record-ledger-offset: participant_id and synchronizer_id are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return RecordLedgerOffsetOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return RecordLedgerOffsetOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		now := time.Now().Unix()
		offset, err := deps.Client.LookupOffsetByTime(ctx, now)
		if err != nil {
			return RecordLedgerOffsetOutput{}, fmt.Errorf("looking up ledger offset: %w", err)
		}

		deps.Logger.Infow("Recorded ledger offset",
			"participant", in.ParticipantID,
			"offset", offset,
			"timestamp", now,
		)

		return RecordLedgerOffsetOutput{
			ParticipantID:    in.ParticipantID,
			LedgerOffset:     offset,
			TimestampSeconds: now,
		}, nil
	},
)
