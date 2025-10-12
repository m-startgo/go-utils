package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/m-startgo/go-utils/mjson"
	"github.com/m-startgo/go-utils/mstr"
)

var (
	PORT   = 9999
	IPAddr = "127.0.0.1"
)

// 简化版 demo：每秒向目标 UDP 地址发送一条消息。
func main() {
	url := mstr.Join(IPAddr, ":", PORT)

	conn, err := net.Dial("udp", url)
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
		fmt.Println("发送", string(dataByte))

		_, err := conn.Write(dataByte)
		if err != nil {
			fmt.Println("发送失败:", err)
		}

		time.Sleep(time.Second)
	}
}
