package updates

import "github.com/gorilla/websocket"

type (
	Service interface {
		RegisterConnection(userId string, ws *websocket.Conn)
		SendMessageToUser(userId string, message any)
		BroadCast(message any)
	}
)
