package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

func TestMCMS_Timelock(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
	participant := env.Participant(1)

	// Upload DAR
	mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err)
	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar}, participant)
	require.NoError(t, err)
	require.NotEmpty(t, packageIDs)
	mcmsPkgID := packageIDs[0]

	// Create MCMS encoder for this package
	mcmsEncoder := NewMCMSEncoder(mcmsPkgID)

	ccipOwner := participant.Party

	// Shared signer set for all roles (tests can diverge later)
	signers := createSigners(t, 3)
	cfg := New2of3Config(signers)
	sortedSigners := SortSignersByAddress(signers)

	chainID := int64(1)
	baseMcmsID := "mcms-timelock-" + uuid.New().String()[:8]
	mcmsInstanceID := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	// Deploy MCMS with minDelay=0 for testing (self-dispatch tests will update it)
	mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	// Deploy Counter
	counterInstanceID := fmt.Sprintf("counter-timelock-%s@%s", uuid.New().String()[:8], ccipOwner)
	counterCid := createCounter(t, participant, mcmsPkgID, ccipOwner, counterInstanceID)
	counterTargetInstanceID := counterInstanceID

	t.Run("ScheduleAndExecute", func(t *testing.T) {
		t.Parallel()

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceId: types.TEXT(counterTargetInstanceID),
			FunctionName:     types.TEXT("Increment"),
			OperationData:    types.TEXT(""),
		}}
		salt := uuid.New().String()[:8]
		delaySecs := 0

		// Encode schedule params using encoder pattern (gets choice name + operation data)
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(delaySecs),
		}
		scheduleChoice := MustEncodeScheduleBatch(mcmsEncoder, scheduleParams)

		// 1) Set proposer root authorizing schedule_batch with encoded operationData
		proposerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()

		validUntil := time.Now().Add(1 * time.Hour)
		signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)

		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

		// 2) Schedule a batch to increment counter (immediate, delay=0)
		opID := HashTimelockOpId(FromMCMSTimelockCalls(calls), ZeroHash, salt)

		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)

		mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

		// 3) Execute immediately (delay is 0, so it's ready right away)
		mcmsCid = executeScheduledBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, []string{counterCid})

		// 4) Verify counter incremented
		val := queryCounterValue(t, participant, mcmsPkgID, counterInstanceID)
		require.Equal(t, int64(1), val)
	})

	t.Run("CancelBatch", func(t *testing.T) {
		t.Parallel()

		// Fresh MCMS for this subtest (avoid cross-subtest opCount coupling)
		base := "mcms-timelock-cancel-" + uuid.New().String()[:8]
		mcmsInstanceID := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid2 := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 0, nil)

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceId: types.TEXT(counterTargetInstanceID),
			FunctionName:     types.TEXT("Increment"),
			OperationData:    types.TEXT(""),
		}}
		salt := uuid.New().String()[:8]
		opID := HashTimelockOpId(FromMCMSTimelockCalls(calls), ZeroHash, salt)

		// Encode schedule params using encoder pattern
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(0),
		}
		scheduleChoice := MustEncodeScheduleBatch(mcmsEncoder, scheduleParams)

		// Counter is shared; ok.
		proposerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleProposer)
		scheduleProposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := scheduleProposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid2 = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, "Proposer", scheduleProposal, validUntil, sigs)

		opProof, err := scheduleProposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid2 = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, scheduleProposal.Operations[0], opProof)

		// Encode cancel params using encoder pattern
		cancelParams := mcms.CancelBatchParams{OpId: types.TEXT(opID)}
		cancelChoice := MustEncodeCancelBatch(mcmsEncoder, cancelParams)

		// Set canceller root authorizing cancel
		cancellerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleCanceller)
		cancelProposal := NewMCMSProposal(int(chainID), cancellerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, cancelChoice.Choice, cancelChoice.OperationData).
			Build()
		cancelSigs, err := cancelProposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid2 = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, "Canceller", cancelProposal, validUntil, cancelSigs)

		cancelProof, err := cancelProposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid2 = cancelBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, cancelProposal.Operations[0], cancelProof)

		// Should now be gone; ExecuteScheduledBatch should fail with not found.
		err = executeScheduledBatchExpectError(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, opID, calls, ZeroHash, salt, []string{counterCid})
		require.Error(t, err)
		require.Contains(t, err.Error(), "E_OPERATION_NOT_FOUND", "expected E_OPERATION_NOT_FOUND, got: %v", err)
	})

	t.Run("BlockedFunction", func(t *testing.T) {
		t.Parallel()

		base := "mcms-timelock-blocked-" + uuid.New().String()[:8]
		mcmsInstanceID := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid3 := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 1_000_000, []BlockedFunction{
			{TargetInstanceId: counterTargetInstanceID, FunctionName: "dangerous_function"},
		})

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceId: types.TEXT(counterTargetInstanceID),
			FunctionName:     types.TEXT("dangerous_function"),
			OperationData:    types.TEXT(""),
		}}
		salt := uuid.New().String()[:8]
		delaySecs := 1

		// Encode schedule params using encoder pattern
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(delaySecs),
		}
		scheduleChoice := MustEncodeScheduleBatch(mcmsEncoder, scheduleParams)

		proposerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid3 = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid3, "Proposer", proposal, validUntil, sigs)

		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)

		err = scheduleBatchExpectError(t, participant, mcmsPkgID, ccipOwner, mcmsCid3, proposal.Operations[0], opProof)
		require.Error(t, err)
		require.Contains(t, err.Error(), "E_FUNCTION_BLOCKED", "expected E_FUNCTION_BLOCKED, got: %v", err)
	})
}

