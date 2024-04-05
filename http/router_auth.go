package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/config"
)

func NewAuthRouter(svc auth.Service, jwtConf config.JwtOptions) *http.ServeMux {
	router := http.NewServeMux()
	authHandler := newAuthHandler(svc, jwtConf)
	router.HandleFunc("POST /login", authHandler.Login)
	router.HandleFunc("POST /register", authHandler.Register)
	return router
}
