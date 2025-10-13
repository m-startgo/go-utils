package mws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/m-startgo/go-utils/mstr"
)

type (
	OnMessageFunc func(eventName string, data any)
	OnErrorFunc   func(eventName string, err error)
)

type Server struct {
	IP        string
	Port      int
	Path      string
	OnMessage OnMessageFunc
	OnError   OnErrorFunc
	addr      string
}

func NewServer(opt Server) (server *Server, err error) {
	err = nil
	server = &Server{}

	server.IP = opt.IP
	if server.IP == "" {
		server.IP = "127.0.0.1"
	}

	server.Path = opt.Path
	if server.Path == "" {
		server.Path = "/ws"
	}

	server.Port = opt.Port
	if server.Port == 0 {
		err = fmt.Errorf("mudp.NewServer|Port 不能为空")
		return
	}

	server.addr = mstr.Join(server.IP, ":", server.Port)

	server.OnMessage = opt.OnMessage
	if server.OnMessage == nil {
		server.OnMessage = func(eventName string, data any) {
			// 默认空实现，避免 nil 指针异常
		}
	}

	server.OnError = opt.OnError
	if server.OnError == nil {
		server.OnError = func(eventName string, err error) {
			// 默认空实现，避免 nil 指针异常
		}
	}

	return
}

func (server *Server) Start() error {
	http.HandleFunc(server.Path, func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Printf("err:mws.ws_server_demo|Accept|%v", err)
			http.Error(w, "websocket accept error", http.StatusBadRequest)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var v any
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			server.OnError("Read", err)
			return
		}
		server.OnMessage("Read", v)

		for {
			// 回送一个确认消息
			ack := map[string]any{"ok": true, "recv": v}
			err = wsjson.Write(ctx, c, ack)
			if err != nil {
				log.Printf("err:mws.ws_server_demo|Write|%v", err)
				_ = c.Close(websocket.StatusInternalError, "write error")
				return
			}
		}
	})

	return http.ListenAndServe(server.addr, nil)
}
