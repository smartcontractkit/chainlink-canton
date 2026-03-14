package devenv

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
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
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
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

var executeRequiredPackages = []contracts.Package{
	contracts.CCIPCommon,
	contracts.CCIPCommitteeVerifier,
	contracts.CCIPPerPartyRouter,
	contracts.CCIPReceiver,
	contracts.CCIPLockReleaseTokenPool,
	contracts.CCIPTokenAdminRegistry,
	contracts.CCIPOffRamp,
	contracts.CCIPRMN,
}

const (
	perPartyRouterEntityName = "PerPartyRouter"
	createArgFieldPartyOwner = "partyOwner"
	createArgFieldInstanceID = "instanceId"
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

func isRouterCreateConflictError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()

	return strings.Contains(msg, "instanceId already in use") || strings.Contains(msg, "router already exists for this party")
}

func perPartyRouterInstanceID(partyID string) contracts.InstanceID {
	partyHash := crypto.Keccak256([]byte(partyID))
	return contracts.InstanceID(fmt.Sprintf("test-router-%x", partyHash[:6]))
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

// DeployPerPartyRouter uses the PerPartyRouterFactory to create a new PerPartyRouter instance for the given party.
// It returns the address of the newly created PerPartyRouter instance. If a router already exists for the party, it returns the existing router's address.
func (c *Chain) DeployPerPartyRouter(ctx context.Context, participantIdx int, partyId string) (contracts.InstanceAddress, error) {
	routerAddress, _, err := c.deployPerPartyRouterWithContractID(ctx, participantIdx, partyId)
	return routerAddress, err
}

func (c *Chain) deployPerPartyRouterWithContractID(ctx context.Context, participantIdx int, partyId string) (contracts.InstanceAddress, string, error) {
	participant := c.chain.Participants[participantIdx]
	if err := ensureParticipantDarVetted(ctx, participant, contracts.CCIPPerPartyRouter); err != nil {
		return contracts.InstanceAddress{}, "", fmt.Errorf("ensure per-party router dar on participant %d: %w", participantIdx, err)
	}
	if existingCID, existingInstanceID, err := c.findAnyPerPartyRouterForOwner(ctx, participant, partyId); err == nil && existingCID != "" && existingInstanceID != "" {
		return contracts.InstanceID(existingInstanceID).RawInstanceAddress(types.PARTY(partyId)).InstanceAddress(), existingCID, nil
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
		return contracts.InstanceAddress{}, "", fmt.Errorf("failed to get canton per party router factory address ref: %w", err)
	}
	c.logger.Debug().Str("CantonPerPartyRouterFactory", cantonPerPartyRouterFactoryRef.Address).Msg("Resolved per-party router factory address")
	cantonPerPartyRouterFactory := contracts.HexToInstanceAddress(cantonPerPartyRouterFactoryRef.Address)

	// Use deterministic, party-scoped router IDs to avoid cross-party instance collisions.
	routerInstanceID := perPartyRouterInstanceID(partyId)
	routerAddress := routerInstanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()
	fallbackRouterCID, fallbackErr := c.createRouterWithDisclosedFactory(ctx, participantIdx, partyId, routerInstanceID, cantonPerPartyRouterFactory)
	if fallbackErr != nil {
		if _, lookupErr := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, partyId, perpartyrouter.PerPartyRouter{}.GetTemplateID(), routerAddress); lookupErr != nil {
			return contracts.InstanceAddress{}, "", fmt.Errorf("failed to create per-party router: %w", fallbackErr)
		}
	}

	routerCid, err := c.findPerPartyRouterContractID(ctx, participant, partyId, string(routerInstanceID))
	if err != nil {
		if fallbackRouterCID == "" {
			return contracts.InstanceAddress{}, "", fmt.Errorf("resolve per-party router contract id: %w", err)
		}
		routerCid = fallbackRouterCID
	}

	return routerAddress, routerCid, nil
}

func ensureParticipantDarVetted(ctx context.Context, participant canton.Participant, pkg contracts.Package) error {
	dar, err := contracts.GetDar(pkg, contracts.CurrentVersion)
	if err != nil {
		return fmt.Errorf("get %s dar: %w", pkg, err)
	}
	if _, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile:       dar,
		VettingChange: adminv2.UploadDarFileRequest_VETTING_CHANGE_VET_ALL_PACKAGES,
	}); err != nil && !isAlreadyExistsError(err) {
		return fmt.Errorf("upload %s dar: %w", pkg, err)
	}

	return nil
}

func ensureExecuteParticipantDarsVetted(ctx context.Context, participant canton.Participant) error {
	for _, pkg := range executeRequiredPackages {
		if err := ensureParticipantDarVetted(ctx, participant, pkg); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) createRouterWithDisclosedFactory(
	ctx context.Context,
	participantIdx int,
	partyID string,
	routerInstanceID contracts.InstanceID,
	factoryAddr contracts.InstanceAddress,
) (string, error) {
	if len(c.chain.Participants) == 0 {
		return "", fmt.Errorf("no canton participants configured")
	}
	targetParticipant := c.chain.Participants[participantIdx]
	factoryParticipant := c.chain.Participants[0]

	activeFactory, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		factoryParticipant.LedgerServices.State,
		factoryParticipant.PartyID,
		perpartyrouter.PerPartyRouterFactory{}.GetTemplateID(),
		factoryAddr,
	)
	if err != nil {
		return "", fmt.Errorf("resolve disclosed per-party router factory: %w", err)
	}
	created := activeFactory.GetCreatedEvent()
	if created == nil {
		return "", fmt.Errorf("per-party router factory created event is nil")
	}
	if created.GetContractId() == "" || created.GetTemplateId() == nil {
		return "", fmt.Errorf("per-party router factory created event missing contract/template id")
	}
	if len(created.GetCreatedEventBlob()) == 0 {
		return "", fmt.Errorf("per-party router factory disclosure blob is empty")
	}

	disclosedFactory := &apiv2.DisclosedContract{
		TemplateId:       created.GetTemplateId(),
		ContractId:       created.GetContractId(),
		CreatedEventBlob: created.GetCreatedEventBlob(),
	}

	res, err := targetParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: created.GetTemplateId(),
					ContractId: created.GetContractId(),
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.CreateRouter{
						PartyOwner: types.PARTY(partyID),
						InstanceId: types.TEXT(routerInstanceID.String()),
					}.ToMap()),
				}},
			}},
			ActAs:              []string{partyID},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedFactory},
		},
	})
	if err != nil && !isRouterCreateConflictError(err) {
		return "", fmt.Errorf("submit CreateRouter with disclosed factory: %w", err)
	}

	if res != nil && res.GetTransaction() != nil {
		for _, event := range res.GetTransaction().GetEvents() {
			if created, ok := event.GetEvent().(*apiv2.Event_Created); ok {
				createdEvent := created.Created
				if createdEvent == nil || createdEvent.GetTemplateId() == nil {
					continue
				}
				if createdEvent.GetTemplateId().GetEntityName() == perPartyRouterEntityName {
					return createdEvent.GetContractId(), nil
				}
			}
		}
	}

	return "", nil
}

