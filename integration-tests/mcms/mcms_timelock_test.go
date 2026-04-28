package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms/mcmstest"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func TestMCMS_Timelock(t *testing.T) {
	t.Parallel()

	// Use shared environment for participant, package, signers
	sharedEnv := GetSharedEnvironment(t)
	participant := sharedEnv.Participant
	mcmsPkgID := sharedEnv.McmsPkgID
	mcmsTestPkgID := sharedEnv.McmsTestPkgID
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-timelock-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	// Deploy MCMS with minDelay=0 for testing (self-dispatch tests will update it)
	mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	// Deploy Counter
	counterBaseID := "counter-timelock-" + uuid.New().String()[:8]
	counterInstanceAddr := fmt.Sprintf("%s@%s", counterBaseID, ccipOwner)
	counterCid := createCounter(t, participant, mcmsTestPkgID, ccipOwner, counterBaseID)
	counterTargetInstanceAddr := counterInstanceAddr

	t.Run("ScheduleAndExecute", func(t *testing.T) {
		t.Parallel()

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceAddress: types.TEXT(counterTargetInstanceAddr),
			FunctionName:          types.TEXT("Increment"),
			OperationData:         types.TEXT(""),
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
		scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

		// 1) Set proposer root authorizing schedule_batch with encoded operationData
		proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()

		validUntil := time.Now().Add(1 * time.Hour)
		signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)

		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

		// 2) Schedule a batch to increment counter (immediate, delay=0)
		opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)

		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)

		mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

		// 3) Execute immediately (delay is 0, so it's ready right away)
		mcmsCid = executeScheduledBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{counterTargetInstanceAddr: counterCid})

		// 4) Verify counter incremented
		val := queryCounterValue(t, participant, mcmsTestPkgID, counterBaseID)
		require.Equal(t, int64(1), val)
	})

	t.Run("CancelBatch", func(t *testing.T) {
		t.Parallel()

		// Fresh MCMS for this subtest (avoid cross-subtest opCount coupling)
		base := "mcms-timelock-cancel-" + uuid.New().String()[:8]
		mcmsInstanceAddr := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid2 := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 0, nil)

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceAddress: types.TEXT(counterTargetInstanceAddr),
			FunctionName:          types.TEXT("Increment"),
			OperationData:         types.TEXT(""),
		}}
		salt := uuid.New().String()[:8]
		opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)

		// Encode schedule params using encoder pattern
		scheduleParams := mcms.ScheduleBatchParams{
			Calls:       calls,
			Predecessor: types.TEXT(ZeroHash),
			Salt:        types.TEXT(salt),
			DelaySecs:   types.INT64(0),
		}
		scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

		// Counter is shared; ok.
		proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
		scheduleProposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
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
		cancelChoice := MustEncodeCancelBatch(t, mcmsEncoder, cancelParams)

		// Set canceller root authorizing cancel
		cancellerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleCanceller)
		cancelProposal := NewMCMSProposal(int(chainID), cancellerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, cancelChoice.Choice, cancelChoice.OperationData).
			Build()
		cancelSigs, err := cancelProposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid2 = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, "Canceller", cancelProposal, validUntil, cancelSigs)

		cancelProof, err := cancelProposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid2 = cancelBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, cancelProposal.Operations[0], cancelProof)

		// Should now be gone; ExecuteScheduledBatch should fail with not found.
		err = executeScheduledBatchExpectError(t, participant, mcmsPkgID, ccipOwner, mcmsCid2, opID, calls, ZeroHash, salt, map[string]string{counterTargetInstanceAddr: counterCid})
		require.Error(t, err)
		require.Contains(t, err.Error(), "mcms: operation not found in timelock", "expected 'mcms: operation not found in timelock', got: %v", err)
	})

	t.Run("BlockedFunction", func(t *testing.T) {
		t.Parallel()

		base := "mcms-timelock-blocked-" + uuid.New().String()[:8]
		mcmsInstanceAddr := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid3 := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 1_000_000, []mcms.BlockedFunction{
			{TargetInstanceAddress: types.TEXT(counterTargetInstanceAddr), FunctionName: types.TEXT("dangerous_function")},
		})

		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceAddress: types.TEXT(counterTargetInstanceAddr),
			FunctionName:          types.TEXT("dangerous_function"),
			OperationData:         types.TEXT(""),
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
		scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

		proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid3 = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid3, "Proposer", proposal, validUntil, sigs)

		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)

		err = scheduleBatchExpectError(t, participant, mcmsPkgID, ccipOwner, mcmsCid3, proposal.Operations[0], opProof)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mcms: function blocked:", "expected 'mcms: function blocked:', got: %v", err)
	})
}

