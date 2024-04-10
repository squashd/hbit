package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/task"
)

func NewTaskRouter(svc task.Service) *http.ServeMux {
	router := http.NewServeMux()

	userGetter := GetUserIdFromHeader
	authMiddleware := AuthChainMiddleware(userGetter)

	// User-facing endpoints
	userTaskHandler := newTaskHandler(svc)
	router.HandleFunc("GET /", authMiddleware(userTaskHandler.FindAll))
	router.HandleFunc("POST /", authMiddleware(userTaskHandler.Create))
	router.HandleFunc("PUT /", authMiddleware(userTaskHandler.Update))
	router.HandleFunc("DELETE /", authMiddleware(userTaskHandler.Delete))

	// Internal inter-service endpoint (via task-rpg orchestrator)
	taskResolver := newTaskResolutionHandler(svc)
	router.HandleFunc("POST /done", internalAuthMiddleware(taskResolver.Done))
	router.HandleFunc("POST /undo", internalAuthMiddleware(taskResolver.Undone))

	return router
}
