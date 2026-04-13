package devenv

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
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
		depth := types.INT64(finality)
		return common.FinalityConfig{BlockDepth: &depth}, nil
	}
}

func (c *Chain) resolveTransferFactoryForExecute(
	ctx context.Context,
	participant canton.Participant,
	executingParty string,
	message protocol.Message,
) (types.CONTRACT_ID, []*apiv2.DisclosedContract, splice_api_token_metadata_v1.ChoiceContext, error) {
	registryAdmin, err := testhelpers.ResolveRegistryAdmin(ctx, participant)
	if err != nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("resolve registry admin: %w", err)
	}
	if message.TokenTransfer == nil || message.TokenTransfer.Amount == nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("token transfer amount is required")
	}
	_, _, transferClient, err := testhelpers.NewValidatorAPIClients(participant)
	if err != nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("create transfer instruction client: %w", err)
	}
	transferFactoryCIDRaw, transferFactoryDisclosures, transferFactoryChoiceContextValue, err := testhelpers.GetTransferFactory(
		ctx,
		transferClient,
		registryAdmin,
		executingParty,
		executingParty,
	)
	if err != nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("resolve transfer factory from scan-proxy: %w", err)
	}
	transferFactoryCID := types.CONTRACT_ID(transferFactoryCIDRaw)
	return transferFactoryCID, transferFactoryDisclosures, splice_api_token_metadata_v1.ChoiceContext{
		Values: testhelpers.ExtractChoiceContextValues(transferFactoryChoiceContextValue),
	}, nil
}

