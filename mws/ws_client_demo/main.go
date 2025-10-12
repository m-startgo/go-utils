package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/m-startgo/go-utils/mstr"
	"github.com/m-startgo/go-utils/mtime"
	"github.com/m-startgo/go-utils/mws"
)

const (
	port = 9999
	IP   = "127.0.0.1"
)

func main() {
	url := mstr.Join("ws://", IP, ":", port, "/ws")

	// 拨号到 ws 服务器
	conn, _, err := mws.DialContext(context.Background(), url, http.Header{})
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	log.Printf("已连接到服务器：%s", url)

	var i int

	for {
		i++
		msg := []byte(mstr.Join(mtime.NowDefaultString(), "-消息 ", i))
		if err := conn.WriteMessage(1, msg); err != nil {
			log.Printf("发送错误: %v", err)
			return
		}
		log.Printf("已发送: %s", string(msg))

		mt, rmsg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取回应错误: %v", err)
			return
		}
		log.Printf("收到回应(%d): %s", mt, string(rmsg))

		time.Sleep(2 * time.Second)
	}
}
