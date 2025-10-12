package mudp

import (
	"fmt"
	"net"

	"github.com/m-startgo/go-utils/mstr"
)

type Sender struct {
	IP   string
	Port int
	Addr string
	Conn net.Conn
}

func NewSender(opt Sender) (send *Sender, err error) {
	err = nil
	send = &Sender{}

	send.IP = opt.IP
	if send.IP == "" {
		send.IP = "127.0.0.1"
	}
	send.Port = opt.Port
	if send.Port == 0 {
		err = fmt.Errorf("mudp.NewSender|Port 不能为空")
		return
	}

	if opt.Addr == "" {
		send.Addr = mstr.Join(opt.IP, ":", opt.Port)
	} else {
		send.Addr = opt.Addr
	}

	conn, err := net.Dial("udp", send.Addr)
	if err != nil {
		err = fmt.Errorf("mudp.NewSender|net.Dial 失败|%v", err)
		return
	}

	send.Conn = conn

	return
}

func (s *Sender) Write(data []byte) (n int, err error) {
	if s.Conn == nil {
		err = fmt.Errorf("mudp.Write|Conn 不能为空")
		return
	}
	n, err = s.Conn.Write(data)
	if err != nil {
		err = fmt.Errorf("mudp.Write|Write 失败|%v", err)
		return
	}
	return
}

func (s *Sender) Close() (err error) {
	return s.Conn.Close()
}
