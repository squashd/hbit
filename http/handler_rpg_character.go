package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type characterHandler struct {
	charSvc character.UserCharacterService
}

func newCharacterHandler(charSvc character.UserCharacterService) *characterHandler {
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

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	form := character.CreateCharacterForm{
		CreateCharacterParams: data,
		RequestedById:         requestedById,
	}

	// I hate this
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

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
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

// deprecated as character deletion is event-based and handled by the top level
// rpg service
// func (h *characterHandler) CharacterDelete(w http.ResponseWriter, r *http.Request, requestedById string) {
// 	userId := r.PathValue("userId")
// 	if err := h.charSvc.DeleteCharacter(r.Context(), userId); err != nil {
// 		Error(w, r, err)
// 		return
// 	}
// 	respondWithJSON(w, http.StatusOK, nil)
// }
