//go:build !prodledger

package ledgertarget

import (
	latestholding "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	latestmeta "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
)

// AdaptChoiceContext maps caller-side choice context into ledger Send bindings.
func AdaptChoiceContext(ctx latestmeta.ChoiceContext) latestmeta.ChoiceContext {
	return ctx
}

// AdaptExtraArgs maps caller-side token extra args into ledger Send bindings.
func AdaptExtraArgs(args latestmeta.ExtraArgs) latestmeta.ExtraArgs {
	return args
}

// AdaptInstrumentId maps caller-side instrument IDs into ledger Send bindings.
func AdaptInstrumentId(id latestholding.InstrumentId) latestholding.InstrumentId {
	return id
}
