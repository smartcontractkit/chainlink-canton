package tests

import (
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestGlobalConfig_UpgradeV1ToV2 simulates a realistic Smart Contract Upgrade (SCU):
//
//  1. Upload mcms + globalconfig v1.0.0 (no migration capability)
//  2. Create GlobalConfigV1 and verify normal operations work
//  3. SCU: upload globalconfig v2.0.0 (adds MigrateToV2 + GlobalConfigV2)
//  4. Exercise MigrateToV2 via MCMSReceiver_Entrypoint (non-consuming — V1 stays alive)
//  5. Verify V1 and V2 coexist, use V2-specific features
//  6. Archive V1, confirm only V2 remains operational
func TestGlobalConfig_UpgradeV1ToV2(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
	participant := env.Chain.Participants[0]
	ccipOwner := participant.PartyID
	configInstanceID := "globalconfig-" + uuid.New().String()[:8]

	// ===================================================================
	// Phase 1: Deploy initial contracts (mcms + globalconfig v1.0.0)
	//   - V1 has NO MigrateToV2 choice at this point
	// ===================================================================

	mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err)
	gcV1Dar, err := contracts.GetDar(contracts.GlobalConfig, "1.0.0")
	require.NoError(t, err)

	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(
		t.Context(), [][]byte{mcmsDar, gcV1Dar}, participant,
	)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(packageIDs), 2)
	mcmsPkgID := packageIDs[0]
	gcV1PkgID := packageIDs[1]
	t.Logf("Uploaded mcms (pkg=%s) and globalconfig v1.0.0 (pkg=%s)", mcmsPkgID, gcV1PkgID)

	// ===================================================================
	// Phase 2: Create GlobalConfigV1 and verify basic operations
	// ===================================================================

	v1Cid := createGlobalConfigV1(t, participant, ccipOwner, configInstanceID)
	t.Logf("Created GlobalConfigV1: %s", v1Cid)

	instId := queryGlobalConfigInstanceId(t, participant, ccipOwner, v1Cid, "V1")
	require.Equal(t, configInstanceID, instId, "V1 instanceId should match")

	// Exercise UpdateDestChainConfig via MCMSReceiver_Entrypoint to prove V1 works
	instanceAddr := configInstanceID + "@" + ccipOwner
	v1Cid = exerciseMCMSReceiverEntrypoint(
		t, participant, ccipOwner, v1Cid,
		"UpdateDestChainConfig", "", map[string]string{instanceAddr: v1Cid},
	)
	t.Logf("V1 after UpdateDestChainConfig: %s", v1Cid)

	instId = queryGlobalConfigInstanceId(t, participant, ccipOwner, v1Cid, "V1")
	require.Equal(t, configInstanceID, instId, "V1 instanceId should still match after update")

	// ===================================================================
	// Phase 3: SCU Upgrade — upload globalconfig v2.0.0
	//   - Canton recognises this as an upgrade of v1.0.0
	//   - Existing V1 contract gains MigrateToV2 via upgraded code
	// ===================================================================

	gcV2Dar, err := contracts.GetDar(contracts.GlobalConfig, "2.0.0")
	require.NoError(t, err)

	v2PkgIDs, err := testhelpers.UploadDARstoMultipleParticipants(
		t.Context(), [][]byte{gcV2Dar}, participant,
	)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(v2PkgIDs), 1)
	gcV2PkgID := v2PkgIDs[0]
	t.Logf("SCU upgrade: uploaded globalconfig v2.0.0 (pkg=%s)", gcV2PkgID)

	// ===================================================================
	// Phase 4: Migrate V1 → V2 (non-consuming)
	//   - V1 stays alive; V2 is created alongside
	// ===================================================================

	v2Cid := exerciseMCMSReceiverEntrypoint(
		t, participant, ccipOwner, v1Cid,
		"MigrateToV2", "", map[string]string{instanceAddr: v1Cid},
	)
	t.Logf("MigrateToV2 created V2: %s", v2Cid)

	// V2 should have the same instanceId as V1
	instIdV2 := queryGlobalConfigInstanceId(t, participant, ccipOwner, v2Cid, "V2")
	require.Equal(t, configInstanceID, instIdV2, "V2 instanceId should match V1's")

	// V1 should still be alive (non-consuming migration)
	instIdV1 := queryGlobalConfigInstanceId(t, participant, ccipOwner, v1Cid, "V1")
	require.Equal(t, configInstanceID, instIdV1, "V1 should still be queryable after non-consuming migrate")
	t.Log("Confirmed V1 and V2 coexist")

	// ===================================================================
	// Phase 5: Use V2-specific features
	// ===================================================================

	v2Cid = exerciseMCMSReceiverEntrypoint(
		t, participant, ccipOwner, v2Cid,
		"SetNewFeature", "true", map[string]string{instanceAddr: v2Cid},
	)
	t.Logf("V2 after SetNewFeature: %s", v2Cid)

	featureEnabled := queryNewFeatureEnabled(t, participant, ccipOwner, v2Cid)
	require.True(t, featureEnabled, "newFeatureEnabled should be true after SetNewFeature")

	// ===================================================================
	// Phase 6: Archive V1 explicitly, verify only V2 remains
	// ===================================================================

	archiveGlobalConfig(t, participant, ccipOwner, v1Cid, "V1")
	t.Log("Archived GlobalConfigV1")

	// V1 should no longer be queryable
	_, err = tryQueryGlobalConfigInstanceId(t, participant, ccipOwner, v1Cid, "V1")
	require.Error(t, err)
	// Error message is CONTRACT_NOT_FOUND(11,e5da79f1): Contract could not be found with id <id>
	require.Contains(t, err.Error(), "CONTRACT_NOT_FOUND")
	t.Logf("Confirmed V1 is archived (error: %s)", err)

	// V2 should still be fully operational
	instIdV2 = queryGlobalConfigInstanceId(t, participant, ccipOwner, v2Cid, "V2")
	require.Equal(t, configInstanceID, instIdV2, "V2 should still be operational after V1 archival")
	t.Log("SCU upgrade test complete: V1 archived, V2 operational")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func createGlobalConfigV1(
	t *testing.T,
	participant canton.Participant,
	owner, instanceID string,
) string {
	t.Helper()

	emptyIntMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}

	createArgs := &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}}},
			{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
			{Label: "chainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 1}}},
			{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "aa"}}},
			{Label: "destChainConfigs", Value: emptyIntMap},
			{Label: "sourceChainConfigs", Value: emptyIntMap},
		},
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  "#globalconfig",
							ModuleName: "MCMS.GlobalConfig.V1",
							EntityName: "GlobalConfigV1",
						},
						CreateArguments: createArgs,
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

// exerciseMCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint interface
// choice. Because this is an interface choice, the TemplateId references the
// MCMSReceiver interface (not the implementing template).
func exerciseMCMSReceiverEntrypoint(
	t *testing.T,
	participant canton.Participant,
	owner, contractID string,
	functionName, operationData string,
	contractIds map[string]string,
) string {
	t.Helper()

	cidEntries := make([]*apiv2.GenMap_Entry, 0, len(contractIds))
	for instAddr, cid := range contractIds {
		cidEntries = append(cidEntries, &apiv2.GenMap_Entry{
			Key:   &apiv2.Value{Sum: &apiv2.Value_Text{Text: instAddr}},
			Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}},
		})
	}

	choiceArg := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: functionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: operationData}}},
			{Label: "contractIds", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: cidEntries}}}},
		},
	}}}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcms.PackageName,
							ModuleName: "MCMS.MCMSReceiver",
							EntityName: "MCMSReceiver",
						},
						ContractId:     contractID,
						Choice:         "MCMSReceiver_Entrypoint",
						ChoiceArgument: choiceArg,
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil {
			entityName := created.GetTemplateId().GetEntityName()
			if entityName == "GlobalConfigV1" || entityName == "GlobalConfigV2" {
				t.Logf("  MCMSReceiver_Entrypoint created %s: %s", entityName, created.GetContractId())
				return created.GetContractId()
			}
		}
	}
	t.Fatal("no GlobalConfig contract created after MCMSReceiver_Entrypoint")

	return ""
}

