package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Data string `json:"data"`
}

// jsonOfSize returns a JSON body of exactly the given size by padding the "data" field.
// The minimum overhead is for `{"data":""}` which is 12 bytes.
func jsonOfSize(t *testing.T, size int) string {
	t.Helper()
	const overhead = len(`{"data":""}`)
	require.GreaterOrEqual(t, size, overhead, "requested size must be >= JSON overhead")
	padding := strings.Repeat("x", size-overhead)
	b, err := json.Marshal(testRequest{Data: padding})
	require.NoError(t, err)
	require.Len(t, b, size)

	return string(b)
}

func TestRequestSizeLimiterMiddleware(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		limitBytes   int64
		bodySize     int
		customBody   string
		expectStatus int
	}{
		{
			name:         "body exactly at limit is accepted",
			limitBytes:   50,
			bodySize:     50,
			expectStatus: http.StatusOK,
		},
		{
			name:         "body below limit is accepted",
			limitBytes:   50,
			bodySize:     20,
			expectStatus: http.StatusOK,
		},
		{
			name:         "body exceeding limit by one byte is rejected",
			limitBytes:   50,
			bodySize:     51,
			expectStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:         "body significantly exceeding limit is rejected",
			limitBytes:   50,
			bodySize:     1000,
			expectStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:         "zero limit rejects any body",
			limitBytes:   0,
			bodySize:     12, // minimum JSON overhead
			expectStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:         "negative limit rejects any body",
			limitBytes:   -1,
			bodySize:     12, // minimum JSON overhead
			expectStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:         "malformed JSON within limit returns bad request",
			limitBytes:   1000,
			customBody:   "{invalid json",
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			router := gin.New()
			router.Use(RequestSizeLimiterMiddleware(tt.limitBytes))
			router.POST("/test", func(c *gin.Context) {
				var req testRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
						c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
						return
					}
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

					return
				}
				c.Status(http.StatusOK)
			})

			body := tt.customBody
			if body == "" {
				body = jsonOfSize(t, tt.bodySize)
			}
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/test", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)
		})
	}
}
