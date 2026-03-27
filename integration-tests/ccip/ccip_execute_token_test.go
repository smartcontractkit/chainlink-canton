package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
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
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
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

// encodeInstrumentId encodes an InstrumentId to bytes matching Daml encodeInstrumentId.
// Format: UTF-8 bytes of "id@admin" (matches Daml's toHex(id <> "@" <> partyToText admin)).
func encodeInstrumentId(admin, identifier string) []byte {
	return []byte(identifier + "@" + admin)
}

// TestLnRTokenPool_FullReceiveFlow tests the complete CCIP inbound token release flow.
// - Deploy all CCIP contracts including CommitteeVerifier, GlobalConfig, OffRamp, PerPartyRouter
// - Mint tokens to pool (simulating prior locked tokens)
// - Build inbound message with TokenTransfer
// - Generate signatures and call CommitteeVerifier_VerifyMessage to append CCV verification
// - Call PerPartyRouter.Execute to process message and get TokenReceiveTicket
// - Call TokenPool_ReleaseFromTicket to transfer tokens from pool to receiver
// - Verify receiver received the tokens
// - Validate FTF/custom-finality path by:
//   - requiring minBlockConfirmations=2000 on the destination pool
//   - enabling a default inbound limiter with lower capacity
//   - enabling a custom-finality inbound limiter with higher capacity
//     Success proves ReleaseFromTicket selected the custom inbound limiter.
func TestLnRTokenPool_FullReceiveFlow(t *testing.T) {
	t.Parallel()

	runLnRTokenPoolReceiveFlowTest(t, lnrTokenPoolReceiveFlowTestCase{
		tokenAmount:                   big.NewInt(5),
		expectedTransferAmount:        5,
		defaultInboundLimiterCapacity: "1000000",
		customInboundLimiterCapacity:  "10000000",
	})
}

// func TestLnRTokenPool_FullReceiveFlow_DecimalConversion(t *testing.T) {
// 	t.Parallel()

// 	runLnRTokenPoolReceiveFlowTest(t, lnrTokenPoolReceiveFlowTestCase{
// 		tokenAmount:                   big.NewInt(7_000_000_000_000),
// 		sourcePoolData:                encodeUint256Bytes(18),
// 		expectedTransferAmount:        7,
// 		defaultInboundLimiterCapacity: "5",
// 		customInboundLimiterCapacity:  "10",
// 		expectedDefaultLimiterTokens:  "5.",
// 		expectedCustomLimiterTokens:   "3.",
// 	})
// }

