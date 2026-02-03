package changesets

import (
	"encoding/hex"
	"fmt"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

func TestConfigureChainForLanes(t *testing.T) {
	t.Parallel()

	chainSelector := chainsel.CANTON_LOCALNET.Selector

	bc, err := cantonProvider.NewCTFChainProvider(t, chainSelector, cantonProvider.CTFChainProviderConfig{
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

	// Upload Dars
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get common dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), commonDar, "")
	require.NoError(t, err, "failed to upload common dar file")
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get offramp dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), offRampDar, "")
	require.NoError(t, err, "failed to upload offramp dar file")
	onRampDar, err := contracts.GetDar(contracts.CCIPOnRamp, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get onramp dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), onRampDar, "")
	require.NoError(t, err, "failed to upload onramp dar file")
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get token admin registry dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), tokenAdminRegistryDar, "")
	require.NoError(t, err, "failed to upload token admin registry dar file")
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get committee verifier dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), committeeVerifierDar, "")
	require.NoError(t, err, "failed to upload committee verifier dar file")
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get token pool dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), tokenPoolDar, "")
	require.NoError(t, err, "failed to upload token pool dar file")
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get per-party router dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), perPartyRouterDar, "")
	require.NoError(t, err, "failed to upload per-party router dar file")

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
	ccvID := versionTag + "@" + user.PrimaryParty

	// Deploy Chain Contracts
	out, err := DeployChainContracts{}.Apply(*env, CantonCSDeps[DeployChainContractsConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Party:         user.PrimaryParty,
		Config: DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: user.PrimaryParty,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Template: ccvs.CommitteeVerifier{
							Owner:               types.PARTY(user.PrimaryParty),
							CcipOwner:           types.PARTY(user.PrimaryParty),
							VersionTag:          types.TEXT(versionTag),
							MessageSentObserver: types.PARTY(user.PrimaryParty),
							StorageLocation:     "ipfs://test-receive",
							Threshold:           2,
							Signers:             ccvSignerPubKeys,
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: types.NUMERIC(chainSelector),
						OnRampAddress: "", // TODO ?
					},
				},
			},
		},
	})
	require.NoError(t, err)

	err = out.DataStore.Merge(env.DataStore)
	require.NoError(t, err)
	env.DataStore = out.DataStore.Seal()

	addresses := env.DataStore.Addresses().Filter()
	for i, address := range addresses {
		fmt.Printf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier)
	}

	// Resolve contracts
	globalConfig, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	require.NoError(t, err, "failed to get GlobalConfig address")
	feeQuoter, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, ""))
	require.NoError(t, err, "failed to get FeeQuoter address")
	onRamp, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, datastore.ContractType(onramp.ContractType), onramp.Version, ""))
	require.NoError(t, err, "failed to get OnRamp address")
	offRamp, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, datastore.ContractType(offramp.ContractType), offramp.Version, ""))
	require.NoError(t, err, "failed to get OffRamp address")

	// Configure Chain for Lanes
	out, err = ConfigureChainForLanes{}.Apply(*env, CantonCSDeps[ConfigureChainForLanesConfig]{
		ChainSelector: chainSelector,
		Participant:   0,
		UserName:      user.ID,
		Party:         user.PrimaryParty,
		Config: ConfigureChainForLanesConfig{
			Input: sequences.ConfigureChainForLanesInput{
				ChainSelector:      chainSelector,
				GlobalConfig:       contracts.InstanceID(globalConfig.Address),
				FeeQuoter:          contracts.InstanceID(feeQuoter.Address),
				OnRamp:             contracts.InstanceID(onRamp.Address),
				OffRamp:            contracts.InstanceID(offRamp.Address),
				CommitteeVerifiers: nil,
				RemoteChains: map[uint64]adapters.RemoteChainConfig[[]byte, string]{
					chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
						AllowTrafficFrom:         true,
						OnRamps:                  [][]byte{[]byte("")},
						OffRamp:                  nil,
						DefaultInboundCCVs:       nil,
						LaneMandatedInboundCCVs:  []string{ccvID},
						DefaultOutboundCCVs:      nil,
						LaneMandatedOutboundCCVs: nil,
						DefaultExecutor:          "",
						FeeQuoterDestChainConfig: adapters.FeeQuoterDestChainConfig{},
						ExecutorDestChainConfig:  adapters.ExecutorDestChainConfig{},
						AddressBytesLength:       0,
						BaseExecutionGasCost:     0,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	err = out.DataStore.Merge(env.DataStore)
	require.NoError(t, err)
	env.DataStore = out.DataStore.Seal()
}
