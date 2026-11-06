package tests

import (
	"fmt"
	"math/big"
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

	executeScheduledFeeWithdrawal(t, participant, ccipOwner, mcmsCid, opID, calls, salt, map[string]string{
		feeInstanceAddr: treasuryCid,
	})

	authCid, auth := queryFeeWithdrawalAuthorization(t, participant, ccipOwner)
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
) {
	t.Helper()

	executeArgs := mcmsCore.ExecuteScheduledBatch{
		Submitter:   types.PARTY(owner),
		OpId:        types.TEXT(opID),
		Calls:       calls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(salt),
		TargetCids:  toContractIDMap(targetCids),
	}

	_, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
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
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)
}

func queryFeeWithdrawalAuthorization(
	t *testing.T,
	participant canton.Participant,
	owner string,
) (string, *feetreasury.FeeWithdrawalAuthorization) {
	t.Helper()

	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, contracts.IdentifierFromBinding(feetreasury.FeeWithdrawalAuthorization{}))
	require.NoError(t, err)

	for _, c := range activeContracts {
		auth, err := bindings.UnmarshalCreatedEvent[feetreasury.FeeWithdrawalAuthorization](c.GetCreatedEvent())
		require.NoError(t, err)
		if string(auth.FeeOwner) == owner {
			return c.GetCreatedEvent().GetContractId(), auth
		}
	}

	return "", nil
}
