package tests

import (
	"context"
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
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/clientapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
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
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	contractops "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
	rkccip "github.com/smartcontractkit/chainlink-canton/registry-kit/ccip"
	rkledger "github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"

	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"
)

const (
	registrySendInstrumentID   = "CCIP-REGISTRY-SEND"
	registrySendPoolInstanceID = "ccip-registry-send-pool"
	registrySendRLInstanceID   = "ccip-registry-send-rl-out"
	registrySendMintAmount     = "10.0"
)

// TestRegistryTokenPool_FullSendFlow exercises CCIPSender.Send with a Canton Registry Holding
// bridged through a registrar-owned BurnMintTokenPool (hybrid EDS: CCIP/CCV/Executor via EDS,
// Registry pool context assembled manually).
//
//nolint:paralleltest // Holding balance assertions require exclusive CTF env
func TestRegistryTokenPool_FullSendFlow(t *testing.T) {
	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	ccipParticipant := env.Chain.Participants[0]
	senderParticipant := env.Chain.Participants[1]

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	uploadRegistryDARs(t, ccipParticipant, senderParticipant)

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
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPBurnMintTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	ccipTestDar, err := contracts.GetDar(contracts.CCIPTest, contracts.CurrentVersion)
	require.NoError(t, err)

	dars := [][]byte{runtimeDar, coreDar, committeeVerifierDar, ccipSenderDar, ccipExecutorDar, tokenPoolDar, ccipTestDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, senderParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded CCIP DARs to all participants: %v", packageIds)

	partyCCIP := ccipParticipant.PartyID
	partySender := senderParticipant.PartyID
	partyRegistrar := testhelpers.AllocateParty(t, senderParticipant, "registry-send-registrar")
	testhelpers.GrantCanActAs(t, senderParticipant, partyRegistrar)
	t.Logf("Parties: CCIP=%s, Sender=%s, Registrar=%s", partyCCIP, partySender, partyRegistrar)

	ctx := t.Context()
	senderClient := rkledger.NewCTFClient(senderParticipant)
	ccipClient := rkledger.NewCTFClient(ccipParticipant)

	ccvSignerPubKeys := make([]string, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}

	versionTag := "e9a05a20"
	ccvQualifier := devenvcommon.DefaultCommitteeVerifierQualifier
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	scanProxyClient, tokenMetadataClient, transferInstructionClient, err := testhelpers.NewValidatorAPIClients(ccipParticipant)
	require.NoError(t, err)

	registryAdmin, err := testhelpers.GetRegistryAdmin(ctx, tokenMetadataClient)
	require.NoError(t, err)

	nativeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}
	registryInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(partyRegistrar),
		Id:    types.TEXT(registrySendInstrumentID),
	}
	hashedRegistryInstrumentId := contracts.EncodeInstrumentID(registryInstrumentId)

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), reporter)
	cldfEnv := cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{env.Chain}),
		OperationsBundle: bundle,
	}

	tokenPriceExponentUSDCents := uint64(1e6)
	tokenPriceExponentUSD := 1e2 * tokenPriceExponentUSDCents

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
							StorageLocations:             []types.TEXT{"ipfs://test-registry-send"},
							StorageLocationsAdmin:        types.PARTY(partyCCIP),
							PendingStorageLocationsAdmin: types.PARTY(partyCCIP),
							Deps:                         committeeverifier.CommitteeVerifierDeps{},
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
						CcipOwner:     "",
						ChainSelector: types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10)),
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: core.RMNRemote{
						CcipOwner:      "",
						RmnOwner:       types.PARTY(partyCCIP),
						CursedSubjects: nil,
					},
				},
				FeeQuoterConfig: sequences.FeeQuoterParams{
					Template: core.FeeQuoter{
						PriceUpdaters: []types.PARTY{types.PARTY(partyCCIP)},
					},
					USDPerNative: big.NewInt(int64(1 * tokenPriceExponentUSD)),
				},
				NativeInstrumentId: nativeInstrumentId,
			},
		},
	})
	require.NoError(t, err)
	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()

	globalConfigRef, globalConfigAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), global_config.ContractType, global_config.Version, "")
	require.NoError(t, err)
	_, feeQuoterAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), fee_quoter.ContractType, fee_quoter.Version, "")
	require.NoError(t, err)
	_, onRampAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), onramp.ContractType, onramp.Version, "")
	require.NoError(t, err)
	_, offRampAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), offramp.ContractType, offramp.Version, "")
	require.NoError(t, err)
	committeeVerifierRef, committeeVerifierAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), committee_verifier.ContractType, committee_verifier.Version, ccvQualifier)
	require.NoError(t, err)
	_, tokenAdminRegistryAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), token_admin_registry.ContractType, token_admin_registry.Version, "")
	require.NoError(t, err)
	_, rmnRemoteAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), rmn_remote.ContractType, rmn_remote.Version, "")
	require.NoError(t, err)
	_, perPartyRouterFactoryAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), per_party_router_factory.ContractType, per_party_router_factory.Version, "")
	require.NoError(t, err)
	executorRef, executorAddress, err := testhelpers.ResolveAddressFromDatastore(cldfEnv.DataStore, env.Chain.ChainSelector(), executor.ContractType, executor.Version, devenvcommon.DefaultExecutorQualifier)
	require.NoError(t, err)

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
			LinkFeeMultiplierPercent: 100,
			USDPerUnitGas:            big.NewInt(38),
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
								AllowlistEnabled:   false,
								FeeUSDCents:        uint16(ccvFeeUSDCents),
								GasForVerification: 50_000,
								PayloadSizeBytes:   6*64 + 2*32,
								SignatureConfig: lanes.CommitteeVerifierSignatureQuorumConfig{
									Signers:   ccvSignerPubKeys,
									Threshold: 2,
								},
							},
						},
					},
				},
				LaneMandatedOutboundCCVs: []datastore.AddressRef{committeeVerifierRef},
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
	require.NoError(t, err)
	runningDs := datastore.NewMemoryDataStore()
	for _, address := range deployLaneLegReport.Output.Addresses {
		err = runningDs.Addresses().Add(address)
		require.NoError(t, err)
	}
	err = runningDs.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = runningDs.Seal()

	bootstrap, err := registry.BootstrapServices(ctx, senderClient, partyRegistrar, registrySendInstrumentID)
	require.NoError(t, err)

	_, err = registry.MintViaAllocationFactory(ctx, senderClient, bootstrap, partySender, registrySendMintAmount)
	require.NoError(t, err)

	tokenTransferFeeUSDCents := 10
	tokenTransferFeeBps := 500
	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")

	now := time.Now()
	outboundRLAddr, err := rkccip.DeployOutboundRateLimiterForOwner(ctx, senderClient, partyRegistrar, ratelimiter.RateLimiter{
		InstanceId:          types.TEXT(registrySendRLInstanceID),
		PoolInstanceId:      types.TEXT(registrySendPoolInstanceID),
		PoolOwner:           types.PARTY(partyRegistrar),
		RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Outbound,
		Mode:                ratelimiter.RateLimitModeRateLimitMode_DefaultFinality,
		IsEnabled:           true,
		Capacity:            types.NUMERIC("10000000000"),
		Rate:                types.NUMERIC("10000000000"),
		Tokens:              types.NUMERIC("10000000000"),
		LastUpdated:         types.TIMESTAMP(now),
	})
	require.NoError(t, err)

	poolAddr, err := rkccip.DeployBurnMintPoolForOwner(ctx, senderClient, rkccip.PoolDeployDeps{
		CcipOwner:          partyCCIP,
		TokenAdminRegistry: tokenAdminRegistryAddress,
		RMNRemote:          rmnRemoteAddress,
		FeeQuoter:          feeQuoterAddress,
	}, partyRegistrar, registryInstrumentId, registrySendPoolInstanceID,
		map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
			types.NUMERIC(strconv.FormatUint(remoteSelector, 10)): {
				RemotePools:         []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
				RemoteTokenAddress:  types.TEXT(hex.EncodeToString(remoteTokenAddress)),
				FinalityConfig:      ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
				OutboundRateLimiter: outboundRLAddr.Binding(),
			},
		},
		map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{
			types.NUMERIC(strconv.FormatUint(remoteSelector, 10)): {
				IsEnabled:         true,
				DestGasOverhead:   types.INT64(25_000),
				DestBytesOverhead: types.INT64(32),
				FeeUSDCents:       types.NUMERIC(strconv.Itoa(tokenTransferFeeUSDCents)),
				FeeBps:            types.NUMERIC(strconv.Itoa(tokenTransferFeeBps)),
			},
		},
	)
	require.NoError(t, err)

	initialTarCID, err := contractops.FindActiveContractIDByInstanceAddress(
		ctx,
		ccipParticipant.LedgerServices.State,
		[]string{partyCCIP},
		core.TokenAdminRegistry{}.GetTemplateID(),
		tokenAdminRegistryAddress.InstanceAddress(),
	)
	require.NoError(t, err)

	tokenConfigCID, tarCID, err := rkccip.RegisterTokenPoolViaClient(ctx, ccipClient, rkccip.RegisterTokenPoolClientInput{
		TokenAdminRegistryCID: initialTarCID,
		InstrumentId:          registryInstrumentId,
		PoolInstanceID:        registrySendPoolInstanceID,
		CcipParty:             partyCCIP,
		PoolOwnerParty:        partyRegistrar,
		PoolOwnerClient:       senderClient,
	})
	require.NoError(t, err)

	_, err = rkccip.SetBurnMintFactory(ctx, senderClient, rkccip.SetBurnMintFactoryInput{
		TokenAdminRegistryCID: tarCID,
		TokenConfigCID:        tokenConfigCID,
		InstrumentId:          registryInstrumentId,
		BurnMintFactoryCID:    bootstrap.AllocationFactory,
		CcipParty:             partyCCIP,
		PoolOwnerParty:        partyRegistrar,
		CcipClient:            ccipClient,
		PoolOwnerClient:       senderClient,
	})
	require.NoError(t, err)

	edsParticipant := env.Chain.Participants[0]
	edsToken, _ := edsParticipant.TokenSource.Token()
	edsPort := freeport.GetOne(t)
	go func() {
		log.Info().Msg("Running EDS...")
		err := service.RunEDS(ctx, log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel), &config.Config{
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
			// TokenPool API omitted: hybrid send assembles Registry pool context manually.
			// Enabling it would register partyRegistrar with the instrument-holding store on the
			// CCIP participant, which cannot backfill holdings hosted on the sender participant.
		})
		log.Info().Err(err).Msg("EDS terminated")
		if !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Msg("EDS server exited with error")
			t.Fail()
		}
	}()

	globalAPIClient, err := oapiGlobal.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err)
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err)
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err)
	executorAPIClient, err := oapiExecutor.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	perPartyRouterFactoryDisclosure, err := edsTesthelpers.GetPerPartyRouterFactoryDisclosure(ctx, ccipAPIClient, partySender)
	require.NoError(t, err)

	res, err := senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#" + ccipruntime.PackageName, ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: perPartyRouterFactoryDisclosure.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partySender}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-registry-send"}}},
					}}}},
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: perPartyRouterFactoryDisclosure.DisclosedContracts,
		},
	})
	require.NoError(t, err)
	routerCid := extractCreatedContractId(res)
	require.NotEmpty(t, routerCid)

	testPayload := []byte("Hello CCIP Registry - full send test!")
	testPayloadHex := hex.EncodeToString(testPayload)

	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccipsender-registry"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partySender}}},
					}},
				}},
			}},
			ActAs: []string{partySender},
		},
	})
	require.NoError(t, err)
	ccipSenderCid := extractCreatedContractId(res)

	feeTokenHoldingCid, err := testhelpers.MintAMT(ctx, senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, strconv.Itoa(100*int(tokenPriceExponentUSD)))
	require.NoError(t, err)

	holdingRows, err := testhelpers.ListHoldingsForInstrument(ctx, senderParticipant, &registryInstrumentId,
		testhelpers.WithHoldingOwner(partySender), testhelpers.WithUnlockedHoldingsOnly())
	require.NoError(t, err)
	require.NotEmpty(t, holdingRows)
	senderHoldingCids := make([]types.CONTRACT_ID, len(holdingRows))
	for i, row := range holdingRows {
		senderHoldingCids[i] = types.CONTRACT_ID(row.ContractID)
	}

	transferFactoryCid, transferFactoryDisclosures, choiceContextRaw, err := testhelpers.GetTransferFactory(ctx, transferInstructionClient, registryAdmin, partySender, partyCCIP)
	require.NoError(t, err)
	choiceContext, err := testhelpers.ChoiceContextFromData(choiceContextRaw)
	require.NoError(t, err)
	transferFactoryContextValues := testhelpers.ExtractChoiceContextValues(choiceContext)

	const tokenTransferAmountDecimal = "0.0000010000"

	senderBalanceBefore, err := testhelpers.GetHoldingsBalance(ctx, senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	senderRegistryBalanceBefore, err := testhelpers.GetHoldingsBalance(ctx, senderParticipant, &registryInstrumentId, testhelpers.WithHoldingOwner(partySender))
	require.NoError(t, err)

	receiver := hexutil.MustDecode("0xcf8def9adfe3dd90b3dffe42c8eabbf7cd4ee6ca")
	receiverHex := hex.EncodeToString(receiver)

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
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: tokenTransferAmountDecimal,
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(registryInstrumentId.Admin),
				Id:    string(registryInstrumentId.Id),
			},
		},
	}

	poolSendDeps := registryPoolSendDeps{
		Client:                    senderClient,
		CcipClient:                ccipClient,
		RegistrarParty:            partyRegistrar,
		CcipParty:                 partyCCIP,
		Bootstrap:                 bootstrap,
		PoolInstanceID:            registrySendPoolInstanceID,
		RateLimiterInstanceID:     registrySendRLInstanceID,
		PoolAddress:               poolAddr,
		TokenAdminRegistryAddress: tokenAdminRegistryAddress,
		TokenAdminRegistryCID:     tarCID,
		RMNRemoteAddress:          rmnRemoteAddress,
	}

	buildSendBundle := func(enableResultContracts bool) (sender.Send, []*apiv2.DisclosedContract) {
		tokenPoolSendDisclosure := buildRegistryTokenPoolSendDisclosure(t, ctx, senderParticipant, ccipParticipant, ccipAPIClient, poolSendDeps, hashedRegistryInstrumentId, enableResultContracts)
		ccipSendDisclosure, err := edsTesthelpers.GetCCIPSendDisclosure(ctx, ccipAPIClient, msg, nil, tokenPoolSendDisclosure.RequiredCCVs)
		require.NoError(t, err)
		ccvAddressEDS, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
		require.NoError(t, err)
		executorAddressEDS, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
		require.NoError(t, err)
		ccvSendDisclosure, err := edsTesthelpers.GetCCVSendDisclosure(ctx, ccvAPIClient, msg, ccvAddressEDS.InstanceAddress())
		require.NoError(t, err)
		executorSendDisclosure, err := edsTesthelpers.GetExecutorSendDisclosure(ctx, executorAPIClient, msg, executorAddressEDS.InstanceAddress(), ccipSendDisclosure.CCVs)
		require.NoError(t, err)

		sendArgs := sender.Send{
			Context:                  ccipSendDisclosure.ChoiceContext,
			RouterCid:                types.CONTRACT_ID(routerCid),
			DestinationChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
			Message: clientapi.Canton2AnyMessage{
				Receiver: types.TEXT(receiverHex),
				Payload:  types.TEXT(testPayloadHex),
				TokenTransfer: &clientapi.TokenTransfer{
					Token:  registryInstrumentId,
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
					Meta: splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
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
				TokenPoolCid:    types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId),
				Context:         tokenPoolSendDisclosure.ChoiceContext,
			},
			ExecutorInput: &sender.ExecutorInput{
				ExecutorCid: types.CONTRACT_ID(executorSendDisclosure.ContractId),
				Context:     executorSendDisclosure.ChoiceContext,
			},
		}

		sendDisclosures := testhelpers.DeduplicateDisclosedContracts(slices.Concat(
			transferFactoryDisclosures,
			ccipSendDisclosure.DisclosedContracts,
			tokenPoolSendDisclosure.DisclosedContracts,
			ccvSendDisclosure.DisclosedContracts,
			executorSendDisclosure.DisclosedContracts,
		)...)

		return sendArgs, sendDisclosures
	}

	sendArgs, sendDisclosures := buildSendBundle(false)
	quotedFee := quoteCCIPSenderFee(t, senderParticipant, partySender, ccipSenderCid, sendArgs, sendDisclosures)
	feeStr := strings.TrimSuffix(string(quotedFee.FeeTokenAmount), ".")
	poolFeeStr := strings.TrimSuffix(string(quotedFee.PoolFeeTokenAmount), ".")
	require.NotEqual(t, "0", feeStr)
	require.NotEqual(t, "0", poolFeeStr)

	disclosedContracts, err := edsTesthelpers.GetGlobalDisclosureBatch(ctx, globalAPIClient, []contracts.InstanceAddress{
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
	require.Len(t, disclosedContracts, 9)

	var sendRes *apiv2.SubmitAndWaitForTransactionResponse
	var enableResultContracts bool
	for _, enable := range []bool{false, true} {
		sendArgs, sendDisclosures = buildSendBundle(enable)

		sendRes, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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
		if err == nil {
			enableResultContracts = enable
			break
		}
		t.Logf("CCIPSender.Send with enable-result-contracts=%v failed: %v", enable, err)
	}
	require.NoError(t, err, "CCIPSender.Send failed with both enable-result-contracts values")
	t.Logf("CCIPSender.Send succeeded with enable-result-contracts=%v", enableResultContracts)

	var returnedMessageId string
	var returnedEncodedMessage string
	for _, event := range sendRes.GetTransaction().GetEvents() {
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
	require.NotEmpty(t, returnedMessageId)
	require.NotEmpty(t, returnedEncodedMessage)
	require.Equal(t, int64(9500), extractTokenTransferAmountFromEncodedMessageHex(t, types.TEXT(returnedEncodedMessage)))

	senderBalanceAfter, err := testhelpers.GetHoldingsBalance(ctx, senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	senderDelta := new(big.Rat).Sub(senderBalanceBefore, senderBalanceAfter)
	require.Positive(t, senderDelta.Sign())

	quotedFeeAmount, ok := new(big.Rat).SetString(feeStr)
	require.True(t, ok)
	expectedSenderDelta := new(big.Rat).Set(quotedFeeAmount)
	require.Zero(t, senderDelta.Cmp(expectedSenderDelta))

	senderRegistryBalanceAfter, err := testhelpers.GetHoldingsBalance(ctx, senderParticipant, &registryInstrumentId, testhelpers.WithHoldingOwner(partySender))
	require.NoError(t, err)
	registryDelta := new(big.Rat).Sub(senderRegistryBalanceBefore, senderRegistryBalanceAfter)
	tokenTransferAmountRat, ok := new(big.Rat).SetString(tokenTransferAmountDecimal)
	require.True(t, ok)
	require.Positive(t, registryDelta.Sign())
	require.Zero(t, registryDelta.Cmp(tokenTransferAmountRat))

	t.Logf("Registry send completed: messageId=%s", returnedMessageId)
}
