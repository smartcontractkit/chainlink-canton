package tests

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/canton/admin/participant/v30"
)

const (
	UserName = "ledger-api-user"
	Audience = "https://canton.network.global"
)

func getJWT() (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   UserName,
		Audience:  []string{Audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	})
	return t.SignedString([]byte("unsafe"))
}

type Participant struct {
	// Admin API
	PackageServiceClient participantv30.PackageServiceClient

	// Ledger API
	PartyManagementServiceClient admin.PartyManagementServiceClient
	UserManagementServiceClient  admin.UserManagementServiceClient

	StateServiceClient   apiv2.StateServiceClient
	CommandServiceClient apiv2.CommandServiceClient
	UpdateServiceClient  apiv2.UpdateServiceClient
	VersionServiceClient apiv2.VersionServiceClient
}

func NewParticipant(adminApiURL, ledgerApiURL string) (*Participant, error) {
	adminApiClient, err := grpc.NewClient(adminApiURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	ledgerApiClient, err := grpc.NewClient(ledgerApiURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Participant{
		PackageServiceClient: participantv30.NewPackageServiceClient(adminApiClient),

		PartyManagementServiceClient: admin.NewPartyManagementServiceClient(ledgerApiClient),
		UserManagementServiceClient:  admin.NewUserManagementServiceClient(ledgerApiClient),
		StateServiceClient:           apiv2.NewStateServiceClient(ledgerApiClient),
		CommandServiceClient:         apiv2.NewCommandServiceClient(ledgerApiClient),
		UpdateServiceClient:          apiv2.NewUpdateServiceClient(ledgerApiClient),
		VersionServiceClient:         apiv2.NewVersionServiceClient(ledgerApiClient),
	}, nil
}

func UploadDARstoMultipleParticipants(ctx context.Context, dars [][]byte, participants ...*Participant) ([]string, error) {
	var darData []*participantv30.UploadDarRequest_UploadDarData
	for _, dar := range dars {
		darData = append(darData, &participantv30.UploadDarRequest_UploadDarData{
			Bytes: dar,
		})
	}

	var packageIDs []string
	for i, participant := range participants {
		res, err := participant.PackageServiceClient.UploadDar(ctx, &participantv30.UploadDarRequest{
			Dars:               darData,
			VetAllPackages:     true,
			SynchronizeVetting: true,
		})
		if err != nil {
			return nil, fmt.Errorf("uploadDAR to participant %d failed: %w", i+1, err)
		}
		packageIDs = append(packageIDs, res.GetDarIds()...)
	}
	return packageIDs, nil
}

func EnsurePartyOnMultipleParticipants(ctx context.Context, participants ...*Participant) ([]string, error) {
	partyHints := []string{"alice", "bob", "charlie", "dave", "erin"}
	var partyIDs []string
	for i, participant := range participants {
		knownParties, err := participant.PartyManagementServiceClient.ListKnownParties(ctx, &admin.ListKnownPartiesRequest{})
		if err != nil {
			return nil, fmt.Errorf("listing known parties failed on participant %d: %w", i+1, err)
		}
		for _, details := range knownParties.GetPartyDetails() {
			if details.GetIsLocal() && strings.HasPrefix(details.GetParty(), fmt.Sprintf("%s::", partyHints[i])) {
				fmt.Printf("Found existing local party: %s\n", details.GetParty())
				partyIDs = append(partyIDs, details.GetParty())
				break
			}
		}
		if len(partyIDs) <= i {
			res, err := participant.PartyManagementServiceClient.AllocateParty(ctx, &admin.AllocatePartyRequest{
				PartyIdHint: partyHints[i],
			})
			if err != nil {
				return nil, fmt.Errorf("allocating party %s on participant %d failed: %w", partyHints[i], i+1, err)
			}
			fmt.Printf("Allocated new party on participant %v: %s\n", i+1, res.PartyDetails.Party)
			partyIDs = append(partyIDs, res.PartyDetails.Party)
		}

		// Grant user rights on party
		grantUserRightsResult, err := participant.UserManagementServiceClient.GrantUserRights(ctx, &admin.GrantUserRightsRequest{
			UserId: UserName,
			Rights: []*admin.Right{
				{Kind: &admin.Right_CanExecuteAsAnyParty_{CanExecuteAsAnyParty: &admin.Right_CanExecuteAsAnyParty{}}},
				{Kind: &admin.Right_CanReadAsAnyParty_{CanReadAsAnyParty: &admin.Right_CanReadAsAnyParty{}}},
				{Kind: &admin.Right_CanActAs_{CanActAs: &admin.Right_CanActAs{Party: partyIDs[i]}}},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("grantUserRights failed on participant %d: %w", i+1, err)
		}
		fmt.Printf("Granted user %q rights to act as party %q: %v\n", UserName, partyIDs[i], grantUserRightsResult.GetNewlyGrantedRights())
	}
	return partyIDs, nil
}

func GetCurrentOffset(ctx context.Context, participant *Participant) (int64, error) {
	ledgerEndResp, err := participant.StateServiceClient.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get ledger end: %w", err)
	}
	return ledgerEndResp.GetOffset(), nil
}

func GetActiveContracts(ctx context.Context, participant *Participant) ([]*apiv2.CreatedEvent, error) {
	var events []*apiv2.CreatedEvent
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersForAnyParty: &apiv2.Filters{Cumulative: []*apiv2.CumulativeFilter{
				{
					IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
						WildcardFilter: &apiv2.WildcardFilter{
							IncludeCreatedEventBlob: false,
						},
					},
				},
			}},
			Verbose: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			events = append(events, c.ActiveContract.GetCreatedEvent())
		}
	}
	return events, nil
}

