package multinode

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type RPCClientBaseConfig interface {
	NewHeadsPollInterval() time.Duration
	FinalizedBlockPollInterval() time.Duration
}

// RPCClientBase is used to integrate multinode into chain-specific clients.
// For new MultiNode integrations, we wrap the RPC client and inherit from the RPCClientBase
// to get the required RPCClient methods and enable the use of MultiNode.
//
// The RPCClientBase provides chain-agnostic functionality such as head and finalized head
// subscriptions, which are required in each Node lifecycle to execute various
// health checks.
type RPCClientBase[HEAD Head] struct {
	cfg        RPCClientBaseConfig
	log        logger.Logger
	ctxTimeout time.Duration
	subsMu     sync.RWMutex
	subs       map[Subscription]struct{}

	latestBlock          func(ctx context.Context) (HEAD, error)
	latestFinalizedBlock func(ctx context.Context) (HEAD, error)

	// lifeCycleCh can be closed to immediately cancel all in-flight requests on
	// this RPC. Closing and replacing should be serialized through
	// lifeCycleMu since it can happen on state transitions as well as RPCClientBase Close.
	// Also closed when RPC is declared unhealthy.
	lifeCycleMu sync.RWMutex
	lifeCycleCh chan struct{}

	// chainInfoLock protects highestUserObservations and latestChainInfo
	chainInfoLock sync.RWMutex
	// intercepted values seen by callers of the RPCClientBase excluding health check calls. Need to ensure MultiNode provides repeatable read guarantee
	highestUserObservations ChainInfo
	// most recent chain info observed during current lifecycle
	latestChainInfo ChainInfo
}

func NewRPCClientBase[HEAD Head](
	cfg RPCClientBaseConfig, ctxTimeout time.Duration, log logger.Logger,
	latestBlock func(ctx context.Context) (HEAD, error),
	latestFinalizedBlock func(ctx context.Context) (HEAD, error),
) *RPCClientBase[HEAD] {
	return &RPCClientBase[HEAD]{
		cfg:                  cfg,
		log:                  log,
		ctxTimeout:           ctxTimeout,
		latestBlock:          latestBlock,
		latestFinalizedBlock: latestFinalizedBlock,
		subs:                 make(map[Subscription]struct{}),
		lifeCycleCh:          make(chan struct{}),
	}
}

func (m *RPCClientBase[HEAD]) lenSubs() int {
	m.subsMu.RLock()
	defer m.subsMu.RUnlock()
	return len(m.subs)
}

// RegisterSub adds the sub to the RPCClientBase list and returns a managed sub which is removed on unsubscribe
func (m *RPCClientBase[HEAD]) RegisterSub(sub Subscription, lifeCycleCh chan struct{}) (*ManagedSubscription, error) {
	// ensure that the `sub` belongs to current life cycle of the `RPCClientBase` and it should not be killed due to
	// previous `DisconnectAll` call.
	select {
	case <-lifeCycleCh:
		sub.Unsubscribe()
		return nil, fmt.Errorf("failed to register subscription - all in-flight requests were canceled")
	default:
	}
	m.subsMu.Lock()
	defer m.subsMu.Unlock()
	managedSub := &ManagedSubscription{
		sub,
		m.removeSub,
	}
	m.subs[managedSub] = struct{}{}
	return managedSub, nil
}

func (m *RPCClientBase[HEAD]) removeSub(sub Subscription) {
	m.subsMu.Lock()
	defer m.subsMu.Unlock()
	delete(m.subs, sub)
}

func (m *RPCClientBase[HEAD]) SubscribeToHeads(ctx context.Context) (<-chan HEAD, Subscription, error) {
	ctx, cancel, lifeCycleCh := m.AcquireQueryCtx(ctx, m.ctxTimeout)
	defer cancel()

	pollInterval := m.cfg.NewHeadsPollInterval()
	if pollInterval == 0 {
		return nil, nil, errors.New("PollInterval is 0")
	}
	timeout := pollInterval
	poller, channel := NewPoller[HEAD](pollInterval, func(pollRequestCtx context.Context) (HEAD, error) {
		if CtxIsHealthCheckRequest(ctx) {
			pollRequestCtx = CtxAddHealthCheckFlag(pollRequestCtx)
		}
		return m.LatestBlock(pollRequestCtx)
	}, timeout, m.log)

	if err := poller.Start(ctx); err != nil {
		return nil, nil, err
	}

	sub, err := m.RegisterSub(&poller, lifeCycleCh)
	if err != nil {
		sub.Unsubscribe()
		return nil, nil, err
	}
	return channel, sub, nil
}

func (m *RPCClientBase[HEAD]) SubscribeToFinalizedHeads(ctx context.Context) (<-chan HEAD, Subscription, error) {
	ctx, cancel, lifeCycleCh := m.AcquireQueryCtx(ctx, m.ctxTimeout)
	defer cancel()

	finalizedBlockPollInterval := m.cfg.FinalizedBlockPollInterval()
	if finalizedBlockPollInterval == 0 {
		return nil, nil, errors.New("FinalizedBlockPollInterval is 0")
	}
	timeout := finalizedBlockPollInterval
	poller, channel := NewPoller[HEAD](finalizedBlockPollInterval, func(pollRequestCtx context.Context) (HEAD, error) {
		if CtxIsHealthCheckRequest(ctx) {
			pollRequestCtx = CtxAddHealthCheckFlag(pollRequestCtx)
		}
		return m.LatestFinalizedBlock(pollRequestCtx)
	}, timeout, m.log)
	if err := poller.Start(ctx); err != nil {
		return nil, nil, err
	}

	sub, err := m.RegisterSub(&poller, lifeCycleCh)
	if err != nil {
		poller.Unsubscribe()
		return nil, nil, err
	}
	return channel, sub, nil
}

