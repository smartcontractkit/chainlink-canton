package factory

import (
	"context"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
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
