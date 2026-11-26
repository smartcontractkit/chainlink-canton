package tests

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	bindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/feetreasury"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/core"
	splice "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	metadatav1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	transferv1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestFeeTreasury_WithdrawFeesEndToEndViaMCMS exercises the full two-phase design end to end
// against a real Splice (Amulet) instrument on localnet:
//
//	Phase 1 (MCMS-governed): the full MCMS batch (SetRoot -> ExecuteOp(ScheduleBatch) ->
//	  ExecuteScheduledBatch -> processCallsLoopInternal -> mcmsEntrypoint "AuthorizeFeeWithdrawal")
//	  creates a FeeWithdrawalAuthorization with the approved policy.
//	Phase 2 (outside MCMS): anyone submits ExecuteFeeWithdrawal with live Amulet holdings, the real
//	  TransferFactory, and disclosed mining-round context, moving fees from feeOwner to the recipient.
func TestFeeTreasury_WithdrawFeesEndToEndViaMCMS(t *testing.T) {
	t.Parallel()

	env := GetSharedFeeTreasuryEnvironment(t)
	participant := env.Participant
	mcmsEncoder := env.McmsEncoder
	ccipOwner := env.CcipOwner // also the feeOwner / MCMS owner
	cfg := env.Config
	sortedSigners := env.SortedSigners

	// Recipient is hosted on a second participant so its received Amulet balance is readable.
	recipientParticipant := env.RecipientParticipant
	recipient := env.Recipient

	// Splice validator API clients for the real Amulet token standard.
	scanProxyClient, tokenMetadataClient, transferInstructionClient, err := testhelpers.NewValidatorAPIClients(participant)
	require.NoError(t, err)

	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err)

	// Amulet is the native Splice fee instrument; its admin is the DSO/registry admin.
	instrumentID := splice.InstrumentId{Admin: types.PARTY(registryAdmin), Id: types.TEXT("Amulet")}

	chainID := int64(1)
	baseMcmsID := "mcms-feetreasury-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)
	mcmsCid := createMCMSMultiRole(t, participant, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	feeInstanceID := "fee-treasury-itest-" + uuid.New().String()[:8]
	feeInstanceAddr := fmt.Sprintf("%s@%s", feeInstanceID, ccipOwner)
	treasuryCid := createMCMSFeeTreasury(t, participant, ccipOwner, feeInstanceID)

	// Give the feeOwner Amulet fee holdings to later withdraw.
	feeHoldingCid, err := testhelpers.MintAMT(t.Context(), participant, tokenMetadataClient, transferInstructionClient, scanProxyClient, ccipOwner, "10000.0")
	require.NoError(t, err)

	const withdrawAmount = "100.0"

	// ---- Phase 1: authorize via the full MCMS batch ----
	feeContract := feetreasury.NewContract(fmt.Sprintf("#%s", feetreasury.PackageName), "CCIP.FeeTreasury", "MCMSFeeTreasury")
	authorizationID := "fee-withdrawal-itest-" + uuid.New().String()[:8]
	encoded, err := feeContract.Encoder().AuthorizeFeeWithdrawalParams(feetreasury.AuthorizeFeeWithdrawalParams{
		AuthorizationId: types.TEXT(authorizationID),
		Recipient:       types.PARTY(recipient),
		InstrumentId:    instrumentID,
		MaxAmount:       types.NUMERIC("500"),
		ValiditySecs:    types.INT64(3600),
	})
	require.NoError(t, err)

	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(feeInstanceAddr),
		FunctionName:          types.TEXT("AuthorizeFeeWithdrawal"),
		OperationData:         types.TEXT(encoded.OperationData),
	}}
	salt := uuid.New().String()[:8]

	scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, mcmsApi.ScheduleBatchParams{
		Calls:       calls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(salt),
		DelaySecs:   types.INT64(0),
	})

	proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
	proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
		AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
		Build()

	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	mcmsCid = setRootWithRole(t, participant, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

	mcmsCid, _, err = executeScheduledFeeWithdrawal(t, participant, ccipOwner, mcmsCid, opID, calls, salt, map[string]string{
		feeInstanceAddr: treasuryCid,
	}, nil)
	require.NoError(t, err)

	authCid, auth := queryFeeWithdrawalAuthorization(t, participant, ccipOwner, authorizationID)
	require.NotNil(t, auth, "FeeWithdrawalAuthorization should exist after MCMS execution")
	require.Equal(t, authorizationID, string(auth.InstanceId))
	require.Equal(t, recipient, string(auth.Recipient))
	require.Equal(t, instrumentID, auth.InstrumentId)

	// ---- Phase 2: execute the real withdrawal (outside MCMS) ----
	recipientBalanceBefore := holdingsBalance(t, recipientParticipant, &instrumentID, recipient)

	// Probe the registry for the real TransferFactory + disclosed mining-round context.
	tf, err := testhelpers.GetTransferFactoryV2(t.Context(), transferInstructionClient, registryAdmin, transferv1.Transfer{
		Sender:           types.PARTY(ccipOwner),
		Receiver:         types.PARTY(recipient),
		Amount:           types.NUMERIC(withdrawAmount),
		InstrumentId:     instrumentID,
		RequestedAt:      types.TIMESTAMP(time.Now().Add(-time.Hour)),
		ExecuteBefore:    types.TIMESTAMP(time.Now().Add(time.Hour)),
		InputHoldingCids: []types.CONTRACT_ID{types.CONTRACT_ID(feeHoldingCid)},
		Meta:             metadatav1.Metadata{Values: map[string]types.TEXT{}},
	})
	require.NoError(t, err)

	execRes := exerciseExecuteFeeWithdrawal(t, participant, ccipOwner, authCid, feeHoldingCid, withdrawAmount, tf)

	// Without a receiver preapproval the transfer yields a pending AmuletTransferInstruction
	// (plus a change Amulet back to the sender); the recipient must accept it to receive funds.
	var pendingTransferInstructionCid string
	for _, event := range execRes.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "AmuletTransferInstruction" {
				pendingTransferInstructionCid = e.Created.GetContractId()
			}
		}
	}
	require.NotEmpty(t, pendingTransferInstructionCid, "expected a pending AmuletTransferInstruction")

	time.Sleep(500 * time.Millisecond)
	require.NoError(t, testhelpers.AcceptPendingTransferInstruction(t.Context(), recipientParticipant, transferInstructionClient, recipient, pendingTransferInstructionCid))

	recipientBalanceAfter := holdingsBalance(t, recipientParticipant, &instrumentID, recipient)
	delta, _ := new(big.Float).SetRat(new(big.Rat).Sub(recipientBalanceAfter, recipientBalanceBefore)).Float64()
	require.InDelta(t, 100.0, delta, 0.01, "recipient should receive the withdrawn fee amount")

	// ---- Phase 3: recover a pending transfer whose recipient never acts ----
	feeOwnerBalanceBefore := holdingsBalance(t, participant, &instrumentID, ccipOwner)

	recoveryHoldingCid, err := testhelpers.MintAMT(t.Context(), participant, tokenMetadataClient, transferInstructionClient, scanProxyClient, ccipOwner, "150.0")
	require.NoError(t, err)

	const recoveryAmount = "100.0"
	const recoveryValiditySecs = 60
	shortAuthID := "fee-withdrawal-recovery-itest-" + uuid.New().String()[:8]
	shortEncoded, err := feeContract.Encoder().AuthorizeFeeWithdrawalParams(feetreasury.AuthorizeFeeWithdrawalParams{
		AuthorizationId: types.TEXT(shortAuthID),
		Recipient:       types.PARTY(recipient),
		InstrumentId:    instrumentID,
		MaxAmount:       types.NUMERIC(recoveryAmount),
		ValiditySecs:    types.INT64(recoveryValiditySecs),
	})
	require.NoError(t, err)

	shortCalls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(feeInstanceAddr),
		FunctionName:          types.TEXT("AuthorizeFeeWithdrawal"),
		OperationData:         types.TEXT(shortEncoded.OperationData),
	}}
	shortSalt := uuid.New().String()[:8]
	shortSchedule := MustEncodeScheduleBatch(t, mcmsEncoder, mcmsApi.ScheduleBatchParams{
		Calls:       shortCalls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(shortSalt),
		DelaySecs:   types.INT64(0),
	})
	shortProposal := NewMCMSProposal(int(chainID), proposerMultisigID, int(queryProposerOpCount(t, participant, ccipOwner, mcmsCid)), false).
		AddOperation(mcmsInstanceAddr, shortSchedule.Choice, shortSchedule.OperationData).
		Build()
	shortValidUntil := time.Now().Add(1 * time.Hour)
	shortSignatures, err := shortProposal.Sign(shortValidUntil, sortedSigners[:2])
	require.NoError(t, err)
	mcmsCid = setRootWithRole(t, participant, ccipOwner, mcmsCid, "Proposer", shortProposal, shortValidUntil, shortSignatures)

	shortOpID := HashTimelockOpId(UnwrapTimelockCalls(shortCalls), ZeroHash, shortSalt)
	shortOpProof, err := shortProposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, ccipOwner, mcmsCid, shortProposal.Operations[0], shortOpProof)

	var shortRes *apiv2.SubmitAndWaitForTransactionResponse
	mcmsCid, shortRes, err = executeScheduledFeeWithdrawal(t, participant, ccipOwner, mcmsCid, shortOpID, shortCalls, shortSalt, map[string]string{
		feeInstanceAddr: treasuryCid,
	}, nil)
	require.NoError(t, err)

	shortAuthCid, shortAuth := queryFeeWithdrawalAuthorization(t, participant, ccipOwner, shortAuthID)
	require.NotNil(t, shortAuth, "short-lived FeeWithdrawalAuthorization should exist after MCMS execution")

	recoveryTf, err := testhelpers.GetTransferFactoryV2(t.Context(), transferInstructionClient, registryAdmin, transferv1.Transfer{
		Sender:           types.PARTY(ccipOwner),
		Receiver:         types.PARTY(recipient),
		Amount:           types.NUMERIC(recoveryAmount),
		InstrumentId:     instrumentID,
		RequestedAt:      types.TIMESTAMP(time.Now().Add(-time.Hour)),
		ExecuteBefore:    types.TIMESTAMP(time.Now().Add(time.Duration(recoveryValiditySecs) * time.Second)),
		InputHoldingCids: []types.CONTRACT_ID{types.CONTRACT_ID(recoveryHoldingCid)},
		Meta:             metadatav1.Metadata{Values: map[string]types.TEXT{}},
	})
	require.NoError(t, err)

	recoveryExecRes := exerciseExecuteFeeWithdrawal(t, participant, ccipOwner, shortAuthCid, recoveryHoldingCid, recoveryAmount, recoveryTf)

	var pendingRecoveryInstructionCid string
	for _, event := range recoveryExecRes.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "AmuletTransferInstruction" {
				pendingRecoveryInstructionCid = e.Created.GetContractId()
			}
		}
	}
	require.NotEmpty(t, pendingRecoveryInstructionCid, "expected a pending AmuletTransferInstruction for the recovery leg")

	// The withdraw needs the instruction past its executeBefore, which is the
	// authorization's expiry: the authorize batch's ledger time plus validity.
	authorizeAt := shortRes.GetTransaction().GetEffectiveAt().AsTime()
	require.False(t, authorizeAt.IsZero(), "authorize batch should carry a ledger time")
	executeBefore := authorizeAt.Add(time.Duration(recoveryValiditySecs)*time.Second + 5*time.Second)
	if wait := time.Until(executeBefore); wait > 0 {
		time.Sleep(wait)
	}

	expireLock, withdrawDisclosures, err := testhelpers.GetWithdrawContext(t.Context(), transferInstructionClient, pendingRecoveryInstructionCid)
	require.NoError(t, err)
	require.True(t, expireLock, "the scan should report an expirable lock for a pending instruction")

	recoveryCid, recovery := queryFeeTransferRecoveryForInstruction(t, participant, ccipOwner, pendingRecoveryInstructionCid)
	require.NotNil(t, recovery, "FeeTransferRecovery should exist while the transfer is pending")
	recoveryAddr := fmt.Sprintf("%s@%s", string(recovery.InstanceId), ccipOwner)

	withdrawData := "00"
	if expireLock {
		withdrawData = "01"
	}
	withdrawCalls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(recoveryAddr),
		FunctionName:          types.TEXT("Withdraw"),
		OperationData:         types.TEXT(withdrawData),
	}}
	withdrawSalt := uuid.New().String()[:8]
	withdrawSchedule := MustEncodeScheduleBatch(t, mcmsEncoder, mcmsApi.ScheduleBatchParams{
		Calls:       withdrawCalls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(withdrawSalt),
		DelaySecs:   types.INT64(0),
	})
	withdrawProposal := NewMCMSProposal(int(chainID), proposerMultisigID, int(queryProposerOpCount(t, participant, ccipOwner, mcmsCid)), false).
		AddOperation(mcmsInstanceAddr, withdrawSchedule.Choice, withdrawSchedule.OperationData).
		Build()
	withdrawValidUntil := time.Now().Add(1 * time.Hour)
	withdrawSignatures, err := withdrawProposal.Sign(withdrawValidUntil, sortedSigners[:2])
	require.NoError(t, err)
	mcmsCid = setRootWithRole(t, participant, ccipOwner, mcmsCid, "Proposer", withdrawProposal, withdrawValidUntil, withdrawSignatures)

	withdrawOpID := HashTimelockOpId(UnwrapTimelockCalls(withdrawCalls), ZeroHash, withdrawSalt)
	withdrawOpProof, err := withdrawProposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, ccipOwner, mcmsCid, withdrawProposal.Operations[0], withdrawOpProof)

	// The ledger can give the execute an effective time that still trails the
	// deadline even after wall time has passed it, so retry the scheduled op.
	var withdrawRes *apiv2.SubmitAndWaitForTransactionResponse
	for attempt := 1; ; attempt++ {
		var execErr error
		mcmsCid, withdrawRes, execErr = executeScheduledFeeWithdrawal(t, participant, ccipOwner, mcmsCid, withdrawOpID, withdrawCalls, withdrawSalt, map[string]string{
			recoveryAddr: recoveryCid,
		}, withdrawDisclosures)
		if execErr == nil {
			break
		}
		if attempt >= 12 || !strings.Contains(execErr.Error(), "executeBefore has not passed") {
			require.NoError(t, execErr)
		}
		time.Sleep(10 * time.Second)
	}

	var instructionArchived bool
	for _, event := range withdrawRes.GetTransaction().GetEvents() {
		if a, ok := event.GetEvent().(*apiv2.Event_Archived); ok {
			if a.Archived.GetContractId() == pendingRecoveryInstructionCid {
				instructionArchived = true
			}
		}
	}
	require.True(t, instructionArchived, "the withdraw should archive the pending instruction")

	balanceAfter := holdingsBalance(t, participant, &instrumentID, ccipOwner)
	deltaAfter, _ := new(big.Float).SetRat(new(big.Rat).Sub(balanceAfter, feeOwnerBalanceBefore)).Float64()
	require.InDelta(t, 150.0, deltaAfter, 0.01, "feeOwner should keep the minted amount with nothing lost")

	// Phase 2's accepted transfer leaves its own recovery behind for governance
	// cleanup, so only this leg's recovery must be gone.
	remainingCid, remaining := queryFeeTransferRecoveryForInstruction(t, participant, ccipOwner, pendingRecoveryInstructionCid)
	require.Nil(t, remaining, "the recovery should be consumed by the withdraw, found %s", remainingCid)
}

