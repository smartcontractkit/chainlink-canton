package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestMonitoringMiddleware(metrics HTTPAPIMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics.IncrementActiveRequestsCounter(c.Request.Context())

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		metrics.DecrementActiveRequestsCounter(c.Request.Context())

		metrics.RecordHTTPRequestDuration(
			c.Request.Context(),
			duration,
			c.Request.URL.Path,
			c.Request.Method,
			c.Writer.Status(),
		)
	}
}
