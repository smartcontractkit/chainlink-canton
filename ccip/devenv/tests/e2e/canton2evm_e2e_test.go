package canton

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
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

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	lib, err := ccv.NewLibFromCCVEnv(&ccv.Plog, configPath)
	require.NoError(t, err)
	ctx := ccv.Plog.WithContext(t.Context())

	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)
	require.NoError(t, devenvtests.WireVerifierObservationFromLib(lib, chainMap))

	evmChain := devenvtests.GetChainFromMap(t, blockchain.TypeAnvil, in, chainMap)
	cantonChain := devenvtests.GetChainFromMap(t, blockchain.TypeCanton, in, chainMap)
	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	require.True(t, ok, "Canton chain cantonImpl must be *devenv.Chain")

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	// TODO: currently minting 2 holdings of 2k for all the e2e tests. Otherwise one holding might conflict with the other.
	// the tests need to be hardened to not rely on this and instead correctly pick the holding that suffices.
	require.NoError(t, cantonImpl.MintTokens(ctx, uint64(devenvtests.CantonToEVMFeeAmount)*100))
	require.NoError(t, cantonImpl.MintTokens(ctx, uint64(devenvtests.CantonToEVMFeeAmount)*100))

	t.Run("EOA receiver and default committee verifier", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		// Setup message send
		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), 0))

		ds, err := lib.DataStore()
		require.NoError(t, err)
		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr := devenvtests.GetContractAddress(
			t, ds, cantonChain.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		executorAddr := devenvtests.GetContractAddress(
			t, ds, cantonChain.ChainSelector(),
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
			cantonChain.ChainSelector(),
			evmChain.ChainSelector(),
		)

		t.Logf("Sending Canton -> EVM message")
		sendMessageResult, err := cantonChain.SendMessage(
			subtestCtx,
			evmChain.ChainSelector(),
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
		// require.NotEmpty(t, sendMessageResult.ReceiptIssuers)
		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		t.Logf(
			"SendMessage accepted: seqNo=%d",
			seqNo,
			// len(sendMessageResult.ReceiptIssuers),
		)

		t.Logf("Waiting for CCIPMessageSent event: from=%d to=%d seq=%d", cantonChain.ChainSelector(), evmChain.ChainSelector(), seqNo)
		sentEvent, err := cantonChain.ConfirmSendOnSource(subtestCtx, evmChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, 30*time.Second)
		require.NoError(t, err)

		t.Logf("CCIPMessageSent event: %+v", sentEvent)

		t.Logf("Asserting message propagated through aggregator/indexer: messageID=%x", sentEvent.MessageID[:])
		result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, lib, sentEvent.MessageID)
		t.Logf(
			"Message assertion succeeded: aggregated=true indexerResults=%+v",
			result.IndexedVerifications.Results,
		)

		t.Logf("Waiting for execution event on EVM: from=%d seq=%d", cantonChain.ChainSelector(), seqNo)
		ev, err := evmChain.ConfirmExecOnDest(subtestCtx, cantonChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
		require.NoError(t, err)
		assert.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
		t.Logf("Execution event: %+v", ev)
	})

	t.Run("EOA receiver and default committee verifier token transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		// Send params (transfer amount, gas limit, finality) come from token_transfer_config.toml.
		lane := devenvtests.ResolveTokenLane(t, in, lib, chainMap, cantonChain.ChainSelector(), []uint64{evmChain.ChainSelector()})
		tokenTransferAmount := lane.TransferAmount.Uint64()

		// Setup message send
		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), tokenTransferAmount))

		ds, err := lib.DataStore()
		require.NoError(t, err)
		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr := devenvtests.GetContractAddress(
			t, ds, cantonChain.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		executorAddr := devenvtests.GetContractAddress(
			t, ds, cantonChain.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)
		require.NoError(t, err)
		destTokenAddress := lane.DestTokenBySelector[evmChain.ChainSelector()]
		receiverBalanceBefore, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceBefore)

		for sendIdx := range devenvtests.CantonToEVMTokenSequentialSends {
			t.Logf("Token transfer send %d/%d", sendIdx+1, devenvtests.CantonToEVMTokenSequentialSends)
			sendMessageResult, err := cantonChain.SendMessage(
				subtestCtx,
				evmChain.ChainSelector(),
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

			sentEvent, err := cantonChain.ConfirmSendOnSource(subtestCtx, evmChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
			require.NoError(t, err)
			require.NotNil(t, sentEvent.Message)
			require.NotNil(t, sentEvent.Message.TokenTransfer)

			ev, err := evmChain.ConfirmExecOnDest(subtestCtx, cantonChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
			require.NoError(t, err)
			require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
		}

		receiverBalanceAfter, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceAfter)

		expectedTransferPerMessage := new(big.Int).Mul(lane.TransferAmount, big.NewInt(devenvtests.EVMDecimalsScale))
		totalExpectedTransfer := new(big.Int).Mul(expectedTransferPerMessage, big.NewInt(devenvtests.CantonToEVMTokenSequentialSends))
		expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), totalExpectedTransfer)
		t.Logf("EVM receiver token balance: before=%s after=%s totalExpectedTransfer=%s", receiverBalanceBefore.String(), receiverBalanceAfter.String(), totalExpectedTransfer.String())
		require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter)
	})

	// Token transfer with default extraArgs (gasLimit=0, default executor/CCVs from lane) should succeed.
	t.Run("token transfer with default extraArgs", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		lane := devenvtests.ResolveTokenLane(t, in, lib, chainMap, cantonChain.ChainSelector(), []uint64{evmChain.ChainSelector()})
		tokenTransferAmount := lane.TransferAmount.Uint64()

		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), tokenTransferAmount))

		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)

		sendMessageResult, err := cantonChain.SendMessage(
			subtestCtx,
			evmChain.ChainSelector(),
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
		sentEvent, err := cantonChain.ConfirmSendOnSource(subtestCtx, evmChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
		require.NoError(t, err)

		ev, err := evmChain.ConfirmExecOnDest(subtestCtx, cantonChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID}, tests.WaitTimeout(t))
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
	})

	// Token transfer with gasLimit=0 should succeed — the EVM receiver gets the tokens
	// even without gas for a callback since there is no receiver contract to call.
	t.Run("token transfer with zero gas limit", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		lane := devenvtests.ResolveTokenLane(t, in, lib, chainMap, cantonChain.ChainSelector(), []uint64{evmChain.ChainSelector()})
		tokenTransferAmount := lane.TransferAmount.Uint64()

		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), tokenTransferAmount))

		ds, err := lib.DataStore()
		require.NoError(t, err)
		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr := devenvtests.GetContractAddress(
			t, ds, cantonChain.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		executorAddr := devenvtests.GetContractAddress(
			t, ds, cantonChain.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)

		sendMessageResult, err := cantonChain.SendMessage(
			subtestCtx,
			evmChain.ChainSelector(),
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
		sentEvent, err := cantonChain.ConfirmSendOnSource(subtestCtx, evmChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, tests.WaitTimeout(t))
		require.NoError(t, err)

		ev, err := evmChain.ConfirmExecOnDest(subtestCtx, cantonChain.ChainSelector(), cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID}, tests.WaitTimeout(t))
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
	})
}
