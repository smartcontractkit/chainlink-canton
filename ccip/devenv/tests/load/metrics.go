package load

import (
	"slices"
	"sync"
	"testing"
	"time"

	ccvmetrics "github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/metrics"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

// LoadMessageRecord captures per-message phase timings for a successful load test call.
type LoadMessageRecord struct {
	SeqNo               uint64
	SourceChain         uint64
	DestChain           uint64
	MessageID           protocol.Bytes32
	SentTime            time.Time
	ExecutedTime        time.Time
	SendDuration        time.Duration
	ConfirmSendDuration time.Duration
	ConfirmExecDuration time.Duration
	TotalDuration       time.Duration
}

// LoadFailureCounts tracks failures by load test phase.
type LoadFailureCounts struct {
	Send        int
	ConfirmSend int
	ConfirmExec int
}

// LoadMetricsCollector accumulates successful message records and phase failure counts.
type LoadMetricsCollector struct {
	mu       sync.Mutex
	records  []LoadMessageRecord
	failures LoadFailureCounts
}

func (c *LoadMetricsCollector) appendRecord(r LoadMessageRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records = append(c.records, r)
}

func (c *LoadMetricsCollector) incrementSendFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures.Send++
}

func (c *LoadMetricsCollector) incrementConfirmSendFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures.ConfirmSend++
}

func (c *LoadMetricsCollector) incrementConfirmExecFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures.ConfirmExec++
}

func (c *LoadMetricsCollector) snapshot() ([]LoadMessageRecord, LoadFailureCounts) {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]LoadMessageRecord, len(c.records))
	copy(out, c.records)
	return out, c.failures
}

// PercentileStats holds common percentile values for a set of durations.
type PercentileStats struct {
	Min time.Duration
	Max time.Duration
	P90 time.Duration
	P95 time.Duration
	P99 time.Duration
}

func calculatePercentiles(durations []time.Duration) PercentileStats {
	if len(durations) == 0 {
		return PercentileStats{}
	}

	slices.Sort(durations)

	p90Index := int(float64(len(durations)) * 0.90)
	p95Index := int(float64(len(durations)) * 0.95)
	p99Index := int(float64(len(durations)) * 0.99)

	if p90Index >= len(durations) {
		p90Index = len(durations) - 1
	}
	if p95Index >= len(durations) {
		p95Index = len(durations) - 1
	}
	if p99Index >= len(durations) {
		p99Index = len(durations) - 1
	}

	return PercentileStats{
		Min: durations[0],
		Max: durations[len(durations)-1],
		P90: durations[p90Index],
		P95: durations[p95Index],
		P99: durations[p99Index],
	}
}

func logPhasePercentiles(t *testing.T, label string, stats PercentileStats, count int) {
	t.Helper()
	t.Logf("%s (n=%d):\n"+
		"  Min: %v\n"+
		"  Max: %v\n"+
		"  P90: %v\n"+
		"  P95: %v\n"+
		"  P99: %v",
		label, count,
		stats.Min, stats.Max, stats.P90, stats.P95, stats.P99,
	)
}

// PrintPhaseMetricsSummary prints per-phase timing percentiles for successful load messages.
func PrintPhaseMetricsSummary(t *testing.T, records []LoadMessageRecord, failures LoadFailureCounts, skipExecConfirm bool) {
	t.Helper()

	t.Logf("\n" +
		"========================================\n" +
		"         Load Phase Metrics            \n" +
		"========================================")

	if len(records) == 0 {
		t.Logf("No successful load messages recorded")
	} else {
		sendDurations := make([]time.Duration, len(records))
		confirmSendDurations := make([]time.Duration, len(records))
		for i, r := range records {
			sendDurations[i] = r.SendDuration
			confirmSendDurations[i] = r.ConfirmSendDuration
		}

		logPhasePercentiles(t, "Send", calculatePercentiles(sendDurations), len(records))
		logPhasePercentiles(t, "Confirm Send", calculatePercentiles(confirmSendDurations), len(records))

		if skipExecConfirm {
			sourceConfirmed := make([]time.Duration, len(records))
			for i, r := range records {
				sourceConfirmed[i] = r.SendDuration + r.ConfirmSendDuration
			}
			logPhasePercentiles(t, "Source Confirmed (Send → Confirm Send)", calculatePercentiles(sourceConfirmed), len(records))
		} else {
			confirmExecDurations := make([]time.Duration, len(records))
			totalDurations := make([]time.Duration, len(records))
			for i, r := range records {
				confirmExecDurations[i] = r.ConfirmExecDuration
				totalDurations[i] = r.TotalDuration
			}
			logPhasePercentiles(t, "Confirm Exec", calculatePercentiles(confirmExecDurations), len(records))
			logPhasePercentiles(t, "Total (E2E)", calculatePercentiles(totalDurations), len(records))
		}
	}

	totalFailures := failures.Send + failures.ConfirmSend + failures.ConfirmExec
	if totalFailures > 0 {
		t.Logf("----------------------------------------\n"+
			"Failures: send=%d confirm_send=%d confirm_exec=%d",
			failures.Send, failures.ConfirmSend, failures.ConfirmExec)
	}

	t.Logf("========================================")
}

// ToCCVMessageMetrics maps load records with ExecutedTime set to CCV message metrics.
func ToCCVMessageMetrics(records []LoadMessageRecord) []ccvmetrics.MessageMetrics {
	out := make([]ccvmetrics.MessageMetrics, 0, len(records))
	for _, r := range records {
		if r.ExecutedTime.IsZero() {
			continue
		}
		out = append(out, ccvmetrics.MessageMetrics{
			SeqNo:           r.SeqNo,
			MessageID:       r.MessageID.String(),
			SourceChain:     r.SourceChain,
			DestChain:       r.DestChain,
			SentTime:        r.SentTime,
			ExecutedTime:    r.ExecutedTime,
			LatencyDuration: r.ExecutedTime.Sub(r.SentTime),
		})
	}
	return out
}
