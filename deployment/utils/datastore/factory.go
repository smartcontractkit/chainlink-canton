package datastore

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
)

// Qualifiers distinguish multiple CCIPFactory instances on the same chain.
const (
	QualifierCCIP  = "ccip"
	QualifierCCV   = "ccv"
	QualifierCore  = "core" // deprecated; use QualifierCCIP
	QualifierPools = "pools"
)

// FactoryAddressRefFromRefs finds a factory ref by qualifier in the given ref list.
func FactoryAddressRefFromRefs(chainSelector uint64, qualifier string, refs []datastore.AddressRef) (datastore.AddressRef, error) {
	for _, ref := range refs {
		if ref.ChainSelector == chainSelector &&
			ref.Type == datastore.ContractType(factoryops.ContractType) &&
			ref.Qualifier == qualifier {
			return validateFactoryAddressRef(ref)
		}
	}

	return datastore.AddressRef{}, fmt.Errorf("missing CCIPFactory ref for chain %d qualifier %q", chainSelector, qualifier)
}

// FactoryAddressRef returns the CCIPFactory address ref for a chain and qualifier.
func FactoryAddressRef(ds datastore.DataStore, chainSelector uint64, qualifier string) (datastore.AddressRef, error) {
	if qualifier == "" {
		return datastore.AddressRef{}, fmt.Errorf("factory qualifier is required")
	}

	ref, err := ds.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(factoryops.ContractType),
		factoryops.Version,
		qualifier,
	))
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("CCIPFactory ref for chain %d qualifier %q: %w", chainSelector, qualifier, err)
	}

	return validateFactoryAddressRef(ref)
}

// FirstCommitteeVerifierRawAddress returns the raw instance address label of any CommitteeVerifier on the chain.
func FirstCommitteeVerifierRawAddress(ds datastore.DataStore, chainSelector uint64, qualifier string) (contracts.RawInstanceAddress, error) {
	if ds == nil {
		return "", fmt.Errorf("datastore is nil")
	}

	refs := ds.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
	for _, ref := range refs {
		if ref.Type != datastore.ContractType(committee_verifier.ContractType) {
			continue
		}
		if qualifier != "" && ref.Qualifier != qualifier {
			continue
		}
		raw, err := GetRawInstanceAddressFromAddressRef(ref)
		if err != nil {
			return "", err
		}

		return raw, nil
	}

	return "", fmt.Errorf("no CommitteeVerifier found for chain %d", chainSelector)
}

// FirstCommitteeVerifierRawAddressFromRefs finds a CommitteeVerifier in the given ref list.
func FirstCommitteeVerifierRawAddressFromRefs(chainSelector uint64, refs []datastore.AddressRef, qualifier string) (contracts.RawInstanceAddress, error) {
	for _, ref := range refs {
		if ref.ChainSelector != chainSelector {
			continue
		}
		if ref.Type != datastore.ContractType(committee_verifier.ContractType) {
			continue
		}
		if qualifier != "" && ref.Qualifier != qualifier {
			continue
		}
		raw, err := GetRawInstanceAddressFromAddressRef(ref)
		if err != nil {
			return "", err
		}

		return raw, nil
	}

	return "", fmt.Errorf("no CommitteeVerifier found for chain %d", chainSelector)
}

// RmnRemoteRawAddress returns the raw instance address of the chain RMNRemote.
func RmnRemoteRawAddress(ds datastore.DataStore, chainSelector uint64) (contracts.RawInstanceAddress, error) {
	if ds == nil {
		return "", fmt.Errorf("datastore is nil")
	}

	ref, err := ds.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(rmn_remote.ContractType),
		rmn_remote.Version,
		"",
	))
	if err != nil {
		return "", fmt.Errorf("RMNRemote for chain %d: %w", chainSelector, err)
	}

	return GetRawInstanceAddressFromAddressRef(ref)
}

func validateFactoryAddressRef(ref datastore.AddressRef) (datastore.AddressRef, error) {
	labels := ref.Labels.List()
	if len(labels) == 0 {
		return datastore.AddressRef{}, fmt.Errorf("CCIPFactory address ref is missing raw instance address label")
	}
	if _, err := contracts.RawInstanceAddressFromString(labels[0]); err != nil {
		return datastore.AddressRef{}, fmt.Errorf("parse CCIPFactory raw instance address label %q: %w", labels[0], err)
	}

	return ref, nil
}
