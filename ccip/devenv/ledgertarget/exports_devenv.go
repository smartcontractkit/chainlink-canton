//go:build !prodledger

package ledgertarget

import (
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	rt "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/clientapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/events"
	ccipreceiver "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/receiver"
	ccipsender "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
)

type (
	CreateRouter   = rt.CreateRouter
	PerPartyRouter = rt.PerPartyRouter

	RawInstanceAddress = chainlinkapi.RawInstanceAddress

	FinalityConfig = ccipcodec.FinalityConfig

	ExecutionStateChanged = events.ExecutionStateChanged
	CCIPMessageSent       = events.CCIPMessageSent

	Canton2AnyMessage  = clientapi.Canton2AnyMessage
	ExtraArgs          = clientapi.ExtraArgs
	GenericExtraArgsV3 = clientapi.GenericExtraArgsV3
	ExecutorExtraArg   = clientapi.ExecutorExtraArg
	ExecutorUseDefault = clientapi.ExecutorUseDefault
	TokenTransfer      = clientapi.TokenTransfer
	CCVExtraArg        = clientapi.CCVExtraArg

	ApplyPriceUpdatersUpdate = core.ApplyPriceUpdatersUpdate
	UpdatePrices             = core.UpdatePrices
	PriceUpdates             = core.PriceUpdates
	TokenPriceUpdate         = core.TokenPriceUpdate
	FeeQuoter                = core.FeeQuoter
	FeeQuoterDestChainConfig = core.FeeQuoterDestChainConfig

	CCIPReceiver       = ccipreceiver.CCIPReceiver
	CCIPSender         = ccipsender.CCIPSender
	Send               = ccipsender.Send
	FeeTokenInput      = ccipsender.FeeTokenInput
	TokenTransferInput = ccipsender.TokenTransferInput
	CCVSendInput       = ccipsender.CCVSendInput
	ExecutorInput      = ccipsender.ExecutorInput

	ReceiverExecute            = ccipreceiver.Execute
	ReceiverCCVInput           = ccipreceiver.CCVInput
	ReceiverTokenTransferInput = ccipreceiver.TokenTransferInput
)

const (
	MessageExecutionStateUNTOUCHED   = ccipapi.MessageExecutionStateUNTOUCHED
	MessageExecutionStateIN_PROGRESS = ccipapi.MessageExecutionStateIN_PROGRESS
	MessageExecutionStateSUCCESS     = ccipapi.MessageExecutionStateSUCCESS
	MessageExecutionStateFAILURE     = ccipapi.MessageExecutionStateFAILURE
)
