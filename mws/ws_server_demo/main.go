package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/m-startgo/go-utils/mjson"
	"github.com/m-startgo/go-utils/mstr"
	"github.com/m-startgo/go-utils/mws"
)

const (
	port = 9999
	IP   = "127.0.0.1"
)

func main() {
	// 监听 /ws
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := mws.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "升级失败", http.StatusBadRequest)
			return
		}
		defer c.Close()

		log.Printf("客户端已连接：%s", r.RemoteAddr)
		go func() {
			for {
				mt, msg, err := c.ReadMessage()
				if err != nil {
					log.Printf("读取错误: %v", err)
					break
				}
				timeNow := time.Now().UnixNano()

				log.Println("收到消息:", mt, string(msg), timeNow)
				// 立即回应客户端
				if mt > 0 {
					sendMsg := []byte("pong")
					if err := c.WriteMessage(1, sendMsg); err != nil {
						log.Printf("写入错误: %v", err)
						break
					}
				}
			}
		}()
		// 主动给客户端发消息
		var i int
		for {
			i++
			timeNow := time.Now().UnixNano()
			data := map[string]string{
				"time": strconv.FormatInt(timeNow, 10),
				"id":   strconv.Itoa(i),
				"msg":  "hello ws-client",
			}
			dataByte, _ := mjson.ToByte(data)
			c.WriteMessage(2, dataByte)
			time.Sleep(4 * time.Second)
		}
	})

	// 启动服务器
	addr := mstr.Join(IP, ":", port)
	log.Printf("ws 服务器正在监听 %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
