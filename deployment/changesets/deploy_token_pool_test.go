package changesets

import (
	"sync"
	"testing"

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
	"github.com/smartcontractkit/go-daml/pkg/auth"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestDeployTokenPool(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	token, err := bc.(*canton.Chain).Participants[0].JWTProvider.Token(t.Context())
	require.NoError(t, err)

	insecureCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	adminApiClient, err := grpc.NewClient(bc.(*canton.Chain).Participants[0].Endpoints.AdminAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	require.NoError(t, err, "Failed to dial gRPC admin API")
	ledgerApiClient, err := grpc.NewClient(bc.(*canton.Chain).Participants[0].Endpoints.GRPCLedgerAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	require.NoError(t, err, "Failed to dial gRPC ledger API")

	packageServiceClient := participantv30.NewPackageServiceClient(adminApiClient)
	userManagementServiceClient := admin.NewUserManagementServiceClient(ledgerApiClient)

	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get ccip-common dar file")
	poolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get lockreleasetokenpool dar file")
	_, err = packageServiceClient.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: commonDar},
			{Bytes: poolDar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload dar files")

	userResp, err := userManagementServiceClient.GetUser(t.Context(), &admin.GetUserRequest{
		UserId: "user-participant1",
	})
	require.NoError(t, err, "failed to get user")
	user := userResp.GetUser()
	party := user.GetPrimaryParty()

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

	config := CantonCSDeps[DeployTokenPoolConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		UserName:      user.GetId(),
		Party:         party,
		Config: DeployTokenPoolConfig{
			CcipOwner: party,
			PoolOwner: party,
			InstrumentId: lockreleasetokenpool.InstrumentId{
				Admin: types.PARTY(party),
				Id:    types.TEXT("AMT"),
			},
			Decimals:  6,
			Qualifier: "AMT",
		},
	}

	deployTokenPool := DeployTokenPool{}

	output, err := deployTokenPool.Apply(*env, config)
	require.NoError(t, err)
	addresses := output.DataStore.Addresses().Filter()
	require.NotEmpty(t, addresses)
}
