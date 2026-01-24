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

	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
	"github.com/smartcontractkit/chainlink-canton-internal/integration-tests/testhelpers"
	"github.com/smartcontractkit/chainlink-canton-internal/openapi/gen/transferInstructionV1"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

func encodeInstrumentId(admin, identifier string) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteByte(byte(len(admin)))
	buf.WriteString(admin)
	buf.WriteByte(byte(len(identifier)))
	buf.WriteString(identifier)
	return buf.Bytes(), nil
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

	// Instrument ID for AMT
	instrumentIdAmt := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: registryAdmin}}},
		{Label: "id", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "Amulet"}}},
	}}}}

	// CCV Setup
	// Generate signer keys for CommitteeVerifier
	var ccvSignerKeys []*ecdsa.PrivateKey
	var ccvSignerPubKeys []string
	for i := 0; i < 3; i++ {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}
	t.Logf("Generated %d CCV signer keys", len(ccvSignerKeys))

	// Deploy CCVRegistry
	res, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.CCVRegistry", EntityName: "CCVRegistry"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccvregistry-receive"}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	ccvRegistryCid := extractCreatedContractId(res)
	t.Logf("Deployed CCVRegistry: %s", ccvRegistryCid)

	// Deploy CommitteeVerifier
	versionTag := "49ff34ed"
	ccvId := versionTag + "@" + partyCCIP
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccv-receive"}}},
						{Label: "versionTag", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: versionTag}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "messageSentObserver", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "storageLocation", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ipfs://test-receive"}}},
						{Label: "threshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2}}},
						{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[0]}},
							{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[1]}},
							{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[2]}},
						}}}}},
					}},
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

	// Deploy TokenAdminRegistry
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-tar-receive"}}},
						{Label: "tokenConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid := extractCreatedContractId(res)
	t.Logf("Deployed TokenAdminRegistry: %s", tokenAdminRegistryCid)

	// Deploy OffRamp
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-offramp-receive"}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	offRampCid := extractCreatedContractId(res)
	t.Logf("Deployed OffRamp: %s", offRampCid)

	// Deploy PerPartyRouterFactory
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-factory-receive"}}},
						{Label: "registeredRouters", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
					}},
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
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-receiver"}}},
					}}}},
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
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "decimals", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 6}}},
						{Label: "chainCCVRequirements", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "poolReceiveContext", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
						}}}}},
						// Use RelativeHours 24 for Amulet tokens (have expiry constraints)
						{Label: "transferTimeout", Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
							Constructor: "RelativeHours",
							Value:       &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 24}},
						}}}},
					}},
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
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "newAdmin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
					}}}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid = extractCreatedContractId(res)
	t.Log("Called ProposeAdministrator")

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
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid = extractCreatedContractId(res)
	t.Log("Called AcceptAdminRole")

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
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "optTokenPoolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}}}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
					}}}},
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
	encodedInstrumentId, err := encodeInstrumentId(registryAdmin, "Amulet")
	require.NoError(t, err)
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

	// Build MessageV1 Daml record for CommitteeVerifier
	// TokenTransferV1 record
	tokenTransferRecord := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "amount", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tokenAmount.String()}}},
		{Label: "sourcePoolAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString([]byte(partyTokenPoolOwner))}}},
		{Label: "sourceTokenAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
		{Label: "destTokenAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: encodedInstrumentIdHex}}},
		{Label: "tokenReceiver", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString([]byte(partyReceiver))}}},
		{Label: "extraData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
	}}}}

	// MessageV1 record
	messageV1Record := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "sourceChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "123"}}},
		{Label: "destChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "456"}}},
		{Label: "sequenceNumber", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "1"}}},
		{Label: "executionGasLimit", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 200000}}},
		{Label: "ccipReceiveGasLimit", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 100000}}},
		{Label: "finality", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2000}}},
		{Label: "ccvAndExecutorHash", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0000000000000000000000000000000000000000000000000000000000000000"}}},
		{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString([]byte("0000000000000000000000000000000000000001"))}}},
		{Label: "offRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString([]byte("0000000000000000000000000000000000000002"))}}},
		{Label: "sender", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString([]byte("0000000000000000000000000000000000000003"))}}},
		{Label: "receiver", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(EncodePartyID(partyReceiver))}}},
		{Label: "destBlob", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
		{Label: "tokenTransfer", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: tokenTransferRecord}}}},
		{Label: "messageData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
	}}}}

	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier"},
					ContractId: disclosedCCV.ContractId,
					Choice:     "CommitteeVerifier_VerifyMessage",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccvRegistryCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedCCVRegistry.ContractId}}},
						{Label: "message", Value: messageV1Record},
						{Label: "messageId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: messageHashHex}}},
						{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: verifierResultsHex}}},
						{Label: "receiver", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
					}}}},
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
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "offRampCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedOffRamp.ContractId}}},
						{Label: "globalConfigCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedGlobalConfig.ContractId}}},
						{Label: "tokenAdminRegistryCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedTar.ContractId}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: encodedMessageHex}}},
						{Label: "ccvVerifyTickets", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_ContractId{ContractId: ccvVerifyTicketCid}},
						}}}}},
						{Label: "tokenPoolCCVTicket", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}},
						{Label: "receiverRequiredCCVIds", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}}}},
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

		var acceptDisclosures []*apiv2.DisclosedContract
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

		res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
