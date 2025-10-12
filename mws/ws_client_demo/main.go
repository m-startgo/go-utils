package main

import (
	"context"
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
	url := mstr.Join("ws://", IP, ":", port, "/ws")

	// 拨号到 ws 服务器
	conn, _, err := mws.DialContext(context.Background(), url, http.Header{})
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	log.Printf("已连接到服务器：%s", url)

	var i int

	for {
		i++
		timeNow := time.Now().UnixNano()
		data := map[string]string{
			"time": strconv.FormatInt(timeNow, 10),
			"id":   strconv.Itoa(i),
			"msg":  "hello ws-server",
		}
		dataByte, _ := mjson.ToByte(data)
		if err := conn.WriteMessage(2, dataByte); err != nil {
			log.Printf("发送错误: %v", err)
			return
		}
		mt, rmsg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取回应错误: %v", err)
			return
		}
		timeNow2 := time.Now().UnixNano()
		log.Println("收到回应:", mt, string(rmsg), timeNow2)

		time.Sleep(5 * time.Second)
	}
}
