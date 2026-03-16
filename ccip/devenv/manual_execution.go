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
	"strconv"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
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
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"

	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

func uniqueParties(parties ...string) []string {
	seen := make(map[string]struct{}, len(parties))
	out := make([]string, 0, len(parties))
	for _, party := range parties {
		if party == "" {
			continue
		}
		if _, ok := seen[party]; ok {
			continue
		}
		seen[party] = struct{}{}
		out = append(out, party)
	}

	return out
}

func isAlreadyExistsError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already exists")
}

func newTransferInstructionClient(participant canton.Participant) (*transferInstructionV1.ClientWithResponses, error) {
	requestEditor := func(ctx context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("retrieve participant token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		return nil
	}

	client, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("create transfer instruction client: %w", err)
	}

	return client, nil
}

func emptyMetadataValue() *apiv2.Value {
	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
	}}}}
}

func perPartyRouterQualifier(partyID string) string {
	return fmt.Sprintf("party:%s", partyID)
}

// DeployPerPartyRouter uses the PerPartyRouterFactory to create a new PerPartyRouter instance for the given party.
// It returns the address of the newly created PerPartyRouter instance. If a router already exists for the party, it returns the existing router's address.
func (c *Chain) DeployPerPartyRouter(ctx context.Context, partyId string) (contracts.InstanceAddress, error) {
	participantIdx := c.participantIndexForParty(partyId)
	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: participantIdx,
	}
	var err error
	// Create PerPartyRouter using the topology-deployed default factory.
	cantonPerPartyRouterFactoryRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		c.chainDetails.ChainSelector,
		datastore.ContractType(per_party_router_factory.ContractType),
		per_party_router_factory.Version,
		"",
	))
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to get per-party router factory address ref: %w", err)
	}
	c.logger.Debug().Str("CantonPerPartyRouterFactory", cantonPerPartyRouterFactoryRef.Address).Msg("Resolved per-party router factory address")
	cantonPerPartyRouterFactory := contracts.HexToInstanceAddress(cantonPerPartyRouterFactoryRef.Address)

	// Fixed instance ID keeps address derivation deterministic; InstanceAddress is still unique by party.
	routerInstanceID := contracts.InstanceID("test-router")
	// Ignore only idempotent create errors, since the router might already exist for this party.
	_, err = operations.ExecuteOperation(c.e.OperationsBundle, per_party_router_factory.CreateRouter, deps, contract.ChoiceInput[perpartyrouter.CreateRouter]{
		ChainSelector:   c.chainDetails.ChainSelector,
		InstanceAddress: cantonPerPartyRouterFactory,
		ActAs:           []string{partyId},
		Args: perpartyrouter.CreateRouter{
			PartyOwner: types.PARTY(partyId),
			InstanceId: types.TEXT(routerInstanceID.String()),
		},
	})
	if err != nil && !isAlreadyExistsError(err) {
		return contracts.InstanceAddress{}, fmt.Errorf("failed to create per-party router: %w", err)
	}
	selectedAddress := routerInstanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()
	updatedDataStore := datastore.NewMemoryDataStore()
	if err = updatedDataStore.Merge(c.e.DataStore); err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("clone datastore for per-party router persistence: %w", err)
	}
	if err = updatedDataStore.AddressRefStore.Upsert(datastore.AddressRef{
		Address:       selectedAddress.Hex(),
		Type:          datastore.ContractType("PerPartyRouter"),
		Version:       per_party_router_factory.Version,
		Qualifier:     perPartyRouterQualifier(partyId),
		ChainSelector: c.chainDetails.ChainSelector,
	}); err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("persist per-party router address ref: %w", err)
	}
	c.e.DataStore = updatedDataStore.Seal()

	return selectedAddress, nil
}

