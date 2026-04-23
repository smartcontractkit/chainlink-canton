package topology

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// ProposeP2POp authorizes the PartyToParticipant mapping for a single
// participant during initial onboarding. Each participant runs this independently.
//
// Canton accumulates proposals from each participant automatically; once the
// threshold is reached, the mapping is activated.
var ProposeP2POp = operations.NewOperation(
	"canton-ceremony/topology/propose-p2p",
	semver.MustParse("1.0.0"),
	"Propose PartyToParticipant mapping for a single participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeP2PInput) (ProposeP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.Members) == 0 {
			return ProposeP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-p2p: participant_id, party_id, and members are required"),
			)
		}

		ctx := b.GetContext()

		// This operation can only run on the participant that owns the signing key.
		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ProposeP2POutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s",
				in.ParticipantID, currentUID)
		}

		threshold := in.Threshold
		if threshold <= 0 {
			threshold = len(in.Members)/2 + 1
		}

		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.Members))
		for i, m := range in.Members {
			hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: m.ParticipantUID,
				Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
			}
		}

		// Collect each member's PROTOCOL (DAML) signing key so Canton can
		// verify transaction signatures submitted via InteractiveSubmissionService.
		damlKeys := make([]*cryptov30.SigningPublicKey, 0, len(in.Members))
		for _, m := range in.Members {
			if m.DamlKeyB64 == "" {
				continue
			}
			damlKeyBytes, decErr := base64.StdEncoding.DecodeString(m.DamlKeyB64)
			if decErr != nil {
				return ProposeP2POutput{}, fmt.Errorf("decoding DAML key for %s: %w", m.ParticipantID, decErr)
			}
			var damlKey cryptov30.SigningPublicKey
			if unmErr := proto.Unmarshal(damlKeyBytes, &damlKey); unmErr != nil {
				return ProposeP2POutput{}, fmt.Errorf("unmarshalling DAML key for %s: %w", m.ParticipantID, unmErr)
			}
			damlKeys = append(damlKeys, &damlKey)
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_PartyToParticipant{
				PartyToParticipant: &protov30.PartyToParticipant{
					Party:        in.PartyID,
					Threshold:    uint32(threshold),
					Participants: hostingParticipants,
					PartySigningKeys: &cryptov30.SigningKeysWithThreshold{
						Keys:      damlKeys,
						Threshold: uint32(threshold),
					},
				},
			},
		}

		// mustFullyAuthorize=false: Canton stores the proposal and accumulates
		// signatures from other participants as they also propose.
		_, err = deps.Client.Authorize(ctx, 1, mapping, in.SynchronizerID, false)
		if err != nil {
			return ProposeP2POutput{}, fmt.Errorf("authorizing P2P mapping: %w", err)
		}

		deps.Logger.Infow("P2P proposal submitted",
			"participant", in.ParticipantID, "party", in.PartyID)

		return ProposeP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
		}, nil
	},
)

// ProposeKickP2POp authorizes the updated PartyToParticipant mapping with the
// kicked participant removed. Each remaining participant runs this independently.
var ProposeKickP2POp = operations.NewOperation(
	"canton-ceremony/topology/propose-kick-p2p",
	semver.MustParse("1.0.0"),
	"Propose updated PartyToParticipant mapping with kicked participant removed",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeKickP2PInput) (ProposeKickP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.RemainingParticipants) == 0 {
			return ProposeKickP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-kick-p2p: participant_id, party_id, and remaining_participants are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeKickP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}

		if currentUID != in.ParticipantID {
			return ProposeKickP2POutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.RemainingParticipants))
		for i, uid := range in.RemainingParticipants {
			hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: uid,
				Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
			}
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_PartyToParticipant{
				PartyToParticipant: &protov30.PartyToParticipant{
					Party:        in.PartyID,
					Threshold:    uint32(in.NewP2PThreshold),
					Participants: hostingParticipants,
				},
			},
		}

		_, err = deps.Client.Authorize(ctx, uint32(in.CurrentP2PSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return ProposeKickP2POutput{}, fmt.Errorf("authorizing kick P2P proposal: %w", err)
		}

		deps.Logger.Infow("Kick P2P proposal submitted",
			"participant", in.ParticipantID,
			"party", in.PartyID,
			"remaining_count", len(in.RemainingParticipants),
		)

		return ProposeKickP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
		}, nil
	},
)

