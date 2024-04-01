package http

import (
	"log"
	"net/http"
	"time"

	"github.com/SQUASHD/hbit"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, middleware ...Middleware) http.Handler {
	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, http.StatusOK}
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		crw := NewCustomResponseWriter(w)
		next.ServeHTTP(crw, r)
		log.Println(r.Method, r.URL.Path, crw.statusCode, time.Since(start))
	})
}

type AuthedHandler func(w http.ResponseWriter, r *http.Request, userId string)

func (s *ServerMonolith) AuthMiddleware(next AuthedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := s.authenticateUser(w, r)
		if err != nil {
			ClearTokensFromCookie(w)
			Error(w, r, err)
			return
		}
		next(w, r, userId)
	})
}

func (s *ServerMonolith) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := s.authenticateUser(w, r)
		if err != nil {
			ClearTokensFromCookie(w)
			Error(w, r, err)
			return
		}
		_, err = s.authSvc.IsAdmin(r.Context(), userId)
		if err != nil {
			Error(w, r, err)
			return
		} else {
			next(w, r)
		}
	})
}

func (s *ServerMonolith) AdminRouterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := s.authenticateUser(w, r)
		if err != nil {
			ClearTokensFromCookie(w)
			Error(w, r, err)
			return
		}
		_, err = s.authSvc.IsAdmin(r.Context(), userId)
		if err != nil {
			Error(w, r, err)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (s *ServerMonolith) authenticateUser(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	refreshToken := GetRefreshTokenFromCookie(r)
	accessToken := GetAccessTokenFromCookie(r)
	if refreshToken == "" {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "Missing refresh token"}
	}
	if accessToken == "" {
		accessToken, userID, err = s.authSvc.RefreshToken(r.Context(), refreshToken)
		if err != nil {
			return "", err
		}
		SetAccessCookie(w, accessToken, s.jwtConf.AccessDuration)
	} else {
		userID, err = s.authSvc.AuthenticateUser(r.Context(), accessToken)
		if err != nil {
			return "", err
		}
	}
	return userID, nil
}
