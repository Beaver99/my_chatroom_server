package protocol

import (
	"context"
	"fmt"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type LoginMessage struct {
	MessageType string // 1
	Username    string
	Password    string
	Mode        string // indicate one-to-one chatroom 0 or group chatroom 1
	ID          string // uid of the peer that this user want to connect to
}

type SendMessage struct {
	MessageType string // 2
	Sendername string
	Message     string
}

type ReplyMessage struct{
	MessageType string // 0
	State bool // true for success
	Err string
}

func SendMsg(ctx context.Context, c *websocket.Conn, sendData interface{}) {
	err := wsjson.Write(ctx, c, sendData)
	if err != nil {
		log.Println(err)
	}
}

// ## ??
func recvMsg(ctx context.Context, c *websocket.Conn) map[string]interface{} {
	var v interface{}
	err := wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	msg := v.(map[string]interface{})
	return msg
}

// ## ??
func printMsg(msg map[string]interface{}) {
	fmt.Println(msg["Username"])
	fmt.Println(": ")
	fmt.Println(msg["Message"])
}

