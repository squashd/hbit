package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
)

func NewRPGRouter(charSvc character.CharacterService, questSvc quest.QuestService) *http.ServeMux {
	rpgRouter := http.NewServeMux()
	rpgHandler := newCharacterHandler(charSvc)

	userIdGetter := GetUserIdFromHeader
	userAuthMiddleware := AuthChainMiddleware(userIdGetter)

	rpgRouter.HandleFunc("GET /characters", userAuthMiddleware(rpgHandler.CharacterGet))
	rpgRouter.HandleFunc("POST /characters", userAuthMiddleware(rpgHandler.CharacterCreate))
	rpgRouter.HandleFunc("PUT /characters", userAuthMiddleware(rpgHandler.CharacterUpdate))

	questHandler := newQuestHandler(questSvc)
	rpgRouter.HandleFunc("GET /quests", userAuthMiddleware(questHandler.GetAll))

	return rpgRouter
}
