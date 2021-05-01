package logic

import "nhooyr.io/websocket"

// store those global variables who are used by many other modules.
var userAccountDB = initUserInfoRedis()

// TODO: conn register management
var userConnRegister = make(map[UUID]*websocket.Conn)
