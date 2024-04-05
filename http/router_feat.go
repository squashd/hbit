package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/feat"
)

func NewFeatRouter(svc feat.Service) *http.ServeMux {
	server := &featRouter{svc: svc}
	return server.RegisterRoutes()
}

type featRouter struct {
	svc feat.Service
}

func (s *featRouter) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	handler := newFeatHandler(s.svc)
	userGetter := GetUserIdFromHeader
	AuthMiddleware := AuthChainMiddleware(userGetter)
	mux.HandleFunc("/", AuthMiddleware(handler.FeatsGet))
	return mux
}
