package common

import (
	"context"
	"time"
)

type EDSMonitoring interface {
	Metrics() EDSMetricLabeler
}

type EDSMetricLabeler interface {
	With(keyValues ...string) EDSMetricLabeler

	// HTTP API
	IncrementActiveRequestsCounter(ctx context.Context)
	DecrementActiveRequestsCounter(ctx context.Context)
	RecordHTTPRequestDuration(ctx context.Context, duration time.Duration, path, method string, status int)

	// Contract Store
	IncrementStoreSubscriptionUptime(ctx context.Context)
	IncrementStoreUpdatesCounter(ctx context.Context)
	RecordStoreLedgerEndGauge(ctx context.Context, ledgerEnd int64)
}
