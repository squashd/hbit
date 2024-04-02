package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/auth"
)

func NewAuthRouter(svc auth.Service) *http.ServeMux {
	router := http.NewServeMux()
	authHandler := newAuthHandler(svc)
	router.HandleFunc("POST /login", authHandler.Login)
	router.HandleFunc("POST /register", authHandler.Register)
	return router
}
