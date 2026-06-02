package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcmscore "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/rmn"
	mcms_bindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestMCMS_ChangesetProposalE2E validates the complete pipeline:
//
//	changeset → MCMS proposal generation → on-chain MCMS signing/execution → state verification.
//
// Flow:
//  1. Deploy CCIP chain contracts (GlobalConfig, CCV, RMN, etc.)
//  2. Deploy MCMS contract with 2-of-3 signer config
//  3. Call ConfigureGlobalConfig changeset with TimelockConfig to produce an mcms.TimelockProposal
//  4. Validate proposal structure (metadata, operations, transactions, AdditionalFields)
//  5. Sign and execute the proposal on-chain via ScheduleBatch → ExecuteScheduledBatch
//  6. Verify the GlobalConfig dest chain config was applied
func TestMCMS_ChangesetProposalE2E(t *testing.T) {
	// if os.Getenv("INTEGRATION_TEST") == "" {
	//	t.Skip("Skipping integration test. Set INTEGRATION_TEST=1 to run.")
	// }

	t.Parallel()

	// ========================================================================
	// Phase 1: Environment Setup
	// ========================================================================

	sharedEnv := GetSharedCCIPMCMSEnvironment(t)
	participant := sharedEnv.Participant
	party := sharedEnv.CcipOwner
	config := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainSelector := chainsel.CANTON_LOCALNET.Selector
	chainSelectorNumeric := types.NUMERIC(strconv.FormatUint(chainSelector, 10))
	chainID := int64(1)

	// Build CLDF environment from the shared participant
	cantonChain := &canton.Chain{
		ChainMetadata: canton.ChainMetadata{Selector: chainSelector},
		Participants:  []canton.Participant{participant},
	}
	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())
	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{cantonChain}),
		OperationsBundle: bundle,
	}

	// Deploy CCIP chain contracts
	t.Log("Deploying CCIP chain contracts...")
	deployOut, err := cld_ops.ExecuteSequence(bundle, sequences.DeployChainContracts, *cantonChain, sequences.DeployChainContractsParams{
		CCIPOwnerParty: party,
		RMNOwnerParty:  party,
		CommitteeVerifiers: []sequences.CommitteeVerifierParams{{
			Template: ccvs.CommitteeVerifier{
				Owner:                        types.PARTY(party),
				CcipOwner:                    types.PARTY(party),
				VersionTag:                   "test-v1",
				StorageLocations:             []types.TEXT{"ipfs://test"},
				StorageLocationsAdmin:        types.PARTY(party),
				PendingStorageLocationsAdmin: types.PARTY(party),
			},
		}},
		GlobalConfig: sequences.GlobalConfigParams{
			Template: common.GlobalConfig{
				ChainSelector: chainSelectorNumeric,
			},
		},
		RMNRemote: sequences.RMNRemoteParams{
			Template: rmn.RMNRemote{
				RmnOwner: types.PARTY(party),
			},
		},
		NativeInstrumentId: splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(party),
			Id:    "LINK",
		},
	})
	require.NoError(t, err, "deploy chain contracts")

	// Extract GlobalConfig and CommitteeVerifier raw instance addresses from deploy output
	var gcRaw contracts.RawInstanceAddress
	var ccvRaw mcms_bindings.RawInstanceAddress
	for _, addr := range deployOut.Output.Addresses {
		raw, err := contracts.RawInstanceAddressFromString(addr.Labels.List()[0])
		require.NoError(t, err)
		switch string(addr.Type) {
		case "CantonGlobalConfig":
			gcRaw = raw
		case "CommitteeVerifier":
			ccvRaw = mcms_bindings.RawInstanceAddress{Unpack: types.TEXT(raw.String())}
		}
	}
	require.NotEmpty(t, gcRaw.String(), "GlobalConfig address not found in deploy output")
	require.NotEmpty(t, ccvRaw.Unpack, "CommitteeVerifier address not found in deploy output")
	t.Logf("GlobalConfig address: %s (hex: %s)", gcRaw.String(), gcRaw.InstanceAddress().Hex())

	// Deploy MCMS contract with 2-of-3 signer config
	baseMcmsID := "mcms-cs-e2e-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, party)
	mcmsRaw, err := contracts.RawInstanceAddressFromString(mcmsInstanceAddr)
	require.NoError(t, err)

	t.Log("Deploying MCMS contract with 2-of-3 config...")
	mcmsCid := createMCMSMultiRole(t, participant, party, chainID, baseMcmsID, config, 0, nil)
	t.Logf("MCMS contract created: %s", mcmsCid)

	// ========================================================================
	// Phase 2: Generate Proposal via ConfigureGlobalConfig Changeset
	// ========================================================================

	destChainSelector := "999"

	t.Log("Generating MCMS proposal via ConfigureGlobalConfig changeset...")
	csOut, err := changesets.ConfigureGlobalConfig{}.Apply(*env, changesets.CantonCSDeps[changesets.ConfigureGlobalConfigConfig]{
		ChainSelector: chainSelector,
		Participant:   0,
		Config: changesets.ConfigureGlobalConfigConfig{
			InstanceAddress:    gcRaw.InstanceAddress(),
			RawInstanceAddress: gcRaw.String(),
			DestChainUpdates: []common.DestChainConfigArgs{{
				DestChainSelector:         types.NUMERIC(destChainSelector),
				IsEnabled:                 true,
				AddressBytesLength:        20,
				TokenReceiverAllowed:      true,
				BaseExecutionGasCost:      21000,
				OffRampAddress:            "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
				DefaultCCVs:               []mcms_bindings.RawInstanceAddress{ccvRaw},
				MessageNetworkFeeUSDCents: "100",
				TokenNetworkFeeUSDCents:   "50",
			}},
			TimelockConfig: &changesets.MCMSTimelockConfig{
				MinDelay:         0,
				Description:      "E2E: configure global config via MCMS changeset",
				OverridePrevRoot: true,
				Action:           mcms_types.TimelockActionSchedule,
				MCMSContract: cantonmcms.MCMSContractInfo{
					RawInstanceAddress: mcmsRaw,
					InstanceAddress:    mcmsRaw.InstanceAddress(),
				},
				Role: cantonsdk.TimelockRoleProposer,
			},
		},
	})
	require.NoError(t, err, "ConfigureGlobalConfig changeset")
	require.Len(t, csOut.MCMSTimelockProposals, 1, "expected exactly one proposal")

	proposal := csOut.MCMSTimelockProposals[0]

	// ========================================================================
	// Phase 2b: Validate Proposal Structure
	// ========================================================================

	t.Log("Validating proposal structure...")
	assert.Equal(t, "v1", proposal.Version)
	assert.Equal(t, "E2E: configure global config via MCMS changeset", proposal.Description)
	assert.Equal(t, mcms_types.TimelockActionSchedule, proposal.Action)
	require.Len(t, proposal.Operations, 1)

	batchOp := proposal.Operations[0]
	assert.Equal(t, mcms_types.ChainSelector(chainSelector), batchOp.ChainSelector)
	require.Len(t, batchOp.Transactions, 1, "one dest chain update → one transaction")

	tx := batchOp.Transactions[0]
	assert.Equal(t, gcRaw.InstanceAddress().Hex(), tx.To)
	assert.NotEmpty(t, tx.Data, "encoded data should not be empty")

	var af cantonsdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
	assert.Equal(t, gcRaw.String(), af.TargetInstanceAddress)
	assert.Equal(t, "ApplyDestChainConfigUpdates", af.FunctionName)
	assert.NotEmpty(t, af.OperationData)
	t.Logf("Proposal validated: 1 batch, 1 tx targeting %s via %s", af.TargetInstanceAddress, af.FunctionName)

	// ========================================================================
	// Phase 3: Execute Proposal On-Chain via MCMS SDK
	// ========================================================================

	// Step 3a: Convert TimelockProposal -> MCMS Proposal
	t.Log("Converting TimelockProposal via SDK...")
	converter := cantonsdk.NewTimelockConverter()
	convertersMap := map[mcms_types.ChainSelector]sdk.TimelockConverter{
		mcms_types.ChainSelector(chainSelector): converter,
	}
	mcmsProposal, _, err := proposal.Convert(t.Context(), convertersMap)
	require.NoError(t, err, "convert timelock proposal")

	// Step 3b: Sign via SDK (2-of-3 quorum)
	t.Log("Signing proposal via SDK...")
	inspector := cantonsdk.NewInspector(participant.LedgerServices.State, []string{party}, cantonsdk.TimelockRoleProposer)
	inspectorsMap := map[mcms_types.ChainSelector]sdk.Inspector{
		mcms_types.ChainSelector(chainSelector): inspector,
	}
	signable, err := mcmscore.NewSignable(&mcmsProposal, inspectorsMap)
	require.NoError(t, err, "create signable")
	_, err = signable.SignAndAppend(mcmscore.NewPrivateKeySigner(sortedSigners[0].PrivateKey))
	require.NoError(t, err, "sign with signer 0")
	_, err = signable.SignAndAppend(mcmscore.NewPrivateKeySigner(sortedSigners[1].PrivateKey))
	require.NoError(t, err, "sign with signer 1")
	quorumMet, err := signable.ValidateSignatures(t.Context())
	require.NoError(t, err, "validate signatures")
	require.True(t, quorumMet, "quorum not met")

	// Step 3c: SetRoot + Execute (ScheduleBatch) via SDK
	t.Log("SetRoot + Execute via SDK...")
	encoders, err := mcmsProposal.GetEncoders()
	require.NoError(t, err, "get encoders")
	encoder := encoders[mcms_types.ChainSelector(chainSelector)].(*cantonsdk.Encoder)
	executor, err := cantonsdk.NewExecutor(
		encoder, inspector,
		participant.LedgerServices.Command, party, []string{party},
		cantonsdk.TimelockRoleProposer,
	)
	require.NoError(t, err, "create executor")
	executors := map[mcms_types.ChainSelector]sdk.Executor{
		mcms_types.ChainSelector(chainSelector): executor,
	}
	executable, err := mcmscore.NewExecutable(&mcmsProposal, executors)
	require.NoError(t, err, "create executable")

	_, err = executable.SetRoot(t.Context(), mcms_types.ChainSelector(chainSelector))
	require.NoError(t, err, "SetRoot")
	t.Log("SetRoot succeeded")

	for i := range mcmsProposal.Operations {
		_, err = executable.Execute(t.Context(), i)
		require.NoError(t, err, "execute operation %d", i)
	}
	t.Log("ScheduleBatch succeeded")

	// Step 3d: ExecuteScheduledBatch via TimelockExecutor
	// TargetCid resolution is now handled automatically by the SDK using TargetTemplateID.
	t.Log("ExecuteScheduledBatch via SDK...")
	timelockExecutor := cantonsdk.NewTimelockExecutor(
		participant.LedgerServices.Command, participant.LedgerServices.State, party, []string{party},
	)
	timelockExecutors := map[mcms_types.ChainSelector]sdk.TimelockExecutor{
		mcms_types.ChainSelector(chainSelector): timelockExecutor,
	}
	timelockExecutable, err := mcmscore.NewTimelockExecutable(t.Context(), &proposal, timelockExecutors)
	require.NoError(t, err, "create timelock executable")

	for i := range proposal.Operations {
		_, err = timelockExecutable.Execute(t.Context(), i)
		require.NoError(t, err, "timelock execute operation %d", i)
	}
	t.Log("ExecuteScheduledBatch succeeded")

	// ========================================================================
	// Phase 4: Verify GlobalConfig Was Updated
	// ========================================================================

	t.Log("Verifying GlobalConfig dest chain config...")

	// Re-resolve the GlobalConfig (it may have a new contract ID after the update)
	updatedGC, err := opcontract.FindActiveContractByInstanceAddress(
		t.Context(), participant.LedgerServices.State, []string{party},
		common.GlobalConfig{}.GetTemplateID(), gcRaw.InstanceAddress(),
	)
	require.NoError(t, err, "find updated GlobalConfig")

	// Parse destChainConfigs GenMap and validate individual field values
	var foundDestChainConfigs bool
	for _, field := range updatedGC.GetCreatedEvent().GetCreateArguments().GetFields() {
		if field.GetLabel() != "destChainConfigs" {
			continue
		}
		genMap := field.GetValue().GetGenMap()
		require.NotNil(t, genMap, "destChainConfigs should be a GenMap")
		require.Len(t, genMap.GetEntries(), 1, "expected exactly one dest chain config entry")

		entry := genMap.GetEntries()[0]
		assert.Equal(t, destChainSelector+".", entry.GetKey().GetNumeric(), "dest chain selector key (Daml Numeric includes trailing dot)")

		record := entry.GetValue().GetRecord()
		require.NotNil(t, record, "dest chain config value should be a Record")
		t.Logf("GlobalConfig destChainConfig record has %d fields", len(record.GetFields()))

		for _, f := range record.GetFields() {
			switch f.GetLabel() {
			case "isEnabled":
				assert.True(t, f.GetValue().GetBool(), "isEnabled")
			case "addressBytesLength":
				assert.Equal(t, int64(20), f.GetValue().GetInt64(), "addressBytesLength")
			case "tokenReceiverAllowed":
				assert.True(t, f.GetValue().GetBool(), "tokenReceiverAllowed")
			case "baseExecutionGasCost":
				assert.Equal(t, int64(21000), f.GetValue().GetInt64(), "baseExecutionGasCost")
			case "offRampAddress":
				assert.Equal(t, hex.EncodeToString([]byte("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef")), f.GetValue().GetText(), "offRampAddress (stored as hex-encoded ASCII)")
			case "messageNetworkFeeUSDCents":
				assert.Equal(t, "100.", f.GetValue().GetNumeric(), "messageNetworkFeeUSDCents")
			case "tokenNetworkFeeUSDCents":
				assert.Equal(t, "50.", f.GetValue().GetNumeric(), "tokenNetworkFeeUSDCents")
			}
		}
		foundDestChainConfigs = true
	}
	require.True(t, foundDestChainConfigs, "GlobalConfig should have destChainConfigs field")

	t.Log("=== TestMCMS_ChangesetProposalE2E PASSED ===")
	t.Log("Validated: changeset → proposal → SDK sign/execute → state update with correct field values")
}
