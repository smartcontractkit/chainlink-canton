package datastore

import (
	"fmt"

	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

// ProposerMCMSAddressRef resolves the MCMS instance used for proposal scheduling.
func ProposerMCMSAddressRef(ds datastore.DataStore, chainSelector uint64, qualifier string) (datastore.AddressRef, error) {
	if qualifier == "" {
		return datastore.AddressRef{}, fmt.Errorf("MCMS qualifier is required")
	}

	refs := ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig)),
		datastore.AddressRefByQualifier(qualifier),
	)
	switch len(refs) {
	case 1:
		return refs[0], nil
	case 0:
		return datastore.AddressRef{}, fmt.Errorf("no ProposerManyChainMultiSig ref for chain %d qualifier %q", chainSelector, qualifier)
	default:
		return datastore.AddressRef{}, fmt.Errorf("multiple ProposerManyChainMultiSig refs for chain %d qualifier %q", chainSelector, qualifier)
	}
}

// MCMSRawInstanceAddress returns the raw instance address label for an MCMS owner qualifier.
func MCMSRawInstanceAddress(ds datastore.DataStore, chainSelector uint64, qualifier string) (contracts.RawInstanceAddress, error) {
	ref, err := ProposerMCMSAddressRef(ds, chainSelector, qualifier)
	if err != nil {
		return "", err
	}

	labels := ref.Labels.List()
	if len(labels) == 0 {
		return "", fmt.Errorf("MCMS ref for qualifier %q is missing raw instance address label", qualifier)
	}

	return contracts.RawInstanceAddressFromString(labels[0])
}
