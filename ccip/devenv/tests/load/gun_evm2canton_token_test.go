package load

import (
	"math/big"
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // registers Canton via init
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestEVM2Canton_TokenLoad runs WASP RPS=1 against the EVM→Canton token transfer path.
//
// Requires a running devenv and env-canton-evm-out.toml (devenv only).
//
//nolint:paralleltest // single-flight exec on Canton dest; shares env with e2e.
func TestEVM2Canton_TokenLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM→Canton token load test in short mode")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	boot := devenvtests.BootstrapE2E(t, env)
	boot.SkipIfRemote(t, "token load not on prod-testnet")

	ctx := ccv.Plog.WithContext(t.Context())
	boot.SetupCantonReceive(t, ctx)

	lane := devenvtests.ResolveTokenLane(t, boot.Cfg, boot.Lib, boot.ChainMap, boot.EVM.ChainSelector(), []uint64{boot.Canton.ChainSelector()})
	t.Logf("Token lane: pool=%s transfer=%s srcToken=%x",
		lane.PoolRef.Qualifier,
		lane.TransferAmount.String(),
		lane.SrcToken)

	receiverParticipant, _, err := boot.Canton.ClientParticipant()
	require.NoError(t, err)
	require.NotEmpty(t, receiverParticipant.PartyID)

	receiver, err := boot.Canton.GetEOAReceiverAddress()
	require.NoError(t, err)
	cantonDest := cantonTokenLoadDestination(boot.Canton, receiver, lane)

	sched := loadSchedule(t)
	estimatedMessages := estimateMessages(sched)
	evmSender, err := boot.EVM.GetEOAReceiverAddress()
	require.NoError(t, err)
	senderBalance, err := boot.EVM.GetTokenBalance(ctx, evmSender, lane.SrcToken)
	require.NoError(t, err)
	requiredBalance := new(big.Int).Mul(lane.TransferAmount, big.NewInt(int64(estimatedMessages)))
	t.Logf("EVM sender token balance=%s requiredForRun=%s (estimatedMessages=%d; devenv pre-funds sender)",
		senderBalance.String(), requiredBalance.String(), estimatedMessages)
	if senderBalance.Cmp(requiredBalance) < 0 {
		t.Logf("warning: EVM sender balance may be insufficient for full run")
	}

	ccvAddr, executorAddr := resolveEVMSourceAddrs(t, boot.Lib, boot.EVM.ChainSelector())

	gun, err := NewCCIPLoadGun(
		boot.EVM,
		[]Destination{cantonDest},
		ccvAddr,
		executorAddr,
		LoadGunOptions{
			ConfirmSend:        EVMSourceConfirmSend(boot),
			ConfirmExecTimeout: devenvtests.ConfirmExecTimeout(t),
		},
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-evm2canton-token", sched, "token_transfer")

	totalHoldingsRat, err := testhelpers.GetHoldingsBalance(ctx, receiverParticipant, nil)
	require.NoError(t, err)
	totalHoldingsFloat, _ := new(big.Float).SetRat(totalHoldingsRat).Float64()
	t.Logf("Canton receiver total holdings after load: %.10f (calls=%d)", totalHoldingsFloat, gun.CallCount())
}
