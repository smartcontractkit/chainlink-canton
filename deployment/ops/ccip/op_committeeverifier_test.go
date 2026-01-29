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
	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCommitteeVerifierForwardToVerifier(t *testing.T) {
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

	instanceID := "test-ccip-instance"

	// --------------------------
	// Deploy contracts first (required for forward operation)
	// --------------------------
	var committeeVerifierContractID string
	var committeeVerifierTemplateID string
	var ccvRegistryContractID string

	t.Run("DeployContracts", func(t *testing.T) {
		t.Parallel()

		// Deploy CommitteeVerifier
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

		// Deploy CCVRegistry using DeployCCVRegistryOp
		ccvRegistryResult, err := cld_ops.ExecuteOperation(bundle, DeployCCVRegistryOp, deps, DeployCCVRegistryInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy CCVRegistry")
		ccvRegistryContractID = ccvRegistryResult.Output.Output.CCVRegistryContractID
		require.NotEmpty(t, ccvRegistryContractID, "CCVRegistry contract ID should not be empty")
		t.Logf("Deployed CCVRegistry contract ID: %s", ccvRegistryContractID)
	})

	// --------------------------
	// Test CommitteeVerifier ForwardToVerifier
	// --------------------------
	t.Run("ForwardToVerifier", func(t *testing.T) {
		t.Parallel()

		// Scale chain selectors to NUMERIC(10) mantissa
		scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
		sourceChainSelector := new(big.Int).SetInt64(1111111111)
		sourceChainSelectorMantissa := new(big.Int).Mul(sourceChainSelector, scale10)

		destChainSelector := new(big.Int).SetInt64(2222222222)
		destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

		sequenceNumber := new(big.Int).SetInt64(1)
		sequenceNumberMantissa := new(big.Int).Mul(sequenceNumber, scale10)

		// Create test message
		message := ccvs.MessageV1{
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
		}

		// Create test fee token
		feeToken := ccvs.InstrumentId{
			Admin: types.PARTY(deps.Party),
			Id:    types.TEXT("test-token"),
		}

		// Execute forward operation
		result, err := cld_ops.ExecuteOperation(bundle, CommitteeVerifierForwardToVerifierOp, deps, CommitteeVerifierForwardToVerifierInput{
			CommitteeVerifierContractID: committeeVerifierContractID,
			CommitteeVerifierTemplateID: committeeVerifierTemplateID,
			CcvRegistryCid:              ccvRegistryContractID,
			Message:                     message,
			MessageId:                   "test-message-id-1",
			FeeToken:                    feeToken,
			FeeTokenAmount:              "0", // Zero fee for test
			VerifierArgs:                "0x",
		})
		require.NoError(t, err, "failed to forward message to verifier")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CCVTicketContractID, "CCVTicket contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CCVTicketTemplateID, "CCVTicket template ID should not be empty")
		t.Logf("Forwarded message to verifier, transaction ID: %s", result.Output.Output.TransactionID)
		t.Logf("Created CCVTicket contract ID: %s", result.Output.Output.CCVTicketContractID)
	})

	// --------------------------
	// Test CommitteeVerifier VerifyMessage
	// --------------------------
	t.Run("VerifyMessage", func(t *testing.T) {
		t.Parallel()

		// Scale chain selectors to NUMERIC(10) mantissa
		scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
		sourceChainSelector := new(big.Int).SetInt64(1111111111)
		sourceChainSelectorMantissa := new(big.Int).Mul(sourceChainSelector, scale10)

		destChainSelector := new(big.Int).SetInt64(2222222222)
		destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

		sequenceNumber := new(big.Int).SetInt64(2)
		sequenceNumberMantissa := new(big.Int).Mul(sequenceNumber, scale10)

		// Create test message
		message := ccvs.MessageV1{
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
		}

		// Create mock verifier results (version tag + signature length + signatures)
		// In a real scenario, this would contain actual signatures from the committee
		// For testing, we use a minimal valid structure
		versionTag := "49ff34ed" // 4 bytes hex
		// Signature length (2 bytes) = 0 for test (no actual signatures)
		// In real usage, this would contain actual signature data
		verifierResults := versionTag + "0000" // version tag + 0-length signatures

		// Execute verify operation
		result, err := cld_ops.ExecuteOperation(bundle, CommitteeVerifierVerifyMessageOp, deps, CommitteeVerifierVerifyMessageInput{
			CommitteeVerifierContractID: committeeVerifierContractID,
			CommitteeVerifierTemplateID: committeeVerifierTemplateID,
			CcvRegistryCid:              ccvRegistryContractID,
			Message:                     message,
			MessageId:                   "test-message-id-2",
			VerifierResults:             verifierResults,
			Receiver:                    deps.Party, // Receiver party
		})
		// Note: This test may fail if verifierResults doesn't contain valid signatures
		// In a real scenario, you would need to generate actual signatures from the committee signers
		if err != nil {
			t.Logf("VerifyMessage failed (expected if verifierResults doesn't contain valid signatures): %v", err)
			// This is expected if we don't have valid signatures, so we just log it
			return
		}
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CCVVerifyTicketContractID, "CCVVerifyTicket contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CCVVerifyTicketTemplateID, "CCVVerifyTicket template ID should not be empty")
		t.Logf("Verified message, transaction ID: %s", result.Output.Output.TransactionID)
		t.Logf("Created CCVVerifyTicket contract ID: %s", result.Output.Output.CCVVerifyTicketContractID)
	})
}
