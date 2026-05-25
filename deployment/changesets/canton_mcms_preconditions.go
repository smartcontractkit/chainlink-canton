package changesets

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

func requireCCIPOwnerMCMSRef(e cldf.Environment, sel uint64) error {
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, sel, cantonmcms.QualifierCCIPOwner); err != nil {
		return fmt.Errorf("ccipOwner MCMS must be deployed first (mcms-ccip): %w", err)
	}

	return nil
}

func requireDualMCMSRefs(e cldf.Environment, sel uint64) error {
	if err := requireCCIPOwnerMCMSRef(e, sel); err != nil {
		return fmt.Errorf("ccipOwner MCMS must be deployed first (canton_devnet_deploy_mcms_ccip.yaml): %w", err)
	}
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, sel, cantonmcms.QualifierCCVOwner); err != nil {
		return fmt.Errorf("ccvOwner MCMS must be deployed first (canton_devnet_deploy_mcms_ccv.yaml): %w", err)
	}

	return nil
}

func requireCantonCCIPOwnerMCMSForChains(e cldf.Environment, chainSelectors []uint64) error {
	for _, sel := range chainSelectors {
		family, err := chainsel.GetSelectorFamily(sel)
		if err != nil {
			return fmt.Errorf("chain %d: %w", sel, err)
		}
		if family != chainsel.FamilyCanton {
			continue
		}
		if err := requireCCIPOwnerMCMSRef(e, sel); err != nil {
			return err
		}
	}

	return nil
}
