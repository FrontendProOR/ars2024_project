package middleware

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
)

var idempotencyMap = sync.Map{}

func IdempotencyCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idempotencyKey := r.Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			// Generate a UUID
			idempotencyKey = uuid.New().String()
			// Add the UUID to the request header
			r.Header.Set("Idempotency-Key", idempotencyKey)
		}

		if _, ok := idempotencyMap.Load(idempotencyKey); ok {
			http.Error(w, "Duplicate request", http.StatusConflict)
			return
		}
		idempotencyMap.Store(idempotencyKey, struct{}{})
		next(w, r)
	}
}
