package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/transferInstructionV1"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

func ChoiceContextFromData(choiceContextData map[string]any) (*apiv2.Value, error) {
	values, ok := choiceContextData["values"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no values found in choice context")
	}

	var fields []*apiv2.TextMap_Entry
	for k, v := range values {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value *apiv2.Value
		switch tag {
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: valueString}}
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return nil, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: valueBool}}
		case "AV_Int":
			// JSON numbers come as float64
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(valueFloat)}}
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Text{Text: valueString}}
		default:
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

func GetRegistryAdmin(ctx context.Context, metadataClient *tokenMetadataV1.ClientWithResponses) (string, error) {
	registryInfoResponse, err := metadataClient.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting registry info: %w", err)
	}
	if registryInfoResponse.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d: %v", registryInfoResponse.StatusCode(), registryInfoResponse.Body)
	}
	return registryInfoResponse.JSON200.AdminId, nil
}

func GetTransferFactory(ctx context.Context, transferInstructionClient *transferInstructionV1.ClientWithResponses, registryAdmin, sender, receiver string) (string, []*apiv2.DisclosedContract, *apiv2.Value, error) {
	transferFactoryResponse, err := transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": registryAdmin,
			"transfer": map[string]any{
				"sender":   sender,
				"receiver": receiver,
				"amount":   "100.00",
				"instrumentId": map[string]any{
					"admin": registryAdmin,
					"id":    "Amulet",
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(time.Hour * -1).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				"inputHoldingCids": []string{},
				"meta": map[string]any{
					"values": map[string]any{},
				},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{
					"values": map[string]any{},
				},
				"meta": map[string]any{
					"values": map[string]any{},
				},
			},
		},
		ExcludeDebugFields: nil,
	})
	if err != nil {
		return "", nil, nil, fmt.Errorf("error getting transferFactory response: %w", err)
	}
	if transferFactoryResponse.StatusCode() != http.StatusOK {
		return "", nil, nil, fmt.Errorf("unexpected status code: %d: %v", transferFactoryResponse.StatusCode(), transferFactoryResponse.Body)
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		if err != nil {
			return "", nil, nil, fmt.Errorf("failed to parse template id: %w", err)
		}
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if err != nil {
			return "", nil, nil, fmt.Errorf("failed to decode created event blob: %w", err)
		}
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}
	choiceContext, err := ChoiceContextFromData(transferFactoryResponse.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	factoryCid := transferFactoryResponse.JSON200.FactoryId

	return factoryCid, disclosedContracts, choiceContext, nil
}

func GetAmuletRulesContract(ctx context.Context, scanProxyClient *scanProxy.ClientWithResponses) (string, *apiv2.Identifier, error) {
	// Get Amulet Rules contract
	amuletRulesResponse, err := scanProxyClient.GetAmuletRulesWithResponse(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("error getting amulet rules response: %w", err)
	}
	if amuletRulesResponse.StatusCode() != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d: %v", amuletRulesResponse.StatusCode(), amuletRulesResponse.Body)
	}
	amuletRulesId, err := TemplateIdFromString(amuletRulesResponse.JSON200.AmuletRules.Contract.TemplateId)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse amulet rules template id: %w", err)
	}
	return amuletRulesResponse.JSON200.AmuletRules.Contract.ContractId, amuletRulesId, nil
}

