//go:build disabled

package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

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

	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.ConfigureChainForLanes, chain, config.Config.Input)
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
