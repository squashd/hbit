package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/quest/database"
)

func (s *Server) handleQuestGetAll(w http.ResponseWriter, r *http.Request, userId string) {
	quests, err := s.questSvc.ListQuests(r.Context(), userId)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, quests)
}

func (s *Server) handleQuestCreate(w http.ResponseWriter, r *http.Request) {
	var data database.CreateQuestParams
	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	quest, err := s.questSvc.CreateQuest(r.Context(), data)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, quest)
}

func (s *Server) handleQuestUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var data database.UpdateQuestParams
	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	data.ID = id

	quest, err := s.questSvc.UpdateQuest(r.Context(), data)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, quest)
}

func (s *Server) handleQuestDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.questSvc.DeleteQuest(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