// ProposeAddP2POp authorizes the updated PartyToParticipant mapping with the
// new participant added. Each existing participant runs this independently.
var ProposeAddP2POp = operations.NewOperation(
	"canton-ceremony/topology/propose-add-p2p",
	semver.MustParse("1.0.0"),
	"Propose updated PartyToParticipant mapping with new participant added",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeAddP2PInput) (ProposeAddP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.AllParticipantUIDs) == 0 {
			return ProposeAddP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-add-p2p: participant_id, party_id, and all_participant_uids are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeAddP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}

		if currentUID != in.ParticipantID {
			return ProposeAddP2POutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.AllParticipantUIDs))
		for i, uid := range in.AllParticipantUIDs {
			hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: uid,
				Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
			}
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_PartyToParticipant{
				PartyToParticipant: &protov30.PartyToParticipant{
					Party:        in.PartyID,
					Threshold:    uint32(in.NewP2PThreshold),
					Participants: hostingParticipants,
				},
			},
		}

		_, err = deps.Client.Authorize(ctx, uint32(in.CurrentP2PSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return ProposeAddP2POutput{}, fmt.Errorf("authorizing add P2P proposal: %w", err)
		}

		deps.Logger.Infow("Add P2P proposal submitted",
			"participant", in.ParticipantID,
			"party", in.PartyID,
			"total_participants", len(in.AllParticipantUIDs),
		)

		return ProposeAddP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
		}, nil
	},
)

// ProposeRotationP2POp authorizes the updated PartyToParticipant mapping with
// the target's old DAML key replaced by the new one.
var ProposeRotationP2POp = operations.NewOperation(
	"canton-ceremony/topology/propose-rotation-p2p",
	semver.MustParse("1.0.0"),
	"Propose updated PartyToParticipant mapping with rotated DAML signing key",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeRotationP2PInput) (ProposeRotationP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.AllParticipantUIDs) == 0 {
			return ProposeRotationP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-rotation-p2p: participant_id, party_id, and all_participant_uids are required"),
			)
		}
		if in.OldDamlKeyB64 == "" || in.NewDamlKeyB64 == "" {
			return ProposeRotationP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-rotation-p2p: old_daml_key_b64 and new_daml_key_b64 are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeRotationP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ProposeRotationP2POutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		// Build hosting participants (same membership, key rotation doesn't change participants).
		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.AllParticipantUIDs))
		for i, uid := range in.AllParticipantUIDs {
			hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: uid,
				Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
			}
		}

		// Build updated signing keys by replacing the old DAML key with the new one.
		updatedKeys := make([]*cryptov30.SigningPublicKey, 0, len(in.CurrentSigningKeysB64))
		for _, skB64 := range in.CurrentSigningKeysB64 {
			if skB64 == in.OldDamlKeyB64 {
				// Replace with new key.
				newKeyBytes, decErr := base64.StdEncoding.DecodeString(in.NewDamlKeyB64)
				if decErr != nil {
					return ProposeRotationP2POutput{}, fmt.Errorf("decoding new DAML key: %w", decErr)
				}
				var newKey cryptov30.SigningPublicKey
				if unmErr := proto.Unmarshal(newKeyBytes, &newKey); unmErr != nil {
					return ProposeRotationP2POutput{}, fmt.Errorf("unmarshalling new DAML key: %w", unmErr)
				}
				updatedKeys = append(updatedKeys, &newKey)
			} else {
				// Keep existing key.
				keyBytes, decErr := base64.StdEncoding.DecodeString(skB64)
				if decErr != nil {
					return ProposeRotationP2POutput{}, fmt.Errorf("decoding existing signing key: %w", decErr)
				}
				var existingKey cryptov30.SigningPublicKey
				if unmErr := proto.Unmarshal(keyBytes, &existingKey); unmErr != nil {
					return ProposeRotationP2POutput{}, fmt.Errorf("unmarshalling existing signing key: %w", unmErr)
				}
				updatedKeys = append(updatedKeys, &existingKey)
			}
		}

		threshold := in.NewP2PThreshold
		if threshold <= 0 {
			threshold = len(in.AllParticipantUIDs)/2 + 1
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_PartyToParticipant{
				PartyToParticipant: &protov30.PartyToParticipant{
					Party:        in.PartyID,
					Threshold:    uint32(threshold),
					Participants: hostingParticipants,
					PartySigningKeys: &cryptov30.SigningKeysWithThreshold{
						Keys:      updatedKeys,
						Threshold: in.SigningKeysThreshold,
					},
				},
			},
		}

		_, err = deps.Client.Authorize(ctx, uint32(in.CurrentP2PSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return ProposeRotationP2POutput{}, fmt.Errorf("authorizing rotation P2P proposal: %w", err)
		}

		deps.Logger.Infow("Rotation P2P proposal submitted",
			"participant", in.ParticipantID,
			"party", in.PartyID,
		)

		return ProposeRotationP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
			ProposedAt:    time.Now().UTC(),
		}, nil
	},
)
