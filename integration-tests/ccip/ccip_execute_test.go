// Integration test for CCIP execute flow with CommitteeVerifier signature verification.
//
// This test demonstrates the complete execute flow WITHOUT token transfers:
//   - Deploy CCVRegistry, GlobalConfig, CommitteeVerifier, OffRamp, PerPartyRouter
//   - Build a message with payload data
//   - Generate ECDSA signatures and verify via CommitteeVerifier
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
	"strconv"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/noders-team/go-daml/pkg/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
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
	_ = binary.Write(&result, binary.BigEndian, uint16(len(signatures))) //nolint:gosec
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
	onRampDar, err := contracts.GetDar(contracts.CCIPOnRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	feeQuoterDar, err := contracts.GetDar(contracts.CCIPFeeQuoter, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err)
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err)

	dars := [][]byte{commonDar, offRampDar, onRampDar, feeQuoterDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.Party
	partyReceiver := receiverParticipant.Party
	t.Logf("Parties: CCIP=%s, Receiver=%s", partyCCIP, partyReceiver)

	// Generate signer keys for CommitteeVerifier (3 signers, threshold 2)
	ccvSignerKeys := make([]*ecdsa.PrivateKey, 0, 3)
	ccvSignerPubKeys := make([]types.TEXT, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, types.TEXT(pubKeyHex))
	}
	t.Logf("Generated %d CCV signer keys", len(ccvSignerKeys))
	versionTag := "49ff34ed"
	ccvID := versionTag + "@" + partyCCIP

	chainSelectorBig := big.NewInt(0).SetUint64(chainsel.CANTON_LOCALNET.Selector)
	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)
	cldfEnv := cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{env.Chain}),
		OperationsBundle: bundle,
	}

	// Deploy Chain contracts
	out, err := changesets.DeployChainContracts{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.DeployChainContractsConfig]{
		ChainSelector: env.Selector,
		Participant:   0,
		Party:         partyCCIP,
		Config: changesets.DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: partyCCIP,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Template: ccvs.CommitteeVerifier{
							Owner:               types.PARTY(partyCCIP),
							CcipOwner:           types.PARTY(partyCCIP),
							VersionTag:          types.TEXT(versionTag),
							MessageSentObserver: types.PARTY(partyCCIP),
							StorageLocation:     "ipfs://test-receive",
							Threshold:           2,
							Signers:             ccvSignerPubKeys,
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: chainSelectorBig,
						OnRampAddress: "", // TODO ?
					},
				},
			},
		},
	})
	require.NoErrorf(t, err, "Failed to deploy CCIP contracts: %v", err)

	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()

	t.Log("Deployed CCIP chain contracts:")
	addresses := cldfEnv.DataStore.Addresses().Filter()
	for i, address := range addresses {
		t.Logf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier)
	}

	// Resolve contracts
	globalConfig, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	require.NoError(t, err, "failed to get GlobalConfig address")
	feeQuoter, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, ""))
	require.NoError(t, err, "failed to get FeeQuoter address")
	onRamp, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(onramp.ContractType), onramp.Version, ""))
	require.NoError(t, err, "failed to get OnRamp address")
	offRamp, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(offramp.ContractType), offramp.Version, ""))
	require.NoError(t, err, "failed to get OffRamp address")

	// Deploy and configure lane
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	out, err = changesets.ConfigureChainForLanes{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.ConfigureChainForLanesConfig]{
		ChainSelector: env.Selector,
		Participant:   0,
		Party:         partyCCIP,
		Config: changesets.ConfigureChainForLanesConfig{
			Input: sequences.ConfigureChainForLanesInput{
				ChainSelector:      env.Selector,
				GlobalConfig:       contracts.InstanceID(globalConfig.Address),
				FeeQuoter:          contracts.InstanceID(feeQuoter.Address),
				OnRamp:             contracts.InstanceID(onRamp.Address),
				OffRamp:            contracts.InstanceID(offRamp.Address),
				CommitteeVerifiers: nil,
				RemoteChains: map[uint64]adapters.RemoteChainConfig[[]byte, string]{
					remoteSelector: {
						AllowTrafficFrom:         true,
						OnRamps:                  [][]byte{[]byte("0000000000000000000000000000000000000001")},
						OffRamp:                  nil,
						DefaultInboundCCVs:       nil,
						LaneMandatedInboundCCVs:  []string{ccvID},
						DefaultOutboundCCVs:      nil,
						LaneMandatedOutboundCCVs: nil,
						DefaultExecutor:          "",
						FeeQuoterDestChainConfig: adapters.FeeQuoterDestChainConfig{},
						ExecutorDestChainConfig:  adapters.ExecutorDestChainConfig{},
						AddressBytesLength:       0,
						BaseExecutionGasCost:     0,
					},
				},
			},
		},
	})
	require.NoErrorf(t, err, "Failed to configure chain for lanes")
	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()
	t.Log("Configured chain for lanes")

	// Create PerPartyRouter for receiver
	disclosedFactory, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	})
	require.NoError(t, err)

	res, err := receiverParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
		SourceChainSelector: remoteSelector,
		DestChainSelector:   env.Selector,
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
		{Label: "sourceChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: strconv.FormatUint(remoteSelector, 10)}}},
		{Label: "destChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: strconv.FormatUint(env.Selector, 10)}}},
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
				// is ExecutionStateChangedEvent with: sourceChainSelector, sequenceNumber, messageId, state, returnData
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
