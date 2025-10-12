package main

import (
	"context"
	"fmt"
	"log"
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
			log.Printf("err:mws.ws_client_demo|Dial|status=%s", resp.Status)
		}
		log.Fatalf("err:mws.ws_client_demo|Dial|%v", err)
	}
	defer c.CloseNow()

	go func() {
		for {
			// 发送一条简单消息
			msg := map[string]any{"msg": "hello from client"}
			err = wsjson.Write(context.Background(), c, msg)
			if err != nil {
				log.Fatalf("err:mws.ws_client_demo|Write|%v", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		var v any
		err = wsjson.Read(context.Background(), c, &v)
		if err != nil {
			// 不要直接 Fatal，优雅处理连接关闭或 EOF 情况
			log.Printf("err:mws.ws_client_demo|Read|%v", err)
			break
		}
		fmt.Println("recv:", v)

	}
}
