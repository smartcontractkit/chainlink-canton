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
	NodeBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{Name: "node_balance", Help: "Account balances"},
		[]string{"account", "chainID", "chainFamily"},
	)
)

type GenericBalanceMetrics interface {
	RecordNodeBalance(ctx context.Context, account string, balance float64)
}

var _ GenericBalanceMetrics = &balanceMetrics{}

type balanceMetrics struct {
	network     string
	chainID     string
	nodeBalance metric.Float64Gauge
}

func NewGenericBalanceMetrics(network string, chainID string) (GenericBalanceMetrics, error) {
	nodeBalance, err := beholder.GetMeter().Float64Gauge("node_balance")
	if err != nil {
		return nil, fmt.Errorf("failed to register node balance metric: %w", err)
	}

	return &balanceMetrics{
		network:     network,
		chainID:     chainID,
		nodeBalance: nodeBalance,
	}, nil
}

func (m *balanceMetrics) RecordNodeBalance(ctx context.Context, account string, balance float64) {
	NodeBalance.WithLabelValues(account, m.chainID, m.network).Set(balance)
	m.nodeBalance.Record(ctx, balance, metric.WithAttributes(
		attribute.String("chainFamily", m.network), // leaving chainFamily for backward compatibility with prom metric
		attribute.String("chainID", m.chainID),
		attribute.String("account", account)))
}
