package tests

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"https://github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"https://github.com/smartcontractkit/go-daml/pkg/types"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/lockreleasetokenpool"
	offrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/perpartyrouter"
	tokenadminregistryBinding "github.com/smartcontractkit/chainlink-canton/bindings/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

func encodeInstrumentId(admin, identifier string) []byte {
	var buf bytes.Buffer

	buf.WriteByte(byte(len(admin)))
	buf.WriteString(admin)
	buf.WriteByte(byte(len(identifier)))
	buf.WriteString(identifier)

	return buf.Bytes()
}

// TestLnRTokenPool_FullReceiveFlow tests the complete CCIP inbound token release flow.
// - Deploy all CCIP contracts including CommitteeVerifier, GlobalConfig, OffRamp, PerPartyRouter
// - Mint tokens to pool (simulating prior locked tokens)
// - Build inbound message with TokenTransfer
// - Generate signatures and call CommitteeVerifier_VerifyMessage to get CCVVerifyTicket
// - Call PerPartyRouter.Execute to process message and get TokenReceiveTicket
// - Call TokenPool_ReleaseFromTicket to transfer tokens from pool to receiver
// - Verify receiver received the tokens
func TestLnRTokenPool_FullReceiveFlow(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(3))

	// Setup participants
	ccipParticipant := env.Participant(1)
	receiverParticipant := env.Participant(2)
	tokenPoolOwnerParticipant := env.Participant(3)

	// DAR Uploading
	// Read DARs
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err)

	// Upload DARs to all participants
	dars := [][]byte{commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, tokenPoolDar, perPartyRouterDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant, tokenPoolOwnerParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.Party
	partyReceiver := receiverParticipant.Party
	partyTokenPoolOwner := tokenPoolOwnerParticipant.Party
	t.Logf("Parties: CCIP=%s, Receiver=%s, PoolOwner=%s", partyCCIP, partyReceiver, partyTokenPoolOwner)

	// Get DSO Admin Party (registry admin)
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), env.Splice.TokenMetadataClient)
	require.NoError(t, err)

	// Token Setup
	// Mint tokens to Token Pool Owner (these will be "locked" in the pool)
	poolHoldingCid, err := testhelpers.MintAMT(t.Context(), tokenPoolOwnerParticipant, env.Splice.TokenMetadataClient, env.Splice.TransferInstructionClient, tokenPoolOwnerParticipant.ScanProxyClient, partyTokenPoolOwner, "100.00")
	require.NoError(t, err)
	t.Logf("Minted 100 AMT to Pool Owner, Holding CID: %s", poolHoldingCid)

	// Instrument ID for AMT using bindings
	instrumentIdAmt := tokenadminregistryBinding.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}

	// CCV Setup
	// Generate signer keys for CommitteeVerifier
	ccvSignerKeys := make([]*ecdsa.PrivateKey, 0, 3)
	ccvSignerPubKeys := make([]string, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}
	t.Logf("Generated %d CCV signer keys", len(ccvSignerKeys))

	// Deploy CCVRegistry using bindings
	res, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.CCVRegistry", EntityName: "CCVRegistry"},
					CreateArguments: ledger.ConvertToRecord(common.CCVRegistry{
						CcipOwner:  types.PARTY(partyCCIP),
						InstanceId: types.TEXT("test-ccvregistry-receive"),
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	ccvRegistryCid := extractCreatedContractId(res)
	t.Logf("Deployed CCVRegistry: %s", ccvRegistryCid)

	// Deploy CommitteeVerifier using bindings
	versionTag := "49ff34ed"
	ccvId := versionTag + "@" + partyCCIP
	ccvSignerPubKeysTypes := make([]types.TEXT, len(ccvSignerPubKeys))
	for i, pk := range ccvSignerPubKeys {
		ccvSignerPubKeysTypes[i] = types.TEXT(pk)
	}
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier"},
					CreateArguments: ledger.ConvertToRecord(ccvs.CommitteeVerifier{
						Owner:               types.PARTY(partyCCIP),
						InstanceId:          types.TEXT("test-ccv-receive"),
						CcipOwner:           types.PARTY(partyCCIP),
						VersionTag:          types.TEXT(versionTag),
						MessageSentObserver: types.PARTY(partyCCIP),
						StorageLocation:     types.TEXT("ipfs://test-receive"),
						Threshold:           types.INT64(2),
						Signers:             ccvSignerPubKeysTypes,
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	ccvCid := extractCreatedContractId(res)
	t.Logf("Deployed CommitteeVerifier: %s (ccvId: %s)", ccvCid, ccvId)

	// CCIP Deployment
	sourceChainSelector := "123"
	destChainSelector := "456"

	// Deploy GlobalConfig with source chain config including the CCV
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-globalconfig-receive"}}},
						{Label: "chainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector}}},
						{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
						{Label: "destChainConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "sourceChainConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{
							{
								Key: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}},
								Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
									{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}}},
									{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0000000000000000000000000000000000000001"}}},
									{Label: "laneMandatedCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_Text{Text: ccvId}},
									}}}}},
									{Label: "defaultCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
								}}}},
							},
						}}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	globalConfigCid := extractCreatedContractId(res)
	t.Logf("Deployed GlobalConfig: %s", globalConfigCid)

	// Deploy TokenAdminRegistry using bindings
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					CreateArguments: ledger.ConvertToRecord(tokenadminregistryBinding.TokenAdminRegistry{
						Owner:        types.PARTY(partyCCIP),
						InstanceId:   types.TEXT("test-tar-receive"),
						TokenConfigs: types.GENMAP{},
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid := extractCreatedContractId(res)
	t.Logf("Deployed TokenAdminRegistry: %s", tokenAdminRegistryCid)

	// Deploy OffRamp using bindings
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp"},
					CreateArguments: ledger.ConvertToRecord(offrampBinding.OffRamp{
						CcipOwner:  types.PARTY(partyCCIP),
						InstanceId: types.TEXT("test-offramp-receive"),
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	offRampCid := extractCreatedContractId(res)
	t.Logf("Deployed OffRamp: %s", offRampCid)

	// Deploy PerPartyRouterFactory using bindings
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					CreateArguments: ledger.ConvertToRecord(perpartyrouter.PerPartyRouterFactory{
						CcipOwner:         types.PARTY(partyCCIP),
						InstanceId:        types.TEXT("test-factory-receive"),
						RegisteredRouters: types.GENMAP{},
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	factoryCid := extractCreatedContractId(res)
	t.Logf("Deployed PerPartyRouterFactory: %s", factoryCid)

	// Create PerPartyRouter for receiver
	disclosedFactory, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	})
	require.NoError(t, err)

	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: disclosedFactory.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.CreateRouter{
						PartyOwner: types.PARTY(partyReceiver),
						InstanceId: types.TEXT("test-router-receiver"),
					}),
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedFactory},
		},
	})
	require.NoError(t, err)
	routerCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "PerPartyRouter" {
				routerCid = e.Created.ContractId
				break
			}
		}
	}
	require.NotEmpty(t, routerCid)
	t.Logf("Created PerPartyRouter for receiver: %s", routerCid)

	// Token Pool Setup
	// Deploy LockReleaseTokenPool
	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-lockreleasetokenpool", ModuleName: "CCIP.LockReleaseTokenPool", EntityName: "LockReleaseTokenPool"},
					CreateArguments: ledger.ConvertToRecord(lockreleasetokenpool.LockReleaseTokenPool{
						CcipOwner:            types.PARTY(partyCCIP),
						PoolOwner:            types.PARTY(partyTokenPoolOwner),
						InstanceId:           types.TEXT("test-pool-receive"),
						InstrumentId:         lockreleasetokenpool.InstrumentId{Admin: instrumentIdAmt.Admin, Id: instrumentIdAmt.Id},
						Decimals:             types.INT64(6),
						ChainCCVRequirements: types.GENMAP{},
						PoolReceiveContext:   lockreleasetokenpool.ChoiceContext{Values: types.TEXTMAP{}},
						TransferTimeout:      lockreleasetokenpool.TransferTimeout{RelativeHours: func() *types.INT64 { i := types.INT64(24); return &i }()},
					}),
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	tokenPoolCid := extractCreatedContractId(res)
	t.Logf("Deployed LockReleaseTokenPool: %s", tokenPoolCid)

	// Register pool in TokenAdminRegistry (ProposeAdministrator -> AcceptAdminRole -> SetPool)
	// Step 1: ProposeAdministrator
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: tokenAdminRegistryCid,
					Choice:     "TokenAdminRegistry_ProposeAdministrator",
					ChoiceArgument: ledger.MapToValue(tokenadminregistryBinding.TokenAdminRegistryProposeAdministrator{
						InstrumentId: instrumentIdAmt,
						NewAdmin:     types.PARTY(partyTokenPoolOwner),
						Caller:       types.PARTY(partyCCIP),
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid = extractCreatedContractId(res)
	t.Log("Called ProposeAdministrator, tokenAdminRegistryCid:", tokenAdminRegistryCid)

	// Step 2: AcceptAdminRole
	time.Sleep(500 * time.Millisecond)
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_AcceptAdminRole",
					ChoiceArgument: ledger.MapToValue(tokenadminregistryBinding.TokenAdminRegistryAcceptAdminRole{
						InstrumentId: instrumentIdAmt,
						Caller:       types.PARTY(partyTokenPoolOwner),
					}),
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid = extractCreatedContractId(res)
	t.Log("Called AcceptAdminRole, tokenAdminRegistryCid:", tokenAdminRegistryCid)

	// Step 3: SetPool
	time.Sleep(500 * time.Millisecond)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	res, err = tokenPoolOwnerParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_SetPool",
					ChoiceArgument: ledger.MapToValue(tokenadminregistryBinding.TokenAdminRegistrySetPool{
						InstrumentId:      instrumentIdAmt,
						OptTokenPoolOwner: func() *types.PARTY { p := types.PARTY(partyTokenPoolOwner); return &p }(),
						Caller:            types.PARTY(partyTokenPoolOwner),
					}),
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid = extractCreatedContractId(res)
	t.Logf("Called SetPool, TokenAdminRegistry: %s", tokenAdminRegistryCid)

	// Build Message
	// Encode instrumentId for destTokenAddress
	encodedInstrumentId := encodeInstrumentId(registryAdmin, "Amulet")
	encodedInstrumentIdHex := hex.EncodeToString(encodedInstrumentId)

	// Build token transfer (5 AMT in Splice Decimal format)
	tokenAmount := big.NewInt(5)
	encodedTokenTransfer := buildTokenTransferV1(tokenAmount, partyTokenPoolOwner, encodedInstrumentIdHex, partyReceiver)

	// Build message
	msg := &MessageV1{
		SourceChainSelector: 123,
		DestChainSelector:   456,
		SequenceNumber:      1,
		ExecutionGasLimit:   200000,
		CCIPReceiveGasLimit: 100000,
		Finality:            2000,
		CCVAndExecutorHash:  [32]byte{},
		OnRampAddress:       []byte("0000000000000000000000000000000000000001"),
		OffRampAddress:      []byte("0000000000000000000000000000000000000002"),
		Sender:              []byte("0000000000000000000000000000000000000003"),
		Receiver:            EncodePartyID(partyReceiver),
		DestBlob:            []byte{},
		TokenTransfer:       encodedTokenTransfer,
		MessageData:         []byte{},
	}
	encodedMessage, err := EncodeMessageV1(msg)
	require.NoError(t, err)
	encodedMessageHex := hex.EncodeToString(encodedMessage)
	messageHash := crypto.Keccak256(encodedMessage)
	messageHashHex := hex.EncodeToString(messageHash)
	t.Logf("Message hash: %s", messageHashHex)

	// CCV Verification
	// Generate verifierResults with 2 of 3 signatures
	verifierResults, err := GenerateVerifierResults(encodedMessage, ccvSignerKeys[:2])
	require.NoError(t, err)
	verifierResultsHex := hex.EncodeToString(verifierResults)

	// Call CommitteeVerifier_VerifyMessage to get CCVVerifyTicket
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err)
	disclosedCCVRegistry, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.CCVRegistry", EntityName: "CCVRegistry",
	})
	require.NoError(t, err)

	// Build MessageV1 Daml record for CommitteeVerifier using bindings
	// TokenTransferV1 record
	tokenTransfer := &ccvs.TokenTransferV1{
		Amount:             types.NUMERIC(tokenAmount.String()),
		SourcePoolAddress:  types.TEXT(hex.EncodeToString([]byte(partyTokenPoolOwner))),
		SourceTokenAddress: types.TEXT(""),
		DestTokenAddress:   types.TEXT(encodedInstrumentIdHex),
		TokenReceiver:      types.TEXT(hex.EncodeToString([]byte(partyReceiver))),
		ExtraData:          types.TEXT(""),
	}

	// MessageV1 record
	messageV1 := ccvs.MessageV1{
		SourceChainSelector: types.NUMERIC("123"),
		DestChainSelector:   types.NUMERIC("456"),
		SequenceNumber:      types.NUMERIC("1"),
		ExecutionGasLimit:   types.INT64(200000),
		CcipReceiveGasLimit: types.INT64(100000),
		Finality:            types.INT64(2000),
		CcvAndExecutorHash:  types.TEXT("0000000000000000000000000000000000000000000000000000000000000000"),
		OnRampAddress:       types.TEXT(hex.EncodeToString([]byte("0000000000000000000000000000000000000001"))),
		OffRampAddress:      types.TEXT(hex.EncodeToString([]byte("0000000000000000000000000000000000000002"))),
		Sender:              types.TEXT(hex.EncodeToString([]byte("0000000000000000000000000000000000000003"))),
		Receiver:            types.TEXT(hex.EncodeToString(EncodePartyID(partyReceiver))),
		DestBlob:            types.TEXT(""),
		TokenTransfer:       tokenTransfer,
		MessageData:         types.TEXT(""),
	}

	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier"},
					ContractId: disclosedCCV.ContractId,
					Choice:     "CommitteeVerifier_VerifyMessage",
					ChoiceArgument: ledger.MapToValue(ccvs.CommitteeVerifierVerifyMessage{
						CcvRegistryCid:  types.CONTRACT_ID(disclosedCCVRegistry.ContractId),
						Message:         messageV1,
						MessageId:       types.TEXT(messageHashHex),
						VerifierResults: types.TEXT(verifierResultsHex),
						Receiver:        types.PARTY(partyReceiver),
						Caller:          types.PARTY(partyReceiver),
					}),
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedCCV, disclosedCCVRegistry},
		},
	})
	require.NoError(t, err)
	ccvVerifyTicketCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "CCVVerifyTicket" {
				ccvVerifyTicketCid = e.Created.ContractId
				break
			}
		}
	}
	require.NotEmpty(t, ccvVerifyTicketCid)
	t.Logf("Got CCVVerifyTicket: %s", ccvVerifyTicketCid)

	// Execute Message
	// Get all disclosures needed for execute
	time.Sleep(500 * time.Millisecond)
	disclosedRouter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	})
	require.NoError(t, err)
	disclosedOffRamp, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp",
	})
	require.NoError(t, err)
	disclosedGlobalConfig, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig",
	})
	require.NoError(t, err)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)
	disclosedCCVVerifyTicket, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.Tickets", EntityName: "CCVVerifyTicket",
	})
	require.NoError(t, err)

	// Call PerPartyRouter.Execute
	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter"},
					ContractId: disclosedRouter.ContractId,
					Choice:     "Execute",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.Execute{
						OffRampCid:             types.CONTRACT_ID(disclosedOffRamp.ContractId),
						GlobalConfigCid:        types.CONTRACT_ID(disclosedGlobalConfig.ContractId),
						TokenAdminRegistryCid:  types.CONTRACT_ID(disclosedTar.ContractId),
						EncodedMessage:         types.TEXT(encodedMessageHex),
						CcvVerifyTickets:       []types.CONTRACT_ID{types.CONTRACT_ID(ccvVerifyTicketCid)},
						TokenPoolCCVTicket:     nil,
						ReceiverRequiredCCVIds: []types.TEXT{},
					}),
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedRouter, disclosedOffRamp, disclosedGlobalConfig, disclosedTar, disclosedCCVVerifyTicket},
		},
	})
	require.NoError(t, err)

	// Extract TokenReceiveTicket from result
	tokenReceiveTicketCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "TokenReceiveTicket" {
				tokenReceiveTicketCid = e.Created.ContractId
				break
			}
		}
	}
	require.NotEmpty(t, tokenReceiveTicketCid, "Expected TokenReceiveTicket to be created")
	t.Logf("Execute completed, TokenReceiveTicket: %s", tokenReceiveTicketCid)

	// Release Tokens from Pool
	// Get TransferFactory and disclosures for release
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := testhelpers.GetTransferFactory(t.Context(), env.Splice.TransferInstructionClient, registryAdmin, partyTokenPoolOwner, partyReceiver)
	require.NoError(t, err)

	disclosedPool, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), tokenPoolOwnerParticipant, &apiv2.Identifier{
		PackageId: "#ccip-lockreleasetokenpool", ModuleName: "CCIP.LockReleaseTokenPool", EntityName: "LockReleaseTokenPool",
	})
	require.NoError(t, err)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)
	disclosedTokenReceiveTicket, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.Tickets", EntityName: "TokenReceiveTicket",
	})
	require.NoError(t, err)

	// Get pool's holdings
	poolHoldings, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), tokenPoolOwnerParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	require.NotEmpty(t, poolHoldings, "Pool should have holdings")

	// Filter for unlocked holdings only (Amulet tokens may be locked until next mining round)
	// HoldingV1.view has: owner, instrumentId, amount, lock, meta
	// lock field (index 3) is Optional Lock - None means unlocked
	var poolHoldingCids []*apiv2.Value
	var poolHoldingDisclosures []*apiv2.DisclosedContract
	var lockedCount, unlockedCount int
	for _, h := range poolHoldings {
		// Check if holding is locked by examining the interface view
		viewFields := h.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()
		lockField := viewFields[3].GetValue() // lock : Optional Lock
		isLocked := lockField.GetOptional().GetValue() != nil

		if isLocked {
			lockedCount++
			continue
		}
		unlockedCount++

		poolHoldingCids = append(poolHoldingCids, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: h.GetCreatedEvent().GetContractId()}})
		// Build disclosure for each holding
		poolHoldingDisclosures = append(poolHoldingDisclosures, &apiv2.DisclosedContract{
			ContractId:       h.GetCreatedEvent().GetContractId(),
			TemplateId:       h.GetCreatedEvent().GetTemplateId(),
			CreatedEventBlob: h.GetCreatedEvent().GetCreatedEventBlob(),
		})
	}
	t.Logf("Pool has %d holdings (%d unlocked, %d locked)", len(poolHoldings), unlockedCount, lockedCount)

	if unlockedCount == 0 {
		t.Skip("SKIPPING: All pool holdings are locked (Amulet tokens are locked until next mining round). " +
			"This is expected on fresh localnet. Either wait for mining round to complete, or use TestToken instead of Amulet.")
	}

	// Build disclosures for release (include pool holdings)
	releaseDisclosures := slices.Concat(
		[]*apiv2.DisclosedContract{disclosedPool, disclosedTar, disclosedTokenReceiveTicket},
		transferFactoryDisclosures,
		poolHoldingDisclosures,
	)

	// Capture receiver's balance before release
	receiverHoldingsBefore, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceBefore := getHoldingsBalance(receiverHoldingsBefore)

	// Call TokenPool_ReleaseFromTicket - receiver triggers their own release
	// This is the real flow: receiver calls from their participant with their authorization
	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenpool-interfaces", ModuleName: "CCIP.Interfaces.TokenPool", EntityName: "ITokenPool"},
					ContractId: disclosedPool.ContractId,
					Choice:     "TokenPool_ReleaseFromTicket",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "tokenReceiveTicketCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: tokenReceiveTicketCid}}},
						{Label: "tokenAdminRegistryCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedTar.ContractId}}},
						{Label: "tokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}}},
							{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "context", Value: choiceContext},
								{Label: "meta", Value: emptyMetadata},
							}}}}},
							{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: poolHoldingCids}}}},
						}}}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: releaseDisclosures,
		},
	})
	require.NoError(t, err)

	// Determine if transfer completed directly or created a pending TransferInstruction.
	// Amulet transfers without receiver preapproval create a pending AmuletTransferInstruction.
	var pendingTransferInstructionCid string
	var releaseCompleted bool
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			templateName := e.Created.GetTemplateId().GetEntityName()
			if templateName == "AmuletTransferInstruction" {
				pendingTransferInstructionCid = e.Created.GetContractId()
			}
			if templateName == "Amulet" || templateName == "LockedAmulet" {
				releaseCompleted = true
			}
		}
	}

	// If pending, receiver must accept the TransferInstruction to complete the transfer
	if !releaseCompleted && pendingTransferInstructionCid != "" {
		// Brief delay for contract propagation before querying
		time.Sleep(500 * time.Millisecond)

		// Fetch fresh accept context - OpenMiningRound changes every round so we must get
		// current disclosures immediately before exercising the accept choice.
		// https://github.com/hyperledger-labs/splice/blob/9a40442f3cf09d19285f8466f35ebaa07a4e5d58/daml/splice-amulet/daml/Splice/Amulet/TwoStepTransfer.daml#L124
		acceptContextResp, err := env.Splice.TransferInstructionClient.GetTransferInstructionAcceptContextWithResponse(t.Context(), pendingTransferInstructionCid, transferInstructionV1.GetChoiceContextRequest{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, acceptContextResp.StatusCode(), "Failed to get accept context: %s", string(acceptContextResp.Body))

		acceptDisclosures := make([]*apiv2.DisclosedContract, 0, len(acceptContextResp.JSON200.DisclosedContracts))
		for _, contract := range acceptContextResp.JSON200.DisclosedContracts {
			id, err := testhelpers.TemplateIdFromString(contract.TemplateId)
			require.NoError(t, err)
			createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
			require.NoError(t, err)
			acceptDisclosures = append(acceptDisclosures, &apiv2.DisclosedContract{
				TemplateId:       id,
				ContractId:       contract.ContractId,
				CreatedEventBlob: createdEventBlob,
			})
		}
		acceptContext, err := testhelpers.ChoiceContextFromData(acceptContextResp.JSON200.ChoiceContextData)
		require.NoError(t, err)

		_, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.Must(uuid.NewUUID()).String(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
						ContractId: pendingTransferInstructionCid,
						Choice:     "TransferInstruction_Accept",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "context", Value: acceptContext},
								{Label: "meta", Value: emptyMetadata},
							}}}}},
						}}}},
					}},
				}},
				ActAs:              []string{partyReceiver},
				DisclosedContracts: acceptDisclosures,
			},
		})
		require.NoError(t, err)
	}

	// Verify receiver's balance increased by the expected transfer amount
	receiverHoldingsAfter, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceAfter := getHoldingsBalance(receiverHoldingsAfter)

	expectedTransferAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	actualTransferAmount := receiverBalanceAfter - receiverBalanceBefore
	require.InDelta(t, expectedTransferAmount, actualTransferAmount, 0.01, "Receiver balance should increase by transfer amount")
	t.Logf("Receiver balance: %.2f -> %.2f AMT (transferred %.2f)", receiverBalanceBefore, receiverBalanceAfter, actualTransferAmount)
}

