package testhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	splice_api_token_transfer_instruction_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

// ListedHolding is one active HoldingV1 contract that matched list filters, with parsed amount.
type ListedHolding struct {
	ContractID string
	View       splice_api_token_holding_v1.HoldingView
	Amount     *big.Rat
}

// Filter selects HoldingV1 rows using the ledger contract id (created event) and the interface view.
type Filter func(contractID string, hv splice_api_token_holding_v1.HoldingView) bool

// WithHoldingOwner restricts the sum to holdings whose interface view owner equals party.
func WithHoldingOwner(party string) Filter {
	p := strings.TrimSpace(party)
	if p == "" {
		return func(string, splice_api_token_holding_v1.HoldingView) bool { return true }
	}

	return func(_ string, hv splice_api_token_holding_v1.HoldingView) bool {
		return string(hv.Owner) == p
	}
}

// WithUnlockedHoldingsOnly skips holdings with a non-nil Lock in the interface view.
func WithUnlockedHoldingsOnly() Filter {
	return func(_ string, hv splice_api_token_holding_v1.HoldingView) bool {
		return hv.Lock == nil
	}
}

// WithUnlockedHoldingsOnly skips holdings with a non-nil Lock in the interface view.
func ExcludeCIDs(cids []string) Filter {
	return func(contractID string, _ splice_api_token_holding_v1.HoldingView) bool {
		return !slices.Contains(cids, contractID)
	}
}

// GetHoldingsBalance returns the sum of holding amounts in ledger Numeric units (*big.Rat).
// If instrument is nil, amounts for all instruments are summed (filters still apply).
func GetHoldingsBalance(
	ctx context.Context,
	participant canton.Participant,
	instrument *splice_api_token_holding_v1.InstrumentId,
	filters ...Filter,
) (*big.Rat, error) {
	rows, err := listHoldingsForInstrument(ctx, participant, instrument, filters...)
	if err != nil {
		return nil, err
	}

	total := new(big.Rat)
	for _, row := range rows {
		total.Add(total, row.Amount)
	}

	return total, nil
}

// PickHoldingsForInstrument returns len(minAmounts) distinct HoldingV1 rows for instrument.
// The i-th returned row corresponds to minAmounts[i]: each chosen holding has Amount >= minAmounts[i]
// (nil entries in minAmounts are treated as zero minimum).
//
// Callers choose how many holdings they need by the length of minAmounts (e.g. one entry for a single fee leg,
// two entries for two same-instrument legs). For another instrument, call again with that instrument’s ID and filters.
//
// filters should narrow the ACS row set (e.g. WithHoldingOwner(party), WithUnlockedHoldingsOnly(), WithHoldingContractID(cid)).
// Candidates are sorted by amount descending; slots with larger minimums are satisfied first, then each slot
// gets the largest still-unused holding that meets its minimum.
func PickHoldingsForInstrument(
	ctx context.Context,
	participant canton.Participant,
	instrument splice_api_token_holding_v1.InstrumentId,
	minAmounts []*big.Rat,
	filters ...Filter,
) ([]ListedHolding, error) {
	if len(minAmounts) == 0 {
		return nil, fmt.Errorf("minAmounts must be non-empty (each entry is one required holding)")
	}

	rows, err := listHoldingsForInstrument(ctx, participant, &instrument, filters...)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(rows, func(a, b ListedHolding) int {
		return b.Amount.Cmp(a.Amount)
	})

	type slot struct {
		idx int
		min *big.Rat
	}
	slots := make([]slot, len(minAmounts))
	for i := range minAmounts {
		m := minAmounts[i]
		if m == nil {
			m = big.NewRat(0, 1)
		}
		slots[i] = slot{idx: i, min: m}
	}
	slices.SortFunc(slots, func(a, b slot) int {
		return b.min.Cmp(a.min)
	})

	out := make([]ListedHolding, len(minAmounts))
	used := make(map[string]struct{})
	for _, s := range slots {
		var chosen ListedHolding
		var found bool
		for _, row := range rows {
			if _, ok := used[row.ContractID]; ok {
				continue
			}
			if row.Amount.Cmp(s.min) < 0 {
				continue
			}
			chosen, found = row, true
			break
		}
		if !found {
			return nil, fmt.Errorf(
				"could not pick holding for slot %d (min=%s): insufficient distinct holdings for instrument",
				s.idx, s.min.FloatString(8),
			)
		}
		used[chosen.ContractID] = struct{}{}
		out[s.idx] = chosen
	}

	return out, nil
}

