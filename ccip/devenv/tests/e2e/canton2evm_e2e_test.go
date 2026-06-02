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
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
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

const (
	cantonToEVMFeeAmount            = int64(2_000)
	cantonToEVMTokenTransferAmount  = int64(1_000)                     // 1000 tokens (1000 * 10^10) on canton decimals
	evmDecimalsScale                = int64(1_000_000_000_000_000_000) // EVM 18 decimals
	cantonToEVMTokenSequentialSends = 2
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

	t.Run("EOA receiver and default committee verifier", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		// Setup message send
		require.NoError(t, cantonImpl.MintTokens(ctx, uint64(cantonToEVMFeeAmount)))
		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(cantonToEVMFeeAmount), 0))

		ds, err := lib.DataStore()
		require.NoError(t, err)
		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr, err := tcapi.GetContractAddress(
			ds,
			cantonChain.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		require.NoError(t, err)
		executorAddr, err := tcapi.GetContractAddress(
			ds,
			cantonChain.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)
		require.NoError(t, err)
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

		lane := defaultDevenvTokenLane(t, lib, in, cantonChain.ChainSelector(), evmChain.ChainSelector())

		// Setup message send
		require.NoError(t, cantonImpl.MintTokens(ctx, cantonToEVMTokenSequentialSends*uint64(cantonToEVMFeeAmount)))           // Holdings for fee
		require.NoError(t, cantonImpl.MintTokens(ctx, cantonToEVMTokenSequentialSends*uint64(cantonToEVMTokenTransferAmount))) // Holdings for token transfer
		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(cantonToEVMFeeAmount), uint64(cantonToEVMTokenTransferAmount)))    // Setup with fee and token transfer amounts

		ds, err := lib.DataStore()
		require.NoError(t, err)
		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr, err := tcapi.GetContractAddress(
			ds,
			cantonChain.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		require.NoError(t, err)
		executorAddr, err := tcapi.GetContractAddress(
			ds,
			cantonChain.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)
		require.NoError(t, err)
		destTokenAddress := lane.DestToken
		receiverBalanceBefore, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceBefore)

		for sendIdx := range cantonToEVMTokenSequentialSends {
			t.Logf("Token transfer send %d/%d", sendIdx+1, cantonToEVMTokenSequentialSends)
			sendMessageResult, err := cantonChain.SendMessage(
				subtestCtx,
				evmChain.ChainSelector(),
				cciptestinterfaces.MessageFields{
					Receiver: receiver,
					Data:     []byte("canton2evm token transfer"),
					TokenAmount: cciptestinterfaces.TokenAmount{
						Amount: big.NewInt(cantonToEVMTokenTransferAmount),
					},
				},
				cciptestinterfaces.MessageOptions{
					ExecutionGasLimit: 500_000,
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

		expectedTransferPerMessage := new(big.Int).Mul(big.NewInt(cantonToEVMTokenTransferAmount), big.NewInt(evmDecimalsScale))
		totalExpectedTransfer := new(big.Int).Mul(expectedTransferPerMessage, big.NewInt(cantonToEVMTokenSequentialSends))
		expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), totalExpectedTransfer)
		t.Logf("EVM receiver token balance: before=%s after=%s totalExpectedTransfer=%s", receiverBalanceBefore.String(), receiverBalanceAfter.String(), totalExpectedTransfer.String())
		require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter)
	})
}