func (c *Chain) DeployCCIPReceiver(ctx context.Context, partyId string) (contracts.InstanceAddress, error) {
	participantIdx := c.participantIndexForParty(partyId)

	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: participantIdx,
	}
	qualifier := perPartyRouterQualifier(partyId)

	// Deploy receiver contract
	out, err := operations.ExecuteOperation(c.e.OperationsBundle, receiver.Deploy, deps, contract.DeployInput[ccipreceiver.CCIPReceiver]{
		ChainSelector: c.chainDetails.ChainSelector,
		Qualifier:     &qualifier,
		ActAs:         []string{partyId},
		Template: ccipreceiver.CCIPReceiver{
			Owner:        types.PARTY(partyId),
			RequiredCCVs: nil,
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
	receiverParty, err := c.resolvePartyFromHashedAddress(ctx, message.Receiver)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("resolve executing party from receiver %s: %w", hex.EncodeToString(message.Receiver), err)
	}
	// Execute on the receiver party so receiver hash/party checks line up with the message payload.
	executingParty := receiverParty
	participantIdx := c.participantIndexForParty(executingParty)
	participant := c.chain.Participants[participantIdx]

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

	activeRouter, err := c.findLatestActiveContractByInstanceAddress(
		ctx,
		participant,
		perpartyrouter.PerPartyRouter{}.GetTemplateID(),
		routerAddress,
	)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get router contract ID: %w", err)
	}
	routerCid := activeRouter.GetCreatedEvent().GetContractId()
	c.logger.Debug().Str("InstanceAddress", routerAddress.String()).Str("ContractId", routerCid).Msg("Resolved PerPartyRouter contract")

	activeReceiver, err := c.findLatestActiveContractByInstanceAddress(
		ctx,
		participant,
		ccipreceiver.CCIPReceiver{}.GetTemplateID(),
		receiverAddress,
	)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get receiver contract ID: %w", err)
	}
	receiverCid := activeReceiver.GetCreatedEvent().GetContractId()
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
		var poolTokenAdminRegistryCID string
		tokenTransferValue, tokenDisclosures, poolTokenAdminRegistryCID, err = c.buildManualExecuteTokenTransferInput(
			ctx,
			participant,
			receiverParty,
			&message,
			func(templateID string, address contracts.InstanceAddress) (string, error) {
				active, findErr := c.findLatestActiveContractByInstanceAddress(ctx, participant, templateID, address)
				if findErr != nil {
					return "", findErr
				}

				return active.GetCreatedEvent().GetContractId(), nil
			},
		)
		if err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("build token transfer execute input: %w", err)
		}
		if poolTokenAdminRegistryCID != "" {
			if err := overrideChoiceContextContractID(choiceContext, "token-admin-registry", poolTokenAdminRegistryCID); err != nil {
				return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("override token-admin-registry in execute context: %w", err)
			}
		}
		executeDisclosures = append(executeDisclosures, tokenDisclosures...)
	}
	executeDisclosures = dedupeDisclosedContractsByID(executeDisclosures)
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
	if message.TokenTransfer != nil {
		if err := c.acceptPendingTransferInstruction(ctx, participant, executingParty, res.GetTransaction().GetEvents()); err != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("accept pending token transfer instruction: %w", err)
		}
	}

	// Get Update
	updateRes, err := participant.LedgerServices.Update.GetUpdateById(ctx, &apiv2.GetUpdateByIdRequest{
		UpdateId: res.GetTransaction().GetUpdateId(),
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						executingParty: {
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

func (c *Chain) acceptPendingTransferInstruction(
	ctx context.Context,
	participant canton.Participant,
	executingParty string,
	events []*apiv2.Event,
) error {
	var pendingTransferInstructionCID string
	for _, event := range events {
		created := event.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		if created.GetTemplateId().GetEntityName() == "AmuletTransferInstruction" {
			pendingTransferInstructionCID = created.GetContractId()
		}
	}
	if pendingTransferInstructionCID == "" {
		return nil
	}

	transferClient, err := newTransferInstructionClient(participant)
	if err != nil {
		return err
	}

	acceptCtxResp, err := transferClient.GetTransferInstructionAcceptContextWithResponse(
		ctx,
		pendingTransferInstructionCID,
		transferInstructionV1.GetChoiceContextRequest{},
	)
	if err != nil {
		return fmt.Errorf("get transfer instruction accept context: %w", err)
	}
	if acceptCtxResp.StatusCode() != http.StatusOK || acceptCtxResp.JSON200 == nil {
		return fmt.Errorf("transfer instruction accept context response status %d: %s", acceptCtxResp.StatusCode(), string(acceptCtxResp.Body))
	}
	acceptContext, err := ChoiceContextFromData(acceptCtxResp.JSON200.ChoiceContextData)
	if err != nil {
		return fmt.Errorf("convert transfer instruction accept context: %w", err)
	}

	acceptDisclosures := make([]*apiv2.DisclosedContract, 0, len(acceptCtxResp.JSON200.DisclosedContracts))
	for _, contract := range acceptCtxResp.JSON200.DisclosedContracts {
		id, parseErr := TemplateIdFromString(contract.TemplateId)
		if parseErr != nil {
			return fmt.Errorf("parse transfer instruction accept disclosure template id: %w", parseErr)
		}
		createdEventBlob, decodeErr := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if decodeErr != nil {
			return fmt.Errorf("decode transfer instruction accept disclosure created event blob: %w", decodeErr)
		}
		acceptDisclosures = append(acceptDisclosures, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}

	emptyMetadata := emptyMetadataValue()
	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  "#splice-api-token-transfer-instruction-v1",
						ModuleName: "Splice.Api.Token.TransferInstructionV1",
						EntityName: "TransferInstruction",
					},
					ContractId: pendingTransferInstructionCID,
					Choice:     "TransferInstruction_Accept",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "context", Value: acceptContext},
							{Label: "meta", Value: emptyMetadata},
						}}}}},
					}}}},
				}},
			}},
			ActAs:              uniqueParties(executingParty, participant.PartyID),
			DisclosedContracts: acceptDisclosures,
		},
	})
	if err != nil {
		return fmt.Errorf("submit transfer instruction accept: %w", err)
	}

	return nil
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

