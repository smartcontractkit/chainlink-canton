package devenv

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	ccipreceiver "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	ccipsender "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func instanceIDFromEnv(key, defaultID string) contracts.InstanceID {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return contracts.InstanceID(v)
	}

	return contracts.InstanceID(defaultID)
}

// DeployPerPartyRouter uses the PerPartyRouterFactory to create a PerPartyRouter instance for the given party.
// It returns the partyOwner-based instance address used in CCIP protocol fields, reusing an existing router on ledger when present.
func (c *Chain) DeployPerPartyRouter(ctx context.Context, participant canton.Participant, partyId string) (routerAddress contracts.InstanceAddress, err error) {
	var unset contracts.InstanceAddress
	routerInstanceID := instanceIDFromEnv("CANTON_ROUTER_INSTANCE_ID", "test-router")
	if c.routerAddress != unset {
		if found, _, ok, findErr := c.findPerPartyRouterByParty(ctx, participant, partyId, routerInstanceID); findErr != nil {
			return contracts.InstanceAddress{}, fmt.Errorf("verify cached per-party router: %w", findErr)
		} else if ok && found == c.routerAddress {
			return c.routerAddress, nil
		}
		c.routerAddress = unset
	}

	if found, _, ok, findErr := c.findPerPartyRouterByParty(ctx, participant, partyId, routerInstanceID); findErr != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("find existing per-party router: %w", findErr)
	} else if ok {
		c.routerAddress = found
		return found, nil
	}

	perPartyRouterFactoryDisclosure, err := c.GetPerPartyRouterFactoryDisclosure(ctx, partyId)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to get canton per party router factory disclosure: %w", err)
	}
	c.logger.Debug().Str("ContractId", perPartyRouterFactoryDisclosure.ContractId).Msg("Resolved per-party router factory address")

	_, err = operations.ExecuteOperation(
		c.e.OperationsBundle,
		per_party_router_factory.CreateRouter,
		c.chain,
		contract.ChoiceInput[ccipruntime.CreateRouter]{
			InstanceAddress:    perPartyRouterFactoryDisclosure.Address.InstanceAddress(),
			ContractID:         perPartyRouterFactoryDisclosure.ContractId,
			ParticipantIndex:   c.clientParticipantIndex(),
			DisclosedContracts: contract.DisclosedContractsFromProto(perPartyRouterFactoryDisclosure.DisclosedContracts),
			Args: ccipruntime.CreateRouter{
				PartyOwner: types.PARTY(partyId),
				InstanceId: types.TEXT(routerInstanceID.String()),
			},
		},
		operations.WithForceExecute[contract.ChoiceInput[ccipruntime.CreateRouter], canton.Chain](),
	)
	if err != nil {
		if found, _, ok, findErr := c.findPerPartyRouterByParty(ctx, participant, partyId, routerInstanceID); findErr != nil {
			return contracts.InstanceAddress{}, fmt.Errorf("create per-party router: %w (find after failure: %w)", err, findErr)
		} else if ok {
			c.routerAddress = found
			return found, nil
		}

		return contracts.InstanceAddress{}, fmt.Errorf("create per-party router: %w", err)
	}

	found, _, ok, findErr := c.findPerPartyRouterByParty(ctx, participant, partyId, routerInstanceID)
	if findErr != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("find per-party router after create: %w", findErr)
	}
	if !ok {
		return contracts.InstanceAddress{}, fmt.Errorf("per-party router not found after create for party %s", partyId)
	}

	c.routerAddress = found

	return found, nil
}

func (c *Chain) findPerPartyRouterByParty(
	ctx context.Context,
	participant canton.Participant,
	partyId string,
	preferredInstanceID contracts.InstanceID,
) (contracts.InstanceAddress, string, bool, error) {
	templateID := contracts.TemplateIDFromBinding(ccipruntime.PerPartyRouter{}).ToLedgerIdentifier()
	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, templateID)
	if err != nil {
		return contracts.InstanceAddress{}, "", false, err
	}

	var fallback contracts.InstanceAddress
	var fallbackCid string
	var hasFallback bool
	for _, ac := range activeContracts {
		created := ac.GetCreatedEvent()
		if created == nil {
			continue
		}
		instanceId, partyOwner, ok := perPartyRouterFieldsFromCreated(created)
		if !ok {
			c.logger.Debug().Msg("Skipping PerPartyRouter active contract missing instanceId or partyOwner")
			continue
		}
		if partyOwner != partyId {
			continue
		}
		addr := contracts.InstanceID(instanceId).RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()
		cid := created.GetContractId()
		if contracts.InstanceID(instanceId) == preferredInstanceID {
			return addr, cid, true, nil
		}
		if !hasFallback {
			fallback = addr
			fallbackCid = cid
			hasFallback = true
		}
	}
	if hasFallback {
		return fallback, fallbackCid, true, nil
	}

	return contracts.InstanceAddress{}, "", false, nil
}

