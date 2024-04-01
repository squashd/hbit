package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/user"
	"github.com/SQUASHD/hbit/user/database"
)

type UserHandler struct {
	userSvc user.Service
}

func NewUserHandler(userSvc user.Service) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) SettingsGet(w http.ResponseWriter, r *http.Request, requestedById string) {
	settings, err := h.userSvc.GetSettings(r.Context(), requestedById)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, settings)
}

func (h *UserHandler) SettingsUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")
	var data database.UpdateUserSettingsParams

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := user.UpdateSettingsForm{
		UpdateUserSettingsParams: data,
		UserId:                   id,
		RequestedById:            requestedById,
	}

	settings, err := h.userSvc.UpdateSettings(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, settings)
}
