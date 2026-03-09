package store

import (
	"context"
)

type StoreMetrics interface {
	IncrementStoreSubscriptionUptime(ctx context.Context)
	IncrementStoreUpdatesCounter(ctx context.Context)
	RecordStoreLedgerEndGauge(ctx context.Context, ledgerEnd int64)
}
