package load

import (
	"math/big"
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// TestCanton2EVM_TokenLoad runs WASP RPS=1 against the Canton→EVM token transfer path.
//
// Devenv: requires a running devenv and env-canton-evm-out.toml; pre-mints fee + transfer
// holdings via SetupCantonTokenSend before WASP starts.
//
// Prod-testnet: requires a pre-funded Canton party (Amulet + transfer instrument) and
// PRIVATE_KEY wallet with Sepolia ETH for execution gas. Full exec confirm; no EVM balance
// assert on prod — verify delivery via indexer/CCIP ops.
//
//nolint:paralleltest // Canton holdings must stay 1-wide; shares env with e2e.
func TestCanton2EVM_TokenLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton→EVM token load test in short mode")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	boot := devenvtests.BootstrapE2E(t, env)

	ctx := ccv.Plog.WithContext(t.Context())

	evmSelectors := discoverEVMTokenSelectors(t, boot.Cfg)
	require.NotEmpty(t, evmSelectors, "need at least one EVM token destination in the env file")
	lane := devenvtests.ResolveTokenLane(t, boot.Env, boot.Cfg, boot.Lib, boot.ChainMap, boot.Canton.ChainSelector(), evmSelectors)
	t.Logf("Token lane: pool=%s transfer=%s", lane.PoolRef.Qualifier, lane.TransferAmount.String())

	destinations := discoverEVMTokenDestinationsFromBoot(t, boot, lane)
	require.NotEmpty(t, destinations, "need at least one EVM token destination in the env file")
	t.Logf("Canton→EVM token load destinations: %d EVM chain(s)", len(destinations))
	for _, d := range destinations {
		t.Logf("  - selector=%d receiver=%x", d.Chain.ChainSelector(), d.Receiver)
	}

	evmReceiver := boot.ResolveEVMReceiver(t)
	firstDest := destinations[0]
	destToken := lane.DestTokenBySelector[firstDest.Chain.ChainSelector()]
	receiverBalanceBefore, err := firstDest.Chain.GetTokenBalance(ctx, evmReceiver, destToken)
	require.NoError(t, err)
	require.NotNil(t, receiverBalanceBefore)

	sched := loadSchedule(t)
	estimatedMessages := estimateMessages(sched)
	boot.SetupCantonTokenSend(t, ctx, lane, int(estimatedMessages))

	ccvAddr, executorAddr := resolveCantonSourceAddrs(t, boot.Lib, boot.Canton.ChainSelector())

	gun, err := NewCCIPLoadGun(
		boot.Canton,
		destinations,
		ccvAddr,
		executorAddr,
		LoadGunOptions{
			ConfirmSend:        CantonSourceConfirmSend(boot),
			ConfirmExecTimeout: devenvtests.ConfirmExecTimeout(t),
			SkipExecConfirm:    false,
		},
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-canton2evm-token", sched, "token_transfer", false, boot.Cfg.IndexerEndpoints)

	receiverBalanceAfter, err := firstDest.Chain.GetTokenBalance(ctx, evmReceiver, destToken)
	require.NoError(t, err)
	require.NotNil(t, receiverBalanceAfter)

	expectedPerMessage := new(big.Int).Mul(lane.TransferAmount, big.NewInt(devenv.CantonFixedPointToEVMScale))
	expectedDelta := new(big.Int).Mul(expectedPerMessage, big.NewInt(gun.CallCount()))
	t.Logf("EVM receiver token balance: before=%s after=%s expectedDelta=%s calls=%d",
		receiverBalanceBefore.String(), receiverBalanceAfter.String(), expectedDelta.String(), gun.CallCount())
	if !boot.Env.IsRemote() {
		expectedBalance := new(big.Int).Add(new(big.Int).Set(receiverBalanceBefore), expectedDelta)
		require.Equal(t, expectedBalance, receiverBalanceAfter)
	}
}
