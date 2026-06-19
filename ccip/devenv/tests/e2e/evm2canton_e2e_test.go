package canton

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/proxy"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/versioned_verifier_resolver"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	ccldf "github.com/smartcontractkit/chainlink-ccv/build/devenv/cldf"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // register Canton ImplFactory
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

//nolint:paralleltest // we won't run this in parallel.
func TestEVM2Canton_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM2Canton_Basic test in short mode")
	}

	boot := devenvtests.BootstrapE2E(t, devenvtests.ParseEnvFromFlag(t))
	ctx := ccv.Plog.WithContext(t.Context())
	boot.SetupCantonReceive(t, ctx)

	srcSelector := boot.EVM.ChainSelector()
	dstSelector := boot.Canton.ChainSelector()
	_, opsEnv, err := ccldf.NewCLDFOperationsEnvironment(boot.Cfg.Blockchains, boot.Cfg.CLDF.DataStore)
	require.NoError(t, err)
	var receiverParticipant canton.Participant
	if chains := opsEnv.BlockChains.CantonChains(); len(chains[dstSelector].Participants) > 0 {
		receiverParticipant = chains[dstSelector].Participants[0]
	}
	require.NotEmpty(t, receiverParticipant.PartyID)

	receiver, err := boot.Canton.GetEOAReceiverAddress()
	require.NoError(t, err)

	ds, err := boot.Lib.DataStore()
	require.NoError(t, err)
	ccvAddr := devenvtests.GetContractAddress(
		t, ds, srcSelector,
		datastore.ContractType(versioned_verifier_resolver.CommitteeVerifierResolverType),
		versioned_verifier_resolver.Version.String(),
		common.DefaultCommitteeVerifierQualifier,
		"source committee verifier",
	)
	executorAddress := devenvtests.GetContractAddress(
		t, ds, srcSelector,
		datastore.ContractType(sequences.ExecutorProxyType),
		proxy.Deploy.Version(),
		common.DefaultExecutorQualifier,
		"source executor",
	)

	t.Run("message transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		seqNo, err := boot.EVM.GetExpectedNextSequenceNumber(subtestCtx, dstSelector)
		require.NoError(t, err)
		sendMessageResult, err := boot.EVM.SendMessage(subtestCtx, dstSelector, cciptestinterfaces.MessageFields{
			Receiver: receiver,
			Data:     []byte("Hello message transfer from EVM!"),
		}, devenvtests.EVMToCantonMessageOptions(200_000, executorAddress, ccvAddr), 3)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)

		sentEvent, err := boot.EVM.ConfirmSendOnSource(subtestCtx, dstSelector, cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, 15*time.Second)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.Nil(t, sentEvent.Message.TokenTransfer)

		execKey := cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID}
		executionStateChangedEvent, err := boot.Canton.ConfirmExecOnDest(subtestCtx, srcSelector, execKey, devenvtests.ConfirmExecTimeout(t))
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State)

		idempotentEvent, err := boot.Canton.ConfirmExecOnDest(subtestCtx, srcSelector, execKey, devenvtests.ConfirmExecTimeout(t))
		require.NoError(t, err)
		require.Equal(t, executionStateChangedEvent.State, idempotentEvent.State)
		require.Equal(t, executionStateChangedEvent.MessageNumber, idempotentEvent.MessageNumber)
		require.Equal(t, executionStateChangedEvent.SourceChainSelector, idempotentEvent.SourceChainSelector)
	})

	t.Run("token transfer", func(t *testing.T) {
		boot.SkipIfRemote(t, "token e2e not on prod-testnet")

		subtestCtx := ccv.Plog.WithContext(t.Context())

		lane := devenvtests.ResolveTokenLane(t, boot.Cfg, boot.Lib, boot.ChainMap, srcSelector, []uint64{dstSelector})
		srcToken := lane.SrcToken
		srcSender, err := boot.EVM.GetEOAReceiverAddress()
		require.NoError(t, err)
		seqNo, err := boot.EVM.GetExpectedNextSequenceNumber(subtestCtx, dstSelector)
		require.NoError(t, err)
		sendMessageResult, err := boot.EVM.SendMessage(subtestCtx, dstSelector, cciptestinterfaces.MessageFields{
			Receiver: receiver,
			Data:     []byte("Hello token transfer from EVM!"),
			TokenAmount: cciptestinterfaces.TokenAmount{
				Amount:       lane.TransferAmount,
				TokenAddress: srcToken,
			},
		}, cciptestinterfaces.MessageOptions{
			ExecutionGasLimit: lane.ExecutionGasLimit,
			FinalityConfig:    lane.FinalityConfig,
			Executor:          executorAddress,
			CCVs: []protocol.CCV{
				{
					CCVAddress: ccvAddr,
					Args:       []byte{},
					ArgsLen:    0,
				},
			},
		}, 3)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		require.NotNil(t, sendMessageResult.Message.TokenTransfer)

		sentEvent, err := boot.EVM.ConfirmSendOnSource(subtestCtx, dstSelector, cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, 15*time.Second)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.NotNil(t, sentEvent.Message.TokenTransfer)

		result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, boot.Lib, sentEvent.MessageID)
		vr := result.IndexedVerifications.Results[0].VerifierResult
		require.NotNil(t, vr.Message.TokenTransfer)
		require.NotNil(t, vr.Message.TokenTransfer.Amount)
		t.Logf("Canton token transfer amount from verifier result: %s", vr.Message.TokenTransfer.Amount.String())
		require.Positive(t, vr.Message.TokenTransfer.Amount.Cmp(big.NewInt(0)), "token transfer amount must be positive")

		execKey := cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID}
		executionStateChangedEvent, err := boot.Canton.ConfirmExecOnDest(subtestCtx, srcSelector, execKey, devenvtests.ConfirmExecTimeout(t))
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State)

		totalHoldingsRat, err := testhelpers.GetHoldingsBalance(subtestCtx, receiverParticipant, nil)
		require.NoError(t, err)
		totalHoldingsFloat, _ := new(big.Float).SetRat(totalHoldingsRat).Float64()
		t.Logf("Canton receiver total holdings after execute: %.10f", totalHoldingsFloat)

		srcBalanceAfter, err := boot.EVM.GetTokenBalance(subtestCtx, srcSender, srcToken)
		require.NoError(t, err)
		require.NotNil(t, srcBalanceAfter)
		dstBalanceAfter, err := boot.Canton.GetTokenBalance(subtestCtx, receiver, nil)
		require.NoError(t, err)
		require.NotNil(t, dstBalanceAfter)
		t.Logf("Token balances after execute: evm_sender=%s canton_receiver=%s", srcBalanceAfter.String(), dstBalanceAfter.String())
	})
}
