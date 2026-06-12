package adapters

import (
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

// CantonCCIPExecutorConfigAdapter exposes Canton executor destination addresses through
// chainlink-ccip's apply-executor-config registry (Canton uses Executor, not ExecutorProxy).
type CantonCCIPExecutorConfigAdapter struct {
	inner CantonExecutorConfigAdapter
}

var _ ccipadapters.ExecutorConfigAdapter = (*CantonCCIPExecutorConfigAdapter)(nil)

func (a *CantonCCIPExecutorConfigAdapter) GetDeployedChains(ds datastore.DataStore, qualifier string) []uint64 {
	return a.inner.GetDeployedChains(ds, qualifier)
}

func (a *CantonCCIPExecutorConfigAdapter) BuildChainConfig(ds datastore.DataStore, chainSelector uint64, qualifier string) (ccipadapters.ExecutorChainConfig, error) {
	cfg, err := a.inner.BuildChainConfig(ds, chainSelector, qualifier)
	if err != nil {
		return ccipadapters.ExecutorChainConfig{}, err
	}

	return ccipadapters.ExecutorChainConfig{
		OffRampAddress:       cfg.OffRampAddress,
		RmnAddress:           cfg.RmnAddress,
		ExecutorProxyAddress: cfg.DefaultExecutorAddress,
	}, nil
}
