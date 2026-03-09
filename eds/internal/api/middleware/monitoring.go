package middleware

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-canton/eds/common"
)

var _ HTTPAPIMetrics = common.EDSMetricLabeler(nil)

// HTTPAPIMetrics captures a subset of metrics that are used for monitoring the HTTP API.
// It is implemented by the global common.EDSMetricLabeler
type HTTPAPIMetrics interface {
	IncrementActiveRequestsCounter(ctx context.Context)
	DecrementActiveRequestsCounter(ctx context.Context)
	RecordHTTPRequestDuration(ctx context.Context, duration time.Duration, path, method string, status int)
}
