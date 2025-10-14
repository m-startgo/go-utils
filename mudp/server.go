package mudp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// OnMessageFunc 消息回调函数: remoteAddr, data
type OnMessageFunc func(remote string, data []byte)

// Server 基于 UDP 的接收器
// 支持 Start/Stop, 并发安全
// 示例:
//  srv := mudp.NewServer(mudp.Server{IP: "0.0.0.0", Port: 9000, OnMessage: func(r string, d []byte){ fmt.Println(r, string(d)) }})
//  srv.Start()
//  defer srv.Stop()

type Server struct {
	IP        string
	Port      int
	BufferLen int           // 每次读取的缓冲区大小, 默认 65536
	Timeout   time.Duration // Read deadline, 0 表示无超时
	OnMessage OnMessageFunc

	ln     *net.UDPConn
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.Mutex
	closed bool
}

// NewServer 构造 Server
func NewServer(opt *Server) (s *Server, err error) {
	if opt == nil {
		err = fmt.Errorf("err:mudp.NewServer|opt 不能为空")
		return
	}

	s = &Server{}

	if opt.IP == "" {
		opt.IP = "0.0.0.0"
	}
	if opt.Port == 0 {
		err = fmt.Errorf("err:mudp.NewServer|Port 不能为空")
		return
	}
	if opt.BufferLen <= 0 {
		opt.BufferLen = 65536
	}
	if opt.OnMessage == nil {
		opt.OnMessage = func(remote string, data []byte) {}
	}

	s.IP = opt.IP
	s.Port = opt.Port
	s.BufferLen = opt.BufferLen
	s.Timeout = opt.Timeout
	s.OnMessage = opt.OnMessage

	return
}

// Start 启动服务器并在后台接收数据
func (s *Server) Start() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		err = fmt.Errorf("err:mudp.Server|already closed")
		return
	}
	if s.ln != nil {
		return nil // already started
	}

	addr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		err = fmt.Errorf("err:mudp.Server|ResolveUDPAddr|%w", err)
		return
	}

	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		err = fmt.Errorf("err:mudp.Server|ListenUDP|%w", err)
		return
	}

	s.ln = ln
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(1)
	go s.recvLoop()
	return
}

func (s *Server) recvLoop() {
	defer s.wg.Done()
	buf := make([]byte, s.BufferLen)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if s.Timeout > 0 {
				s.ln.SetReadDeadline(time.Now().Add(s.Timeout))
			}
			n, raddr, err := s.ln.ReadFromUDP(buf)
			if err != nil {
				ne, ok := err.(net.Error)
				if ok && ne.Timeout() {
					continue
				}
				// 如果因关闭导致的错误，退出
				if s.isClosed() {
					return
				}
				// 其它错误，回调并继续
				go s.OnMessage("", []byte(fmt.Sprintf("err:mudp.recvLoop|ReadFromUDP|%v", err)))
				continue
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			go s.OnMessage(raddr.String(), data)
		}
	}
}

func (s *Server) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// Stop 停止服务器并释放资源
func (s *Server) Stop() (err error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	ln := s.ln
	s.ln = nil
	cancel := s.cancel
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if ln != nil {
		ln.Close()
	}
	// 等待 goroutine 退出
	s.wg.Wait()
	return nil
}
