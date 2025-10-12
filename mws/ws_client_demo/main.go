package main

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, resp, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		if resp != nil {
			log.Printf("err:mws.ws_client_demo|Dial|status=%s", resp.Status)
		}
		log.Fatalf("err:mws.ws_client_demo|Dial|%v", err)
	}
	defer c.CloseNow()

	// 发送一条简单消息
	msg := map[string]any{"msg": "hello from client"}
	err = wsjson.Write(ctx, c, msg)
	if err != nil {
		log.Fatalf("err:mws.ws_client_demo|Write|%v", err)
	}

	// 读取回复
	var v any
	// 使用单独的 context 以便给读取操作一个较短超时
	rctx, rcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer rcancel()
	err = wsjson.Read(rctx, c, &v)
	if err != nil {
		log.Fatalf("err:mws.ws_client_demo|Read|%v", err)
	}
}