// DeployPerPartyRouter uses the PerPartyRouterFactory to create a new PerPartyRouter instance for the given party.
// It returns the address of the newly created PerPartyRouter instance. If a router already exists for the party, it returns the existing router's address.
func (c *Chain) DeployPerPartyRouter(ctx context.Context, participant canton.Participant, partyId string) (routerAddress contracts.InstanceAddress, disclosedRouter *apiv2.DisclosedContract, err error) {
	// Create PerPartyRouter (ignore error if it exists already)
	cantonPerPartyRouterFactory, err := c.getEDSPerPartyRouterFactory()
	if err != nil {
		return contracts.InstanceAddress{}, nil, fmt.Errorf("failed to get canton per party router factory address: %w", err)
	}
	c.logger.Debug().Str("CantonPerPartyRouterFactory", cantonPerPartyRouterFactory.String()).Msg("Resolved per-party router factory address")

	// Fixed instance ID for the router, this makes the InstanceAddress deterministic.
	routerInstanceID := contracts.InstanceID("test-router")
	// Ignore errors, since the router might already exist if this function is called multiple times for the same party. In that case we just want to return the existing router's address.
	_, _ = operations.ExecuteOperation(c.e.OperationsBundle, per_party_router_factory.CreateRouter, c.chain, contract.ChoiceInput[perpartyrouter.CreateRouter]{
		InstanceAddress: cantonPerPartyRouterFactory,
		Args: perpartyrouter.CreateRouter{
			PartyOwner: types.PARTY(partyId),
			InstanceId: types.TEXT(routerInstanceID.String()),
		},
	})
	routerAddress = routerInstanceID.RawInstanceAddress(types.PARTY(partyId)).InstanceAddress()

	// Get disclosure and contract id of the router
	disclosedRouter, err = testhelpers.GetDisclosedContractByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	})
	if err != nil {
		return contracts.InstanceAddress{}, nil, fmt.Errorf("failed to get disclosed router: %w", err)
	}

	return routerAddress, disclosedRouter, nil
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
	routerAddress, _, err := c.DeployPerPartyRouter(ctx, participant, executingParty)
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
	executeDisclosures, err := c.getExecuteDisclosures(ctx, participant, message, ccvs)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get disclosures for execution: %w", err)
	}
	disclosedContracts := executeDisclosures.DisclosedContracts
	tokenTransferArg := &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}
	if message.TokenTransfer != nil {
		transferFactoryCID, transferFactoryDisclosures, transferFactoryChoiceContext, tfErr := c.resolveTransferFactoryForExecute(ctx, participant, executingParty, message)
		if tfErr != nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("resolve transfer factory for execute: %w", tfErr)
		}
		disclosedContracts = append(disclosedContracts, transferFactoryDisclosures...)
		if executeDisclosures.TokenPoolContractID == nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("missing token pool disclosure from EDS")
		}
		if executeDisclosures.PoolExtraContext == nil {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("missing pool extra context from EDS")
		}
		if len(executeDisclosures.TokenPoolHoldingsContractIDs) == 0 {
			return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("missing token pool holdings from EDS")
		}

		tokenTransferRecord := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
			{Label: "tokenPoolCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: executeDisclosures.TokenPoolContractID.ContractId}}},
			{Label: "tokenReceiverParty", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: executingParty}}},
			{Label: "tokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: string(transferFactoryCID)}}},
				{Label: "extraArgs", Value: ledger.MapToValue(splice_api_token_metadata_v1.ExtraArgs{
					Context: transferFactoryChoiceContext,
					Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
				})},
				{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: executeDisclosures.TokenPoolHoldingsContractIDs}}}},
			}}}}},
			{Label: "poolExtraContext", Value: executeDisclosures.PoolExtraContext},
		}}}}
		tokenTransferArg = &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: tokenTransferRecord}}
	}
	disclosedContracts = testhelpers.DeduplicateDisclosedContracts(disclosedContracts...)

	// Resolve all necessary contracts
	routerCid, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, perpartyrouter.PerPartyRouter{}.GetTemplateID(), routerAddress)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get router contract ID: %w", err)
	}
	c.logger.Debug().Str("InstanceAddress", routerAddress.String()).Str("ContractId", routerCid).Msg("Resolved PerPartyRouter contract")

	receiverCid, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, ccipreceiver.CCIPReceiver{}.GetTemplateID(), receiverAddress)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to get receiver contract ID: %w", err)
	}
	c.logger.Debug().Str("InstanceAddress", receiverAddress.String()).Str("ContractId", receiverCid).Msg("Resolved CCIPReceiver contract")

	// Execute message
	c.logger.Debug().
		Str("EncodedMessage", hex.EncodeToString(encodedMessage)).
		Str("VerifierResults", hex.EncodeToString(verifierResults[0])).
		Str("Receiver", hex.EncodeToString(message.Receiver)).
		Msg("Executing message...")

	emptyCCIPCtx := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
	}}}}
	ccvElements := make([]*apiv2.Value, len(verifiers))
	for i, ccvCid := range executeDisclosures.CCVContractIDs {
		ccvElements[i] = &apiv2.Value{
			Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "ccvCid", Value: &apiv2.Value{Sum: ccvCid}},
				{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(verifierResults[i])}}},
				{Label: "ccvExtraContext", Value: emptyCCIPCtx},
			}}},
		}
	}
	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId: receiverCid,
					Choice:     "Execute",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "context", Value: executeDisclosures.ChoiceContext},
						{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCid}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: encodedMessageHex}}},
						{Label: "tokenTransfer", Value: &apiv2.Value{Sum: tokenTransferArg}},
						{Label: "ccvInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: ccvElements}}}},
					}}}},
				}},
			}},
			ActAs:              []string{executingParty},
			DisclosedContracts: disclosedContracts,
		},
	})
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("failed to execute message: %w", err)
	}
	c.logger.Debug().Str("UpdateID", res.GetTransaction().GetUpdateId()).Msg("Executed message")
	if message.TokenTransfer != nil {
		pendingTransferInstructionCID := ""
		for _, event := range res.GetTransaction().GetEvents() {
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

func (c *Chain) getExecuteDisclosures(
	ctx context.Context,
	participant canton.Participant,
	message protocol.Message,
	verifiers []contracts.InstanceAddress,
) (*testhelpers.ExecuteDisclosures, error) {
	offRampDisclosure, err := testhelpers.GetDisclosedContractByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp",
	})
	if err != nil {
		return nil, fmt.Errorf("get offramp disclosure: %w", err)
	}
	globalConfigDisclosure, err := testhelpers.GetDisclosedContractByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig",
	})
	if err != nil {
		return nil, fmt.Errorf("get global config disclosure: %w", err)
	}
	tokenAdminRegistryDisclosure, err := testhelpers.GetDisclosedContractByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	if err != nil {
		return nil, fmt.Errorf("get token admin registry disclosure: %w", err)
	}
	rmnRemoteDisclosure, err := testhelpers.GetDisclosedContractByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote",
	})
	if err != nil {
		return nil, fmt.Errorf("get rmn remote disclosure: %w", err)
	}

	disclosedContracts := []*apiv2.DisclosedContract{
		offRampDisclosure,
		globalConfigDisclosure,
		tokenAdminRegistryDisclosure,
		rmnRemoteDisclosure,
	}
	choiceContext := ledger.MapToValue(common.CCIPContext{
		Values: types.TEXTMAP{
			"off-ramp":             common.AnyValue{AVContractId: contractIDValue(offRampDisclosure.ContractId)},
			"global-config":        common.AnyValue{AVContractId: contractIDValue(globalConfigDisclosure.ContractId)},
			"token-admin-registry": common.AnyValue{AVContractId: contractIDValue(tokenAdminRegistryDisclosure.ContractId)},
			"rmn-remote":           common.AnyValue{AVContractId: contractIDValue(rmnRemoteDisclosure.ContractId)},
		},
	})

	ccvContractIDs := make([]*apiv2.Value_ContractId, len(verifiers))
	for i, verifier := range verifiers {
		activeCCV, err := contract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			participant.PartyID,
			ccvs.CommitteeVerifier{}.GetTemplateID(),
			verifier,
		)
		if err != nil {
			return nil, fmt.Errorf("find active committee verifier %s: %w", verifier.String(), err)
		}
		disclosedCCV := convertToDisclosedContract(activeCCV)
		disclosedContracts = append(disclosedContracts, disclosedCCV)
		ccvContractIDs[i] = &apiv2.Value_ContractId{ContractId: disclosedCCV.ContractId}
	}

	out := &testhelpers.ExecuteDisclosures{
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		CCVContractIDs:     ccvContractIDs,
	}
	if message.TokenTransfer == nil || len(message.TokenTransfer.DestTokenAddress) == 0 {
		out.DisclosedContracts = testhelpers.DeduplicateDisclosedContracts(out.DisclosedContracts...)
		return out, nil
	}

	tokenAdminRegistryRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		c.chainDetails.ChainSelector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	))
	if err != nil {
		return nil, fmt.Errorf("resolve token admin registry ref: %w", err)
	}
	activeTokenAdminRegistry, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		tokenadminregistry.TokenAdminRegistry{}.GetTemplateID(),
		contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address),
	)
	if err != nil {
		return nil, fmt.Errorf("find active token admin registry: %w", err)
	}
	parsedTokenAdminRegistry, err := bindings.UnmarshalCreatedEvent[tokenadminregistry.TokenAdminRegistry](activeTokenAdminRegistry.GetCreatedEvent())
	if err != nil {
		return nil, fmt.Errorf("parse token admin registry created event: %w", err)
	}
	tokenPoolAddress, instrumentID, err := tokenPoolAddressAndInstrumentIDFromRegistry(parsedTokenAdminRegistry, message.TokenTransfer.DestTokenAddress)
	if err != nil {
		return nil, fmt.Errorf("resolve token pool from registry: %w", err)
	}

	activePool, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		tokenPoolAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("find active token pool %s: %w", tokenPoolAddress.String(), err)
	}
	disclosedPool := convertToDisclosedContract(activePool)
	parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return nil, fmt.Errorf("parse token pool created event: %w", err)
	}

	inboundRateLimiterAddress, err := rateLimiterInstanceAddressFromRemoteConfig(parsedPool.RemoteChainConfigs, uint64(message.SourceChainSelector), "inboundRateLimiter")
	if err != nil {
		return nil, fmt.Errorf("resolve inbound rate limiter: %w", err)
	}
	outboundRateLimiterAddress, err := rateLimiterInstanceAddressFromRemoteConfig(parsedPool.RemoteChainConfigs, uint64(message.SourceChainSelector), "outboundRateLimiter")
	if err != nil {
		return nil, fmt.Errorf("resolve outbound rate limiter: %w", err)
	}
	customInboundRateLimiterAddress, err := optionalRateLimiterInstanceAddressFromRemoteConfig(parsedPool.RemoteChainConfigs, uint64(message.SourceChainSelector), "inboundCustomBlockConfirmationsRateLimiter")
	if err != nil {
		return nil, fmt.Errorf("resolve custom inbound rate limiter: %w", err)
	}

	activeInboundRateLimiter, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		common.RateLimiter{}.GetTemplateID(),
		inboundRateLimiterAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("find active inbound rate limiter %s: %w", inboundRateLimiterAddress.String(), err)
	}
	activeOutboundRateLimiter, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		common.RateLimiter{}.GetTemplateID(),
		outboundRateLimiterAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("find active outbound rate limiter %s: %w", outboundRateLimiterAddress.String(), err)
	}

	disclosedContracts = append(disclosedContracts, disclosedPool, convertToDisclosedContract(activeInboundRateLimiter), convertToDisclosedContract(activeOutboundRateLimiter))
	poolExtraValues := types.TEXTMAP{
		"inbound-rate-limiter": common.AnyValue{AVContractId: contractIDValue(activeInboundRateLimiter.GetCreatedEvent().GetContractId())},
	}
	if customInboundRateLimiterAddress != (contracts.InstanceAddress{}) {
		activeCustomInboundRateLimiter, err := contract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			participant.PartyID,
			common.RateLimiter{}.GetTemplateID(),
			customInboundRateLimiterAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("find active custom inbound rate limiter %s: %w", customInboundRateLimiterAddress.String(), err)
		}
		disclosedContracts = append(disclosedContracts, convertToDisclosedContract(activeCustomInboundRateLimiter))
		poolExtraValues["inbound-custom-block-confirmations-rate-limiter"] = common.AnyValue{
			AVContractId: contractIDValue(activeCustomInboundRateLimiter.GetCreatedEvent().GetContractId()),
		}
	}

	poolHoldingDisclosure, err := findUnlockedHoldingDisclosure(ctx, participant, parsedPool.PoolOwner, instrumentID)
	if err != nil {
		return nil, fmt.Errorf("find token pool holding disclosure: %w", err)
	}
	disclosedContracts = append(disclosedContracts, poolHoldingDisclosure)

	out.DisclosedContracts = testhelpers.DeduplicateDisclosedContracts(disclosedContracts...)
	out.PoolExtraContext = ledger.MapToValue(common.CCIPContext{Values: poolExtraValues})
	out.TokenPoolContractID = &apiv2.Value_ContractId{ContractId: disclosedPool.ContractId}
	out.TokenPoolHoldingsContractIDs = []*apiv2.Value{
		{Sum: &apiv2.Value_ContractId{ContractId: poolHoldingDisclosure.ContractId}},
	}

	return out, nil
}

