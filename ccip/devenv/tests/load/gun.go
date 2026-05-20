package load

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/require"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// EVMDestination identifies one EVM destination this gun can target.
// Receiver is captured up-front so the gun does not perform per-call lookups.
type EVMDestination struct {
	Chain    cciptestinterfaces.CCIP17
	Receiver protocol.UnknownAddress
}

// Canton2EVMGun performs one Canton→EVM CCIP message per WASP Call against a round-robin
// EVM destination from a configurable list.
//
// Environment-agnostic: the gun does NOT mint tokens or call SetupSend. Callers
// (devenv test runners) must pre-fund holdings before starting the WASP profile;
// staging/prod runners rely on pre-existing funded accounts. All Canton-holding
// work runs inside Call, serialized with RPS=1 (Canton holdings 1-wide).
type Canton2EVMGun struct {
	mu            sync.Mutex
	inFlight      int32
	maxConcurrent int32
	calls         atomic.Int64

	lib         ccv.Lib
	cantonChain cciptestinterfaces.CCIP17

	destinations []EVMDestination
	destCursor   atomic.Int64

	ccvAddr      protocol.UnknownAddress
	executorAddr protocol.UnknownAddress

	confirmSendTimeout time.Duration
	confirmExecTimeout time.Duration
}

// NewCanton2EVMGun wires references resolved once by the test. destinations must be non-empty;
// the gun selects one per Call using round-robin. ccvAddr and executorAddr are source-side
// (Canton) and shared across all destinations.
func NewCanton2EVMGun(
	lib ccv.Lib,
	cantonChain cciptestinterfaces.CCIP17,
	destinations []EVMDestination,
	ccvAddr, executorAddr protocol.UnknownAddress,
	confirmExecTimeout time.Duration,
) (*Canton2EVMGun, error) {
	if len(destinations) == 0 {
		return nil, fmt.Errorf("Canton2EVMGun: at least one EVM destination is required")
	}
	for i, d := range destinations {
		if d.Chain == nil {
			return nil, fmt.Errorf("Canton2EVMGun: destination[%d].Chain is nil", i)
		}
		if len(d.Receiver) == 0 {
			return nil, fmt.Errorf("Canton2EVMGun: destination[%d].Receiver is empty", i)
		}
	}
	if confirmExecTimeout <= 0 {
		confirmExecTimeout = 5 * time.Minute
	}

	return &Canton2EVMGun{
		lib:                lib,
		cantonChain:        cantonChain,
		destinations:       destinations,
		ccvAddr:            ccvAddr,
		executorAddr:       executorAddr,
		confirmSendTimeout: 30 * time.Second,
		confirmExecTimeout: confirmExecTimeout,
	}, nil
}

// nextDestination picks the next EVM destination round-robin.
func (g *Canton2EVMGun) nextDestination() EVMDestination {
	n := len(g.destinations)
	i := g.destCursor.Add(1) - 1

	return g.destinations[int(i%int64(n))]
}

// Call implements wasp.Gun: send a Canton→EVM message, confirm send, assert verifier, confirm exec on the picked EVM dest.
func (g *Canton2EVMGun) Call(gen *wasp.Generator) *wasp.Response {
	start := time.Now()
	if gen == nil || gen.Cfg == nil || gen.Cfg.T == nil {
		return &wasp.Response{Failed: true, Error: "generator or testing.T missing", Duration: time.Since(start)}
	}
	t := gen.Cfg.T

	var depth int32
	g.mu.Lock()
	g.inFlight++
	depth = g.inFlight
	if depth > g.maxConcurrent {
		g.maxConcurrent = depth
	}
	g.mu.Unlock()

	g.calls.Add(1)
	defer func() {
		g.mu.Lock()
		g.inFlight--
		g.mu.Unlock()
	}()

	if depth > 1 {
		require.FailNow(t, "overlapping Canton2EVMGun.Call",
			"expected single-flight (Canton holdings 1-wide); concurrent depth=%d", depth)
	}

	dest := g.nextDestination()
	subtestCtx := t.Context()
	n := g.calls.Load()
	data := fmt.Appendf(nil, "canton2evm load n=%d dest=%d", n, dest.Chain.ChainSelector())

	sendRes, err := g.cantonChain.SendMessage(
		subtestCtx,
		dest.Chain.ChainSelector(),
		cciptestinterfaces.MessageFields{
			Receiver: dest.Receiver,
			Data:     data,
		},
		cciptestinterfaces.MessageOptions{
			ExecutionGasLimit: 200_000,
			FinalityConfig:    1,
			Executor:          g.executorAddr,
			CCVs: []protocol.CCV{
				{CCVAddress: g.ccvAddr, Args: []byte{}, ArgsLen: 0},
			},
		},
		3,
	)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("SendMessage (dest=%d): %v", dest.Chain.ChainSelector(), err), Duration: time.Since(start)}
	}
	if sendRes.Message == nil {
		return &wasp.Response{Failed: true, Error: "SendMessage returned nil message", Duration: time.Since(start)}
	}
	seqNo := uint64(sendRes.Message.SequenceNumber)

	sentEvent, err := g.cantonChain.ConfirmSendOnSource(
		subtestCtx,
		dest.Chain.ChainSelector(),
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
		g.confirmSendTimeout,
	)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("ConfirmSendOnSource (dest=%d): %v", dest.Chain.ChainSelector(), err), Duration: time.Since(start)}
	}

	devenvtests.AssertSingleVerifierResult(t, subtestCtx, g.lib, sentEvent.MessageID)

	ev, err := dest.Chain.ConfirmExecOnDest(
		subtestCtx,
		g.cantonChain.ChainSelector(),
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
		g.confirmExecTimeout,
	)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("ConfirmExecOnDest (dest=%d): %v", dest.Chain.ChainSelector(), err), Duration: time.Since(start)}
	}
	if ev.State != cciptestinterfaces.ExecutionStateSuccess {
		return &wasp.Response{
			Failed:     true,
			Error:      fmt.Sprintf("execution state=%s (dest=%d)", ev.State.String(), dest.Chain.ChainSelector()),
			Duration:   time.Since(start),
			StatusCode: "500",
		}
	}

	return &wasp.Response{
		Failed:     false,
		StatusCode: "200",
		Duration:   time.Since(start),
	}
}

// MaxConcurrentObserved returns the peak number of overlapping Call invocations seen.
func (g *Canton2EVMGun) MaxConcurrentObserved() int32 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.maxConcurrent
}

// CallCount returns how many times Call ran.
func (g *Canton2EVMGun) CallCount() int64 {
	return g.calls.Load()
}
