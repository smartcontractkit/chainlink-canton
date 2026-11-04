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
	// GetSendDisclosures fetches the transfer factory and its choice context for a
	// pool's lock transfer. poolReceiveMeta carries the pool's poolReceiveMetadata so
	// the token-standard backend returns the right context mode: transfers with a
	// named receiver label (e.g. cip-112/receiver.id) get the two-step context,
	// "" only the direct/preapproved one.
	GetSendDisclosures(ctx context.Context, message oapiCommon.Message, poolReceiveMeta map[string]types.TEXT) (types.CONTRACT_ID, splice_api_token_metadata_v1.ChoiceContext, []oapiCommon.DisclosedContract, error)
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

// BridgeHoldingsDisclosureFactory resolves a pool owner's holdings via DA's
// explicit-disclosure bridge API, when the local holdings store cannot see them.
type BridgeHoldingsDisclosureFactory interface {
	GetBridgeHoldings(
		ctx context.Context,
		owner types.PARTY,
		instrumentId splice_api_token_holding_v1.InstrumentId,
	) ([]oapiCommon.DisclosedContract, error)
}