func perPartyRouterFieldsFromCreated(created *apiv2.CreatedEvent) (instanceId, partyOwner string, ok bool) {
	for _, field := range created.GetCreateArguments().GetFields() {
		switch field.GetLabel() {
		case "instanceId":
			instanceId = field.GetValue().GetText()
		case "partyOwner":
			partyOwner = field.GetValue().GetParty()
		}
	}

	return instanceId, partyOwner, instanceId != "" && partyOwner != ""
}

// findPerPartyRouterCidByParty resolves the ledger contract ID for a PerPartyRouter owned by partyId.
// PerPartyRouter is signed by ccipOwner, so instance-address lookup must use partyOwner (see OnRamp.daml).
func (c *Chain) findPerPartyRouterCidByParty(ctx context.Context, participant canton.Participant, partyId string) (string, error) {
	routerInstanceID := instanceIDFromEnv("CANTON_ROUTER_INSTANCE_ID", "test-router")
	_, cid, ok, err := c.findPerPartyRouterByParty(ctx, participant, partyId, routerInstanceID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("no active PerPartyRouter found for party %s", partyId)
	}

	return cid, nil
}

// DeployCCIPSender returns a sender-owned CCIPSender instance address, deploying only when missing.
func (c *Chain) DeployCCIPSender(ctx context.Context, participant canton.Participant, partyId string) (contracts.InstanceAddress, error) {
	var unset contracts.InstanceAddress
	if c.senderAddress != unset {
		return c.senderAddress, nil
	}

	instanceID := instanceIDFromEnv("CANTON_SENDER_INSTANCE_ID", "e2e-ccipsender")
	senderAddress := instanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()

	if _, err := contract.FindActiveContractIDByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		ccipsender.CCIPSender{}.GetTemplateID(),
		senderAddress,
	); err == nil {
		c.senderAddress = senderAddress
		return senderAddress, nil
	}

	_, err := operations.ExecuteOperation(c.e.OperationsBundle, sender.Deploy, c.chain, contract.DeployInput[ccipsender.CCIPSender]{
		Qualifier:        nil,
		ParticipantIndex: c.clientParticipantIndex(),
		Template: ccipsender.CCIPSender{
			InstanceId: types.TEXT(instanceID),
			Owner:      types.PARTY(partyId),
		},
		OwnerParty: types.PARTY(partyId),
	})
	if err != nil {
		if _, findErr := contract.FindActiveContractIDByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			contract.LedgerQueryParties(participant),
			ccipsender.CCIPSender{}.GetTemplateID(),
			senderAddress,
		); findErr == nil {
			c.senderAddress = senderAddress
			return senderAddress, nil
		}

		return contracts.InstanceAddress{}, fmt.Errorf("failed to deploy ccip sender contract: %w", err)
	}

	c.senderAddress = senderAddress

	return senderAddress, nil
}

func (c *Chain) DeployCCIPReceiver(ctx context.Context, participant canton.Participant, partyId string, receiverFinality int64) (contracts.InstanceAddress, error) {
	var unset contracts.InstanceAddress
	if c.receiverAddress != unset {
		return c.receiverAddress, nil
	}

	instanceID := instanceIDFromEnv("CANTON_RECEIVER_INSTANCE_ID", "e2e-receiver")
	receiverAddress := instanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()

	if _, err := contract.FindActiveContractIDByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		ccipreceiver.CCIPReceiver{}.GetTemplateID(),
		receiverAddress,
	); err == nil {
		c.receiverAddress = receiverAddress
		return receiverAddress, nil
	}

	finalityConfig, err := encodeReceiverFinalityConfig(receiverFinality)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to encode receiver finality config: %w", err)
	}

	_, err = operations.ExecuteOperation(c.e.OperationsBundle, receiver.Deploy, c.chain, contract.DeployInput[ccipreceiver.CCIPReceiver]{
		Qualifier:        nil,
		ParticipantIndex: c.clientParticipantIndex(),
		Template: ccipreceiver.CCIPReceiver{
			InstanceId:             types.TEXT(instanceID),
			Owner:                  types.PARTY(partyId),
			RequiredCCVs:           nil,
			OptionalCCVs:           nil,
			OptionalThreshold:      0,
			ReceiverFinalityConfig: finalityConfig,
		},
		OwnerParty: types.PARTY(partyId),
	})
	if err != nil {
		if _, findErr := contract.FindActiveContractIDByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			contract.LedgerQueryParties(participant),
			ccipreceiver.CCIPReceiver{}.GetTemplateID(),
			receiverAddress,
		); findErr == nil {
			c.receiverAddress = receiverAddress
			return receiverAddress, nil
		}

		return contracts.InstanceAddress{}, fmt.Errorf("failed to deploy receiver contract: %w", err)
	}

	c.receiverAddress = receiverAddress

	return receiverAddress, nil
}

