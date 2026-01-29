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

func TestDeployMCMS(t *testing.T) {
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
	// This is a minimal config for testing
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

	config := mcms.MultisigConfig{
		Signers:      signers,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
	}

	// --------------------------
	// Test MCMS Deployment
	// --------------------------
	t.Run("DeployMCMS", func(t *testing.T) {
		t.Parallel()

		result, err := cld_ops.ExecuteOperation(bundle, DeployMCMSOp, deps, DeployMCMSInput{
			InstanceID: instanceID,
			ChainID:    chainID,
			MCMSID:     mcmsID,
			Role:       "Proposer",
			Config:     config,
		})
		require.NoError(t, err, "failed to deploy MCMS")
		require.NotEmpty(t, result.Output.Output.MCMSContractID, "MCMS contract ID should not be empty")
		require.NotEmpty(t, result.Output.Output.MCMSTemplateID, "MCMS template ID should not be empty")
		t.Logf("Deployed MCMS contract ID: %s", result.Output.Output.MCMSContractID)
		t.Logf("Deployed MCMS template ID: %s", result.Output.Output.MCMSTemplateID)
	})
}
