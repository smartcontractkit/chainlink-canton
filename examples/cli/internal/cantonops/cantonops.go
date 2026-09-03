// Package cantonops contains helpers for looking up or creating the Daml
// contracts the CLI commands depend on (PerPartyRouter, CCIPSender,
// CCIPReceiver) and for extracting a CCIPMessageSent message id from a
// submitted transaction.
package cantonops

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/events"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/ledger/usbwallet"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
)

var (
	perPartyRouterTemplateID        = contracts.IdentifierFromBinding(ccipruntime.PerPartyRouter{})
	perPartyRouterFactoryTemplateID = contracts.IdentifierFromBinding(ccipruntime.PerPartyRouterFactory{})
	ccipSenderTemplateID            = contracts.IdentifierFromBinding(sender.CCIPSender{})
	ccipReceiverTemplateID          = contracts.IdentifierFromBinding(receiver.CCIPReceiver{})
)

// propagationTimeout bounds how long we'll wait for a freshly created
// contract to show up in the active contract set.
const (
	propagationTimeout = time.Minute
	propagationPoll    = 2 * time.Second
)

// WaitForActiveContract polls for an active contract matching templateID
// until one appears, ctx is cancelled, or propagationTimeout elapses.
func WaitForActiveContract(
	ctx context.Context,
	participant canton.Participant,
	templateID *apiv2.Identifier,
) (*apiv2.ActiveContract, error) {
	deadline := time.Now().Add(propagationTimeout)
	for {
		active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, templateID)
		if err != nil {
			return nil, fmt.Errorf("list %s: %w", templateID.GetEntityName(), err)
		}
		if len(active) > 0 {
			return active[0], nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out after %s waiting for %s to become active", propagationTimeout, templateID.GetEntityName())
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(propagationPoll):
		}
	}
}

func ParseDerivationPath(pathOrIndex string) (accounts.DerivationPath, error) {
	index, err := strconv.ParseUint(pathOrIndex, 10, 32)
	if err == nil {
		return accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 6767, 0x80000000 + 0, 0x80000000 + 0, 0x80000000 + uint32(index)}, nil
	}

	return accounts.ParseDerivationPath(pathOrIndex)
}

