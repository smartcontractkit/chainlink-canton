package client

import (
	"context"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
)

// LedgerClient abstracts the Canton Ledger gRPC APIs needed by the
// contract-deploy ceremony. In production this wraps real gRPC stubs;
// for testing it can be backed by a mock.
//
// This is intentionally separate from [CantonClient] which covers the
// Admin API. The two clients connect to different gRPC endpoints and
// serve different responsibilities.
type LedgerClient interface {
	// PartyExists returns true if the given party ID is known on the ledger.
	// Uses PartyManagementService.ListKnownParties with a filter.
	PartyExists(ctx context.Context, partyID string) (bool, error)

	// PrepareSubmission calls InteractiveSubmissionService.PrepareSubmission
	// to prepare a DAML transaction for multi-party signing.
	// Returns the full response including the prepared transaction hash
	// that each signer must sign.
	PrepareSubmission(
		ctx context.Context,
		commands []*apiv2.Command,
		actAs []string,
		readAs []string,
		synchronizerID string,
	) (*interactive.PrepareSubmissionResponse, error)

	// GetActiveContractsByTemplate queries the Active Contract Set for contracts
	// matching the given template (packageID:moduleName:entityName) owned by
	// the specified party. Returns all matching CreatedEvent entries.
	GetActiveContractsByTemplate(
		ctx context.Context,
		partyID string,
		packageID string,
		moduleName string,
		entityName string,
	) ([]*apiv2.CreatedEvent, error)

	// GetActiveContractsByTemplateForParty queries the ACS using a by-party
	// filter rather than the "any party" wildcard. Unlike
	// [LedgerClient.GetActiveContractsByTemplate], this only requires
	// CanReadAs rights for the given party — no participant/IDP admin
	// claims — so it works for users that have been granted per-party
	// rights via [LedgerClient.GrantPartyRights].
	//
	// packageName is the Daml package name (e.g. "test-test"); the client encodes
	// it as "#<name>" in Identifier.package_id for FiltersByParty filters.
	GetActiveContractsByTemplateForParty(
		ctx context.Context,
		partyID string,
		packageName string,
		moduleName string,
		entityName string,
	) ([]*apiv2.CreatedEvent, error)

	// ExecuteSubmission calls InteractiveSubmissionService.ExecuteSubmission
	// to submit a prepared DAML transaction with its party signatures.
	// Returns the contract ID of the first created contract in the committed
	// transaction, or an empty string if no contracts were created.
	ExecuteSubmission(
		ctx context.Context,
		preparedTx *interactive.PreparedTransaction,
		partySignatures *interactive.PartySignatures,
		hashingSchemeVersion interactive.HashingSchemeVersion,
	) (string, error)

	// GrantPartyRights grants actAs and readAs Ledger API rights for partyID
	// to the given userID. This is required in JWT-auth environments where
	// decentralized parties created via topology are not automatically visible
	// to any user. It is idempotent: granting an already-granted right is safe.
	GrantPartyRights(ctx context.Context, userID, partyID string) error
}
