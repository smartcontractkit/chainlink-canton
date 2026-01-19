package tests

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/smartcontractkit/chainlink-canton-internal/generated/coin"
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
	slices.SortFunc(events, func(a, b *apiv2.CreatedEvent) int {
		return a.GetCreatedAt().AsTime().Compare(b.GetCreatedAt().AsTime())
	})
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
	slices.SortFunc(events, func(a, b *apiv2.CreatedEvent) int {
		return a.GetCreatedAt().AsTime().Compare(b.GetCreatedAt().AsTime())
	})
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
	slices.SortFunc(activeContracts, func(a, b *apiv2.ActiveContract) int {
		return a.GetCreatedEvent().GetCreatedAt().AsTime().Compare(b.GetCreatedEvent().GetCreatedAt().AsTime())
	})
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
	slices.SortFunc(activeContracts, func(a, b *apiv2.ActiveContract) int {
		return a.GetCreatedEvent().GetCreatedAt().AsTime().Compare(b.GetCreatedEvent().GetCreatedAt().AsTime())
	})
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

func normalizeTemplateID(tid string) string {
	return strings.TrimPrefix(tid, "#")
}

func QueryDisclosedContractWithBindingClient(
	ctx context.Context,
	cl *client.DamlBindingClient,
	party string,
	contractID string,
	templateID string, // now REQUIRED if you want CreatedEventBlob
) (*model.DisclosedContract, error) {

	templateID = normalizeTemplateID(templateID)
	if templateID == "" {
		return nil, fmt.Errorf("templateID must be non-empty to request CreatedEventBlob")
	}

	le, err := cl.StateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetLedgerEnd failed: %w", err)
	}

	req := &model.GetActiveContractsRequest{
		ActiveAtOffset: le.Offset,
		EventFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				party: {
					Inclusive: &model.InclusiveFilters{
						TemplateFilters: []*model.TemplateFilter{
							{
								TemplateID:              templateID,
								IncludeCreatedEventBlob: true,
							},
						},
					},
				},
			},
			Verbose: true,
		},
	}

	respCh, errCh := cl.StateService.GetActiveContracts(ctx, req)

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("ACS timeout while searching for %s: %w", contractID, ctx.Err())

		case err, ok := <-errCh:
			if ok && err != nil {
				return nil, fmt.Errorf("ACS stream error: %w", err)
			}

		case resp, ok := <-respCh:
			if !ok {
				return nil, fmt.Errorf("active contract not found: %s", contractID)
			}
			if resp == nil || resp.ContractEntry == nil {
				continue
			}

			entry, ok := resp.ContractEntry.(*model.ActiveContractEntry)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.CreatedEvent == nil {
				continue
			}

			created := entry.ActiveContract.CreatedEvent
			if created.ContractID != contractID {
				continue
			}
			if len(created.CreatedEventBlob) == 0 {
				return nil, fmt.Errorf("createdEventBlob missing for %s (IncludeCreatedEventBlob didn't take effect)", contractID)
			}

			return &model.DisclosedContract{
				TemplateID:       created.TemplateID,
				ContractID:       created.ContractID,
				CreatedEventBlob: created.CreatedEventBlob,
				SynchronizerID:   entry.ActiveContract.SynchronizerID,
			}, nil
		}
	}
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

func newBindingClient(t *testing.T, token, ledgerAddr, adminAddr string) *client.DamlBindingClient {
	t.Helper()

	bc, err := client.NewDamlClient(token, ledgerAddr).
		WithAdminAddress(adminAddr).
		Build(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() { bc.Close() })
	return bc
}

func normalizeTemplateKey(tid string) string {
	tid = strings.TrimPrefix(tid, "#")
	parts := strings.Split(tid, ":")
	if len(parts) < 3 {
		return tid
	}
	return parts[len(parts)-2] + ":" + parts[len(parts)-1]
}

