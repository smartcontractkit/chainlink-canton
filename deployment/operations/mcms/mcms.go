package mcms

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/core"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonMCMS")

var Version = semver.MustParse("2.0.0")

var mcmsEncoder = mcmsCore.NewContract("", "MCMS.Main", "MCMS").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[mcmsCore.MCMS]{
	Name:           "canton/mcms/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the MCMS contract on Canton",
	Validate: func(template mcmsCore.MCMS) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}
		if template.ChainId <= 0 {
			return errors.New("chain ID must be greater than zero")
		}

		return nil
	},
	PackageName: string(contracts.MCMSCore),
	Prefix:      "mcms",
})

var SetRoot = contract.NewExercise(contract.ExerciseParams[mcmsCore.SetRoot]{
	Name:         "canton/mcms/set_root",
	Version:      Version,
	Description:  "Sets a merkle root for MCMS",
	ContractType: ContractType,
	Validate: func(input mcmsCore.SetRoot) error {
		// TODO add validation

		return nil
	},
	Template:     mcmsCore.MCMS{},
	Method:       mcmsCore.MCMS{}.SetRoot,
	EncodeMethod: mcmsEncoder.SetRoot,
})

var SetConfig = contract.NewExercise(contract.ExerciseParams[mcmsCore.SetConfig]{
	Name:         "canton/mcms/set_config",
	Version:      Version,
	Description:  "Sets configuration for MCMS",
	ContractType: ContractType,
	Validate: func(input mcmsCore.SetConfig) error {
		// TODO add validation

		return nil
	},
	Template:     mcmsCore.MCMS{},
	Method:       mcmsCore.MCMS{}.SetConfig,
	EncodeMethod: mcmsEncoder.SetConfig,
})
