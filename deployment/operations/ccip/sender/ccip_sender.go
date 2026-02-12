package sender

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/ccipsender"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonCCIPSender")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[ccipsender.CCIPSender]{
	Name:           "canton/ccip/sender/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP Sender contract on Canton",
	Validate: func(template ccipsender.CCIPSender) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPSender),
	Prefix:      "ccipsender",
})

var Send = contract.NewExercise(contract.ExerciseParams[ccipsender.Send]{
	Name:         "canton/ccip/sender/send",
	Version:      Version,
	Description:  "Calls the Send choice on a CCIP Sender contract",
	ContractType: ContractType,
	Validate: func(input ccipsender.Send) error {
		return nil
	},
	Template: ccipsender.CCIPSender{},
	Method:   ccipsender.CCIPSender{}.Send,
})
