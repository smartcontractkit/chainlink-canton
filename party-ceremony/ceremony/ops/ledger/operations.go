package ledger

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/avast/retry-go/v4"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// GrantPartyRightsOp grants actAs and readAs Ledger API rights for the
// decentralized party to the configured user. In no-auth environments
// (deps.UserID is empty) the operation is a no-op and returns immediately.
// In JWT-auth environments, only the participant whose UID matches the
// connected node will execute the grant; all others return an error so
// the caller can track per-participant completion.
var GrantPartyRightsOp = operations.NewOperation(
	"canton-ceremony/ledger/grant-party-rights",
	semver.MustParse("1.0.0"),
	"Grant actAs and readAs Ledger API rights for the decentralized party to the configured user",
	grantRights("actAs and readAs", client.LedgerClient.GrantPartyRights),
)

var GrantReadAsOp = operations.NewOperation(
	"canton-ceremony/ledger/grant-read-as",
	semver.MustParse("1.0.0"),
	"Grant readAs Ledger API rights for the decentralized party to the configured user",
	grantRights("readAs", client.LedgerClient.GrantReadAs),
)

var GrantActAsOp = operations.NewOperation(
	"canton-ceremony/ledger/grant-act-as",
	semver.MustParse("1.0.0"),
	"Grant actAs Ledger API rights for the decentralized party to the configured user",
	grantRights("actAs", client.LedgerClient.GrantActAs),
)

func grantRights(
	rights string,
	grant func(client.LedgerClient, context.Context, string, string) error,
) func(operations.Bundle, ContractDeployDeps, GrantPartyRightsInput) (GrantPartyRightsOutput, error) {
	return func(b operations.Bundle, deps ContractDeployDeps, in GrantPartyRightsInput) (GrantPartyRightsOutput, error) {
		if deps.UserID == "" {
			return GrantPartyRightsOutput{ParticipantID: in.ParticipantID, Granted: false}, nil
		}

		ctx := b.GetContext()

		uid, err := deps.AdminClient.GetParticipantUID(ctx)
		if err != nil {
			return GrantPartyRightsOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if uid != in.ParticipantID {
			return GrantPartyRightsOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s",
				in.ParticipantID, uid)
		}

		if err := grant(deps.LedgerClient, ctx, deps.UserID, in.DecentralizedPartyID); err != nil {
			return GrantPartyRightsOutput{}, fmt.Errorf("granting %s for user %q: %w", rights, deps.UserID, err)
		}

		deps.Logger.Infow("Party rights granted",
			"rights", rights,
			"participant", in.ParticipantID,
			"user", deps.UserID,
			"party", in.DecentralizedPartyID,
		)

		return GrantPartyRightsOutput{
			ParticipantID: in.ParticipantID,
			UserID:        deps.UserID,
			Granted:       true,
		}, nil
	}
}

// FetchParticipantsOp queries the Canton topology store for the participant UIDs
// that currently host the decentralized party.
var FetchParticipantsOp = operations.NewOperation(
	"canton-ceremony/ledger/fetch-participants",
	semver.MustParse("1.0.0"),
	"Fetch participant UIDs for the decentralized party from the topology store",
	func(b operations.Bundle, deps ContractDeployDeps, in FetchParticipantsInput) (FetchParticipantsOutput, error) {
		ctx := b.GetContext()

		p2pState, err := deps.AdminClient.GetP2P(ctx, in.DecentralizedPartyID, in.SynchronizerID)
		if err != nil {
			return FetchParticipantsOutput{}, fmt.Errorf("fetching P2P state for party %q: %w", in.DecentralizedPartyID, err)
		}

		if len(p2pState.Participants) == 0 {
			return FetchParticipantsOutput{}, fmt.Errorf("no participants found for party %q: ensure the onboarding ceremony completed successfully", in.DecentralizedPartyID)
		}

		participants := make([]string, len(p2pState.Participants))
		for i, p := range p2pState.Participants {
			participants[i] = p.ParticipantUID
		}

		deps.Logger.Infow("Fetched participants for decentralized party",
			"party", in.DecentralizedPartyID,
			"count", len(participants),
		)

		out := FetchParticipantsOutput{Participants: participants}
		if p2pState.PartySigningKeys != nil {
			out.PartySigningKeysB64 = p2pState.PartySigningKeys.Keys
			out.PartySigningThreshold = p2pState.PartySigningKeys.Threshold
		}

		return out, nil
	},
)

