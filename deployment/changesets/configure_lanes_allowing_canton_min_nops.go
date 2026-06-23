package changesets

import (
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	v2_0_0 "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ConfigureChainsForLanesFromTopologyAllowingCantonMinNOPs wraps chainlink-ccip's
// ConfigureChainsForLanesFromTopology and skips the production minimum-15-NOP committee
// rule for Canton-family chain selectors (Canton prod runs 4 NOPs).
func ConfigureChainsForLanesFromTopologyAllowingCantonMinNOPs(
	committeeVerifierContractRegistry *adapters.CommitteeVerifierContractRegistry,
	chainFamilyRegistry *adapters.ChainFamilyRegistry,
	mcmsRegistry *changesets.MCMSReaderRegistry,
) deployment.ChangeSetV2[v2_0_0.ConfigureChainsForLanesFromTopologyConfig] {
	inner := v2_0_0.ConfigureChainsForLanesFromTopology(
		committeeVerifierContractRegistry,
		chainFamilyRegistry,
		mcmsRegistry,
	)

	return deployment.CreateChangeSet(
		inner.Apply,
		func(e deployment.Environment, cfg v2_0_0.ConfigureChainsForLanesFromTopologyConfig) error {
			return WithCantonProductionMinNOPCheckBypassed(cfg.Topology, func() error {
				return inner.VerifyPreconditions(e, cfg)
			})
		},
	)
}

// ConfigureCantonCommitteeVerifierForLanesFromTopologyAllowingCantonMinNOPs wraps Run 2
// (Canton CommitteeVerifier lane configure) with the same production minimum-NOP bypass as Run 1.
func ConfigureCantonCommitteeVerifierForLanesFromTopologyAllowingCantonMinNOPs(
	committeeVerifierContractRegistry *adapters.CommitteeVerifierContractRegistry,
	chainFamilyRegistry *adapters.ChainFamilyRegistry,
	mcmsRegistry *changesets.MCMSReaderRegistry,
) deployment.ChangeSetV2[v2_0_0.ConfigureChainsForLanesFromTopologyConfig] {
	inner := ConfigureCantonCommitteeVerifierForLanesFromTopology(
		committeeVerifierContractRegistry,
		chainFamilyRegistry,
		mcmsRegistry,
	)

	return deployment.CreateChangeSet(
		inner.Apply,
		func(e deployment.Environment, cfg v2_0_0.ConfigureChainsForLanesFromTopologyConfig) error {
			return WithCantonProductionMinNOPCheckBypassed(cfg.Topology, func() error {
				return inner.VerifyPreconditions(e, cfg)
			})
		},
	)
}
