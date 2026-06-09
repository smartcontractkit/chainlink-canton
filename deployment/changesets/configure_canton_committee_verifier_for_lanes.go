package changesets

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// committeeVerifierLaneChainFamily delegates contract resolution to the Canton adapter but
// routes lane configure dispatch to CommitteeVerifier-only sequences (Run 2).
type committeeVerifierLaneChainFamily struct {
	ccipadapters.ChainFamily
	cv *adapters.CantonChainFamilyAdapter
}

func (w *committeeVerifierLaneChainFamily) ConfigureChainForLanes() *cldfops.Sequence[ccipadapters.ConfigureChainForLanesInput, ccipseq.OnChainOutput, cldfchain.BlockChains] {
	return w.cv.ConfigureCommitteeVerifierForLanes()
}

// noopConfigureLanesChainFamily skips lane configure on non-Canton chains during Run 2.
// Lane expansion still produces partial configs for remote families; this avoids
// configuring CCIP core contracts on those chains in the CV-only pass.
type noopConfigureLanesChainFamily struct {
	ccipadapters.ChainFamily
}

func (n *noopConfigureLanesChainFamily) ConfigureChainForLanes() *cldfops.Sequence[ccipadapters.ConfigureChainForLanesInput, ccipseq.OnChainOutput, cldfchain.BlockChains] {
	return cldfops.NewSequence(
		"canton/noop-configure-chain-for-lanes",
		semver.MustParse("2.0.0"),
		"No-op lane configure for non-Canton chains during CommitteeVerifier-only pass",
		func(_ cldfops.Bundle, _ cldfchain.BlockChains, _ ccipadapters.ConfigureChainForLanesInput) (ccipseq.OnChainOutput, error) {
			return ccipseq.OnChainOutput{}, nil
		},
	)
}

// ConfigureCantonCommitteeVerifierForLanesFromTopology is Run 2: Canton CommitteeVerifier lane
// configure only. Emits mcms-ccv timelock proposals; does not configure CCIP core contracts.
func ConfigureCantonCommitteeVerifierForLanesFromTopology(
	committeeVerifierContractRegistry *ccipadapters.CommitteeVerifierContractRegistry,
	chainFamilyRegistry *ccipadapters.ChainFamilyRegistry,
	mcmsRegistry *ccipchangesets.MCMSReaderRegistry,
) cldf.ChangeSetV2[v2cs.ConfigureChainsForLanesFromTopologyConfig] {
	validate := func(e cldf.Environment, cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) error {
		base := v2cs.ConfigureChainsForLanesFromTopology(committeeVerifierContractRegistry, chainFamilyRegistry, mcmsRegistry)
		if err := base.VerifyPreconditions(e, cfg); err != nil {
			return err
		}
		if !hasCantonChainsInConfig(cfg) {
			return fmt.Errorf("at least one Canton chain is required")
		}
		for _, sel := range cantonChainSelectorsInLanes(cfg.Lanes) {
			if err := requireTripleMCMSRefs(e, sel); err != nil {
				return err
			}
		}

		return nil
	}

	apply := func(e cldf.Environment, cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) (cldf.ChangesetOutput, error) {
		cvRegistry, err := chainFamilyRegistryForCommitteeVerifierOnly(chainFamilyRegistry, cfg)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		cvBase := v2cs.ConfigureChainsForLanesFromTopology(committeeVerifierContractRegistry, cvRegistry, mcmsRegistry)

		cfgCopy := cantonOnlyConfigureLanesConfig(cfg)
		cfgCopy.MCMS.Qualifier = cantonmcms.QualifierCCVOwner

		return cvBase.Apply(e, cfgCopy)
	}

	return cldf.CreateChangeSet(apply, validate)
}

func hasCantonChainsInConfig(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) bool {
	return len(cantonChainSelectorsInLanes(cfg.Lanes)) > 0
}

func cantonOnlyConfigureLanesConfig(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) v2cs.ConfigureChainsForLanesFromTopologyConfig {
	cfgCopy := cfg
	cfgCopy.Lanes = nil
	for _, lane := range cfg.Lanes {
		if laneInvolvesCanton(lane) {
			cfgCopy.Lanes = append(cfgCopy.Lanes, lane)
		}
	}

	return cfgCopy
}

func laneInvolvesCanton(lane v2cs.CrossFamilyLanePair) bool {
	for _, sel := range []uint64{lane.ChainA, lane.ChainB} {
		family, err := chainsel.GetSelectorFamily(sel)
		if err == nil && family == chainsel.FamilyCanton {
			return true
		}
	}

	return false
}

func cantonChainSelectorsInLanes(lanes []v2cs.CrossFamilyLanePair) []uint64 {
	seen := make(map[uint64]struct{})
	var selectors []uint64
	for _, lane := range lanes {
		for _, sel := range []uint64{lane.ChainA, lane.ChainB} {
			if _, ok := seen[sel]; ok {
				continue
			}
			family, err := chainsel.GetSelectorFamily(sel)
			if err != nil || family != chainsel.FamilyCanton {
				continue
			}
			seen[sel] = struct{}{}
			selectors = append(selectors, sel)
		}
	}

	return selectors
}

func chainFamilyRegistryForCommitteeVerifierOnly(
	base *ccipadapters.ChainFamilyRegistry,
	cfg v2cs.ConfigureChainsForLanesFromTopologyConfig,
) (*ccipadapters.ChainFamilyRegistry, error) {
	reg := ccipadapters.NewChainFamilyRegistry()
	cantonCV := &adapters.CantonChainFamilyAdapter{}

	for family := range familiesInConfigureLanesConfig(cfg) {
		adapter, ok := base.GetChainFamily(family)
		if !ok {
			return nil, fmt.Errorf("no adapter registered for chain family %q", family)
		}
		if family == chainsel.FamilyCanton {
			reg.RegisterChainFamily(family, &committeeVerifierLaneChainFamily{
				ChainFamily: adapter,
				cv:          cantonCV,
			})

			continue
		}
		reg.RegisterChainFamily(family, &noopConfigureLanesChainFamily{ChainFamily: adapter})
	}

	return reg, nil
}

func familiesInConfigureLanesConfig(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) map[string]struct{} {
	families := make(map[string]struct{})
	for _, lane := range cfg.Lanes {
		for _, sel := range []uint64{lane.ChainA, lane.ChainB} {
			if family, err := chainsel.GetSelectorFamily(sel); err == nil {
				families[family] = struct{}{}
			}
		}
	}

	return families
}