// ManuallyExecuteMessage implements cciptestinterfaces.CCIP17.
func (c *Chain) ManuallyExecuteMessage(ctx context.Context, message protocol.Message, gasLimit uint64, verifiers []protocol.UnknownAddress, verifierResults [][]byte) (cciptestinterfaces.ExecutionStateChangedEvent, error) {
	participant, clientIdx, err := c.ClientParticipant()
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("no canton participants configured: %w", err)
	}

	// Ensure that the message receiver is the party we're executing with
	executingParty := participant.PartyID
	if contracts.HashedPartyFromString(executingParty) != contracts.BytesToHashedParty(message.Receiver.Bytes()) {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("message receiver %s does not match executing party %s (%s)", hex.EncodeToString(message.Receiver), contracts.HashedPartyFromString(executingParty).String(), executingParty)
	}

	routerAddress := c.routerAddress
	if routerAddress == (contracts.InstanceAddress{}) {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf(
			"per-party router not deployed; call SetupReceive or SetupSend on the client participant before executing messages",
		)
	}

	// Deploy CCIPReceiver contract
	receiverAddress, err := c.DeployCCIPReceiver(ctx, participant, executingParty, int64(message.Finality))
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to deploy CCIPReceiver contract: %w", err)
	}
	c.logger.Debug().Str("ReceiverAddress", receiverAddress.String()).Msg("Deployed CCIPReceiver")

	encodedMessage, err := message.Encode()
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to encode message: %w", err)
	}
	ccvs := make([]contracts.InstanceAddress, len(verifiers))
	for i, verifier := range verifiers {
		ccvs[i] = contracts.HexToInstanceAddress(verifier.String())
	}
	encodedMessageHex := hex.EncodeToString(encodedMessage)

	routerCid, err := c.findPerPartyRouterCidByParty(ctx, participant, executingParty)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get router contract ID: %w", err)
	}
	c.logger.Debug().Str("InstanceAddress", routerAddress.String()).Str("ContractId", routerCid).Msg("Resolved PerPartyRouter contract")

	// Collect disclosures
	// CCIP
	ccipExecuteDisclosure, err := c.GetCCIPExecuteDisclosure(ctx, encodedMessageHex)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get CCIP execute disclosure: %w", err)
	}
	executeArgs := ccipreceiver.Execute{
		Context:        ccipExecuteDisclosure.ChoiceContext,
		RouterCid:      types.CONTRACT_ID(routerCid),
		EncodedMessage: types.TEXT(encodedMessageHex),
		TokenTransfer:  nil,
		CcvInputs:      make([]ccipreceiver.CCVInput, len(verifiers)),
	}
	disclosedContracts := ccipExecuteDisclosure.DisclosedContracts

	// CCVs
	if len(verifierResults) != len(verifiers) {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("verifierResults length %d does not match verifiers length %d", len(verifierResults), len(verifiers))
	}
	for i, vr := range verifierResults {
		verifier := ccvs[i]
		ccvExecuteDisclosure, err := c.GetCCVExecuteDisclosure(ctx, encodedMessageHex, verifier)
		if err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get CCV execute disclosure for verifier %s: %w", verifier.String(), err)
		}

		executeArgs.CcvInputs[i] = ccipreceiver.CCVInput{
			CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
			VerifierResults: types.TEXT(hex.EncodeToString(vr)),
			CcvExtraContext: ccvExecuteDisclosure.ChoiceContext,
		}
		disclosedContracts = append(disclosedContracts, ccvExecuteDisclosure.DisclosedContracts...)
	}

	// Token Pool
	if message.TokenTransfer != nil {
		hashedInstrumentId := contracts.BytesToEncodedInstrumentID(message.TokenTransfer.DestTokenAddress)
		tokenPoolAddress, err := c.GetTokenPoolForToken(ctx, hashedInstrumentId)
		if err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get token pool for token %s: %w", hashedInstrumentId.String(), err)
		}

		tokenPoolDisclosure, err := c.GetTokenPoolExecuteDisclosure(ctx, encodedMessageHex, tokenPoolAddress.InstanceAddress())
		if err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get token pool execute disclosure: %w", err)
		}

		executeArgs.TokenTransfer = &ccipreceiver.TokenTransferInput{
			TokenPoolCid:       types.CONTRACT_ID(tokenPoolDisclosure.ContractId),
			TokenReceiverParty: types.PARTY(executingParty),
			PoolExtraContext:   tokenPoolDisclosure.ChoiceContext,
		}
		disclosedContracts = append(disclosedContracts, tokenPoolDisclosure.DisclosedContracts...)
	}

	// Execute message
	c.logger.Debug().
		Str("EncodedMessage", hex.EncodeToString(encodedMessage)).
		Str("VerifierResults", hex.EncodeToString(verifierResults[0])).
		Str("Receiver", hex.EncodeToString(message.Receiver)).
		Msg("Executing message...")

	executeReport, err := operations.ExecuteOperation(c.e.OperationsBundle, receiver.Execute, c.chain, contract.ChoiceInput[ccipreceiver.Execute]{
		InstanceAddress:    receiverAddress,
		ParticipantIndex:   clientIdx,
		Args:               executeArgs,
		DisclosedContracts: contract.DisclosedContractsFromProto(disclosedContracts),
	})
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to execute message: %w", err)
	}

	// TODO: refactor to use our standard unmarshaler
	c.logger.Info().Str("UpdateID", executeReport.Output.ExecInfo.UpdateID).Msg("Message executed")
	update, err := participant.LedgerServices.Update.GetUpdateById(ctx, &apiv2.GetUpdateByIdRequest{
		UpdateId: executeReport.Output.ExecInfo.UpdateID,
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						participant.PartyID: {
							Cumulative: []*apiv2.CumulativeFilter{
								{
									IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
										WildcardFilter: &apiv2.WildcardFilter{
											IncludeCreatedEventBlob: false,
										},
									},
								},
							},
						},
					},
					Verbose: true,
				},
			},
		},
	})
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get update %q: %w", executeReport.Output.ExecInfo.UpdateID, err)
	}
	c.logger.Debug().Str("UpdateID", update.GetTransaction().GetUpdateId()).Msg("Executed message")
	if message.TokenTransfer != nil {
		pendingTransferInstructionCID := ""
		for _, event := range update.GetTransaction().GetEvents() {
			if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
				templateName := e.Created.GetTemplateId().GetEntityName()
				if strings.Contains(templateName, "TransferInstruction") {
					pendingTransferInstructionCID = e.Created.GetContractId()
				}
			}
		}
		if pendingTransferInstructionCID != "" {
			_, _, transferClient, err := testhelpers.NewValidatorAPIClients(participant)
			if err != nil {
				return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("create transfer instruction client: %w", err)
			}
			if err := testhelpers.AcceptPendingTransferInstruction(ctx, participant, transferClient, executingParty, pendingTransferInstructionCID); err != nil {
				return cciptestinterfaces.ExecutionStateChangedEvent{}, err
			}
		}
	}

	// Get ExecutionStateChangedEvent from events
	expectedTemplateID := core.ExecutionStateChanged{}.GetTemplateID()
	for _, event := range update.GetTransaction().GetEvents() {
		//nolint:nestif // need to check if all of these are nil
		if createdEvent := event.GetCreated(); createdEvent != nil {
			if templateId := createdEvent.GetTemplateId(); templateId != nil {
				gotTemplateId := fmt.Sprintf("#%s:%s:%s", createdEvent.GetPackageName(), templateId.GetModuleName(), templateId.GetEntityName())
				if gotTemplateId == expectedTemplateID {
					// Found the event, parse it
					c.logger.Debug().Int64("Offset", createdEvent.GetOffset()).Str("ContractId", createdEvent.GetContractId()).Msg("Found ExecutionStateChanged event")
					parsedEvent, err := parseExecutionStateChangedEvent(createdEvent)
					if err != nil {
						return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to parse ExecutionStateChanged event: %w", err)
					}

					return parsedEvent, nil
				}
			}
		}
	}

	// No event found in the update, return an error
	return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("no ExecutionStateChanged event found in update %s", update.GetTransaction().GetUpdateId())
}