func CantonSubmit(
	ctx context.Context,
	participant canton.Participant,
	ledgerFlag string,
	commands []*apiv2.Command,
	disclosedContracts []*apiv2.DisclosedContract,
) (*apiv2.Transaction, error) {
	// Parse ledger flag to determine submission path
	if ledgerFlag == "" {
		// Direct submission using local party
		resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId:          uuid.NewString(),
				Commands:           commands,
				ActAs:              []string{participant.PartyID},
				DisclosedContracts: disclosedContracts,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("submitAndWaitForTransaction: %w", err)
		}

		return resp.GetTransaction(), nil
	}
	// Interactive Ledger signing flow

	// Parse derivation path
	derivationPath, err := ParseDerivationPath(ledgerFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid ledger derivation path: %w", err)
	}

	// Connect Ledger
	fmt.Println("Looking for connected Ledger devices...")
	ledgerHub, err := usbwallet.NewCantonLedgerHub()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ledger hub: %w", err)
	}
	wallets := ledgerHub.Wallets()
	if len(wallets) == 0 {
		return nil, fmt.Errorf("no Ledger wallets found. Please connect a Ledger device and try again")
	}
	fmt.Printf("Found %d connected Ledger wallets\n", len(wallets))
	wallet := wallets[0]

	err = wallet.Open("")
	if err != nil {
		return nil, fmt.Errorf("failed to open Ledger wallet: %w", err)
	}
	defer wallet.Close()

	fmt.Printf("Getting public key for derivation path: %v\n", derivationPath.String())

	// Get the public key at the derivation path and compare with expected key
	pubkey, err := wallet.GetPublicKey(derivationPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from Ledger: %w", err)
	}
	fingerprint, err := usbwallet.Fingerprint(pubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to get fingerprint from Ledger: %w", err)
	}
	fmt.Printf("Using key with fingerprint %s to sign transaction\n", fingerprint)
	splitParty := strings.Split(participant.PartyID, "::")
	if len(splitParty) != 2 {
		return nil, fmt.Errorf("wrong partyId format, expected hint::fingerprint, got: %s", participant.PartyID)
	}
	if splitParty[1] != fingerprint {
		return nil, fmt.Errorf("fingerprint mismatch: Ledger device has %s, but participant partyId has %s", fingerprint, splitParty[1])
	}

	// Prepare transaction
	preparedSubmission, err := participant.LedgerServices.InteractiveSubmission.PrepareSubmission(ctx, &interactive.PrepareSubmissionRequest{
		CommandId:          uuid.NewString(),
		Commands:           commands,
		ActAs:              []string{participant.PartyID},
		ReadAs:             []string{participant.PartyID},
		DisclosedContracts: disclosedContracts,
	})
	if err != nil {
		return nil, fmt.Errorf("prepareSubmission: %w", err)
	}

	// Sign transaction
	fmt.Printf("📜 Signing transaction, tx hash=%v\n", strings.ToUpper(hex.EncodeToString(preparedSubmission.GetPreparedTransactionHash())))
	signature, err := signPreparedTransaction(wallet, derivationPath, preparedSubmission)
	if err != nil {
		return nil, fmt.Errorf("failed to sign prepared transaction: %w", err)
	}
	fmt.Println("✅ Transaction signed.")

	// Submit transaction
	resp, err := participant.LedgerServices.InteractiveSubmission.ExecuteSubmissionAndWaitForTransaction(ctx, &interactive.ExecuteSubmissionAndWaitForTransactionRequest{
		SubmissionId:        uuid.NewString(),
		PreparedTransaction: preparedSubmission.GetPreparedTransaction(),
		PartySignatures: &interactive.PartySignatures{
			Signatures: []*interactive.SinglePartySignatures{
				{
					Party: participant.PartyID,
					Signatures: []*apiv2.Signature{{
						Format:               apiv2.SignatureFormat_SIGNATURE_FORMAT_CONCAT,
						Signature:            signature,
						SignedBy:             fingerprint,
						SigningAlgorithmSpec: apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_ED25519,
					}},
				},
			},
		},
		HashingSchemeVersion: preparedSubmission.GetHashingSchemeVersion(),
	})
	if err != nil {
		return nil, fmt.Errorf("executeSubmissionAndWaitForTransaction: %w", err)
	}

	return resp.GetTransaction(), nil
}

// signPreparedTransaction streams the prepared transaction to the device so it recomputes the
// transaction hash itself.
//
// The Canton app can only do that for transactions that fit its on-device parser. Large or deeply
// nested transactions, such as a CCIP send, overrun the tiny dynamic memory pool or the fixed
// component buffer of the app and are rejected with SW_TX_PARSING_FAIL / SW_WRONG_TX_LENGTH. That
// is not the graceful "template unknown, blind sign it" path of the app, it aborts the signing
// session outright, so we retry by signing the transaction hash reported by the participant
// instead.
func signPreparedTransaction(
	wallet usbwallet.CantonWallet,
	derivationPath accounts.DerivationPath,
	preparedSubmission *interactive.PrepareSubmissionResponse,
) ([]byte, error) {
	signature, err := wallet.SignPreparedTransaction(derivationPath, preparedSubmission.GetPreparedTransaction())
	if err == nil {
		return signature, nil
	}
	if !errors.Is(err, usbwallet.ErrTransactionParsingFailed) && !errors.Is(err, usbwallet.ErrTransactionTooLong) {
		return nil, err
	}

	fmt.Printf("⚠️ The Ledger device could not process this transaction (%v).\n", err)
	fmt.Println("⚠️ Falling back to blind signing: the device signs the transaction hash reported by")
	fmt.Println("⚠️ the participant instead of recomputing it, so it cannot verify what it signs.")
	fmt.Println("⚠️ Blind signing has to be enabled in the Canton app settings for this to work.")

	signature, err = wallet.SignHash(derivationPath, preparedSubmission.GetPreparedTransactionHash())
	if err != nil {
		if errors.Is(err, usbwallet.ErrBlindSigningDisabled) {
			return nil, fmt.Errorf("blind signing is disabled in the Canton app settings: %w", err)
		}

		return nil, fmt.Errorf("blind signing the transaction hash: %w", err)
	}

	return signature, nil
}

