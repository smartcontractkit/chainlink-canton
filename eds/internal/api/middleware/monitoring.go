package middleware

import (
	"context"
	"time"
)

type HTTPAPIMetrics interface {
	IncrementActiveRequestsCounter(ctx context.Context)
	DecrementActiveRequestsCounter(ctx context.Context)
	RecordHTTPRequestDuration(ctx context.Context, duration time.Duration, path, method string, status int)
}
