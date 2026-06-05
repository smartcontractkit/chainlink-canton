package rmn_remote

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("RMNRemote")

var Version = semver.MustParse("1.6.0")

var rmnEncoder = core.NewContract("", "CCIP.RMNRemote", "RMNRemote").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[core.RMNRemote]{
	Name:           "canton/ccip/rmn_remote/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the RMN Remote contract on Canton",
	Validate: func(template core.RMNRemote) error {
		if template.CcipOwner == "" {
			return errors.New("CcipOwner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPRMN),
	Prefix:      "rmn_remote",
})

var Curse = contract.NewExercise(contract.ExerciseParams[core.Curse]{
	Name:         "canton/ccip/rmn_remote/curse",
	Version:      Version,
	Description:  "Curses a subject on the RMNRemote contract",
	ContractType: ContractType,
	Validate: func(input core.Curse) error {
		if input.Subject == "" {
			return errors.New("subject cannot be empty")
		}

		return nil
	},
	Template:     core.RMNRemote{},
	Method:       core.RMNRemote{}.Curse,
	EncodeMethod: rmnEncoder.Curse,
})

var Uncurse = contract.NewExercise(contract.ExerciseParams[core.Uncurse]{
	Name:         "canton/ccip/rmn_remote/uncurse",
	Version:      Version,
	Description:  "Uncurses a subject on the RMNRemote contract",
	ContractType: ContractType,
	Validate: func(input core.Uncurse) error {
		if input.Subject == "" {
			return errors.New("subject cannot be empty")
		}

		return nil
	},
	Template:     core.RMNRemote{},
	Method:       core.RMNRemote{}.Uncurse,
	EncodeMethod: rmnEncoder.Uncurse,
})
