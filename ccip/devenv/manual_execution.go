package devenv

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"

	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

// DeployPerPartyRouter uses the PerPartyRouterFactory to create a new PerPartyRouter instance for the given party.
// It returns the address of the newly created PerPartyRouter instance. If a router already exists for the party, it returns the existing router's address.
func (c *Chain) DeployPerPartyRouter(ctx context.Context, partyId string) (contracts.InstanceAddress, error) {
	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: 0,
	}

	// Create PerPartyRouter (ignore error if it exists already)
	cantonPerPartyRouterFactoryRef, err := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(per_party_router_factory.ContractType),
			per_party_router_factory.Version,
			"",
		),
	)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to get canton per party router factory address ref: %w", err)
	}
	c.logger.Debug().Str("CantonPerPartyRouterFactory", cantonPerPartyRouterFactoryRef.Address).Msg("Resolved per-party router factory address")
	cantonPerPartyRouterFactory := contracts.HexToInstanceAddress(cantonPerPartyRouterFactoryRef.Address)

	// Fixed instance ID for the router, this makes the InstanceAddress deterministic.
	routerInstanceID := contracts.InstanceID("test-router")
	// Ignore errors, since the router might already exist if this function is called multiple times for the same party. In that case we just want to return the existing router's address.
	_, _ = operations.ExecuteOperation(c.e.OperationsBundle, per_party_router_factory.CreateRouter, deps, contract.ChoiceInput[perpartyrouter.CreateRouter]{
		ChainSelector:   c.chainDetails.ChainSelector,
		InstanceAddress: cantonPerPartyRouterFactory,
		ActAs:           []string{partyId},
		Args: perpartyrouter.CreateRouter{
			PartyOwner: types.PARTY(partyId),
			InstanceId: types.TEXT(routerInstanceID.String()),
		},
	})
	routerAddress := routerInstanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()

	return routerAddress, nil
}

