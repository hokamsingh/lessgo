package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type Caching struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCaching(redisAddr string, ttl time.Duration) *Caching {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr, // e.g., "localhost:6379"
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return &Caching{
		client: client,
		ttl:    ttl,
	}
}

func (c *Caching) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method == http.MethodGet {
			// Try to get the cached response from Redis
			data, err := c.client.Get(ctx, r.RequestURI).Result()
			if err == nil {
				// If found in cache, write it directly to the response
				w.Header().Set("X-Cache-Hit", "true")
				w.Write([]byte(data))
				return
			} else if err != redis.Nil {
				// Log any errors retrieving from Redis
				log.Printf("Error retrieving from cache: %v", err)
			}
		}

		// Create a response writer to capture the response
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		if r.Method == http.MethodGet {
			// Cache the response in Redis
			err := c.client.Set(ctx, r.RequestURI, rec.body, c.ttl).Err()
			if err != nil {
				log.Printf("Error setting cache: %v", err)
			}
		}
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rec *responseRecorder) Write(p []byte) (int, error) {
	rec.body = append(rec.body, p...)
	return rec.ResponseWriter.Write(p)
}

// Implement the Flush method
func (rec *responseRecorder) Flush() {
	if flusher, ok := rec.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
