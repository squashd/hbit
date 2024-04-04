package http

import (
	"log"
	"net/http"
	"time"
)

// Middleware type definition
type Middleware func(http.Handler) http.Handler

// ChainMiddleware chains multiple middleware functions together
func ChainMiddleware(handler http.Handler, middleware ...Middleware) http.Handler {
	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}

// customResponseWriter is a wrapper around an http.ResponseWriter that keeps track of the response status code
// deprecated
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewCustomResponseWriter was middleware used to note the status code of the response
// Howver, it causes issues when using an API Gateway
func NewCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, http.StatusOK}
}

// WriteHeader implements the http.ResponseWriter interface
func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// LoggerMiddleware logs the request method, URL path, status code, and duration of the request
// deprecated
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		crw := NewCustomResponseWriter(w)
		next.ServeHTTP(crw, r)
		log.Println(r.Method, r.URL.Path, crw.statusCode, time.Since(start))
	})
}

// AuthedHandler is a type definition for a handler that requires authentication
// Most routes go through this middleware
type AuthedHandler func(w http.ResponseWriter, r *http.Request, userId string)

func AuthMiddleware(next AuthedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := getUserIdFromHeader(r)
		next(w, r, userId)
	})
}

func getUserIdFromHeader(r *http.Request) string {
	return r.Header.Get("X-User-Id")
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		next.ServeHTTP(w, r)
	})
}
