package replication

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ExportAcsOp exports the Active Contract Set from the source participant.
// The ACS is gzipped and base64-encoded for transport through reports.json.
var ExportAcsOp = operations.NewOperation(
	"canton-ceremony/replication/export-acs",
	semver.MustParse("1.0.0"),
	"Export ACS from source participant as base64-encoded gzip",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ExportAcsInput) (ExportAcsOutput, error) {
		if in.ParticipantID == "" || len(in.PartyIDs) == 0 || in.SynchronizerID == "" {
			return ExportAcsOutput{}, operations.NewUnrecoverableError(
				errors.New("export-acs: participant_id, party_ids, and synchronizer_id are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ExportAcsOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ExportAcsOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		raw, err := deps.Client.ExportAcs(ctx, in.PartyIDs, in.SynchronizerID, in.LedgerOffset)
		if err != nil {
			return ExportAcsOutput{}, fmt.Errorf("exporting ACS: %w", err)
		}

		hash := sha256.Sum256(raw)

		var compressed bytes.Buffer
		gz := gzip.NewWriter(&compressed)
		if _, err := gz.Write(raw); err != nil {
			return ExportAcsOutput{}, fmt.Errorf("gzip compressing ACS: %w", err)
		}
		if err := gz.Close(); err != nil {
			return ExportAcsOutput{}, fmt.Errorf("gzip close: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())

		deps.Logger.Infow("ACS exported",
			"participant", in.ParticipantID,
			"raw_size", len(raw),
			"compressed_size", compressed.Len(),
			"sha256", hex.EncodeToString(hash[:]),
		)

		return ExportAcsOutput{
			ParticipantID:  in.ParticipantID,
			AcsSnapshotB64: encoded,
			SHA256:         hex.EncodeToString(hash[:]),
			SizeBytes:      len(raw),
		}, nil
	},
)
