package middleware

import (
	"context"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	Timeout time.Duration
}

// NewTimeoutMiddleware creates a new instance of timeout middleware
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{Timeout: timeout}
}

// Handle adds a timeout to the request context
func (tm *TimeoutMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), tm.Timeout)
		defer cancel()

		// Replace the request context with a new context with a timeout
		r = r.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		// End the request when the timeout is reached or after the main handler completes
		select {
		case <-done:
			// The handler completed its work before the timeout
		case <-ctx.Done():
			// Timeout: cancel the request and return 504 status
			http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		}
	})
}
