package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
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

func TestCCIPSend(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	fmt.Println("JWT: ", jwToken)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	userParticipant, err := NewParticipant("participant2.admin-api.localhost:8080", "participant2.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

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
				"amount":   100,
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
	for i, event := range response.GetTransaction().GetEvents() {
		fmt.Printf("Event %v: %v\n", i, event)
	}

	fmt.Println("All active contract:")
	ac, err := GetActiveContracts(ctx, userParticipant)
	require.NoError(t, err)
	for i, contract := range ac {
		fmt.Printf("Active Contract %v: %v\n", i, contract)
	}

	fmt.Println("Active contracts for user ", partyUser)
	activeContracts, err := GetActiveContractsForParty(ctx, userParticipant, partyUser)
	require.NoError(t, err)
	for i, contract := range activeContracts {
		fmt.Printf("Active contract %v: %v\n", i, contract)
	}

	fmt.Println("Active interfaces for user ", partyUser)
	activeContracts2, err := GetActiveContractsForPartyI(ctx, userParticipant, partyUser)
	require.NoError(t, err)
	for i, contract := range activeContracts2 {
		fmt.Printf("Active contract %v: %v\n", i, contract)
	}
}
