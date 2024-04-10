package http

import (
	"net/http"
	"os"
)

func NewTaskOrchestrationRouter(client *http.Client) *http.ServeMux {
	router := http.NewServeMux()

	// I dinnae like this
	orchestrator := NewTaskOrchestrator(
		os.Getenv("TASK_SVC_URL"),
		os.Getenv("RPG_SVC_URL"),
		client,
	)

	userGetter := GetUserIdFromHeader
	AuthMiddleware := AuthChainMiddleware(userGetter)

	router.HandleFunc("POST /{id}/done", AuthMiddleware(orchestrator.OrchestrateTaskDone))
	router.HandleFunc("POST /{id}/undo", AuthMiddleware(orchestrator.OrchestrateTaskUndo))

	return router
}
