package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit/auth"
)

type AuthHandler struct {
	authSvc auth.Service
}

func NewAuthHandler(authSvc auth.Service) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	ClearTokensFromCookie(w)
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
}

func (h *AuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) registerAuthRoutes(router *http.ServeMux) http.Handler {
	return router
}