func holdingsBalance(t *testing.T, participant canton.Participant, instrument *splice.InstrumentId, owner string) *big.Rat {
	t.Helper()
	bal, err := testhelpers.GetHoldingsBalance(t.Context(), participant, instrument, testhelpers.WithHoldingOwner(owner))
	require.NoError(t, err)

	return bal
}

func exerciseExecuteFeeWithdrawal(
	t *testing.T,
	participant canton.Participant,
	submitter string,
	authCid string,
	feeHoldingCid string,
	amount string,
	tf *testhelpers.TransferFactoryV2,
) *apiv2.SubmitAndWaitForTransactionResponse {
	t.Helper()

	choiceContext, err := testhelpers.ChoiceContextFromData(tf.ChoiceContextData)
	require.NoError(t, err)

	execArgs := feetreasury.ExecuteFeeWithdrawal{
		Submitter:          types.PARTY(submitter),
		TransferFactoryCid: types.CONTRACT_ID(tf.FactoryID),
		InputHoldingCids:   []types.CONTRACT_ID{types.CONTRACT_ID(feeHoldingCid)},
		Amount:             types.NUMERIC(amount),
		ExtraArgs: metadatav1.ExtraArgs{
			Context: metadatav1.ChoiceContext{Values: testhelpers.ExtractChoiceContextValues(choiceContext)},
			Meta:    metadatav1.Metadata{Values: map[string]types.TEXT{}},
		},
		RequestedAt: types.TIMESTAMP(time.Now().Add(-time.Minute)),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(feetreasury.FeeWithdrawalAuthorization{}),
						ContractId:     authCid,
						Choice:         "ExecuteFeeWithdrawal",
						ChoiceArgument: ledger.MapToValue(execArgs),
					},
				},
			}},
			ActAs:              []string{submitter},
			DisclosedContracts: tf.DisclosedContracts,
		},
	})
	require.NoError(t, err)

	return res
}