func contractIDValue(contractID string) *types.CONTRACT_ID {
	if contractID == "" {
		return nil
	}
	cid := types.CONTRACT_ID(contractID)
	return &cid
}

func tokenPoolAddressAndInstrumentIDFromRegistry(
	tokenAdminRegistryCreatedEvent *tokenadminregistry.TokenAdminRegistry,
	destTokenAddress []byte,
) (contracts.InstanceAddress, splice_api_token_holding_v1.InstrumentId, error) {
	expectedHashedInstrumentID := contracts.BytesToInstanceAddress(destTokenAddress)
	tokenConfigValue, ok := tokenAdminRegistryCreatedEvent.TokenConfigs[strings.TrimPrefix(expectedHashedInstrumentID.Hex(), "0x")]
	if !ok {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("missing token config for hashed instrument id %s", expectedHashedInstrumentID.Hex())
	}
	tokenConfigMap, ok := tokenConfigValue.(map[string]any)
	if !ok {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("unexpected token config type %T", tokenConfigValue)
	}

	var tokenConfig tokenadminregistry.TokenConfig
	if err := ledger.MapToStruct(tokenConfigMap, &tokenConfig); err != nil {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("parse token config: %w", err)
	}

	return contracts.InstanceID(tokenConfig.TokenPool.PoolInstanceId).
		RawInstanceAddress(tokenConfig.TokenPool.PoolOwner).
		InstanceAddress(), tokenConfig.InstrumentId, nil
}

