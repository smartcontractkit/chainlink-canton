package canton

import (
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/executor"
	mockreceiver "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/mock_receiver"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantoncommitteeverifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
)

//nolint:paralleltest // we won't run this in parallel.
func TestCanton2EVM_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton2EVM_Basic test in short mode")
	}

	// Register the canton impl factory for the canton family.
	ccv.RegisterImplFactory(chainsel.FamilyCanton, cantondevenv.NewImplFactory())

	ctx := ccv.Plog.WithContext(t.Context())
	l := zerolog.Ctx(ctx)

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	var cantonChain *blockchain.Input
	for _, bc := range in.Blockchains {
		if bc.Type == blockchain.TypeCanton {
			cantonChain = bc
			break
		}
	}
	require.NotNil(t, cantonChain, "need at least one canton chain for this test")

	var evmChain *blockchain.Input
	for _, bc := range in.Blockchains {
		if bc.Type == blockchain.TypeAnvil {
			evmChain = bc
			break
		}
	}
	require.NotNil(t, evmChain, "need at least one evm chain for this test")

	cantonDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(cantonChain.ChainID, chainsel.FamilyCanton)
	require.NoError(t, err)
	evmDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(evmChain.ChainID, chainsel.FamilyEVM)
	require.NoError(t, err)

	_, e, err := ccv.NewCLDFOperationsEnvironment(in.Blockchains, in.CLDF.DataStore)
	require.NoError(t, err)
	b := ccv.NewDefaultCLDFBundle(e)
	e.OperationsBundle = b

	lib, err := ccv.NewLib(l, configPath, chainsel.FamilyEVM, chainsel.FamilyCanton)
	require.NoError(t, err)
	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)

	srcSelector := cantonDetails.ChainSelector
	srcChain := chainMap[srcSelector]
	require.NotNil(t, srcChain)
	dstSelector := evmDetails.ChainSelector
	dstChain := chainMap[dstSelector]
	require.NotNil(t, dstChain)

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	// Destination receiver on EVM (deployed MockReceiver contract).
	receiverRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			evmDetails.ChainSelector,
			datastore.ContractType(mockreceiver.ContractType),
			mockreceiver.Version,
			common.DefaultCommitteeVerifierQualifier,
		),
	)
	require.NoError(t, err)
	receiverAddress := protocol.UnknownAddress(gethcommon.HexToAddress(receiverRef.Address).Bytes())

	// Canton default CCV and executor addresses.
	cantonCCVRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			cantonDetails.ChainSelector,
			datastore.ContractType(cantoncommitteeverifier.ContractType),
			cantoncommitteeverifier.Version,
			common.DefaultCommitteeVerifierQualifier,
		),
	)
	require.NoError(t, err)
	defaultCCVAddress := protocol.UnknownAddress(contracts.HexToInstanceAddress(cantonCCVRef.Address).Bytes())

	cantonExecutorRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			cantonDetails.ChainSelector,
			datastore.ContractType(executor.ProxyType),
			semver.MustParse(executor.DeployProxy.Version()),
			common.DefaultExecutorQualifier,
		),
	)
	require.NoError(t, err)
	normalizedExecutorAddress := contracts.HexToInstanceAddress(cantonExecutorRef.Address) // Left-pad to 32 bytes if needed.
	normalizedExecutorBytes := normalizedExecutorAddress.Bytes()
	require.Len(t, normalizedExecutorBytes, 32, "executor address must be 32 bytes after normalization")
	defaultExecutorAddress := protocol.UnknownAddress(normalizedExecutorBytes)

	sendMessageResult, err := srcChain.SendMessage(ctx, dstSelector, cciptestinterfaces.MessageFields{
		Receiver:    receiverAddress,
		Data:        []byte("Hello from Canton!"),
		TokenAmount: cciptestinterfaces.TokenAmount{},
		FeeToken:    nil,
	}, cciptestinterfaces.MessageOptions{
		Version:             3,
		ExecutionGasLimit:   100_000,
		OutOfOrderExecution: false,
		CCVs: []protocol.CCV{
			{
				CCVAddress: defaultCCVAddress,
				Args:       []byte{},
				ArgsLen:    0,
			},
		},
		FinalityConfig: 0,
		Executor:       defaultExecutorAddress,
		ExecutorArgs:   nil,
		TokenArgs:      nil,
	})

	require.NoError(t, err, "failed to send message from Canton chain")
	require.Lenf(t, sendMessageResult.ReceiptIssuers, 3, "expected 3 receipt issuers for the message")
	require.NotNil(t, sendMessageResult.Message, "SendMessage result must include decoded message")
	seqNo := uint64(sendMessageResult.Message.SequenceNumber)
	l.Info().Uint64("SeqNo", seqNo).Msg("Using sequence number from sent message")

	sentEvent, err := srcChain.WaitOneSentEventBySeqNo(ctx, dstSelector, seqNo, 10*time.Second)
	require.NoError(t, err)
	messageID := sentEvent.MessageID

	indexerClient, err := lib.Indexer()
	require.NoError(t, err)
	indexerMonitor, err := ccv.NewIndexerMonitor(
		zerolog.Ctx(ctx).With().Str("component", "indexer-client").Logger(),
		indexerClient,
	)
	require.NoError(t, err)

	aggregatorClients := make(map[string]*ccv.AggregatorClient)
	for qualifier := range in.AggregatorEndpoints {
		client, err := in.NewAggregatorClientForCommittee(
			zerolog.Ctx(ctx).With().Str("component", fmt.Sprintf("aggregator-client-%s", qualifier)).Logger(),
			qualifier,
		)
		require.NoError(t, err)
		require.NotNil(t, client)
		aggregatorClients[qualifier] = client
		t.Cleanup(func() {
			client.Close()
		})
	}
	defaultAggregatorClient := aggregatorClients[common.DefaultCommitteeVerifierQualifier]

	testCtx := e2e.NewTestingContext(t, t.Context(), chainMap, defaultAggregatorClient, indexerMonitor)
	result, err := testCtx.AssertMessage(messageID, e2e.AssertMessageOptions{
		TickInterval:            time.Second,
		Timeout:                 tests.WaitTimeout(t),
		ExpectedVerifierResults: 1,
		AssertVerifierLogs:      false,
		AssertExecutorLogs:      false,
	})
	require.NoError(t, err)
	require.NotNil(t, result.AggregatedResult)
	require.Len(t, result.IndexedVerifications.Results, 1)
	message := result.IndexedVerifications.Results[0].VerifierResult.Message
	expectedExecutedMessageID := message.MustMessageID()
	expectedExecutedMessageNumber := uint64(message.SequenceNumber)

	executionStateChangedEvent, err := dstChain.WaitOneExecEventBySeqNo(ctx, srcSelector, seqNo, tests.WaitTimeout(t))
	require.NoError(t, err)
	l.Info().Interface("ExecutionStateChangedEvent", executionStateChangedEvent).Msg("executionStateChangedEvent")

	require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State, "expected execution state to be success")
	require.EqualValues(t, srcSelector, message.SourceChainSelector, "expected source chain selector to match in executed message")
	if executionStateChangedEvent.SourceChainSelector != 0 {
		require.EqualValues(t, srcSelector, executionStateChangedEvent.SourceChainSelector, "expected source chain selector to match")
	}
	require.EqualValues(t, expectedExecutedMessageID, executionStateChangedEvent.MessageID, "expected message ID to match")
	require.Equal(t, expectedExecutedMessageNumber, executionStateChangedEvent.MessageNumber, "expected message number to match")
}
