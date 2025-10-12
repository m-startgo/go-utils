package main

import (
	"fmt"

	"github.com/m-startgo/go-utils/mstr"
)

var (
	PROT = 9999
	IP   = "127.0.0.1"
)

func main() {
	url := mstr.Join("udp://", IP, ":", PROT)

	fmt.Println("server exit:", url)
}
