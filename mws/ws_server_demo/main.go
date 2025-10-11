package main

import (
	"log"
	"net/http"

	"github.com/m-startgo/go-utils/mstr"
	"github.com/m-startgo/go-utils/mws"
)

const (
	port = 9999
	IP   = "127.0.0.1"
)

func main() {
	http.Handle("/ws", mws.SimpleEchoHandler())
	addr := mstr.Join(IP, ":", port)
	log.Printf("ws server listening %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