func (c *Chain) DeployCCIPReceiver(ctx context.Context, partyId string) (contracts.InstanceAddress, error) {
	// Use only a single participant for now
	participant := c.chain.Participants[0]
	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: 0,
	}

	// Upload the necessary Dar
	receiverDar, err := contracts.GetDar(contracts.CCIPReceiver, contracts.CurrentVersion)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to get receiver dar: %w", err)
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile:       receiverDar,
		VettingChange: adminv2.UploadDarFileRequest_VETTING_CHANGE_VET_ALL_PACKAGES,
	})
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to upload receiver dar file: %w", err)
	}

	// Deploy receiver contract
	out, err := operations.ExecuteOperation(c.e.OperationsBundle, receiver.Deploy, deps, contract.DeployInput[ccipreceiver.CCIPReceiver]{
		ChainSelector: c.chainDetails.ChainSelector,
		Qualifier:     nil,
		ActAs:         []string{participant.PartyID},
		Template: ccipreceiver.CCIPReceiver{
			Owner:        types.PARTY(participant.PartyID),
			RequiredCCVs: nil,
		},
		OwnerParty: types.PARTY(participant.PartyID),
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
	routerAddress, err := c.DeployPerPartyRouter(ctx, executingParty)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to deploy per-party router: %w", err)
	}
	c.logger.Debug().Str("RouterAddress", routerAddress.String()).Msg("Deployed PerPartyRouter")

	// Deploy CCIPReceiver contract
	receiverAddress, err := c.DeployCCIPReceiver(ctx, executingParty)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to deploy CCIPReceiver contract: %w", err)
	}
	c.logger.Debug().Str("ReceiverAddress", receiverAddress.String()).Msg("Deployed CCIPReceiver")

	// Get disclosures for execution using EDS
	ccvs := make([]contracts.InstanceAddress, len(verifiers))
	for i, verifier := range verifiers {
		ccvs[i] = contracts.HexToInstanceAddress(verifier.String())
	}
	disclosedContracts, choiceContext, ccvContractIDs, err := c.GetDisclosuresForExecution(ctx, ccvs)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get disclosures for execution: %w", err)
	}

	// Resolve all necessary contracts
	resolveActiveContractIDByAddress := func(templateID string, address contracts.InstanceAddress) (string, error) {
		cid, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, templateID, address)
		if err == nil {
			return cid, nil
		}
		if !strings.Contains(err.Error(), "multiple active contracts found") {
			return "", err
		}

		parts := strings.Split(templateID, ":")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid template ID %q", templateID)
		}
		ledgerEnd, endErr := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
		if endErr != nil {
			return "", fmt.Errorf("get ledger end for duplicate contract fallback: %w", endErr)
		}
		stream, streamErr := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
			ActiveAtOffset: ledgerEnd.GetOffset(),
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					participant.PartyID: {
						Cumulative: []*apiv2.CumulativeFilter{
							{
								IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
									TemplateFilter: &apiv2.TemplateFilter{
										TemplateId: &apiv2.Identifier{
											PackageId:  parts[0],
											ModuleName: parts[1],
											EntityName: parts[2],
										},
										IncludeCreatedEventBlob: true,
									},
								},
							},
						},
					},
				},
				Verbose: true,
			},
		})
		if streamErr != nil {
			return "", fmt.Errorf("get active contracts for duplicate contract fallback: %w", streamErr)
		}
		defer stream.CloseSend()
		var latest *apiv2.ActiveContract
		for {
			resp, recvErr := stream.Recv()
			if recvErr != nil {
				if errors.Is(recvErr, io.EOF) {
					break
				}
				return "", fmt.Errorf("receive active contracts for duplicate contract fallback: %w", recvErr)
			}
			entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
				continue
			}
			if latest == nil || entry.ActiveContract.GetCreatedEvent().GetOffset() > latest.GetCreatedEvent().GetOffset() {
				latest = entry.ActiveContract
			}
		}
		if latest == nil || latest.GetCreatedEvent() == nil || latest.GetCreatedEvent().GetContractId() == "" {
			return "", fmt.Errorf("no active contracts found for duplicate fallback")
		}

		return latest.GetCreatedEvent().GetContractId(), nil
	}

	routerCid, err := resolveActiveContractIDByAddress(perpartyrouter.PerPartyRouter{}.GetTemplateID(), routerAddress)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get router contract ID: %w", err)
	}
	c.logger.Debug().Str("InstanceAddress", routerAddress.String()).Str("ContractId", routerCid).Msg("Resolved PerPartyRouter contract")

	receiverCid, err := resolveActiveContractIDByAddress(ccipreceiver.CCIPReceiver{}.GetTemplateID(), receiverAddress)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get receiver contract ID: %w", err)
	}
	c.logger.Debug().Str("InstanceAddress", receiverAddress.String()).Str("ContractId", receiverCid).Msg("Resolved CCIPReceiver contract")

	emptyCCIPCtx := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
	}}}}
	ccvElements := make([]*apiv2.Value, len(verifiers))
	for i, ccvCid := range ccvContractIDs {
		ccvElements[i] = &apiv2.Value{
			Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "ccvCid", Value: &apiv2.Value{Sum: ccvCid}},
				{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(verifierResults[i])}}},
				{Label: "ccvExtraContext", Value: emptyCCIPCtx},
			}}},
		}
	}
	executeDisclosures := disclosedContracts
	tokenTransferValue := &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}
	if message.TokenTransfer != nil {
		var tokenDisclosures []*apiv2.DisclosedContract
		tokenTransferValue, tokenDisclosures, err = c.buildManualExecuteTokenTransferInput(
			ctx,
			participant,
			executingParty,
			message,
			resolveActiveContractIDByAddress,
		)
		if err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("build token transfer execute input: %w", err)
		}
		executeDisclosures = append(executeDisclosures, tokenDisclosures...)
	}
	// Execute message
	encodedMessage, err := message.Encode()
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to encode message: %w", err)
	}
	c.logger.Debug().
		Str("EncodedMessage", hex.EncodeToString(encodedMessage)).
		Str("VerifierResults", hex.EncodeToString(verifierResults[0])).
		Str("Receiver", hex.EncodeToString(message.Receiver)).
		Msg("Executing message...")

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId: receiverCid,
					Choice:     "Execute",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "context", Value: choiceContext},
						{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCid}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(encodedMessage)}}},
						{Label: "tokenTransfer", Value: tokenTransferValue},
						{Label: "ccvInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: ccvElements}}}},
						{Label: "additionalRequiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}}}},
				}},
			}},
			ActAs:              []string{executingParty},
			DisclosedContracts: executeDisclosures,
		},
	})
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to execute message: %w", err)
	}
	c.logger.Debug().Str("UpdateID", res.GetTransaction().GetUpdateId()).Msg("Executed message")

	// Get Update
	updateRes, err := participant.LedgerServices.Update.GetUpdateById(ctx, &apiv2.GetUpdateByIdRequest{
		UpdateId: res.GetTransaction().GetUpdateId(),
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
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get update by UpdateId %q: %w", res.GetTransaction().GetUpdateId(), err)
	}

	// Get ExecutionStateChangedEvent from events
	expectedTemplateID := common.ExecutionStateChanged{}.GetTemplateID()
	for _, event := range updateRes.GetTransaction().GetEvents() {
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
	return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("no ExecutionStateChanged event found in update %s", res.GetTransaction().GetUpdateId())
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

func (c *Chain) buildManualExecuteTokenTransferInput(
	ctx context.Context,
	participant canton.Participant,
	executingParty string,
	message protocol.Message,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (*apiv2.Value, []*apiv2.DisclosedContract, error) {
	sourceSelectorKey := fmt.Sprintf("%d", message.SourceChainSelector)
	sourceSelectorNumericKey := sourceSelectorKey + "."
	sourcePoolHex := strings.ToLower(hex.EncodeToString(message.TokenTransfer.SourcePoolAddress))
	if len(sourcePoolHex) > 40 {
		sourcePoolHex = sourcePoolHex[len(sourcePoolHex)-40:]
	}
	sourcePoolHex = strings.TrimPrefix(sourcePoolHex, "0x")
	destTokenHex := strings.ToLower(hex.EncodeToString(message.TokenTransfer.DestTokenAddress))
	destTokenHex = strings.TrimPrefix(destTokenHex, "0x")
	requireInstrumentMatch := len(destTokenHex) > 40

	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("get ledger end for token pool lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &apiv2.TemplateFilter{
									TemplateId: &apiv2.Identifier{
										PackageId:  "#ccip-lockreleasetokenpool",
										ModuleName: "CCIP.LockReleaseTokenPool",
										EntityName: "LockReleaseTokenPool",
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("query active lock/release pools: %w", err)
	}
	defer stream.CloseSend()

	var selectedPool *lockreleasetokenpool.LockReleaseTokenPool
	var selectedPoolContractID string
	var tokenMatchedPool *lockreleasetokenpool.LockReleaseTokenPool
	var tokenMatchedPoolContractID string
	var fallbackPool *lockreleasetokenpool.LockReleaseTokenPool
	var fallbackPoolContractID string
	debugPoolCandidates := make([]string, 0, 8)
	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, nil, fmt.Errorf("receive lock/release pools: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}
		parsed, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](entry.ActiveContract.GetCreatedEvent())
		if parseErr != nil {
			return nil, nil, fmt.Errorf("parse lock/release pool: %w", parseErr)
		}
		chainPoolCfgAny, ok := findChainPoolConfigBySelector(parsed.ChainPoolConfigs, sourceSelectorKey)
		if !ok {
			continue
		}
		remoteTokenMatch := false
		instrumentCombined := fmt.Sprintf("%s@%s", string(parsed.InstrumentId.Id), string(parsed.InstrumentId.Admin))
		instrumentRawHex := strings.ToLower(hex.EncodeToString([]byte(instrumentCombined)))
		instrumentKeccakCombinedHex := strings.ToLower(hex.EncodeToString(crypto.Keccak256([]byte(instrumentCombined))))
		instrumentKeccakTextHex := strings.ToLower(hex.EncodeToString(crypto.Keccak256([]byte(instrumentRawHex))))
		instrumentTokenMatch := instrumentRawHex == destTokenHex ||
			instrumentKeccakCombinedHex == destTokenHex ||
			instrumentKeccakTextHex == destTokenHex ||
			strings.HasSuffix(instrumentRawHex, destTokenHex) ||
			strings.HasSuffix(destTokenHex, instrumentRawHex) ||
			strings.HasSuffix(instrumentKeccakCombinedHex, destTokenHex) ||
			strings.HasSuffix(destTokenHex, instrumentKeccakCombinedHex) ||
			strings.HasSuffix(instrumentKeccakTextHex, destTokenHex) ||
			strings.HasSuffix(destTokenHex, instrumentKeccakTextHex)
		if instrumentTokenMatch && tokenMatchedPool == nil {
			tokenMatchedPool = parsed
			tokenMatchedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
		}
		remoteTokenAny, ok := parsed.RemoteTokens[sourceSelectorKey]
		if !ok {
			remoteTokenAny, ok = parsed.RemoteTokens[sourceSelectorNumericKey]
		}
		if ok {
			remoteTokenHex := strings.ToLower(strings.TrimPrefix(fmt.Sprint(remoteTokenAny), "0x"))
			if remoteTokenHex == destTokenHex || strings.HasSuffix(remoteTokenHex, destTokenHex) || strings.HasSuffix(destTokenHex, remoteTokenHex) {
				remoteTokenMatch = true
			}
		}
		if !remoteTokenMatch && strings.Contains(strings.ToLower(fmt.Sprint(parsed.RemoteTokens)), destTokenHex) {
			remoteTokenMatch = true
		}
		remotePools := extractRemotePools(chainPoolCfgAny)
		if len(debugPoolCandidates) < 8 {
			debugPoolCandidates = append(debugPoolCandidates, fmt.Sprintf(
				"poolCID=%s instrumentRaw=%s instrumentKeccak=%s remotePools=%v",
				entry.ActiveContract.GetCreatedEvent().GetContractId(),
				instrumentRawHex,
				instrumentKeccakCombinedHex,
				remotePools,
			))
		}
		remotePoolMatch := false
		if len(remotePools) == 0 {
			if requireInstrumentMatch {
				continue
			}
			remotePoolMatch = strings.Contains(strings.ToLower(fmt.Sprint(chainPoolCfgAny)), sourcePoolHex)
		}
		if len(remotePools) == 0 {
			if remoteTokenMatch && tokenMatchedPool == nil {
				tokenMatchedPool = parsed
				tokenMatchedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
			}
			if fallbackPool == nil {
				fallbackPool = parsed
				fallbackPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
			}
			continue
		}
		for _, remotePool := range remotePools {
			remotePoolHex := strings.ToLower(strings.TrimPrefix(remotePool, "0x"))
			if remotePoolHex == sourcePoolHex || strings.HasSuffix(remotePoolHex, sourcePoolHex) || strings.HasSuffix(sourcePoolHex, remotePoolHex) {
				remotePoolMatch = true
				break
			}
		}
		if !remotePoolMatch {
			if remoteTokenMatch && tokenMatchedPool == nil {
				tokenMatchedPool = parsed
				tokenMatchedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
			}
			if fallbackPool == nil {
				fallbackPool = parsed
				fallbackPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
			}
			continue
		}
		if instrumentTokenMatch {
			selectedPool = parsed
			selectedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
			break
		}
		if requireInstrumentMatch {
			continue
		}
		selectedPool = parsed
		selectedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
		break
	}
	if selectedPool == nil && tokenMatchedPool != nil && !requireInstrumentMatch {
		selectedPool = tokenMatchedPool
		selectedPoolContractID = tokenMatchedPoolContractID
	}
	if selectedPool == nil && tokenMatchedPool != nil && requireInstrumentMatch {
		updatedPool, updatedCID, ensureErr := ensureManualExecuteSourcePoolAllowed(
			ctx,
			participant,
			tokenMatchedPool,
			tokenMatchedPoolContractID,
			sourceSelectorKey,
			sourcePoolHex,
			resolveActiveContractIDByAddress,
		)
		if ensureErr == nil {
			selectedPool = updatedPool
			selectedPoolContractID = updatedCID
		}
	}
	if selectedPool == nil && fallbackPool != nil && !requireInstrumentMatch {
		selectedPool = fallbackPool
		selectedPoolContractID = fallbackPoolContractID
	}
	if selectedPool == nil {
		if requireInstrumentMatch {
			return nil, nil, fmt.Errorf(
				"no lock/release pool found with instrument hash matching dest token %s for source selector %s and source pool %s; candidates=%v",
				destTokenHex,
				sourceSelectorKey,
				sourcePoolHex,
				debugPoolCandidates,
			)
		}
		return nil, nil, fmt.Errorf("no lock/release pool found for source selector %s and source pool %s", sourceSelectorKey, sourcePoolHex)
	}

	rateLimiterCID, rateLimiterDisclosure, err := resolveRateLimiterForManualExecute(
		ctx,
		participant,
		selectedPool,
		sourceSelectorKey,
		sourceSelectorNumericKey,
		resolveActiveContractIDByAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve rate limiter disclosure: %w", err)
	}

	transferClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(func(reqCtx context.Context, req *http.Request) error {
			token, tokenErr := participant.TokenSource.Token()
			if tokenErr != nil {
				return fmt.Errorf("retrieve participant token: %w", tokenErr)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			return nil
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create transfer instruction client: %w", err)
	}
	transferFactoryResp, err := transferClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": string(selectedPool.InstrumentId.Admin),
			"transfer": map[string]any{
				"sender":   string(selectedPool.PoolOwner),
				"receiver": executingParty,
				"amount":   message.TokenTransfer.Amount.String(),
				"instrumentId": map[string]any{
					"admin": string(selectedPool.InstrumentId.Admin),
					"id":    string(selectedPool.InstrumentId.Id),
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"inputHoldingCids": []string{},
				"meta":             map[string]any{"values": map[string]any{}},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{"values": map[string]any{}},
				"meta":    map[string]any{"values": map[string]any{}},
			},
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("get transfer factory: %w", err)
	}
	if transferFactoryResp.StatusCode() != http.StatusOK {
		return nil, nil, fmt.Errorf("transfer factory response status %d: %s", transferFactoryResp.StatusCode(), string(transferFactoryResp.Body))
	}
	transferFactoryCID := transferFactoryResp.JSON200.FactoryId
	transferFactoryCtx, err := ChoiceContextFromData(transferFactoryResp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return nil, nil, fmt.Errorf("convert transfer factory context: %w", err)
	}
	transferFactoryDisclosures := make([]*apiv2.DisclosedContract, 0, len(transferFactoryResp.JSON200.ChoiceContext.DisclosedContracts))
	for _, d := range transferFactoryResp.JSON200.ChoiceContext.DisclosedContracts {
		id, idErr := TemplateIdFromString(d.TemplateId)
		if idErr != nil {
			return nil, nil, fmt.Errorf("parse transfer factory disclosure template id: %w", idErr)
		}
		createdEventBlob, decodeErr := base64.StdEncoding.DecodeString(d.CreatedEventBlob)
		if decodeErr != nil {
			return nil, nil, fmt.Errorf("decode transfer factory disclosure created event blob: %w", decodeErr)
		}
		transferFactoryDisclosures = append(transferFactoryDisclosures, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       d.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   d.SynchronizerId,
		})
	}

	collectPoolHoldings := func() ([]string, []*apiv2.DisclosedContract, error) {
		holdings, holdingsErr := listHoldingContracts(ctx, participant)
		if holdingsErr != nil {
			return nil, nil, fmt.Errorf("list pool holdings: %w", holdingsErr)
		}
		poolHoldings := make([]string, 0, len(holdings))
		poolHoldingDisclosures := make([]*apiv2.DisclosedContract, 0, len(holdings))
		for _, holding := range holdings {
			views := holding.GetCreatedEvent().GetInterfaceViews()
			if len(views) == 0 || views[0].GetViewValue() == nil {
				continue
			}
			fields := views[0].GetViewValue().GetFields()
			if len(fields) < 4 {
				continue
			}
			ownerParty := fields[0].GetValue().GetParty()
			instrumentRecord := fields[1].GetValue().GetRecord()
			if instrumentRecord == nil || len(instrumentRecord.GetFields()) < 2 {
				continue
			}
			instrumentAdmin := instrumentRecord.GetFields()[0].GetValue().GetParty()
			instrumentID := instrumentRecord.GetFields()[1].GetValue().GetText()
			isLocked := fields[3].GetValue().GetOptional().GetValue() != nil
			if isLocked {
				continue
			}
			if ownerParty != string(selectedPool.PoolOwner) {
				continue
			}
			if instrumentAdmin != string(selectedPool.InstrumentId.Admin) || instrumentID != string(selectedPool.InstrumentId.Id) {
				continue
			}
			poolHoldings = append(poolHoldings, holding.GetCreatedEvent().GetContractId())
			poolHoldingDisclosures = append(poolHoldingDisclosures, &apiv2.DisclosedContract{
				TemplateId:       holding.GetCreatedEvent().GetTemplateId(),
				ContractId:       holding.GetCreatedEvent().GetContractId(),
				CreatedEventBlob: holding.GetCreatedEvent().GetCreatedEventBlob(),
				SynchronizerId:   holding.GetSynchronizerId(),
			})
		}
		return poolHoldings, poolHoldingDisclosures, nil
	}

	poolHoldings, poolHoldingDisclosures, err := collectPoolHoldings()
	if err != nil {
		return nil, nil, err
	}
	if len(poolHoldings) == 0 {
		// Devenv token-transfer tests may not pre-seed release liquidity; create a minimal pool-owned holding.
		_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.New().String(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{PackageId: "#ccip-test", ModuleName: "TestToken", EntityName: "TestHolding"},
						CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(selectedPool.PoolOwner)}}},
							{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(selectedPool.InstrumentId.Admin)}}},
							{Label: "instrumentId", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(selectedPool.InstrumentId.Admin)}}},
								{Label: "id", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(selectedPool.InstrumentId.Id)}}},
							}}}}},
							{Label: "amount", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: message.TokenTransfer.Amount.String()}}},
						}},
					}},
				}},
				ActAs: []string{string(selectedPool.PoolOwner), string(selectedPool.InstrumentId.Admin)},
			},
		})
		if err != nil {
			return nil, nil, fmt.Errorf("create pool holding for token execution liquidity: %w", err)
		}
		poolHoldings, poolHoldingDisclosures, err = collectPoolHoldings()
		if err != nil {
			return nil, nil, err
		}
		if len(poolHoldings) == 0 {
			return nil, nil, fmt.Errorf("no unlocked pool holdings found for pool owner %s", selectedPool.PoolOwner)
		}
	}

	poolExtraContext, err := ChoiceContextFromData(map[string]any{
		"values": map[string]any{
			"rate-limiter": map[string]any{
				"tag":   "AV_ContractId",
				"value": rateLimiterCID,
			},
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("build pool extra context: %w", err)
	}
	emptyMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
	}}}}
	tokenPoolHoldingValues := make([]*apiv2.Value, 0, len(poolHoldings))
	for _, cid := range poolHoldings {
		tokenPoolHoldingValues = append(tokenPoolHoldingValues, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}})
	}

	tokenTransferValue := &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "tokenPoolCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: selectedPoolContractID}}},
		{Label: "tokenReceiverParty", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: executingParty}}},
		{Label: "tokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
			{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCID}}},
			{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "context", Value: transferFactoryCtx},
				{Label: "meta", Value: emptyMetadata},
			}}}}},
			{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: tokenPoolHoldingValues}}}},
		}}}}},
		{Label: "poolExtraContext", Value: poolExtraContext},
	}}}}}}}

	disclosures := []*apiv2.DisclosedContract{rateLimiterDisclosure}
	disclosures = append(disclosures, transferFactoryDisclosures...)
	disclosures = append(disclosures, poolHoldingDisclosures...)

	return tokenTransferValue, disclosures, nil
}

