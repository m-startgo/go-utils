# mws — 轻量的 WebSocket 封装

这是一个基于 `github.com/gorilla/websocket` 的轻量封装，放在仓库的 `mws` 包中，目标是提供简单、易用的 WebSocket 客户端/服务端升级与读写接口，适合在示例或小型服务中直接使用。

主要特性：

- 提供 `Conn` 封装（包裹 `*websocket.Conn`），包含常用的 JSON/text 读写方法。
- 提供 `Upgrade` 帮助函数，将 HTTP 请求升级为 WebSocket 连接（基于 `websocket.Upgrader`）。
- 提供 `DialContext` 以 context 为入参的客户端连接方法。
- 简单的 `SimpleEchoHandler` 示例 handler，用于 demo。

## 快速开始

```bash
# 启动 demo server（在仓库根目录下运行）：
go run ./mws/ws_server_demo

# 在另一个终端运行 client：
go run ./mws/ws_client_demo
```

（注意：示例里 server 监听在 `127.0.0.1:9999`，client 默认连接 `ws://127.0.0.1:9999/ws`）

## API 概览

包内实现的主要符号：

- `type Conn` — 封装的连接类型，方法包括：

  - `WriteJSON(v any) error` — 使用默认发送超时写入 JSON。
  - `ReadJSON(v any) error` — 读取 JSON 消息并解码到 `v`。
  - `WriteMessage(msgType int, data []byte) error` — 写入文本/二进制消息。
  - `ReadMessage() (int, []byte, error)` — 读取原始消息（返回 message type + bytes）。
  - `SetDeadlines(send, read time.Duration)` — 设置默认的发送与读取超时时间（0 表示不使用超时）。
  - `Close() error` — 关闭连接。

- `func NewConn(ws *websocket.Conn) *Conn` — 包装已存在的 `*websocket.Conn`。
- `func DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*Conn, *http.Response, error)` — 使用默认 dialer 发起客户端连接。
- `var Upgrader websocket.Upgrader` — 提供默认 `Upgrader`（默认允许所有 origin，示例便于开发测试）。
- `func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error)` — 将 HTTP 请求升级为 WebSocket。
- `func SimpleEchoHandler() http.Handler` — 一个简单的 echo handler，用于示例和快速验证。

示例（服务端）：

```go
http.Handle("/ws", mws.SimpleEchoHandler())
log.Fatal(http.ListenAndServe(":8080", nil))
```

示例（客户端）：

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

conn, _, err := mws.DialContext(ctx, "ws://localhost:8080/ws", http.Header{})
defer conn.Close()

_ = conn.WriteMessage(websocket.TextMessage, []byte("hello"))
_, msg, _ := conn.ReadMessage()
fmt.Println(string(msg))
```

## 生产注意事项

- Origin 校验：当前 `Upgrader.CheckOrigin` 默认为允许所有来源（方便开发和 demo）。在生产环境中，请务必实现更严格的校验，或覆盖 `Upgrader.CheckOrigin`。
- 并发写：`github.com/gorilla/websocket` 要求写操作必须串行化（不能并发写）。当前封装为轻量示例，并没有内部自动序列化写操作；如果你的应用会有多个 goroutine 同时写入，请为 `Conn` 添加写入队列或外部同步（例如单个写协程或 mutex）。
- 超时设置：默认发送超时为 10 秒，读超时默认禁用。请根据网络与业务特性调整 `SetDeadlines`。
- 资源回收：确保在连接不再使用时调用 `Close()`，并对长连接做心跳/超时管理以避免资源泄漏。
