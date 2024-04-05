package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/config"
)

type authHandler struct {
	authSvc auth.Service
	jwtConf config.JwtOptions
}

func newAuthHandler(authSvc auth.Service, jwtConf config.JwtOptions) *authHandler {
	return &authHandler{
		authSvc: authSvc,
		jwtConf: jwtConf,
	}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var form auth.CreateUserForm

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		Error(w, r, err)
		return
	}

	dto, err := h.authSvc.Register(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	setAccessCookie(w, dto.AccessToken, h.jwtConf.AccessDuration)
	setRefreshCookie(w, dto.RefreshToken, h.jwtConf.RefreshDuration)
	respondWithJSON(w, http.StatusOK, dto)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var form auth.LoginForm

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		Error(w, r, err)
		return
	}

	loginDto, err := h.authSvc.Login(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	setAccessCookie(w, loginDto.AccessToken, h.jwtConf.AccessDuration)
	setRefreshCookie(w, loginDto.RefreshToken, h.jwtConf.RefreshDuration)

	respondWithJSON(w, http.StatusOK, loginDto)
}

func (h *authHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	clearTokensFromCookie(w)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
}

func (h *authHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	var form auth.RevokeTokenForm

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		Error(w, r, err)
		return
	}

	err = h.authSvc.RevokeToken(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully revoked token"})
}
func (h *authHandler) Verify(w http.ResponseWriter, r *http.Request) {
	userId, err := authenticateUser(w, r, h.authSvc, h.jwtConf)
	if err != nil {
		Error(w, r, err)
		return
	}
	r.Header.Set("X-User-Id", userId)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully verified token"})
}

func (h *authHandler) AdminRouterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := h.authSvc.AuthenticateUser(r.Context(), getAccessTokenFromCookie(r))
		if err != nil {
			clearTokensFromCookie(w)
			Error(w, r, err)
			return
		}
		_, err = h.authSvc.IsAdmin(r.Context(), userId)
		if err != nil {
			Error(w, r, err)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
