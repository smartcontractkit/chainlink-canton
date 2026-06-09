package token_admin_registry

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/bind"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("TokenAdminRegistry")

var Version = semver.MustParse("2.0.0")

var tarEncoder = core.NewContract("", "CCIP.TokenAdminRegistry", "TokenAdminRegistry").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[core.TokenAdminRegistry]{
	Name:           "canton/ccip/token_admin_registry/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP TokenAdminRegistry contract on Canton",
	Validate: func(template core.TokenAdminRegistry) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPTokenAdminRegistry),
	Prefix:      "core",
})

var ProposeAdministrator = contract.NewExercise(contract.ExerciseParams[core.ProposeAdministrator]{
	Name:         "canton/ccip/token_admin_registry/propose_administrator",
	Version:      Version,
	Description:  "Proposes a new administrator for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input core.ProposeAdministrator) error {
		return nil
	},
	Template: core.TokenAdminRegistry{},
	Method:   core.TokenAdminRegistry{}.ProposeAdministrator,
	EncodeMethod: func(args core.ProposeAdministrator) (*bind.EncodedChoice, error) {
		return tarEncoder.ProposeAdministratorMCMSParams(core.ProposeAdministratorMCMSParams{
			TokenConfigCid: args.TokenConfigCid,
			InstrumentId:   args.InstrumentId,
			NewAdmin:       args.NewAdmin,
		})
	},
})

var AcceptAdminRole = contract.NewExercise(contract.ExerciseParams[core.AcceptAdminRole]{
	Name:         "canton/ccip/token_admin_registry/accept_admin_role",
	Version:      Version,
	Description:  "Accepts the admin role for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input core.AcceptAdminRole) error {
		return nil
	},
	Template: core.TokenAdminRegistry{},
	Method:   core.TokenAdminRegistry{}.AcceptAdminRole,
	EncodeMethod: func(args core.AcceptAdminRole) (*bind.EncodedChoice, error) {
		return tarEncoder.AcceptAdminRoleMCMSParams(core.AcceptAdminRoleMCMSParams{
			TokenConfigCid: args.TokenConfigCid,
			InstrumentId:   args.InstrumentId,
		})
	},
})

var SetPool = contract.NewExercise(contract.ExerciseParams[core.SetPool]{
	Name:         "canton/ccip/token_admin_registry/set_pool",
	Version:      Version,
	Description:  "Sets the token pool owner for a token in the TokenAdminRegistry",
	ContractType: ContractType,
	Validate: func(input core.SetPool) error {
		return nil
	},
	Template: core.TokenAdminRegistry{},
	Method:   core.TokenAdminRegistry{}.SetPool,
	EncodeMethod: func(args core.SetPool) (*bind.EncodedChoice, error) {
		return tarEncoder.SetPoolMCMSParams(core.SetPoolMCMSParams{
			TokenConfigCid: args.TokenConfigCid,
			InstrumentId:   args.InstrumentId,
			TokenPool:      args.TokenPool,
		})
	},
})
