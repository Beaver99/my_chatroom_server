package protocol

import (
	"context"
	"errors"
	"nhooyr.io/websocket"
)

type Msg map[string]interface{}

func (msg Msg) readReply(ctx context.Context, conn *websocket.Conn, feild string) (interface{}, error) {
	if value, ok := msg[feild]; ok {
		return value, nil
	} else {
		SendMsg(ctx, conn, ReplyMessage{
			MessageType: "0",
			State:       false,
			Err:         "Malformed message! Your message lacks some feilds:" + feild,
		})
		return nil, errors.New("this client send malformed message without" + feild + "feild.")
	}
}
func (msg Msg) ReadReply(ctx context.Context, conn *websocket.Conn, feild string) (interface{}, error){
	return msg.readReply(ctx, conn, feild)
}
