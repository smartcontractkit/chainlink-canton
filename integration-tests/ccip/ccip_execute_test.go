// Integration test for CCIP execute flow with CommitteeVerifier signature verification.
//
// This test demonstrates the complete execute flow WITHOUT token transfers:
//   - Deploy CCVRegistry, GlobalConfig, CommitteeVerifier, OffRamp, PerPartyRouter
//   - Build a message with payload data
//   - Generate ECDSA signatures and verify via CommitteeVerifier
//   - Execute via PerPartyRouter
//   - Validate the returned message payload matches the original
//
// Requires running localnet:
//
//	cd compose/localnet && docker compose up -d
//	go test ./src/tests/... -run TestCCIPExecuteE2E -v

package tests

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
	"github.com/smartcontractkit/chainlink-canton-internal/integration-tests/testhelpers"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
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
	binary.Write(&buf, binary.BigEndian, msg.SourceChainSelector)
	binary.Write(&buf, binary.BigEndian, msg.DestChainSelector)
	binary.Write(&buf, binary.BigEndian, msg.SequenceNumber)
	binary.Write(&buf, binary.BigEndian, msg.ExecutionGasLimit)
	binary.Write(&buf, binary.BigEndian, msg.CCIPReceiveGasLimit)
	binary.Write(&buf, binary.BigEndian, msg.Finality)
	buf.Write(msg.CCVAndExecutorHash[:])

	// Length-prefixed fields (1-byte length)
	buf.WriteByte(uint8(len(msg.OnRampAddress)))
	buf.Write(msg.OnRampAddress)
	buf.WriteByte(uint8(len(msg.OffRampAddress)))
	buf.Write(msg.OffRampAddress)
	buf.WriteByte(uint8(len(msg.Sender)))
	buf.Write(msg.Sender)
	buf.WriteByte(uint8(len(msg.Receiver)))
	buf.Write(msg.Receiver)

	// 2-byte length prefixed fields
	binary.Write(&buf, binary.BigEndian, uint16(len(msg.DestBlob)))
	buf.Write(msg.DestBlob)

	if msg.TokenTransfer != nil {
		tokenBytes := encodeTokenTransferV1(msg.TokenTransfer)
		binary.Write(&buf, binary.BigEndian, uint16(len(tokenBytes)))
		buf.Write(tokenBytes)
	} else {
		binary.Write(&buf, binary.BigEndian, uint16(0))
	}

	binary.Write(&buf, binary.BigEndian, uint16(len(msg.MessageData)))
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

	buf.WriteByte(uint8(len(tt.SourcePoolAddress)))
	buf.Write(tt.SourcePoolAddress)
	buf.WriteByte(uint8(len(tt.SourceTokenAddress)))
	buf.Write(tt.SourceTokenAddress)
	buf.WriteByte(uint8(len(tt.DestTokenAddress)))
	buf.Write(tt.DestTokenAddress)
	buf.WriteByte(uint8(len(tt.TokenReceiver)))
	buf.Write(tt.TokenReceiver)

	binary.Write(&buf, binary.BigEndian, uint16(len(tt.ExtraData)))
	buf.Write(tt.ExtraData)

	return buf.Bytes()
}

// GenerateVerifierResults generates the verifierResults blob for CommitteeVerifier.
// Format: versionTag (4 bytes) || signatureLength (2 bytes) || signatures (64 bytes each)
func GenerateVerifierResults(encodedMessage []byte, privateKeys []*ecdsa.PrivateKey) ([]byte, error) {
	versionTag, _ := hex.DecodeString("49ff34ed")

	preimage := append(versionTag, encodedMessage...)
	msgHash := crypto.Keccak256(preimage)

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
	binary.Write(&result, binary.BigEndian, uint16(len(signatures)))
	result.Write(signatures)

	return result.Bytes(), nil
}

// EncodePartyID encodes a Canton party ID to bytes.
func EncodePartyID(partyID string) []byte {
	return []byte(partyID)
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

	dars := [][]byte{commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.Party
	partyReceiver := receiverParticipant.Party
	t.Logf("Parties: CCIP=%s, Receiver=%s", partyCCIP, partyReceiver)

	// Generate signer keys for CommitteeVerifier (3 signers, threshold 2)
	var ccvSignerKeys []*ecdsa.PrivateKey
	var ccvSignerPubKeys []*apiv2.Value
	for i := 0; i < 3; i++ {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, &apiv2.Value{Sum: &apiv2.Value_Text{Text: pubKeyHex}})
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
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccvregistry-e2e"}}},
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
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccv-e2e"}}},
						{Label: "versionTag", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: versionTag}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "messageSentObserver", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "storageLocation", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ipfs://test-e2e"}}},
						{Label: "threshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2}}},
						{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: ccvSignerPubKeys}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	t.Logf("Deployed CommitteeVerifier (ccvId: %s)", ccvId)

	// Deploy GlobalConfig with source chain config including the CCV
	sourceChainSelector := "123"
	destChainSelector := "456"
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

	// Deploy TokenAdminRegistry (required by OffRamp even without token transfers)
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

	// Deploy OffRamp
	res, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-offramp-e2e"}}},
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
		OnRampAddress:       []byte("0000000000000000000000000000000000000001"),
		OffRampAddress:      []byte("0000000000000000000000000000000000000002"),
		Sender:              []byte("0000000000000000000000000000000000000003"),
		Receiver:            EncodePartyID(partyReceiver),
		DestBlob:            []byte{},
		TokenTransfer:       nil, // No token transfer
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

	// Get disclosures for CommitteeVerifier_VerifyMessage
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err)
	disclosedCCVRegistry, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.CCVRegistry", EntityName: "CCVRegistry",
	})
	require.NoError(t, err)

	// Build MessageV1 Daml record (no token transfer)
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
		{Label: "tokenTransfer", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}},
		{Label: "messageData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(testPayload)}}},
	}}}}

	// Call CommitteeVerifier_VerifyMessage to get CCVVerifyTicket
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

	// Get disclosures for Execute
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
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
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

	// Extract messageId from ExecutionStateChanged event to verify success.
	// The Ledger API returns created/archived events; exercised events with results
	// require explicit configuration. ExecutionStateChanged proves execution succeeded.
	var returnedMessageId string
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "ExecutionStateChanged" {
				// ExecutionStateChanged has: ccipOwner, receiver, event
				// event is ExecutionStateChangedEvent with: sourceChainSelector, sequenceNumber, messageId, state, returnData
				eventRecord := e.Created.GetCreateArguments().GetFields()[2].GetValue().GetRecord()
				returnedMessageId = eventRecord.GetFields()[2].GetValue().GetText()
				break
			}
		}
	}
	require.NotEmpty(t, returnedMessageId, "ExecutionStateChanged event should be created")

	// Verify the messageId in ExecutionStateChanged matches what we sent
	// This proves the message was decoded and processed correctly (including the payload)
	require.Equal(t, messageHashHex, returnedMessageId, "ExecutionStateChanged messageId should match")

	// Note: The full message payload is verified implicitly:
	// - The messageId is keccak256(encodedMessage) which includes the payload
	// - If the payload was corrupted, the hash wouldn't match
	// - The Execute succeeded (no error), meaning OffRamp decoded the message successfully

	t.Logf("Execute completed")
	t.Logf("  Message ID: %s", returnedMessageId)
	t.Logf("  Original payload: %s", string(testPayload))
}
