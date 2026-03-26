package tests

import (
	"context"
	"strconv"
	"strings"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

func quoteCCIPSenderFee(
	t *testing.T,
	participant canton.Participant,
	partySender string,
	ccipSenderCid string,
	sendArgs ccipsender.Send,
	disclosures []*apiv2.DisclosedContract,
) ccipsender.GetFeeResult {
	t.Helper()

	getFeeArgs := ccipsender.GetFee2{
		DestinationChainSelector: sendArgs.DestinationChainSelector,
		Message:                  sendArgs.Message,
		Ccvs:                     sendArgs.Ccvs,
		Context:                  sendArgs.Context,
		RouterCid:                sendArgs.RouterCid,
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					ContractId:     ccipSenderCid,
					Choice:         "GetFee",
					ChoiceArgument: ledger.MapToValue(getFeeArgs),
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: disclosures,
		},
		TransactionFormat: &apiv2.TransactionFormat{
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					partySender: {},
				},
			},
			// GetFee is a nonconsuming choice, and you want the exercised event with its return value
			TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
		},
	})
	require.NoError(t, err)

	var resultRecord *apiv2.Record
	for _, event := range res.GetTransaction().GetEvents() {
		exercised := event.GetExercised()
		if exercised == nil || exercised.GetChoice() != "GetFee" {
			continue
		}
		resultRecord = exercised.GetExerciseResult().GetRecord()
		break
	}
	require.NotNil(t, resultRecord, "GetFee should return an exercised record result")

	fields := resultRecord.GetFields()
	require.NotEmpty(t, fields, "GetFee should return fee fields")

	quote := ccipsender.GetFeeResult{
		FeeTokenAmount: types.NUMERIC(fields[0].GetValue().GetNumeric()),
	}
	if len(fields) > 1 {
		quote.PoolFeeTokenAmount = types.NUMERIC(fields[1].GetValue().GetNumeric())
	}
	return quote
}

func getHoldingsBalanceNumeric0(t *testing.T, ctx context.Context, participant canton.Participant) int64 {
	t.Helper()

	holdings, err := testhelpers.ListActiveContractsByInterfaceId(ctx, participant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)

	var total int64
	for _, h := range holdings {
		views := h.GetCreatedEvent().GetInterfaceViews()
		if len(views) == 0 {
			continue
		}
		fields := views[0].GetViewValue().GetFields()
		if len(fields) < 3 {
			continue
		}
		amountStr := fields[2].GetValue().GetNumeric()
		intPart, fracPart, hasFrac := strings.Cut(amountStr, ".")
		if hasFrac {
			require.Emptyf(t, strings.Trim(fracPart, "0"), "expected integer Numeric 0, got %q", amountStr)
		}
		amt, err := strconv.ParseInt(intPart, 10, 64)
		require.NoError(t, err)
		total += amt
	}

	return total
}
