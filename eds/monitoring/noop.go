package monitoring

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-canton/eds/common"
)

var _ common.EDSMetricLabeler = &NoopEDSMetricLabeler{}

type NoopEDSMetricLabeler struct{}

func (n NoopEDSMetricLabeler) With(keyValues ...string) common.EDSMetricLabeler { return n }

func (n NoopEDSMetricLabeler) IncrementActiveRequestsCounter(ctx context.Context) {}

func (n NoopEDSMetricLabeler) DecrementActiveRequestsCounter(ctx context.Context) {}

func (n NoopEDSMetricLabeler) RecordHTTPRequestDuration(ctx context.Context, duration time.Duration, path, method string, status int) {
}
func (n NoopEDSMetricLabeler) IncrementStoreSubscriptionUptime(ctx context.Context)           {}
func (n NoopEDSMetricLabeler) IncrementStoreUpdatesCounter(ctx context.Context)               {}
func (n NoopEDSMetricLabeler) RecordStoreLedgerEndGauge(ctx context.Context, ledgerEnd int64) {}
