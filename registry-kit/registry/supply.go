package registry

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	registryholding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_holding_v0"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// HoldingRow is one Registry Holding contract for an instrument.
type HoldingRow struct {
	ContractID string
	Owner      string
	Amount     decimal.Decimal
}

// QuerySupply aggregates Registry Holding template amounts for an instrument.
func QuerySupply(ctx context.Context, client ledger.Client, registrarParty, instrumentID string) (decimal.Decimal, []HoldingRow, error) {
	holdingTpl := contracts.IdentifierFromBinding(registryholding.Holding{})
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, client.ForParty(registrarParty), holdingTpl)
	if err != nil {
		return decimal.Zero, nil, fmt.Errorf("list holdings: %w", err)
	}

	total := decimal.Zero
	rows := make([]HoldingRow, 0)
	for _, ac := range active {
		holding, err := bindings.UnmarshalCreatedEvent[registryholding.Holding](ac.GetCreatedEvent())
		if err != nil {
			return decimal.Zero, nil, fmt.Errorf("unmarshal holding: %w", err)
		}
		if string(holding.Instrument.Id) != instrumentID {
			continue
		}
		amount, err := parseDecimal("holding amount", string(holding.Amount))
		if err != nil {
			return decimal.Zero, nil, err
		}
		total = total.Add(amount)
		rows = append(rows, HoldingRow{
			ContractID: ac.GetCreatedEvent().GetContractId(),
			Owner:      string(holding.Owner),
			Amount:     amount,
		})
	}

	return total, rows, nil
}