func extractRemotePools(chainPoolCfgAny any) []string {
	switch cfg := chainPoolCfgAny.(type) {
	case lockreleasetokenpool.ChainPoolConfig:
		out := make([]string, 0, len(cfg.RemotePools))
		for _, rp := range cfg.RemotePools {
			out = append(out, string(rp))
		}
		return out
	case map[string]any:
		m := cfg
		if data, ok := cfg["data"].(map[string]any); ok {
			m = data
		}
		raw, ok := m["remotePools"]
		if !ok || raw == nil {
			return nil
		}
		switch pools := raw.(type) {
		case []any:
			out := make([]string, 0, len(pools))
			for _, v := range pools {
				out = append(out, fmt.Sprint(v))
			}
			return out
		case []string:
			return pools
		}
	}
	return nil
}

func parseRawInstanceAddress(v any) (string, error) {
	switch rv := v.(type) {
	case common.RawInstanceAddress:
		return string(rv.Unpack), nil
	case map[string]any:
		m := rv
		if data, ok := rv["data"].(map[string]any); ok {
			m = data
		}
		if unpack, ok := m["unpack"].(string); ok && unpack != "" {
			return unpack, nil
		}
	}
	return "", fmt.Errorf("unexpected raw instance address type %T", v)
}

func normalizeNumericText(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return v
	}
	if strings.Contains(v, ".") {
		parts := strings.SplitN(v, ".", 2)
		frac := strings.TrimRight(parts[1], "0")
		if frac == "" {
			return parts[0]
		}
	}
	if strings.ContainsAny(v, "eE") {
		if f, _, err := big.ParseFloat(v, 10, 256, big.ToZero); err == nil {
			if i, _ := f.Int(nil); i != nil {
				return i.String()
			}
		}
	}
	return strings.TrimSuffix(v, ".")
}

