package contractdeploy

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// UploadDarsOp uploads all configured DARs to a single participant.
// Each participant runs this operation independently. The operation is
// idempotent: re-uploading the same DAR returns the same package ID.
var UploadDarsOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/upload-dars",
	semver.MustParse("1.0.0"),
	"Upload DAR files to a participant via PackageService",
	func(b operations.Bundle, deps ContractDeployDeps, in UploadDarsInput) (UploadDarsOutput, error) {
		ctx := b.GetContext()

		// Only the participant that owns this node can upload DARs for it.
		uid, err := deps.AdminClient.GetParticipantUID(ctx)
		if err != nil {
			return UploadDarsOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}
		if uid != in.ParticipantID {
			return UploadDarsOutput{}, retry.Unrecoverable(
				fmt.Errorf("participant mismatch: connected to %q, expected %q", uid, in.ParticipantID))
		}

		var packageIDs []string
		for _, pkg := range in.Packages {
			darBytes, err := deps.DARLoader(pkg.Name, pkg.Version)
			if err != nil {
				return UploadDarsOutput{}, fmt.Errorf("loading DAR %q@%s: %w", pkg.Name, pkg.Version, err)
			}

			pkgID, err := deps.AdminClient.UploadDar(ctx, darBytes)
			if err != nil {
				return UploadDarsOutput{}, fmt.Errorf("uploading DAR %q@%s: %w", pkg.Name, pkg.Version, err)
			}

			deps.Logger.Infow("DAR uploaded",
				"participant", in.ParticipantID,
				"package", pkg.Name,
				"version", pkg.Version,
				"package_id", pkgID,
			)
			packageIDs = append(packageIDs, pkgID)
		}

		return UploadDarsOutput{
			ParticipantID: in.ParticipantID,
			PackageIDs:    packageIDs,
		}, nil
	},
)

// VerifyPartyOp checks that the decentralized party is visible on the Ledger API.
// This confirms that the onboarding ceremony completed successfully and the
// party can be used for contract transactions.
var VerifyPartyOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/verify-party",
	semver.MustParse("1.0.0"),
	"Verify decentralized party exists on the Ledger API",
	func(b operations.Bundle, deps ContractDeployDeps, in VerifyPartyInput) (VerifyPartyOutput, error) {
		ctx := b.GetContext()

		exists, err := deps.LedgerClient.PartyExists(ctx, in.DecentralizedPartyID)
		if err != nil {
			return VerifyPartyOutput{}, fmt.Errorf("checking party %q: %w", in.DecentralizedPartyID, err)
		}

		if !exists {
			return VerifyPartyOutput{}, fmt.Errorf("party %q not found on ledger; ensure onboarding is complete", in.DecentralizedPartyID)
		}

		deps.Logger.Infow("Party verified on ledger", "party", in.DecentralizedPartyID)

		return VerifyPartyOutput{
			PartyID:  in.DecentralizedPartyID,
			Verified: true,
		}, nil
	},
)

// PrepareSubmissionOp prepares a contract creation transaction via the
// InteractiveSubmissionService. The coordinator builds a CreateCommand from
// the contract arguments and calls PrepareSubmission. The returned hash is
// what each signer must sign in the next step.
var PrepareSubmissionOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/prepare-submission",
	semver.MustParse("1.0.0"),
	"Prepare contract creation via InteractiveSubmissionService",
	func(b operations.Bundle, deps ContractDeployDeps, in PrepareSubmissionInput) (PrepareSubmissionOutput, error) {
		ctx := b.GetContext()

		// Parse contract arguments from proto JSON into a DAML Record.
		var record apiv2.Record
		if err := protojson.Unmarshal([]byte(in.ContractArgs), &record); err != nil {
			return PrepareSubmissionOutput{}, fmt.Errorf("parsing contract args: %w", err)
		}

		commands := []*apiv2.Command{{
			Command: &apiv2.Command_Create{
				Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  in.PackageID,
						ModuleName: in.TemplateModule,
						EntityName: in.TemplateEntity,
					},
					CreateArguments: &record,
				},
			},
		}}

		resp, err := deps.LedgerClient.PrepareSubmission(
			ctx,
			commands,
			[]string{in.DecentralizedPartyID},
			[]string{in.DecentralizedPartyID},
			in.SynchronizerID,
		)
		if err != nil {
			return PrepareSubmissionOutput{}, fmt.Errorf("preparing submission: %w", err)
		}

		hashHex := hex.EncodeToString(resp.GetPreparedTransactionHash())

		// Serialize the prepared transaction for distribution to signers.
		txBytes, err := proto.Marshal(resp.GetPreparedTransaction())
		if err != nil {
			return PrepareSubmissionOutput{}, fmt.Errorf("marshalling prepared transaction: %w", err)
		}
		txB64 := base64.StdEncoding.EncodeToString(txBytes)

		deps.Logger.Infow("Submission prepared",
			"hash", hashHex,
			"template", in.TemplateModule+":"+in.TemplateEntity,
		)

		return PrepareSubmissionOutput{
			PreparedTransactionHash: hashHex,
			PreparedTxB64:           txB64,
		}, nil
	},
)

// SignSubmissionOp is a placeholder for the signing step.
// Each participant would sign the prepared transaction hash with their DAML key.
//
// TODO: Implement once a signing mechanism is available. Options:
//   - VaultService.ExportKeyPair + external Ed25519 signing
//   - Future Canton API for vault-based transaction signing
var SignSubmissionOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/sign-submission",
	semver.MustParse("1.0.0"),
	"Sign prepared transaction hash with participant's DAML key (not yet implemented)",
	func(_ operations.Bundle, _ ContractDeployDeps, _ SignSubmissionInput) (SignSubmissionOutput, error) {
		return SignSubmissionOutput{}, ErrSigningNotImplemented
	},
)

// ExecuteSubmissionOp is a placeholder for the execution step.
// The coordinator would aggregate signatures and submit.
//
// TODO: Implement once signing is available.
var ExecuteSubmissionOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/execute-submission",
	semver.MustParse("1.0.0"),
	"Execute signed submission via InteractiveSubmissionService (not yet implemented)",
	func(_ operations.Bundle, _ ContractDeployDeps, _ ExecuteSubmissionInput) (ExecuteSubmissionOutput, error) {
		return ExecuteSubmissionOutput{}, ErrSigningNotImplemented
	},
)

// VerifyContractOp is a placeholder for the verification step.
// It would check that the contract appears in the Active Contract Set.
//
// TODO: Implement once signing and execution are available.
var VerifyContractOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/verify-contract",
	semver.MustParse("1.0.0"),
	"Verify contract exists in Active Contract Set (not yet implemented)",
	func(_ operations.Bundle, _ ContractDeployDeps, _ VerifyContractInput) (VerifyContractOutput, error) {
		return VerifyContractOutput{}, ErrSigningNotImplemented
	},
)
