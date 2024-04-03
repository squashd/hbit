package updates

import (
	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn) // userId -> connection

func registerConnection(userId string, ws *websocket.Conn) {
	clients[userId] = ws
}

func sendMessageToUser(userId string, message interface{}) {
	if conn, ok := clients[userId]; ok {
		conn.WriteJSON(message)
	}
}
