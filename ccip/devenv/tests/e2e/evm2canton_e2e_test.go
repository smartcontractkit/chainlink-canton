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

type evmToCantonHarness struct {
	cfg                    *ccv.Cfg
	logger                 zerolog.Logger
	srcSelector            uint64
	dstSelector            uint64
	srcChain               cciptestinterfaces.CCIP17
	dstChain               cciptestinterfaces.CCIP17
	receiver               contracts.HashedParty
	receiverParty          string
	defaultCCVAddress      protocol.UnknownAddress
	executorAddress        protocol.UnknownAddress
	assertSingleVerifierFn func(t *testing.T, messageID [32]byte) (protocol.Message, protocol.UnknownAddress, []byte)
}

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

func defaultEVMToCantonMessageOptions(
	defaultCCVAddress protocol.UnknownAddress,
	executorAddress protocol.UnknownAddress,
	executionGasLimit uint32,
) cciptestinterfaces.MessageOptions {
	return defaultEVMToCantonMessageOptionsWithFinality(defaultCCVAddress, executorAddress, executionGasLimit, 0)
}

func defaultEVMToCantonMessageOptionsWithFinality(
	defaultCCVAddress protocol.UnknownAddress,
	executorAddress protocol.UnknownAddress,
	executionGasLimit uint32,
	finality uint16,
) cciptestinterfaces.MessageOptions {
	return cciptestinterfaces.MessageOptions{
		Version:             3,
		ExecutionGasLimit:   executionGasLimit,
		OutOfOrderExecution: false,
		CCVs: []protocol.CCV{
			{
				CCVAddress: defaultCCVAddress,
				Args:       []byte{},
				ArgsLen:    0,
			},
		},
		FinalityConfig: finality,
		Executor:       executorAddress,
		ExecutorArgs:   nil,
		TokenArgs:      nil,
	}
}

func setupEVMToCantonHarness(t *testing.T) *evmToCantonHarness {
	t.Helper()

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
	poolOwnerParty := chain.Participants[0].PartyID
	require.GreaterOrEqual(t, len(chain.Participants), 2, "Canton chain must have at least 2 participants")
	receiverParty := chain.Participants[1].PartyID
	require.NotEqual(t, poolOwnerParty, receiverParty, "receiver party must differ from pool owner party")

	receiver := contracts.HashedPartyFromString(receiverParty)
	t.Logf("Message receiver: %s (party=%s)", receiver.Hex(), receiverParty)

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

	executorAddress := protocol.UnknownAddress(gethcommon.HexToAddress("0xEBa517d200000000000000000000000000000000").Bytes())

	indexerClient, err := lib.Indexer()
	require.NoError(t, err)
	indexerMonitor, err := ccv.NewIndexerMonitor(
		zerolog.Ctx(ctx).With().Str("component", "indexer-client").Logger(),
		indexerClient)
	require.NoError(t, err)
	require.NotNil(t, indexerMonitor)

	aggregatorClients := newAggregatorClients(t, ctx, in)
	defaultAggregatorClient := aggregatorClients[common.DefaultCommitteeVerifierQualifier]
	assertSingleVerifier := func(t *testing.T, messageID [32]byte) (protocol.Message, protocol.UnknownAddress, []byte) {
		t.Helper()
		testCtx := e2e.NewTestingContext(t, t.Context(), chainMap, defaultAggregatorClient, indexerMonitor)
		res, err := testCtx.AssertMessage(messageID, e2e.AssertMessageOptions{
			TickInterval:            time.Second,
			Timeout:                 tests.WaitTimeout(t),
			ExpectedVerifierResults: 1,
			AssertVerifierLogs:      false,
			AssertExecutorLogs:      false,
		})
		require.NoError(t, err)
		require.NotNil(t, res.AggregatedResult)
		require.Len(t, res.IndexedVerifications.Results, 1)
		vr := res.IndexedVerifications.Results[0].VerifierResult

		return vr.Message, vr.VerifierDestAddress, vr.CCVData
	}

	return &evmToCantonHarness{
		cfg:                    in,
		logger:                 *l,
		srcSelector:            srcSelector,
		dstSelector:            dstSelector,
		srcChain:               srcChain,
		dstChain:               dstChain,
		receiver:               receiver,
		receiverParty:          receiverParty,
		defaultCCVAddress:      defaultCCVAddress,
		executorAddress:        executorAddress,
		assertSingleVerifierFn: assertSingleVerifier,
	}
}