func createSigners(t *testing.T, count int) []*MCMSSigner {
	t.Helper()
	signers := make([]*MCMSSigner, count)
	for i := range count {
		signer, err := NewMCMSSigner()
		require.NoError(t, err)
		signers[i] = signer
	}

	return signers
}

func createCounter(t *testing.T, participant testhelpers.Participant, mcmsPkgID, owner, instanceID string) string {
	t.Helper()
	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Counter",
							EntityName: "Counter",
						},
						CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
							{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}}},
							{Label: "value", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
						}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

func createMCMSMultiRole(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	chainID int64,
	baseMcmsID string,
	config MCMSConfig,
	minDelayMicros int64,
	blockedFunctions []BlockedFunction,
) string {
	t.Helper()
	instanceID := fmt.Sprintf("%s@%s", baseMcmsID, owner)

	// Build signer info values
	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
		},
	}}}

	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
		},
	}}}

	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}

	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: minDelayMicros}}},
		},
	}}}

	blockedFuncValues := make([]*apiv2.Value, 0, len(blockedFunctions))
	for _, bf := range blockedFunctions {
		blockedFuncValues = append(blockedFuncValues, &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: bf.TargetInstanceId}}},
				{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: bf.FunctionName}}},
			},
		}}})
	}
	blockedFunctionsValue := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: blockedFuncValues}}}

	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
							{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}}},
							{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainID}}},
							{Label: "proposer", Value: roleStateValue},
							{Label: "canceller", Value: roleStateValue},
							{Label: "bypasser", Value: roleStateValue},
							{Label: "minDelay", Value: minDelayValue},
							{Label: "blockedFunctions", Value: blockedFunctionsValue},
							{Label: "timelockTimestamps", Value: emptyMap},
						}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

func setRootWithRole(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	roleConstructor string,
	proposal *MCMSProposal,
	validUntil time.Time,
	signatures []RawSignature,
) string {
	t.Helper()

	signatureValues := make([]*apiv2.Value, len(signatures))
	for i, sig := range signatures {
		signatureValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.PublicKey}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.R}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.S}}},
			},
		}}}
	}

	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}
	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.Metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: proposal.Metadata.OverridePreviousRoot}}},
		},
	}}}

	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "SetRoot",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: roleConstructor}}}},
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.GetRoot()}}},
								{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntil.UnixMicro()}}},
								{Label: "metadata", Value: metadataValue},
								{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
								{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}
	t.Fatal("no MCMS contract created after SetRoot")

	return ""
}