func (c *Chain) DeployCCIPReceiver(ctx context.Context, participantIdx int, partyId string) (contracts.InstanceAddress, error) {
	participant := c.chain.Participants[participantIdx]
	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: participantIdx,
	}

	if err := ensureParticipantDarVetted(ctx, participant, contracts.CCIPReceiver); err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("ensure receiver dar on participant %d: %w", participantIdx, err)
	}

	// Deploy receiver contract
	out, err := operations.ExecuteOperation(c.e.OperationsBundle, receiver.Deploy, deps, contract.DeployInput[ccipreceiver.CCIPReceiver]{
		ChainSelector: c.chainDetails.ChainSelector,
		Qualifier:     nil,
		ActAs:         []string{partyId},
		Template: ccipreceiver.CCIPReceiver{
			Owner:         types.PARTY(partyId),
			RequiredCCVs:  nil,
			MinBlockDepth: types.INT64(0),
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
	executingParty := receiverParty
	participantIdx := c.participantIndexForParty(executingParty)
	participant := c.chain.Participants[participantIdx]
	if err := ensureExecuteParticipantDarsVetted(ctx, participant); err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("ensure execute participant dars: %w", err)
	}

	// Deploy PerPartyRouter for the receiver party
	routerAddress, routerCid, err := c.deployPerPartyRouterWithContractID(ctx, participantIdx, executingParty)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to deploy per-party router: %w", err)
	}
	c.logger.Debug().Str("RouterAddress", routerAddress.String()).Msg("Deployed PerPartyRouter")

	// Deploy CCIPReceiver contract
	receiverAddress, err := c.DeployCCIPReceiver(ctx, participantIdx, executingParty)
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
		cid, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, executingParty, templateID, address)
		if err == nil {
			return cid, nil
		}
		if !strings.Contains(err.Error(), "multiple active contracts found") && !strings.Contains(err.Error(), "no active contract found") {
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
		latestForParty := func(party string) (*apiv2.ActiveContract, error) {
			stream, streamErr := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
				ActiveAtOffset: ledgerEnd.GetOffset(),
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						party: {
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
				return nil, fmt.Errorf("get active contracts for duplicate contract fallback: %w", streamErr)
			}
			defer stream.CloseSend()
			var latest *apiv2.ActiveContract
			for {
				resp, recvErr := stream.Recv()
				if recvErr != nil {
					if errors.Is(recvErr, io.EOF) {
						break
					}

					return nil, fmt.Errorf("receive active contracts for duplicate contract fallback: %w", recvErr)
				}
				entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
				if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
					continue
				}
				if latest == nil || entry.ActiveContract.GetCreatedEvent().GetOffset() > latest.GetCreatedEvent().GetOffset() {
					latest = entry.ActiveContract
				}
			}

			return latest, nil
		}

		latest, latestErr := latestForParty(executingParty)
		if latestErr != nil {
			return "", latestErr
		}
		if latest == nil && executingParty != participant.PartyID {
			latest, latestErr = latestForParty(participant.PartyID)
			if latestErr != nil {
				return "", latestErr
			}
		}
		if latest == nil || latest.GetCreatedEvent() == nil || latest.GetCreatedEvent().GetContractId() == "" {
			return "", fmt.Errorf("no active contracts found for duplicate fallback")
		}

		return latest.GetCreatedEvent().GetContractId(), nil
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
			&message,
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
			ActAs:              uniqueParties(executingParty, participant.PartyID),
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

