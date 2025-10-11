package mudp

import (
	"fmt"
	"net"
)

// NewReceiver binds a UDP listener to localAddr (e.g., ":9000" or "127.0.0.1:9000").
func NewReceiver(localAddr string) (r *Receiver, err error) {
	r = &Receiver{}
	laddr, e := net.ResolveUDPAddr("udp", localAddr)
	if e != nil {
		err = fmt.Errorf("err:utils.udp|NewReceiver|resolve local addr: %w", e)
		return
	}
	conn, e := net.ListenUDP("udp", laddr)
	if e != nil {
		err = fmt.Errorf("err:utils.udp|NewReceiver|listen udp: %w", e)
		return
	}
	r.conn = conn
	return
}
