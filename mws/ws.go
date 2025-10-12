package mws

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

// Conn 是对 websocket.Conn 的轻量封装，提供并发安全的写操作、消息回调与安全关闭。
// 简单契约：
// - Dial  创建客户端连接并返回 *Conn
// - ServeUpgrade 在 fasthttp 上处理升级并在连接建立后回调
// - Conn 支持 SetOnMessage/SetOnClose 注册回调，SendText/SendBinary 发送消息，Close 关闭连接

type Conn struct {
	ws        *websocket.Conn
	writeMu   sync.Mutex
	closed    int32
	onMessage func(messageType int, data []byte)
	onClose   func(err error)
}

// NewConn 包装 websocket.Conn 并启动读取协程
func NewConn(ws *websocket.Conn) *Conn {
	c := &Conn{ws: ws}
	go c.readLoop()
	return c
}

// Dial 建立到服务器的 websocket 连接，返回封装后的 *Conn
func Dial(urlStr string) (c *Conn, resp *http.Response, err error) {
	var d websocket.Dialer
	ws, r, e := d.Dial(urlStr, nil)
	if e != nil {
		err = e
		resp = r
		return
	}
	c = NewConn(ws)
	resp = r
	return
}

// ServeUpgrade 在 fasthttp 上处理 websocket 升级。onConnect 回调在连接建立后被调用。
func ServeUpgrade(ctx *fasthttp.RequestCtx, onConnect func(*Conn)) error {
	upgrader := websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true },
	}
	return upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		c := NewConn(ws)
		if onConnect != nil {
			onConnect(c)
		}
	})
}

// SetOnMessage 注册消息回调
func (c *Conn) SetOnMessage(fn func(messageType int, data []byte)) {
	c.onMessage = fn
}

// SetOnClose 注册连接关闭回调
func (c *Conn) SetOnClose(fn func(err error)) {
	c.onClose = fn
}

// SendText 发送文本消息（并发安全）
func (c *Conn) SendText(msg string) error {
	return c.writeMessage(websocket.TextMessage, []byte(msg))
}

// SendBinary 发送二进制消息（并发安全）
func (c *Conn) SendBinary(b []byte) error {
	return c.writeMessage(websocket.BinaryMessage, b)
}

func (c *Conn) writeMessage(mt int, data []byte) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return websocket.ErrCloseSent
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.ws.WriteMessage(mt, data)
}

// Close 安全关闭连接并触发 onClose 回调
func (c *Conn) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}
	err := c.ws.Close()
	if c.onClose != nil {
		c.onClose(err)
	}
	return err
}

func (c *Conn) readLoop() {
	for {
		mt, msg, err := c.ws.ReadMessage()
		if err != nil {
			// trigger close callback
			if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
				if c.onClose != nil {
					c.onClose(err)
				}
			}
			return
		}
		if c.onMessage != nil {
			c.onMessage(mt, msg)
		}
	}
}
