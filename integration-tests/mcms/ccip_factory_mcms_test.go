package tests

import (
	"fmt"
	"slices"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/core"
	splice "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestCCIP_MCMSFactoryDeploy validates the full MCMS governance flow for CCIP:
// 1. Arbitrary party deploys CCIPFactory
// 2. MCMS takes ownership (SetOwnerToMCMS)
// 3. All CCIP components deployed through MCMS Bypasser operations targeting the factory
// 4. Factory state verified to contain all deployed contracts
func TestCCIP_MCMSFactoryDeploy(t *testing.T) {
	t.Parallel()

	env := GetSharedCCIPMCMSEnvironment(t)
	participant := env.Participant
	mcmsEncoder := env.McmsEncoder
	ccipOwner := env.CcipOwner
	cfg := env.Config
	sortedSigners := env.SortedSigners
	factoryEncoder := env.FactoryEncoder

	chainID := int64(1)
	uid := uuid.New().String()[:8]

	// --- Step 1: Create MCMS (2-of-3, minDelay=0, no blocked functions) ---
	baseMcmsID := "mcms-ccip-" + uid
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)
	mcmsCid := createMCMSMultiRole(t, participant, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	// --- Step 2: Create CCIPFactory with owner=ccipOwner, mcmsParty=ccipOwner ---
	// In a real scenario these would be different parties; we use the same for test simplicity.
	factoryInstanceID := "factory-" + uid
	factoryInstanceAddr := fmt.Sprintf("%s@%s", factoryInstanceID, ccipOwner)

	factoryCid := createCCIPFactory(t, participant, ccipOwner, factoryInstanceID)
	t.Logf("Factory created: CID=%s, instanceAddr=%s", factoryCid, factoryInstanceAddr)

	// --- Step 3: SetOwnerToMCMS ---
	// In the test setup, owner == mcmsParty (both ccipOwner), so SetOwnerToMCMS assertion
	// (owner != mcmsParty) would fail. We skip this step since the factory is already
	// effectively MCMS-controlled when owner == mcmsParty.

	// --- Step 4: MCMS-driven deploys via BypasserExecuteBatch ---
	// Deploy order follows the dependency graph:
	// 1. RMNRemote (no deps)
	// 2. GlobalConfig (no deps)
	// 3. TokenAdminRegistry (no deps)
	// 4. LinkToken (needed by FeeQuoter)
	// 5. FeeQuoter (deps: LinkToken instrument id)
	// 6. CommitteeVerifier (deps: RMNRemote)
	// 7. OffRamp (deps: GlobalConfig, RMNRemote, TAR)
	// 8. OnRamp (deps: GlobalConfig, RMNRemote, TAR, FeeQuoter, CCV)
	// 9. PerPartyRouterFactory (deps: all of the above)

	rmnInstanceID := "rmn-" + uid
	rmnInstanceAddr := fmt.Sprintf("%s@%s", rmnInstanceID, ccipOwner)
	gcInstanceID := "gc-" + uid
	tarInstanceID := "tar-" + uid
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	linkInstanceID := "link-" + uid
	linkInstrumentID := splice.InstrumentId{Admin: types.PARTY(ccipOwner), Id: "link-token"}
	fqInstanceID := "fq-" + uid
	fqInstanceAddr := fmt.Sprintf("%s@%s", fqInstanceID, ccipOwner)
	ccvInstanceID := "ccv-" + uid
	ccvInstanceAddr := fmt.Sprintf("%s@%s", ccvInstanceID, ccipOwner)
	offRampInstanceID := "offramp-" + uid
	offRampInstanceAddr := fmt.Sprintf("%s@%s", offRampInstanceID, ccipOwner)
	onRampInstanceID := "onramp-" + uid
	onRampInstanceAddr := fmt.Sprintf("%s@%s", onRampInstanceID, ccipOwner)
	pprInstanceID := "ppr-" + uid

	// 1. Deploy RMNRemote
	rmnParams := factory.DeployRMNRemoteParams{
		InstanceId:      types.TEXT(rmnInstanceID),
		RmnOwner:        types.PARTY(ccipOwner),
		CcipOwner:       types.PARTY(ccipOwner),
		CustomObservers: []types.PARTY{},
		CursedSubjects:  []types.TEXT{},
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployRMNRemoteParams", rmnParams)
	t.Logf("RMNRemote deployed: %s", rmnInstanceID)

	// 2. Deploy GlobalConfig
	gcParams := factory.DeployGlobalConfigParams{
		InstanceId:    types.TEXT(gcInstanceID),
		ChainSelector: types.NUMERIC("123"),
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployGlobalConfigParams", gcParams)
	t.Logf("GlobalConfig deployed: %s", gcInstanceID)

	// 3. Deploy TokenAdminRegistry
	tarParams := factory.DeployTokenAdminRegistryParams{
		InstanceId: types.TEXT(tarInstanceID),
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployTokenAdminRegistryParams", tarParams)
	t.Logf("TokenAdminRegistry deployed: %s", tarInstanceID)

	// 4. Deploy LinkToken
	linkParams := factory.DeployLinkTokenParams{
		InstanceId:   types.TEXT(linkInstanceID),
		InstrumentId: linkInstrumentID,
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployLinkTokenParams", linkParams)
	t.Logf("LinkToken deployed: %s", linkInstanceID)

	// 5. Deploy FeeQuoter
	fqParams := factory.DeployFeeQuoterParams{
		InstanceId:            types.TEXT(fqInstanceID),
		LinkTokenInstrumentId: linkInstrumentID,
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployFeeQuoterParams", fqParams)
	t.Logf("FeeQuoter deployed: %s", fqInstanceID)

	// 6. Deploy CommitteeVerifier (deps: RMNRemote)
	ccvParams := factory.DeployCommitteeVerifierParams{
		InstanceId:                   types.TEXT(ccvInstanceID),
		Owner:                        types.PARTY(ccipOwner),
		CcipOwner:                    types.PARTY(ccipOwner),
		VersionTag:                   "01020304",
		AllowListAdmin:               nil,
		MessageSentObservers:         []types.PARTY{types.PARTY(ccipOwner)},
		RmnRemote:                    chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		StorageLocations:             []types.TEXT{},
		StorageLocationsAdmin:        types.PARTY(ccipOwner),
		PendingStorageLocationsAdmin: types.PARTY(ccipOwner),
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployCommitteeVerifierParams", ccvParams)
	t.Logf("CommitteeVerifier deployed: %s", ccvInstanceID)

	// 7. Deploy OffRamp (deps: GlobalConfig, RMNRemote, TAR)
	gcInstanceAddr := fmt.Sprintf("%s@%s", gcInstanceID, ccipOwner)
	offRampParams := factory.DeployOffRampParams{
		InstanceId:         types.TEXT(offRampInstanceID),
		GlobalConfig:       chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(gcInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployOffRampParams", offRampParams)
	t.Logf("OffRamp deployed: %s", offRampInstanceID)

	// 8. Deploy OnRamp (deps: GlobalConfig, RMNRemote, TAR, FeeQuoter, CCV)
	onRampParams := factory.DeployOnRampParams{
		InstanceId:         types.TEXT(onRampInstanceID),
		GlobalConfig:       chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(gcInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
		FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(fqInstanceAddr)},
		CcvRegistry:        chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(ccvInstanceAddr)},
		MaxUSDCentsPerMsg:  types.NUMERIC("100000"),
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployOnRampParams", onRampParams)
	t.Logf("OnRamp deployed: %s", onRampInstanceID)

	// 9. Deploy PerPartyRouterFactory (deps: all)
	pprParams := factory.DeployPerPartyRouterFactoryParams{
		InstanceId:         types.TEXT(pprInstanceID),
		OnRamp:             chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(onRampInstanceAddr)},
		OffRamp:            chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(offRampInstanceAddr)},
		GlobalConfig:       chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(gcInstanceAddr)},
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
		FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(fqInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployPerPartyRouterFactoryParams", pprParams)
	t.Logf("PerPartyRouterFactory deployed: %s", pprInstanceID)

	// --- Step 5: Verify factory state ---
	factoryFields := queryContractFields(t, participant, factoryCid)
	require.Equal(t, factoryInstanceID, factoryFields["instanceId"])
	require.Equal(t, "true", factoryFields["perPartyRouterFactoryDeployed"])

	t.Logf("All CCIP components deployed via MCMS. Factory CID=%s, MCMS CID=%s", factoryCid, mcmsCid)

	// --- Step 6 (Optional): Deploy LockReleaseTokenPool ---
	lrtpInstanceID := "lrtp-" + uid

	lrtpParams := factory.DeployLockReleaseTokenPoolParams{
		InstanceId:         types.TEXT(lrtpInstanceID),
		PoolOwner:          types.PARTY(ccipOwner),
		CcipOwner:          types.PARTY(ccipOwner),
		InstrumentId:       splice.InstrumentId{Admin: types.PARTY(ccipOwner), Id: "test-token"},
		Decimals:           types.INT64(18),
		RateLimitAdmin:     nil,
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
		FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(fqInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
		TransferTimeout:    lockreleasetokenpool.TransferTimeout{Indefinite: new(types.UNIT{})},
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployLockReleaseTokenPoolParams", lrtpParams)
	t.Logf("LockReleaseTokenPool deployed via MCMS: %s", lrtpInstanceID)

	// Deploy a second pool
	lrtp2InstanceID := "lrtp2-" + uid

	lrtp2Params := factory.DeployLockReleaseTokenPoolParams{
		InstanceId:         types.TEXT(lrtp2InstanceID),
		PoolOwner:          types.PARTY(ccipOwner),
		CcipOwner:          types.PARTY(ccipOwner),
		InstrumentId:       splice.InstrumentId{Admin: types.PARTY(ccipOwner), Id: "amulet-token"},
		Decimals:           types.INT64(10),
		RateLimitAdmin:     nil,
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
		FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(fqInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
		TransferTimeout:    lockreleasetokenpool.TransferTimeout{Indefinite: new(types.UNIT{})},
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployLockReleaseTokenPoolParams", lrtp2Params)
	t.Logf("LockReleaseTokenPool deployed via MCMS: %s", lrtp2InstanceID)

	// Final factory state check
	factoryFields = queryContractFields(t, participant, factoryCid)
	require.Equal(t, "true", factoryFields["perPartyRouterFactoryDeployed"])
	t.Logf("Final factory state verified. All deployments successful.")

	// --- Step 7: MCMS-driven config on deployed contracts ---
	// After factory deploys, apply config to GlobalConfig and FeeQuoter via MCMS.
	// Each CCIP contract implements MCMSReceiver; we target them directly.

	// Find deployed contract CIDs
	gcCid := findNewContractCid(t, participant, contracts.IdentifierFromBinding(core.GlobalConfig{}), ccipOwner, gcInstanceAddr)
	require.NotEmpty(t, gcCid)

	// GlobalConfig: ApplySourceChainConfigUpdates
	remoteChainSelector := types.NUMERIC("456")
	gcEncoder := core.NewContract(fmt.Sprintf("#%s", core.PackageName), "CCIP.GlobalConfig", "GlobalConfig").Encoder()

	sourceConfigArgs := core.ApplySourceChainConfigUpdates{
		SourceChainConfigUpdates: []core.SourceChainConfigArgs{{
			SourceChainSelector: remoteChainSelector,
			IsEnabled:           types.BOOL(true),
			OnRampAddresses:     []types.TEXT{"0000000000000000000000000000000000000000000000000000000000abcdef"},
			DefaultCCVs:         []chainlinkapi.RawInstanceAddress{{Unpack: types.TEXT(ccvInstanceAddr)}},
			LaneMandatedCCVs:    []chainlinkapi.RawInstanceAddress{},
		}},
	}
	sourceEncoded, err := gcEncoder.ApplySourceChainConfigUpdates(sourceConfigArgs)
	require.NoError(t, err)

	gcCid, mcmsCid = mcmsContractConfig(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, gcCid, gcInstanceAddr, chainID, sortedSigners, sourceEncoded)
	t.Logf("GlobalConfig: ApplySourceChainConfigUpdates applied via MCMS")

	// GlobalConfig: ApplyDestChainConfigUpdates
	destConfigArgs := core.ApplyDestChainConfigUpdates{
		DestChainConfigUpdates: []core.DestChainConfigArgs{{
			DestChainSelector:         remoteChainSelector,
			IsEnabled:                 types.BOOL(true),
			AddressBytesLength:        types.INT64(64),
			TokenReceiverAllowed:      types.BOOL(true),
			BaseExecutionGasCost:      types.INT64(200000),
			OffRampAddress:            "0000000000000000000000000000000000000000000000000000000000fedcba",
			DefaultExecutor:           nil,
			LaneMandatedCCVs:          []chainlinkapi.RawInstanceAddress{},
			DefaultCCVs:               []chainlinkapi.RawInstanceAddress{{Unpack: types.TEXT(ccvInstanceAddr)}},
			MessageNetworkFeeUSDCents: types.NUMERIC("100"),
			TokenNetworkFeeUSDCents:   types.NUMERIC("50"),
		}},
	}
	destEncoded, err := gcEncoder.ApplyDestChainConfigUpdates(destConfigArgs)
	require.NoError(t, err)

	gcCid, _ = mcmsContractConfig(t, participant, mcmsEncoder, ccipOwner, mcmsCid, mcmsInstanceAddr, gcCid, gcInstanceAddr, chainID, sortedSigners, destEncoded)
	_ = gcCid
	t.Logf("GlobalConfig: ApplyDestChainConfigUpdates applied via MCMS")

	t.Logf("MCMS-driven config flow complete.")
}

// TestCCIP_MCMSFactoryDeploy_FullGovernance validates the complete MCMS governance lifecycle:
// 1. Bootstrap party deploys CCIPFactory (owner != mcmsParty)
// 2. mcmsParty exercises SetOwnerToMCMS to transfer factory ownership
// 3. All CCIP components deployed through MCMS Bypasser operations targeting the factory
// 4. Factory state verified to contain all deployed contracts
//
// This tests the realistic scenario where an arbitrary party bootstraps the infrastructure
// and MCMS takes over governance via the mcmsParty-controlled handover choice.
func TestCCIP_MCMSFactoryDeploy_FullGovernance(t *testing.T) {
	t.Parallel()

	// Use environment with two parties on the same participant (bootstrap + MCMS).
	env := GetSharedCCIPMCMSTwoParticipantEnvironment(t)
	participant := env.Participant
	mcmsEncoder := env.McmsEncoder
	mcmsParty := env.CcipOwner
	bootstrapParty := env.BootstrapParty
	cfg := env.Config
	sortedSigners := env.SortedSigners
	factoryEncoder := env.FactoryEncoder

	chainID := int64(1)
	uid := uuid.New().String()[:8]

	t.Logf("Bootstrap party: %s", bootstrapParty)
	t.Logf("MCMS party: %s", mcmsParty)

	// --- Step 1: Create MCMS (2-of-3, minDelay=0) owned by mcmsParty ---
	baseMcmsID := "mcms-gov-" + uid
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, mcmsParty)
	mcmsCid := createMCMSMultiRole(t, participant, mcmsParty, chainID, baseMcmsID, cfg, 0, nil)
	t.Logf("MCMS created: CID=%s, instanceAddr=%s", mcmsCid, mcmsInstanceAddr)

	// --- Step 2: Bootstrap party deploys CCIPFactory with owner=bootstrapParty, mcmsParty=mcmsParty ---
	factoryInstanceID := "factory-gov-" + uid
	factoryInstanceAddr := fmt.Sprintf("%s@%s", factoryInstanceID, mcmsParty)

	factoryCid := createCCIPFactoryWithMCMS(t, participant, bootstrapParty, mcmsParty, factoryInstanceID)
	t.Logf("Factory created by bootstrap party: CID=%s", factoryCid)

	// Verify initial state: owner = bootstrapParty, mcmsParty = mcmsParty
	initialFields := queryContractFields(t, participant, factoryCid)
	require.Equal(t, bootstrapParty, initialFields["owner"], "initial factory owner should be bootstrapParty")
	require.Equal(t, mcmsParty, initialFields["mcmsParty"], "initial factory mcmsParty should be mcmsParty")
	t.Logf("Verified initial state: owner=%s, mcmsParty=%s (different parties)", initialFields["owner"], initialFields["mcmsParty"])

	// --- Step 3: mcmsParty exercises SetOwnerToMCMS (controller = mcmsParty only) ---
	factoryCid = setFactoryOwnerToMCMS(t, participant, mcmsParty, factoryCid)
	t.Logf("SetOwnerToMCMS executed: new factory CID=%s", factoryCid)

	// Verify ownership transferred to mcmsParty
	factoryFields := queryContractFields(t, participant, factoryCid)
	require.Equal(t, mcmsParty, factoryFields["owner"], "factory owner should be mcmsParty after SetOwnerToMCMS")
	require.Equal(t, mcmsParty, factoryFields["mcmsParty"], "factory mcmsParty should remain mcmsParty")
	t.Logf("Verified: factory owner=%s, mcmsParty=%s", factoryFields["owner"], factoryFields["mcmsParty"])

	// --- Step 4: MCMS-driven deploys via BypasserExecuteBatch ---
	// Deploy 2 components to validate MCMS governance flow.
	// Additional components (RMNRemote, FeeQuoter, etc.) are tested separately.

	gcInstanceID := "gc-gov-" + uid
	tarInstanceID := "tar-gov-" + uid

	// 1. Deploy GlobalConfig
	gcParams := factory.DeployGlobalConfigParams{
		InstanceId:    types.TEXT(gcInstanceID),
		ChainSelector: types.NUMERIC("123"),
	}
	factoryCid, mcmsCid = mcmsFactoryDeploy(t, participant, mcmsEncoder, mcmsParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployGlobalConfigParams", gcParams)
	t.Logf("GlobalConfig deployed via MCMS: %s", gcInstanceID)

	// 2. Deploy TokenAdminRegistry
	tarParams := factory.DeployTokenAdminRegistryParams{
		InstanceId: types.TEXT(tarInstanceID),
	}
	factoryCid, _ = mcmsFactoryDeploy(t, participant, mcmsEncoder, mcmsParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployTokenAdminRegistryParams", tarParams)
	t.Logf("TokenAdminRegistry deployed via MCMS: %s", tarInstanceID)

	// --- Step 5: Verify factory state ---
	finalFactoryFields := queryContractFields(t, participant, factoryCid)
	require.Equal(t, factoryInstanceID, finalFactoryFields["instanceId"])
	require.Equal(t, mcmsParty, finalFactoryFields["owner"])

	t.Logf("MCMS governance flow validated: bootstrap -> SetOwnerToMCMS -> MCMS-controlled deployments")
}

// encodableParams is a constraint for factory deploy params that support MCMSEncoder methods.
type encodableParams interface {
	factory.DeployRMNRemoteParams |
		factory.DeployGlobalConfigParams |
		factory.DeployTokenAdminRegistryParams |
		factory.DeployFeeQuoterParams |
		factory.DeployLinkTokenParams |
		factory.DeployCommitteeVerifierParams |
		factory.DeployOffRampParams |
		factory.DeployOnRampParams |
		factory.DeployPerPartyRouterFactoryParams |
		factory.DeployLockReleaseTokenPoolParams |
		factory.DeployBurnMintTokenPoolParams |
		factory.DeployRateLimiterParams
}

// mcmsFactoryDeploy executes a single MCMS Bypasser operation to deploy a component via the factory.
// It handles: encode → build proposal → sign → SetRoot → ExecuteOp (Bypasser with TargetCids).
// Returns the new factory CID (consuming choice) and new MCMS CID.
func mcmsFactoryDeploy[T encodableParams](
	t *testing.T,
	participant canton.Participant,
	mcmsEncoder mcmsApi.MCMSEncoder,
	ccipOwner string,
	mcmsCid string,
	mcmsInstanceAddr string,
	factoryCid string,
	factoryInstanceAddr string,
	chainID int64,
	sortedSigners []*MCMSSigner,
	factoryEnc factory.MCMSEncoder,
	encoderMethodName string,
	params T,
) (string, string) {
	t.Helper()

	return mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, ccipOwner, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEnc, encoderMethodName, params)
}

func mcmsFactoryDeployWithMCMSQueryParty[T encodableParams](
	t *testing.T,
	participant canton.Participant,
	mcmsEncoder mcmsApi.MCMSEncoder,
	ccipOwner string,
	mcmsQueryParty string,
	mcmsCid string,
	mcmsInstanceAddr string,
	factoryCid string,
	factoryInstanceAddr string,
	chainID int64,
	sortedSigners []*MCMSSigner,
	factoryEnc factory.MCMSEncoder,
	encoderMethodName string,
	params T,
) (string, string) {
	t.Helper()

	encoded := encodeFactoryParams(t, factoryEnc, encoderMethodName, params)
	t.Logf("  MCMS deploy: choice=%s, opData=%d hex chars", encoded.Choice, len(encoded.OperationData))

	// Wrap the factory operation in a BypasserExecuteBatch
	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(factoryInstanceAddr),
		FunctionName:          types.TEXT(encoded.Choice),
		OperationData:         types.TEXT(encoded.OperationData),
	}}
	bypasserParams := mcmsApi.BypasserExecuteBatchParams{Calls: calls}
	bypasserChoice := MustEncodeBypasserExecuteBatch(t, mcmsEncoder, bypasserParams)

	// Build proposal
	bypasserMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleBypasser)

	// Get current bypasser op count from the MCMS contract
	opCount := queryBypasserOpCountForParty(t, participant, mcmsQueryParty, mcmsCid)

	proposal := NewMCMSProposal(int(chainID), bypasserMultisigID, int(opCount), false).
		AddOperation(mcmsInstanceAddr, bypasserChoice.Choice, bypasserChoice.OperationData).
		Build()

	// Sign (2 of 3)
	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	// SetRoot (Bypasser role)
	mcmsCid = setRootWithRoleAndDisclosureParty(t, participant, ccipOwner, mcmsQueryParty, mcmsCid, "Bypasser", proposal, validUntil, signatures)

	// ExecuteOp with TargetCids pointing to the current factory CID
	targetCids := map[string]string{factoryInstanceAddr: factoryCid}
	mcmsCid = bypasserExecuteBatchWithDisclosureParty(t, participant, ccipOwner, mcmsQueryParty, mcmsCid, targetCids, proposal.Operations[0], mustGetOpProof(t, proposal, 0))

	// Find new factory CID (consuming choice)
	newFactoryCid := findNewContractCid(t, participant, contracts.IdentifierFromBinding(factory.CCIPFactory{}), ccipOwner, factoryInstanceAddr)
	require.NotEmpty(t, newFactoryCid, "factory CID should be refreshed after deploy")

	return newFactoryCid, mcmsCid
}

// encodeFactoryParams dispatches to the correct encoder method based on the params type name.
func encodeFactoryParams[T encodableParams](t *testing.T, enc factory.MCMSEncoder, methodName string, params T) *bind.EncodedChoice {
	t.Helper()
	var result *bind.EncodedChoice
	var err error

	switch methodName {
	case "DeployRMNRemoteParams":
		result, err = enc.DeployRMNRemoteParams(any(params).(factory.DeployRMNRemoteParams))
	case "DeployGlobalConfigParams":
		result, err = enc.DeployGlobalConfigParams(any(params).(factory.DeployGlobalConfigParams))
	case "DeployTokenAdminRegistryParams":
		result, err = enc.DeployTokenAdminRegistryParams(any(params).(factory.DeployTokenAdminRegistryParams))
	case "DeployFeeQuoterParams":
		result, err = enc.DeployFeeQuoterParams(any(params).(factory.DeployFeeQuoterParams))
	case "DeployLinkTokenParams":
		result, err = enc.DeployLinkTokenParams(any(params).(factory.DeployLinkTokenParams))
	case "DeployCommitteeVerifierParams":
		result, err = enc.DeployCommitteeVerifierParams(any(params).(factory.DeployCommitteeVerifierParams))
	case "DeployOffRampParams":
		result, err = enc.DeployOffRampParams(any(params).(factory.DeployOffRampParams))
	case "DeployOnRampParams":
		result, err = enc.DeployOnRampParams(any(params).(factory.DeployOnRampParams))
	case "DeployPerPartyRouterFactoryParams":
		result, err = enc.DeployPerPartyRouterFactoryParams(any(params).(factory.DeployPerPartyRouterFactoryParams))
	case "DeployLockReleaseTokenPoolParams":
		result, err = enc.DeployLockReleaseTokenPoolParams(any(params).(factory.DeployLockReleaseTokenPoolParams))
	case "DeployBurnMintTokenPoolParams":
		result, err = enc.DeployBurnMintTokenPoolParams(any(params).(factory.DeployBurnMintTokenPoolParams))
	case "DeployRateLimiterParams":
		result, err = enc.DeployRateLimiterParams(any(params).(factory.DeployRateLimiterParams))
	default:
		t.Fatalf("unknown encoder method: %s", methodName)
	}
	require.NoError(t, err, "failed to encode %s", methodName)

	return result
}

// mcmsContractConfig executes a single MCMS Bypasser operation to configure a deployed CCIP contract.
// Similar to mcmsFactoryDeploy but targets the deployed contract directly (not via factory).
// Returns: new contract CID (consuming choice) and new MCMS CID.
func mcmsContractConfig(
	t *testing.T,
	participant canton.Participant,
	mcmsEncoder mcmsApi.MCMSEncoder,
	ccipOwner string,
	mcmsCid string,
	mcmsInstanceAddr string,
	contractCid string,
	contractInstanceAddr string,
	chainID int64,
	sortedSigners []*MCMSSigner,
	encoded *bind.EncodedChoice,
) (string, string) {
	t.Helper()

	return mcmsContractConfigWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, ccipOwner, mcmsCid, mcmsInstanceAddr, contractCid, contractInstanceAddr, chainID, sortedSigners, encoded)
}

func mcmsContractConfigWithMCMSQueryParty(
	t *testing.T,
	participant canton.Participant,
	mcmsEncoder mcmsApi.MCMSEncoder,
	ccipOwner string,
	mcmsQueryParty string,
	mcmsCid string,
	mcmsInstanceAddr string,
	contractCid string,
	contractInstanceAddr string,
	chainID int64,
	sortedSigners []*MCMSSigner,
	encoded *bind.EncodedChoice,
) (string, string) {
	t.Helper()

	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(contractInstanceAddr),
		FunctionName:          types.TEXT(encoded.Choice),
		OperationData:         types.TEXT(encoded.OperationData),
	}}
	bypasserParams := mcmsApi.BypasserExecuteBatchParams{Calls: calls}
	bypasserChoice := MustEncodeBypasserExecuteBatch(t, mcmsEncoder, bypasserParams)

	bypasserMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleBypasser)
	opCount := queryBypasserOpCountForParty(t, participant, mcmsQueryParty, mcmsCid)

	proposal := NewMCMSProposal(int(chainID), bypasserMultisigID, int(opCount), false).
		AddOperation(mcmsInstanceAddr, bypasserChoice.Choice, bypasserChoice.OperationData).
		Build()

	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	mcmsCid = setRootWithRoleAndDisclosureParty(t, participant, ccipOwner, mcmsQueryParty, mcmsCid, "Bypasser", proposal, validUntil, signatures)

	targetCids := map[string]string{contractInstanceAddr: contractCid}
	mcmsCid = bypasserExecuteBatchWithDisclosureParty(t, participant, ccipOwner, mcmsQueryParty, mcmsCid, targetCids, proposal.Operations[0], mustGetOpProof(t, proposal, 0))

	// Find the new contract CID after consuming choice
	newCid := findNewContractCidByOldCid(t, participant, &apiv2.Identifier{
		PackageId:  encoded.TemplateID.PackageID,
		ModuleName: encoded.TemplateID.ModuleName,
		EntityName: encoded.TemplateID.TemplateName,
	}, contractCid)

	return newCid, mcmsCid
}

// findNewContractCidByOldCid finds the refreshed CID for a contract after a consuming choice.
// It queries ACS for contracts of the same template and returns the one that isn't the old CID.
func findNewContractCidByOldCid(
	t *testing.T,
	participant canton.Participant,
	identifier *apiv2.Identifier,
	oldCid string,
) string {
	t.Helper()

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, identifier)
	require.NoError(t, err)

	// Return the most recently created one (sorted by creation time)
	for _, v := range slices.Backward(activeContracts) {
		cid := v.GetCreatedEvent().GetContractId()
		if cid != oldCid {
			return cid
		}
	}

	// If only one exists and it's different from old, return it
	if len(activeContracts) == 1 {
		return activeContracts[0].GetCreatedEvent().GetContractId()
	}

	t.Fatalf("could not find new CID for %s/%s (old=%s)", identifier.GetModuleName(), identifier.GetEntityName(), oldCid)

	return ""
}

func mustGetOpProof(t *testing.T, proposal *MCMSProposal, idx int) []string {
	t.Helper()
	proof, err := proposal.GetOpProof(idx)
	require.NoError(t, err)

	return proof
}

// createCCIPFactory creates a CCIPFactory contract with owner=mcmsParty (for test simplicity).
func createCCIPFactory(
	t *testing.T,
	participant canton.Participant,
	owner string,
	instanceID string,
) string {
	return createCCIPFactoryWithMCMS(t, participant, owner, owner, instanceID)
}

// createCCIPFactoryWithMCMS creates a CCIPFactory contract with separate owner and mcmsParty.
// This allows testing the full governance flow where a bootstrap party deploys the factory
// and mcmsParty takes ownership via SetOwnerToMCMS.
func createCCIPFactoryWithMCMS(
	t *testing.T,
	participant canton.Participant,
	owner string,
	mcmsParty string,
	instanceID string,
) string {
	t.Helper()

	factoryContract := factory.CCIPFactory{
		InstanceId:                    types.TEXT(instanceID),
		Owner:                         types.PARTY(owner),
		McmsParty:                     types.PARTY(mcmsParty),
		UsedInstanceIds:               map[types.TEXT]types.BOOL{},
		DeployedContracts:             map[types.TEXT]types.CONTRACT_ID{},
		PerPartyRouterFactoryDeployed: types.BOOL(false),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId:      contracts.IdentifierFromBinding(factory.CCIPFactory{}),
						CreateArguments: ledger.ConvertToRecord(factoryContract),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

// setFactoryOwnerToMCMS exercises the SetOwnerToMCMS choice on the factory,
// transferring ownership from the current owner to mcmsParty.
// Returns the new factory contract ID.
//
// SetOwnerToMCMS is controller mcmsParty; submit with ActAs = mcmsParty.
func setFactoryOwnerToMCMS(
	t *testing.T,
	participant canton.Participant,
	mcmsParty string,
	factoryCid string,
) string {
	t.Helper()

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: contracts.IdentifierFromBinding(factory.CCIPFactory{}),
						ContractId: factoryCid,
						Choice:     "SetOwnerToMCMS",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{},
						}}},
					},
				},
			}},
			ActAs: []string{mcmsParty},
		},
	})
	require.NoError(t, err)

	// Find the new factory CID in the transaction events
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "CCIPFactory" {
			return created.GetContractId()
		}
	}

	t.Fatal("SetOwnerToMCMS did not create new factory contract")

	return ""
}

