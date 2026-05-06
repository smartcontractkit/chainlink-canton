package canton

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/proxy"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/versioned_verifier_resolver"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

const (
	// 1e11 (10-decimal units) gives a stable non-dust transfer in this lane after fee handling.
	evmToCantonTransferAmount = int64(100_000_000_000)
)

//nolint:paralleltest // we won't run this in parallel.
func TestEVM2Canton_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM2Canton_Basic test in short mode")
	}

	// Register the Canton impl factory so the shared CCV harness can resolve the
	// Canton family to this repo's devenv/test implementation.
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

	srcChain := devenvtests.GetChain(t, blockchain.TypeAnvil, in, harness)
	dstChain := devenvtests.GetChain(t, blockchain.TypeCanton, in, harness)

	for _, client := range harness.AggregatorClients {
		t.Cleanup(func() {
			client.Close()
		})
	}

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	srcSelector := srcChain.ChainSelector()
	dstSelector := dstChain.ChainSelector()
	_, opsEnv, err := ccv.NewCLDFOperationsEnvironment(in.Blockchains, in.CLDF.DataStore)
	require.NoError(t, err)
	var receiverParticipant canton.Participant
	if chains := opsEnv.BlockChains.CantonChains(); len(chains[dstSelector].Participants) > 0 {
		receiverParticipant = chains[dstSelector].Participants[0]
	}
	require.NotEmpty(t, receiverParticipant.PartyID)

	receiver, err := dstChain.GetEOAReceiverAddress()
	require.NoError(t, err)

	ccvAddr, err := tcapi.GetContractAddress(
		in,
		srcSelector,
		datastore.ContractType(versioned_verifier_resolver.CommitteeVerifierResolverType),
		versioned_verifier_resolver.Version.String(),
		common.DefaultCommitteeVerifierQualifier,
		"source committee verifier",
	)
	require.NoError(t, err)

	executorAddress, err := tcapi.GetContractAddress(
		in,
		srcSelector,
		datastore.ContractType(sequences.ExecutorProxyType),
		proxy.Deploy.Version(),
		common.DefaultExecutorQualifier,
		"source executor",
	)
	require.NoError(t, err)

	t.Run("message transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		seqNo, err := srcChain.GetExpectedNextSequenceNumber(subtestCtx, dstSelector)
		require.NoError(t, err)
		sendMessageResult, err := srcChain.SendMessage(subtestCtx, dstSelector, cciptestinterfaces.MessageFields{
			Receiver: receiver,
			Data:     []byte("Hello message transfer from EVM!"),
		}, cciptestinterfaces.MessageOptions{
			ExecutionGasLimit: 200_000,
			FinalityConfig:    0,
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

		sentEvent, err := srcChain.ConfirmSendOnSource(subtestCtx, dstSelector, cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, 15*time.Second)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.Nil(t, sentEvent.Message.TokenTransfer)

		result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, &harness, sentEvent.MessageID)
		vr := result.IndexedVerifications.Results[0].VerifierResult
		message, verifierDestAddress, ccvData := vr.Message, vr.VerifierDestAddress, vr.CCVData
		require.Nil(t, message.TokenTransfer)
		executionStateChangedEvent, err := dstChain.ManuallyExecuteMessage(subtestCtx, message, 0, []protocol.UnknownAddress{verifierDestAddress}, [][]byte{ccvData})
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State)
	})

	t.Run("token transfer", func(t *testing.T) {
		subtestCtx := ccv.Plog.WithContext(t.Context())

		tokenRef, err := in.CLDF.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				srcSelector,
				datastore.ContractType("BurnMintERC20WithDripToken"),
				semver.MustParse("1.0.0"),
				burnMint20ToLockRelease20TokenQualifier(t),
			),
		)
		require.NoError(t, err)
		srcToken := protocol.UnknownAddress(gethcommon.HexToAddress(tokenRef.Address).Bytes())
		srcSender, err := srcChain.GetEOAReceiverAddress()
		require.NoError(t, err)
		seqNo, err := srcChain.GetExpectedNextSequenceNumber(subtestCtx, dstSelector)
		require.NoError(t, err)
		sendMessageResult, err := srcChain.SendMessage(subtestCtx, dstSelector, cciptestinterfaces.MessageFields{
			Receiver: receiver,
			Data:     []byte("Hello token transfer from EVM!"),
			TokenAmount: cciptestinterfaces.TokenAmount{
				Amount:       big.NewInt(evmToCantonTransferAmount),
				TokenAddress: srcToken,
			},
		}, cciptestinterfaces.MessageOptions{
			ExecutionGasLimit: 200_000,
			FinalityConfig:    0,
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

		sentEvent, err := srcChain.ConfirmSendOnSource(subtestCtx, dstSelector, cciptestinterfaces.MessageEventKey{SeqNum: seqNo}, 15*time.Second)
		require.NoError(t, err)
		require.NotNil(t, sentEvent.Message)
		require.NotNil(t, sentEvent.Message.TokenTransfer)

		result := devenvtests.AssertSingleVerifierResult(t, subtestCtx, &harness, sentEvent.MessageID)
		vr := result.IndexedVerifications.Results[0].VerifierResult
		message, verifierDestAddress, ccvData := vr.Message, vr.VerifierDestAddress, vr.CCVData
		require.NotNil(t, message.TokenTransfer)
		require.NotNil(t, message.TokenTransfer.Amount)
		t.Logf("Canton token transfer amount from verifier result: %s", message.TokenTransfer.Amount.String())
		require.Positive(t, message.TokenTransfer.Amount.Cmp(big.NewInt(0)), "token transfer amount must be positive")
		executionStateChangedEvent, err := dstChain.ManuallyExecuteMessage(subtestCtx, message, 0, []protocol.UnknownAddress{verifierDestAddress}, [][]byte{ccvData})
		require.NoError(t, err)
		require.Equal(t, cciptestinterfaces.ExecutionStateSuccess, executionStateChangedEvent.State)

		totalHoldingsRat, err := testhelpers.GetHoldingsBalance(subtestCtx, receiverParticipant, nil)
		require.NoError(t, err)
		totalHoldingsFloat, _ := new(big.Float).SetRat(totalHoldingsRat).Float64()
		t.Logf("Canton receiver total holdings after execute: %.10f", totalHoldingsFloat)

		srcBalanceAfter, err := srcChain.GetTokenBalance(subtestCtx, srcSender, srcToken)
		require.NoError(t, err)
		require.NotNil(t, srcBalanceAfter)
		dstBalanceAfter, err := dstChain.GetTokenBalance(subtestCtx, receiver, nil)
		require.NoError(t, err)
		require.NotNil(t, dstBalanceAfter)
		t.Logf("Token balances after execute: evm_sender=%s canton_receiver=%s", srcBalanceAfter.String(), dstBalanceAfter.String())
	})
}
