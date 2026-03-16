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
	perPartyRouterEntityName  = "PerPartyRouter"
	createArgFieldPartyOwner  = "partyOwner"
	createArgFieldInstanceID  = "instanceId"
	rateLimiterModuleName     = "CCIP.RateLimiter"
	rateLimiterEntityName     = "RateLimiter"
	lockReleasePoolModuleName = "CCIP.LockReleaseTokenPool"
	lockReleasePoolEntityName = "LockReleaseTokenPool"
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

func (c *Chain) DeployCCIPReceiver(ctx context.Context, participantIdx int, partyId string, minBlockDepth int64) (contracts.InstanceAddress, error) {
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
			MinBlockDepth: types.INT64(minBlockDepth),
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

	// Deploy CCIPReceiver after token transfer prep so its min depth matches the effective message finality.
	receiverAddress, err := c.DeployCCIPReceiver(ctx, participantIdx, executingParty, int64(message.Finality))
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to deploy CCIPReceiver contract: %w", err)
	}
	c.logger.Debug().Str("ReceiverAddress", receiverAddress.String()).Msg("Deployed CCIPReceiver")

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
// - injecting inbound rate limiter CIDs in poolExtraContext.
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
	requireRemotePoolMatch := len(sourcePoolHex) > 0

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
	var selectorMatchedPool *lockreleasetokenpool.LockReleaseTokenPool
	var selectorMatchedPoolContractID string
	selectorMatchedPoolCount := 0
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
		rawRemotePools, rawChainPoolCfgFound := extractRemotePoolsFromCreateArgs(entry.ActiveContract.GetCreatedEvent(), sourceSelectorKey)
		if !ok {
			if !rawChainPoolCfgFound {
				continue
			}
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
		if len(rawRemotePools) > 0 {
			remotePools = rawRemotePools
		}
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
			if instrumentTokenMatch {
				selectorMatchedPoolCount++
				if selectorMatchedPool == nil {
					selectorMatchedPool = parsed
					selectorMatchedPoolContractID = entry.ActiveContract.GetCreatedEvent().GetContractId()
				}
			}
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
	if selectedPool == nil && tokenMatchedPool != nil && requireInstrumentMatch && !requireRemotePoolMatch {
		selectedPool = tokenMatchedPool
		selectedPoolContractID = tokenMatchedPoolContractID
	}
	if selectedPool == nil && selectorMatchedPool != nil && requireInstrumentMatch && selectorMatchedPoolCount == 1 {
		selectedPool = selectorMatchedPool
		selectedPoolContractID = selectorMatchedPoolContractID
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
			return nil, nil, fmt.Errorf(
				"no lock/release pool found with instrument hash matching dest token %s and remote pool %s for source selector %s; candidates=%v",
				destTokenHex,
				sourcePoolHex,
				sourceSelectorKey,
				debugPoolCandidates,
			)
		}

		return nil, nil, fmt.Errorf("no lock/release pool found for source selector %s and source pool %s", sourceSelectorKey, sourcePoolHex)
	}
	if len(c.chain.Participants) > 0 &&
		participant.PartyID != c.chain.Participants[0].PartyID &&
		string(selectedPool.PoolOwner) == c.chain.Participants[0].PartyID {
		return c.buildManualExecuteTokenTransferInput(
			ctx,
			c.chain.Participants[0],
			tokenReceiverParty,
			message,
			resolveActiveContractIDByAddress,
		)
	}
	selectedPool, selectedPoolContractID, err = ensureManualExecuteSourcePoolConfig(
		ctx,
		participant,
		selectedPool,
		selectedPoolContractID,
		sourceSelectorKey,
		sourcePoolHex,
		message.Finality,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("ensure source pool config for manual execute: %w", err)
	}
	selectedPool, selectedPoolContractID, err = ensureManualExecuteInboundRateLimiters(
		ctx,
		participant,
		selectedPool,
		selectedPoolContractID,
		sourceSelectorKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("ensure inbound rate limiters for manual execute: %w", err)
	}
	defaultRateLimiterCID, defaultRateLimiterDisclosure, err := resolveRateLimiterForManualExecute(
		ctx,
		participant,
		selectedPool,
		sourceSelectorKey,
		sourceSelectorNumericKey,
		selectedPool.InboundRateLimiters,
		common.RateLimitModeRateLimitMode_DefaultFinality,
		resolveActiveContractIDByAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve default-finality rate limiter disclosure: %w", err)
	}
	customRateLimiterCID, customRateLimiterDisclosure, err := resolveRateLimiterForManualExecute(
		ctx,
		participant,
		selectedPool,
		sourceSelectorKey,
		sourceSelectorNumericKey,
		selectedPool.InboundCustomRateLimiters,
		common.RateLimitModeRateLimitMode_CustomFinality,
		resolveActiveContractIDByAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve custom-finality rate limiter disclosure: %w", err)
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
			"inbound-rate-limiter": map[string]any{
				"tag":   "AV_ContractId",
				"value": defaultRateLimiterCID,
			},
			"inbound-custom-rate-limiter": map[string]any{
				"tag":   "AV_ContractId",
				"value": customRateLimiterCID,
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
	disclosures := make([]*apiv2.DisclosedContract, 0, 3+len(transferFactoryDisclosures)+len(poolHoldingDisclosures))
	disclosures = append(disclosures, defaultRateLimiterDisclosure, customRateLimiterDisclosure, poolDisclosure)
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
		Msg("Prepared manual execute token transfer input")

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

func extractRemotePoolsFromCreateArgs(created *apiv2.CreatedEvent, sourceSelectorKey string) ([]string, bool) {
	if created == nil || created.GetCreateArguments() == nil {
		return nil, false
	}

	var chainPoolConfigs *apiv2.GenMap
	for _, f := range created.GetCreateArguments().GetFields() {
		if f.GetLabel() == "chainPoolConfigs" && f.GetValue() != nil {
			chainPoolConfigs = f.GetValue().GetGenMap()
			break
		}
	}
	if chainPoolConfigs == nil {
		return nil, false
	}

	sourceSelectorNorm := normalizeNumericText(sourceSelectorKey)
	for _, entry := range chainPoolConfigs.GetEntries() {
		if normalizeNumericText(entry.GetKey().GetNumeric()) != sourceSelectorNorm {
			continue
		}
		record := entry.GetValue().GetRecord()
		if record == nil {
			return nil, true
		}
		for _, field := range record.GetFields() {
			if field.GetLabel() != "remotePools" || field.GetValue() == nil {
				continue
			}
			list := field.GetValue().GetList()
			if list == nil {
				return nil, true
			}
			out := make([]string, 0, len(list.GetElements()))
			for _, elem := range list.GetElements() {
				out = append(out, elem.GetText())
			}

			return out, true
		}

		return nil, true
	}

	return nil, false
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

func ensureManualExecuteSourcePoolConfig(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	selectedPoolContractID string,
	sourceSelectorKey string,
	sourcePoolHex string,
	finality uint16,
) (*lockreleasetokenpool.LockReleaseTokenPool, string, error) {
	updatedChainPoolConfigs := typedChainPoolConfigs(selectedPool.ChainPoolConfigs)
	sourceSelectorNorm := normalizeNumericForCompare(sourceSelectorKey)
	var chainPoolCfg lockreleasetokenpool.ChainPoolConfig
	found := false
	for rawKey, rawCfg := range updatedChainPoolConfigs {
		if normalizeNumericForCompare(rawKey) != sourceSelectorNorm {
			continue
		}
		typedCfg, ok := rawCfg.(lockreleasetokenpool.ChainPoolConfig)
		if !ok {
			return nil, "", fmt.Errorf("unexpected chain pool config type %T for selector %s", rawCfg, rawKey)
		}
		chainPoolCfg = typedCfg
		found = true
		break
	}
	if !found {
		return nil, "", fmt.Errorf("missing chain pool config for selector %s", sourceSelectorKey)
	}

	needsUpdate := false
	targetMinBlockDepth := int64(0)
	if finality > 0 {
		targetMinBlockDepth = 1
	}
	if int64(chainPoolCfg.MinBlockDepth) != targetMinBlockDepth {
		chainPoolCfg.MinBlockDepth = types.INT64(targetMinBlockDepth)
		needsUpdate = true
	}
	if sourcePoolHex != "" && len(chainPoolCfg.RemotePools) == 0 {
		chainPoolCfg.RemotePools = []types.TEXT{types.TEXT(canonicalCantonRemotePoolHex(sourcePoolHex))}
		needsUpdate = true
	}
	if !needsUpdate {
		return selectedPool, selectedPoolContractID, nil
	}

	for rawKey := range updatedChainPoolConfigs {
		if normalizeNumericForCompare(rawKey) == sourceSelectorNorm {
			delete(updatedChainPoolConfigs, rawKey)
		}
	}
	updatedChainPoolConfigs[sourceSelectorKey] = chainPoolCfg
	updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateChainPoolConfigs{
		NewChainPoolConfigs: updatedChainPoolConfigs,
	}
	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  "#ccip-lockreleasetokenpool",
						ModuleName: "CCIP.LockReleaseTokenPool",
						EntityName: "LockReleaseTokenPool",
					},
					ContractId:     selectedPoolContractID,
					Choice:         "LockReleaseTokenPool_UpdateChainPoolConfigs",
					ChoiceArgument: ledger.MapToValue(updateArgs.ToMap()),
				}},
			}},
			ActAs: []string{string(selectedPool.PoolOwner)},
		},
	})
	if err != nil {
		return nil, "", err
	}

	for _, event := range res.GetTransaction().GetEvents() {
		created := event.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		if created.GetTemplateId().GetEntityName() != "LockReleaseTokenPool" {
			continue
		}
		updatedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
		if parseErr != nil {
			return nil, "", fmt.Errorf("parse updated lock/release pool: %w", parseErr)
		}
		if updatedPool.InstanceId != selectedPool.InstanceId {
			continue
		}

		return updatedPool, created.GetContractId(), nil
	}

	return nil, "", fmt.Errorf("updated lock/release pool create event not found in transaction")
}

