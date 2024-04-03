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

	RespondWithJSON(w, http.StatusOK, dto)
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

	RespondWithJSON(w, http.StatusOK, loginDto)
}

func (h *authHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	ClearTokensFromCookie(w)
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
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

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully revoked token"})
}
func (h *authHandler) Verify(w http.ResponseWriter, r *http.Request) {
	userId, err := authenticateUser(w, r, h.authSvc, h.jwtConf)
	if err != nil {
		Error(w, r, err)
		return
	}
	r.Header.Set("X-User-Id", userId)
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully verified token"})
}

func (h *authHandler) AdminRouterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := h.authSvc.AuthenticateUser(r.Context(), GetAccessTokenFromCookie(r))
		if err != nil {
			ClearTokensFromCookie(w)
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

func (h *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := authenticateUser(w, r, h.authSvc, h.jwtConf)
		if err != nil {
			Error(w, r, err)
			return
		}
		w.Header().Set("X-User-Id", userId)
		next.ServeHTTP(w, r)
	})
}

func authenticateUser(w http.ResponseWriter, r *http.Request, svc auth.Service, jwtConf config.JwtOptions) (userID string, err error) {
	refreshToken := GetRefreshTokenFromCookie(r)
	accessToken := GetAccessTokenFromCookie(r)
	if refreshToken == "" {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "Missing refresh token"}
	}
	if accessToken == "" {
		accessToken, userID, err = svc.RefreshToken(r.Context(), refreshToken)
		if err != nil {
			return "", err
		}
		SetAccessCookie(w, accessToken, jwtConf.AccessDuration)
	} else {
		userID, err = svc.AuthenticateUser(r.Context(), accessToken)
		if err != nil {
			return "", err
		}
	}
	return userID, nil
}
