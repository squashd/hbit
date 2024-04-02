package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
)

func NewRPGRouter(routers ...*http.ServeMux) *http.ServeMux {
	router := http.NewServeMux()
	for _, r := range routers {
		router.Handle("/", r)
	}
	return router
}

func NewCharacterRouter(svc character.Service) *http.ServeMux {
	router := http.NewServeMux()
	charHandler := newCharacterHandler(svc)
	router.HandleFunc("GET /characters/{id}", AuthMiddleware(charHandler.CharacterGet))
	router.HandleFunc("POST /characters", AuthMiddleware(charHandler.CharacterCreate))
	router.HandleFunc("PUT /characters/{id}", AuthMiddleware(charHandler.CharacterUpdate))
	router.HandleFunc("DELETE /characters/{id}", AuthMiddleware(charHandler.CharacterDelete))
	return router
}

func NewQuestRouter(svc quest.Service) *http.ServeMux {
	router := http.NewServeMux()
	questHandler := newQuestHandler(svc)
	router.HandleFunc("GET /quests", AuthMiddleware(questHandler.GetAll))
	return router
}
