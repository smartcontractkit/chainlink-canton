package devenv

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"

	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

func receiverWaitForFinalityConfig() common.FinalityConfig {
	return common.FinalityConfig{WaitForFinality: &types.UNIT{}}
}

func receiverWaitForSafeConfig() common.FinalityConfig {
	return common.FinalityConfig{WaitForSafe: &types.UNIT{}}
}

func encodeReceiverFinalityConfig(finality int64) (common.FinalityConfig, error) {
	switch {
	case finality < 0:
		return common.FinalityConfig{}, fmt.Errorf("invalid finality %d: must be non-negative", finality)
	case finality == 0:
		return receiverWaitForFinalityConfig(), nil
	case finality == 0x00010000:
		return receiverWaitForSafeConfig(), nil
	case finality > 0xFFFF:
		return common.FinalityConfig{}, fmt.Errorf("invalid finality %d: max supported block depth is 65535", finality)
	default:
		return common.FinalityConfig{BlockDepth: new(types.INT64(finality))}, nil
	}
}

// DeployPerPartyRouter uses the PerPartyRouterFactory to create a new PerPartyRouter instance for the given party.
// It returns the address of the newly created PerPartyRouter instance. If a router already exists for the party, it returns the existing router's address.
func (c *Chain) DeployPerPartyRouter(ctx context.Context, participant canton.Participant, partyId string) (routerAddress contracts.InstanceAddress, err error) {
	// Create PerPartyRouter (ignore error if it exists already)
	perPartyRouterFactoryDisclosure, err := c.GetPerPartyRouterFactoryDisclosure(ctx, partyId)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to get canton per party router factory disclosure: %w", err)
	}
	c.logger.Debug().Str("ContractId", perPartyRouterFactoryDisclosure.ContractId).Msg("Resolved per-party router factory address")

	// Fixed instance ID for the router, this makes the InstanceAddress deterministic.
	routerInstanceID := contracts.InstanceID("test-router")
	// Ignore errors, since the router might already exist if this function is called multiple times for the same party. In that case we just want to return the existing router's address.
	_, _ = operations.ExecuteOperation(c.e.OperationsBundle, per_party_router_factory.CreateRouter, c.chain, contract.ChoiceInput[perpartyrouter.CreateRouter]{
		InstanceAddress: perPartyRouterFactoryDisclosure.Address.InstanceAddress(),
		Args: perpartyrouter.CreateRouter{
			PartyOwner: types.PARTY(partyId),
			InstanceId: types.TEXT(routerInstanceID.String()),
		},
	})
	routerAddress = routerInstanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()

	// Try to get the PerPartyRouter, to ensure that is actually exists
	_, err = testhelpers.GetDisclosedContractByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	})
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to find per-party router: %w", err)
	}

	return routerAddress, nil
}

func (c *Chain) DeployCCIPReceiver(partyId string, receiverFinality int64) (contracts.InstanceAddress, error) {
	finalityConfig, err := encodeReceiverFinalityConfig(receiverFinality)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to encode receiver finality config: %w", err)
	}

	// Deploy receiver contract
	out, err := operations.ExecuteOperation(c.e.OperationsBundle, receiver.Deploy, c.chain, contract.DeployInput[ccipreceiver.CCIPReceiver]{
		Qualifier: nil,
		Template: ccipreceiver.CCIPReceiver{
			Owner:                  types.PARTY(partyId),
			RequiredCCVs:           nil,
			OptionalCCVs:           nil,
			OptionalThreshold:      0,
			ReceiverFinalityConfig: finalityConfig,
		},
		OwnerParty: types.PARTY(partyId),
	})
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to deploy receiver contract: %w", err)
	}
	receiverAddress := contracts.HexToInstanceAddress(out.Output.Address)

	return receiverAddress, nil
}

// ManuallyExecuteMessage implements cciptestinterfaces.CCIP17.
func (c *Chain) ManuallyExecuteMessage(ctx context.Context, message protocol.Message, gasLimit uint64, verifiers []protocol.UnknownAddress, verifierResults [][]byte) (cciptestinterfaces.ExecutionStateChangedEvent, error) {
	// Use only a single participant for now
	participant := c.chain.Participants[0]

	// Ensure that the message receiver is the party we're executing with
	executingParty := participant.PartyID
	if contracts.HashedPartyFromString(executingParty) != contracts.BytesToHashedParty(message.Receiver.Bytes()) {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("message receiver %s does not match executing party %s (%s)", hex.EncodeToString(message.Receiver), contracts.HashedPartyFromString(executingParty).String(), executingParty)
	}

	// Deploy PerPartyRouter for the receiver party
	routerAddress, err := c.DeployPerPartyRouter(ctx, participant, executingParty)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to deploy per-party router: %w", err)
	}
	c.logger.Debug().Str("RouterAddress", routerAddress.String()).Msg("Deployed PerPartyRouter")

	// Deploy CCIPReceiver contract
	receiverAddress, err := c.DeployCCIPReceiver(executingParty, int64(message.Finality))
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

	routerCid, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, perpartyrouter.PerPartyRouter{}.GetTemplateID(), routerAddress)
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
	executeArgs := ccipreceiver.Execute2{
		Context:        ccipExecuteDisclosure.ChoiceContext,
		RouterCid:      types.CONTRACT_ID(routerCid),
		EncodedMessage: types.TEXT(encodedMessageHex),
		TokenTransfer:  nil,
		CcvInputs:      make([]ccipreceiver.CCVInput, len(verifiers)),
	}
	disclosedContracts := ccipExecuteDisclosure.DisclosedContracts

	// CCVs
	for i, verifier := range ccvs {
		ccvExecuteDisclosure, err := c.GetCCVExecuteDisclosure(ctx, encodedMessageHex, verifier)
		if err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get CCV execute disclosure for verifier %s: %w", verifier.String(), err)
		}

		executeArgs.CcvInputs[i] = ccipreceiver.CCVInput{
			CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
			VerifierResults: types.TEXT(hex.EncodeToString(verifierResults[i])),
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
			TokenInput:         tokenPoolDisclosure.TokenInput,
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

	executeReport, err := operations.ExecuteOperation(c.e.OperationsBundle, receiver.Execute, c.chain, contract.ChoiceInput[ccipreceiver.Execute2]{
		InstanceAddress:    receiverAddress,
		Args:               executeArgs,
		DisclosedContracts: contract.DisclosedContractsFromProto(disclosedContracts),
	})
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to execute message: %w", err)
	}
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
	expectedTemplateID := common.ExecutionStateChanged{}.GetTemplateID()
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
	executionStateChanged, err := bindings.UnmarshalCreatedEvent[common.ExecutionStateChanged](event)
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
	case common.MessageExecutionStateUNTOUCHED:
		executionState = cciptestinterfaces.ExecutionStateUntouched
	case common.MessageExecutionStateIN_PROGRESS:
		executionState = cciptestinterfaces.ExecutionStateInProgress
	case common.MessageExecutionStateSUCCESS:
		executionState = cciptestinterfaces.ExecutionStateSuccess
	case common.MessageExecutionStateFAILURE:
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
