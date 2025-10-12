package main

import (
	"fmt"
	"log"

	"github.com/m-startgo/go-utils/mws"
	"github.com/valyala/fasthttp"
)

func main() {
	h := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/ws":
			err := mws.ServeUpgrade(ctx, func(c *mws.Conn) {
				fmt.Println("client connected")
				c.SetOnMessage(func(mt int, data []byte) {
					fmt.Printf("recv: %s\n", string(data))
					// echo back
					_ = c.SendText("echo: " + string(data))
				})
				c.SetOnClose(func(err error) {
					fmt.Println("closed:", err)
				})
			})
			if err != nil {
				log.Println("upgrade error:", err)
			}
		default:
			ctx.SetStatusCode(404)
		}
	}

	fmt.Println("ws server listening :8080")
	if err := fasthttp.ListenAndServe(":8080", h); err != nil {
		log.Fatal(err)
	}
}
