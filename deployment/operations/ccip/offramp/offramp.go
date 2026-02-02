package offramp

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("OffRamp")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[offramp.OffRamp]{
	Name:           "canton/ccip/offramp/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP OffRamp contract on Canton",
	Validate: func(template offramp.OffRamp) error {
		if template.CcipOwner == "" {
			return errors.New("CcipOwner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPOffRamp),
	Prefix:      "offramp",
})
