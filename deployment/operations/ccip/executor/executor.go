package executor

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("Executor")

var Version = semver.MustParse("0.1.0")

var executorEncoder = executor.NewContract("", "CCIP.Executor", "Executor").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[executor.Executor]{
	Name:           "canton/ccip/executor/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP Executor contract on Canton",
	Validate: func(template executor.Executor) error {
		if template.Owner == "" {
			return fmt.Errorf("owner cannot be empty")
		}
		if len(template.RemoteChainConfigs) != 0 {
			return fmt.Errorf("remote chain configs should not be set during deployment")
		}

		return nil
	},
	PackageName: string(contracts.CCIPExecutor),
	Prefix:      "executor",
})

var ApplyDestChainUpdates = contract.NewExercise(contract.ExerciseParams[executor.ApplyDestChainUpdates]{
	Name:         "canton/ccip/executor/apply_dest_chain_updates",
	Version:      Version,
	Description:  "Applies dest chain config updates to a Canton Executor",
	ContractType: ContractType,
	Validate: func(input executor.ApplyDestChainUpdates) error {

		return nil
	},
	Template:     executor.Executor{},
	Method:       executor.Executor{}.ApplyDestChainUpdates,
	EncodeMethod: executorEncoder.ApplyDestChainUpdates,
})

var SetDynamicConfig = contract.NewExercise(contract.ExerciseParams[executor.SetDynamicConfig]{
	Name:         "canton/ccip/executor/set_dynamic_config",
	Version:      Version,
	Description:  "Updates the dynamic config of a Canton Executor",
	ContractType: ContractType,
	Validate: func(input executor.SetDynamicConfig) error {

		return nil
	},
	Template:     executor.Executor{},
	Method:       executor.Executor{}.SetDynamicConfig,
	EncodeMethod: executorEncoder.SetDynamicConfig,
})

var ApplyAllowedCCVUpdates = contract.NewExercise(contract.ExerciseParams[executor.ApplyAllowedCCVUpdates]{
	Name:         "canton/ccip/executor/apply_allowed_ccv_updates",
	Version:      Version,
	Description:  "Applies allowed CCV updates to a Canton Executor",
	ContractType: ContractType,
	Validate: func(input executor.ApplyAllowedCCVUpdates) error {

		return nil
	},
	Template:     executor.Executor{},
	Method:       executor.Executor{}.ApplyAllowedCCVUpdates,
	EncodeMethod: executorEncoder.ApplyAllowedCCVUpdates,
})
