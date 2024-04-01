package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/quest"
	"github.com/SQUASHD/hbit/quest/database"
)

type QuestHandler struct {
	questSvc quest.Service
}

func NewQuestHandler(questSvc quest.Service) *QuestHandler {
	return &QuestHandler{questSvc: questSvc}
}

func (h *QuestHandler) GetAll(w http.ResponseWriter, r *http.Request, userId string) {
	quests, err := h.questSvc.ListQuests(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, quests)
}

func (h *QuestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data database.CreateQuestParams
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

func (h *QuestHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var data database.UpdateQuestParams
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

func (h *QuestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.questSvc.DeleteQuest(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
