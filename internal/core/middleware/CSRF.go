package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
)

type CSRFProtection struct{}

func NewCSRFProtection() *CSRFProtection {
	return &CSRFProtection{}
}

func (csrf *CSRFProtection) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Generate and set CSRF token for GET requests
			token, err := GenerateCSRFToken()
			if err != nil {
				http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
				return
			}
			SetCSRFCookie(w, token)
		} else if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			// Validate CSRF token for state-changing requests
			if !ValidateCSRFToken(r) {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func GenerateCSRFToken() (string, error) {
	token := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

// SetCSRFCookie sets a CSRF token as a secure cookie.
func SetCSRFCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // Prevent access from JavaScript
		Secure:   true, // Ensure the cookie is only sent over HTTPS
	})
}

func ValidateCSRFToken(r *http.Request) bool {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return false
	}
	csrfToken := r.Header.Get("X-CSRF-Token") // Or retrieve from form data
	return csrfToken == cookie.Value
}
