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
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/transferInstructionV1"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-canton-internal/src/protocol"
)

func TestLnRTokenPool_ReleaseAndLock(t *testing.T) {
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
	onRampDar, err := os.ReadFile("../../contracts/ccip/onramp/.daml/dist/ccip-onramp-1.0.0.dar")
	require.NoError(t, err)
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
			onRampDar,
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
			onRampDar,
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
			onRampDar,
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
	fmt.Printf("Minted 100 AMT to Token Pool Owner, Token Holding Cid: %s\n", tokenHoldingCid)

	// Mint Tokens to User Party (used to pay for fees later)
	userTokenHoldingCid, err := MintAMT(ctx, userParticipant, metadataClient, transferInstructionClient, scanProxyClient, partyUser, "1.00")
	require.NoError(t, err)
	fmt.Printf("Minted 1 AMT to User, Token Holding Cid: %s\n", userTokenHoldingCid)

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
									Label: "messageSentObserver",
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
									Label: "messageSentObserver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}},
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
		}, {
			Label: "offRamp",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}},
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

	// Deploy FeeQuoter
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
									Label: "tokenTransferFeeConfigs",
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
	time.Sleep(time.Second)
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
	var userHoldingCids []*apiv2.Value
	fmt.Printf("User Party %s has %d token holdings after Execution\n", partyUser, len(userHoldings))
	for i, holding := range userHoldings {
		balance := holding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		fmt.Printf(" %v - Token Holding Cid: %s, Balance: %s\n", i, holding.GetCreatedEvent().GetContractId(), balance)
		// Skip the initially minted holding, as that will be used to pay for fees
		if holding.GetCreatedEvent().GetContractId() == userTokenHoldingCid {
			continue
		}
		userHoldingCids = append(userHoldingCids, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: holding.GetCreatedEvent().GetContractId()}})
	}

	// ============================
	// |        CCIP Send         |
	// ============================
	disclosedTokenAdminRegistry, err = ccipApi.GetTokenAdminRegistry(ctx)
	require.NoError(t, err)
	disclosedOnRamp, err := ccipApi.GetOnRamp(ctx)
	require.NoError(t, err)
	disclosedFeeQuoter, err := ccipApi.GetFeeQuoter(ctx)
	require.NoError(t, err)
	disclosedCCV1, err = ccipApi.GetCCV(ctx)
	require.NoError(t, err)
	disclosedCCV2, err = tokenPoolApi.GetCCV(ctx)
	require.NoError(t, err)
	disclosedTokenPool, err = tokenPoolApi.GetTokenPool(ctx)
	require.NoError(t, err)
	disclosures = slices.Concat([]*apiv2.DisclosedContract{disclosedTokenAdminRegistry, disclosedOnRamp, disclosedFeeQuoter, disclosedCCV1, disclosedCCV2, disclosedTokenPool}, transferFactoryDisclosures)

	fmt.Println("==== Sending CCIP Message from Token Pool ====")
	response, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
							Choice:     "TokenPool_LockOrBurn",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "destChainSelector",
									Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector1}},
								}, {
									Label: "message",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "receiver",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "B6D4805bf6943c5875C0C7b67EDa24b2bDACBF6e"}},
										}, {
											Label: "payload",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}},
										}, {
											Label: "feeToken",
											Value: instrumentIdAmt,
										}, {
											Label: "extraArgs",
											Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}},
										}, {
											Label: "tokenAmounts",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{{
												Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
													{
														Label: "instrumentId",
														Value: instrumentIdAmt,
													}, {
														Label: "amount",
														Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "5000000"}}, // 5 AMT with 6 decimals
													},
												}}},
											}}}}},
										},
									}}}},
								}, {
									Label: "tokenInput",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "transferFactory",
											Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}},
										}, {
											Label: "tokenPoolHoldings",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}},
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
									Label: "feeInput",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "instrumentId",
											Value: instrumentIdAmt,
										}, {
											Label: "transferFactory",
											Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}},
										}, {
											Label: "inputHoldingCids",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{{Sum: &apiv2.Value_ContractId{ContractId: userTokenHoldingCid}}}}}},
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
									Label: "senderInputCids",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: userHoldingCids}}},
								}, {
									Label: "onRampCid",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedOnRamp.GetContractId()}},
								}, {
									Label: "feeQuoterCid",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedFeeQuoter.GetContractId()}},
								}, {
									Label: "tokenAdminRegistryCid",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedTokenAdminRegistry.GetContractId()}},
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
			DisclosedContracts: disclosures,
		},
	})
	require.NoError(t, err)
	fmt.Printf("Sent CCIP Message, Update: %v\n", response.GetTransaction().GetUpdateId())

	// Assertions
	// Query User Holdings
	userHoldings, err = GetActiveContractsForPartyInterface(ctx, userParticipant, partyUser, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	require.NoError(t, err)
	fmt.Printf("User Party %s has %d token holdings after Execution\n", partyUser, len(userHoldings))
	for i, holding := range userHoldings {
		balance := holding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		fmt.Printf(" %v - Token Holding Cid: %s, Balance: %s\n", i, holding.GetCreatedEvent().GetContractId(), balance)
	}

	// Query Token Pool Holdings
	tokenPoolHoldings, err := GetActiveContractsForPartyInterface(ctx, tokenPoolOwnerParticipant, partyTokenPoolOwner, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	require.NoError(t, err)
	fmt.Printf("Token Pool Party %s has %d token holdings after Execution\n", partyTokenPoolOwner, len(tokenPoolHoldings))
	for i, holding := range tokenPoolHoldings {
		balance := holding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		fmt.Printf(" %v - Token Holding Cid: %s, Balance: %s\n", i, holding.GetCreatedEvent().GetContractId(), balance)
	}

	// Query CCIP Owner Holdings
	ccipOwnerHoldings, err := GetActiveContractsForPartyInterface(ctx, ccipParticipant, partyCCIP, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	require.NoError(t, err)
	fmt.Printf("CCIP Owner Party %s has %d token holdings after Execution\n", partyCCIP, len(ccipOwnerHoldings))
	for i, holding := range ccipOwnerHoldings {
		balance := holding.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		fmt.Printf(" %v - Token Holding Cid: %s, Balance: %s\n", i, holding.GetCreatedEvent().GetContractId(), balance)
	}

	// Query CCIPMessageSent events
	messageSentEvents, err := GetActiveContractsForPartyTemplateId(ctx, ccipParticipant, partyCCIP, &apiv2.Identifier{
		PackageId:  "#ccip-onramp",
		ModuleName: "CCIP.OnRamp",
		EntityName: "CCIPMessageSent",
	})
	require.NoError(t, err)
	fmt.Printf("CCIP Owner Party %s has %d CCIPMessageSent events after Send\n", partyCCIP, len(messageSentEvents))
	for i, msgEvent := range messageSentEvents {
		eventData := msgEvent.GetCreatedEvent().GetCreateArguments().GetFields()[3].GetValue().GetSum().(*apiv2.Value_Record).Record
		j, err := protojson.Marshal(eventData)
		require.NoError(t, err)
		fmt.Printf(" %v - CCIPMessageSent Cid: %s, data: %s\n", i, msgEvent.GetCreatedEvent().GetContractId(), string(j))
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
