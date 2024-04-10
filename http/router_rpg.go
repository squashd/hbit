package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
)

func NewRPGRouter(
	charSvc character.Service,
	questSvc quest.Service,
	rpgSvc rpg.Service,
) *http.ServeMux {
	rpgRouter := http.NewServeMux()
	charHandler := newCharacterHandler(charSvc)

	userIdGetter := GetUserIdFromHeader
	userAuthMiddleware := AuthChainMiddleware(userIdGetter)

	// TODO: Implement character creation and updating
	rpgRouter.HandleFunc("GET /characters", userAuthMiddleware(charHandler.CharacterGet))
	rpgRouter.HandleFunc("POST /characters", userAuthMiddleware(charHandler.CharacterCreate))
	rpgRouter.HandleFunc("PUT /characters", userAuthMiddleware(charHandler.CharacterUpdate))

	// TODO: Implement quests
	questHandler := newQuestHandler(questSvc)
	rpgRouter.HandleFunc("GET /quests", userAuthMiddleware(questHandler.GetAll))

	rpgHandler := newRPGHandler(rpgSvc)
	rpgRouter.HandleFunc("POST /rewards/calculate", internalAuthMiddleware(rpgHandler.CalculateRewards))
	rpgRouter.HandleFunc("POST /rewards/undo", internalAuthMiddleware(rpgHandler.UndoRewards))

	return rpgRouter
}
