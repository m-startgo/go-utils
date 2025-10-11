package mudp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Sender 表示一个 UDP 发送者
type Sender struct {
	conn *net.UDPConn
	mu   sync.Mutex
}

// NewSender 创建一个 UDP 发送者，可选绑定本地地址（"ip:port" 或 ":0"）。
// 如果 localAddr 为空，操作系统会选择本地地址。
func NewSender(localAddr string) (s *Sender, err error) {
	s = &Sender{}
	var laddr *net.UDPAddr
	if localAddr != "" {
		laddr, err = net.ResolveUDPAddr("udp", localAddr)
		if err != nil {
			err = fmt.Errorf("err:utils.udp|NewSender|resolve local addr: %w", err)
			return
		}
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		err = fmt.Errorf("err:utils.udp|NewSender|listen udp: %w", err)
		return
	}
	s.conn = conn
	return
}

// Send 将数据发送到 remoteAddr（"ip:port"）。支持通过 ctx 取消和可选的超时参数。
func (s *Sender) Send(ctx context.Context, remoteAddr string, data []byte, timeout time.Duration) (n int, err error) {
	if s == nil || s.conn == nil {
		err = errors.New("err:utils.udp|Send|nil sender")
		return
	}
	raddr, e := net.ResolveUDPAddr("udp", remoteAddr)
	if e != nil {
		err = fmt.Errorf("err:utils.udp|Send|resolve remote addr: %w", e)
		return
	}

	// 在锁内获取连接副本，然后在执行阻塞 I/O 前释放锁。
	// 这样可以避免在网络操作期间持有互斥锁，从而防止与 Close() 产生死锁。
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		err = errors.New("err:utils.udp|Send|nil sender")
		return
	}

	// 将提供的 timeout 与 ctx 的 deadline 合并（以较早的为准）。
	var dl time.Time
	if timeout > 0 {
		dl = time.Now().Add(timeout)
	}
	if ctxDead, ok := ctx.Deadline(); ok {
		if dl.IsZero() || ctxDead.Before(dl) {
			dl = ctxDead
		}
	}
	if !dl.IsZero() {
		_ = conn.SetWriteDeadline(dl)
	} else {
		_ = conn.SetWriteDeadline(time.Time{})
	}

	type res struct {
		n   int
		err error
	}
	ch := make(chan res, 1)
	go func() {
		ni, er := conn.WriteToUDP(data, raddr)
		ch <- res{n: ni, err: er}
	}()

	select {
	case <-ctx.Done():
		err = fmt.Errorf("err:utils.udp|Send|ctx canceled: %w", ctx.Err())
		return
	case r := <-ch:
		if r.err != nil {
			err = fmt.Errorf("err:utils.udp|Send|write: %w", r.err)
			return
		}
		n = r.n
		return
	}
}

// Close 关闭发送者的连接。
func (s *Sender) Close() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return nil
	}
	err := s.conn.Close()
	s.conn = nil
	return err
}

// Receiver 在绑定的地址上接收 UDP 数据报并将其分发到 handler。
// 支持通过 ctx 取消，ctx 取消或发生错误时返回。
type Receiver struct {
	conn *net.UDPConn
}

// Handler 在每个收到的数据报上被调用，addr 为发送者地址。
type Handler func(data []byte, addr *net.UDPAddr)

// Listen 开始接收消息并为每个数据包调用 handler。
// bufferSize 控制读取的最大数据包大小（如果 <=0 则默认为 4096）。
func (r *Receiver) Listen(ctx context.Context, handler Handler, bufferSize int) error {
	if r == nil || r.conn == nil {
		return errors.New("err:utils.udp|Listen|nil receiver")
	}
	if handler == nil {
		return errors.New("err:utils.udp|Listen|nil handler")
	}
	if bufferSize <= 0 {
		bufferSize = 4096
	}

	buf := make([]byte, bufferSize)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_ = r.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			n, addr, err := r.conn.ReadFromUDP(buf)
			if err != nil {
				// 超时或临时错误；检查 ctx 并继续
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					// 再次检查 ctx
					select {
					case <-ctx.Done():
						return nil
					default:
						continue
					}
				}
				return fmt.Errorf("err:utils.udp|Listen|read: %w", err)
			}
			// 为 handler 拷贝一份数据切片
			data := make([]byte, n)
			copy(data, buf[:n])
			go handler(data, addr)
		}
	}
}

// Close 关闭接收者的连接。
func (r *Receiver) Close() error {
	if r == nil || r.conn == nil {
		return nil
	}
	return r.conn.Close()
}