func createMCMSFeeTreasury(
	t *testing.T,
	participant canton.Participant,
	owner string,
	instanceID string,
) string {
	t.Helper()

	treasury := feetreasury.MCMSFeeTreasury{
		InstanceId:     types.TEXT(instanceID),
		FeeOwner:       types.PARTY(owner),
		McmsController: types.PARTY(owner),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId:      contracts.IdentifierFromBinding(feetreasury.MCMSFeeTreasury{}),
						CreateArguments: ledger.ConvertToRecord(treasury),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

func executeScheduledFeeWithdrawal(
	t *testing.T,
	participant canton.Participant,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcmsApi.TimelockCall,
	salt string,
	targetCids map[string]string,
	disclosedContracts []*apiv2.DisclosedContract,
) (string, *apiv2.SubmitAndWaitForTransactionResponse, error) {
	t.Helper()

	executeArgs := mcmsCore.ExecuteScheduledBatch{
		Submitter:   types.PARTY(owner),
		OpId:        types.TEXT(opID),
		Calls:       calls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(salt),
		TargetCids:  toContractIDMap(targetCids),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(mcmsCore.MCMS{}),
						ContractId:     mcmsCid,
						Choice:         "ExecuteScheduledBatch",
						ChoiceArgument: ledger.MapToValue(executeArgs),
					},
				},
			}},
			ActAs:              []string{owner},
			DisclosedContracts: disclosedContracts,
		},
	})
	if err != nil {
		return mcmsCid, nil, fmt.Errorf("execute scheduled batch: %w", err)
	}
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			return created.GetContractId(), res, nil
		}
	}

	return mcmsCid, res, fmt.Errorf("no MCMS contract created after ExecuteScheduledBatch")
}

