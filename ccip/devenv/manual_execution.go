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

	// Deploy CCIPReceiver contract
	receiverAddress, err := c.DeployCCIPReceiver(ctx, participantIdx, executingParty, int64(message.Finality))
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
// - selecting the correct preconfigured lock/release pool,
// - matching by instrument hash vs destTokenAddress,
// - requiring exact source pool membership in remotePools,
// - resolving inbound rate limiters,
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
	selectorInfo := newManualExecuteSelectorInfo(message)
	executionParticipant := participant
	resolveByAddress := resolveActiveContractIDByAddress

	selectedPool, err := c.selectManualExecuteTokenPool(ctx, executionParticipant, message, selectorInfo)
	if err != nil {
		return nil, nil, err
	}
	if poolOwnerParty := string(selectedPool.pool.PoolOwner); poolOwnerParty != "" && executionParticipant.PartyID != poolOwnerParty {
		executionParticipant = c.chain.Participants[c.participantIndexForParty(poolOwnerParty)]
		resolveByAddress = func(templateID string, address contracts.InstanceAddress) (string, error) {
			return contract.FindActiveContractIDByInstanceAddress(
				ctx,
				executionParticipant.LedgerServices.State,
				executionParticipant.PartyID,
				templateID,
				address,
			)
		}
		selectedPool, err = c.selectManualExecuteTokenPool(ctx, executionParticipant, message, selectorInfo)
		if err != nil {
			return nil, nil, err
		}
	}
	defaultRateLimiterCID, defaultRateLimiterDisclosure, err := resolveRateLimiterForManualExecute(
		ctx,
		executionParticipant,
		selectedPool.pool,
		selectorInfo.sourceSelectorKey,
		selectorInfo.sourceSelectorNumericKey,
		selectedPool.pool.InboundRateLimiters,
		resolveByAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve default-finality rate limiter disclosure: %w", err)
	}
	customRateLimiterCID, customRateLimiterDisclosure, err := resolveRateLimiterForManualExecute(
		ctx,
		executionParticipant,
		selectedPool.pool,
		selectorInfo.sourceSelectorKey,
		selectorInfo.sourceSelectorNumericKey,
		selectedPool.pool.InboundCustomRateLimiters,
		resolveByAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve custom-finality rate limiter disclosure: %w", err)
	}

	instrumentAdmin := string(selectedPool.pool.InstrumentId.Admin)
	instrumentID := string(selectedPool.pool.InstrumentId.Id)
	expectedTransferAdmin := instrumentAdmin
	transferSenderParty := string(selectedPool.pool.PoolOwner)
	poolHoldings, poolHoldingDisclosures, err := collectUnlockedPoolHoldings(ctx, executionParticipant, transferSenderParty, instrumentAdmin, instrumentID)
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

	transferClient, err := newTransferInstructionClient(executionParticipant)
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
	tokenTransferValue := buildManualExecuteTokenTransferValue(
		selectedPool.contractID,
		tokenReceiverParty,
		transferFactoryCID,
		transferFactoryCtx,
		poolHoldings,
		poolExtraContext,
	)

	poolDisclosure, err := getDisclosedContractByID(ctx, executionParticipant, selectedPool.contractID)
	if err != nil {
		return nil, nil, fmt.Errorf("get token pool disclosure: %w", err)
	}
	disclosures := make([]*apiv2.DisclosedContract, 0, 3+len(transferFactoryDisclosures)+len(poolHoldingDisclosures))
	disclosures = append(disclosures, defaultRateLimiterDisclosure, customRateLimiterDisclosure, poolDisclosure)
	disclosures = append(disclosures, transferFactoryDisclosures...)
	disclosures = append(disclosures, poolHoldingDisclosures...)
	finalRemotePools := []string{}
	if selectedCfgAny, ok := findChainPoolConfigBySelector(selectedPool.pool.ChainPoolConfigs, selectorInfo.sourceSelectorKey); ok {
		finalRemotePools = extractRemotePools(selectedCfgAny)
	}
	instrumentCombined := fmt.Sprintf("%s@%s", instrumentID, instrumentAdmin)
	c.logger.Debug().
		Str("SelectedTokenPoolCID", selectedPool.contractID).
		Str("MessageSourcePoolHex", selectorInfo.sourcePoolHex).
		Any("FinalRemotePools", finalRemotePools).
		Str("SelectedInstrumentCombined", instrumentCombined).
		Str("SelectedInstrumentCombinedHex", hex.EncodeToString([]byte(instrumentCombined))).
		Int("SelectedInstrumentCombinedLen", len([]byte(instrumentCombined))).
		Str("SelectedInstrumentHashHex", strings.ToLower(hex.EncodeToString(crypto.Keccak256([]byte(instrumentCombined))))).
		Str("SelectedPoolInstrumentAdmin", string(selectedPool.pool.InstrumentId.Admin)).
		Str("ExpectedTransferAdmin", expectedTransferAdmin).
		Str("MessageDestTokenHex", selectorInfo.destTokenHex).
		Str("SelectedTokenPoolPackageID", selectedPool.packageID).
		Str("SelectedInstrumentFromCreateArgs", selectedPool.instrumentCreate).
		Msg("Prepared manual execute token transfer input")

	return tokenTransferValue, disclosures, nil
}

