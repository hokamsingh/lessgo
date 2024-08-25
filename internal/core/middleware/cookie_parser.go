package middleware

import (
	"context"
	"net/http"
)

type CookieParser struct{}

func NewCookieParser() *CookieParser {
	return &CookieParser{}
}

func (cp *CookieParser) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}
		r = r.WithContext(context.WithValue(r.Context(), "cookies", cookieMap))
		next.ServeHTTP(w, r)
	})
}
