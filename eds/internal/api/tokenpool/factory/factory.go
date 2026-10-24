package factory

import (
	"context"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/eds/api/common"
)

type DisclosureFactory interface {
	GetSendDisclosures(ctx context.Context, message oapiCommon.Message) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error)
	GetExecuteDisclosures(
		ctx context.Context,
		message *protocol.Message,
		instrumentId splice_api_token_holding_v1.InstrumentId,
		inputHoldingCids []types.CONTRACT_ID,
		receiver types.PARTY,
	) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error)
}

// FeeTokenDisclosureFactory fetches the transfer factory disclosures for a fee token,
// used by the CCIP API to return the fee token's transfer factory disclosures to clients.
type FeeTokenDisclosureFactory interface {
	GetFeeTokenDisclosures(
		ctx context.Context,
		instrumentId splice_api_token_holding_v1.InstrumentId,
		sender types.PARTY,
		receiver types.PARTY,
	) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error)
}

// TransferDisclosureFactory is a DisclosureFactory backed by a token-standard TransferFactory,
// which can additionally resolve fee token disclosures.
type TransferDisclosureFactory interface {
	DisclosureFactory
	FeeTokenDisclosureFactory
}
