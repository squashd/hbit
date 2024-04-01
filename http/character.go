package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/character"
)

func (s *Server) handleCharacterGetAll(w http.ResponseWriter, r *http.Request, requestedById string) {
	characters, err := s.charSvc.List(r.Context())
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, characters)
}

func (s *Server) handleCharacterGet(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	form := character.ReadCharacterForm{
		RequestedById: requestedById,
		CharacterId:   id,
	}

	character, err := s.charSvc.Read(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, character)
}

func (s *Server) handleCharacterCreate(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data character.CreateCharacterData

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.CreateCharacterForm{
		CreateCharacterData: data,
		RequestedById:       requestedById,
	}

	character, err := s.charSvc.Create(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, character)
}

func (s *Server) handleCharacterUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	var data character.UpdateCharacterData
	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.UpdateCharacterForm{
		UpdateCharacterData: data,
		RequestedById:       requestedById,
		CharacterId:         id,
	}

	character, err := s.charSvc.Update(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, character)
}

func (s *Server) handleCharacterDelete(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	form := character.DeleteCharacterForm{
		RequestedById: requestedById,
		CharacterId:   id,
	}

	if err := s.charSvc.Delete(r.Context(), form); err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
