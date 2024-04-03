package http

import (
	"net/http"

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

	RespondWithJSON(w, http.StatusOK, settings)
}

func (h *userHandler) SettingsUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data userdb.UpdateUserSettingsParams

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := user.UpdateSettingsForm{
		UpdateUserSettingsParams: data,
		RequestedById:            requestedById,
	}

	settings, err := h.userSvc.UpdateSettings(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, settings)
}
