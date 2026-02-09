// Integration test for CCIP execute flow with CommitteeVerifier signature verification.
//
// This test demonstrates the complete execute flow WITHOUT token transfers:
//   - Deploy RMNRemote, GlobalConfig, CommitteeVerifier, OffRamp, PerPartyRouter
//   - Build a message with payload data
//   - PrepareExecute to create ExecutingMessageV1
//   - Generate ECDSA signatures and verify via CommitteeVerifier (appends CCV verification)
//   - Execute via PerPartyRouter
//   - Validate the returned message payload matches the original

package tests

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

// MessageV1 matches the Daml CCIP.MessageCodecV1.MessageV1 structure.
type MessageV1 struct {
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	ExecutionGasLimit   uint32
	CCIPReceiveGasLimit uint32
	Finality            uint16
	CCVAndExecutorHash  [32]byte
	OnRampAddress       []byte
	OffRampAddress      []byte
	Sender              []byte
	Receiver            []byte
	DestBlob            []byte
	TokenTransfer       *TokenTransferV1
	MessageData         []byte
}

// TokenTransferV1 matches the Daml CCIP.MessageCodecV1.TokenTransferV1 structure.
type TokenTransferV1 struct {
	Amount             *big.Int
	SourcePoolAddress  []byte
	SourceTokenAddress []byte
	DestTokenAddress   []byte
	TokenReceiver      []byte
	ExtraData          []byte
}

// EncodeMessageV1 encodes a MessageV1 to bytes matching the Daml MessageCodecV1.encodeMessageV1 format.
func EncodeMessageV1(msg *MessageV1) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteByte(0x01) // Version
	_ = binary.Write(&buf, binary.BigEndian, msg.SourceChainSelector)
	_ = binary.Write(&buf, binary.BigEndian, msg.DestChainSelector)
	_ = binary.Write(&buf, binary.BigEndian, msg.SequenceNumber)
	_ = binary.Write(&buf, binary.BigEndian, msg.ExecutionGasLimit)
	_ = binary.Write(&buf, binary.BigEndian, msg.CCIPReceiveGasLimit)
	_ = binary.Write(&buf, binary.BigEndian, msg.Finality)
	buf.Write(msg.CCVAndExecutorHash[:])

	// Length-prefixed fields (1-byte length)
	buf.WriteByte(uint8(len(msg.OnRampAddress))) //nolint:gosec
	buf.Write(msg.OnRampAddress)
	buf.WriteByte(uint8(len(msg.OffRampAddress))) //nolint:gosec
	buf.Write(msg.OffRampAddress)
	buf.WriteByte(uint8(len(msg.Sender))) //nolint:gosec
	buf.Write(msg.Sender)
	buf.WriteByte(uint8(len(msg.Receiver))) //nolint:gosec
	buf.Write(msg.Receiver)

	// 2-byte length prefixed fields
	_ = binary.Write(&buf, binary.BigEndian, uint16(len(msg.DestBlob))) //nolint:gosec
	buf.Write(msg.DestBlob)

	if msg.TokenTransfer != nil {
		tokenBytes := encodeTokenTransferV1(msg.TokenTransfer)
		_ = binary.Write(&buf, binary.BigEndian, uint16(len(tokenBytes))) //nolint:gosec
		buf.Write(tokenBytes)
	} else {
		_ = binary.Write(&buf, binary.BigEndian, uint16(0))
	}

	_ = binary.Write(&buf, binary.BigEndian, uint16(len(msg.MessageData))) //nolint:gosec
	buf.Write(msg.MessageData)

	return buf.Bytes(), nil
}

func encodeTokenTransferV1(tt *TokenTransferV1) []byte {
	var buf bytes.Buffer

	buf.WriteByte(0x01) // Version

	amountBytes := make([]byte, 32)
	if tt.Amount != nil {
		tt.Amount.FillBytes(amountBytes)
	}
	buf.Write(amountBytes)

	buf.WriteByte(uint8(len(tt.SourcePoolAddress))) //nolint:gosec
	buf.Write(tt.SourcePoolAddress)
	buf.WriteByte(uint8(len(tt.SourceTokenAddress))) //nolint:gosec
	buf.Write(tt.SourceTokenAddress)
	buf.WriteByte(uint8(len(tt.DestTokenAddress))) //nolint:gosec
	buf.Write(tt.DestTokenAddress)
	buf.WriteByte(uint8(len(tt.TokenReceiver))) //nolint:gosec
	buf.Write(tt.TokenReceiver)

	_ = binary.Write(&buf, binary.BigEndian, uint16(len(tt.ExtraData))) //nolint:gosec
	buf.Write(tt.ExtraData)

	return buf.Bytes()
}

