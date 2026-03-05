package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

var (
	promNumBroadcasted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_broadcasted",
		Help: "The number of transactions broadcasted",
	}, []string{"chainID"})
	promTimeUntilBroadcast = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tx_manager_time_until_tx_broadcast",
		Help: "The amount of time elapsed from when a transaction is enqueued to until it is broadcast.",
		Buckets: []float64{
			float64(500 * time.Millisecond),
			float64(time.Second),
			float64(5 * time.Second),
			float64(15 * time.Second),
			float64(30 * time.Second),
			float64(time.Minute),
			float64(2 * time.Minute),
		},
	}, []string{"chainID"})
	promNumGasBumps = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	}, []string{"chainID"})

	promGasBumpExceedsLimit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	}, []string{"chainID"})
	promNumConfirmedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_confirmed_transactions",
		Help: "Total number of confirmed transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"chainID"})
	promTimeUntilTxConfirmed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tx_manager_time_until_tx_confirmed",
		Help: "The amount of time elapsed from a transaction being broadcast to being included in a block.",
		Buckets: []float64{
			float64(500 * time.Millisecond),
			float64(time.Second),
			float64(5 * time.Second),
			float64(15 * time.Second),
			float64(30 * time.Second),
			float64(time.Minute),
			float64(2 * time.Minute),
			float64(5 * time.Minute),
			float64(10 * time.Minute),
		},
	}, []string{"chainID"})
	promBlocksUntilTxConfirmed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tx_manager_blocks_until_tx_confirmed",
		Help: "The amount of blocks that have been mined from a transaction being broadcast to being included in a block.",
		Buckets: []float64{
			float64(1),
			float64(5),
			float64(10),
			float64(20),
			float64(50),
			float64(100),
		},
	}, []string{"chainID"})
)

type GenericTXMMetrics interface {
	IncrementNumBroadcastedTxs(ctx context.Context)
	RecordTimeUntilTxBroadcast(ctx context.Context, duration float64)
	IncrementNumGasBumps(ctx context.Context)
	IncrementGasBumpExceedsLimit(ctx context.Context)
	IncrementNumConfirmedTxs(ctx context.Context, confirmedTransactions int)
	RecordTimeUntilTxConfirmed(ctx context.Context, duration float64)
	RecordBlocksUntilTxConfirmed(ctx context.Context, blocksElapsed float64)
}

type txmMetrics struct {
	chainID                string
	numBroadcastedTxs      metric.Int64Counter
	timeUntilBroadcast     metric.Float64Histogram
	numGasBumps            metric.Int64Counter
	gasBumpExceedsLimit    metric.Int64Counter
	numConfirmedTxs        metric.Int64Counter
	timeUntilTxConfirmed   metric.Float64Histogram
	blocksUntilTxConfirmed metric.Float64Histogram
}

func NewGenericTxmMetrics(chainID string) (GenericTXMMetrics, error) {
	numBroadcastedTxs, err := beholder.GetMeter().Int64Counter("tx_manager_num_broadcasted")
	if err != nil {
		return nil, fmt.Errorf("failed to register broadcasted txs number metric: %w", err)
	}

	timeUntilBroadcast, err := beholder.GetMeter().Float64Histogram("tx_manager_time_until_tx_broadcast")
	if err != nil {
		return nil, fmt.Errorf("failed to register time until broadcast metric: %w", err)
	}

	numGasBumps, err := beholder.GetMeter().Int64Counter("tx_manager_num_gas_bumps")
	if err != nil {
		return nil, fmt.Errorf("failed to register number of gas bumps metric: %w", err)
	}

	gasBumpExceedsLimit, err := beholder.GetMeter().Int64Counter("tx_manager_gas_bump_exceeds_limit")
	if err != nil {
		return nil, fmt.Errorf("failed to register gas bump exceeds limit metric: %w", err)
	}

	numConfirmedTxs, err := beholder.GetMeter().Int64Counter("tx_manager_num_confirmed_transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to register confirmed txs number metric: %w", err)
	}

	timeUntilTxConfirmed, err := beholder.GetMeter().Float64Histogram("tx_manager_time_until_tx_confirmed")
	if err != nil {
		return nil, fmt.Errorf("failed to register time until tx confirmed metric: %w", err)
	}

	blocksUntilTxConfirmed, err := beholder.GetMeter().Float64Histogram("tx_manager_blocks_until_tx_confirmed")
	if err != nil {
		return nil, fmt.Errorf("failed to register blocks until tx confirmed metric: %w", err)
	}

	return &txmMetrics{
		chainID:                chainID,
		numBroadcastedTxs:      numBroadcastedTxs,
		timeUntilBroadcast:     timeUntilBroadcast,
		numGasBumps:            numGasBumps,
		gasBumpExceedsLimit:    gasBumpExceedsLimit,
		numConfirmedTxs:        numConfirmedTxs,
		timeUntilTxConfirmed:   timeUntilTxConfirmed,
		blocksUntilTxConfirmed: blocksUntilTxConfirmed,
	}, nil
}

func (m *txmMetrics) IncrementNumBroadcastedTxs(ctx context.Context) {
	promNumBroadcasted.WithLabelValues(m.chainID).Add(float64(1))
	m.numBroadcastedTxs.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *txmMetrics) RecordTimeUntilTxBroadcast(ctx context.Context, duration float64) {
	promTimeUntilBroadcast.WithLabelValues(m.chainID).Observe(duration)
	m.timeUntilBroadcast.Record(ctx, duration, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *txmMetrics) IncrementNumGasBumps(ctx context.Context) {
	promNumGasBumps.WithLabelValues(m.chainID).Add(float64(1))
	m.numGasBumps.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *txmMetrics) IncrementGasBumpExceedsLimit(ctx context.Context) {
	promGasBumpExceedsLimit.WithLabelValues(m.chainID).Add(float64(1))
	m.gasBumpExceedsLimit.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *txmMetrics) IncrementNumConfirmedTxs(ctx context.Context, confirmedTransactions int) {
	promNumConfirmedTxs.WithLabelValues(m.chainID).Add(float64(confirmedTransactions))
	m.numConfirmedTxs.Add(ctx, int64(confirmedTransactions), metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *txmMetrics) RecordTimeUntilTxConfirmed(ctx context.Context, duration float64) {
	promTimeUntilTxConfirmed.WithLabelValues(m.chainID).Observe(duration)
	m.timeUntilTxConfirmed.Record(ctx, duration, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *txmMetrics) RecordBlocksUntilTxConfirmed(ctx context.Context, blocksElapsed float64) {
	promBlocksUntilTxConfirmed.WithLabelValues(m.chainID).Observe(blocksElapsed)
	m.blocksUntilTxConfirmed.Record(ctx, blocksElapsed, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}
