package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/m-startgo/go-utils/mstr"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

// 简单的 WebSocket echo/ack 服务示例
func main() {
	addr := mstr.Join(IP, ":", PORT)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Printf("err:mws.ws_server_demo|Accept|%v", err)
			http.Error(w, "websocket accept error", http.StatusBadRequest)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var v any
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Printf("err:mws.ws_server_demo|Read|%v", err)
			// 尝试发送错误状态并返回
			_ = c.Close(websocket.StatusInternalError, "read error")
			return
		}

		log.Printf("received: %v", v)

		// 回送一个确认消息
		ack := map[string]any{"ok": true, "recv": v}
		err = wsjson.Write(ctx, c, ack)
		if err != nil {
			log.Printf("err:mws.ws_server_demo|Write|%v", err)
			_ = c.Close(websocket.StatusInternalError, "write error")
			return
		}
		// 正常关闭连接
		_ = c.Close(websocket.StatusNormalClosure, "bye")
	})

	log.Printf("ws server listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("err:mws.ws_server_demo|ListenAndServe|%v", err)
	}
}
