package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	contractops "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"

	// Import to register lane adapters
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"

	_ "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
)

type lnrTokenPoolReceiveFlowTestCase struct {
	tokenAmount                   *big.Int
	sourcePoolData                []byte
	expectedTransferAmount        float64
	defaultInboundLimiterCapacity string
	customInboundLimiterCapacity  string
	expectedDefaultLimiterTokens  string
	expectedCustomLimiterTokens   string
}

// TestLnRTokenPool_FullReceiveFlow tests the complete CCIP inbound token release flow.
// - Deploy all CCIP contracts including CommitteeVerifier, GlobalConfig, OffRamp, PerPartyRouter
// - Mint tokens to pool (simulating prior locked tokens)
// - Build inbound message with TokenTransfer
// - Generate signatures and call VerifyMessage to append CCV verification
// - Call PerPartyRouter.Execute to process message and get TokenReceiveTicket
// - Call TokenPool_ReleaseFromTicket to transfer tokens from pool to receiver
// - Verify receiver received the tokens
// - Validate FTF/custom-finality path by:
//   - requiring finalityConfig=2000-block-depth on the destination pool
//   - enabling a default inbound limiter with lower capacity
//   - enabling a custom-finality inbound limiter with higher capacity
//     Success proves ReleaseFromTicket selected the custom inbound limiter.
func TestLnRTokenPool_FullReceiveFlow(t *testing.T) {
	t.Parallel()

	runLnRTokenPoolReceiveFlowTest(t, lnrTokenPoolReceiveFlowTestCase{
		tokenAmount:                   big.NewInt(50_000_000_000),
		expectedTransferAmount:        5,
		defaultInboundLimiterCapacity: "1000000000000",
		customInboundLimiterCapacity:  "10000000000000",
	})
}

func TestLnRTokenPool_FullReceiveFlow_DecimalConversion(t *testing.T) {
	t.Parallel()

	runLnRTokenPoolReceiveFlowTest(t, lnrTokenPoolReceiveFlowTestCase{
		tokenAmount:                   new(big.Int).SetUint64(7_000_000_000_000_000_000),
		sourcePoolData:                encodeUint256Bytes(18),
		expectedTransferAmount:        7,
		defaultInboundLimiterCapacity: "50000000000",
		customInboundLimiterCapacity:  "100000000000",
		expectedDefaultLimiterTokens:  "50000000000.",
		expectedCustomLimiterTokens:   "30000000000.",
	})
}

