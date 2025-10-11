package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-startgo/go-utils/mws"
)

var upgrader = websocket.Upgrader{}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		return
	}
	c := mws.NewFromWS(ws)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = c.Listen(ctx, func(mt int, data []byte) {
			fmt.Println("server recv:", string(data))
			// echo back
			_, _ = c.Send(context.Background(), data, time.Second)
		})
	}()
}

func main() {
	http.HandleFunc("/echo", echoHandler)
	addr := ":9999"
	fmt.Println("ws server listening", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("server error:", err)
	}
}
