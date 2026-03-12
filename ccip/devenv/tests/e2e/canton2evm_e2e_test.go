package canton

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/executor"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
)

const (
	cantonToEVMBasicPayload = "Hello from Canton!"
	cantonToEVMTokenPayload = "Hello token transfer from Canton!"

	// Token qualifier for Canton lock/release -> EVM burn/mint topology.
	cantonToEVMDestTokenQualifier  = "TEST (BurnMintTokenPool 1.7.0 [default] to LockReleaseTokenPool 1.7.0 [default])"
	cantonToEVMTokenTransferAmount = int64(1000)

	cantonToEVMSentEventTimeout = 2 * time.Minute
	cantonToEVMExecutionTimeout = 4 * time.Minute
)

func resolveCantonSendContracts(t *testing.T, in *ccv.Cfg, chainSelector uint64) (protocol.UnknownAddress, protocol.UnknownAddress) {
	t.Helper()

	selectRef := func(contractType datastore.ContractType, preferredQualifier string) datastore.AddressRef {
		candidates := in.CLDF.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByType(contractType),
		)
		require.NotEmptyf(t, candidates, "missing %s address for chain selector %d", contractType, chainSelector)
		for _, ref := range candidates {
			if ref.Qualifier == preferredQualifier {
				return ref
			}
		}
		return candidates[0]
	}

	ccvRef := selectRef(datastore.ContractType(canton_committee_verifier.ContractType), devenvcommon.DefaultCommitteeVerifierQualifier)
	executorRef := selectRef(datastore.ContractType(executor.ContractType), devenvcommon.DefaultExecutorQualifier)

	return protocol.UnknownAddress(contracts.HexToInstanceAddress(ccvRef.Address).Bytes()),
		protocol.UnknownAddress(contracts.HexToInstanceAddress(executorRef.Address).Bytes())
}

