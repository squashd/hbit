package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit"
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

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	dto, err := h.authSvc.Register(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	setAccessCookie(w, dto.AccessToken, h.jwtConf.AccessDuration)
	setRefreshCookie(w, dto.RefreshToken, h.jwtConf.RefreshDuration)

	// Since we're setting the tokens in the cookie, we can just return the username
	// in the http response
	respondWithJSON(w, http.StatusOK, dto.Username)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var form auth.LoginForm

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	loginDto, err := h.authSvc.Login(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	setAccessCookie(w, loginDto.AccessToken, h.jwtConf.AccessDuration)
	setRefreshCookie(w, loginDto.RefreshToken, h.jwtConf.RefreshDuration)

	// Since we're setting the tokens in the cookie, we can just return the username
	// in the http response
	respondWithJSON(w, http.StatusOK, loginDto.Username)
}

func (h *authHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	clearTokensFromCookie(w)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
}

func (h *authHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	var form auth.RevokeTokenForm

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	err := h.authSvc.RevokeToken(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully revoked token"})
}

func Verify(svc auth.Service, jwtConf config.JwtOptions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := authenticateUser(w, r, svc, jwtConf)
		if err != nil {
			Error(w, r, err)
			return
		}

		authRes := struct {
			message string
		}{
			message: "Successfully verified token",
		}

		respondWithJSON(w, http.StatusOK, authRes)
	}
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
