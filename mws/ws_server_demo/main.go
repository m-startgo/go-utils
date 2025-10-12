package main

import (
	"context"
	"fmt"
	"net/http"

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
			fmt.Println("err:mws.ws_server_demo|Accept", err)
			return
		}
		// Close with a normal closure when the handler exits
		defer func() { _ = c.Close(websocket.StatusNormalClosure, "bye") }()

		for {
			var v any
			err := wsjson.Read(context.Background(), c, &v)
			if err != nil {
				fmt.Println("read-err:", v)
				break
			}
			fmt.Println("read:", v)

			ack := map[string]any{"ok": true, "recv": v}
			err = wsjson.Write(context.Background(), c, ack)
			if err != nil {
				fmt.Println("err:mws.ws_server_demo|Write|", err)
				break
			}
		}
	})

	fmt.Println("ws server listening on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("err:mws.ws_server_demo|ListenAndServe", err)
	}
}
