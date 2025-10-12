package mudp

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/panjf2000/gnet/v2"
)

// HandlerFunc 接收来自远端 addr 的数据包
type HandlerFunc func(addr net.Addr, data []byte)

// Server 是基于 gnet 的 UDP 服务器封装
type Server struct {
	protoAddr string
	opts      []gnet.Option
	handler   HandlerFunc

	wg   sync.WaitGroup
	once sync.Once
	rerr error
}

// NewServer 创建一个新的 UDP Server。addr 支持带或不带协议前缀，
// 例如 ":9000" 或 "udp://:9000"。
func NewServer(addr string, handler HandlerFunc, opts ...gnet.Option) *Server {
	proto := addr
	if !strings.Contains(addr, "://") {
		proto = "udp://" + addr
	}
	return &Server{
		protoAddr: proto,
		opts:      opts,
		handler:   handler,
	}
}

// Start 阻塞方式启动服务器（直接调用 gnet.Run，会阻塞当前 goroutine）
func (s *Server) Start() error {
	eh := &udpEventHandler{handler: s.handler}
	s.rerr = gnet.Run(eh, s.protoAddr, s.opts...)
	return s.rerr
}

// StartAsync 在独立 goroutine 启动服务器（非阻塞）
func (s *Server) StartAsync() {
	s.once.Do(func() {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			eh := &udpEventHandler{handler: s.handler}
			s.rerr = gnet.Run(eh, s.protoAddr, s.opts...)
		}()
	})
}

// Stop 停止服务器。注意：gnet 包提供的 Stop 已被标记为 DEPRECATED，但
// 在这里用于优雅关闭监听地址（如果失败会返回错误）。
func (s *Server) Stop() error {
	// 使用 context.Background，这里不携带额外的上下文。
	err := gnet.Stop(context.Background(), s.protoAddr)
	s.wg.Wait()
	return err
}

// udpEventHandler 将 gnet 事件映射到用户提供的 HandlerFunc
type udpEventHandler struct {
	gnet.BuiltinEventEngine
	handler HandlerFunc
}

// OnTraffic 在收到数据时调用，读取可用数据并交给用户回调处理
func (h *udpEventHandler) OnTraffic(c gnet.Conn) gnet.Action {
	if h.handler == nil {
		return gnet.None
	}
	// 为 UDP 准备一个足够大的缓冲区
	buf := make([]byte, 65536)
	n, err := c.Read(buf)
	if err != nil && err != io.EOF {
		return gnet.None
	}
	if n > 0 {
		// 复制数据传递给回调，避免后续复用缓冲区造成数据竞争
		data := make([]byte, n)
		copy(data, buf[:n])
		// 注意：对于 UDP，Conn.RemoteAddr() 表示当前包的远端地址
		h.handler(c.RemoteAddr(), data)
	}
	return gnet.None
}
