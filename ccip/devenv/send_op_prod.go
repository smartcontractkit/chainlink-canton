//go:build prodledger

package devenv

import (
	ccipsender "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var prodSenderEncoder = ccipsender.NewContract("", "CCIP.CCIPSender", "CCIPSender").Encoder()

var ccipSenderSendOperation = contract.NewExercise(contract.ExerciseParams[ccipsender.Send]{
	Name:         "canton/ccip/sender/send",
	Version:      sender.Version,
	Description:  "Calls the Send choice on a CCIP Sender contract",
	ContractType: sender.ContractType,
	Validate: func(input ccipsender.Send) error {
		return nil
	},
	Template:     ccipsender.CCIPSender{},
	Method:       ccipsender.CCIPSender{}.Send,
	EncodeMethod: prodSenderEncoder.Send,
})
