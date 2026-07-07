package onramp

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("OnRamp")

var Version = semver.MustParse("2.0.0")

var Deploy = contract.NewDeploy(contract.DeployParams[ccipruntime.OnRamp]{
	Name:           "canton/ccip/onramp/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP OnRamp contract on Canton",
	Validate: func(template ccipruntime.OnRamp) error {
		if template.CcipOwner == "" {
			return errors.New("CcipOwner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPRuntime),
	Prefix:      "onramp",
})