// GetOrCreateRouter returns the PerPartyRouter contract id for the
// participant's party, creating one via the factory if it doesn't exist.
func GetOrCreateRouter(ctx context.Context, participant canton.Participant, ccipEdsClient oapiCCIP.ClientWithResponsesInterface, useLedger string) (string, error) {
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, perPartyRouterTemplateID)
	if err != nil {
		return "", fmt.Errorf("list PerPartyRouter: %w", err)
	}
	if len(active) > 0 {
		return active[0].GetCreatedEvent().GetContractId(), nil
	}

	fmt.Printf("⚠️ No active PerPartyRouter found for party %s, deploying one...\n", participant.PartyID)
	disclosures, err := eds.GetPerPartyRouterFactoryDisclosure(ctx, ccipEdsClient, participant.PartyID)
	if err != nil {
		return "", fmt.Errorf("get factory disclosure: %w", err)
	}

	tx, err := CantonSubmit(
		ctx,
		participant,
		useLedger,
		[]*apiv2.Command{{
			Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
				TemplateId: perPartyRouterFactoryTemplateID,
				ContractId: disclosures.ContractId,
				Choice:     "CreateRouter",
				ChoiceArgument: ledger.MapToValue(ccipruntime.CreateRouter{
					PartyOwner: types.PARTY(participant.PartyID),
					InstanceId: types.TEXT(fmt.Sprintf("router-%s", participant.PartyID)),
				}),
			}},
		}},
		disclosures.DisclosedContracts,
	)
	if err != nil {
		return "", fmt.Errorf("submit CreateRouter: %w", err)
	}
	fmt.Printf("Created PerPartyRouter in update: %s\n", tx.GetUpdateId())

	contract, err := WaitForActiveContract(ctx, participant, perPartyRouterTemplateID)
	if err != nil {
		return "", err
	}

	return contract.GetCreatedEvent().GetContractId(), nil
}

// GetOrCreateSender returns the CCIPSender contract id for the party,
// creating one if missing.
func GetOrCreateSender(ctx context.Context, participant canton.Participant, useLedger string) (string, error) {
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, ccipSenderTemplateID)
	if err != nil {
		return "", fmt.Errorf("list CCIPSender: %w", err)
	}
	if len(active) > 0 {
		return active[0].GetCreatedEvent().GetContractId(), nil
	}

	fmt.Printf("⚠️ No active CCIPSender found for party %s, deploying one...\n", participant.PartyID)
	tx, err := CantonSubmit(
		ctx,
		participant,
		useLedger,
		[]*apiv2.Command{{
			Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
				TemplateId: ccipSenderTemplateID,
				CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ccipsender"}}},
					{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: participant.PartyID}}},
				}},
			}},
		}},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("submit Create CCIPSender: %w", err)
	}
	fmt.Printf("Created CCIPSender in update: %s\n", tx.GetUpdateId())

	contract, err := WaitForActiveContract(ctx, participant, ccipSenderTemplateID)
	if err != nil {
		return "", err
	}

	return contract.GetCreatedEvent().GetContractId(), nil
}

func finalityConfigEqual(a, b ccipcodec.FinalityConfig) bool {
	switch a.GetVariantTag() {
	case "WaitForFinality":
		return b.GetVariantTag() == "WaitForFinality"
	case "WaitForSafe":
		return b.GetVariantTag() == "WaitForSafe"
	case "BlockDepth":
		aDepth, aOk := a.GetVariantValue().(*types.INT64)
		bDepth, bOk := b.GetVariantValue().(*types.INT64)
		if !aOk || !bOk || aDepth == nil || bDepth == nil {
			return false
		}

		return *aDepth == *bDepth
	default:
		return false
	}
}

func receiverFinalityLabel(cfg ccipcodec.FinalityConfig) string {
	switch cfg.GetVariantTag() {
	case "BlockDepth":
		depth, ok := cfg.GetVariantValue().(*types.INT64)
		if !ok || depth == nil {
			return "BlockDepth"
		}

		return fmt.Sprintf("BlockDepth(%d)", int64(*depth))
	default:
		return cfg.GetVariantTag()
	}
}