func GetFirstOpenMiningRound(ctx context.Context, scanProxyClient *scanProxy.ClientWithResponses) (string, error) {
	openMiningRoundResponse, err := scanProxyClient.GetOpenAndIssuingMiningRoundsWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting open mining rounds response: %w", err)
	}
	if openMiningRoundResponse.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d: %v", openMiningRoundResponse.StatusCode(), openMiningRoundResponse.Body)
	}
	slices.SortFunc(openMiningRoundResponse.JSON200.OpenMiningRounds, func(a, b scanProxy.ContractWithState) int {
		aOpen, _ := time.Parse(time.RFC3339, a.Contract.Payload["opensAt"].(string))
		bOpen, _ := time.Parse(time.RFC3339, b.Contract.Payload["opensAt"].(string))
		return int(aOpen.UnixMilli() - bOpen.UnixMilli())
	})
	var openMiningRoundCid string
	for _, round := range openMiningRoundResponse.JSON200.OpenMiningRounds {
		opensAt, err := time.Parse(time.RFC3339, round.Contract.Payload["opensAt"].(string))
		if err != nil {
			return "", fmt.Errorf("failed to parse opensAt %q: %w", round.Contract.Payload["opensAt"], err)
		}
		targetClosesAt, err := time.Parse(time.RFC3339, round.Contract.Payload["targetClosesAt"].(string))
		if err != nil {
			return "", fmt.Errorf("failed to parse targetClosesAt %q: %w", round.Contract.Payload["targetClosesAt"], err)
		}
		if opensAt.Before(time.Now()) && targetClosesAt.After(time.Now()) {
			openMiningRoundCid = round.Contract.ContractId
		}
	}
	return openMiningRoundCid, nil
}

func MintAMT(
	ctx context.Context,
	participant *Participant,
	metadataClient *tokenMetadataV1.ClientWithResponses,
	transferInstructionClient *transferInstructionV1.ClientWithResponses,
	scanProxyClient *scanProxy.ClientWithResponses,
	toParty string,
	amount string,
) (string, error) {
	// Get Instrument Admin
	registryAdmin, err := GetRegistryAdmin(ctx, metadataClient)
	if err != nil {
		return "", fmt.Errorf("failed to get registry admin: %w", err)
	}

	// Get AmuletRules Contract
	amuletRulesCid, amuletRulesId, err := GetAmuletRulesContract(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get amulet rules contract: %w", err)
	}

	// Create Transfer Factory
	_, disclosedContracts, _, err := GetTransferFactory(ctx, transferInstructionClient, registryAdmin, registryAdmin, toParty)
	if err != nil {
		return "", fmt.Errorf("failed to get transfer factory: %w", err)
	}

	// Get open mining round
	openMiningRoundCid, err := GetFirstOpenMiningRound(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get open mining round: %w", err)
	}

	// Mint AMT
	response, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: amuletRulesId,
							ContractId: amuletRulesCid,
							Choice:     "AmuletRules_DevNet_Tap",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: toParty}},
								}, {
									Label: "amount",
									Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: amount}},
								}, {
									Label: "openRound",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: openMiningRoundCid}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{toParty},
			DisclosedContracts: disclosedContracts,
		},
		TransactionFormat: nil,
	})
	if err != nil {
		return "", fmt.Errorf("failed to mint AMT: %w", err)
	}

	tokenHoldingCid := ""
	for _, event := range response.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenHoldingCid = e.Created.ContractId
		}
	}
	return tokenHoldingCid, nil
}

