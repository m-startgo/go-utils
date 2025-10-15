package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/m-startgo/go-utils/mjson"
	"github.com/m-startgo/go-utils/mtcp"
)

var (
	PORT = 9999
	IP   = "127.0.0.1"
)

// 简化版 demo：每秒向目标 TCP 地址发送一条消息。
func main() {
	sender, err := mtcp.NewSender(mtcp.Sender{
		IP:   IP,
		Port: PORT,
	})
	if err != nil {
		fmt.Println("创建 Sender 失败:", err)
		return
	}

	i := 0
	for {
		i++
		timeNow := time.Now().UnixNano()
		data := map[string]string{
			"time": strconv.FormatInt(timeNow, 10),
			"id":   strconv.Itoa(i),
			"msg":  "hello tcp",
		}
		dataByte, _ := mjson.ToByte(data)
		fmt.Println("发送", string(dataByte))

		_, err := sender.Write(dataByte)
		if err != nil {
			fmt.Println("发送失败:", err)
		}

		time.Sleep(time.Second)
	}
}