// SplitHoldingSelfTransfer splits holding into two equal halves via a self-transfer (same sender and
// receiver). The input holding CID is consumed; two new HoldingV1 contracts with half the amount each
// are returned (first two ACS rows with that amount, excluding the original CID).
//
// registryAdmin is the instrument registry admin (see GetRegistryAdmin). If the factory workflow
// creates a pending transfer instruction, it is accepted as the holding owner.
func SplitHoldingSelfTransfer(
	ctx context.Context,
	participant canton.Participant,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	registryAdmin string,
	holding ListedHolding,
) (ListedHolding, ListedHolding, error) {
	if holding.Amount == nil || holding.Amount.Sign() <= 0 {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("holding amount must be positive")
	}

	owner := strings.TrimSpace(string(holding.View.Owner))
	if owner == "" {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("holding has empty owner")
	}

	half := new(big.Rat).Quo(new(big.Rat).Set(holding.Amount), big.NewRat(2, 1))
	halfStr := half.FloatString(10)
	halfWire, ok := new(big.Rat).SetString(halfStr)
	if !ok {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("split: invalid half numeric %q", halfStr)
	}

	tf, err := GetTransferFactoryV2(ctx, transferInstructionClient, registryAdmin, splice_api_token_transfer_instruction_v1.Transfer{
		Sender:       types.PARTY(owner),
		Receiver:     types.PARTY(owner),
		Amount:       types.NUMERIC(halfStr),
		InstrumentId: holding.View.InstrumentId,
		InputHoldingCids: []types.CONTRACT_ID{
			types.CONTRACT_ID(holding.ContractID),
		},
		Meta: splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	})
	if err != nil {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("get transfer factory: %w", err)
	}

	extraArgs := splice_api_token_metadata_v1.ExtraArgs{
		Context: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
		Meta:    splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	}
	if len(tf.ChoiceContextData) > 0 {
		raw, mErr := json.Marshal(tf.ChoiceContextData)
		if mErr != nil {
			return ListedHolding{}, ListedHolding{}, fmt.Errorf("marshal choice context: %w", mErr)
		}
		if uErr := json.Unmarshal(raw, &extraArgs.Context); uErr != nil {
			return ListedHolding{}, ListedHolding{}, fmt.Errorf("unmarshal choice context: %w", uErr)
		}
	}

	exercise := splice_api_token_transfer_instruction_v1.TransferFactoryTransfer{
		ExpectedAdmin: types.PARTY(registryAdmin),
		Transfer: splice_api_token_transfer_instruction_v1.Transfer{
			Sender:       types.PARTY(owner),
			Receiver:     types.PARTY(owner),
			Amount:       types.NUMERIC(halfStr),
			InstrumentId: holding.View.InstrumentId,
			InputHoldingCids: []types.CONTRACT_ID{
				types.CONTRACT_ID(holding.ContractID),
			},
			Meta: splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
		},
		ExtraArgs: extraArgs,
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  "#" + splice_api_token_transfer_instruction_v1.PackageName,
						ModuleName: "Splice.Api.Token.TransferInstructionV1",
						EntityName: "TransferFactory",
					},
					ContractId:     tf.FactoryID,
					Choice:         "TransferFactory_Transfer",
					ChoiceArgument: ledger.MapToValue(exercise),
				}},
			}},
			ActAs:              []string{owner},
			DisclosedContracts: tf.DisclosedContracts,
		},
	})
	if err != nil {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("submit transfer: %w", err)
	}

	var pendingTI string
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			name := e.Created.GetTemplateId().GetEntityName()
			if strings.Contains(name, "TransferInstruction") {
				pendingTI = e.Created.GetContractId()
			}
		}
	}
	if pendingTI != "" {
		if err := AcceptPendingTransferInstruction(ctx, participant, transferInstructionClient, owner, pendingTI); err != nil {
			return ListedHolding{}, ListedHolding{}, err
		}
	}

	rows, err := listHoldingsForInstrument(ctx, participant, &holding.View.InstrumentId, WithHoldingOwner(owner), WithUnlockedHoldingsOnly())
	if err != nil {
		return ListedHolding{}, ListedHolding{}, err
	}

	return firstTwoHoldingsAfterHalfSplit(rows, holding.ContractID, halfWire)
}

