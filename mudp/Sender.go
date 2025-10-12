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
		return
	}

	send.Conn = conn

	return
}
