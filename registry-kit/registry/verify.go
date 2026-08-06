package registry

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	registryholding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_holding_v0"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// VerifyHolding checks the Registry Holding template and HoldingV1 interface view.
// owner is the Holding owner; registrar is the instrument admin.
func VerifyHolding(ctx context.Context, client ledger.Client, owner, registrar, holdingCID, instrumentID, expectedAmount string) error {
	holdingTpl := contracts.IdentifierFromBinding(registryholding.Holding{})
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, client.ForParty(owner), holdingTpl)
	if err != nil {
		return fmt.Errorf("list holdings: %w", err)
	}

	var found *registryholding.Holding
	for _, ac := range active {
		if ac.GetCreatedEvent().GetContractId() != holdingCID {
			continue
		}
		found, err = bindings.UnmarshalCreatedEvent[registryholding.Holding](ac.GetCreatedEvent())
		if err != nil {
			return fmt.Errorf("unmarshal holding: %w", err)
		}
	}
	if found == nil {
		return fmt.Errorf("registry Holding %s not in ACS", holdingCID)
	}
	if types.PARTY(owner) != found.Owner {
		return fmt.Errorf("owner: expected %s got %s", owner, found.Owner)
	}
	expected, err := parseDecimal("expected amount", expectedAmount)
	if err != nil {
		return err
	}
	actual, err := parseDecimal("holding amount", string(found.Amount))
	if err != nil {
		return err
	}
	if !expected.Equal(actual) {
		return fmt.Errorf("amount: expected %s got %s", expectedAmount, found.Amount)
	}
	if types.PARTY(registrar) != found.Registrar {
		return fmt.Errorf("registrar: expected %s got %s", registrar, found.Registrar)
	}

	instrument := holdingsInstrumentQuery(registrar, instrumentID)
	partyParticipant := client.ForParty(owner)
	holdings, err := testhelpers.ListHoldingsForInstrument(ctx, partyParticipant, instrument, testhelpers.WithHoldingOwner(owner))
	if err != nil {
		return fmt.Errorf("list HoldingV1: %w", err)
	}

	var holdingV1 *testhelpers.ListedHolding
	for i := range holdings {
		if holdings[i].ContractID == holdingCID {
			holdingV1 = &holdings[i]
			break
		}
	}
	if holdingV1 == nil {
		return fmt.Errorf("HoldingV1 interface view for %s not found", holdingCID)
	}
	if string(holdingV1.View.InstrumentId.Admin) != registrar {
		return fmt.Errorf("HoldingV1 instrumentId.admin: expected %s got %s", registrar, holdingV1.View.InstrumentId.Admin)
	}
	holdingV1Amount, err := parseDecimal("HoldingV1 amount", string(holdingV1.View.Amount))
	if err != nil {
		return err
	}
	if !expected.Equal(holdingV1Amount) {
		return fmt.Errorf("HoldingV1 amount: expected %s got %s", expectedAmount, holdingV1.View.Amount)
	}

	balance, err := testhelpers.GetHoldingsBalance(ctx, partyParticipant, instrument, testhelpers.WithHoldingOwner(owner))
	if err != nil {
		return fmt.Errorf("aggregate balance: %w", err)
	}
	expectedRat := new(big.Rat)
	if _, ok := expectedRat.SetString(expectedAmount); !ok {
		return fmt.Errorf("aggregate balance: invalid expected amount %q", expectedAmount)
	}
	if expectedRat.Cmp(balance) != 0 {
		return fmt.Errorf("aggregate balance: expected %s got %s", expectedAmount, balance.FloatString(10))
	}

	return nil
}

func holdingsInstrumentQuery(registrar, instrumentID string) *splice_api_token_holding_v1.InstrumentId {
	return &splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registrar),
		Id:    types.TEXT(instrumentID),
	}
}
