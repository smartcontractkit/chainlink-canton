package keys

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// CreateMemberKeyOp generates a namespace signing key and a DAML (protocol)
// signing key for a single participant, and fetches the participant's UID.
//
// Canton equivalent:
//
//	val nsKey = participant.keys.secret.generate_signing_key("ns", SigningKeyUsage.NamespaceOnly)
//	val damlKey = participant.keys.secret.generate_signing_key("daml", SigningKeyUsage.ProtocolOnly)
var CreateMemberKeyOp = operations.NewOperation(
	"canton-ceremony/keys/create-member-key",
	semver.MustParse("1.0.0"),
	"Generate namespace and DAML signing keys for a ceremony participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateMemberKeyInput) (CreateMemberKeyOutput, error) {
		if in.ParticipantID == "" || in.NamespaceName == "" {
			return CreateMemberKeyOutput{}, operations.NewUnrecoverableError(
				errors.New("create-member-key: participant_id and namespace_name are required"),
			)
		}
		ctx := b.GetContext()

		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("fetching client participant ID: %w", err)
		}

		if in.ParticipantID != pid {
			return CreateMemberKeyOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s", pid, in.ParticipantID)
		}

		key, err := deps.Client.GenerateSigningKey(ctx, in.NamespaceName, []cryptov30.SigningKeyUsage{
			cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
		})
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("generating namespace key: %w", err)
		}

		// Also generate a PROTOCOL (DAML) signing key for the participant.
		damlKey, err := deps.Client.GenerateSigningKey(ctx, in.NamespaceName+"-protocol", []cryptov30.SigningKeyUsage{
			cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
		})
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("generating protocol (DAML) signing key: %w", err)
		}

		uid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		fingerprint, err := helpers.GetPublicKeyFingerprint(key.GetPublicKey())
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("computing public key fingerprint: %w", err)
		}

		damlFingerprint, err := helpers.GetPublicKeyFingerprint(damlKey.GetPublicKey())
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("computing DAML key fingerprint: %w", err)
		}

		// Preserve the complete SigningPublicKey proto (Format, Scheme, KeySpec,
		// Usage) so ProposeNamespaceDelegationOp can send it back to Canton
		// verbatim, avoiding PROTO_DESERIALIZATION_FAILURE from reconstructed keys.
		keyProtoBytes, err := proto.Marshal(key)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("marshalling signing public key proto: %w", err)
		}
		keyB64 := base64.StdEncoding.EncodeToString(keyProtoBytes)

		damlKeyProtoBytes, err := proto.Marshal(damlKey)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("marshalling DAML signing public key proto: %w", err)
		}
		damlKeyB64 := base64.StdEncoding.EncodeToString(damlKeyProtoBytes)

		return CreateMemberKeyOutput{
			ParticipantID:        in.ParticipantID,
			ParticipantUID:       uid,
			NamespaceFingerprint: fingerprint,
			SigningKeyB64:        keyB64,
			DamlKeyB64:           damlKeyB64,
			DamlKeyFingerprint:   damlFingerprint,
		}, nil
	},
)
