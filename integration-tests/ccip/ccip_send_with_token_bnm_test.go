package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipsender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccvs"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/client"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/common"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_burn_mint_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
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
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"

	// Import to register adapters
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"
)

// TestBnMTokenPool_FullSendFlow tests full send flow with token transfer.
// Validates LockOrBurn deducts proportional feeBps from encoded token amount.
//
//nolint:paralleltest // We can't run this test in parallel as that would mix up the holding calculations
func TestBnMTokenPool_FullSendFlow(t *testing.T) {
	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	ccipParticipant := env.Chain.Participants[0]
	senderParticipant := env.Chain.Participants[1]

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	// Upload DARs
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	onRampDar, err := contracts.GetDar(contracts.CCIPOnRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	feeQuoterDar, err := contracts.GetDar(contracts.CCIPFeeQuoter, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err)
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err)
	rmnDar, err := contracts.GetDar(contracts.CCIPRMN, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipSenderDar, err := contracts.GetDar(contracts.CCIPSender, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipExecutorDar, err := contracts.GetDar(contracts.CCIPExecutor, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPBurnMintTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	linkTokenDar, err := contracts.GetDar(contracts.Link, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipTestDar, err := contracts.GetDar(contracts.CCIPTest, contracts.CurrentVersion)
	require.NoError(t, err)

	dars := [][]byte{commonDar, offRampDar, onRampDar, feeQuoterDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar, rmnDar, ccipSenderDar, ccipExecutorDar, tokenPoolDar, linkTokenDar, ccipTestDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, senderParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.PartyID
	partySender := senderParticipant.PartyID
	t.Logf("Parties: CCIP=%s, Sender=%s", partyCCIP, partySender)

	// CCV Setup
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
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	// Create Scan and Registry API clients
	// Using the scanProxy endpoint of the 0-th participant, all participants are able to forward requests using the BFT Scan Proxy, it doesn't matter which one we use
	scanProxyClient, tokenMetadataClient, transferInstructionClient, err := testhelpers.NewValidatorAPIClients(ccipParticipant)
	require.NoError(t, err, "Failed to create validator API clients")

	// Setup Amulet token as fee token
	// Get registry admin for Amulet tokens
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err, "failed to get registry admin")

	// Native is Amulet
	nativeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}
	linkInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(partyCCIP),
		Id:    "ChainLink",
	}
	hashedLinkInstrumentId := contracts.EncodeInstrumentID(linkInstrumentId)

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

	tokenPriceExponentUSDCents := uint64(1e6)                 // 6 decimals for USD cents
	tokenPriceExponentUSD := 1e2 * tokenPriceExponentUSDCents // FeeQuoter usd8 scale: $1 == 1e8

	// Deploy Chain contracts
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
						Template: ccvs.CommitteeVerifier{
							Owner:                        types.PARTY(partyCCIP),
							CcipOwner:                    types.PARTY(partyCCIP),
							VersionTag:                   types.TEXT(versionTag),
							MessageSentObservers:         nil,
							StorageLocations:             []types.TEXT{"ipfs://test-send"},
							StorageLocationsAdmin:        types.PARTY(partyCCIP),
							PendingStorageLocationsAdmin: types.PARTY(partyCCIP),
							Deps:                         ccvs.CommitteeVerifierDeps{}, // Set by sequence
						},
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
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10)),
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						CcipOwner:      "", // Populated by the sequence
						RmnOwner:       types.PARTY(partyCCIP),
						CursedSubjects: nil,
					},
				},
				FeeQuoterConfig: sequences.FeeQuoterParams{
					Template: feequoter.FeeQuoter{
						PriceUpdaters: []types.PARTY{types.PARTY(partyCCIP)},
					},

					USDPerNative: big.NewInt(int64(1 * tokenPriceExponentUSD)), // $1
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
	addresses := cldfEnv.DataStore.Addresses().Filter()
	for i, address := range addresses {
		t.Logf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s, Labels=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier, address.Labels.String())
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
	executorRef, executorAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), executor.ContractType, executor.Version, devenvcommon.DefaultExecutorQualifier)
	require.NoError(t, err, "failed to get Executor address")

	// Deploy and configure lane for outbound sends

	// 8 cents outgoing CCV verification fee
	ccvFeeUSDCents := 7

	feeQuoterDestChainConfig := lanes.FeeQuoterDestChainConfig{
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         binary.BigEndian.Uint32([]byte{0x28, 0x12, 0xd5, 0x2c}),
		DefaultTokenFeeUSDCents:     10,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          25,
		V2Params: &lanes.FeeQuoterV2Params{
			LinkFeeMultiplierPercent: 100, // Not used when paying in native
			// Integer scaled by 1e10 in cantonFeeQuoterUSDPerUnitGas (38 -> 0.0000000038).
			USDPerUnitGas: big.NewInt(38),
		},
	}

	deployLaneLegReport, err := cld_ops.ExecuteSequence(cldfEnv.OperationsBundle, sequences.ConfigureLaneLegAsSourceWithInput, cldfEnv.BlockChains, sequences.ConfigureLaneLegInput{
		DataStore: cldfEnv.DataStore,
		Lane: lanes.UpdateLanesInput{
			Source: &lanes.ChainDefinition{
				Selector: env.Chain.ChainSelector(),
				CommitteeVerifiers: []lanes.CommitteeVerifierConfig[datastore.AddressRef]{
					{
						CommitteeVerifier: []datastore.AddressRef{committeeVerifierRef},
						RemoteChains: map[uint64]lanes.CommitteeVerifierRemoteChainConfig{
							remoteSelector: {
								AllowlistEnabled:          false,
								AddedAllowlistedSenders:   nil,
								RemovedAllowlistedSenders: nil,
								FeeUSDCents:               uint16(ccvFeeUSDCents),
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
				LaneMandatedOutboundCCVs: []datastore.AddressRef{committeeVerifierRef},
				DefaultOutboundCCVs:      nil,
				CantonLaneConfig: &lanes.CantonLaneConfig{
					GlobalConfig: globalConfigRef,
				},
				DefaultExecutor: executorRef,
				FeeQuoter:       feeQuoterAddress.InstanceAddress().Bytes(),
				OnRamp:          onRampAddress.InstanceAddress().Bytes(),
				OffRamp:         offRampAddress.InstanceAddress().Bytes(),
			},
			Dest: &lanes.ChainDefinition{
				Selector:                 remoteSelector,
				AddressBytesLength:       20,
				FeeQuoterDestChainConfig: feeQuoterDestChainConfig,
				ExecutorDestChainConfig: lanes.ExecutorDestChainConfig{
					USDCentsFee: 50,
					Enabled:     true,
				},
				OnRamp:  hexutil.MustDecode("0xf6eced5e96fff2de4f0ecd722beb57556fc443fd"),
				OffRamp: hexutil.MustDecode("0xd8c9ec8cad3fb34aeca3ddbebfabe9f28a9bfaed"),
				Router:  hexutil.MustDecode("0xe3ddcb2fde5d27a33c450fddc54a3f9bb2ecaa9f"),
			},
			IsDisabled:   false,
			TestRouter:   false,
			ExtraConfigs: lanes.ExtraConfigs{},
		},
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
		OwnerParty: types.PARTY(partyCCIP),
		Template: link.LinkRegistry{
			RegistryAdmin:        types.PARTY(partyCCIP),
			RegistryInstrumentId: linkInstrumentId,
			RegistryMeta:         splice_api_token_metadata_v1.Metadata{},
			TransferPreapprovals: nil,
		},
	})
	require.NoError(t, err)
	linkRegistryAddress, err := contracts.RawInstanceAddressFromString(linkTokenOut.Output.Labels.List()[0])
	require.NoError(t, err)

	// Setup token pool for outbound token transfer in Send.
	poolInstanceID := "test-pool-send"
	outboundRateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployOutbound, env.Chain, contractops.DeployInput[common.RateLimiter]{
		OwnerParty: types.PARTY(partyCCIP),
		Template: common.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceID),
			PoolOwner:           types.PARTY(partyCCIP),
			RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
			Direction:           common.RateLimitDirectionRateLimitDirection_Outbound,
			Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
			IsEnabled:           true,
			Capacity:            types.NUMERIC("10000000000"),
			Rate:                types.NUMERIC("10000000000"),
			Tokens:              types.NUMERIC("10000000000"),
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
	})
	require.NoError(t, err, "failed to deploy outbound rate limiter")
	outboundRateLimiterRawAddr := outboundRateLimiterOut.Output.Labels.List()[0]
	outboundRateLimiterAddr, err := contracts.RawInstanceAddressFromString(outboundRateLimiterRawAddr)
	require.NoError(t, err, "failed to parse outbound rate limiter raw address")

	// Pool transfer amounts use the token's smallest units; with Amulet's 10 token
	// decimals, a transfer amount of 100 means 100 local units and a 5% bps fee
	// floors to 5.
	tokenTransferFeeUSDCents := 10
	tokenTransferFeeBps := 500 // 5%

	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")
	out, err = changesets.DeployBurnMintTokenPool{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.DeployBurnMintTokenPoolConfig]{
		ChainSelector: env.Chain.ChainSelector(),
		Participant:   1,
		Config: changesets.DeployBurnMintTokenPoolConfig{
			CcipOwner:    partyCCIP,
			PoolOwner:    partyCCIP,
			InstrumentId: linkInstrumentId,
			Decimals:     10,
			InstanceID:   poolInstanceID,
			RemoteChainConfigs: map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
				types.NUMERIC(strconv.FormatUint(remoteSelector, 10)): burnminttokenpool.RemoteChainConfig{
					RemotePools:        []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
					RemoteTokenAddress: types.TEXT(hex.EncodeToString(remoteTokenAddress)),
					InboundCCVs:        []mcms.RawInstanceAddress{},
					OutboundCCVs:       []mcms.RawInstanceAddress{},
					FinalityConfig:     common.FinalityConfig{WaitForFinality: &types.UNIT{}},
					InboundRateLimiter: outboundRateLimiterAddr.Binding(),
					InboundCustomBlockConfirmationsRateLimiter: outboundRateLimiterAddr.Binding(),
					OutboundRateLimiter:                        outboundRateLimiterAddr.Binding(),
				},
			},
			// Set a custom token transfer fee config
			TokenTransferFeeConfigs: map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{
				types.NUMERIC(strconv.FormatUint(remoteSelector, 10)): {
					IsEnabled:         types.BOOL(true),
					DestGasOverhead:   types.INT64(25_000),
					DestBytesOverhead: types.INT64(32),
					FeeUSDCents:       types.NUMERIC(strconv.Itoa(tokenTransferFeeUSDCents)),
					FeeBps:            types.NUMERIC(strconv.Itoa(tokenTransferFeeBps)),
				},
			},
			PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
			TransferTimeout: burnminttokenpool.TransferTimeout{
				RelativeHours: func(v types.INT64) *types.INT64 { return &v }(types.INT64(24)),
			},
			Deps: burnminttokenpool.BurnMintTokenPoolDeps{
				TokenAdminRegistry: tokenAdminRegistryAddress.Binding(),
				RmnRemote:          rmnRemoteAddress.Binding(),
				FeeQuoter:          feeQuoterAddress.Binding(),
			},
			// By setting the TAR address, the CS will automatically register the newly deployed pool with the TAR
			TokenAdminRegistryInstanceAddress: tokenAdminRegistryAddress.InstanceAddress(),
		},
	})
	require.NoError(t, err, "failed to deploy burn mint token pool via changeset")
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
			ExecutorAPIConfig: config.ExecutorAPIConfig{
				Enabled: true,
				Executors: []config.Executor{
					{
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         partyCCIP,
							InstanceAddress: executorAddress.InstanceAddress(),
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
							Party:           new(partyCCIP),
							InstanceAddress: new(linkRegistryAddress.InstanceAddress()),
						},
					},
				},
			},
			TokenStandardAPIConfig: config.TokenStandardAPIConfig{
				Enabled: true,
				Admin:   partyCCIP,
				Registries: map[string]config.Registry{
					linkRegistryAddress.InstanceAddress().Hex(): {
						ContractIdentifier: config.ContractIdentifier{
							PartyID:         partyCCIP,
							InstanceAddress: linkRegistryAddress.InstanceAddress(),
						},
						TokenType: config.TokenTypeLINK,
						TokenId:   "LINK",
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
	globalAPIClient, err := oapiGlobal.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create GlobalConfig API client")
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create CCIP API client")
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create CCV API client")
	tokenPoolAPIClient, err := oapiTokenPool.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create Token Pool API client")
	executorAPIClient, err := oapiExecutor.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create Executor API client")

	// wait for EDS to start up
	time.Sleep(1 * time.Second)

	// Create PerPartyRouter for sender via EDS
	perPartyRouterFactoryDisclosure, err := edsTesthelpers.GetPerPartyRouterFactoryDisclosure(t.Context(), ccipAPIClient, partySender)
	require.NoError(t, err)
	t.Logf("disclosed contracts: %v", perPartyRouterFactoryDisclosure.DisclosedContracts)
	t.Logf("per party router factory cid: %s", perPartyRouterFactoryDisclosure.ContractId)

	res, err := senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#" + perpartyrouter.PackageName, ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: perPartyRouterFactoryDisclosure.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partySender}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-receiver"}}},
					}}}},
				}},
			}},
			ActAs:              []string{partySender},
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
	t.Logf("Created PerPartyRouter for sender: %s", routerCid)

	// Build test payload
	testPayload := []byte("Hello CCIP - this is a test send message!")
	testPayloadHex := hex.EncodeToString(testPayload)
	t.Logf("Test payload: %s", string(testPayload))

	// Deploy CCIPSender for sender
	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccipsender"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partySender}}},
					}},
				}},
			}},
			ActAs: []string{partySender},
		},
	})
	require.NoError(t, err)
	ccipSenderCid := extractCreatedContractId(res)
	t.Logf("Deployed CCIPSender: %s", ccipSenderCid)

	// Prepare receiver address (destination party encoded as keccak256)
	receiver := hexutil.MustDecode("0xcf8def9adfe3dd90b3dffe42c8eabbf7cd4ee6ca")
	receiverHex := hex.EncodeToString(receiver)

	// Fund separate holdings for fee payment and token transfer input.
	// The mint amount here follows the existing test's usd8-sized quantity setup.
	feeTokenHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, strconv.Itoa(100*int(tokenPriceExponentUSD)))
	require.NoError(t, err, "failed to mint Amulet tokens to sender")
	t.Logf("Minted fee-token Amulet holding to sender, Holding CID: %s", feeTokenHoldingCid)

	// Mint LINK
	linkRegistryCid, err := contractops.FindActiveContractIDByInstanceAddress(t.Context(), ccipParticipant.LedgerServices.State, []string{partyCCIP}, contracts.TemplateIDFromBinding(link.LinkRegistry{}).String(), linkRegistryAddress.InstanceAddress())
	require.NoError(t, err)
	t.Logf("LinkRegistry ContractId: %v", linkRegistryCid)
	_, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-burn-mint-v1", ModuleName: "Splice.Api.Token.BurnMintV1", EntityName: "BurnMintFactory"},
					ContractId: linkRegistryCid,
					Choice:     "BurnMintFactory_BurnMint",
					ChoiceArgument: ledger.MapToValue(splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint{
						ExpectedAdmin:    types.PARTY(partyCCIP),
						InstrumentId:     linkInstrumentId,
						InputHoldingCids: nil,
						Outputs: []splice_api_token_burn_mint_v1.BurnMintOutput{
							{
								Owner:   types.PARTY(partyCCIP),
								Amount:  "100.0",
								Context: splice_api_token_metadata_v1.ChoiceContext{},
							},
						},
						ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{},
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	linkHoldings, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), ccipParticipant, contracts.TemplateIDFromBinding(link.LinkHolding{}).ToLedgerIdentifier())
	require.NoError(t, err)
	require.NotEmpty(t, linkHoldings, "expected at least one active LinkHolding after minting")
	t.Logf("Minted 100 LINK to ccipOwner. Current LINK holdings %v", len(linkHoldings))
	for _, holding := range linkHoldings {
		t.Logf("- Holding CID: %s", holding.GetCreatedEvent().GetContractId())
	}
	// Transfer LINK to sender
	_, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferFactory"},
					ContractId: linkRegistryCid,
					Choice:     "TransferFactory_Transfer",
					ChoiceArgument: ledger.MapToValue(splice_api_token_transfer_instruction_v1.TransferFactoryTransfer{
						ExpectedAdmin: types.PARTY(partyCCIP),
						Transfer: splice_api_token_transfer_instruction_v1.Transfer{
							Sender:        types.PARTY(partyCCIP),
							Receiver:      types.PARTY(partySender),
							Amount:        "50.0",
							InstrumentId:  linkInstrumentId,
							RequestedAt:   types.TIMESTAMP(time.Now()),
							ExecuteBefore: types.TIMESTAMP(time.Now().Add(1 * time.Hour)),
							// Only using the last holding here - since that's the one that should have been minted just now
							InputHoldingCids: []types.CONTRACT_ID{types.CONTRACT_ID(linkHoldings[len(linkHoldings)-1].GetCreatedEvent().GetContractId())},
							Meta:             splice_api_token_metadata_v1.Metadata{},
						},
						ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{},
					}),
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	// Wait for transaction to propagate
	time.Sleep(time.Second * 5)
	linkTransferInstructions, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), senderParticipant, contracts.TemplateIDFromBinding(link.LinkTransferInstruction{}).ToLedgerIdentifier())
	require.NoError(t, err)
	require.NotEmpty(t, linkTransferInstructions, "expected at least one active LinkTransferInstruction after initiating transfer")
	t.Logf("Initiated transfer of 50 LINK to sender party. Current TransferInstructions %v", len(linkTransferInstructions))
	// Accept incoming transfer
	_, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
					ContractId: linkTransferInstructions[0].GetCreatedEvent().GetContractId(),
					Choice:     "TransferInstruction_Accept",
					ChoiceArgument: ledger.MapToValue(splice_api_token_transfer_instruction_v1.TransferInstructionAccept{
						ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{},
					}),
				}},
			}},
			ActAs: []string{partySender},
		},
	})
	require.NoError(t, err)
	senderLinkHoldings, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), senderParticipant, contracts.TemplateIDFromBinding(link.LinkHolding{}).ToLedgerIdentifier())
	require.NoError(t, err)
	t.Logf("Accepted transfer of 50 LINK to sender party. Current holdings:")
	senderHoldingCids := make([]types.CONTRACT_ID, len(linkHoldings))
	for i, holding := range senderLinkHoldings {
		t.Logf("- Holding CID: %s", holding.GetCreatedEvent().GetContractId())
		senderHoldingCids[i] = types.CONTRACT_ID(holding.GetCreatedEvent().GetContractId())
	}

	senderBalanceBefore, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)

	// Get transfer factory for Amulet tokens (sender to CCIP owner)
	transferFactoryCid, transferFactoryDisclosures, choiceContextRaw, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, registryAdmin, partySender, partyCCIP)
	require.NoError(t, err, "failed to get transfer factory")
	require.NotNil(t, choiceContextRaw, "choiceContext should not be nil")

	choiceContext, err := testhelpers.ChoiceContextFromData(choiceContextRaw)
	require.NoError(t, err, "failed to parse choice context")
	// Verify choiceContext has the expected structure (Record with "values" field)
	require.NotNil(t, choiceContext.GetRecord(), "choiceContext should be a Record")
	require.NotEmpty(t, choiceContext.GetRecord().GetFields(), "choiceContext should have fields")

	// Verify the "values" field exists and contains entries
	valuesField := choiceContext.GetRecord().GetFields()[0]
	require.Equal(t, "values", valuesField.GetLabel(), "first field should be 'values'")
	require.NotNil(t, valuesField.GetValue().GetTextMap(), "values field should be a TextMap")

	// Check if "amulet-rules" entry exists
	textMap := valuesField.GetValue().GetTextMap()
	hasAmuletRules := false
	for _, entry := range textMap.GetEntries() {
		if entry.GetKey() == "amulet-rules" {
			hasAmuletRules = true
			t.Logf("Found 'amulet-rules' entry in choiceContext")

			break
		}
	}
	if !hasAmuletRules {
		t.Logf("WARNING: 'amulet-rules' entry not found in choiceContext")
		t.Logf("choiceContext entries: %v", func() []string {
			entries := textMap.GetEntries()
			keys := make([]string, 0, len(entries))
			for _, entry := range entries {
				keys = append(keys, entry.GetKey())
			}

			return keys
		}())
	}

	// Extract transfer factory context values (e.g. amulet-rules) for the fee token input
	transferFactoryContextValues := testhelpers.ExtractChoiceContextValues(choiceContext)

	const tokenTransferAmountDecimal = "0.0000010000"

	senderLinkBalanceBefore, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &linkInstrumentId)
	require.NoError(t, err)

	executorRawOrHashedAddress := oapiCommon.RawOrHashedAddress{}
	_ = executorRawOrHashedAddress.FromRawInstanceAddress(executorAddress.String())
	msg := oapiCommon.Message{
		DestinationChainSelector: strconv.FormatUint(remoteSelector, 10),
		Executor: struct {
			Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiCommon.MessageExecutorType `json:"type"`
		}{
			Type:    oapiCommon.WithAddress,
			Address: &executorRawOrHashedAddress,
		},
		FeeToken: oapiCommon.InstrumentId{
			Admin: oapiCommon.PartyId(nativeInstrumentId.Admin),
			Id:    string(nativeInstrumentId.Id),
		},
		Payload:  "",
		Receiver: "",
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: tokenTransferAmountDecimal,
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(linkInstrumentId.Admin),
				Id:    string(linkInstrumentId.Id),
			},
		},
	}

	tokenPoolAddressEDS, err := edsTesthelpers.GetTokenPoolForToken(t.Context(), ccipAPIClient, hashedLinkInstrumentId)
	require.NoError(t, err)
	tokenPoolSendDisclosure, err := edsTesthelpers.GetTokenPoolSendDisclosure(t.Context(), tokenPoolAPIClient, msg, tokenPoolAddressEDS.InstanceAddress())
	require.NoError(t, err)
	ccipSendDisclosure, err := edsTesthelpers.GetCCIPSendDisclosure(t.Context(), ccipAPIClient, msg, nil, tokenPoolSendDisclosure.RequiredCCVs)
	require.NoError(t, err)
	ccvAddressEDS, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
	require.NoError(t, err)
	executorAddressEDS, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
	require.NoError(t, err)
	ccvSendDisclosure, err := edsTesthelpers.GetCCVSendDisclosure(t.Context(), ccvAPIClient, msg, ccvAddressEDS.InstanceAddress())
	require.NoError(t, err)
	executorSendDisclosure, err := edsTesthelpers.GetExecutorSendDisclosure(t.Context(), executorAPIClient, msg, executorAddressEDS.InstanceAddress(), ccipSendDisclosure.CCVs)
	require.NoError(t, err)

	// Sanity check - validate that all disclosures can be queried from the Global EDS API endpoint
	// TODO: these EDS-specific tests should be separated
	disclosedContracts, err := edsTesthelpers.GetGlobalDisclosureBatch(t.Context(), globalAPIClient, []contracts.InstanceAddress{
		perPartyRouterFactoryAddress.InstanceAddress(),
		globalConfigAddress.InstanceAddress(),
		feeQuoterAddress.InstanceAddress(),
		onRampAddress.InstanceAddress(),
		offRampAddress.InstanceAddress(),
		tokenAdminRegistryAddress.InstanceAddress(),
		rmnRemoteAddress.InstanceAddress(),
		committeeVerifierAddress.InstanceAddress(),
		executorAddress.InstanceAddress(),
		linkRegistryAddress.InstanceAddress(),
	})
	require.NoError(t, err)
	require.Lenf(t, disclosedContracts, 10, "expected to retrieve disclosures for all queried addresses")

	// Pool takes a token amount cut at LockOrBurn: feeBps = 500 (5%).
	// Message uses Decimal token amount 0.0000010000 → 10,000 smallest units;
	// after 5% pool fee the bridged amount should be 9,500.
	sendArgs := ccipsender.Send{
		Context:                  ccipSendDisclosure.ChoiceContext,
		RouterCid:                types.CONTRACT_ID(routerCid),
		DestinationChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		Message: ccipclient.Canton2AnyMessage{
			Receiver: types.TEXT(receiverHex),
			Payload:  types.TEXT(testPayloadHex),
			TokenTransfer: &ccipclient.TokenTransfer{
				Token:  linkInstrumentId,
				Amount: types.NUMERIC(tokenTransferAmountDecimal),
			},
			FeeToken: nativeInstrumentId,
			ExtraArgs: ccipclient.ExtraArgs{
				V3: &ccipclient.GenericExtraArgsV3{
					GasLimit: 0,
					Ccvs: []ccipclient.CCVExtraArg{
						{
							CcvAddress: committeeVerifierAddress.Binding(),
							CcvArgs:    types.TEXT(""),
						},
					},
					Executor: ccipclient.ExecutorExtraArg{
						ExecutorWithAddress: &ccipclient.ExecutorWithAddress{
							ExecutorAddress: executorAddress.Binding(),
							ExecutorArgs:    types.TEXT(""),
						},
					},
					TokenReceiver: types.TEXT(""),
					TokenArgs:     types.TEXT(""),
				},
			},
		},
		FeeTokenInput: ccipsender.FeeTokenInput{
			SenderInputCids:         []types.CONTRACT_ID{types.CONTRACT_ID(feeTokenHoldingCid)},
			FeeTokenConfigCid:       contractID(ccipSendDisclosure.FeeTokenConfigCid),
			FeeTokenTransferFactory: types.CONTRACT_ID(transferFactoryCid),
			FeeTokenExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
				Context: splice_api_token_metadata_v1.ChoiceContext{
					Values: transferFactoryContextValues,
				},
				Meta: splice_api_token_metadata_v1.Metadata{
					Values: map[string]types.TEXT{},
				},
			},
		},
		CcvSendInputs: []ccipsender.CCVSendInput{
			{
				CcvAddress:      ccvSendDisclosure.Address.Binding(),
				CcvCid:          types.CONTRACT_ID(ccvSendDisclosure.ContractId),
				CcvExtraContext: ccvSendDisclosure.ChoiceContext,
			},
		},
		TokenTransferInput: &ccipsender.TokenTransferInput{
			SenderInputCids:  senderHoldingCids,
			TokenPoolCid:     types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId),
			PoolExtraContext: tokenPoolSendDisclosure.ChoiceContext,
		},
		ExecutorInput: &ccipsender.ExecutorInput{
			ExecutorCid:          types.CONTRACT_ID(executorSendDisclosure.ContractId),
			ExecutorExtraContext: executorSendDisclosure.ChoiceContext,
		},
	}

	sendDisclosures := testhelpers.DeduplicateDisclosedContracts(slices.Concat(
		transferFactoryDisclosures,
		ccipSendDisclosure.DisclosedContracts,
		tokenPoolSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorSendDisclosure.DisclosedContracts,
	)...)
	quotedFee := quoteCCIPSenderFee(t, senderParticipant, partySender, ccipSenderCid, sendArgs, sendDisclosures)
	feeStr := strings.TrimSuffix(string(quotedFee.FeeTokenAmount), ".")
	poolFeeStr := strings.TrimSuffix(string(quotedFee.PoolFeeTokenAmount), ".")
	t.Logf("GetFee: feeTokenAmount=%s poolFeeTokenAmount=%s", feeStr, poolFeeStr)
	require.NotEqual(t, "0", feeStr, "GetFee should return a positive fee")
	require.NotEqual(t, "0", poolFeeStr, "GetFee should return a positive pool fee")

	// quoteCCIPSenderFee uses TRANSACTION_SHAPE_LEDGER_EFFECTS; refresh EDS disclosures so the
	// follow-up Send exercises current contract witnesses (avoids LOCAL_VERDICT_INACTIVE_CONTRACTS).
	tokenPoolSendDisclosure, err = edsTesthelpers.GetTokenPoolSendDisclosure(t.Context(), tokenPoolAPIClient, msg, tokenPoolAddressEDS.InstanceAddress())
	require.NoError(t, err)
	ccipSendDisclosure, err = edsTesthelpers.GetCCIPSendDisclosure(t.Context(), ccipAPIClient, msg, nil, tokenPoolSendDisclosure.RequiredCCVs)
	require.NoError(t, err)
	ccvAddressEDS, err = contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
	require.NoError(t, err)
	executorAddressEDS, err = contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
	require.NoError(t, err)
	ccvSendDisclosure, err = edsTesthelpers.GetCCVSendDisclosure(t.Context(), ccvAPIClient, msg, ccvAddressEDS.InstanceAddress())
	require.NoError(t, err)
	executorSendDisclosure, err = edsTesthelpers.GetExecutorSendDisclosure(t.Context(), executorAPIClient, msg, executorAddressEDS.InstanceAddress(), ccipSendDisclosure.CCVs)
	require.NoError(t, err)
	sendArgs.Context = ccipSendDisclosure.ChoiceContext
	sendArgs.FeeTokenInput.FeeTokenConfigCid = contractID(ccipSendDisclosure.FeeTokenConfigCid)
	sendArgs.CcvSendInputs[0].CcvCid = types.CONTRACT_ID(ccvSendDisclosure.ContractId)
	sendArgs.CcvSendInputs[0].CcvExtraContext = ccvSendDisclosure.ChoiceContext
	sendArgs.TokenTransferInput.TokenPoolCid = types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId)
	sendArgs.TokenTransferInput.PoolExtraContext = tokenPoolSendDisclosure.ChoiceContext
	sendArgs.ExecutorInput.ExecutorCid = types.CONTRACT_ID(executorSendDisclosure.ContractId)
	sendArgs.ExecutorInput.ExecutorExtraContext = executorSendDisclosure.ChoiceContext
	sendDisclosures = testhelpers.DeduplicateDisclosedContracts(slices.Concat(
		transferFactoryDisclosures,
		ccipSendDisclosure.DisclosedContracts,
		tokenPoolSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorSendDisclosure.DisclosedContracts,
	)...)

	// CCIPSender.Send: PrepareSend + CCV tickets + Send in one transaction.
	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					ContractId:     ccipSenderCid,
					Choice:         "Send",
					ChoiceArgument: ledger.MapToValue(sendArgs),
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: sendDisclosures,
		},
	})
	require.NoError(t, err)

	// Extract messageId from CCIPMessageSent to verify success
	var returnedMessageId string
	var returnedEncodedMessage string
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
				fields := e.Created.GetCreateArguments().GetFields()
				if len(fields) >= 5 {
					eventField := fields[4].GetValue().GetRecord()
					if eventField != nil {
						for _, field := range eventField.Fields {
							if field.GetLabel() == "messageId" {
								returnedMessageId = field.GetValue().GetText()
							}
							if field.GetLabel() == "encodedMessage" {
								returnedEncodedMessage = field.GetValue().GetText()
							}
						}
					}
				}

				break
			}
		}
	}
	require.NotEmpty(t, returnedMessageId, "CCIPMessageSent should be created")
	require.NotEmpty(t, returnedEncodedMessage, "CCIPMessageSent should contain encoded message")

	// Verify pool feeBps haircut: 10,000 smallest units with 5% feeBps => 9,500 bridged.
	require.Equal(t, int64(9500), extractTokenTransferAmountFromEncodedMessageHex(t, returnedEncodedMessage), "encoded token amount should be net after 5% feeBps")

	senderBalanceAfter, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	senderDelta := new(big.Rat).Sub(senderBalanceBefore, senderBalanceAfter)

	t.Logf(
		"Sender balance: before=%s after=%s deducted=%s",
		senderBalanceBefore,
		senderBalanceAfter,
		senderDelta,
	)
	require.Positive(t, senderDelta.Sign(), "sender balance should decrease after send")

	quotedFeeAmount, ok := new(big.Rat).SetString(feeStr)
	require.True(t, ok, "quoted fee should parse as a decimal value")
	// Fee token is Amulet (native); the bridged asset is LINK. Native holdings only pay
	// the CCIP fee — unlike ccip_send_with_token_lnr_test where fee and transfer share
	// the same instrument.
	expectedSenderDelta := new(big.Rat).Set(quotedFeeAmount)
	require.Zero(t, senderDelta.Cmp(expectedSenderDelta), "sender fee-token deduction should equal GetFee feeTokenAmount")

	senderLinkBalanceAfter, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &linkInstrumentId)
	require.NoError(t, err)
	linkDelta := new(big.Rat).Sub(senderLinkBalanceBefore, senderLinkBalanceAfter)
	tokenTransferAmountRat, ok := new(big.Rat).SetString(tokenTransferAmountDecimal)
	require.True(t, ok, "token transfer amount should parse as a decimal value")
	t.Logf(
		"Sender LINK: before=%s after=%s deducted=%s",
		senderLinkBalanceBefore,
		senderLinkBalanceAfter,
		linkDelta,
	)
	require.Positive(t, linkDelta.Sign(), "sender LINK balance should decrease after send")
	require.Zero(t, linkDelta.Cmp(tokenTransferAmountRat), "sender LINK deduction should equal message token transfer amount")

	t.Logf("Send completed")
	t.Logf("  Message ID: %s", returnedMessageId)
	t.Logf("  Original payload: %s", string(testPayload))

	t.Logf("✅ Success")
}