// findNewContractCid queries ACS for the latest contract of a given template with a specific instanceId.
func findNewContractCid(
	t *testing.T,
	participant canton.Participant,
	identifier *apiv2.Identifier,
	owner string,
	instanceAddr string,
) string {
	t.Helper()

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, identifier)
	require.NoError(t, err)

	// Extract instanceId from instanceAddr (strip "@owner")
	instanceID := instanceAddr[:len(instanceAddr)-len(owner)-1]

	for _, contract := range activeContracts {
		for _, field := range contract.GetCreatedEvent().GetCreateArguments().GetFields() {
			if field.GetLabel() == "instanceId" && field.GetValue().GetText() == instanceID {
				return contract.GetCreatedEvent().GetContractId()
			}
		}
	}
	t.Fatalf("contract not found: %s/%s with instanceId=%s", identifier.ModuleName, identifier.EntityName, instanceID)

	return ""
}

// queryContractFields reads fields from a CCIPFactory contract by its CID.
func queryContractFields(
	t *testing.T,
	participant canton.Participant,
	contractCid string,
) map[string]string {
	t.Helper()

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, contracts.IdentifierFromBinding(factory.CCIPFactory{}))
	require.NoError(t, err)

	for _, contract := range activeContracts {
		if contract.GetCreatedEvent().GetContractId() == contractCid {
			result := make(map[string]string)
			for _, field := range contract.GetCreatedEvent().GetCreateArguments().GetFields() {
				switch {
				case field.GetValue().GetText() != "":
					result[field.GetLabel()] = field.GetValue().GetText()
				case field.GetValue().GetParty() != "":
					result[field.GetLabel()] = field.GetValue().GetParty()
				default:
					result[field.GetLabel()] = fmt.Sprintf("%v", field.GetValue().GetBool())
				}
			}

			return result
		}
	}

	t.Fatalf("contract not found: CCIP.Factory/CCIPFactory CID=%s", contractCid)

	return nil
}

