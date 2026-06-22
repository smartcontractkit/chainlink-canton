package load

import (
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // registers Canton via init
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

const (
	// mintBuffer is a safety multiplier on top of estimated message count so a slow fee
	// burn does not starve the run mid-flight (1.5x).
	mintBufferNumerator   uint64 = 3
	mintBufferDenominator uint64 = 2
)

// TestCanton2EVM_Load runs WASP RPS=1 against the real Canton→EVM path (message-only),
// round-robining across every EVM destination found in the env file.
//
// Devenv: requires a running devenv and env-canton-evm-out.toml; pre-mints fee holdings
// and calls SetupSend once before WASP starts.
//
// Prod-testnet: send-only load (Canton send + confirm send, no ConfirmExecOnDest on EVM).
// Set CANTON_GRPC_URL, CANTON_PARTY_ID, CANTON_AUTH_*, PRIVATE_KEY (EVM message receiver),
// CANTON_LOAD_SKIP_EXEC_CONFIRM=true, and pre-fund the Canton party (~50 Amulet per message).
// Verify delivery via indexer/CCIP ops — the test does not assert EVM execution on prod.
//
//nolint:paralleltest // Canton holdings must stay 1-wide; shares env with e2e.
func TestCanton2EVM_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton→EVM load test in short mode")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	t.Logf("env: %s", env)
	boot := devenvtests.BootstrapE2E(t, env)
	ctx := ccv.Plog.WithContext(t.Context())

	skipExec := loadSkipExecConfirm(t)
	// if boot.Env.IsRemote() && !skipExec {
	// 	t.Skip("prod-testnet requires CANTON_LOAD_SKIP_EXEC_CONFIRM=true (EVM executor not available)")
	// }

	destinations := discoverEVMDestinationsFromBoot(t, boot)
	require.NotEmpty(t, destinations, "need at least one EVM destination in the env file")
	t.Logf("Canton→EVM load destinations: %d EVM chain(s)", len(destinations))
	for _, d := range destinations {
		t.Logf("  - selector=%d receiver=%x", d.Chain.ChainSelector(), d.Receiver)
	}

	ccvAddr, executorAddr := resolveCantonSourceAddrs(t, boot.Lib, boot.Canton.ChainSelector())

	sched := loadSchedule(t)

	estimatedMessages := estimateMessages(sched)
	requiredAmulet := estimatedMessages * uint64(devenvtests.CantonToEVMFeeAmount) * mintBufferNumerator / mintBufferDenominator
	if boot.Env.IsRemote() {
		t.Logf("Prod: ensure Canton party holds at least %d Amulet (estimatedMessages=%d feePerMessage=%d)",
			requiredAmulet, estimatedMessages, devenvtests.CantonToEVMFeeAmount)
	} else {
		t.Logf("Pre-mint: estimatedMessages=%d feePerMessage=%d totalFeeMint=%d",
			estimatedMessages, devenvtests.CantonToEVMFeeAmount, requiredAmulet)
		require.NoError(t, boot.Canton.MintTokens(ctx, requiredAmulet))
	}
	require.NoError(t, boot.Canton.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), 0))

	gun, err := NewCCIPLoadGun(
		boot.Canton,
		destinations,
		ccvAddr,
		executorAddr,
		LoadGunOptions{
			ConfirmSend:        CantonSourceConfirmSend(boot),
			ConfirmExecTimeout: devenvtests.ConfirmExecTimeout(t),
			SkipExecConfirm:    skipExec,
		},
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-canton2evm", sched, "message_only", skipExec, boot.Cfg.IndexerEndpoints)
}
