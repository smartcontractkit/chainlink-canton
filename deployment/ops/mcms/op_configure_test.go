package mcms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/mcms"
	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

//nolint:paralleltest // Cannot run in parallel due to shared state (mcmsContractID, mcmsTemplateID)
func TestConfigureMCMS(t *testing.T) {
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

	t.Cleanup(func() { setupResult.BindingClient.Close() })

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

	instanceID := "local-mcms"
	chainID := int64(1)
	mcmsID := "test-mcms-" + instanceID

	// Create a simple 2-of-3 multisig config
	signers := []mcms.SignerInfo{
		{
			SignerAddress: types.TEXT("0x1111111111111111111111111111111111111111"),
			SignerIndex:   types.INT64(0),
			SignerGroup:   types.INT64(0),
		},
		{
			SignerAddress: types.TEXT("0x2222222222222222222222222222222222222222"),
			SignerIndex:   types.INT64(1),
			SignerGroup:   types.INT64(0),
		},
		{
			SignerAddress: types.TEXT("0x3333333333333333333333333333333333333333"),
			SignerIndex:   types.INT64(2),
			SignerGroup:   types.INT64(0),
		},
	}

	// Create group quorums (32 groups, 2-of-3 for group 0)
	groupQuorums := make([]types.INT64, 32)
	groupQuorums[0] = types.INT64(2) // 2-of-3 for group 0

	// Create group parents (32 groups, all 0 for flat structure)
	groupParents := make([]types.INT64, 32)

	initialConfig := mcms.MultisigConfig{
		Signers:      signers,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
	}

	// --------------------------
	// Deploy MCMS first (required for configuration)
	// --------------------------
	var mcmsContractID string
	var mcmsTemplateID string

	t.Run("DeployMCMS", func(t *testing.T) {
		result, err := cld_ops.ExecuteOperation(bundle, DeployMCMSOp, deps, DeployMCMSInput{
			InstanceID: instanceID,
			ChainID:    chainID,
			MCMSID:     mcmsID,
			Role:       "Proposer",
			Config:     initialConfig,
		})
		require.NoError(t, err, "failed to deploy MCMS")
		require.NotEmpty(t, result.Output.Output.MCMSContractID, "MCMS contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.MCMSTemplateID, "MCMS template ID should not be empty")
		mcmsContractID = result.Output.Output.MCMSContractID
		mcmsTemplateID = result.Output.Output.MCMSTemplateID
		t.Logf("Deployed MCMS contract ID: %s", mcmsContractID)
	})

	// --------------------------
	// Test SetConfig Operation
	// --------------------------
	t.Run("SetConfig", func(t *testing.T) {
		// Create updated config with a new signer
		newSigners := []mcms.SignerInfo{
			{
				SignerAddress: types.TEXT("0x1111111111111111111111111111111111111111"),
				SignerIndex:   types.INT64(0),
				SignerGroup:   types.INT64(0),
			},
			{
				SignerAddress: types.TEXT("0x2222222222222222222222222222222222222222"),
				SignerIndex:   types.INT64(1),
				SignerGroup:   types.INT64(0),
			},
			{
				SignerAddress: types.TEXT("0x3333333333333333333333333333333333333333"),
				SignerIndex:   types.INT64(2),
				SignerGroup:   types.INT64(0),
			},
			{
				SignerAddress: types.TEXT("0x4444444444444444444444444444444444444444"),
				SignerIndex:   types.INT64(3),
				SignerGroup:   types.INT64(0),
			},
		}

		// Update quorum to 3-of-4
		newGroupQuorums := make([]types.INT64, 32)
		newGroupQuorums[0] = types.INT64(3) // 3-of-4 for group 0

		setConfigArgs := mcms.SetConfig{
			NewSigners:      newSigners,
			NewGroupQuorums: newGroupQuorums,
			NewGroupParents: groupParents,      // Keep same parent structure
			ClearRoot:       types.BOOL(false), // Don't clear root
		}

		result, err := cld_ops.ExecuteOperation(bundle, SetConfigOp, deps, SetConfigInput{
			MCMSContractID: mcmsContractID,
			MCMSTemplateID: mcmsTemplateID,
			Config:         setConfigArgs,
		})
		require.NoError(t, err, "failed to set MCMS config")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")
		require.NotEmpty(t, result.Output.Output.NewMCMSContractID, "New MCMS contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.NewMCMSTemplateID, "New MCMS template ID should not be empty")

		// Update contract ID for subsequent operations
		mcmsContractID = result.Output.Output.NewMCMSContractID
		mcmsTemplateID = result.Output.Output.NewMCMSTemplateID

		t.Logf("SetConfig succeeded, tx=%s newCID=%s",
			result.Output.Output.TransactionID,
			result.Output.Output.NewMCMSContractID,
		)
	})
}
