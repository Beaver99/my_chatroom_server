package logic

import (
	"nhooyr.io/websocket"
	"sync"
)

type userConnMap struct {
	uc sync.Map
}

func UserConnMapLoad(userID UUID) ( *websocket.Conn, bool) {
	ctemp, ok := userConnRegister.uc.Load(userID)
	c := ctemp.(*websocket.Conn)
	return c, ok
}

func UserConnMapStore(userID UUID, conn *websocket.Conn) {
	userConnRegister.uc.Store(userID, conn)
}

func UserConnMapDelete(userID UUID) {
	userConnRegister.uc.Delete(userID)
}
