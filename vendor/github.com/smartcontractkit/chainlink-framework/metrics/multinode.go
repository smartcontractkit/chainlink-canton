package metrics

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

var (
	// Node States
	promMultiNodeRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "multi_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})

	// Node Verification
	promNodeClientVersion = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pool_rpc_node_client_version",
		Help: "Tracks node RPC client versions",
	}, []string{"network", "chainID", "nodeName", "rpcClientVersion"})
	promPoolRPCNodeVerifies = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies",
		Help: "The total number of chain ID verifications for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeVerifiesFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies_failed",
		Help: "The total number of failed chain ID verifications for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeVerifiesSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies_success",
		Help: "The total number of successful chain ID verifications for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})

	// Node State Transitions
	promPoolRPCNodeTransitionsToAlive = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_alive",
		Help: "Total number of times node has transitioned to Alive",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeTransitionsToInSync = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_in_sync",
		Help: "Total number of times node has transitioned from OutOfSync to Alive",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeTransitionsToOutOfSync = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_out_of_sync",
		Help: "Total number of times node has transitioned to OutOfSync",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeTransitionsToUnreachable = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_unreachable",
		Help: "Total number of times node has transitioned to Unreachable",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeTransitionsToInvalidChainID = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_invalid_chain_id",
		Help: "Total number of times node has transitioned to InvalidChainID",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeTransitionsToUnusable = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_unusable",
		Help: "Total number of times node has transitioned to Unusable",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeTransitionsToSyncing = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_syncing",
		Help: "Total number of times node has transitioned to Syncing",
	}, []string{"network", "chainID", "nodeName"})

	// Node Lifecycle
	promPoolRPCNodeHighestSeenBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pool_rpc_node_highest_seen_block",
		Help: "The highest seen block for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	PromPoolRPCNodeHighestFinalizedBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pool_rpc_node_highest_finalized_block",
		Help: "The highest seen finalized block for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeNumSeenBlocks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_seen_blocks",
		Help: "The total number of new blocks seen by the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodePolls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_polls_total",
		Help: "The total number of poll checks for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodePollsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_polls_failed",
		Help: "The total number of failed poll checks for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodePollsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_polls_success",
		Help: "The total number of successful poll checks for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})

	// Transaction Sender
	promMultiNodeInvariantViolations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "multi_node_invariant_violations",
		Help: "The number of invariant violations",
	}, []string{"network", "chainId", "invariant"})
)

type GenericMultiNodeMetrics interface {
	RecordNodeStates(ctx context.Context, state string, count int64)
	RecordNodeClientVersion(ctx context.Context, nodeName string, version string)
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
	IncrementInvariantViolations(ctx context.Context, invariant string)
	SetHighestSeenBlock(ctx context.Context, nodeName string, blockNumber int64)
	SetHighestFinalizedBlock(ctx context.Context, nodeName string, blockNumber int64)
	IncrementSeenBlocks(ctx context.Context, nodeName string)
	IncrementPolls(ctx context.Context, nodeName string)
	IncrementPollsFailed(ctx context.Context, nodeName string)
	IncrementPollsSuccess(ctx context.Context, nodeName string)
}

var _ GenericMultiNodeMetrics = &multiNodeMetrics{}

type multiNodeMetrics struct {
	network                         string
	chainID                         string
	nodeStates                      metric.Int64Gauge
	nodeClientVersion               metric.Int64Gauge
	nodeVerifies                    metric.Int64Counter
	nodeVerifiesFailed              metric.Int64Counter
	nodeVerifiesSuccess             metric.Int64Counter
	nodeTransitionsToAlive          metric.Int64Counter
	nodeTransitionsToInSync         metric.Int64Counter
	nodeTransitionsToOutOfSync      metric.Int64Counter
	nodeTransitionsToUnreachable    metric.Int64Counter
	nodeTransitionsToInvalidChainID metric.Int64Counter
	nodeTransitionsToUnusable       metric.Int64Counter
	nodeTransitionsToSyncing        metric.Int64Counter
	highestSeenBlock                metric.Int64Gauge
	highestFinalizedBlock           metric.Int64Gauge
	seenBlocks                      metric.Int64Counter
	polls                           metric.Int64Counter
	pollsFailed                     metric.Int64Counter
	pollsSuccess                    metric.Int64Counter
	invariantViolations             metric.Int64Counter
}

