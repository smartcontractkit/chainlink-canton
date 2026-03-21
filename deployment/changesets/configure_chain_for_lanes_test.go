package changesets

import (
	"encoding/hex"
	"strconv"
	"sync"
	"testing"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
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

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	ccipOwnerParty := participant.PartyID

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

	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars:               darData,
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload DAR files")

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
	ccvSignerPubKeys := make([]string, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}
	versionTag := "49ff34ed"
	ccvQualifier := "default"

	// Deploy Chain Contracts
	out, err := DeployChainContracts{}.Apply(*env, CantonCSDeps[DeployChainContractsConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: ccipOwnerParty,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Template: ccvs.CommitteeVerifier{
							Owner:                        types.PARTY(ccipOwnerParty),
							CcipOwner:                    types.PARTY(ccipOwnerParty),
							VersionTag:                   types.TEXT(versionTag),
							MessageSentObservers:         nil,
							StorageLocations:             []types.TEXT{"ipfs://test-receive"},
							StorageLocationsAdmin:        types.PARTY(ccipOwnerParty),
							PendingStorageLocationsAdmin: types.PARTY(ccipOwnerParty),
						},
						Qualifier: ccvQualifier,
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: types.NUMERIC(strconv.FormatUint(chainSelector, 10)),
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						CcipOwner:      "", // Populated by the sequence
						RmnOwner:       types.PARTY(ccipOwnerParty),
						CursedSubjects: nil,
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
		t.Logf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s, Labels=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier, address.Labels.String())
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
	committeeVerifier, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, datastore.ContractType(committee_verifier.ContractType), committee_verifier.Version, ccvQualifier))
	require.NoError(t, err, "failed to get CommitteeVerifier address")

	// Configure Chain for Lanes
	committeeVerifierRawAddr, err := contracts.RawInstanceAddressFromString(committeeVerifier.Labels.List()[0])
	require.NoError(t, err, "failed to parse CommitteeVerifier raw address")
	out, err = ConfigureChainForLanes{}.Apply(*env, CantonCSDeps[ConfigureChainForLanesConfig]{
		ChainSelector: chainSelector,
		Participant:   0,
		Config: ConfigureChainForLanesConfig{
			Input: sequences.ConfigureChainForLanesInput{
				ChainSelector: chainSelector,
				GlobalConfig:  contracts.HexToInstanceAddress(globalConfig.Address),
				FeeQuoter:     contracts.HexToInstanceAddress(feeQuoter.Address),
				OnRamp:        contracts.HexToInstanceAddress(onRamp.Address),
				OffRamp:       contracts.HexToInstanceAddress(offRamp.Address),
				CommitteeVerifiers: []adapters.CommitteeVerifierConfig[contracts.InstanceAddress]{
					{
						CommitteeVerifier: []contracts.InstanceAddress{contracts.HexToInstanceAddress(committeeVerifier.Address)},
						RemoteChains: map[uint64]adapters.CommitteeVerifierRemoteChainConfig{
							chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
								AllowlistEnabled:   false,
								FeeUSDCents:        50,
								GasForVerification: 50_000,
								PayloadSizeBytes:   6*64 + 2*32,
								SignatureConfig: adapters.CommitteeVerifierSignatureQuorumConfig{
									Signers:   ccvSignerPubKeys,
									Threshold: 2,
								},
							},
						},
					},
				},
				RemoteChains: map[uint64]adapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]{
					chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
						AllowTrafficFrom:         true,
						OnRamps:                  [][]byte{hexutil.MustDecode("0x0b5c6e23bb6f8e3abe5fcbdd406f3df8b96b8e1c")},
						OffRamp:                  hexutil.MustDecode("0xbd8f6f5f14d9efbde3c72cc1affc968a5f49a2b3"),
						DefaultInboundCCVs:       []contracts.RawInstanceAddress{committeeVerifierRawAddr},
						LaneMandatedInboundCCVs:  nil,
						DefaultOutboundCCVs:      []contracts.RawInstanceAddress{committeeVerifierRawAddr},
						LaneMandatedOutboundCCVs: nil,
						DefaultExecutor:          "",
						FeeQuoterDestChainConfig: adapters.FeeQuoterDestChainConfig{
							IsEnabled:                   true,
							MaxDataBytes:                50000,
							MaxPerMsgGasLimit:           4000000,
							DestGasOverhead:             300000,
							DestGasPerPayloadByteBase:   16,
							ChainFamilySelector:         [4]byte{0x28, 0x12, 0xd5, 0x2c},
							DefaultTxGasLimit:           200000,
							DefaultTokenFeeUSDCents:     0,
							DefaultTokenDestGasOverhead: 34000,
							NetworkFeeUSDCents:          0,
						},
						ExecutorDestChainConfig: adapters.ExecutorDestChainConfig{},
						AddressBytesLength:      20,
						BaseExecutionGasCost:    0,
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