func runLnRTokenPoolReceiveFlowTest(t *testing.T, tc lnrTokenPoolReceiveFlowTestCase) {
	t.Helper()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(3))

	// Setup participants
	ccipParticipant := env.Chain.Participants[0]
	receiverParticipant := env.Chain.Participants[1]
	tokenPoolOwnerParticipant := env.Chain.Participants[0]

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	// DAR Uploading
	// Read DARs
	rmnDar, err := contracts.GetDar(contracts.CCIPRMN, contracts.CurrentVersion)
	require.NoError(t, err)
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipReceiverDar, err := contracts.GetDar(contracts.CCIPReceiver, contracts.CurrentVersion)
	require.NoError(t, err)
	onRampDar, err := contracts.GetDar(contracts.CCIPOnRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	feeQuoterDar, err := contracts.GetDar(contracts.CCIPFeeQuoter, contracts.CurrentVersion)
	require.NoError(t, err)
	executorDar, err := contracts.GetDar(contracts.CCIPExecutor, contracts.CurrentVersion)
	require.NoError(t, err)

	// Upload DARs to all participants
	dars := [][]byte{rmnDar, commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, tokenPoolDar, perPartyRouterDar, ccipReceiverDar, onRampDar, feeQuoterDar, executorDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant, tokenPoolOwnerParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.PartyID
	partyReceiver := receiverParticipant.PartyID
	partyTokenPoolOwner := tokenPoolOwnerParticipant.PartyID
	t.Logf("Parties: CCIP=%s, Receiver=%s, PoolOwner=%s", partyCCIP, partyReceiver, partyTokenPoolOwner)

	// Create Scan and Registry API clients
	// Using the scanProxy endpoint of the 0-th participant, all participants are able to forward requests using the BFT Scan Proxy, it doesn't matter which one we use
	scanProxyClient, tokenMetadataClient, transferInstructionClient, err := testhelpers.NewValidatorAPIClients(ccipParticipant)
	require.NoError(t, err, "Failed to create validator API clients")

	// Get DSO Admin Party (registry admin)
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err)

	// Token Setup
	// Mint tokens to Token Pool Owner (these will be "locked" in the pool)
	poolHoldingCid, err := testhelpers.MintAMT(t.Context(), tokenPoolOwnerParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partyTokenPoolOwner, "100.00")
	require.NoError(t, err)
	t.Logf("Minted 100 AMT to Pool Owner, Holding CID: %s", poolHoldingCid)

	// Instrument ID for AMT
	nativeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}
	hashedInstrumentId := contracts.EncodeInstrumentID(nativeInstrumentId)

	// CCIP Deployment
	sourceChainSelector := fmt.Sprintf("%d", chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector)

	// CCV Setup
	// Generate signer keys for CommitteeVerifier
	ccvSignerKeys := make([]*ecdsa.PrivateKey, 0, 3)
	ccvSignerPubKeys := make([]string, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}
	t.Logf("Generated %d CCV signer keys", len(ccvSignerKeys))

	versionTag := "e9a05a20"
	ccvQualifier := devenvcommon.DefaultCommitteeVerifierQualifier

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)
	cldfEnv := cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{env.Chain}),
		OperationsBundle: bundle,
	}

	// Deploy Chain contracts via changeset
	out, err := changesets.DeployChainContracts{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.DeployChainContractsConfig]{
		ChainSelector: env.Chain.ChainSelector(),
		Participant:   0,
		Config: changesets.DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: partyCCIP,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Qualifier: ccvQualifier,
						Template: ccvs.CommitteeVerifier{
							Owner:                        types.PARTY(partyCCIP),
							CcipOwner:                    types.PARTY(partyCCIP),
							VersionTag:                   types.TEXT(versionTag),
							MessageSentObservers:         nil,
							StorageLocations:             []types.TEXT{"ipfs://test-receive"},
							StorageLocationsAdmin:        types.PARTY(partyCCIP),
							PendingStorageLocationsAdmin: types.PARTY(partyCCIP),
							Deps:                         ccvs.CommitteeVerifierDeps{},
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "",
						ChainSelector: types.NUMERIC(strconv.FormatUint(env.Chain.ChainSelector(), 10)),
					},
				},
				Executors: []sequences.ExecutorParams{
					{
						Qualifier: devenvcommon.DefaultExecutorQualifier,
						Template: executorBinding.Executor{
							Owner:         types.PARTY(partyCCIP),
							MaxCCVsPerMsg: 10,
							DynamicConfig: executorBinding.DynamicConfig{
								FeeAggregator:         nil,
								AllowedFinalityConfig: common.FinalityConfig{WaitForFinality: &types.UNIT{}},
								CcvAllowlistEnabled:   false,
							},
							AllowedCCVs: nil,
						},
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						CcipOwner:      "",
						RmnOwner:       types.PARTY(partyCCIP),
						CursedSubjects: nil,
					},
				},
				NativeInstrumentId: nativeInstrumentId,
			},
		},
	})
	require.NoErrorf(t, err, "Failed to deploy CCIP contracts: %v", err)

	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()

	t.Log("Deployed CCIP chain contracts:")
	for i, address := range cldfEnv.DataStore.Addresses().Filter() {
		t.Logf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s, Labels=%s",
			i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier, address.Labels.String())
	}

	// Resolve contracts
	globalConfigRef, globalConfigAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), global_config.ContractType, global_config.Version, "")
	require.NoError(t, err, "failed to get GlobalConfig address")
	_, feeQuoterAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), fee_quoter.ContractType, fee_quoter.Version, "")
	require.NoError(t, err, "failed to get FeeQuoter address")
	_, onRampAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), onramp.ContractType, onramp.Version, "")
	require.NoError(t, err, "failed to get OnRamp address")
	_, offRampAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), offramp.ContractType, offramp.Version, "")
	require.NoError(t, err, "failed to get OffRamp address")
	committeeVerifierRef, committeeVerifierAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), committee_verifier.ContractType, committee_verifier.Version, ccvQualifier)
	require.NoError(t, err, "failed to get CommitteeVerifier address")
	_, tokenAdminRegistryAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), token_admin_registry.ContractType, token_admin_registry.Version, "")
	require.NoError(t, err, "failed to get TokenAdminRegistry address")
	_, rmnRemoteAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), rmn_remote.ContractType, rmn_remote.Version, "")
	require.NoError(t, err, "failed to get RMNRemote address")
	_, perPartyRouterFactoryAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), per_party_router_factory.ContractType, per_party_router_factory.Version, "")
	require.NoError(t, err, "failed to get PerPartyRouterFactory address")

	// Deploy and configure lane (CCIP 2.0 lanes sequence, same as ccip_execute_test)
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	cantonAdapter, ok := lanes.GetLaneAdapterRegistry().GetLaneAdapter(chainsel.FamilyCanton, semver.MustParse("2.0.0"))
	require.Truef(t, ok, "failed to get Canton Lane adapter")
	deployLaneLegReport, err := cld_ops.ExecuteSequence(cldfEnv.OperationsBundle, cantonAdapter.ConfigureLaneLegAsDest(), cldfEnv.BlockChains, lanes.UpdateLanesInput{
		Source: &lanes.ChainDefinition{
			Selector: remoteSelector,
			OnRamp:   hexutil.MustDecode("0xf6eced5e96fff2de4f0ecd722beb57556fc443fd"),
			OffRamp:  hexutil.MustDecode("0xd8c9ec8cad3fb34aeca3ddbebfabe9f28a9bfaed"),
		},
		Dest: &lanes.ChainDefinition{
			Selector: env.Chain.ChainSelector(),
			CommitteeVerifiers: []lanes.CommitteeVerifierConfig[datastore.AddressRef]{
				{
					CommitteeVerifier: []datastore.AddressRef{committeeVerifierRef},
					RemoteChains: map[uint64]lanes.CommitteeVerifierRemoteChainConfig{
						remoteSelector: {
							AllowlistEnabled:          false,
							AddedAllowlistedSenders:   nil,
							RemovedAllowlistedSenders: nil,
							FeeUSDCents:               50,
							GasForVerification:        50_000,
							PayloadSizeBytes:          6*64 + 2*32,
							SignatureConfig: lanes.CommitteeVerifierSignatureQuorumConfig{
								Signers:   ccvSignerPubKeys,
								Threshold: 2,
							},
						},
					},
				},
			},
			LaneMandatedInboundCCVs: []datastore.AddressRef{committeeVerifierRef},
			DefaultInboundCCVs:      nil,
			CantonLaneConfig: &lanes.CantonLaneConfig{
				GlobalConfig: globalConfigRef,
			},
		},
		IsDisabled:   false,
		TestRouter:   false,
		ExtraConfigs: lanes.ExtraConfigs{},
	})
	require.NoErrorf(t, err, "Failed to configure chain for lanes")
	runningDs := datastore.NewMemoryDataStore()
	for _, address := range deployLaneLegReport.Output.Addresses {
		err = runningDs.Addresses().Add(address)
		require.NoErrorf(t, err, "Failed to add address %s", address.Address)
	}
	err = runningDs.Merge(cldfEnv.DataStore)
	require.NoErrorf(t, err, "Failed to merge datastore")
	cldfEnv.DataStore = runningDs.Seal()
	t.Log("Configured chain for lanes")

	// Token Pool Setup
	// Deploy default inbound RateLimiter required by ReleaseFromTicket receive flow.
	// Keep it enabled but lower-capacity so the test fails if the default-finality limiter
	// is selected for this FTF transfer instead of the custom-finality limiter.
	poolInstanceId := "test-pool-receive"
	inboundRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployInbound, env.Chain, contractops.DeployInput[common.RateLimiter]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: common.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceId),
			PoolOwner:           types.PARTY(partyTokenPoolOwner),
			RemoteChainSelector: types.NUMERIC(sourceChainSelector),
			Direction:           common.RateLimitDirectionRateLimitDirection_Inbound,
			Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
			IsEnabled:           true,
			Capacity:            types.NUMERIC(tc.defaultInboundLimiterCapacity),
			Rate:                types.NUMERIC(tc.defaultInboundLimiterCapacity),
			Tokens:              types.NUMERIC(tc.defaultInboundLimiterCapacity),
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
	})
	require.NoError(t, err, "failed to deploy inbound rate limiter")
	inboundRateLimiterRawAddr := inboundRateLimiterOut.Output.Labels.List()[0]
	inboundRateLimiterAddr, err := contracts.RawInstanceAddressFromString(inboundRateLimiterRawAddr)
	require.NoError(t, err, "failed to parse inbound rate limiter address")

	inboundCustomRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployInbound, env.Chain, contractops.DeployInput[common.RateLimiter]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: common.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceId),
			PoolOwner:           types.PARTY(partyTokenPoolOwner),
			RemoteChainSelector: types.NUMERIC(sourceChainSelector),
			Direction:           common.RateLimitDirectionRateLimitDirection_Inbound,
			Mode:                common.RateLimitModeRateLimitMode_CustomFinality,
			IsEnabled:           true,
			Capacity:            types.NUMERIC(tc.customInboundLimiterCapacity),
			Rate:                types.NUMERIC(tc.customInboundLimiterCapacity),
			Tokens:              types.NUMERIC(tc.customInboundLimiterCapacity),
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
	})
	require.NoError(t, err, "failed to deploy inbound custom rate limiter")
	inboundCustomRateLimiterRawAddr := inboundCustomRateLimiterOut.Output.Labels.List()[0]
	inboundCustomRateLimiterAddr, err := contracts.RawInstanceAddressFromString(inboundCustomRateLimiterRawAddr)
	require.NoError(t, err, "failed to parse inbound custom rate limiter address")

	outboundRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployOutbound, env.Chain, contractops.DeployInput[common.RateLimiter]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: common.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceId),
			PoolOwner:           types.PARTY(partyTokenPoolOwner),
			RemoteChainSelector: types.NUMERIC(sourceChainSelector),
			Direction:           common.RateLimitDirectionRateLimitDirection_Outbound,
			Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
			IsEnabled:           false,
			Capacity:            types.NUMERIC("0"),
			Rate:                types.NUMERIC("0"),
			Tokens:              types.NUMERIC("0"),
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
	})
	require.NoError(t, err, "failed to deploy outbound rate limiter")
	outboundRateLimiterRawAddr := outboundRateLimiterOut.Output.Labels.List()[0]
	outboundRateLimiterAddr, err := contracts.RawInstanceAddressFromString(outboundRateLimiterRawAddr)
	require.NoError(t, err, "failed to parse outbound rate limiter address")

	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")
	out, err = changesets.DeployTokenPool{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.DeployTokenPoolConfig]{
		ChainSelector: env.Chain.ChainSelector(),
		Participant:   0,
		Config: changesets.DeployTokenPoolConfig{
			CcipOwner:          partyCCIP,
			PoolOwner:          partyTokenPoolOwner,
			InstrumentId:       nativeInstrumentId,
			Decimals:           10,
			InstanceID:         poolInstanceId,
			PoolReceiveContext: common.CCIPContext{Values: map[string]common.AnyValue{}},
			TransferTimeout: lockreleasetokenpool.TransferTimeout{
				RelativeHours: func(v types.INT64) *types.INT64 { return &v }(types.INT64(24)),
			},
			RemoteChainConfigs: map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig{
				types.NUMERIC(sourceChainSelector): {
					RemotePools:        []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
					RemoteTokenAddress: types.TEXT(hex.EncodeToString(remoteTokenAddress)),
					InboundCCVs:        []mcms.RawInstanceAddress{},
					OutboundCCVs:       []mcms.RawInstanceAddress{},
					FinalityConfig: common.FinalityConfig{
						BlockDepth: new(types.INT64(2000)),
					},
					InboundRateLimiter:                         inboundRateLimiterAddr.Binding(),
					InboundCustomBlockConfirmationsRateLimiter: inboundCustomRateLimiterAddr.Binding(),
					OutboundRateLimiter:                        outboundRateLimiterAddr.Binding(),
				},
			},
			TokenTransferFeeConfigs: map[types.NUMERIC]lockreleasetokenpool.TokenTransferFeeConfig2{},
			Deps: lockreleasetokenpool.LockReleaseTokenPoolDeps{
				TokenAdminRegistry: tokenAdminRegistryAddress.Binding(),
				RmnRemote:          rmnRemoteAddress.Binding(),
				FeeQuoter:          feeQuoterAddress.Binding(),
			},
			// By setting the TAR address, the CS will automatically register the newly deployed pool with the TAR
			TokenAdminRegistryInstanceAddress: tokenAdminRegistryAddress.InstanceAddress(),
		},
	})
	require.NoError(t, err, "Failed to deploy Token Pool")
	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()
	_, tokenPoolAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), lock_release_token_pool.ContractType, lock_release_token_pool.Version, "")
	require.NoError(t, err, "failed to get Token Pool address")

	// Run EDS
	edsParticipant := env.Chain.Participants[0]
	edsToken, _ := edsParticipant.TokenSource.Token()
	edsPort := freeport.GetOne(t)
	go func() {
		log.Info().Msg("Running EDS...")
		err := service.RunEDS(t.Context(), log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel), &config.Config{
			ChainSelector: strconv.FormatUint(env.Chain.ChainSelector(), 10),
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: uint16(edsPort),
			},
			Node: config.NodeConfig{
				URL: edsParticipant.Endpoints.GRPCLedgerAPIURL,
				AuthConfig: commonconfig.AuthConfig{
					Type:   commonconfig.AuthTypeInsecureStatic,
					UserID: edsParticipant.UserID,
					JWT:    edsToken.AccessToken,
				},
				MaxRetries: 0,
			},
			CCIPAPIConfig: config.CCIPAPIConfig{
				Enabled: true,
				PerPartyRouterFactory: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: perPartyRouterFactoryAddress.InstanceAddress(),
				},
				OnRamp: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: onRampAddress.InstanceAddress(),
				},
				OffRamp: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: offRampAddress.InstanceAddress(),
				},
				GlobalConfig: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: globalConfigAddress.InstanceAddress(),
				},
				TokenAdminRegistry: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: tokenAdminRegistryAddress.InstanceAddress(),
				},
				RMNRemote: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: rmnRemoteAddress.InstanceAddress(),
				},
				FeeQuoter: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: feeQuoterAddress.InstanceAddress(),
				},
			},
			CCVAPIConfig: config.CCVAPIConfig{
				Enabled: true,
				CCVs: []config.CCV{
					{
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         partyCCIP,
							InstanceAddress: committeeVerifierAddress.InstanceAddress(),
						},
					},
				},
			},
			TokenPoolAPIConfig: config.TokenPoolAPIConfig{
				Enabled: true,
				TokenPools: []config.TokenPool{
					{
						Type: config.TokenPoolTypeLockRelease,
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         partyCCIP,
							InstanceAddress: tokenPoolAddress.InstanceAddress(),
						},
						PoolOwner: partyCCIP,
						// By setting the TokenStandard info, the Toke Pool API will return the necessary factory disclosures
						TokenStandardURL: new(fmt.Sprintf("%s/v0/scan-proxy", ccipParticipant.Endpoints.ValidatorAPIURL)),
						TokenStandardAuthConfig: &commonconfig.AuthConfig{
							Type: commonconfig.AuthTypeInsecureStatic,
							JWT:  edsToken.AccessToken,
						},
					},
				},
			},
		})
		log.Info().Err(err).Msg("EDS terminated")
		if !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Msg("EDS server exited with error")
			t.Fail()
			return
		}
	}()

	// Create EDS clients
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create CCIP API client")
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create CCV API client")
	tokenPoolAPIClient, err := oapiTokenPool.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create Token Pool API client")

	// wait for EDS to start up
	time.Sleep(1 * time.Second)

	// Create PerPartyRouter for receiver via EDS
	perPartyRouterFactoryDisclosure, err := edsTesthelpers.GetPerPartyRouterFactoryDisclosure(t.Context(), ccipAPIClient, partyReceiver)
	require.NoError(t, err)
	t.Logf("disclosed contracts: %v", perPartyRouterFactoryDisclosure.DisclosedContracts)
	t.Logf("per party router factory cid: %s", perPartyRouterFactoryDisclosure.ContractId)

	res, err := receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: perPartyRouterFactoryDisclosure.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-receiver"}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: perPartyRouterFactoryDisclosure.DisclosedContracts,
		},
	})
	require.NoError(t, err)
	routerCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "PerPartyRouter" {
				routerCid = e.Created.ContractId
				break
			}
		}
	}
	require.NotEmpty(t, routerCid)
	t.Logf("Created PerPartyRouter for receiver: %s", routerCid)

	// Build Message
	// Build token transfer (5 AMT in Splice Decimal format)
	encodedTokenTransfer := buildTokenTransferV1(tc.tokenAmount, remotePoolAddress, remoteTokenAddress, hashedInstrumentId, partyReceiver, tc.sourcePoolData)

	// Build message
	msg := &MessageV1{
		SourceChainSelector: remoteSelector,
		DestChainSelector:   env.Chain.ChainSelector(),
		SequenceNumber:      1,
		ExecutionGasLimit:   200000,
		CCIPReceiveGasLimit: 100000,
		Finality:            finalityConfigFromBlockConfirmations(2000),
		CCVAndExecutorHash:  [32]byte{},
		OnRampAddress:       gethcommon.LeftPadBytes(hexutil.MustDecode("0xf6eced5e96fff2de4f0ecd722beb57556fc443fd"), 32),
		OffRampAddress:      offRampAddress.InstanceAddress().Bytes(),
		Sender:              hexToBytes("0000000000000000000000000000000000000003"),
		Receiver:            EncodePartyID(partyReceiver),
		DestBlob:            []byte{},
		TokenTransfer:       encodedTokenTransfer,
		MessageData:         []byte{},
	}
	encodedMessage, err := EncodeMessageV1(msg)
	require.NoError(t, err)
	encodedMessageHex := hex.EncodeToString(encodedMessage)
	messageHash := crypto.Keccak256(encodedMessage)
	messageHashHex := hex.EncodeToString(messageHash)
	t.Logf("Message hash: %s", messageHashHex)

	// CCV Verification
	// Generate verifierResults with 2 of 3 signatures
	verifierResults, err := GenerateVerifierResults(encodedMessage, ccvSignerKeys[:2])
	require.NoError(t, err)
	verifierResultsHex := hex.EncodeToString(verifierResults)

	// Deploy CCIPReceiver for receiver
	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccipreceiver"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "receiverFinalityConfig", Value: finalityConfigValueFromBlockConfirmations(2000)},
						{Label: "requiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "optionalCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "optionalThreshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
					}},
				}},
			}},
			ActAs: []string{partyReceiver},
		},
	})
	require.NoError(t, err)
	ccipReceiverCid := extractCreatedContractId(res)
	t.Logf("Deployed CCIPReceiver: %s", ccipReceiverCid)

	// Capture receiver's balance before execute
	receiverHoldingsBefore, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceBefore := getHoldingsBalance(receiverHoldingsBefore)

	tokenPoolAddressEDS, err := edsTesthelpers.GetTokenPoolForToken(t.Context(), ccipAPIClient, hashedInstrumentId)
	require.NoError(t, err)
	ccipExecuteDisclosure, err := edsTesthelpers.GetCCIPExecuteDisclosure(t.Context(), ccipAPIClient, encodedMessageHex)
	require.NoError(t, err)
	ccvExecuteDisclosure, err := edsTesthelpers.GetCCVExecuteDisclosure(t.Context(), ccvAPIClient, encodedMessageHex, committeeVerifierAddress.InstanceAddress())
	require.NoError(t, err)
	tokenPoolDisclosure, err := edsTesthelpers.GetTokenPoolExecuteDisclosure(t.Context(), tokenPoolAPIClient, encodedMessageHex, tokenPoolAddressEDS.InstanceAddress())
	require.NoError(t, err)

	executeArgs := ccipreceiver.Execute2{
		Context:        ccipExecuteDisclosure.ChoiceContext,
		RouterCid:      types.CONTRACT_ID(routerCid),
		EncodedMessage: types.TEXT(encodedMessageHex),
		TokenTransfer: &ccipreceiver.TokenTransferInput{
			TokenPoolCid:       types.CONTRACT_ID(tokenPoolDisclosure.ContractId),
			TokenReceiverParty: types.PARTY(partyReceiver),
			TokenInput:         tokenPoolDisclosure.TokenInput,
			PoolExtraContext:   tokenPoolDisclosure.ChoiceContext,
		},
		CcvInputs: []ccipreceiver.CCVInput{
			{
				CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
				VerifierResults: types.TEXT(verifierResultsHex),
				CcvExtraContext: ccvExecuteDisclosure.ChoiceContext,
			},
		},
	}

	// CCIPReceiver.Execute: PrepareExecute + CCV + Pool Verify + Execute + Release
	// in one receiver-authored transaction with disclosed shared dependencies.
	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId:     ccipReceiverCid,
					Choice:         "Execute",
					ChoiceArgument: ledger.MapToValue(executeArgs),
				}},
			}},
			ActAs: []string{partyReceiver},
			DisclosedContracts: slices.Concat(
				tokenPoolDisclosure.DisclosedContracts,
				ccipExecuteDisclosure.DisclosedContracts,
				ccvExecuteDisclosure.DisclosedContracts,
			),
		},
	})
	require.NoError(t, err)
	t.Log("CCIPReceiver.Execute completed")

	// Check for pending AmuletTransferInstruction (Amulet transfers without receiver preapproval)
	var pendingTransferInstructionCid string
	var releaseCompleted bool
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			templateName := e.Created.GetTemplateId().GetEntityName()
			if templateName == "AmuletTransferInstruction" {
				pendingTransferInstructionCid = e.Created.GetContractId()
			}
			if templateName == "Amulet" || templateName == "LockedAmulet" {
				releaseCompleted = true
			}
		}
	}

	// If pending, receiver must accept the TransferInstruction to complete the transfer
	if !releaseCompleted && pendingTransferInstructionCid != "" {
		time.Sleep(500 * time.Millisecond)
		err := testhelpers.AcceptPendingTransferInstruction(t.Context(), receiverParticipant, transferInstructionClient, partyReceiver, pendingTransferInstructionCid)
		require.NoError(t, err)
	}

	// Verify receiver's balance increased by the expected transfer amount
	receiverHoldingsAfter, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceAfter := getHoldingsBalance(receiverHoldingsAfter)

	actualTransferAmount := receiverBalanceAfter - receiverBalanceBefore
	require.InDelta(t, tc.expectedTransferAmount, actualTransferAmount, 0.01, "Receiver balance should increase by transfer amount")
	t.Logf("Receiver balance: %.2f -> %.2f AMT (transferred %.2f)", receiverBalanceBefore, receiverBalanceAfter, actualTransferAmount)

	if tc.expectedDefaultLimiterTokens != "" {
		defaultRateLimiter, err := findActiveRateLimiterByInstanceID(t.Context(), tokenPoolOwnerParticipant, inboundRateLimiterAddr.InstanceID())
		require.NoError(t, err)
		require.Equal(t, tc.expectedDefaultLimiterTokens, getRateLimiterTokens(defaultRateLimiter))
	}
	if tc.expectedCustomLimiterTokens != "" {
		customRateLimiter, err := findActiveRateLimiterByInstanceID(t.Context(), tokenPoolOwnerParticipant, inboundCustomRateLimiterAddr.InstanceID())
		require.NoError(t, err)
		require.Equal(t, tc.expectedCustomLimiterTokens, getRateLimiterTokens(customRateLimiter))
	}

	t.Logf("✅ Success")
}

