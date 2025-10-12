package main

import (
	"fmt"
	"log"
	"time"

	"github.com/m-startgo/go-utils/mws"
)

func main() {
	// 注意：server 必须先启动
	url := "ws://127.0.0.1:8080/ws"
	conn, _, err := mws.Dial(url)
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
