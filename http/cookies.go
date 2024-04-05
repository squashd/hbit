package http

import (
	"net/http"
	"time"
)

const (
	ACCESS_TOKEN  = "access_jwt"
	REFRESH_TOKEN = "refresh_jwt"
)

func clearAccessCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     ACCESS_TOKEN,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(time.Duration(-1) * time.Second),
	})
}

func clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_TOKEN,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(time.Duration(-1) * time.Second),
	})
}

func clearTokensFromCookie(w http.ResponseWriter) {
	clearAccessCookie(w)
	clearRefreshCookie(w)
}

func setAccessCookie(w http.ResponseWriter, token string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     ACCESS_TOKEN,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Secure:   true,
		Domain:   "localhost",
		Expires:  time.Now().UTC().Add(duration),
	})
}

func setRefreshCookie(w http.ResponseWriter, token string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_TOKEN,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().UTC().Add(duration),
	})
}

func getAccessTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(ACCESS_TOKEN)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getRefreshTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(REFRESH_TOKEN)
	if err != nil || cookie == nil {
		return ""
	}
	return cookie.Value
}
