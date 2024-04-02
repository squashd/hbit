package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type characterHandler struct {
	charSvc character.Service
}

func newCharacterHandler(charSvc character.Service) *characterHandler {
	return &characterHandler{charSvc: charSvc}
}

func (h *characterHandler) CharacterGetAll(w http.ResponseWriter, r *http.Request, requestedById string) {
	characters, err := h.charSvc.List(r.Context())
	if err != nil {
		Error(w, r, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, characters)
}

func (h *characterHandler) CharacterGet(w http.ResponseWriter, r *http.Request, requestedById string) {
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

func (h *characterHandler) CharacterCreate(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data rpgdb.CreateCharacterParams

	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.CreateCharacterForm{
		CreateCharacterParams: data,
		RequestedById:         requestedById,
	}

	character, err := h.charSvc.Create(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, character)
}

func (h *characterHandler) CharacterUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	var data rpgdb.UpdateCharacterParams
	if err := Decode(r, &data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.UpdateCharacterForm{
		UpdateCharacterParams: data,
		RequestedById:         requestedById,
		CharacterId:           id,
	}

	character, err := h.charSvc.Update(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, character)
}

func (h *characterHandler) CharacterDelete(w http.ResponseWriter, r *http.Request, requestedById string) {
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
