package load

import (
	"fmt"
	"os"
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
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/require"

	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
)

const (
	envMessageRate      = "CANTON_LOAD_MESSAGE_RATE"
	envLoadDuration     = "CANTON_LOAD_DURATION"
	defaultMessageRate  = "1/1s"
	defaultLoadDuration = 90 * time.Second
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

func runWASP(t *testing.T, gun *CCIPLoadGun, genName string, sched scheduleConfig) {
	t.Helper()

	p := wasp.NewProfile().Add(wasp.NewGenerator(&wasp.Config{
		T:        t,
		LoadType: wasp.RPS,
		GenName:  genName,
		Schedule: wasp.Combine(
			wasp.Plain(sched.rate, sched.duration),
		),
		RateLimitUnitDuration: sched.rateUnit,
		Gun:                   gun,
		Labels: map[string]string{
			"go_test_name":  genName,
			"branch":        "test",
			"commit":        "test",
			"message_rate":  sched.messageRate,
			"load_duration": sched.duration.String(),
		},
		LokiConfig: nil,
	}))

	_, err := p.Run(true)
	require.NoError(t, err)
	p.Wait()

	require.Positive(t, gun.CallCount(), "gun should have completed at least one message")
	require.LessOrEqual(t, gun.MaxConcurrentObserved(), int32(1),
		"Gun.Call must not overlap (single-flight)")
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