func createSigners(t *testing.T) []*MCMSSigner {
	t.Helper()
	const count = 3
	signers := make([]*MCMSSigner, count)
	for i := range count {
		signer, err := NewMCMSSigner()
		require.NoError(t, err)
		signers[i] = signer
	}

	return signers
}

func createCounter(t *testing.T, participant canton.Participant, mcmsTestPkgID, owner, instanceID string) string {
	t.Helper()

	// Use bindings to create Counter - provides type safety and tests bindings work correctly
	counter := mcmstest.Counter{
		Owner:      types.PARTY(owner),
		InstanceId: types.TEXT(instanceID),
		Value:      types.INT64(0),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsTestPkgID,
							ModuleName: "MCMS.Mock.Counter",
							EntityName: "Counter",
						},
						CreateArguments: ledger.ConvertToRecord(counter),
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
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	chainID int64,
	baseMcmsID string,
	config MCMSConfig,
	minDelayMicros int64,
	blockedFunctions []mcms.BlockedFunction,
) string {
	t.Helper()
	mcmsContract := buildMCMSMultiRoleContract(owner, chainID, baseMcmsID, config, minDelayMicros, blockedFunctions)

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						CreateArguments: ledger.ConvertToRecord(mcmsContract),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

func buildMCMSMultiRoleContract(
	owner string,
	chainID int64,
	baseMcmsID string,
	config MCMSConfig,
	minDelayMicros int64,
	blockedFunctions []mcms.BlockedFunction,
) mcms.MCMS {
	signerInfos := make([]mcms.SignerInfo, len(config.Signers))
	for i, si := range config.Signers {
		signerInfos[i] = mcms.SignerInfo{
			SignerAddress: types.TEXT(si.SignerAddress),
			SignerIndex:   types.INT64(si.SignerIndex),
			SignerGroup:   types.INT64(si.SignerGroup),
		}
	}

	// Convert group quorums and parents
	groupQuorums := make([]types.INT64, NumGroups)
	groupParents := make([]types.INT64, NumGroups)
	for i := range NumGroups {
		groupQuorums[i] = types.INT64(config.GroupQuorums[i])
		groupParents[i] = types.INT64(config.GroupParents[i])
	}

	// Build MultisigConfig
	multisigConfig := mcms.MultisigConfig{
		Signers:      signerInfos,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
	}

	// Build empty RoleState (same for all roles initially)
	roleState := mcms.RoleState{
		Config:     multisigConfig,
		SeenHashes: map[types.TEXT]types.TIMESTAMP{},
		ExpiringRoot: mcms.ExpiringRoot{
			Root:       types.TEXT(""),
			ValidUntil: types.TIMESTAMP(time.Unix(0, 0)),
			OpCount:    types.INT64(0),
		},
		RootMetadata: mcms.RootMetadata{
			ChainId:              types.INT64(0),
			MultisigId:           types.TEXT(""),
			PreOpCount:           types.INT64(0),
			PostOpCount:          types.INT64(0),
			OverridePreviousRoot: types.BOOL(false),
		},
	}

	return mcms.MCMS{
		Owner:              types.PARTY(owner),
		InstanceId:         types.TEXT(baseMcmsID),
		ChainId:            types.INT64(chainID),
		Proposer:           roleState,
		Canceller:          roleState,
		Bypasser:           roleState,
		MinDelay:           types.RELTIME(time.Duration(minDelayMicros) * time.Microsecond),
		BlockedFunctions:   blockedFunctions,
		TimelockTimestamps: map[types.TEXT]types.TIMESTAMP{},
	}
}

func setRootWithRole(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	roleConstructor string,
	proposal *MCMSProposal,
	validUntil time.Time,
	signatures []RawSignature,
) string {
	t.Helper()

	return setRootWithRoleAndDisclosureParty(t, participant, mcmsPkgID, owner, owner, mcmsCid, roleConstructor, proposal, validUntil, signatures)
}

func setRootWithRoleAndDisclosureParty(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	submitter string,
	disclosureParty string,
	mcmsCid string,
	roleConstructor string,
	proposal *MCMSProposal,
	validUntil time.Time,
	signatures []RawSignature,
) string {
	t.Helper()

	// Use bindings for type safety - convert local types to binding types

	// Convert signatures to binding type
	bindingSignatures := make([]mcms.RawSignature, len(signatures))
	for i, sig := range signatures {
		bindingSignatures[i] = mcms.RawSignature{
			PublicKey: types.TEXT(sig.PublicKey),
			R:         types.TEXT(sig.R),
			S:         types.TEXT(sig.S),
		}
	}

	// Get metadata proof
	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	metadataProofTexts := make([]types.TEXT, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofTexts[i] = types.TEXT(p)
	}

	// Build SetRoot choice argument using bindings
	setRootArgs := mcms.SetRoot{
		TargetRole: mcms.Role(roleConstructor),
		Submitter:  types.PARTY(submitter),
		NewRoot:    types.TEXT(proposal.GetRoot()),
		ValidUntil: types.TIMESTAMP(validUntil),
		Metadata: mcms.RootMetadata{
			ChainId:              types.INT64(proposal.Metadata.ChainId),
			MultisigId:           types.TEXT(proposal.Metadata.MultisigId),
			PreOpCount:           types.INT64(proposal.Metadata.PreOpCount),
			PostOpCount:          types.INT64(proposal.Metadata.PostOpCount),
			OverridePreviousRoot: types.BOOL(proposal.Metadata.OverridePreviousRoot),
		},
		MetadataProof: metadataProofTexts,
		Signatures:    bindingSignatures,
	}

	disclosedContracts := mcmsDisclosureForSubmitter(t, participant, mcmsPkgID, submitter, disclosureParty, mcmsCid)

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "SetRoot",
						ChoiceArgument: ledger.MapToValue(setRootArgs),
					},
				},
			}},
			ActAs:              []string{submitter},
			DisclosedContracts: disclosedContracts,
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}

	return findRefreshedMCMSCid(t, participant, disclosureParty, mcmsPkgID, mcmsCid)
}

