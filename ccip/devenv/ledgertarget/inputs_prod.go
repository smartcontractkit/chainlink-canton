//go:build prodledger

package ledgertarget

import (
	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	latestmeta "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

func rawInstanceAddressBinding(addr contracts.RawInstanceAddress) chainlinkapi.RawInstanceAddress {
	return chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(addr)}
}

func NewCCVSendInput(ccvAddress contracts.RawInstanceAddress, ccvCid types.CONTRACT_ID, ctx latestmeta.ChoiceContext) CCVSendInput {
	return CCVSendInput{
		CcvAddress:      rawInstanceAddressBinding(ccvAddress),
		CcvCid:          ccvCid,
		CcvExtraContext: AdaptChoiceContext(ctx),
	}
}

func NewCCVExtraArg(ccvAddress contracts.RawInstanceAddress, ccvArgs string) CCVExtraArg {
	return CCVExtraArg{
		CcvAddress: rawInstanceAddressBinding(ccvAddress),
		CcvArgs:    types.TEXT(ccvArgs),
	}
}

func NewExecutorInput(executorCid types.CONTRACT_ID, ctx latestmeta.ChoiceContext) ExecutorInput {
	return ExecutorInput{
		ExecutorCid:          executorCid,
		ExecutorExtraContext: AdaptChoiceContext(ctx),
	}
}

func NewTokenTransferInput(senderInputCids []types.CONTRACT_ID, tokenPoolCid types.CONTRACT_ID, ctx latestmeta.ChoiceContext) TokenTransferInput {
	return TokenTransferInput{
		SenderInputCids:  senderInputCids,
		TokenPoolCid:     tokenPoolCid,
		PoolExtraContext: AdaptChoiceContext(ctx),
	}
}

func NewReceiverCCVInput(ccvCid types.CONTRACT_ID, verifierResults types.TEXT, ctx latestmeta.ChoiceContext) ReceiverCCVInput {
	return ReceiverCCVInput{
		CcvCid:          ccvCid,
		VerifierResults: verifierResults,
		CcvExtraContext: AdaptChoiceContext(ctx),
	}
}

func NewReceiverTokenTransferInput(
	tokenPoolCid types.CONTRACT_ID,
	tokenReceiverParty types.PARTY,
	ctx latestmeta.ChoiceContext,
) ReceiverTokenTransferInput {
	return ReceiverTokenTransferInput{
		TokenPoolCid:       tokenPoolCid,
		TokenReceiverParty: tokenReceiverParty,
		PoolExtraContext:   AdaptChoiceContext(ctx),
	}
}
