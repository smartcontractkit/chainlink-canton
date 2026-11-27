package common

import (
	"context"
	"time"
)

// EDSMetricLabeler provides all metric recording functionality for the Explicit Disclosure Server.
type EDSMetricLabeler interface {
	// With returns a new metrics labeler with the given key-value pairs.
	With(keyValues ...string) EDSMetricLabeler

	// HTTP API

	// IncrementActiveRequestsCounter increments the active requests counter.
	IncrementActiveRequestsCounter(ctx context.Context)
	// DecrementActiveRequestsCounter decrements the active requests counter.
	DecrementActiveRequestsCounter(ctx context.Context)
	// RecordHTTPRequestDuration records the HTTP request duration.
	RecordHTTPRequestDuration(ctx context.Context, duration time.Duration, path, method string, status int)

	// Contract Store

	// IncrementStoreSubscriptionUptime increments the store uptime counter.
	IncrementStoreSubscriptionUptime(ctx context.Context)
	// IncrementStoreUpdatesCounter increments the store updates counter.
	IncrementStoreUpdatesCounter(ctx context.Context)
	// IncrementStoreReconnectCounter increments the store reconnects counter.
	IncrementStoreReconnectCounter(ctx context.Context)
	// RecordStoreLedgerEndGauge records the store's current ledger end.
	RecordStoreLedgerEndGauge(ctx context.Context, ledgerEnd int64)

	// Digital Asset (DA) dependencies

	// RecordDARequest records the outcome and duration of one outgoing request to a DA endpoint.
	// operation must be a fixed route name (never a URL), to keep label cardinality bounded.
	RecordDARequest(ctx context.Context, operation string, success bool, duration time.Duration)
}
