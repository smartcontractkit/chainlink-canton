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
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
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
	rkccip "github.com/smartcontractkit/chainlink-canton/registry-kit/ccip"
	rkledger "github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"

	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"

	_ "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
)

const (
	registryExecuteInstrumentID      = "CCIP-REGISTRY-EXECUTE"
	registryExecutePoolInstanceID    = "ccip-registry-execute-pool"
	registryExecuteDefaultRLInstance = "ccip-registry-execute-rl-in-default"
	registryExecuteCustomRLInstance  = "ccip-registry-execute-rl-in-custom"
)

// TestRegistryTokenPool_FullReceiveFlow exercises CCIPReceiver.Execute with a Canton Registry Holding
// minted through a registrar-owned BurnMintTokenPool (hybrid EDS: CCIP/CCV via EDS, Registry pool
// context assembled manually).
//
//nolint:paralleltest // Holding balance assertions require exclusive CTF env
func TestRegistryTokenPool_FullReceiveFlow(t *testing.T) {
	runRegistryTokenPoolReceiveFlowTest(t, bnmTokenPoolReceiveFlowTestCase{
		tokenAmount:                   big.NewInt(50e10),
		expectedTransferAmount:        50,
		defaultInboundLimiterCapacity: "1000000000000",
		customInboundLimiterCapacity:  "10000000000000",
		expectedCustomLimiterTokens:   "9500000000000.",
	})
}

