package factory

import (
	"context"
	"time"
)

// Fixed operation names for the DA request metrics. They deliberately identify the
// route, not the URL: the URLs carry party IDs and query parameters.
const (
	OperationBurnMintFactory = "burn_mint_factory"
	OperationTransferFactory = "transfer_factory"
	OperationBridgeHoldings  = "bridge_holdings"
)

// Metrics records the outcome of outgoing requests to DA endpoints.
// It is implemented by the global common.EDSMetricLabeler.
type Metrics interface {
	RecordDARequest(ctx context.Context, operation string, success bool, duration time.Duration)
}

// Option configures the factories created by NewBurnMintFactory and NewTransferFactory.
type Option func(*options)

type options struct {
	metrics Metrics
}

// WithMetrics makes the factories record their DA requests. Without it nothing is recorded.
func WithMetrics(metrics Metrics) Option {
	return func(o *options) { o.metrics = metrics }
}

func applyOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	return o
}

// recordDARequest records one DA request, if metrics are configured.
func recordDARequest(ctx context.Context, m Metrics, operation string, started time.Time, success bool) {
	if m == nil {
		return
	}
	m.RecordDARequest(ctx, operation, success, time.Since(started))
}
