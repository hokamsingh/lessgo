package middleware

import (
	"log"
	"net/http"
	"time"
)

// ProfilingMiddleware represents a structure for profiling requests
type ProfilingMiddleware struct {
}

// NewProfilingMiddleware creates a new instance of ProfilingMiddleware
func NewProfilingMiddleware() *ProfilingMiddleware {
	return &ProfilingMiddleware{}
}

// Handle processes the requests and measures their execution time
func (pm *ProfilingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Record the start time of the request

		// Pass the request to the next handler in the middleware chain
		next.ServeHTTP(w, r)

		// After the request is completed, measure the execution time
		duration := time.Since(start)

		// Log information about the request and its execution time
		log.Printf("Request: %s %s | Time: %v | From IP: %s", r.Method, r.URL.Path, duration, getClientIP(r))
	})
}

// getClientIP extracts the client's IP address from the request
func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
