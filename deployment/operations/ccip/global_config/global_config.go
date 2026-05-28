package global_config

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/common"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	internalparse "github.com/smartcontractkit/chainlink-canton/internal/parse"
)

var ContractType = deployment.ContractType("CantonGlobalConfig")

var Version = semver.MustParse("0.1.0")

var globalConfigEncoder = common.NewContract("", "CCIP.GlobalConfig", "GlobalConfig").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[common.GlobalConfig]{
	Name:           "canton/ccip/global_config/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP GlobalConfig contract on Canton",
	Validate: func(template common.GlobalConfig) error {
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
	PackageName: string(contracts.CCIPCommon),
	Prefix:      "globalconfig",
})

var ApplyDestChainConfigUpdates = contract.NewExercise(contract.ExerciseParams[common.ApplyDestChainConfigUpdates]{
	Name:         "canton/ccip/global_config/apply_dest_chain_config_updates",
	Version:      Version,
	Description:  "Updates the GlobalConfig's destination chain configuration",
	ContractType: ContractType,
	Validate: func(input common.ApplyDestChainConfigUpdates) error {
		// TODO add validation
		return nil
	},
	Template:     common.GlobalConfig{},
	Method:       common.GlobalConfig{}.ApplyDestChainConfigUpdates,
	EncodeMethod: globalConfigEncoder.ApplyDestChainConfigUpdates,
})

var ApplySourceChainConfigUpdates = contract.NewExercise(contract.ExerciseParams[common.ApplySourceChainConfigUpdates]{
	Name:         "canton/ccip/global_config/apply_source_chain_config_updates",
	Version:      Version,
	Description:  "Updates the GlobalConfig's source chain configuration",
	ContractType: ContractType,
	Validate: func(input common.ApplySourceChainConfigUpdates) error {
		// TODO add validation
		return nil
	},
	Template:     common.GlobalConfig{},
	Method:       common.GlobalConfig{}.ApplySourceChainConfigUpdates,
	EncodeMethod: globalConfigEncoder.ApplySourceChainConfigUpdates,
})
