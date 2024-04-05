package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type characterHandler struct {
	charSvc character.CharacterManagement
}

func newCharacterHandler(charSvc character.CharacterManagement) *characterHandler {
	return &characterHandler{charSvc: charSvc}
}

func (h *characterHandler) CharacterGet(w http.ResponseWriter, r *http.Request, requestedById string) {

	character, err := h.charSvc.GetCharacter(r.Context(), requestedById)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *characterHandler) CharacterCreate(w http.ResponseWriter, r *http.Request, requestedById string) {
	var data rpgdb.CreateCharacterParams

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.CreateCharacterForm{
		CreateCharacterParams: data,
		RequestedById:         requestedById,
	}

	form.CreateCharacterParams.UserID = requestedById

	character, err := h.charSvc.CreateCharacter(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, character)
}

func (h *characterHandler) CharacterUpdate(w http.ResponseWriter, r *http.Request, requestedById string) {
	id := r.PathValue("id")

	var data rpgdb.UpdateCharacterParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		Error(w, r, err)
		return
	}

	form := character.UpdateCharacterForm{
		UpdateCharacterParams: data,
		RequestedById:         requestedById,
		CharacterId:           id,
	}

	character, err := h.charSvc.UpdateCharacter(r.Context(), form)
	if err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *characterHandler) CharacterDelete(w http.ResponseWriter, r *http.Request, requestedById string) {

	if err := h.charSvc.DeleteCharacter(r.Context(), requestedById); err != nil {
		Error(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