func ensureCanReadAs(ctx context.Context, cl *client.DamlBindingClient, userID, p string) error {
	rights, err := cl.UserMng.ListUserRights(ctx, userID)
	if err != nil {
		return err
	}
	for _, r := range rights {
		if v, ok := r.Type.(model.CanReadAs); ok && v.Party == p {
			return nil // already has it
		}
	}

	_, err = cl.UserMng.GrantUserRights(ctx, userID, "", []*model.Right{
		{Type: model.CanReadAs{Party: p}},
	})
	return err
}

func normTmplID(s string) string {
	return strings.TrimPrefix(s, "#")
}

func partyForUserOrAllocate(
	t *testing.T,
	ctx context.Context,
	p *client.DamlBindingClient,
	hint string,
) string {
	// 1) Try user -> primary party
	u, err := p.UserMng.GetUser(ctx, UserName)
	if err == nil && u != nil && u.PrimaryParty != "" {
		return u.PrimaryParty
	}

	// 2) Fallback: try find existing party
	parties, err := p.PartyMng.ListKnownParties(ctx, "", 1000, "")
	require.NoError(t, err)

	// Prefer display name match if available in your wrapper
	for _, pd := range parties.PartyDetails {
		if pd.LocalMetadata != nil && pd.LocalMetadata["displayName"] == hint {
			return pd.Party
		}
	}
	// Fallback: match by prefix "Bob::"
	for _, pd := range parties.PartyDetails {
		if pd.Party == hint || strings.HasPrefix(pd.Party, hint+"::") {
			return pd.Party
		}
	}

	// 3) Allocate if missing
	resp, err := p.PartyMng.AllocateParty(ctx, hint, map[string]string{"type": "testing"}, "")
	require.NoError(t, err)
	return resp.Party
}

// TestCoin

