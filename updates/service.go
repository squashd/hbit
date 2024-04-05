package updates

import (
	"github.com/SQUASHD/hbit"
	"github.com/gorilla/websocket"
)

type (
	Service struct {
		clients map[hbit.UserId]*websocket.Conn
	}
)

func NewService() *Service {
	return &Service{
		clients: make(map[hbit.UserId]*websocket.Conn),
	}
}

func (s *Service) RegisterConnection(userId hbit.UserId, ws *websocket.Conn) {
	s.clients[userId] = ws
	s.SendMessageToUser(userId, "connected")
}

func (s *Service) SendMessageToUser(userId hbit.UserId, message any) {
	conn, ok := s.clients[userId]
	if !ok {
		return
	}

	tagged := tagPayload(message)
	if tagged.tag == UNKNOWN {
		conn.WriteJSON("unknown message type")
	}
	conn.WriteJSON(message)
}

func (s *Service) BroadCast(pm *websocket.PreparedMessage) {
	for _, conn := range s.clients {
		conn.WritePreparedMessage(pm)
	}
}
