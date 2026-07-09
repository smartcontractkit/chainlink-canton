package adapters

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/versioned_verifier_resolver"
	dsutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	ccvadapters "github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
)

type CantonVerifierJobConfigAdapter struct{}

var _ ccvadapters.VerifierConfigAdapter = (*CantonVerifierJobConfigAdapter)(nil)

// GetSignerAddressFamily implements [adapters.VerifierConfigAdapter].
func (a *CantonVerifierJobConfigAdapter) GetSignerAddressFamily() string {
	return chainsel.FamilyCanton
}

// ResolveVerifierContractAddresses implements [adapters.VerifierConfigAdapter].
func (a *CantonVerifierJobConfigAdapter) ResolveVerifierContractAddresses(
	ds datastore.DataStore,
	chainSelector uint64,
	committeeQualifier string,
	executorQualifier string,
) (*ccvadapters.VerifierContractAddresses, error) {
	toAddress := func(ref datastore.AddressRef) (string, error) { return ref.Address, nil }

	committeeVerifierAddr, err := dsutils.FindAndFormatFirstRef(ds, chainSelector, toAddress,
		datastore.AddressRef{
			Type:      datastore.ContractType(versioned_verifier_resolver.CommitteeVerifierResolverType),
			Qualifier: committeeQualifier,
		},
		datastore.AddressRef{
			Type:      datastore.ContractType(committee_verifier.ContractType),
			Qualifier: committeeQualifier,
			Version:   committee_verifier.Version,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get committee verifier address for chain %d: %w", chainSelector, err)
	}

	onRampAddr, err := dsutils.FindAndFormatRef(ds, datastore.AddressRef{
		Type:    datastore.ContractType(onramp.ContractType),
		Version: onramp.Version,
	}, chainSelector, toAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get on ramp address for chain %d: %w", chainSelector, err)
	}

	rmnRemoteAddr, err := dsutils.FindAndFormatRef(ds, datastore.AddressRef{
		Type:    datastore.ContractType(rmn_remote.ContractType),
		Version: rmn_remote.Version,
	}, chainSelector, toAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get rmn remote address for chain %d: %w", chainSelector, err)
	}

	return &ccvadapters.VerifierContractAddresses{
		CommitteeVerifierAddress: committeeVerifierAddr,
		OnRampAddress:            onRampAddr,
		RMNRemoteAddress:         rmnRemoteAddr,
	}, nil
}