// getHoldingsBalance returns total balance across all holdings
func getHoldingsBalance(holdings []*apiv2.ActiveContract) float64 {
	var total float64
	for _, h := range holdings {
		balanceStr := h.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		balance, _ := new(big.Float).SetString(balanceStr)
		balanceFloat, _ := balance.Float64()
		total += balanceFloat
	}

	return total
}

// extractCreatedContractId returns the first created contract ID from a transaction response
func extractCreatedContractId(res *apiv2.SubmitAndWaitForTransactionResponse) string {
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			return e.Created.ContractId
		}
	}

	return ""
}

// buildTokenTransferV1 builds the encoded token transfer bytes for the message
func buildTokenTransferV1(amount *big.Int, sourcePoolOwner, destTokenAddressHex, tokenReceiverParty string) *TokenTransferV1 {
	sourcePoolAddress := []byte(sourcePoolOwner)
	sourceTokenAddress := []byte{} // Empty for Canton

	// Decode destTokenAddress from hex
	destTokenAddress, _ := hex.DecodeString(destTokenAddressHex)

	tokenReceiver := []byte(tokenReceiverParty)

	return &TokenTransferV1{
		Amount:             amount,
		SourcePoolAddress:  sourcePoolAddress,
		SourceTokenAddress: sourceTokenAddress,
		DestTokenAddress:   destTokenAddress,
		TokenReceiver:      tokenReceiver,
		ExtraData:          []byte{},
	}
}

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}
