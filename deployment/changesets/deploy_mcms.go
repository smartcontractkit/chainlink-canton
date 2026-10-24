package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

type DeployAndConfigureMCMSConfig struct {
	Params sequences.DeployAndConfigureMCMSParams `json:"params" yaml:"params"`
}

type DeployAndConfigureMCMS struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployAndConfigureMCMSConfig]] = DeployAndConfigureMCMS{}

func (d DeployAndConfigureMCMS) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployAndConfigureMCMSConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	params := config.Config.Params
	if params.OwnerParty == "" {
		return fmt.Errorf("owner party is required")
	}
	if params.ChainID <= 0 {
		return fmt.Errorf("chain ID must be greater than zero")
	}

	return nil
}

func (d DeployAndConfigureMCMS) Apply(e cldf.Environment, config CantonCSDeps[DeployAndConfigureMCMSConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.DeployAndConfigureMCMS, chain, config.Config.Params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployAndConfigureMCMS sequence: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{DataStore: ds, Reports: []operations.Report[any, any]{}}, nil
}
