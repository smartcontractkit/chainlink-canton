//go:build !prodledger

package ledgertarget

import (
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/chainlink/chainlinkapi"
	latestmeta "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
)

func rawInstanceAddressBinding(addr contracts.RawInstanceAddress) chainlinkapi.RawInstanceAddress {
	return chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(addr)}
}

func NewCCVSendInput(ccvAddress contracts.RawInstanceAddress, ccvCid types.CONTRACT_ID, ctx latestmeta.ChoiceContext) CCVSendInput {
	return CCVSendInput{
		CcvAddress: rawInstanceAddressBinding(ccvAddress),
		CcvCid:     ccvCid,
		Context:    AdaptChoiceContext(ctx),
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
		ExecutorCid: executorCid,
		Context:     AdaptChoiceContext(ctx),
	}
}

func NewTokenTransferInput(senderInputCids []types.CONTRACT_ID, tokenPoolCid types.CONTRACT_ID, ctx latestmeta.ChoiceContext) TokenTransferInput {
	return TokenTransferInput{
		SenderInputCids: senderInputCids,
		TokenPoolCid:    tokenPoolCid,
		Context:         AdaptChoiceContext(ctx),
	}
}

func NewReceiverCCVInput(ccvCid types.CONTRACT_ID, verifierResults types.TEXT, ctx latestmeta.ChoiceContext) ReceiverCCVInput {
	return ReceiverCCVInput{
		CcvCid:          ccvCid,
		VerifierResults: verifierResults,
		Context:         AdaptChoiceContext(ctx),
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
		Context:            AdaptChoiceContext(ctx),
	}
}
