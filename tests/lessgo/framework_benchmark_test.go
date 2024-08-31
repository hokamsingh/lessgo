package lessgo_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

// BenchmarkHandler benchmarks the /ping handler in the lessgo framework.
func BenchmarkHandler(b *testing.B) {
	// Initialize the lessgo app with necessary middlewares
	corsOptions := LessGo.NewCorsOptions(
		[]string{"*"},
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		[]string{"Content-Type", "Authorization"},
	)

	size := LessGo.ConvertToBytes(int64(1024), LessGo.Kilobytes)
	parserOptions := LessGo.NewParserOptions(size * 5)

	rClient := LessGo.NewRedisClient("localhost:6379")
	App := LessGo.App(
		LessGo.WithCORS(*corsOptions),
		LessGo.WithJSONParser(*parserOptions),
		LessGo.WithCookieParser(),
		LessGo.WithCsrf(),
		LessGo.WithXss(),
		LessGo.WithCaching(rClient, 5*time.Minute, true),
		LessGo.WithRedisRateLimiter(rClient, 100, 1*time.Second),
	)

	// Add a simple /ping route
	App.Get("/ping", func(ctx *LessGo.Context) {
		ctx.Send("pong")
	})

	// Create a request to benchmark the /ping handler
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	log.Println("Starting benchmark")
	for i := 0; i < b.N; i++ {
		log.Printf("Iteration: %d", i)
		App.Mux.ServeHTTP(w, req)
	}
	log.Println("Benchmark completed")
}
