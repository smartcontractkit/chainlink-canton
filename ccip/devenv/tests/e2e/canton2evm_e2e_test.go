package canton

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
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
	cantonToEVMDestTokenQualifier  = "TEST (BurnMintTokenPool 2.0.0 [default] to LockReleaseTokenPool 2.0.0 [default])"
	cantonToEVMTokenTransferAmount = int64(1000)
	cantonToEVMDecimalsScale       = int64(100_000_000) // Canton 10 decimals -> EVM 18 decimals
)

//nolint:paralleltest // we won't run this in parallel.
func TestCanton2EVM_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton2EVM_Basic test in short mode")
	}

	ccv.RegisterImplFactory(chainsel.FamilyCanton, cantondevenv.NewImplFactory())

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	ctx := ccv.Plog.WithContext(t.Context())
	harness, err := tcapi.NewTestHarness(
		ctx,
		configPath,
		in,
		chainsel.FamilyEVM,
		chainsel.FamilyCanton,
	)
	require.NoError(t, err)

	evmChain := devenvtests.GetChain(t, blockchain.TypeAnvil, in, harness)
	cantonChain := devenvtests.GetChain(t, blockchain.TypeCanton, in, harness)

	for _, client := range harness.AggregatorClients {
		t.Cleanup(func() {
			client.Close()
		})
	}

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	// t.Run("EOA receiver and default committee verifier", func(t *testing.T) {
	// 	subtestCtx := ccv.Plog.WithContext(t.Context())

	// 	receiver, err := evmChain.GetEOAReceiverAddress()
	// 	require.NoError(t, err)
	// 	ccvAddr, err := tcapi.GetContractAddress(
	// 		in,
	// 		cantonChain.ChainSelector(),
	// 		datastore.ContractType(canton_committee_verifier.ContractType),
	// 		canton_committee_verifier.Version.String(),
	// 		devenvcommon.DefaultCommitteeVerifierQualifier,
	// 		"canton committee verifier",
	// 	)
	// 	require.NoError(t, err)
	// 	executorAddr, err := tcapi.GetContractAddress(
	// 		in,
	// 		cantonChain.ChainSelector(),
	// 		datastore.ContractType(executor.ContractType),
	// 		executor.Version.String(),
	// 		devenvcommon.DefaultExecutorQualifier,
	// 		"source executor",
	// 	)
	// 	require.NoError(t, err)
	// 	t.Logf(
	// 		"Resolved contracts: receiver=%x cantonCCV=%x cantonExecutor=%x srcSelector=%d dstSelector=%d",
	// 		receiver,
	// 		ccvAddr,
	// 		executorAddr,
	// 		cantonChain.ChainSelector(),
	// 		evmChain.ChainSelector(),
	// 	)

	// 	t.Logf("Sending Canton -> EVM message")
	// 	msgFields := cciptestinterfaces.MessageFields{
	// 		Receiver: receiver,
	// 		Data:     []byte("canton2evm tcapi test"),
	// 	}
	// 	msgOpts := cciptestinterfaces.MessageOptions{
	// 		Version:           3,
	// 		ExecutionGasLimit: 200_000,
	// 		FinalityConfig:    1,
	// 		Executor:          executorAddr,
	// 		CCVs: []protocol.CCV{
	// 			{
	// 				CCVAddress: ccvAddr,
	// 				Args:       []byte{},
	// 				ArgsLen:    0,
	// 			},
	// 		},
	// 	}
	// 	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	// 	require.True(t, ok, "Canton chain implementation must be *devenv.Chain")
	// 	require.NoError(t, cantonImpl.PrepareSendPrerequisites(subtestCtx, msgFields))
	// 	sendMessageResult, err := cantonChain.SendMessage(subtestCtx, evmChain.ChainSelector(), msgFields, msgOpts)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, sendMessageResult.Message)
	// 	// require.NotEmpty(t, sendMessageResult.ReceiptIssuers)
	// 	seqNo := uint64(sendMessageResult.Message.SequenceNumber)
	// 	t.Logf(
	// 		"SendMessage accepted: seqNo=%d",
	// 		seqNo,
	// 		// len(sendMessageResult.ReceiptIssuers),
	// 	)

	// 	t.Logf("Waiting for CCIPMessageSent event: from=%d to=%d seq=%d", cantonChain.ChainSelector(), evmChain.ChainSelector(), seqNo)
	// 	sentEvent, err := cantonChain.WaitOneSentEventBySeqNo(subtestCtx, evmChain.ChainSelector(), seqNo, 30*time.Second)
	// 	require.NoError(t, err)

	// 	t.Logf("CCIPMessageSent event: %+v", sentEvent)

	// 	t.Logf("Asserting message propagated through aggregator/indexer: messageID=%x", sentEvent.MessageID)
	// 	result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, &harness, sentEvent.MessageID)
	// 	t.Logf(
	// 		"Message assertion succeeded: aggregated=true indexerResults=%+v",
	// 		result.IndexedVerifications.Results,
	// 	)

	// 	t.Logf("Waiting for execution event on EVM: from=%d seq=%d", cantonChain.ChainSelector(), seqNo)
	// 	ev, err := evmChain.WaitOneExecEventBySeqNo(subtestCtx, cantonChain.ChainSelector(), seqNo, tests.WaitTimeout(t))
	// 	require.NoError(t, err)
	// 	assert.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
	// 	t.Logf("Execution event: %+v", ev)
	// })

	t.Run("EOA receiver and default committee verifier token transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr, err := tcapi.GetContractAddress(
			in,
			cantonChain.ChainSelector(),
			datastore.ContractType(canton_committee_verifier.ContractType),
			canton_committee_verifier.Version.String(),
			devenvcommon.DefaultCommitteeVerifierQualifier,
			"canton committee verifier",
		)
		require.NoError(t, err)
		executorAddr, err := tcapi.GetContractAddress(
			in,
			cantonChain.ChainSelector(),
			datastore.ContractType(executor.ContractType),
			executor.Version.String(),
			devenvcommon.DefaultExecutorQualifier,
			"source executor",
		)
		require.NoError(t, err)
		destTokenRef, err := in.CLDF.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				evmChain.ChainSelector(),
				datastore.ContractType("BurnMintERC20WithDripToken"),
				semver.MustParse("1.0.0"),
				cantonToEVMDestTokenQualifier,
			),
		)
		require.NoError(t, err)
		destTokenAddress := protocol.UnknownAddress(gethcommon.HexToAddress(destTokenRef.Address).Bytes())
		receiverBalanceBefore, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceBefore)

		tokenMsgFields := cciptestinterfaces.MessageFields{
			Receiver: receiver,
			Data:     []byte("canton2evm token transfer"),
			TokenAmount: cciptestinterfaces.TokenAmount{
				Amount: big.NewInt(cantonToEVMTokenTransferAmount),
			},
		}
		tokenMsgOpts := cciptestinterfaces.MessageOptions{
			Version:           3,
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
		}
		cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
		require.True(t, ok, "Canton chain implementation must be *devenv.Chain")

		const numTokenTransfers = 2
		for i := range numTokenTransfers {
			tokenMsgFields.Data = []byte(fmt.Sprintf("canton2evm token transfer #%d", i+1))
			require.NoError(t, cantonImpl.PrepareSendPrerequisites(subtestCtx))
			sendMessageResult, err := cantonChain.SendMessage(subtestCtx, evmChain.ChainSelector(), tokenMsgFields, tokenMsgOpts)
			require.NoError(t, err)
			require.NotNil(t, sendMessageResult.Message)
			require.NotNil(t, sendMessageResult.Message.TokenTransfer)
			seqNo := uint64(sendMessageResult.Message.SequenceNumber)

			sentEvent, err := cantonChain.WaitOneSentEventBySeqNo(subtestCtx, evmChain.ChainSelector(), seqNo, tests.WaitTimeout(t))
			require.NoError(t, err)
			require.NotNil(t, sentEvent.Message)
			require.NotNil(t, sentEvent.Message.TokenTransfer)

			ev, err := evmChain.WaitOneExecEventBySeqNo(subtestCtx, cantonChain.ChainSelector(), seqNo, tests.WaitTimeout(t))
			require.NoError(t, err)
			assert.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
		}

		receiverBalanceAfter, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceAfter)
		expectedTransferAmount := new(big.Int).Mul(big.NewInt(cantonToEVMTokenTransferAmount), big.NewInt(cantonToEVMDecimalsScale))
		totalExpected := new(big.Int).Mul(expectedTransferAmount, big.NewInt(numTokenTransfers))
		expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), totalExpected)
		require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter)
		t.Logf("EVM receiver token balance: before=%s after=%s (%d transfers)", receiverBalanceBefore.String(), receiverBalanceAfter.String(), numTokenTransfers)
	})
}
