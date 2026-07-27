package canton

import (
	"math/big"
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
)

//nolint:paralleltest // we won't run this in parallel.
func TestCanton2EVM_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton2EVM_Basic test in short mode")
	}

	boot := devenvtests.BootstrapE2E(t, devenvtests.ParseEnvFromFlag(t))
	ctx := ccv.Plog.WithContext(t.Context())

	t.Run("EOA receiver and default committee verifier", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		boot.SetupCantonSend(t, ctx, 0)
		receiver := boot.ResolveEVMReceiver(t)

		ds, err := boot.Lib.DataStore()
		require.NoError(t, err)
		ccvAddr := devenvtests.GetContractAddress(
			t, ds, boot.Canton.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		executorAddr := devenvtests.GetContractAddress(
			t, ds, boot.Canton.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)
		t.Logf(
			"Resolved contracts: receiver=%x cantonCCV=%x cantonExecutor=%x srcSelector=%d dstSelector=%d",
			receiver,
			ccvAddr,
			executorAddr,
			boot.Canton.ChainSelector(),
			boot.EVM.ChainSelector(),
		)

		t.Logf("Sending Canton -> EVM message")
		sendMessageResult, err := boot.Canton.SendMessage(
			subtestCtx,
			boot.EVM.ChainSelector(),
			cciptestinterfaces.MessageFields{
				Receiver: receiver,
				Data:     []byte("canton2evm tcapi test"),
			},
			cciptestinterfaces.MessageOptions{
				ExecutionGasLimit: 200_000,
				FinalityConfig:    1,
				Executor:          executorAddr,
				CCVs: []protocol.CCV{
					{
						CCVAddress: ccvAddr,
						Args:       []byte{},
						ArgsLen:    0,
					},
				},
			},
			3,
		)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		t.Logf("SendMessage accepted: seqNo=%d", seqNo)

		t.Logf("Waiting for CCIPMessageSent event: from=%d to=%d seq=%d", boot.Canton.ChainSelector(), boot.EVM.ChainSelector(), seqNo)
		sentEvent, err := boot.Canton.ConfirmSendOnSource(subtestCtx, boot.EVM.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, devenvtests.ConfirmSendTimeout(t, boot.Env))
		require.NoError(t, err)

		t.Logf("CCIPMessageSent event: %+v", sentEvent)

		t.Logf("Asserting message propagated through aggregator/indexer: messageID=%x", sentEvent.MessageID[:])
		result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, boot.Lib, sentEvent.MessageID)
		t.Logf(
			"Message assertion succeeded: aggregated=true indexerResults=%+v",
			result.IndexedVerifications.Results,
		)

		t.Logf("Waiting for execution event on EVM: from=%d seq=%d", boot.Canton.ChainSelector(), seqNo)
		ev, err := boot.EVM.ConfirmExecOnDest(subtestCtx, boot.Canton.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
		require.NoError(t, err)
		assert.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
		t.Logf("Execution event: %+v", ev)
	})

	t.Run("EOA receiver and default committee verifier token transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		lane := devenvtests.ResolveTokenLane(t, boot.Env, boot.Cfg, boot.Lib, boot.ChainMap, boot.Canton.ChainSelector(), []uint64{boot.EVM.ChainSelector()})
		boot.SetupCantonTokenSend(t, ctx, lane, cantondevenv.CantonToEVMTokenSequentialSends)

		receiver := boot.ResolveEVMReceiver(t)

		ds, err := boot.Lib.DataStore()
		require.NoError(t, err)
		ccvAddr := devenvtests.GetContractAddress(
			t, ds, boot.Canton.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		executorAddr := devenvtests.GetContractAddress(
			t, ds, boot.Canton.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)
		destTokenAddress := lane.DestTokenBySelector[boot.EVM.ChainSelector()]
		receiverBalanceBefore, err := boot.EVM.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceBefore)

		for sendIdx := range cantondevenv.CantonToEVMTokenSequentialSends {
			t.Logf("Token transfer send %d/%d", sendIdx+1, cantondevenv.CantonToEVMTokenSequentialSends)
			sendMessageResult, err := boot.Canton.SendMessage(
				subtestCtx,
				boot.EVM.ChainSelector(),
				cciptestinterfaces.MessageFields{
					Receiver: receiver,
					Data:     []byte("canton2evm token transfer"),
					TokenAmount: cciptestinterfaces.TokenAmount{
						Amount: lane.TransferAmount,
					},
				},
				cciptestinterfaces.MessageOptions{
					ExecutionGasLimit: lane.ExecutionGasLimit,
					FinalityConfig:    lane.FinalityConfig,
					Executor:          executorAddr,
					CCVs: []protocol.CCV{
						{
							CCVAddress: ccvAddr,
							Args:       []byte{},
							ArgsLen:    0,
						},
					},
				},
				3,
			)
			require.NoError(t, err)
			require.NotNil(t, sendMessageResult.Message)
			require.NotNil(t, sendMessageResult.Message.TokenTransfer)
			seqNo := uint64(sendMessageResult.Message.SequenceNumber)

			sentEvent, err := boot.Canton.ConfirmSendOnSource(subtestCtx, boot.EVM.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, devenvtests.ConfirmSendTimeout(t, boot.Env))
			require.NoError(t, err)
			require.NotNil(t, sentEvent.Message)
			require.NotNil(t, sentEvent.Message.TokenTransfer)

			t.Logf("Asserting message propagated through aggregator/indexer: messageID=%x", sentEvent.MessageID[:])
			result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, boot.Lib, sentEvent.MessageID)
			t.Logf(
				"Message assertion succeeded: aggregated=true indexerResults=%+v",
				result.IndexedVerifications.Results,
			)

			ev, err := boot.EVM.ConfirmExecOnDest(subtestCtx, boot.Canton.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
			require.NoError(t, err)
			require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
		}

		receiverBalanceAfter, err := boot.EVM.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceAfter)

		expectedTransferPerMessage := new(big.Int).Mul(lane.TransferAmount, big.NewInt(cantondevenv.CantonFixedPointToEVMScale))
		totalExpectedTransfer := new(big.Int).Mul(expectedTransferPerMessage, big.NewInt(cantondevenv.CantonToEVMTokenSequentialSends))
		expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), totalExpectedTransfer)
		t.Logf("EVM receiver token balance: before=%s after=%s totalExpectedTransfer=%s", receiverBalanceBefore.String(), receiverBalanceAfter.String(), totalExpectedTransfer.String())
		require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter)
	})

	// Token transfer with default extraArgs (gasLimit=0, default executor/CCVs from lane) should succeed.
	t.Run("token transfer with default extraArgs", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		lane := devenvtests.ResolveTokenLane(t, boot.Env, boot.Cfg, boot.Lib, boot.ChainMap, boot.Canton.ChainSelector(), []uint64{boot.EVM.ChainSelector()})
		boot.SetupCantonTokenSend(t, ctx, lane, 1)

		receiver := boot.ResolveEVMReceiver(t)

		sendMessageResult, err := boot.Canton.SendMessage(
			subtestCtx,
			boot.EVM.ChainSelector(),
			cciptestinterfaces.MessageFields{
				Receiver: receiver,
				Data:     []byte("canton2evm token transfer default extraArgs"),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount: lane.TransferAmount,
				},
			},
			cciptestinterfaces.MessageOptions{
				ExecutionGasLimit: 0,
				FinalityConfig:    lane.FinalityConfig,
			},
			3,
		)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		require.NotNil(t, sendMessageResult.Message.TokenTransfer)

		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		sentEvent, err := boot.Canton.ConfirmSendOnSource(subtestCtx, boot.EVM.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, devenvtests.ConfirmSendTimeout(t, boot.Env))
		require.NoError(t, err)

		ev, err := boot.EVM.ConfirmExecOnDest(subtestCtx, boot.Canton.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID}, tests.WaitTimeout(t))
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
	})

	// Token transfer with gasLimit=0 should succeed — the EVM receiver gets the tokens
	// even without gas for a callback since there is no receiver contract to call.
	t.Run("token transfer with zero gas limit", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		lane := devenvtests.ResolveTokenLane(t, boot.Env, boot.Cfg, boot.Lib, boot.ChainMap, boot.Canton.ChainSelector(), []uint64{boot.EVM.ChainSelector()})
		boot.SetupCantonTokenSend(t, ctx, lane, 1)

		receiver := boot.ResolveEVMReceiver(t)

		ds, err := boot.Lib.DataStore()
		require.NoError(t, err)
		ccvAddr := devenvtests.GetContractAddress(
			t, ds, boot.Canton.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		executorAddr := devenvtests.GetContractAddress(
			t, ds, boot.Canton.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)

		sendMessageResult, err := boot.Canton.SendMessage(
			subtestCtx,
			boot.EVM.ChainSelector(),
			cciptestinterfaces.MessageFields{
				Receiver: receiver,
				Data:     []byte("canton2evm token transfer zero gas"),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount: lane.TransferAmount,
				},
			},
			cciptestinterfaces.MessageOptions{
				ExecutionGasLimit: 0,
				FinalityConfig:    lane.FinalityConfig,
				Executor:          executorAddr,
				CCVs: []protocol.CCV{
					{
						CCVAddress: ccvAddr,
						Args:       []byte{},
						ArgsLen:    0,
					},
				},
			},
			3,
		)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)

		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		sentEvent, err := boot.Canton.ConfirmSendOnSource(subtestCtx, boot.EVM.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, devenvtests.ConfirmSendTimeout(t, boot.Env))
		require.NoError(t, err)

		ev, err := boot.EVM.ConfirmExecOnDest(subtestCtx, boot.Canton.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID}, tests.WaitTimeout(t))
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
	})
}
