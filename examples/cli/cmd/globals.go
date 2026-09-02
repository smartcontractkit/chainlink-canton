// Package cmd holds the Cobra commands for the CCIP demo CLI.
package cmd

import (
	"context"

	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/clients"
	cfgpkg "github.com/smartcontractkit/chainlink-canton/examples/cli/internal/config"
)

// Globals are flags shared by every subcommand.
type Globals struct {
	ConfigPath *string
	Network    *string
}

// Resolve loads the config + profile and constructs the client bundle.
// Can optionally disable Canton if not required.
func (g *Globals) Resolve(ctx context.Context, disableCanton bool) (*clients.Bundle, error) {
	cfg, err := cfgpkg.Load(*g.ConfigPath)
	if err != nil {
		return nil, err
	}
	profile, err := cfgpkg.Get(*g.Network)
	if err != nil {
		return nil, err
	}

	if disableCanton {
		cfg.Canton.Disabled = true
	}

	return clients.New(ctx, profile, cfg)
}
