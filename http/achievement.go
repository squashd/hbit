package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/achievement"
)

type AchievementHandler struct {
	achSvc achievement.Service
}

func NewAchievementHandler(svc achievement.Service) *AchievementHandler {
	return &AchievementHandler{achSvc: svc}
}

func (h *AchievementHandler) AchievementsGet(w http.ResponseWriter, r *http.Request, userId string) {
	achievements, err := h.achSvc.GetUserAchievements(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, achievements)
}
