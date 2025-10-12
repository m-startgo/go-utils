package main

import (
	"fmt"
	"log"
	"net"

	"github.com/m-startgo/go-utils/mudp"
)

func main() {
	addr := "udp://127.0.0.1:9999"
	srv := mudp.NewServer(addr, func(remoteAddr net.Addr, data []byte) {
		log.Printf("recv from %s: %s", remoteAddr.String(), string(data))
	})
	err := srv.Start()
	if err != nil {
		fmt.Println("server start failed:", err)
		return
	}
}
