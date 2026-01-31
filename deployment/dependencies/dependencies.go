package dependencies

import (
	"github.com/noders-team/go-daml/pkg/client"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

// TODO: use the CLDF cldf_chain.BlockChains directly instead once client & party have been added there

type CantonDeps = struct {
	Chain         canton.Chain
	BindingClient *client.DamlBindingClient
	Party         string
}
