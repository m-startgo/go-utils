package main

import (
	"fmt"
	"log"
	"time"

	"github.com/m-startgo/go-utils/mstr"
	"github.com/m-startgo/go-utils/mws"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

func main() {
	// 注意：server 必须先启动
	wsUrl := mstr.Join("ws://", IP, ":", PORT)
	conn, _, err := mws.Dial(wsUrl)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	conn.SetOnMessage(func(mt int, data []byte) {
		fmt.Println("server reply:", string(data))
	})
	conn.SetOnClose(func(err error) {
		fmt.Println("closed:", err)
	})

	_ = conn.SendText("hello from client")

	time.Sleep(1 * time.Second)
	_ = conn.Close()
	// wait a bit to let close callbacks print
	time.Sleep(200 * time.Millisecond)
}