// GenerateVerifierResults generates the verifierResults blob for CommitteeVerifier.
// Format: versionTag (4 bytes) || signatureLength (2 bytes) || signatures (64 bytes each)
// Matches EVM: signers sign keccak256(versionTag || messageId) where messageId = keccak256(encodedMessage).
func GenerateVerifierResults(encodedMessage []byte, privateKeys []*ecdsa.PrivateKey) ([]byte, error) {
	versionTag, _ := hex.DecodeString("49ff34ed")

	messageId := crypto.Keccak256(encodedMessage)
	msgHash := crypto.Keccak256(append(versionTag, messageId...))

	var signatures []byte
	for _, pk := range privateKeys {
		sig, err := crypto.Sign(msgHash, pk)
		if err != nil {
			return nil, fmt.Errorf("failed to sign: %w", err)
		}
		signatures = append(signatures, sig[:64]...) // r || s, drop v
	}

	var result bytes.Buffer
	result.Write(versionTag)
	_ = binary.Write(&result, binary.BigEndian, uint16(len(signatures))) //nolint:gosec
	result.Write(signatures)

	return result.Bytes(), nil
}

// EncodePartyID encodes a Canton party ID as a 32-byte keccak256 address.
// Matches Daml encodePartyAddress: keccak256(toHex(partyToText party)),
// which is equivalent to keccak256(partyBytes) since keccak256 hex-decodes its input.
func EncodePartyID(partyID string) []byte {
	return crypto.Keccak256([]byte(partyID))
}

// hexToBytes decodes a hex string to raw bytes. Panics on invalid hex.
func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex: " + s)
	}
	return b
}

// rawInstanceAddress wraps a text value as a Daml RawInstanceAddress newtype for the gRPC API.
func rawInstanceAddress(text string) *apiv2.Value {
	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: text}}},
	}}}}
}

