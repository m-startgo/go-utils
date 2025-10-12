package main

import (
	"fmt"
	"log"

	"github.com/m-startgo/go-utils/mstr"
	"github.com/panjf2000/gnet/v2"
)

var (
	PORT   = 9999
	IPAddr = "127.0.0.1"
)

type echoServer struct {
	gnet.BuiltinEventEngine
	Eng       gnet.Engine
	Addr      string
	Multicore bool
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.Eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", es.Multicore, es.Addr)
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}

func main() {
	multicore := true

	UDPAddr := mstr.Join("tcp://", IPAddr, ":", PORT)

	echo := &echoServer{Addr: UDPAddr, Multicore: multicore}

	err := gnet.Run(echo, echo.Addr, gnet.WithMulticore(multicore))
	if err != nil {
		fmt.Println("服务启动失败:", err)
	}
}
