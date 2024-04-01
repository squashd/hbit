package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/user"
	"github.com/SQUASHD/hbit/user/database"
)

func (s *ServerMonolith) handleSettingsGet(w http.ResponseWriter, r *http.Request, requestedById string) {
	settings, err := s.userSvc.GetSettings(r.Context(), requestedById)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, settings)
}

func (s *ServerMonolith) handleSettingsUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
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

	settings, err := s.userSvc.UpdateSettings(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, settings)
}