func scheduleBatch(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()

	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}

	proofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		proofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "ExecuteOp",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "op", Value: opValue},
								{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: proofValues}}}},
								{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}
	t.Fatal("no MCMS contract created after ExecuteOp(schedule_batch)")

	return ""
}

func scheduleBatchExpectError(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	op MCMSOp,
	opProof []string,
) error {
	t.Helper()

	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}
	proofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		proofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}
	_, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "ExecuteOp",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "op", Value: opValue},
								{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: proofValues}}}},
								{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})

	return err
}

//nolint:unparam // predecessor is always ZeroHash in tests but kept for API consistency
func executeScheduledBatch(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcms.TimelockCall,
	predecessor string,
	salt string,
	targetCids []string,
) string {
	t.Helper()
	callValues := make([]*apiv2.Value, len(calls))
	for i, call := range calls {
		callValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.TargetInstanceId)}}},
				{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.FunctionName)}}},
				{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.OperationData)}}},
			},
		}}}
	}
	targetValues := make([]*apiv2.Value, len(targetCids))
	for i, cid := range targetCids {
		targetValues[i] = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}}
	}
	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "ExecuteScheduledBatch",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "opId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: opID}}},
								{Label: "calls", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: callValues}}}},
								{Label: "predecessor", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: predecessor}}},
								{Label: "salt", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: salt}}},
								{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: targetValues}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}
	t.Fatal("no MCMS contract created after ExecuteScheduledBatch")

	return ""
}

func executeScheduledBatchExpectError(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcms.TimelockCall,
	predecessor string,
	salt string,
	targetCids []string,
) error {
	t.Helper()
	callValues := make([]*apiv2.Value, len(calls))
	for i, call := range calls {
		callValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.TargetInstanceId)}}},
				{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.FunctionName)}}},
				{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.OperationData)}}},
			},
		}}}
	}
	targetValues := make([]*apiv2.Value, len(targetCids))
	for i, cid := range targetCids {
		targetValues[i] = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}}
	}
	_, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "ExecuteScheduledBatch",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "opId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: opID}}},
								{Label: "calls", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: callValues}}}},
								{Label: "predecessor", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: predecessor}}},
								{Label: "salt", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: salt}}},
								{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: targetValues}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})

	return err
}

func cancelBatch(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()
	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}
	proofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		proofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}
	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "ExecuteOp",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Canceller"}}}},
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "op", Value: opValue},
								{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: proofValues}}}},
								{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}
	t.Fatal("no MCMS contract created after ExecuteOp(cancel_batch)")

	return ""
}

func queryCounterValue(t *testing.T, participant testhelpers.Participant, mcmsPkgID, instanceID string) int64 {
	t.Helper()
	counterContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{
		PackageId:  mcmsPkgID,
		ModuleName: "MCMS.Counter",
		EntityName: "Counter",
	})
	require.NoError(t, err)
	for _, contract := range counterContracts {
		var inst string
		var val int64
		for _, field := range contract.GetCreatedEvent().GetCreateArguments().GetFields() {
			switch field.GetLabel() {
			case "instanceId":
				inst = field.GetValue().GetText()
			case "value":
				val = field.GetValue().GetInt64()
			}
		}
		if inst == instanceID {
			return val
		}
	}
	t.Fatalf("counter with instanceId %s not found", instanceID)

	return -1
}

