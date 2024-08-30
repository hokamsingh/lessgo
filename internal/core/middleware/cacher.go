package middleware

import (
	"bytes"
	"context"
	"encoding/gob"
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

func NewCaching(client *redis.Client, ttl time.Duration, cacheControl bool) *Caching {
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
				// Cache hit: deserialize cached response
				var cachedResponse cachedResponse
				decoder := gob.NewDecoder(bytes.NewReader([]byte(data)))
				err := decoder.Decode(&cachedResponse)
				if err != nil {
					log.Printf("Error decoding cached response: %v", err)
					next.ServeHTTP(w, r)
					return
				}

				// Write cached headers
				for key, values := range cachedResponse.Headers {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}

				// Write cached body
				w.Header().Set("X-Cache-Hit", "true")
				io.WriteString(w, cachedResponse.Body)
				return
			} else if err != redis.Nil {
				log.Printf("Error retrieving from cache: %v", err)
			}
		}

		// Capture response
		rec := &ResponseRecorder{ResponseWriter: w, StatusCode: http.StatusOK, Body: new(bytes.Buffer)}
		next.ServeHTTP(rec, r)

		// Cache only successful responses (status code 200)
		if r.Method == http.MethodGet && rec.StatusCode == http.StatusOK {
			cachedResponse := cachedResponse{
				Headers: rec.Header(),
				Body:    rec.Body.String(),
			}

			var buffer bytes.Buffer
			encoder := gob.NewEncoder(&buffer)
			err := encoder.Encode(cachedResponse)
			if err != nil {
				log.Printf("Error encoding cached response: %v", err)
				return
			}

			err = c.client.Set(ctx, r.RequestURI, buffer.Bytes(), c.ttl).Err()
			if err != nil {
				log.Printf("Error setting cache: %v", err)
			}
		}
	})
}

// cachedResponse stores both headers and body
type cachedResponse struct {
	Headers http.Header
	Body    string
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       *bytes.Buffer
}

func (rec *ResponseRecorder) Write(p []byte) (int, error) {
	rec.Body.Write(p)                  // Write to the buffer
	return rec.ResponseWriter.Write(p) // Stream response to client
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Implement the Flush method
func (rec *ResponseRecorder) Flush() {
	if flusher, ok := rec.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func init() {
	// Register the cachedResponse type with gob so it can be encoded/decoded
	gob.Register(cachedResponse{})
}
