package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	splice "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
	ceremonyledger "github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	ceremonyruntime "github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/runtime"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestCCIP_MCMSFactoryDeploy_DecentralizedParty validates the same factory
// deployment flow as TestCCIP_MCMSFactoryDeploy, except the MCMS contract is
// created by the party-ceremony contract-deploy flow with a decentralized party
// as MCMS owner/signatory
func TestCCIP_MCMSFactoryDeploy_DecentralizedParty(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(3))
	participants := env.Chain.Participants
	ceremonyParticipants := toCeremonyRuntimeParticipants(t, participants)
	participant := participants[0]
	ccipOwner := participant.PartyID

	synchronizerID, err := ceremonyruntime.DiscoverSynchronizerID(t.Context(), ceremonyParticipants[0])
	require.NoError(t, err)

	participantIDs := make([]string, len(participants))
	for i, p := range ceremonyParticipants {
		participantIDs[i], err = ceremonyruntime.ParticipantUID(t.Context(), p)
		require.NoError(t, err, "participant %d UID", i+1)
	}

	// --- Step 1: Onboard decentralized party via ceremony ---
	onboardingReporter := operations.NewMemoryReporter()
	decentralizedParty := onboardDecentralizedPartyForMCMSFactoryTest(t, ceremonyParticipants, participantIDs, synchronizerID, onboardingReporter)
	damlFingerprints := damlFingerprintsFromOnboardingReports(t, onboardingReporter)
	t.Logf("Decentralized party onboarded: %s", decentralizedParty)

	// --- Step 2: Upload DARs and prepare encoders ---
	uploadCCIPMCMSFactoryDARs(t, participants)
	mcmsEncoder := NewMCMSEncoder()
	factoryEncoder := factory.NewContract(fmt.Sprintf("#%s", factory.PackageName), "CCIP.FactoryV2", "CCIPFactory").Encoder()

	signers := createSigners(t)
	sortedSigners := SortSignersByAddress(signers)
	cfg := New2of3Config(signers)

	chainID := int64(1)
	uid := uuid.New().String()[:8]

	// --- Step 3: Deploy MCMS (2-of-3, minDelay=0) via decentralized party ceremony ---
	baseMcmsID := "mcms-dec-" + uid
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, decentralizedParty)
	mcmsCid := deployMCMSWithDecentralizedParty(t, ceremonyParticipants, participantIDs, damlFingerprints, decentralizedParty, synchronizerID, chainID, baseMcmsID, cfg)
	require.NotEmpty(t, mcmsCid)
	require.Equal(t, int64(0), queryBypasserOpCountForParty(t, participant, decentralizedParty, mcmsCid))
	t.Logf("MCMS created via ceremony: CID=%s, instanceAddr=%s", mcmsCid, mcmsInstanceAddr)

	// --- Step 4: Create CCIPFactory with owner=ccipOwner, mcmsParty=ccipOwner ---
	// Factory's mcmsParty is ccipOwner (not decentralizedParty) because factory ownership
	// stays with the local participant. The decentralized party governs only the MCMS contract,
	// which in turn controls the factory via BypasserExecuteBatch.
	factoryInstanceID := "factory-dec-" + uid
	factoryInstanceAddr := fmt.Sprintf("%s@%s", factoryInstanceID, ccipOwner)
	factoryCid := createCCIPFactory(t, participant, ccipOwner, factoryInstanceID)
	t.Logf("Factory created: CID=%s, instanceAddr=%s", factoryCid, factoryInstanceAddr)

	rmnInstanceID := "rmn-dec-" + uid
	rmnInstanceAddr := fmt.Sprintf("%s@%s", rmnInstanceID, ccipOwner)
	gcInstanceID := "gc-dec-" + uid
	gcInstanceAddr := fmt.Sprintf("%s@%s", gcInstanceID, ccipOwner)
	tarInstanceID := "tar-dec-" + uid
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	fqInstanceID := "fq-dec-" + uid
	fqInstanceAddr := fmt.Sprintf("%s@%s", fqInstanceID, ccipOwner)
	ccvInstanceID := "ccv-dec-" + uid
	ccvInstanceAddr := fmt.Sprintf("%s@%s", ccvInstanceID, ccipOwner)
	offRampInstanceID := "offramp-dec-" + uid
	offRampInstanceAddr := fmt.Sprintf("%s@%s", offRampInstanceID, ccipOwner)
	onRampInstanceID := "onramp-dec-" + uid
	onRampInstanceAddr := fmt.Sprintf("%s@%s", onRampInstanceID, ccipOwner)
	pprInstanceID := "ppr-dec-" + uid

	// --- Step 5: MCMS-driven deploys via BypasserExecuteBatch ---
	// Deploy order follows dependency graph (same as TestCCIP_MCMSFactoryDeploy).

	// 1. Deploy RMNRemote (no deps)
	rmnParams := factory.DeployRMNRemoteParams{
		InstanceId:      types.TEXT(rmnInstanceID),
		RmnOwner:        types.PARTY(ccipOwner),
		CcipOwner:       types.PARTY(ccipOwner),
		CustomObservers: []types.PARTY{},
		CursedSubjects:  []types.TEXT{},
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployRMNRemoteParams", rmnParams)
	t.Logf("RMNRemote deployed: %s", rmnInstanceID)

	// 2. Deploy GlobalConfig (no deps)
	gcParams := factory.DeployGlobalConfigParams{
		InstanceId:    types.TEXT(gcInstanceID),
		ChainSelector: types.NUMERIC("123"),
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployGlobalConfigParams", gcParams)
	t.Logf("GlobalConfig deployed: %s", gcInstanceID)

	// 3. Deploy TokenAdminRegistry (no deps)
	tarParams := factory.DeployTokenAdminRegistryParams{
		InstanceId: types.TEXT(tarInstanceID),
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployTokenAdminRegistryParams", tarParams)
	t.Logf("TokenAdminRegistry deployed: %s", tarInstanceID)

	// 4. Deploy FeeQuoter (no deps beyond link token)
	fqParams := factory.DeployFeeQuoterParams{
		InstanceId:            types.TEXT(fqInstanceID),
		LinkTokenInstrumentId: splice.InstrumentId{Admin: types.PARTY(ccipOwner), Id: "link-token"},
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployFeeQuoterParams", fqParams)
	t.Logf("FeeQuoter deployed: %s", fqInstanceID)

	// 5. Deploy CommitteeVerifier (deps: RMNRemote)
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
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployCommitteeVerifierParams", ccvParams)
	t.Logf("CommitteeVerifier deployed: %s", ccvInstanceID)

	// 6. Deploy OffRamp (deps: GlobalConfig, RMNRemote, TAR)
	offRampParams := factory.DeployOffRampParams{
		InstanceId:         types.TEXT(offRampInstanceID),
		GlobalConfig:       chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(gcInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployOffRampParams", offRampParams)
	t.Logf("OffRamp deployed: %s", offRampInstanceID)

	// 7. Deploy OnRamp (deps: GlobalConfig, RMNRemote, TAR, FeeQuoter, CCV)
	onRampParams := factory.DeployOnRampParams{
		InstanceId:         types.TEXT(onRampInstanceID),
		GlobalConfig:       chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(gcInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
		FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(fqInstanceAddr)},
		MaxUSDCentsPerMsg:  types.NUMERIC("100000"),
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployOnRampParams", onRampParams)
	t.Logf("OnRamp deployed: %s", onRampInstanceID)

	// 8. Deploy PerPartyRouterFactory (deps: all)
	pprParams := factory.DeployPerPartyRouterFactoryParams{
		InstanceId:         types.TEXT(pprInstanceID),
		OnRamp:             chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(onRampInstanceAddr)},
		OffRamp:            chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(offRampInstanceAddr)},
		GlobalConfig:       chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(gcInstanceAddr)},
		TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(tarInstanceAddr)},
		FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(fqInstanceAddr)},
		RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(rmnInstanceAddr)},
	}
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployPerPartyRouterFactoryParams", pprParams)
	t.Logf("PerPartyRouterFactory deployed: %s", pprInstanceID)

	// --- Step 6: Verify factory state ---
	factoryFields := queryContractFields(t, participant, factoryCid)
	require.Equal(t, factoryInstanceID, factoryFields["instanceId"])
	require.Equal(t, ccipOwner, factoryFields["owner"])
	require.Equal(t, "true", factoryFields["perPartyRouterFactoryDeployed"])

	t.Logf("All CCIP components deployed via MCMS. Factory CID=%s, MCMS CID=%s", factoryCid, mcmsCid)

	// --- Step 7 (Optional): Deploy LockReleaseTokenPool ---
	lrtpInstanceID := "lrtp-dec-" + uid
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
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployLockReleaseTokenPoolParams", lrtpParams)
	t.Logf("LockReleaseTokenPool deployed via MCMS: %s", lrtpInstanceID)

	// Deploy a second pool
	lrtp2InstanceID := "lrtp2-dec-" + uid
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
	factoryCid, mcmsCid = mcmsFactoryDeployWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, factoryCid, factoryInstanceAddr, chainID, sortedSigners, factoryEncoder, "DeployLockReleaseTokenPoolParams", lrtp2Params)
	t.Logf("LockReleaseTokenPool deployed via MCMS: %s", lrtp2InstanceID)

	// Final factory state check
	factoryFields = queryContractFields(t, participant, factoryCid)
	require.Equal(t, "true", factoryFields["perPartyRouterFactoryDeployed"])
	t.Logf("Final factory state verified. All deployments successful.")

	// --- Step 8: MCMS-driven config on deployed contracts ---
	gcCid := findNewContractCid(t, participant, contracts.IdentifierFromBinding(core.GlobalConfig{}), ccipOwner, gcInstanceAddr)
	require.NotEmpty(t, gcCid)

	remoteChainSelector := types.NUMERIC("456")
	gcEncoder := core.NewContract(fmt.Sprintf("#%s", core.PackageName), "CCIP.CoreV2.GlobalConfig", "GlobalConfig").Encoder()

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

	gcCid, mcmsCid = mcmsContractConfigWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, gcCid, gcInstanceAddr, chainID, sortedSigners, sourceEncoded)
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

	gcCid, _ = mcmsContractConfigWithMCMSQueryParty(t, participant, mcmsEncoder, ccipOwner, decentralizedParty, mcmsCid, mcmsInstanceAddr, gcCid, gcInstanceAddr, chainID, sortedSigners, destEncoded)
	require.NotEmpty(t, gcCid)
	t.Logf("GlobalConfig: ApplyDestChainConfigUpdates applied via MCMS")

	t.Logf("MCMS-driven config flow complete with decentralized party.")
}

