// middleware/middleware.go
package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware returns a mux.MiddlewareFunc that enforces rate limiting.
func RateLimitMiddleware(limiter *rate.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded, please try again later.(5 requests per second)", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