//nolint:paralleltest // we won't run this in parallel.
func TestEVM2Canton_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM2Canton_Basic test in short mode")
	}
	h := setupEVMToCantonHarness(t)

	t.Run("message transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())
		seqNo, err := h.srcChain.GetExpectedNextSequenceNumber(subtestCtx, h.dstSelector)
		require.NoError(t, err)
		h.logger.Info().Uint64("SeqNo", seqNo).Msg("Expecting sequence number")

		sendMessageResult, err := h.srcChain.SendMessage(subtestCtx, h.dstSelector, cciptestinterfaces.MessageFields{
			Receiver:    h.receiver.Bytes(),
			Data:        []byte(evmToCantonBasicPayload),
			TokenAmount: cciptestinterfaces.TokenAmount{},
			FeeToken:    nil,
		}, defaultEVMToCantonMessageOptions(h.defaultCCVAddress, h.executorAddress, 100_000))
		require.NoError(t, err, "failed to send message from EVM chain")
		require.Lenf(t, sendMessageResult.ReceiptIssuers, 3, "expected 3 receipt issuers for the message")
		require.NotNil(t, sendMessageResult.Message, "expected send message result to include message payload")
		t.Logf(
			"SendMessage accepted (basic): srcSelector=%d dstSelector=%d seqNo=%d receipts=%d issuers=%x",
			h.srcSelector,
			h.dstSelector,
			sendMessageResult.Message.SequenceNumber,
			len(sendMessageResult.ReceiptIssuers),
			sendMessageResult.ReceiptIssuers,
		)

		sentEvent, err := h.srcChain.WaitOneSentEventBySeqNo(subtestCtx, h.dstSelector, seqNo, time.Second*10)
		require.NoError(t, err)
		messageID := sentEvent.MessageID
		t.Logf("Message sent with ID: %s", hexutil.Encode(messageID[:]))

		message, verifierDestAddress, ccvData := h.assertSingleVerifierFn(t, messageID)
		executionStateChangedEvent, err := h.dstChain.ManuallyExecuteMessage(subtestCtx, message, 0, []protocol.UnknownAddress{verifierDestAddress}, [][]byte{ccvData})
		require.NoError(t, err, "failed to manually execute message on Canton chain")
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State, "expected message execution to succeed")
		require.EqualValues(t, h.srcSelector, executionStateChangedEvent.SourceChainSelector, "expected source chain selector to match")
		require.Equal(t, messageID, executionStateChangedEvent.MessageID, "expected message ID to match")
		require.Equal(t, seqNo, executionStateChangedEvent.MessageNumber, "expected message number to match")
		require.Equal(t, []byte{}, executionStateChangedEvent.ReturnData, "expected empty return data from message execution")
	})

	t.Run("token transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())
		tokenRefAddress, err := h.cfg.CLDF.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				h.srcSelector,
				datastore.ContractType("BurnMintERC20WithDrip"),
				semver.MustParse("1.5.0"),
				evmToCantonTokenQualifier,
			),
		)
		require.NoError(t, err, "failed to resolve source token address for token transfer e2e")
		srcToken := protocol.UnknownAddress(gethcommon.HexToAddress(tokenRefAddress.Address).Bytes())
		cantonReceiver := protocol.UnknownAddress(h.receiver.Bytes())
		receiverBalanceBefore, err := h.dstChain.GetTokenBalance(subtestCtx, cantonReceiver, srcToken)
		require.NoError(t, err, "failed to read receiver token balance on Canton before execution")
		require.NotNil(t, receiverBalanceBefore)
		t.Logf(
			"Receiver balance before token execution (Canton): receiver=%s token=%x balance=%s",
			h.receiverParty,
			srcToken,
			receiverBalanceBefore.String(),
		)

		seqNo, err := h.srcChain.GetExpectedNextSequenceNumber(subtestCtx, h.dstSelector)
		require.NoError(t, err)
		sendMessageResult, err := h.srcChain.SendMessage(
			subtestCtx,
			h.dstSelector,
			cciptestinterfaces.MessageFields{
				Receiver: h.receiver.Bytes(),
				Data:     []byte(evmToCantonTokenPayload),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount:       big.NewInt(evmToCantonTransferAmount),
					TokenAddress: srcToken,
				},
			},
			defaultEVMToCantonMessageOptions(h.defaultCCVAddress, h.executorAddress, 200_000),
		)
		require.NoError(t, err, "failed to send token transfer message from EVM chain")
		require.NotNil(t, sendMessageResult.Message)
		require.NotNil(t, sendMessageResult.Message.TokenTransfer, "token transfer should be populated in sent message")
		require.Lenf(t, sendMessageResult.ReceiptIssuers, 4, "expected 4 receipt issuers for token transfer message")
		t.Logf(
			"SendMessage accepted (token transfer): srcSelector=%d dstSelector=%d seqNo=%d receipts=%d issuers=%x tokenTransferPresent=%t",
			h.srcSelector,
			h.dstSelector,
			sendMessageResult.Message.SequenceNumber,
			len(sendMessageResult.ReceiptIssuers),
			sendMessageResult.ReceiptIssuers,
			sendMessageResult.Message.TokenTransfer != nil,
		)

		sentEvent, err := h.srcChain.WaitOneSentEventBySeqNo(subtestCtx, h.dstSelector, seqNo, time.Second*15)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.NotNil(t, sentEvent.Message.TokenTransfer, "token transfer should be present in sent event")

		message, verifierDestAddress, ccvData := h.assertSingleVerifierFn(t, sentEvent.MessageID)
		require.NotNil(t, message.TokenTransfer, "indexed message should include token transfer")
		require.EqualValues(t, h.receiver.Bytes(), message.TokenTransfer.TokenReceiver, "token transfer receiver should match Canton message receiver")
		// Assert the receiver's token amount (amount credited to receiver on Canton) matches the EVM->Canton transfer amount.
		require.Equal(t, 0, message.TokenTransfer.Amount.Cmp(big.NewInt(evmToCantonTransferAmount)), "receiver token amount should match transfer amount")

		// Note: verifierDestAddress is the canton verifier not EVM.
		executionStateChangedEvent, err := h.dstChain.ManuallyExecuteMessage(subtestCtx, message, 0, []protocol.UnknownAddress{verifierDestAddress}, [][]byte{ccvData})
		require.NoError(t, err, "failed to manually execute token transfer message on Canton chain")
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State, "expected token transfer message execution to succeed")

		receiverBalanceAfter, err := h.dstChain.GetTokenBalance(subtestCtx, cantonReceiver, srcToken)
		require.NoError(t, err, "failed to read receiver token balance on Canton after execution")
		require.NotNil(t, receiverBalanceAfter)
		transferred := new(big.Int).Sub(receiverBalanceAfter, receiverBalanceBefore)
		expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), big.NewInt(evmToCantonTransferAmount))
		require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter, "receiver final token balance should equal initial balance plus transfer amount")
		require.Equal(t, big.NewInt(evmToCantonTransferAmount), transferred, "receiver token balance should increase by transfer amount")
		t.Logf(
			"Receiver balance after token execution (Canton): receiver=%s before=%s after=%s delta=%s",
			h.receiverParty,
			receiverBalanceBefore.String(),
			receiverBalanceAfter.String(),
			transferred.String(),
		)
	})
}

