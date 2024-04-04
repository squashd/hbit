package http

import (
	"net/http"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/updates"
	"github.com/gorilla/websocket"
)

func NewUpdatesRouter(s *updates.Service) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Error(w, r, &hbit.Error{Code: hbit.EINTERNAL, Message: err.Error()})
			return
		}
		s.RegisterConnection("1", conn)
	})
	r.HandleFunc("GET /send", func(w http.ResponseWriter, r *http.Request) {
		type Reward struct {
			Gold int `json:"gold"`
			Exp  int `json:"exp"`
			Mana int `json:"mana"`
		}
		reward := Reward{
			Gold: 100,
			Exp:  100,
			Mana: 100,
		}
		advancedMsg := struct {
			UserId string `json:"userId"`
			Reward Reward `json:"reward"`
		}{
			UserId: hbit.NewUUID(),
			Reward: reward,
		}

		s.SendMessageToUser("1", advancedMsg)
	})
	return r
}