func queryGlobalConfigInstanceId(
	t *testing.T,
	participant canton.Participant,
	owner, contractID string,
	version string,
) string {
	t.Helper()

	result, err := tryQueryGlobalConfigInstanceId(t, participant, owner, contractID, version)
	require.NoError(t, err)

	return result
}

// tryQueryGlobalConfigInstanceId is the error-returning variant of queryGlobalConfigInstanceId.
// Used to verify that an archived contract can no longer be exercised.
func tryQueryGlobalConfigInstanceId(
	t *testing.T,
	participant canton.Participant,
	owner, contractID string,
	version string,
) (string, error) {
	t.Helper()

	var moduleName, entityName, choiceName string
	if version == "V1" {
		moduleName = "MCMS.GlobalConfig.V1"
		entityName = "GlobalConfigV1"
		choiceName = "GetInstanceIdV1"
	} else {
		moduleName = "MCMS.GlobalConfig.V2"
		entityName = "GlobalConfigV2"
		choiceName = "GetInstanceIdV2"
	}

	choiceArg := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "viewer", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
		},
	}}}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  "#globalconfig",
							ModuleName: moduleName,
							EntityName: entityName,
						},
						ContractId:     contractID,
						Choice:         choiceName,
						ChoiceArgument: choiceArg,
					},
				},
			}},
			ActAs: []string{owner},
		},
		TransactionFormat: &apiv2.TransactionFormat{
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					owner: {},
				},
			},
			TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
		},
	})
	if err != nil {
		return "", err
	}

	exerciseResult := res.GetTransaction().GetEvents()[0].GetExercised().GetExerciseResult()

	return exerciseResult.GetText(), nil
}

func queryNewFeatureEnabled(
	t *testing.T,
	participant canton.Participant,
	owner, contractID string,
) bool {
	t.Helper()

	choiceArg := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
		},
	}}}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  "#globalconfig",
							ModuleName: "MCMS.GlobalConfig.V2",
							EntityName: "GlobalConfigV2",
						},
						ContractId:     contractID,
						Choice:         "GetNewFeatureEnabled",
						ChoiceArgument: choiceArg,
					},
				},
			}},
			ActAs: []string{owner},
		},
		TransactionFormat: &apiv2.TransactionFormat{
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					owner: {},
				},
			},
			TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
		},
	})
	require.NoError(t, err)

	exerciseResult := res.GetTransaction().GetEvents()[0].GetExercised().GetExerciseResult()
	optionalValue := exerciseResult.GetOptional()
	if optionalValue == nil || optionalValue.GetValue() == nil {
		return false
	}

	return optionalValue.GetValue().GetBool()
}

// archiveGlobalConfig exercises the ArchiveGlobalConfigV1/V2 choice directly.
func archiveGlobalConfig(
	t *testing.T,
	participant canton.Participant,
	owner, contractID string,
	version string,
) {
	t.Helper()

	var moduleName, entityName, choiceName string
	if version == "V1" {
		moduleName = "MCMS.GlobalConfig.V1"
		entityName = "GlobalConfigV1"
		choiceName = "ArchiveGlobalConfigV1"
	} else {
		moduleName = "MCMS.GlobalConfig.V2"
		entityName = "GlobalConfigV2"
		choiceName = "ArchiveGlobalConfigV2"
	}

	choiceArg := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{},
	}}}

	_, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  "#globalconfig",
							ModuleName: moduleName,
							EntityName: entityName,
						},
						ContractId:     contractID,
						Choice:         choiceName,
						ChoiceArgument: choiceArg,
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
}
