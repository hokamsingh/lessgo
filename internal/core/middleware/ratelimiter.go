package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	interval time.Duration
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		interval: interval,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		key := r.RemoteAddr
		if rl.requests[key] >= rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			rl.mu.Unlock()
			return
		}
		rl.requests[key]++
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.interval)
		rl.mu.Lock()
		rl.requests = make(map[string]int)
		rl.mu.Unlock()
	}
}
