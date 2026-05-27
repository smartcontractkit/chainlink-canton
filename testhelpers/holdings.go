package testhelpers

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
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

// HoldingV1InterfaceID is the ledger Identifier for the Splice HoldingV1 interface (package name + module + entity).
func HoldingV1InterfaceID() *apiv2.Identifier {
	return &apiv2.Identifier{
		PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
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
	rows, err := ListHoldingsForInstrument(ctx, participant, instrument, filters...)
	if err != nil {
		return nil, err
	}

	total := new(big.Rat)
	for _, row := range rows {
		total.Add(total, row.Amount)
	}

	return total, nil
}

// SelectHoldingsForInstrument assigns len(minAmounts) different contract IDs from rows to output slots.
//
// Rows must already be limited to a single instrument and any caller-side filters (for example via
// ListHoldingsForInstrument with the desired instrument and Filter values).
//
// Greedy: sort rows by amount ascending, sort slots by required minimum
// descending, then for each slot choose the smallest unused row whose amount meets that minimum.
// Output index i corresponds to minAmounts[i]; a nil entry in minAmounts is treated as zero.
// TODO: should receive instrument as parameter to allow instrument-specific selection
func SelectHoldingsForInstrument(rows []ListedHolding, minAmounts []*big.Rat) ([]ListedHolding, error) {
	if len(minAmounts) == 0 {
		return nil, fmt.Errorf("minAmounts must be non-empty (each entry is one required holding)")
	}

	rows = slices.Clone(rows)
	slices.SortFunc(rows, func(a, b ListedHolding) int {
		return a.Amount.Cmp(b.Amount)
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

// ListHoldingsForInstrument loads active contracts that implement HoldingV1, parses each Holding view,
// applies instrument and Filter predicates, and returns ListedHolding rows with parseable amounts.
// instrument nil matches every instrument (same rule as GetHoldingsBalance).
func ListHoldingsForInstrument(
	ctx context.Context,
	participant canton.Participant,
	instrument *splice_api_token_holding_v1.InstrumentId,
	filters ...Filter,
) ([]ListedHolding, error) {
	contracts, err := ListActiveContractsByInterfaceId(ctx, participant, HoldingV1InterfaceID())
	if err != nil {
		return nil, err
	}

	var out []ListedHolding
	for _, ac := range contracts {
		var err error
		out, err = appendListedHoldingIfSelected(out, ac.GetCreatedEvent(), instrument, filters)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ListedHoldingsFromTransactionEventsForInstrument returns HoldingV1 rows parsed from Created events in a
// transaction that match inst and pass all filters (same rules as ListHoldingsForInstrument).
func ListedHoldingsFromTransactionEventsForInstrument(events []*apiv2.Event, inst splice_api_token_holding_v1.InstrumentId, filters ...Filter) ([]ListedHolding, error) {
	var out []ListedHolding
	for _, ev := range events {
		if ev == nil {
			continue
		}
		var err error
		out, err = appendListedHoldingIfSelected(out, ev.GetCreated(), &inst, filters)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// appendListedHoldingIfSelected parses HoldingV1 from created; when instrument is non-nil, only that
// instrument is kept. Appends to out when filters pass and amount parses; skips nil/empty creates.
func appendListedHoldingIfSelected(
	out []ListedHolding,
	created *apiv2.CreatedEvent,
	instrument *splice_api_token_holding_v1.InstrumentId,
	filters []Filter,
) ([]ListedHolding, error) {
	if created == nil || strings.TrimSpace(created.GetContractId()) == "" {
		return out, nil
	}
	contractID := strings.TrimSpace(created.GetContractId())

	hv, ok, err := parseHoldingV1View(created)
	if err != nil {
		return nil, err
	}
	if !ok {
		return out, nil
	}
	if instrument != nil && hv.InstrumentId != *instrument {
		return out, nil
	}
	if !holdingPassesFilters(contractID, hv, filters) {
		return out, nil
	}

	amountRaw := strings.TrimSpace(string(hv.Amount))
	if amountRaw == "" {
		return out, nil
	}
	amt, ok := new(big.Rat).SetString(amountRaw)
	if !ok {
		return nil, fmt.Errorf("invalid holding amount %q", amountRaw)
	}

	return append(out, ListedHolding{
		ContractID: contractID,
		View:       hv,
		Amount:     amt,
	}), nil
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
	want := HoldingV1InterfaceID()
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
