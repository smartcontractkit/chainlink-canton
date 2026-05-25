package changesets

import (
	"fmt"

	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	cciptokens "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// ConfigureCantonTokensForTransfers configures Canton token pools for cross-chain transfers and
// emits one MCMS timelock proposal against mcms-ccip (ccipOwner). The proposal may include TAR
// registration, rate limiter deploys (via ccip factory), and token pool ApplyChainUpdates.
//
// Run after DeployCantonTokenExpansion when pools are deployed and lane configure has completed.
func ConfigureCantonTokensForTransfers(
	mcmsRegistry *ccipchangesets.MCMSReaderRegistry,
) cldf.ChangeSetV2[cciptokens.ConfigureTokensForTransfersConfig] {
	base := cciptokens.ConfigureTokensForTransfers(cciptokens.GetTokenAdapterRegistry(), mcmsRegistry)

	validate := func(e cldf.Environment, cfg cciptokens.ConfigureTokensForTransfersConfig) error {
		if err := cfg.MCMS.Validate(); err != nil {
			return err
		}
		if len(cfg.Tokens) == 0 {
			return fmt.Errorf("at least one token transfer config is required")
		}

		selectors := make([]uint64, 0, len(cfg.Tokens))
		for _, token := range cfg.Tokens {
			selectors = append(selectors, token.ChainSelector)
		}

		return requireCantonCCIPOwnerMCMSForChains(e, selectors)
	}

	apply := func(e cldf.Environment, cfg cciptokens.ConfigureTokensForTransfersConfig) (cldf.ChangesetOutput, error) {
		cfgCopy := cfg
		cfgCopy.MCMS.Qualifier = cantonmcms.QualifierCCIPOwner

		return base.Apply(e, cfgCopy)
	}

	return cldf.CreateChangeSet(apply, validate)
}
