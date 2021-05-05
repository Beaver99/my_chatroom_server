package server

import (
	"context"
	"github.com/GGBooy/message"
	"log"
	"my_chatroom_server/logic"
	"my_chatroom_server/protocol"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func userHandle(ctx context.Context, conn *websocket.Conn, userID logic.UUID) {
	// tell the client that he/she has logged in
	protocol.SendMsg(ctx, conn, message.ReplyMessage{
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
			protocol.SendMsg(ctx, conn, message.ReplyMessage{
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
		// TODO: UUID system
		couterIDTemp := counterIDTemp.(string)
		couterID := logic.UUID(couterIDTemp)
		if mode == "0" {
			logic.SingleChat(ctx, conn, couterID, SendMsg, userID)
			return

		} else if mode == "1" {
			// TODO: group chat
			logic.GroupChat(ctx, conn, couterID, SendMsg, userID)
			return
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
