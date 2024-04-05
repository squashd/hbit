package http

import (
	"encoding/json"
	"net/http"

	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/rpg/character"
)

// TODO: Aggregate results from both services and return a single response
func Register(
	authSvc auth.Service,
	charSvc character.CharacterService,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerForm auth.CreateUserForm

		decode := json.NewDecoder(r.Body)
		if err := decode.Decode(&registerForm); err != nil {
			Error(w, r, err)
			return
		}

		user, err := authSvc.Register(r.Context(), registerForm)
		if err != nil {
			Error(w, r, err)
			return
		}
		char, err := charSvc.CreateCharacter(r.Context(), character.CreateCharacterForm{})
		if err != nil {
			Error(w, r, err)
			return
		}
		resp := struct {
			User      auth.AuthDTO
			Character character.DTO
		}{
			User:      user,
			Character: char,
		}

		respondWithJSON(w, http.StatusCreated, resp)

	}
}
