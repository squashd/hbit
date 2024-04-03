package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
)

func NewRPGRouter(charSvc character.CharacterService, questSvc quest.QuestService) *http.ServeMux {
	router := http.NewServeMux()
	charHandler := newCharacterHandler(charSvc)

	router.HandleFunc("GET /characters/{id}", AuthMiddleware(charHandler.CharacterGet))
	router.HandleFunc("POST /characters", AuthMiddleware(charHandler.CharacterCreate))
	router.HandleFunc("PUT /characters/{id}", AuthMiddleware(charHandler.CharacterUpdate))
	router.HandleFunc("DELETE /characters/{id}", AuthMiddleware(charHandler.CharacterDelete))

	questHandler := newQuestHandler(questSvc)
	router.HandleFunc("GET /quests", AuthMiddleware(questHandler.GetAll))

	rpgRouter := http.NewServeMux()
	rpgRouter.Handle("/rpg/", http.StripPrefix("/rpg", router))

	return rpgRouter
}
