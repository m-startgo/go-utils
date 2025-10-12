package main

import (
	"fmt"

	"github.com/m-startgo/go-utils/mstr"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

func main() {
	wsUrl := mstr.Join("ws://", IP, ":", PORT, "/ws")

	fmt.Println("wsUrl:", wsUrl)
}
