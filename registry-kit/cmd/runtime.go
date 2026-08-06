package cmd

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/registry-kit/config"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/operator"
)

// Runtime holds loaded config, state, and live ledger connections for one command.
type Runtime struct {
	ConfigPath string
	StatePath  string
	Config     config.Config
	State      config.State
}

func loadRuntime(configPath, statePath string) (Runtime, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return Runtime{}, err
	}
	if statePath == "" {
		statePath = config.StatePathNextTo(configPath)
	}
	st, err := config.LoadState(statePath)
	if err != nil {
		return Runtime{}, err
	}

	return Runtime{
		ConfigPath: configPath,
		StatePath:  statePath,
		Config:     cfg,
		State:      st,
	}, nil
}

func (rt *Runtime) saveState() error {
	return rt.State.Save(rt.StatePath)
}

func (rt *Runtime) connect(ctx context.Context, partyRole string) (ledger.Client, canton.Participant, error) {
	party, err := rt.Config.ActingParty(partyRole)
	if err != nil {
		return nil, canton.Participant{}, err
	}
	client, participant, err := ledger.ConnectDevnet(ctx, rt.Config, party)
	if err != nil {
		return nil, canton.Participant{}, fmt.Errorf("connect devnet as %s: %w", party, err)
	}

	return client, participant, nil
}

func (rt *Runtime) operatorClient() *operator.Client {
	return operator.NewClient(rt.Config.Operator.BaseURL)
}
