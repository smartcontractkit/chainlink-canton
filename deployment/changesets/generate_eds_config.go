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

// GenerateEDSConfigConfig is the input for GenerateEDSConfig.
// edsConfig.Config carries server/monitoring settings and optional per-pool overlays.
// LockReleaseTransferPreapproval applies to all discovered lockRelease pools when set.
type GenerateEDSConfigConfig struct {
	edsConfig.Config
	LockReleaseTransferPreapproval *edsConfig.TransferPreapproval `json:"lock_release_transfer_preapproval,omitempty" toml:"lock_release_transfer_preapproval,omitempty"`
}

type GenerateEDSConfig struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[GenerateEDSConfigConfig]] = GenerateEDSConfig{}

func (c GenerateEDSConfig) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[GenerateEDSConfigConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (c GenerateEDSConfig) Apply(e cldf.Environment, config CantonCSDeps[GenerateEDSConfigConfig]) (cldf.ChangesetOutput, error) {
	deps := eds.BuildConfigDeps{
		Env: e,
	}

	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, eds.BuildConfig, deps, eds.GenerateEDSConfigInput{
		ChainSelector:                  config.ChainSelector,
		Participant:                    config.Participant,
		LockReleaseTransferPreapproval: config.Config.LockReleaseTransferPreapproval,
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

	config.Config.Node = out.Output.NodeConfig
	config.Config.CCIPAPIConfig = out.Output.CCIPAPIConfig
	config.Config.CCVAPIConfig = out.Output.CCVAPIConfig
	config.Config.ExecutorAPIConfig = out.Output.ExecutorAPIConfig

	generatedPools := &edsConfig.Config{
		TokenPoolAPIConfig: out.Output.TokenPoolAPIConfig,
	}
	mergedPools, err := generatedPools.Merge(&edsConfig.Config{
		TokenPoolAPIConfig: config.Config.TokenPoolAPIConfig,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge token pool EDS config overlay: %w", err)
	}
	config.Config.TokenPoolAPIConfig = mergedPools.TokenPoolAPIConfig

	if err := deployment.SaveEDSConfig(ds, &config.Config.Config); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save EDS config: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
	}, nil
}
