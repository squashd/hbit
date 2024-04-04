package updates

import (
	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn) // userId -> connection

func registerConnection(userId string, ws *websocket.Conn) {
	clients[userId] = ws
}

func sendMessageToUser(userId string, message any) {
	conn, ok := clients[userId]
	if !ok {
		return
	}

	tagged := tagPayload(message)
	if tagged.tag == "unknown" {
		conn.WriteJSON(message)
	}
	conn.WriteJSON(message)
}
