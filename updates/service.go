package updates

import "github.com/gorilla/websocket"

type (
	Service struct {
		clients map[string]*websocket.Conn
	}
)

func NewService() *Service {
	return &Service{
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *Service) RegisterConnection(userId string, ws *websocket.Conn) {
	s.clients[userId] = ws
}

func (s *Service) SendMessageToUser(userId string, message any) {
	if conn, ok := s.clients[userId]; ok {
		conn.WriteJSON(message)
	}
}

func (s *Service) BroadCast(pm *websocket.PreparedMessage) {
	for _, conn := range s.clients {
		conn.WritePreparedMessage(pm)
	}
}