type manualExecuteSelectorInfo struct {
	sourceSelectorKey        string
	sourceSelectorNumericKey string
	sourcePoolHex            string
	sourcePoolHexTail40      string
	destTokenHex             string
	requireInstrumentMatch   bool
}

type manualExecutePoolCandidate struct {
	pool             *lockreleasetokenpool.LockReleaseTokenPool
	contractID       string
	offset           int64
	packageID        string
	instrumentCreate string
}

func newManualExecuteSelectorInfo(message *protocol.Message) manualExecuteSelectorInfo {
	sourceSelectorKey := fmt.Sprintf("%d", message.SourceChainSelector)
	sourcePoolHex := strings.TrimPrefix(strings.ToLower(hex.EncodeToString(message.TokenTransfer.SourcePoolAddress)), "0x")
	sourcePoolHexTail40 := sourcePoolHex
	if len(sourcePoolHexTail40) > 40 {
		sourcePoolHexTail40 = sourcePoolHexTail40[len(sourcePoolHexTail40)-40:]
	}
	destTokenHex := strings.TrimPrefix(strings.ToLower(hex.EncodeToString(message.TokenTransfer.DestTokenAddress)), "0x")

	return manualExecuteSelectorInfo{
		sourceSelectorKey:        sourceSelectorKey,
		sourceSelectorNumericKey: sourceSelectorKey + ".",
		sourcePoolHex:            sourcePoolHex,
		sourcePoolHexTail40:      sourcePoolHexTail40,
		destTokenHex:             destTokenHex,
		requireInstrumentMatch:   destTokenHex != "",
	}
}

func (c *Chain) selectManualExecuteTokenPool(
	ctx context.Context,
	participant canton.Participant,
	message *protocol.Message,
	info manualExecuteSelectorInfo,
) (*manualExecutePoolCandidate, error) {
	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end for token pool lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{{
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
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query active lock/release pools: %w", err)
	}
	defer stream.CloseSend()

	var selected *manualExecutePoolCandidate
	debugPoolCandidates := make([]string, 0, 8)
	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, fmt.Errorf("receive lock/release pools: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		candidate, err := buildManualExecutePoolCandidate(entry.ActiveContract.GetCreatedEvent())
		if err != nil {
			return nil, err
		}
		chainPoolCfgAny, ok := findChainPoolConfigBySelector(candidate.pool.ChainPoolConfigs, info.sourceSelectorKey)
		if !ok {
			continue
		}
		instrumentMatch, remotePoolMatch, finalityCompatible, debugText := manualExecutePoolMatches(message, info, candidate.pool, chainPoolCfgAny)
		if len(debugPoolCandidates) < cap(debugPoolCandidates) {
			debugPoolCandidates = append(debugPoolCandidates, fmt.Sprintf("poolCID=%s %s", candidate.contractID, debugText))
		}
		if !remotePoolMatch {
			continue
		}
		if !finalityCompatible {
			continue
		}
		if info.requireInstrumentMatch && !instrumentMatch {
			continue
		}
		selected = newerManualExecutePoolCandidate(selected, candidate)
	}
	if selected == nil {
		if info.requireInstrumentMatch {
			return nil, fmt.Errorf(
				"no lock/release pool found with instrument hash matching dest token %s for source selector %s and source pool %s; candidates=%v",
				info.destTokenHex,
				info.sourceSelectorKey,
				info.sourcePoolHex,
				debugPoolCandidates,
			)
		}
		return nil, fmt.Errorf("no lock/release pool found for source selector %s and source pool %s", info.sourceSelectorKey, info.sourcePoolHex)
	}

	return selected, nil
}

func buildManualExecutePoolCandidate(created *apiv2.CreatedEvent) (*manualExecutePoolCandidate, error) {
	pool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](created)
	if err != nil {
		return nil, fmt.Errorf("parse lock/release pool: %w", err)
	}
	candidate := &manualExecutePoolCandidate{
		pool:             pool,
		contractID:       created.GetContractId(),
		offset:           created.GetOffset(),
		instrumentCreate: extractInstrumentCombinedFromCreateArgs(created),
	}
	if tid := created.GetTemplateId(); tid != nil {
		candidate.packageID = tid.GetPackageId()
	}

	return candidate, nil
}

