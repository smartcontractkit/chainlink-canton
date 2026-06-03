package load

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	utilstests "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // registers Canton via init
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// TestCanton2EVM_TokenLoad runs WASP RPS=1 against the Canton→EVM token transfer path.
//
// Requires a running devenv and ../../env-canton-evm-out.toml.
//
//nolint:paralleltest // Canton holdings must stay 1-wide; shares env with e2e.
func TestCanton2EVM_TokenLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton→EVM token load test in short mode")
	}

	configPath := "../../env-canton-evm-out.toml"
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("skipping Canton→EVM token load test: %v (start devenv to generate %s)", err, configPath)
	}

	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	ctx := ccv.Plog.WithContext(t.Context())
	lib, err := ccv.NewLibFromCCVEnv(&ccv.Plog, configPath, chainsel.FamilyEVM, chainsel.FamilyCanton)
	require.NoError(t, err)

	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)
	require.NoError(t, devenvtests.WireVerifierObservationFromLib(lib, chainMap))

	cantonChain := devenvtests.GetChainFromMap(t, blockchain.TypeCanton, in, chainMap)
	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	require.True(t, ok, "Canton chain must be *cantondevenv.Chain")

	tokenInput := devenvtests.LoadTokenTransferInput(t, devenvtests.DirectionCantonToEVM)
	evmSelectors := discoverEVMTokenSelectors(t, in)
	require.NotEmpty(t, evmSelectors, "need at least one EVM token destination in the env file")
	lane := devenvtests.ResolveTokenLane(t, in, lib, chainMap, cantonChain.ChainSelector(), evmSelectors, tokenInput)
	t.Logf("Token lane: pool=%s transfer=%s", lane.PoolRef.Qualifier, lane.TransferAmount.String())

	destinations := discoverEVMTokenDestinations(t, in, chainMap, lane)
	require.NotEmpty(t, destinations, "need at least one EVM token destination in the env file")
	t.Logf("Canton→EVM token load destinations: %d EVM chain(s)", len(destinations))

	firstDest := destinations[0]
	destToken := lane.DestTokenBySelector[firstDest.Chain.ChainSelector()]
	receiverBalanceBefore, err := firstDest.Chain.GetTokenBalance(ctx, firstDest.Receiver, destToken)
	require.NoError(t, err)
	require.NotNil(t, receiverBalanceBefore)

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	ccvAddr, executorAddr := resolveCantonSourceAddrs(t, lib, cantonChain.ChainSelector())
	sched := loadSchedule(t)

	estimatedMessages := estimateMessages(sched)
	feeMint, transferMint := cantonTokenPreMintAmounts(estimatedMessages, lane)
	t.Logf("Pre-mint: estimatedMessages=%d feeMint=%d transferMint=%d",
		estimatedMessages, feeMint, transferMint)
	require.NoError(t, cantonImpl.MintTokens(ctx, feeMint))
	require.NoError(t, cantonImpl.MintTokens(ctx, transferMint))
	require.NoError(t, cantonImpl.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), lane.TransferAmount.Uint64()))

	gun, err := NewCCIPLoadGun(
		cantonChain,
		destinations,
		ccvAddr,
		executorAddr,
		utilstests.WaitTimeout(t),
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-canton2evm-token", sched, "token_transfer")

	receiverBalanceAfter, err := firstDest.Chain.GetTokenBalance(ctx, firstDest.Receiver, destToken)
	require.NoError(t, err)
	require.NotNil(t, receiverBalanceAfter)

	expectedPerMessage := new(big.Int).Mul(lane.TransferAmount, big.NewInt(devenvtests.EVMDecimalsScale))
	expectedDelta := new(big.Int).Mul(expectedPerMessage, big.NewInt(gun.CallCount()))
	expectedBalance := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), expectedDelta)
	t.Logf("EVM receiver token balance: before=%s after=%s expectedDelta=%s calls=%d",
		receiverBalanceBefore.String(), receiverBalanceAfter.String(), expectedDelta.String(), gun.CallCount())
	require.Equal(t, expectedBalance, receiverBalanceAfter)
}
