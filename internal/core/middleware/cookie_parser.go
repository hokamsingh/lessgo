package middleware

import (
	"context"
	"net/http"
)

type CookieParser struct{}

func NewCookieParser() *CookieParser {
	return &CookieParser{}
}

type Cookies string

func (cp *CookieParser) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}
		cookiesKey := Cookies("cookies")
		r = r.WithContext(context.WithValue(r.Context(), cookiesKey, cookieMap))
		next.ServeHTTP(w, r)
	})
}
