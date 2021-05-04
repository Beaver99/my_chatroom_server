package logic

import (
	"context"
	"log"
	"my_chatroom_server/protocol"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var chFile = make(chan map[string]interface{}, 64)

// TODO: for loop compatiblity
func SingleChat(ctx context.Context, conn *websocket.Conn, counterID UUID, SendMsg map[string]interface{},
	userID UUID) {
	userAccountDB := GetUserAccountDB()
	isValidID, _ := userAccountDB.Exists(ctx, string(counterID)).Result()
	if isValidID == 1 {
		// FIXME: delete this later
		protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
			MessageType: "0",
			State:       true,
			Err:         "",
		})
		// now we can constantly chat with the peer until the client switch to another peer or log out
		// TODO: send file

		//couterConn, isLogin := UserConnMapLoad(UUID(counterID))
		// read from the client and decide what to do.
		for {
			// TOFIGUREOUT: why this works
			SendMsg = nil
			err := wsjson.Read(ctx, conn, &SendMsg)
			if err != nil {
				log.Println("read json error:", err)
				continue
			}
			// TOFIGUREOUT: why add this SendMsgtemp
			SendMsgtemp := SendMsg
			sendMsg := protocol.Msg(SendMsg)
			// check msgType
			msgTypeTemp, err := sendMsg.ReadReply(ctx, conn, "MessageType")
			if err != nil {
				log.Println(err)
				continue
			}
			msgType := msgTypeTemp.(string)
			counterConn, isLogin := UserConnMapLoad(UUID(counterID));

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
			} else if msgType == "4" {
				SendFile(ctx, conn, SendMsgtemp)
			} else if msgType == "6" {
				if SendMsgtemp["Offset"].(string) == "0" {
					// FIXME: offline?
					counterConn, _ := UserConnMapLoad(UUID(counterID))
					go RecvSeg(ctx, counterConn, SendMsgtemp)
				} else {
					chFile <- SendMsgtemp
				}
			} else if msgType == "7" {
				// TODO: switch peer
				log.Println("the client switch to another peer, but this feature is not support yet.")
				return
			} else if msgType == "8" {
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
