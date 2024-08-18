package middleware

import (
	"context"
	"net/http"
)

func CookieParser(next http.Handler) http.Handler {
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
