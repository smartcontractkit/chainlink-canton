package global_config

import (
	"errors"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonGlobalConfig")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[common.GlobalConfig]{
	Name:           "canton/ccip/global_config/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP GlobalConfig contract on Canton",
	Validate: func(template common.GlobalConfig) error {
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}
		if (*big.Int)(template.ChainSelector).Cmp(big.NewInt(0)) <= 0 {
			return errors.New("chain selector must be greater than zero")
		}
		// TODO - what's this field?
		// if template.OnRampAddress == "" {
		// 	return errors.New("on ramp address cannot be empty")
		// }

		return nil
	},
	PackageName: string(contracts.CCIPCommon),
	Prefix:      "globalconfig",
})

var UpdateDestChainConfig = contract.NewWrite(contract.WriteParams[common.UpdateDestChainConfig]{
	Name:         "canton/ccip/global_config/update_dest_chain_config",
	Version:      Version,
	Description:  "Updates the GlobalConfig's destination chain configuration",
	ContractType: ContractType,
	Validate: func(input common.UpdateDestChainConfig) error {
		// TODO add validation
		return nil
	},
	Template: common.GlobalConfig{},
	Method:   common.GlobalConfig{}.UpdateDestChainConfig,
})

var UpdateSourceChainConfig = contract.NewWrite(contract.WriteParams[common.UpdateSourceChainConfig]{
	Name:         "canton/ccip/global_config/update_source_chain_config",
	Version:      Version,
	Description:  "Updates the GlobalConfig's source chain configuration",
	ContractType: ContractType,
	Validate: func(input common.UpdateSourceChainConfig) error {
		// TODO add validation
		return nil
	},
	Template: common.GlobalConfig{},
	Method:   common.GlobalConfig{}.UpdateSourceChainConfig,
})
