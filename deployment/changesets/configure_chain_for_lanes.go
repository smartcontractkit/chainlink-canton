package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	deploymentadapters "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

type ConfigureChainForLanesConfig struct {
	Input sequences.ConfigureChainForLanesInput
}

var _ cldf.ChangeSetV2[CantonCSDeps[ConfigureChainForLanesConfig]] = ConfigureChainForLanes{}

type ConfigureChainForLanes struct{}

func (c ConfigureChainForLanes) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[ConfigureChainForLanesConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (c ConfigureChainForLanes) Apply(e cldf.Environment, config CantonCSDeps[ConfigureChainForLanesConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chainFamily := deploymentadapters.NewCantonChainFamilyAdapter()
	ccipInput := sequences.ToCCIPConfigureChainForLanesInput(config.Config.Input)
	out, err := operations.ExecuteSequence(e.OperationsBundle, chainFamily.ConfigureChainForLanes(), e.BlockChains, ccipInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployChainContracts sequence: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []operations.Report[any, any]{},
	}, nil
}
