package ccv_registry

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CCVRegistry")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[common.CCVRegistry]{
	Name:           "canton/ccip/ccv_registry/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP CCV Registry contract on Canton",
	Validate: func(template common.CCVRegistry) error {
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommon),
	Prefix:      "ccvregistry",
})