func receiverInstanceID(cfg ccipcodec.FinalityConfig) string {
	switch cfg.GetVariantTag() {
	case "WaitForFinality":
		return "ccipreceiver-WaitForFinality"
	case "WaitForSafe":
		return "ccipreceiver-WaitForSafe"
	case "BlockDepth":
		depth, ok := cfg.GetVariantValue().(*types.INT64)
		if !ok || depth == nil {
			return "ccipreceiver-BlockDepth"
		}

		return fmt.Sprintf("ccipreceiver-BlockDepth-%d", int64(*depth))
	default:
		return fmt.Sprintf("ccipreceiver-%s", cfg.GetVariantTag())
	}
}

func rawInstanceAddressValue(addr contracts.RawInstanceAddress) *apiv2.Value {
	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "unpack", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: addr.String()}}},
	}}}}
}

func requiredCCVsListValue(requiredCCVs []contracts.RawInstanceAddress) *apiv2.Value {
	elements := make([]*apiv2.Value, len(requiredCCVs))
	for i, ccv := range requiredCCVs {
		elements[i] = rawInstanceAddressValue(ccv)
	}

	return &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: elements}}}
}

func receiverRequiredCCVConfigured(recv *receiver.CCIPReceiver, requiredCCV contracts.RawInstanceAddress) bool {
	for _, ccv := range recv.RequiredCCVs {
		if string(ccv.Unpack) == requiredCCV.String() {
			return true
		}
	}

	return false
}

func findActiveReceiverByFinality(
	ctx context.Context,
	participant canton.Participant,
	receiverFinality ccipcodec.FinalityConfig,
) (string, *receiver.CCIPReceiver, error) {
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, ccipReceiverTemplateID)
	if err != nil {
		return "", nil, fmt.Errorf("list CCIPReceiver: %w", err)
	}

	for _, ac := range active {
		recv, err := bindings.UnmarshalCreatedEvent[receiver.CCIPReceiver](ac.GetCreatedEvent())
		if err != nil {
			return "", nil, fmt.Errorf("unmarshal CCIPReceiver: %w", err)
		}
		if string(recv.Owner) != participant.PartyID {
			continue
		}
		if finalityConfigEqual(recv.ReceiverFinalityConfig, receiverFinality) {
			return ac.GetCreatedEvent().GetContractId(), recv, nil
		}
	}

	return "", nil, nil
}

func updateReceiverRequiredCCVs(
	ctx context.Context,
	participant canton.Participant,
	receiverCid string,
	requiredCCVs []contracts.RawInstanceAddress,
	useLedger string,
) (string, error) {
	bindingCCVs := make([]chainlinkapi.RawInstanceAddress, len(requiredCCVs))
	for i, ccv := range requiredCCVs {
		bindingCCVs[i] = chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(ccv.String())}
	}

	tx, err := CantonSubmit(
		ctx,
		participant,
		useLedger,
		[]*apiv2.Command{{
			Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
				TemplateId: ccipReceiverTemplateID,
				ContractId: receiverCid,
				Choice:     "UpdateRequiredCCVs",
				ChoiceArgument: ledger.MapToValue(receiver.UpdateRequiredCCVs{
					NewRequiredCCVs: bindingCCVs,
				}),
			}},
		}},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("submit UpdateRequiredCCVs: %w", err)
	}

	newCid, err := extractCreatedReceiverCID(tx)
	if err != nil {
		return "", err
	}
	fmt.Printf("Updated CCIPReceiver required CCVs in update: %s\n", tx.GetUpdateId())

	return newCid, nil
}

func extractCreatedReceiverCID(tx *apiv2.Transaction) (string, error) {
	for _, event := range tx.GetEvents() {
		e, ok := event.GetEvent().(*apiv2.Event_Created)
		if !ok {
			continue
		}
		if e.Created.GetTemplateId().GetEntityName() != "CCIPReceiver" {
			continue
		}

		return e.Created.GetContractId(), nil
	}

	return "", fmt.Errorf("CCIPReceiver created event not found in transaction events")
}

func waitForReceiverWithFinality(
	ctx context.Context,
	participant canton.Participant,
	receiverFinality ccipcodec.FinalityConfig,
) (string, error) {
	deadline := time.Now().Add(propagationTimeout)
	for {
		receiverCid, _, err := findActiveReceiverByFinality(ctx, participant, receiverFinality)
		if err != nil {
			return "", err
		}
		if receiverCid != "" {
			return receiverCid, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf(
				"timed out after %s waiting for CCIPReceiver with %s finality",
				propagationTimeout,
				receiverFinalityLabel(receiverFinality),
			)
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(propagationPoll):
		}
	}
}

