package rmn_remote

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("RMNRemote")

var Version = semver.MustParse("1.6.0")

var Deploy = contract.NewDeploy(contract.DeployParams[rmn.RMNRemote]{
	Name:           "canton/ccip/rmn_remote/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the RMN Remote contract on Canton",
	Validate: func(template rmn.RMNRemote) error {
		if template.CcipOwner == "" {
			return errors.New("CcipOwner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPRMN),
	Prefix:      "rmn_remote",
})

var Curse = contract.NewExercise(contract.ExerciseParams[rmn.Curse]{
	Name:         "canton/ccip/rmn_remote/curse",
	Version:      Version,
	Description:  "Curses a subject on the RMNRemote contract",
	ContractType: ContractType,
	Validate: func(input rmn.Curse) error {
		if input.Subject == "" {
			return errors.New("Subject cannot be empty")
		}
		return nil
	},
	Template: rmn.RMNRemote{},
	Method:   rmn.RMNRemote{}.Curse,
})

var Uncurse = contract.NewExercise(contract.ExerciseParams[rmn.Uncurse]{
	Name:         "canton/ccip/rmn_remote/uncurse",
	Version:      Version,
	Description:  "Uncurses a subject on the RMNRemote contract",
	ContractType: ContractType,
	Validate: func(input rmn.Uncurse) error {
		if input.Subject == "" {
			return errors.New("Subject cannot be empty")
		}
		return nil
	},
	Template: rmn.RMNRemote{},
	Method:   rmn.RMNRemote{}.Uncurse,
})
