package contracts

import (
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_metadata_v1"
)

func ChoiceContextFromData(contextData map[string]any) (splice_api_token_metadata_v1.ChoiceContext, error) {
	ctx := splice_api_token_metadata_v1.ChoiceContext{}
	err := ledger.MapToStruct(contextData, &ctx)
	if err != nil {
		return splice_api_token_metadata_v1.ChoiceContext{}, err
	}

	return ctx, nil
}
