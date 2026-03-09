package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/smartcontractkit/chainlink-canton/eds/common"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type EDSMetrics struct {
	activeRequestsUpDownCounter metric.Int64UpDownCounter
	requestDurationSeconds      metric.Float64Histogram

	storeSubscriptionUptime metric.Int64Counter
	storeUpdatesCounter     metric.Int64Counter
	storeLedgerEndGauge     metric.Int64Gauge
}

func InitMetrics() (*EDSMetrics, error) {
	m := &EDSMetrics{}
	var err error

	m.activeRequestsUpDownCounter, err = beholder.GetMeter().Int64UpDownCounter(
		"eds_active_requests",
		metric.WithDescription("Total number of active requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds_active_requests: %w", err)
	}

	m.requestDurationSeconds, err = beholder.GetMeter().Float64Histogram(
		"eds_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("seconds"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds_request_duration_seconds: %w", err)
	}

	m.storeSubscriptionUptime, err = beholder.GetMeter().Int64Counter(
		"eds_store_ledger_api_subscription_uptime",
		metric.WithDescription("Counter increasing while EDS is successfully subscribed to Ledger API updates"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds_store_ledger_api_subscription_uptime: %w", err)
	}

	m.storeUpdatesCounter, err = beholder.GetMeter().Int64Counter(
		"eds_store_updates_total",
		metric.WithDescription("Total number of CreateEvent updates processed by EDS"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds_store_updates_total: %w", err)
	}

	m.storeLedgerEndGauge, err = beholder.GetMeter().Int64Gauge(
		"eds_store_ledger_end",
		metric.WithDescription("The current ledger end that the EDS has processed. Should not be used for liveness checks as it will only increase if there are contract updates to be processed."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds_store_ledger_end: %w", err)
	}

	return m, nil
}

func MetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "eds_request_duration_seconds"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			}},
		),
	}
}

var _ common.EDSMetricLabeler = &EDSMetricLabeler{}

type EDSMetricLabeler struct {
	metrics.Labeler
	m *EDSMetrics
}

func NewEDSMetricLabeler(labeler metrics.Labeler, m *EDSMetrics) *EDSMetricLabeler {
	return &EDSMetricLabeler{
		Labeler: labeler,
		m:       m,
	}
}

func (l *EDSMetricLabeler) With(keyValues ...string) common.EDSMetricLabeler {
	return &EDSMetricLabeler{Labeler: l.Labeler.With(keyValues...), m: l.m}
}

func (l *EDSMetricLabeler) IncrementActiveRequestsCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(l.Labels).AsStringAttributes()
	l.m.activeRequestsUpDownCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l *EDSMetricLabeler) DecrementActiveRequestsCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(l.Labels).AsStringAttributes()
	l.m.activeRequestsUpDownCounter.Add(ctx, -1, metric.WithAttributes(otelLabels...))
}

func (l *EDSMetricLabeler) RecordHTTPRequestDuration(ctx context.Context, duration time.Duration, path, method string, status int) {
	otelLabels := beholder.OtelAttributes(l.Labels).AsStringAttributes()
	l.m.requestDurationSeconds.Record(ctx, duration.Seconds(), metric.WithAttributes([]attribute.KeyValue{
		attribute.String("path", path),
		attribute.String("method", method),
		attribute.Int("status", status),
	}...), metric.WithAttributes(otelLabels...))
}

func (l *EDSMetricLabeler) IncrementStoreSubscriptionUptime(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(l.Labels).AsStringAttributes()
	l.m.storeSubscriptionUptime.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l *EDSMetricLabeler) IncrementStoreUpdatesCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(l.Labels).AsStringAttributes()
	l.m.storeUpdatesCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l *EDSMetricLabeler) RecordStoreLedgerEndGauge(ctx context.Context, ledgerEnd int64) {
	otelLabels := beholder.OtelAttributes(l.Labels).AsStringAttributes()
	l.m.storeLedgerEndGauge.Record(ctx, ledgerEnd, metric.WithAttributes(otelLabels...))
}
