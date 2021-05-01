package logic

import "nhooyr.io/websocket"

func GetUserConnMap() map[UUID]*websocket.Conn {
	return userConnRegister
}
