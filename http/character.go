package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/character"
)

type CharacterHandler struct {
	charSvc character.Service
}

func NewCharacterHandler(charSvc character.Service) *CharacterHandler {
	return &CharacterHandler{charSvc: charSvc}
}

func (h *CharacterHandler) CharacterGetAll(w http.ResponseWriter, r *http.Request, requestedById string) {
	characters, err := h.charSvc.List(r.Context())
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, characters)
}

func (h *CharacterHandler) CharacterGet(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	form := character.ReadCharacterForm{
		RequestedById: requestedById,
		CharacterId:   id,
	}

	character, err := h.charSvc.Read(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) CharacterCreate(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data character.CreateCharacterData

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.CreateCharacterForm{
		CreateCharacterData: data,
		RequestedById:       requestedById,
	}

	character, err := h.charSvc.Create(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, character)
}

func (h *CharacterHandler) CharacterUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
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

	character, err := h.charSvc.Update(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) CharacterDelete(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	form := character.DeleteCharacterForm{
		RequestedById: requestedById,
		CharacterId:   id,
	}

	if err := h.charSvc.Delete(r.Context(), form); err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
