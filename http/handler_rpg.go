package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg"
)

type (
	rpgHandler struct {
		rpgSvc rpg.Service
	}
)

func newRPGHandler(rpgSvc rpg.Service) *rpgHandler {
	return &rpgHandler{rpgSvc: rpgSvc}
}

func (h *rpgHandler) CalculateRewards(w http.ResponseWriter, r *http.Request) {
	var req rpg.TaskRewardRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	rewards, err := h.rpgSvc.CalculateRewards(r.Context(), req)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, rewards)
}

func (h *rpgHandler) UndoRewards(w http.ResponseWriter, r *http.Request) {
	var req rpg.TaskRewardRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	rewards, err := h.rpgSvc.UndoRewards(r.Context(), req)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, rewards)
}
