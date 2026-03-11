package canton

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/committee_verifier"
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
)

const (
	evmToCantonBasicPayload = "Hello from EVM!"
	evmToCantonTokenPayload = "Hello token transfer from EVM!"

	evmToCantonTokenQualifier = "TEST (BurnMintTokenPool 1.7.0 [default] to LockReleaseTokenPool 1.7.0 [default])"
	evmToCantonTransferAmount = int64(1000)
)

func mustGetBlockchainInputByType(t *testing.T, cfg *ccv.Cfg, chainType string) *blockchain.Input {
	t.Helper()
	for _, bc := range cfg.Blockchains {
		if bc.Type == chainType {
			return bc
		}
	}
	require.FailNowf(t, "missing chain", "need at least one %s chain for this test", chainType)
	return nil
}

func newAggregatorClients(
	t *testing.T,
	ctx context.Context,
	cfg *ccv.Cfg,
) map[string]*ccv.AggregatorClient {
	t.Helper()
	clients := make(map[string]*ccv.AggregatorClient)
	for qualifier := range cfg.AggregatorEndpoints {
		client, err := cfg.NewAggregatorClientForCommittee(
			zerolog.Ctx(ctx).With().Str("component", fmt.Sprintf("aggregator-client-%s", qualifier)).Logger(),
			qualifier,
		)
		require.NoError(t, err)
		require.NotNil(t, client)
		clients[qualifier] = client
		t.Cleanup(func() {
			client.Close()
		})
	}
	return clients
}

