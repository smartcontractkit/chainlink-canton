package onramp

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("OnRamp")

var Version = semver.MustParse("1.7.0")

var Deploy = contract.NewDeploy(contract.DeployParams[onramp.OnRamp]{
	Name:           "canton/ccip/onramp/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP OnRamp contract on Canton",
	Validate: func(template onramp.OnRamp) error {
		if template.CcipOwner == "" {
			return errors.New("CcipOwner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPOnRamp),
	Prefix:      "onramp",
})
