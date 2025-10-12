package mudp

import (
	"net"
)

// SendTo 使用标准库 net.UDPConn 发送单个 UDP 包到目标地址。
// 这个函数与 gnet 无需耦合，便于在客户端短连接场景中使用。
func SendTo(remote string, data []byte) error {
	raddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	return err
}
