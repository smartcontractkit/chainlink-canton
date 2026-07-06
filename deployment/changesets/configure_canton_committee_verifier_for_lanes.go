package changesets

import (
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/adapters"
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
		filteredCfg := topologyForCantonCommitteeVerifierConfigure(cantonOnlyConfigureLanesConfig(cfg))
		if err := base.VerifyPreconditions(e, filteredCfg); err != nil {
			return err
		}
		if !hasCantonChainsInConfig(cfg) {
			return fmt.Errorf("at least one Canton chain is required")
		}
		for _, sel := range cantonChainSelectorsInLanes(cfg.Lanes) {
			if err := requireCommitteeVerifierLaneMCMSRefs(e, sel); err != nil {
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

		cfgCopy := topologyForCantonCommitteeVerifierConfigure(cantonOnlyConfigureLanesConfig(cfg))
		cantonSelectors := cantonChainSelectorsInLanes(cfgCopy.Lanes)
		if len(cantonSelectors) == 0 {
			return cldf.ChangesetOutput{}, fmt.Errorf("at least one Canton chain is required")
		}
		qualifier, err := committeeVerifierLaneMCMSQualifier(e, cantonSelectors[0])
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		cfgCopy.MCMS.Qualifier = qualifier

		return cvBase.Apply(e, cfgCopy)
	}

	return cldf.CreateChangeSet(apply, validate)
}

// ConfigureCantonBuildLanesCrossFamilyFromTopology is Run 1 for envs whose topology
// includes EVM-only committees (e.g. staging_testnet "secondary": Sepolia chain_config
// but no Canton CV). expandLanesToPartialChainConfigs would otherwise require a
// CommitteeVerifier ref on Canton for every committee with remote chain_config.
// prod_testnet/prod_mainnet have no secondary committee and can use stock chainlink-ccip.
func ConfigureCantonBuildLanesCrossFamilyFromTopology(
	committeeVerifierContractRegistry *ccipadapters.CommitteeVerifierContractRegistry,
	chainFamilyRegistry *ccipadapters.ChainFamilyRegistry,
	mcmsRegistry *ccipchangesets.MCMSReaderRegistry,
) cldf.ChangeSetV2[v2cs.ConfigureChainsForLanesFromTopologyConfig] {
	base := v2cs.ConfigureChainsForLanesFromTopology(
		committeeVerifierContractRegistry,
		chainFamilyRegistry,
		mcmsRegistry,
	)

	return cldf.CreateChangeSet(
		func(e cldf.Environment, cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) (cldf.ChangesetOutput, error) {
			return base.Apply(e, topologyForCantonCommitteeVerifierConfigure(cfg))
		},
		func(e cldf.Environment, cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) error {
			return base.VerifyPreconditions(e, topologyForCantonCommitteeVerifierConfigure(cfg))
		},
	)
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

// topologyForCantonCommitteeVerifierConfigure drops EVM-only committees from the in-memory
// topology copy used for Run 2. expandLanesToPartialChainConfigs selects a committee
// qualifier whenever it has chain_config for the remote selector; on Canton that can
// include committees (e.g. secondary) that verify Sepolia traffic on EVM but are not
// members on Canton. Those duplicate MCMS ops against the same CommitteeVerifier instance.
// Mainnet Canton↔Eth Run 2 only emits ops for committees with Canton chain_config (see
// chainlink-deployments PR #15556).
func topologyForCantonCommitteeVerifierConfigure(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) v2cs.ConfigureChainsForLanesFromTopologyConfig {
	if cfg.Topology == nil || cfg.Topology.NOPTopology == nil {
		return cfg
	}

	cantonSelectors := cantonChainSelectorsInLanes(cfg.Lanes)
	if len(cantonSelectors) == 0 {
		return cfg
	}

	remotesByCanton := remotesByCantonSelector(cfg.Lanes, cantonSelectors)
	filteredCommittees := make(map[string]offchain.CommitteeConfig, len(cfg.Topology.NOPTopology.Committees))
	for name, committee := range cfg.Topology.NOPTopology.Committees {
		if keepCommitteeForCantonCommitteeVerifierConfigure(committee, cantonSelectors, remotesByCanton) {
			filteredCommittees[name] = committee
		}
	}

	topoCopy := *cfg.Topology
	nopCopy := *cfg.Topology.NOPTopology
	nopCopy.Committees = filteredCommittees
	topoCopy.NOPTopology = &nopCopy

	cfgCopy := cfg
	cfgCopy.Topology = &topoCopy

	return cfgCopy
}

func remotesByCantonSelector(lanes []v2cs.CrossFamilyLanePair, cantonSelectors []uint64) map[uint64]map[uint64]struct{} {
	cantonSet := make(map[uint64]struct{}, len(cantonSelectors))
	for _, sel := range cantonSelectors {
		cantonSet[sel] = struct{}{}
	}

	out := make(map[uint64]map[uint64]struct{})
	for _, lane := range lanes {
		for cantonSel := range cantonSet {
			var remote uint64
			switch {
			case lane.ChainA == cantonSel:
				remote = lane.ChainB
			case lane.ChainB == cantonSel:
				remote = lane.ChainA
			default:
				continue
			}
			if out[cantonSel] == nil {
				out[cantonSel] = make(map[uint64]struct{})
			}
			out[cantonSel][remote] = struct{}{}
		}
	}

	return out
}

func keepCommitteeForCantonCommitteeVerifierConfigure(
	committee offchain.CommitteeConfig,
	cantonSelectors []uint64,
	remotesByCanton map[uint64]map[uint64]struct{},
) bool {
	for _, cantonSel := range cantonSelectors {
		localKey := strconv.FormatUint(cantonSel, 10)
		_, hasLocalConfig := committee.ChainConfigs[localKey]

		for remoteSel := range remotesByCanton[cantonSel] {
			remoteKey := strconv.FormatUint(remoteSel, 10)
			if _, hasRemoteConfig := committee.ChainConfigs[remoteKey]; hasRemoteConfig && !hasLocalConfig {
				return false
			}
		}
	}

	return true
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
