package http

import (
	"net/http"

	"github.com/SQUASHD/hbit/updates"
)

func NewUpdatesRouter(s *updates.Service) *http.ServeMux {
	r := http.NewServeMux()
	handler := newUpdatesServiceHandler(s)
	r.HandleFunc("/ws", AuthMiddleware(handler.ConnectToWS))
	r.HandleFunc("GET /send", AuthMiddleware(handler.TestConncetion))

	return r
}
