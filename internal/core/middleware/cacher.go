package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type Caching struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCaching(redisAddr string, ttl time.Duration) *Caching {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr, // e.g., "localhost:6379"
	})
	return &Caching{
		client: client,
		ttl:    ttl,
	}
}

func (c *Caching) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Try to get the cached response from Redis
		data, err := c.client.Get(ctx, r.RequestURI).Result()
		if err == nil {
			// If found in cache, write it directly to the response
			w.Write([]byte(data))
			return
		}

		// If not cached, create a response writer to capture the response
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		// Cache the response in Redis
		c.client.Set(ctx, r.RequestURI, rec.body, c.ttl)
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