func (c *Chain) findPerPartyRouterContractID(
	ctx context.Context,
	participant canton.Participant,
	partyOwner string,
	instanceID string,
) (string, error) {
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", fmt.Errorf("get ledger end for router lookup: %w", err)
	}

	findForParty := func(filterParty string) (string, error) {
		stream, streamErr := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
			ActiveAtOffset: ledgerEnd.GetOffset(),
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					filterParty: {
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
				Verbose: true,
			},
		})
		if streamErr != nil {
			return "", fmt.Errorf("query active routers: %w", streamErr)
		}
		defer stream.CloseSend()

		var latestCID string
		var latestOffset int64 = -1
		for {
			resp, recvErr := stream.Recv()
			if recvErr != nil {
				if errors.Is(recvErr, io.EOF) {
					break
				}

				return "", fmt.Errorf("receive active routers: %w", recvErr)
			}
			entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
				continue
			}
			created := entry.ActiveContract.GetCreatedEvent()
			if created.GetTemplateId() == nil {
				continue
			}
			if created.GetTemplateId().GetModuleName() != "CCIP.PerPartyRouter" || created.GetTemplateId().GetEntityName() != perPartyRouterEntityName {
				continue
			}
			args := created.GetCreateArguments()
			if args == nil {
				continue
			}

			var gotOwner, gotInstanceID string
			for _, field := range args.GetFields() {
				if field == nil || field.GetValue() == nil {
					continue
				}
				switch field.GetLabel() {
				case createArgFieldPartyOwner:
					gotOwner = field.GetValue().GetParty()
				case createArgFieldInstanceID:
					gotInstanceID = field.GetValue().GetText()
				}
			}
			if gotOwner != partyOwner || gotInstanceID != instanceID {
				continue
			}
			if created.GetOffset() > latestOffset {
				latestOffset = created.GetOffset()
				latestCID = created.GetContractId()
			}
		}

		return latestCID, nil
	}

	for _, p := range uniqueParties(partyOwner, participant.PartyID) {
		cid, cidErr := findForParty(p)
		if cidErr != nil {
			return "", cidErr
		}
		if cid != "" {
			return cid, nil
		}
	}

	return "", fmt.Errorf("no active PerPartyRouter found for owner %s and instanceId %s", partyOwner, instanceID)
}

