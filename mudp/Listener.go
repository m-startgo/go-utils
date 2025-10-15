package mudp

import (
	"fmt"

	"github.com/m-startgo/go-utils/mstr"
	"github.com/panjf2000/gnet/v2"
)

// OnMessageFunc 是接收到数据时的回调函数签名。
// eventName: 事件名(例如 "OnBoot"/"OnTraffic")
// data: 原始字节数据
type OnMessageFunc func(eventName string, data []byte)

// Listener 表示一个 UDP 监听器配置。
//
// 示例:
//
//	l, err := mudp.NewListener(mudp.Listener{IP: "127.0.0.1", Port: 9999, OnMessage: fn})
type Listener struct {
	IP        string
	Port      int
	MultiCore bool
	OnMessage OnMessageFunc
	addr      string
}

// echoServer 是内部用于 gnet 的事件引擎实现。
type echoServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	addr      string
	multiCore bool
	onMessage OnMessageFunc
}

// NewListener 根据传入配置创建并返回一个 Listener 实例。
// 返回的 Listener 仅为配置容器；调用者需要调用 Start() 启动服务。
func NewListener(opt Listener) (l *Listener, err error) {
	err = nil
	l = &Listener{}

	l.IP = opt.IP
	if l.IP == "" {
		l.IP = "127.0.0.1"
	}
	l.Port = opt.Port
	if l.Port == 0 {
		err = fmt.Errorf("err:mudp.NewListener|Port|不能为空")
		return
	}

	// gnet 使用的地址格式示例: udp://127.0.0.1:9999
	l.addr = mstr.Join("udp://", l.IP, ":", l.Port)

	l.OnMessage = opt.OnMessage
	if l.OnMessage == nil {
		// 提供默认空实现以避免 nil 调用
		l.OnMessage = func(eventName string, data []byte) {}
	}

	l.MultiCore = opt.MultiCore

	return
}

// Start 阻塞启动监听，返回非 nil 错误表示启动失败或运行期间出错。
func (l *Listener) Start() error {
	echo := &echoServer{
		addr:      l.addr,
		multiCore: l.MultiCore,
		onMessage: l.OnMessage,
	}
	err := gnet.Run(echo, echo.addr, gnet.WithMulticore(l.MultiCore))
	return err
}

// 引擎启动回调：向外部报告 OnBoot 事件
func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	msg := mstr.Join("listening on:", es.addr)
	go es.onMessage("OnBoot", []byte(msg)) // 异步回调以避免阻塞
	return gnet.None
}

// 数据到达回调：复制数据并异步上报
func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, err := c.Next(-1)
	if err != nil {
		return gnet.Close
	}
	// 复制 buf 到新的切片，避免 gnet 重用底层缓冲区导致竞争
	data := make([]byte, len(buf))
	copy(data, buf)
	go es.onMessage("OnTraffic", data) // 异步回调
	return gnet.None
}
