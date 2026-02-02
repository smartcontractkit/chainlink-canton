package changesets

import (
	"sync"
	"testing"

	"github.com/noders-team/go-daml/pkg/client"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestDeployCoin(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	// TODO Expose via CLDF/CTF
	token, err := bc.(*canton.Chain).Participants[0].JWTProvider.Token(t.Context())
	require.NoError(t, err)
	bindingClient, err := client.NewDamlClient(token, bc.(*canton.Chain).Participants[0].Endpoints.GRPCLedgerAPIURL).
		WithAdminAddress(bc.(*canton.Chain).Participants[0].Endpoints.AdminAPIURL).
		Build(t.Context())
	require.NoError(t, err, "failed to create Daml binding client")
	t.Cleanup(bindingClient.Close)

	// Upload Dar
	coinDar, err := contracts.GetDar(contracts.Coin, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get coin dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), coinDar, "")
	require.NoError(t, err, "failed to upload coin dar file")

	// Get primary party
	user, err := bindingClient.UserMng.GetUser(t.Context(), "user-participant1")
	require.NoError(t, err, "failed to get user")

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)
	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	config := CantonCSDeps[DeployCoinConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		UserName:      user.ID,
		Party:         user.PrimaryParty,
		Config: DeployCoinConfig{
			Symbol: "LINK",
		},
	}

	deployCoin := DeployCoin{}

	_, err = deployCoin.Apply(*env, config)
	require.NoError(t, err)
}
