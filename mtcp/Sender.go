package mtcp

import (
	"fmt"
	"net"

	"github.com/m-startgo/go-utils/mstr"
)

// Sender 简单的 TCP 客户端封装。
type Sender struct {
	IP   string
	Port int
	Conn net.Conn
	addr string
}

// NewSender 创建并连接到目标 TCP 地址。
func NewSender(opt Sender) (send *Sender, err error) {
	err = nil
	send = &Sender{}

	send.IP = opt.IP
	if send.IP == "" {
		send.IP = "127.0.0.1"
	}
	send.Port = opt.Port
	if send.Port == 0 {
		err = fmt.Errorf("err:mtcp.NewSender|Port|不能为空")
		return
	}

	send.addr = mstr.Join(send.IP, ":", send.Port)

	conn, err := net.Dial("tcp", send.addr)
	if err != nil {
		err = fmt.Errorf("err:mtcp.NewSender|net.Dial|%v", err)
		return
	}

	send.Conn = conn

	return
}

// Write 发送数据到已连接的 TCP 目标。
func (s *Sender) Write(data []byte) (n int, err error) {
	if s.Conn == nil {
		err = fmt.Errorf("err:mtcp.Write|Conn|不能为空")
		return
	}
	n, err = s.Conn.Write(data)
	if err != nil {
		err = fmt.Errorf("err:mtcp.Write|Write|%v", err)
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