func findChainPoolConfigBySelector(chainPoolConfigs map[string]any, sourceSelectorKey string) (any, bool) {
	sourceSelectorNorm := normalizeNumericText(sourceSelectorKey)
	for rawKey, cfg := range chainPoolConfigs {
		if normalizeNumericText(rawKey) == sourceSelectorNorm {
			return cfg, true
		}
	}
	return nil, false
}

func chainPoolConfigFromAny(v any) (lockreleasetokenpool.ChainPoolConfig, bool) {
	switch cfg := v.(type) {
	case lockreleasetokenpool.ChainPoolConfig:
		return cfg, true
	case map[string]any:
		m := cfg
		if data, ok := cfg["data"].(map[string]any); ok {
			m = data
		}
		out := lockreleasetokenpool.ChainPoolConfig{
			InboundCCVs:  []common.RawInstanceAddress{},
			OutboundCCVs: []common.RawInstanceAddress{},
			RemotePools:  []types.TEXT{},
		}
		parseCCVs := func(raw any) []common.RawInstanceAddress {
			outCCV := make([]common.RawInstanceAddress, 0)
			rawList, ok := raw.([]any)
			if !ok {
				return outCCV
			}
			for _, item := range rawList {
				unpack, err := parseRawInstanceAddress(item)
				if err != nil || unpack == "" {
					continue
				}
				outCCV = append(outCCV, common.RawInstanceAddress{Unpack: types.TEXT(unpack)})
			}
			return outCCV
		}
		out.InboundCCVs = parseCCVs(m["inboundCCVs"])
		out.OutboundCCVs = parseCCVs(m["outboundCCVs"])
		if remoteRaw, ok := m["remotePools"].([]any); ok {
			for _, rp := range remoteRaw {
				out.RemotePools = append(out.RemotePools, types.TEXT(fmt.Sprint(rp)))
			}
		}
		return out, true
	default:
		return lockreleasetokenpool.ChainPoolConfig{}, false
	}
}

