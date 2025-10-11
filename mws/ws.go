package mws

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Conn 是对 websocket.Conn 的薄封装，提供
// JSON/文本 的读写、简单的超时辅助和安全关闭。
type Conn struct {
	ws           *websocket.Conn
	sendDeadline time.Duration
	readDeadline time.Duration
}

// NewConn 用来封装一个已存在的 *websocket.Conn
func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{ws: ws, sendDeadline: 10 * time.Second, readDeadline: 0}
}

// SetDeadlines 设置发送和读取操作的默认截止时间。
func (c *Conn) SetDeadlines(send, read time.Duration) {
	c.sendDeadline = send
	c.readDeadline = read
}

// WriteJSON 使用配置的发送截止时间将 v 写为 JSON。
func (c *Conn) WriteJSON(v any) error {
	if c.sendDeadline > 0 {
		_ = c.ws.SetWriteDeadline(time.Now().Add(c.sendDeadline))
	}
	return c.ws.WriteJSON(v)
}

// ReadJSON 将下一个 JSON 消息读取到 v 中。如果设置了读取截止时间，则会应用该截止时间。
func (c *Conn) ReadJSON(v any) error {
	if c.readDeadline > 0 {
		_ = c.ws.SetReadDeadline(time.Now().Add(c.readDeadline))
	} else {
		// Clear deadline
		_ = c.ws.SetReadDeadline(time.Time{})
	}
	return c.ws.ReadJSON(v)
}

// WriteMessage 写入一个文本或二进制消息。msgType 应为 websocket.TextMessage 或 websocket.BinaryMessage。
/*
websocket.TextMessage （值通常为 1）—— 文本消息（UTF-8 编码）
websocket.BinaryMessage （值通常为 2）—— 二进制消息
websocket.CloseMessage （控制帧，值通常为 8）
websocket.PingMessage、websocket.PongMessage（控制帧）
*/
func (c *Conn) WriteMessage(msgType int, data []byte) error {
	if c.sendDeadline > 0 {
		_ = c.ws.SetWriteDeadline(time.Now().Add(c.sendDeadline))
	}
	return c.ws.WriteMessage(msgType, data)
}

// ReadMessage 从连接中读取单个消息。
func (c *Conn) ReadMessage() (int, []byte, error) {
	if c.readDeadline > 0 {
		_ = c.ws.SetReadDeadline(time.Now().Add(c.readDeadline))
	} else {
		_ = c.ws.SetReadDeadline(time.Time{})
	}
	return c.ws.ReadMessage()
}

// Close 关闭底层连接。
func (c *Conn) Close() error {
	return c.ws.Close()
}

// DialContext 连接（拨号）到一个 websocket 服务器并返回封装后的 Conn。
func DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*Conn, *http.Response, error) {
	d := websocket.DefaultDialer
	// respect context by setting net.Dialer in the dialer if provided via ctx (left default here)
	ws, resp, err := d.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		return nil, resp, err
	}
	return NewConn(ws), resp, nil
}

// Upgrader 使用合理的默认值封装了 websocket.Upgrader。
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// default: allow all origins. Caller can override.
		return true
	},
}

// Upgrade 将 HTTP 请求升级为 websocket 连接并返回封装后的 Conn。
func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	ws, err := Upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return NewConn(ws), nil
}

// SimpleEchoHandler 返回一个 http.Handler，会将请求升级为 websocket 并将文本消息回显（作为相同消息）。
func SimpleEchoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade failed", http.StatusBadRequest)
			return
		}
		defer c.Close()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				return
			}
			// echo
			if mt == websocket.TextMessage {
				_ = c.WriteMessage(websocket.TextMessage, msg)
			} else {
				_ = c.WriteMessage(mt, msg)
			}
		}
	})
}

// ErrInvalidURL 在提供的 URL 为空时返回。
var ErrInvalidURL = errors.New("invalid websocket url")
