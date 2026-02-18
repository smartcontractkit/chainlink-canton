package dependencies

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

type CantonDeps = struct {
	Chain canton.Chain
	// Which participant in the Chain to use (0-indexed, i.e. defaults to the first participant)
	Participant int
}