//nolint:paralleltest // we won't run this in parallel.
func TestEVM2Canton_NonZeroFinality(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM2Canton_NonZeroFinality test in short mode")
	}

	h := setupEVMToCantonHarness(t)
	subtestCtx := ccv.Plog.WithContext(t.Context())
	seqNo, err := h.srcChain.GetExpectedNextSequenceNumber(subtestCtx, h.dstSelector)
	require.NoError(t, err)

	sendMessageResult, err := h.srcChain.SendMessage(subtestCtx, h.dstSelector, cciptestinterfaces.MessageFields{
		Receiver:    h.receiver.Bytes(),
		Data:        []byte(evmToCantonBasicPayload + " with FTF"),
		TokenAmount: cciptestinterfaces.TokenAmount{},
		FeeToken:    nil,
	}, defaultEVMToCantonMessageOptionsWithFinality(h.defaultCCVAddress, h.executorAddress, 100_000, 1))
	require.NoError(t, err, "failed to send non-zero-finality message from EVM chain")
	require.NotNil(t, sendMessageResult.Message, "expected send message result to include message payload")
	require.EqualValues(t, 1, sendMessageResult.Message.Finality, "expected sent message finality to be non-zero")
	t.Logf(
		"SendMessage accepted (non-zero finality): srcSelector=%d dstSelector=%d seqNo=%d finality=%d receipts=%d issuers=%x",
		h.srcSelector,
		h.dstSelector,
		sendMessageResult.Message.SequenceNumber,
		sendMessageResult.Message.Finality,
		len(sendMessageResult.ReceiptIssuers),
		sendMessageResult.ReceiptIssuers,
	)

	sentEvent, err := h.srcChain.WaitOneSentEventBySeqNo(subtestCtx, h.dstSelector, seqNo, time.Second*10)
	require.NoError(t, err)
	messageID := sentEvent.MessageID
	t.Logf("Message sent with non-zero finality and ID: %s", hexutil.Encode(messageID[:]))

	message, verifierDestAddress, ccvData := h.assertSingleVerifierFn(t, messageID)
	require.EqualValues(t, 1, message.Finality, "indexed message should preserve non-zero finality")

	executionStateChangedEvent, err := h.dstChain.ManuallyExecuteMessage(subtestCtx, message, 0, []protocol.UnknownAddress{verifierDestAddress}, [][]byte{ccvData})
	require.NoError(t, err, "failed to manually execute non-zero-finality message on Canton chain")
	require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State, "expected non-zero-finality message execution to succeed")
	require.EqualValues(t, h.srcSelector, executionStateChangedEvent.SourceChainSelector, "expected source chain selector to match")
	require.Equal(t, messageID, executionStateChangedEvent.MessageID, "expected message ID to match")
	require.Equal(t, seqNo, executionStateChangedEvent.MessageNumber, "expected message number to match")
	require.Equal(t, []byte{}, executionStateChangedEvent.ReturnData, "expected empty return data from message execution")
}

