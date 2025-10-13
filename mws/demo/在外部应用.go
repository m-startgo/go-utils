package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/m-startgo/go-utils/mstr"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

// ...existing code...

// 简单的 WebSocket echo/ack 服务示例，且支持在 handler 之外发送消息
var (
	clients   = make(map[*websocket.Conn]chan any)
	clientsMu sync.Mutex
)

// SendToAll 向所有已注册的连接广播消息（非阻塞，若通道满则丢弃该连接的这条消息）。
func SendToAll(msg any) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for _, ch := range clients {
		select {
		case ch <- msg:
		default:
			// 如果通道满，跳过，避免阻塞
		}
	}
}

// registerClient 注册连接并启动写入 goroutine
func registerClient(c *websocket.Conn) {
	sendCh := make(chan any, 16)

	clientsMu.Lock()
	clients[c] = sendCh
	clientsMu.Unlock()

	go func() {
		// 该 goroutine 专门负责把 sendCh 的消息写到 websocket
		for v := range sendCh {
			// 写超时可通过 context.WithTimeout 实现，此处使用背景上下文
			if err := wsjson.Write(context.Background(), c, v); err != nil {
				fmt.Println("err:mws.ws_server_demo|writer|", err)
				break
			}
		}
		// 退出时确保连接被关闭
		_ = c.Close(websocket.StatusNormalClosure, "bye")
	}()
}

// unregisterClient 注销连接并清理资源
func unregisterClient(c *websocket.Conn) {
	clientsMu.Lock()
	ch, ok := clients[c]
	if ok {
		delete(clients, c)
		close(ch)
	}
	clientsMu.Unlock()
}

// main 函数，演示如何在 handler 外部推送消息
func main() {
	addr := mstr.Join(IP, ":", PORT)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			fmt.Println("err:mws.ws_server_demo|Accept", err)
			return
		}

		// 注册连接（启动写 goroutine）
		registerClient(c)
		// 注销与关闭由读取循环/写循环决定
		defer unregisterClient(c)

		// Close with a normal closure when the handler exits
		defer func() { _ = c.Close(websocket.StatusNormalClosure, "bye") }()

		for {
			var v any
			err := wsjson.Read(context.Background(), c, &v)
			if err != nil {
				fmt.Println("err:mws.ws_server_demo|Read|", err)
				break
			}
			fmt.Println("read:", v)

			ack := map[string]any{"ok": true, "recv": v}
			// 直接通过写通道发送响应也可以（这里示例直接写）
			err = wsjson.Write(context.Background(), c, ack)
			if err != nil {
				fmt.Println("err:mws.ws_server_demo|Write|", err)
				break
			}
		}
	})

	// 示例：在 handler 外部每 5 秒向所有客户端广播一次消息
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		i := 0
		for range ticker.C {
			i++
			SendToAll(map[string]any{"push": fmt.Sprintf("server push %d", i)})
		}
	}()

	fmt.Println("ws server listening on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("err:mws.ws_server_demo|ListenAndServe", err)
		return
	}
}
