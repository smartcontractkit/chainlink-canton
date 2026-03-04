package main

import (
	"github.com/Masterminds/semver/v3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	evmadapters "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/adapters"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cli"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/registry"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/chainconfig"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/committeeverifier"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	cantonadapters "github.com/smartcontractkit/chainlink-canton/ccip/devenv/adapters"
)

// TODO: should this be defined elsewhere?
var tokenPoolVersions = []string{
	"1.6.1",
	"1.7.0",
}

func init() {
	// Register the canton modifier for the canton family.
	committeeverifier.RegisterModifier(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierModifier)
	// Register the canton chain config loader for the canton family.
	chainconfig.RegisterChainConfigLoader(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierConfigLoader)
	// Register the canton chain family adapter for the canton family.
	registry.RegisterChainFamilyAdapter(chain_selectors.FamilyCanton, cantonadapters.NewChainFamilyAdapter(&evmadapters.ChainFamilyAdapter{}))
	// Register the canton impl factory for the canton family.
	registry.RegisterImplFactory(chain_selectors.FamilyCanton, cantondevenv.NewImplFactory())
	// Register the canton token adapter for the canton family.
	for _, version := range tokenPoolVersions {
		tokenscore.GetTokenAdapterRegistry().RegisterTokenAdapter(chain_selectors.FamilyCanton, semver.MustParse(version), cantonadapters.NewTokenAdapter(&evmadapters.TokenAdapter{}))
	}
}

func main() {
	cli.RunCLI()
}