func (c *Chain) findAnyPerPartyRouterForOwner(
	ctx context.Context,
	participant canton.Participant,
	partyOwner string,
) (string, string, error) {
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", "", fmt.Errorf("get ledger end for any-router lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				partyOwner: {
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
			Verbose: true,
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("query any router for owner: %w", err)
	}
	defer stream.CloseSend()

	var latestCID string
	var latestInstanceID string
	var latestOffset int64 = -1
	for {
		resp, recvErr := stream.Recv()
		if recvErr != nil {
			if errors.Is(recvErr, io.EOF) {
				break
			}

			return "", "", fmt.Errorf("receive any router for owner: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		created := entry.ActiveContract.GetCreatedEvent()
		if created.GetTemplateId() == nil || created.GetTemplateId().GetModuleName() != "CCIP.PerPartyRouter" || created.GetTemplateId().GetEntityName() != perPartyRouterEntityName {
			continue
		}
		args := created.GetCreateArguments()
		if args == nil {
			continue
		}
		var gotOwner, gotInstanceID string
		for _, field := range args.GetFields() {
			if field == nil || field.GetValue() == nil {
				continue
			}
			switch field.GetLabel() {
			case createArgFieldPartyOwner:
				gotOwner = field.GetValue().GetParty()
			case createArgFieldInstanceID:
				gotInstanceID = field.GetValue().GetText()
			}
		}
		if gotOwner != partyOwner || gotInstanceID == "" {
			continue
		}
		if created.GetOffset() > latestOffset {
			latestOffset = created.GetOffset()
			latestCID = created.GetContractId()
			latestInstanceID = gotInstanceID
		}
	}

	return latestCID, latestInstanceID, nil
}

// BuildManualExecuteTokenTransferInput builds the input for the manual execute token transfer command.
// - selecting the correct lock/release pool,
// - matching by instrument hash vs destTokenAddress,
// - ensuring exact source pool membership in remotePools,
// - ensuring inbound rate limiter exists/resolves,
// - collecting sender pool holdings CIDs + disclosures,
// - fetching transfer factory + choice context/disclosures via scan-proxy,
// - injecting rate limiter CID in poolExtraContext.
func (c *Chain) buildManualExecuteTokenTransferInput(
	ctx context.Context,
	participant canton.Participant,
	tokenReceiverParty string,
	message *protocol.Message,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (*apiv2.Value, []*apiv2.DisclosedContract, error) {
	if message == nil || message.TokenTransfer == nil {
		return nil, nil, fmt.Errorf("token transfer message is required")
	}
	sourceSelectorKey := fmt.Sprintf("%d", message.SourceChainSelector)
	sourceSelectorNumericKey := sourceSelectorKey + "."
	sourcePoolHex := strings.ToLower(hex.EncodeToString(message.TokenTransfer.SourcePoolAddress))
	sourcePoolHex = strings.TrimPrefix(sourcePoolHex, "0x")
	sourcePoolHexTail40 := sourcePoolHex
	if len(sourcePoolHexTail40) > 40 {
		sourcePoolHexTail40 = sourcePoolHexTail40[len(sourcePoolHexTail40)-40:]
	}
	destTokenHex := strings.ToLower(hex.EncodeToString(message.TokenTransfer.DestTokenAddress))
	destTokenHex = strings.TrimPrefix(destTokenHex, "0x")
	requireInstrumentMatch := len(destTokenHex) > 0

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
	var selectedPoolPackageID string
	var selectedInstrumentFromCreateArgs string
	var tokenMatchedPool *lockreleasetokenpool.LockReleaseTokenPool
	var tokenMatchedPoolContractID string
	var fallbackPool *lockreleasetokenpool.LockReleaseTokenPool
	var fallbackPoolContractID string
	var ensureSourcePoolErr error
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
		instrumentTokenMatch := instrumentRawHex == destTokenHex ||
			instrumentKeccakCombinedHex == destTokenHex ||
			strings.HasSuffix(instrumentRawHex, destTokenHex) ||
			strings.HasSuffix(destTokenHex, instrumentRawHex) ||
			strings.HasSuffix(instrumentKeccakCombinedHex, destTokenHex) ||
			strings.HasSuffix(destTokenHex, instrumentKeccakCombinedHex)
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
			cfgText := strings.ToLower(fmt.Sprint(chainPoolCfgAny))
			remotePoolMatch = strings.Contains(cfgText, sourcePoolHex) || strings.Contains(cfgText, sourcePoolHexTail40)
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
			if remotePoolHex == sourcePoolHex ||
				remotePoolHex == sourcePoolHexTail40 ||
				strings.HasSuffix(remotePoolHex, sourcePoolHex) ||
				strings.HasSuffix(sourcePoolHex, remotePoolHex) ||
				strings.HasSuffix(remotePoolHex, sourcePoolHexTail40) ||
				strings.HasSuffix(sourcePoolHexTail40, remotePoolHex) {
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
			if tid := entry.ActiveContract.GetCreatedEvent().GetTemplateId(); tid != nil {
				selectedPoolPackageID = tid.GetPackageId()
			}
			selectedInstrumentFromCreateArgs = extractInstrumentCombinedFromCreateArgs(entry.ActiveContract.GetCreatedEvent())

			break
		}
		if requireInstrumentMatch {
			continue
		}
		selectedPool = parsed
		selectedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
		if tid := entry.ActiveContract.GetCreatedEvent().GetTemplateId(); tid != nil {
			selectedPoolPackageID = tid.GetPackageId()
		}
		selectedInstrumentFromCreateArgs = extractInstrumentCombinedFromCreateArgs(entry.ActiveContract.GetCreatedEvent())

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
			sourceSelectorNumericKey,
			sourcePoolHex,
			resolveActiveContractIDByAddress,
		)
		if ensureErr == nil {
			selectedPool = updatedPool
			selectedPoolContractID = updatedCID
			if d, dErr := getDisclosedContractByID(ctx, participant, updatedCID); dErr == nil && d.GetTemplateId() != nil {
				selectedPoolPackageID = d.GetTemplateId().GetPackageId()
			}
			selectedInstrumentFromCreateArgs = ""
		} else {
			ensureSourcePoolErr = ensureErr
		}
	}
	if selectedPool == nil && fallbackPool != nil && !requireInstrumentMatch {
		selectedPool = fallbackPool
		selectedPoolContractID = fallbackPoolContractID
	}
	if selectedPool == nil {
		if len(debugPoolCandidates) == 0 && len(c.chain.Participants) > 0 && participant.PartyID != c.chain.Participants[0].PartyID {
			return c.buildManualExecuteTokenTransferInput(
				ctx,
				c.chain.Participants[0],
				tokenReceiverParty,
				message,
				resolveActiveContractIDByAddress,
			)
		}
		if requireInstrumentMatch {
			ensureSourcePoolErrMsg := "<nil>"
			if ensureSourcePoolErr != nil {
				ensureSourcePoolErrMsg = ensureSourcePoolErr.Error()
			}

			return nil, nil, fmt.Errorf(
				"no lock/release pool found with instrument hash matching dest token %s for source selector %s and source pool %s; ensureSourcePoolErr=%s; candidates=%v",
				destTokenHex,
				sourceSelectorKey,
				sourcePoolHex,
				ensureSourcePoolErrMsg,
				debugPoolCandidates,
			)
		}

		return nil, nil, fmt.Errorf("no lock/release pool found for source selector %s and source pool %s", sourceSelectorKey, sourcePoolHex)
	}
	if selectedCfgAny, ok := findChainPoolConfigBySelector(selectedPool.ChainPoolConfigs, sourceSelectorKey); ok {
		remotePools := extractRemotePools(selectedCfgAny)
		hasSourcePool := false
		for _, rp := range remotePools {
			rpHex := strings.ToLower(strings.TrimPrefix(rp, "0x"))
			// Must be exact: on-ledger VerifyInboundMessage checks exact byte equality.
			if rpHex == sourcePoolHex {
				hasSourcePool = true
				break
			}
		}
		if !hasSourcePool {
			updatedPool, updatedCID, ensureErr := ensureManualExecuteSourcePoolAllowed(
				ctx,
				participant,
				selectedPool,
				selectedPoolContractID,
				sourceSelectorKey,
				sourceSelectorNumericKey,
				sourcePoolHex,
				resolveActiveContractIDByAddress,
			)
			if ensureErr != nil {
				return nil, nil, fmt.Errorf("ensure source pool is allowed for selected lock/release pool: %w", ensureErr)
			}
			selectedPool = updatedPool
			selectedPoolContractID = updatedCID
		}
	}
	selectedPool, selectedPoolContractID, err = ensureManualExecuteInboundRateLimiterConfigured(
		ctx,
		participant,
		selectedPool,
		selectedPoolContractID,
		sourceSelectorKey,
		sourceSelectorNumericKey,
		resolveActiveContractIDByAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("ensure inbound rate limiter is configured for selected lock/release pool: %w", err)
	}
	// Inbound limiter reconciliation may select a different active pool version.
	// Re-enforce exact 32-byte sourcePoolAddress membership in remotePools on that final version.
	if selectedCfgAny, ok := findChainPoolConfigBySelector(selectedPool.ChainPoolConfigs, sourceSelectorKey); ok {
		remotePools := extractRemotePools(selectedCfgAny)
		hasExactSourcePool := false
		for _, rp := range remotePools {
			if strings.EqualFold(strings.TrimPrefix(rp, "0x"), sourcePoolHex) {
				hasExactSourcePool = true
				break
			}
		}
		if !hasExactSourcePool {
			updatedPool, updatedCID, ensureErr := ensureManualExecuteSourcePoolAllowed(
				ctx,
				participant,
				selectedPool,
				selectedPoolContractID,
				sourceSelectorKey,
				sourceSelectorNumericKey,
				sourcePoolHex,
				resolveActiveContractIDByAddress,
			)
			if ensureErr != nil {
				return nil, nil, fmt.Errorf("re-ensure source pool is allowed on inbound-limiter-selected lock/release pool: %w", ensureErr)
			}
			selectedPool = updatedPool
			selectedPoolContractID = updatedCID
		}
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

	instrumentAdmin := string(selectedPool.InstrumentId.Admin)
	instrumentID := string(selectedPool.InstrumentId.Id)
	expectedTransferAdmin := instrumentAdmin
	transferSenderParty := string(selectedPool.PoolOwner)

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
			instrumentRecord := fields[1].GetValue().GetRecord()
			if instrumentRecord == nil || len(instrumentRecord.GetFields()) < 2 {
				continue
			}
			var holdingInstrumentAdmin, holdingInstrumentID string
			holdingOwner := fields[0].GetValue().GetParty()
			if holdingOwner != transferSenderParty {
				continue
			}
			for _, instrumentField := range instrumentRecord.GetFields() {
				if instrumentField == nil || instrumentField.GetValue() == nil {
					continue
				}
				switch instrumentField.GetLabel() {
				case "admin":
					holdingInstrumentAdmin = instrumentField.GetValue().GetParty()
				case "id":
					holdingInstrumentID = instrumentField.GetValue().GetText()
				}
			}
			isLocked := fields[3].GetValue().GetOptional().GetValue() != nil
			if isLocked {
				continue
			}
			disclosure := &apiv2.DisclosedContract{
				TemplateId:       holding.GetCreatedEvent().GetTemplateId(),
				ContractId:       holding.GetCreatedEvent().GetContractId(),
				CreatedEventBlob: holding.GetCreatedEvent().GetCreatedEventBlob(),
				SynchronizerId:   holding.GetSynchronizerId(),
			}
			if holdingInstrumentAdmin != instrumentAdmin || holdingInstrumentID != instrumentID {
				continue
			}
			poolHoldings = append(poolHoldings, holding.GetCreatedEvent().GetContractId())
			poolHoldingDisclosures = append(poolHoldingDisclosures, disclosure)
		}

		return poolHoldings, poolHoldingDisclosures, nil
	}

	poolHoldings, poolHoldingDisclosures, err := collectPoolHoldings()
	if err != nil {
		return nil, nil, err
	}
	if len(poolHoldings) == 0 {
		return nil, nil, fmt.Errorf(
			"no unlocked pool holdings found for transfer sender %s and instrument %s/%s",
			transferSenderParty,
			instrumentAdmin,
			instrumentID,
		)
	}

	transferClient, err := newTransferInstructionClient(participant)
	if err != nil {
		return nil, nil, err
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

	poolDisclosure, err := getDisclosedContractByID(ctx, participant, selectedPoolContractID)
	if err != nil {
		return nil, nil, fmt.Errorf("get token pool disclosure: %w", err)
	}
	disclosures := make([]*apiv2.DisclosedContract, 0, 2+len(transferFactoryDisclosures)+len(poolHoldingDisclosures))
	disclosures = append(disclosures, rateLimiterDisclosure, poolDisclosure)
	disclosures = append(disclosures, transferFactoryDisclosures...)
	disclosures = append(disclosures, poolHoldingDisclosures...)
	finalRemotePools := []string{}
	if selectedCfgAny, ok := findChainPoolConfigBySelector(selectedPool.ChainPoolConfigs, sourceSelectorKey); ok {
		finalRemotePools = extractRemotePools(selectedCfgAny)
	}
	instrumentCombined := fmt.Sprintf("%s@%s", instrumentID, instrumentAdmin)
	c.logger.Debug().
		Str("SelectedTokenPoolCID", selectedPoolContractID).
		Str("MessageSourcePoolHex", sourcePoolHex).
		Any("FinalRemotePools", finalRemotePools).
		Str("SelectedInstrumentCombined", instrumentCombined).
		Str("SelectedInstrumentCombinedHex", hex.EncodeToString([]byte(instrumentCombined))).
		Int("SelectedInstrumentCombinedLen", len([]byte(instrumentCombined))).
		Str("SelectedInstrumentHashHex", strings.ToLower(hex.EncodeToString(crypto.Keccak256([]byte(instrumentCombined))))).
		Str("SelectedPoolInstrumentAdmin", string(selectedPool.InstrumentId.Admin)).
		Str("ExpectedTransferAdmin", expectedTransferAdmin).
		Str("MessageDestTokenHex", destTokenHex).
		Str("SelectedTokenPoolPackageID", selectedPoolPackageID).
		Str("SelectedInstrumentFromCreateArgs", selectedInstrumentFromCreateArgs).
		Msg("Prepared manual execute token transfer input")

	return tokenTransferValue, disclosures, nil
}

func ensureManualExecuteInboundRateLimiterConfigured(
	ctx context.Context,
	participant canton.Participant,
	pool *lockreleasetokenpool.LockReleaseTokenPool,
	poolCID string,
	sourceSelectorKey string,
	sourceSelectorNumericKey string,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (*lockreleasetokenpool.LockReleaseTokenPool, string, error) {
	if pool == nil || poolCID == "" {
		return nil, "", fmt.Errorf("pool and poolCID are required")
	}
	tryResolveConfigured := func(candidatePool *lockreleasetokenpool.LockReleaseTokenPool, rawAny any) bool {
		_, _, err := resolveRateLimiterFromRawAddress(
			ctx,
			participant,
			rawAny,
			candidatePool,
			sourceSelectorKey,
			resolveActiveContractIDByAddress,
		)

		return err == nil
	}
	if raw, ok := pool.InboundRateLimiters[sourceSelectorKey]; ok && tryResolveConfigured(pool, raw) {
		return pool, poolCID, nil
	}
	if raw, ok := pool.InboundRateLimiters[sourceSelectorNumericKey]; ok && tryResolveConfigured(pool, raw) {
		return pool, poolCID, nil
	}

	// Multiple active versions of the same pool instance may exist; prefer the newest one that already has a resolvable inbound limiter.
	ledgerEndForPoolLookup, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, "", fmt.Errorf("get ledger end for lock/release pool inbound limiter lookup: %w", err)
	}
	poolStream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndForPoolLookup.GetOffset(),
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
		return nil, "", fmt.Errorf("query active lock/release pools for inbound limiter lookup: %w", err)
	}
	defer poolStream.CloseSend()
	var bestPool *lockreleasetokenpool.LockReleaseTokenPool
	var bestCID string
	var bestOffset int64
	for {
		resp, recvErr := poolStream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, "", fmt.Errorf("receive active lock/release pools for inbound limiter lookup: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		created := entry.ActiveContract.GetCreatedEvent()
		parsedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
		if parseErr != nil {
			continue
		}
		if string(parsedPool.InstanceId) != string(pool.InstanceId) || string(parsedPool.PoolOwner) != string(pool.PoolOwner) {
			continue
		}
		if raw, ok := parsedPool.InboundRateLimiters[sourceSelectorKey]; ok && tryResolveConfigured(parsedPool, raw) {
			if created.GetOffset() >= bestOffset {
				bestOffset = created.GetOffset()
				bestPool = parsedPool
				bestCID = created.GetContractId()
			}

			continue
		}
		if raw, ok := parsedPool.InboundRateLimiters[sourceSelectorNumericKey]; ok && tryResolveConfigured(parsedPool, raw) {
			if created.GetOffset() >= bestOffset {
				bestOffset = created.GetOffset()
				bestPool = parsedPool
				bestCID = created.GetContractId()
			}
		}
	}
	if bestPool != nil && bestCID != "" {
		return bestPool, bestCID, nil
	}

	selectorNorm := normalizeNumericText(sourceSelectorKey)
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, "", fmt.Errorf("get ledger end for inbound rate limiter lookup: %w", err)
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
			Verbose: true,
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("query active inbound rate limiters: %w", err)
	}
	defer stream.CloseSend()

	var replacementRaw common.RawInstanceAddress
	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, "", fmt.Errorf("receive active inbound rate limiters: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		tmpl := entry.ActiveContract.GetCreatedEvent().GetTemplateId()
		if tmpl == nil || tmpl.GetModuleName() != "CCIP.RateLimiter" || tmpl.GetEntityName() != "RateLimiter" {
			continue
		}
		parsed, parseErr := bindings.UnmarshalCreatedEvent[common.RateLimiter](entry.ActiveContract.GetCreatedEvent())
		if parseErr != nil {
			continue
		}
		if parsed.Direction != common.RateLimitDirectionRateLimitDirection_Inbound {
			continue
		}
		if string(parsed.PoolOwner) != string(pool.PoolOwner) || string(parsed.PoolInstanceId) != string(pool.InstanceId) {
			continue
		}
		if normalizeNumericText(string(parsed.RemoteChainSelector)) != selectorNorm {
			continue
		}
		replacementRaw = contracts.InstanceID(string(parsed.InstanceId)).RawInstanceAddress(parsed.PoolOwner).Binding()

		break
	}
	if replacementRaw == (common.RawInstanceAddress{}) {
		selectorForInstanceID := strings.ReplaceAll(sourceSelectorKey, ".", "-")
		instanceID := fmt.Sprintf("manualexec-inbound-rl-%s-%s", selectorForInstanceID, uuid.NewString()[:8])
		selectorNumeric := sourceSelectorKey
		if !strings.Contains(selectorNumeric, ".") {
			selectorNumeric += "."
		}
		nowMicro := time.Now().UnixMicro()
		_, createErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.New().String(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  "#ccip-common",
							ModuleName: "CCIP.RateLimiter",
							EntityName: "RateLimiter",
						},
						CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: createArgFieldInstanceID, Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}}},
							{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(pool.InstanceId)}}},
							{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(pool.PoolOwner)}}},
							{Label: "remoteChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: selectorNumeric}}},
							{Label: "direction", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
							{Label: "mode", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitMode_DefaultFinality"}}}},
							{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}}},
							{Label: "capacity", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "999999999999999999"}}},
							{Label: "rate", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "999999999999999999"}}},
							{Label: "tokens", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
							{Label: "lastUpdated", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: nowMicro}}},
						}},
					}},
				}},
				ActAs: []string{string(pool.PoolOwner)},
			},
		})
		if createErr != nil {
			return nil, "", fmt.Errorf("create inbound rate limiter for selector %s and pool %s@%s: %w", sourceSelectorKey, pool.InstanceId, pool.PoolOwner, createErr)
		}
		replacementRaw = contracts.MustNewInstanceID(instanceID).RawInstanceAddress(pool.PoolOwner).Binding()
	}

	updatedInbound := types.GENMAP{}
	maps.Copy(updatedInbound, pool.InboundRateLimiters)
	updatedInbound[sourceSelectorKey] = replacementRaw
	updatedInbound[sourceSelectorNumericKey] = replacementRaw
	updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateRateLimiters{
		NewOutboundRateLimiters: pool.OutboundRateLimiters,
		NewInboundRateLimiters:  updatedInbound,
	}
	updateRes, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
					Choice:         "LockReleaseTokenPool_UpdateRateLimiters",
					ChoiceArgument: ledger.MapToValue(updateArgs),
				}},
			}},
			ActAs: []string{string(pool.PoolOwner)},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("patch inbound rate limiter for manual execute: %w", err)
	}
	for _, ev := range updateRes.GetTransaction().GetEvents() {
		created := ev.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		tmpl := created.GetTemplateId()
		if tmpl.GetModuleName() != "CCIP.LockReleaseTokenPool" || tmpl.GetEntityName() != "LockReleaseTokenPool" {
			continue
		}
		updatedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
		if parseErr != nil {
			continue
		}
		if string(updatedPool.InstanceId) == string(pool.InstanceId) && string(updatedPool.PoolOwner) == string(pool.PoolOwner) {
			return updatedPool, created.GetContractId(), nil
		}
	}

	ledgerEndAfterUpdate, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, "", fmt.Errorf("get ledger end for updated lock/release pool lookup after inbound rate limiter patch: %w", err)
	}
	streamAfterUpdate, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndAfterUpdate.GetOffset(),
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
		return nil, "", fmt.Errorf("query active lock/release pools after inbound rate limiter patch: %w", err)
	}
	defer streamAfterUpdate.CloseSend()

	expectedLimiterRaw, err := contracts.RawInstanceAddressFromString(string(replacementRaw.Unpack))
	if err != nil {
		return nil, "", fmt.Errorf("parse replacement inbound rate limiter raw address %q: %w", replacementRaw.Unpack, err)
	}
	expectedLimiterAddr := expectedLimiterRaw.InstanceAddress().String()
	var latestMatchingPool *lockreleasetokenpool.LockReleaseTokenPool
	var latestMatchingPoolCID string
	var latestMatchingOffset int64
	for {
		resp, recvErr := streamAfterUpdate.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, "", fmt.Errorf("receive active lock/release pools after inbound rate limiter patch: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		created := entry.ActiveContract.GetCreatedEvent()
		parsedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
		if parseErr != nil {
			continue
		}
		if string(parsedPool.InstanceId) != string(pool.InstanceId) || string(parsedPool.PoolOwner) != string(pool.PoolOwner) {
			continue
		}
		rawAny, ok := parsedPool.InboundRateLimiters[sourceSelectorKey]
		if !ok {
			rawAny, ok = parsedPool.InboundRateLimiters[sourceSelectorNumericKey]
		}
		if !ok {
			continue
		}
		rawText, parseRawErr := parseRawInstanceAddress(rawAny)
		if parseRawErr != nil {
			continue
		}
		rawAddr, rawAddrErr := contracts.RawInstanceAddressFromString(rawText)
		if rawAddrErr != nil || rawAddr.InstanceAddress().String() != expectedLimiterAddr {
			continue
		}
		if created.GetOffset() >= latestMatchingOffset {
			latestMatchingOffset = created.GetOffset()
			latestMatchingPool = parsedPool
			latestMatchingPoolCID = created.GetContractId()
		}
	}
	if latestMatchingPool != nil && latestMatchingPoolCID != "" {
		return latestMatchingPool, latestMatchingPoolCID, nil
	}

	raw := contracts.InstanceID(string(pool.InstanceId)).RawInstanceAddress(pool.PoolOwner)
	updatedCID, err := resolveActiveContractIDByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), raw.InstanceAddress())
	if err != nil {
		return nil, "", fmt.Errorf("resolve updated lock/release pool contract id after inbound rate limiter patch: %w", err)
	}
	updatedActive, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		raw.InstanceAddress(),
	)
	if err != nil {
		return nil, "", fmt.Errorf("find active updated lock/release pool contract after inbound rate limiter patch: %w", err)
	}
	updatedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](updatedActive.GetCreatedEvent())
	if err != nil {
		return nil, "", fmt.Errorf("parse updated lock/release pool after inbound rate limiter patch: %w", err)
	}

	return updatedPool, updatedCID, nil
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

