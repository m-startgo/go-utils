package mudp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Sender represents a UDP sender
type Sender struct {
	conn *net.UDPConn
	mu   sync.Mutex
}

// NewSender creates a UDP sender bound to an optional localAddr ("ip:port" or ":0").
// If localAddr is empty, the OS will pick the local address.
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

// Send sends data to remoteAddr ("ip:port"). It respects ctx for cancellation and an optional timeout.
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

	s.mu.Lock()
	// ensure conn hasn't been closed between the earlier nil-check and the lock
	if s.conn == nil {
		s.mu.Unlock()
		err = errors.New("err:utils.udp|Send|nil sender")
		return
	}

	// Merge provided timeout with context deadline (use the earlier one).
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
		_ = s.conn.SetWriteDeadline(dl)
	} else {
		_ = s.conn.SetWriteDeadline(time.Time{})
	}

	type res struct {
		n   int
		err error
	}
	ch := make(chan res, 1)
	go func() {
		ni, er := s.conn.WriteToUDP(data, raddr)
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

// Close closes the sender's connection.
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

// Receiver receives UDP datagrams on a bound address and dispatches them to a handler.
// It supports cancellation via ctx and returns when ctx is done or on error.
type Receiver struct {
	conn *net.UDPConn
}

// Handler is a callback invoked for each received datagram. addr is the sender address.
type Handler func(data []byte, addr *net.UDPAddr)

// Listen starts receiving messages and calls handler for each datagram.
// bufferSize controls the maximum datagram size to read (default 4096 if <=0).
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
				// timeout or temporary; check context and continue
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					// check ctx again
					select {
					case <-ctx.Done():
						return nil
					default:
						continue
					}
				}
				return fmt.Errorf("err:utils.udp|Listen|read: %w", err)
			}
			// copy the data slice for handler
			data := make([]byte, n)
			copy(data, buf[:n])
			go handler(data, addr)
		}
	}
}

// Close closes the receiver's connection.
func (r *Receiver) Close() error {
	if r == nil || r.conn == nil {
		return nil
	}
	return r.conn.Close()
}
