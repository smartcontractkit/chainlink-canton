package changesets

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// ConfigureCantonChainsForLanesFromTopology configures CCIP lanes and emits separate MCMS timelock
// proposals for ccipOwner (_0) and ccvOwner (_1) on Canton chains. Each proposal uses a distinct
// MCMS root (mcms-ccip vs mcms-ccv) and they may be signed and scheduled in parallel after deploy
// has written CV addresses to the datastore.
//
// Canton contract refs (GlobalConfig, Router, OnRamp, OffRamp, FeeQuoter, Executor, CommitteeVerifier)
// are resolved inside CantonChainFamilyAdapter.ConfigureChainForLanes from the runtime datastore.
// This changeset runs topology enrichment via the generic CCIP changeset, then re-splits adapter
// batch ops into dual proposals instead of emitting a single merged MCMS proposal.
func ConfigureCantonChainsForLanesFromTopology(
	committeeVerifierContractRegistry *ccipadapters.CommitteeVerifierContractRegistry,
	chainFamilyRegistry *ccipadapters.ChainFamilyRegistry,
	mcmsRegistry *ccipchangesets.MCMSReaderRegistry,
) cldf.ChangeSetV2[v2cs.ConfigureChainsForLanesFromTopologyConfig] {
	base := v2cs.ConfigureChainsForLanesFromTopology(committeeVerifierContractRegistry, chainFamilyRegistry, mcmsRegistry)

	validate := func(e cldf.Environment, cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) error {
		if err := base.VerifyPreconditions(e, cfg); err != nil {
			return err
		}
		for _, chainCfg := range cfg.Chains {
			family, err := chainsel.GetSelectorFamily(chainCfg.ChainSelector)
			if err != nil {
				return fmt.Errorf("chain %d: %w", chainCfg.ChainSelector, err)
			}
			if family == chainsel.FamilyCanton {
				if err := requireDualMCMSRefs(e, chainCfg.ChainSelector); err != nil {
					return err
				}
			}
		}

		return nil
	}

	apply := func(e cldf.Environment, cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) (cldf.ChangesetOutput, error) {
		adapters.SetRuntimeDataStore(e.DataStore)

		if !needsCantonDualLaneProposals(cfg) {
			return base.Apply(e, cfg)
		}

		// The generic changeset enriches topology input and dispatches to the Canton adapter.
		// MCMS.Qualifier only affects proposal metadata here; adapter batch ops are re-split below.
		cfgCopy := cfg
		cfgCopy.MCMS.Qualifier = cantonmcms.QualifierCCIPOwner

		baseOut, err := base.Apply(e, cfgCopy)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		return buildDualConfigureLaneProposals(e, cfg, mcmsRegistry, baseOut)
	}

	return cldf.CreateChangeSet(apply, validate)
}

func needsCantonDualLaneProposals(cfg v2cs.ConfigureChainsForLanesFromTopologyConfig) bool {
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

func buildDualConfigureLaneProposals(
	e cldf.Environment,
	cfg v2cs.ConfigureChainsForLanesFromTopologyConfig,
	mcmsRegistry *ccipchangesets.MCMSReaderRegistry,
	baseOut cldf.ChangesetOutput,
) (cldf.ChangesetOutput, error) {
	if len(baseOut.MCMSTimelockProposals) == 0 {
		return baseOut, nil
	}

	batchOps := cantonmcms.BatchOpsFromTimelockProposals(baseOut.MCMSTimelockProposals)
	desc := cfg.MCMS.Description

	var proposals []mcms.TimelockProposal
	var evmBatchOps []mcms_types.BatchOperation

	for _, op := range batchOps {
		if len(op.Transactions) == 0 {
			continue
		}

		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: %w", op.ChainSelector, err)
		}
		if family != chainsel.FamilyCanton {
			evmBatchOps = append(evmBatchOps, op)
			continue
		}

		chain, ok := e.BlockChains.CantonChains()[uint64(op.ChainSelector)]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("canton chain %d not loaded", op.ChainSelector)
		}

		chainProposals, err := cantonmcms.BuildDualTimelockProposalsFromBatchOps(
			e.GetContext(), e, chain, uint64(op.ChainSelector), []mcms_types.BatchOperation{op}, cfg.MCMS, desc,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		proposals = append(proposals, chainProposals...)
	}

	if len(evmBatchOps) > 0 {
		evmOut, err := ccipchangesets.NewOutputBuilder(e, mcmsRegistry).
			WithReports(baseOut.Reports).
			WithDataStore(baseOut.DataStore).
			WithBatchOps(evmBatchOps).
			Build(cfg.MCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("build non-Canton lane proposal: %w", err)
		}
		proposals = append(proposals, evmOut.MCMSTimelockProposals...)
	}

	baseOut.MCMSTimelockProposals = proposals

	return baseOut, nil
}
