package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/feat"
)

type featHandler struct {
	featSvc feat.UserFeatService
}

func newFeatHandler(svc feat.UserFeatService) *featHandler {
	return &featHandler{featSvc: svc}
}

func (h *featHandler) FeatsGet(w http.ResponseWriter, r *http.Request, userId string) {
	achievements, err := h.featSvc.GetUserFeats(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, achievements)
}