func scheduleBatch(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()

	// Use bindings for type safety
	executeOpArgs := mcms.ExecuteOp{
		TargetRole: mcms.RoleProposer,
		Submitter:  types.PARTY(owner),
		Op: mcms.Op{
			ChainId:               types.INT64(op.ChainId),
			MultisigId:            types.TEXT(op.MultisigId),
			Nonce:                 types.INT64(op.Nonce),
			TargetInstanceAddress: types.TEXT(op.TargetInstanceAddress),
			FunctionName:          types.TEXT(op.FunctionName),
			OperationData:         types.TEXT(op.OperationData),
		},
		OpProof:    toTextSlice(opProof),
		TargetCids: map[types.TEXT]types.CONTRACT_ID{},
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "ExecuteOp",
						ChoiceArgument: ledger.MapToValue(executeOpArgs),
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

// toTextSlice converts a []string to []types.TEXT for binding compatibility
func toTextSlice(s []string) []types.TEXT {
	result := make([]types.TEXT, len(s))
	for i, v := range s {
		result[i] = types.TEXT(v)
	}

	return result
}

func scheduleBatchExpectError(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	op MCMSOp,
	opProof []string,
) error {
	t.Helper()

	// Use bindings for type safety
	executeOpArgs := mcms.ExecuteOp{
		TargetRole: mcms.RoleProposer,
		Submitter:  types.PARTY(owner),
		Op: mcms.Op{
			ChainId:               types.INT64(op.ChainId),
			MultisigId:            types.TEXT(op.MultisigId),
			Nonce:                 types.INT64(op.Nonce),
			TargetInstanceAddress: types.TEXT(op.TargetInstanceAddress),
			FunctionName:          types.TEXT(op.FunctionName),
			OperationData:         types.TEXT(op.OperationData),
		},
		OpProof:    toTextSlice(opProof),
		TargetCids: map[types.TEXT]types.CONTRACT_ID{},
	}

	_, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "ExecuteOp",
						ChoiceArgument: ledger.MapToValue(executeOpArgs),
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
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcms.TimelockCall,
	predecessor string,
	salt string,
	targetCids map[string]string,
) string {
	t.Helper()

	// Use bindings for type safety - calls is already []mcms.TimelockCall
	executeArgs := mcms.ExecuteScheduledBatch{
		Submitter:   types.PARTY(owner),
		OpId:        types.TEXT(opID),
		Calls:       calls,
		Predecessor: types.TEXT(predecessor),
		Salt:        types.TEXT(salt),
		TargetCids:  toContractIDMap(targetCids),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "ExecuteScheduledBatch",
						ChoiceArgument: ledger.MapToValue(executeArgs),
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

// toContractIDMap converts a map[string]string to the typed TargetCids binding map.
// The key is the instanceAddress and the value is the contract ID wrapped as CONTRACT_ID type
// so it serializes as Value_ContractId rather than Value_Text.
func toContractIDMap(m map[string]string) map[types.TEXT]types.CONTRACT_ID {
	if m == nil {
		return map[types.TEXT]types.CONTRACT_ID{}
	}
	result := make(map[types.TEXT]types.CONTRACT_ID, len(m))
	for instanceID, cid := range m {
		result[types.TEXT(instanceID)] = types.CONTRACT_ID(cid)
	}

	return result
}

func executeScheduledBatchExpectError(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcms.TimelockCall,
	predecessor string,
	salt string,
	targetCids map[string]string,
) error {
	t.Helper()

	executeArgs := mcms.ExecuteScheduledBatch{
		Submitter:   types.PARTY(owner),
		OpId:        types.TEXT(opID),
		Calls:       calls,
		Predecessor: types.TEXT(predecessor),
		Salt:        types.TEXT(salt),
		TargetCids:  toContractIDMap(targetCids),
	}

	_, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "ExecuteScheduledBatch",
						ChoiceArgument: ledger.MapToValue(executeArgs),
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
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()

	executeOpArgs := mcms.ExecuteOp{
		TargetRole: mcms.RoleCanceller,
		Submitter:  types.PARTY(owner),
		Op: mcms.Op{
			ChainId:               types.INT64(op.ChainId),
			MultisigId:            types.TEXT(op.MultisigId),
			Nonce:                 types.INT64(op.Nonce),
			TargetInstanceAddress: types.TEXT(op.TargetInstanceAddress),
			FunctionName:          types.TEXT(op.FunctionName),
			OperationData:         types.TEXT(op.OperationData),
		},
		OpProof:    toTextSlice(opProof),
		TargetCids: map[types.TEXT]types.CONTRACT_ID{},
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "ExecuteOp",
						ChoiceArgument: ledger.MapToValue(executeOpArgs),
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

func queryCounterValue(t *testing.T, participant canton.Participant, mcmsTestPkgID, instanceID string) int64 {
	t.Helper()
	counterContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{
		PackageId:  mcmsTestPkgID,
		ModuleName: "MCMS.Mock.Counter",
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
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	targetCids map[string]string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()

	return bypasserExecuteBatchWithDisclosureParty(t, participant, mcmsPkgID, owner, owner, mcmsCid, targetCids, op, opProof)
}

func bypasserExecuteBatchWithDisclosureParty(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	submitter string,
	disclosureParty string,
	mcmsCid string,
	targetCids map[string]string,
	op MCMSOp,
	opProof []string,
) string {
	t.Helper()

	// BypasserExecuteBatch is dispatched via ExecuteOp with targetRole=Bypasser
	executeOpArgs := mcms.ExecuteOp{
		TargetRole: mcms.RoleBypasser,
		Submitter:  types.PARTY(submitter),
		Op: mcms.Op{
			ChainId:               types.INT64(op.ChainId),
			MultisigId:            types.TEXT(op.MultisigId),
			Nonce:                 types.INT64(op.Nonce),
			TargetInstanceAddress: types.TEXT(op.TargetInstanceAddress),
			FunctionName:          types.TEXT(op.FunctionName),
			OperationData:         types.TEXT(op.OperationData),
		},
		OpProof:    toTextSlice(opProof),
		TargetCids: toContractIDMap(targetCids),
	}

	disclosedContracts := mcmsDisclosureForSubmitter(t, participant, mcmsPkgID, submitter, disclosureParty, mcmsCid)

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "ExecuteOp",
						ChoiceArgument: ledger.MapToValue(executeOpArgs),
					},
				},
			}},
			ActAs:              []string{submitter},
			DisclosedContracts: disclosedContracts,
		},
	})
	require.NoError(t, err)
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId()
		}
	}

	return findRefreshedMCMSCid(t, participant, disclosureParty, mcmsPkgID, mcmsCid)
}