func assertExecutionStateSuccess(t *testing.T, ev cciptestinterfaces.ExecutionStateChangedEvent, context string) {
	t.Helper()
	require.Equalf(t, cciptestinterfaces.ExecutionStateSuccess, ev.State, "execution state should be success for %s", context)
}

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

	t.Run("EOA receiver and default committee verifier", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		ccvAddr, executorAddr := resolveCantonSendContracts(t, in, cantonChain.ChainSelector())
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
				Data:     []byte(cantonToEVMBasicPayload),
			},
			cciptestinterfaces.MessageOptions{
				Version: 3,
				// Token transfer execution on EVM is materially more expensive than a data-only message.
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
		)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		require.NotEmpty(t, sendMessageResult.ReceiptIssuers)
		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		t.Logf(
			"SendMessage accepted (basic): srcSelector=%d dstSelector=%d seqNo=%d receipts=%d issuers=%x tokenTransferPresent=%t",
			cantonChain.ChainSelector(),
			evmChain.ChainSelector(),
			seqNo,
			len(sendMessageResult.ReceiptIssuers),
			sendMessageResult.ReceiptIssuers,
			sendMessageResult.Message.TokenTransfer != nil,
		)

		t.Logf("Waiting for CCIPMessageSent event: from=%d to=%d seq=%d", cantonChain.ChainSelector(), evmChain.ChainSelector(), seqNo)
		sentEvent, err := cantonChain.WaitOneSentEventBySeqNo(subtestCtx, evmChain.ChainSelector(), seqNo, cantonToEVMSentEventTimeout)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.Equal(t, seqNo, uint64(sentEvent.Message.SequenceNumber), "sent event sequence number should match send result")
		require.Nil(t, sentEvent.Message.TokenTransfer, "basic message should not include token transfer")
		t.Logf(
			"CCIPMessageSent observed (basic): messageID=%x seqNo=%d tokenTransferPresent=%t",
			sentEvent.MessageID,
			sentEvent.Message.SequenceNumber,
			sentEvent.Message.TokenTransfer != nil,
		)

		chainMap, err := harness.Lib.ChainsMap(subtestCtx)
		require.NoError(t, err)
		testCtx, cleanupFn := tcapi.NewTestingContext(subtestCtx, chainMap, harness.AggregatorClients[devenvcommon.DefaultCommitteeVerifierQualifier], harness.IndexerMonitor)
		defer cleanupFn()
		t.Logf("Asserting message propagated through aggregator/indexer: messageID=%x", sentEvent.MessageID)
		result, err := testCtx.AssertMessage(sentEvent.MessageID, tcapi.AssertMessageOptions{
			TickInterval:            1 * time.Second,
			ExpectedVerifierResults: 1,
			Timeout:                 tests.WaitTimeout(t),
			AssertVerifierLogs:      false,
			AssertExecutorLogs:      false,
		})
		require.NoError(t, err)
		require.NotNil(t, result.AggregatedResult)
		require.Len(t, result.IndexedVerifications.Results, 1)
		t.Logf(
			"Message assertion succeeded: aggregated=true indexerResults=%+v",
			result.IndexedVerifications.Results,
		)

		t.Logf("Waiting for execution event on EVM: from=%d seq=%d", cantonChain.ChainSelector(), seqNo)
		ev, err := evmChain.WaitOneExecEventBySeqNo(subtestCtx, cantonChain.ChainSelector(), seqNo, cantonToEVMExecutionTimeout)
		require.NoError(t, err)
		assertExecutionStateSuccess(t, ev, "basic")
		t.Logf("Execution event: %+v", ev)
	})

	t.Run("EOA receiver and default committee verifier with token transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		receiver, err := evmChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		destTokenRef, err := in.CLDF.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				evmChain.ChainSelector(),
				datastore.ContractType("BurnMintERC20WithDrip"),
				semver.MustParse("1.5.0"),
				cantonToEVMDestTokenQualifier,
			),
		)
		require.NoError(t, err, "failed to resolve destination EVM token address for Canton->EVM token transfer")
		destTokenAddress := protocol.UnknownAddress(gethcommon.HexToAddress(destTokenRef.Address).Bytes())
		receiverBalanceBefore, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceBefore)
		t.Logf(
			"Receiver balance before token execution (EVM): receiver=%x token=%x balance=%s",
			receiver,
			destTokenAddress,
			receiverBalanceBefore.String(),
		)
		ccvAddr, executorAddr := resolveCantonSendContracts(t, in, cantonChain.ChainSelector())

		t.Logf("Sending Canton -> EVM token transfer message")
		sendMessageResult, err := cantonChain.SendMessage(
			subtestCtx,
			evmChain.ChainSelector(),
			cciptestinterfaces.MessageFields{
				Receiver: receiver,
				Data:     []byte(cantonToEVMTokenPayload),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount: big.NewInt(cantonToEVMTokenTransferAmount),
				},
			},
			cciptestinterfaces.MessageOptions{
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
			},
		)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		require.NotNil(t, sendMessageResult.Message.TokenTransfer, "token transfer should be populated for Canton source sends")
		require.GreaterOrEqual(t, len(sendMessageResult.ReceiptIssuers), 4, "token transfer should include additional receipt issuer")
		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		t.Logf(
			"SendMessage accepted (token transfer): srcSelector=%d dstSelector=%d seqNo=%d receipts=%d issuers=%x tokenTransferPresent=%t",
			cantonChain.ChainSelector(),
			evmChain.ChainSelector(),
			seqNo,
			len(sendMessageResult.ReceiptIssuers),
			sendMessageResult.ReceiptIssuers,
			sendMessageResult.Message.TokenTransfer != nil,
		)

		sentEvent, err := cantonChain.WaitOneSentEventBySeqNo(subtestCtx, evmChain.ChainSelector(), seqNo, cantonToEVMSentEventTimeout)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.NotNil(t, sentEvent.Message.TokenTransfer, "sent event should include token transfer")
		require.Equal(t, seqNo, uint64(sentEvent.Message.SequenceNumber), "token transfer sent event sequence number should match send result")
		t.Logf(
			"CCIPMessageSent observed (token transfer): messageID=%x seqNo=%d tokenTransferPresent=%t",
			sentEvent.MessageID,
			sentEvent.Message.SequenceNumber,
			sentEvent.Message.TokenTransfer != nil,
		)
		t.Logf("Waiting for execution event on EVM for token transfer: from=%d seq=%d", cantonChain.ChainSelector(), seqNo)
		ev, err := evmChain.WaitOneExecEventBySeqNo(subtestCtx, cantonChain.ChainSelector(), seqNo, cantonToEVMExecutionTimeout)
		require.NoError(t, err)
		assertExecutionStateSuccess(t, ev, "token transfer")
		receiverBalanceAfter, err := evmChain.GetTokenBalance(subtestCtx, receiver, destTokenAddress)
		require.NoError(t, err)
		require.NotNil(t, receiverBalanceAfter)
		expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), big.NewInt(cantonToEVMTokenTransferAmount))
		require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter, "receiver final token balance should equal initial balance plus transfer amount")
		transferred := new(big.Int).Sub(receiverBalanceAfter, receiverBalanceBefore)
		require.Equal(t, big.NewInt(cantonToEVMTokenTransferAmount), transferred, "receiver token balance should increase by transfer amount")
		t.Logf(
			"Receiver balance after token execution (EVM): receiver=%x before=%s after=%s delta=%s",
			receiver,
			receiverBalanceBefore.String(),
			receiverBalanceAfter.String(),
			transferred.String(),
		)
		t.Logf("Token transfer sent and execution observed successfully with sequence %d", seqNo)
	})
}
