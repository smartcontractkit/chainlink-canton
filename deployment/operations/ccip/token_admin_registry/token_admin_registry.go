package token_admin_registry

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/bind"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/tokenadminregistry"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("TokenAdminRegistry")

var Version = semver.MustParse("0.1.0")

var tarEncoder = tokenadminregistry.NewContract("", "CCIP.TokenAdminRegistry", "TokenAdminRegistry").Encoder()

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

var ProposeAdministrator = contract.NewExercise(contract.ExerciseParams[tokenadminregistry.ProposeAdministrator]{
	Name:         "canton/ccip/token_admin_registry/propose_administrator",
	Version:      Version,
	Description:  "Proposes a new administrator for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input tokenadminregistry.ProposeAdministrator) error {
		return nil
	},
	Template: tokenadminregistry.TokenAdminRegistry{},
	Method:   tokenadminregistry.TokenAdminRegistry{}.ProposeAdministrator,
	EncodeMethod: func(args tokenadminregistry.ProposeAdministrator) (*bind.EncodedChoice, error) {
		return tarEncoder.ProposeAdministratorMCMSParams(tokenadminregistry.ProposeAdministratorMCMSParams{
			TokenConfigCid: args.TokenConfigCid,
			InstrumentId:   args.InstrumentId,
			NewAdmin:       args.NewAdmin,
		})
	},
})

var AcceptAdminRole = contract.NewExercise(contract.ExerciseParams[tokenadminregistry.AcceptAdminRole]{
	Name:         "canton/ccip/token_admin_registry/accept_admin_role",
	Version:      Version,
	Description:  "Accepts the admin role for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input tokenadminregistry.AcceptAdminRole) error {
		return nil
	},
	Template: tokenadminregistry.TokenAdminRegistry{},
	Method:   tokenadminregistry.TokenAdminRegistry{}.AcceptAdminRole,
	EncodeMethod: func(args tokenadminregistry.AcceptAdminRole) (*bind.EncodedChoice, error) {
		return tarEncoder.AcceptAdminRoleMCMSParams(tokenadminregistry.AcceptAdminRoleMCMSParams{
			TokenConfigCid: args.TokenConfigCid,
			InstrumentId:   args.InstrumentId,
		})
	},
})

var SetPool = contract.NewExercise(contract.ExerciseParams[tokenadminregistry.SetPool]{
	Name:         "canton/ccip/token_admin_registry/set_pool",
	Version:      Version,
	Description:  "Sets the token pool owner for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input tokenadminregistry.SetPool) error {
		return nil
	},
	Template: tokenadminregistry.TokenAdminRegistry{},
	Method:   tokenadminregistry.TokenAdminRegistry{}.SetPool,
	EncodeMethod: func(args tokenadminregistry.SetPool) (*bind.EncodedChoice, error) {
		return tarEncoder.SetPoolMCMSParams(tokenadminregistry.SetPoolMCMSParams{
			TokenConfigCid: args.TokenConfigCid,
			InstrumentId:   args.InstrumentId,
			TokenPool:      args.TokenPool,
		})
	},
})
