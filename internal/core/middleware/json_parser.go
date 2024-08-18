package middleware

import (
	"context"
	"encoding/json"
	"net/http"
)

func JSONParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			key := "jsonBody"
			r = r.WithContext(context.WithValue(r.Context(), key, body))
		}
		next.ServeHTTP(w, r)
	})
}
