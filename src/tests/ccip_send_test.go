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

	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/transferInstructionV1"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
	participantv30 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/canton/admin/participant/v30"
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
		valueString := f["value"].(string)

		var value *apiv2.Value
		switch tag { // TODO Add remaining cases
		case "AV_ContractId":
			value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: valueString}}
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
	feeQuoterDar, err := os.ReadFile("../../contracts/ccip/feequoter/.daml/dist/feequoter-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{feeQuoterDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP FeeQuoter to ccipParticipant: %s\n", packageIDs)
	onRampDar, err := os.ReadFile("../../contracts/ccip/onramp/.daml/dist/onramp-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{onRampDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP OnRamp to ccipParticipant: %s\n", packageIDs)

	parties, err := EnsurePartyOnMultipleParticipants(ctx, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyCCIP := parties[0]
	partyUser := parties[1]
	_ = partyCCIP
	_ = partyUser

	// List packages
	pkgs, err := userParticipant.PackageServiceClient.ListPackages(ctx, &participantv30.ListPackagesRequest{
		Limit:      0,
		FilterName: "",
	})
	require.NoError(t, err)
	for i, description := range pkgs.GetPackageDescriptions() {
		fmt.Printf("Package %d: %s\n", i, description)
	}

	// Acquire AMT
	registryUrl := "http://scan.localhost:8080"
	metadataClient, err := tokenMetadataV1.NewClientWithResponses(registryUrl)
	require.NoError(t, err)
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(registryUrl)
	require.NoError(t, err)
	authProvider, err := securityprovider.NewSecurityProviderBearerToken(jwToken)
	require.NoError(t, err)
	scanProxyClient, err := scanProxy.NewClientWithResponses("http://sv.wallet.localhost:8080/api/validator", scanProxy.WithRequestEditorFn(authProvider.Intercept))
	require.NoError(t, err)

	// Get Instrument Admin
	registryInfoResponse, err := metadataClient.GetRegistryInfoWithResponse(t.Context())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, registryInfoResponse.StatusCode())
	registryAdmin := registryInfoResponse.JSON200.AdminId
	fmt.Printf("AMT Admin: %v\n", registryAdmin)

	// Create Transfer Factory
	transferFactoryResponse, err := transferInstructionClient.GetTransferFactoryWithResponse(t.Context(), transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": registryAdmin,
			"transfer": map[string]any{
				"sender":   registryAdmin,
				"receiver": partyUser,
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
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, transferFactoryResponse.StatusCode())
	fmt.Println(string(transferFactoryResponse.Body))

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		require.NoError(t, err)
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		require.NoError(t, err)
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}
	fmt.Println("Disclosed Contracts: ", disclosedContracts)

	// Get Amulet Rules contract
	amuletRulesResponse, err := scanProxyClient.GetAmuletRulesWithResponse(t.Context())
	require.NoError(t, err)
	fmt.Println(string(amuletRulesResponse.Body))
	require.Equal(t, http.StatusOK, amuletRulesResponse.StatusCode())
	amuletRulesId, err := TemplateIdFromString(amuletRulesResponse.JSON200.AmuletRules.Contract.TemplateId)
	require.NoError(t, err)

	// Get open mining round
	openMiningRoundResponse, err := scanProxyClient.GetOpenAndIssuingMiningRoundsWithResponse(t.Context())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, openMiningRoundResponse.StatusCode())
	fmt.Println(string(openMiningRoundResponse.Body))
	slices.SortFunc(openMiningRoundResponse.JSON200.OpenMiningRounds, func(a, b scanProxy.ContractWithState) int {
		aOpen, err := time.Parse(time.RFC3339, a.Contract.Payload["opensAt"].(string))
		require.NoError(t, err)
		bOpen, err := time.Parse(time.RFC3339, b.Contract.Payload["opensAt"].(string))
		require.NoError(t, err)
		return int(aOpen.UnixMilli() - bOpen.UnixMilli())
	})
	var openMiningRoundCid string
	for _, round := range openMiningRoundResponse.JSON200.OpenMiningRounds {
		opensAt, err := time.Parse(time.RFC3339, round.Contract.Payload["opensAt"].(string))
		require.NoError(t, err)
		targetClosesAt, err := time.Parse(time.RFC3339, round.Contract.Payload["targetClosesAt"].(string))
		require.NoError(t, err)
		if opensAt.Before(time.Now()) && targetClosesAt.After(time.Now()) {
			openMiningRoundCid = round.Contract.ContractId
		}
	}
	require.NotZero(t, openMiningRoundCid)
	fmt.Println("Using open mining round: ", openMiningRoundCid)

	response, err := userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: amuletRulesId,
							ContractId: amuletRulesResponse.JSON200.AmuletRules.Contract.ContractId,
							Choice:     "AmuletRules_DevNet_Tap",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "amount",
									Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "100"}},
								}, {
									Label: "openRound",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: openMiningRoundCid}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyUser},
			DisclosedContracts: disclosedContracts,
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	tokenHoldingCid := ""
	for i, event := range response.GetTransaction().GetEvents() {
		fmt.Printf("Event %v: %v\n", i, event)
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Minted AMT to: %s\n", tokenHoldingCid)

	// CCIP Party deploys CCIP contracts
	chainSelector := "1111111111"
	destChainSelector1 := "2222222222"
	destChainSelector2 := "3333333333"
	instrumentIdAmt := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "admin",
			Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: registryAdmin}},
		}, {
			Label: "id",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "Amulet"}},
		},
	}}}}
	_ = destChainSelector1
	_ = destChainSelector2
	_ = instrumentIdAmt

	// FeeQuoter
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

	// OnRamp
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

	// Apply price updates for destChainSelector1 and AMT
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

	// Create TransferInstruction to ccip party
	transferFactoryResponse, err = transferInstructionClient.GetTransferFactoryWithResponse(t.Context(), transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": registryAdmin,
			"transfer": map[string]any{
				"sender":   partyUser,
				"receiver": partyCCIP,
				"amount":   "2.00",
				"instrumentId": map[string]any{
					"admin": registryAdmin,
					"id":    "Amulet",
				},
				"lock":          nil,
				"requestedAt":   time.Now().Add(time.Minute * -1).Format(time.RFC3339),
				"executeBefore": time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				"inputHoldingCids": []string{
					tokenHoldingCid,
				},
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
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, transferFactoryResponse.StatusCode())

	transferFactoryCid := transferFactoryResponse.JSON200.FactoryId
	fmt.Println("TransferFactory Cid: ", transferFactoryCid)

	disclosedContracts = nil
	for _, contract := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		require.NoError(t, err)
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		require.NoError(t, err)
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}
	fmt.Println("Disclosed Contracts: ", disclosedContracts)

	choiceContext, err := ChoiceContextFromData(transferFactoryResponse.JSON200.ChoiceContext.ChoiceContextData)
	require.NoError(t, err)
	response, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice-api-token-transfer-instruction-v1",
								ModuleName: "Splice.Api.Token.TransferInstructionV1",
								EntityName: "TransferFactory",
							},
							ContractId: transferFactoryCid,
							Choice:     "TransferFactory_Transfer",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "expectedAdmin",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: registryAdmin}},
								}, {
									Label: "transfer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "sender",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
										}, {
											Label: "receiver",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
										}, {
											Label: "amount",
											Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "2.00"}},
										}, {
											Label: "instrumentId",
											Value: instrumentIdAmt,
										}, {
											Label: "requestedAt",
											Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().Add(time.Minute * -1).UnixMicro()}},
										}, {
											Label: "executeBefore",
											Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().Add(time.Hour * 24).UnixMicro()}},
										}, {
											Label: "inputHoldingCids",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{
													Sum: &apiv2.Value_ContractId{ContractId: tokenHoldingCid},
												},
											}}}},
										}, {
											Label: "meta",
											Value: emptyMetadata,
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
						},
					},
				},
			},
			ActAs:              []string{partyUser},
			DisclosedContracts: disclosedContracts,
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	transferInstructionCid := ""
	for _, event := range response.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			transferInstructionCid = e.Created.ContractId
		}
	}
	fmt.Printf("Created TransferInstruction to CCIP Party: %s\n", transferInstructionCid)

	// Call CCIPSend on the OnRamp

	// Query disclosures for CCIP contract
	feeQuoterDisclosure, err := QueryDisclosedContract(ctx, feeQuoterCid, ccipParticipant)
	require.NoError(t, err)
	onRampDisclosure, err := QueryDisclosedContract(ctx, onRampCid, ccipParticipant)
	require.NoError(t, err)

	response, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
							Choice:     "CCIPSend",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "destChainSelector",
									Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector1}},
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
											Value: instrumentIdAmt,
										}, {
											Label: "extraArgs",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0x1234"}},
										},
									}}}},
								}, {
									Label: "feeInput",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "transfer",
											Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferInstructionCid}},
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
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "feeQuoter",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: feeQuoterCid}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyUser},
			DisclosedContracts: append(disclosedContracts, feeQuoterDisclosure, onRampDisclosure),
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
	fmt.Printf("Sent CCIP Message, event: %v\n", ccipMessageSentCid)
}

func TestMessageSentListener(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	onRampDar, err := os.ReadFile("../../contracts/ccip/onramp/.daml/dist/onramp-1.0.0.dar")
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
						fmt.Printf("CCIPMessageSent: %+v\n", event.GetCreated())
					}
				}
			}
		}
	}()
	<-t.Context().Done()
	fmt.Println("Exiting.")
}
