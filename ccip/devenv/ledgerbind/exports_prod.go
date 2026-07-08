//go:build prodledger

package ledgerbind

import (
	rt "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/core"
	ccipreceiver "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/receiver"
	ccipsender "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/sender"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/chainlink/chainlinkapi"
)

type (
	CreateRouter = rt.CreateRouter
	PerPartyRouter = rt.PerPartyRouter

	RawInstanceAddress = chainlinkapi.RawInstanceAddress

	FinalityConfig = core.FinalityConfig

	ExecutionStateChanged = core.ExecutionStateChanged
	CCIPMessageSent       = core.CCIPMessageSent

	Canton2AnyMessage  = core.Canton2AnyMessage
	ExtraArgs          = core.ExtraArgs
	GenericExtraArgsV3 = core.GenericExtraArgsV3
	ExecutorExtraArg   = core.ExecutorExtraArg
	ExecutorUseDefault = core.ExecutorUseDefault
	TokenTransfer      = core.TokenTransfer
	CCVExtraArg        = core.CCVExtraArg

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

	ReceiverExecute             = ccipreceiver.Execute
	ReceiverCCVInput            = ccipreceiver.CCVInput
	ReceiverTokenTransferInput  = ccipreceiver.TokenTransferInput
)

const (
	MessageExecutionStateUNTOUCHED   = core.MessageExecutionStateUNTOUCHED
	MessageExecutionStateIN_PROGRESS = core.MessageExecutionStateIN_PROGRESS
	MessageExecutionStateSUCCESS     = core.MessageExecutionStateSUCCESS
	MessageExecutionStateFAILURE     = core.MessageExecutionStateFAILURE
)