// What we're testing?
// - Deploy/ensure code is available on all participiants (upload DAR)
// - Alice (P1) creates coinRegistry contract for an instrument LINK
// - Bob (P2) creates a MintPreapproval (permission for minting)
// - (P1) Alice exercises a mintChoice (BurnMintFactory_BurnMint) on the registry to mint to Bob
//   - This mint produces;
//   - a Coin holding wrapper contract (i.e. our Coin app contract)
//   - also an Splice.Api.Token.HoldingV1:Holding on the "real token holding" contract
//
// - (P2) Bob transfers part of his holdings to charlie using splice TransferFactory interface choice on registry
//   - this transfer requires inputHoldingsCids of type ContractID Splice.Api.Token.HoldingV1:Holding
//
// - (P3) charlie accepts the transfer from bob
// - (P1) Alice grants mint rights to Dave
// - (P5) Erin grants MintPreapproval from Alice
// - (P5) Dave uses MintRole to mint to Erin
func TestCoin(t *testing.T) {
	ctx := context.Background()

	jwtToken, err := getJWT()
	require.NoError(t, err)

	p1 := newBindingClient(t, jwtToken, "participant1.grpc-ledger-api.localhost:8080", "participant1.admin-api.localhost:8080")
	p2 := newBindingClient(t, jwtToken, "participant2.grpc-ledger-api.localhost:8080", "participant2.admin-api.localhost:8080")
	p3 := newBindingClient(t, jwtToken, "participant3.grpc-ledger-api.localhost:8080", "participant3.admin-api.localhost:8080")
	p4 := newBindingClient(t, jwtToken, "participant4.grpc-ledger-api.localhost:8080", "participant4.admin-api.localhost:8080")
	p5 := newBindingClient(t, jwtToken, "participant5.grpc-ledger-api.localhost:8080", "participant5.admin-api.localhost:8080")

	_, err = p1.VersionService.GetLedgerAPIVersion(ctx, &model.GetLedgerAPIVersionRequest{})
	require.NoError(t, err)

	// Upload DARs (Admin service: PackageMng)
	darContent, err := os.ReadFile("../../contracts/coin/.daml/dist/coin-0.0.1.dar")
	require.NoError(t, err)

	upload := func(p *client.DamlBindingClient) {
		submissionID := "validate-" + time.Now().Format("20060102150405")
		err = p.PackageMng.ValidateDarFile(ctx, darContent, submissionID)
		require.NoError(t, err)

		uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")

		err := p.PackageMng.UploadDarFile(ctx, darContent, uploadSubmissionID)
		require.NoError(t, err)
	}
	upload(p1)
	upload(p2)
	upload(p3)
	upload(p4)
	upload(p5)

	// ******************* TEST0 Alice creates coinRegistry contract for an instrument LINK *******************
	// Get primary party for this user
	u, err := p1.UserMng.GetUser(ctx, UserName)
	require.NoError(t, err)
	partyAlice := u.PrimaryParty
	require.NotEmpty(t, partyAlice)
	require.NoError(t, ensureCanReadAs(ctx, p1, UserName, partyAlice))

	parties, err := p1.PartyMng.ListKnownParties(ctx, "", 100, "")
	require.NoError(t, err)
	for i, pd := range parties.PartyDetails {
		fmt.Printf("%d: party=%q display=%q\n", i, pd.Party, pd.IdentityProviderID)
	}

	partyBob := partyForUserOrAllocate(t, ctx, p2, "Bob")
	partyCharlie := partyForUserOrAllocate(t, ctx, p3, "Charlie")
	partyDave := partyForUserOrAllocate(t, ctx, p4, "Dave")
	partyErin := partyForUserOrAllocate(t, ctx, p5, "Erin")

	// Create registry contract using our generated bindings
	reg := coin.CoinRegistry{
		Issuer: types.PARTY(partyAlice),
		InstrumentId: coin.InstrumentId{
			Admin: types.PARTY(partyAlice),
			Id:    "LINK",
		},
		Meta: coin.Metadata{
			Values: types.TEXTMAP{},
		},
	}

	// Submit via binding client's CommandService (ledger service)
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "coin-test",
			UserID:     UserName,
			CommandID:  uuid.Must(uuid.NewUUID()).String(),
			ActAs:      []string{partyAlice},
			Commands:   []*model.Command{{Command: reg.CreateCommand()}},
		},
	}

	submitRespNew, err := p1.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	require.NoError(t, err)

	// retrieve the contract ID and template ID from the create event
	registryContractID := ""
	registryTemplateID := ""
	for _, event := range submitRespNew.Transaction.Events {
		if event.Created == nil {
			continue
		}

		if normalizeTemplateKey(event.Created.TemplateID) == "Coin.Registry:CoinRegistry" {
			registryContractID = event.Created.ContractID
			registryTemplateID = event.Created.TemplateID
			break
		}

	}
	require.NotEmpty(t, registryContractID)
	fmt.Println("Registry Contract ID: ", registryContractID)

	disclosedRegistry, err := QueryDisclosedContractWithBindingClient(
		ctx, p1, partyAlice, registryContractID, registryTemplateID,
	)
	require.NoError(t, err)

	// ******************* TEST1 Bob creates MintPreapproval *******************
	mintPreapproval := coin.MintPreapproval{
		Receiver: types.PARTY(partyBob),
		Sender:   types.PARTY(partyAlice),
	}

	cmds = &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "coin-test",
			UserID:     UserName,
			CommandID:  uuid.Must(uuid.NewUUID()).String(),
			ActAs:      []string{partyBob}, // Bob party
			Commands:   []*model.Command{{Command: mintPreapproval.CreateCommand()}},
		},
	}

	submitResp, err := p2.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	require.NoError(t, err)

	bobMintPreapprovalCid := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}
		if normalizeTemplateKey(event.Created.TemplateID) == "Coin.Registry:MintPreapproval" {
			bobMintPreapprovalCid = event.Created.ContractID
			break
		}
	}

	require.NotEmpty(t, bobMintPreapprovalCid, "no MintPreapproval create found in transaction")
	t.Logf("Bob created MintPreapproval, CID: %s", bobMintPreapprovalCid)

	// ******************* TEST2 Alice mints to Bob (USING GENERATED BINDINGS + SubmitAndWait) *******************
	mintArgs := coin.BurnMintFactoryBurnMint{
		ExpectedAdmin: types.PARTY(partyAlice),
		InstrumentId: coin.InstrumentId{
			Admin: types.PARTY(partyAlice),
			Id:    "LINK",
		},

		InputHoldingCids: []types.CONTRACT_ID{},

		Outputs: []coin.BurnMintOutput{
			{
				Owner:  types.PARTY(partyBob),
				Amount: types.NUMERIC(big.NewInt(4213)),
				Context: coin.ChoiceContext{
					Values: types.TEXTMAP{
						"mint-preapproval": coin.AnyValue{
							AVContractId: func() *types.CONTRACT_ID {
								cid := types.CONTRACT_ID(bobMintPreapprovalCid)
								return &cid
							}(),
						},
					},
				},
			},
		},

		ExtraActors: []types.PARTY{},
		ExtraArgs: coin.ExtraArgs{
			Context: coin.ChoiceContext{Values: types.TEXTMAP{}},
			Meta:    coin.Metadata{Values: types.TEXTMAP{}},
		},
	}

	// Exercise the choice via binding (on the registry contract id)
	exerciseCmd := coin.CoinRegistry{}.BurnMintFactoryBurnMint(registryContractID, mintArgs)

	pkgs, err := p1.PackageMng.ListKnownPackages(ctx)
	require.NoError(t, err)

	var burnMintPkgID string
	for _, p := range pkgs {
		if strings.Contains(strings.ToLower(p.Name), "splice-api-token-burn-mint") {
			burnMintPkgID = p.PackageID
			break
		}
	}
	require.NotEmpty(t, burnMintPkgID, "burn/mint package not found on ledger (did the DAR include it?)")

	t.Logf("burnMintPkg   id=%s", burnMintPkgID)

	mintReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "coin-test-mint",
			UserID:     UserName,
			CommandID:  uuid.Must(uuid.NewUUID()).String(),
			ActAs:      []string{partyAlice},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				TemplateID: fmt.Sprintf("%s:%s:%s", burnMintPkgID, "Splice.Api.Token.BurnMintV1", "BurnMintFactory"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
		},
	}

	mintResp, err := p1.CommandService.SubmitAndWaitForTransaction(ctx, mintReq)
	require.NoError(t, err)

	bobHoldingCid := ""
	for _, event := range mintResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		if normalizeTemplateKey(event.Created.TemplateID) == "Coin.Holding:CoinHolding" {
			bobHoldingCid = event.Created.ContractID
			break
		}
	}

	require.NotEmpty(t, bobHoldingCid, "no CoinHolding create found in mint transaction")
	t.Logf("Alice minted to Bob, CoinHolding CID: %s", bobHoldingCid)

	// ******************* TEST3: Bob transfers part of their holdings to Charlie (no Charlie rights needed) *******************

	transferArgs := coin.TransferFactoryTransfer{
		ExpectedAdmin: types.PARTY(partyAlice),
		Transfer: coin.Transfer22{
			Sender:        types.PARTY(partyBob),
			Receiver:      types.PARTY(partyCharlie),
			Amount:        types.NUMERIC(big.NewInt(10)),
			InstrumentId:  coin.InstrumentId{Admin: types.PARTY(partyAlice), Id: "LINK"},
			RequestedAt:   types.TIMESTAMP(time.Now()),
			ExecuteBefore: types.TIMESTAMP(time.Now().Add(24 * time.Hour)),
			InputHoldingCids: []types.CONTRACT_ID{
				types.CONTRACT_ID(bobHoldingCid),
			},
			Meta: coin.Metadata{Values: types.TEXTMAP{}},
		},
		ExtraArgs: coin.ExtraArgs{
			Context: coin.ChoiceContext{Values: types.TEXTMAP{}},
			Meta:    coin.Metadata{Values: types.TEXTMAP{}},
		},
	}

	// Exercise command via your generated binding (interface choice)
	exerciseCmd = coin.CoinRegistry{}.TransferFactoryTransfer(registryContractID, transferArgs)

	pkgsP2, err := p2.PackageMng.ListKnownPackages(ctx)
	require.NoError(t, err)

	var transferInstructionPkgID string
	for _, p := range pkgsP2 {
		if strings.Contains(strings.ToLower(p.Name), "splice-api-token-transfer-instruction") {
			transferInstructionPkgID = p.PackageID
			break
		}
	}
	require.NotEmpty(t, transferInstructionPkgID, "transfer instruction package not found on ledger (did the DAR include it?)")
	t.Logf("transferInstructionPkgID   id=%s", transferInstructionPkgID)

	//  Interactive submission for disclosure
	cmdID := uuid.Must(uuid.NewUUID()).String()

	prepReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			UserID:    UserName,
			CommandID: cmdID,
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				TemplateID: fmt.Sprintf("%s:%s:%s", transferInstructionPkgID, "Splice.Api.Token.TransferInstructionV1", "TransferFactory"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
			ActAs:  []string{partyBob},
			ReadAs: []string{partyBob},
			DisclosedContracts: []*model.DisclosedContract{
				disclosedRegistry,
			},
		},
	}

	transferResp, err := p2.CommandService.SubmitAndWaitForTransaction(ctx, prepReq)
	require.NoError(t, err)

	fmt.Println("transferResp: ", transferResp)

	transferInstructionCid := ""
	for _, event := range transferResp.Transaction.Events {
		if event.Created != nil {
			tid := strings.TrimPrefix(event.Created.TemplateID, "#")
			// only capture the TransferInstruction CID
			// tid format: <pkgid>:<Module>:<Entity>
			parts := strings.Split(tid, ":")
			if len(parts) >= 3 {
				module := parts[len(parts)-2]
				entity := parts[len(parts)-1]
				if module == "Coin.Transfer" && entity == "CoinTransferInstruction" {
					transferInstructionCid = event.Created.ContractID
				}
			}
		}
	}
	require.NotEmpty(t, transferInstructionCid)

	require.NotEmpty(t, transferInstructionCid, "no CoinTransferInstruction created in transfer tx")
	t.Logf("TransferInstruction CID: %s", transferInstructionCid)

	pkgs, err = p3.PackageMng.ListKnownPackages(ctx)
	require.NoError(t, err)

	var TransferInstructionPkgID string
	for _, p := range pkgs {
		if strings.Contains(strings.ToLower(p.Name), "splice-api-token-transfer-instruction") {
			TransferInstructionPkgID = p.PackageID
			break
		}
	}
	require.NotEmpty(t, TransferInstructionPkgID, "transfer instructions package not found on ledger (did the DAR include it?)")

	t.Logf("TransferInstructionPkgID   id=%s", TransferInstructionPkgID)

	// ******************* TEST4 Charlie accepts transfer from Bob *******************
	acceptArgs := coin.TransferInstructionAccept{
		ExtraArgs: coin.ExtraArgs{
			Context: coin.ChoiceContext{Values: types.TEXTMAP{}},
			Meta:    coin.Metadata{Values: types.TEXTMAP{}},
		},
	}

	acceptCmd := coin.CoinTransferInstruction{}.TransferInstructionAccept(transferInstructionCid, acceptArgs)

	// Submit as Charlie on participant3
	acceptReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "coin-test-accept",
			UserID:     UserName,
			CommandID:  uuid.Must(uuid.NewUUID()).String(),
			ActAs:      []string{partyCharlie},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				TemplateID: fmt.Sprintf("%s:%s:%s", transferInstructionPkgID, "Splice.Api.Token.TransferInstructionV1", "TransferInstruction"),
				ContractID: transferInstructionCid,
				Choice:     acceptCmd.Choice,
				Arguments:  acceptCmd.Arguments,
			}}},
			DisclosedContracts: []*model.DisclosedContract{
				disclosedRegistry,
			},
		},
	}

	acceptResp, err := p3.CommandService.SubmitAndWaitForTransaction(ctx, acceptReq)
	require.NoError(t, err)
	fmt.Println("acceptResp:", acceptResp)

	charlieHoldingCid := ""
	for _, event := range acceptResp.Transaction.Events {
		if event.Created == nil {
			continue
		}
		charlieHoldingCid = event.Created.ContractID
	}

	t.Logf("Charlie accepted transfer, holding CID: %s", charlieHoldingCid)

	// ******************* TEST 5 Alice grant mint rights to Dave *******************
	roleArgs := coin.MintRole{
		Issuer:   types.PARTY(partyAlice),
		Minter:   types.PARTY(partyDave),
		Registry: types.CONTRACT_ID(registryContractID),
	}

	roleCmdID := uuid.Must(uuid.NewUUID()).String()

	roleReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "coin-test-grant-mint-rights",
			UserID:     UserName,
			CommandID:  roleCmdID,
			ActAs:      []string{partyAlice},
			Commands:   []*model.Command{{Command: roleArgs.CreateCommand()}},
		},
	}

	resp, err := p1.CommandService.SubmitAndWaitForTransaction(ctx, roleReq)
	require.NoError(t, err)

	daveMintRoleCid := ""
	for _, event := range resp.Transaction.Events {
		if event.Created == nil {
			continue
		}
		daveMintRoleCid = event.Created.ContractID
	}

	require.NotEmpty(t, daveMintRoleCid, "MintRole not created")
	t.Logf("Alice granted MintRole to Dave, CID: %s", daveMintRoleCid)

	// ******************* TEST6: Erin grants MintPreapproval *******************

	preapprovalArgs := coin.MintPreapproval{
		Receiver: types.PARTY(partyErin),
		Sender:   types.PARTY(partyAlice),
	}

	mintPreApprovalmdID := uuid.Must(uuid.NewUUID()).String()
	preapprovalReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "coin-test-grant-mint-preapproval",
			UserID:     UserName,
			CommandID:  mintPreApprovalmdID,
			ActAs:      []string{partyErin},
			Commands:   []*model.Command{{Command: preapprovalArgs.CreateCommand()}},
		},
	}

	resp, err = p5.CommandService.SubmitAndWaitForTransaction(ctx, preapprovalReq)
	require.NoError(t, err)

	erinMintPreapprovalCid := ""
	erinMintPreapprovalTemplateID := ""
	for _, event := range resp.Transaction.Events {
		if event.Created == nil {
			continue
		}
		erinMintPreapprovalTemplateID = event.Created.TemplateID
		erinMintPreapprovalCid = event.Created.ContractID
	}

	require.NotEmpty(t, erinMintPreapprovalCid, "MintPreapproval not created")
	t.Logf("Erin created MintPreapproval, CID: %s", erinMintPreapprovalCid)

	disclosedErinPreApproval, err := QueryDisclosedContractWithBindingClient(
		ctx, p5, partyErin, erinMintPreapprovalCid, erinMintPreapprovalTemplateID,
	)
	require.NoError(t, err)

	// ******************* TEST7: Dave uses the MintRole to mint to Erin *******************
	// AFTER submit, wait for Erin to observe it async because this event hasn't happened yet:

	leBeforeErinMint, err := p5.StateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	require.NoError(t, err)

	holdingTmpl := normTmplID(coin.CoinHolding{}.GetTemplateID())
	// 1) Start updates stream on Erin's participant (p5), filtered to Erin + CoinHolding template
	streamCtx, cancelStream := context.WithCancel(context.Background())
	defer cancelStream()

	updCh, errCh := p5.UpdateService.GetUpdates(streamCtx, &model.GetUpdatesRequest{
		BeginExclusive: leBeforeErinMint.Offset,
		Filter: &model.TransactionFilter{
			FiltersByParty: map[string]*model.Filters{
				partyErin: {
					Inclusive: &model.InclusiveFilters{
						TemplateFilters: []*model.TemplateFilter{
							{
								TemplateID:              coin.CoinHolding{}.GetTemplateID(),
								IncludeCreatedEventBlob: true,
							},
						},
					},
				},
			},
		},
		UpdateFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				partyErin: {
					Inclusive: &model.InclusiveFilters{
						TemplateFilters: []*model.TemplateFilter{
							{
								TemplateID:              coin.CoinHolding{}.GetTemplateID(),
								IncludeCreatedEventBlob: true,
							},
						},
					},
				},
			},
			Verbose: true,
		},
		Verbose: true,
	})
	require.NoError(t, err)

	// 2) Channel to return Erin's holding CID
	erinCidChan := make(chan string, 1)

	go func() {
		defer close(erinCidChan)

		for resp := range updCh {
			upd := resp.Update
			if upd == nil || upd.Transaction == nil {
				continue
			}
			tx := upd.Transaction

			for _, ev := range tx.Events {
				if ev == nil || ev.Created == nil {
					continue
				}
				if normTmplID(ev.Created.TemplateID) != holdingTmpl {
					continue
				}

				// Found Erin’s CoinHolding create
				erinCidChan <- ev.Created.ContractID
				cancelStream()
				return
			}
		}
	}()

	cid := types.CONTRACT_ID(erinMintPreapprovalCid)

	// Dave uses the MintRole to mint to Erin
	daveMintRoleArgs := coin.MintRoleMint{
		InstrumentId: coin.InstrumentId{
			Admin: types.PARTY(partyAlice),
			Id:    "LINK",
		},
		Outputs: []coin.BurnMintOutput{
			{
				Owner:  types.PARTY(partyErin),
				Amount: types.NUMERIC(big.NewInt(7777)), // example
				Context: coin.ChoiceContext{
					Values: types.TEXTMAP{
						"mint-preapproval": coin.AnyValue{
							AVContractId: &cid,
						},
					},
				},
			},
		},
	}

	daveMintRoleCmd := coin.MintRole{}.MintRoleMint(daveMintRoleCid, daveMintRoleArgs)

	daveMintRoleCmdID := uuid.Must(uuid.NewUUID()).String()
	daveMintRoleReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID:  daveMintRoleCmdID,
			WorkflowID: "coin-test-mint-to-erin",
			Commands:   []*model.Command{{Command: daveMintRoleCmd}},
			ActAs:      []string{partyDave},
			DisclosedContracts: []*model.DisclosedContract{
				disclosedRegistry,
				disclosedErinPreApproval,
			},
		},
	}

	resp, err = p4.CommandService.SubmitAndWaitForTransaction(ctx, daveMintRoleReq)
	require.NoError(t, err)
	fmt.Println("dave mint to erin resp: ", resp)

	// 4) Wait for Erin’s holding CID from the async stream
	select {
	case erinHoldingCid := <-erinCidChan:
		require.NotEmpty(t, erinHoldingCid, "stream ended without observing CoinHolding")
		t.Logf("Erin observed CoinHolding CID: %s", erinHoldingCid)

	case err := <-errCh:
		t.Fatalf("updates stream error: %v", err)

	case <-time.After(30 * time.Second):
		t.Fatal("Timeout waiting for Erin CoinHolding create via updates stream")
	}
}

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}
