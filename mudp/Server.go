package mudp

import (
	"github.com/m-startgo/go-utils/mstr"
	"github.com/panjf2000/gnet/v2"
)

type OnMessageFunc func(eventName string, data []byte)

type Server struct {
	Port      int
	IP        string
	MultiCore bool
	OnMessage OnMessageFunc
}

type echoServer struct {
	gnet.BuiltinEventEngine
	Eng       gnet.Engine
	Addr      string
	MultiCore bool
	OnMessage OnMessageFunc
}

func NewServer(opt Server) *Server {
	var c Server
	c.IP = opt.IP
	if c.IP == "" {
		c.IP = "127.0.0.1"
	}
	c.Port = opt.Port
	if c.Port == 0 {
		c.Port = 9000
	}

	c.OnMessage = opt.OnMessage
	if c.OnMessage == nil {
		c.OnMessage = func(eventName string, data []byte) {
			// 默认空实现，避免 nil 指针异常
		}
	}

	c.MultiCore = opt.MultiCore
	return &c
}

func (c *Server) Start() error {
	UDPAddr := mstr.Join("udp://", c.IP, ":", c.Port)

	echo := &echoServer{
		Addr:      UDPAddr,
		MultiCore: c.MultiCore,
		OnMessage: c.OnMessage,
	}

	err := gnet.Run(echo, echo.Addr, gnet.WithMulticore(c.MultiCore))

	return err
}

// 引擎启动准备好接收数据时
func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.Eng = eng
	go es.OnMessage("OnBoot", []byte("server is ready")) // 异步调用避免阻塞
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	go es.OnMessage("OnTraffic", buf) // 异步调用避免阻塞
	return gnet.None
}
