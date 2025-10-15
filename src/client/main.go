package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/canton/admin/participant/v30"
)

func getJWT() (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   "ledger-api-user",
		Audience:  []string{"https://canton.network.global"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	})
	return t.SignedString([]byte("unsafe"))
}

func main() {
	jwtToken, err := getJWT()
	if err != nil {
		panic(err)
	}
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwtToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ledgerApiClient, err := grpc.NewClient("participant2.grpc-ledger-api.localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	adminApiClient, err := grpc.NewClient("participant2.admin-api.localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	vClient := apiv2.NewVersionServiceClient(ledgerApiClient)
	versionResp, err := vClient.GetLedgerApiVersion(ctx, &apiv2.GetLedgerApiVersionRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Ledger API version: ", versionResp.Version)

	packageClient := apiv2.NewPackageServiceClient(ledgerApiClient)
	packageResp, err := packageClient.ListPackages(ctx, &apiv2.ListPackagesRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Available Package IDs: ", packageResp.GetPackageIds())

	adminPackageServiceClient := participantv30.NewPackageServiceClient(adminApiClient)
	listPackagesResponse, err := adminPackageServiceClient.ListPackages(ctx, &participantv30.ListPackagesRequest{})
	if err != nil {
		panic(err)
	}

	for i, description := range listPackagesResponse.GetPackageDescriptions() {
		fmt.Printf("Package Description %v: %v\n", i, description)
	}

	// Upload and vet DAR
	dar, err := os.ReadFile("../../coin/.daml/dist/coin-0.0.1.dar")
	if err != nil {
		panic(err)
	}
	uploadResp, err := adminPackageServiceClient.UploadDar(ctx, &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{
				Bytes:                 dar,
				Description:           nil,
				ExpectedMainPackageId: nil,
			},
		},
		VetAllPackages:     true,
		SynchronizeVetting: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Uploaded DARs: %v", uploadResp.GetDarIds())

	partyManagementClient := admin.NewPartyManagementServiceClient(ledgerApiClient)

	partyId := "testparty1::12206d4c777e200933120adee2d41eab669a12640ff94501e614b6d2daca35bab374"

	found := false
	knownParties, err := partyManagementClient.ListKnownParties(ctx, &admin.ListKnownPartiesRequest{})
	if err != nil {
		panic(err)
	}
	for i, details := range knownParties.GetPartyDetails() {
		fmt.Printf("Party Details %v: %v\n", i, details)
		if details.Party == partyId {
			found = true
		}
	}

	if !found {
		resp, err := partyManagementClient.AllocateParty(ctx, &admin.AllocatePartyRequest{
			PartyIdHint:        "testparty1",
			LocalMetadata:      nil,
			IdentityProviderId: "",
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Allocated new party: %v\n", resp.PartyDetails.Party)
		partyId = resp.PartyDetails.Party
	}

	fmt.Printf("Using party: %v\n", partyId)

	userManagementServiceClient := admin.NewUserManagementServiceClient(ledgerApiClient)
	userList, err := userManagementServiceClient.ListUsers(ctx, &admin.ListUsersRequest{})
	if err != nil {
		panic(err)
	}
	for i, user := range userList.GetUsers() {
		fmt.Printf("User %v: %v\n", i, user)
	}

	grantRightsResult, err := userManagementServiceClient.GrantUserRights(ctx, &admin.GrantUserRightsRequest{
		UserId: "ledger-api-user",
		Rights: []*admin.Right{
			{Kind: &admin.Right_CanExecuteAsAnyParty_{CanExecuteAsAnyParty: &admin.Right_CanExecuteAsAnyParty{}}},
			{Kind: &admin.Right_CanReadAsAnyParty_{CanReadAsAnyParty: &admin.Right_CanReadAsAnyParty{}}},
			{Kind: &admin.Right_CanActAs_{CanActAs: &admin.Right_CanActAs{Party: partyId}}},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Granted rights for user ledger-api-user: %v\n", grantRightsResult.NewlyGrantedRights)

	listRights, err := userManagementServiceClient.ListUserRights(ctx, &admin.ListUserRightsRequest{
		UserId: "ledger-api-user",
	})
	if err != nil {
		panic(err)
	}
	for i, right := range listRights.GetRights() {
		fmt.Printf("Right %v: %v\n", i, right)
	}

	commandServiceClient := apiv2.NewCommandServiceClient(ledgerApiClient)

	res, err := commandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "CoinRegistry",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "issuer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyId}},
								}, {
									Label: "instrumentId",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "admin",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyId}},
										}, {
											Label: "id",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "LINK"}},
										},
									}}}},
								}, {
									Label: "meta",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "values",
											Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
										},
									}}}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyId},
		},
	})
	if err != nil {
		panic(err)
	}

	registryCid := ""
	fmt.Println("Transaction result:")
	fmt.Printf("%v\n", res)
	fmt.Println("\t Events:")
	for i, event := range res.GetTransaction().GetEvents() {
		fmt.Printf("\tEvent %v: %v\n", i, event.GetEvent())
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			registryCid = e.Created.ContractId
		}
	}
	fmt.Println()
	fmt.Printf("Registry CID: %v\n", registryCid)

}
