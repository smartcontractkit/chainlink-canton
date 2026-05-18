package devenv

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/registry"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/chainconfig"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/committeeverifier"
)

func init() {
	ccv.RegisterLauncher(chain_selectors.FamilyCanton, NewLauncher())
	// Register the canton modifier for the canton family.
	committeeverifier.RegisterModifier(chain_selectors.FamilyCanton, CommitteeVerifierModifier)
	// Register the canton chain config loader for the canton family.
	chainconfig.RegisterChainConfigLoader(chain_selectors.FamilyCanton, CommitteeVerifierConfigLoader)
	// Register the canton impl factory for the canton family.
	ccv.RegisterImplFactory(chain_selectors.FamilyCanton, NewImplFactory())
	// Register the canton CLDF provider factory for the canton family.
	registry.GetGlobalCLDFProviderRegistry().Register(chain_selectors.FamilyCanton, NewCLDF)

	// Register the Canton impl factory so the shared CCV harness can resolve the
	// Canton family to this repo's devenv/test implementation.
	ccv.RegisterImplFactory(chain_selectors.FamilyCanton, NewImplFactory())

	// The other Canton adapters are registered via the init function in the adapters package
}
