package main

import (
	"context"
	"fmt"
	"time"

	"github.com/m-startgo/go-utils/mws"
)

func main() {
	url := "ws://127.0.0.1:9999/echo"
	c, err := mws.NewClient(url, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.Listen(ctx, func(mt int, data []byte) {
			fmt.Println("client recv:", string(data))
		})
	}()

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("hello ws %d", i)
		_, err := c.Send(context.Background(), []byte(msg), time.Second)
		if err != nil {
			fmt.Println("send err:", err)
			return
		}
		time.Sleep(300 * time.Millisecond)
	}

	// 发送完成后取消 Listen 的上下文并等待其返回，避免主 goroutine 永久阻塞导致 runtime deadlock
	cancel()
	if err := <-errCh; err != nil {
		fmt.Println("listen exit with error:", err)
	} else {
		fmt.Println("listen exited")
	}
}
