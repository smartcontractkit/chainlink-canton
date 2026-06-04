package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestSizeLimiterMiddleware returns a middleware that limits all incoming request sized to be below (or equal to)
// the configured limitBytes.
// Negative limits are treated as `0`, see http.MaxBytesReader for more details.
// Must be used before any middleware that reads the request body, otherwise the limit will not strictly be enforced.
func RequestSizeLimiterMiddleware(limitBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitBytes)
		c.Next()
	}
}
