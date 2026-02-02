package token_admin_registry

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("TokenAdminRegistry")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[tokenadminregistry.TokenAdminRegistry]{
	Name:           "canton/ccip/token_admin_registry/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP TokenAdminRegistry contract on Canton",
	Validate: func(template tokenadminregistry.TokenAdminRegistry) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPTokenAdminRegistry),
	Prefix:      "tokenadminregistry",
})
