package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/m-startgo/go-utils/mws"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := "ws://127.0.0.1:9999/ws"
	conn, _, err := mws.DialContext(ctx, url, http.Header{})
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// send text
	if err := conn.WriteMessage(1, []byte("hello from client")); err != nil {
		log.Fatalf("write: %v", err)
	}

	// read echo
	mt, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("read: %v", err)
	}
	log.Printf("recv(%d): %s", mt, string(msg))
}