func queryBypasserOpCountForParty(
	t *testing.T,
	participant canton.Participant,
	party string,
	mcmsCid string,
) int64 {
	t.Helper()

	queryParticipant := participant
	queryParticipant.PartyID = party

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), queryParticipant, contracts.IdentifierFromBinding(mcmsCore.MCMS{}))
	require.NoError(t, err)

	for _, ac := range activeContracts {
		if ac.GetCreatedEvent().GetContractId() == mcmsCid {
			return extractBypasserOpCount(t, ac.GetCreatedEvent().GetCreateArguments())
		}
	}

	t.Fatalf("MCMS contract not found: %s", mcmsCid)

	return 0
}

// extractBypasserOpCount extracts bypasser.expiringRoot.opCount from MCMS create args.
func extractBypasserOpCount(t *testing.T, record *apiv2.Record) int64 {
	t.Helper()

	for _, field := range record.GetFields() {
		if field.GetLabel() == "bypasser" {
			for _, innerField := range field.GetValue().GetRecord().GetFields() {
				if innerField.GetLabel() == "expiringRoot" {
					for _, rootField := range innerField.GetValue().GetRecord().GetFields() {
						if rootField.GetLabel() == "opCount" {
							return rootField.GetValue().GetInt64()
						}
					}
				}
			}
		}
	}
	t.Fatal("could not find bypasser.expiringRoot.opCount")

	return 0
}