func NewGenericMultiNodeMetrics(network string, chainID string) (GenericMultiNodeMetrics, error) {
	nodeStates, err := beholder.GetMeter().Int64Gauge("multi_node_states")
	if err != nil {
		return nil, fmt.Errorf("failed to register multinode states metric: %w", err)
	}

	nodeClientVersion, err := beholder.GetMeter().Int64Gauge("pool_rpc_node_client_version")
	if err != nil {
		return nil, fmt.Errorf("failed to register client version metric: %w", err)
	}

	nodeVerifies, err := beholder.GetMeter().Int64Counter("pool_rpc_node_verifies")
	if err != nil {
		return nil, fmt.Errorf("failed to register node verifies metric: %w", err)
	}

	nodeVerifiesFailed, err := beholder.GetMeter().Int64Counter("pool_rpc_node_verifies_failed")
	if err != nil {
		return nil, fmt.Errorf("failed to register node verifies failed metric: %w", err)
	}

	nodeVerifiesSuccess, err := beholder.GetMeter().Int64Counter("pool_rpc_node_verifies_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register node verifies success metric: %w", err)
	}

	nodeTransitionsToAlive, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_alive")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to alive metric: %w", err)
	}

	nodeTransitionsToInSync, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_in_sync")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to in sync metric: %w", err)
	}

	nodeTransitionsToOutOfSync, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_out_of_sync")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to out of sync metric: %w", err)
	}

	nodeTransitionsToUnreachable, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_unreachable")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to unreachable metric: %w", err)
	}

	nodeTransitionsToInvalidChainID, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_invalid_chain_id")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to invalid chain id metric: %w", err)
	}

	nodeTransitionsToUnusable, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_unusable")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to unusable metric: %w", err)
	}

	nodeTransitionsToSyncing, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_transitions_to_syncing")
	if err != nil {
		return nil, fmt.Errorf("failed to register node transitions to syncing metric: %w", err)
	}

	highestSeenBlock, err := beholder.GetMeter().Int64Gauge("pool_rpc_node_highest_seen_block")
	if err != nil {
		return nil, fmt.Errorf("failed to register highest seen block metric: %w", err)
	}

	highestFinalizedBlock, err := beholder.GetMeter().Int64Gauge("pool_rpc_node_highest_finalized_block")
	if err != nil {
		return nil, fmt.Errorf("failed to register highest finalized block metric: %w", err)
	}

	seenBlocks, err := beholder.GetMeter().Int64Counter("pool_rpc_node_num_seen_blocks")
	if err != nil {
		return nil, fmt.Errorf("failed to register seen blocks counter: %w", err)
	}

	polls, err := beholder.GetMeter().Int64Counter("pool_rpc_node_polls_total")
	if err != nil {
		return nil, fmt.Errorf("failed to register node polls metric: %w", err)
	}

	pollsFailed, err := beholder.GetMeter().Int64Counter("pool_rpc_node_polls_failed")
	if err != nil {
		return nil, fmt.Errorf("failed to register node polls failed metric: %w", err)
	}

	pollsSuccess, err := beholder.GetMeter().Int64Counter("pool_rpc_node_polls_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register node polls success metric: %w", err)
	}

	invariantViolations, err := beholder.GetMeter().Int64Counter("multi_node_invariant_violations")
	if err != nil {
		return nil, fmt.Errorf("failed to register invariant violations metric: %w", err)
	}

	return &multiNodeMetrics{
		network:                         network,
		chainID:                         chainID,
		nodeStates:                      nodeStates,
		nodeClientVersion:               nodeClientVersion,
		nodeVerifies:                    nodeVerifies,
		nodeVerifiesFailed:              nodeVerifiesFailed,
		nodeVerifiesSuccess:             nodeVerifiesSuccess,
		nodeTransitionsToAlive:          nodeTransitionsToAlive,
		nodeTransitionsToInSync:         nodeTransitionsToInSync,
		nodeTransitionsToOutOfSync:      nodeTransitionsToOutOfSync,
		nodeTransitionsToUnreachable:    nodeTransitionsToUnreachable,
		nodeTransitionsToInvalidChainID: nodeTransitionsToInvalidChainID,
		nodeTransitionsToUnusable:       nodeTransitionsToUnusable,
		nodeTransitionsToSyncing:        nodeTransitionsToSyncing,
		highestSeenBlock:                highestSeenBlock,
		highestFinalizedBlock:           highestFinalizedBlock,
		seenBlocks:                      seenBlocks,
		polls:                           polls,
		pollsFailed:                     pollsFailed,
		pollsSuccess:                    pollsSuccess,
		invariantViolations:             invariantViolations,
	}, nil
}

