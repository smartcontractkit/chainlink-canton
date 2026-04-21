package adapters

import (
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	mcmsreaderapi "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
)

var tokenPoolVersions = []string{
	"1.6.1",
	"2.0.0",
}

func init() {
	ccipadapters.GetDeployChainContractsRegistry().Register(chainsel.FamilyCanton, &CantonDeployChainContractsAdapter{})
	ccipadapters.GetChainFamilyRegistry().RegisterChainFamily(chainsel.FamilyCanton, &CantonChainFamilyAdapter{})
	lanes.GetLaneAdapterRegistry().RegisterLaneAdapter(chainsel.FamilyCanton, semver.MustParse("2.0.0"), CantonLaneAdapter{})
	mcmsreaderapi.GetRegistry().RegisterMCMSReader(chainsel.FamilyCanton, &CantonMCMSReader{})
	ccipadapters.GetCommitteeVerifierContractRegistry().Register(chainsel.FamilyCanton, &CantonCommitteeVerifierContractAdapter{})
	ccipadapters.GetAggregatorConfigRegistry().Register(chainsel.FamilyCanton, &CantonAggregatorConfigAdapter{})
	ccipadapters.GetIndexerConfigRegistry().Register(chainsel.FamilyCanton, &CantonIndexerConfigAdapter{})
	ccipadapters.GetVerifierJobConfigRegistry().Register(chainsel.FamilyCanton, &CantonVerifierJobConfigAdapter{})
	ccipadapters.GetExecutorConfigRegistry().Register(chainsel.FamilyCanton, &CantonExecutorConfigAdapter{})
	// Register the canton token adapter for the canton family.
	for _, version := range tokenPoolVersions {
		tokenscore.GetTokenAdapterRegistry().RegisterTokenAdapter(chainsel.FamilyCanton, semver.MustParse(version), CantonTokenAdapter{})
	}
}
