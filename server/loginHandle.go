package server

import (
	"github.com/go-redis/redis/v8"
	"log"
	"my_chatroom_server/protocol"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// /login
func loginHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 从客户端接受 WebSocket 握手，并将连接升级到 WebSocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: false})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	ctx := req.Context()

	LoginReq := make(map[string]interface{})
	err = wsjson.Read(ctx, conn, &LoginReq) // what if res??
	if err != nil {
		log.Println("read json error:", err)
		return
	}

	loginReq := protocol.Msg(LoginReq) // ## api as recv / arg??


	// check msgType
	msgTypeTemp, err := loginReq.ReadReply(ctx, conn, "MessageType")
	if err != nil{
		log.Println(err)
		return
	}
	msgType := msgTypeTemp.(string)
	if msgType != "1" {
		protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
			MessageType: "0",
			State:       false,
			Err:         "Wrong mesage type! You are supposed to send 'LoginMessage' type message.",
		})
		log.Println("this client is supposed to send 'login' type message but send:", msgType, "type message")
		return
	}


	// handle login request:
	// if username does not occur in DB, insert the username-psw pair into DB;
	// else check if the psw is right.
	usernameTemp, err := loginReq.ReadReply(ctx, conn, "Username")
	if err != nil{
		log.Println(err)
		return
	}
	username := usernameTemp.(string)


	pswTemp, err:= loginReq.ReadReply(ctx, conn, "Password")
	if err != nil{
		log.Println(err)
		return
	}
	psw := pswTemp.(string)

	password, err := userInfoDB.Get(ctx, username).Result()
	if err == redis.Nil{// user does not exists yet, create user-psw
		userInfoDB.Set(ctx, username, psw, 0)
	}else if err ==nil{//user exists, check its psw
		//log.Println("psw: ", psw, "password: ", password)
		if psw != password{
			protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
				MessageType: "0",
				State:       false,
				Err:         "Wrong Username-Password pair!",
			})
			log.Println("Wrong Username-Password pair!")
			return
		}
	}


	protocol.SendMsg(ctx, conn, protocol.ReplyMessage{
		MessageType: "0",
		State:       true,
		Err:         "",
	})


	//

	//// 1. 新用户进来，构建该用户的实例
	//token := req.FormValue("token")
	//nickname := req.FormValue("nickname")
	//if l := len(nickname); l < 2 || l > 20 {
	//	log.Println("nickname illegal: ", nickname)
	//	wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：2-20")) // ??
	//	conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
	//	return
	//}
	//if !logic.Broadcaster.CanEnterRoom(nickname) {
	//	log.Println("昵称已经存在：", nickname)
	//	wsjson.Write(req.Context(), conn, logic.NewErrorMessage("该昵称已经已存在！"))
	//	conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
	//	return
	//}
	//
	//userHasToken := logic.NewUser(conn, token, nickname, req.RemoteAddr) // ?? token

	//// 2. 开启给用户发送消息的 goroutine
	//go userHasToken.SendMessage(req.Context())
	//
	//// 3. 给当前用户发送欢迎信息
	//userHasToken.MessageChannel <- logic.NewWelcomeMessage(userHasToken)
	//
	//// 避免 token 泄露
	//tmpUser := *userHasToken
	//user := &tmpUser
	//user.Token = ""
	//
	//// 给所有用户告知新用户到来
	//msg := logic.NewUserEnterMessage(user)
	//logic.Broadcaster.Broadcast(msg)
	//
	//// 4. 将该用户加入广播器的用列表中
	//logic.Broadcaster.UserEntering(user)
	//log.Println("user:", nickname, "joins chat") // ?? log println
	//
	//// 5. 接收用户消息
	//err = user.ReceiveMessage(req.Context())
	//
	//// 6. 用户离开
	//logic.Broadcaster.UserLeaving(user)
	//msg = logic.NewUserLeaveMessage(user)
	//logic.Broadcaster.Broadcast(msg)
	//log.Println("user:", nickname, "leaves chat")

	// 根据读取时的错误执行不同的 Close
	//if err == nil {
	//	conn.Close(websocket.StatusNormalClosure, "")
	//} else {
	//	log.Println("read from client error:", err)
	//	conn.Close(websocket.StatusInternalError, "Read from client error")
	//}
	conn.Close(websocket.StatusNormalClosure, "")
}