func firstTwoHoldingsAfterHalfSplit(rows []ListedHolding, excludeCID string, half *big.Rat) (ListedHolding, ListedHolding, error) {
	var first ListedHolding
	var foundFirst bool
	for _, row := range rows {
		if row.ContractID == excludeCID {
			continue
		}
		if row.Amount.Cmp(half) != 0 {
			continue
		}
		if !foundFirst {
			first, foundFirst = row, true
			continue
		}

		return first, row, nil
	}

	return ListedHolding{}, ListedHolding{}, fmt.Errorf("split: need 2 holdings at half amount %s (after excluding original cid)", half.FloatString(10))
}

// listHoldingsForInstrument walks active HoldingV1 contracts, applies optional instrument and filters,
// and returns rows with a non-empty contract id and a parseable amount (including zero).
// If instrument is nil, all instruments match (same semantics as GetHoldingsBalance).
func listHoldingsForInstrument(
	ctx context.Context,
	participant canton.Participant,
	instrument *splice_api_token_holding_v1.InstrumentId,
	filters ...Filter,
) ([]ListedHolding, error) {
	contracts, err := ListActiveContractsByInterfaceId(ctx, participant, holdingV1InterfaceID())
	if err != nil {
		return nil, err
	}

	var out []ListedHolding
	for _, ac := range contracts {
		created := ac.GetCreatedEvent()
		if created == nil || strings.TrimSpace(created.GetContractId()) == "" {
			continue
		}
		contractID := strings.TrimSpace(created.GetContractId())

		hv, ok, err := parseHoldingV1View(ac)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if instrument != nil && hv.InstrumentId != *instrument {
			continue
		}
		if !holdingPassesFilters(contractID, hv, filters) {
			continue
		}

		amountRaw := strings.TrimSpace(string(hv.Amount))
		if amountRaw == "" {
			continue
		}
		amt, ok := new(big.Rat).SetString(amountRaw)
		if !ok {
			return nil, fmt.Errorf("invalid holding amount %q", amountRaw)
		}

		out = append(out, ListedHolding{
			ContractID: contractID,
			View:       hv,
			Amount:     amt,
		})
	}

	return out, nil
}

func holdingPassesFilters(contractID string, hv splice_api_token_holding_v1.HoldingView, filters []Filter) bool {
	for _, f := range filters {
		if !f(contractID, hv) {
			return false
		}
	}

	return true
}

func parseHoldingV1View(ac *apiv2.ActiveContract) (splice_api_token_holding_v1.HoldingView, bool, error) {
	created := ac.GetCreatedEvent()
	if created == nil {
		return splice_api_token_holding_v1.HoldingView{}, false, nil
	}
	want := holdingV1InterfaceID()
	for _, iv := range created.GetInterfaceViews() {
		iface := iv.GetInterfaceId()
		if iface.GetModuleName() != want.GetModuleName() || iface.GetEntityName() != want.GetEntityName() {
			continue
		}
		vv := iv.GetViewValue()
		if vv == nil {
			continue
		}
		var hv splice_api_token_holding_v1.HoldingView
		if err := ledger.RecordToStruct(vv, &hv); err != nil {
			return splice_api_token_holding_v1.HoldingView{}, false, fmt.Errorf("parse HoldingV1 view: %w", err)
		}

		return hv, true, nil
	}

	return splice_api_token_holding_v1.HoldingView{}, false, nil
}

// holdingV1InterfaceID is the canonical HoldingV1 ledger interface id (single place for module/entity strings).
func holdingV1InterfaceID() *apiv2.Identifier {
	return &apiv2.Identifier{
		PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	}
}
