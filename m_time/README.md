# m_time

m_time 是一个轻量的 Go 时间工具包，提供常见的解析、格式化和便捷方法（StartOfDay/EndOfDay、AddDays、AddHours 等）。

主要功能

- Parse 支持多种输入：日期字符串、数字字符串（秒/毫秒/微秒/纳秒）、整数/浮点数时间戳。
- Format 支持简单的 token，例如 `YYYY-MM-DD HH:mm:ss.SSSSSS ±HH:MM`。
- 可配置数字时间戳解析后的默认时区：`SetDefaultLocation(loc *time.Location)`，传入 `nil` 恢复为 UTC（默认）。
- 提供 `Time` 的轻量封装，方便链式调用和常用操作。

示例

```go
import (
    "fmt"
    "time"

    "github.com/m-startgo/go-utils/m_time"
)

func example() {
    // 解析字符串
    t, err := m_time.Parse("2021-01-02 15:04:05")
    if err != nil {
        panic(err)
    }
    fmt.Println(t.Format(m_time.DefaultToken))

    // 解析 13 位毫秒时间戳
    t2, _ := m_time.Parse(1609459200000)
    fmt.Println(t2.ToTime())

    // 将数字时间戳解析为本地时区
    m_time.SetDefaultLocation(time.Local)
    t3, _ := m_time.Parse(1609459200000)
    fmt.Println(t3.ToTime())
    // 恢复默认 UTC
    m_time.SetDefaultLocation(nil)
}
```

注意

- `MustParse` 不再 panic；它在解析失败时返回零值 `Time{}`。建议在需要错误信息时使用 `Parse`。
- 底层字符串解析依赖 `github.com/araddon/dateparse`，更多复杂格式可由该库支持。
