# mudp 使用说明

这是 `mudp` 包的简单使用示例。示例包含如何创建 UDP Receiver（监听）和 Sender（发送），以及常见的超时、关闭和并发注意事项。

## Receiver 使用示例

```bash

# Receiver 监听在本地 9000 端口，打印收到的消息内容和发送者地址
go run mudp/udp_receiver_demo/main.go

# 在另一个终端运行 Sender 发送消息
go run mudp/udp_sender_demo/main.go
```

## 注意事项

- Receiver 的 `Listen` 会为每个收到的数据包启动一个 goroutine 来执行 handler，若 handler 有较重的工作，建议在 handler 内使用有界的 worker 池以避免无限制并发。
- Sender 的 `Send` 支持 context 取消与 timeout（最后以最先到期的 deadline 为准）。发送时会将连接的写超时设置为合并后的 deadline，从而避免长期阻塞。
- 在高并发场景下，频繁创建大量 goroutine 进行 Write 可能导致资源压力，可考虑重用发送器并使用限流方案。

```text
示例依赖：
  - Go 1.25+
  - 包路径按项目相对导入或将模块名替换为你项目的模块路径
```

如需我为 `mudp` 添加单元测试来验证并发发送/关闭的行为，我可以继续添加并运行测试。
