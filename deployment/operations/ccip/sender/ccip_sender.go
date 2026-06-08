package sender

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonCCIPSender")

var Version = semver.MustParse("0.1.0")

var senderEncoder = sender.NewContract("", "CCIP.CCIPSender", "CCIPSender").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[sender.CCIPSender]{
	Name:           "canton/ccip/sender/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP Sender contract on Canton",
	Validate: func(template sender.CCIPSender) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPSender),
	Prefix:      "ccipsender",
})

var Send = contract.NewExercise(contract.ExerciseParams[sender.Send]{
	Name:         "canton/ccip/sender/send",
	Version:      Version,
	Description:  "Calls the Send choice on a CCIP Sender contract",
	ContractType: ContractType,
	Validate: func(input sender.Send) error {
		return nil
	},
	Template:     sender.CCIPSender{},
	Method:       sender.CCIPSender{}.Send,
	EncodeMethod: senderEncoder.Send,
})
