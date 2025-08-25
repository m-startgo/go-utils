## 目的

为 AI 编码代理提供可执行、项目特定的快速上手指引，使其能在本仓库中迅速定位关键约定、运行测试并提交符合仓库风格的更改。

## 快速事实（可直接引用）

- 模块路径：`github.com/m-startgo/go-utils`（见 `go.mod`）
- Go 版本：1.25（见 `go.mod`）
- 仓库类型：通用 Go 工具包，按功能分包，目录前缀均为 `m_*`（例如 `m_time`, `m_str`, `m_log`）
- 测试：各包包含 `_test.go`，可以运行 `go test ./...` 来执行全部测试
- 关键依赖（可在 `go.mod` 验证）：`github.com/araddon/dateparse`, `github.com/robfig/cron/v3`, `github.com/shopspring/decimal`

## 大局观（为什么这样组织）

- 这是一个以“功能分包”为主的工具库（每个 `m_*` 包提供一类常用工具函数或轻量封装）。
- 设计目标是轻量、可复用、无服务态（library-only），因此改动通常只影响单独包，跨包改动需谨慎并增加测试。

## 项目约定与要点（务必遵守）

- 包命名：以 `m_` 前缀分隔功能域（示例：`m_time` 实现时间解析/格式化，`m_str` 提供类型转字符串工具）。
- 注释与文档：导出函数/类型应有注释（参考 `m_time` 中大量注释）。
- 语言使用：UI/提示/日志可使用中文；代码示例、命令与技术文档优先使用英文以保证可复现性（来源：`.github/prompts/通用规则.prompt.md`）。
- 错误格式：仓库偏好包含可追踪信息，示例格式为 `err:<pkg.func>|<场景>|<错误信息>`（参照通用规则）。

## 包级/API 约定（具体、可复制的例子）

- m_time
  - 解析 API `Parse(v any)` 支持多种输入（string、int、float 等）；`MustParse` 在解析失败时不再 panic，而返回零值 `Time{}`（见 `m_time/time.go`）。
  - 默认时区由 `SetDefaultLocation(loc *time.Location)` 控制，传 `nil` 会恢复为 UTC（见 `m_time/time.go`）。
  - 默认 token：`DefaultToken = "YYYY-MM-DD HH:mm:ss.SSSSSS ±HH:MM"`；库内部将 token 映射为 Go layout（参见 `tokenToLayout` 函数）。
- m_str
  - `ToStr(any)` 将多种基本类型转为字符串（参见 `m_str/ToStr.go`），实现时优先处理字节切片与数值类型。

## 常用命令（在本仓库中实际可用）

- 运行全部测试：`go test ./...` 或更详细 `go test -v ./...`
- 运行单个包测试（示例）：`go test ./m_time -v`
- 格式化/静态检查（推荐）：`gofmt -w .`、`go vet ./...`、`golangci-lint run`（如增设 linter 配置则遵循之）

## 外部依赖与集成点

- 主要第三方库见 `go.mod`；新增依赖时请把 `go.mod` 保持干净（使用 `go get`/`go mod tidy` 并提交变化）。
- 库为纯 Go 包，无运行时服务；集成点通常是被上游项目 `import "github.com/m-startgo/go-utils/m_xxx"`。

## 提交 / PR 指南（AI 代理执行修改时遵循）

- 小步提交，接口变更必须包含向后兼容说明与测试。若删除或修改导出 API，请在 PR 描述中说明替代方案并标注迁移步骤。
- 变更依赖请更新 `go.mod` 并运行 `go test ./...` 确保没有回归。

## 可检查的代表文件（定位示例）

- 顶层 README: `README.md`（安装、示例用法）
- 时间工具与说明: `m_time/time.go`, `m_time/README.md`（解析/格式化/时区约定）
- 字符串工具: `m_str/ToStr.go`, `m_str/TplFormat.go`（模板格式化示例在顶层 README 中演示）
- 模块声明与依赖: `go.mod`

## 交互提示示例（AI 要如何提问或提交）

- 当不确定预期行为时，先列出 2 个合理实现选项并说明兼容性影响。例如："要把 X 方法的错误类型从 string 改成 error，A) 只修改内部实现（兼容），B) 修改签名并升级依赖（破坏性）；请选择。"
- 如果修改公共 API，自动生成或更新至少一个单元测试并在 PR 描述里说明回归测试步骤。
