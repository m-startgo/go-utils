package mws

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Conn is a thin wrapper around websocket.Conn that provides
// JSON/text read/write with simple timeout helpers and safe close.
type Conn struct {
	ws           *websocket.Conn
	sendDeadline time.Duration
	readDeadline time.Duration
}

// NewConn wraps an existing *websocket.Conn
func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{ws: ws, sendDeadline: 10 * time.Second, readDeadline: 0}
}

// SetDeadlines sets the default deadlines for send and read operations.
func (c *Conn) SetDeadlines(send, read time.Duration) {
	c.sendDeadline = send
	c.readDeadline = read
}

// WriteJSON writes v as JSON with the configured send deadline.
func (c *Conn) WriteJSON(v any) error {
	if c.sendDeadline > 0 {
		_ = c.ws.SetWriteDeadline(time.Now().Add(c.sendDeadline))
	}
	return c.ws.WriteJSON(v)
}

// ReadJSON reads the next JSON message into v. If a read deadline is set it will be applied.
func (c *Conn) ReadJSON(v any) error {
	if c.readDeadline > 0 {
		_ = c.ws.SetReadDeadline(time.Now().Add(c.readDeadline))
	} else {
		// Clear deadline
		_ = c.ws.SetReadDeadline(time.Time{})
	}
	return c.ws.ReadJSON(v)
}

// WriteMessage writes a text or binary message. msgType should be websocket.TextMessage or websocket.BinaryMessage.
func (c *Conn) WriteMessage(msgType int, data []byte) error {
	if c.sendDeadline > 0 {
		_ = c.ws.SetWriteDeadline(time.Now().Add(c.sendDeadline))
	}
	return c.ws.WriteMessage(msgType, data)
}

// ReadMessage reads a single message from the connection.
func (c *Conn) ReadMessage() (int, []byte, error) {
	if c.readDeadline > 0 {
		_ = c.ws.SetReadDeadline(time.Now().Add(c.readDeadline))
	} else {
		_ = c.ws.SetReadDeadline(time.Time{})
	}
	return c.ws.ReadMessage()
}

// Close closes the underlying connection.
func (c *Conn) Close() error {
	return c.ws.Close()
}

// DialContext dials a websocket server and returns a wrapped Conn.
func DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*Conn, *http.Response, error) {
	d := websocket.DefaultDialer
	// respect context by setting net.Dialer in the dialer if provided via ctx (left default here)
	ws, resp, err := d.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		return nil, resp, err
	}
	return NewConn(ws), resp, nil
}

// Upgrader wraps websocket.Upgrader with sane defaults.
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// default: allow all origins. Caller can override.
		return true
	},
}

// Upgrade upgrades an HTTP request to a websocket connection and returns a wrapped Conn.
func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	ws, err := Upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return NewConn(ws), nil
}

// SimpleEchoHandler returns an http.Handler that upgrades and echos back text messages as JSON {msg:...}.
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

// ErrInvalidURL is returned when the provided URL is empty.
var ErrInvalidURL = errors.New("invalid websocket url")
