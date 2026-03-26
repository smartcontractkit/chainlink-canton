package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	contractops "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"

	// Import to register adapters
	_ "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/adapters"

	_ "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
)

// TestCCIPSendWithTokenTransferFeeBps tests full send flow with token transfer.
// Validates LockOrBurn deducts proportional feeBps from encoded token amount.
//
//nolint:paralleltest // We can't run this test in parallel as that would mix up the holding calculations
func TestCCIPSendWithTokenTransferFeeBps(t *testing.T) {
	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	ccipParticipant := env.Chain.Participants[0]
	senderParticipant := env.Chain.Participants[1]

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
	lockReleaseTokenPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipTestDar, err := contracts.GetDar(contracts.CCIPTest, contracts.CurrentVersion)
	require.NoError(t, err)

	dars := [][]byte{commonDar, offRampDar, onRampDar, feeQuoterDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar, rmnDar, ccipSenderDar, ccipExecutorDar, lockReleaseTokenPoolDar, ccipTestDar}
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

	versionTag := "49ff34ed"
	ccvQualifier := devenvcommon.DefaultExecutorQualifier
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	// Create Scan and Registry API clients
	// Using the scanProxy endpoint of the 0-th participant, all participants are able to forward requests using the BFT Scan Proxy, it doesn't matter which one we use
	tokenSource := ccipParticipant.TokenSource
	interceptor := func(ctx context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		return nil
	}
	scanProxyClient, err := scanProxy.NewClientWithResponses(ccipParticipant.Endpoints.ValidatorAPIURL, scanProxy.WithRequestEditorFn(interceptor))
	require.NoError(t, err, "Failed to create scan proxy client")
	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", ccipParticipant.Endpoints.ValidatorAPIURL), tokenMetadataV1.WithRequestEditorFn(interceptor))
	require.NoError(t, err, "Failed to create token metadata client")
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", ccipParticipant.Endpoints.ValidatorAPIURL), transferInstructionV1.WithRequestEditorFn(interceptor))
	require.NoError(t, err, "Failed to create transfer instruction client")

	// Setup Amulet token as fee token
	// Get registry admin for Amulet tokens
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err, "failed to get registry admin")

	// Native is Amulet
	nativeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}

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

	tokenPriceExponentUSDCents := uint64(1e6) // 6 decimals for USDCents
	tokenPriceExponentUSD := 1e2 * tokenPriceExponentUSDCents

	// Deploy Chain contracts
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
								MinBlockConfirmations: 0,
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
					//nolint:gosec
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
	globalConfig, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	require.NoError(t, err, "failed to get GlobalConfig address")
	feeQuoter, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, ""))
	require.NoError(t, err, "failed to get FeeQuoter address")
	onRamp, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(onramp.ContractType), onramp.Version, ""))
	require.NoError(t, err, "failed to get OnRamp address")
	offRamp, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(offramp.ContractType), offramp.Version, ""))
	require.NoError(t, err, "failed to get OffRamp address")
	committeeVerifier, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(committee_verifier.ContractType), committee_verifier.Version, ccvQualifier))
	require.NoError(t, err, "failed to get CommitteeVerifier address")
	committeeVerifierRawAddr, err := contracts.RawInstanceAddressFromString(committeeVerifier.Labels.List()[0])
	require.NoError(t, err, "failed to parse CommitteeVerifier raw address")
	tokenAdminRegistry, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(token_admin_registry.ContractType), token_admin_registry.Version, ""))
	require.NoError(t, err, "failed to get TokenAdminRegistry address")
	rmnRemote, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	require.NoError(t, err, "failed to get RMNRemote address")
	executorAddress, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(executor.ContractType), executor.Version, devenvcommon.DefaultExecutorQualifier))
	require.NoError(t, err, "failed to get Executor address")

	// Deploy and configure lane for outbound sends
	cantonAdapter, ok := lanes.GetLaneAdapterRegistry().GetLaneAdapter(chainsel.FamilyCanton, semver.MustParse("2.0.0"))
	require.Truef(t, ok, "failed to get Canton Lane adapter")

	// 8 cents outgoing CCV verification fee
	ccvFeeUSDCents := 8

	feeQuoterDestChainConfig := lanes.FeeQuoterDestChainConfig{
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         binary.BigEndian.Uint32([]byte{0x28, 0x12, 0xd5, 0x2c}),
		DefaultTokenFeeUSDCents:     20,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          11,
		V2Params: &lanes.FeeQuoterV2Params{
			LinkFeeMultiplierPercent: 100,            // Not used, since we're paying in native
			USDPerUnitGas:            big.NewInt(38), // not used, since exec is disabled
		},
	}

	deployLaneLegReport, err := cld_ops.ExecuteSequence(cldfEnv.OperationsBundle, cantonAdapter.ConfigureLaneLegAsSource(), cldfEnv.BlockChains, lanes.UpdateLanesInput{
		Source: &lanes.ChainDefinition{
			Selector: env.Chain.ChainSelector(),
			CommitteeVerifiers: []lanes.CommitteeVerifierConfig[datastore.AddressRef]{
				{
					CommitteeVerifier: []datastore.AddressRef{committeeVerifier},
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
			LaneMandatedOutboundCCVs: []datastore.AddressRef{committeeVerifier},
			DefaultOutboundCCVs:      nil,
			CantonLaneConfig: &lanes.CantonLaneConfig{
				GlobalConfig: globalConfig,
			},
			DefaultExecutor: executorAddress,
			FeeQuoter:       contracts.HexToInstanceAddress(feeQuoter.Address).Bytes(),
			OnRamp:          contracts.HexToInstanceAddress(onRamp.Address).Bytes(),
			OffRamp:         contracts.HexToInstanceAddress(offRamp.Address).Bytes(),
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

	// Setup token pool for outbound token transfer in Send.
	tokenAdminRegistryRawAddr, err := contracts.RawInstanceAddressFromString(tokenAdminRegistry.Labels.List()[0])
	require.NoError(t, err, "failed to parse TokenAdminRegistry raw address")
	rmnRemoteRawAddr, err := contracts.RawInstanceAddressFromString(rmnRemote.Labels.List()[0])
	require.NoError(t, err, "failed to parse RMNRemote raw address")
	feeQuoterRawAddr, err := contracts.RawInstanceAddressFromString(feeQuoter.Labels.List()[0])
	require.NoError(t, err, "failed to parse FeeQuoter raw address")

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

	// Create TransferPreapproval to be set in the pool's PoolReceiveContext
	ccipOwnerHoldingCid, err := testhelpers.MintAMT(t.Context(), ccipParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partyCCIP, "100")
	require.NoError(t, err, "failed to mint AMT for CCIP owner")
	t.Logf("Minted 100 Amulet to ccipOwner, Holding CID: %s", ccipOwnerHoldingCid)
	preapprovalCid, err := testhelpers.CreateTransferPreapproval(t.Context(), ccipParticipant, scanProxyClient, partyCCIP, ccipOwnerHoldingCid)
	require.NoError(t, err, "failed to create preapproval")
	t.Logf("Created preapproval CID: %s", preapprovalCid)

	tokenTransferFeeUSDCents := 10
	tokenTransferFeeBps := 500 // 5%

	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")
	out, err = changesets.DeployTokenPool{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.DeployTokenPoolConfig]{
		ChainSelector: env.Chain.ChainSelector(),
		Participant:   1,
		Config: changesets.DeployTokenPoolConfig{
			CcipOwner:    partyCCIP,
			PoolOwner:    partyCCIP,
			InstrumentId: nativeInstrumentId,
			Decimals:     8,
			InstanceID:   poolInstanceID,
			Qualifier:    "test-pool-send",
			RemoteChainConfigs: types.GENMAP{
				strconv.FormatUint(remoteSelector, 10): lockreleasetokenpool.RemoteChainConfig{
					RemotePools:           []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
					RemoteTokenAddress:    types.TEXT(hex.EncodeToString(remoteTokenAddress)),
					InboundCCVs:           []mcms.RawInstanceAddress{},
					OutboundCCVs:          []mcms.RawInstanceAddress{},
					MinBlockConfirmations: types.INT64(0),
					InboundRateLimiter:    outboundRateLimiterAddr.Binding(),
					InboundCustomBlockConfirmationsRateLimiter: outboundRateLimiterAddr.Binding(),
					OutboundRateLimiter:                        outboundRateLimiterAddr.Binding(),
				},
			},
			// Set a custom token transfer fee config
			TokenTransferFeeConfigs: types.GENMAP{
				strconv.FormatUint(remoteSelector, 10): map[string]any{
					"isEnabled":         true,
					"destGasOverhead":   int64(25_000),
					"destBytesOverhead": int64(32),
					"feeUSDCents":       types.NUMERIC(strconv.Itoa(tokenTransferFeeUSDCents)),
					"feeBps":            types.NUMERIC(strconv.Itoa(tokenTransferFeeBps)),
				},
			},
			PoolReceiveContext: common.CCIPContext{Values: types.TEXTMAP{
				// This transfer preapproval must be specified for transfer to the pool's owner to be automatically be accepted.
				// The pool's EDS must also provide explicit disclosures for this preapproval contract, else the send will fail.
				"transfer-preapproval": common.AnyValue{AVContractId: &preapprovalCid},
			}},
			TransferTimeout: lockreleasetokenpool.TransferTimeout{
				RelativeHours: func(v types.INT64) *types.INT64 { return &v }(types.INT64(24)),
			},
			// By setting the TAR address, the CS will automatically register the newly deployed pool with the TAR
			TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistry.Address),
			Deps: lockreleasetokenpool.LockReleaseTokenPoolDeps{
				TokenAdminRegistry: tokenAdminRegistryRawAddr.Binding(),
				RmnRemote:          rmnRemoteRawAddr.Binding(),
				FeeQuoter:          feeQuoterRawAddr.Binding(),
			},
		},
	})
	require.NoError(t, err, "failed to deploy lock release token pool via changeset")
	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()

	// Create PerPartyRouter for sender
	var res *apiv2.SubmitAndWaitForTransactionResponse
	disclosedFactory, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	})
	require.NoError(t, err)

	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: disclosedFactory.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.CreateRouter{
						PartyOwner: types.PARTY(partySender),
						InstanceId: "router-sender",
					}),
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedFactory},
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
			CommandId: uuid.Must(uuid.NewUUID()).String(),
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

	// Get disclosures for CCIPSender.Send
	disclosedCCIPSender, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), senderParticipant, &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"})
	require.NoError(t, err)
	disclosedRouter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), senderParticipant, &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter"})
	require.NoError(t, err)
	disclosedOnRamp, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{PackageId: "#ccip-onramp", ModuleName: "CCIP.OnRamp", EntityName: "OnRamp"})
	require.NoError(t, err)
	disclosedGlobalConfig, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig"})
	require.NoError(t, err)
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"})
	require.NoError(t, err)
	disclosedRmnRemote, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote"})
	require.NoError(t, err)
	disclosedFeeQuoter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter"})
	require.NoError(t, err)
	disclosedPreapproval, err := testhelpers.GetDisclosedContractById(t.Context(), ccipParticipant, string(preapprovalCid))
	require.NoError(t, err)

	// Prepare receiver address (destination party encoded as keccak256)
	receiver := hexutil.MustDecode("0xcf8def9adfe3dd90b3dffe42c8eabbf7cd4ee6ca")
	receiverHex := hex.EncodeToString(receiver)

	require.NotEmpty(t, disclosedOnRamp.ContractId, "OnRamp disclosure missing/empty")
	require.NotEmpty(t, disclosedGlobalConfig.ContractId, "GlobalConfig disclosure missing/empty")
	require.NotEmpty(t, disclosedTar.ContractId, "TAR disclosure missing/empty")
	require.NotEmpty(t, disclosedRmnRemote.ContractId, "RMNRemote disclosure missing/empty")
	require.NotEmpty(t, disclosedFeeQuoter.ContractId, "FeeQuoter disclosure missing/empty")
	require.NotEmpty(t, disclosedPreapproval.ContractId, "Preapproval disclosure missing/empty")

	require.NotEmpty(t, disclosedCCIPSender.ContractId, "CCIPSender disclosure missing/empty")
	require.NotEmpty(t, disclosedRouter.ContractId, "Router disclosure missing/empty")

	t.Logf("disclosedCCIPSender.ContractId=%q", disclosedCCIPSender.ContractId)
	t.Logf("disclosedRouter.ContractId=%q", disclosedRouter.ContractId)

	t.Logf("disclosedOnRamp.ContractId=%q", disclosedOnRamp.ContractId)
	t.Logf("disclosedGlobalConfig.ContractId=%q", disclosedGlobalConfig.ContractId)
	t.Logf("disclosedTar.ContractId=%q", disclosedTar.ContractId)
	t.Logf("disclosedRmnRemote.ContractId=%q", disclosedRmnRemote.ContractId)
	t.Logf("disclosedFeeQuoter.ContractId=%q", disclosedFeeQuoter.ContractId)

	t.Logf("partySender=%q", partySender)

	// Mint 100 whole AMT in local 1e8 units so the sender can cover non-zero fees.
	feeTokenHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, strconv.Itoa(100*int(tokenPriceExponentUSD))) //nolint:gosec
	require.NoError(t, err, "failed to mint Amulet tokens to sender")
	t.Logf("Minted 100 whole Amulet tokens to sender, Holding CID: %s", feeTokenHoldingCid)
	tokenTransferHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, strconv.Itoa(100*int(tokenPriceExponentUSD))) //nolint:gosec
	require.NoError(t, err, "failed to mint Amulet tokens for token transfer")
	t.Logf("Minted token-transfer Amulet holding, CID: %s", tokenTransferHoldingCid)
	senderBalanceBefore := getHoldingsBalanceNumeric(t, t.Context(), senderParticipant)
	ccipOwnerBalanceBefore := getHoldingsBalanceNumeric(t, t.Context(), ccipParticipant)

	// Get disclosed contract for the fee token holding
	disclosedFeeTokenHolding, err := testhelpers.GetDisclosedContractById(t.Context(), senderParticipant, feeTokenHoldingCid)
	require.NoError(t, err, "failed to get disclosed contract for fee token holding")
	disclosedTokenTransferHolding, err := testhelpers.GetDisclosedContractById(t.Context(), senderParticipant, tokenTransferHoldingCid)
	require.NoError(t, err, "failed to get disclosed contract for token transfer holding")

	// Get disclosed CommitteeVerifier contract for CCV send inputs
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err, "failed to get disclosed CommitteeVerifier")
	disclosedPool, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-lockreleasetokenpool", ModuleName: "CCIP.LockReleaseTokenPool", EntityName: "LockReleaseTokenPool",
	})
	require.NoError(t, err, "failed to get disclosed LockReleaseTokenPool")
	disclosedOutboundRateLimiter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter",
	})
	require.NoError(t, err, "failed to get disclosed outbound RateLimiter")

	// Get transfer factory for Amulet tokens (sender to CCIP owner)
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, registryAdmin, partySender, partyCCIP)
	require.NoError(t, err, "failed to get transfer factory")
	require.NotNil(t, choiceContext, "choiceContext should not be nil")

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
	transferFactoryContextValues := make(types.TEXTMAP)
	if choiceContextRecord := choiceContext.GetRecord(); choiceContextRecord != nil && len(choiceContextRecord.Fields) > 0 {
		valuesField := choiceContextRecord.Fields[0]
		if valuesField.GetLabel() == "values" && valuesField.GetValue().GetTextMap() != nil {
			for _, entry := range valuesField.GetValue().GetTextMap().GetEntries() {
				if v := entry.GetValue().GetVariant(); v != nil {
					cid := types.CONTRACT_ID(v.GetValue().GetContractId())
					transferFactoryContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVContractId: &cid}
				} else if entry.GetValue().GetText() != "" {
					txt := types.TEXT(entry.GetValue().GetText())
					transferFactoryContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVText: &txt}
				}
			}
		}
	}

	extraArgs := splice_api_token_metadata_v1.ExtraArgs{
		Context: splice_api_token_metadata_v1.ChoiceContext{
			Values: transferFactoryContextValues,
		},
		Meta: splice_api_token_metadata_v1.Metadata{
			Values: types.TEXTMAP{},
		},
	}

	feeTokenInput := interfaces.TokenInput{
		TransferFactory:   types.CONTRACT_ID(transferFactoryCid),
		ExtraArgs:         extraArgs,
		TokenPoolHoldings: []types.CONTRACT_ID{},
	}

	tokenTransferFactoryCid, tokenTransferFactoryDisclosures, tokenTransferChoiceContext, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, registryAdmin, partySender, partySender)
	require.NoError(t, err, "failed to get token transfer factory")
	tokenTransferContextValues := make(types.TEXTMAP)
	if tokenCtxRecord := tokenTransferChoiceContext.GetRecord(); tokenCtxRecord != nil && len(tokenCtxRecord.Fields) > 0 {
		tokenValuesField := tokenCtxRecord.Fields[0]
		if tokenValuesField.GetLabel() == "values" && tokenValuesField.GetValue().GetTextMap() != nil {
			for _, entry := range tokenValuesField.GetValue().GetTextMap().GetEntries() {
				if v := entry.GetValue().GetVariant(); v != nil {
					cid := types.CONTRACT_ID(v.GetValue().GetContractId())
					tokenTransferContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVContractId: &cid}
				} else if entry.GetValue().GetText() != "" {
					txt := types.TEXT(entry.GetValue().GetText())
					tokenTransferContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVText: &txt}
				}
			}
		}
	}
	tokenTransferInput := interfaces.TokenInput{
		TransferFactory: types.CONTRACT_ID(tokenTransferFactoryCid),
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: splice_api_token_metadata_v1.ChoiceContext{Values: tokenTransferContextValues},
			Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
		},
		TokenPoolHoldings: []types.CONTRACT_ID{},
	}

	// Build the main Send context with CCIP contract IDs (matching execute test pattern)
	onRampCid := types.CONTRACT_ID(disclosedOnRamp.ContractId)
	globalConfigCid := types.CONTRACT_ID(disclosedGlobalConfig.ContractId)
	tarCid := types.CONTRACT_ID(disclosedTar.ContractId)
	feeQuoterCid := types.CONTRACT_ID(disclosedFeeQuoter.ContractId)
	rmnRemoteCid := types.CONTRACT_ID(disclosedRmnRemote.ContractId)
	sendContext := common.CCIPContext{
		Values: types.TEXTMAP{
			"on-ramp":              common.AnyValue{AVContractId: &onRampCid},
			"global-config":        common.AnyValue{AVContractId: &globalConfigCid},
			"token-admin-registry": common.AnyValue{AVContractId: &tarCid},
			"fee-quoter":           common.AnyValue{AVContractId: &feeQuoterCid},
			"rmn-remote":           common.AnyValue{AVContractId: &rmnRemoteCid},
		},
	}
	//nolint:gosec
	tokenTransferAmount := int64(10 * tokenPriceExponentUSD) // Transferring 10 units
	outboundRateLimiterContractID := types.CONTRACT_ID(disclosedOutboundRateLimiter.ContractId)

	sendArgs := ccipsender.Send{
		Context:                  sendContext,
		RouterCid:                types.CONTRACT_ID(routerCid),
		DestinationChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		Message: ccipclient.Canton2AnyMessage{
			Receiver: types.TEXT(receiverHex),
			Payload:  types.TEXT(testPayloadHex),
			TokenTransfer: &ccipclient.TokenTransferInput{
				Token:           nativeInstrumentId,
				Amount:          types.NUMERIC(strconv.FormatInt(tokenTransferAmount, 10)),
				SenderInputCids: []types.CONTRACT_ID{types.CONTRACT_ID(tokenTransferHoldingCid)},
				TokenPoolCid:    types.CONTRACT_ID(disclosedPool.ContractId),
				TokenInput:      tokenTransferInput,
				PoolExtraContext: common.CCIPContext{
					Values: types.TEXTMAP{
						"rate-limiter": common.AnyValue{AVContractId: &outboundRateLimiterContractID},
					},
				},
			},
			FeeToken: ccipclient.FeeTokenInput{
				Token:           nativeInstrumentId,
				TokenInput:      feeTokenInput,
				SenderInputCids: []types.CONTRACT_ID{types.CONTRACT_ID(feeTokenHoldingCid)},
			},
			ExtraArgs: ccipclient.ExtraArgs{
				V3: &ccipclient.GenericExtraArgsV3{
					GasLimit:           0,
					BlockConfirmations: 0,
					Ccvs: []mcms.RawInstanceAddress{
						{Unpack: types.TEXT(committeeVerifierRawAddr.String())},
					},
					Executor:      nil, // Not specifying the ExecutorInput is the same as the no-exec flag on EVM
					TokenReceiver: types.TEXT(""),
					TokenArgs:     types.TEXT(""),
				},
			},
		},
		Ccvs: []ccipclient.CCVSendInput{
			{
				CcvCid:          types.CONTRACT_ID(disclosedCCV.ContractId),
				VerifierArgs:    types.TEXT(""),
				CcvExtraContext: common.CCIPContext{},
			},
		},
	}

	ccipSendArgs := ledger.MapToValue(sendArgs)
	sendDisclosures := testhelpers.DeduplicateDisclosedContracts(slices.Concat(
		[]*apiv2.DisclosedContract{
			disclosedCCIPSender,
			disclosedRouter,
			disclosedOnRamp,
			disclosedGlobalConfig,
			disclosedTar,
			disclosedRmnRemote,
			disclosedFeeQuoter,
			disclosedFeeTokenHolding,
			disclosedTokenTransferHolding,
			disclosedCCV,
			disclosedPool,
			disclosedOutboundRateLimiter,
			disclosedPreapproval,
		},
		transferFactoryDisclosures,
		tokenTransferFactoryDisclosures,
	)...)

	// CCIPSender.Send: PrepareSend + CCV tickets + Send in one transaction
	time.Sleep(500 * time.Millisecond)
	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					ContractId:     ccipSenderCid,
					Choice:         "Send",
					ChoiceArgument: ccipSendArgs,
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: sendDisclosures,
		},
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			t.Logf("gRPC error: code=%s message=%s", s.Code(), s.Message())
			t.Logf("Error details: %v", s.Details())
		}
		require.NoError(t, err)
	}

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

	// Verify pool feeBps haircut was applied to token transfer amount in encoded message:
	// 5% feeBps deducted from the transferred amount
	require.Equal(t, int64(float64(tokenTransferAmount)*(1-(float64(tokenTransferFeeBps)/10_000))), extractTokenTransferAmountFromEncodedMessageHex(t, returnedEncodedMessage), "encoded token amount should be net after 5% feeBps")

	senderBalanceAfter := getHoldingsBalanceNumeric(t, t.Context(), senderParticipant)
	senderDelta := new(big.Float).Sub(senderBalanceBefore, senderBalanceAfter)

	// Validate that the sender was charged the correct total amount:
	// - CCV Fee
	// - (Network) Transfer Fee
	// - The amount of tokens that they sent
	// Execution fees should not be charged, as the message is sent with no-execution
	ccvFee := big.NewFloat(float64(ccvFeeUSDCents) * float64(tokenPriceExponentUSDCents))
	transferFee := big.NewFloat(float64(tokenTransferFeeUSDCents) * float64(tokenPriceExponentUSDCents))
	tokenTransferTotal := big.NewFloat(float64(tokenTransferAmount))
	expectedTokenDelta := new(big.Float).Add(new(big.Float).Add(ccvFee, transferFee), tokenTransferTotal)

	t.Logf(
		"Sender balance (local units): before=%s after=%s deducted=%s",
		senderBalanceBefore.String(),
		senderBalanceAfter.String(),
		senderDelta.String(),
	)
	assert.Equal(t, 0, expectedTokenDelta.Cmp(senderDelta), "Sender should be: %s, got: %s", expectedTokenDelta.String(), senderDelta.String())

	// Validate the CCIP Owner holdings - the amount credited should be exactly equal to the amount that the sender has paid, no tokens are lost.
	ccipOwnerBalanceAfter := getHoldingsBalanceNumeric(t, t.Context(), ccipParticipant)
	ccipOwnerDelta := new(big.Float).Sub(ccipOwnerBalanceAfter, ccipOwnerBalanceBefore)
	t.Logf(
		"CCIP owner balance (local units): before=%s after=%s credited=%s",
		ccipOwnerBalanceBefore.String(),
		ccipOwnerBalanceAfter.String(),
		ccipOwnerDelta.String(),
	)
	assert.Equal(t, 0, ccipOwnerDelta.Cmp(senderDelta), "Sender and Owner delta should be equal, got: ccipOwnerDelta: %s, senderDelta: %s", ccipOwnerDelta.String(), senderDelta.String())

	t.Logf("Send completed")
	t.Logf("  Message ID: %s", returnedMessageId)
	t.Logf("  Original payload: %s", string(testPayload))
}

