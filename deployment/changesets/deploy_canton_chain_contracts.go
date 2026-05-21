package changesets

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	v2cs "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// DeployCantonChainContracts deploys Canton CCIP contracts and emits separate MCMS timelock
// proposals for ccipOwner and ccvOwner (same signers, distinct MCMS instance addresses).
func DeployCantonChainContracts(
	registry *ccipadapters.DeployChainContractsRegistry,
) cldf.ChangeSetV2[ccipchangesets.WithMCMS[v2cs.DeployChainContractsCfg]] {
	validate := func(e cldf.Environment, cfg ccipchangesets.WithMCMS[v2cs.DeployChainContractsCfg]) error {
		if err := cfg.MCMS.Validate(); err != nil {
			return err
		}
		if cfg.Cfg.Topology == nil {
			return fmt.Errorf("topology is required")
		}
		if len(cfg.Cfg.ChainSelectors) == 0 {
			return fmt.Errorf("at least one chain selector is required")
		}
		for _, sel := range cfg.Cfg.ChainSelectors {
			family, err := chainsel.GetSelectorFamily(sel)
			if err != nil {
				return fmt.Errorf("chain %d: %w", sel, err)
			}
			if family != chainsel.FamilyCanton {
				return fmt.Errorf("chain %d: deploy canton chain contracts only supports Canton (got %s)", sel, family)
			}
			if err := requireDualMCMSRefs(e, sel); err != nil {
				return err
			}
		}

		return nil
	}

	apply := func(e cldf.Environment, cfg ccipchangesets.WithMCMS[v2cs.DeployChainContractsCfg]) (cldf.ChangesetOutput, error) {
		ds := datastore.NewMemoryDataStore()
		var allReports []operations.Report[any, any]
		var proposals []mcms.TimelockProposal

		adapter, ok := registry.Get(chainsel.FamilyCanton)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("no Canton deploy chain contracts adapter registered")
		}

		for _, sel := range cfg.Cfg.ChainSelectors {
			committeeVerifiers, err := v2cs.BuildCommitteeVerifierParams(cfg.Cfg.Topology, sel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("build committee verifier params for chain %d: %w", sel, err)
			}

			perChain := cfg.Cfg.DefaultCfg
			if override, ok := cfg.Cfg.ChainCfgs[sel]; ok {
				perChain = override
			}
			existingAddresses := e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(sel))

			input := ccipadapters.DeployChainContractsInput{
				ChainSelector:     sel,
				DeployerContract:  perChain.DeployerContract,
				DeployTestRouter:  perChain.DeployTestRouter,
				ExistingAddresses: existingAddresses,
				DeployerKeyOwned:  perChain.DeployerKeyOwned,
				ContractParams: ccipadapters.DeployContractParams{
					RMNRemote:          perChain.RMNRemote,
					OffRamp:            perChain.OffRamp,
					CommitteeVerifiers: committeeVerifiers,
					OnRamp:             perChain.OnRamp,
					FeeQuoter:          perChain.FeeQuoter,
					Executors:          perChain.Executors,
					MockReceivers:      perChain.MockReceivers,
				},
			}

			report, err := operations.ExecuteSequence(e.OperationsBundle, adapter.DeployChainContracts(), e.BlockChains, input)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("deploy canton chain contracts on %d: %w", sel, err)
			}

			for _, ref := range report.Output.Addresses {
				if err := ds.Addresses().Add(ref); err != nil {
					return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("store address ref: %w", err)
				}
			}
			if err := ccipseq.WriteMetadataToDatastore(ds, report.Output.Metadata); err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("write metadata: %w", err)
			}
			allReports = append(allReports, report.ExecutionReports...)

			chain := e.BlockChains.CantonChains()[sel]
			chainProposals, err := cantonmcms.BuildDualTimelockProposalsFromBatchOps(
				e.GetContext(), e, chain, sel, report.Output.BatchOps, cfg.MCMS, cfg.MCMS.Description,
			)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, err
			}
			proposals = append(proposals, chainProposals...)
		}

		return cldf.ChangesetOutput{
			DataStore:             ds,
			Reports:               allReports,
			MCMSTimelockProposals: proposals,
		}, nil
	}

	return cldf.CreateChangeSet(apply, validate)
}

func requireDualMCMSRefs(e cldf.Environment, sel uint64) error {
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, sel, cantonmcms.QualifierCCIPOwner); err != nil {
		return fmt.Errorf("ccipOwner MCMS must be deployed first (canton_devnet_deploy_mcms_ccip.yaml): %w", err)
	}
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, sel, cantonmcms.QualifierCCVOwner); err != nil {
		return fmt.Errorf("ccvOwner MCMS must be deployed first (canton_devnet_deploy_mcms_ccv.yaml): %w", err)
	}

	return nil
}