func newerManualExecutePoolCandidate(current, next *manualExecutePoolCandidate) *manualExecutePoolCandidate {
	if next == nil {
		return current
	}
	if current == nil || next.offset > current.offset {
		return next
	}

	return current
}

func manualExecutePoolMatches(
	message *protocol.Message,
	info manualExecuteSelectorInfo,
	pool *lockreleasetokenpool.LockReleaseTokenPool,
	chainPoolCfgAny any,
) (bool, bool, bool, string) {
	instrumentCombined := fmt.Sprintf("%s@%s", string(pool.InstrumentId.Id), string(pool.InstrumentId.Admin))
	instrumentRawHex := strings.ToLower(hex.EncodeToString([]byte(instrumentCombined)))
	instrumentKeccakHex := strings.ToLower(hex.EncodeToString(crypto.Keccak256([]byte(instrumentCombined))))
	instrumentMatch := instrumentRawHex == info.destTokenHex ||
		instrumentKeccakHex == info.destTokenHex ||
		strings.HasSuffix(instrumentRawHex, info.destTokenHex) ||
		strings.HasSuffix(info.destTokenHex, instrumentRawHex) ||
		strings.HasSuffix(instrumentKeccakHex, info.destTokenHex) ||
		strings.HasSuffix(info.destTokenHex, instrumentKeccakHex)

	remotePools := extractRemotePools(chainPoolCfgAny)
	minBlockDepth := extractMinBlockDepth(chainPoolCfgAny)

	remotePoolMatch := manualExecuteRemotePoolMatches(info, chainPoolCfgAny, remotePools)
	finalityCompatible := (message.Finality == 0 && minBlockDepth == 0) ||
		(message.Finality > 0 && minBlockDepth > 0 && minBlockDepth <= int64(message.Finality))

	return instrumentMatch, remotePoolMatch, finalityCompatible, fmt.Sprintf(
		"instrumentRaw=%s instrumentKeccak=%s minBlockDepth=%d remotePools=%v",
		instrumentRawHex,
		instrumentKeccakHex,
		minBlockDepth,
		remotePools,
	)
}

func manualExecuteRemotePoolMatches(info manualExecuteSelectorInfo, chainPoolCfgAny any, remotePools []string) bool {
	if len(remotePools) == 0 {
		if info.requireInstrumentMatch {
			return false
		}
		cfgText := strings.ToLower(fmt.Sprint(chainPoolCfgAny))
		return strings.Contains(cfgText, info.sourcePoolHex) || strings.Contains(cfgText, info.sourcePoolHexTail40)
	}
	for _, remotePool := range remotePools {
		remotePoolHex := strings.ToLower(strings.TrimPrefix(remotePool, "0x"))
		if remotePoolHex == info.sourcePoolHex ||
			remotePoolHex == info.sourcePoolHexTail40 ||
			strings.HasSuffix(remotePoolHex, info.sourcePoolHex) ||
			strings.HasSuffix(info.sourcePoolHex, remotePoolHex) ||
			strings.HasSuffix(remotePoolHex, info.sourcePoolHexTail40) ||
			strings.HasSuffix(info.sourcePoolHexTail40, remotePoolHex) {
			return true
		}
	}

	return false
}

func collectUnlockedPoolHoldings(
	ctx context.Context,
	participant canton.Participant,
	transferSenderParty string,
	instrumentAdmin string,
	instrumentID string,
) ([]string, []*apiv2.DisclosedContract, error) {
	holdings, err := listHoldingContracts(ctx, participant)
	if err != nil {
		return nil, nil, fmt.Errorf("list pool holdings: %w", err)
	}
	poolHoldings := make([]string, 0, len(holdings))
	poolHoldingDisclosures := make([]*apiv2.DisclosedContract, 0, len(holdings))
	for _, holding := range holdings {
		created := holding.GetCreatedEvent()
		if created == nil || !holdingMatchesPoolInstrument(created, transferSenderParty, instrumentAdmin, instrumentID) {
			continue
		}
		poolHoldings = append(poolHoldings, created.GetContractId())
		poolHoldingDisclosures = append(poolHoldingDisclosures, &apiv2.DisclosedContract{
			TemplateId:       created.GetTemplateId(),
			ContractId:       created.GetContractId(),
			CreatedEventBlob: created.GetCreatedEventBlob(),
			SynchronizerId:   holding.GetSynchronizerId(),
		})
	}

	return poolHoldings, poolHoldingDisclosures, nil
}

