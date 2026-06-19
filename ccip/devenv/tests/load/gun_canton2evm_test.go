package load

import (
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	utilstests "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
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
// Requires a running devenv and ../../env-canton-evm-out.toml (same as the basic e2e test).
//
// Devenv-specific: this test pre-mints fee holdings and calls SetupSend once before WASP
// starts. The CCIPLoadGun itself is environment-agnostic so the same gun can be reused by
// a future staging/prod runner that assumes pre-funded accounts.
//
//nolint:paralleltest // Canton holdings must stay 1-wide; shares env with e2e.
func TestCanton2EVM_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton→EVM load test in short mode")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	boot := devenvtests.BootstrapE2E(t, env)
	ctx := ccv.Plog.WithContext(t.Context())

	destinations := discoverEVMDestinations(t, boot.Cfg, boot.ChainMap)
	require.NotEmpty(t, destinations, "need at least one EVM destination in the env file")
	t.Logf("Canton→EVM load destinations: %d EVM chain(s)", len(destinations))
	for _, d := range destinations {
		t.Logf("  - selector=%d receiver=%x", d.Chain.ChainSelector(), d.Receiver)
	}

	ccvAddr, executorAddr := resolveCantonSourceAddrs(t, boot.Lib, boot.Canton.ChainSelector())

	sched := loadSchedule(t)

	estimatedMessages := uint64(sched.rate) * uint64(sched.duration/sched.rateUnit)
	if estimatedMessages == 0 {
		estimatedMessages = 1
	}
	mintAmount := estimatedMessages * uint64(devenvtests.CantonToEVMFeeAmount) * mintBufferNumerator / mintBufferDenominator
	t.Logf("Pre-mint: estimatedMessages=%d feePerMessage=%d totalFeeMint=%d",
		estimatedMessages, devenvtests.CantonToEVMFeeAmount, mintAmount)
	if !boot.Env.IsRemote() {
		require.NoError(t, boot.Canton.MintTokens(ctx, mintAmount))
	}
	require.NoError(t, boot.Canton.SetupSend(ctx, uint64(devenvtests.CantonToEVMFeeAmount), 0))

	gun, err := NewCCIPLoadGun(
		boot.Canton,
		destinations,
		ccvAddr,
		executorAddr,
		LoadGunOptions{
			ConfirmSend:        CantonSourceConfirmSend(boot),
			ConfirmExecTimeout: utilstests.WaitTimeout(t),
		},
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-canton2evm", sched, "message_only")
}
