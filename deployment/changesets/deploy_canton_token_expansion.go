package changesets

import (
	"fmt"

	cciptokens "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// DeployCantonTokenExpansion deploys Canton token pools (and optionally resolves token refs) and
// emits one MCMS timelock proposal against mcms-ccip (ccipOwner).
//
// Run after DeployCantonChainContracts and ConfigureCantonChainsForLanesFromTopology. Omit
// TokenTransferConfig from per-chain input for deploy-only; use ConfigureCantonTokensForTransfers
// to wire remote lanes, rate limiters, and TAR registration for transfers.
func DeployCantonTokenExpansion() cldf.ChangeSetV2[cciptokens.TokenExpansionInput] {
	base := cciptokens.TokenExpansion()

	validate := func(e cldf.Environment, cfg cciptokens.TokenExpansionInput) error {
		if err := cfg.MCMS.Validate(); err != nil {
			return err
		}
		if len(cfg.TokenExpansionInputPerChain) == 0 {
			return fmt.Errorf("tokenExpansionInputPerChain is required")
		}

		selectors := make([]uint64, 0, len(cfg.TokenExpansionInputPerChain))
		for sel := range cfg.TokenExpansionInputPerChain {
			selectors = append(selectors, sel)
		}

		return requireCantonCCIPOwnerMCMSForChains(e, selectors)
	}

	apply := func(e cldf.Environment, cfg cciptokens.TokenExpansionInput) (cldf.ChangesetOutput, error) {
		cfgCopy := cfg
		cfgCopy.MCMS.Qualifier = cantonmcms.QualifierCCIPOwner

		return base.Apply(e, cfgCopy)
	}

	return cldf.CreateChangeSet(apply, validate)
}