func TestCCIPSend(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	userParticipant, err := NewParticipant("participant2.admin-api.localhost:8080", "participant2.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	// Dars
	// Uploading the CCIP dars to all participants
	feeQuoterDar, err := os.ReadFile("../../contracts/ccip/feequoter/.daml/dist/ccip-feequoter-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{feeQuoterDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP FeeQuoter to ccipParticipant: %s\n", packageIDs)
	onRampDar, err := os.ReadFile("../../contracts/ccip/onramp/.daml/dist/ccip-onramp-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{onRampDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP OnRamp to ccipParticipant: %s\n", packageIDs)
	routerDar, err := os.ReadFile("../../contracts/ccip/router/.daml/dist/ccip-router-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{routerDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP Router to ccipParticipant: %s\n", packageIDs)

	parties, err := EnsurePartyOnMultipleParticipants(ctx, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyCCIP := parties[0]
	partyUser := parties[1]
	_ = partyCCIP
	_ = partyUser

	// HTTP Clients
	registryUrl := "http://scan.localhost:8080"
	metadataClient, err := tokenMetadataV1.NewClientWithResponses(registryUrl)
	require.NoError(t, err)
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(registryUrl)
	require.NoError(t, err)
	authProvider, err := securityprovider.NewSecurityProviderBearerToken(jwToken)
	require.NoError(t, err)
	scanProxyClient, err := scanProxy.NewClientWithResponses("http://sv.wallet.localhost:8080/api/validator", scanProxy.WithRequestEditorFn(authProvider.Intercept))
	require.NoError(t, err)

	// Get DSO Admin Party
	registryAdmin, err := GetRegistryAdmin(ctx, metadataClient)
	require.NoError(t, err)

	// Mint Tokens to User Party
	tokenHoldingCid, err := MintAMT(ctx, userParticipant, metadataClient, transferInstructionClient, scanProxyClient, partyUser, "100.00")
	require.NoError(t, err)
	fmt.Printf("Minted 100 AMT, Token Holding Cid: %s\n", tokenHoldingCid)

	// ========================
	// |   CCIP Deployment    |
	// ========================

	// CCIP Party deploys CCIP contracts
	chainSelector := "1111111111"
	destChainSelector1 := "2222222222"
	instrumentIdAmt := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "admin",
			Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: registryAdmin}},
		}, {
			Label: "id",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "Amulet"}},
		},
	}}}}

	// Deploy FeeQuoter
	res, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-feequoter",
								ModuleName: "CCIP.FeeQuoter",
								EntityName: "FeeQuoter",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "feeTokens",
									Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
								}, {
									Label: "destChainConfigs",
									Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
								}, {
									Label: "usdPerUnitGasByDestChainSelector",
									Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
								}, {
									Label: "usdPerToken",
									Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
								}, {
									Label: "priceUpdaters",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{
											Sum: &apiv2.Value_Party{Party: partyCCIP},
										},
									}}}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	feeQuoterCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			feeQuoterCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed FeeQuoter to: %s\n", feeQuoterCid)

	// Deploy OnRamp
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-onramp",
								ModuleName: "CCIP.OnRamp",
								EntityName: "OnRamp",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "staticConfig",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "chainSelector",
											Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: chainSelector}},
										},
									}}}},
								}, {
									Label: "destChainConfigs",
									Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	onRampCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			onRampCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed OnRamp to: %s\n", onRampCid)

	// Deploy Router
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-router",
								ModuleName: "CCIP.Router",
								EntityName: "Router",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	routerCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			routerCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed Router to: %s\n", routerCid)

	// Apply configs

	// Apply DestChainConfig to FeeQuoter
	feeQuoterDestChainConfig1 := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "isEnabled",
			Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}},
		}, {
			Label: "maxDataBytes",
			Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 5000}},
		}, {
			Label: "maxPerMsgGasLimit",
			Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2000000}},
		}, {
			Label: "destGasOverhead",
			Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 30000}},
		}, {
			Label: "destGasPerPayloadByteBase",
			Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 12}},
		}, {
			Label: "chainFamilySelector",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "2812d52c"}},
		}, {
			Label: "defaultTxGasLimit",
			Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 500000}},
		}, {
			Label: "networkFeeUSD",
			Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "1.00"}},
		}, {
			Label: "defaultTokenFeeUSD",
			Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0.50"}},
		}, {
			Label: "defaultTokenDestGasOverhead",
			Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}},
		},
	}}}}
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-feequoter",
								ModuleName: "CCIP.FeeQuoter",
								EntityName: "FeeQuoter",
							},
							ContractId: feeQuoterCid,
							Choice:     "ApplyDestChainConfigUpdates",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "destChainConfigArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										&apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
											{
												Label: "destChainSelector",
												Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector1}},
											}, {
												Label: "destChainConfig",
												Value: feeQuoterDestChainConfig1,
											},
										}}}},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			feeQuoterCid = e.Created.ContractId
		}
	}
	fmt.Printf("Applied DestChainConfigUpdates to FeeQuoter: %v\n", feeQuoterCid)

	// Apply Price Updates to FeeQuoter
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-feequoter",
								ModuleName: "CCIP.FeeQuoter",
								EntityName: "FeeQuoter",
							},
							ContractId: feeQuoterCid,
							Choice:     "UpdatePrices",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "priceUpdates",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "tokenPriceUpdates",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{
													Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
														{
															Label: "instrumentId",
															Value: instrumentIdAmt,
														}, {
															Label: "usdPerToken",
															Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0.50"}},
														},
													}}},
												},
											}}}},
										}, {
											Label: "gasPriceUpdates",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{
													Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
														{
															Label: "destChainSelector",
															Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector1}},
														}, {
															Label: "usdPerUnitGas",
															Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0.01"}},
														},
													}}},
												},
											}}}},
										},
									}}}},
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			feeQuoterCid = e.Created.ContractId
		}
	}
	fmt.Printf("Applied Price Updates to FeeQuoter: %v\n", feeQuoterCid)

	// Apply DestChainConfig to OnRamp
	onRampDestChainConfig1 := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "sequenceNumber",
			Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}},
		}, {
			Label: "defaultExecutor",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}},
		}, {
			Label: "laneMandatedCCVs",
			Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}},
		}, {
			Label: "defaultCCVs",
			Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}},
		},
	}}}}
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-onramp",
								ModuleName: "CCIP.OnRamp",
								EntityName: "OnRamp",
							},
							ContractId: onRampCid,
							Choice:     "ApplyDestChainConfigUpdates",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "destChainConfigArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										&apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
											{
												Label: "destChainSelector",
												Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector1}},
											}, {
												Label: "destChainConfig",
												Value: onRampDestChainConfig1,
											},
										}}}},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			onRampCid = e.Created.ContractId
		}
	}
	fmt.Printf("Applied DestChainConfigUpdates to OnRamp: %v\n", onRampCid)

	// =========================
	// |     Call CCIPSend     |
	// =========================

	ccipApi := CCIPApi{
		Participant: ccipParticipant,
		CCIPParty:   partyCCIP,
	}

	for range 3 {
		CCIPSend(t, ccipApi, userParticipant, partyUser, destChainSelector1, instrumentIdAmt, metadataClient, transferInstructionClient)
	}

	// Check token holdings for the CCIP party
	tokenHoldings, err := GetActiveContractsForPartyInterface(ctx, ccipParticipant, partyCCIP, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	require.NoError(t, err)
	fmt.Printf("CCIP Party %s has %d token holdings after CCIPSend\n", partyCCIP, len(tokenHoldings))
	for _, holding := range tokenHoldings {
		balance := holding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		fmt.Printf(" - Token Holding Cid: %s, Balance: %s\n", holding.GetCreatedEvent().GetContractId(), balance)
	}
}