// BuildManualExecuteTokenTransferInput builds the input for the manual execute token transfer command.
// Assumes deploy/configure already produced a correct fresh environment and follows the direct path:
// - find a lock/release pool with the source selector config and exact source pool match;
// - if dest token address is provided, require exact remoteTokenAddress match;
// - resolve configured inbound rate limiter and gather required disclosures;
// - build transfer input using sender holdings plus transfer-factory context.
func (c *Chain) buildManualExecuteTokenTransferInput(
	ctx context.Context,
	participant canton.Participant,
	tokenReceiverParty string,
	message *protocol.Message,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (*apiv2.Value, []*apiv2.DisclosedContract, string, error) {
	if message == nil || message.TokenTransfer == nil {
		return nil, nil, "", fmt.Errorf("token transfer message is required")
	}
	const tokenPoolQualifier = "TEST (LockReleaseTokenPool 1.7.0 [default] to BurnMintTokenPool 1.7.0 [default])"
	sourceSelectorKey := fmt.Sprintf("%d", message.SourceChainSelector)
	sourcePoolHex := strings.ToLower(hex.EncodeToString(message.TokenTransfer.SourcePoolAddress))
	sourcePoolHex = strings.TrimPrefix(sourcePoolHex, "0x")

	var tokenPoolRef datastore.AddressRef
	foundPoolRef := false
	candidates := c.e.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(c.chainDetails.ChainSelector),
		datastore.AddressRefByQualifier(tokenPoolQualifier),
	)
	for _, candidate := range candidates {
		if candidate.Type == datastore.ContractType("LockReleaseTokenPool") ||
			candidate.Type == datastore.ContractType(lock_release_token_pool.ContractType) {
			tokenPoolRef = candidate
			foundPoolRef = true
			break
		}
	}
	if !foundPoolRef {
		return nil, nil, "", fmt.Errorf("resolve source lock/release pool from datastore: no matching address ref for qualifier %q", tokenPoolQualifier)
	}
	tokenPoolAddress := contracts.HexToInstanceAddress(tokenPoolRef.Address)
	selectedPoolContractID, err := resolveActiveContractIDByAddress(
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		tokenPoolAddress,
	)
	if err != nil {
		return nil, nil, "", fmt.Errorf("resolve lock/release pool contract ID by datastore address: %w", err)
	}
	activePool, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		tokenPoolAddress,
	)
	if err != nil {
		return nil, nil, "", fmt.Errorf("find lock/release pool active contract by datastore address: %w", err)
	}
	selectedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return nil, nil, "", fmt.Errorf("parse selected lock/release pool: %w", err)
	}
	poolTokenAdminRegistryCID := ""
	var poolTARDisclosure *apiv2.DisclosedContract
	// Resolve TAR from the selected pool deps when present to avoid mismatches
	// with globally looked-up/default TAR contracts.
	if poolTARRaw, rawErr := parseRawInstanceAddress(selectedPool.Deps.TokenAdminRegistry); rawErr == nil && strings.TrimSpace(poolTARRaw) != "" {
		poolTARRawAddr, err := contracts.RawInstanceAddressFromString(poolTARRaw)
		if err != nil {
			return nil, nil, "", fmt.Errorf("parse selected pool token admin registry address %q: %w", poolTARRaw, err)
		}
		poolTARActive, err := contract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			participant.PartyID,
			tokenadminregistry.TokenAdminRegistry{}.GetTemplateID(),
			poolTARRawAddr.InstanceAddress(),
		)
		if err != nil {
			return nil, nil, "", fmt.Errorf("find selected pool token admin registry contract by address: %w", err)
		}
		poolTokenAdminRegistryCID = poolTARActive.GetCreatedEvent().GetContractId()
		poolTARDisclosure = &apiv2.DisclosedContract{
			TemplateId:       poolTARActive.GetCreatedEvent().GetTemplateId(),
			ContractId:       poolTARActive.GetCreatedEvent().GetContractId(),
			CreatedEventBlob: poolTARActive.GetCreatedEvent().GetCreatedEventBlob(),
			SynchronizerId:   poolTARActive.GetSynchronizerId(),
		}
	}

	chainPoolCfgAny, ok := findChainPoolConfigBySelector(selectedPool.RemoteChainConfigs, sourceSelectorKey)
	if !ok {
		return nil, nil, "", fmt.Errorf("selected lock/release pool has no remote config for source selector %s", sourceSelectorKey)
	}
	cfg, cfgOK := chainPoolConfigFromAny(chainPoolCfgAny)
	if !cfgOK {
		return nil, nil, "", fmt.Errorf("selected lock/release pool remote config is invalid for source selector %s", sourceSelectorKey)
	}
	remotePoolMatch := false
	for _, remotePool := range cfg.RemotePools {
		remotePoolHex := strings.ToLower(strings.TrimPrefix(string(remotePool), "0x"))
		if remotePoolHex == sourcePoolHex {
			remotePoolMatch = true
			break
		}
	}
	if !remotePoolMatch {
		return nil, nil, "", fmt.Errorf(
			"selected lock/release pool is not configured for source pool %s on source selector %s",
			sourcePoolHex,
			sourceSelectorKey,
		)
	}

	rateLimiterCID, rateLimiterDisclosure, err := resolveRateLimiterForManualExecute(
		ctx,
		participant,
		selectedPool,
		sourceSelectorKey,
		resolveActiveContractIDByAddress,
	)
	if err != nil {
		return nil, nil, "", fmt.Errorf("resolve rate limiter disclosure: %w", err)
	}

	instrumentAdmin := string(selectedPool.InstrumentId.Admin)
	instrumentID := string(selectedPool.InstrumentId.Id)
	expectedTransferAdmin := instrumentAdmin
	transferSenderParty := string(selectedPool.PoolOwner)

	holdings, err := listHoldingContracts(ctx, participant)
	if err != nil {
		return nil, nil, "", fmt.Errorf("list pool holdings: %w", err)
	}
	poolHoldingCIDs, poolHoldingDisclosures := selectUnlockedHoldingCIDs(
		holdings,
		transferSenderParty,
		instrumentAdmin,
		instrumentID,
	)
	poolHoldings := make([]string, 0, len(poolHoldingCIDs))
	for _, cid := range poolHoldingCIDs {
		poolHoldings = append(poolHoldings, string(cid))
	}
	if len(poolHoldings) == 0 {
		return nil, nil, "", fmt.Errorf(
			"no unlocked pool holdings found for transfer sender %s and instrument %s/%s",
			transferSenderParty,
			instrumentAdmin,
			instrumentID,
		)
	}

	transferClient, err := newTransferInstructionClient(participant)
	if err != nil {
		return nil, nil, "", err
	}
	transferFactoryResp, err := transferClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": expectedTransferAdmin,
			"transfer": map[string]any{
				"sender":   transferSenderParty,
				"receiver": tokenReceiverParty,
				"amount":   message.TokenTransfer.Amount.String(),
				"instrumentId": map[string]any{
					"admin": instrumentAdmin,
					"id":    instrumentID,
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"inputHoldingCids": poolHoldings,
				"meta":             map[string]any{"values": map[string]any{}},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{"values": map[string]any{}},
				"meta":    map[string]any{"values": map[string]any{}},
			},
		},
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("get transfer factory: %w", err)
	}
	if transferFactoryResp.StatusCode() != http.StatusOK {
		return nil, nil, "", fmt.Errorf("transfer factory response status %d: %s", transferFactoryResp.StatusCode(), string(transferFactoryResp.Body))
	}
	transferFactoryCID := transferFactoryResp.JSON200.FactoryId
	transferFactoryCtx, err := ChoiceContextFromData(transferFactoryResp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return nil, nil, "", fmt.Errorf("convert transfer factory context: %w", err)
	}
	transferFactoryDisclosures := make([]*apiv2.DisclosedContract, 0, len(transferFactoryResp.JSON200.ChoiceContext.DisclosedContracts))
	for _, d := range transferFactoryResp.JSON200.ChoiceContext.DisclosedContracts {
		id, idErr := TemplateIdFromString(d.TemplateId)
		if idErr != nil {
			return nil, nil, "", fmt.Errorf("parse transfer factory disclosure template id: %w", idErr)
		}
		createdEventBlob, decodeErr := base64.StdEncoding.DecodeString(d.CreatedEventBlob)
		if decodeErr != nil {
			return nil, nil, "", fmt.Errorf("decode transfer factory disclosure created event blob: %w", decodeErr)
		}
		transferFactoryDisclosures = append(transferFactoryDisclosures, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       d.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   d.SynchronizerId,
		})
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
		return nil, nil, "", fmt.Errorf("build pool extra context: %w", err)
	}
	emptyMetadata := emptyMetadataValue()
	tokenPoolHoldingValues := make([]*apiv2.Value, 0, len(poolHoldings))
	for _, cid := range poolHoldings {
		tokenPoolHoldingValues = append(tokenPoolHoldingValues, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}})
	}

	tokenTransferValue := &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "tokenPoolCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: selectedPoolContractID}}},
		{Label: "tokenReceiverParty", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: tokenReceiverParty}}},
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

	poolDisclosure := &apiv2.DisclosedContract{
		TemplateId:       activePool.GetCreatedEvent().GetTemplateId(),
		ContractId:       activePool.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: activePool.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   activePool.GetSynchronizerId(),
	}
	disclosures := make([]*apiv2.DisclosedContract, 0, 3+len(transferFactoryDisclosures)+len(poolHoldingDisclosures))
	disclosures = append(disclosures, rateLimiterDisclosure, poolDisclosure)
	if poolTARDisclosure != nil {
		disclosures = append(disclosures, poolTARDisclosure)
	}
	disclosures = append(disclosures, transferFactoryDisclosures...)
	disclosures = append(disclosures, poolHoldingDisclosures...)
	c.logger.Debug().
		Str("SelectedTokenPoolCID", selectedPoolContractID).
		Str("MessageSourcePoolHex", sourcePoolHex).
		Str("SelectedPoolInstrumentAdmin", string(selectedPool.InstrumentId.Admin)).
		Str("ExpectedTransferAdmin", expectedTransferAdmin).
		Msg("Prepared manual execute token transfer input")

	return tokenTransferValue, disclosures, poolTokenAdminRegistryCID, nil
}

