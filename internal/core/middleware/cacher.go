package middleware

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type Caching struct {
	client       *redis.Client
	ttl          time.Duration
	cacheControl bool
}

func NewCaching(redisAddr string, ttl time.Duration, cacheControl bool) *Caching {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr, // e.g., "localhost:6379"
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return &Caching{
		client:       client,
		ttl:          ttl,
		cacheControl: cacheControl,
	}
}

func (c *Caching) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Respect Cache-Control: no-store
		if c.cacheControl && r.Header.Get("Cache-Control") == "no-store" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet {
			data, err := c.client.Get(ctx, r.RequestURI).Result()
			if err == nil {
				// Cache hit
				w.Header().Set("X-Cache-Hit", "true")
				io.WriteString(w, data)
				return
			} else if err != redis.Nil {
				log.Printf("Error retrieving from cache: %v", err)
			}
		}

		// Capture response
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: new(bytes.Buffer)}
		next.ServeHTTP(rec, r)

		// Cache only successful responses (status code 200)
		if r.Method == http.MethodGet && rec.statusCode == http.StatusOK {
			err := c.client.Set(ctx, r.RequestURI, rec.body.String(), c.ttl).Err()
			if err != nil {
				log.Printf("Error setting cache: %v", err)
			}
		}
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rec *responseRecorder) Write(p []byte) (int, error) {
	rec.body.Write(p)                  // Write to the buffer
	return rec.ResponseWriter.Write(p) // Stream response to client
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Implement the Flush method
func (rec *responseRecorder) Flush() {
	if flusher, ok := rec.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
