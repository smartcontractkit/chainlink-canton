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
	damlledger "github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/clientapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/events"
	executorBinding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/executor"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	registryrl "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_burn_mint_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
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
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"

	// Import to register adapters
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"
)

// TestRegistryBurnMintTokenPool_FullSendFlow proves that a registry-pools
// BurnMintTokenPool deployed via the atomic Initialize choice (CreateAndExercise,
// no factory, no separate register step) is fully usable in a real CCIPSender.Send
// - the same LockOrBurn/VerifyOutboundCCVs/CalculateFee path a production pool goes
// through, driven by exactly the rate limiter and TAR registration Initialize itself
// produced.
//
// EDS's TokenPoolAPI does not know about the registry-pools package (it is hardcoded
// to the production burnminttokenpool binding), so the pool-specific choice context
// and disclosures are hand-built here instead of fetched through EDS - mirroring
// exactly what eds/internal/api/tokenpool/tokenpool.go's burnMintTokenPoolSend
// computes for the production pool. Everything else (CCIP core deploy, lane config,
// EDS for the CCIP/CCV/Executor side, PerPartyRouter, CCIPSender, Send) is unchanged
// from TestBnMTokenPool_FullSendFlow.
//
//nolint:paralleltest // We can't run this test in parallel as that would mix up the holding calculations
func TestRegistryBurnMintTokenPool_FullSendFlow(t *testing.T) {
	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	ccipParticipant := env.Chain.Participants[0]
	senderParticipant := env.Chain.Participants[1]

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	// Upload DARs
	runtimeDar, err := contracts.GetDar(contracts.CCIPRuntimeV2, contracts.DevVersion)
	require.NoError(t, err)
	coreDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifierV2, contracts.DevVersion)
	require.NoError(t, err)
	ccipSenderDar, err := contracts.GetDar(contracts.CCIPSenderV2, contracts.DevVersion)
	require.NoError(t, err)
	ccipExecutorDar, err := contracts.GetDar(contracts.CCIPExecutorV2, contracts.DevVersion)
	require.NoError(t, err)
	registryTokenPoolDar, err := contracts.GetDar(contracts.CCIPRegistryBurnMintTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	registryRateLimiterDar, err := contracts.GetDar(contracts.CCIPRegistryRateLimiterV2, contracts.DevVersion)
	require.NoError(t, err)
	linkTokenDar, err := contracts.GetDar(contracts.Link, contracts.DevVersion)
	require.NoError(t, err)

	dars := [][]byte{runtimeDar, coreDar, committeeVerifierDar, ccipSenderDar, ccipExecutorDar, registryTokenPoolDar, registryRateLimiterDar, linkTokenDar}
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
		Admin: types.PARTY(partySender),
		Id:    "ChainLink",
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

	tokenPriceExponentUSDCents := uint64(1e6)                 // 6 decimals for USD cents
	tokenPriceExponentUSD := 1e2 * tokenPriceExponentUSDCents // FeeQuoter usd8 scale: $1 == 1e8

	// Deploy Chain contracts (GlobalConfig, RMNRemote, FeeQuoter, CommitteeVerifier, Executor, TAR)
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
							StorageLocations:             []types.TEXT{"ipfs://test-send"},
							StorageLocationsAdmin:        types.PARTY(partyCCIP),
							PendingStorageLocationsAdmin: types.PARTY(partyCCIP),
							Deps:                         committeeverifier.CommitteeVerifierDeps{}, // Set by sequence
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
								AllowedFinalityConfig: ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
								CcvAllowlistEnabled:   false,
							},
							AllowedCCVs: nil,
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: core.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10)),
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: core.RMNRemote{
						CcipOwner:      "", // Populated by the sequence
						RmnOwner:       types.PARTY(partyCCIP),
						CursedSubjects: nil,
					},
				},
				FeeQuoterConfig: sequences.FeeQuoterParams{
					Template: core.FeeQuoter{
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
		ParticipantIndex: 1,
		OwnerParty:       types.PARTY(partySender),
		Template: link.LinkRegistry{
			RegistryAdmin:        types.PARTY(partySender),
			RegistryInstrumentId: linkInstrumentId,
			RegistryMeta:         splice_api_token_metadata_v1.Metadata{},
		},
	})
	require.NoError(t, err)
	linkRegistryAddress, err := contracts.RawInstanceAddressFromString(linkTokenOut.Output.Labels.List()[0])
	require.NoError(t, err)
	linkRegistryCid, err := contractops.FindActiveContractIDByInstanceAddress(t.Context(), senderParticipant.LedgerServices.State, []string{partySender}, contracts.TemplateIDFromBinding(link.LinkRegistry{}).String(), linkRegistryAddress.InstanceAddress())
	require.NoError(t, err)
	t.Logf("LinkRegistry ContractId: %v", linkRegistryCid)

	// ------------------------------------------------------------------------
	// Deploy the registry-pools BurnMintTokenPool via a single atomic
	// CreateAndExercise + Initialize - deploy, propose+accept admin, deploy the
	// lane's 3 rate limiters, wire the lane, and register with TAR all in one
	// ledger command. This replaces the production pool's separate
	// DeployBurnMintTokenPool + RegisterTokenPool changesets.
	// ------------------------------------------------------------------------
	poolInstanceID := "test-registry-pool-send"

	// Pool transfer amounts use the token's smallest units; with 10 pool decimals,
	// a transfer amount of 1.0 means 1e10 smallest units, and a 5% bps fee floors
	// to 5e8.
	tokenTransferFeeUSDCents := 10
	tokenTransferFeeBps := 500 // 5%

	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")

	pool := registrybnm.BurnMintTokenPool{
		InstanceId:   types.TEXT(poolInstanceID),
		PoolOwner:    types.PARTY(partySender),
		CcipOwner:    types.PARTY(partyCCIP),
		InstrumentId: linkInstrumentId,
		Decimals:     10,
		Observers:    nil,
		// Set directly at creation - Initialize only manages remoteChainConfigs/rate
		// limiters/TAR registration, not the token-transfer fee overrides.
		TokenTransferFeeConfigs: map[types.NUMERIC]registrybnm.TokenTransferFeeConfig2{
			types.NUMERIC(strconv.FormatUint(remoteSelector, 10)): {
				IsEnabled:         true,
				DestGasOverhead:   25_000,
				DestBytesOverhead: 32,
				FeeUSDCents:       types.NUMERIC(strconv.Itoa(tokenTransferFeeUSDCents)),
				FeeBps:            types.NUMERIC(strconv.Itoa(tokenTransferFeeBps)),
			},
		},
		PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{},
		TransferTimeout: registrybnm.TransferTimeout{
			RelativeHours: func(v types.INT64) *types.INT64 { return &v }(24),
		},
		Deps: registrybnm.BurnMintTokenPoolDeps{
			TokenAdminRegistry: tokenAdminRegistryAddress.Binding(),
			RmnRemote:          rmnRemoteAddress.Binding(),
			FeeQuoter:          feeQuoterAddress.Binding(),
		},
	}

	lane := registrybnm.LaneDeploySpec{
		RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		RemotePools:         []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
		RemoteTokenAddress:  types.TEXT(hex.EncodeToString(remoteTokenAddress)),
		InboundCCVs:         nil,
		OutboundCCVs:        nil,
		FinalityConfig:      ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
		Inbound: registrybnm.RateLimiterDeploySpec{
			InstanceId: "test-registry-pool-send-in",
			IsEnabled:  true,
			Capacity:   "10000000000",
			Rate:       "10000000000",
		},
		Outbound: registrybnm.RateLimiterDeploySpec{
			InstanceId: "test-registry-pool-send-out",
			IsEnabled:  true,
			Capacity:   "10000000000",
			Rate:       "10000000000",
		},
		InboundCustomFinality: registrybnm.RateLimiterDeploySpec{
			InstanceId: "test-registry-pool-send-in-custom",
			IsEnabled:  false,
			Capacity:   "0",
			Rate:       "0",
		},
	}

	// TAR lives on ccipParticipant (ccipOwner=partyCCIP), not senderParticipant.
	// Initialize's authorization only needs partySender's signature (self-issued
	// isAdmin path), but the submitting participant (senderParticipant) still needs
	// to *see* the TAR contract to build the transaction - that's a cross-participant
	// visibility gap, separate from authorization, so it must be disclosed explicitly.
	tarCid, err := contractops.FindActiveContractIDByInstanceAddress(t.Context(), ccipParticipant.LedgerServices.State, []string{partyCCIP}, contracts.TemplateIDFromBinding(core.TokenAdminRegistry{}).String(), tokenAdminRegistryAddress.InstanceAddress())
	require.NoError(t, err, "find TAR contract id from ccip participant")
	tarDisclosure, err := testhelpers.GetDisclosedContractById(t.Context(), ccipParticipant, tarCid)
	require.NoError(t, err, "disclose TAR contract to sender's participant")

	initArgs := registrybnm.Initialize{
		TokenAdminRegistryCid:  types.CONTRACT_ID(tarCid),
		ExistingTokenConfigCid: nil,
		Admin:                  types.PARTY(partySender),
		Lanes:                  []registrybnm.LaneDeploySpec{lane},
	}
	_, err = submitCreateAndExercise(t.Context(), senderParticipant, partySender, pool, "Initialize", initArgs, tarDisclosure)
	require.NoError(t, err, "create registry pool and exercise Initialize atomically")

	tokenPoolAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(poolInstanceID), types.PARTY(partySender))
	poolCid, err := contractops.FindActiveContractIDByInstanceAddress(t.Context(), senderParticipant.LedgerServices.State, []string{partySender}, pool.GetTemplateID(), tokenPoolAddress.InstanceAddress())
	require.NoError(t, err, "find registry pool contract id")
	outboundRateLimiterCid, err := contractops.FindActiveContractIDByInstanceAddress(t.Context(), senderParticipant.LedgerServices.State, []string{partySender}, registryrl.RateLimiter{}.GetTemplateID(), contracts.NewRawInstanceAddress(contracts.InstanceID(lane.Outbound.InstanceId), types.PARTY(partySender)).InstanceAddress())
	require.NoError(t, err, "find outbound rate limiter contract id")

	// Run EDS for the CCIP-global/CCV/Executor side only. There's no pool-side EDS
	// instance: TokenPoolAPI doesn't know the registry-pools template, and (per
	// eds/service/service.go) TokenStandardAPI's supply computation only works when
	// TokenPoolAPI is also enabled, since that's what starts the shared
	// instrumentHoldingStore's backfill - so a pool-only EDS instance here would
	// serve nothing. The pool's own send disclosure is hand-built below instead,
	// mirroring what EDS's burnMintTokenPoolSend computes.
	edsCCIPPort := freeport.GetOne(t)

	ccipEDSToken, _ := ccipParticipant.TokenSource.Token()
	edsCCIPCtx, edsCCIPCancel := context.WithCancel(t.Context())
	t.Cleanup(edsCCIPCancel)
	go func() {
		log.Info().Int("port", edsCCIPPort).Msg("Running EDS-ccip...")
		err := service.RunEDS(edsCCIPCtx, log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel), &config.Config{
			ChainSelector: strconv.FormatUint(env.Chain.ChainSelector(), 10),
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: uint16(edsCCIPPort),
			},
			Node: config.NodeConfig{
				URL: ccipParticipant.Endpoints.GRPCLedgerAPIURL,
				AuthConfig: commonconfig.AuthConfig{
					Type:   commonconfig.AuthTypeInsecureStatic,
					UserID: ccipParticipant.UserID,
					JWT:    ccipEDSToken.AccessToken,
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
				Enabled: false,
			},
			TokenStandardAPIConfig: config.TokenStandardAPIConfig{
				Enabled: false,
			},
		})
		log.Info().Err(err).Msg("EDS-ccip terminated")
		if !errors.Is(err, context.Canceled) {
			t.Errorf("EDS-ccip server exited with error: %v", err)
		}
	}()

	// Create EDS clients — CCIP global/CCV/executor on P0.
	globalAPIClient, err := oapiGlobal.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsCCIPPort))
	require.NoError(t, err, "Failed to create GlobalConfig API client")
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsCCIPPort))
	require.NoError(t, err, "Failed to create CCIP API client")
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsCCIPPort))
	require.NoError(t, err, "Failed to create CCV API client")
	executorAPIClient, err := oapiExecutor.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsCCIPPort))
	require.NoError(t, err, "Failed to create Executor API client")

	waitForEDSListening(t, edsCCIPPort)

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
					TemplateId: contracts.IdentifierFromBinding(ccipruntime.PerPartyRouterFactory{}),
					ContractId: perPartyRouterFactoryDisclosure.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partySender}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-receiver"}}},
						{Label: "feeTransferLifetime", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}},
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
					TemplateId: contracts.IdentifierFromBinding(sender.CCIPSender{}),
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

	// Prepare EVM receiver address
	receiver := hexutil.MustDecode("0xcf8def9adfe3dd90b3dffe42c8eabbf7cd4ee6ca")
	receiverHex := hex.EncodeToString(receiver)

	// Fund separate holdings for fee payment and token transfer input.
	feeTokenHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, "1000.0")
	require.NoError(t, err, "failed to mint Amulet tokens to sender")
	t.Logf("Minted fee-token Amulet holding to sender, Holding CID: %s", feeTokenHoldingCid)

	// Mint LINK to sender (pool owner / token admin)
	_, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-burn-mint-v1", ModuleName: "Splice.Api.Token.BurnMintV1", EntityName: "BurnMintFactory"},
					ContractId: linkRegistryCid,
					Choice:     "BurnMintFactory_BurnMint",
					ChoiceArgument: damlledger.MapToValue(splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint{
						ExpectedAdmin:    types.PARTY(partySender),
						InstrumentId:     linkInstrumentId,
						InputHoldingCids: nil,
						Outputs: []splice_api_token_burn_mint_v1.BurnMintOutput{
							{
								Owner:   types.PARTY(partySender),
								Amount:  "50.0",
								Context: splice_api_token_metadata_v1.ChoiceContext{},
							},
						},
						ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{},
					}),
				}},
			}},
			ActAs: []string{partySender},
		},
	})
	require.NoError(t, err)
	var senderLinkHoldings []*apiv2.ActiveContract
	require.Eventually(t, func() bool {
		var listErr error
		senderLinkHoldings, listErr = testhelpers.ListActiveContractsByTemplateId(t.Context(), senderParticipant, contracts.IdentifierFromBinding(link.LinkHolding{}))
		if listErr != nil {
			return false
		}

		return len(senderLinkHoldings) > 0
	}, 15*time.Second, 200*time.Millisecond, "expected at least one active LinkHolding after minting")
	t.Logf("Minted 50 LINK to sender. Current holdings: %v", len(senderLinkHoldings))
	senderHoldingCids := make([]types.CONTRACT_ID, len(senderLinkHoldings))
	for i, holding := range senderLinkHoldings {
		t.Logf("- Holding CID: %s", holding.GetCreatedEvent().GetContractId())
		senderHoldingCids[i] = types.CONTRACT_ID(holding.GetCreatedEvent().GetContractId())
	}

	// Get transfer factory for Amulet tokens (sender to CCIP owner)
	transferFactoryCid, transferFactoryDisclosures, choiceContextRaw, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, registryAdmin, partySender, partyCCIP)
	require.NoError(t, err, "failed to get transfer factory")
	require.NotNil(t, choiceContextRaw, "choiceContext should not be nil")

	choiceContext, err := testhelpers.ChoiceContextFromData(choiceContextRaw)
	require.NoError(t, err, "failed to parse choice context")
	transferFactoryContextValues := testhelpers.ExtractChoiceContextValues(choiceContext)

	const tokenTransferAmountDecimal = "1.0"

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
		Payload:  testPayloadHex,
		Sender:   partySender,
		Receiver: receiverHex,
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: tokenTransferAmountDecimal,
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(linkInstrumentId.Admin),
				Id:    string(linkInstrumentId.Id),
			},
		},
	}

	// Hand-built equivalent of edsTesthelpers.GetTokenPoolSendDisclosure, since EDS's
	// TokenPoolAPI is hardcoded to the production burnminttokenpool package and does
	// not know about the registry-pools template. Mirrors exactly what
	// eds/internal/api/tokenpool/tokenpool.go's burnMintTokenPoolSend computes: the
	// outbound rate limiter under "rate-limiter" and the BurnMintFactory (here,
	// LinkRegistry itself) under "burn-mint-factory". Our lane's OutboundCCVs is
	// empty, so RequiredCCVs is empty too.
	tokenPoolChoiceContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: map[string]splice_api_token_metadata_v1.AnyValue{
			string(registryrl.RateLimiterContextKey):      {AVContractId: func(v types.CONTRACT_ID) *types.CONTRACT_ID { return &v }(types.CONTRACT_ID(outboundRateLimiterCid))},
			string(registrybnm.BurnMintFactoryContextKey): {AVContractId: func(v types.CONTRACT_ID) *types.CONTRACT_ID { return &v }(types.CONTRACT_ID(linkRegistryCid))},
		},
	}
	var requiredCCVs []string

	ccipSendDisclosure, err := edsTesthelpers.GetCCIPSendDisclosure(t.Context(), ccipAPIClient, msg, nil, requiredCCVs)
	require.NoError(t, err)
	ccvAddressEDS, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
	require.NoError(t, err)
	executorAddressEDS, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
	require.NoError(t, err)
	ccvSendDisclosure, err := edsTesthelpers.GetCCVSendDisclosure(t.Context(), ccvAPIClient, msg, ccvAddressEDS.InstanceAddress())
	require.NoError(t, err)
	executorSendDisclosure, err := edsTesthelpers.GetExecutorSendDisclosure(t.Context(), executorAPIClient, msg, executorAddressEDS.InstanceAddress(), ccipSendDisclosure.CCVs)
	require.NoError(t, err)

	// Sanity check - CCIP global batch covers P0 contracts only (LinkRegistry/pool live on senderParticipant).
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
	})
	require.NoError(t, err)
	require.Lenf(t, disclosedContracts, 9, "expected to retrieve disclosures for all CCIP global addresses")

	// Pool takes a token amount cut at LockOrBurn: feeBps = 500 (5%).
	// Message uses Decimal token amount 0.0000010000 → 10,000 smallest units;
	// after 5% pool fee the bridged amount should be 9,500.
	sendArgs := sender.Send{
		Context:                  ccipSendDisclosure.ChoiceContext,
		RouterCid:                types.CONTRACT_ID(routerCid),
		DestinationChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		Message: clientapi.Canton2AnyMessage{
			Receiver: types.TEXT(receiverHex),
			Payload:  types.TEXT(testPayloadHex),
			TokenTransfer: &clientapi.TokenTransfer{
				Token:  linkInstrumentId,
				Amount: types.NUMERIC(tokenTransferAmountDecimal),
			},
			FeeToken: nativeInstrumentId,
			ExtraArgs: clientapi.ExtraArgs{
				V3: &clientapi.GenericExtraArgsV3{
					GasLimit: 0,
					Ccvs: []clientapi.CCVExtraArg{
						{
							CcvAddress: committeeVerifierAddress.Binding(),
							CcvArgs:    types.TEXT(""),
						},
					},
					Executor: clientapi.ExecutorExtraArg{
						ExecutorWithAddress: &clientapi.ExecutorWithAddress{
							ExecutorAddress: executorAddress.Binding(),
							ExecutorArgs:    types.TEXT(""),
						},
					},
					TokenReceiver: types.TEXT(""),
					TokenArgs:     types.TEXT(""),
				},
			},
		},
		FeeTokenInput: sender.FeeTokenInput{
			SenderInputCids:         []types.CONTRACT_ID{types.CONTRACT_ID(feeTokenHoldingCid)},
			FeeTokenConfigCid:       types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid),
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
		CcvSendInputs: []sender.CCVSendInput{
			{
				CcvAddress: ccvSendDisclosure.Address.Binding(),
				CcvCid:     types.CONTRACT_ID(ccvSendDisclosure.ContractId),
				Context:    ccvSendDisclosure.ChoiceContext,
			},
		},
		TokenTransferInput: &sender.TokenTransferInput{
			SenderInputCids: senderHoldingCids,
			TokenPoolCid:    types.CONTRACT_ID(poolCid),
			Context:         tokenPoolChoiceContext,
		},
		ExecutorInput: &sender.ExecutorInput{
			ExecutorCid: types.CONTRACT_ID(executorSendDisclosure.ContractId),
			Context:     executorSendDisclosure.ChoiceContext,
		},
	}

	sendDisclosures := testhelpers.DeduplicateDisclosedContracts(slices.Concat(
		transferFactoryDisclosures,
		ccipSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorSendDisclosure.DisclosedContracts,
	)...)
	quotedFee := quoteCCIPSenderFee(t, senderParticipant, partySender, ccipSenderCid, sendArgs, sendDisclosures)
	feeStr := string(quotedFee.FeeTokenAmount)
	poolFeeStr := string(quotedFee.PoolFeeTokenAmount)
	t.Logf("GetFee: feeTokenAmount=%s poolFeeTokenAmount=%s", feeStr, poolFeeStr)
	require.NotEqual(t, "0", feeStr, "GetFee should return a positive fee")
	require.NotEqual(t, "0", poolFeeStr, "GetFee should return a positive pool fee")

	// quoteCCIPSenderFee uses TRANSACTION_SHAPE_LEDGER_EFFECTS; refresh EDS disclosures so the
	// follow-up Send exercises current contract witnesses (avoids LOCAL_VERDICT_INACTIVE_CONTRACTS).
	ccipSendDisclosure, err = edsTesthelpers.GetCCIPSendDisclosure(t.Context(), ccipAPIClient, msg, nil, requiredCCVs)
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
	sendArgs.FeeTokenInput.FeeTokenConfigCid = types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid)
	sendArgs.CcvSendInputs[0].CcvCid = types.CONTRACT_ID(ccvSendDisclosure.ContractId)
	sendArgs.CcvSendInputs[0].Context = ccvSendDisclosure.ChoiceContext
	sendArgs.ExecutorInput.ExecutorCid = types.CONTRACT_ID(executorSendDisclosure.ContractId)
	sendArgs.ExecutorInput.Context = executorSendDisclosure.ChoiceContext
	sendDisclosures = testhelpers.DeduplicateDisclosedContracts(slices.Concat(
		transferFactoryDisclosures,
		ccipSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorSendDisclosure.DisclosedContracts,
	)...)

	senderBalanceBefore, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	senderLinkBalanceBefore, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &linkInstrumentId)
	require.NoError(t, err)

	// CCIPSender.Send: PrepareSend + CCV tickets + Send in one transaction.
	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     contracts.IdentifierFromBinding(sender.CCIPSender{}),
					ContractId:     ccipSenderCid,
					Choice:         "Send",
					ChoiceArgument: damlledger.MapToValue(sendArgs),
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: sendDisclosures,
		},
	})
	require.NoError(t, err)

	// Extract CCIPMessageSent event to verify success
	eventTemplateId := events.CCIPMessageSent{}.GetTemplateID()
	var ccipMessageSent *events.CCIPMessageSent
	for _, event := range res.GetTransaction().GetEvents() {
		if createdEvent := event.GetCreated(); createdEvent != nil {
			if templateId := createdEvent.GetTemplateId(); templateId != nil {
				gotTemplateId := fmt.Sprintf("#%s:%s:%s", createdEvent.GetPackageName(), templateId.GetModuleName(), templateId.GetEntityName())
				if gotTemplateId == eventTemplateId {
					ccipMessageSent, err = bindings.UnmarshalCreatedEvent[events.CCIPMessageSent](createdEvent)
					require.NoError(t, err)
					break
				}
			}
		}
	}
	require.NotNil(t, ccipMessageSent, "CCIPMessageSent event not found")
	require.NotEmpty(t, ccipMessageSent.Event.MessageId, "CCIPMessageSent should be created")
	require.NotEmpty(t, ccipMessageSent.Event.EncodedMessage, "CCIPMessageSent should contain encoded message")
	require.Equal(t, nativeInstrumentId, ccipMessageSent.Event.FeeToken, "CCIPMessageSent should contain feeToken InstrumentId")

	// Verify pool feeBps haircut: 1e10 LINK, minus 5% feeBps => 0.95e10 bridged.
	require.Equal(t, int64(9500000000), extractTokenTransferAmountFromEncodedMessageHex(t, ccipMessageSent.Event.EncodedMessage), "encoded token amount should be net after 5% feeBps")

	wantAmount, ok := new(big.Rat).SetString(tokenTransferAmountDecimal)
	require.True(t, ok, "token transfer amount should parse as a decimal value")
	gotAmount, ok := new(big.Rat).SetString(string(ccipMessageSent.Event.TokenAmountBeforeTokenPoolFees))
	require.True(t, ok, "token amount before token pool fees should parse as a decimal value")
	require.Equal(t, 0, wantAmount.Cmp(gotAmount), "token amount before token pool fees should match original transfer amount")

	quotedFeeAmount, ok := new(big.Rat).SetString(feeStr)
	require.True(t, ok, "quoted fee should parse as a decimal value")
	poolFeeAmountRat, ok := new(big.Rat).SetString(poolFeeStr)
	require.True(t, ok, "pool fee should parse as a decimal value")
	tokenTransferAmountRat, ok := new(big.Rat).SetString(tokenTransferAmountDecimal)
	require.True(t, ok, "token transfer amount should parse as a decimal value")
	// When poolOwner == sender, pool feeBps is retained by the sender's pool; only net LINK leaves total balance.
	netLinkTransferAmountRat := new(big.Rat).Sub(tokenTransferAmountRat, poolFeeAmountRat)

	require.Eventually(t, func() bool {
		senderBalanceAfter, balanceErr := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &nativeInstrumentId)
		if balanceErr != nil {
			return false
		}
		senderDelta := new(big.Rat).Sub(senderBalanceBefore, senderBalanceAfter)

		return senderDelta.Cmp(quotedFeeAmount) == 0
	}, 15*time.Second, 200*time.Millisecond, "sender fee-token deduction should equal GetFee feeTokenAmount")

	require.Eventually(t, func() bool {
		senderLinkBalanceAfter, balanceErr := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &linkInstrumentId)
		if balanceErr != nil {
			return false
		}
		linkDelta := new(big.Rat).Sub(senderLinkBalanceBefore, senderLinkBalanceAfter)

		return linkDelta.Cmp(netLinkTransferAmountRat) == 0
	}, 15*time.Second, 200*time.Millisecond, "sender LINK deduction should equal net token transfer amount after pool feeBps")

	t.Logf("Send completed")
	t.Logf("  Message ID: %s", ccipMessageSent.Event.MessageId)
	t.Logf("  Original payload: %s", string(testPayload))
	t.Logf("Success")
}