func runRegistryTokenPoolReceiveFlowTest(t *testing.T, tc bnmTokenPoolReceiveFlowTestCase) {
	t.Helper()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(3))

	ccipParticipant := env.Chain.Participants[0]
	receiverParticipant := env.Chain.Participants[1]
	tokenPoolOwnerParticipant := env.Chain.Participants[0]

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	uploadRegistryDARs(t, ccipParticipant, receiverParticipant, tokenPoolOwnerParticipant)

	runtimeDar, err := contracts.GetDar(contracts.CCIPRuntimeV2, contracts.DevVersion)
	require.NoError(t, err)
	coreDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifierV2, contracts.DevVersion)
	require.NoError(t, err)
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPBurnMintTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	ccipReceiverDar, err := contracts.GetDar(contracts.CCIPReceiverV2, contracts.DevVersion)
	require.NoError(t, err)
	executorDar, err := contracts.GetDar(contracts.CCIPExecutorV2, contracts.DevVersion)
	require.NoError(t, err)

	dars := [][]byte{runtimeDar, coreDar, committeeVerifierDar, tokenPoolDar, ccipReceiverDar, executorDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant, tokenPoolOwnerParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded CCIP DARs to all participants: %v", packageIds)

	ctx := t.Context()
	partyCCIP := ccipParticipant.PartyID
	partyReceiver := receiverParticipant.PartyID
	partyRegistrar := testhelpers.AllocateParty(t, ccipParticipant, "registry-execute-registrar")
	testhelpers.GrantCanActAs(t, ccipParticipant, partyRegistrar)
	t.Logf("Parties: CCIP=%s, Receiver=%s, Registrar=%s", partyCCIP, partyReceiver, partyRegistrar)

	ccipClient := rkledger.NewCTFClient(ccipParticipant)
	registrarClient := rkledger.NewCTFClient(ccipParticipant)

	_, tokenMetadataClient, _, err := testhelpers.NewValidatorAPIClients(ccipParticipant)
	require.NoError(t, err)

	registryAdmin, err := testhelpers.GetRegistryAdmin(ctx, tokenMetadataClient)
	require.NoError(t, err)

	nativeInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}
	registryInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(partyRegistrar),
		Id:    types.TEXT(registryExecuteInstrumentID),
	}
	hashedRegistryInstrumentId := contracts.EncodeInstrumentID(registryInstrumentId)

	sourceChainSelector := fmt.Sprintf("%d", chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector)
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	ccvSignerKeys := make([]*ecdsa.PrivateKey, 0, 3)
	ccvSignerPubKeys := make([]string, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}

	versionTag := "e9a05a20"
	ccvQualifier := devenvcommon.DefaultCommitteeVerifierQualifier

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), reporter)
	cldfEnv := cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{env.Chain}),
		OperationsBundle: bundle,
	}

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
							StorageLocations:             []types.TEXT{"ipfs://test-registry-execute"},
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
								AllowedFinalityConfig: ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
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

	cantonAdapter, ok := lanes.GetLaneAdapterRegistry().GetLaneAdapter(chainsel.FamilyCanton, semver.MustParse("2.0.0"))
	require.True(t, ok)
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
	require.NoError(t, err)
	runningDs := datastore.NewMemoryDataStore()
	for _, address := range deployLaneLegReport.Output.Addresses {
		err = runningDs.Addresses().Add(address)
		require.NoError(t, err)
	}
	err = runningDs.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = runningDs.Seal()

	bootstrap, err := registry.BootstrapServices(ctx, registrarClient, partyRegistrar, registryExecuteInstrumentID)
	require.NoError(t, err)

	now := time.Now()
	inboundRateLimiterAddr, err := rkccip.DeployInboundRateLimiterForOwner(ctx, registrarClient, partyRegistrar, ratelimiter.RateLimiter{
		InstanceId:          types.TEXT(registryExecuteDefaultRLInstance),
		PoolInstanceId:      types.TEXT(registryExecutePoolInstanceID),
		PoolOwner:           types.PARTY(partyRegistrar),
		RemoteChainSelector: types.NUMERIC(sourceChainSelector),
		Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Inbound,
		Mode:                ratelimiter.RateLimitModeRateLimitMode_DefaultFinality,
		IsEnabled:           true,
		Capacity:            types.NUMERIC(tc.defaultInboundLimiterCapacity),
		Rate:                types.NUMERIC(tc.defaultInboundLimiterCapacity),
		Tokens:              types.NUMERIC(tc.defaultInboundLimiterCapacity),
		LastUpdated:         types.TIMESTAMP(now),
	})
	require.NoError(t, err)

	inboundCustomRateLimiterAddr, err := rkccip.DeployInboundRateLimiterForOwner(ctx, registrarClient, partyRegistrar, ratelimiter.RateLimiter{
		InstanceId:          types.TEXT(registryExecuteCustomRLInstance),
		PoolInstanceId:      types.TEXT(registryExecutePoolInstanceID),
		PoolOwner:           types.PARTY(partyRegistrar),
		RemoteChainSelector: types.NUMERIC(sourceChainSelector),
		Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Inbound,
		Mode:                ratelimiter.RateLimitModeRateLimitMode_CustomFinality,
		IsEnabled:           true,
		Capacity:            types.NUMERIC(tc.customInboundLimiterCapacity),
		Rate:                types.NUMERIC(tc.customInboundLimiterCapacity),
		Tokens:              types.NUMERIC(tc.customInboundLimiterCapacity),
		LastUpdated:         types.TIMESTAMP(now),
	})
	require.NoError(t, err)

	outboundRateLimiterAddr, err := rkccip.DeployOutboundRateLimiterForOwner(ctx, registrarClient, partyRegistrar, ratelimiter.RateLimiter{
		InstanceId:          types.TEXT("ccip-registry-execute-rl-out"),
		PoolInstanceId:      types.TEXT(registryExecutePoolInstanceID),
		PoolOwner:           types.PARTY(partyRegistrar),
		RemoteChainSelector: types.NUMERIC(sourceChainSelector),
		Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Outbound,
		Mode:                ratelimiter.RateLimitModeRateLimitMode_DefaultFinality,
		IsEnabled:           false,
		Capacity:            types.NUMERIC("0"),
		Rate:                types.NUMERIC("0"),
		Tokens:              types.NUMERIC("0"),
		LastUpdated:         types.TIMESTAMP(now),
	})
	require.NoError(t, err)

	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")

	poolAddr, err := rkccip.DeployBurnMintPoolForOwner(ctx, registrarClient, rkccip.PoolDeployDeps{
		CcipOwner:          partyCCIP,
		TokenAdminRegistry: tokenAdminRegistryAddress,
		RMNRemote:          rmnRemoteAddress,
		FeeQuoter:          feeQuoterAddress,
	}, partyRegistrar, registryInstrumentId, registryExecutePoolInstanceID,
		map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
			types.NUMERIC(sourceChainSelector): {
				RemotePools:        []types.TEXT{types.TEXT(hex.EncodeToString(remotePoolAddress))},
				RemoteTokenAddress: types.TEXT(hex.EncodeToString(remoteTokenAddress)),
				InboundCCVs:        []chainlinkapi.RawInstanceAddress{},
				OutboundCCVs:       []chainlinkapi.RawInstanceAddress{},
				FinalityConfig: ccipcodec.FinalityConfig{
					BlockDepth: new(types.INT64(2000)),
				},
				InboundRateLimiter:                         inboundRateLimiterAddr.Binding(),
				InboundCustomBlockConfirmationsRateLimiter: inboundCustomRateLimiterAddr.Binding(),
				OutboundRateLimiter:                        outboundRateLimiterAddr.Binding(),
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
		PoolInstanceID:        registryExecutePoolInstanceID,
		CcipParty:             partyCCIP,
		PoolOwnerParty:        partyRegistrar,
		PoolOwnerClient:       registrarClient,
	})
	require.NoError(t, err)

	tokenConfigCID, err = rkccip.SetBurnMintFactory(ctx, registrarClient, rkccip.SetBurnMintFactoryInput{
		TokenAdminRegistryCID: tarCID,
		TokenConfigCID:        tokenConfigCID,
		InstrumentId:          registryInstrumentId,
		BurnMintFactoryCID:    bootstrap.AllocationFactory,
		CcipParty:             partyCCIP,
		PoolOwnerParty:        partyRegistrar,
		CcipClient:            ccipClient,
		PoolOwnerClient:       registrarClient,
	})
	require.NoError(t, err)
	_ = tokenConfigCID

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
		})
		log.Info().Err(err).Msg("EDS terminated")
		if !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Msg("EDS server exited with error")
			t.Fail()
		}
	}()

	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err)
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	perPartyRouterFactoryDisclosure, err := edsTesthelpers.GetPerPartyRouterFactoryDisclosure(ctx, ccipAPIClient, partyReceiver)
	require.NoError(t, err)

	res, err := receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#" + ccipruntime.PackageName, ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: perPartyRouterFactoryDisclosure.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-registry-execute"}}},
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

	tokenTransfer := buildTokenTransferV1(tc.tokenAmount, remotePoolAddress, remoteTokenAddress, hashedRegistryInstrumentId, partyReceiver, tc.sourcePoolData)

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

	verifierResults, err := GenerateVerifierResults(encodedMessage, ccvSignerKeys[:2])
	require.NoError(t, err)
	verifierResultsHex := hex.EncodeToString(verifierResults)

	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccipreceiver-registry"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "receiverFinalityConfig", Value: finalityConfigValueFromBlockConfirmations()},
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

	receiverBalanceRatBefore, err := testhelpers.GetHoldingsBalance(ctx, receiverParticipant, &registryInstrumentId, testhelpers.WithHoldingOwner(partyReceiver))
	require.NoError(t, err)
	receiverBalanceBefore, _ := new(big.Float).SetRat(receiverBalanceRatBefore).Float64()

	ccipExecuteDisclosure, err := edsTesthelpers.GetCCIPExecuteDisclosure(ctx, ccipAPIClient, encodedMessageHex)
	require.NoError(t, err)
	ccvExecuteDisclosure, err := edsTesthelpers.GetCCVExecuteDisclosure(ctx, ccvAPIClient, encodedMessageHex, committeeVerifierAddress.InstanceAddress())
	require.NoError(t, err)

	poolExecuteDeps := registryPoolExecuteDeps{
		Client:                       registrarClient,
		CcipClient:                   ccipClient,
		RegistrarParty:               partyRegistrar,
		CcipParty:                    partyCCIP,
		Bootstrap:                    bootstrap,
		PoolInstanceID:               registryExecutePoolInstanceID,
		DefaultRateLimiterInstanceID: registryExecuteDefaultRLInstance,
		CustomRateLimiterInstanceID:  registryExecuteCustomRLInstance,
		PoolAddress:                  poolAddr,
		TokenAdminRegistryCID:        tarCID,
		RMNRemoteAddress:             rmnRemoteAddress,
	}
	tokenPoolDisclosure := buildRegistryTokenPoolExecuteDisclosure(t, ctx, ccipParticipant, ccipParticipant, ccipAPIClient, poolExecuteDeps, hashedRegistryInstrumentId, true)

	executeArgs := receiver.Execute{
		Context:        ccipExecuteDisclosure.ChoiceContext,
		RouterCid:      types.CONTRACT_ID(routerCid),
		EncodedMessage: types.TEXT(encodedMessageHex),
		TokenTransfer: &receiver.TokenTransferInput{
			TokenPoolCid:       types.CONTRACT_ID(tokenPoolDisclosure.ContractId),
			TokenReceiverParty: types.PARTY(partyReceiver),
			Context:            tokenPoolDisclosure.ChoiceContext,
		},
		CcvInputs: []receiver.CCVInput{
			{
				CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
				VerifierResults: types.TEXT(verifierResultsHex),
				Context:         ccvExecuteDisclosure.ChoiceContext,
			},
		},
	}

	_, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
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

	receiverBalanceRatAfter, err := testhelpers.GetHoldingsBalance(ctx, receiverParticipant, &registryInstrumentId, testhelpers.WithHoldingOwner(partyReceiver))
	require.NoError(t, err)
	receiverBalanceAfter, _ := new(big.Float).SetRat(receiverBalanceRatAfter).Float64()

	actualTransferAmount := receiverBalanceAfter - receiverBalanceBefore
	require.InDelta(t, tc.expectedTransferAmount, actualTransferAmount, 0.01, "Receiver balance should increase by transfer amount")
	t.Logf("Receiver balance: %.2f -> %.2f (transferred %.2f)", receiverBalanceBefore, receiverBalanceAfter, actualTransferAmount)

	registrarQueryParticipant := ccipParticipant
	registrarQueryParticipant.PartyID = partyRegistrar

	if tc.expectedDefaultLimiterTokens != "" {
		defaultRateLimiter, err := findActiveRateLimiterByInstanceID(ctx, registrarQueryParticipant, registryExecuteDefaultRLInstance)
		require.NoError(t, err)
		require.Equal(t, tc.expectedDefaultLimiterTokens, getRateLimiterTokens(defaultRateLimiter))
	}
	if tc.expectedCustomLimiterTokens != "" {
		customRateLimiter, err := findActiveRateLimiterByInstanceID(ctx, registrarQueryParticipant, registryExecuteCustomRLInstance)
		require.NoError(t, err)
		require.Equal(t, tc.expectedCustomLimiterTokens, getRateLimiterTokens(customRateLimiter))
	}

	t.Log("Success")
}