func ensureManualExecuteSourcePoolAllowed(
	ctx context.Context,
	participant canton.Participant,
	pool *lockreleasetokenpool.LockReleaseTokenPool,
	poolCID string,
	sourceSelectorKey string,
	sourcePoolHex string,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (*lockreleasetokenpool.LockReleaseTokenPool, string, error) {
	updatedChainPoolConfigs := types.GENMAP{}
	for k, v := range pool.ChainPoolConfigs {
		updatedChainPoolConfigs[k] = v
	}
	existingCfgAny, _ := findChainPoolConfigBySelector(pool.ChainPoolConfigs, sourceSelectorKey)
	existingCfg, ok := chainPoolConfigFromAny(existingCfgAny)
	if !ok {
		existingCfg = lockreleasetokenpool.ChainPoolConfig{
			InboundCCVs:  []common.RawInstanceAddress{},
			OutboundCCVs: []common.RawInstanceAddress{},
			RemotePools:  []types.TEXT{},
		}
	}
	existingCfg.RemotePools = []types.TEXT{types.TEXT(strings.TrimPrefix(strings.ToLower(sourcePoolHex), "0x"))}
	updatedChainPoolConfigs[sourceSelectorKey] = existingCfg

	updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateChainPoolConfigs{
		NewChainPoolConfigs: updatedChainPoolConfigs,
	}
	_, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  "#ccip-lockreleasetokenpool",
						ModuleName: "CCIP.LockReleaseTokenPool",
						EntityName: "LockReleaseTokenPool",
					},
					ContractId:     poolCID,
					Choice:         "LockReleaseTokenPool_UpdateChainPoolConfigs",
					ChoiceArgument: ledger.MapToValue(updateArgs.ToMap()),
				}},
			}},
			ActAs: []string{string(pool.PoolOwner)},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("patch lock/release pool source remote pool for manual execute: %w", err)
	}

	raw := contracts.MustNewInstanceID(string(pool.InstanceId)).RawInstanceAddress(pool.PoolOwner)
	updatedCID, err := resolveActiveContractIDByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), raw.InstanceAddress())
	if err != nil {
		return nil, "", fmt.Errorf("resolve updated lock/release pool contract id: %w", err)
	}
	updatedPool := *pool
	updatedPool.ChainPoolConfigs = updatedChainPoolConfigs
	return &updatedPool, updatedCID, nil
}