func onboardDecentralizedPartyForMCMSFactoryTest(
	t *testing.T,
	participants []ceremonyruntime.Participant,
	participantIDs []string,
	synchronizerID string,
	reporter operations.Reporter,
) string {
	t.Helper()

	require.Len(t, participants, 3)
	require.Len(t, participantIDs, 3)

	namespaceName := "mcms-factory-" + uuid.New().String()[:8]
	input := onboarding.OnboardingInput{
		NamespaceName:  namespaceName,
		PartyPrefix:    "mcms-factory",
		Participants:   participantIDs,
		SynchronizerID: synchronizerID,
		Threshold:      3,
	}

	deps := make([]ceremony.CantonDeps, len(participants))
	for i, participant := range participants {
		dep, cleanup, err := ceremonyruntime.NewOnboardingDeps(t.Context(), participant, logger.Test(t), nil)
		require.NoError(t, err)
		t.Cleanup(func() { _ = cleanup() })
		deps[i] = dep
	}

	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), reporter)
	}

	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[0], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[1], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[2], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[0], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[1], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[2], input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())

	sr, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps[0], input)
	require.NoError(t, err)
	require.NotEmpty(t, sr.Output.PartyID)

	return sr.Output.PartyID
}

func deployMCMSWithDecentralizedParty(
	t *testing.T,
	participants []ceremonyruntime.Participant,
	participantIDs []string,
	damlFingerprints map[string]string,
	decentralizedParty string,
	synchronizerID string,
	chainID int64,
	baseMcmsID string,
	cfg MCMSConfig,
) string {
	t.Helper()

	input := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: decentralizedParty,
		SynchronizerID:       synchronizerID,
		Packages: []contractdeploy.PackageRef{{
			Name:    string(contracts.MCMSCore),
			Version: contracts.DevVersion,
		}},
		TemplateModule: "MCMS.Main",
		TemplateEntity: "MCMS",
		ContractArgs:   mcmsContractArgsJSON(t, decentralizedParty, chainID, baseMcmsID, cfg),
	}

	deps := make([]ceremonyledger.ContractDeployDeps, len(participants))
	for i, participant := range participants {
		fp, ok := damlFingerprints[participantIDs[i]]
		require.True(t, ok, "DAML fingerprint missing for participant %s", participantIDs[i])

		dep, cleanup, err := ceremonyruntime.NewContractDeployDeps(t.Context(), participant, embeddedDARLoader, fp, logger.Test(t), nil)
		require.NoError(t, err)
		t.Cleanup(func() { _ = cleanup() })
		deps[i] = dep
	}

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), reporter)
	}

	_, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[1], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[2], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[1], input)
	require.NoError(t, err)
	require.NotEmpty(t, sr.Output.ContractID)

	return sr.Output.ContractID
}

