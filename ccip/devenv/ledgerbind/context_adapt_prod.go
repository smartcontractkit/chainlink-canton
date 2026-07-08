//go:build prodledger

package ledgerbind

import (
	"encoding/json"
	"fmt"

	latestholding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	latestmeta "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	prodholding "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_holding_v1"
	prodmeta "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_metadata_v1"
)

func adaptViaJSON[T any](in any) (T, error) {
	var out T
	data, err := json.Marshal(in)
	if err != nil {
		return out, fmt.Errorf("marshal for ledger binding adapt: %w", err)
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("unmarshal for ledger binding adapt: %w", err)
	}

	return out, nil
}

func mustAdapt[T any](in any) T {
	out, err := adaptViaJSON[T](in)
	if err != nil {
		panic(err)
	}

	return out
}

// AdaptChoiceContext maps caller-side choice context into v1_0_0 ledger Send bindings.
func AdaptChoiceContext(ctx latestmeta.ChoiceContext) prodmeta.ChoiceContext {
	return mustAdapt[prodmeta.ChoiceContext](ctx)
}

// AdaptExtraArgs maps caller-side token extra args into v1_0_0 ledger Send bindings.
func AdaptExtraArgs(args latestmeta.ExtraArgs) prodmeta.ExtraArgs {
	return mustAdapt[prodmeta.ExtraArgs](args)
}

// AdaptInstrumentId maps caller-side instrument IDs into v1_0_0 ledger Send bindings.
func AdaptInstrumentId(id latestholding.InstrumentId) prodholding.InstrumentId {
	return mustAdapt[prodholding.InstrumentId](id)
}
