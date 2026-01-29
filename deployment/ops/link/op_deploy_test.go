package linkops

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/noders-team/go-daml/pkg/client"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

func TestDeployAndMintLink(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)
	chain := bc.(*canton.Chain)

	// TODO: use a proper JWT provider for the bindings, the BindingClient doesn't yet allow using a auth.TokenProvider
	token, err := chain.Participants[0].JWTProvider.Token(t.Context())
	require.NoError(t, err)

	bindingClient, err := client.NewDamlClient(token, chain.Participants[0].Endpoints.GRPCLedgerAPIURL).
		WithAdminAddress(chain.Participants[0].Endpoints.AdminAPIURL).
		Build(ctx)
	require.NoError(t, err, "failed to create Daml binding client")
	t.Cleanup(bindingClient.Close)

	// Upload Dar
	coinDar, err := contracts.GetDar(contracts.Coin, contracts.CurrentVersion)
	err = bindingClient.PackageMng.UploadDarFile(ctx, coinDar, "")
	require.NoError(t, err, "failed to upload coin dar file")

	// Get primary party
	user, err := bindingClient.UserMng.GetUser(ctx, "user-participant1")
	require.NoError(t, err, "failed to get user")

	deps := compileClient.CantonOpDeps{
		BindingClient: bindingClient,
		Party:         user.PrimaryParty,
	}

	reporter := cld_ops.NewMemoryReporter()

	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	result, err := cld_ops.ExecuteOperation(bundle, DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	fmt.Println("Created LINK token registry contract:")
	fmt.Println("UpdateId: ", result.Output.UpdateID)
	fmt.Println("InstanceId: ", result.Output.Output.RegistryInstanceID)
	fmt.Println("ContractId: ", result.Output.Output.RegistryContractID)

	_, err = cld_ops.ExecuteOperation(bundle, MintLINKPreApprovalOp, deps, MintLinkTokenInput{
		RegistryContractID: result.Output.Output.RegistryContractID,
		ReceiverParty:      user.PrimaryParty, // approve preapproval to mint for self
		Amount:             "100000",
	})
	require.NoError(t, err, "failed to mint LINK token")
}