func ensureManualExecuteInboundRateLimiters(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	selectedPoolContractID string,
	sourceSelectorKey string,
) (*lockreleasetokenpool.LockReleaseTokenPool, string, error) {
	defaultRaw, _, _, err := resolveActiveRateLimiterForPoolSelectorAndMode(
		ctx,
		participant,
		selectedPool,
		sourceSelectorKey,
		common.RateLimitModeRateLimitMode_DefaultFinality,
	)
	if err != nil {
		return nil, "", err
	}
	customRaw, _, _, err := resolveActiveRateLimiterForPoolSelectorAndMode(
		ctx,
		participant,
		selectedPool,
		sourceSelectorKey,
		common.RateLimitModeRateLimitMode_CustomFinality,
	)
	if err != nil {
		return nil, "", err
	}

	updatedInbound := types.GENMAP{}
	maps.Copy(updatedInbound, selectedPool.InboundRateLimiters)
	updatedInboundCustom := types.GENMAP{}
	maps.Copy(updatedInboundCustom, selectedPool.InboundCustomRateLimiters)
	defaultRawAddr, err := contracts.RawInstanceAddressFromString(defaultRaw)
	if err != nil {
		return nil, "", fmt.Errorf("parse default inbound rate limiter raw address: %w", err)
	}
	customRawAddr, err := contracts.RawInstanceAddressFromString(customRaw)
	if err != nil {
		return nil, "", fmt.Errorf("parse custom inbound rate limiter raw address: %w", err)
	}
	sourceSelectorNorm := normalizeNumericForCompare(sourceSelectorKey)
	inboundNeedsUpdate := false
	inboundFound := false
	for rawKey, rawValue := range updatedInbound {
		if normalizeNumericForCompare(rawKey) != sourceSelectorNorm {
			continue
		}
		inboundFound = true
		configuredRaw, parseErr := parseRawInstanceAddress(rawValue)
		if parseErr != nil || configuredRaw != defaultRaw {
			inboundNeedsUpdate = true
		}
	}
	for rawKey := range updatedInbound {
		if normalizeNumericForCompare(rawKey) == sourceSelectorNorm {
			delete(updatedInbound, rawKey)
		}
	}
	if !inboundFound {
		inboundNeedsUpdate = true
	}
	updatedInbound[sourceSelectorKey] = defaultRawAddr.Binding()
	inboundCustomNeedsUpdate := false
	inboundCustomFound := false
	for rawKey, rawValue := range updatedInboundCustom {
		if normalizeNumericForCompare(rawKey) != sourceSelectorNorm {
			continue
		}
		inboundCustomFound = true
		configuredRaw, parseErr := parseRawInstanceAddress(rawValue)
		if parseErr != nil || configuredRaw != customRaw {
			inboundCustomNeedsUpdate = true
		}
	}
	for rawKey := range updatedInboundCustom {
		if normalizeNumericForCompare(rawKey) == sourceSelectorNorm {
			delete(updatedInboundCustom, rawKey)
		}
	}
	if !inboundCustomFound {
		inboundCustomNeedsUpdate = true
	}
	updatedInboundCustom[sourceSelectorKey] = customRawAddr.Binding()
	if !inboundNeedsUpdate && !inboundCustomNeedsUpdate {
		return selectedPool, selectedPoolContractID, nil
	}

	updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateRateLimiters{
		NewOutboundRateLimiters:      selectedPool.OutboundRateLimiters,
		NewInboundRateLimiters:       updatedInbound,
		NewInboundCustomRateLimiters: updatedInboundCustom,
	}
	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  "#ccip-lockreleasetokenpool",
						ModuleName: "CCIP.LockReleaseTokenPool",
						EntityName: "LockReleaseTokenPool",
					},
					ContractId:     selectedPoolContractID,
					Choice:         "LockReleaseTokenPool_UpdateRateLimiters",
					ChoiceArgument: ledger.MapToValue(updateArgs.ToMap()),
				}},
			}},
			ActAs: []string{string(selectedPool.PoolOwner)},
		},
	})
	if err != nil {
		return nil, "", err
	}

	for _, event := range res.GetTransaction().GetEvents() {
		created := event.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		if created.GetTemplateId().GetEntityName() != "LockReleaseTokenPool" {
			continue
		}
		updatedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
		if parseErr != nil {
			return nil, "", fmt.Errorf("parse updated lock/release pool after rate limiter repair: %w", parseErr)
		}
		if updatedPool.InstanceId != selectedPool.InstanceId {
			continue
		}

		return updatedPool, created.GetContractId(), nil
	}

	return nil, "", fmt.Errorf("updated lock/release pool create event not found after rate limiter repair")
}