func (m *multiNodeMetrics) RecordNodeStates(ctx context.Context, state string, count int64) {
	promMultiNodeRPCNodeStates.WithLabelValues(m.network, m.chainID, state).Set(float64(count))
	m.nodeStates.Record(ctx, count, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("state", state)))
}

func (m *multiNodeMetrics) RecordNodeClientVersion(ctx context.Context, nodeName string, version string) {
	promNodeClientVersion.WithLabelValues(m.network, m.chainID, nodeName, version).Set(1)
	m.nodeClientVersion.Record(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName),
		attribute.String("rpcClientVersion", version)))
}

func (m *multiNodeMetrics) IncrementNodeVerifies(ctx context.Context, nodeName string) {
	promPoolRPCNodeVerifies.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeVerifies.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeVerifiesFailed(ctx context.Context, nodeName string) {
	promPoolRPCNodeVerifiesFailed.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeVerifiesFailed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeVerifiesSuccess(ctx context.Context, nodeName string) {
	promPoolRPCNodeVerifiesSuccess.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeVerifiesSuccess.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToAlive(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToAlive.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToAlive.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToInSync(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToInSync.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToInSync.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToOutOfSync(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToOutOfSync.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToOutOfSync.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToUnreachable(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToUnreachable.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToInvalidChainID(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToInvalidChainID.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToUnusable(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToUnusable.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToUnusable.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementNodeTransitionsToSyncing(ctx context.Context, nodeName string) {
	promPoolRPCNodeTransitionsToSyncing.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.nodeTransitionsToSyncing.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) SetHighestSeenBlock(ctx context.Context, nodeName string, blockNumber int64) {
	promPoolRPCNodeHighestSeenBlock.WithLabelValues(m.network, m.chainID, nodeName).Set(float64(blockNumber))
	m.highestSeenBlock.Record(ctx, blockNumber, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) SetHighestFinalizedBlock(ctx context.Context, nodeName string, blockNumber int64) {
	PromPoolRPCNodeHighestFinalizedBlock.WithLabelValues(m.network, m.chainID, nodeName).Set(float64(blockNumber))
	m.highestFinalizedBlock.Record(ctx, blockNumber, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementSeenBlocks(ctx context.Context, nodeName string) {
	promPoolRPCNodeNumSeenBlocks.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.seenBlocks.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementPolls(ctx context.Context, nodeName string) {
	promPoolRPCNodePolls.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.polls.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementPollsFailed(ctx context.Context, nodeName string) {
	promPoolRPCNodePollsFailed.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.pollsFailed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementPollsSuccess(ctx context.Context, nodeName string) {
	promPoolRPCNodePollsSuccess.WithLabelValues(m.network, m.chainID, nodeName).Inc()
	m.pollsSuccess.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("nodeName", nodeName)))
}

func (m *multiNodeMetrics) IncrementInvariantViolations(ctx context.Context, invariant string) {
	promMultiNodeInvariantViolations.WithLabelValues(m.network, m.chainID, invariant).Inc()
	m.invariantViolations.Add(ctx, 1, metric.WithAttributes(
		attribute.String("network", m.network),
		attribute.String("chainID", m.chainID),
		attribute.String("invariant", invariant)))
}
