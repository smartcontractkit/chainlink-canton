package adapters

import (
	"fmt"

	ccvadapters "github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
)

type CantonIndexerConfigAdapter struct{}

var _ ccvadapters.IndexerConfigAdapter = (*CantonIndexerConfigAdapter)(nil)

func (a *CantonIndexerConfigAdapter) ResolveVerifierAddresses(
	ds datastore.DataStore, chainSelector uint64, qualifier string, kind ccvadapters.VerifierKind,
) ([]string, error) {
	switch kind {
	case ccvadapters.CommitteeVerifierKind:
		return a.resolveCommitteeVerifierAddresses(ds, chainSelector, qualifier)
	case ccvadapters.CCTPVerifierKind, ccvadapters.LombardVerifierKind:
		return nil, &ccvadapters.MissingIndexerVerifierAddressesError{
			Kind:          kind,
			ChainSelector: chainSelector,
			Qualifier:     qualifier,
		}
	default:
		return nil, fmt.Errorf("unknown verifier kind %q", kind)
	}
}

func (a *CantonIndexerConfigAdapter) resolveCommitteeVerifierAddresses(
	ds datastore.DataStore, chainSelector uint64, qualifier string,
) ([]string, error) {
	refs := ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByQualifier(qualifier),
		datastore.AddressRefByType(datastore.ContractType(committee_verifier.ContractType)),
		datastore.AddressRefByVersion(committee_verifier.Version),
	)

	if len(refs) == 0 {
		return nil, &ccvadapters.MissingIndexerVerifierAddressesError{
			Kind:          ccvadapters.CommitteeVerifierKind,
			ChainSelector: chainSelector,
			Qualifier:     qualifier,
		}
	}

	addresses := make([]string, 0, len(refs))
	for _, r := range refs {
		addresses = append(addresses, r.Address)
	}

	return addresses, nil
}
