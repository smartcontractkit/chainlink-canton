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

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/clientapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/events"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
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
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"

	// Import to register adapters
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"
)

// TestCCIPSendE2E tests the full send flow without token transfers.
// Validates that CCIPMessageSent is created with the correct message.
func TestCCIPSend(t *testing.T) {
	t.Parallel()

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

	dars := [][]byte{runtimeDar, coreDar, committeeVerifierDar, ccipSenderDar, ccipExecutorDar}
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
					USDPerNative: big.NewInt(100_000_000), // $1 with a decimals
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
	feeQuoterDestChainConfig := evmFeeQuoterDestChainConfigForLane()
	feeQuoterDestChainConfig.V2Params.USDPerUnitGas = big.NewInt(38)

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
		})
		log.Info().Err(err).Msg("EDS terminated")
		if !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Msg("EDS server exited with error")
			t.Fail()
			return
		}
	}()

	time.Sleep(1 * time.Second)

	// Create EDS clients
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create CCIP API client")
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create CCV API client")
	executorAPIClient, err := oapiExecutor.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create Executor API client")

	// Apply FeeQuoter dest chain config (needed by OnRamp.FinalizeFeeFromRouter)
	// Create PerPartyRouter for sender
	// Create PerPartyRouter for receiver
	perPartyRouterFactoryDisclosure, err := edsTesthelpers.GetPerPartyRouterFactoryDisclosure(t.Context(), ccipAPIClient, partySender)
	require.NoError(t, err)

	// TODO use operation to deploy PerPartyRouter
	res, err := senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: contracts.IdentifierFromBinding(ccipruntime.PerPartyRouterFactory{}),
					ContractId: perPartyRouterFactoryDisclosure.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(ccipruntime.CreateRouter{
						PartyOwner: types.PARTY(partySender),
						InstanceId: "router-sender",
					}),
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

	// Get disclosures for CCIPSender.Send
	disclosedCCIPSender, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), senderParticipant, contracts.IdentifierFromBinding(sender.CCIPSender{}))
	require.NoError(t, err)
	disclosedRouter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), senderParticipant, contracts.IdentifierFromBinding(ccipruntime.PerPartyRouter{}))
	require.NoError(t, err)

	// Prepare EVM receiver address
	receiver := hexutil.MustDecode("0xcf8def9adfe3dd90b3dffe42c8eabbf7cd4ee6ca")
	receiverHex := hex.EncodeToString(receiver)

	require.NotEmpty(t, disclosedCCIPSender.ContractId, "CCIPSender disclosure missing/empty")
	require.NotEmpty(t, disclosedRouter.ContractId, "Router disclosure missing/empty")

	t.Logf("disclosedCCIPSender.ContractId=%q", disclosedCCIPSender.ContractId)
	t.Logf("disclosedRouter.ContractId=%q", disclosedRouter.ContractId)

	t.Logf("partySender=%q", partySender)

	// Mint Amulet so the sender can cover non-zero fees (amount in holdings Numeric units per env).
	feeTokenHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, "1000.0")
	require.NoError(t, err, "failed to mint Amulet tokens to sender")
	t.Logf("Minted Amulet tokens to sender, Holding CID: %s", feeTokenHoldingCid)

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
	transferFactoryContextValues := make(map[string]splice_api_token_metadata_v1.AnyValue)
	if choiceContextRecord := choiceContext.GetRecord(); choiceContextRecord != nil && len(choiceContextRecord.Fields) > 0 {
		valuesField := choiceContextRecord.Fields[0]
		if valuesField.GetLabel() == "values" && valuesField.GetValue().GetTextMap() != nil {
			for _, entry := range valuesField.GetValue().GetTextMap().GetEntries() {
				if v := entry.GetValue().GetVariant(); v != nil {
					transferFactoryContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVContractId: new(types.CONTRACT_ID(v.GetValue().GetContractId()))}
				} else if entry.GetValue().GetText() != "" {
					transferFactoryContextValues[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVText: new(types.TEXT(entry.GetValue().GetText()))}
				}
			}
		}
	}

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
		Payload:       testPayloadHex,
		Sender:        partySender,
		Receiver:      receiverHex,
		GasLimit:      100_000,
		TokenTransfer: nil,
	}
	ccipSendDisclosure, err := edsTesthelpers.GetCCIPSendDisclosure(t.Context(), ccipAPIClient, msg, nil, nil)
	require.NoError(t, err)
	ccvSendDisclosure, err := edsTesthelpers.GetCCVSendDisclosure(t.Context(), ccvAPIClient, msg, committeeVerifierAddress.InstanceAddress())
	require.NoError(t, err)
	executorSendDisclosure, err := edsTesthelpers.GetExecutorSendDisclosure(t.Context(), executorAPIClient, msg, executorAddress.InstanceAddress(), ccipSendDisclosure.CCVs)
	require.NoError(t, err)

	ccvRawAddr := chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(committeeVerifierAddress.String())}
	execMcmsAddr := chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(executorAddress.String())}
	sendArgs := sender.Send{
		Context:                  ccipSendDisclosure.ChoiceContext,
		RouterCid:                types.CONTRACT_ID(routerCid),
		DestinationChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		Message: clientapi.Canton2AnyMessage{
			Receiver: types.TEXT(receiverHex),
			Payload:  types.TEXT(testPayloadHex),
			FeeToken: nativeInstrumentId,
			ExtraArgs: clientapi.ExtraArgs{
				V3: &clientapi.GenericExtraArgsV3{
					GasLimit: 100_000,
					Ccvs: []clientapi.CCVExtraArg{
						{
							CcvAddress: ccvRawAddr,
							CcvArgs:    types.TEXT(""),
						},
					},
					Executor: clientapi.ExecutorExtraArg{
						ExecutorWithAddress: &clientapi.ExecutorWithAddress{
							ExecutorAddress: execMcmsAddr,
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
				Context:    splice_api_token_metadata_v1.ChoiceContext{},
			},
		},
		ExecutorInput: &sender.ExecutorInput{
			ExecutorCid: types.CONTRACT_ID(executorSendDisclosure.ContractId),
			Context:     splice_api_token_metadata_v1.ChoiceContext{},
		},
	}

	// TODO: not clear what is duplicated here yet
	allDisclosures := slices.Concat(
		transferFactoryDisclosures, // TODO should be returned by EDS?
		ccipSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorSendDisclosure.DisclosedContracts,
	)
	quotedFee := quoteCCIPSenderFee(t, senderParticipant, partySender, ccipSenderCid, sendArgs, allDisclosures)
	require.NotEqual(t, "0.0", string(quotedFee.FeeTokenAmount), "GetFee should return a positive fee")

	// quoteCCIPSenderFee uses TRANSACTION_SHAPE_LEDGER_EFFECTS; refresh EDS disclosures so the
	// follow-up Send exercises current contract witnesses (avoids LOCAL_VERDICT_INACTIVE_CONTRACTS).
	ccipSendDisclosure, err = edsTesthelpers.GetCCIPSendDisclosure(t.Context(), ccipAPIClient, msg, nil, nil)
	require.NoError(t, err)
	ccvSendDisclosure, err = edsTesthelpers.GetCCVSendDisclosure(t.Context(), ccvAPIClient, msg, committeeVerifierAddress.InstanceAddress())
	require.NoError(t, err)
	executorSendDisclosure, err = edsTesthelpers.GetExecutorSendDisclosure(t.Context(), executorAPIClient, msg, executorAddress.InstanceAddress(), ccipSendDisclosure.CCVs)
	require.NoError(t, err)
	sendArgs.Context = ccipSendDisclosure.ChoiceContext
	sendArgs.FeeTokenInput.FeeTokenConfigCid = types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid)
	sendArgs.CcvSendInputs[0].CcvCid = types.CONTRACT_ID(ccvSendDisclosure.ContractId)
	sendArgs.ExecutorInput.ExecutorCid = types.CONTRACT_ID(executorSendDisclosure.ContractId)
	allDisclosures = testhelpers.DeduplicateDisclosedContracts(slices.Concat(
		transferFactoryDisclosures,
		ccipSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorSendDisclosure.DisclosedContracts,
	)...)

	senderBalanceBefore, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	ccipOwnerBalanceBefore, err := testhelpers.GetHoldingsBalance(t.Context(), ccipParticipant, &nativeInstrumentId)
	require.NoError(t, err)

	// CCIPSender.Send: PrepareSend + CCV tickets + Send in one transaction
	res, err = senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     contracts.IdentifierFromBinding(sender.CCIPSender{}),
					ContractId:     ccipSenderCid,
					Choice:         "Send",
					ChoiceArgument: ledger.MapToValue(sendArgs),
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: allDisclosures,
		},
	})
	require.NoError(t, err)

	// Log all events to help debug
	t.Logf("Transaction completed, checking events...")
	for i, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			t.Logf("Event %d: Created %s (CID: %s)", i, e.Created.GetTemplateId().GetEntityName(), e.Created.GetContractId())
		} else if e, ok := event.GetEvent().(*apiv2.Event_Archived); ok {
			t.Logf("Event %d: Archived %s (CID: %s)", i, e.Archived.GetTemplateId().GetEntityName(), e.Archived.GetContractId())
		}
	}

	// Extract CIPMessageSent event to verify success
	eventTemplateId := events.CCIPMessageSent{}.GetTemplateID()
	var ccipMessageSent *events.CCIPMessageSent
	for _, event := range res.GetTransaction().GetEvents() {
		if createdEvent := event.GetCreated(); createdEvent != nil {
			if templateId := createdEvent.GetTemplateId(); templateId != nil {
				gotTemplateId := fmt.Sprintf("#%s:%s:%s", createdEvent.GetPackageName(), templateId.GetModuleName(), templateId.GetEntityName())
				if gotTemplateId == eventTemplateId {
					// Found CCIPMessageSent event
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

	// Verify that the event contains the feeToken InstrumentId
	require.Equal(t, nativeInstrumentId, ccipMessageSent.Event.FeeToken, "CCIPMessageSent should contain feeToken InstrumentId")

	senderBalanceAfter, err := testhelpers.GetHoldingsBalance(t.Context(), senderParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	ccipOwnerBalanceAfter, err := testhelpers.GetHoldingsBalance(t.Context(), ccipParticipant, &nativeInstrumentId)
	require.NoError(t, err)
	senderDelta := new(big.Rat).Sub(senderBalanceBefore, senderBalanceAfter)
	ccipOwnerDelta := new(big.Rat).Sub(ccipOwnerBalanceAfter, ccipOwnerBalanceBefore)

	quotedFeeAmount, ok := new(big.Rat).SetString(strings.TrimSuffix(string(quotedFee.FeeTokenAmount), "."))
	require.True(t, ok, "quoted fee should parse as a decimal value")
	require.Zero(t, senderDelta.Cmp(quotedFeeAmount), "sender deduction should match GetFee quote")
	// In this test the executor is owned by ccipOwner (partyCCIP), so the executor fee
	// and any owner residual both land on ccipOwner — the total still equals the quoted fee.
	require.Zero(t, ccipOwnerDelta.Cmp(quotedFeeAmount), "ccipOwner should receive the full quoted fee (executor is also owned by ccipOwner in this test)")

	t.Logf("Send completed")
	t.Logf("  Message ID: %s", ccipMessageSent.Event.MessageId)
	t.Logf("  Original payload: %s", string(testPayload))
	t.Logf("  Encoded message: %s", ccipMessageSent.Event.EncodedMessage)

	t.Logf("✅ Success")
}