func mcmsDisclosureForSubmitter(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	submitter string,
	disclosureParty string,
	mcmsCid string,
) []*apiv2.DisclosedContract {
	t.Helper()

	if submitter == disclosureParty {
		return nil
	}

	queryParticipant := participant
	queryParticipant.PartyID = disclosureParty
	disclosedMcms, err := testhelpers.GetDisclosedContractById(t.Context(), queryParticipant, mcmsCid)
	require.NoError(t, err, "disclose MCMS contract %s from party %s", mcmsCid, disclosureParty)
	require.Equal(t, mcmsPkgID, disclosedMcms.GetTemplateId().GetPackageId())

	return []*apiv2.DisclosedContract{disclosedMcms}
}

func findRefreshedMCMSCid(
	t *testing.T,
	participant canton.Participant,
	party string,
	mcmsPkgID string,
	oldCid string,
) string {
	t.Helper()

	queryParticipant := participant
	queryParticipant.PartyID = party
	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), queryParticipant, &apiv2.Identifier{
		PackageId:  mcmsPkgID,
		ModuleName: "MCMS.Main",
		EntityName: "MCMS",
	})
	require.NoError(t, err)

	for i := len(activeContracts) - 1; i >= 0; i-- {
		cid := activeContracts[i].GetCreatedEvent().GetContractId()
		if cid != oldCid {
			return cid
		}
	}

	t.Fatalf("no refreshed MCMS contract found after consuming choice (old=%s)", oldCid)

	return ""
}

