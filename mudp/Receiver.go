package mudp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// Receiver 在绑定的地址上接收 UDP 数据报并将其分发到 handler。
// 支持通过 ctx 取消，ctx 取消或发生错误时返回。
type Receiver struct {
	conn *net.UDPConn
}

// NewReceiver 将 UDP 监听器绑定到 localAddr（例如 ":9000" 或 "127.0.0.1:9000"）。
func NewReceiver(localAddr string) (r *Receiver, err error) {
	r = &Receiver{}
	laddr, e := net.ResolveUDPAddr("udp", localAddr)
	if e != nil {
		err = fmt.Errorf("err:mudp.NewReceiver|ResolveUDPAddr: %w", e)
		return
	}
	conn, e := net.ListenUDP("udp", laddr)
	if e != nil {
		err = fmt.Errorf("err:mudp.NewReceiver|ListenUDP: %w", e)
		return
	}
	r.conn = conn
	return
}

// Handler 在每个收到的数据报上被调用，addr 为发送者地址。
type Handler func(data []byte, addr *net.UDPAddr)

// Listen 开始接收消息并为每个数据包调用 handler。
// bufferSize 控制读取的最大数据包大小（如果 <=0 则默认为 4096）。
func (r *Receiver) Listen(ctx context.Context, handler Handler, bufferSize int) error {
	if r == nil || r.conn == nil {
		return errors.New("err:mudp.Listen|nil receiver")
	}
	if handler == nil {
		return errors.New("err:mudp.Listen|nil handler")
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
				return fmt.Errorf("err:mudp.Listen|ReadFromUDP: %w", err)
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
