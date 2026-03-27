package adapters

import (
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
)

func init() {
	v := semver.MustParse("2.0.0")

	lanes.GetLaneAdapterRegistry().RegisterLaneAdapter(chainsel.FamilyCanton, v, &ChainFamilyAdapter{})
	adapters.GetCommitteeVerifierContractRegistry().Register(chainsel.FamilyCanton, &CantonCommitteeVerifierContractAdapter{})
	adapters.GetAggregatorConfigRegistry().Register(chainsel.FamilyCanton, &CantonAggregatorConfigAdapter{})
}
