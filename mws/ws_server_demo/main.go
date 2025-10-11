package main

import (
	"log"
	"net/http"

	"github.com/m-startgo/go-utils/mws"
)

func main() {
	http.Handle("/ws", mws.SimpleEchoHandler())
	addr := "127.0.0.1:9999"
	log.Printf("ws server listening %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
