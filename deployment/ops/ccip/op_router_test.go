package ccip

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/common"
	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

func TestRouterOperations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	setupResult, err := compileClient.Setup(ctx, compileClient.Config{
		LedgerAPIURL:      "participant1.grpc-ledger-api.localhost:8080",
		AdminAPIURL:       "participant1.admin-api.localhost:8080",
		JWTSecret:         "unsafe",
		DeployerParty:     "", // Empty to use primary party or allocate new one
		DeployerPartyHint: "ledger-api-user",
	})
	require.NoError(t, err, "Failed to setup Canton client")

	t.Cleanup(setupResult.BindingClient.Close)

	deps := CantonOpDeps{
		BindingClient: setupResult.BindingClient,
		Party:         setupResult.Party,
		UserID:        setupResult.UserID,
	}

	reporter := cld_ops.NewMemoryReporter()

	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	instanceID := "test-router-instance"
	chainSelectorValue := "1111111111"
	destChainSelectorValue := "2222222222"
	onRampAddress := "0000000000000000000000000000000000000001"

	// --------------------------
	// Deploy contracts first (required for router operations)
	// --------------------------
	var globalConfigContractID string
	var globalConfigTemplateID string
	var onRampContractID string
	var offRampContractID string
	var tokenAdminRegistryContractID string
	var feeQuoterContractID string
	var perPartyRouterContractID string
	var perPartyRouterTemplateID string
	var committeeVerifierContractID string
	var committeeVerifierTemplateID string
	var ccvRegistryContractID string

	// Get CCIPSend tickets from CCV
	var ccvTickets []string
	t.Run("DeployContracts", func(t *testing.T) {
		t.Parallel()

		// Deploy GlobalConfig
		commonResult, err := cld_ops.ExecuteOperation(bundle, DeployCCIPCommonOp, deps, DeployCCIPCommonInput{
			InstanceID:         instanceID,
			ChainSelectorValue: chainSelectorValue,
			OnRampAddress:      onRampAddress,
		})
		require.NoError(t, err, "failed to deploy CCIP Common")
		globalConfigContractID = commonResult.Output.Output.GlobalConfigContractID
		globalConfigTemplateID = commonResult.Output.Output.GlobalConfigTemplateID
		require.NotEmpty(t, globalConfigContractID, "GlobalConfig contract ID should not be empty")
		t.Logf("Deployed GlobalConfig contract ID: %s", globalConfigContractID)

		// Deploy TokenAdminRegistry
		tarResult, err := cld_ops.ExecuteOperation(bundle, DeployTokenAdminRegistryOp, deps, DeployTokenAdminRegistryInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy TokenAdminRegistry")
		tokenAdminRegistryContractID = tarResult.Output.Output.TokenAdminRegistryContractID
		require.NotEmpty(t, tokenAdminRegistryContractID, "TokenAdminRegistry contract ID should not be empty")
		t.Logf("Deployed TokenAdminRegistry contract ID: %s", tokenAdminRegistryContractID)

		// Deploy OnRamp
		onRampResult, err := cld_ops.ExecuteOperation(bundle, DeployOnRampOp, deps, DeployOnRampInput{
			InstanceID:           instanceID,
			DestChainSelector:    destChainSelectorValue,
			DestChainOnRampBytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		})
		require.NoError(t, err, "failed to deploy OnRamp")
		onRampContractID = onRampResult.Output.Output.OnRampContractID
		require.NotEmpty(t, onRampContractID, "OnRamp contract ID should not be empty")
		t.Logf("Deployed OnRamp contract ID: %s", onRampContractID)

		// Deploy OffRamp
		offRampResult, err := cld_ops.ExecuteOperation(bundle, DeployOffRampOp, deps, DeployOffRampInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy OffRamp")
		offRampContractID = offRampResult.Output.Output.OffRampContractID
		require.NotEmpty(t, offRampContractID, "OffRamp contract ID should not be empty")
		t.Logf("Deployed OffRamp contract ID: %s", offRampContractID)

		// Deploy FeeQuoter
		feeQuoterResult, err := cld_ops.ExecuteOperation(bundle, DeployFeeQuoterOp, deps, DeployFeeQuoterInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy FeeQuoter")
		feeQuoterContractID = feeQuoterResult.Output.Output.FeeQuoterContractID
		require.NotEmpty(t, feeQuoterContractID, "FeeQuoter contract ID should not be empty")
		t.Logf("Deployed FeeQuoter contract ID: %s", feeQuoterContractID)

		// Deploy PerPartyRouter
		perPartyRouterResult, err := cld_ops.ExecuteOperation(bundle, DeployPerPartyRouterOp, deps, DeployPerPartyRouterInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy PerPartyRouter")
		perPartyRouterContractID = perPartyRouterResult.Output.Output.PerPartyRouterContractID
		perPartyRouterTemplateID = perPartyRouterResult.Output.Output.PerPartyRouterTemplateID
		require.NotEmpty(t, perPartyRouterContractID, "PerPartyRouter contract ID should not be empty")
		t.Logf("Deployed PerPartyRouter contract ID: %s", perPartyRouterContractID)

		// deploy CommitteeVerifier
		committeeVerifierResult, err := cld_ops.ExecuteOperation(bundle, DeployCommitteeVerifierOp, deps, DeployCommitteeVerifierInput{
			InstanceID:          instanceID,
			VersionTag:          "49ff34ed",
			StorageLocation:     "ipfs://test-ccv",
			Threshold:           2,
			Signers:             []string{"signer1", "signer2", "signer3"},
			MessageSentObserver: "", // Will default to deployer party
		})
		require.NoError(t, err, "failed to deploy CommitteeVerifier")
		committeeVerifierContractID = committeeVerifierResult.Output.Output.CommitteeVerifierContractID
		committeeVerifierTemplateID = committeeVerifierResult.Output.Output.CommitteeVerifierTemplateID
		require.NotEmpty(t, committeeVerifierContractID, "CommitteeVerifier contract ID should not be empty")
		require.NotEmpty(t, committeeVerifierTemplateID, "CommitteeVerifier template ID should not be empty")
		t.Logf("Deployed CommitteeVerifier contract ID: %s", committeeVerifierContractID)

		// Deploy CCVRegistry
		ccvRegistryResult, err := cld_ops.ExecuteOperation(bundle, DeployCCVRegistryOp, deps, DeployCCVRegistryInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy CCVRegistry")
		ccvRegistryContractID = ccvRegistryResult.Output.Output.CCVRegistryContractID
		require.NotEmpty(t, ccvRegistryContractID, "CCVRegistry contract ID should not be empty")
		t.Logf("Deployed CCVRegistry contract ID: %s", ccvRegistryContractID)

		// -------------------------- CONFIGURE CCIP CONTRACTS --------------------------
		// update global config dest chain config
		destChainConfigResult, err := cld_ops.ExecuteOperation(bundle, UpdateGlobalConfigDestChainConfigOp, deps, UpdateGlobalConfigDestChainConfigInput{
			GlobalConfigContractID: globalConfigContractID,
			GlobalConfigTemplateID: globalConfigTemplateID,
			DestChainSelector:      destChainSelectorValue,
			Config: common.DestChainConfig{
				IsEnabled:        types.BOOL(true),
				DefaultExecutor:  types.TEXT("executor-party"),
				OffRampAddress:   types.TEXT("0000000000000000000000000000000000000002"),
				LaneMandatedCCVs: []types.TEXT{},
				DefaultCCVs:      []types.TEXT{},
			},
		})
		require.NoError(t, err, "failed to update GlobalConfig dest chain config")
		require.NotEmpty(t, destChainConfigResult.Output.Output.TransactionID, "Transaction ID should not be empty")
		t.Logf("Updated GlobalConfig dest chain config, transaction ID: %s", destChainConfigResult.Output.Output.TransactionID)

		// carry forward the new CID returned by the update op
		require.NotEmpty(t, destChainConfigResult.Output.Output.NewGlobalConfigContractID, "new GlobalConfig contract ID should not be empty")
		globalConfigContractID = destChainConfigResult.Output.Output.NewGlobalConfigContractID
		globalConfigTemplateID = destChainConfigResult.Output.Output.NewGlobalConfigTemplateID
		t.Logf("Updated GlobalConfig dest chain config, new contract ID: %s", globalConfigContractID)

		// update global config source chain config
		sourceChainConfigResult, err := cld_ops.ExecuteOperation(bundle, UpdateGlobalConfigSourceChainConfigOp, deps, UpdateGlobalConfigSourceChainConfigInput{
			GlobalConfigContractID: globalConfigContractID, // fresh created contractID
			GlobalConfigTemplateID: globalConfigTemplateID,
			SourceChainSelector:    chainSelectorValue, // Use chainSelectorValue instead of undefined evmChainSelector
			Config: common.SourceChainConfig{
				IsEnabled:        types.BOOL(true),
				OnRampAddress:    types.TEXT(onRampAddress),
				LaneMandatedCCVs: []types.TEXT{},
				DefaultCCVs:      []types.TEXT{},
			},
		})
		require.NoError(t, err, "failed to update GlobalConfig source chain config")
		require.NotEmpty(t, sourceChainConfigResult.Output.Output.TransactionID, "Transaction ID should not be empty")
		t.Logf("Updated GlobalConfig source chain config, transaction ID: %s", sourceChainConfigResult.Output.Output.TransactionID)

		// carry forward the new CID returned by the update op
		require.NotEmpty(t, sourceChainConfigResult.Output.Output.NewGlobalConfigContractID, "new GlobalConfig contract ID should not be empty")
		globalConfigContractID = sourceChainConfigResult.Output.Output.NewGlobalConfigContractID
		globalConfigTemplateID = sourceChainConfigResult.Output.Output.NewGlobalConfigTemplateID
		t.Logf("Updated GlobalConfig source chain config, new contract ID: %s", globalConfigContractID)

		// Convert MessageV1Input to ccvs.MessageV1
		scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)

		// Parse source chain selector
		sourceChainSelector, ok := new(big.Int).SetString(chainSelectorValue, 10)
		require.True(t, ok, "failed to parse source chain selector")

		sourceChainSelectorMantissa := new(big.Int).Mul(sourceChainSelector, scale10)

		// Parse dest chain selector
		destChainSelector, ok := new(big.Int).SetString(destChainSelectorValue, 10)
		require.True(t, ok, "failed to parse dest chain selector")

		destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

		// Parse sequence number
		sequenceNumber := new(big.Int).SetInt64(1)
		sequenceNumberMantissa := new(big.Int).Mul(sequenceNumber, scale10)

		require.NotEmpty(t, committeeVerifierContractID, "committeeVerifierContractID not set")
		require.NotEmpty(t, ccvRegistryContractID, "ccvRegistryContractID not set")

		// cerate ticket for CrossChainVerifier_ForwardToVerifier using CommitteeVerifierForwardToVerifierOp
		committeeVerifierForwardToVerifierResult, err := cld_ops.ExecuteOperation(bundle, CommitteeVerifierForwardToVerifierOp, deps, CommitteeVerifierForwardToVerifierInput{
			CommitteeVerifierContractID: committeeVerifierContractID,
			CommitteeVerifierTemplateID: committeeVerifierTemplateID,
			CcvRegistryCid:              ccvRegistryContractID,
			Message: ccvs.MessageV1{
				SourceChainSelector: types.NUMERIC(sourceChainSelectorMantissa),
				DestChainSelector:   types.NUMERIC(destChainSelectorMantissa),
				SequenceNumber:      types.NUMERIC(sequenceNumberMantissa),
				ExecutionGasLimit:   types.INT64(200000),
				CcipReceiveGasLimit: types.INT64(100000),
				Finality:            types.INT64(12),
				CcvAndExecutorHash:  types.TEXT("0x1234567890abcdef"),
				OnRampAddress:       types.TEXT("0000000000000000000000000000000000000001"),
				OffRampAddress:      types.TEXT("0000000000000000000000000000000000000002"),
				Sender:              types.TEXT("0000000000000000000000000000000000000003"),
				Receiver:            types.TEXT("0000000000000000000000000000000000000004"),
				DestBlob:            types.TEXT("0x"),
				MessageData:         types.TEXT("0xdeadbeef"),
				TokenTransfer:       nil, // Optional, nil for this test
			},
			MessageId: "0x01",
			FeeToken: ccvs.InstrumentId{
				Admin: types.PARTY(deps.Party),
				Id:    types.TEXT("test-token"),
			},
			FeeTokenAmount: "0", // Zero fee for test
			VerifierArgs:   "0x00",
		})
		require.NoError(t, err, "failed to create CommitteeVerifierForwardToVerifier")
		require.NotEmpty(t, committeeVerifierForwardToVerifierResult.Output.Output.TransactionID, "Transaction ID should not be empty")
		require.NotEmpty(t, committeeVerifierForwardToVerifierResult.Output.Output.CCVTicketContractID, "CCVTicket contract ID should not be empty")
		require.NotEmpty(t, committeeVerifierForwardToVerifierResult.Output.Output.CCVTicketTemplateID, "CCVTicket template ID should not be empty")
		t.Logf("Created CommitteeVerifierForwardToVerifier contract ID: %s", committeeVerifierForwardToVerifierResult.Output.Output.CCVTicketContractID)
		t.Logf("Created CommitteeVerifierForwardToVerifier template ID: %s", committeeVerifierForwardToVerifierResult.Output.Output.CCVTicketTemplateID)
		ccvTickets = append(ccvTickets, committeeVerifierForwardToVerifierResult.Output.Output.CCVTicketContractID)
	})

	// --------------------------
	// Test Router CCIPSend
	// --------------------------
	t.Run("RouterCCIPSend", func(t *testing.T) {
		t.Parallel()

		// Create test payload
		testPayload := "48656c6c6f2043434950" // "Hello CCIP" in hex

		result, err := cld_ops.ExecuteOperation(bundle, RouterCCIPSendOp, deps, RouterCCIPSendInput{
			PerPartyRouterContractID: perPartyRouterContractID,
			PerPartyRouterTemplateID: perPartyRouterTemplateID,
			OnRampCid:                onRampContractID,
			GlobalConfigCid:          globalConfigContractID,
			TokenAdminRegistryCid:    tokenAdminRegistryContractID,
			DestChainSelector:        destChainSelectorValue,
			Receiver:                 "0000000000000000000000000000000000000004",
			Payload:                  testPayload,
			ExecutionGasLimit:        200000,
			CcipReceiveGasLimit:      100000,
			TokenSendTicket:          nil, // No token transfer for this test
			// Note this works even though we don't have any mandated CCV in this lane, we can also set this to empty string array since there is no manadated ccv on this lane
			CcvTickets: ccvTickets,
		})
		require.NoError(t, err, "failed to execute RouterCCIPSend")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CcipMessageSentCID, "CCIPMessageSent contract ID should not be empty")
		t.Logf("RouterCCIPSend completed, transaction ID: %s", result.Output.Output.TransactionID)
		t.Logf("Created CCIPMessageSent contract ID: %s", result.Output.Output.CcipMessageSentCID)
	})

	// --------------------------
	// Test Router Execute
	// --------------------------
	// t.Run("RouterExecute", func(t *testing.T) {
	// 	// For Execute, we need an encoded message
	// 	// In a real scenario, this would come from a CCIPSend operation
	// 	// For testing, we'll create a minimal encoded message
	// 	// Note: This is a simplified test - in practice, you'd need a properly encoded MessageV1
	// 	encodedMessage := "0x" // Empty encoded message for test

	// 	// Execute the message
	// 	result, err := cld_ops.ExecuteOperation(bundle, RouterExecuteOp, deps, RouterExecuteInput{
	// 		PerPartyRouterContractID: perPartyRouterContractID,
	// 		PerPartyRouterTemplateID: perPartyRouterTemplateID,
	// 		OffRampCid:               offRampContractID,
	// 		GlobalConfigCid:          globalConfigContractID,
	// 		TokenAdminRegistryCid:    tokenAdminRegistryContractID,
	// 		EncodedMessage:           encodedMessage,
	// 		CcvVerifyTickets:         []string{}, // No CCV verify tickets for this test
	// 		TokenPoolCCVTicket:       nil,        // No token pool CCV ticket for this test
	// 		ReceiverRequiredCCVIds:   []string{}, // No required CCV IDs for this test
	// 	})
	// 	// Note: This test may fail if the encoded message is invalid or if required tickets are missing
	// 	// In a real scenario, you would need proper CCV verify tickets and a valid encoded message
	// 	if err != nil {
	// 		t.Logf("RouterExecute failed (expected if encoded message is invalid or tickets are missing): %v", err)
	// 		// This is expected for a minimal test without proper message encoding and tickets
	// 		return
	// 	}
	// 	require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")
	// 	t.Logf("RouterExecute completed, transaction ID: %s", result.Output.Output.TransactionID)
	// 	if result.Output.Output.TokenReceiveTicketID != nil {
	// 		t.Logf("Created TokenReceiveTicket contract ID: %s", *result.Output.Output.TokenReceiveTicketID)
	// 	}
	// })
}
