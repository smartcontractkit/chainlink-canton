// Package cantonops contains helpers for looking up or creating the Daml
// contracts the CLI commands depend on (PerPartyRouter, CCIPSender,
// CCIPReceiver) and for extracting a CCIPMessageSent message id from a
// submitted transaction.
package cantonops

import (
	"context"
	"fmt"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
)

var (
	perPartyRouterTemplateID = &apiv2.Identifier{
		PackageId: "#ccip-runtime", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	}
	perPartyRouterFactoryTemplateID = &apiv2.Identifier{
		PackageId: "#ccip-runtime", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	}
	ccipSenderTemplateID = &apiv2.Identifier{
		PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender",
	}
	ccipReceiverTemplateID = &apiv2.Identifier{
		PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver",
	}
)

// propagationTimeout bounds how long we'll wait for a freshly created
// contract to show up in the active contract set.
const (
	propagationTimeout = time.Minute
	propagationPoll    = 2 * time.Second
)

// waitForActiveContract polls for an active contract matching templateID
// until one appears, ctx is cancelled, or propagationTimeout elapses.
func waitForActiveContract(
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

// GetOrCreateRouter returns the PerPartyRouter contract id for the
// participant's party, creating one via the factory if it doesn't exist.
func GetOrCreateRouter(ctx context.Context, participant canton.Participant, ccipEdsClient oapiCCIP.ClientWithResponsesInterface) (string, error) {
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, perPartyRouterTemplateID)
	if err != nil {
		return "", fmt.Errorf("list PerPartyRouter: %w", err)
	}
	if len(active) > 0 {
		return active[0].GetCreatedEvent().GetContractId(), nil
	}

	fmt.Printf("No active PerPartyRouter found for party %s, deploying one...\n", participant.PartyID)
	disclosures, err := eds.GetPerPartyRouterFactoryDisclosure(ctx, ccipEdsClient, participant.PartyID)
	if err != nil {
		return "", fmt.Errorf("get factory disclosure: %w", err)
	}

	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
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
			ActAs:              []string{participant.PartyID},
			DisclosedContracts: disclosures.DisclosedContracts,
		},
	})
	if err != nil {
		return "", fmt.Errorf("submit CreateRouter: %w", err)
	}
	fmt.Printf("Created PerPartyRouter in update: %s\n", resp.GetTransaction().GetUpdateId())

	contract, err := waitForActiveContract(ctx, participant, perPartyRouterTemplateID)
	if err != nil {
		return "", err
	}

	return contract.GetCreatedEvent().GetContractId(), nil
}

// GetOrCreateSender returns the CCIPSender contract id for the party,
// creating one if missing.
func GetOrCreateSender(ctx context.Context, participant canton.Participant) (string, error) {
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, ccipSenderTemplateID)
	if err != nil {
		return "", fmt.Errorf("list CCIPSender: %w", err)
	}
	if len(active) > 0 {
		return active[0].GetCreatedEvent().GetContractId(), nil
	}

	fmt.Printf("No active CCIPSender found for party %s, deploying one...\n", participant.PartyID)
	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: ccipSenderTemplateID,
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ccipsender"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: participant.PartyID}}},
					}},
				}},
			}},
			ActAs: []string{participant.PartyID},
		},
	})
	if err != nil {
		return "", fmt.Errorf("submit Create CCIPSender: %w", err)
	}
	fmt.Printf("Created CCIPSender in update: %s\n", resp.GetTransaction().GetUpdateId())

	contract, err := waitForActiveContract(ctx, participant, ccipSenderTemplateID)
	if err != nil {
		return "", err
	}

	return contract.GetCreatedEvent().GetContractId(), nil
}

func finalityConfigEqual(a, b core.FinalityConfig) bool {
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

func receiverFinalityLabel(cfg core.FinalityConfig) string {
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

func receiverInstanceID(cfg core.FinalityConfig) string {
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
	receiverFinality core.FinalityConfig,
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
) (string, error) {
	bindingCCVs := make([]chainlinkapi.RawInstanceAddress, len(requiredCCVs))
	for i, ccv := range requiredCCVs {
		bindingCCVs[i] = chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(ccv.String())}
	}

	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: ccipReceiverTemplateID,
					ContractId: receiverCid,
					Choice:     "UpdateRequiredCCVs",
					ChoiceArgument: ledger.MapToValue(receiver.UpdateRequiredCCVs{
						NewRequiredCCVs: bindingCCVs,
					}),
				}},
			}},
			ActAs: []string{participant.PartyID},
		},
	})
	if err != nil {
		return "", fmt.Errorf("submit UpdateRequiredCCVs: %w", err)
	}

	newCid, err := extractCreatedReceiverCID(resp.GetTransaction())
	if err != nil {
		return "", err
	}
	fmt.Printf("Updated CCIPReceiver required CCVs in update: %s\n", resp.GetTransaction().GetUpdateId())

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
	receiverFinality core.FinalityConfig,
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

func receiverFinalityField(cfg core.FinalityConfig) *apiv2.Value {
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
	receiverFinality core.FinalityConfig,
	requiredCCV contracts.RawInstanceAddress,
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
			"Updating CCIPReceiver %s required CCVs to %s...\n",
			receiverCid,
			requiredCCV,
		)

		return updateReceiverRequiredCCVs(ctx, participant, receiverCid, requiredCCVs)
	}

	fmt.Printf(
		"No CCIPReceiver with %s finality found for party %s, deploying one (CCV %s)...\n",
		receiverFinalityLabel(receiverFinality),
		participant.PartyID,
		requiredCCV,
	)
	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
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
			ActAs: []string{participant.PartyID},
		},
	})
	if err != nil {
		return "", fmt.Errorf("submit Create CCIPReceiver: %w", err)
	}
	fmt.Printf("Created CCIPReceiver in update: %s\n", resp.GetTransaction().GetUpdateId())

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
		msg, err := bindings.UnmarshalCreatedEvent[core.CCIPMessageSent](e.Created)
		if err != nil {
			return "", fmt.Errorf("unmarshal CCIPMessageSent: %w", err)
		}

		return string(msg.Event.MessageId), nil
	}

	return "", fmt.Errorf("CCIPMessageSent event not found in transaction events")
}
