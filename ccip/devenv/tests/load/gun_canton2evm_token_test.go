package load

import (
	"math/big"
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // registers Canton via init
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// TestCanton2EVM_TokenLoad runs WASP RPS=1 against the Canton→EVM token transfer path.
//
// Requires a running devenv and env-canton-evm-out.toml (devenv only).
//
//nolint:paralleltest // Canton holdings must stay 1-wide; shares env with e2e.
func TestCanton2EVM_TokenLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton→EVM token load test in short mode")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	boot := devenvtests.BootstrapE2E(t, env)
	boot.SkipIfRemote(t, "token load not on prod-testnet")

	ctx := ccv.Plog.WithContext(t.Context())

	evmSelectors := discoverEVMTokenSelectors(t, boot.Cfg)
	require.NotEmpty(t, evmSelectors, "need at least one EVM token destination in the env file")
	lane := devenvtests.ResolveTokenLane(t, boot.Cfg, boot.Lib, boot.ChainMap, boot.Canton.ChainSelector(), evmSelectors)
	t.Logf("Token lane: pool=%s transfer=%s", lane.PoolRef.Qualifier, lane.TransferAmount.String())

	destinations := discoverEVMTokenDestinations(t, boot.Cfg, boot.ChainMap, lane)
	require.NotEmpty(t, destinations, "need at least one EVM token destination in the env file")
	t.Logf("Canton→EVM token load destinations: %d EVM chain(s)", len(destinations))

	firstDest := destinations[0]
	destToken := lane.DestTokenBySelector[firstDest.Chain.ChainSelector()]
	receiverBalanceBefore, err := firstDest.Chain.GetTokenBalance(ctx, firstDest.Receiver, destToken)
	require.NoError(t, err)
	require.NotNil(t, receiverBalanceBefore)

	ccvAddr, executorAddr := resolveCantonSourceAddrs(t, boot.Lib, boot.Canton.ChainSelector())
	sched := loadSchedule(t)

	setupCantonTokenLoadHoldings(t, ctx, boot.Canton, sched, lane)

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

	runWASP(t, gun, "canton-load-canton2evm-token", sched, "token_transfer", false)

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
