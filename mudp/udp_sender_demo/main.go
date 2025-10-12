package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/m-startgo/go-utils/mjson"
)

// 简化版 demo：每秒向目标 UDP 地址发送一条消息。
func main() {
	conn, err := net.Dial("udp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("创建UDP Sender失败", err)
		return
	}
	defer conn.Close()

	i := 0
	for {
		i++
		timeNow := time.Now().UnixNano()
		data := map[string]string{
			"time": strconv.FormatInt(timeNow, 10),
			"id":   strconv.Itoa(i),
			"msg":  "hello udp",
		}
		dataByte, _ := mjson.ToByte(data)
		_, err := conn.Write(dataByte)
		if err != nil {
			fmt.Println("发送失败:", err)
		}
		time.Sleep(time.Second)
	}
}
