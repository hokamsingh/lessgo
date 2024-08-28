package middleware

import (
	"bytes"
	"compress/gzip"
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

	// Attempt to ping Redis
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Could not connect to Redis: %v. Caching will be disabled.", err)
		client = nil // Mark Redis as unavailable
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

		// Proceed without caching if Redis is unavailable
		if c.client != nil && r.Method == http.MethodGet {
			data, err := c.client.Get(ctx, r.RequestURI).Result()
			if err == nil {
				// Cache hit: decompress data
				reader, err := gzip.NewReader(bytes.NewReader([]byte(data)))
				if err != nil {
					log.Printf("Error decompressing cache data: %v. Proceeding without cache.", err)
					next.ServeHTTP(w, r)
					return
				}
				defer reader.Close()

				// Write decompressed data to response
				w.Header().Set("X-Cache-Hit", "true")
				w.Header().Set("Content-Encoding", "gzip")
				io.Copy(w, reader)
				return
			} else if err != redis.Nil {
				log.Printf("Error retrieving from cache: %v. Proceeding without cache.", err)
			}
		}

		// Capture and compress response
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: new(bytes.Buffer)}
		next.ServeHTTP(rec, r)

		// Cache only successful responses (status codes 200-299) if Redis is available
		if c.client != nil && r.Method == http.MethodGet && rec.statusCode >= http.StatusOK && rec.statusCode < 300 {
			var compressedData bytes.Buffer
			gzipWriter := gzip.NewWriter(&compressedData)
			_, err := gzipWriter.Write(rec.body.Bytes())
			if err != nil {
				log.Printf("Error compressing cache data: %v. Proceeding without cache.", err)
				return
			}
			gzipWriter.Close()

			// Store compressed data in Redis
			err = c.client.Set(ctx, r.RequestURI, compressedData.Bytes(), c.ttl).Err()
			if err != nil {
				log.Printf("Error setting cache: %v. Proceeding without cache.", err)
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
