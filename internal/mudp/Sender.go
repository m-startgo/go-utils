//go:build ignore

package mudp

// import (
// 	"fmt"
// 	"net"

// 	"github.com/m-startgo/go-utils/mstr"
// )

// type Sender struct {
// 	IP   string
// 	Port int
// 	Conn net.Conn
// 	addr string
// }

// func NewSender(opt Sender) (send *Sender, err error) {
// 	err = nil
// 	send = &Sender{}

// 	send.IP = opt.IP
// 	if send.IP == "" {
// 		send.IP = "127.0.0.1"
// 	}
// 	send.Port = opt.Port
// 	if send.Port == 0 {
// 		err = fmt.Errorf("mudp.NewSender|Port 不能为空")
// 		return
// 	}

// 	send.addr = mstr.Join(send.IP, ":", send.Port)

// 	conn, err := net.Dial("udp", send.addr)
// 	if err != nil {
// 		err = fmt.Errorf("mudp.NewSender|net.Dial 失败|%v", err)
// 		return
// 	}

// 	send.Conn = conn

// 	return
// }

// func (s *Sender) Write(data []byte) (n int, err error) {
// 	if s.Conn == nil {
// 		err = fmt.Errorf("mudp.Write|Conn 不能为空")
// 		return
// 	}
// 	n, err = s.Conn.Write(data)
// 	if err != nil {
// 		err = fmt.Errorf("mudp.Write|Write 失败|%v", err)
// 		return
// 	}
// 	return
// }

// func (s *Sender) Close() (err error) {
// 	if s.Conn == nil {
// 		return nil
// 	}
// 	return s.Conn.Close()
// }
