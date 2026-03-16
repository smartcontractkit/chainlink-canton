package canton

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
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

	// We do not call basic.EOAReceiverDefaultVerifier here because its prerequisite
	// lookup assumes EVM contract metadata (committee verifier resolver) on the source chain.
	// For Canton source, we resolve Canton contract types directly and run the same tcapi assertions.
	t.Run("EOA receiver and default committee verifier", func(t *testing.T) {
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
			datastore.ContractType(executor.ProxyType),
			executor.DeployProxy.Version(),
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
				Version:           3,
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
		)
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		require.NotEmpty(t, sendMessageResult.ReceiptIssuers)
		seqNo := uint64(sendMessageResult.Message.SequenceNumber)
		t.Logf(
			"SendMessage accepted: seqNo=%d receipts=%d",
			seqNo,
			len(sendMessageResult.ReceiptIssuers),
		)

		t.Logf("Waiting for CCIPMessageSent event: from=%d to=%d seq=%d", cantonChain.ChainSelector(), evmChain.ChainSelector(), seqNo)
		sentEvent, err := cantonChain.WaitOneSentEventBySeqNo(subtestCtx, evmChain.ChainSelector(), seqNo, 30*time.Second)
		require.NoError(t, err)

		t.Logf("CCIPMessageSent event: %+v", sentEvent)

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
		ev, err := evmChain.WaitOneExecEventBySeqNo(subtestCtx, cantonChain.ChainSelector(), seqNo, tests.WaitTimeout(t))
		require.NoError(t, err)
		assert.Equal(t, cciptestinterfaces.ExecutionStateSuccess, ev.State)
		t.Logf("Execution event: %+v", ev)
	})
}
