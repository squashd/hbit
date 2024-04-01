package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit/auth"
)

func (s *ServerMonolith) handleRegister(w http.ResponseWriter, r *http.Request) {
	var form auth.CreateUserForm

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		Error(w, r, err)
		return
	}

	dto, err := s.authSvc.Register(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, dto)
}

func (s *ServerMonolith) handleLogin(w http.ResponseWriter, r *http.Request) {
	var form auth.LoginForm

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		Error(w, r, err)
		return
	}

	loginDto, err := s.authSvc.Login(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, loginDto)
}

func (s *ServerMonolith) handleSignOut(w http.ResponseWriter, r *http.Request) {
	ClearTokensFromCookie(w)
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
}

func (s *ServerMonolith) handleRevoke(w http.ResponseWriter, r *http.Request) {
	var form auth.RevokeTokenForm

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		Error(w, r, err)
		return
	}

	err = s.authSvc.RevokeToken(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully revoked token"})
}

func (s *ServerMonolith) registerAuthRoutes(router *http.ServeMux) http.Handler {
	return router
}
