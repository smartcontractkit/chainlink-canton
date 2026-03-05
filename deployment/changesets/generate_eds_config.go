package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"

	"github.com/smartcontractkit/chainlink-canton/deployment"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/services/eds"
)

type GenerateEDSConfig struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[edsConfig.ServerConfig]] = GenerateEDSConfig{}

func (c GenerateEDSConfig) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[edsConfig.ServerConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (c GenerateEDSConfig) Apply(e cldf.Environment, config CantonCSDeps[edsConfig.ServerConfig]) (cldf.ChangesetOutput, error) {
	deps := eds.BuildConfigDeps{
		Env: e,
	}

	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, eds.BuildConfig, deps, eds.GenerateEDSConfigInput{
		ChainSelector: config.ChainSelector,
		Participant:   config.Participant,
		ServerConfig:  config.Config,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute BuildConfig operation: %w", err)
	}

	ds := datastore.NewMemoryDataStore()
	if e.DataStore != nil {
		if err := ds.Merge(e.DataStore); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge existing datastore: %w", err)
		}
	}

	edsCfg := out.Output.Config
	if err := deployment.SaveEDSConfig(ds, edsCfg); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save EDS config: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
	}, nil
}