// parseExecutionStateChangedEvent parses a common.ExecutionStateChanged event from a Daml CreatedEvent and converts it to cciptestinterfaces.ExecutionStateChangedEvent.
func parseExecutionStateChangedEvent(event *apiv2.CreatedEvent) (cciptestinterfaces.ExecutionStateChangedEvent, error) {
	executionStateChanged, err := bindings.UnmarshalCreatedEvent[core.ExecutionStateChanged](event)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to unmarshal ExecutionStateChanged event: %w", err)
	}

	// Source chain selector
	sourceChainSelectorFloat, ok := new(big.Float).SetString(string(executionStateChanged.Event.SourceChainSelector))
	if !ok {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to parse source chain selector numeric, input: %s", string(executionStateChanged.Event.SourceChainSelector))
	}
	sourceChainSelector, _ := sourceChainSelectorFloat.Int(nil)
	// Message ID
	messageId, err := hex.DecodeString(string(executionStateChanged.Event.MessageId))
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to decode message ID %q: %w", string(executionStateChanged.Event.MessageId), err)
	}
	// Message number
	sequenceNumberFloat, ok := new(big.Float).SetString(string(executionStateChanged.Event.SequenceNumber))
	if !ok {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to parse sequence number numeric, input: %s", string(executionStateChanged.Event.SequenceNumber))
	}
	sequenceNumber, _ := sequenceNumberFloat.Int(nil)
	// Execution state
	var executionState cciptestinterfaces.MessageExecutionState
	switch executionStateChanged.Event.State {
	case core.MessageExecutionStateUNTOUCHED:
		executionState = cciptestinterfaces.ExecutionStateUntouched
	case core.MessageExecutionStateIN_PROGRESS:
		executionState = cciptestinterfaces.ExecutionStateInProgress
	case core.MessageExecutionStateSUCCESS:
		executionState = cciptestinterfaces.ExecutionStateSuccess
	case core.MessageExecutionStateFAILURE:
		executionState = cciptestinterfaces.ExecutionStateFailure
	default:
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("unknown execution state %q", executionStateChanged.Event.State)
	}
	// Return data
	returnData, err := hex.DecodeString(string(executionStateChanged.Event.ReturnData))
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to decode return data %q: %w", string(executionStateChanged.Event.ReturnData), err)
	}

	return cciptestinterfaces.ExecutionStateChangedEvent{
		SourceChainSelector: protocol.ChainSelector(sourceChainSelector.Uint64()),
		MessageID:           [32]byte(messageId),
		MessageNumber:       sequenceNumber.Uint64(),
		State:               executionState,
		ReturnData:          returnData,
	}, nil
}

