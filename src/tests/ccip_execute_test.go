package tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
	participantv30 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/canton/admin/participant/v30"
	"github.com/smartcontractkit/chainlink-canton-internal/src/protocol"
)

func TestCCIPExecute(t *testing.T) {
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
	// Uploading the CCIP dars to all participants
	offRampDar, err := os.ReadFile("../../contracts/ccip/offramp/.daml/dist/ccip-offramp-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{offRampDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP OffRamp to ccipParticipant: %s\n", packageIDs)
	committeeVerifierDar, err := os.ReadFile("../../contracts/ccip/ccvs/.daml/dist/ccip-committeeverifier-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{committeeVerifierDar}, ccipParticipant, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIP CommitteeVerifier to ccipParticipant: %s\n", packageIDs)

	// Upload the CCIP Receiver Dar only to the userParticipant
	ccipReceiverDar, err := os.ReadFile("../../contracts/ccip/ccipreceiver/.daml/dist/ccip-receiver-1.0.0.dar")
	require.NoError(t, err)
	packageIDs, err = UploadDARstoMultipleParticipants(ctx, [][]byte{ccipReceiverDar}, userParticipant)
	require.NoError(t, err)
	fmt.Printf("Uploaded CCIPReceiver to userParticipant: %s\n", packageIDs)

	// List packages
	pkgs, err := userParticipant.PackageServiceClient.ListPackages(ctx, &participantv30.ListPackagesRequest{
		Limit:      0,
		FilterName: "",
	})
	require.NoError(t, err)
	for i, description := range pkgs.GetPackageDescriptions() {
		fmt.Printf("Package %d: %s\n", i, description)
	}

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

	// ========================
	// |   CCV Deployment     |
	// ========================

	// Deploy CCV1 owned by CCIP Party
	// Create signers
	ccv1Signer1, err := crypto.GenerateKey()
	require.NoError(t, err)
	ccv1Signer1PubKey := hex.EncodeToString(crypto.FromECDSAPub(ccv1Signer1.Public().(*ecdsa.PublicKey)))
	ccv1Signer2, err := crypto.GenerateKey()
	require.NoError(t, err)
	ccv1Signer2PubKey := hex.EncodeToString(crypto.FromECDSAPub(ccv1Signer2.Public().(*ecdsa.PublicKey)))
	ccv1Signer3, err := crypto.GenerateKey()
	require.NoError(t, err)
	ccv1Signer3PubKey := hex.EncodeToString(crypto.FromECDSAPub(ccv1Signer3.Public().(*ecdsa.PublicKey)))
	fmt.Println("CCV1 Signers:")
	fmt.Printf("1: address: %s pubKey: %s\n", crypto.PubkeyToAddress(ccv1Signer1.PublicKey).Hex(), ccv1Signer1PubKey)
	fmt.Printf("2: address: %s pubKey: %s\n", crypto.PubkeyToAddress(ccv1Signer2.PublicKey).Hex(), ccv1Signer2PubKey)
	fmt.Printf("3: address: %s pubKey: %s\n", crypto.PubkeyToAddress(ccv1Signer3.PublicKey).Hex(), ccv1Signer3PubKey)

	res, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-committeeverifier",
								ModuleName: "CCIP.CommitteeVerifier",
								EntityName: "CommitteeVerifier",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "ccipOwner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "storageLocation",
									Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ipfs://example_storage_location"}},
								}, {
									Label: "threshold",
									Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2}},
								}, {
									Label: "signers",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_Text{Text: ccv1Signer1PubKey}},
										{Sum: &apiv2.Value_Text{Text: ccv1Signer2PubKey}},
										{Sum: &apiv2.Value_Text{Text: ccv1Signer3PubKey}},
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
	ccv1Cid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			ccv1Cid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed CCV1 to: %s\n", ccv1Cid)

	// Deploy CCV owned by the user
	// Create signers
	ccv2Signer1, err := crypto.GenerateKey()
	require.NoError(t, err)
	ccv2Signer1PubKey := hex.EncodeToString(crypto.FromECDSAPub(ccv2Signer1.Public().(*ecdsa.PublicKey)))
	ccv2Signer2, err := crypto.GenerateKey()
	require.NoError(t, err)
	ccv2Signer2PubKey := hex.EncodeToString(crypto.FromECDSAPub(ccv2Signer2.Public().(*ecdsa.PublicKey)))
	fmt.Println("CCV2 Signers:")
	fmt.Printf("1: address: %s pubKey: %s\n", crypto.PubkeyToAddress(ccv2Signer1.PublicKey).Hex(), ccv2Signer1PubKey)
	fmt.Printf("2: address: %s pubKey: %s\n", crypto.PubkeyToAddress(ccv2Signer2.PublicKey).Hex(), ccv2Signer2PubKey)

	res, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-committeeverifier",
								ModuleName: "CCIP.CommitteeVerifier",
								EntityName: "CommitteeVerifier",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "ccipOwner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "storageLocation",
									Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ipfs://example_storage_location"}},
								}, {
									Label: "threshold",
									Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2}},
								}, {
									Label: "signers",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_Text{Text: ccv2Signer1PubKey}},
										{Sum: &apiv2.Value_Text{Text: ccv2Signer2PubKey}},
									}}}},
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
	ccv2Cid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			ccv2Cid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed CCV2 to: %s\n", ccv2Cid)

	// =========================
	// |   OffRamp deployment  |
	// =========================

	// CCIP Party deploys CCIP contracts
	sourceChainSelector := uint64(1111111111)
	destChainSelector := uint64(2222222222)
	_ = sourceChainSelector
	_ = destChainSelector

	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-offramp",
								ModuleName: "CCIP.OffRamp",
								EntityName: "OffRamp",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "sourceChainConfigs",
									Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}},
								}, {
									Label: "executionStates",
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
	offRampCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			offRampCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed OffRamp to: %s\n", offRampCid)

	// Apply source chain config to OffRamp
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-offramp",
								ModuleName: "CCIP.OffRamp",
								EntityName: "OffRamp",
							},
							ContractId: offRampCid,
							Choice:     "ApplySourceChainConfigUpdates",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "updates",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										&apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
											{
												Label: "sourceChainSelector",
												Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: strconv.Itoa(int(sourceChainSelector))}},
											}, {
												Label: "sourceChainConfig",
												Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
													{
														Label: "isEnabled",
														Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}},
													}, {
														Label: "onRamp",
														Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "1234"}},
													}, {
														Label: "laneMandatedCCVs",
														Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}},
													}, {
														Label: "defaultCCVs",
														Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
															{
																Sum: &apiv2.Value_ContractId{ContractId: ccv1Cid},
															},
														}}}},
													},
												}}}},
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
			offRampCid = e.Created.ContractId
		}
	}
	fmt.Printf("Applied SourceChainConfig to OffRamp: %v\n", offRampCid)

	// =========================
	// | CCIP Receiver Deploy  |
	// =========================

	// User deploys CCIP Receiver
	res, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-receiver",
								ModuleName: "CCIP.CCIPReceiver",
								EntityName: "CCIPReceiver",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "ccipOwner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "requiredCCVs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_ContractId{ContractId: ccv1Cid}},
										{Sum: &apiv2.Value_ContractId{ContractId: ccv2Cid}},
									}}}},
								}, {
									Label: "optionalCCVs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}},
								}, {
									Label: "threshold",
									Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}},
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
	ccipReceiverCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			ccipReceiverCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed CCIPReceiver to: %s\n", ccipReceiverCid)

	// =========================
	// |   Test Execution     |
	// =========================

	ccipApi := CCIPApi{
		Participant: ccipParticipant,
		CCIPParty:   partyCCIP,
	}
	_ = ccipApi

	message, err := protocol.NewMessage(
		sourceChainSelector,
		destChainSelector,
		0,
		[]byte("123456"),
		[]byte(offRampCid),
		0,
		100_000,
		[]byte("sender"),
		[]byte("receiver"),
		[]byte("destblob"),
		[]byte("data"),
		nil,
	)
	require.NoError(t, err)
	encodedMessage, err := message.Encode()
	require.NoError(t, err)
	fmt.Printf("Encoded message: %s\n", hex.EncodeToString(encodedMessage))

	// Generate signatures

	// CCV1 - sign with 2 of 3 signers
	ccvData1, err := signMessage(encodedMessage, ccv1Signer1, ccv1Signer2, ccv1Signer3)
	require.NoError(t, err)
	fmt.Printf("CCV1 Signature Data: %s\n", hex.EncodeToString(ccvData1))

	// CCV2 - sign with 2 of 2 signers
	ccvData2, err := signMessage(encodedMessage, ccv2Signer1, ccv2Signer2)
	require.NoError(t, err)
	fmt.Printf("CCV2 Signature Data: %s\n", hex.EncodeToString(ccvData2))

	// Call Execute on OffRamp
	disclosedOffRamp, err := ccipApi.GetOffRamp(ctx)
	require.NoError(t, err)
	disclosedCCV1, err := ccipApi.GetCCV(ctx)
	require.NoError(t, err)
	res, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-offramp",
								ModuleName: "CCIP.OffRamp",
								EntityName: "OffRamp",
							},
							ContractId: offRampCid,
							Choice:     "OffRamp_Execute",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "encodedMessage",
									Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(encodedMessage)}},
								}, {
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: ccipReceiverCid}},
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "ccvs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_ContractId{ContractId: ccv1Cid}},
										{Sum: &apiv2.Value_ContractId{ContractId: ccv2Cid}},
									}}}},
								}, {
									Label: "ccvData",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(ccvData1)}},
										{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(ccvData2)}},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyUser},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedOffRamp, disclosedCCV1},
		},
	})
	require.NoError(t, err)
	messageExecutedEventCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			messageExecutedEventCid = e.Created.ContractId
		}
	}
	fmt.Printf("Executed Message: %v\n", messageExecutedEventCid)
}

func signMessage(msg []byte, privateKeys ...*ecdsa.PrivateKey) ([]byte, error) {
	versionTag, err := hex.DecodeString("49ff34ed")
	if err != nil {
		return nil, err
	}

	preimage := append(versionTag, msg...) // Prefix the message with the version tag
	msgHash := crypto.Keccak256Hash(preimage)

	var signatures []byte
	for _, pk := range privateKeys {
		sig, err := crypto.Sign(msgHash[:], pk)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig[:64]...) // Remove v-value
	}

	var ccvData bytes.Buffer

	ccvData.Write(versionTag)
	if err := binary.Write(&ccvData, binary.BigEndian, uint16(len(signatures))); err != nil {
		return nil, err
	}
	ccvData.Write(signatures)
	return ccvData.Bytes(), nil
}
