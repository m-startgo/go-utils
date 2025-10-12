package main

import (
	"fmt"
	"time"

	"github.com/m-startgo/go-utils/mudp"
)

func main() {
	msg := []byte("hello from sender")
	if err := mudp.SendTo("127.0.0.1:9001", msg); err != nil {
		fmt.Println("send error:", err)
		return
	}
	fmt.Println("sent")
	// 等待一小会儿以便接收端打印
	time.Sleep(100 * time.Millisecond)
}
