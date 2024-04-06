package http

import (
	"net/http"
	"os"
)

func NewTaskOrchestrationRouter(client *http.Client) *http.ServeMux {
	router := http.NewServeMux()
	orchestrator := NewTaskOrchestrator(
		os.Getenv("TASK_SVC_URL"),
		os.Getenv("RPG_SVC_URL"),
		client,
	)
	authMiddleware := AuthChainMiddleware(GetUserIdFromHeader)
	router.HandleFunc("POST /{id}/done", authMiddleware(orchestrator.OrchestrateTaskDone))
	router.HandleFunc("POST /{id}/undone", authMiddleware(orchestrator.OrchestrateTaskUndo))

	return router
}
