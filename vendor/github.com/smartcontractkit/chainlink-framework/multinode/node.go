package multinode

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const QueryTimeout = 10 * time.Second

var errInvalidChainID = errors.New("invalid chain id")

type NodeConfig interface {
	PollFailureThreshold() uint32
	PollInterval() time.Duration
	SelectionMode() string
	SyncThreshold() uint32
	NodeIsSyncingEnabled() bool
	FinalizedBlockPollInterval() time.Duration
	EnforceRepeatableRead() bool
	DeathDeclarationDelay() time.Duration
	NewHeadsPollInterval() time.Duration
	VerifyChainID() bool
}

type ChainConfig interface {
	NodeNoNewHeadsThreshold() time.Duration
	NoNewFinalizedHeadsThreshold() time.Duration
	FinalityDepth() uint32
	FinalityTagEnabled() bool
	FinalizedBlockOffset() uint32
}

type nodeMetrics interface {
	IncrementNodeVerifies(ctx context.Context, nodeName string)
	IncrementNodeVerifiesFailed(ctx context.Context, nodeName string)
	IncrementNodeVerifiesSuccess(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToAlive(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToInSync(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToOutOfSync(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToUnreachable(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToInvalidChainID(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToUnusable(ctx context.Context, nodeName string)
	IncrementNodeTransitionsToSyncing(ctx context.Context, nodeName string)
	RecordNodeClientVersion(ctx context.Context, nodeName string, version string)
	SetHighestSeenBlock(ctx context.Context, nodeName string, blockNumber int64)
	SetHighestFinalizedBlock(ctx context.Context, nodeName string, blockNumber int64)
	IncrementSeenBlocks(ctx context.Context, nodeName string)
	IncrementPolls(ctx context.Context, nodeName string)
	IncrementPollsFailed(ctx context.Context, nodeName string)
	IncrementPollsSuccess(ctx context.Context, nodeName string)
}

type Node[
	CHAIN_ID ID,
	RPC any,
] interface {
	// State returns most accurate state of the Node on the moment of call.
	// While some of the checks may be performed in the background and State may return cached value, critical, like
	// `FinalizedBlockOutOfSync`, must be executed upon every call.
	State() nodeState
	// StateAndLatest returns nodeState with the latest ChainInfo observed by Node during current lifecycle.
	StateAndLatest() (nodeState, ChainInfo)
	// HighestUserObservations - returns highest ChainInfo ever observed by underlying RPC excluding results of health check requests
	HighestUserObservations() ChainInfo
	SetPoolChainInfoProvider(PoolChainInfoProvider)
	// Name is a unique identifier for this node.
	Name() string
	// String - returns string representation of the node, useful for debugging (name + URLS used to connect to the RPC)
	String() string
	RPC() RPC
	// UnsubscribeAllExceptAliveLoop - closes all subscriptions except the aliveLoop subscription
	UnsubscribeAllExceptAliveLoop()
	ConfiguredChainID() CHAIN_ID
	// Order - returns priority order configured for the RPC
	Order() int32
	// Start - starts health checks
	Start(context.Context) error
	Close() error
}

type node[
	CHAIN_ID ID,
	HEAD Head,
	RPC RPCClient[CHAIN_ID, HEAD],
] struct {
	services.StateMachine
	lfcLog      logger.Logger
	name        string
	id          int
	chainID     CHAIN_ID
	nodePoolCfg NodeConfig
	chainCfg    ChainConfig
	order       int32
	chainFamily string

	metrics nodeMetrics

	ws   *url.URL
	http *url.URL

	rpc               RPC
	isLoadBalancedRPC bool

	stateMu sync.RWMutex // protects state* fields
	state   nodeState

	poolInfoProvider PoolChainInfoProvider

	stopCh services.StopChan
	// wg waits for subsidiary goroutines
	wg sync.WaitGroup

	healthCheckSubs []Subscription
}

func NewNode[
	CHAIN_ID ID,
	HEAD Head,
	RPC RPCClient[CHAIN_ID, HEAD],
](
	nodeCfg NodeConfig,
	chainCfg ChainConfig,
	lggr logger.Logger,
	metrics nodeMetrics,
	wsuri *url.URL,
	httpuri *url.URL,
	name string,
	id int,
	chainID CHAIN_ID,
	nodeOrder int32,
	rpc RPC,
	chainFamily string,
	isLoadBalancedRPC bool,
) Node[CHAIN_ID, RPC] {
	n := new(node[CHAIN_ID, HEAD, RPC])
	n.name = name
	n.id = id
	n.chainID = chainID
	n.nodePoolCfg = nodeCfg
	n.chainCfg = chainCfg
	n.order = nodeOrder
	n.metrics = metrics
	if wsuri != nil {
		n.ws = wsuri
	}
	if httpuri != nil {
		n.http = httpuri
	}
	n.stopCh = make(services.StopChan)
	lggr = logger.Named(lggr, "Node")
	lggr = logger.With(lggr,
		"nodeTier", Primary.String(),
		"nodeName", name,
		"node", n.String(),
		"chainID", chainID,
		"nodeOrder", n.order,
	)
	n.lfcLog = logger.Named(lggr, "Lifecycle")
	n.rpc = rpc
	n.isLoadBalancedRPC = isLoadBalancedRPC
	n.chainFamily = chainFamily
	return n
}

func (n *node[CHAIN_ID, HEAD, RPC]) String() string {
	s := fmt.Sprintf("(%s)%s", Primary.String(), n.name)
	if n.ws != nil {
		s += fmt.Sprintf(":%s", n.ws.String())
	}
	if n.http != nil {
		s += fmt.Sprintf(":%s", n.http.String())
	}
	return s
}

func (n *node[CHAIN_ID, HEAD, RPC]) ConfiguredChainID() (chainID CHAIN_ID) {
	return n.chainID
}

func (n *node[CHAIN_ID, HEAD, RPC]) Name() string {
	return n.name
}

func (n *node[CHAIN_ID, HEAD, RPC]) RPC() RPC {
	return n.rpc
}

// unsubscribeAllExceptAliveLoop is not thread-safe; it should only be called
// while holding the stateMu lock.
func (n *node[CHAIN_ID, HEAD, RPC]) unsubscribeAllExceptAliveLoop() {
	n.rpc.UnsubscribeAllExcept(n.healthCheckSubs...)
}

func (n *node[CHAIN_ID, HEAD, RPC]) UnsubscribeAllExceptAliveLoop() {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	n.unsubscribeAllExceptAliveLoop()
}

func (n *node[CHAIN_ID, HEAD, RPC]) Close() error {
	return n.StopOnce(n.name, n.close)
}

func (n *node[CHAIN_ID, HEAD, RPC]) close() error {
	defer func() {
		n.wg.Wait()
		n.rpc.Close()
	}()

	n.stateMu.Lock()
	defer n.stateMu.Unlock()

	close(n.stopCh)
	n.state = nodeStateClosed
	return nil
}

// Start dials and verifies the node
// Should only be called once in a node's lifecycle
// Return value is necessary to conform to interface but this will never
// actually return an error.
func (n *node[CHAIN_ID, HEAD, RPC]) Start(startCtx context.Context) error {
	return n.StartOnce(n.name, func() error {
		n.wg.Add(1)
		go n.start()
		return nil
	})
}

// start initially dials the node and verifies chain ID
// This spins off lifecycle goroutines.
// Not thread-safe.
// Node lifecycle is synchronous: only one goroutine should be running at a
// time.
func (n *node[CHAIN_ID, HEAD, RPC]) start() {
	defer n.wg.Done()
	if n.State() != nodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", n.state))
	}
	ctx, cancel := n.stopCh.NewCtx()
	defer cancel()

	if err := n.rpc.Dial(ctx); err != nil {
		n.lfcLog.Errorw("Dial failed: Node is unreachable", "err", err)
		n.declareUnreachable()
		return
	}
	n.setState(nodeStateDialed)

	state := n.verifyConn(ctx, n.lfcLog)
	n.declareState(state)
}

// verifyChainID checks that connection to the node matches the given chain ID
// Not thread-safe
// Pure verifyChainID: does not mutate node "state" field.
func (n *node[CHAIN_ID, HEAD, RPC]) verifyChainID(callerCtx context.Context, lggr logger.Logger) nodeState {
	n.metrics.IncrementNodeVerifies(callerCtx, n.name)
	promFailed := func() {
		n.metrics.IncrementNodeVerifiesFailed(callerCtx, n.name)
	}

	st := n.getCachedState()
	switch st {
	case nodeStateClosed:
		// The node is already closed, and any subsequent transition is invalid.
		// To make spotting such transitions a bit easier, return the invalid node state.
		return nodeStateLen
	case nodeStateDialed, nodeStateOutOfSync, nodeStateInvalidChainID, nodeStateSyncing:
	default:
		panic(fmt.Sprintf("cannot verify node in state %v", st))
	}

	var chainID CHAIN_ID
	var err error
	if chainID, err = n.rpc.ChainID(callerCtx); err != nil {
		promFailed()
		lggr.Errorw("Failed to verify chain ID for node", "err", err, "nodeState", n.getCachedState())
		return nodeStateUnreachable
	} else if chainID.String() != n.chainID.String() {
		promFailed()
		err = fmt.Errorf(
			"rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s: %w",
			chainID.String(),
			n.chainID.String(),
			n.name,
			errInvalidChainID,
		)
		lggr.Errorw("Failed to verify RPC node; remote endpoint returned the wrong chain ID", "err", err, "nodeState", n.getCachedState())
		return nodeStateInvalidChainID
	}

	n.metrics.IncrementNodeVerifiesSuccess(callerCtx, n.name)

	return nodeStateAlive
}

// createVerifiedConn - establishes new connection with the RPC and verifies that it's valid: chainID matches, and it's not syncing.
// Returns desired state if one of the verifications fails. Otherwise, returns nodeStateAlive.
func (n *node[CHAIN_ID, HEAD, RPC]) createVerifiedConn(ctx context.Context, lggr logger.Logger) nodeState {
	if err := n.rpc.Dial(ctx); err != nil {
		n.lfcLog.Errorw("Dial failed: Node is unreachable", "err", err, "nodeState", n.getCachedState())
		return nodeStateUnreachable
	}

	return n.verifyConn(ctx, lggr)
}

// verifyConn - verifies that current connection is valid: chainID matches, and it's not syncing.
// Returns desired state if one of the verifications fails. Otherwise, returns nodeStateAlive.
func (n *node[CHAIN_ID, HEAD, RPC]) verifyConn(ctx context.Context, lggr logger.Logger) nodeState {
	if n.nodePoolCfg.VerifyChainID() {
		state := n.verifyChainID(ctx, lggr)
		if state != nodeStateAlive {
			return state
		}
	}

	if n.nodePoolCfg.NodeIsSyncingEnabled() {
		isSyncing, err := n.rpc.IsSyncing(ctx)
		if err != nil {
			lggr.Errorw("Unexpected error while verifying RPC node synchronization status", "err", err, "nodeState", n.getCachedState())
			return nodeStateUnreachable
		}

		if isSyncing {
			lggr.Errorw("Verification failed: Node is syncing", "nodeState", n.getCachedState())
			return nodeStateSyncing
		}
	}

	return nodeStateAlive
}

func (n *node[CHAIN_ID, HEAD, RPC]) Order() int32 {
	return n.order
}

func (n *node[CHAIN_ID, HEAD, RPC]) newCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := n.stopCh.NewCtx()
	ctx = CtxAddHealthCheckFlag(ctx)
	return ctx, cancel
}