// getHoldingsBalance returns total balance across all holdings
func getHoldingsBalance(holdings []*apiv2.ActiveContract) float64 {
	var total float64
	for _, h := range holdings {
		balanceStr := h.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		balance, _ := new(big.Float).SetString(balanceStr)
		balanceFloat, _ := balance.Float64()
		total += balanceFloat
	}

	return total
}

// extractCreatedContractId returns the first created contract ID from a transaction response
func extractCreatedContractId(res *apiv2.SubmitAndWaitForTransactionResponse) string {
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			return e.Created.ContractId
		}
	}

	return ""
}

// buildTokenTransferV1 builds the encoded token transfer bytes for the message.
// tokenReceiver is keccak256-hashed to match Daml encodePartyAddress.
func buildTokenTransferV1(
	amount *big.Int,
	sourcePoolAddress,
	sourceTokenAddress []byte,
	destTokenAddress contracts.EncodedInstrumentID,
	tokenReceiverParty string,
	extraData []byte,
) *TokenTransferV1 {
	return &TokenTransferV1{
		Amount:             amount,
		SourcePoolAddress:  sourcePoolAddress,
		SourceTokenAddress: sourceTokenAddress,
		DestTokenAddress:   destTokenAddress.Bytes(),
		TokenReceiver:      EncodePartyID(tokenReceiverParty),
		ExtraData:          extraData,
	}
}

func encodeUint256Bytes(value uint64) []byte {
	encoded := make([]byte, 32)
	new(big.Int).SetUint64(value).FillBytes(encoded)

	return encoded
}

func findActiveRateLimiterByInstanceID(ctx context.Context, participant canton.Participant, instanceID string) (*apiv2.ActiveContract, error) {
	rateLimiters, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter",
	})
	if err != nil {
		return nil, err
	}

	for i := len(rateLimiters) - 1; i >= 0; i-- {
		if getRateLimiterInstanceID(rateLimiters[i]) == instanceID {
			return rateLimiters[i], nil
		}
	}

	return nil, fmt.Errorf("rate limiter %s not found", instanceID)
}

func getRateLimiterInstanceID(rateLimiter *apiv2.ActiveContract) string {
	return rateLimiter.GetCreatedEvent().GetCreateArguments().GetFields()[0].GetValue().GetText()
}

func getRateLimiterTokens(rateLimiter *apiv2.ActiveContract) string {
	return rateLimiter.GetCreatedEvent().GetCreateArguments().GetFields()[9].GetValue().GetNumeric()
}
