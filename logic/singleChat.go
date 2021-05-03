package logic

import (
	"context"
	"log"
	"my_chatroom_server/protocol"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)
// TODO: for loop compatiblity
func SingleChat(ctx context.Context, conn *websocket.Conn, couterID UUID, SendMsg map[string]interface{},
userID UUID){
		userAccountDB := GetUserAccountDB()
		isValidID, _ := userAccountDB.Exists(ctx, string(couterID)).Result()
		if isValidID == 1 {
			// FIXME: delete this later
			protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
				MessageType: "0",
				State:       true,
				Err:         "",
			})
			// now we can constantly chat with the peer until the client switch to another peer or log out
			// TODO: send file

			//couterConn, isLogin := UserConnMapLoad(UUID(couterID))
			// read from the client and decide what to do.
			for {
				err := wsjson.Read(ctx, conn, &SendMsg)
				if err != nil {
					log.Println("read json error:", err)
					continue
				}

				sendMsg := protocol.Msg(SendMsg)
				// check msgType
				msgTypeTemp, err := sendMsg.ReadReply(ctx, conn, "MessageType")
				if err != nil {
					log.Println(err)
					continue
				}
				msgType := msgTypeTemp.(string)
				if msgType == "2" || msgType == "3" || msgType == "4" || msgType == "5" || msgType == "6"{
					if 	couterConn, isLogin := UserConnMapLoad(UUID(couterID)); isLogin{
						wsjson.Write(ctx, couterConn, SendMsg)
					}else{
						// write to DB
						OfflineMsgStoreCollection.Insert(OfflineMsg{
							UserID: couterID,
							Msg: SendMsg,
						})
					}
				}else if msgType == "7"{
					// TODO: switch peer
					log.Println("the client switch to another peer, but this feature is not support yet.")
					return
				}else if msgType == "8"{
					// FIXME: invalid logout implementation
					UserConnMapDelete(userID)
					log.Println("the client: ", userID, " has logged out")
					return
				} else {
					log.Println("this client is sending improper type message but send:", msgType, "type message")
					continue
				}
			}


		} else {
			protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
				MessageType: "0",
				State:       false,
				Err:         "you were trying to chat with an entity with invalid user ID.",
			})
			log.Println("this chat request contains a invalid counterpart ID.")
			return
		}
}
