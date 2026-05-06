package testhelpers

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
)

type Filter func(splice_api_token_holding_v1.HoldingView) bool

// WithHoldingOwner restricts the sum to holdings whose interface view owner equals party.
func WithHoldingOwner(party string) Filter {
	p := strings.TrimSpace(party)
	if p == "" {
		return func(splice_api_token_holding_v1.HoldingView) bool { return true }
	}
	return func(hv splice_api_token_holding_v1.HoldingView) bool {
		return string(hv.Owner) == p
	}
}

// WithUnlockedHoldingsOnly skips holdings with a non-nil Lock in the interface view.
func WithUnlockedHoldingsOnly() Filter {
	return func(hv splice_api_token_holding_v1.HoldingView) bool {
		return hv.Lock == nil
	}
}

// GetHoldingsBalance lists HoldingV1 contracts for the participant's party and returns the sum of amounts
// for instrument (ledger Numeric as *big.Rat). Optionally filters by owner and/or unlocked holdings.
func GetHoldingsBalance(
	ctx context.Context,
	participant canton.Participant,
	instrument splice_api_token_holding_v1.InstrumentId,
	filters ...Filter,
) (*big.Rat, error) {
	contracts, err := ListActiveContractsByInterfaceId(ctx, participant, holdingV1InterfaceID())
	if err != nil {
		return nil, err
	}

	total := new(big.Rat)
	for _, ac := range contracts {
		hv, ok, err := parseHoldingV1View(ac)
		if err != nil {
			return nil, err
		}
		if !ok || hv.InstrumentId != instrument {
			continue
		}
		if !holdingPassesFilters(hv, filters) {
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
		total.Add(total, amt)
	}

	return total, nil
}

func holdingPassesFilters(hv splice_api_token_holding_v1.HoldingView, filters []Filter) bool {
	for _, f := range filters {
		if !f(hv) {
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
