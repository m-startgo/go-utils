package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/m-startgo/go-utils/mjson"
	"github.com/m-startgo/go-utils/mudp"
)

func main() {
	// 可选绑定本地地址，例如 ":0" 让系统自动选择端口
	s, err := mudp.NewSender(":0")
	if err != nil {
		panic(err)
	}
	defer s.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

		fmt.Println("sender", string(dataByte))
		_, err := s.Send(ctx, "127.0.0.1:9000", dataByte, time.Second)
		if err != nil {
			fmt.Println("sender error:", err)
			return
		}
		time.Sleep(time.Millisecond * 100) // 每 100 毫秒发送一次
	}
}
