package coin

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/coin"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonCoinRegistry")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[coin.CoinRegistry]{
	Name:           "canton/coin/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CoinRegistry contract on Canton",
	Validate: func(template coin.CoinRegistry) error {
		if template.Issuer == "" {
			return errors.New("issuer cannot be empty")
		}
		if template.InstrumentId == (coin.InstrumentId{}) {
			return errors.New("instrument ID cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.Coin),
	Prefix:      "coinregistry", // TODO might want to make this configurable
})
