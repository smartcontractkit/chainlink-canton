package testhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"time"

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

// ListedHolding is one on-ledger row for a token holding: the contract id, the HoldingV1 view
// from the create event, and Amount parsed as *big.Rat for sums and comparisons.
type ListedHolding struct {
	ContractID string
	View       splice_api_token_holding_v1.HoldingView
	Amount     *big.Rat
}

// Filter decides whether a holding row is included. It receives the contract id and the
// HoldingV1 view parsed from the create event’s interface views.
type Filter func(contractID string, hv splice_api_token_holding_v1.HoldingView) bool

// WithHoldingOwner keeps only holdings whose view owner matches party (after trimming).
func WithHoldingOwner(party string) Filter {
	p := strings.TrimSpace(party)
	if p == "" {
		return func(string, splice_api_token_holding_v1.HoldingView) bool { return true }
	}

	return func(_ string, hv splice_api_token_holding_v1.HoldingView) bool {
		return string(hv.Owner) == p
	}
}

// WithUnlockedHoldingsOnly keeps only holdings with no Lock in the view (unlocked balance).
func WithUnlockedHoldingsOnly() Filter {
	return func(_ string, hv splice_api_token_holding_v1.HoldingView) bool {
		return hv.Lock == nil
	}
}

// ExcludeCIDs drops holdings whose contract id appears in cids.
func ExcludeCIDs(cids []string) Filter {
	return func(contractID string, _ splice_api_token_holding_v1.HoldingView) bool {
		return !slices.Contains(cids, contractID)
	}
}

// GetHoldingsBalance sums Numeric amounts for all holdings that pass filters.
// If instrument is non-nil, only that instrument is included; if nil, every instrument counts.
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

// PickHoldingsForInstrument chooses len(minAmounts) different holdings for the given instrument.
//
// Greedy: sort holdings by amount (largest first), sort slots by required minimum (largest first),
// then for each slot assign the biggest holding that is still unused and meets that minimum.
// Output index i matches minAmounts[i] (nil minimum counts as zero). filters narrow candidates.
//
// Typical use: pass one minAmount for one leg, or several for multiple legs on the same instrument.
// Use a separate call per instrument when you need holdings for different instruments.
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

// SplitHoldingSelfTransfer turns one holding into two by submitting a transfer to yourself for half
// the balance. The input contract is archived; the ledger creates two new Splice.Amulet contracts.
//
// On success the committed transaction must contain exactly two Amulet create events for this owner
// and instrument (DSO must match instrument.Admin). registryAdmin is the registry party used when
// building the transfer factory (see GetRegistryAdmin).
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
	if _, ok := new(big.Rat).SetString(halfStr); !ok {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("split: invalid half numeric %q", halfStr)
	}

	now := time.Now()
	transferParams := splice_api_token_transfer_instruction_v1.Transfer{
		Sender:        types.PARTY(owner),
		Receiver:      types.PARTY(owner),
		Amount:        types.NUMERIC(halfStr),
		InstrumentId:  holding.View.InstrumentId,
		RequestedAt:   types.TIMESTAMP(now.Add(-time.Hour)),
		ExecuteBefore: types.TIMESTAMP(now.Add(24 * time.Hour)),
		InputHoldingCids: []types.CONTRACT_ID{
			types.CONTRACT_ID(holding.ContractID),
		},
		Meta: splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	}

	tf, err := GetTransferFactoryV2(ctx, transferInstructionClient, registryAdmin, transferParams)
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
		Transfer:      transferParams,
		ExtraArgs:     extraArgs,
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

	children, err := splitChildrenFromSelfTransferTransaction(res.GetTransaction(), owner, holding.View.InstrumentId)
	if err != nil {
		return ListedHolding{}, ListedHolding{}, err
	}

	return children[0], children[1], nil
}

// splitChildrenFromSelfTransferTransaction reads the two new Amulet contracts produced by
// SplitHoldingSelfTransfer: same owner, instrument.Admin equal to Amulet dso, amounts from createArguments.
func splitChildrenFromSelfTransferTransaction(
	tx *apiv2.Transaction,
	owner string,
	instrument splice_api_token_holding_v1.InstrumentId,
) ([]ListedHolding, error) {
	if tx == nil {
		return nil, fmt.Errorf("split transfer: nil transaction")
	}

	wantAdmin := strings.TrimSpace(string(instrument.Admin))
	var out []ListedHolding
	for _, event := range tx.GetEvents() {
		created := event.GetCreated()
		if created == nil {
			continue
		}
		tid := created.GetTemplateId()
		if tid == nil || tid.GetModuleName() != "Splice.Amulet" || tid.GetEntityName() != "Amulet" {
			continue
		}
		cid := strings.TrimSpace(created.GetContractId())
		if cid == "" {
			continue
		}
		args := created.GetCreateArguments()
		if args == nil {
			continue
		}

		var dsoParty, ownerParty, initialAmount string
		for _, f := range args.GetFields() {
			switch f.GetLabel() {
			case "dso":
				dsoParty = strings.TrimSpace(f.GetValue().GetParty())
			case "owner":
				ownerParty = strings.TrimSpace(f.GetValue().GetParty())
			case "amount":
				expiring := f.GetValue().GetRecord()
				if expiring == nil {
					continue
				}
				for _, ef := range expiring.GetFields() {
					if ef.GetLabel() == "initialAmount" {
						initialAmount = strings.TrimSpace(ef.GetValue().GetNumeric())
						break
					}
				}
			}
		}
		if ownerParty != owner || dsoParty != wantAdmin {
			continue
		}

		amt, ok := new(big.Rat).SetString(initialAmount)
		if !ok {
			return nil, fmt.Errorf("split transfer: invalid initialAmount %q", initialAmount)
		}

		out = append(out, ListedHolding{
			ContractID: cid,
			View: splice_api_token_holding_v1.HoldingView{
				Owner:        types.PARTY(ownerParty),
				InstrumentId: instrument,
				Amount:       types.NUMERIC(initialAmount),
				Lock:         nil,
				Meta:         splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
			},
			Amount: new(big.Rat).Set(amt),
		})
	}
	if len(out) != 2 {
		return nil, fmt.Errorf("split transfer: expected exactly 2 Amulet child creates (got %d)", len(out))
	}
	return out, nil
}

// listHoldingsForInstrument loads active contracts that implement HoldingV1, parses each Holding view,
// applies instrument and Filter predicates, and returns ListedHolding rows with parseable amounts.
// instrument nil matches every instrument (same rule as GetHoldingsBalance).
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

		hv, ok, err := parseHoldingV1View(ac.GetCreatedEvent())
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

// holdingPassesFilters returns true only if every Filter accepts the row.
func holdingPassesFilters(contractID string, hv splice_api_token_holding_v1.HoldingView, filters []Filter) bool {
	for _, f := range filters {
		if !f(contractID, hv) {
			return false
		}
	}

	return true
}

// parseHoldingV1View finds the HoldingV1 interface view on a create event and decodes it into HoldingView.
// If there is no such view, it returns ok=false and a zero view. A malformed view record returns an error.
func parseHoldingV1View(created *apiv2.CreatedEvent) (splice_api_token_holding_v1.HoldingView, bool, error) {
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

// holdingV1InterfaceID is the ledger Identifier for the Splice HoldingV1 interface (package name + module + entity).
func holdingV1InterfaceID() *apiv2.Identifier {
	return &apiv2.Identifier{
		PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	}
}
