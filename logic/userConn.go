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
	if ok{
		return ctemp.(*websocket.Conn), ok
	}else{
		return nil, ok
	}
}

func UserConnMapStore(userID UUID, conn *websocket.Conn) {
	userConnRegister.uc.Store(userID, conn)
}

func UserConnMapDelete(userID UUID) {
	userConnRegister.uc.Delete(userID)
}
