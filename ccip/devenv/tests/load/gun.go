package load

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// ConfirmSendFunc confirms a CCIP send on the source chain after SendMessage returns.
type ConfirmSendFunc func(
	t *testing.T,
	ctx context.Context,
	destSelector uint64,
	seqNo uint64,
	sendResult cciptestinterfaces.MessageSentEvent,
) (cciptestinterfaces.MessageSentEvent, error)

// LoadGunOptions configures send confirmation and exec timeout for CCIPLoadGun.
type LoadGunOptions struct {
	ConfirmSend        ConfirmSendFunc
	ConfirmExecTimeout time.Duration
	SkipExecConfirm    bool
}

type loadMessageBuilder func(
	source cciptestinterfaces.CCIP17,
	callNum int64,
	ccvAddr, executorAddr protocol.UnknownAddress,
) (cciptestinterfaces.MessageFields, cciptestinterfaces.MessageOptions, error)

// Destination identifies one load target chain, its receiver, and how to build load messages for it.
// buildMessage is set by discoverEVMDestinationsFromBoot or discoverCantonDest.
type Destination struct {
	Chain        cciptestinterfaces.CCIP17
	Receiver     protocol.UnknownAddress
	TokenLane    *devenvtests.TokenLane // nil = message-only; non-nil = token test
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
	flightReady   sync.Cond
	inFlight      int32
	maxConcurrent int32
	calls         atomic.Int64
	messageIDs    []protocol.Bytes32

	source       cciptestinterfaces.CCIP17
	destinations []Destination
	destCursor   atomic.Int64

	ccvAddr      protocol.UnknownAddress
	executorAddr protocol.UnknownAddress

	confirmSend        ConfirmSendFunc
	confirmExecTimeout time.Duration
	skipExecConfirm    bool

	metricsCollector *LoadMetricsCollector
}

// NewCCIPLoadGun wires a CCIP source with one or more destinations for load testing.
func NewCCIPLoadGun(
	source cciptestinterfaces.CCIP17,
	destinations []Destination,
	ccvAddr, executorAddr protocol.UnknownAddress,
	opts LoadGunOptions,
) (*CCIPLoadGun, error) {
	if source == nil {
		return nil, fmt.Errorf("CCIPLoadGun: source is nil")
	}
	if opts.ConfirmSend == nil {
		return nil, fmt.Errorf("CCIPLoadGun: ConfirmSend is nil")
	}
	if len(destinations) == 0 {
		return nil, fmt.Errorf("CCIPLoadGun: at least one destination is required")
	}
	tokenMode := destinations[0].TokenLane != nil
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
		if (d.TokenLane != nil) != tokenMode {
			return nil, fmt.Errorf("CCIPLoadGun: destination[%d] mixes token and message-only destinations", i)
		}
	}
	confirmExecTimeout := opts.ConfirmExecTimeout
	if confirmExecTimeout <= 0 {
		confirmExecTimeout = 5 * time.Minute
	}

	g := &CCIPLoadGun{
		source:             source,
		destinations:       destinations,
		ccvAddr:            ccvAddr,
		executorAddr:       executorAddr,
		confirmSend:        opts.ConfirmSend,
		confirmExecTimeout: confirmExecTimeout,
		skipExecConfirm:    opts.SkipExecConfirm,
		metricsCollector:   &LoadMetricsCollector{},
	}
	g.flightReady.L = &g.mu

	return g, nil
}

func (g *CCIPLoadGun) acquireSingleFlight() {
	g.mu.Lock()
	for g.inFlight >= 1 {
		g.flightReady.Wait()
	}
	g.inFlight++
	if g.inFlight > g.maxConcurrent {
		g.maxConcurrent = g.inFlight
	}
	g.mu.Unlock()
}

func (g *CCIPLoadGun) releaseSingleFlight() {
	g.mu.Lock()
	g.inFlight--
	if g.inFlight == 0 {
		g.flightReady.Broadcast()
	}
	g.mu.Unlock()
}

