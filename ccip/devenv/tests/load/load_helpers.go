package load

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/proxy"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/versioned_verifier_resolver"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	ccvload "github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/load"
	ccvmetrics "github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/metrics"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
)

const (
	envMessageRate         = "CANTON_LOAD_MESSAGE_RATE"
	envLoadDuration        = "CANTON_LOAD_DURATION"
	envLoadCallTimeout     = "CANTON_LOAD_CALL_TIMEOUT"
	defaultMessageRate     = "1/1s"
	defaultLoadDuration    = 90 * time.Second
	defaultLoadCallPadding = 2 * time.Minute
	confirmSendTimeout     = 30 * time.Second
)

type scheduleConfig struct {
	messageRate string
	rate        int64
	rateUnit    time.Duration
	duration    time.Duration
}

func loadSchedule(t *testing.T) scheduleConfig {
	t.Helper()

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

	return scheduleConfig{
		messageRate: messageRate,
		rate:        rate,
		rateUnit:    rateUnit,
		duration:    duration,
	}
}

func waspCallTimeout(t *testing.T, gun *CCIPLoadGun, sched scheduleConfig) time.Duration {
	t.Helper()

	if v := strings.TrimSpace(os.Getenv(envLoadCallTimeout)); v != "" {
		parsed, err := time.ParseDuration(v)
		require.NoError(t, err, "%s=%q invalid", envLoadCallTimeout, v)
		return parsed
	}

	return gun.ConfirmExecTimeout() + sched.rateUnit + defaultLoadCallPadding
}

func printLoadMetrics(t *testing.T, gun *CCIPLoadGun) {
	t.Helper()

	records := gun.Metrics()
	failures := gun.FailureCounts()
	PrintPhaseMetricsSummary(t, records, failures, false)

	ccvMetrics := ToCCVMessageMetrics(records)
	if len(ccvMetrics) > 0 {
		totals := ccvmetrics.MessageTotals{
			Sent:     len(ccvMetrics),
			Received: len(ccvMetrics),
		}
		summary := ccvmetrics.CalculateMetricsSummary(ccvMetrics, totals)
		ccvmetrics.PrintMetricsSummary(t, summary)
	}
}

func logLoadMessageSummary(t *testing.T, gun *CCIPLoadGun) {
	t.Helper()

	ids := gun.MessageIDs()
	lggr := ccv.Plog
	lggr.Info().Int("count", len(ids)).Msg("Load message summary")

	for i, id := range ids {
		lggr.Info().Int("index", i+1).Str("messageID", id.String()).Msg("Load message sent")
	}
}

func runWASP(t *testing.T, gun *CCIPLoadGun, genName string, sched scheduleConfig, scenario string) {
	t.Helper()
	defer logLoadMessageSummary(t, gun)

	callTimeout := waspCallTimeout(t, gun, sched)
	ccv.Plog.Info().
		Str("messageRate", sched.messageRate).
		Dur("rateUnit", sched.rateUnit).
		Dur("loadDuration", sched.duration).
		Dur("callTimeout", callTimeout).
		Dur("confirmExecTimeout", gun.ConfirmExecTimeout()).
		Msg("WASP load schedule")

	labels := map[string]string{
		"go_test_name":  genName,
		"branch":        "test",
		"commit":        "test",
		"message_rate":  sched.messageRate,
		"load_duration": sched.duration.String(),
	}
	if scenario != "" {
		labels["scenario"] = scenario
	}

	p := wasp.NewProfile().Add(wasp.NewGenerator(&wasp.Config{
		T:        t,
		LoadType: wasp.RPS,
		GenName:  genName,
		Schedule: wasp.Combine(
			wasp.Plain(sched.rate, sched.duration),
		),
		RateLimitUnitDuration: sched.rateUnit,
		CallTimeout:           callTimeout,
		Gun:                   gun,
		Labels:                labels,
		LokiConfig:            nil,
	}))

	_, err := p.Run(true)
	require.NoError(t, err)
	p.Wait()

	require.Positive(t, gun.CallCount(), "gun should have completed at least one message")
	require.LessOrEqual(t, gun.MaxConcurrentObserved(), int32(1),
		"Gun.Call must not overlap (single-flight)")

	printLoadMetrics(t, gun)
}

func discoverEVMDestinations(t *testing.T, in *ccv.Cfg, chainMap map[uint64]cciptestinterfaces.CCIP17) []Destination {
	t.Helper()

	dests := make([]Destination, 0)
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

		dests = append(dests, evmLoadDestination(chain, receiver))
		seen[details.ChainSelector] = struct{}{}
	}

	return dests
}

// discoverEVMTokenSelectors returns the chain selectors of every EVM chain in the
// env file. Callers resolve the token lane over these selectors before building
// destinations.
func discoverEVMTokenSelectors(t *testing.T, in *ccv.Cfg) []uint64 {
	t.Helper()

	selectors := make([]uint64, 0)
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
		selectors = append(selectors, details.ChainSelector)
		seen[details.ChainSelector] = struct{}{}
	}

	return selectors
}