func receiverFinalityField(cfg ccipcodec.FinalityConfig) *apiv2.Value {
	switch cfg.GetVariantTag() {
	case "WaitForFinality":
		return &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
			Constructor: "WaitForFinality",
			Value:       &apiv2.Value{Sum: &apiv2.Value_Unit{}},
		}}}
	case "WaitForSafe":
		return &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
			Constructor: "WaitForSafe",
			Value:       &apiv2.Value{Sum: &apiv2.Value_Unit{}},
		}}}
	case "BlockDepth":
		depth, ok := cfg.GetVariantValue().(*types.INT64)
		if !ok || depth == nil {
			panic("invalid BlockDepth finality config")
		}

		return &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
			Constructor: "BlockDepth",
			Value:       &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(*depth)}},
		}}}
	default:
		panic(fmt.Sprintf("unsupported receiver finality config %q", cfg.GetVariantTag()))
	}
}

// GetOrCreateReceiver returns the CCIPReceiver contract id for the party with
// the requested finality config and attesting CCV, creating or updating one as
// needed. Multiple receivers per party are supported (e.g. full finality and FTF).
func GetOrCreateReceiver(
	ctx context.Context,
	participant canton.Participant,
	receiverFinality ccipcodec.FinalityConfig,
	requiredCCV contracts.RawInstanceAddress,
	useLedger string,
) (string, error) {
	requiredCCVs := []contracts.RawInstanceAddress{requiredCCV}

	receiverCid, recv, err := findActiveReceiverByFinality(ctx, participant, receiverFinality)
	if err != nil {
		return "", err
	}
	if receiverCid != "" {
		if receiverRequiredCCVConfigured(recv, requiredCCV) {
			fmt.Printf(
				"Using CCIPReceiver %s (%s finality, CCV %s)\n",
				receiverCid,
				receiverFinalityLabel(receiverFinality),
				requiredCCV,
			)

			return receiverCid, nil
		}

		fmt.Printf(
			"⚠️ Updating CCIPReceiver %s required CCVs to %s...\n",
			receiverCid,
			requiredCCV,
		)

		return updateReceiverRequiredCCVs(ctx, participant, receiverCid, requiredCCVs, useLedger)
	}

	fmt.Printf(
		"⚠️ No CCIPReceiver with %s finality found for party %s, deploying one (CCV %s)...\n",
		receiverFinalityLabel(receiverFinality),
		participant.PartyID,
		requiredCCV,
	)
	tx, err := CantonSubmit(
		ctx,
		participant,
		useLedger,
		[]*apiv2.Command{{
			Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
				TemplateId: ccipReceiverTemplateID,
				CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: receiverInstanceID(receiverFinality)}}},
					{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: participant.PartyID}}},
					{Label: "receiverFinalityConfig", Value: receiverFinalityField(receiverFinality)},
					{Label: "requiredCCVs", Value: requiredCCVsListValue(requiredCCVs)},
					{Label: "optionalCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					{Label: "optionalThreshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
				}},
			}},
		}},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("submit Create CCIPReceiver: %w", err)
	}
	fmt.Printf("Created CCIPReceiver in update: %s\n", tx.GetUpdateId())

	return waitForReceiverWithFinality(ctx, participant, receiverFinality)
}

// GetMessageIdFromTransaction scans the transaction events for a created
// CCIPMessageSent contract and returns its message id, or an error.
func GetMessageIdFromTransaction(tx *apiv2.Transaction) (string, error) {
	for _, event := range tx.GetEvents() {
		e, ok := event.GetEvent().(*apiv2.Event_Created)
		if !ok {
			continue
		}
		if e.Created.GetTemplateId().GetEntityName() != "CCIPMessageSent" {
			continue
		}
		msg, err := bindings.UnmarshalCreatedEvent[events.CCIPMessageSent](e.Created)
		if err != nil {
			return "", fmt.Errorf("unmarshal CCIPMessageSent: %w", err)
		}

		return string(msg.Event.MessageId), nil
	}

	return "", fmt.Errorf("CCIPMessageSent event not found in transaction events")
}