func bypasserExecuteBatch(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	calls []mcms.TimelockCall,
	targetCids []string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()
	callValues := make([]*apiv2.Value, len(calls))
	for i, call := range calls {
		callValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.TargetInstanceId)}}},
				{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.FunctionName)}}},
				{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(call.OperationData)}}},
			},
		}}}
	}
	targetValues := make([]*apiv2.Value, len(targetCids))
	for i, cid := range targetCids {
		targetValues[i] = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: cid}}
	}
	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}
	proofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		proofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}
	// Use ExecuteOp with targetRole=Bypasser (BypasserExecuteBatch is dispatched via ExecuteOp)
	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "ExecuteOp",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Bypasser"}}}},
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "op", Value: opValue},
								{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: proofValues}}}},
								{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: targetValues}}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}
	t.Fatal("no MCMS contract created after ExecuteOp (Bypasser)")

	return ""
}

func queryMinDelay(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
) int64 {
	t.Helper()
	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "GetMinDelay",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
		// Use LEDGER_EFFECTS to include exercise events for non-consuming choices
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
	// GetMinDelay returns RelTime which is a record with microseconds field
	exerciseResult := res.GetTransaction().GetEvents()[0].GetExercised().GetExerciseResult()

	return exerciseResult.GetRecord().GetFields()[0].GetValue().GetInt64()
}

func queryBlockedFunctionsCount(
	t *testing.T,
	participant testhelpers.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
) int64 {
	t.Helper()
	res, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId: mcmsCid,
						Choice:     "GetBlockedFunctionsCount",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
							Fields: []*apiv2.RecordField{
								{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
							},
						}}},
					},
				},
			}},
			ActAs: []string{owner},
		},
		// Use LEDGER_EFFECTS to include exercise events for non-consuming choices
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

	return exerciseResult.GetInt64()
}

