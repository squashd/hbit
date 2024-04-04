package http

import (
	"net/http"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/updates"
	"github.com/gorilla/websocket"
)

type updatesServiceHandler struct {
	svc *updates.Service
}

func newUpdatesServiceHandler(svc *updates.Service) *updatesServiceHandler {
	return &updatesServiceHandler{svc: svc}
}

func (h *updatesServiceHandler) ConnectToWS(w http.ResponseWriter, r *http.Request, userId string) {
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
	h.svc.RegisterConnection(hbit.UserId(userId), conn)
}

func (h *updatesServiceHandler) TestConncetion(w http.ResponseWriter, r *http.Request, userId string) {
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

	h.svc.SendMessageToUser(hbit.UserId(userId), advancedMsg)
}
