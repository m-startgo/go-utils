package mudp

import (
	"net"

	"github.com/panjf2000/gnet/v2"
)

// OnMessageFunc 在服务器接收到 UDP 包时被调用。
// addr: 发送方的远程地址；data: 接收到的负载的拷贝。
type OnMessageFunc func(addr net.Addr, data []byte)

// Server 是基于 gnet 的简单封装，用于运行 UDP 服务。
// 它提供 Start/Stop 生命周期管理和每个数据包的回调。
type Server struct {
	addr      string
	onMessage OnMessageFunc
	opts      []gnet.Option
	started   bool
	engine    gnet.Engine
}

// gnetHandler 实现了 gnet.EventHandler，在 OnTraffic 中读取到达的 UDP 包
// 并转发给用户提供的回调。
type gnetHandler struct {
	gnet.BuiltinEventEngine
	onMessage OnMessageFunc
	srv       *Server
}

// OnTraffic 在 gnet 的连接（UDP 套接字）上有数据到达时被调用。
func (h *gnetHandler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	if h.onMessage == nil {
		return gnet.None
	}
	n := c.InboundBuffered()
	if n <= 0 {
		return gnet.None
	}
	buf, err := c.Next(n)
	if err != nil {
		return gnet.None
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	// 异步调用以避免阻塞事件循环
	go h.onMessage(c.RemoteAddr(), data)
	return gnet.None
}

// OnBoot 在引擎启动时被调用，保存 Engine 到 Server，以便 Stop 时调用。
func (h *gnetHandler) OnBoot(eng gnet.Engine) (action gnet.Action) {
	if h.srv != nil {
		h.srv.engine = eng
	}
	return gnet.None
}

// NewServer 创建一个在 addr 上监听的 Server（例如 "udp://:9001" 或 "udp://127.0.0.1:9001"）。
// 为每个接收到的数据包调用传入的回调。
func NewServer(addr string, onMessage OnMessageFunc, opts ...gnet.Option) *Server {
	return &Server{addr: addr, onMessage: onMessage, opts: opts}
}

// Start 在 goroutine 中启动 gnet 引擎。该方法快速返回，引擎在后台运行。
// 调用 Stop() 停止服务器。
func (s *Server) Start() error {
	if s.started {
		return nil
	}
	handler := &gnetHandler{onMessage: s.onMessage, srv: s}
	s.started = true
	return gnet.Run(handler, s.addr, s.opts...)
}

// Stop 关闭为配置地址提供服务的引擎。
// 默认使用 3 秒的上下文超时。
func (s *Server) Stop() error {
	s.engine.Stop()
	return nil
}
