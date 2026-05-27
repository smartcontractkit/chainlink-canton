package receiver

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonCCIPReceiver")

var Version = semver.MustParse("0.1.0")

var receiverEncoder = ccipreceiver.NewContract("", "CCIP.CCIPReceiver", "CCIPReceiver").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[ccipreceiver.CCIPReceiver]{
	Name:           "canton/ccip/receiver/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP Receiver contract on Canton",
	Validate: func(template ccipreceiver.CCIPReceiver) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPReceiver),
	Prefix:      "ccipreceiver",
})

var Execute = contract.NewExercise(contract.ExerciseParams[ccipreceiver.Execute]{
	Name:         "canton/ccip/receiver/execute",
	Version:      Version,
	Description:  "Calls the Execute choice on a CCIP Receiver contract",
	ContractType: ContractType,
	Validate: func(input ccipreceiver.Execute) error {
		return nil
	},
	Template:     ccipreceiver.CCIPReceiver{},
	Method:       ccipreceiver.CCIPReceiver{}.Execute,
	EncodeMethod: receiverEncoder.Execute,
})
