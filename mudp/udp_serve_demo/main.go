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
	udpServer, err := mudp.NewServer(mudp.Server{
		Port:      PORT,
		IP:        IP,
		MultiCore: true,
		OnMessage: func(eventName string, data []byte) {
			fmt.Println(eventName, string(data))
		},
	})
	if err != nil {
		fmt.Println("服务创建失败", err)
	}

	err = udpServer.Start() // 阻塞启动
	if err != nil {
		fmt.Println("服务启动失败", err)
	}
}
