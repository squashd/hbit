package http

import (
	"log"
	"net/http"
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/config"
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

// LoggerMiddleware logs the request method, URL path, and duration of the request
// deprecated
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

// AuthedHandler is a type definition for a handler that requires authentication
// Most routes go through this middleware
type AuthedHandler func(w http.ResponseWriter, r *http.Request, userId string)

// AuthChainMiddleware is a higher order function that returns a middleware function that authenticates users
func AuthChainMiddleware(userIdGetter func(r *http.Request) (string, error)) func(next AuthedHandler) http.HandlerFunc {
	return func(next AuthedHandler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, err := userIdGetter(r)
			if err != nil {
				Error(w, r, err)
				return
			}
			next(w, r, userId)
		})
	}
}

// GetUserIdFromHeader is a helper function that extracts the X-User-Id header from a request
func GetUserIdFromHeader(r *http.Request) (string, error) {
	userId := r.Header.Get("X-User-Id")
	if userId == "" {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "Missing user id header"}
	}
	return userId, nil
}

// CORSMiddlware... the name's on the tin
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		next.ServeHTTP(w, r)
	})
}

// Refactor from being a method of authHandler
func JwtAuthMiddleware(svc auth.JwtAuth, jwtConf config.JwtOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, err := authenticateUser(w, r, svc, jwtConf)
			if err != nil {
				Error(w, r, err)
				return
			}
			r.Header.Add("X-User-Id", userId)
			next.ServeHTTP(w, r)
		})
	}
}

// authenticateUser is a helper function to JwtAuthMiddleware
func authenticateUser(w http.ResponseWriter, r *http.Request, svc auth.JwtAuth, jwtConf config.JwtOptions) (string, error) {
	refreshToken := getRefreshTokenFromCookie(r)
	accessToken := getAccessTokenFromCookie(r)
	// If refresh token is missing, clear all tokens and return error
	if refreshToken == "" {
		clearTokensFromCookie(w)
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "Missing refresh token"}
	}
	// If access token is missing, refresh token
	if accessToken == "" {
		accessToken, userId, err := svc.RefreshToken(r.Context(), refreshToken)
		if err != nil {
			return "", err
		}
		setAccessCookie(w, accessToken, jwtConf.AccessDuration)
		return userId, nil
	}
	// If both tokens are present, authenticate user
	userId, err := svc.AuthenticateUser(r.Context(), accessToken)
	if err != nil {
		return "", err
	}
	return userId, nil
}
