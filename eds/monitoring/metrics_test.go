package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// collectDAMetrics records the given requests through the real labeler and returns what an OTEL
// exporter would see. It swaps the global beholder client, so tests using it must not run in parallel.
func collectDAMetrics(t *testing.T, record func(l *EDSMetricLabeler)) map[string]metricdata.Metrics {
	t.Helper()

	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(MetricViews()...),
	)

	previous := beholder.GetClient()
	client := beholder.NewNoopClient()
	client.Meter = provider.Meter("test")
	beholder.SetClient(client)
	t.Cleanup(func() { beholder.SetClient(previous) })

	edsMetrics, err := InitMetrics()
	require.NoError(t, err)
	record(NewEDSMetricLabeler(metrics.NewLabeler(), edsMetrics))

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	out := map[string]metricdata.Metrics{}
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			out[m.Name] = m
		}
	}

	return out
}

func attrs(kv ...string) attribute.Set {
	var set []attribute.KeyValue
	for i := 0; i < len(kv); i += 2 {
		set = append(set, attribute.String(kv[i], kv[i+1]))
	}

	return attribute.NewSet(set...)
}

//nolint:paralleltest // swaps the global beholder client, so it cannot run in parallel
func TestRecordDARequestExportedMetrics(t *testing.T) {
	got := collectDAMetrics(t, func(l *EDSMetricLabeler) {
		ctx := context.Background()
		l.RecordDARequest(ctx, "burn_mint_factory", true, 120*time.Millisecond)
		l.RecordDARequest(ctx, "burn_mint_factory", true, 300*time.Millisecond)
		l.RecordDARequest(ctx, "burn_mint_factory", false, 3*time.Second)
		l.RecordDARequest(ctx, "bridge_holdings", true, 80*time.Millisecond)
	})

	t.Run("requests counter splits by operation and result", func(t *testing.T) {
		t.Parallel()
		m, ok := got["eds_da_requests_total"]
		require.True(t, ok, "eds_da_requests_total was not exported")
		sum, ok := m.Data.(metricdata.Sum[int64])
		require.True(t, ok)

		counts := map[attribute.Set]int64{}
		for _, dp := range sum.DataPoints {
			counts[dp.Attributes] = dp.Value
		}
		assert.Equal(t, map[attribute.Set]int64{
			attrs("operation", "burn_mint_factory", "result", "success"): 2,
			attrs("operation", "burn_mint_factory", "result", "failure"): 1,
			attrs("operation", "bridge_holdings", "result", "success"):   1,
		}, counts)
	})

	t.Run("up gauge holds the latest result per operation", func(t *testing.T) {
		t.Parallel()
		m, ok := got["eds_da_up"]
		require.True(t, ok, "eds_da_up was not exported")
		gauge, ok := m.Data.(metricdata.Gauge[int64])
		require.True(t, ok)

		values := map[attribute.Set]int64{}
		for _, dp := range gauge.DataPoints {
			values[dp.Attributes] = dp.Value
		}
		assert.Equal(t, map[attribute.Set]int64{
			attrs("operation", "burn_mint_factory"): 0, // last call failed
			attrs("operation", "bridge_holdings"):   1,
		}, values)
	})

	t.Run("duration histogram is in seconds with the configured buckets", func(t *testing.T) {
		t.Parallel()
		m, ok := got["eds_da_request_duration_seconds"]
		require.True(t, ok, "eds_da_request_duration_seconds was not exported")
		assert.Equal(t, "seconds", m.Unit)
		hist, ok := m.Data.(metricdata.Histogram[float64])
		require.True(t, ok)

		for _, dp := range hist.DataPoints {
			assert.Equal(t, []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20}, dp.Bounds)
			if dp.Attributes == attrs("operation", "burn_mint_factory") {
				assert.EqualValues(t, 3, dp.Count)
				assert.InDelta(t, 3.42, dp.Sum, 0.001)
			}
		}
		assert.Len(t, hist.DataPoints, 2)
	})

	t.Run("no label carries a URL or path", func(t *testing.T) {
		t.Parallel()
		for _, name := range []string{"eds_da_requests_total", "eds_da_up", "eds_da_request_duration_seconds"} {
			for _, key := range exportedLabelKeys(got[name]) {
				assert.Contains(t, []string{"operation", "result"}, key, "unexpected label on %s", name)
			}
		}
	})
}

func exportedLabelKeys(m metricdata.Metrics) []string {
	var sets []attribute.Set
	switch d := m.Data.(type) {
	case metricdata.Sum[int64]:
		for _, dp := range d.DataPoints {
			sets = append(sets, dp.Attributes)
		}
	case metricdata.Gauge[int64]:
		for _, dp := range d.DataPoints {
			sets = append(sets, dp.Attributes)
		}
	case metricdata.Histogram[float64]:
		for _, dp := range d.DataPoints {
			sets = append(sets, dp.Attributes)
		}
	}

	var keys []string
	for _, set := range sets {
		iter := set.Iter()
		for iter.Next() {
			keys = append(keys, string(iter.Attribute().Key))
		}
	}

	return keys
}

// Nothing recorded means nothing exported, so an idle service has no series to alert on.
//
//nolint:paralleltest // swaps the global beholder client, so it cannot run in parallel
func TestRecordDARequestIdleExportsNoSamples(t *testing.T) {
	got := collectDAMetrics(t, func(*EDSMetricLabeler) {})

	for _, name := range []string{"eds_da_requests_total", "eds_da_up", "eds_da_request_duration_seconds"} {
		_, ok := got[name]
		assert.False(t, ok, "%s must not be exported without traffic", name)
	}
}
