package updates

import "github.com/gorilla/websocket"

type (
	service struct {
		clients map[string]*websocket.Conn
	}
)

func NewService() Service {
	return &service{
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *service) RegisterConnection(userId string, ws *websocket.Conn) {
	s.clients[userId] = ws
}

func (s *service) SendMessageToUser(userId string, message any) {
	if conn, ok := s.clients[userId]; ok {
		conn.WriteJSON(message)
	}
}

func (s *service) BroadCast(message any) {
	for _, conn := range s.clients {
		conn.WriteJSON(message)
	}
}
