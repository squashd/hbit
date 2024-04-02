package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type questHandler struct {
	questSvc quest.Service
}

func newQuestHandler(questSvc quest.Service) *questHandler {
	return &questHandler{questSvc: questSvc}
}

func (h *questHandler) GetAll(w http.ResponseWriter, r *http.Request, userId string) {
	quests, err := h.questSvc.ListQuests(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, quests)
}

func (h *questHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data rpgdb.CreateQuestParams
	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	quest, err := h.questSvc.CreateQuest(r.Context(), data)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, quest)
}

func (h *questHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var data rpgdb.UpdateQuestParams
	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	data.ID = id

	quest, err := h.questSvc.UpdateQuest(r.Context(), data)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, quest)
}

func (h *questHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.questSvc.DeleteQuest(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
