//go:build prodledger

package devenv

import (
	ccipreceiver "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var prodReceiverEncoder = ccipreceiver.NewContract("", "CCIP.CCIPReceiver", "CCIPReceiver").Encoder()

var ccipReceiverExecuteOperation = contract.NewExercise(contract.ExerciseParams[ccipreceiver.Execute]{
	Name:         "canton/ccip/receiver/execute",
	Version:      receiver.Version,
	Description:  "Calls the Execute choice on a CCIP Receiver contract",
	ContractType: receiver.ContractType,
	Validate: func(input ccipreceiver.Execute) error {
		return nil
	},
	Template:     ccipreceiver.CCIPReceiver{},
	Method:       ccipreceiver.CCIPReceiver{}.Execute,
	EncodeMethod: prodReceiverEncoder.Execute,
})