// lockForParty serializes manual-execute work per receiver party. Two ConfirmExecOnDest
// invocations targeting the same party will queue, because the receiver's PerPartyRouter
// contract is consumed and recreated on every execute (line ~146 of this file): a parallel
// caller would race on routerCid resolution and one would fail with "already archived".
func (c *Chain) lockForParty(party string) func() {
	c.partyMutexesMu.Lock()
	m, ok := c.partyMutexes[party]
	if !ok {
		m = &sync.Mutex{}
		c.partyMutexes[party] = m
	}
	c.partyMutexesMu.Unlock()
	m.Lock()

	return m.Unlock
}

// findExistingExecutionState scans active ExecutionStateChanged contracts on the receiver
// party and returns the parsed event matching (sourceChainSelector, seqNo, messageID) if
// present. Used by ConfirmExecOnDest for idempotency: if a message has already been
// executed, repeated calls return the same event instead of attempting to re-execute
// (which would fail because the underlying PerPartyRouter contract has been consumed).
func (c *Chain) findExistingExecutionState(
	ctx context.Context, sourceChainSelector, seqNo uint64, messageID protocol.Bytes32,
) (cciptestinterfaces.ExecutionStateChangedEvent, bool, error) {
	if len(c.chain.Participants) == 0 {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, false, fmt.Errorf("findExistingExecutionState: no participants on chain")
	}
	participant := c.chain.Participants[0]

	templateID := contracts.TemplateIDFromBinding(core.ExecutionStateChanged{}).ToLedgerIdentifier()
	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, templateID)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, false, fmt.Errorf("findExistingExecutionState: list active ExecutionStateChanged contracts: %w", err)
	}
	for _, ac := range activeContracts {
		created := ac.GetCreatedEvent()
		if created == nil {
			continue
		}
		ev, err := parseExecutionStateChangedEvent(created)
		if err != nil {
			c.logger.Debug().Err(err).Msg("Skipping unparseable ExecutionStateChanged active contract")
			continue
		}
		if uint64(ev.SourceChainSelector) == sourceChainSelector &&
			ev.MessageNumber == seqNo &&
			ev.MessageID == messageID {
			return ev, true, nil
		}
	}

	return cciptestinterfaces.ExecutionStateChangedEvent{}, false, nil
}

