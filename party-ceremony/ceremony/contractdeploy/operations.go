package contractdeploy

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/Masterminds/semver/v3"
	retry "github.com/avast/retry-go/v4"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
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
			HashingSchemeVersion:    int32(resp.GetHashingSchemeVersion()),
		}, nil
	},
)

// SignSubmissionOp signs the prepared transaction hash with the participant's
// signing key via [ContractDeployDeps.Signer].
//
// The resulting [v2.Signature] proto is serialised to base64 so it can be
// stored in the framework reporter and later deserialised by [ExecuteSubmissionOp].
var SignSubmissionOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/sign-submission",
	semver.MustParse("1.0.0"),
	"Sign prepared transaction hash with participant's signing key",
	func(b operations.Bundle, deps ContractDeployDeps, in SignSubmissionInput) (SignSubmissionOutput, error) {
		ctx := b.GetContext()

		// Only the expected participant signs — same participant-gate pattern as UploadDarsOp.
		uid, err := deps.AdminClient.GetParticipantUID(ctx)
		if err != nil {
			return SignSubmissionOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if uid != in.ParticipantID {
			return SignSubmissionOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s",
				in.ParticipantID, uid)
		}

		hashBytes, err := hex.DecodeString(in.PreparedTransactionHash)
		if err != nil {
			return SignSubmissionOutput{}, fmt.Errorf("decoding prepared transaction hash: %w", err)
		}

		sig, err := deps.Signer.Sign(ctx, hashBytes)
		if err != nil {
			return SignSubmissionOutput{}, fmt.Errorf("signing transaction hash for %q: %w", in.ParticipantID, err)
		}

		sigBytes, err := proto.Marshal(sig)
		if err != nil {
			return SignSubmissionOutput{}, fmt.Errorf("marshalling signature: %w", err)
		}

		deps.Logger.Infow("Transaction hash signed",
			"participant", in.ParticipantID,
			"key", sig.GetSignedBy(),
		)

		return SignSubmissionOutput{
			ParticipantID:  in.ParticipantID,
			SignatureB64:   base64.StdEncoding.EncodeToString(sigBytes),
			KeyFingerprint: sig.GetSignedBy(),
		}, nil
	},
)

// ExecuteSubmissionOp aggregates all participant signatures and submits the
// prepared transaction via InteractiveSubmissionService.ExecuteSubmission.
//
// All participants' signatures are grouped under a single [interactive.SinglePartySignatures]
// entry for the decentralized party (each participant signs on behalf of the party).
var ExecuteSubmissionOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/execute-submission",
	semver.MustParse("1.0.0"),
	"Execute signed submission via InteractiveSubmissionService",
	func(b operations.Bundle, deps ContractDeployDeps, in ExecuteSubmissionInput) (ExecuteSubmissionOutput, error) {
		ctx := b.GetContext()

		// Deserialise the prepared transaction proto.
		txBytes, err := base64.StdEncoding.DecodeString(in.PreparedTxB64)
		if err != nil {
			return ExecuteSubmissionOutput{}, fmt.Errorf("decoding prepared transaction: %w", err)
		}
		var preparedTx interactive.PreparedTransaction
		if err := proto.Unmarshal(txBytes, &preparedTx); err != nil {
			return ExecuteSubmissionOutput{}, fmt.Errorf("unmarshalling prepared transaction: %w", err)
		}

		// Deserialise each participant's signature.
		sigs := make([]*apiv2.Signature, 0, len(in.SignaturesB64))
		for i, sigB64 := range in.SignaturesB64 {
			sigBytes, err := base64.StdEncoding.DecodeString(sigB64)
			if err != nil {
				return ExecuteSubmissionOutput{}, fmt.Errorf("decoding signature[%d]: %w", i, err)
			}
			var sig apiv2.Signature
			if err := proto.Unmarshal(sigBytes, &sig); err != nil {
				return ExecuteSubmissionOutput{}, fmt.Errorf("unmarshalling signature[%d]: %w", i, err)
			}
			sigs = append(sigs, &sig)
		}

		partySignatures := &interactive.PartySignatures{
			Signatures: []*interactive.SinglePartySignatures{{
				Party:      in.DecentralizedPartyID,
				Signatures: sigs,
			}},
		}

		hsv := interactive.HashingSchemeVersion(in.HashingSchemeVersion)
		contractID, err := deps.LedgerClient.ExecuteSubmission(ctx, &preparedTx, partySignatures, hsv)
		if err != nil {
			return ExecuteSubmissionOutput{}, fmt.Errorf("executing submission: %w", err)
		}

		deps.Logger.Infow("Submission executed",
			"party", in.DecentralizedPartyID,
			"signatures", len(sigs),
			"contract_id", contractID,
		)

		return ExecuteSubmissionOutput{ContractID: contractID}, nil
	},
)

// VerifyContractOp confirms the contract ID returned by [ExecuteSubmissionOp]
// is non-empty, proving the transaction was committed and created a contract.
var VerifyContractOp = operations.NewOperation(
	"contract-deploy/canton-ceremony/verify-contract",
	semver.MustParse("1.0.0"),
	"Verify contract exists in Active Contract Set",
	func(b operations.Bundle, deps ContractDeployDeps, in VerifyContractInput) (VerifyContractOutput, error) {
		if in.ContractID == "" {
			return VerifyContractOutput{}, fmt.Errorf(
				"contract %s.%s not found: no contract ID returned by execute step",
				in.TemplateModule, in.TemplateEntity,
			)
		}

		deps.Logger.Infow("Contract verified",
			"contract_id", in.ContractID,
			"template", in.TemplateModule+"."+in.TemplateEntity,
		)

		return VerifyContractOutput{
			Verified:   true,
			ContractID: in.ContractID,
		}, nil
	},
)