func rateLimiterInstanceAddressFromRemoteConfig(
	remoteChainConfigs map[string]any,
	remoteSelector uint64,
	field string,
) (contracts.InstanceAddress, error) {
	remoteConfig, ok := remoteChainConfigs[fmt.Sprintf("%d.", remoteSelector)]
	if !ok {
		return contracts.InstanceAddress{}, fmt.Errorf("missing remote chain config for selector %d", remoteSelector)
	}

	switch cfg := remoteConfig.(type) {
	case lockreleasetokenpool.RemoteChainConfig:
		return rateLimiterInstanceAddressFromRawField(rawRateLimiterField(cfg, field), field)
	case map[string]any:
		rawField, ok := cfg[field]
		if !ok {
			return contracts.InstanceAddress{}, fmt.Errorf("missing %s in remote chain config for %d", field, remoteSelector)
		}
		rawMap, ok := rawField.(map[string]any)
		if !ok {
			return contracts.InstanceAddress{}, fmt.Errorf("unexpected %s type %T", field, rawField)
		}
		unpack, ok := rawMap["unpack"].(string)
		if !ok || unpack == "" {
			return contracts.InstanceAddress{}, fmt.Errorf("missing %s address in remote chain config for %d", field, remoteSelector)
		}
		return rateLimiterInstanceAddressFromRawField(unpack, field)
	default:
		return contracts.InstanceAddress{}, fmt.Errorf("unexpected remote chain config type %T", remoteConfig)
	}
}

