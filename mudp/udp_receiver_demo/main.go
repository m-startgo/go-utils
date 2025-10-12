package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m-startgo/go-utils/mudp"
)

func main() {
	srv := mudp.NewServer(":9001", func(addr net.Addr, data []byte) {
		fmt.Printf("recv from %s: %s\n", addr.String(), string(data))
	})

	go func() {
		if err := srv.Start(); err != nil {
			fmt.Println("server exit:", err)
		}
	}()

	// 等待 Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("shutting down")
	_ = srv.Stop()
	// 给 gnet 一点时间清理
	time.Sleep(200 * time.Millisecond)
}