//nolint:paralleltest // we won't run this in parallel.
func TestEVM2Canton_TokenTransfer_NonZeroFinality(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM2Canton_TokenTransfer_NonZeroFinality test in short mode")
	}

	h := setupEVMToCantonHarness(t)
	subtestCtx := ccv.Plog.WithContext(t.Context())
	tokenRef, err := h.cfg.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			h.srcSelector,
			datastore.ContractType("BurnMintERC20WithDrip"),
			semver.MustParse("1.5.0"),
			evmToCantonTokenQualifier,
		),
	)
	require.NoError(t, err, "failed to resolve source token address for FTF token transfer e2e")
	srcToken := protocol.UnknownAddress(gethcommon.HexToAddress(tokenRef.Address).Bytes())
	cantonReceiver := protocol.UnknownAddress(h.receiver.Bytes())
	receiverBalanceBefore, err := h.dstChain.GetTokenBalance(subtestCtx, cantonReceiver, srcToken)
	require.NoError(t, err, "failed to read receiver token balance on Canton before FTF execution")
	require.NotNil(t, receiverBalanceBefore)

	seqNo, err := h.srcChain.GetExpectedNextSequenceNumber(subtestCtx, h.dstSelector)
	require.NoError(t, err)
	sendMessageResult, err := h.srcChain.SendMessage(
		subtestCtx,
		h.dstSelector,
		cciptestinterfaces.MessageFields{
			Receiver: h.receiver.Bytes(),
			Data:     []byte(evmToCantonTokenPayload + " with FTF"),
			TokenAmount: cciptestinterfaces.TokenAmount{
				Amount:       big.NewInt(evmToCantonTransferAmount),
				TokenAddress: srcToken,
			},
		},
		defaultEVMToCantonMessageOptionsWithFinality(h.defaultCCVAddress, h.executorAddress, 200_000, 1),
	)
	require.NoError(t, err, "failed to send non-zero-finality token transfer message from EVM chain")
	require.NotNil(t, sendMessageResult.Message)
	require.NotNil(t, sendMessageResult.Message.TokenTransfer, "token transfer should be populated in sent FTF message")
	require.EqualValues(t, 1, sendMessageResult.Message.Finality, "expected sent token transfer finality to be non-zero")

	sentEvent, err := h.srcChain.WaitOneSentEventBySeqNo(subtestCtx, h.dstSelector, seqNo, time.Second*15)
	require.NoError(t, err)
	require.NotNil(t, sentEvent.Message)
	require.NotNil(t, sentEvent.Message.TokenTransfer, "token transfer should be present in sent FTF event")

	message, verifierDestAddress, ccvData := h.assertSingleVerifierFn(t, sentEvent.MessageID)
	require.NotNil(t, message.TokenTransfer, "indexed FTF message should include token transfer")
	require.EqualValues(t, 1, message.Finality, "indexed FTF token transfer should preserve non-zero finality")
	require.EqualValues(t, h.receiver.Bytes(), message.TokenTransfer.TokenReceiver, "token transfer receiver should match Canton message receiver")
	require.Equal(t, 0, message.TokenTransfer.Amount.Cmp(big.NewInt(evmToCantonTransferAmount)), "receiver token amount should match transfer amount")

	executionStateChangedEvent, err := h.dstChain.ManuallyExecuteMessage(subtestCtx, message, 0, []protocol.UnknownAddress{verifierDestAddress}, [][]byte{ccvData})
	require.NoError(t, err, "failed to manually execute non-zero-finality token transfer message on Canton chain")
	require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State, "expected FTF token transfer execution to succeed")

	receiverBalanceAfter, err := h.dstChain.GetTokenBalance(subtestCtx, cantonReceiver, srcToken)
	require.NoError(t, err, "failed to read receiver token balance on Canton after FTF execution")
	require.NotNil(t, receiverBalanceAfter)
	transferred := new(big.Int).Sub(receiverBalanceAfter, receiverBalanceBefore)
	expectedReceiverBalanceAfter := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), big.NewInt(evmToCantonTransferAmount))
	require.Equal(t, expectedReceiverBalanceAfter, receiverBalanceAfter, "receiver final token balance should equal initial balance plus transfer amount")
	require.Equal(t, big.NewInt(evmToCantonTransferAmount), transferred, "receiver token balance should increase by transfer amount")
}
