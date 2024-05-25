// The TestRateLimiter function tests the RateLimiter middleware by simulating multiple requests and
// checking the response status codes.
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"project/api/middleware"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limiter := middleware.RateLimiter(handler)

	// Create a new request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Call the RateLimiter middleware 6 times
	for i := 0; i < 6; i++ {
		// Create a new response recorder for each request
		res := httptest.NewRecorder()

		limiter.ServeHTTP(res, req)

		// If this is the 6th request, check if the response status code is 429
		if i == 5 {
			assert.Equal(t, http.StatusTooManyRequests, res.Code)
		} else {
			// For the first 5 requests, check if the response status code is OK
			assert.Equal(t, http.StatusOK, res.Code)
		}
	}
}
