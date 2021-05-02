package server

import (
	"context"
	"log"
	"my_chatroom_server/logic"
	"my_chatroom_server/protocol"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func userHandle(ctx context.Context, conn *websocket.Conn, userID logic.UUID) {
	// tell the client that he/she has logged in
	protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
		MessageType: "0",
		State:       true,
		Err:         "",
	})

	// TODO: set deadline for userHandle

	ChatReq := make(map[string]interface{}, 32)
	SendMsg := make(map[string]interface{}, 8)

	for {

		//ChatReq := make(map[string]interface{})
		err := wsjson.Read(ctx, conn, &ChatReq)
		if err != nil {
			log.Println("read json error when reading ChatReq:", err)
			continue
		}

		chatReq := protocol.Msg(ChatReq)

		// check msgType
		msgTypeTemp, err := chatReq.ReadReply(ctx, conn, "MessageType")
		if err != nil {
			log.Println(err)
			continue
		}
		msgType := msgTypeTemp.(string)
		if msgType != "7" {
			protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
				MessageType: "0",
				State:       false,
				Err:         "Wrong mesage type! You are supposed to send 'ChatRequestMessage' type message.",
			})
			log.Println("this client is supposed to send 'ChatRequest' type message but send:", msgType, "type message")
			continue
		}

		// handle chatRequest msg

		// judge the mode, "0" for single chat, "1" for group chat
		modeTemp, err := chatReq.ReadReply(ctx, conn, "Mode")
		if err != nil {
			log.Println(err)
			continue
		}
		mode := modeTemp.(string)

		counterIDTemp, err := chatReq.ReadReply(ctx, conn, "ID")
		if err != nil {
			log.Println(err)
			continue
		}
		couterID := counterIDTemp.(string)

		if mode == "0" {
			userAccountDB := logic.GetUserAccountDB()
			isValidID, _ := userAccountDB.Exists(ctx, couterID).Result()
			if isValidID == 1 {
				// FIXME: delete this later
				protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
					MessageType: "0",
					State:       true,
					Err:         "",
				})
				// now we can constantly chat with the peer until the client switch to another peer or log out
				// TODO: offline chat
				// TODO: send file

				//if couterConn, ok := logic.UserConnMapLoad(logic.UUID(couterID)); ok {
				//	// read from the client and decide what to do.
				//	for {
				//		err := wsjson.Read(ctx, conn, &SendMsg)
				//		if err != nil {
				//			log.Println("read json error:", err)
				//			continue
				//		}
				//
				//		// check msgType
				//		msgTypeTemp, err := chatReq.ReadReply(ctx, conn, "MessageType")
				//		if err != nil {
				//			log.Println(err)
				//			continue
				//		}
				//		msgType := msgTypeTemp.(string)
				//		if msgType == "2" {
				//			wsjson.Write(ctx, couterConn, SendMsg)
				//		}else if msgType == "7"{
				//			// TODO: switch peer
				//			log.Println("the client switch to another peer, but this feature is not support yet.")
				//			continue
				//		}else if msgType == "8"{
				//			logic.UserConnMapDelete(userID)
				//			log.Println("the client: ", userID, "has logged out.")
				//			return
				//		} else {
				//			log.Println("this client is sending improper type message but send:", msgType, "type message")
				//			continue
				//		}
				//	}
				//
				//} else {
				//	log.Println("the counterpart has not logged yet.")
				//	continue
				//}

				couterConn, isLogin := logic.UserConnMapLoad(logic.UUID(couterID))
				// read from the client and decide what to do.
				for {
					err := wsjson.Read(ctx, conn, &SendMsg)
					if err != nil {
						log.Println("read json error:", err)
						continue
					}

					// check msgType
					msgTypeTemp, err := chatReq.ReadReply(ctx, conn, "MessageType")
					if err != nil {
						log.Println(err)
						continue
					}
					msgType := msgTypeTemp.(string)
					if msgType == "2" {
						if isLogin{
							wsjson.Write(ctx, couterConn, SendMsg)
						}else{
							// write to DB
							logic.OfflineMsgStoreCollection.Insert(logic.OfflineMsg{
								UserID: logic.UUID(couterID),
								Msg: SendMsg,
							})
						}
					}else if msgType == "7"{
						// TODO: switch peer
						log.Println("the client switch to another peer, but this feature is not support yet.")
						continue
					}else if msgType == "8"{
						logic.UserConnMapDelete(userID)
						log.Println("the client: ", userID, "has logged out.")
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
				continue
			}

		} else if mode == "1" {
			// TODO: group chat
			continue
		} else {
			log.Println("Illegal 'mode' value: ", mode)
			continue
		}
	}

	// 根据读取时的错误执行不同的 Close
	//if err == nil {
	//	conn.Close(websocket.StatusNormalClosure, "")
	//} else {
	//	log.Println("read from client error:", err)
	//	conn.Close(websocket.StatusInternalError, "Read from client error")
	//}
	//conn.Close(websocket.StatusNormalClosure, "")
}
