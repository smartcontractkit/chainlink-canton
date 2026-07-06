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

func requireRMNOwnerMCMSRef(e cldf.Environment, sel uint64) error {
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, sel, cantonmcms.QualifierRMNOwner); err != nil {
		return fmt.Errorf("rmnOwner MCMS must be deployed first (mcms-rmn): %w", err)
	}

	return nil
}

func requireCCVOwnerMCMSRef(e cldf.Environment, sel uint64) error {
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, sel, cantonmcms.QualifierCCVOwner); err != nil {
		return fmt.Errorf("ccvOwner MCMS must be deployed first (mcms-ccv): %w", err)
	}

	return nil
}

// requireCommitteeVerifierLaneMCMSRefs allows Run 2 on prod (mcms-ccv@ccvOwner) or staging
// single-MCMS (mcms-ccip@ccipOwner only, CV contracts also under ccipOwner).
func requireCommitteeVerifierLaneMCMSRefs(e cldf.Environment, sel uint64) error {
	if err := requireCCVOwnerMCMSRef(e, sel); err == nil {
		return nil
	}

	return requireCCIPOwnerMCMSRef(e, sel)
}

// committeeVerifierLaneMCMSQualifier picks the MCMS root for CV lane proposals.
// Uses ccvOwner when mcms-ccv is deployed; otherwise falls back to CLLCCIP (mcms-ccip).
func committeeVerifierLaneMCMSQualifier(e cldf.Environment, sel uint64) (string, error) {
	if err := requireCCVOwnerMCMSRef(e, sel); err == nil {
		return cantonmcms.QualifierCCVOwner, nil
	}
	if err := requireCCIPOwnerMCMSRef(e, sel); err != nil {
		return "", fmt.Errorf("committee verifier lane configure requires mcms-ccv or mcms-ccip on Canton: %w", err)
	}

	return cantonmcms.QualifierCCIPOwner, nil
}

// requireTripleMCMSRefs ensures mcms-ccip, mcms-rmn, and mcms-ccv are in the datastore.
func requireTripleMCMSRefs(e cldf.Environment, sel uint64) error {
	if err := requireCCIPOwnerMCMSRef(e, sel); err != nil {
		return fmt.Errorf("ccipOwner MCMS must be deployed first (canton_devnet_deploy_mcms_ccip.yaml): %w", err)
	}
	if err := requireRMNOwnerMCMSRef(e, sel); err != nil {
		return fmt.Errorf("rmnOwner MCMS must be deployed first (canton_devnet_deploy_mcms_rmn.yaml): %w", err)
	}
	if err := requireCCVOwnerMCMSRef(e, sel); err != nil {
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