func TestMCMS_SelfDispatch(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
	participant := env.Participant(1)

	mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err)
	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar}, participant)
	require.NoError(t, err)
	require.NotEmpty(t, packageIDs)
	mcmsPkgID := packageIDs[0]

	// Create MCMS encoder for this package
	mcmsEncoder := NewMCMSEncoder(mcmsPkgID)

	ccipOwner := participant.Party

	signers := createSigners(t, 3)
	cfg := New2of3Config(signers)
	sortedSigners := SortSignersByAddress(signers)
	chainID := int64(1)

	t.Run("SelfDispatchUpdateMinDelay", func(t *testing.T) {
		t.Parallel()

		base := "mcms-sd-delay-" + uuid.New().String()[:8]
		mcmsInstanceID := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 1_000_000, nil)

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceId: types.TEXT(mcmsInstanceID),
			FunctionName:     types.TEXT("UpdateMinDelay"),
			OperationData:    types.TEXT("1"),
		}}
		salt := uuid.New().String()[:8]
		delaySecs := 1

		// Encode schedule params using encoder pattern
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(delaySecs),
		}
		scheduleChoice := MustEncodeScheduleBatch(mcmsEncoder, scheduleParams)

		// 1) Set proposer root authorizing schedule_batch with encoded operationData
		proposerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, sigs)

		// 2) Schedule batch: self-dispatch update_min_delay to 1s
		opID := HashTimelockOpId(FromMCMSTimelockCalls(calls), ZeroHash, salt)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

		// 3) Wait for delay
		time.Sleep(1500 * time.Millisecond)

		// 4) Execute - empty targetCids (all self-dispatch)
		mcmsCid = executeScheduledBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, nil)

		// 5) Verify minDelay changed to 1s (1_000_000 microseconds)
		delay := queryMinDelay(t, participant, mcmsPkgID, ccipOwner, mcmsCid)
		require.Equal(t, int64(1_000_000), delay, "minDelay should be 1s after self-dispatch")
	})

	t.Run("SelfDispatchBlockFunction", func(t *testing.T) {
		t.Parallel()

		base := "mcms-sd-block-" + uuid.New().String()[:8]
		mcmsInstanceID := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 1_000_000, nil)

		// Define calls and salt first so we can encode them
		bf := BlockedFunction{TargetInstanceId: "some-target@owner", FunctionName: "dangerous_op"}
		calls := []mcms.TimelockCall{{
			TargetInstanceId: types.TEXT(mcmsInstanceID),
			FunctionName:     types.TEXT("BlockFunction"),
			OperationData:    types.TEXT(EncodeBlockedFunction(bf)),
		}}
		salt := uuid.New().String()[:8]
		delaySecs := 1

		// Encode schedule params using encoder pattern
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(delaySecs),
		}
		scheduleChoice := MustEncodeScheduleBatch(mcmsEncoder, scheduleParams)

		proposerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, sigs)

		opID := HashTimelockOpId(FromMCMSTimelockCalls(calls), ZeroHash, salt)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

		time.Sleep(1500 * time.Millisecond)
		mcmsCid = executeScheduledBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, nil)

		count := queryBlockedFunctionsCount(t, participant, mcmsPkgID, ccipOwner, mcmsCid)
		require.Equal(t, int64(1), count, "should have 1 blocked function after self-dispatch")
	})

	t.Run("BypasserSelfDispatch", func(t *testing.T) {
		t.Parallel()

		base := "mcms-sd-bypass-" + uuid.New().String()[:8]
		mcmsInstanceID := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 0, nil)

		// Define calls first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceId: types.TEXT(mcmsInstanceID),
			FunctionName:     types.TEXT("UpdateMinDelay"),
			OperationData:    types.TEXT("2"),
		}}

		// Encode bypasser execute params using encoder pattern
		bypasserParams := mcms.BypasserExecuteBatchParams{Calls: calls}
		bypasserChoice := MustEncodeBypasserExecuteBatch(mcmsEncoder, bypasserParams)

		// Set bypasser root with encoded operationData
		bypasserMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleBypasser)
		proposal := NewMCMSProposal(int(chainID), bypasserMultisigID, 0, false).
			AddOperation(mcmsInstanceID, bypasserChoice.Choice, bypasserChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Bypasser", proposal, validUntil, sigs)

		// Execute immediately (no delay)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid = bypasserExecuteBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, calls, nil, proposal.Operations[0], opProof)

		delay := queryMinDelay(t, participant, mcmsPkgID, ccipOwner, mcmsCid)
		require.Equal(t, int64(2_000_000), delay, "minDelay should be 2s after bypasser self-dispatch")
	})

	t.Run("MixedBatch", func(t *testing.T) {
		t.Parallel()

		base := "mcms-sd-mixed-" + uuid.New().String()[:8]
		mcmsInstanceID := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 0, nil)

		counterInstanceID := fmt.Sprintf("counter-mixed-%s@%s", uuid.New().String()[:8], ccipOwner)
		counterCid := createCounter(t, participant, mcmsPkgID, ccipOwner, counterInstanceID)
		counterTargetInstanceID := counterInstanceID

		// Define mixed calls and salt first so we can encode them
		calls := []mcms.TimelockCall{
			{TargetInstanceId: types.TEXT(mcmsInstanceID), FunctionName: types.TEXT("UpdateMinDelay"), OperationData: types.TEXT("1")},
			{TargetInstanceId: types.TEXT(counterTargetInstanceID), FunctionName: types.TEXT("Increment"), OperationData: types.TEXT("")},
		}
		salt := uuid.New().String()[:8]
		delaySecs := 0

		// Encode schedule params using encoder pattern
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(delaySecs),
		}
		scheduleChoice := MustEncodeScheduleBatch(mcmsEncoder, scheduleParams)

		proposerMultisigID := MakeMcmsId(mcmsInstanceID, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceID, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, sigs)

		opID := HashTimelockOpId(FromMCMSTimelockCalls(calls), ZeroHash, salt)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

		time.Sleep(1500 * time.Millisecond)

		// targetCids only has 1 entry (for the external counter call)
		mcmsCid = executeScheduledBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, []string{counterCid})

		delay := queryMinDelay(t, participant, mcmsPkgID, ccipOwner, mcmsCid)
		require.Equal(t, int64(1_000_000), delay, "minDelay should be 1s after mixed batch")

		val := queryCounterValue(t, participant, mcmsPkgID, counterInstanceID)
		require.Equal(t, int64(1), val, "counter should be incremented")
	})
}