func (m *RPCClientBase[HEAD]) LatestBlock(ctx context.Context) (HEAD, error) {
	// capture lifeCycleCh to ensure we are not updating chainInfo with observations related to previous life cycle
	ctx, cancel, lifeCycleCh := m.AcquireQueryCtx(ctx, m.ctxTimeout)
	defer cancel()

	head, err := m.latestBlock(ctx)
	if err != nil {
		return head, err
	}

	if !head.IsValid() {
		return head, errors.New("invalid head")
	}

	m.OnNewHead(ctx, lifeCycleCh, head)
	return head, nil
}

func (m *RPCClientBase[HEAD]) LatestFinalizedBlock(ctx context.Context) (HEAD, error) {
	ctx, cancel, lifeCycleCh := m.AcquireQueryCtx(ctx, m.ctxTimeout)
	defer cancel()

	head, err := m.latestFinalizedBlock(ctx)
	if err != nil {
		return head, err
	}

	if !head.IsValid() {
		return head, errors.New("invalid head")
	}

	m.OnNewFinalizedHead(ctx, lifeCycleCh, head)
	return head, nil
}

func (m *RPCClientBase[HEAD]) OnNewHead(ctx context.Context, requestCh <-chan struct{}, head HEAD) {
	if !head.IsValid() {
		return
	}

	m.chainInfoLock.Lock()
	defer m.chainInfoLock.Unlock()
	blockNumber := head.BlockNumber()
	totalDifficulty := head.GetTotalDifficulty()
	if !CtxIsHealthCheckRequest(ctx) {
		m.highestUserObservations.BlockNumber = max(m.highestUserObservations.BlockNumber, blockNumber)
		m.highestUserObservations.TotalDifficulty = MaxTotalDifficulty(m.highestUserObservations.TotalDifficulty, totalDifficulty)
	}
	select {
	case <-requestCh: // no need to update latestChainInfo, as RPCClientBase already started new life cycle
		return
	default:
		m.latestChainInfo.BlockNumber = blockNumber
		m.latestChainInfo.TotalDifficulty = totalDifficulty
	}
}

func (m *RPCClientBase[HEAD]) OnNewFinalizedHead(ctx context.Context, requestCh <-chan struct{}, head HEAD) {
	if !head.IsValid() {
		return
	}

	m.chainInfoLock.Lock()
	defer m.chainInfoLock.Unlock()
	if !CtxIsHealthCheckRequest(ctx) {
		m.highestUserObservations.FinalizedBlockNumber = max(m.highestUserObservations.FinalizedBlockNumber, head.BlockNumber())
	}
	select {
	case <-requestCh: // no need to update latestChainInfo, as RPCClientBase already started new life cycle
		return
	default:
		m.latestChainInfo.FinalizedBlockNumber = head.BlockNumber()
	}
}

// makeQueryCtx returns a context that cancels if:
// 1. Passed in ctx cancels
// 2. Passed in channel is closed
// 3. Default timeout is reached (queryTimeout)
func makeQueryCtx(ctx context.Context, ch services.StopChan, timeout time.Duration) (context.Context, context.CancelFunc) {
	var chCancel, timeoutCancel context.CancelFunc
	ctx, chCancel = ch.Ctx(ctx)
	ctx, timeoutCancel = context.WithTimeout(ctx, timeout)
	cancel := func() {
		chCancel()
		timeoutCancel()
	}
	return ctx, cancel
}

func (m *RPCClientBase[HEAD]) AcquireQueryCtx(parentCtx context.Context, timeout time.Duration) (ctx context.Context, cancel context.CancelFunc,
	lifeCycleCh chan struct{}) {
	// Need to wrap in mutex because state transition can cancel and replace context
	m.lifeCycleMu.RLock()
	lifeCycleCh = m.lifeCycleCh
	m.lifeCycleMu.RUnlock()
	ctx, cancel = makeQueryCtx(parentCtx, lifeCycleCh, timeout)
	return
}

func (m *RPCClientBase[HEAD]) UnsubscribeAllExcept(subs ...Subscription) {
	m.subsMu.Lock()
	keepSubs := map[Subscription]struct{}{}
	for _, sub := range subs {
		keepSubs[sub] = struct{}{}
	}

	var unsubs []Subscription
	for sub := range m.subs {
		if _, keep := keepSubs[sub]; !keep {
			unsubs = append(unsubs, sub)
		}
	}
	m.subsMu.Unlock()

	for _, sub := range unsubs {
		sub.Unsubscribe()
	}
}

// CancelLifeCycle closes and replaces the lifeCycleCh
func (m *RPCClientBase[HEAD]) CancelLifeCycle() {
	m.lifeCycleMu.Lock()
	defer m.lifeCycleMu.Unlock()
	close(m.lifeCycleCh)
	m.lifeCycleCh = make(chan struct{})
}

func (m *RPCClientBase[HEAD]) resetLatestChainInfo() {
	m.chainInfoLock.Lock()
	m.latestChainInfo = ChainInfo{}
	m.chainInfoLock.Unlock()
}

func (m *RPCClientBase[HEAD]) Close() {
	m.CancelLifeCycle()
	m.UnsubscribeAllExcept()
	m.resetLatestChainInfo()
}

func (m *RPCClientBase[HEAD]) GetInterceptedChainInfo() (latest, highestUserObservations ChainInfo) {
	m.chainInfoLock.RLock()
	defer m.chainInfoLock.RUnlock()
	return m.latestChainInfo, m.highestUserObservations
}
