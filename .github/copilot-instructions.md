<!--
为 AI 编程代理（如 Copilot）准备的项目特定指令。
目的：帮助代理快速理解本仓库结构、编码约定、常用工作流与可复用模式，便于高质量 PR/补丁的生成。
-->

# Copilot / AI 代理使用说明（简明）

本项目是一个小型的 Go 实用函数库，按功能分成若干包（例如 `m_str`, `m_math`, `m_cycle`, `m_cron`），目标是提供开箱即用的工具函数。

以下说明帮助 AI 代理在修改、添加或重构代码时保持一致性并生成高质量变更。

## 要点速览

- 模块路径：`module github.com/m-startgo/go-utils`（见 `go.mod`）。
- Go 版本：`go 1.25`（在 `go.mod` 中声明）。
- 主要依赖：`github.com/shopspring/decimal`（高精度小数封装），`github.com/robfig/cron/v3`（cron 调度）。
- 风格：轻量、单文件包实现；函数与类型以驼峰/首字母大写导出，注释使用中文或英文并紧跟导出的符号（同一风格）；测试文件以 `_test.go` 存放。

## 项目架构与设计意图（“为什么”）

- 每个子目录实现一个小工具集合，目标是低耦合、零依赖或少依赖地提供常用功能。例如：

  - `m_str`：字符串格式化与通用 ToStr/Join 工具（见 `m_str/TplFormat.go`, `m_str/ToStr.go`）。
  - `m_math`：数值工具，当前使用 `shopspring/decimal` 封装（`m_math/decimal.go`）。
  - `m_cycle`：简单定时/周期任务封装（`m_cycle/cycle.go`）。
  - `m_cron`：基于 `robfig/cron/v3` 的轻量 Cron 封装（`m_cron/New.go`）。

- 设计偏向：易用的封装、chainable API（例如 `m_math.Decimal` 支持链式算术），以及尽量在小包内隐藏第三方细节。

## 代理应遵循的具体规则

1. 保持 API 稳定：不要随意修改导出函数/类型的签名，除非有强烈理由并在 PR 描述中说明兼容性影响。
2. 注释与示例：导出符号应有注释（中文可接受），并尽量在文件顶部保留简短用法示例（多数文件已有）。
3. 引入依赖：谨慎引入新第三方库；优先考虑标准库或已经在 `go.mod` 中的依赖。
4. 错误处理：返回错误而不是 panic，除非函数以 `Must` 命名并在注释中声明会 panic（参考 `MustFromString` 模式）。
5. JSON/DB 兼容：对于数值类型（如 `Decimal`），优先使用字符串编码以保全精度；实现 `encoding.TextMarshaler`/`json.Marshaler` 和 `database/sql` 的 `Scanner`/`Valuer` 是受欢迎的改进（如果需要，请在修改中添加对应测试）。

## 常用工作流 & 命令（代理可在 PR 描述中建议）

- 本地构建/检查：

  - go mod tidy
  - go vet ./...
  - go test ./... -v

- 仅运行某个包的测试：`go test ./m_math -v`。

注意：仓库没有特殊构建脚本或 CI 文件可见；在建议更改时，包含上述命令作为验证步骤。

## 代码模式与示例（摘录）

- 链式值对象（`m_math/decimal.go`）

  - Pattern: 封装第三方类型为值对象 `type Decimal struct { d decimal.Decimal }`，暴露链式方法（`Add`, `Mul`, `Round` 等）并提供构造器 `NewFromString`, `MustFromString`。
  - 建议改进：为 JSON 与 database/sql 添加 `MarshalJSON`/`UnmarshalJSON` 与 `Scan`/`Value`，以避免精度丢失。

- 文本模版替换（`m_str/TplFormat.go`）

  - Pattern: 使用 `os.Expand` 结合 map[string]string 替换 `${key}` 风格占位符；缺少键时返回空字符串。

- 周期任务与 Cron（`m_cycle`, `m_cron`）
  - Pattern: 为常见定时场景提供轻量封装，导出 `New`、`Start`/`Stop`、`SetInterval` 等方法；`m_cron` 使用 `cron.WithSeconds()` 并在 `New` 中立即 `Start()`。

## 测试 & 质量门（项目可见实践）

- 大多数包包含 `_test.go` 示例测试；生成或修改代码时，应添加覆盖关键逻辑的单元测试（happy path + 边界，如空输入、零除、nil 回调）。
- 建议对对外导出的行为（特别是 JSON/DB 接口）增加 round-trip 测试。

## PR / 提交建议

- 为每个改动写一个简短说明，包含：变更动机、向后兼容性影响、如何验证（包含 go test 命令或示例）。
- 小的 API 扩展（例如增加接口实现）应同时添加对应测试。

## 参考文件（仓库中示例）

- `go.mod` — 模块名与依赖
- `README.md` — 项目总体说明与使用示例
- `m_math/decimal.go` — 封装第三方库的典型模式
- `m_str/TplFormat.go`, `m_str/ToStr.go` — 字符串工具示例
- `m_cycle/cycle.go`, `m_cron/New.go` — 定时/调度工具示例

---

如果你希望我把这份说明直接提交到仓库 (.github/copilot-instructions.md)，我可以创建文件并运行快速测试（go vet / go test）。请确认是否写入或先在 PR 分支中保存。