// UploadDarsOp uploads all configured DARs to a single participant.
// Each participant runs this operation independently. The operation is
// idempotent: re-uploading the same DAR returns the same package ID.
var UploadDarsOp = operations.NewOperation(
	"canton-ceremony/ledger/upload-dars",
	semver.MustParse("1.0.0"),
	"Upload DAR files to a participant via PackageService",
	func(b operations.Bundle, deps ContractDeployDeps, in UploadDarsInput) (UploadDarsOutput, error) {
		ctx := b.GetContext()

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
//
// After the P2P mapping is confirmed in the topology state, the party may take
// a moment to become visible on the deploy runner's Ledger API (each
// participant materializes the party asynchronously). We therefore poll until
// the party shows up, using the same bounds as the P2P confirmation loop in
// the onboarding ceremony.
var VerifyPartyOp = operations.NewOperation(
	"canton-ceremony/ledger/verify-party",
	semver.MustParse("1.0.0"),
	"Verify decentralized party exists on the Ledger API",
	func(b operations.Bundle, deps ContractDeployDeps, in VerifyPartyInput) (VerifyPartyOutput, error) {
		ctx := b.GetContext()

		err := retry.Do(
			func() error {
				exists, qErr := deps.LedgerClient.PartyExists(ctx, in.DecentralizedPartyID)
				if qErr != nil {
					return retry.Unrecoverable(fmt.Errorf("checking party %q: %w", in.DecentralizedPartyID, qErr))
				}
				if !exists {
					return fmt.Errorf("party %q not yet visible on ledger", in.DecentralizedPartyID)
				}

				return nil
			},
			retry.Context(ctx),
			retry.Attempts(30),
			retry.Delay(500*time.Millisecond),
		)
		if err != nil {
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
// InteractiveSubmissionService.
var PrepareSubmissionOp = operations.NewOperation(
	"canton-ceremony/ledger/prepare-submission",
	semver.MustParse("1.0.0"),
	"Prepare contract creation via InteractiveSubmissionService",
	func(b operations.Bundle, deps ContractDeployDeps, in PrepareSubmissionInput) (PrepareSubmissionOutput, error) {
		var record apiv2.Record
		if err := protojson.Unmarshal([]byte(in.ContractArgs), &record); err != nil {
			return PrepareSubmissionOutput{}, fmt.Errorf("parsing contract args: %w", err)
		}

		command := &apiv2.Command{
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
		}

		return prepare(b, deps, command, in.DecentralizedPartyID, in.SynchronizerID, nil, in.TemplateModule+":"+in.TemplateEntity)
	},
)

// PrepareExerciseOp prepares the exercise of a choice on an existing contract
// via the InteractiveSubmissionService.
var PrepareExerciseOp = operations.NewOperation(
	"canton-ceremony/ledger/prepare-exercise",
	semver.MustParse("1.0.0"),
	"Prepare a choice exercise via InteractiveSubmissionService",
	func(b operations.Bundle, deps ContractDeployDeps, in PrepareExerciseInput) (PrepareSubmissionOutput, error) {
		var arg apiv2.Value
		if err := protojson.Unmarshal([]byte(in.ChoiceArgument), &arg); err != nil {
			return PrepareSubmissionOutput{}, fmt.Errorf("parsing choice argument: %w", err)
		}

		disclosed := make([]*apiv2.DisclosedContract, 0, len(in.DisclosedContracts))
		for i, d := range in.DisclosedContracts {
			dc, err := d.Proto()
			if err != nil {
				return PrepareSubmissionOutput{}, fmt.Errorf("disclosed contract %d: %w", i, err)
			}
			disclosed = append(disclosed, dc)
		}

		command := &apiv2.Command{
			Command: &apiv2.Command_Exercise{
				Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  in.PackageID,
						ModuleName: in.TemplateModule,
						EntityName: in.TemplateEntity,
					},
					ContractId:     in.ContractID,
					Choice:         in.Choice,
					ChoiceArgument: &arg,
				},
			},
		}

		return prepare(b, deps, command, in.DecentralizedPartyID, in.SynchronizerID, disclosed, in.TemplateModule+":"+in.TemplateEntity+"."+in.Choice)
	},
)

func prepare(
	b operations.Bundle,
	deps ContractDeployDeps,
	command *apiv2.Command,
	partyID string,
	synchronizerID string,
	disclosed []*apiv2.DisclosedContract,
	action string,
) (PrepareSubmissionOutput, error) {
	resp, err := deps.LedgerClient.PrepareSubmission(
		b.GetContext(),
		[]*apiv2.Command{command},
		[]string{partyID},
		[]string{partyID},
		synchronizerID,
		disclosed,
	)
	if err != nil {
		return PrepareSubmissionOutput{}, fmt.Errorf("preparing submission: %w", err)
	}

	hashHex := hex.EncodeToString(resp.GetPreparedTransactionHash())

	txBytes, err := proto.Marshal(resp.GetPreparedTransaction())
	if err != nil {
		return PrepareSubmissionOutput{}, fmt.Errorf("marshalling prepared transaction: %w", err)
	}

	deps.Logger.Infow("Submission prepared",
		"hash", hashHex,
		"action", action,
	)

	return PrepareSubmissionOutput{
		PreparedTransactionHash: hashHex,
		PreparedTxB64:           base64.StdEncoding.EncodeToString(txBytes),
		HashingSchemeVersion:    int32(resp.GetHashingSchemeVersion()),
	}, nil
}

// SignSubmissionOp signs the prepared transaction hash with the participant's
// signing key via [ContractDeployDeps.Signer].
var SignSubmissionOp = operations.NewOperation(
	"canton-ceremony/ledger/sign-submission",
	semver.MustParse("1.0.0"),
	"Sign prepared transaction hash with participant's signing key",
	func(b operations.Bundle, deps ContractDeployDeps, in SignSubmissionInput) (SignSubmissionOutput, error) {
		ctx := b.GetContext()

		uid, err := deps.AdminClient.GetParticipantUID(ctx)
		if err != nil {
			return SignSubmissionOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if uid != in.ParticipantID {
			return SignSubmissionOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s",
				in.ParticipantID, uid)
		}

		if in.Expected != nil {
			if err := in.Expected.Check(in.PreparedTxB64); err != nil {
				return SignSubmissionOutput{}, operations.NewUnrecoverableError(err)
			}
		}
		hashBytes, err := verifyPreparedTxHash(in)
		if err != nil {
			return SignSubmissionOutput{}, operations.NewUnrecoverableError(err)
		}

		if deps.Confirmer != nil {
			detail := ceremony.DAMLSignDetail{
				TransactionHash: in.PreparedTransactionHash,
				SignerIdentity:  in.ParticipantID,
			}
			if in.Expected != nil {
				detail.Template = in.Expected.TemplateModule + ":" + in.Expected.TemplateEntity
				detail.ContractID = in.Expected.ContractID
				detail.Choice = in.Expected.Choice
				detail.ChoiceArgument = in.Expected.ChoiceArgument
			}
			if cErr := deps.Confirmer.ConfirmDAMLSign(ctx, detail); cErr != nil {
				return SignSubmissionOutput{}, operations.NewUnrecoverableError(cErr)
			}
		}

		signer := deps.Signer
		if signer == nil {
			if deps.SignerFactory == nil {
				return SignSubmissionOutput{}, fmt.Errorf("no transaction signer configured")
			}
			signer, err = deps.SignerFactory(ctx, in.ParticipantID, in.KnownSigningKeysB64)
			if err != nil {
				return SignSubmissionOutput{}, fmt.Errorf("creating transaction signer for %q: %w", in.ParticipantID, err)
			}
		}

		sig, err := signer.Sign(ctx, hashBytes)
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
var ExecuteSubmissionOp = operations.NewOperation(
	"canton-ceremony/ledger/execute-submission",
	semver.MustParse("1.0.0"),
	"Execute signed submission via InteractiveSubmissionService",
	func(b operations.Bundle, deps ContractDeployDeps, in ExecuteSubmissionInput) (ExecuteSubmissionOutput, error) {
		ctx := b.GetContext()

		txBytes, err := base64.StdEncoding.DecodeString(in.PreparedTxB64)
		if err != nil {
			return ExecuteSubmissionOutput{}, fmt.Errorf("decoding prepared transaction: %w", err)
		}
		var preparedTx interactive.PreparedTransaction
		if err := proto.Unmarshal(txBytes, &preparedTx); err != nil {
			return ExecuteSubmissionOutput{}, fmt.Errorf("unmarshalling prepared transaction: %w", err)
		}

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
		result, err := deps.LedgerClient.ExecuteSubmission(ctx, &preparedTx, partySignatures, hsv)
		if err != nil {
			return ExecuteSubmissionOutput{}, fmt.Errorf("executing submission: %w", err)
		}

		deps.Logger.Infow("Submission executed",
			"party", in.DecentralizedPartyID,
			"signatures", len(sigs),
			"update_id", result.UpdateID,
			"created", len(result.Created),
		)

		return ExecuteSubmissionOutput{
			ContractID:       result.FirstContractID(),
			UpdateID:         result.UpdateID,
			CreatedContracts: result.Created,
		}, nil
	},
)

// VerifyContractOp confirms the contract ID returned by [ExecuteSubmissionOp]
// is non-empty.
var VerifyContractOp = operations.NewOperation(
	"canton-ceremony/ledger/verify-contract",
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

// ErrPreparedTxMismatch is returned when the prepared transaction is not the expected exercise.
var ErrPreparedTxMismatch = errors.New("prepared transaction does not match the workflow")

// ErrPreparedTxHashMismatch is returned when the hash to sign is not the hash of the prepared transaction.
var ErrPreparedTxHashMismatch = errors.New("prepared transaction hash does not match the transaction")

// verifyPreparedTxHash recomputes the hash from the prepared transaction so a
// signer never signs a hash it was only told about.
func verifyPreparedTxHash(in SignSubmissionInput) ([]byte, error) {
	if v := interactive.HashingSchemeVersion(in.HashingSchemeVersion); v != interactive.HashingSchemeVersion_HASHING_SCHEME_VERSION_V2 {
		return nil, fmt.Errorf("%w: hashing scheme %s is not supported", ErrPreparedTxHashMismatch, v)
	}
	preparedTx, err := decodePreparedTx(in.PreparedTxB64)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPreparedTxHashMismatch, err)
	}
	computed, err := HashPreparedTransactionV2(preparedTx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPreparedTxHashMismatch, err)
	}
	given, err := hex.DecodeString(in.PreparedTransactionHash)
	if err != nil {
		return nil, fmt.Errorf("%w: decoding prepared transaction hash: %w", ErrPreparedTxHashMismatch, err)
	}
	if !bytes.Equal(computed, given) {
		return nil, fmt.Errorf("%w: computed %x, given %x", ErrPreparedTxHashMismatch, computed, given)
	}

	return computed, nil
}

func decodePreparedTx(preparedTxB64 string) (*interactive.PreparedTransaction, error) {
	txBytes, err := base64.StdEncoding.DecodeString(preparedTxB64)
	if err != nil {
		return nil, fmt.Errorf("decoding prepared transaction: %w", err)
	}
	var preparedTx interactive.PreparedTransaction
	if err := proto.Unmarshal(txBytes, &preparedTx); err != nil {
		return nil, fmt.Errorf("unmarshalling prepared transaction: %w", err)
	}

	return &preparedTx, nil
}

// Check returns an error wrapping [ErrPreparedTxMismatch] unless preparedTxB64 is exactly the expected exercise.
func (e *ExpectedExercise) Check(preparedTxB64 string) error {
	preparedTx, err := decodePreparedTx(preparedTxB64)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPreparedTxMismatch, err)
	}

	tx := preparedTx.GetTransaction()
	roots := tx.GetRoots()
	if len(roots) != 1 {
		return fmt.Errorf("%w: %d root nodes, want 1", ErrPreparedTxMismatch, len(roots))
	}
	var exercise *v1.Exercise
	for _, n := range tx.GetNodes() {
		if n.GetNodeId() == roots[0] {
			exercise = n.GetV1().GetExercise()
		}
	}
	if exercise == nil {
		return fmt.Errorf("%w: root node is not an exercise", ErrPreparedTxMismatch)
	}

	tid, iid := exercise.GetTemplateId(), exercise.GetInterfaceId()
	tplMatch := tid.GetModuleName() == e.TemplateModule && tid.GetEntityName() == e.TemplateEntity
	ifaceMatch := iid.GetModuleName() == e.TemplateModule && iid.GetEntityName() == e.TemplateEntity
	if !tplMatch && !ifaceMatch {
		return fmt.Errorf("%w: template %s:%s, want %s:%s", ErrPreparedTxMismatch,
			tid.GetModuleName(), tid.GetEntityName(), e.TemplateModule, e.TemplateEntity)
	}
	matched := tid
	if !tplMatch {
		matched = iid
	}
	if name, isRef := strings.CutPrefix(e.PackageID, "#"); isRef && exercise.GetPackageName() != name {
		return fmt.Errorf("%w: package name %s, want %s", ErrPreparedTxMismatch, exercise.GetPackageName(), name)
	} else if !isRef && matched.GetPackageId() != e.PackageID {
		return fmt.Errorf("%w: package %s, want %s", ErrPreparedTxMismatch, matched.GetPackageId(), e.PackageID)
	}
	if exercise.GetContractId() != e.ContractID {
		return fmt.Errorf("%w: contract %s, want %s", ErrPreparedTxMismatch, exercise.GetContractId(), e.ContractID)
	}
	if exercise.GetChoiceId() != e.Choice {
		return fmt.Errorf("%w: choice %s, want %s", ErrPreparedTxMismatch, exercise.GetChoiceId(), e.Choice)
	}
	var wantArg apiv2.Value
	if err := protojson.Unmarshal([]byte(e.ChoiceArgument), &wantArg); err != nil {
		return fmt.Errorf("%w: parsing choice argument: %w", ErrPreparedTxMismatch, err)
	}
	gotArg := proto.Clone(exercise.GetChosenValue()).(*apiv2.Value)
	dropPairedLabelledNones(&wantArg, gotArg)
	alignLabelledFields(&wantArg, gotArg)
	if !proto.Equal(normalizeValue(&wantArg), normalizeValue(gotArg)) {
		return fmt.Errorf("%w: choice argument %s, want %s", ErrPreparedTxMismatch,
			protojson.Format(exercise.GetChosenValue()), e.ChoiceArgument)
	}

	return nil
}

// dropPairedLabelledNones drops None fields only where both records are labelled, so positional fields stay aligned.
func dropPairedLabelledNones(want, got *apiv2.Value) {
	isNoneField := func(f *apiv2.RecordField) bool { return isNone(f.GetValue()) }
	switch w := want.GetSum().(type) {
	case *apiv2.Value_Record:
		g := got.GetRecord()
		if g == nil {
			return
		}
		if allLabelled(w.Record) && allLabelled(g) {
			w.Record.Fields = slices.DeleteFunc(w.Record.Fields, isNoneField)
			g.Fields = slices.DeleteFunc(g.Fields, isNoneField)
			byLabel := make(map[string]*apiv2.RecordField, len(g.GetFields()))
			for _, f := range g.GetFields() {
				byLabel[f.GetLabel()] = f
			}
			for _, f := range w.Record.GetFields() {
				if gf, ok := byLabel[f.GetLabel()]; ok {
					dropPairedLabelledNones(f.GetValue(), gf.GetValue())
				}
			}

			return
		}
		for i := range min(len(w.Record.GetFields()), len(g.GetFields())) {
			dropPairedLabelledNones(w.Record.GetFields()[i].GetValue(), g.GetFields()[i].GetValue())
		}
	case *apiv2.Value_Variant:
		dropPairedLabelledNones(w.Variant.GetValue(), got.GetVariant().GetValue())
	case *apiv2.Value_Optional:
		dropPairedLabelledNones(w.Optional.GetValue(), got.GetOptional().GetValue())
	case *apiv2.Value_List:
		gl := got.GetList().GetElements()
		for i := range min(len(w.List.GetElements()), len(gl)) {
			dropPairedLabelledNones(w.List.GetElements()[i], gl[i])
		}
	case *apiv2.Value_TextMap:
		forEachTextMapPair(w.TextMap, got.GetTextMap(), dropPairedLabelledNones)
	case *apiv2.Value_GenMap:
		ge := got.GetGenMap().GetEntries()
		for i := range min(len(w.GenMap.GetEntries()), len(ge)) {
			dropPairedLabelledNones(w.GenMap.GetEntries()[i].GetKey(), ge[i].GetKey())
			dropPairedLabelledNones(w.GenMap.GetEntries()[i].GetValue(), ge[i].GetValue())
		}
	}
}

// normalizeValue drops what the ledger may add or strip: type IDs, labels, trailing Nones, numeric padding.
func normalizeValue(v *apiv2.Value) *apiv2.Value {
	if v == nil {
		return nil
	}
	out := proto.Clone(v).(*apiv2.Value)
	normalizeInPlace(out)

	return out
}

func normalizeInPlace(v *apiv2.Value) {
	switch s := v.GetSum().(type) {
	case *apiv2.Value_Numeric:
		if strings.Contains(s.Numeric, ".") {
			s.Numeric = strings.TrimSuffix(strings.TrimRight(s.Numeric, "0"), ".")
		}
	case *apiv2.Value_Record:
		s.Record.RecordId = nil
		fields := s.Record.GetFields()
		for len(fields) > 0 && isNone(fields[len(fields)-1].GetValue()) {
			fields = fields[:len(fields)-1]
		}
		s.Record.Fields = fields
		for _, f := range fields {
			f.Label = ""
			normalizeInPlace(f.GetValue())
		}
	case *apiv2.Value_Variant:
		s.Variant.VariantId = nil
		normalizeInPlace(s.Variant.GetValue())
	case *apiv2.Value_Enum:
		s.Enum.EnumId = nil
	case *apiv2.Value_Optional:
		if s.Optional.GetValue() != nil {
			normalizeInPlace(s.Optional.GetValue())
		}
	case *apiv2.Value_List:
		for _, e := range s.List.GetElements() {
			normalizeInPlace(e)
		}
	case *apiv2.Value_TextMap:
		slices.SortFunc(s.TextMap.Entries, func(a, b *apiv2.TextMap_Entry) int {
			return strings.Compare(a.GetKey(), b.GetKey())
		})
		for _, e := range s.TextMap.GetEntries() {
			normalizeInPlace(e.GetValue())
		}
	case *apiv2.Value_GenMap:
		for _, e := range s.GenMap.GetEntries() {
			normalizeInPlace(e.GetKey())
			normalizeInPlace(e.GetValue())
		}
	}
}

// alignLabelledFields puts got's labelled record fields in want's label order: the ledger accepts labelled fields in any order.
func alignLabelledFields(want, got *apiv2.Value) {
	switch w := want.GetSum().(type) {
	case *apiv2.Value_Record:
		g := got.GetRecord()
		if g == nil {
			return
		}
		if allLabelled(w.Record) && allLabelled(g) {
			byLabel := make(map[string]*apiv2.RecordField, len(g.GetFields()))
			for _, f := range g.GetFields() {
				byLabel[f.GetLabel()] = f
			}
			ordered := make([]*apiv2.RecordField, 0, len(g.GetFields()))
			for _, f := range w.Record.GetFields() {
				if gf, ok := byLabel[f.GetLabel()]; ok {
					ordered = append(ordered, gf)
					delete(byLabel, f.GetLabel())
				}
			}
			for _, f := range g.GetFields() {
				if _, ok := byLabel[f.GetLabel()]; ok {
					ordered = append(ordered, f)
				}
			}
			g.Fields = ordered
		}
		for i := range min(len(w.Record.GetFields()), len(g.GetFields())) {
			alignLabelledFields(w.Record.GetFields()[i].GetValue(), g.GetFields()[i].GetValue())
		}
	case *apiv2.Value_Variant:
		alignLabelledFields(w.Variant.GetValue(), got.GetVariant().GetValue())
	case *apiv2.Value_Optional:
		alignLabelledFields(w.Optional.GetValue(), got.GetOptional().GetValue())
	case *apiv2.Value_List:
		gl := got.GetList().GetElements()
		for i := range min(len(w.List.GetElements()), len(gl)) {
			alignLabelledFields(w.List.GetElements()[i], gl[i])
		}
	case *apiv2.Value_TextMap:
		forEachTextMapPair(w.TextMap, got.GetTextMap(), alignLabelledFields)
	case *apiv2.Value_GenMap:
		ge := got.GetGenMap().GetEntries()
		for i := range min(len(w.GenMap.GetEntries()), len(ge)) {
			alignLabelledFields(w.GenMap.GetEntries()[i].GetKey(), ge[i].GetKey())
			alignLabelledFields(w.GenMap.GetEntries()[i].GetValue(), ge[i].GetValue())
		}
	}
}

// forEachTextMapPair pairs entries by key: the ledger returns a TextMap sorted by key, not in the args file's order.
func forEachTextMapPair(want, got *apiv2.TextMap, f func(want, got *apiv2.Value)) {
	byKey := make(map[string]*apiv2.Value, len(got.GetEntries()))
	for _, e := range got.GetEntries() {
		byKey[e.GetKey()] = e.GetValue()
	}
	for _, e := range want.GetEntries() {
		if g, ok := byKey[e.GetKey()]; ok {
			f(e.GetValue(), g)
		}
	}
}

func allLabelled(r *apiv2.Record) bool {
	for _, f := range r.GetFields() {
		if f.GetLabel() == "" {
			return false
		}
	}

	return len(r.GetFields()) > 0
}

func isNone(v *apiv2.Value) bool {
	o, ok := v.GetSum().(*apiv2.Value_Optional)

	return ok && o.Optional.GetValue() == nil
}

func (d DisclosedContract) Proto() (*apiv2.DisclosedContract, error) {
	if d.ContractID == "" {
		return nil, errors.New("contractId is empty")
	}
	blob, err := base64.StdEncoding.DecodeString(d.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("createdEventBlob of %s is not base64: %w", d.ContractID, err)
	}
	if len(blob) == 0 {
		return nil, fmt.Errorf("createdEventBlob of %s is empty", d.ContractID)
	}

	out := &apiv2.DisclosedContract{
		ContractId:       d.ContractID,
		CreatedEventBlob: blob,
		SynchronizerId:   d.SynchronizerID,
	}
	if d.TemplateID != "" {
		parts := strings.Split(d.TemplateID, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("templateId %q of %s is not <package>:<module>:<entity>", d.TemplateID, d.ContractID)
		}
		out.TemplateId = &apiv2.Identifier{PackageId: parts[0], ModuleName: parts[1], EntityName: parts[2]}
	}

	return out, nil
}
