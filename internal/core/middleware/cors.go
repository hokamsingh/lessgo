package middleware

import (
	"net/http"
)

// CORSOptions defines the configuration for the CORS middleware
type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// CORSMiddleware is the middleware that handles CORS
type CORSMiddleware struct {
	options CORSOptions
}

// NewCORSMiddleware creates a new instance of CORSMiddleware
func NewCORSMiddleware(options CORSOptions) *CORSMiddleware {
	return &CORSMiddleware{options: options}
}

// Handle sets the CORS headers on the response and restricts methods
func (cm *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedMethods := cm.getAllowedMethods()
		allowedMethodsMap := make(map[string]bool)
		for _, method := range allowedMethods {
			allowedMethodsMap[method] = true
		}

		if _, ok := allowedMethodsMap[r.Method]; !ok && r.Method != http.MethodOptions {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", cm.getAllowedOrigins())
		w.Header().Set("Access-Control-Allow-Methods", cm.getAllowedMethodsHeader())
		w.Header().Set("Access-Control-Allow-Headers", cm.getAllowedHeaders())

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cm *CORSMiddleware) getAllowedOrigins() string {
	if len(cm.options.AllowedOrigins) == 0 {
		return "*"
	}
	return stringJoin(cm.options.AllowedOrigins, ", ")
}

func (cm *CORSMiddleware) getAllowedMethods() []string {
	if len(cm.options.AllowedMethods) == 0 {
		return []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	return cm.options.AllowedMethods
}

func (cm *CORSMiddleware) getAllowedMethodsHeader() string {
	return stringJoin(cm.getAllowedMethods(), ", ")
}

func (cm *CORSMiddleware) getAllowedHeaders() string {
	if len(cm.options.AllowedHeaders) == 0 {
		return "Content-Type, Authorization"
	}
	return stringJoin(cm.options.AllowedHeaders, ", ")
}

func stringJoin(elems []string, sep string) string {
	result := ""
	for i, elem := range elems {
		if i > 0 {
			result += sep
		}
		result += elem
	}
	return result
}

// NewCorsOptions creates a new CORSOptions instance
func NewCorsOptions(origins []string, methods []string, headers []string) *CORSOptions {
	return &CORSOptions{
		AllowedOrigins: origins,
		AllowedMethods: methods,
		AllowedHeaders: headers,
	}
}
