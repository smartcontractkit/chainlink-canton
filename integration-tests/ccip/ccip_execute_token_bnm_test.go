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
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkregistry"
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

type bnmTokenPoolReceiveFlowTestCase struct {
	tokenAmount                   *big.Int
	sourcePoolData                []byte
	expectedTransferAmount        float64
	defaultInboundLimiterCapacity string
	customInboundLimiterCapacity  string
	expectedDefaultLimiterTokens  string
	expectedCustomLimiterTokens   string
}

func TestBnMTokenPool_FullReceiveFlow(t *testing.T) {
	t.Parallel()

	runBnMTokenPoolReceiveFlowTest(t, bnmTokenPoolReceiveFlowTestCase{
		tokenAmount:                   big.NewInt(50e10),
		expectedTransferAmount:        50,
		defaultInboundLimiterCapacity: "1000000000000",
		customInboundLimiterCapacity:  "10000000000000",
	})
}

func runBnMTokenPoolReceiveFlowTest(t *testing.T, tc bnmTokenPoolReceiveFlowTestCase) {
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
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPBurnMintTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	linkTokenDar, err := contracts.GetDar(contracts.Link, contracts.CurrentVersion)
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
	dars := [][]byte{rmnDar, commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, tokenPoolDar, linkTokenDar, perPartyRouterDar, ccipReceiverDar, onRampDar, feeQuoterDar, executorDar}
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
	_, tokenMetadataClient, _, err := testhelpers.NewValidatorAPIClients(ccipParticipant)
	require.NoError(t, err, "Failed to create validator API clients")

	// Get DSO Admin Party (registry admin)
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err)

	// Instrument ID for AMT
	nativeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}
	linkInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(partyCCIP),
		Id:    "ChainLink", // TODO: we don't seem to have a standard name for link token instrumentID, do we?
	}
	hashedLinkInstrumentId := contracts.EncodeInstrumentID(linkInstrumentId)

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
				RMNOwnerParty:  partyCCIP,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Qualifier: ccvQualifier,
						Template: committeeverifier.CommitteeVerifier{
							Owner:                        types.PARTY(partyCCIP),
							CcipOwner:                    types.PARTY(partyCCIP),
							VersionTag:                   types.TEXT(versionTag),
							MessageSentObservers:         nil,
							StorageLocations:             []types.TEXT{"ipfs://test-receive"},
							StorageLocationsAdmin:        types.PARTY(partyCCIP),
							PendingStorageLocationsAdmin: types.PARTY(partyCCIP),
							Deps:                         committeeverifier.CommitteeVerifierDeps{},
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: core.GlobalConfig{
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
								AllowedFinalityConfig: core.FinalityConfig{WaitForFinality: &types.UNIT{}},
								CcvAllowlistEnabled:   false,
							},
							AllowedCCVs: nil,
						},
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: core.RMNRemote{
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

	// Deploy Link Token
	linkTokenOut, err := cld_ops.ExecuteOperation(bundle, linkregistry.Deploy, env.Chain, contractops.DeployInput[link.LinkRegistry]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: link.LinkRegistry{
			RegistryAdmin:        types.PARTY(partyTokenPoolOwner),
			RegistryInstrumentId: linkInstrumentId,
			RegistryMeta:         splice_api_token_metadata_v1.Metadata{},
		},
	})
	require.NoError(t, err)
	linkRegistryAddress, err := contracts.RawInstanceAddressFromString(linkTokenOut.Output.Labels.List()[0])
	require.NoError(t, err)

	// Token Pool Setup
	// Deploy default inbound RateLimiter required by ReleaseFromTicket receive flow.
	// Keep it enabled but lower-capacity so the test fails if the default-finality limiter
	// is selected for this FTF transfer instead of the custom-finality limiter.
	poolInstanceId := "test-pool-receive"
	inboundRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployInbound, env.Chain, contractops.DeployInput[core.RateLimiter]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: core.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceId),
			PoolOwner:           types.PARTY(partyTokenPoolOwner),
			RemoteChainSelector: types.NUMERIC(sourceChainSelector),
			Direction:           core.RateLimitDirectionRateLimitDirection_Inbound,
			Mode:                core.RateLimitModeRateLimitMode_DefaultFinality,
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

	inboundCustomRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployInbound, env.Chain, contractops.DeployInput[core.RateLimiter]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: core.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceId),
			PoolOwner:           types.PARTY(partyTokenPoolOwner),
			RemoteChainSelector: types.NUMERIC(sourceChainSelector),
			Direction:           core.RateLimitDirectionRateLimitDirection_Inbound,
			Mode:                core.RateLimitModeRateLimitMode_CustomFinality,
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

	outboundRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployOutbound, env.Chain, contractops.DeployInput[core.RateLimiter]{
		OwnerParty: types.PARTY(partyTokenPoolOwner),
		Template: core.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceId),
			PoolOwner:           types.PARTY(partyTokenPoolOwner),
			RemoteChainSelector: types.NUMERIC(sourceChainSelector),
			Direction:           core.RateLimitDirectionRateLimitDirection_Outbound,
			Mode:                core.RateLimitModeRateLimitMode_DefaultFinality,
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
	out, err = changesets.DeployBurnMintTokenPool{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.DeployBurnMintTokenPoolConfig]{
		ChainSelector: env.Chain.ChainSelector(),
		Participant:   0,
		Config: changesets.DeployBurnMintTokenPoolConfig{
			CcipOwner:          partyCCIP,
			PoolOwner:          partyTokenPoolOwner,
			InstrumentId:       linkInstrumentId,
			Decimals:           10,
			InstanceID:         poolInstanceId,
			PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
			TransferTimeout: burnminttokenpool.TransferTimeout{
				RelativeHours: func(v types.INT64) *types.INT64 { return &v }(types.INT64(24)),
			},
			RemoteChainConfigs: map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
				types.NUMERIC(sourceChainSelector): {
					RemotePools:        []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
					RemoteTokenAddress: types.TEXT(hex.EncodeToString(remoteTokenAddress)),
					InboundCCVs:        []chainlinkapi.RawInstanceAddress{},
					OutboundCCVs:       []chainlinkapi.RawInstanceAddress{},
					FinalityConfig: core.FinalityConfig{
						BlockDepth: new(types.INT64(2000)),
					},
					InboundRateLimiter:                         inboundRateLimiterAddr.Binding(),
					InboundCustomBlockConfirmationsRateLimiter: inboundCustomRateLimiterAddr.Binding(),
					OutboundRateLimiter:                        outboundRateLimiterAddr.Binding(),
				},
			},
			TokenTransferFeeConfigs: map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{},
			Deps: burnminttokenpool.BurnMintTokenPoolDeps{
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
	_, tokenPoolAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), burn_mint_token_pool.ContractType, burn_mint_token_pool.Version, "")
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
				TokenPools: map[string]config.TokenPool{
					tokenPoolAddress.InstanceAddress().Hex(): {
						Type: config.TokenPoolTypeBurnMint,
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         partyCCIP,
							InstanceAddress: tokenPoolAddress.InstanceAddress(),
						},
						PoolOwner: partyCCIP,
						// By setting the TokenStandard info, the Token Pool API will return the necessary factory disclosures
						BurnMintFactory: &config.BurnMintFactory{
							Type:            config.FactoryTypeAddress,
							TemplateId:      new(link.LinkRegistry{}.GetTemplateID()),
							Party:           new(partyTokenPoolOwner),
							InstanceAddress: new(linkRegistryAddress.InstanceAddress()),
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
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#" + ccipruntime.PackageName, ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
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
	// Build token transfer (5 LINK in Splice Decimal format)
	tokenTransfer := buildTokenTransferV1(tc.tokenAmount, remotePoolAddress, remoteTokenAddress, hashedLinkInstrumentId, partyReceiver, tc.sourcePoolData)

	// Build message
	msg := protocol.Message{
		Version:              1,
		SourceChainSelector:  protocol.ChainSelector(remoteSelector),
		DestChainSelector:    protocol.ChainSelector(env.Chain.ChainSelector()),
		SequenceNumber:       1,
		ExecutionGasLimit:    200000,
		CcipReceiveGasLimit:  100000,
		Finality:             protocol.NewFinality().WithBlockDepth(2000),
		CcvAndExecutorHash:   [32]byte{},
		OnRampAddress:        gethcommon.LeftPadBytes(gethcommon.HexToAddress("0xf6eced5e96fff2de4f0ecd722beb57556fc443fd").Bytes(), 32),
		OnRampAddressLength:  32,
		OffRampAddress:       offRampAddress.InstanceAddress().Bytes(),
		OffRampAddressLength: 32,
		Sender:               gethcommon.HexToAddress("0000000000000000000000000000000000000003").Bytes(),
		SenderLength:         20,
		Receiver:             contracts.HashedPartyFromString(partyReceiver).Bytes(),
		ReceiverLength:       32,
		DestBlob:             nil,
		DestBlobLength:       0,
		TokenTransfer:        tokenTransfer,
		Data:                 nil,
		DataLength:           0,
	}
	encodedMessage, err := msg.Encode()
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
			CommandId: uuid.NewString(),
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

	// Capture receiver's balance before execute (LINK for this receiver).
	receiverBalanceRatBefore, err := testhelpers.GetHoldingsBalance(t.Context(), receiverParticipant, &linkInstrumentId, testhelpers.WithHoldingOwner(partyReceiver))
	require.NoError(t, err)
	receiverBalanceBefore, _ := new(big.Float).SetRat(receiverBalanceRatBefore).Float64()

	tokenPoolAddressEDS, err := edsTesthelpers.GetTokenPoolForToken(t.Context(), ccipAPIClient, hashedLinkInstrumentId)
	require.NoError(t, err)
	ccipExecuteDisclosure, err := edsTesthelpers.GetCCIPExecuteDisclosure(t.Context(), ccipAPIClient, encodedMessageHex)
	require.NoError(t, err)
	ccvExecuteDisclosure, err := edsTesthelpers.GetCCVExecuteDisclosure(t.Context(), ccvAPIClient, encodedMessageHex, committeeVerifierAddress.InstanceAddress())
	require.NoError(t, err)
	tokenPoolDisclosure, err := edsTesthelpers.GetTokenPoolExecuteDisclosure(t.Context(), tokenPoolAPIClient, encodedMessageHex, tokenPoolAddressEDS.InstanceAddress())
	require.NoError(t, err)

	executeArgs := receiver.Execute{
		Context:        ccipExecuteDisclosure.ChoiceContext,
		RouterCid:      types.CONTRACT_ID(routerCid),
		EncodedMessage: types.TEXT(encodedMessageHex),
		TokenTransfer: &receiver.TokenTransferInput{
			TokenPoolCid:       types.CONTRACT_ID(tokenPoolDisclosure.ContractId),
			TokenReceiverParty: types.PARTY(partyReceiver),
			PoolExtraContext:   tokenPoolDisclosure.ChoiceContext,
		},
		CcvInputs: []receiver.CCVInput{
			{
				CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
				VerifierResults: types.TEXT(verifierResultsHex),
				CcvExtraContext: ccvExecuteDisclosure.ChoiceContext,
			},
		},
	}

	// CCIPReceiver.Execute: PrepareExecute + CCV + Pool Verify + Execute + Release
	// in one receiver-authored transaction with disclosed shared dependencies.
	_, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
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

	// Verify receiver's balance increased by the expected transfer amount
	receiverBalanceRatAfter, err := testhelpers.GetHoldingsBalance(t.Context(), receiverParticipant, &linkInstrumentId, testhelpers.WithHoldingOwner(partyReceiver))
	require.NoError(t, err)
	receiverBalanceAfter, _ := new(big.Float).SetRat(receiverBalanceRatAfter).Float64()

	actualTransferAmount := receiverBalanceAfter - receiverBalanceBefore
	require.InDelta(t, tc.expectedTransferAmount, actualTransferAmount, 0.01, "Receiver balance should increase by transfer amount")
	t.Logf("Receiver balance: %.2f -> %.2f LINK (transferred %.2f)", receiverBalanceBefore, receiverBalanceAfter, actualTransferAmount)

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