func CCIPSend(
	t *testing.T,
	ccipApi CCIPApi,
	participant *Participant,
	party string,
	destChainSelector string,
	feeToken *apiv2.Value,

	// Token Standard APIs
	metadataClient *tokenMetadataV1.ClientWithResponses,
	transferInstructionClient *transferInstructionV1.ClientWithResponses,
) {
	// Add authorization to the context
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Query disclosures for CCIP contract - these would be served off-chain
	ccipParty := ccipApi.GetCCIPParty()
	feeQuoterDisclosure, err := ccipApi.GetFeeQuoter(t.Context())
	require.NoError(t, err)
	onRampDisclosure, err := ccipApi.GetOnRamp(t.Context())
	require.NoError(t, err)
	routerDisclosure, err := ccipApi.GetRouter(t.Context())
	require.NoError(t, err)

	// User queries their current TokenHoldings
	tokenHoldings, err := GetActiveContractsForPartyInterface(ctx, participant, party, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	require.NoError(t, err)
	// Use the first token Holding we find
	tokenHolding := tokenHoldings[0]
	tokenHoldingCid := tokenHolding.GetCreatedEvent().GetContractId()
	fmt.Printf("Using token holding cid: %s for CCIPSend, balance: %v\n", tokenHoldingCid, tokenHolding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric)

	// Query TransferFactory from AMT registry
	registryAdmin, err := GetRegistryAdmin(ctx, metadataClient)
	require.NoError(t, err)
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := GetTransferFactory(ctx, transferInstructionClient, registryAdmin, party, ccipParty)
	require.NoError(t, err)

	// Call CCIPSend from User Participant on Router
	response, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-router",
								ModuleName: "CCIP.Router",
								EntityName: "Router",
							},
							ContractId: routerDisclosure.ContractId,
							Choice:     "CCIPSend",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "destChainSelector",
									Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector}},
								}, {
									Label: "message",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "receiver",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0xReceiverAddress"}},
										}, {
											Label: "payload",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0xPayloadData"}},
										}, {
											Label: "feeToken",
											Value: feeToken,
										}, {
											Label: "extraArgs",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0x1234"}},
										},
									}}}},
								}, {
									Label: "feeInput",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "transferFactory",
											Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}},
										},
										{
											Label: "inputHoldingCids",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{
													Sum: &apiv2.Value_ContractId{ContractId: tokenHoldingCid},
												},
											}}}},
										}, {
											Label: "extraArgs",
											Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
												{
													Label: "context",
													Value: choiceContext,
												}, {
													Label: "meta",
													Value: emptyMetadata,
												},
											}}}},
										},
									}}}},
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: party}},
								}, {
									Label: "onRamp",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: onRampDisclosure.ContractId}},
								}, {
									Label: "feeQuoter",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: feeQuoterDisclosure.ContractId}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{party},
			DisclosedContracts: append(transferFactoryDisclosures, feeQuoterDisclosure, onRampDisclosure, routerDisclosure),
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	ccipMessageSentCid := ""
	for _, event := range response.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
				ccipMessageSentCid = e.Created.ContractId
			}
		}
	}
	fmt.Printf("Sent CCIP Message using Router, event: %v\n", ccipMessageSentCid)
}