func getHoldingsBalanceNumeric(t *testing.T, ctx context.Context, participant canton.Participant) *big.Float {
	t.Helper()

	holdings, err := testhelpers.ListActiveContractsByInterfaceId(ctx, participant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)

	total := big.NewFloat(0)
	for _, h := range holdings {
		views := h.GetCreatedEvent().GetInterfaceViews()
		if len(views) == 0 {
			continue
		}
		fields := views[0].GetViewValue().GetFields()
		if len(fields) < 3 {
			continue
		}
		amountStr := fields[2].GetValue().GetNumeric()
		amountBig, ok := new(big.Float).SetString(amountStr)
		require.Truef(t, ok, "failed to parse amount: %s", amountStr)
		total.Add(total, amountBig)
	}

	return total
}

// extractTokenTransferAmountFromEncodedMessageHex decodes encodedMessage and returns
// tokenTransfer.amount (uint256 big-endian) from the token-transfer payload.
// It skips the fixed CCIP message prefix, short variable fields, and destination blob,
// then reads the amount from bytes [1:33] of tokenTransfer (byte 0 is version/tag).
func extractTokenTransferAmountFromEncodedMessageHex(t *testing.T, encodedMessageHex string) int64 {
	t.Helper()

	b, err := hex.DecodeString(encodedMessageHex)
	require.NoError(t, err, "decode encodedMessage")

	i := 1 + 8 + 8 + 8 + 4 + 4 + 2 + 32
	for range 4 {
		i += 1 + int(b[i])
	}
	destBlobLen := int(b[i])<<8 | int(b[i+1])
	i += 2 + destBlobLen
	i += 2

	amt := int64(0)
	for _, x := range b[i+25 : i+33] {
		amt = (amt << 8) | int64(x)
	}

	return amt
}
