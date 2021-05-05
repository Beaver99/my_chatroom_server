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
	logic.OfflineMsgStoreCollection.Find(bson.M{"userid":string(userID)}).All(&msgs)
	// FIXME: ## send and delete!
	for i := range msgs{
		//// TODO: this is ulgy and slow!
		err = wsjson.Write(ctx, conn, msgs[i]["msg"])
		if err != nil{
			log.Println(err)
		}
	}
	change, err := logic.OfflineMsgStoreCollection.RemoveAll(bson.M{"userid":string(userID)})
	if err != nil{
		log.Println(err)
		log.Println(change)
	}
	//context.WithCancel(ctx)
}

