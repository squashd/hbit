package http

import (
	"net/http"
)

func (s *Server) handleAchievementsGet(w http.ResponseWriter, r *http.Request, userId string) {
	achievements, err := s.achSvc.GetUserAchievements(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, achievements)
}
