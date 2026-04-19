package keyrotation

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/internal/helpers"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// ── Step 2: GenerateRotatedKeyOp ────────────────────────────────────────────

// GenerateRotatedKeyOp generates new namespace and/or DAML signing keys for the
// target participant. When RotateDamlKey is true, it also auto-discovers the
// target's old DAML key by cross-referencing vault keys against the party's
// current signing keys.
//
// Only runs on the target participant's node (UID check enforced).
var GenerateRotatedKeyOp = operations.NewOperation(
	"keyrotation/canton-ceremony/generate-rotated-key",
	semver.MustParse("1.0.0"),
	"Generate new signing key(s) for the target participant and discover old DAML key",
	func(b operations.Bundle, deps ceremony.CantonDeps, in GenerateRotatedKeyInput) (GenerateRotatedKeyOutput, error) {
		if in.ParticipantID == "" {
			return GenerateRotatedKeyOutput{}, operations.NewUnrecoverableError(
				errors.New("generate-rotated-key: participant_id is required"),
			)
		}
		if !in.RotateNamespaceKey && !in.RotateDamlKey {
			return GenerateRotatedKeyOutput{}, operations.NewUnrecoverableError(
				errors.New("generate-rotated-key: at least one of rotate_namespace_key or rotate_daml_key must be true"),
			)
		}

		ctx := b.GetContext()

		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return GenerateRotatedKeyOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}
		if in.ParticipantID != pid {
			return GenerateRotatedKeyOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, pid,
			)
		}

		// Auto-discover the existing namespace key name from the vault.
		nsKeyName, err := deps.Client.GetNamespaceKeyName(ctx, in.SynchronizerID, in.DNSOwners)
		if err != nil {
			return GenerateRotatedKeyOutput{}, fmt.Errorf("discovering namespace key name from vault: %w", err)
		}

		out := GenerateRotatedKeyOutput{
			ParticipantID:  in.ParticipantID,
			ParticipantUID: pid,
		}

		// Generate new namespace key if requested.
		if in.RotateNamespaceKey {
			key, genErr := deps.Client.GenerateSigningKey(ctx, nsKeyName+"-rotated", []cryptov30.SigningKeyUsage{
				cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
			})
			if genErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("generating rotated namespace key: %w", genErr)
			}

			fp, fpErr := helpers.GetPublicKeyFingerprint(key.GetPublicKey())
			if fpErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("computing namespace key fingerprint: %w", fpErr)
			}

			keyBytes, mErr := proto.Marshal(key)
			if mErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("marshalling namespace key proto: %w", mErr)
			}

			out.NewNamespaceKeyB64 = base64.StdEncoding.EncodeToString(keyBytes)
			out.NewNamespaceFingerprint = fp

			deps.Logger.Infow("Rotated namespace key generated",
				"participant", in.ParticipantID,
				"new_fingerprint", fp,
			)
		}

		// Generate new DAML key if requested.
		if in.RotateDamlKey {
			damlKey, genErr := deps.Client.GenerateSigningKey(ctx, nsKeyName+"-protocol-rotated", []cryptov30.SigningKeyUsage{
				cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
			})
			if genErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("generating rotated DAML key: %w", genErr)
			}

			fp, fpErr := helpers.GetPublicKeyFingerprint(damlKey.GetPublicKey())
			if fpErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("computing DAML key fingerprint: %w", fpErr)
			}

			keyBytes, mErr := proto.Marshal(damlKey)
			if mErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("marshalling DAML key proto: %w", mErr)
			}

			out.NewDamlKeyB64 = base64.StdEncoding.EncodeToString(keyBytes)
			out.NewDamlKeyFingerprint = fp

			// Auto-discover the old DAML key by cross-referencing vault with party signing keys.
			if len(in.KnownSigningKeysB64) > 0 {
				oldFP, oldKeyB64, discoverErr := deps.Client.GetProtocolKeyFingerprint(ctx, in.KnownSigningKeysB64)
				if discoverErr != nil {
					return GenerateRotatedKeyOutput{}, fmt.Errorf("discovering old DAML key: %w", discoverErr)
				}
				out.OldDamlKeyFingerprint = oldFP
				out.OldDamlKeyB64 = oldKeyB64
			}

			deps.Logger.Infow("Rotated DAML key generated",
				"participant", in.ParticipantID,
				"new_fingerprint", fp,
				"old_fingerprint", out.OldDamlKeyFingerprint,
			)
		}

		return out, nil
	},
)