// verifierResult is a minimal view of a CCIP verifier result, with the verifier
// destination address already translated to the hashed Canton instance address.
type verifierResult struct {
	Message                protocol.Message
	HashedVerifierDestAddr protocol.UnknownAddress
	CCVData                []byte
}

// verifierAssertTimeout maps ConfirmExecOnDest's timeout to the verifier poll window.
func verifierAssertTimeout(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return 5 * time.Minute
	}

	return timeout
}

// fetchVerifierResult queries the indexer (aggregator optional) for the verifier
// output for messageID. Caller must have wired [VerifierObservation] on the chain first.
func (c *Chain) fetchVerifierResult(ctx context.Context, messageID protocol.Bytes32, timeout time.Duration) (verifierResult, error) {
	if !c.verifierObs.wired() {
		return verifierResult{}, fmt.Errorf("verifier observation not wired")
	}

	res, err := AssertMessageWithVerifierObservation(ctx, c.verifierObs, messageID, tcapi.AssertMessageOptions{
		TickInterval:            time.Second,
		Timeout:                 verifierAssertTimeout(timeout),
		ExpectedVerifierResults: 1,
		AssertVerifierLogs:      false,
		AssertExecutorLogs:      false,
	})
	if err != nil {
		return verifierResult{}, fmt.Errorf("assertMessage: %w", err)
	}
	if len(res.IndexedVerifications.Results) != 1 {
		return verifierResult{}, fmt.Errorf("expected 1 indexed verifier result, got %d", len(res.IndexedVerifications.Results))
	}
	vr := res.IndexedVerifications.Results[0].VerifierResult

	hashedDestAddr, err := hashInstanceAddress(vr.VerifierDestAddress)
	if err != nil {
		return verifierResult{}, fmt.Errorf("hashInstanceAddress: %w", err)
	}

	return verifierResult{
		Message:                vr.Message,
		HashedVerifierDestAddr: hashedDestAddr,
		CCVData:                vr.CCVData,
	}, nil
}

func encodeReceiverFinalityConfig(finality int64) (core.FinalityConfig, error) {
	switch {
	case finality < 0:
		return core.FinalityConfig{}, fmt.Errorf("invalid finality %d: must be non-negative", finality)
	case finality == 0:
		return core.FinalityConfig{WaitForFinality: &types.UNIT{}}, nil
	case finality == 0x00010000:
		return core.FinalityConfig{WaitForSafe: &types.UNIT{}}, nil
	case finality > 0xFFFF:
		return core.FinalityConfig{}, fmt.Errorf("invalid finality %d: max supported block depth is 65535", finality)
	default:
		return core.FinalityConfig{BlockDepth: new(types.INT64(finality))}, nil
	}
}

// hashInstanceAddress resolves a verifier result's VerifierDestAddress to the hashed
// Canton instance address used for EDS disclosure lookup.
//
// The indexer may return either format depending on environment:
//   - devenv: raw instance address string bytes ("instanceId@owner")
//   - prod-testnet: already-hashed 32-byte InstanceAddress (from verifier config)
func hashInstanceAddress(addr protocol.UnknownAddress) (protocol.UnknownAddress, error) {
	if len(addr) == 0 {
		return nil, fmt.Errorf("empty verifier dest address")
	}

	if raw, err := contracts.RawInstanceAddressFromString(string(addr)); err == nil {
		return protocol.UnknownAddress(raw.InstanceAddress().Bytes()), nil
	}

	if len(addr) == contracts.InstanceAddressLength {
		return addr, nil
	}

	if hexAddr, err := protocol.NewUnknownAddressFromHex(string(addr)); err == nil && len(hexAddr) == contracts.InstanceAddressLength {
		return hexAddr, nil
	}

	return nil, fmt.Errorf("unrecognized verifier dest address format (len=%d)", len(addr))
}
