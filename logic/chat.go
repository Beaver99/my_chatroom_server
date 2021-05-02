package logic

import "nhooyr.io/websocket"

type UUID string

type Chat struct{
	id UUID

	conn *websocket.Conn

	// indicate whether the counterpart is a single user("0") or a group("1")
	mode string

	counterID UUID
}

