package ccip

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	compileClient "github.com/smartcontractkit/chainlink-canton/deployment/client"
)

func TestDeployCCIPContracts(t *testing.T) {
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
	}

	reporter := cld_ops.NewMemoryReporter()

	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	instanceID := "test-ccip-instance"
	chainSelectorValue := "1111111111"
	onRampAddress := "0000000000000000000000000000000000000001"

	// --------------------------
	// Test CCIP Common Deployment
	// --------------------------
	t.Run("DeployCCIPCommon", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployCCIPCommonOp, deps, DeployCCIPCommonInput{
			InstanceID:         instanceID,
			ChainSelectorValue: chainSelectorValue,
			OnRampAddress:      onRampAddress,
		})
		require.NoError(t, err, "failed to deploy CCIP Common")
		require.NotEmpty(t, result.Output.Output.GlobalConfigContractID, "GlobalConfig contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.GlobalConfigTemplateID, "GlobalConfig template ID should not be empty")
		t.Logf("Deployed GlobalConfig contract ID: %s", result.Output.Output.GlobalConfigContractID)
	})

	// --------------------------
	// Test Token Admin Registry Deployment
	// --------------------------
	t.Run("DeployTokenAdminRegistry", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployTokenAdminRegistryOp, deps, DeployTokenAdminRegistryInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy TokenAdminRegistry")
		require.NotEmpty(t, result.Output.Output.TokenAdminRegistryContractID, "TokenAdminRegistry contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.TokenAdminRegistryTemplateID, "TokenAdminRegistry template ID should not be empty")
		t.Logf("Deployed TokenAdminRegistry contract ID: %s", result.Output.Output.TokenAdminRegistryContractID)
	})

	// --------------------------
	// Test Committee Verifier Deployment
	// --------------------------
	t.Run("DeployCommitteeVerifier", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployCommitteeVerifierOp, deps, DeployCommitteeVerifierInput{
			InstanceID:          instanceID,
			VersionTag:          "1.0.0",
			StorageLocation:     "ipfs://test-ccv",
			Threshold:           2,
			Signers:             []string{"signer1", "signer2", "signer3"},
			MessageSentObserver: "", // Will default to deployer party
		})
		require.NoError(t, err, "failed to deploy CommitteeVerifier")
		require.NotEmpty(t, result.Output.Output.CommitteeVerifierContractID, "CommitteeVerifier contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CommitteeVerifierTemplateID, "CommitteeVerifier template ID should not be empty")
		t.Logf("Deployed CommitteeVerifier contract ID: %s", result.Output.Output.CommitteeVerifierContractID)
	})

	// --------------------------
	// Test CCV Registry Deployment
	// --------------------------
	t.Run("DeployCCVRegistry", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployCCVRegistryOp, deps, DeployCCVRegistryInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy CCVRegistry")
		require.NotEmpty(t, result.Output.Output.CCVRegistryContractID, "CCVRegistry contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.CCVRegistryTemplateID, "CCVRegistry template ID should not be empty")
		t.Logf("Deployed CCVRegistry contract ID: %s", result.Output.Output.CCVRegistryContractID)
	})

	// --------------------------
	// Test Fee Quoter Deployment
	// --------------------------
	t.Run("DeployFeeQuoter", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployFeeQuoterOp, deps, DeployFeeQuoterInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy FeeQuoter")
		require.NotEmpty(t, result.Output.Output.FeeQuoterContractID, "FeeQuoter contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.FeeQuoterTemplateID, "FeeQuoter template ID should not be empty")
		t.Logf("Deployed FeeQuoter contract ID: %s", result.Output.Output.FeeQuoterContractID)
	})

	// --------------------------
	// Test OffRamp Deployment
	// --------------------------
	t.Run("DeployOffRamp", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployOffRampOp, deps, DeployOffRampInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy OffRamp")
		require.NotEmpty(t, result.Output.Output.OffRampContractID, "OffRamp contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.OffRampTemplateID, "OffRamp template ID should not be empty")
		t.Logf("Deployed OffRamp contract ID: %s", result.Output.Output.OffRampContractID)
	})

	// --------------------------
	// Test PerPartyRouter Deployment
	// --------------------------
	t.Run("DeployPerPartyRouter", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployPerPartyRouterOp, deps, DeployPerPartyRouterInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy PerPartyRouter")
		require.NotEmpty(t, result.Output.Output.PerPartyRouterContractID, "PerPartyRouter contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.PerPartyRouterTemplateID, "PerPartyRouter template ID should not be empty")
		t.Logf("Deployed PerPartyRouter contract ID: %s", result.Output.Output.PerPartyRouterContractID)
	})

	// --------------------------
	// Test OnRamp Deployment
	// --------------------------
	t.Run("DeployOnRamp", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployOnRampOp, deps, DeployOnRampInput{
			InstanceID:           instanceID,
			DestChainSelector:    "2222222222",
			DestChainOnRampBytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		})
		require.NoError(t, err, "failed to deploy OnRamp")
		require.NotEmpty(t, result.Output.Output.OnRampContractID, "OnRamp contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.OnRampTemplateID, "OnRamp template ID should not be empty")
		t.Logf("Deployed OnRamp contract ID: %s", result.Output.Output.OnRampContractID)
	})
}
