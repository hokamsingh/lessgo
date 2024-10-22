package middleware

import (
	"net/http"
	"regexp"
	"strings"
)

type XSSProtection struct{}

// Creates a new middleware for XSS protection
func NewXSSProtection() *XSSProtection {
	return &XSSProtection{}
}

// Regular expression to detect potentially harmful XSS patterns, including encoded variants
var unsafePattern = regexp.MustCompile(`(?i)<script.*?>|javascript:|data:text/html|onerror=|onload=|onclick=|<iframe>|<img src=|<object>|<embed>|eval\(|%3Cscript%3E|&#60;script&#62;|&#x3C;script&#x3E;|&lt;script&gt;`)

// Middleware to handle requests and check for XSS attacks
func (xss *XSSProtection) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if containsXSS(r) {
			http.Error(w, "XSS detected", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// containsXSS checks various parts of the request for XSS payloads
func containsXSS(r *http.Request) bool {
	// Check URL query parameters
	for _, values := range r.URL.Query() {
		for _, value := range values {
			if isXSS(value) {
				return true
			}
		}
	}

	// Check form values
	if err := r.ParseForm(); err == nil {
		for _, values := range r.Form {
			for _, value := range values {
				if isXSS(value) {
					return true
				}
			}
		}
	}

	// Check cookie values
	for _, cookie := range r.Cookies() {
		if isXSS(cookie.Value) {
			return true
		}
	}

	// Check header values
	for _, values := range r.Header {
		for _, value := range values {
			if isXSS(value) {
				return true
			}
		}
	}

	return false
}

// isXSS checks if a string contains potentially harmful XSS payloads using regular expressions
func isXSS(value string) bool {
	valueLower := strings.ToLower(value)
	return unsafePattern.MatchString(valueLower)
}
