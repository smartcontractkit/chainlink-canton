package load

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type loadMessageBuilder func(
	source cciptestinterfaces.CCIP17,
	callNum int64,
	ccvAddr, executorAddr protocol.UnknownAddress,
) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error)

// Destination identifies one load target chain, its receiver, and how to build load messages for it.
// buildMessage is set by discoverEVMDestinations or discoverCantonDest.
type Destination struct {
	Chain        cciptestinterfaces.CCIP17
	Receiver     protocol.UnknownAddress
	buildMessage loadMessageBuilder
}

// BuildMessage returns CCIP message fields and options for this destination.
func (d Destination) BuildMessage(
	source cciptestinterfaces.CCIP17,
	callNum int64,
	ccvAddr, executorAddr protocol.UnknownAddress,
) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error) {
	if d.buildMessage == nil {
		return cciptestinterfaces.MessageFields{}, cciptestinterfaces.MessageOptions{},
			fmt.Errorf("destination %d: buildMessage not configured", d.Chain.ChainSelector())
	}

	return d.buildMessage(source, callNum, ccvAddr, executorAddr)
}

// CCIPLoadGun performs one CCIP message per WASP Call: send on source, confirm send,
// confirm exec on a round-robin destination.
//
// Environment-agnostic: the gun does NOT mint tokens or call SetupSend. Callers
// (devenv test runners) must pre-fund holdings before starting the WASP profile when
// required; staging/prod runners rely on pre-existing funded accounts.
type CCIPLoadGun struct {
	mu            sync.Mutex
	inFlight      int32
	maxConcurrent int32
	calls         atomic.Int64

	source       cciptestinterfaces.CCIP17
	destinations []Destination
	destCursor   atomic.Int64

	ccvAddr      protocol.UnknownAddress
	executorAddr protocol.UnknownAddress

	confirmSendTimeout time.Duration
	confirmExecTimeout time.Duration
}

// NewCCIPLoadGun wires a CCIP source with one or more destinations for load testing.
func NewCCIPLoadGun(
	source cciptestinterfaces.CCIP17,
	destinations []Destination,
	ccvAddr, executorAddr protocol.UnknownAddress,
	confirmExecTimeout time.Duration,
) (*CCIPLoadGun, error) {
	if source == nil {
		return nil, fmt.Errorf("CCIPLoadGun: source is nil")
	}
	if len(destinations) == 0 {
		return nil, fmt.Errorf("CCIPLoadGun: at least one destination is required")
	}
	for i, d := range destinations {
		if d.Chain == nil {
			return nil, fmt.Errorf("CCIPLoadGun: destination[%d].Chain is nil", i)
		}
		if len(d.Receiver) == 0 {
			return nil, fmt.Errorf("CCIPLoadGun: destination[%d].Receiver is empty", i)
		}
		if d.buildMessage == nil {
			return nil, fmt.Errorf("CCIPLoadGun: destination[%d].buildMessage is nil", i)
		}
	}
	if confirmExecTimeout <= 0 {
		confirmExecTimeout = 5 * time.Minute
	}

	return &CCIPLoadGun{
		source:             source,
		destinations:       destinations,
		ccvAddr:            ccvAddr,
		executorAddr:       executorAddr,
		confirmSendTimeout: 30 * time.Second,
		confirmExecTimeout: confirmExecTimeout,
	}, nil
}

func (g *CCIPLoadGun) nextDestination() Destination {
	n := len(g.destinations)
	i := g.destCursor.Add(1) - 1

	return g.destinations[int(i%int64(n))]
}

// Call implements wasp.Gun: send, confirm send on source, confirm exec on the picked destination.
func (g *CCIPLoadGun) Call(gen *wasp.Generator) *wasp.Response {
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
		require.FailNow(t, "overlapping CCIPLoadGun.Call",
			"expected single-flight; concurrent depth=%d", depth)
	}

	dest := g.nextDestination()
	destSelector := dest.Chain.ChainSelector()
	subtestCtx := t.Context()
	callNum := g.calls.Load()
	fields, opts, err := dest.BuildMessage(g.source, callNum, g.ccvAddr, g.executorAddr)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("BuildMessage (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}

	sendRes, err := g.source.SendMessage(
		subtestCtx,
		destSelector,
		fields,
		opts,
		3,
	)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("SendMessage (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}
	if sendRes.Message == nil {
		return &wasp.Response{Failed: true, Error: "SendMessage returned nil message", Duration: time.Since(start)}
	}
	seqNo := uint64(sendRes.Message.SequenceNumber)

	sentEvent, err := g.source.ConfirmSendOnSource(
		subtestCtx,
		destSelector,
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
		g.confirmSendTimeout,
	)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("ConfirmSendOnSource (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}

	ev, err := dest.Chain.ConfirmExecOnDest(
		subtestCtx,
		g.source.ChainSelector(),
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID},
		g.confirmExecTimeout,
	)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("ConfirmExecOnDest (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}
	if ev.State != cciptestinterfaces.ExecutionStateSuccess {
		return &wasp.Response{
			Failed:     true,
			Error:      fmt.Sprintf("execution state=%s (dest=%d)", ev.State.String(), destSelector),
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
func (g *CCIPLoadGun) MaxConcurrentObserved() int32 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.maxConcurrent
}

// CallCount returns how many times Call ran.
func (g *CCIPLoadGun) CallCount() int64 {
	return g.calls.Load()
}
