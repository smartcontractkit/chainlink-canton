package mcms

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/mcms"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonMCMS")

var Version = semver.MustParse("1.0.0")

var mcmsEncoder = mcms.NewContract("", "MCMS.Main", "MCMS").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[mcms.MCMS]{
	Name:           "canton/mcms/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the MCMS contract on Canton",
	Validate: func(template mcms.MCMS) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}
		if template.ChainId <= 0 {
			return errors.New("chain ID must be greater than zero")
		}

		return nil
	},
	PackageName: string(contracts.MCMS),
	Prefix:      "mcms",
})

var SetRoot = contract.NewExercise(contract.ExerciseParams[mcms.SetRoot]{
	Name:         "canton/mcms/set_root",
	Version:      Version,
	Description:  "Sets a merkle root for MCMS",
	ContractType: ContractType,
	Validate: func(input mcms.SetRoot) error {
		// TODO add validation

		return nil
	},
	Template:     mcms.MCMS{},
	Method:       mcms.MCMS{}.SetRoot,
	EncodeMethod: mcmsEncoder.SetRoot,
})

var SetConfig = contract.NewExercise(contract.ExerciseParams[mcms.SetConfig]{
	Name:         "canton/mcms/set_config",
	Version:      Version,
	Description:  "Sets configuration for MCMS",
	ContractType: ContractType,
	Validate: func(input mcms.SetConfig) error {
		// TODO add validation

		return nil
	},
	Template:     mcms.MCMS{},
	Method:       mcms.MCMS{}.SetConfig,
	EncodeMethod: mcmsEncoder.SetConfig,
})
