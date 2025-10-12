package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/m-startgo/go-utils/mstr"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

func main() {
	addr := mstr.Join("ws://", IP, ":", PORT, "/ws")

	c, resp, err := websocket.Dial(context.Background(), addr, nil)
	if err != nil {
		if resp != nil {
			fmt.Println("err:mws.ws_client_demo|Dial|status", resp.Status)
			return
		}
		fmt.Println("err:mws.ws_client_demo|Dial|", err)
		return
	}
	defer c.CloseNow()

	go func() {
		for {
			// 发送一条简单消息
			msg := map[string]any{"msg": "hello from client"}
			err := wsjson.Write(context.Background(), c, msg)
			if err != nil {
				fmt.Println("err:mws.ws_client_demo|Write|", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		var v any
		err := wsjson.Read(context.Background(), c, &v)
		if err != nil {
			fmt.Println("err:mws.ws_server_demo|Read|", err)
			break
		}
		fmt.Println("recv:", v)
	}
}