func queryFeeWithdrawalAuthorization(
	t *testing.T,
	participant canton.Participant,
	owner string,
	instanceID string,
) (string, *feetreasury.FeeWithdrawalAuthorization) {
	t.Helper()

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, contracts.IdentifierFromBinding(feetreasury.FeeWithdrawalAuthorization{}))
	require.NoError(t, err)

	for _, c := range activeContracts {
		auth, err := bindings.UnmarshalCreatedEvent[feetreasury.FeeWithdrawalAuthorization](c.GetCreatedEvent())
		require.NoError(t, err)
		if string(auth.FeeOwner) == owner && string(auth.InstanceId) == instanceID {
			return c.GetCreatedEvent().GetContractId(), auth
		}
	}

	return "", nil
}

func queryFeeTransferRecoveryForInstruction(
	t *testing.T,
	participant canton.Participant,
	owner string,
	instructionCid string,
) (string, *feetreasury.FeeTransferRecovery) {
	t.Helper()

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, contracts.IdentifierFromBinding(feetreasury.FeeTransferRecovery{}))
	require.NoError(t, err)

	for _, c := range activeContracts {
		recovery, err := bindings.UnmarshalCreatedEvent[feetreasury.FeeTransferRecovery](c.GetCreatedEvent())
		require.NoError(t, err)
		if string(recovery.FeeOwner) == owner && string(recovery.TransferInstructionCid) == instructionCid {
			return c.GetCreatedEvent().GetContractId(), recovery
		}
	}

	return "", nil
}

func queryProposerOpCount(
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
		if ac.GetCreatedEvent().GetContractId() != mcmsCid {
			continue
		}
		for _, field := range ac.GetCreatedEvent().GetCreateArguments().GetFields() {
			if field.GetLabel() != "proposer" {
				continue
			}
			for _, innerField := range field.GetValue().GetRecord().GetFields() {
				if innerField.GetLabel() != "expiringRoot" {
					continue
				}
				for _, rootField := range innerField.GetValue().GetRecord().GetFields() {
					if rootField.GetLabel() == "opCount" {
						return rootField.GetValue().GetInt64()
					}
				}
			}
		}
	}

	t.Fatalf("MCMS contract not found or proposer op count missing: %s", mcmsCid)

	return 0
}