// ConfirmExecTimeout returns the exec confirmation timeout configured for this gun.
func (g *CCIPLoadGun) ConfirmExecTimeout() time.Duration {
	return g.confirmExecTimeout
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

	g.acquireSingleFlight()
	defer g.releaseSingleFlight()

	g.calls.Add(1)

	dest := g.nextDestination()
	destSelector := dest.Chain.ChainSelector()
	subtestCtx := t.Context()
	callNum := g.calls.Load()
	fields, opts, err := dest.BuildMessage(g.source, callNum, g.ccvAddr, g.executorAddr)
	if err != nil {
		g.metricsCollector.incrementSendFailure()
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("BuildMessage (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}

	sentTime := time.Now()
	sendRes, err := g.source.SendMessage(
		subtestCtx,
		destSelector,
		fields,
		opts,
		3,
	)
	sendDuration := time.Since(sentTime)
	if err != nil {
		g.metricsCollector.incrementSendFailure()
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("SendMessage (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}
	if sendRes.Message == nil {
		g.metricsCollector.incrementSendFailure()
		return &wasp.Response{Failed: true, Error: "SendMessage returned nil message", Duration: time.Since(start)}
	}
	seqNo := uint64(sendRes.Message.SequenceNumber)

	confirmSendStart := time.Now()
	sentEvent, err := g.confirmSend(t, subtestCtx, destSelector, seqNo, sendRes)
	confirmSendDuration := time.Since(confirmSendStart)
	if err != nil {
		g.metricsCollector.incrementConfirmSendFailure()
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("ConfirmSend (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}

	g.mu.Lock()
	g.messageIDs = append(g.messageIDs, sentEvent.MessageID)
	g.mu.Unlock()

	ccv.Plog.Info().
		Str("messageID", sentEvent.MessageID.String()).
		Uint64("seqNo", seqNo).
		Uint64("destSelector", destSelector).
		Msg("Load message confirmed on source")

	sourceChain := g.source.ChainSelector()
	record := LoadMessageRecord{
		SeqNo:               seqNo,
		SourceChain:         sourceChain,
		DestChain:           destSelector,
		MessageID:           sentEvent.MessageID,
		SentTime:            sentTime,
		SendDuration:        sendDuration,
		ConfirmSendDuration: confirmSendDuration,
	}

	if g.skipExecConfirm {
		record.TotalDuration = time.Since(sentTime)
		g.metricsCollector.appendRecord(record)
		return &wasp.Response{
			Failed:     false,
			StatusCode: "200",
			Duration:   time.Since(start),
		}
	}

	confirmExecStart := time.Now()
	ev, err := dest.Chain.ConfirmExecOnDest(
		subtestCtx,
		g.source.ChainSelector(),
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo, MessageID: sentEvent.MessageID},
		g.confirmExecTimeout,
	)
	confirmExecDuration := time.Since(confirmExecStart)
	if err != nil {
		g.metricsCollector.incrementConfirmExecFailure()
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("ConfirmExecOnDest (dest=%d): %v", destSelector, err), Duration: time.Since(start)}
	}
	if ev.State != cciptestinterfaces.ExecutionStateSuccess {
		g.metricsCollector.incrementConfirmExecFailure()
		return &wasp.Response{
			Failed:     true,
			Error:      fmt.Sprintf("execution state=%s (dest=%d)", ev.State.String(), destSelector),
			Duration:   time.Since(start),
			StatusCode: "500",
		}
	}

	record.ConfirmExecDuration = confirmExecDuration
	record.ExecutedTime = time.Now()
	record.TotalDuration = record.ExecutedTime.Sub(sentTime)
	g.metricsCollector.appendRecord(record)

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

// Metrics returns a copy of successful load message timing records.
func (g *CCIPLoadGun) Metrics() []LoadMessageRecord {
	records, _ := g.metricsCollector.snapshot()
	return records
}

// FailureCounts returns phase failure counters from failed WASP calls.
func (g *CCIPLoadGun) FailureCounts() LoadFailureCounts {
	_, failures := g.metricsCollector.snapshot()
	return failures
}

// MessageIDs returns a copy of message IDs collected after successful ConfirmSend.
func (g *CCIPLoadGun) MessageIDs() []protocol.Bytes32 {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]protocol.Bytes32, len(g.messageIDs))
	copy(out, g.messageIDs)

	return out
}
