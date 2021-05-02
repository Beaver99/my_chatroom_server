package server

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"log"
	"my_chatroom_server/logic"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func sendOfflineMsg(ctx context.Context, conn *websocket.Conn, userID logic.UUID){
	var msgs []map[string]interface{}
	var err error
	logic.OfflineMsgStoreCollection.Find(bson.M{"userID":string(userID)}).All(&msgs)
	// FIXME: ## send and delete!
	for i := range msgs{
		err = wsjson.Write(ctx, conn, msgs[i])
		if err != nil{
			log.Println(err)
		}
	}
}
