package client

import (
	"context"

	v2crypto "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
)

// TransactionSigner signs a prepared DAML transaction hash and returns the
// corresponding Ledger API [v2crypto.Signature] proto.
//
// Implementations:
//   - [VaultSigner]  — exports keys from Canton's VaultService and signs with Ed25519
//   - Future: KMSSigner (AWS/GCP KMS), LedgerSigner (hardware), etc.
//
// The returned [v2crypto.Signature] is consumed directly by
// [ExecuteSubmissionRequest.PartySignatures].
type TransactionSigner interface {
	// Sign signs the given hash bytes and returns a Ledger API Signature proto.
	// hash is the raw bytes of PrepareSubmissionResponse.PreparedTransactionHash.
	// For ECDSA signers this value is still treated as the message to sign;
	// the algorithm-specific SHA-256/SHA-384 hashing is applied by the signer.
	Sign(ctx context.Context, hash []byte) (*v2crypto.Signature, error)
}

// ProtocolKeyResolver resolves the local participant's active PROTOCOL signing
// key from the current PartyToParticipant signing-key set.
type ProtocolKeyResolver interface {
	GetProtocolKeyFingerprint(ctx context.Context, knownSigningKeys []string) (fingerprint string, keyB64 string, err error)
}

// TransactionSignerFactory creates the signer for the connected participant
// once the current PartyToParticipant signing keys are known.
type TransactionSignerFactory func(ctx context.Context, participantID string, knownSigningKeysB64 []string) (TransactionSigner, error)