func resolveRateLimiterForManualExecute(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	sourceSelectorNumericKey string,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	if inboundRateLimiterRaw, ok := selectedPool.InboundRateLimiters[sourceSelectorKey]; ok {
		cid, disclosure, err := resolveRateLimiterFromRawAddress(ctx, participant, inboundRateLimiterRaw, resolveActiveContractIDByAddress)
		if err == nil {
			return cid, disclosure, nil
		}
	}
	if inboundRateLimiterRaw, ok := selectedPool.InboundRateLimiters[sourceSelectorNumericKey]; ok {
		cid, disclosure, err := resolveRateLimiterFromRawAddress(ctx, participant, inboundRateLimiterRaw, resolveActiveContractIDByAddress)
		if err == nil {
			return cid, disclosure, nil
		}
	}

	selectorNorm := normalizeNumericText(sourceSelectorKey)
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", nil, fmt.Errorf("get ledger end for fallback rate limiter lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &apiv2.TemplateFilter{
									TemplateId: &apiv2.Identifier{
										PackageId:  "#ccip-common",
										ModuleName: "CCIP.RateLimiter",
										EntityName: "RateLimiter",
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return "", nil, fmt.Errorf("query active rate limiters for fallback lookup: %w", err)
	}
	defer stream.CloseSend()

	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return "", nil, fmt.Errorf("receive active rate limiters for fallback lookup: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		parsed, parseErr := bindings.UnmarshalCreatedEvent[common.RateLimiter](entry.ActiveContract.GetCreatedEvent())
		if parseErr != nil {
			continue
		}
		if parsed.Direction != common.RateLimitDirectionRateLimitDirection_Inbound {
			continue
		}
		if string(parsed.PoolOwner) != string(selectedPool.PoolOwner) || string(parsed.PoolInstanceId) != string(selectedPool.InstanceId) {
			continue
		}
		if normalizeNumericText(string(parsed.RemoteChainSelector)) != selectorNorm {
			continue
		}
		cid := entry.ActiveContract.GetCreatedEvent().GetContractId()
		disclosure := &apiv2.DisclosedContract{
			TemplateId:       entry.ActiveContract.GetCreatedEvent().GetTemplateId(),
			ContractId:       cid,
			CreatedEventBlob: entry.ActiveContract.GetCreatedEvent().GetCreatedEventBlob(),
			SynchronizerId:   entry.ActiveContract.GetSynchronizerId(),
		}
		return cid, disclosure, nil
	}

	return "", nil, fmt.Errorf("missing inbound rate limiter for source selector %s", sourceSelectorKey)
}

func resolveRateLimiterFromRawAddress(
	ctx context.Context,
	participant canton.Participant,
	inboundRateLimiterRaw any,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	rawRateLimiter, err := parseRawInstanceAddress(inboundRateLimiterRaw)
	if err != nil {
		return "", nil, err
	}
	rateLimiterRawAddr, err := contracts.RawInstanceAddressFromString(rawRateLimiter)
	if err != nil {
		return "", nil, fmt.Errorf("parse rate limiter raw instance address: %w", err)
	}
	rateLimiterCID, err := resolveActiveContractIDByAddress(common.RateLimiter{}.GetTemplateID(), rateLimiterRawAddr.InstanceAddress())
	if err != nil {
		return "", nil, fmt.Errorf("resolve rate limiter contract id: %w", err)
	}
	rateLimiterDisclosure, err := getDisclosedContractByID(ctx, participant, rateLimiterCID)
	if err != nil {
		return "", nil, fmt.Errorf("resolve rate limiter disclosure: %w", err)
	}
	return rateLimiterCID, rateLimiterDisclosure, nil
}

func listHoldingContracts(ctx context.Context, participant canton.Participant) ([]*apiv2.ActiveContract, error) {
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end for holdings: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{
								InterfaceFilter: &apiv2.InterfaceFilter{
									InterfaceId: &apiv2.Identifier{
										PackageId:  "#splice-api-token-holding-v1",
										ModuleName: "Splice.Api.Token.HoldingV1",
										EntityName: "Holding",
									},
									IncludeInterfaceView:    true,
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query holding contracts: %w", err)
	}
	defer stream.CloseSend()

	var out []*apiv2.ActiveContract
	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, fmt.Errorf("receive holding contracts: %w", recvErr)
		}
		if entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			out = append(out, entry.ActiveContract)
		}
	}

	return out, nil
}

