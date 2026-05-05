package adapters

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	dsutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	ccvadapters "github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	ccvexecutor "github.com/smartcontractkit/chainlink-ccv/executor"
	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
)

type CantonExecutorConfigAdapter struct{}

var _ ccvadapters.ExecutorConfigAdapter = (*CantonExecutorConfigAdapter)(nil)

func (a *CantonExecutorConfigAdapter) GetDeployedChains(ds datastore.DataStore, qualifier string) []uint64 {
	if ds == nil {
		return nil
	}
	refs := ds.Addresses().Filter(
		datastore.AddressRefByQualifier(qualifier),
		datastore.AddressRefByType(datastore.ContractType(executor.ContractType)),
	)
	seen := make(map[uint64]struct{}, len(refs))
	chains := make([]uint64, 0, len(refs))
	for _, ref := range refs {
		family, err := chainsel.GetSelectorFamily(ref.ChainSelector)
		if err != nil || family != chainsel.FamilyCanton {
			continue
		}
		if _, exists := seen[ref.ChainSelector]; !exists {
			seen[ref.ChainSelector] = struct{}{}
			chains = append(chains, ref.ChainSelector)
		}
	}

	return chains
}

func (a *CantonExecutorConfigAdapter) BuildChainConfig(ds datastore.DataStore, chainSelector uint64, qualifier string) (ccvexecutor.ChainConfiguration, error) {
	toAddress := func(ref datastore.AddressRef) (string, error) { return ref.Address, nil }

	offRampAddr, err := dsutils.FindAndFormatRef(ds, datastore.AddressRef{
		Type:    datastore.ContractType(offramp.ContractType),
		Version: offramp.Version,
	}, chainSelector, toAddress)
	if err != nil {
		return ccvexecutor.ChainConfiguration{}, fmt.Errorf("failed to get off ramp address for chain %d: %w", chainSelector, err)
	}

	rmnRemoteAddr, err := dsutils.FindAndFormatRef(ds, datastore.AddressRef{
		Type:    datastore.ContractType(rmn_remote.ContractType),
		Version: rmn_remote.Version,
	}, chainSelector, toAddress)
	if err != nil {
		return ccvexecutor.ChainConfiguration{}, fmt.Errorf("failed to get rmn remote address for chain %d: %w", chainSelector, err)
	}

	executorAddr, err := dsutils.FindAndFormatRef(ds, datastore.AddressRef{
		Type:      datastore.ContractType(executor.ContractType),
		Qualifier: qualifier,
		Version:   executor.Version,
	}, chainSelector, toAddress)
	if err != nil {
		return ccvexecutor.ChainConfiguration{}, fmt.Errorf("failed to get executor address for chain %d: %w", chainSelector, err)
	}

	return ccvexecutor.ChainConfiguration{
		DestinationChainConfig: chainaccess.DestinationChainConfig{
			OffRampAddress: offRampAddr,
			RmnAddress:     rmnRemoteAddr,
		},
		DefaultExecutorAddress: executorAddr,
	}, nil
}
