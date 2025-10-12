# mudp 使用说明

提供基于 `github.com/panjf2000/gnet/v2` 的简单 UDP 封装：

- `NewServer(addr, handler, opts...)`：创建 UDP Server（addr 可不带协议，如 `:9001`）
- `Server.Start()` / `Server.StartAsync()` / `Server.Stop()`：生命周期管理
- `SendTo(remote, data)`：使用标准库发送单个 UDP 包

示例：

```bash
# 1. 启动接收端
go run ./mudp/udp_serve_demo

# 2. 运行发送端（在另一个终端）
go run ./mudp/udp_send_demo

```

接收端将打印收到的消息。
