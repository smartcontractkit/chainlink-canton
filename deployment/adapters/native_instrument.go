package adapters

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/nativeinstrument"
)

func lookupNativeInstrumentID(
	ctx context.Context,
	participant canton.Participant,
	ds datastore.DataStore,
	chainSelector uint64,
) (splice_api_token_holding_v1.InstrumentId, error) {
	return nativeinstrument.ResolveNativeInstrumentID(ctx, participant, ds, chainSelector)
}

func instrumentPriceKey(admin types.PARTY, id types.TEXT) string {
	return fmt.Sprintf("%s:%s", admin, id)
}