// TestCCIPExecuteE2E tests the full execute flow without token transfers.
// Validates that the message payload returned from Execute matches the original.
func TestCCIPExecuteE2E(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	ccipParticipant := env.Participant(1)
	receiverParticipant := env.Participant(2)

	// Upload DARs
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err)
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipReceiverDar, err := contracts.GetDar(contracts.CCIPReceiver, contracts.CurrentVersion)
	require.NoError(t, err)

	dars := [][]byte{commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar, ccipReceiverDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.Party
	partyReceiver := receiverParticipant.Party
	t.Logf("Parties: CCIP=%s, Receiver=%s", partyCCIP, partyReceiver)

	// CCV Setup
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

	sourceChainSelector := "123"
	destChainSelector := "456"
	versionTag := "49ff34ed"
	ccvID := "test-ccv-e2e@" + partyCCIP

	// Deploy RMNRemote
	res, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-rmn-e2e"}}},
						{Label: "rmnOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "cursedSubjects", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	rmnRemoteCid := extractCreatedContractId(res)
	t.Logf("Deployed RMNRemote: %s", rmnRemoteCid)

	// Deploy CommitteeVerifier
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccv-e2e"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "versionTag", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: versionTag}}},
						{Label: "messageSentObserver", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "storageLocation", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ipfs://test-receive"}}},
						{Label: "threshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2}}},
						{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[0]}},
							{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[1]}},
							{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[2]}},
						}}}}},
						{Label: "rmnRemoteInstanceAddress", Value: rawInstanceAddress("test-rmn-e2e@" + partyCCIP)},
						{Label: "remoteChainFeeConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	ccvCid := extractCreatedContractId(res)
	t.Logf("Deployed CommitteeVerifier: %s (ccvId: %s)", ccvCid, ccvID)

	// Deploy GlobalConfig with source chain config including the CCV
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-globalconfig-e2e"}}},
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
										rawInstanceAddress(ccvID),
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
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-tar-e2e"}}},
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

	// Deploy OffRamp with instance addresses for GlobalConfig, RMNRemote, TokenAdminRegistry
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-offramp-e2e"}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "globalConfigInstanceAddress", Value: rawInstanceAddress("test-globalconfig-e2e@" + partyCCIP)},
						{Label: "rmnRemoteInstanceAddress", Value: rawInstanceAddress("test-rmn-e2e@" + partyCCIP)},
						{Label: "tokenAdminRegistryInstanceAddress", Value: rawInstanceAddress("test-tar-e2e@" + partyCCIP)},
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
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-factory-e2e"}}},
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
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-e2e"}}},
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

	// Build message (no token transfer, just payload data)
	testPayload := []byte("Hello CCIP - this is a test message payload!")
	msg := &MessageV1{
		SourceChainSelector: 123,
		DestChainSelector:   456,
		SequenceNumber:      1,
		ExecutionGasLimit:   200000,
		CCIPReceiveGasLimit: 100000,
		Finality:            2000,
		CCVAndExecutorHash:  [32]byte{},
		OnRampAddress:       hexToBytes("0000000000000000000000000000000000000001"),
		OffRampAddress:      hexToBytes("0000000000000000000000000000000000000002"),
		Sender:              hexToBytes("0000000000000000000000000000000000000003"),
		Receiver:            EncodePartyID(partyReceiver),
		DestBlob:            []byte{},
		TokenTransfer:       nil,
		MessageData:         testPayload,
	}
	encodedMessage, err := EncodeMessageV1(msg)
	require.NoError(t, err)
	encodedMessageHex := hex.EncodeToString(encodedMessage)
	messageHash := crypto.Keccak256(encodedMessage)
	messageHashHex := hex.EncodeToString(messageHash)
	t.Logf("Message hash: %s", messageHashHex)
	t.Logf("Message payload: %s", string(testPayload))

	// Generate verifierResults with 2 of 3 signatures
	verifierResults, err := GenerateVerifierResults(encodedMessage, ccvSignerKeys[:2])
	require.NoError(t, err)
	verifierResultsHex := hex.EncodeToString(verifierResults)

	// Deploy CCIPReceiver for receiver
	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "requiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyReceiver},
		},
	})
	require.NoError(t, err)
	ccipReceiverCid := extractCreatedContractId(res)
	t.Logf("Deployed CCIPReceiver: %s", ccipReceiverCid)

	// Get disclosures for CCIPReceiver.Execute
	disclosedCCIPReceiver, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver",
	})
	require.NoError(t, err)
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
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)
	disclosedRmnRemote, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote",
	})
	require.NoError(t, err)
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err)

	// CCIPReceiver.Execute: PrepareExecute + CCV verification + Execute in one transaction
	res, err = receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId: ccipReceiverCid,
					Choice:     "Execute",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCid}}},
						{Label: "offRampCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: offRampCid}}},
						{Label: "globalConfigCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: globalConfigCid}}},
						{Label: "tokenAdminRegistryCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: tokenAdminRegistryCid}}},
						{Label: "rmnRemoteCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: rmnRemoteCid}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: encodedMessageHex}}},
						{Label: "tokenTransfer", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}},
						{Label: "ccvInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "ccvCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedCCV.ContractId}}},
								{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: verifierResultsHex}}},
							}}}},
						}}}}},
						{Label: "additionalRequiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedCCIPReceiver, disclosedRouter, disclosedOffRamp, disclosedGlobalConfig, disclosedTar, disclosedRmnRemote, disclosedCCV},
		},
	})
	require.NoError(t, err)

	// Extract messageId from CCIPMessageReceived to verify success
	var returnedMessageId string
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageReceived" {
				// CCIPMessageReceived: owner, router, messageId, message, tokenReleaseResult
				returnedMessageId = e.Created.GetCreateArguments().GetFields()[2].GetValue().GetText()
				break
			}
		}
	}
	require.NotEmpty(t, returnedMessageId, "CCIPMessageReceived should be created")

	require.Equal(t, messageHashHex, returnedMessageId, "CCIPMessageReceived messageId should match")

	t.Logf("Execute completed")
	t.Logf("  Message ID: %s", returnedMessageId)
	t.Logf("  Original payload: %s", string(testPayload))
}
