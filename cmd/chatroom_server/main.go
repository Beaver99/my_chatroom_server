package main

import (
	"fmt"
	"log"
	"net/http"

	//"my_chatroom_server/server"
	"my_chatroom_server/server"
)

var (
	addr   = "0.0.0.0:20229"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s

`
)

func main() {
	fmt.Println(banner, addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
