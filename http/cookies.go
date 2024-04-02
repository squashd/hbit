package http

import (
	"net/http"
	"time"
)

const (
	ACCESS_TOKEN  = "access_jwt"
	REFRESH_TOKEN = "refresh_jwt"
)

func ClearAccessCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     ACCESS_TOKEN,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(time.Duration(-1) * time.Second),
	})
}

func ClearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_TOKEN,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(time.Duration(-1) * time.Second),
	})
}

func ClearTokensFromCookie(w http.ResponseWriter) {
	ClearAccessCookie(w)
	ClearRefreshCookie(w)
}

func SetAccessCookie(w http.ResponseWriter, token string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     ACCESS_TOKEN,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Secure:   true,
		Expires:  time.Now().Add(duration),
	})
}

func SetRefreshCookie(w http.ResponseWriter, token string, duration int) {
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_TOKEN,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(duration) * time.Second),
	})
}

func GetAccessTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(ACCESS_TOKEN)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func GetRefreshTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(REFRESH_TOKEN)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func GetCookieValue(r *http.Request, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
