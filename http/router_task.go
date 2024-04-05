package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/task"
)

func NewTaskRouter(svc task.UserTaskService) *http.ServeMux {
	router := http.NewServeMux()
	handler := newTaskHandler(svc)
	userGetter := GetUserIdFromHeader
	AuthMiddleware := AuthChainMiddleware(userGetter)

	router.HandleFunc("GET /", AuthMiddleware(handler.FindAll))
	router.HandleFunc("POST /", AuthMiddleware(handler.Create))
	router.HandleFunc("PUT /", AuthMiddleware(handler.Update))
	router.HandleFunc("DELETE /", AuthMiddleware(handler.Delete))
	return router
}