func mcmsContractArgsJSON(t *testing.T, owner string, chainID int64, baseMcmsID string, cfg MCMSConfig) string {
	t.Helper()

	mcmsContract := buildMCMSMultiRoleContract(owner, chainID, baseMcmsID, cfg, 0, nil)
	record := ledger.ConvertToRecord(mcmsContract)
	data, err := protojson.Marshal(record)
	require.NoError(t, err)

	return string(data)
}

func damlFingerprintsFromOnboardingReports(t *testing.T, reporter operations.Reporter) map[string]string {
	t.Helper()

	reports, err := reporter.GetReports()
	require.NoError(t, err)

	result := map[string]string{}
	for _, report := range reports {
		if report.Def.ID != "canton-ceremony/keys/create-member-key" {
			continue
		}

		if out, ok := report.Output.(keys.CreateMemberKeyOutput); ok {
			result[out.ParticipantID] = out.DamlKeyFingerprint

			continue
		}

		if out, ok := ceremony.Rehydrate[keys.CreateMemberKeyOutput](report.Output); ok {
			result[out.ParticipantID] = out.DamlKeyFingerprint
		}
	}

	return result
}

func uploadCCIPMCMSFactoryDARs(t *testing.T, participants []canton.Participant) {
	t.Helper()

	darPackages := []contracts.Package{
		contracts.MCMSCore,
		contracts.MCMSTest,
		contracts.CCIPRuntimeV2,
		contracts.CCIPCoreV2,
		contracts.CCIPCommitteeVerifierV2,
		contracts.CCIPExtensionAPIV2,
		contracts.CCIPLockReleaseTokenPoolV2,
		contracts.CCIPExecutorV2,
		contracts.CCIPFactoryV2,
	}

	darBytes := make([][]byte, 0, len(darPackages))
	for _, pkg := range darPackages {
		dar, err := contracts.GetDar(pkg, contracts.DevVersion)
		require.NoError(t, err)
		darBytes = append(darBytes, dar)
	}

	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), darBytes, participants...)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(packageIDs), len(darPackages))
}

func embeddedDARLoader(packageName, version string) ([]byte, error) {
	return contracts.GetDar(contracts.Package(packageName), version)
}

func toCeremonyRuntimeParticipants(t *testing.T, participants []canton.Participant) []ceremonyruntime.Participant {
	t.Helper()

	result := make([]ceremonyruntime.Participant, len(participants))
	for i, participant := range participants {
		runtimeParticipant, err := ceremonyruntime.FromCantonParticipant(participant)
		require.NoError(t, err, "runtime participant %s", participant.Name)
		result[i] = runtimeParticipant
	}

	return result
}
