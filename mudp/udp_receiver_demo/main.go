package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/m-startgo/go-utils/mudp"
)

func main() {
	// 在本机 9000 端口监听
	addr := "127.0.0.1:9000"
	r, err := mudp.NewReceiver(addr)
	if err != nil {
		err = fmt.Errorf("udp 监听失败: %w", err)
		panic(err)
	}
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动接收器
	go func() {
		fmt.Println("监听已启动", addr)
		err := r.Listen(ctx, func(data []byte, addr *net.UDPAddr) {
			timeNow := time.Now().UnixNano()

			fmt.Println("receiver", addr.String(), string(data), timeNow)
		}, 4096)
		if err != nil {
			fmt.Println("receiver error:", err)
		}
	}()

	select {}
}
