package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/task"
)

func NewTaskRouter(svc task.UserTaskService) *http.ServeMux {
	router := http.NewServeMux()
	handler := newTaskHandler(svc)
	router.HandleFunc("GET /tasks", AuthMiddleware(handler.FindAll))
	router.HandleFunc("POST /tasks", AuthMiddleware(handler.Create))
	router.HandleFunc("PUT /tasks", AuthMiddleware(handler.Update))
	router.HandleFunc("DELETE /tasks", AuthMiddleware(handler.Delete))
	router.HandleFunc("GET /test", handler.Test)
	return router
}
