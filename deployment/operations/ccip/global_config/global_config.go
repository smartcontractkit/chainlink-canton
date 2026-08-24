package global_config

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	internalparse "github.com/smartcontractkit/chainlink-canton/internal/parse"
)

var ContractType = deployment.ContractType("CantonGlobalConfig")

var Version = semver.MustParse("2.0.0")

var globalConfigEncoder = core.NewContract("", "CCIP.CoreV2.GlobalConfig", "GlobalConfig").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[core.GlobalConfig]{
	Name:           "canton/ccip/global_config/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP GlobalConfig contract on Canton",
	Validate: func(template core.GlobalConfig) error {
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}

		chainSelector, err := internalparse.Uint64Checked(string(template.ChainSelector))
		if err != nil {
			return err
		}

		if chainSelector == 0 {
			return errors.New("chain selector must be greater than zero")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCoreV2),
	Prefix:      "globalconfig",
})

var ApplyDestChainConfigUpdates = contract.NewExercise(contract.ExerciseParams[core.ApplyDestChainConfigUpdates]{
	Name:         "canton/ccip/global_config/apply_dest_chain_config_updates",
	Version:      Version,
	Description:  "Updates the GlobalConfig's destination chain configuration",
	ContractType: ContractType,
	Validate: func(input core.ApplyDestChainConfigUpdates) error {
		// TODO add validation
		return nil
	},
	Template:     core.GlobalConfig{},
	Method:       core.GlobalConfig{}.ApplyDestChainConfigUpdates,
	EncodeMethod: globalConfigEncoder.ApplyDestChainConfigUpdates,
})

var ApplySourceChainConfigUpdates = contract.NewExercise(contract.ExerciseParams[core.ApplySourceChainConfigUpdates]{
	Name:         "canton/ccip/global_config/apply_source_chain_config_updates",
	Version:      Version,
	Description:  "Updates the GlobalConfig's source chain configuration",
	ContractType: ContractType,
	Validate: func(input core.ApplySourceChainConfigUpdates) error {
		// TODO add validation
		return nil
	},
	Template:     core.GlobalConfig{},
	Method:       core.GlobalConfig{}.ApplySourceChainConfigUpdates,
	EncodeMethod: globalConfigEncoder.ApplySourceChainConfigUpdates,
})
