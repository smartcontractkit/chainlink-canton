package replication

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ImportAcsOp imports an ACS snapshot into the target participant.
// It disconnects from the synchronizer, imports the ACS, then reconnects.
var ImportAcsOp = operations.NewOperation(
	"canton-ceremony/replication/import-acs",
	semver.MustParse("1.0.0"),
	"Import ACS into target participant (disconnect → import → reconnect)",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ImportAcsInput) (ImportAcsOutput, error) {
		if in.ParticipantID == "" || in.SynchronizerID == "" || in.SynchronizerAlias == "" || in.AcsSnapshotB64 == "" {
			return ImportAcsOutput{}, operations.NewUnrecoverableError(
				errors.New("import-acs: participant_id, synchronizer_id, synchronizer_alias, and acs_snapshot_b64 are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ImportAcsOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ImportAcsOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		// Decode base64 → gzip → raw ACS.
		compressed, err := base64.StdEncoding.DecodeString(in.AcsSnapshotB64)
		if err != nil {
			return ImportAcsOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding ACS base64: %w", err),
			)
		}

		gz, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			return ImportAcsOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("creating gzip reader: %w", err),
			)
		}
		raw, err := io.ReadAll(gz)
		if err != nil {
			return ImportAcsOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decompressing ACS: %w", err),
			)
		}
		gz.Close()

		// Verify integrity.
		if in.ExpectedSHA256 != "" {
			hash := sha256.Sum256(raw)
			actual := hex.EncodeToString(hash[:])
			if actual != in.ExpectedSHA256 {
				return ImportAcsOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("ACS integrity check failed: expected SHA256 %s, got %s",
						in.ExpectedSHA256, actual),
				)
			}
		}

		// Disconnect → Import → Reconnect.
		deps.Logger.Infow("Disconnecting from synchronizer for ACS import",
			"participant", in.ParticipantID,
			"alias", in.SynchronizerAlias,
		)
		if err := deps.Client.DisconnectSynchronizer(ctx, in.SynchronizerAlias); err != nil {
			return ImportAcsOutput{}, fmt.Errorf("disconnecting synchronizer: %w", err)
		}

		deps.Logger.Infow("Importing ACS",
			"participant", in.ParticipantID,
			"raw_size", len(raw),
		)
		if err := deps.Client.ImportAcs(ctx, raw, in.SynchronizerID); err != nil {
			// Attempt reconnect even if import fails.
			deps.Logger.Errorw("ACS import failed, attempting reconnect",
				"err", err,
				"alias", in.SynchronizerAlias,
			)
			_ = deps.Client.ReconnectSynchronizer(ctx, in.SynchronizerAlias)

			return ImportAcsOutput{}, fmt.Errorf("importing ACS: %w", err)
		}

		deps.Logger.Infow("Reconnecting to synchronizer after ACS import",
			"participant", in.ParticipantID,
			"alias", in.SynchronizerAlias,
		)
		if err := deps.Client.ReconnectSynchronizer(ctx, in.SynchronizerAlias); err != nil {
			return ImportAcsOutput{}, fmt.Errorf("reconnecting synchronizer: %w", err)
		}

		deps.Logger.Infow("ACS imported successfully",
			"participant", in.ParticipantID,
		)

		return ImportAcsOutput{
			ParticipantID: in.ParticipantID,
			Imported:      true,
		}, nil
	},
)