func extractInstrumentCombinedFromCreateArgs(created *apiv2.CreatedEvent) string {
	if created == nil || created.GetCreateArguments() == nil {
		return ""
	}
	var instrumentRecord *apiv2.Record
	for _, f := range created.GetCreateArguments().GetFields() {
		if f.GetLabel() == "instrumentId" && f.GetValue() != nil {
			instrumentRecord = f.GetValue().GetRecord()
			break
		}
	}
	if instrumentRecord == nil {
		return ""
	}
	var admin string
	var id string
	for _, f := range instrumentRecord.GetFields() {
		switch f.GetLabel() {
		case "admin":
			if f.GetValue() != nil {
				admin = f.GetValue().GetParty()
			}
		case "id":
			if f.GetValue() != nil {
				id = f.GetValue().GetText()
			}
		}
	}
	if id == "" || admin == "" {
		return ""
	}

	return id + "@" + admin
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
	sourceSelectorNumericKey string,
	sourcePoolHex string,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (*lockreleasetokenpool.LockReleaseTokenPool, string, error) {
	updatedChainPoolConfigs := types.GENMAP{}
	maps.Copy(updatedChainPoolConfigs, pool.ChainPoolConfigs)
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
	updatedChainPoolConfigs[sourceSelectorNumericKey] = existingCfg

	updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateChainPoolConfigs{
		NewChainPoolConfigs: updatedChainPoolConfigs,
	}
	updateRes, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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

	var latestCreated *apiv2.CreatedEvent
	var identityMatchedPool *lockreleasetokenpool.LockReleaseTokenPool
	var identityMatchedCID string
	for _, ev := range updateRes.GetTransaction().GetEvents() {
		created := ev.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		tmpl := created.GetTemplateId()
		if tmpl.GetModuleName() != "CCIP.LockReleaseTokenPool" || tmpl.GetEntityName() != "LockReleaseTokenPool" {
			continue
		}
		if latestCreated == nil || created.GetOffset() > latestCreated.GetOffset() {
			latestCreated = created
		}
		parsedCreated, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
		if parseErr != nil {
			continue
		}
		if string(parsedCreated.InstanceId) == string(pool.InstanceId) && string(parsedCreated.PoolOwner) == string(pool.PoolOwner) {
			identityMatchedPool = parsedCreated
			identityMatchedCID = created.GetContractId()
		}
	}
	if identityMatchedPool != nil && identityMatchedCID != "" {
		return identityMatchedPool, identityMatchedCID, nil
	}
	if latestCreated != nil {
		updatedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](latestCreated)
		if parseErr == nil {
			return updatedPool, latestCreated.GetContractId(), nil
		}
	}

	raw := contracts.InstanceID(string(pool.InstanceId)).RawInstanceAddress(pool.PoolOwner)
	updatedCID, err := resolveActiveContractIDByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), raw.InstanceAddress())
	if err == nil {
		updatedPool := *pool
		updatedPool.ChainPoolConfigs = updatedChainPoolConfigs

		return &updatedPool, updatedCID, nil
	}

	return nil, "", fmt.Errorf("resolve updated lock/release pool contract id: %w", err)
}

