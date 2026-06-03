package load

import (
	"fmt"
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

const (
	// cantonToEVMFeeAmount matches the message-only path in canton2evm_e2e_test.go.
	cantonToEVMFeeAmount int64 = 2_000

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

	configPath := "../../env-canton-evm-out.toml"
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("skipping Canton→EVM load test: %v (start devenv to generate %s)", err, configPath)
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

	destinations := discoverEVMDestinations(t, in, chainMap)
	require.NotEmpty(t, destinations, "need at least one EVM destination in the env file")
	t.Logf("Canton→EVM load destinations: %d EVM chain(s)", len(destinations))
	for _, d := range destinations {
		t.Logf("  - selector=%d receiver=%x", d.Chain.ChainSelector(), d.Receiver)
	}

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	ccvAddr, executorAddr := resolveCantonSourceAddrs(t, lib, cantonChain.ChainSelector())

	sched := loadSchedule(t)

	estimatedMessages := uint64(sched.rate) * uint64(sched.duration/sched.rateUnit)
	if estimatedMessages == 0 {
		estimatedMessages = 1
	}
	mintAmount := estimatedMessages * uint64(cantonToEVMFeeAmount) * mintBufferNumerator / mintBufferDenominator
	t.Logf("Pre-mint: estimatedMessages=%d feePerMessage=%d totalFeeMint=%d",
		estimatedMessages, cantonToEVMFeeAmount, mintAmount)
	require.NoError(t, cantonImpl.MintTokens(ctx, mintAmount))
	require.NoError(t, cantonImpl.SetupSend(ctx, uint64(cantonToEVMFeeAmount), 0))

	gun, err := NewCCIPLoadGun(
		cantonChain,
		destinations,
		ccvAddr,
		executorAddr,
		utilstests.WaitTimeout(t),
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-canton2evm", sched)
}
