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
			err = fmt.Errorf("err:mudp.NewSender|ResolveUDPAddr: %w", err)
			return
		}
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		err = fmt.Errorf("err:mudp.NewSender|ListenUDP: %w", err)
		return
	}
	s.conn = conn
	return
}

// Send 将数据发送到 remoteAddr（"ip:port"）。支持通过 ctx 取消和可选的超时参数。
func (s *Sender) Send(ctx context.Context, remoteAddr string, data []byte, timeout time.Duration) (n int, err error) {
	if s == nil || s.conn == nil {
		err = errors.New("err:mudp.Send|nil sender")
		return
	}
	raddr, e := net.ResolveUDPAddr("udp", remoteAddr)
	if e != nil {
		err = fmt.Errorf("err:mudp.Send|resolve addr: %w", e)
		return
	}

	// 在锁内获取连接副本，然后在执行阻塞 I/O 前释放锁。
	// 这样可以避免在网络操作期间持有互斥锁，从而防止与 Close() 产生死锁。
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		err = errors.New("err:mudp.Send|conn nil")
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
		err = fmt.Errorf("err:mudp.Send|ctx.Done(): %w", ctx.Err())
		return
	case r := <-ch:
		if r.err != nil {
			err = fmt.Errorf("err:mudp.Send|WriteToUDP: %w", r.err)
			return
		}
		n = r.n
		return
	}
}

// Close 关闭发送者的连接。
func (s *Sender) Close() error {
	if s == nil {
		return fmt.Errorf("err:mudp.Close|Send|nil sender")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return fmt.Errorf("err:mudp.Close|Send|nil conn")
	}
	err := s.conn.Close()
	s.conn = nil
	return fmt.Errorf("err:mudp.Close|Send|close conn: %w", err)
}
