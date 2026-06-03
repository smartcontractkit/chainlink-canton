package replication

import (
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// RecordTargetLedgerOffsetOp records the current ledger offset on the target
// (new) participant. The offset is later used as BeginOffsetExclusive for
// ClearPartyOnboardingFlag, so the forward search skips past any prior
// activation of the same party on this participant (e.g. before a previous
// kick) and lands on the new activation that carries the onboarding flag.
//
// Must run on the target participant during its first ceremony resume,
// before the new P2P topology with the onboarding flag is broadcast.
var RecordTargetLedgerOffsetOp = operations.NewOperation(
	"canton-ceremony/replication/record-target-ledger-offset",
	semver.MustParse("1.0.0"),
	"Record current ledger offset on target participant to seed ClearPartyOnboardingFlag search",
	func(b operations.Bundle, deps ceremony.CantonDeps, in RecordTargetLedgerOffsetInput) (RecordTargetLedgerOffsetOutput, error) {
		if in.ParticipantID == "" || in.SynchronizerID == "" {
			return RecordTargetLedgerOffsetOutput{}, operations.NewUnrecoverableError(
				errors.New("record-target-ledger-offset: participant_id and synchronizer_id are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return RecordTargetLedgerOffsetOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return RecordTargetLedgerOffsetOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		now := time.Now().Unix()
		offset, err := deps.Client.LookupOffsetByTime(ctx, now)
		if err != nil {
			return RecordTargetLedgerOffsetOutput{}, fmt.Errorf("looking up ledger offset: %w", err)
		}

		deps.Logger.Infow("Recorded target ledger offset",
			"participant", in.ParticipantID,
			"offset", offset,
			"timestamp", now,
		)

		return RecordTargetLedgerOffsetOutput{
			ParticipantID:    in.ParticipantID,
			LedgerOffset:     offset,
			TimestampSeconds: now,
		}, nil
	},
)
