package store

import (
	"context"

	"github.com/smartcontractkit/chainlink-canton/eds/common"
)

var _ Metrics = common.EDSMetricLabeler(nil)

// Metrics captures a subset of metrics that are used for monitoring a ContractStore.
// // It is implemented by the global common.EDSMetricLabeler
type Metrics interface {
	IncrementStoreSubscriptionUptime(ctx context.Context)
	IncrementStoreUpdatesCounter(ctx context.Context)
	RecordStoreLedgerEndGauge(ctx context.Context, ledgerEnd int64)
}