//nolint:paralleltest // we won't run this in parallel.
func TestEVM2Canton_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM2Canton_Basic test in short mode")
	}

	// Register the canton impl factory for the canton family.
	ccv.RegisterImplFactory(chainsel.FamilyCanton, cantondevenv.NewImplFactory())

	ctx := ccv.Plog.WithContext(t.Context())
	l := zerolog.Ctx(ctx)

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	cantonChain := mustGetBlockchainInputByType(t, in, blockchain.TypeCanton)
	evmChain := mustGetBlockchainInputByType(t, in, blockchain.TypeAnvil)

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

	srcSelector := evmDetails.ChainSelector
	srcChain := chainMap[srcSelector]
	require.NotNil(t, srcChain)
	dstSelector := cantonDetails.ChainSelector
	dstChain := chainMap[dstSelector]
	require.NotNil(t, dstChain)

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	chain := e.BlockChains.CantonChains()[cantonDetails.ChainSelector]
	participant := chain.Participants[0]
	party := participant.PartyID

	// Hash receiver party
	receiver := contracts.HashedPartyFromString(party)
	t.Logf("Message receiver: %s", receiver.Hex())

	// Get EVM CCV
	ref, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			evmDetails.ChainSelector,
			datastore.ContractType(committee_verifier.ResolverType),
			semver.MustParse(committee_verifier.Deploy.Version()),
			common.DefaultCommitteeVerifierQualifier,
		),
	)
	require.NoError(t, err, "failed to get EVM committee verifier address from datastore")
	defaultCCVAddress := protocol.UnknownAddress(gethcommon.HexToAddress(ref.Address).Bytes())

	// No-execution tag
	executorAddress := protocol.UnknownAddress(gethcommon.HexToAddress("0xEBa517d200000000000000000000000000000000").Bytes())

	// Send message
	seqNo, err := srcChain.GetExpectedNextSequenceNumber(ctx, dstSelector)
	require.NoError(t, err)
	l.Info().Uint64("SeqNo", seqNo).Msg("Expecting sequence number")
	sendMessageResult, err := srcChain.SendMessage(ctx, dstSelector, cciptestinterfaces.MessageFields{
		Receiver:    receiver.Bytes(),
		Data:        []byte(evmToCantonBasicPayload),
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
		Executor:       executorAddress,
		ExecutorArgs:   nil,
		TokenArgs:      nil,
	})
	require.NoError(t, err, "failed to send message from EVM chain")
	require.Lenf(t, sendMessageResult.ReceiptIssuers, 3, "expected 3 receipt issuers for the message")
	require.NotNil(t, sendMessageResult.Message, "expected send message result to include message payload")
	t.Logf(
		"SendMessage accepted (basic): srcSelector=%d dstSelector=%d seqNo=%d receipts=%d issuers=%x",
		srcSelector,
		dstSelector,
		sendMessageResult.Message.SequenceNumber,
		len(sendMessageResult.ReceiptIssuers),
		sendMessageResult.ReceiptIssuers,
	)
	sentEvent, err := srcChain.WaitOneSentEventBySeqNo(ctx, dstSelector, seqNo, time.Second*10)
	require.NoError(t, err)
	messageID := sentEvent.MessageID
	t.Logf("Message sent with ID: %s", hexutil.Encode(messageID[:]))

	var indexerMonitor *ccv.IndexerMonitor
	indexerClient, err := lib.Indexer()
	require.NoError(t, err)
	indexerMonitor, err = ccv.NewIndexerMonitor(
		zerolog.Ctx(ctx).With().Str("component", "indexer-client").Logger(),
		indexerClient)
	require.NoError(t, err)
	require.NotNil(t, indexerMonitor)

	aggregatorClients := newAggregatorClients(t, ctx, in)
	defaultAggregatorClient := aggregatorClients[common.DefaultCommitteeVerifierQualifier]

	t.Run("token transfer", func(t *testing.T) {
		tokenRef, err := in.CLDF.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				srcSelector,
				datastore.ContractType("BurnMintERC20WithDrip"),
				semver.MustParse("1.5.0"),
				evmToCantonTokenQualifier,
			),
		)
		require.NoError(t, err, "failed to resolve source token address for token transfer e2e")
		srcToken := protocol.UnknownAddress(gethcommon.HexToAddress(tokenRef.Address).Bytes())
		senderAddress, err := srcChain.GetSenderAddress()
		require.NoError(t, err)
		senderBalanceBefore, err := srcChain.GetTokenBalance(ctx, senderAddress, srcToken)
		require.NoError(t, err)
		require.NotNil(t, senderBalanceBefore)

		seqNo, err := srcChain.GetExpectedNextSequenceNumber(ctx, dstSelector)
		require.NoError(t, err)
		sendMessageResult, err := srcChain.SendMessage(
			ctx,
			dstSelector,
			cciptestinterfaces.MessageFields{
				Receiver: receiver.Bytes(),
				Data:     []byte(evmToCantonTokenPayload),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount:       big.NewInt(evmToCantonTransferAmount),
					TokenAddress: srcToken,
				},
			},
			cciptestinterfaces.MessageOptions{
				Version:             3,
				ExecutionGasLimit:   200_000,
				OutOfOrderExecution: false,
				CCVs: []protocol.CCV{
					{
						CCVAddress: defaultCCVAddress,
						Args:       []byte{},
						ArgsLen:    0,
					},
				},
				FinalityConfig: 0,
				Executor:       executorAddress,
				ExecutorArgs:   nil,
				TokenArgs:      nil,
			},
		)
		require.NoError(t, err, "failed to send token transfer message from EVM chain")
		require.NotNil(t, sendMessageResult.Message)
		require.NotNil(t, sendMessageResult.Message.TokenTransfer, "token transfer should be populated in sent message")
		require.Lenf(t, sendMessageResult.ReceiptIssuers, 4, "expected 4 receipt issuers for token transfer message")
		t.Logf(
			"SendMessage accepted (token transfer): srcSelector=%d dstSelector=%d seqNo=%d receipts=%d issuers=%x tokenTransferPresent=%t",
			srcSelector,
			dstSelector,
			sendMessageResult.Message.SequenceNumber,
			len(sendMessageResult.ReceiptIssuers),
			sendMessageResult.ReceiptIssuers,
			sendMessageResult.Message.TokenTransfer != nil,
		)

		sentEvent, err := srcChain.WaitOneSentEventBySeqNo(ctx, dstSelector, seqNo, time.Second*15)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.NotNil(t, sentEvent.Message.TokenTransfer, "token transfer should be present in sent event")

		testCtx := e2e.NewTestingContext(t, t.Context(), chainMap, defaultAggregatorClient, indexerMonitor)
		res, err := testCtx.AssertMessage(sentEvent.MessageID, e2e.AssertMessageOptions{
			TickInterval:            time.Second,
			Timeout:                 tests.WaitTimeout(t),
			ExpectedVerifierResults: 1,
			AssertVerifierLogs:      false,
			AssertExecutorLogs:      false,
		})
		require.NoError(t, err)
		require.NotNil(t, res.AggregatedResult)
		require.Len(t, res.IndexedVerifications.Results, 1)

		message := res.IndexedVerifications.Results[0].VerifierResult.Message
		require.NotNil(t, message.TokenTransfer, "indexed message should include token transfer")
		senderBalanceAfter, err := srcChain.GetTokenBalance(ctx, senderAddress, srcToken)
		require.NoError(t, err)
		require.NotNil(t, senderBalanceAfter)
		spent := new(big.Int).Sub(senderBalanceBefore, senderBalanceAfter)
		require.Equal(t, big.NewInt(evmToCantonTransferAmount), spent, "sender EVM token balance should decrease by transfer amount")
	})

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

	// Manually execute
	executionStateChangedEvent, err := dstChain.ManuallyExecuteMessage(ctx, message, 0, []protocol.UnknownAddress{result.IndexedVerifications.Results[0].VerifierResult.VerifierDestAddress}, [][]byte{result.IndexedVerifications.Results[0].VerifierResult.CCVData})
	require.NoError(t, err, "failed to manually execute message on Canton chain")
	require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State, "expected message execution to succeed")
	require.EqualValues(t, srcSelector, executionStateChangedEvent.SourceChainSelector, "expected source chain selector to match")
	require.Equal(t, messageID, executionStateChangedEvent.MessageID, "expected message ID to match")
	require.Equal(t, seqNo, executionStateChangedEvent.MessageNumber, "expected message number to match")
	require.Equal(t, []byte{}, executionStateChangedEvent.ReturnData, "expected empty return data from message execution")
}
