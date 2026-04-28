package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func getPathLabel(c *gin.Context, unmatchedLabel string) string {
	if fullPath := c.FullPath(); fullPath != "" {
		return fullPath
	}

	return unmatchedLabel
}

func RequestMonitoringMiddleware(metrics HTTPAPIMetrics) gin.HandlerFunc {
	const unmatchedPathLabel = "unmatched"
	return func(c *gin.Context) {
		metrics.IncrementActiveRequestsCounter(c.Request.Context())

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		metrics.DecrementActiveRequestsCounter(c.Request.Context())

		metrics.RecordHTTPRequestDuration(
			c.Request.Context(),
			duration,
			getPathLabel(c, unmatchedPathLabel),
			c.Request.Method,
			c.Writer.Status(),
		)
	}
}
