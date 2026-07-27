//go:build prodledger

package ledgertarget

import (
	"errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	rt "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/ccipruntime"
	ccipreceiver "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/receiver"
	ccipsender "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	pprof "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	receiverop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	senderop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var prodCreateRouterEncoder = rt.NewContract("", "CCIP.PerPartyRouter", "PerPartyRouterFactory").Encoder()

var CreateRouterOperation = contract.NewExercise(contract.ExerciseParams[rt.CreateRouter]{
	Name:         "canton/ccip/per_party_router_factory/create_router",
	Version:      pprof.Version,
	Description:  "Creates a new PerPartyRouter using the PerPartyRouterFactory",
	ContractType: pprof.ContractType,
	Validate: func(input rt.CreateRouter) error {
		if input.InstanceId == "" {
			return errors.New("instance ID cannot be empty")
		}
		if input.PartyOwner == "" {
			return errors.New("router owner cannot be empty")
		}

		return nil
	},
	Template:     rt.PerPartyRouterFactory{},
	Method:       rt.PerPartyRouterFactory{}.CreateRouter,
	EncodeMethod: prodCreateRouterEncoder.CreateRouter,
})

var ReceiverDeployOperation = contract.NewDeploy(contract.DeployParams[ccipreceiver.CCIPReceiver]{
	Name:           "canton/ccip/receiver/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(receiverop.ContractType, *receiverop.Version),
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

var SenderDeployOperation = contract.NewDeploy(contract.DeployParams[ccipsender.CCIPSender]{
	Name:           "canton/ccip/sender/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(senderop.ContractType, *senderop.Version),
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

var prodReceiverEncoder = ccipreceiver.NewContract("", "CCIP.CCIPReceiver", "CCIPReceiver").Encoder()

var ReceiverExecuteOperation = contract.NewExercise(contract.ExerciseParams[ccipreceiver.Execute]{
	Name:         "canton/ccip/receiver/execute",
	Version:      receiverop.Version,
	Description:  "Calls the Execute choice on a CCIP Receiver contract",
	ContractType: receiverop.ContractType,
	Validate: func(input ccipreceiver.Execute) error {
		return nil
	},
	Template:     ccipreceiver.CCIPReceiver{},
	Method:       ccipreceiver.CCIPReceiver{}.Execute,
	EncodeMethod: prodReceiverEncoder.Execute,
})

var prodSenderEncoder = ccipsender.NewContract("", "CCIP.CCIPSender", "CCIPSender").Encoder()

var SenderSendOperation = contract.NewExercise(contract.ExerciseParams[ccipsender.Send]{
	Name:         "canton/ccip/sender/send",
	Version:      senderop.Version,
	Description:  "Calls the Send choice on a CCIP Sender contract",
	ContractType: senderop.ContractType,
	Validate: func(input ccipsender.Send) error {
		return nil
	},
	Template:     ccipsender.CCIPSender{},
	Method:       ccipsender.CCIPSender{}.Send,
	EncodeMethod: prodSenderEncoder.Send,
})
