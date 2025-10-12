package mudp

import (
	"fmt"

	"github.com/m-startgo/go-utils/mstr"
	"github.com/panjf2000/gnet/v2"
)

type OnMessageFunc func(eventName string, data []byte)

type Server struct {
	IP        string
	Port      int
	MultiCore bool
	OnMessage OnMessageFunc
	addr      string
}

type echoServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	addr      string
	multiCore bool
	onMessage OnMessageFunc
}

func NewServer(opt Server) (server *Server, err error) {
	err = nil
	server = &Server{}

	server.IP = opt.IP
	if server.IP == "" {
		server.IP = "127.0.0.1"
	}
	server.Port = opt.Port
	if server.Port == 0 {
		err = fmt.Errorf("mudp.NewServer|Port 不能为空")
		return
	}

	if opt.addr == "" {
		server.addr = mstr.Join("udp://", server.IP, ":", server.Port)
	} else {
		server.addr = opt.addr
	}

	server.OnMessage = opt.OnMessage
	if server.OnMessage == nil {
		server.OnMessage = func(eventName string, data []byte) {
			// 默认空实现，避免 nil 指针异常
		}
	}

	server.MultiCore = opt.MultiCore

	return
}

func (c *Server) Start() error {
	echo := &echoServer{
		addr:      c.addr,
		multiCore: c.MultiCore,
		onMessage: c.OnMessage,
	}
	err := gnet.Run(echo, echo.addr, gnet.WithMulticore(c.MultiCore))
	return err
}

// 引擎启动准备好接收数据时
func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	go es.onMessage("OnBoot", []byte("server is ready")) // 异步调用避免阻塞
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, err := c.Next(-1)
	if err != nil {
		return gnet.Close
	}
	// 复制 buf 到新切片以避免 gnet 中的缓冲区被复用后产生数据竞争
	data := make([]byte, len(buf))
	copy(data, buf)
	go es.onMessage("OnTraffic", data) // 异步调用避免阻塞
	return gnet.None
}
