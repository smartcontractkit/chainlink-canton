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
	Sign(ctx context.Context, hash []byte) (*v2crypto.Signature, error)
}
