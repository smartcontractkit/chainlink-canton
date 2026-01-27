package linkops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

func TestDeployAndMintLink(t *testing.T) {
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

	deps := compileClient.CantonOpDeps{
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

	result, err := cld_ops.ExecuteOperation(bundle, DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	_, err = cld_ops.ExecuteOperation(bundle, MintLINKPreApprovalOp, deps, MintLinkTokenInput{
		RegistryContractID: result.Output.Output.RegistryContractID,
		ReceiverParty:      setupResult.Party, // approve preapproval to mint for self
		Amount:             "100000",
	})
	require.NoError(t, err, "failed to mint LINK token")
}
