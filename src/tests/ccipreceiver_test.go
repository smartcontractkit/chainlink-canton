package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

func TestCCIPReceiver(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	userParticipant, err := NewParticipant("participant2.admin-api.localhost:8080", "participant2.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	// Dars
	// Uploading the CCIP dar to the ccipParticipant
	// Uploading the CCIP Receiver dar to the userParticipant
	ccipDar, err := os.ReadFile("../../contracts/ccip/.daml/dist/ccip-1.0.0.dar")
	require.NoError(t, err)
	ccipReceiverDar, err := os.ReadFile("../../contracts/ccipreceiver/.daml/dist/ccipreceiver-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{ccipDar}, ccipParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP to ccipParticipant: %s\n", packageIDs)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{ccipReceiverDar}, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIPReceiver to userParticipant: %s\n", packageIDs)

	parties, err := EnsurePartyOnMultipleParticipants(ctx, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyCCIP := parties[0]
	partyUser := parties[1]
	_ = partyCCIP
	_ = partyUser

	// User deploys the CCIP Receiver contract
	res, err := userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccipreceiver",
								ModuleName: "CCIPReceiver.Receiver",
								EntityName: "ReceiverImplementation",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "ccip",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
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
	fmt.Printf("Deployed CCIP Receiver, CID: %v\n", receiverCid)
	// Query the contract for explicit disclosure
	disclosedReceiver, err := QueryDisclosedContract(ctx, receiverCid, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Queried registry for disclosure: %v\n", disclosedReceiver.GetContractId())

	// Upload the receiver DAR to the ccip participant
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{ccipReceiverDar}, ccipParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP Receiver DAR to ccipParticipant: %s\n", packageIDs)

	// CCIP calls the CCIPReceiver_Receive choice on the deployed receiver via the CCIPReceiver interface
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip",
								ModuleName: "CCIP.Receiver",
								EntityName: "CCIPReceiver",
							},
							ContractId: receiverCid,
							Choice:     "CCIPReceiver_Receive",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "input",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "ccip",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
										}, {
											Label: "messageId",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0x1234567890"}},
										},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyCCIP},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedReceiver,
			},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	fmt.Printf("Called CCIPReceiver_Receive in tx: %s\n", res.GetTransaction().GetUpdateId())
}
