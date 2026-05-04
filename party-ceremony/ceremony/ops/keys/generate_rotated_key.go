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

// GenerateRotatedKeyOp generates or registers locally configured KMS namespace
// and/or DAML signing keys for the
// target participant. When RotateDamlKey is true, it also auto-discovers the
// target's old DAML key by cross-referencing vault keys against the party's
// current signing keys.
//
// Only runs on the target participant's node (UID check enforced).
var GenerateRotatedKeyOp = operations.NewOperation(
	"canton-ceremony/keys/generate-rotated-key",
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
			key, genErr := obtainSigningKey(ctx, deps.Client, deps.KMS.NamespaceKeyID, nsKeyName+"-rotated", []cryptov30.SigningKeyUsage{
				cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
			})
			if genErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("obtaining rotated namespace key: %w", genErr)
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
			damlKey, genErr := obtainSigningKey(ctx, deps.Client, deps.KMS.ProtocolKeyID, nsKeyName+"-protocol-rotated", []cryptov30.SigningKeyUsage{
				cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
			})
			if genErr != nil {
				return GenerateRotatedKeyOutput{}, fmt.Errorf("obtaining rotated DAML key: %w", genErr)
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
