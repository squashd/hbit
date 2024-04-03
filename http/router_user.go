package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/user"
)

func NewUserRouter(svc user.Service) *http.ServeMux {
	router := http.NewServeMux()
	handler := NewUserHandler(svc)
	router.HandleFunc("GET /settings", AuthMiddleware(handler.SettingsGet))
	router.HandleFunc("PUT /settings", AuthMiddleware(handler.SettingsUpdate))
	return router
}
