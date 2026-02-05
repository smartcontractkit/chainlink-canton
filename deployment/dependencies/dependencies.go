package dependencies

import (
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

// TODO: use the CLDF cldf_chain.BlockChains directly instead once client & party have been added there

type CantonDeps = struct {
	Chain              canton.Chain
	CommandServiceClient apiv2.CommandServiceClient
	StateServiceClient   apiv2.StateServiceClient
	Party                string
}
