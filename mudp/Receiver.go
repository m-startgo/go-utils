package mudp

import (
	"net"
)

// NewReceiver 将 UDP 监听器绑定到 localAddr（例如 ":9000" 或 "127.0.0.1:9000"）。
func NewReceiver(localAddr string) (r *Receiver, err error) {
	r = &Receiver{}
	laddr, e := net.ResolveUDPAddr("udp", localAddr)
	if e != nil {
		err = e
		return
	}
	conn, e := net.ListenUDP("udp", laddr)
	if e != nil {
		err = e
		return
	}
	r.conn = conn
	return
}