func runLnRTokenPoolReceiveFlowTest(t *testing.T, tc lnrTokenPoolReceiveFlowTestCase) {
	t.Helper()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(3))

	// Setup participants
	ccipParticipant := env.Chain.Participants[0]
	receiverParticipant := env.Chain.Participants[1]
	tokenPoolOwnerParticipant := env.Chain.Participants[0]

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

	// Upload DARs to all participants
	dars := [][]byte{rmnDar, commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, tokenPoolDar, perPartyRouterDar, ccipReceiverDar, onRampDar, feeQuoterDar}
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

	// Get DSO Admin Party (registry admin)
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err)

	// Token Setup
	// Mint tokens to Token Pool Owner (these will be "locked" in the pool)
	poolHoldingCid, err := testhelpers.MintAMT(t.Context(), tokenPoolOwnerParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partyTokenPoolOwner, "100.00")
	require.NoError(t, err)
	t.Logf("Minted 100 AMT to Pool Owner, Holding CID: %s", poolHoldingCid)

	// Instrument ID for AMT
	instrumentIdAmt := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: registryAdmin}}},
		{Label: "id", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "Amulet"}}},
	}}}}

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

	versionTag := "49ff34ed"
	ccvQualifier := "default"

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
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						CcipOwner:      "",
						RmnOwner:       types.PARTY(partyCCIP),
						CursedSubjects: nil,
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
	for i, address := range cldfEnv.DataStore.Addresses().Filter() {
		t.Logf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s, Labels=%s",
			i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier, address.Labels.String())
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
	tokenAdminRegistry, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(token_admin_registry.ContractType), token_admin_registry.Version, ""))
	require.NoError(t, err, "failed to get TokenAdminRegistry address")
	rmnRemote, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	require.NoError(t, err, "failed to get RMNRemote address")
	perPartyRouterFactory, err := cldfEnv.DataStore.Addresses().Get(datastore.NewAddressRefKey(env.Chain.ChainSelector(), datastore.ContractType(per_party_router_factory.ContractType), per_party_router_factory.Version, ""))
	require.NoError(t, err, "failed to get PerPartyRouterFactory address")

	// Token Pool Setup
	// Deploy default inbound RateLimiter required by ReleaseFromTicket receive flow.
	// Keep it enabled but lower-capacity so the test fails if the default-finality limiter
	// is selected for this FTF transfer instead of the custom-finality limiter.
	inboundRateLimiterInstanceID := "test-pool-receive-inbound-rl"
	res, err := tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: inboundRateLimiterInstanceID}}},
						{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "remoteChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}}},
						{Label: "direction", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
						{Label: "mode", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitMode_DefaultFinality"}}}},
						{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}}},
						{Label: "capacity", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tc.defaultInboundLimiterCapacity}}},
						{Label: "rate", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tc.defaultInboundLimiterCapacity}}},
						{Label: "tokens", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tc.defaultInboundLimiterCapacity}}},
						{Label: "lastUpdated", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().UnixMicro()}}},
					}},
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	inboundRateLimiterCid := extractCreatedContractId(res)
	t.Logf("Deployed default inbound RateLimiter: %s", inboundRateLimiterCid)
	inboundCustomBlockConfirmationsRateLimiterInstanceID := "test-pool-receive-inbound-custom-rl"
	res, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: inboundCustomBlockConfirmationsRateLimiterInstanceID}}},
						{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "remoteChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}}},
						{Label: "direction", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
						{Label: "mode", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitMode_CustomFinality"}}}},
						{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}}},
						{Label: "capacity", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tc.customInboundLimiterCapacity}}},
						{Label: "rate", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tc.customInboundLimiterCapacity}}},
						{Label: "tokens", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: tc.customInboundLimiterCapacity}}},
						{Label: "lastUpdated", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().UnixMicro()}}},
					}},
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	inboundCustomBlockConfirmationsRateLimiterCid := extractCreatedContractId(res)
	t.Logf("Deployed custom-finality inbound RateLimiter: %s", inboundCustomBlockConfirmationsRateLimiterCid)
	outboundRateLimiterInstanceID := "test-pool-receive-outbound-rl"
	res, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: outboundRateLimiterInstanceID}}},
						{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "remoteChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}}},
						{Label: "direction", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
						{Label: "mode", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitMode_DefaultFinality"}}}},
						{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
						{Label: "capacity", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
						{Label: "rate", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
						{Label: "tokens", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
						{Label: "lastUpdated", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().UnixMicro()}}},
					}},
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	outboundRateLimiterCid := extractCreatedContractId(res)
	t.Logf("Deployed outbound RateLimiter: %s", outboundRateLimiterCid)

	remotePoolAddress := hexutil.MustDecode("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3")
	remoteTokenAddress := hexutil.MustDecode("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d")

	// Deploy LockReleaseTokenPool
	lrtpInstanceAddress := contracts.InstanceID("test-pool-receive").RawInstanceAddress(types.PARTY(partyTokenPoolOwner)).InstanceAddress()
	res, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-lockreleasetokenpool", ModuleName: "CCIP.LockReleaseTokenPool", EntityName: "LockReleaseTokenPool"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "decimals", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 6}}},
						{Label: "rateLimitAdmin", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{}}}},
						{Label: "remoteChainConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{{
							Key: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}},
							Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "remotePools", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
									{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(remotePoolAddress)}},
								}}}}},
								{Label: "remoteTokenAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(remoteTokenAddress)}}},
								{Label: "inboundCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
								{Label: "outboundCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
								{Label: "minBlockConfirmations", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2000}}},
								{Label: "inboundRateLimiter", Value: rawInstanceAddress(inboundRateLimiterInstanceID + "@" + partyTokenPoolOwner)},
								{Label: "inboundCustomBlockConfirmationsRateLimiter", Value: rawInstanceAddress(inboundCustomBlockConfirmationsRateLimiterInstanceID + "@" + partyTokenPoolOwner)},
								{Label: "outboundRateLimiter", Value: rawInstanceAddress(outboundRateLimiterInstanceID + "@" + partyTokenPoolOwner)},
							}}}},
						}}}}}},
						{Label: "tokenTransferFeeConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "poolReceiveContext", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
						}}}}},
						// Use RelativeHours 24 for Amulet tokens (have expiry constraints)
						{Label: "transferTimeout", Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
							Constructor: "RelativeHours",
							Value:       &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 24}},
						}}}},
						{Label: "deps", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "tokenAdminRegistry", Value: rawInstanceAddress(tokenAdminRegistry.Labels.List()[0])},
							{Label: "rmnRemote", Value: rawInstanceAddress(rmnRemote.Labels.List()[0])},
							{Label: "feeQuoter", Value: rawInstanceAddress(feeQuoter.Labels.List()[0])},
						}}}}},
					}},
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	tokenPoolCid := extractCreatedContractId(res)
	t.Logf("Deployed LockReleaseTokenPool: %s", tokenPoolCid)

	// Register pool in TokenAdminRegistry (ProposeAdministrator -> AcceptAdminRole -> SetPool)
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)
	// Step 1: ProposeAdministrator
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_ProposeAdministrator",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "newAdmin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
					}}}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	t.Log("Called ProposeAdministrator")

	// Step 2: AcceptAdminRole
	time.Sleep(500 * time.Millisecond)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	_, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_AcceptAdminRole",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	t.Log("Called AcceptAdminRole")

	// Step 3: SetPool
	time.Sleep(500 * time.Millisecond)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	_, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_SetPool",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "tokenPool", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
							{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						}}}}}}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	t.Log("Called SetPool")

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
				Port: uint16(edsPort), //nolint:gosec // this is a port number
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
			Contracts: config.Contracts{
				PerPartyRouterFactory: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(perPartyRouterFactory.Address),
				},
				OnRamp: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(onRamp.Address),
				},
				OffRamp: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(offRamp.Address),
				},
				GlobalConfig: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(globalConfig.Address),
				},
				TokenAdminRegistry: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistry.Address),
				},
				RMNRemote: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(rmnRemote.Address),
				},
				FeeQuoter: config.ContractIdentifier{
					PartyID:         partyCCIP,
					InstanceAddress: contracts.HexToInstanceAddress(feeQuoter.Address),
				},
				CCVs: []config.ContractIdentifier{
					{
						PartyID:         partyCCIP,
						InstanceAddress: contracts.HexToInstanceAddress(committeeVerifier.Address),
					},
				},
				PoolOwner: partyTokenPoolOwner,
				TokenPools: []config.ContractIdentifier{
					{
						PartyID:         partyCCIP,
						InstanceAddress: lrtpInstanceAddress,
					},
				},
			},
		})
		log.Info().Msg("EDS terminated")
		if err != nil {
			log.Error().Err(err).Msg("EDS server exited with error")
		}
	}()
	// Create EDS client
	edsClient, err := edsv1.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "Failed to create EDS client")

	// Deploy and configure lane
	committeeVerifierRawAddr, err := contracts.RawInstanceAddressFromString(committeeVerifier.Labels.List()[0])
	require.NoError(t, err, "failed to parse CommitteeVerifier raw address")
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	out, err = changesets.ConfigureChainForLanes{}.Apply(cldfEnv, changesets.CantonCSDeps[changesets.ConfigureChainForLanesConfig]{
		ChainSelector: env.Chain.ChainSelector(),
		Participant:   0,
		Config: changesets.ConfigureChainForLanesConfig{
			Input: sequences.ConfigureChainForLanesInput{
				ChainSelector: env.Chain.ChainSelector(),
				GlobalConfig:  contracts.HexToInstanceAddress(globalConfig.Address),
				FeeQuoter:     contracts.HexToInstanceAddress(feeQuoter.Address),
				OnRamp:        contracts.HexToInstanceAddress(onRamp.Address),
				OffRamp:       contracts.HexToInstanceAddress(offRamp.Address),
				CommitteeVerifiers: []adapters.CommitteeVerifierConfig[contracts.InstanceAddress]{
					{
						CommitteeVerifier: []contracts.InstanceAddress{contracts.HexToInstanceAddress(committeeVerifier.Address)},
						RemoteChains: map[uint64]adapters.CommitteeVerifierRemoteChainConfig{
							remoteSelector: {
								AllowlistEnabled:   false,
								FeeUSDCents:        50,
								GasForVerification: 50_000,
								PayloadSizeBytes:   6*64 + 2*32,
								SignatureConfig: adapters.CommitteeVerifierSignatureQuorumConfig{
									Signers:   ccvSignerPubKeys,
									Threshold: 2,
								},
							},
						},
					},
				},
				RemoteChains: map[uint64]adapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]{
					remoteSelector: {
						AllowTrafficFrom:         true,
						OnRamps:                  [][]byte{hexutil.MustDecode("0xf6eced5e96fff2de4f0ecd722beb57556fc443fd")},
						OffRamp:                  hexutil.MustDecode("0xd8c9ec8cad3fb34aeca3ddbebfabe9f28a9bfaed"),
						DefaultInboundCCVs:       []contracts.RawInstanceAddress{committeeVerifierRawAddr},
						LaneMandatedInboundCCVs:  nil,
						DefaultOutboundCCVs:      []contracts.RawInstanceAddress{committeeVerifierRawAddr},
						LaneMandatedOutboundCCVs: nil,
						DefaultExecutor:          "",
						FeeQuoterDestChainConfig: adapters.FeeQuoterDestChainConfig{
							IsEnabled:                   true,
							MaxDataBytes:                50000,
							MaxPerMsgGasLimit:           4000000,
							DestGasOverhead:             300000,
							DestGasPerPayloadByteBase:   16,
							ChainFamilySelector:         [4]byte{0x28, 0x12, 0xd5, 0x2c},
							DefaultTxGasLimit:           200000,
							LinkFeeMultiplierPercent:    90,
							DefaultTokenFeeUSDCents:     0,
							DefaultTokenDestGasOverhead: 34000,
							NetworkFeeUSDCents:          0,
						},
						ExecutorDestChainConfig: adapters.ExecutorDestChainConfig{},
						AddressBytesLength:      20,
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

	// wait for EDS to start up
	time.Sleep(20 * time.Second)

	// Create PerPartyRouter for receiver via EDS
	perPartyRouterFactoryCid, disclosedContracts, err := testhelpers.GetPerPartyRouterFactoryDisclosures(t.Context(), edsClient)
	require.NoError(t, err)
	t.Logf("disclosed contracts: %v", disclosedContracts)
	t.Logf("per party router factory cid: %s", perPartyRouterFactoryCid)

	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: perPartyRouterFactoryCid,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-receiver"}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: disclosedContracts,
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
	// Encode instrumentId for destTokenAddress
	encodedInstrumentId := encodeInstrumentId(registryAdmin, "Amulet")
	hashedInstrumentId := crypto.Keccak256(encodedInstrumentId)

	// Build token transfer (5 AMT in Splice Decimal format)
	encodedTokenTransfer := buildTokenTransferV1(tc.tokenAmount, remotePoolAddress, remoteTokenAddress, hashedInstrumentId, partyReceiver, tc.sourcePoolData)

	// Build message
	msg := &MessageV1{
		SourceChainSelector: remoteSelector,
		DestChainSelector:   env.Chain.ChainSelector(),
		SequenceNumber:      1,
		ExecutionGasLimit:   200000,
		CCIPReceiveGasLimit: 100000,
		Finality:            2000,
		CCVAndExecutorHash:  [32]byte{},
		OnRampAddress:       gethcommon.LeftPadBytes(hexutil.MustDecode("0xf6eced5e96fff2de4f0ecd722beb57556fc443fd"), 32),
		OffRampAddress:      contracts.HexToInstanceAddress(offRamp.Address).Bytes(),
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
						{Label: "minBlockConfirmations", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2000}}},
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

	// Get TransferFactory and pool holdings (needed by CCIPReceiver.Execute for token release)
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, registryAdmin, partyTokenPoolOwner, partyReceiver)
	require.NoError(t, err)

	poolHoldings, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), tokenPoolOwnerParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	require.NotEmpty(t, poolHoldings, "Pool should have holdings")

	// Filter for unlocked holdings only (Amulet tokens may be locked until next mining round)
	// var poolHoldingCids []*apiv2.Value
	// var poolHoldingDisclosures []*apiv2.DisclosedContract
	// var lockedCount, unlockedCount int
	// for _, h := range poolHoldings {
	// 	viewFields := h.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()
	// 	lockField := viewFields[3].GetValue()
	// 	isLocked := lockField.GetOptional().GetValue() != nil

	// 	if isLocked {
	// 		lockedCount++
	// 		continue
	// 	}
	// 	unlockedCount++

	// 	poolHoldingCids = append(poolHoldingCids, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: h.GetCreatedEvent().GetContractId()}})
	// 	poolHoldingDisclosures = append(poolHoldingDisclosures, &apiv2.DisclosedContract{
	// 		ContractId:       h.GetCreatedEvent().GetContractId(),
	// 		TemplateId:       h.GetCreatedEvent().GetTemplateId(),
	// 		CreatedEventBlob: h.GetCreatedEvent().GetCreatedEventBlob(),
	// 	})
	// }
	// t.Logf("Pool has %d holdings (%d unlocked, %d locked)", len(poolHoldings), unlockedCount, lockedCount)

	// if unlockedCount == 0 {
	// 	t.Skip("SKIPPING: All pool holdings are locked (Amulet tokens are locked until next mining round). " +
	// 		"This is expected on fresh localnet. Either wait for mining round to complete, or use TestToken instead of Amulet.")
	// }

	// Capture receiver's balance before execute
	receiverHoldingsBefore, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceBefore := getHoldingsBalance(receiverHoldingsBefore)

	// Get disclosures for CCIPReceiver.Execute. The execute submission itself stays
	// receiver-only; pool/ccip-owned contracts are only supplied via disclosure.
	// disclosedCCIPReceiver, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
	// 	PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver",
	// })
	// require.NoError(t, err)
	// disclosedRouter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
	// 	PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	// })
	// require.NoError(t, err)

	// TODO: EDS call to get execute disclosures
	executeDisclosuresEDS, err := testhelpers.GetCCIPExecuteDisclosures(
		t.Context(),
		encodedMessageHex,
		edsClient,
		[]contracts.InstanceAddress{
			contracts.HexToInstanceAddress(committeeVerifier.Address),
		},
	)
	require.NoError(t, err)
	require.Len(t, executeDisclosuresEDS.CCVContractIDs, 1)

	executeDisclosures := slices.Concat(
		executeDisclosuresEDS.DisclosedContracts,
		transferFactoryDisclosures, // not from EDS
	)

	// CCIPReceiver.Execute: PrepareExecute + CCV + Pool Verify + Execute + Release
	// in one receiver-authored transaction with disclosed shared dependencies.
	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId: ccipReceiverCid,
					Choice:     "Execute",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "context", Value: choiceContext},
						{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCid}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: encodedMessageHex}}},
						{Label: "tokenTransfer", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "tokenPoolCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: executeDisclosuresEDS.TokenPoolContractID.ContractId}}},
							{Label: "tokenReceiverParty", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
							{Label: "tokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}}},
								{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
									{Label: "context", Value: choiceContext},
									{Label: "meta", Value: emptyMetadata},
								}}}}},
								{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: executeDisclosuresEDS.TokenPoolHoldingsContractIDs}}}},
							}}}}},
							{Label: "poolExtraContext", Value: executeDisclosuresEDS.PoolExtraContext},
						}}}}}}}},
						{Label: "ccvInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "ccvCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: executeDisclosuresEDS.CCVContractIDs[0].ContractId}}},
								{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: verifierResultsHex}}},
								{Label: "ccvExtraContext", Value: emptyCCIPContext},
							}}}},
						}}}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: executeDisclosures,
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

		acceptContextResp, err := transferInstructionClient.GetTransferInstructionAcceptContextWithResponse(t.Context(), pendingTransferInstructionCid, transferInstructionV1.GetChoiceContextRequest{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, acceptContextResp.StatusCode(), "Failed to get accept context: %s", string(acceptContextResp.Body))

		acceptDisclosures := make([]*apiv2.DisclosedContract, 0, len(acceptContextResp.JSON200.DisclosedContracts))
		for _, contract := range acceptContextResp.JSON200.DisclosedContracts {
			id, err := testhelpers.TemplateIdFromString(contract.TemplateId)
			require.NoError(t, err)
			createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
			require.NoError(t, err)
			acceptDisclosures = append(acceptDisclosures, &apiv2.DisclosedContract{
				TemplateId:       id,
				ContractId:       contract.ContractId,
				CreatedEventBlob: createdEventBlob,
			})
		}
		acceptContext, err := testhelpers.ChoiceContextFromData(acceptContextResp.JSON200.ChoiceContextData)
		require.NoError(t, err)

		_, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.Must(uuid.NewUUID()).String(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
						ContractId: pendingTransferInstructionCid,
						Choice:     "TransferInstruction_Accept",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "context", Value: acceptContext},
								{Label: "meta", Value: emptyMetadata},
							}}}}},
						}}}},
					}},
				}},
				ActAs:              []string{partyReceiver},
				DisclosedContracts: acceptDisclosures,
			},
		})
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
		defaultRateLimiter, err := findActiveRateLimiterByInstanceID(t.Context(), tokenPoolOwnerParticipant, inboundRateLimiterInstanceID)
		require.NoError(t, err)
		require.Equal(t, tc.expectedDefaultLimiterTokens, getRateLimiterTokens(defaultRateLimiter))
	}
	if tc.expectedCustomLimiterTokens != "" {
		customRateLimiter, err := findActiveRateLimiterByInstanceID(t.Context(), tokenPoolOwnerParticipant, inboundCustomBlockConfirmationsRateLimiterInstanceID)
		require.NoError(t, err)
		require.Equal(t, tc.expectedCustomLimiterTokens, getRateLimiterTokens(customRateLimiter))
	}
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
	sourceTokenAddress,
	destTokenAddress []byte,
	tokenReceiverParty string,
	extraData []byte,
) *TokenTransferV1 {
	return &TokenTransferV1{
		Amount:             amount,
		SourcePoolAddress:  sourcePoolAddress,
		SourceTokenAddress: sourceTokenAddress,
		DestTokenAddress:   destTokenAddress,
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

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}

var emptyCCIPContext = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}
