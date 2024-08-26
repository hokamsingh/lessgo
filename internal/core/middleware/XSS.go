package middleware

import (
	"html"
	"net/http"
	"strings"
)

type XSSProtection struct{}

func NewXSSProtection() *XSSProtection {
	return &XSSProtection{}
}

func (xss *XSSProtection) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if containsXSS(r) {
			http.Error(w, "XSS detected", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// containsXSS checks various parts of the request for XSS payloads.
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

	// Check cookies
	for _, cookie := range r.Cookies() {
		if isXSS(cookie.Value) {
			return true
		}
	}

	// Check headers
	for _, values := range r.Header {
		for _, value := range values {
			if isXSS(value) {
				return true
			}
		}
	}

	return false
}

// isXSS checks if a string contains potentially harmful XSS payloads.
func isXSS(value string) bool {
	// Define a more comprehensive list of unsafe patterns
	unsafePatterns := []string{
		"<script>",
		"javascript:",
		"data:text/html",
		"vbscript:",
		"mocha:",
		"onerror=",
		"onload=",
		"onclick=",
		"onmouseover=",
		"onfocus=",
		"onchange=",
		"onsubmit=",
		"onreset=",
		"onabort=",
		"<iframe>",
		"<img src=",
		"<object>",
		"<embed>",
		"<style>",
		"<link>",
		"<meta>",
		"document.cookie",
		"window.location",
		"self.location",
		"eval(",
		"<!--",
		"--!>",
		"<![CDATA[",
		"svg/onload=",
		"math:xmlns",
		"data:",
		"ftp://",
		"file://",
		"%3Cscript%3E",
		"&#60;script&#62;",
		"&#x3C;script&#x3E;",
	}

	valueLower := strings.ToLower(value)
	for _, pattern := range unsafePatterns {
		if strings.Contains(valueLower, pattern) {
			return true
		}
	}

	// Use HTML escaping as an additional check
	escaped := html.EscapeString(value)
	return escaped != value
}
