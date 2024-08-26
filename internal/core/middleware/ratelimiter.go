package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests        map[string][]time.Time
	mu              sync.Mutex
	limit           int
	interval        time.Duration
	cleanupInterval time.Duration
}

func NewRateLimiter(limit int, interval, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:        make(map[string][]time.Time),
		limit:           limit,
		interval:        interval,
		cleanupInterval: cleanupInterval,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		key := r.RemoteAddr
		now := time.Now()

		// Filter out expired timestamps
		requests := rl.requests[key]
		var newRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.interval {
				newRequests = append(newRequests, reqTime)
			}
		}
		rl.requests[key] = newRequests

		if len(newRequests) >= rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			rl.mu.Unlock()
			return
		}

		// Add current request timestamp
		rl.requests[key] = append(rl.requests[key], now)
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.cleanupInterval)
		rl.mu.Lock()
		now := time.Now()
		for key, timestamps := range rl.requests {
			var validTimestamps []time.Time
			for _, reqTime := range timestamps {
				if now.Sub(reqTime) < rl.interval {
					validTimestamps = append(validTimestamps, reqTime)
				}
			}
			if len(validTimestamps) > 0 {
				rl.requests[key] = validTimestamps
			} else {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}