func GetActiveContractsForParty(ctx context.Context, participant *Participant, party string) ([]*apiv2.CreatedEvent, error) {
	var events []*apiv2.CreatedEvent
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			events = append(events, c.ActiveContract.GetCreatedEvent())
		}
	}
	return events, nil
}

func GetActiveContractsForPartyTemplateId(ctx context.Context, participant *Participant, party string, templateId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
								TemplateId:              templateId,
								IncludeCreatedEventBlob: true,
							}},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			activeContracts = append(activeContracts, c.ActiveContract)
		}
	}
	return activeContracts, nil
}

func GetActiveContractsForPartyInterface(ctx context.Context, participant *Participant, party string, interfaceId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &apiv2.InterfaceFilter{
								InterfaceId:             interfaceId,
								IncludeInterfaceView:    true,
								IncludeCreatedEventBlob: true,
							}},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			activeContracts = append(activeContracts, c.ActiveContract)
		}
	}
	return activeContracts, nil
}

func QueryDisclosedContract(ctx context.Context, contractId string, participant *Participant) (*apiv2.DisclosedContract, error) {
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersForAnyParty: &apiv2.Filters{Cumulative: []*apiv2.CumulativeFilter{
				{
					IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
						WildcardFilter: &apiv2.WildcardFilter{
							IncludeCreatedEventBlob: true,
						},
					},
				},
			}},
			Verbose: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			if c.ActiveContract.GetCreatedEvent().ContractId == contractId {
				return &apiv2.DisclosedContract{
					TemplateId:       c.ActiveContract.GetCreatedEvent().GetTemplateId(),
					ContractId:       c.ActiveContract.GetCreatedEvent().GetContractId(),
					CreatedEventBlob: c.ActiveContract.GetCreatedEvent().GetCreatedEventBlob(),
					SynchronizerId:   c.ActiveContract.GetSynchronizerId(),
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to find active contract with id %s", contractId)
}

func TemplateIdFromString(s string) (*apiv2.Identifier, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid template id format: %s", s)
	}
	return &apiv2.Identifier{
		PackageId:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

func TestCoin(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	participant1, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	participant2, err := NewParticipant("participant2.admin-api.localhost:8080", "participant2.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	participant3, err := NewParticipant("participant3.admin-api.localhost:8080", "participant3.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	participant4, err := NewParticipant("participant4.admin-api.localhost:8080", "participant4.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	participant5, err := NewParticipant("participant5.admin-api.localhost:8080", "participant5.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	version, err := participant1.VersionServiceClient.GetLedgerApiVersion(ctx, &apiv2.GetLedgerApiVersionRequest{})
	require.NoError(t, err)
	fmt.Println(version.Version)

	// Upload the DARs to all participants
	coinDar, err := os.ReadFile("../../contracts/coin/.daml/dist/coin-0.0.1.dar")
	require.NoError(t, err)
	spliceDar, err := os.ReadFile("../../contracts/splice/.daml/dist/splice-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{coinDar, spliceDar}, participant1, participant2, participant3, participant4, participant5)
	require.NoError(t, err)
	fmt.Printf("Uploaded coin DARs to all participants: %v\n", packageIDs)

	parties, err := EnsurePartyOnMultipleParticipants(ctx, participant1, participant2, participant3, participant4, participant5)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyAlice := parties[0]
	partyBob := parties[1]
	partyCharlie := parties[2]
	partyDave := parties[3]
	partyErin := parties[4]
	_ = partyAlice
	_ = partyBob
	_ = partyCharlie
	_ = partyDave
	_ = partyErin

	// Alice is the issuer, creating registry contract
	instrumentId := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "admin",
			Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
		}, {
			Label: "id",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "LINK"}},
		},
	}}}}
	res, err := participant1.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "instrumentId",
									Value: instrumentId,
								}, {
									Label: "meta",
									Value: emptyMetadata,
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyAlice},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	registryCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			registryCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed registry, CID: %v\n", registryCid)

	// Query the contract for explicit disclosure
	disclosedRegistry, err := QueryDisclosedContract(ctx, registryCid, participant1)
	require.NoError(t, err)
	fmt.Printf("Queried registry for disclosure: %v\n", disclosedRegistry.GetContractId())

	// Bob creates MintPreapproval
	res, err = participant2.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintPreapproval",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyBob}},
								}, {
									Label: "sender",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyBob},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	bobMintPreapprovalCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			bobMintPreapprovalCid = e.Created.ContractId
		}
	}
	fmt.Printf("Bob created MintPreapproval, CID: %v\n", bobMintPreapprovalCid)

	// Alice mints to Bob
	res, err = participant1.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice",
								ModuleName: "Splice.Api.Token.BurnMintV1",
								EntityName: "BurnMintFactory",
							},
							ContractId: registryCid,
							Choice:     "BurnMintFactory_BurnMint",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "expectedAdmin",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "sender",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "instrumentId",
									Value: instrumentId,
								}, {
									Label: "inputHoldingCids",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}},
								}, {
									Label: "outputs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{
											Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
												{
													Label: "owner",
													Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyBob}},
												}, {
													Label: "amount",
													Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "42.13"}},
												}, {
													Label: "context",
													Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
														{
															Label: "values",
															Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{
																{
																	Key: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "mint-preapproval"}},
																	Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
																		Constructor: "AV_ContractId",
																		Value:       &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: bobMintPreapprovalCid}},
																	}}},
																},
															}}}},
														},
													}}}},
												},
											}}},
										},
									}}}},
								}, {
									Label: "extraActors",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}},
								}, {
									Label: "extraArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "context",
											Value: emptyChoiceContext,
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
			ActAs: []string{partyAlice},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	bobCoinHoldingCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			bobCoinHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Alice minted to Bob, CID: %v\n", bobCoinHoldingCid)

	// Bob transfers part of their holdings to charlie
	res, err = participant2.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice",
								ModuleName: "Splice.Api.Token.TransferInstructionV1",
								EntityName: "TransferFactory",
							},
							ContractId: registryCid,
							Choice:     "TransferFactory_Transfer",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "expectedAdmin",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "transfer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "sender",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyBob}},
										}, {
											Label: "receiver",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCharlie}},
										}, {
											Label: "amount",
											Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "10.00"}},
										}, {
											Label: "instrumentId",
											Value: instrumentId,
										}, {
											Label: "requestedAt",
											Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().UnixMicro()}},
										}, {
											Label: "executeBefore",
											Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().Add(time.Hour * 24).UnixMicro()}},
										}, {
											Label: "inputHoldingCids",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{
													Sum: &apiv2.Value_ContractId{ContractId: bobCoinHoldingCid},
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
											Value: emptyChoiceContext,
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
			ActAs: []string{partyBob},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedRegistry,
			},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	transferInstructionCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			// There are multiple contracts created in this transaction, since the change will be returned to the sender (bob)
			// We're interested in the TransferInstruction, since that's the one being sent to charlie
			if e.Created.GetTemplateId().GetEntityName() == "CoinTransferInstruction" {
				transferInstructionCid = e.Created.ContractId
				break
			}
		}
	}
	fmt.Printf("Bob transferred to Charlie, CID: %v\n", transferInstructionCid)

	// Charlie accepts transfer from Bob
	res, err = participant3.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice",
								ModuleName: "Splice.Api.Token.TransferInstructionV1",
								EntityName: "TransferInstruction",
							},
							ContractId: transferInstructionCid,
							Choice:     "TransferInstruction_Accept",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "extraArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "context",
											Value: emptyChoiceContext,
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
			ActAs: []string{partyCharlie},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedRegistry,
			},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	charlieHoldingCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			charlieHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Charlie accepted transfer, holding CID: %v\n", charlieHoldingCid)

	// Alice grants mint rights to Dave
	res, err = participant1.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintRole",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "issuer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "minter",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyDave}},
								}, {
									Label: "registry",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: registryCid}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyAlice},
		},
	})
	require.NoError(t, err)
	daveMintRoleCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			daveMintRoleCid = e.Created.ContractId
		}
	}
	fmt.Printf("Alice granted MintRole to Dave, CID: %v\n", daveMintRoleCid)

	// Erin grants MintPreapproval
	res, err = participant5.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintPreapproval",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyErin}},
								}, {
									Label: "sender",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyErin},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	erinPreApprovalCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			erinPreApprovalCid = e.Created.ContractId
		}
	}
	fmt.Printf("Erin created MintPreapproval, CID: %v\n", erinPreApprovalCid)
	disclosedErinPreApproval, err := QueryDisclosedContract(ctx, erinPreApprovalCid, participant5)
	require.NoError(t, err)
	fmt.Printf("Queried MintPreapproval for disclosure: %v\n", disclosedErinPreApproval.GetContractId())

	// Dave uses the MintRole to mint to Erin

	// Asynchronously, listen to all updates that have Erin on Participant 5 as a stakeholder
	currentOffset, err := GetCurrentOffset(ctx, participant5)
	require.NoError(t, err)
	updateStream, err := participant5.UpdateServiceClient.GetUpdates(ctx, &apiv2.GetUpdatesRequest{
		BeginExclusive: currentOffset,
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						partyErin: {
							Cumulative: []*apiv2.CumulativeFilter{
								{
									IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{},
								},
							},
						},
					},
				},
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
			},
			IncludeReassignments:  nil,
			IncludeTopologyEvents: nil,
		},
	})
	require.NoError(t, err)
	defer updateStream.CloseSend()
	erinHoldingCid := ""
	go func() {
		for {
			update, err := updateStream.Recv()
			if err == io.EOF {
				return
			}
			require.NoError(t, err)
			fmt.Printf("Received update on Participant 5: %v\n", update.GetTransaction())
			for _, event := range update.GetTransaction().GetEvents() {
				if c, ok := event.GetEvent().(*apiv2.Event_Created); ok {
					erinHoldingCid = c.Created.GetContractId()
				}
			}
		}
	}()

	res, err = participant4.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintRole",
							},
							ContractId: daveMintRoleCid,
							Choice:     "MintRole_Mint",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "instrumentId",
									Value: instrumentId,
								}, {
									Label: "outputs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{
											Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
												{
													Label: "owner",
													Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyErin}},
												}, {
													Label: "amount",
													Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "77.77"}},
												}, {
													Label: "context",
													Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
														{
															Label: "values",
															Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{
																{
																	Key: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "mint-preapproval"}},
																	Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
																		Constructor: "AV_ContractId",
																		Value:       &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: erinPreApprovalCid}},
																	}}},
																},
															}}}},
														},
													}}}},
												},
											}}},
										},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyDave},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedRegistry,
				disclosedErinPreApproval,
			},
		},
	})
	require.NoError(t, err)
	// Since the submitting party's (dave) node (participant4) isn't a stakeholder on the created contract (CoinHolding),
	// it doesn't actually receive the CreatedEvent as part of the transaction output.
	// This check here therefore won't result in a contract ID, despite us being a witness on the created contract.
	erinCoinHoldingCid := ""
	for i, event := range res.GetTransaction().GetEvents() {
		fmt.Printf("\tEvent %v: %v\n", i, event.GetEvent())
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			erinCoinHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Dave minted to Erin, CID: %v\n", erinCoinHoldingCid) // Will be empty
	fmt.Printf("Transaction: %+v\n", res.GetTransaction())

	// Wait for a couple of seconds, to receive the update on participant 5
	time.Sleep(time.Second * 3)
	fmt.Printf("Received CoinHolding creation update on Participant 5, CID: %v\n", erinHoldingCid)
	require.NotEmpty(t, erinHoldingCid)
}

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}

var emptyChoiceContext = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}
