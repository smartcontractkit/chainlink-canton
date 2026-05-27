package devenv

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/chainreg"
)

func init() {
	if err := chainreg.Register(chain_selectors.FamilyCanton, chainreg.Registration{
		ImplFactory:       NewImplFactory(),
		CLDFProvider:      NewCLDF,
		ChainConfigLoader: CommitteeVerifierConfigLoader,
		Launcher:          NewLauncher(),
		VerifierModifier:  CommitteeVerifierModifier,
	}); err != nil {
		panic("canton chainreg: " + err.Error())
	}

	// The other Canton adapters are registered via the init function in the adapters package
}
