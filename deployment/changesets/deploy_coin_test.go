package changesets

import (
	"sync"
	"testing"

	"https://github.com/smartcontractkit/go-daml/pkg/auth"

	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestDeployCoin(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	token, err := bc.(*canton.Chain).Participants[0].JWTProvider.Token(t.Context())
	require.NoError(t, err)

	// Create gRPC clients
	insecureCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	adminApiClient, err := grpc.NewClient(bc.(*canton.Chain).Participants[0].Endpoints.AdminAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	require.NoError(t, err, "Failed to dial gRPC admin API")
	ledgerApiClient, err := grpc.NewClient(bc.(*canton.Chain).Participants[0].Endpoints.GRPCLedgerAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	require.NoError(t, err, "Failed to dial gRPC ledger API")

	packageServiceClient := participantv30.NewPackageServiceClient(adminApiClient)
	userManagementServiceClient := admin.NewUserManagementServiceClient(ledgerApiClient)

	// Upload Dar
	coinDar, err := contracts.GetDar(contracts.Coin, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get coin dar file")
	_, err = packageServiceClient.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: coinDar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload coin dar file")

	// Get primary party
	userResp, err := userManagementServiceClient.GetUser(t.Context(), &admin.GetUserRequest{
		UserId: "user-participant1",
	})
	require.NoError(t, err, "failed to get user")
	user := userResp.GetUser()

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
		UserName:      user.GetId(),
		Party:         user.GetPrimaryParty(),
		Config: DeployCoinConfig{
			Symbol: "LINK",
		},
	}

	deployCoin := DeployCoin{}

	_, err = deployCoin.Apply(*env, config)
	require.NoError(t, err)
}
