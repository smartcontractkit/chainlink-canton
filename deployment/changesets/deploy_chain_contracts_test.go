package changesets

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/ethereum/go-ethereum/crypto"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-canton/bindings/auth"
	"github.com/smartcontractkit/chainlink-canton/bindings/types"
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

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

func TestDeployChainContracts(t *testing.T) {
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

	// Upload Dars
	dars := []struct {
		contract contracts.Package
		name     string
	}{
		{contracts.CCIPCommon, "common"},
		{contracts.CCIPOffRamp, "offramp"},
		{contracts.CCIPOnRamp, "onramp"},
		{contracts.CCIPTokenAdminRegistry, "token admin registry"},
		{contracts.CCIPCommitteeVerifier, "committee verifier"},
		{contracts.CCIPLockReleaseTokenPool, "token pool"},
		{contracts.CCIPPerPartyRouter, "per-party router"},
	}

	darData := make([]*participantv30.UploadDarRequest_UploadDarData, 0, len(dars))
	for _, dar := range dars {
		darBytes, err := contracts.GetDar(dar.contract, contracts.CurrentVersion)
		require.NoError(t, err, "failed to get %s dar file", dar.name)
		darData = append(darData, &participantv30.UploadDarRequest_UploadDarData{Bytes: darBytes})
	}

	_, err = packageServiceClient.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars:               darData,
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload DAR files")

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

	// CCV Setup
	// Generate signer keys for CommitteeVerifier
	ccvSignerPubKeys := make([]types.TEXT, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, types.TEXT(pubKeyHex))
	}
	versionTag := "49ff34ed"
	// ccvID := versionTag + "@" + user.PrimaryParty

	chainSelector := types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10))
	config := CantonCSDeps[DeployChainContractsConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Party:         user.GetPrimaryParty(),
		Config: DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: user.GetPrimaryParty(),
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Template: ccvs.CommitteeVerifier{
							Owner:               types.PARTY(user.GetPrimaryParty()),
							CcipOwner:           types.PARTY(user.GetPrimaryParty()),
							VersionTag:          types.TEXT(versionTag),
							MessageSentObserver: types.PARTY(user.GetPrimaryParty()),
							StorageLocation:     "ipfs://test-receive",
							Threshold:           2,
							Signers:             ccvSignerPubKeys,
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: chainSelector,
						OnRampAddress: "", // TODO ?
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						CcipOwner:      "", // Populated by the sequence
						RmnOwner:       types.PARTY(user.GetPrimaryParty()),
						CursedSubjects: nil,
					},
				},
			},
		},
	}

	deployChainContracts := DeployChainContracts{}

	out, err := deployChainContracts.Apply(*env, config)
	require.NoError(t, err)

	addresses := out.DataStore.Addresses().Filter()
	for i, address := range addresses {
		fmt.Printf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier)
	}
}
