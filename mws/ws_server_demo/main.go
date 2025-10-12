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
	wsUrl := mstr.Join(IP, ":", PORT)
	fmt.Println("wsUrl:", wsUrl)
}