func resolveRateLimiterForManualExecute(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	sourceSelectorNumericKey string,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	configuredInboundErrs := make([]string, 0, 2)
	if inboundRateLimiterRaw, ok := selectedPool.InboundRateLimiters[sourceSelectorKey]; ok {
		cid, disclosure, err := resolveRateLimiterFromRawAddress(
			ctx,
			participant,
			inboundRateLimiterRaw,
			selectedPool,
			sourceSelectorKey,
			resolveActiveContractIDByAddress,
		)
		if err == nil {
			return cid, disclosure, nil
		}
		configuredInboundErrs = append(configuredInboundErrs, fmt.Sprintf("%s: %v", sourceSelectorKey, err))
	}
	if inboundRateLimiterRaw, ok := selectedPool.InboundRateLimiters[sourceSelectorNumericKey]; ok {
		cid, disclosure, err := resolveRateLimiterFromRawAddress(
			ctx,
			participant,
			inboundRateLimiterRaw,
			selectedPool,
			sourceSelectorKey,
			resolveActiveContractIDByAddress,
		)
		if err == nil {
			return cid, disclosure, nil
		}
		configuredInboundErrs = append(configuredInboundErrs, fmt.Sprintf("%s: %v", sourceSelectorNumericKey, err))
	}
	if len(configuredInboundErrs) == 0 {
		return "", nil, fmt.Errorf(
			"missing configured inbound rate limiter entry for source selector %s (keys: %s, %s)",
			sourceSelectorKey,
			sourceSelectorKey,
			sourceSelectorNumericKey,
		)
	}

	return "", nil, fmt.Errorf(
		"resolve configured inbound rate limiter for source selector %s (pool=%s@%s): %s",
		sourceSelectorKey,
		selectedPool.InstanceId,
		selectedPool.PoolOwner,
		strings.Join(configuredInboundErrs, " | "),
	)
}

