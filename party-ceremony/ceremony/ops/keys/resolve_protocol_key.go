package keys

import (
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// ResolveProtocolSigningKeyOp resolves the local participant's active PROTOCOL
// signing key from the current PartyToParticipant signing-key set.
var ResolveProtocolSigningKeyOp = operations.NewOperation(
	"canton-ceremony/keys/resolve-protocol-signing-key",
	semver.MustParse("1.0.0"),
	"Resolve the local participant's active protocol signing key from P2P topology",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ResolveProtocolSigningKeyInput) (ResolveProtocolSigningKeyOutput, error) {
		if in.ParticipantID == "" || len(in.KnownSigningKeysB64) == 0 {
			return ResolveProtocolSigningKeyOutput{}, operations.NewUnrecoverableError(
				errors.New("resolve-protocol-signing-key: participant_id and known_signing_keys_b64 are required"),
			)
		}

		ctx := b.GetContext()
		uid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ResolveProtocolSigningKeyOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}
		if uid != in.ParticipantID {
			return ResolveProtocolSigningKeyOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s", in.ParticipantID, uid)
		}

		fp, keyB64, err := deps.Client.GetProtocolKeyFingerprint(ctx, in.KnownSigningKeysB64)
		if err != nil {
			return ResolveProtocolSigningKeyOutput{}, fmt.Errorf("resolving protocol signing key: %w", err)
		}
		idx := slices.Index(in.KnownSigningKeysB64, keyB64)
		if idx < 0 {
			return ResolveProtocolSigningKeyOutput{}, fmt.Errorf("resolved protocol signing key is not in known signing-key set")
		}

		return ResolveProtocolSigningKeyOutput{
			ParticipantID:      in.ParticipantID,
			KeyB64:             keyB64,
			KeyFingerprint:     fp,
			KnownSigningKeyIdx: idx,
		}, nil
	},
)
