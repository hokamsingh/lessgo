package middleware

import (
	"net/http"
	"strings"
)

type XSSProtection struct{}

func NewXSSProtection() *XSSProtection {
	return &XSSProtection{}
}

func (xss *XSSProtection) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, values := range r.URL.Query() {
			for _, value := range values {
				if strings.Contains(value, "<script>") {
					http.Error(w, "XSS detected", http.StatusBadRequest)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