func resolveRateLimiterFromRawAddress(
	ctx context.Context,
	participant canton.Participant,
	inboundRateLimiterRaw any,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	_ func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	rawRateLimiter, err := parseRawInstanceAddress(inboundRateLimiterRaw)
	if err != nil {
		return "", nil, err
	}
	rateLimiterRawAddr, err := contracts.RawInstanceAddressFromString(rawRateLimiter)
	if err != nil {
		return "", nil, fmt.Errorf("parse rate limiter raw instance address: %w", err)
	}

	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", nil, fmt.Errorf("get ledger end for rate limiter lookup: %w", err)
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
			Verbose: true,
		},
	})
	if err != nil {
		return "", nil, fmt.Errorf("query active rate limiters by raw address: %w", err)
	}
	defer stream.CloseSend()
	selectorNorm := normalizeNumericText(sourceSelectorKey)
	poolSelectorCandidates := make([]string, 0, 6)
	instanceCandidates := make([]string, 0, 3)

	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return "", nil, fmt.Errorf("receive active rate limiters by raw address: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		tmpl := entry.ActiveContract.GetCreatedEvent().GetTemplateId()
		if tmpl == nil || tmpl.GetModuleName() != "CCIP.RateLimiter" || tmpl.GetEntityName() != "RateLimiter" {
			continue
		}
		parsed, parseErr := bindings.UnmarshalCreatedEvent[common.RateLimiter](entry.ActiveContract.GetCreatedEvent())
		if parseErr != nil {
			continue
		}
		candidateRaw := contracts.InstanceID(string(parsed.InstanceId)).RawInstanceAddress(parsed.PoolOwner)
		if parsed.Direction == common.RateLimitDirectionRateLimitDirection_Inbound &&
			(selectedPool == nil || (string(parsed.PoolOwner) == string(selectedPool.PoolOwner) && string(parsed.PoolInstanceId) == string(selectedPool.InstanceId))) &&
			(selectorNorm == "" || normalizeNumericText(string(parsed.RemoteChainSelector)) == selectorNorm) &&
			len(poolSelectorCandidates) < cap(poolSelectorCandidates) {
			poolSelectorCandidates = append(poolSelectorCandidates, fmt.Sprintf("%s=>%s", parsed.InstanceId, candidateRaw.InstanceAddress().String()))
		}
		if candidateRaw.InstanceAddress().String() != rateLimiterRawAddr.InstanceAddress().String() {
			continue
		}
		if len(instanceCandidates) < cap(instanceCandidates) {
			instanceCandidates = append(instanceCandidates, fmt.Sprintf("instanceId=%s pool=%s@%s selector=%s dir=%s", parsed.InstanceId, parsed.PoolInstanceId, parsed.PoolOwner, parsed.RemoteChainSelector, parsed.Direction))
		}
		if parsed.Direction != common.RateLimitDirectionRateLimitDirection_Inbound {
			continue
		}
		if selectedPool != nil {
			if string(parsed.PoolOwner) != string(selectedPool.PoolOwner) || string(parsed.PoolInstanceId) != string(selectedPool.InstanceId) {
				continue
			}
		}
		if selectorNorm != "" && normalizeNumericText(string(parsed.RemoteChainSelector)) != selectorNorm {
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

	return "", nil, fmt.Errorf(
		"no active inbound rate limiter matched raw=%s instance address %s for pool %s@%s selector %s (poolSelectorCandidates=%v instanceCandidates=%v)",
		rawRateLimiter,
		rateLimiterRawAddr.InstanceAddress().String(),
		selectedPool.InstanceId,
		selectedPool.PoolOwner,
		sourceSelectorKey,
		poolSelectorCandidates,
		instanceCandidates,
	)
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
