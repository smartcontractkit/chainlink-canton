package adapters

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
)

type CantonCommitteeVerifierContractAdapter struct{}

var _ adapters.CommitteeVerifierContractAdapter = (*CantonCommitteeVerifierContractAdapter)(nil)

func (c CantonCommitteeVerifierContractAdapter) ResolveCommitteeVerifierContracts(ds datastore.DataStore, chainSelector uint64, qualifier string) ([]datastore.AddressRef, error) {
	verifier, err := ds.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(committee_verifier.ContractType),
		committee_verifier.Version,
		qualifier,
	))
	if err != nil {
		return nil, fmt.Errorf("committee verifier not found for chain %d qualifier %q: %w", chainSelector, qualifier, err)
	}

	return []datastore.AddressRef{verifier}, nil
}

// GetCommitteeVerifierResolver implements [adapters.CommitteeVerifierContractAdapter].
// Canton has no separate versioned verifier resolver; the CommitteeVerifier is the lane CCV.
func (c CantonCommitteeVerifierContractAdapter) GetCommitteeVerifierResolver(ds datastore.DataStore, chainSelector uint64, qualifier string) ([]datastore.AddressRef, error) {
	verifier, err := ds.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(committee_verifier.ContractType),
		committee_verifier.Version,
		qualifier,
	))
	if err != nil {
		return nil, fmt.Errorf("committee verifier not found for chain %d qualifier %q: %w", chainSelector, qualifier, err)
	}

	return []datastore.AddressRef{verifier}, nil
}
