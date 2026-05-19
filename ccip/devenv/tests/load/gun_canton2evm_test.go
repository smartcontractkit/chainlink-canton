package load

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/chainimpl"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	ccvload "github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/load"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	utilstests "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
)

const (
	// cantonToEVMFeeAmount matches the message-only path in canton2evm_e2e_test.go.
	cantonToEVMFeeAmount int64 = 2_000

	// Schedule env vars (consumed in TestCanton2EVM_Load):
	//   CANTON_LOAD_MESSAGE_RATE: "<int>/<duration>" (e.g. "1/1s", "1/20s", "10/5m")
	//     Maps to wasp.Plain rate + wasp.Config.RateLimitUnitDuration.
	//   CANTON_LOAD_DURATION: total wall-clock runtime (Go time.Duration, e.g. "90s", "10m")
	//     Maps to the duration argument of wasp.Plain.
	envMessageRate      = "CANTON_LOAD_MESSAGE_RATE"
	envLoadDuration     = "CANTON_LOAD_DURATION"
	defaultMessageRate  = "1/1s"
	defaultLoadDuration = 90 * time.Second

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
// starts. The Canton2EVMGun itself is environment-agnostic so the same gun can be reused by
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

	chainimpl.RegisterImplFactory(chainsel.FamilyCanton, cantondevenv.NewImplFactory())

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

	cantonChain := devenvtests.GetChain(t, blockchain.TypeCanton, in, harness)
	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	require.True(t, ok, "Canton chain must be *cantondevenv.Chain")

	destinations := discoverEVMDestinations(t, ctx, in, harness)
	require.NotEmpty(t, destinations, "need at least one EVM destination in the env file")
	t.Logf("Canton→EVM load destinations: %d EVM chain(s)", len(destinations))
	for _, d := range destinations {
		t.Logf("  - selector=%d receiver=%x", d.Chain.ChainSelector(), d.Receiver)
	}

	for _, client := range harness.AggregatorClients {
		t.Cleanup(func() {
			client.Close()
		})
	}

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	ccvAddr, err := tcapi.GetContractAddress(
		in,
		cantonChain.ChainSelector(),
		datastore.ContractType(canton_committee_verifier.ContractType),
		canton_committee_verifier.Version.String(),
		common.DefaultCommitteeVerifierQualifier,
		"canton committee verifier",
	)
	require.NoError(t, err)
	executorAddr, err := tcapi.GetContractAddress(
		in,
		cantonChain.ChainSelector(),
		datastore.ContractType(executor.ContractType),
		executor.Version.String(),
		common.DefaultExecutorQualifier,
		"source executor",
	)
	require.NoError(t, err)

	// Resolve the load schedule from env vars (or defaults). Same "N/T" shape as CCV.
	messageRate := os.Getenv(envMessageRate)
	if messageRate == "" {
		messageRate = defaultMessageRate
	}
	rate, rateUnit := ccvload.ParseMessageRate(messageRate)
	require.NotZero(t, rate, "%s=%q invalid (expected e.g. '1/20s')", envMessageRate, messageRate)
	require.Positive(t, rateUnit, "%s=%q invalid (expected e.g. '1/20s')", envMessageRate, messageRate)

	duration := defaultLoadDuration
	if d := os.Getenv(envLoadDuration); d != "" {
		parsed, err := time.ParseDuration(d)
		require.NoError(t, err, "%s=%q invalid", envLoadDuration, d)
		duration = parsed
	}
	t.Logf("Load schedule: rate=%d unit=%s totalDuration=%s", rate, rateUnit, duration)

	// Devenv-only pre-funding: mint enough fee holdings for the whole profile, then SetupSend
	// once. The gun does not mint during Call (staging/prod rely on pre-funded accounts).
	estimatedMessages := uint64(rate) * uint64(duration/rateUnit)
	if estimatedMessages == 0 {
		estimatedMessages = 1
	}
	mintAmount := estimatedMessages * uint64(cantonToEVMFeeAmount) * mintBufferNumerator / mintBufferDenominator
	t.Logf("Pre-mint: estimatedMessages=%d feePerMessage=%d totalFeeMint=%d",
		estimatedMessages, cantonToEVMFeeAmount, mintAmount)
	require.NoError(t, cantonImpl.MintTokens(ctx, mintAmount))
	require.NoError(t, cantonImpl.SetupSend(ctx, uint64(cantonToEVMFeeAmount), 0))

	gun, err := NewCanton2EVMGun(
		&harness,
		cantonChain,
		destinations,
		ccvAddr,
		executorAddr,
		utilstests.WaitTimeout(t),
	)
	require.NoError(t, err)

	p := wasp.NewProfile().Add(wasp.NewGenerator(&wasp.Config{
		T:        t,
		LoadType: wasp.RPS,
		GenName:  "canton-load-canton2evm",
		Schedule: wasp.Combine(
			wasp.Plain(rate, duration),
		),
		RateLimitUnitDuration: rateUnit,
		Gun:                   gun,
		Labels: map[string]string{
			"go_test_name":  "canton-load-canton2evm",
			"branch":        "test",
			"commit":        "test",
			"message_rate":  messageRate,
			"load_duration": duration.String(),
		},
		LokiConfig: nil,
	}))

	_, err = p.Run(true)
	require.NoError(t, err)
	p.Wait()

	require.Greater(t, gun.CallCount(), int64(0), "gun should have completed at least one message")
	require.LessOrEqual(t, gun.MaxConcurrentObserved(), int32(1),
		"Gun.Call must not overlap (Canton holdings 1-wide)")
}

// discoverEVMDestinations enumerates every Anvil blockchain in the env, resolves its CCIP17
// implementation via the harness's chain map, and captures an EOA receiver per chain.
func discoverEVMDestinations(t *testing.T, ctx context.Context, in *ccv.Cfg, harness tcapi.TestHarness) []EVMDestination {
	t.Helper()
	chainMap, err := harness.Lib.ChainsMap(ctx)
	require.NoError(t, err)

	dests := make([]EVMDestination, 0)
	seen := make(map[uint64]struct{})
	for _, bc := range in.Blockchains {
		if bc.Type != blockchain.TypeAnvil {
			continue
		}
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(bc.ChainID, chainsel.FamilyEVM)
		require.NoError(t, err, "resolve chain selector for chainID=%s", bc.ChainID)
		if _, dup := seen[details.ChainSelector]; dup {
			continue
		}
		chain, ok := chainMap[details.ChainSelector]
		require.True(t, ok, "EVM chain %d not in harness chain map", details.ChainSelector)

		receiver, err := chain.GetEOAReceiverAddress()
		require.NoError(t, err)

		dests = append(dests, EVMDestination{Chain: chain, Receiver: receiver})
		seen[details.ChainSelector] = struct{}{}
	}
	return dests
}
