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
func (g *Globals) Resolve(ctx context.Context) (*clients.Bundle, error) {
	cfg, err := cfgpkg.Load(*g.ConfigPath)
	if err != nil {
		return nil, err
	}
	profile, err := cfgpkg.Get(*g.Network)
	if err != nil {
		return nil, err
	}

	return clients.New(ctx, profile, cfg)
}
