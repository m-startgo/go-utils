package mws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Conn 是 websocket 连接的简单封装，提供 Send / Listen / Close 接口，风格参考 mudp 包。
type Conn struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// NewClient 建立到 websocket 服务端的连接。requestHeader 可为 nil。
func NewClient(urlStr string, requestHeader http.Header) (c *Conn, err error) {
	ws, _, e := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if e != nil {
		return nil, fmt.Errorf("err:mws.NewClient|Dial: %w", e)
	}
	return &Conn{conn: ws}, nil
}

// Send 将 data 作为 text frame 发送到对端。返回写入的字节长度或错误。
// 支持通过 ctx 取消，以及可选的 timeout（写超时）。
func (c *Conn) Send(ctx context.Context, data []byte, timeout time.Duration) (n int, err error) {
	if c == nil || c.conn == nil {
		return 0, errors.New("err:mws.Send|nil conn")
	}
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return 0, errors.New("err:mws.Send|conn nil")
	}

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

	type res struct{ err error }
	ch := make(chan res, 1)
	go func() {
		ch <- res{err: conn.WriteMessage(websocket.TextMessage, data)}
	}()

	select {
	case <-ctx.Done():
		return 0, fmt.Errorf("err:mws.Send|ctx.Done(): %w", ctx.Err())
	case r := <-ch:
		if r.err != nil {
			return 0, fmt.Errorf("err:mws.Send|WriteMessage: %w", r.err)
		}
		return len(data), nil
	}
}

// Handler 用于处理收到的消息，messageType 与 websocket.MessageType 一致。
type Handler func(messageType int, data []byte)

// Listen 循环读取消息并交给 handler 处理。通过 ctx 取消监听。
// 内部使用短超时周期性检查 ctx，以便可以响应取消。
func (c *Conn) Listen(ctx context.Context, handler Handler) error {
	if c == nil || c.conn == nil {
		return errors.New("err:mws.Listen|nil conn")
	}
	if handler == nil {
		return errors.New("err:mws.Listen|nil handler")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// 设置短读超时以便可以周期性检查 ctx.Done()
			_ = c.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			mt, data, err := c.conn.ReadMessage()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					select {
					case <-ctx.Done():
						return nil
					default:
						continue
					}
				}
				return fmt.Errorf("err:mws.Listen|ReadMessage: %w", err)
			}
			// 复制数据，避免回调中修改底层切片
			d := make([]byte, len(data))
			copy(d, data)
			go handler(mt, d)
		}
	}
}

// Close 关闭 websocket 连接。
func (c *Conn) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.conn.Close()
	c.conn = nil
	if err != nil {
		return fmt.Errorf("err:mws.Close|close conn: %w", err)
	}
	return nil
}

// NewFromWS 用于将已升级的 *websocket.Conn 包装为 *Conn（服务器端使用）。
func NewFromWS(ws *websocket.Conn) *Conn {
	return &Conn{conn: ws}
}