func queryMinDelay(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
) int64 {
	t.Helper()

	getMinDelayArgs := mcms.GetMinDelay{
		Submitter: types.PARTY(owner),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "GetMinDelay",
						ChoiceArgument: ledger.MapToValue(getMinDelayArgs),
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
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
) int64 {
	t.Helper()

	getBlockedFunctionsCountArgs := mcms.GetBlockedFunctionsCount{
		Submitter: types.PARTY(owner),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
						ContractId:     mcmsCid,
						Choice:         "GetBlockedFunctionsCount",
						ChoiceArgument: ledger.MapToValue(getBlockedFunctionsCountArgs),
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

	// Use shared environment for participant, package, signers
	sharedEnv := GetSharedEnvironment(t)
	participant := sharedEnv.Participant
	mcmsPkgID := sharedEnv.McmsPkgID
	mcmsTestPkgID := sharedEnv.McmsTestPkgID
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners
	chainID := int64(1)

	t.Run("SelfDispatchUpdateMinDelay", func(t *testing.T) {
		t.Parallel()

		base := "mcms-sd-delay-" + uuid.New().String()[:8]
		mcmsInstanceAddr := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 1_000_000, nil)

		// Define calls and salt first so we can encode them
		// Define calls and salt first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceAddress: types.TEXT(mcmsInstanceAddr),
			FunctionName:          types.TEXT("UpdateMinDelay"),
			OperationData:         types.TEXT(EncodeMinDelay(1)), // 1 second
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
		scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

		// 1) Set proposer root authorizing schedule_batch with encoded operationData
		proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, sigs)

		// 2) Schedule batch: self-dispatch update_min_delay to 1s
		opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
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
		mcmsInstanceAddr := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 1_000_000, nil)

		// Define calls and salt first so we can encode them
		bf := mcms.BlockedFunction{TargetInstanceAddress: "some-target@owner", FunctionName: "dangerous_op"}
		calls := []mcms.TimelockCall{{
			TargetInstanceAddress: types.TEXT(mcmsInstanceAddr),
			FunctionName:          types.TEXT("BlockFunction"),
			OperationData:         types.TEXT(MustEncodeBlockedFunction(t, bf)),
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
		scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

		proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, sigs)

		opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
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
		mcmsInstanceAddr := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 0, nil)

		// Define calls first so we can encode them
		// Define calls first so we can encode them
		calls := []mcms.TimelockCall{{
			TargetInstanceAddress: types.TEXT(mcmsInstanceAddr),
			FunctionName:          types.TEXT("UpdateMinDelay"),
			OperationData:         types.TEXT(EncodeMinDelay(2)), // 2 seconds
		}}

		// Encode bypasser execute params using encoder pattern
		bypasserParams := mcms.BypasserExecuteBatchParams{Calls: calls}
		bypasserChoice := MustEncodeBypasserExecuteBatch(t, mcmsEncoder, bypasserParams)

		// Set bypasser root with encoded operationData
		bypasserMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleBypasser)
		proposal := NewMCMSProposal(int(chainID), bypasserMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, bypasserChoice.Choice, bypasserChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Bypasser", proposal, validUntil, sigs)

		// Execute immediately (no delay)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid = bypasserExecuteBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, nil, proposal.Operations[0], opProof)

		delay := queryMinDelay(t, participant, mcmsPkgID, ccipOwner, mcmsCid)
		require.Equal(t, int64(2_000_000), delay, "minDelay should be 2s after bypasser self-dispatch")
	})

	t.Run("MixedBatch", func(t *testing.T) {
		t.Parallel()

		base := "mcms-sd-mixed-" + uuid.New().String()[:8]
		mcmsInstanceAddr := fmt.Sprintf("%s@%s", base, ccipOwner)
		mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, base, cfg, 0, nil)

		counterBase := "counter-mixed-" + uuid.New().String()[:8]
		counterCid := createCounter(t, participant, mcmsTestPkgID, ccipOwner, counterBase)
		counterTargetInstanceAddr := fmt.Sprintf("%s@%s", counterBase, ccipOwner)

		// Define mixed calls and salt first so we can encode them
		calls := []mcms.TimelockCall{
			{TargetInstanceAddress: types.TEXT(mcmsInstanceAddr), FunctionName: types.TEXT("UpdateMinDelay"), OperationData: types.TEXT(EncodeMinDelay(1))},
			{TargetInstanceAddress: types.TEXT(counterTargetInstanceAddr), FunctionName: types.TEXT("Increment"), OperationData: types.TEXT("")},
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
		scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

		proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
		proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
			AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
			Build()
		validUntil := time.Now().Add(1 * time.Hour)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)
		mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, sigs)

		opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

		time.Sleep(1500 * time.Millisecond)

		// targetCids only has 1 entry (for the external counter call)
		mcmsCid = executeScheduledBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{counterTargetInstanceAddr: counterCid})

		delay := queryMinDelay(t, participant, mcmsPkgID, ccipOwner, mcmsCid)
		require.Equal(t, int64(1_000_000), delay, "minDelay should be 1s after mixed batch")

		val := queryCounterValue(t, participant, mcmsTestPkgID, counterBase)
		require.Equal(t, int64(1), val, "counter should be incremented")
	})
}