func holdingMatchesPoolInstrument(created *apiv2.CreatedEvent, transferSenderParty, instrumentAdmin, instrumentID string) bool {
	views := created.GetInterfaceViews()
	if len(views) == 0 || views[0].GetViewValue() == nil {
		return false
	}
	fields := views[0].GetViewValue().GetFields()
	if len(fields) < 4 || fields[0].GetValue().GetParty() != transferSenderParty {
		return false
	}
	if fields[3].GetValue().GetOptional().GetValue() != nil {
		return false
	}
	instrumentRecord := fields[1].GetValue().GetRecord()
	if instrumentRecord == nil || len(instrumentRecord.GetFields()) < 2 {
		return false
	}
	var holdingInstrumentAdmin, holdingInstrumentID string
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

	return holdingInstrumentAdmin == instrumentAdmin && holdingInstrumentID == instrumentID
}

func buildManualExecuteTokenTransferValue(
	selectedPoolContractID string,
	tokenReceiverParty string,
	transferFactoryCID string,
	transferFactoryCtx *apiv2.Value,
	poolHoldings []string,
	poolExtraContext *apiv2.Value,
) *apiv2.Value {
	tokenPoolHoldingValues := make([]*apiv2.Value, 0, len(poolHoldings))
	for _, cid := range poolHoldings {
		tokenPoolHoldingValues = append(tokenPoolHoldingValues, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}})
	}
	emptyMetadata := emptyMetadataValue()

	return &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
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

func extractMinBlockDepth(chainPoolCfgAny any) int64 {
	switch cfg := chainPoolCfgAny.(type) {
	case lockreleasetokenpool.ChainPoolConfig:
		return int64(cfg.MinBlockDepth)
	case map[string]any:
		m := cfg
		if data, ok := cfg["data"].(map[string]any); ok {
			m = data
		}
		switch raw := m["minBlockDepth"].(type) {
		case types.INT64:
			return int64(raw)
		case int64:
			return raw
		case int:
			return int64(raw)
		case float64:
			return int64(raw)
		}
	}

	return 0
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

func resolveRateLimiterForManualExecute(
	ctx context.Context,
	participant canton.Participant,
	selectedPool *lockreleasetokenpool.LockReleaseTokenPool,
	sourceSelectorKey string,
	sourceSelectorNumericKey string,
	selectedRateLimiters types.GENMAP,
	resolveActiveContractIDByAddress func(templateID string, address contracts.InstanceAddress) (string, error),
) (string, *apiv2.DisclosedContract, error) {
	configuredInboundErrs := make([]string, 0, 2)
	selectorCandidates := []string{sourceSelectorKey}
	if sourceSelectorNumericKey != "" && sourceSelectorNumericKey != sourceSelectorKey {
		selectorCandidates = append(selectorCandidates, sourceSelectorNumericKey)
	}
	for _, selectorCandidate := range selectorCandidates {
		inboundRateLimiterRaw, ok := selectedRateLimiters[selectorCandidate]
		if !ok {
			continue
		}
		cid, disclosure, err := resolveRateLimiterFromRawAddress(
			ctx,
			participant,
			inboundRateLimiterRaw,
			selectedPool,
			resolveActiveContractIDByAddress,
		)
		if err == nil {
			return cid, disclosure, nil
		}
		configuredInboundErrs = append(configuredInboundErrs, fmt.Sprintf("%s: %v", selectorCandidate, err))
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
		return "", nil, fmt.Errorf(
			"resolve inbound rate limiter by address %s for pool %s@%s: %w",
			rateLimiterRawAddr.InstanceAddress().String(),
			selectedPool.InstanceId,
			selectedPool.PoolOwner,
			err,
		)
	}
	disclosure, err := getDisclosedContractByID(ctx, participant, rateLimiterCID)
	if err != nil {
		return "", nil, fmt.Errorf("get disclosed inbound rate limiter %s: %w", rateLimiterCID, err)
	}

	return rateLimiterCID, disclosure, nil
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
