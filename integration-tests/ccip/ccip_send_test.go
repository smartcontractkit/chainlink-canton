package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"testing"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	mcmsbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"

	// Import to register adapters
	_ "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/adapters"

	_ "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
)

// TestCCIPSendE2E tests the full send flow without token transfers.
// Validates that CCIPMessageSent is created with the correct message.
func TestCCIPSend(t *testing.T) {
	t.Parallel()

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
	ccipTestDar, err := contracts.GetDar(contracts.CCIPTest, contracts.CurrentVersion)
	require.NoError(t, err)

	dars := [][]byte{commonDar, offRampDar, onRampDar, feeQuoterDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar, rmnDar, ccipSenderDar, ccipExecutorDar, ccipTestDar}
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
	ccvQualifier := devenvcommon.DefaultCommitteeVerifierQualifier
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
	executorAddress, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(executor.ContractType), executor.Version, devenvcommon.DefaultExecutorQualifier))
	require.NoError(t, err, "failed to get Executor address")

	// Deploy and configure lane for outbound sends
	cantonAdapter, ok := lanes.GetLaneAdapterRegistry().GetLaneAdapter(chainsel.FamilyCanton, semver.MustParse("2.0.0"))
	require.Truef(t, ok, "failed to get Canton Lane adapter")
	feeQuoterDestChainConfig := lanes.DefaultFeeQuoterDestChainConfig(true, remoteSelector)
	feeQuoterDestChainConfig.V2Params.USDPerUnitGas = big.NewInt(38)

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

	// Create PerPartyRouter for sender
	disclosedFactory, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	})
	require.NoError(t, err)

	// TODO use operation to deploy PerPartyRouter
	res, err := senderParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
	disclosedCCIPSender, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), senderParticipant, &apiv2.Identifier{
		PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender",
	})
	require.NoError(t, err)
	disclosedRouter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), senderParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	})
	require.NoError(t, err)
	disclosedOnRamp, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-onramp", ModuleName: "CCIP.OnRamp", EntityName: "OnRamp",
	})
	require.NoError(t, err)
	disclosedGlobalConfig, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig",
	})
	require.NoError(t, err)
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)
	disclosedRmnRemote, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote",
	})
	require.NoError(t, err)
	disclosedFeeQuoter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter",
	})
	require.NoError(t, err)

	// Prepare receiver address (destination party encoded as keccak256)
	receiver := hexutil.MustDecode("0xcf8def9adfe3dd90b3dffe42c8eabbf7cd4ee6ca")
	receiverHex := hex.EncodeToString(receiver)

	require.NotEmpty(t, disclosedOnRamp.ContractId, "OnRamp disclosure missing/empty")
	require.NotEmpty(t, disclosedGlobalConfig.ContractId, "GlobalConfig disclosure missing/empty")
	require.NotEmpty(t, disclosedTar.ContractId, "TAR disclosure missing/empty")
	require.NotEmpty(t, disclosedRmnRemote.ContractId, "RMNRemote disclosure missing/empty")
	require.NotEmpty(t, disclosedFeeQuoter.ContractId, "FeeQuoter disclosure missing/empty")

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
	feeTokenHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partySender, "10000000000")
	require.NoError(t, err, "failed to mint Amulet tokens to sender")
	t.Logf("Minted 100 whole Amulet tokens to sender, Holding CID: %s", feeTokenHoldingCid)

	// Get disclosed contract for the fee token holding
	disclosedFeeTokenHolding, err := testhelpers.GetDisclosedContractById(t.Context(), senderParticipant, feeTokenHoldingCid)
	require.NoError(t, err, "failed to get disclosed contract for fee token holding")

	// Get disclosed CommitteeVerifier contract for CCV send inputs
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err, "failed to get disclosed CommitteeVerifier")

	disclosedExecutor, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-executor", ModuleName: "CCIP.Executor", EntityName: "Executor",
	})
	require.NoError(t, err, "failed to get disclosed Executor")

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

	// Build the main Send context with CCIP contract IDs (matching execute test pattern)
	onRampCid := types.CONTRACT_ID(disclosedOnRamp.ContractId)
	globalConfigCid := types.CONTRACT_ID(disclosedGlobalConfig.ContractId)
	tarCid := types.CONTRACT_ID(disclosedTar.ContractId)
	feeQuoterCid := types.CONTRACT_ID(disclosedFeeQuoter.ContractId)
	rmnRemoteCid := types.CONTRACT_ID(disclosedRmnRemote.ContractId)
	executorCid := types.CONTRACT_ID(disclosedExecutor.ContractId)
	sendContext := common.CCIPContext{
		Values: types.TEXTMAP{
			"on-ramp":              common.AnyValue{AVContractId: &onRampCid},
			"global-config":        common.AnyValue{AVContractId: &globalConfigCid},
			"token-admin-registry": common.AnyValue{AVContractId: &tarCid},
			"fee-quoter":           common.AnyValue{AVContractId: &feeQuoterCid},
			"rmn-remote":           common.AnyValue{AVContractId: &rmnRemoteCid},
		},
	}

	sendArgs := ccipsender.Send{
		Context:                  sendContext,
		RouterCid:                types.CONTRACT_ID(routerCid),
		DestinationChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
		Message: ccipclient.Canton2AnyMessage{
			Receiver: types.TEXT(receiverHex),
			Payload:  types.TEXT(testPayloadHex),
			FeeToken: ccipclient.FeeTokenInput{
				Token:           nativeInstrumentId,
				TokenInput:      feeTokenInput,
				SenderInputCids: []types.CONTRACT_ID{types.CONTRACT_ID(feeTokenHoldingCid)},
			},
			ExtraArgs: ccipclient.ExtraArgs{
				V3: &ccipclient.GenericExtraArgsV3{
					GasLimit:           100_000,
					BlockConfirmations: 0,
					Ccvs: []mcmsbindings.RawInstanceAddress{
						{Unpack: types.TEXT(committeeVerifierRawAddr.String())},
					},
					Executor: &ccipclient.ExecutorInput{
						ExecutorCid:  executorCid,
						ExecutorArgs: types.TEXT(""),
					},
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

	// CCIPSender.Send: PrepareSend + CCV tickets + Send in one transaction
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
			ActAs: []string{partySender},
			DisclosedContracts: slices.Concat(
				[]*apiv2.DisclosedContract{disclosedCCIPSender, disclosedRouter, disclosedOnRamp, disclosedGlobalConfig, disclosedTar, disclosedRmnRemote, disclosedFeeQuoter, disclosedFeeTokenHolding, disclosedCCV, disclosedExecutor},
				transferFactoryDisclosures,
			),
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

	// Extract messageId from CCIPMessageSent to verify success
	var returnedMessageId string
	var returnedEncodedMessage string
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
				fields := e.Created.GetCreateArguments().GetFields()
				if len(fields) >= 5 {
					// fields[4] is the "event" field (CCIPMessageSentEvent)
					// (ccipOwner, ccvOwners, sender, observers, event)
					eventField := fields[4].GetValue().GetRecord()
					if eventField != nil && len(eventField.Fields) >= 4 {
						// eventField.Fields[2] is messageId, eventField.Fields[3] is encodedMessage
						returnedMessageId = eventField.Fields[2].GetValue().GetText()
						if len(eventField.Fields) >= 4 {
							returnedEncodedMessage = eventField.Fields[3].GetValue().GetText()
						}
					}
				}

				break
			}
		}
	}
	require.NotEmpty(t, returnedMessageId, "CCIPMessageSent should be created")
	require.NotEmpty(t, returnedEncodedMessage, "CCIPMessageSent should contain encoded message")

	t.Logf("Send completed")
	t.Logf("  Message ID: %s", returnedMessageId)
	t.Logf("  Original payload: %s", string(testPayload))
	t.Logf("  Encoded message: %s", returnedEncodedMessage)
}
