package main

import (
	"fmt"

	"github.com/m-startgo/go-utils/mudp"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

func main() {
	server := mudp.NewServer(mudp.Server{
		Port:      PORT,
		IP:        IP,
		MultiCore: true,
		OnMessage: func(eventName string, data []byte) {
			fmt.Println(eventName, string(data))
		},
	})

	err := server.Start()
	if err != nil {
		fmt.Println("服务启动失败", err)
	}
}
