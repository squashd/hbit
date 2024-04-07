package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/user"
	"github.com/SQUASHD/hbit/user/userdb"
)

type userHandler struct {
	userSvc user.Service
}

func NewUserHandler(userSvc user.Service) *userHandler {
	return &userHandler{userSvc: userSvc}
}

func (h *userHandler) SettingsGet(w http.ResponseWriter, r *http.Request, requestedById string) {
	settings, err := h.userSvc.GetSettings(r.Context(), requestedById)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, settings)
}

func (h *userHandler) SettingsUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
	var request userdb.UpdateUserSettingsParams

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	form := user.UpdateSettingsForm{
		UpdateUserSettingsParams: request,
		RequestedById:            requestedById,
	}

	settings, err := h.userSvc.UpdateSettings(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, settings)
}
