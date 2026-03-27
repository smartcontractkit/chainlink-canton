package main

import (
	"github.com/Masterminds/semver/v3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cli"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/chainconfig"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/committeeverifier"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	"github.com/smartcontractkit/chainlink-canton/deployment/adapters"
)

// TODO: should this be defined elsewhere?
var tokenPoolVersions = []string{
	"1.6.1",
	"2.0.0",
}

func init() {
	ccv.RegisterLauncher(chain_selectors.FamilyCanton, cantondevenv.NewLauncher())
	// Register the canton modifier for the canton family.
	committeeverifier.RegisterModifier(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierModifier)
	// Register the canton chain config loader for the canton family.
	chainconfig.RegisterChainConfigLoader(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierConfigLoader)
	// Register the canton impl factory for the canton family.
	ccv.RegisterImplFactory(chain_selectors.FamilyCanton, cantondevenv.NewImplFactory())

	// The other Canton adapters are registered via the init function in the adapters package

	// Register the canton token adapter for the canton family.
	for _, version := range tokenPoolVersions {
		tokenscore.GetTokenAdapterRegistry().RegisterTokenAdapter(chain_selectors.FamilyCanton, semver.MustParse(version), adapters.CantonTokenAdapter{})
	}
}

func main() {
	cli.RunCLI()
}