func discoverEVMTokenDestinations(
	t *testing.T,
	in *ccv.Cfg,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
	lane devenvtests.TokenLane,
) []Destination {
	t.Helper()

	selectors := discoverEVMTokenSelectors(t, in)
	dests := make([]Destination, 0, len(selectors))
	for _, selector := range selectors {
		chain, ok := chainMap[selector]
		require.True(t, ok, "EVM chain %d not in harness chain map", selector)

		receiver, err := chain.GetEOAReceiverAddress()
		require.NoError(t, err)

		dests = append(dests, evmTokenLoadDestination(chain, receiver, lane))
	}

	return dests
}

func estimateMessages(sched scheduleConfig) uint64 {
	// rate and rateUnit are validated positive in loadSchedule.
	if sched.rate <= 0 {
		return 1
	}
	estimated := uint64(sched.rate) * uint64(sched.duration/sched.rateUnit) //nolint:gosec // rate > 0
	if estimated == 0 {
		return 1
	}

	return estimated
}

// setupCantonTokenLoadHoldings pre-mints two separate Amulet holdings (fee + transfer) and
// calls SetupSend once, matching the e2e canton2evm token transfer pattern.
func setupCantonTokenLoadHoldings(
	t *testing.T,
	ctx context.Context,
	cantonImpl *cantondevenv.Chain,
	sched scheduleConfig,
	lane devenvtests.TokenLane,
) {
	t.Helper()

	estimated := estimateMessages(sched)
	fee := uint64(cantondevenv.CantonToEVMTokenTransferFeeAmount)
	feeTotal := new(big.Rat).SetUint64(estimated * fee)
	transferTotalFP := new(big.Int).Mul(lane.TransferAmount, big.NewInt(int64(estimated)))
	transferTotal := new(big.Rat).SetFrac(transferTotalFP, big.NewInt(cantondevenv.CantonFixedPointScale))
	transferPerSend := new(big.Rat).SetFrac(lane.TransferAmount, big.NewInt(cantondevenv.CantonFixedPointScale))
	t.Logf("Pre-mint: estimatedMessages=%d feeTotal=%s transferTotal=%s",
		estimated, feeTotal.FloatString(10), transferTotal.FloatString(10))
	require.NoError(t, cantonImpl.MintTokens(ctx, feeTotal))
	require.NoError(t, cantonImpl.MintTokens(ctx, transferTotal))
	if lane.TransferInstrument.Admin != "" {
		require.NoError(t, cantonImpl.SetupSend(ctx, fee, transferPerSend, lane.TransferInstrument))
	} else {
		require.NoError(t, cantonImpl.SetupSend(ctx, fee, transferPerSend))
	}
}

func evmLoadDestination(chain cciptestinterfaces.CCIP17, receiver protocol.UnknownAddress) Destination {
	destSelector := chain.ChainSelector()
	return Destination{
		Chain:    chain,
		Receiver: receiver,
		buildMessage: func(_ cciptestinterfaces.CCIP17, callNum int64, ccvAddr, executorAddr protocol.UnknownAddress) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error) {
			return cciptestinterfaces.MessageFields{
					Receiver: receiver,
					Data:     fmt.Appendf(nil, "canton2evm load n=%d dest=%d", callNum, destSelector),
				}, cciptestinterfaces.MessageOptions{
					ExecutionGasLimit: 200_000,
					FinalityConfig:    1,
					Executor:          executorAddr,
					CCVs: []protocol.CCV{
						{CCVAddress: ccvAddr, Args: []byte{}, ArgsLen: 0},
					},
				}, nil
		},
	}
}

func evmTokenLoadDestination(chain cciptestinterfaces.CCIP17, receiver protocol.UnknownAddress, lane devenvtests.TokenLane) Destination {
	destSelector := chain.ChainSelector()
	laneCopy := lane
	return Destination{
		Chain:     chain,
		Receiver:  receiver,
		TokenLane: &laneCopy,
		buildMessage: func(_ cciptestinterfaces.CCIP17, callNum int64, ccvAddr, executorAddr protocol.UnknownAddress) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error) {
			return cciptestinterfaces.MessageFields{
					Receiver: receiver,
					Data:     fmt.Appendf(nil, "canton2evm token load n=%d dest=%d", callNum, destSelector),
					TokenAmount: cciptestinterfaces.TokenAmount{
						Amount: new(big.Int).Set(lane.TransferAmount),
					},
				}, cciptestinterfaces.MessageOptions{
					ExecutionGasLimit: lane.ExecutionGasLimit,
					FinalityConfig:    lane.FinalityConfig,
					Executor:          executorAddr,
					CCVs: []protocol.CCV{
						{CCVAddress: ccvAddr, Args: []byte{}, ArgsLen: 0},
					},
				}, nil
		},
	}
}

func discoverCantonDest(t *testing.T, in *ccv.Cfg, chainMap map[uint64]cciptestinterfaces.CCIP17) Destination {
	t.Helper()

	chain := devenvtests.GetChainFromMap(t, blockchain.TypeCanton, in, chainMap)
	receiver, err := chain.GetEOAReceiverAddress()
	require.NoError(t, err)

	return cantonLoadDestination(chain, receiver)
}

