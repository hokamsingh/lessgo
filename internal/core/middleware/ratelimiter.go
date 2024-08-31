package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiterType defines the type of rate limiter (InMemory or RedisBacked).
type RateLimiterType int

const (
	InMemory RateLimiterType = iota
	RedisBacked
)

// RateLimiter is a middleware that limits the number of requests
// a client can make to your server within a specified interval.
type RateLimiter struct {
	limiterType     RateLimiterType
	limit           int
	interval        time.Duration
	redisClient     *redis.Client
	shards          []*shard
	numShards       int
	cleanupInterval time.Duration
	bufferPool      sync.Pool
}

// shard represents a partition of the request map to reduce lock contention.
type shard struct {
	requests map[string]*circularBuffer
	mu       sync.RWMutex
}

// circularBuffer is a fixed-size buffer for storing timestamps of requests.
type circularBuffer struct {
	timestamps []time.Time
	size       int
	start      int
	end        int
	full       bool
}

// NewRateLimiter creates and returns a new RateLimiter instance based on the provided configuration.
//
// The limiterType parameter determines whether an in-memory or Redis-backed rate limiter is used.
// The config parameter is either an InMemoryConfig or RedisConfig, depending on the limiterType.
func NewRateLimiter(limiterType RateLimiterType, config interface{}) *RateLimiter {
	switch limiterType {
	case InMemory:
		cfg := config.(InMemoryConfig)
		rl := &RateLimiter{
			limiterType:     InMemory,
			limit:           cfg.Limit,
			interval:        cfg.Interval,
			cleanupInterval: cfg.CleanupInterval,
			numShards:       cfg.NumShards,
			shards:          make([]*shard, cfg.NumShards),
			bufferPool: sync.Pool{
				New: func() interface{} {
					return &circularBuffer{
						timestamps: make([]time.Time, cfg.Limit),
						size:       cfg.Limit,
					}
				},
			},
		}
		for i := 0; i < cfg.NumShards; i++ {
			rl.shards[i] = &shard{
				requests: make(map[string]*circularBuffer),
			}
		}
		go rl.cleanup()
		return rl

	case RedisBacked:
		ctx := context.Background()
		cfg := config.(RedisConfig)
		client := &cfg.Client
		_, err := client.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Could not connect to Redis: %v", err)
		}
		return &RateLimiter{
			limiterType: RedisBacked,
			limit:       cfg.Limit,
			interval:    cfg.Interval,
			redisClient: client,
		}

	default:
		panic("Unsupported rate limiter type")
	}
}

// InMemoryConfig is the configuration for the in-memory rate limiter.
type InMemoryConfig struct {
	NumShards       int
	Limit           int
	Interval        time.Duration
	CleanupInterval time.Duration
}

func NewInMemoryConfig(NumShards int, Limit int, Interval time.Duration, CleanupInterval time.Duration) *InMemoryConfig {
	return &InMemoryConfig{
		NumShards:       NumShards,
		Limit:           Limit,
		Interval:        Interval,
		CleanupInterval: CleanupInterval,
	}
}

// RedisConfig is the configuration for the Redis-backed rate limiter.
type RedisConfig struct {
	Client   redis.Client
	Limit    int
	Interval time.Duration
}

func NewRedisConfig(client *redis.Client, limit int, interval time.Duration) *RedisConfig {
	return &RedisConfig{
		Client:   *client,
		Limit:    limit,
		Interval: interval,
	}
}

// Handle is the middleware function that processes incoming HTTP requests.
//
// It applies rate limiting based on the specified rate limiter type (in-memory or Redis-backed).
func (rl *RateLimiter) Handle(next http.Handler) http.Handler {
	switch rl.limiterType {
	case InMemory:
		return rl.handleInMemory(next)
	case RedisBacked:
		return rl.handleRedis(next)
	default:
		panic("Unsupported rate limiter type")
	}
}

// handleInMemory handles rate limiting using an in-memory approach.
//
// It uses a circular buffer to store timestamps of requests and a sync.Pool to reuse buffers.
func (rl *RateLimiter) handleInMemory(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		now := time.Now()

		sh := rl.getShard(key)
		sh.mu.Lock()

		cb, exists := sh.requests[key]
		if !exists {
			cb = rl.bufferPool.Get().(*circularBuffer)
			sh.requests[key] = cb
		}

		count := 0
		for i := 0; i < cb.size; i++ {
			if cb.timestamps[i].IsZero() {
				break
			}
			if now.Sub(cb.timestamps[i]) < rl.interval {
				count++
			}
		}

		if count >= rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			sh.mu.Unlock()
			return
		}

		cb.add(now)
		sh.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// handleRedis handles rate limiting using a Redis-backed approach.
//
// It uses Redis sorted sets to store timestamps of requests and ensures rate limiting across distributed systems.
func (rl *RateLimiter) handleRedis(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		now := time.Now().UnixNano()
		ctx := context.Background()

		windowStart := now - rl.interval.Nanoseconds()

		pipe := rl.redisClient.TxPipeline()
		pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
		pipe.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(windowStart, 10))
		pipe.ZCard(ctx, key)
		pipe.Expire(ctx, key, rl.interval)

		_, err := pipe.Exec(ctx)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		reqCount, err := rl.redisClient.ZCard(ctx, key).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if int(reqCount) > rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getShard returns the shard corresponding to the provided key.
//
// Sharding helps in distributing the requests across multiple shards to reduce lock contention.
func (rl *RateLimiter) getShard(key string) *shard {
	hash := fnv32(key)
	return rl.shards[int(hash)%rl.numShards]
}

// add inserts a new timestamp into the circular buffer.
//
// It updates the buffer's pointers and handles buffer wrapping.
func (cb *circularBuffer) add(t time.Time) {
	cb.timestamps[cb.end] = t
	cb.end = (cb.end + 1) % cb.size
	if cb.end == cb.start {
		cb.start = (cb.start + 1) % cb.size
		cb.full = true
	}
}

// cleanup periodically removes expired entries from the in-memory rate limiter.
//
// Buffers that are no longer in use are returned to the buffer pool.
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.cleanupInterval)
		for _, sh := range rl.shards {
			sh.mu.Lock()
			for key, cb := range sh.requests {
				count := 0
				now := time.Now()
				for i := 0; i < cb.size; i++ {
					if cb.timestamps[i].IsZero() {
						break
					}
					if now.Sub(cb.timestamps[i]) < rl.interval {
						count++
					}
				}
				if count == 0 {
					rl.bufferPool.Put(cb)
					delete(sh.requests, key)
				}
			}
			sh.mu.Unlock()
		}
	}
}

// fnv32 is a hash function that computes a 32-bit FNV-1a hash.
//
// It is used to map keys to shards.
func fnv32(key string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)

	hash := uint32(offset32)
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= prime32
	}
	return hash
}
