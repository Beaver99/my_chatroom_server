package logic

import (
	"context"
	"nhooyr.io/websocket/wsjson"
)

func deliveryMan(ctx context.Context, counterID UUID, msgType string, SendMsg map[string]interface{}){

	counterConn, isLogin := UserConnMapLoad(counterID)

	if msgType == "2" {
		if isLogin {
			wsjson.Write(ctx, counterConn, SendMsg)
		} else {
			// write to DB
			OfflineMsgStoreCollection.Insert(OfflineMsg{
				UserID: counterID,
				Msg:    SendMsg,
			})
		}
	}
}