func cantonLoadDestination(chain cciptestinterfaces.CCIP17, receiver protocol.UnknownAddress) Destination {
	destSelector := chain.ChainSelector()
	return Destination{
		Chain:    chain,
		Receiver: receiver,
		buildMessage: func(_ cciptestinterfaces.CCIP17, callNum int64, ccvAddr, executorAddr protocol.UnknownAddress) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error) {
			return cciptestinterfaces.MessageFields{
					Receiver: receiver,
					Data:     fmt.Appendf(nil, "evm2canton load n=%d dest=%d", callNum, destSelector),
				}, cciptestinterfaces.MessageOptions{
					ExecutionGasLimit: 200_000,
					FinalityConfig:    0,
					Executor:          executorAddr,
					CCVs: []protocol.CCV{
						{CCVAddress: ccvAddr, Args: []byte{}, ArgsLen: 0},
					},
				}, nil
		},
	}
}

func cantonTokenLoadDestination(chain cciptestinterfaces.CCIP17, receiver protocol.UnknownAddress, lane devenvtests.TokenLane) Destination {
	destSelector := chain.ChainSelector()
	laneCopy := lane
	return Destination{
		Chain:     chain,
		Receiver:  receiver,
		TokenLane: &laneCopy,
		buildMessage: func(_ cciptestinterfaces.CCIP17, callNum int64, ccvAddr, executorAddr protocol.UnknownAddress) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error) {
			return cciptestinterfaces.MessageFields{
					Receiver: receiver,
					Data:     fmt.Appendf(nil, "evm2canton token load n=%d dest=%d", callNum, destSelector),
					TokenAmount: cciptestinterfaces.TokenAmount{
						Amount:       new(big.Int).Set(lane.TransferAmount),
						TokenAddress: lane.SrcToken,
					},
				}, cciptestinterfaces.MessageOptions{
					ExecutionGasLimit: lane.ExecutionGasLimit,
					FinalityConfig:    lane.FinalityConfig,
					Executor:          executorAddr,
					CCVs: []protocol.CCV{
						{CCVAddress: ccvAddr, Args: []byte{}, ArgsLen: 0},
					},
				}, nil
		},
	}
}

func resolveCantonSourceAddrs(t *testing.T, lib ccv.Lib, cantonSelector uint64) (protocol.UnknownAddress, protocol.UnknownAddress) {
	t.Helper()

	ds, err := lib.DataStore()
	require.NoError(t, err)

	ccvAddr := devenvtests.GetContractAddress(
		t, ds, cantonSelector,
		datastore.ContractType(canton_committee_verifier.ContractType),
		canton_committee_verifier.Version.String(),
		common.DefaultCommitteeVerifierQualifier,
		"canton committee verifier",
	)
	executorAddr := devenvtests.GetContractAddress(
		t, ds, cantonSelector,
		datastore.ContractType(executor.ContractType),
		executor.Version.String(),
		common.DefaultExecutorQualifier,
		"source executor",
	)

	return ccvAddr, executorAddr
}

func resolveEVMSourceAddrs(t *testing.T, lib ccv.Lib, evmSelector uint64) (protocol.UnknownAddress, protocol.UnknownAddress) {
	t.Helper()

	ds, err := lib.DataStore()
	require.NoError(t, err)

	ccvAddr := devenvtests.GetContractAddress(
		t, ds, evmSelector,
		datastore.ContractType(versioned_verifier_resolver.CommitteeVerifierResolverType),
		versioned_verifier_resolver.Version.String(),
		common.DefaultCommitteeVerifierQualifier,
		"source committee verifier",
	)
	executorAddr := devenvtests.GetContractAddress(
		t, ds, evmSelector,
		datastore.ContractType(sequences.ExecutorProxyType),
		proxy.Deploy.Version(),
		common.DefaultExecutorQualifier,
		"source executor",
	)

	return ccvAddr, executorAddr
}

func cantonSourceConfirmSend(source cciptestinterfaces.CCIP17) ConfirmSendFunc {
	return func(
		t *testing.T,
		ctx context.Context,
		destSelector uint64,
		seqNo uint64,
		_ cciptestinterfaces.MessageSentEvent,
	) (cciptestinterfaces.MessageSentEvent, error) {
		return source.ConfirmSendOnSource(
			ctx,
			destSelector,
			cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
			confirmSendTimeout,
		)
	}
}

// TODO: this is needed because EVM impls method has a bug on prod-testnet checks.
// currently just a simple wrapper around the method.
func evmSourceConfirmSend(source cciptestinterfaces.CCIP17) ConfirmSendFunc {
	return func(
		t *testing.T,
		ctx context.Context,
		destSelector uint64,
		seqNo uint64,
		_ cciptestinterfaces.MessageSentEvent,
	) (cciptestinterfaces.MessageSentEvent, error) {
		return source.ConfirmSendOnSource(
			ctx,
			destSelector,
			cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
			confirmSendTimeout,
		)
	}
}
