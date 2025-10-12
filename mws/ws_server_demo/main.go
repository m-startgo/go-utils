package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-startgo/go-utils/mstr"
	"github.com/m-startgo/go-utils/mtime"
	"github.com/m-startgo/go-utils/mws"
)

const (
	port = 9999
	IP   = "127.0.0.1"
)

func main() {
	// 监听 /ws
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := mws.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "升级失败", http.StatusBadRequest)
			return
		}
		defer c.Close()

		log.Printf("客户端已连接：%s", r.RemoteAddr)
		go func() {
			for {
				mt, msg, err := c.ReadMessage()
				if err != nil {
					log.Printf("读取错误: %v", err)
					break
				}

				if mt == websocket.TextMessage {
					log.Printf("接收到了文本：%s", string(msg))
				} else if mt == websocket.BinaryMessage {
					log.Printf("接收到了二进制数据（%d 字节）", len(msg))
				} else {
					log.Printf("接收到了消息类型 %d（%d 字节）", mt, len(msg))
				}
				// 立即回应客户端
				if mt > 0 {
					sendMsg := []byte(mstr.Join(mtime.NowDefaultString(), "收到消息: ", msg))
					if err := c.WriteMessage(websocket.TextMessage, sendMsg); err != nil {
						log.Printf("写入错误: %v", err)
						break
					}
				}

				fmt.Printf("----\n")
			}
		}()
		// 主动给客户端发消息
		for {
			sendMsg := []byte(mstr.Join(mtime.NowDefaultString(), "服务端发送"))
			c.WriteMessage(1, []byte(sendMsg))
			time.Sleep(4 * time.Second)
		}
	})

	// 启动服务器
	addr := mstr.Join(IP, ":", port)
	log.Printf("ws 服务器正在监听 %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
