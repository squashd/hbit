package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/task"
)

func NewTaskRouter(svc task.UserTaskService) *http.ServeMux {
	router := http.NewServeMux()
	handler := newTaskHandler(svc)
	userGetter := GetUserIdFromHeader
	authMiddleware := AuthChainMiddleware(userGetter)

	router.HandleFunc("GET /", authMiddleware(handler.FindAll))
	router.HandleFunc("POST /", authMiddleware(handler.Create))
	router.HandleFunc("PUT /", authMiddleware(handler.Update))
	router.HandleFunc("DELETE /", authMiddleware(handler.Delete))

	router.HandleFunc("POST /done", internalAuthMiddleware(handler.Done))
	router.HandleFunc("POST /undo", internalAuthMiddleware(handler.Undone))
	return router
}
