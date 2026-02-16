package token_admin_registry

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"

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

var ProposeAdministrator = contract.NewExercise(contract.ExerciseParams[tokenadminregistry.TokenAdminRegistryProposeAdministrator]{
	Name:         "canton/ccip/token_admin_registry/propose_administrator",
	Version:      Version,
	Description:  "Proposes a new administrator for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input tokenadminregistry.TokenAdminRegistryProposeAdministrator) error {
		return nil
	},
	Template: tokenadminregistry.TokenAdminRegistry{},
	Method:   tokenadminregistry.TokenAdminRegistry{}.TokenAdminRegistryProposeAdministrator,
})

var AcceptAdminRole = contract.NewExercise(contract.ExerciseParams[tokenadminregistry.TokenAdminRegistryAcceptAdminRole]{
	Name:         "canton/ccip/token_admin_registry/accept_admin_role",
	Version:      Version,
	Description:  "Accepts the admin role for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input tokenadminregistry.TokenAdminRegistryAcceptAdminRole) error {
		return nil
	},
	Template: tokenadminregistry.TokenAdminRegistry{},
	Method:   tokenadminregistry.TokenAdminRegistry{}.TokenAdminRegistryAcceptAdminRole,
})

var SetPool = contract.NewExercise(contract.ExerciseParams[tokenadminregistry.TokenAdminRegistrySetPool]{
	Name:         "canton/ccip/token_admin_registry/set_pool",
	Version:      Version,
	Description:  "Sets the token pool owner for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input tokenadminregistry.TokenAdminRegistrySetPool) error {
		return nil
	},
	Template: tokenadminregistry.TokenAdminRegistry{},
	Method:   tokenadminregistry.TokenAdminRegistry{}.TokenAdminRegistrySetPool,
})