// ── Step 4: CreateRotationDNSProposalOp ─────────────────────────────────────

// CreateRotationDNSProposalOp builds the updated DecentralizedNamespaceDefinition
// with the target's old namespace fingerprint replaced by the new one, then calls
// Authorize to create the proposal. The serial is explicitly set to
// currentSerial+1 to safely update the existing mapping.
//
// Canton equivalent:
//
//	val updatedDns = DecentralizedNamespaceDefinition.tryCreate(
//	    currentNamespace, currentThreshold, currentOwners - oldNamespace + newNamespace)
//	participant.topology.decentralized_namespaces.propose(updatedDns, store = syncId)
var CreateRotationDNSProposalOp = operations.NewOperation(
	"keyrotation/canton-ceremony/create-rotation-dns-proposal",
	semver.MustParse("1.0.0"),
	"Create updated DecentralizedNamespaceDefinition proposal with rotated owner fingerprint",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateRotationDNSProposalInput) (CreateRotationDNSProposalOutput, error) {
		if in.DecentralizedNamespace == "" {
			return CreateRotationDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-rotation-dns-proposal: decentralized_namespace is required"),
			)
		}
		if in.OldNamespaceFingerprint == "" || in.NewNamespaceFingerprint == "" {
			return CreateRotationDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-rotation-dns-proposal: old and new namespace fingerprints are required"),
			)
		}

		// Build new owners list by replacing the old fingerprint with the new one.
		newOwners := make([]string, 0, len(in.CurrentOwners))
		replaced := false
		for _, owner := range in.CurrentOwners {
			if owner == in.OldNamespaceFingerprint {
				newOwners = append(newOwners, in.NewNamespaceFingerprint)
				replaced = true
			} else {
				newOwners = append(newOwners, owner)
			}
		}
		if !replaced {
			return CreateRotationDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("create-rotation-dns-proposal: old namespace fingerprint %q not found in current owners",
					in.OldNamespaceFingerprint),
			)
		}

		ctx := b.GetContext()

		threshold := in.CurrentThreshold
		if threshold <= 0 {
			threshold = len(newOwners)/2 + 1
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_DecentralizedNamespaceDefinition{
				DecentralizedNamespaceDefinition: &protov30.DecentralizedNamespaceDefinition{
					DecentralizedNamespace: in.DecentralizedNamespace,
					Threshold:              int32(threshold),
					Owners:                 newOwners,
				},
			},
		}

		tx, err := deps.Client.Authorize(ctx, uint32(in.CurrentSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return CreateRotationDNSProposalOutput{}, fmt.Errorf("authorizing rotation DNS proposal: %w", err)
		}

		txBytes, err := proto.Marshal(tx)
		if err != nil {
			return CreateRotationDNSProposalOutput{}, fmt.Errorf("marshalling rotation DNS proposal: %w", err)
		}
		txB64 := base64.StdEncoding.EncodeToString(txBytes)

		hash := sha256.Sum256(txBytes)
		proposalHash := fmt.Sprintf("%x", hash)

		deps.Logger.Infow("Rotation DNS proposal created",
			"namespace", in.DecentralizedNamespace,
			"threshold", threshold,
			"old_fingerprint", in.OldNamespaceFingerprint,
			"new_fingerprint", in.NewNamespaceFingerprint,
			"proposal_hash", proposalHash,
		)

		return CreateRotationDNSProposalOutput{
			DNSTxB64:           txB64,
			ProposalHashSHA256: proposalHash,
			NewOwners:          newOwners,
			RequiredSigners:    in.AllParticipantIDs,
		}, nil
	},
)

// ── Step 7: ProposeRotationP2POp ────────────────────────────────────────────

// ProposeRotationP2POp authorizes the updated PartyToParticipant mapping with
// the target's old DAML key replaced by the new one. Each participant runs this
// independently; Canton accumulates proposals.
//
// Canton equivalent:
//
//	participant.topology.party_to_participant_mappings.propose(
//	    partyId, participants, threshold,
//	    signingKeys = updatedSigningKeysWithThreshold, store = syncId)
var ProposeRotationP2POp = operations.NewOperation(
	"keyrotation/canton-ceremony/propose-rotation-p2p",
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
