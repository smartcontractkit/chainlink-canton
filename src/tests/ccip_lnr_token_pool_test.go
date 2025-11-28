package tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/transferInstructionV1"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-canton-internal/src/protocol"
)

func TestLnRTokenPool_Release(t *testing.T) {
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ccipParticipant, err := NewParticipant("participant1.admin-api.localhost:8080", "participant1.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	userParticipant, err := NewParticipant("participant2.admin-api.localhost:8080", "participant2.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)
	tokenPoolOwnerParticipant, err := NewParticipant("participant3.admin-api.localhost:8080", "participant3.grpc-ledger-api.localhost:8080")
	require.NoError(t, err)

	// ========================
	// |   DAR Uploading     |
	// ========================

	// Read Dars
	offRampDar, err := os.ReadFile("../../contracts/ccip/offramp/.daml/dist/ccip-offramp-1.0.0.dar")
	require.NoError(t, err)
	tokenAdminRegistryDar, err := os.ReadFile("../../contracts/ccip/tokenAdminRegistry/.daml/dist/ccip-tokenadminregistry-1.0.0.dar")
	require.NoError(t, err)
	committeeVerifierDar, err := os.ReadFile("../../contracts/ccip/ccvs/.daml/dist/ccip-committeeverifier-1.0.0.dar")
	require.NoError(t, err)
	tokenPoolDar, err := os.ReadFile("../../contracts/ccip/pools/lockReleaseTokenPool/.daml/dist/ccip-lockreleasetokenpool-1.0.0.dar")
	require.NoError(t, err)

	// Upload Dars
	// CCIP Participant
	packageIDs, err := UploadDARstoMultipleParticipants(ctx,
		[][]byte{
			offRampDar,
			tokenAdminRegistryDar,
			committeeVerifierDar,
		},
		ccipParticipant,
	)
	require.NoError(t, err)
	fmt.Printf("Uploaded Dars to ccipParticipant: %s\n", packageIDs)

	// Token Pool Owner Participant
	packageIDs, err = UploadDARstoMultipleParticipants(ctx,
		[][]byte{
			offRampDar,
			tokenAdminRegistryDar,
			committeeVerifierDar,
			tokenPoolDar,
		},
		tokenPoolOwnerParticipant,
	)
	require.NoError(t, err)
	fmt.Printf("Uploaded Dars to tokenPoolOwnerParticipant: %s\n", packageIDs)

	// User Participant
	packageIDs, err = UploadDARstoMultipleParticipants(ctx,
		[][]byte{
			offRampDar,
			tokenAdminRegistryDar,
			committeeVerifierDar,
			tokenPoolDar,
		},
		userParticipant,
	)
	require.NoError(t, err)
	fmt.Printf("Uploaded Dars to userParticipant: %s\n", packageIDs)

	parties, err := EnsurePartyOnMultipleParticipants(ctx, ccipParticipant, userParticipant, tokenPoolOwnerParticipant)
	require.NoError(t, err)
	fmt.Printf("Allocated parties on participants: %v\n", parties)
	partyCCIP := parties[0]
	partyUser := parties[1]
	partyTokenPoolOwner := parties[2]
	_ = partyCCIP
	_ = partyUser
	_ = partyTokenPoolOwner

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

	// Mint Tokens to TokenPoolOwner Party
	tokenHoldingCid, err := MintAMT(ctx, tokenPoolOwnerParticipant, metadataClient, transferInstructionClient, scanProxyClient, partyTokenPoolOwner, "100.00")
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
	encodedInstrumentIdAmt, err := encodeInstrumentId(registryAdmin, "Amulet")
	require.NoError(t, err)
	_ = chainSelector
	_ = destChainSelector1
	_ = instrumentIdAmt

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

	// Deploy CCV owned by the Token Pool Owner Party
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

	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}},
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
			ActAs: []string{partyTokenPoolOwner},
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

	// ==============================
	// |     CCIP Deployment        |
	// ==============================

	// Deploy TokenAdminRegistry
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-tokenadminregistry",
								ModuleName: "CCIP.TokenAdminRegistry",
								EntityName: "TokenAdminRegistry",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "owner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "tokenConfigs",
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
	tokenAdminRegistryCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenAdminRegistryCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed TokenAdminRegistry to: %s\n", tokenAdminRegistryCid)

	// Deploy OffRamp
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

	ccipApi := &CCIPApi{
		CCIPParty:   partyCCIP,
		Participant: ccipParticipant,
	}

	// TokenPoolOwner deploys LockReleaseTokenPool
	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-lockreleasetokenpool",
								ModuleName: "CCIP.LockReleaseTokenPool",
								EntityName: "LockReleaseTokenPool",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "ccipOwner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}},
								}, {
									Label: "poolOwner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}},
								}, {
									Label: "instrumentId",
									Value: instrumentIdAmt,
								}, {
									Label: "decimals",
									Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 6}},
								}, {
									Label: "requiredCCVs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_ContractId{ContractId: ccv1Cid}},
										{Sum: &apiv2.Value_ContractId{ContractId: ccv2Cid}},
									}}}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	tokenPoolCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenPoolCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed LockReleaseTokenPool to: %s\n", tokenPoolCid)

	// =============================
	// |   Token Admin Setup       |
	// =============================

	// CCIP Owner calls ProposeAdministrator to add the tokenPoolOwner as an admin
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-tokenadminregistry",
								ModuleName: "CCIP.TokenAdminRegistry",
								EntityName: "TokenAdminRegistry",
							},
							ContractId: tokenAdminRegistryCid,
							Choice:     "TokenAdminRegistry_ProposeAdministrator",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "instrumentId",
									Value: instrumentIdAmt,
								}, {
									Label: "newAdmin",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}},
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
			tokenAdminRegistryCid = e.Created.ContractId
		}
	}
	fmt.Printf("Called ProposeAdministrator on TokenAdminRegistry: %v\n", tokenAdminRegistryCid)

	// TokenPoolOwner calls AcceptAdministrator to accept admin role
	disclosedTokenAdminRegistry, err := ccipApi.GetTokenAdminRegistry(ctx)
	require.NoError(t, err)
	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-tokenadminregistry",
								ModuleName: "CCIP.TokenAdminRegistry",
								EntityName: "TokenAdminRegistry",
							},
							ContractId: disclosedTokenAdminRegistry.ContractId,
							Choice:     "TokenAdminRegistry_AcceptAdminRole",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "instrumentId",
									Value: instrumentIdAmt,
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTokenAdminRegistry},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenAdminRegistryCid = e.Created.ContractId
		}
	}
	fmt.Printf("Called AcceptAdminRole on TokenAdminRegistry: %v\n", tokenAdminRegistryCid)

	// TokenPoolOwner calls SetPool to set themselves as the token pool owner
	disclosedTokenAdminRegistry, err = ccipApi.GetTokenAdminRegistry(ctx)
	require.NoError(t, err)
	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-tokenadminregistry",
								ModuleName: "CCIP.TokenAdminRegistry",
								EntityName: "TokenAdminRegistry",
							},
							ContractId: disclosedTokenAdminRegistry.ContractId,
							Choice:     "TokenAdminRegistry_SetPool",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "instrumentId",
									Value: instrumentIdAmt,
								}, {
									Label: "optTokenPoolOwner",
									Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}}}},
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTokenAdminRegistry},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			tokenAdminRegistryCid = e.Created.ContractId
		}
	}
	fmt.Printf("Called SetPool on TokenAdminRegistry: %v\n", tokenAdminRegistryCid)

	tokenPoolApi := &TokenPoolApi{
		Participant: tokenPoolOwnerParticipant,
		AdminParty:  partyTokenPoolOwner,
	}

	// =============================
	// |        Execution          |
	// =============================

	// User calls ReleaseOrMint on the token pool to execute the message
	disclosedTokenAdminRegistry, err = ccipApi.GetTokenAdminRegistry(ctx)
	require.NoError(t, err)
	disclosedOffRamp, err := ccipApi.GetOffRamp(ctx)
	require.NoError(t, err)
	disclosedCCV1, err := ccipApi.GetCCV(ctx)
	require.NoError(t, err)
	disclosedCCV2, err := tokenPoolApi.GetCCV(ctx)
	require.NoError(t, err)
	disclosedTokenPool, err := tokenPoolApi.GetTokenPool(ctx)
	require.NoError(t, err)
	disclosedHoldings, err := tokenPoolApi.GetHoldings(ctx)
	require.NoError(t, err)
	disclosedHolding := disclosedHoldings[len(disclosedHoldings)-1]
	fmt.Printf("Using Holding: %s\n", disclosedHolding.ContractId)
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := GetTransferFactory(ctx, transferInstructionClient, registryAdmin, partyTokenPoolOwner, partyUser)
	require.NoError(t, err)
	disclosures := slices.Concat([]*apiv2.DisclosedContract{disclosedTokenAdminRegistry, disclosedOffRamp, disclosedCCV1, disclosedCCV2, disclosedTokenPool, disclosedHolding}, transferFactoryDisclosures)

	message, err := protocol.NewMessage(
		11111,
		22222,
		0,
		[]byte("123456"),
		[]byte(offRampCid),
		0,
		100_000,
		[]byte("sender"),
		[]byte(partyUser),
		[]byte("destblob"),
		[]byte("data"),
		protocol.NewTokenTransfer(
			big.NewInt(1e7), // 10 AMT with 6 decimals
			[]byte("sourcePoolAddress"),
			[]byte("sourceTokenAddress"),
			encodedInstrumentIdAmt,
			[]byte(partyUser),
			[]byte("extraData"),
		),
	)
	require.NoError(t, err)
	encodedMessage, err := message.Encode()
	require.NoError(t, err)
	fmt.Printf("Encoded message: %s\n", hex.EncodeToString(encodedMessage))

	// CCV1 - sign with 2 of 3 signers
	ccvData1, err := signMessage(encodedMessage, ccv1Signer1, ccv1Signer2, ccv1Signer3)
	require.NoError(t, err)
	fmt.Printf("CCV1 Signature Data: %s\n", hex.EncodeToString(ccvData1))

	// CCV2 - sign with 2 of 2 signers
	ccvData2, err := signMessage(encodedMessage, ccv2Signer1, ccv2Signer2)
	require.NoError(t, err)
	fmt.Printf("CCV2 Signature Data: %s\n", hex.EncodeToString(ccvData2))

	fmt.Println("==== Executing CCIP Message on Token Pool ====")
	response, err := userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#ccip-tokenpool-interfaces",
								ModuleName: "CCIP.Interfaces.TokenPool",
								EntityName: "ITokenPool",
							},
							ContractId: disclosedTokenPool.ContractId,
							Choice:     "TokenPool_ReleaseOrMint",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "encodedMessage",
									Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(encodedMessage)}},
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
								}, {
									Label: "caller",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyUser}},
								}, {
									Label: "tokenInput",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "transferFactory",
											Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}},
										}, {
											Label: "tokenPoolHoldings",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{Sum: &apiv2.Value_ContractId{ContractId: disclosedHolding.GetContractId()}},
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
									Label: "offRampCid",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedOffRamp.ContractId}},
								}, {
									Label: "tokenAdminRegistryCid",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedTokenAdminRegistry.ContractId}},
								},
							}}}},
						},
					},
				},
			},
			ActAs:              []string{partyUser},
			DisclosedContracts: disclosures,
		},
	})
	require.NoError(t, err)
	fmt.Printf("Executed CCIP Message, Update: %v\n", response.GetTransaction().GetUpdateId())

	// Query User Holdings
	userHoldings, err := GetActiveContractsForPartyInterface(ctx, userParticipant, partyUser, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	require.NoError(t, err)
	fmt.Printf("User Party %s has %d token holdings after Execution\n", partyCCIP, len(userHoldings))
	for i, holding := range userHoldings {
		balance := holding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		fmt.Printf(" %v - Token Holding Cid: %s, Balance: %s\n", i, holding.GetCreatedEvent().GetContractId(), balance)
	}
}

func encodeInstrumentId(admin, identifier string) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteByte(byte(len(admin)))
	buf.WriteString(admin)
	buf.WriteByte(byte(len(identifier)))
	buf.WriteString(identifier)
	return buf.Bytes(), nil
}
