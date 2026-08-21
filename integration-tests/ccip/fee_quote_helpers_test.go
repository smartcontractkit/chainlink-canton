package tests

import (
	"math/big"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/sender"
)

// getFeeChoiceArgumentMap builds the GetFee choice argument from Send: same encoding as
// Send.ToMap() for every field the GetFee choice accepts (GetFee has no feeTokenInput).
func getFeeChoiceArgumentMap(sendArgs sender.Send) map[string]any {
	m := sendArgs.ToMap()
	delete(m, "feeTokenInput")

	return m
}

// quoteCCIPSenderFee submits CCIPSender.GetFee with TRANSACTION_SHAPE_LEDGER_EFFECTS so the
// response includes the exercised return record. Callers that immediately Submit Send with the
// same sendArgs must re-fetch EDS disclosures and patch sendArgs + disclosed contracts before Send,
// or Canton may reject Send with LOCAL_VERDICT_INACTIVE_CONTRACTS.
func quoteCCIPSenderFee(
	t *testing.T,
	participant canton.Participant,
	partySender string,
	ccipSenderCid string,
	sendArgs sender.Send,
	disclosures []*apiv2.DisclosedContract,
) sender.GetFeeResult2 {
	t.Helper()

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     contracts.IdentifierFromBinding(sender.CCIPSender{}),
					ContractId:     ccipSenderCid,
					Choice:         "GetFee",
					ChoiceArgument: ledger.MapToValue(getFeeChoiceArgumentMap(sendArgs)),
				}},
			}},
			ActAs:              []string{partySender},
			DisclosedContracts: disclosures,
		},
		// Ledger-effects shape is required for exercised choice results to appear in the
		// transaction event stream; without it GetFee returns no matching Exercised events.
		TransactionFormat: &apiv2.TransactionFormat{
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					partySender: {},
				},
			},
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

	quote := sender.GetFeeResult2{
		FeeTokenAmount: types.NUMERIC(fields[0].GetValue().GetNumeric()),
	}
	if len(fields) > 1 {
		quote.PoolFeeTokenAmount = types.NUMERIC(fields[1].GetValue().GetNumeric())
	}

	return quote
}

// evmFeeQuoterDestChainConfigForLane returns the FeeQuoter dest-chain defaults that the EVM 2.0
// chain-family adapter used to expose through the lanes.LaneAdapter registry. It no longer
// registers there, because CCIP 2.0 lanes are configured through
// ConfigureChainsForLanesFromTopology. The Canton lane tests still build lanes.UpdateLanesInput
// by hand, so the values are pinned here.
func evmFeeQuoterDestChainConfigForLane() lanes.FeeQuoterDestChainConfig {
	// bytes4(keccak256("CCIP ChainFamilySelector EVM")) = 0x2812d52c
	const evmChainFamilySelector uint32 = 0x2812d52c
	return lanes.FeeQuoterDestChainConfig{
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         evmChainFamilySelector,
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		V2Params: &lanes.FeeQuoterV2Params{
			LinkFeeMultiplierPercent: 90,
			USDPerUnitGas:            big.NewInt(1e6),
		},
	}
}