func parseRawInstanceAddress(v any) (string, error) {
	switch rv := v.(type) {
	case common.RawInstanceAddress:
		return string(rv.Unpack), nil
	case string:
		if rv == "" {
			return "", fmt.Errorf("empty raw instance address string")
		}

		return rv, nil
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

func dedupeDisclosedContractsByID(disclosures []*apiv2.DisclosedContract) []*apiv2.DisclosedContract {
	out := make([]*apiv2.DisclosedContract, 0, len(disclosures))
	seen := make(map[string]struct{}, len(disclosures))
	for _, dc := range disclosures {
		if dc == nil {
			continue
		}
		cid := strings.TrimSpace(dc.GetContractId())
		if cid == "" {
			continue
		}
		if _, ok := seen[cid]; ok {
			continue
		}
		seen[cid] = struct{}{}
		out = append(out, dc)
	}

	return out
}

func overrideChoiceContextContractID(choiceContext *apiv2.Value, key string, contractID string) error {
	if choiceContext == nil || choiceContext.GetRecord() == nil {
		return fmt.Errorf("choice context record is nil")
	}
	fields := choiceContext.GetRecord().GetFields()
	for _, field := range fields {
		if field == nil || field.GetLabel() != "values" || field.GetValue() == nil || field.GetValue().GetTextMap() == nil {
			continue
		}
		entries := field.GetValue().GetTextMap().GetEntries()
		for _, entry := range entries {
			if entry == nil || entry.GetKey() != key {
				continue
			}
			entry.Value = &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
				Constructor: "AV_ContractId",
				Value:       &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: contractID}},
			}}}

			return nil
		}
	}

	return fmt.Errorf("choice context key %q not found", key)
}