func optionalRateLimiterInstanceAddressFromRemoteConfig(
	remoteChainConfigs map[string]any,
	remoteSelector uint64,
	field string,
) (contracts.InstanceAddress, error) {
	remoteConfig, ok := remoteChainConfigs[fmt.Sprintf("%d.", remoteSelector)]
	if !ok {
		return contracts.InstanceAddress{}, fmt.Errorf("missing remote chain config for selector %d", remoteSelector)
	}

	switch cfg := remoteConfig.(type) {
	case lockreleasetokenpool.RemoteChainConfig:
		if raw := rawRateLimiterField(cfg, field); raw != "" {
			return rateLimiterInstanceAddressFromRawField(raw, field)
		}
		return contracts.InstanceAddress{}, nil
	case map[string]any:
		rawField, ok := cfg[field]
		if !ok || rawField == nil {
			return contracts.InstanceAddress{}, nil
		}
		rawMap, ok := rawField.(map[string]any)
		if !ok {
			return contracts.InstanceAddress{}, fmt.Errorf("unexpected %s type %T", field, rawField)
		}
		unpack, _ := rawMap["unpack"].(string)
		if unpack == "" {
			return contracts.InstanceAddress{}, nil
		}
		return rateLimiterInstanceAddressFromRawField(unpack, field)
	default:
		return contracts.InstanceAddress{}, fmt.Errorf("unexpected remote chain config type %T", remoteConfig)
	}
}

func rawRateLimiterField(cfg lockreleasetokenpool.RemoteChainConfig, field string) string {
	switch field {
	case "inboundRateLimiter":
		return string(cfg.InboundRateLimiter.Unpack)
	case "inboundCustomBlockConfirmationsRateLimiter":
		return string(cfg.InboundCustomBlockConfirmationsRateLimiter.Unpack)
	case "outboundRateLimiter":
		return string(cfg.OutboundRateLimiter.Unpack)
	default:
		return ""
	}
}

func rateLimiterInstanceAddressFromRawField(raw string, field string) (contracts.InstanceAddress, error) {
	rawAddress, err := contracts.RawInstanceAddressFromString(raw)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("parse %s raw instance address: %w", field, err)
	}
	return rawAddress.InstanceAddress(), nil
}

func findUnlockedHoldingDisclosure(
	ctx context.Context,
	participant canton.Participant,
	owner types.PARTY,
	instrumentID splice_api_token_holding_v1.InstrumentId,
) (*apiv2.DisclosedContract, error) {
	holdingContracts, err := testhelpers.ListActiveContractsByInterfaceId(ctx, participant, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	if err != nil {
		return nil, fmt.Errorf("list active holdings: %w", err)
	}

	for i := len(holdingContracts) - 1; i >= 0; i-- {
		view, err := holdingViewFromActiveContract(holdingContracts[i])
		if err != nil {
			return nil, err
		}
		if view.Owner != owner {
			continue
		}
		if view.InstrumentId.Admin != instrumentID.Admin || view.InstrumentId.Id != instrumentID.Id {
			continue
		}
		if view.Lock != nil {
			continue
		}
		amount, ok := new(big.Rat).SetString(string(view.Amount))
		if !ok || amount.Sign() <= 0 {
			continue
		}
		return convertToDisclosedContract(holdingContracts[i]), nil
	}

	return nil, fmt.Errorf("no unlocked holding found for %s/%s", instrumentID.Admin, instrumentID.Id)
}

func holdingViewFromActiveContract(active *apiv2.ActiveContract) (*splice_api_token_holding_v1.HoldingView, error) {
	for _, interfaceView := range active.GetCreatedEvent().GetInterfaceViews() {
		if interfaceView.GetInterfaceId().GetModuleName() != "Splice.Api.Token.HoldingV1" ||
			interfaceView.GetInterfaceId().GetEntityName() != "Holding" {
			continue
		}
		var view splice_api_token_holding_v1.HoldingView
		if err := ledger.RecordToStruct(interfaceView.GetViewValue(), &view); err != nil {
			return nil, fmt.Errorf("parse holding view: %w", err)
		}
		return &view, nil
	}

	return nil, fmt.Errorf("holding interface view not found")
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
