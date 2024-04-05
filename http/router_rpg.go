package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
)

func NewRPGRouter(charSvc character.CharacterService, questSvc quest.QuestService) *http.ServeMux {
	rpgRouter := http.NewServeMux()
	rpgHandler := newCharacterHandler(charSvc)

	userGetter := GetUserIdFromHeader
	AuthMiddleware := AuthChainMiddleware(userGetter)

	rpgRouter.HandleFunc("GET /characters", AuthMiddleware(rpgHandler.CharacterGet))
	rpgRouter.HandleFunc("POST /characters", AuthMiddleware(rpgHandler.CharacterCreate))
	rpgRouter.HandleFunc("PUT /characters", AuthMiddleware(rpgHandler.CharacterUpdate))
	rpgRouter.HandleFunc("DELETE /characters", AuthMiddleware(rpgHandler.CharacterDelete))

	questHandler := newQuestHandler(questSvc)
	rpgRouter.HandleFunc("GET /quests", AuthMiddleware(questHandler.GetAll))

	return rpgRouter
}
