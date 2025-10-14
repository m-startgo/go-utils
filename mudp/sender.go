package mudp

import (
	"fmt"
	"net"
	"time"
)

// Sender 基于 UDP 的简易发送器
// 输入: IP(可选, 默认 127.0.0.1), Port(必填), Timeout(可选)
// 输出: 写入字节数, 错误
// 错误模式: 参数校验错误, 网络错误
// 成功: 返回实际写入字节数
// 示例:
//  s, err := mudp.NewSender(mudp.Sender{IP: "127.0.0.1", Port: 9000})
//  defer s.Close()
//  s.Write([]byte("hello"))

type Sender struct {
	IP      string
	Port    int
	Timeout time.Duration // 可选, 0 表示无超时
	conn    *net.UDPConn
	addr    *net.UDPAddr
}

// NewSender 创建 Sender 并建立 UDP 连接
func NewSender(opt Sender) (s *Sender, err error) {
	s = &Sender{}
	if opt.IP == "" {
		opt.IP = "127.0.0.1"
	}
	if opt.Port == 0 {
		err = fmt.Errorf("err:mudp.NewSender|Port 不能为空")
		return
	}
	if opt.Timeout < 0 {
		err = fmt.Errorf("err:mudp.NewSender|Timeout 无效")
		return
	}

	s.IP = opt.IP
	s.Port = opt.Port
	s.Timeout = opt.Timeout

	s.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		err = fmt.Errorf("err:mudp.NewSender|ResolveUDPAddr|%w", err)
		return
	}

	// 使用 nil 本地地址，让内核选择
	conn, err := net.DialUDP("udp", nil, s.addr)
	if err != nil {
		err = fmt.Errorf("err:mudp.NewSender|DialUDP|%w", err)
		return
	}

	s.conn = conn
	return
}

// Write 发送数据
func (s *Sender) Write(data []byte) (n int, err error) {
	if s == nil || s.conn == nil {
		err = fmt.Errorf("err:mudp.Write|Conn 为空")
		return
	}
	if s.Timeout > 0 {
		err = s.conn.SetWriteDeadline(time.Now().Add(s.Timeout))
		if err != nil {
			err = fmt.Errorf("err:mudp.Write|SetWriteDeadline|%w", err)
			return
		}
	}

	n, err = s.conn.Write(data)
	if err != nil {
		err = fmt.Errorf("err:mudp.Write|Write|%w", err)
		return
	}
	return
}

// Close 关闭连接
func (s *Sender) Close() (err error) {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