func getDisclosedContractByID(ctx context.Context, participant canton.Participant, contractID string) (*apiv2.DisclosedContract, error) {
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end for disclosure lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
								WildcardFilter: &apiv2.WildcardFilter{
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query contracts for disclosure lookup: %w", err)
	}
	defer stream.CloseSend()

	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, fmt.Errorf("receive disclosure lookup contracts: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}
		created := entry.ActiveContract.GetCreatedEvent()
		if created.GetContractId() != contractID {
			continue
		}

		return &apiv2.DisclosedContract{
			TemplateId:       created.GetTemplateId(),
			ContractId:       created.GetContractId(),
			CreatedEventBlob: created.GetCreatedEventBlob(),
			SynchronizerId:   entry.ActiveContract.GetSynchronizerId(),
		}, nil
	}

	return nil, fmt.Errorf("contract id %s not found for disclosure lookup", contractID)
}

// This is copied from chainlink-canton, replace with EDS client once available.
func ChoiceContextFromData(choiceContextData map[string]any) (*apiv2.Value, error) {
	values, ok := choiceContextData["values"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no values found in choice context")
	}

	// ref: https://docs.digitalasset.com/build/3.5/reference/json-api/lf-value-specification.html
	// AnyValue is a variant
	fields := make([]*apiv2.TextMap_Entry, 0, len(values))
	for k, v := range values {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value *apiv2.Value
		switch tag {
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Text{Text: valueString}}
		case "AV_Int":
			// JSON numbers come as float64
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(valueFloat)}}
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: valueString}}
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return nil, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: valueBool}}
		case "AV_Date":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Date{Date: int32(t.Unix() / 86400)}} //nolint:gosec
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Time value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: t.UnixMicro()}}
		case "AV_RelTime":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_RelTime value is not a number: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(valueFloat)}}},
			}}}}
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: valueString}}
		default:
			// Add lists and maps
			return nil, fmt.Errorf("unimplemented tag: %v", tag)
		}

		fields = append(fields, &apiv2.TextMap_Entry{
			Key: k,
			Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
				Constructor: tag,
				Value:       value,
			}}},
		})
	}

	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "values",
			Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: fields}}},
		},
	}}}}, nil
}