func normalizeSelectorKey(v string) (string, bool) {
	v = strings.TrimSpace(v)
	if v == "" || strings.ContainsAny(v, "eE") {
		return "", false
	}

	parts := strings.SplitN(v, ".", 2)
	intPart := parts[0]
	if intPart == "" {
		return "", false
	}
	if len(parts) == 2 && strings.Trim(parts[1], "0") != "" {
		return "", false
	}

	n, err := strconv.ParseUint(intPart, 10, 64)
	if err != nil {
		return "", false
	}

	return strconv.FormatUint(n, 10), true
}

func findChainPoolConfigBySelector(chainPoolConfigs map[string]any, sourceSelectorKey string) (any, bool) {
	sourceSelectorNorm, ok := normalizeSelectorKey(sourceSelectorKey)
	if !ok {
		return nil, false
	}
	for rawKey, cfg := range chainPoolConfigs {
		rawKeyNorm, rawKeyOK := normalizeSelectorKey(rawKey)
		if rawKeyOK && rawKeyNorm == sourceSelectorNorm {
			return cfg, true
		}
	}

	return nil, false
}

func chainPoolConfigFromAny(v any) (lockreleasetokenpool.RemoteChainConfig, bool) {
	switch cfg := v.(type) {
	case lockreleasetokenpool.RemoteChainConfig:
		return cfg, true
	case map[string]any:
		m := cfg
		if data, ok := cfg["data"].(map[string]any); ok {
			m = data
		}
		out := lockreleasetokenpool.RemoteChainConfig{
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
		if raw, ok := m["inboundRateLimiter"]; ok {
			if unpack, err := parseRawInstanceAddress(raw); err == nil && unpack != "" {
				out.InboundRateLimiter = common.RawInstanceAddress{Unpack: types.TEXT(unpack)}
			}
		}
		if raw, ok := m["outboundRateLimiter"]; ok {
			if unpack, err := parseRawInstanceAddress(raw); err == nil && unpack != "" {
				out.OutboundRateLimiter = common.RawInstanceAddress{Unpack: types.TEXT(unpack)}
			}
		}
		if remoteTokenAddress, ok := m["remoteTokenAddress"].(string); ok {
			out.RemoteTokenAddress = types.TEXT(remoteTokenAddress)
		}

		return out, true
	default:
		return lockreleasetokenpool.RemoteChainConfig{}, false
	}
}

func resolveRateLimiterForManualExecute(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	cfgAny, ok := findChainPoolConfigBySelector(selectedPool.RemoteChainConfigs, sourceSelectorKey)
	if !ok {
		return "", nil, fmt.Errorf("missing configured inbound rate limiter entry for source selector %s", sourceSelectorKey)
	}
	cfg, ok := chainPoolConfigFromAny(cfgAny)
	if !ok {
		return "", nil, fmt.Errorf("invalid configured inbound rate limiter entry for source selector %s", sourceSelectorKey)
	}
	if strings.TrimSpace(string(cfg.InboundRateLimiter.Unpack)) == "" {
		return "", nil, fmt.Errorf("missing configured inbound rate limiter entry for source selector %s", sourceSelectorKey)
	}
	rawRateLimiter, err := parseRawInstanceAddress(cfg.InboundRateLimiter)
	if err != nil {
		return "", nil, err
	}
	rateLimiterRawAddr, err := contracts.RawInstanceAddressFromString(rawRateLimiter)
	if err != nil {
		return "", nil, fmt.Errorf("parse rate limiter raw instance address: %w", err)
	}
	activeRateLimiter, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		common.RateLimiter{}.GetTemplateID(),
		rateLimiterRawAddr.InstanceAddress(),
	)
	if err != nil {
		return "", nil, fmt.Errorf("find inbound rate limiter active contract by configured raw address: %w", err)
	}
	cid, err := resolveActiveContractIDByAddress(common.RateLimiter{}.GetTemplateID(), rateLimiterRawAddr.InstanceAddress())
	if err != nil {
		return "", nil, fmt.Errorf("resolve inbound rate limiter contract ID by configured raw address: %w", err)
	}

	return cid, &apiv2.DisclosedContract{
		TemplateId:       activeRateLimiter.GetCreatedEvent().GetTemplateId(),
		ContractId:       activeRateLimiter.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: activeRateLimiter.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   activeRateLimiter.GetSynchronizerId(),
	}, nil
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