func resolveRateLimiterForManualExecute(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	sourceSelectorNumericKey string,
	selectedRateLimiters types.GENMAP,
	expectedMode common.RateLimitMode,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	configuredInboundErrs := make([]string, 0, 2)
	if inboundRateLimiterRaw, ok := selectedRateLimiters[sourceSelectorKey]; ok {
		cid, disclosure, err := resolveRateLimiterFromRawAddress(
			ctx,
			participant,
			inboundRateLimiterRaw,
			selectedPool,
			sourceSelectorKey,
			expectedMode,
			resolveActiveContractIDByAddress,
		)
		if err == nil {
			return cid, disclosure, nil
		}
		configuredInboundErrs = append(configuredInboundErrs, fmt.Sprintf("%s: %v", sourceSelectorKey, err))
	}
	if inboundRateLimiterRaw, ok := selectedRateLimiters[sourceSelectorNumericKey]; ok {
		cid, disclosure, err := resolveRateLimiterFromRawAddress(
			ctx,
			participant,
			inboundRateLimiterRaw,
			selectedPool,
			sourceSelectorKey,
			expectedMode,
			resolveActiveContractIDByAddress,
		)
		if err == nil {
			return cid, disclosure, nil
		}
		configuredInboundErrs = append(configuredInboundErrs, fmt.Sprintf("%s: %v", sourceSelectorNumericKey, err))
	}
	if len(configuredInboundErrs) == 0 {
		_, cid, disclosure, err := resolveActiveRateLimiterForPoolSelectorAndMode(
			ctx,
			participant,
			selectedPool,
			sourceSelectorKey,
			expectedMode,
		)
		if err == nil {
			return cid, disclosure, nil
		}

		return "", nil, fmt.Errorf(
			"missing configured inbound rate limiter entry for source selector %s (keys: %s, %s): %w",
			sourceSelectorKey,
			sourceSelectorKey,
			sourceSelectorNumericKey,
			err,
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

func resolveActiveRateLimiterForPoolSelectorAndMode(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	expectedMode common.RateLimitMode,
) (string, string, *apiv2.DisclosedContract, error) {
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", "", nil, fmt.Errorf("get ledger end for rate limiter lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{{
						IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
							WildcardFilter: &apiv2.WildcardFilter{
								IncludeCreatedEventBlob: true,
							},
						},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("query active rate limiters by pool/mode: %w", err)
	}
	defer stream.CloseSend()

	selectorNorm := normalizeNumericText(sourceSelectorKey)
	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return "", "", nil, fmt.Errorf("receive active rate limiters by pool/mode: %w", recvErr)
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
		if parsed.Direction != common.RateLimitDirectionRateLimitDirection_Inbound || parsed.Mode != expectedMode {
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
		rawAddr := contracts.MustNewInstanceID(string(parsed.InstanceId)).RawInstanceAddress(parsed.PoolOwner).String()
		return rawAddr, cid, &apiv2.DisclosedContract{
			TemplateId:       entry.ActiveContract.GetCreatedEvent().GetTemplateId(),
			ContractId:       cid,
			CreatedEventBlob: entry.ActiveContract.GetCreatedEvent().GetCreatedEventBlob(),
			SynchronizerId:   entry.ActiveContract.GetSynchronizerId(),
		}, nil
	}

	return "", "", nil, fmt.Errorf(
		"no active inbound rate limiter found for pool %s@%s selector %s mode %s",
		selectedPool.InstanceId,
		selectedPool.PoolOwner,
		sourceSelectorKey,
		expectedMode,
	)
}

func resolveRateLimiterFromRawAddress(
	ctx context.Context,
	participant canton.Participant,
	inboundRateLimiterRaw any,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	expectedMode common.RateLimitMode,
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
	var matchingTemplate *common.RateLimiter

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
			parsed.Mode == expectedMode &&
			(selectedPool == nil || (string(parsed.PoolOwner) == string(selectedPool.PoolOwner) && string(parsed.PoolInstanceId) == string(selectedPool.InstanceId))) &&
			(selectorNorm == "" || normalizeNumericText(string(parsed.RemoteChainSelector)) == selectorNorm) &&
			len(poolSelectorCandidates) < cap(poolSelectorCandidates) {
			poolSelectorCandidates = append(poolSelectorCandidates, fmt.Sprintf("%s=>%s", parsed.InstanceId, candidateRaw.InstanceAddress().String()))
			if matchingTemplate == nil {
				candidate := parsed
				matchingTemplate = candidate
			}
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
		if parsed.Mode != expectedMode {
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

	if matchingTemplate != nil {
		createdCID, createdDisclosure, createErr := createExactInboundRateLimiterForRawAddress(
			ctx,
			participant,
			rawRateLimiter,
			selectedPool,
			sourceSelectorKey,
			expectedMode,
			matchingTemplate,
		)
		if createErr == nil {
			return createdCID, createdDisclosure, nil
		}
		instanceCandidates = append(instanceCandidates, fmt.Sprintf("createErr=%v", createErr))
	}

	return "", nil, fmt.Errorf(
		"no active inbound rate limiter matched raw=%s instance address %s for pool %s@%s selector %s mode %s (poolSelectorCandidates=%v instanceCandidates=%v)",
		rawRateLimiter,
		rateLimiterRawAddr.InstanceAddress().String(),
		selectedPool.InstanceId,
		selectedPool.PoolOwner,
		sourceSelectorKey,
		expectedMode,
		poolSelectorCandidates,
		instanceCandidates,
	)
}

func createExactInboundRateLimiterForRawAddress(
	ctx context.Context,
	participant canton.Participant,
	rawRateLimiter string,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	expectedMode common.RateLimitMode,
	template *common.RateLimiter,
) (string, *apiv2.DisclosedContract, error) {
	rawAddr, err := contracts.RawInstanceAddressFromString(rawRateLimiter)
	if err != nil {
		return "", nil, fmt.Errorf("parse expected raw rate limiter address: %w", err)
	}
	parts := strings.SplitN(rawAddr.String(), "@", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid expected raw rate limiter address: %s", rawRateLimiter)
	}
	instanceID := parts[0]
	ownerParty := parts[1]
	if selectedPool != nil && string(selectedPool.PoolOwner) != ownerParty {
		return "", nil, fmt.Errorf("expected raw owner %s does not match selected pool owner %s", ownerParty, selectedPool.PoolOwner)
	}
	selectorValue := sourceSelectorKey
	if template != nil && string(template.RemoteChainSelector) != "" {
		selectorValue = string(template.RemoteChainSelector)
	}
	isEnabled := false
	capacity := types.NUMERIC("999999999999999999")
	rate := types.NUMERIC("999999999999999999")
	tokens := types.NUMERIC("999999999999999999")
	modeCtor := "RateLimitMode_DefaultFinality"
	if expectedMode == common.RateLimitModeRateLimitMode_CustomFinality {
		modeCtor = "RateLimitMode_CustomFinality"
	}
	nowMicro := time.Now().UnixMicro()
	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}}},
						{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(selectedPool.InstanceId)}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ownerParty}}},
						{Label: "remoteChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: selectorValue}}},
						{Label: "direction", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
						{Label: "mode", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: modeCtor}}}},
						{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: isEnabled}}},
						{Label: "capacity", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: string(capacity)}}},
						{Label: "rate", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: string(rate)}}},
						{Label: "tokens", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: string(tokens)}}},
						{Label: "lastUpdated", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: nowMicro}}},
					}},
				}},
			}},
			ActAs: []string{ownerParty},
		},
	})
	if err != nil {
		return "", nil, err
	}
	for _, event := range res.GetTransaction().GetEvents() {
		created := event.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		if created.GetTemplateId().GetModuleName() != "CCIP.RateLimiter" || created.GetTemplateId().GetEntityName() != "RateLimiter" {
			continue
		}
		disclosure, disclosureErr := getDisclosedContractByID(ctx, participant, created.GetContractId())
		if disclosureErr != nil {
			return "", nil, fmt.Errorf("get disclosed contract for exact inbound rate limiter: %w", disclosureErr)
		}

		return created.GetContractId(), disclosure, nil
	}

	return "", nil, fmt.Errorf("created exact inbound rate limiter event not found")
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
