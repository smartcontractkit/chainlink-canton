package main

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cli"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/chainconfig"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/committeeverifier"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
)

func init() {
	ccv.RegisterLauncher(chain_selectors.FamilyCanton, cantondevenv.NewLauncher())
	// Register the canton modifier for the canton family.
	committeeverifier.RegisterModifier(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierModifier)
	// Register the canton chain config loader for the canton family.
	chainconfig.RegisterChainConfigLoader(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierConfigLoader)
	// Register the canton impl factory for the canton family.
	ccv.RegisterImplFactory(chain_selectors.FamilyCanton, cantondevenv.NewImplFactory())

	// The other Canton adapters are registered via the init function in the adapters package
}

func main() {
	cli.RunCLI()
}
