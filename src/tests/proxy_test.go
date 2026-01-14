package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
)

func TestProxy(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	userParticipant, err := NewParticipant("participant2.admin-api.localhost:8080", "participant2.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	// ========================
	// |   DAR Uploading     |
	// ========================
	// Uploading the Proxy Dar to all participants
	proxyDar, err := os.ReadFile("../../contracts/test/proxy/.daml/dist/test-proxy-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{proxyDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded Proxy to ccipParticipant: %s\n", packageIDs)

	// Upload the Receiver Dar only to the userParticipant
	receiverDat, err := os.ReadFile("../../contracts/test/receiver/.daml/dist/test-receiver-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{receiverDat}, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded Receiver to userParticipant: %s\n", packageIDs)

	// ========================
	// |   Party Allocation   |
	// ========================
	// Allocate parties on both participants
	parties, err := EnsurePartyOnMultipleParticipants(ctx, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyCCIP := parties[0]
	partyUser := parties[1]
	_ = partyCCIP
	_ = partyUser

	// Deploy Proxy on ccipParticipant using partyCCIP

	res, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#test-proxy",
								ModuleName: "Proxy",
								EntityName: "Proxy",
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
	proxyCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			proxyCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed Proxy to: %s\n", proxyCid)

	// Deploy Receiver on userParticipant using partyUser

	res, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#test-receiver",
								ModuleName: "Receiver",
								EntityName: "ReceiverImplementation",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "value",
									Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 42}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyUser},
		},
	})
	require.NoError(t, err)
	receiverCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			receiverCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed Receiver to: %s\n", receiverCid)

	// ========================
	// |     Proxy Call       |
	// ========================
	disclosedProxy, err := QueryDisclosedContract(ctx, proxyCid, ccipParticipant)
	require.NoError(t, err)
	res, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#test-proxy",
								ModuleName: "Proxy",
								EntityName: "Proxy",
							},
							ContractId: proxyCid,
							Choice:     "Proxy_EmitEvent",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: receiverCid}},
								}, {
									Label: "eventValue",
									Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 1234}},
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyUser},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedProxy},
		},
	})
	require.NoError(t, err)
	fmt.Printf("Called Proxy to emit event, transaction ID: %s\n", res.GetTransaction().GetUpdateId())
}