func TestMessageSentListener(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	onRampDar, err := os.ReadFile("../../contracts/ccip/onramp/.daml/dist/ccip-onramp-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{onRampDar}, ccipParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP OnRamp to ccipParticipant: %s\n", packageIDs)

	parties, err := EnsurePartyOnMultipleParticipants(ctx, ccipParticipant)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyCCIP := parties[0]

	currentOffset, err := GetCurrentOffset(ctx, ccipParticipant)
	require.NoError(t, err)
	updateStream, err := ccipParticipant.UpdateServiceClient.GetUpdates(ctx, &apiv2.GetUpdatesRequest{
		BeginExclusive: currentOffset,
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						partyCCIP: {
							Cumulative: []*apiv2.CumulativeFilter{
								{
									IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
										TemplateId: &apiv2.Identifier{
											PackageId:  "#ccip-onramp",
											ModuleName: "CCIP.OnRamp",
											EntityName: "CCIPMessageSent",
										},
										IncludeCreatedEventBlob: false,
									}},
								},
							},
						},
					},
					Verbose: true,
				},
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
			},
		},
	})
	require.NoError(t, err)
	defer updateStream.CloseSend()
	fmt.Println("Listening for CCIPMessageSent events...")
	go func() {
		for {
			update, err := updateStream.Recv()
			if err == io.EOF {
				return
			}
			require.NoError(t, err)
			// fmt.Printf("Received update on Participant 5: %v\n", update.GetTransaction())
			for _, event := range update.GetTransaction().GetEvents() {
				if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
					if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
						fmt.Println("CCIPMessageSent event received:")
						marshalledEvent := protojson.Format(event)
						fmt.Println(string(marshalledEvent))
						fmt.Println("==============================")
					}
				}
			}
		}
	}()
	<-t.Context().Done()
	fmt.Println("Exiting.")
}
