package registry

import (
	"context"
	"fmt"

	registryapp "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_app_v0"
	registryv0 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_v0"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
)

// FactoryDiscovery lists Registry service contracts visible to a party.
type FactoryDiscovery struct {
	ProviderService         []ContractRef
	RegistrarService        []ContractRef
	AllocationFactory       []ContractRef
	TransferRule            []ContractRef
	InstrumentConfiguration []ContractRef
}

// DiscoverFactories queries ACS for Registry utility templates.
func DiscoverFactories(ctx context.Context, client ledger.Client, party string) (FactoryDiscovery, error) {
	var out FactoryDiscovery
	var err error

	if out.ProviderService, err = FindContractsByEntity(ctx, client, party, registryapp.ProviderService{}, "ProviderService"); err != nil {
		return FactoryDiscovery{}, fmt.Errorf("provider service: %w", err)
	}
	if out.RegistrarService, err = FindContractsByEntity(ctx, client, party, registryapp.RegistrarService{}, "RegistrarService"); err != nil {
		return FactoryDiscovery{}, fmt.Errorf("registrar service: %w", err)
	}
	if out.AllocationFactory, err = FindContractsByEntity(ctx, client, party, registryapp.AllocationFactory{}, "AllocationFactory"); err != nil {
		return FactoryDiscovery{}, fmt.Errorf("allocation factory: %w", err)
	}
	if out.TransferRule, err = FindContractsByEntity(ctx, client, party, registryv0.TransferRule{}, "TransferRule"); err != nil {
		return FactoryDiscovery{}, fmt.Errorf("transfer rule: %w", err)
	}
	if out.InstrumentConfiguration, err = FindContractsByEntity(ctx, client, party, registryv0.InstrumentConfiguration{}, "InstrumentConfiguration"); err != nil {
		return FactoryDiscovery{}, fmt.Errorf("instrument configuration: %w", err)
	}

	return out, nil
}
