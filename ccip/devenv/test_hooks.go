package devenv

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// ClearFeeHoldingForTest clears the next fee holding CID so SendMessage fails without a fee input.
func (c *Chain) ClearFeeHoldingForTest() {
	c.nextFeeCID = ""
}

// SetNextFeeCIDForTest overrides the fee holding CID used on the next send.
func (c *Chain) SetNextFeeCIDForTest(cid string) {
	c.nextFeeCID = cid
}

// MintTokensReturningCID mints Amulet and returns the new holding contract ID.
func (c *Chain) MintTokensReturningCID(ctx context.Context, amount string) (string, error) {
	participant := c.chain.Participants[0]
	party := participant.PartyID

	validatorAPIClients, err := c.getValidatorAPIClients()
	if err != nil {
		return "", fmt.Errorf("get validator API clients: %w", err)
	}

	cid, err := testhelpers.MintAMT(
		ctx,
		participant,
		validatorAPIClients.metadataClient,
		validatorAPIClients.transferClient,
		validatorAPIClients.scanClient,
		party,
		amount,
	)
	if err != nil {
		return "", fmt.Errorf("mint tokens: %w", err)
	}

	return cid, nil
}

// SetFeeTokenInstrumentForTest overrides the fee token instrument used on the next send.
func (c *Chain) SetFeeTokenInstrumentForTest(inst splice_api_token_holding_v1.InstrumentId) {
	c.feeTokenInstrument = inst
}

// SetTransferTokenInstrumentForTest overrides the transfer token instrument used on the next send.
func (c *Chain) SetTransferTokenInstrumentForTest(inst splice_api_token_holding_v1.InstrumentId) {
	c.transferTokenInstrument = &inst
}
