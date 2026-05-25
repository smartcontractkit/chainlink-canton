package changesets

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
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
		for _, chainCfg := range cfg.Chains {
			family, err := chainsel.GetSelectorFamily(chainCfg.ChainSelector)
			if err != nil {
				return fmt.Errorf("chain %d: %w", chainCfg.ChainSelector, err)
			}
			if family != chainsel.FamilyCanton {
				continue
			}
			if err := requireDualMCMSRefs(e, chainCfg.ChainSelector); err != nil {
				return err
			}
			if len(chainCfg.CommitteeVerifiers) == 0 {
				return fmt.Errorf("chain %d: committeeverifiers required for CV lane configure", chainCfg.ChainSelector)
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
	for _, chainCfg := range cfg.Chains {
		family, err := chainsel.GetSelectorFamily(chainCfg.ChainSelector)
		if err != nil {
			continue
		}
		if family == chainsel.FamilyCanton {
			return true
		}
	}

	return false
}

func cantonOnlyConfigureLanesConfig(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) v2cs.ConfigureChainsForLanesFromTopologyConfig {
	cfgCopy := cfg
	cfgCopy.Chains = nil
	for _, chainCfg := range cfg.Chains {
		family, err := chainsel.GetSelectorFamily(chainCfg.ChainSelector)
		if err != nil || family != chainsel.FamilyCanton {
			continue
		}
		cfgCopy.Chains = append(cfgCopy.Chains, chainCfg)
	}

	return cfgCopy
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
		reg.RegisterChainFamily(family, adapter)
	}

	return reg, nil
}

func familiesInConfigureLanesConfig(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) map[string]struct{} {
	families := make(map[string]struct{})
	for _, chainCfg := range cfg.Chains {
		if family, err := chainsel.GetSelectorFamily(chainCfg.ChainSelector); err == nil {
			families[family] = struct{}{}
		}
		for _, cv := range chainCfg.CommitteeVerifiers {
			for remoteSelector := range cv.RemoteChains {
				if family, err := chainsel.GetSelectorFamily(remoteSelector); err == nil {
					families[family] = struct{}{}
				}
			}
		}
		for remoteSelector := range chainCfg.RemoteChains {
			if family, err := chainsel.GetSelectorFamily(remoteSelector); err == nil {
				families[family] = struct{}{}
			}
		}
	}

	return families
}
