package mudp

import (
	"fmt"
	"net"

	"github.com/m-startgo/go-utils/mstr"
)

// Sender 表示一个简易 UDP 发送客户端。
type Sender struct {
	IP   string
	Port int
	Conn net.Conn
	addr string
}

// NewSender 创建并连接到目标 UDP 地址。
// NOTE: 该函数会立即进行网络拨号，调用者负责在不需要时调用 Close()
func NewSender(opt Sender) (send *Sender, err error) {
	err = nil
	send = &Sender{}

	send.IP = opt.IP
	if send.IP == "" {
		send.IP = "127.0.0.1"
	}
	send.Port = opt.Port
	if send.Port == 0 {
		err = fmt.Errorf("err:mudp.NewSender|Port|不能为空")
		return
	}

	send.addr = mstr.Join(send.IP, ":", send.Port)

	conn, err := net.Dial("udp", send.addr)
	if err != nil {
		err = fmt.Errorf("err:mudp.NewSender|net.Dial|%v", err)
		return
	}

	send.Conn = conn

	return
}

// Write 发送数据到已连接的 UDP 目标。
func (s *Sender) Write(data []byte) (n int, err error) {
	if s.Conn == nil {
		err = fmt.Errorf("err:mudp.Write|Conn|不能为空")
		return
	}
	n, err = s.Conn.Write(data)
	if err != nil {
		err = fmt.Errorf("err:mudp.Write|Write|%v", err)
		return
	}
	return
}

// Close 关闭底层连接。
func (s *Sender) Close() (err error) {
	if s.Conn == nil {
		return nil
	}
	return s.Conn.Close()
}
