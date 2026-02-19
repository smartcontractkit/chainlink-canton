package tests

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

// TestCCIPSendE2E tests the full send flow without token transfers.
// Validates that CCIPMessageSent is created with the correct message.
func TestCCIPSendE2E(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	ccipParticipant := env.Participant(1)
	senderParticipant := env.Participant(2)

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

	dars := [][]byte{commonDar, offRampDar, onRampDar, feeQuoterDar, tokenAdminRegistryDar, committeeVerifierDar, perPartyRouterDar, rmnDar, ccipSenderDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, senderParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.Party
	partySender := senderParticipant.Party
	t.Logf("Parties: CCIP=%s, Sender=%s", partyCCIP, partySender)

	// CCV Setup
	ccvSignerKeys := make([]*ecdsa.PrivateKey, 0, 3)
	ccvSignerPubKeys := make([]types.TEXT, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, types.TEXT(pubKeyHex))
	}
	t.Logf("Generated %d CCV signer keys", len(ccvSignerKeys))

	versionTag := "49ff34ed"
	ccvQualifier := "default"
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

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
		ChainSelector: env.Selector,
		Participant:   0,
		Party:         partyCCIP,
		Config: changesets.DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: partyCCIP,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Qualifier: ccvQualifier,
						Template: ccvs.CommitteeVerifier{
							Owner:                    types.PARTY(partyCCIP),
							CcipOwner:                types.PARTY(partyCCIP),
							VersionTag:               types.TEXT(versionTag),
							MessageSentObserver:      types.PARTY(partyCCIP),
							StorageLocation:          "ipfs://test-send",
							Threshold:                2,
							Signers:                  ccvSignerPubKeys,
							RmnRemoteInstanceAddress: common.RawInstanceAddress{}, // Set by sequence
							// MUST be a real GENMAP, not a Go map.
							RemoteChainFeeConfigs: types.GENMAP{
								strconv.FormatUint(remoteSelector, 10): ccvs.CCVFeeConfig{
									FeeUSDCents:        types.NUMERIC("0"),
									GasForVerification: types.INT64(0),
									PayloadSizeBytes:   types.INT64(0),
								}.ToMap(),
							},
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10)),
						OnRampAddress: "", // populated by the sequence
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
				},
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
	globalConfig, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	require.NoError(t, err, "failed to get GlobalConfig address")
	feeQuoter, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, ""))
	require.NoError(t, err, "failed to get FeeQuoter address")
	onRamp, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(onramp.ContractType), onramp.Version, ""))
	require.NoError(t, err, "failed to get OnRamp address")
	offRamp, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(offramp.ContractType), offramp.Version, ""))
	require.NoError(t, err, "failed to get OffRamp address")
	committeeVerifier, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Selector, datastore.ContractType(committee_verifier.ContractType), committee_verifier.Version, ccvQualifier))
	require.NoError(t, err, "failed to get CommitteeVerifier address")

	// Deploy and configure lane for outbound sends
	committeeVerifierRawAddr, err := contracts.RawInstanceAddressFromString(committeeVerifier.Labels.List()[0])
	require.NoError(t, err, "failed to parse CommitteeVerifier raw address")
	out, err = changesets.ConfigureChainForLanes{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.ConfigureChainForLanesConfig]{
		ChainSelector: env.Selector,
		Participant:   0,
		Party:         partyCCIP,
		Config: changesets.ConfigureChainForLanesConfig{
			Input: sequences.ConfigureChainForLanesInput{
				ChainSelector:      env.Selector,
				GlobalConfig:       contracts.HexToInstanceAddress(globalConfig.Address),
				FeeQuoter:          contracts.HexToInstanceAddress(feeQuoter.Address),
				OnRamp:             contracts.HexToInstanceAddress(onRamp.Address),
				OffRamp:            contracts.HexToInstanceAddress(offRamp.Address),
				CommitteeVerifiers: nil,
				RemoteChains: map[uint64]adapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]{
					remoteSelector: {
						AllowTrafficFrom:         true,
						OnRamps:                  [][]byte{[]byte("0000000000000000000000000000000000000001")}, // remote chain onRamp
						OffRamp:                  nil,
						DefaultInboundCCVs:       nil,
						LaneMandatedInboundCCVs:  nil,
						DefaultOutboundCCVs:      []contracts.RawInstanceAddress{committeeVerifierRawAddr},
						LaneMandatedOutboundCCVs: nil,
						DefaultExecutor:          contracts.RawInstanceAddress(committeeVerifierRawAddr.String()), // random executor
						FeeQuoterDestChainConfig: adapters.FeeQuoterDestChainConfig{
							NetworkFeeUSDCents:      0,
							DefaultTokenFeeUSDCents: 0,
						},
						ExecutorDestChainConfig: adapters.ExecutorDestChainConfig{},
						AddressBytesLength:      0,
						BaseExecutionGasCost:    0,
					},
				},
			},
		},
	})
	require.NoErrorf(t, err, "Failed to configure chain for lanes")
	err = out.DataStore.Merge(cldfEnv.DataStore)
	require.NoError(t, err)
	cldfEnv.DataStore = out.DataStore.Seal()
	t.Log("Configured chain for lanes")

	// Setup Amulet token as fee token
	// Get registry admin for Amulet tokens
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), env.Splice.TokenMetadataClient)
	require.NoError(t, err, "failed to get registry admin")

	// Fee token is Amulet
	feeTokenInstrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}

	// Get disclosed FeeQuoter contract
	disclosedFeeQuoterForConfig, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter",
	})
	require.NoError(t, err)

	// Configure FeeToken: Add FeeToken to FeeQuoter

	// ApplyFeeTokenUpdates: Add the fee token with premium multiplier of 1.0 (no premium)
	premiumMultiplier := "1.0" // 1.0 means no premium/discount
	_, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter"},
					ContractId: disclosedFeeQuoterForConfig.ContractId,
					Choice:     "ApplyFeeTokenUpdates",
					ChoiceArgument: ledger.MapToValue(feequoter.ApplyFeeTokenUpdates{
						FeeTokensToRemove: []splice_api_token_holding_v1.InstrumentId{},
						FeeTokensToAdd: []feequoter.FeeTokenArgs{
							{
								InstrumentId:      feeTokenInstrumentId,
								PremiumMultiplier: types.NUMERIC(premiumMultiplier),
							},
						},
						Caller: types.PARTY(partyCCIP),
					}),
				}},
			}},
			ActAs:              []string{partyCCIP},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedFeeQuoterForConfig},
		},
	})
	require.NoError(t, err, "failed to apply fee token updates")
	t.Logf("Applied FeeToken updates to FeeQuoter")

	// Get the updated FeeQuoter cID
	disclosedFeeQuoterForConfig, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter",
	})
	require.NoError(t, err)
	t.Logf("Updated FeeQuoter cID: %s", disclosedFeeQuoterForConfig.ContractId)

	// UpdatePrices: Set price for FeeToken (e.g., $1.00 per token)
	usdPerToken := "1.00"
	_, err = ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter"},
					ContractId: disclosedFeeQuoterForConfig.ContractId,
					Choice:     "UpdatePrices",
					ChoiceArgument: ledger.MapToValue(feequoter.UpdatePrices{
						PriceUpdates: feequoter.PriceUpdates{
							TokenPriceUpdates: []feequoter.TokenPriceUpdate{
								{
									InstrumentId: feeTokenInstrumentId,
									UsdPerToken:  types.NUMERIC(usdPerToken),
								},
							},
							GasPriceUpdates: []feequoter.GasPriceUpdate{}, // No gas price updates for this test
						},
						Caller: types.PARTY(partyCCIP),
					}),
				}},
			}},
			ActAs:              []string{partyCCIP},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedFeeQuoterForConfig},
		},
	})
	require.NoError(t, err, "failed to update prices")
	t.Logf("Updated FeeToken price to $%s per token", usdPerToken)

	// Create PerPartyRouter for sender
	disclosedFactory, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	})
	require.NoError(t, err)

	res, err := senderParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: disclosedFactory.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.CreateRouter{
						PartyOwner: types.PARTY(partySender),
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
	res, err = senderParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
	receiverPartyID := "receiver-party"
	receiverBytes := EncodePartyID(receiverPartyID)
	receiverHex := hex.EncodeToString(receiverBytes)

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

	// Mint Amulet tokens to sender so they can pay the fee
	feeTokenHoldingCid, err := testhelpers.MintAMT(t.Context(), senderParticipant, env.Splice.TokenMetadataClient, env.Splice.TransferInstructionClient, senderParticipant.ScanProxyClient, partySender, "100.00")
	require.NoError(t, err, "failed to mint Amulet tokens to sender")
	t.Logf("Minted 100 Amulet tokens to sender, Holding CID: %s", feeTokenHoldingCid)

	// Get disclosed contract for the fee token holding
	disclosedFeeTokenHolding, err := testhelpers.GetDisclosedContractById(t.Context(), senderParticipant, feeTokenHoldingCid)
	require.NoError(t, err, "failed to get disclosed contract for fee token holding")

	// Get disclosed CommitteeVerifier contract for CCV send inputs
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err, "failed to get disclosed CommitteeVerifier")

	// Get transfer factory for Amulet tokens (sender to CCIP owner)
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := testhelpers.GetTransferFactory(t.Context(), env.Splice.TransferInstructionClient, registryAdmin, partySender, partyCCIP)
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
			var keys []string
			for _, entry := range textMap.GetEntries() {
				keys = append(keys, entry.GetKey())
			}
			return keys
		}())
	}

	// Log choiceContext structure for debugging
	choiceContextJSON, _ := json.MarshalIndent(choiceContext, "", "  ")
	t.Logf("choiceContext structure:\n%s", string(choiceContextJSON))

	// Strip "0x" prefix from transferFactoryCid if present (Canton contract IDs shouldn't have it, but be safe)
	transferFactoryCid = strings.TrimPrefix(transferFactoryCid, "0x")
	transferFactoryCid = strings.TrimPrefix(transferFactoryCid, "0X")

	// Build command arguments manually using apiv2.Value structures (like ccip_execute_token_test.go)
	// Use choiceContext directly as apiv2.Value to preserve the structure
	var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "values",
			Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
		},
	}}}}

	// Build feeTokenHoldingCids list
	feeTokenHoldingCids := []*apiv2.Value{
		{Sum: &apiv2.Value_ContractId{ContractId: feeTokenHoldingCid}},
	}

	// Build FeeToken InstrumentId
	feeTokenInstrumentIdValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(feeTokenInstrumentId.Admin)}}},
		{Label: "id", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(feeTokenInstrumentId.Id)}}},
	}}}}

	ccipSendArgs := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCid}}},
		{Label: "onRampCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedOnRamp.ContractId}}},
		{Label: "globalConfigCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedGlobalConfig.ContractId}}},
		{Label: "tokenAdminRegistryCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedTar.ContractId}}},
		{Label: "rmnRemoteCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedRmnRemote.ContractId}}},
		{Label: "feeQuoterCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedFeeQuoter.ContractId}}},
		{Label: "destChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: strconv.FormatUint(remoteSelector, 10)}}},
		{Label: "receiver", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: receiverHex}}},
		{Label: "payload", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: testPayloadHex}}},
		{Label: "ccipReceiveGasLimit", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 100000}}},
		{Label: "senderRequiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
		{Label: "feeToken", Value: feeTokenInstrumentIdValue},
		{Label: "feeTokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
			{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}}},
			{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "context", Value: choiceContext},
				{Label: "meta", Value: emptyMetadata},
			}}}}},
			{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
		}}}}},
		{Label: "feeTokenHoldingCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: feeTokenHoldingCids}}}},
		{Label: "tokenTransfer", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{}}}},
		{Label: "ccvSendInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
			{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "ccvCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedCCV.ContractId}}},
				{Label: "ccvRawAddress", Value: rawInstanceAddress(committeeVerifierRawAddr.String())},
				{Label: "verifierArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			}}}},
		}}}}},
	}}}}

	b, _ := json.MarshalIndent(ccipSendArgs, "", "  ")
	require.NoError(t, err)
	t.Logf("Send arg:\n%s", string(b))

	// CCIPSender.Send: PrepareSend + CCV tickets + Send in one transaction
	res, err = senderParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
			DisclosedContracts: append(
				append(
					[]*apiv2.DisclosedContract{disclosedCCIPSender, disclosedRouter, disclosedOnRamp, disclosedGlobalConfig, disclosedTar, disclosedRmnRemote, disclosedFeeQuoter, disclosedFeeTokenHolding, disclosedCCV},
					transferFactoryDisclosures...,
				),
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
				// CCIPMessageSent structure:
				// - ccipOwner (Party)
				// - sender (Party)
				// - observers ([Party])
				// - event (CCIPMessageSentEvent)
				//   - destChainSelector
				//   - sequenceNumber
				//   - messageId
				//   - encodedMessage
				//   - verifierBlobs
				//   - receipts
				fields := e.Created.GetCreateArguments().GetFields()
				if len(fields) >= 4 {
					// fields[3] is the "event" field (CCIPMessageSentEvent)
					eventField := fields[3].GetValue().GetRecord()
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
