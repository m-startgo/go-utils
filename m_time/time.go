package m_time

/*
我想封装一些关于时间方面的工具函数库

它的使用方式和API以及常用功能最好是可以和 dayjs 对齐

你可以使用 github.com/araddon/dateparse 库来实现日期字符串解析相关的操作，这样可以帮我们节省大量的时间和精力

*/

import (
	"time"

	"github.com/shopspring/decimal"
)

// Format 格式化时间
// 现在时间: m_time.New().Format("2006-01-02 15:04:05") // 输出: 2023-10-15 13:45:26
// 指定时间: m_time.NewFromTime(time.Date(2023, 10, 15, 0, 0, 0, 0, time.Local)).Format("2006-01-02") // 输出: 2023-10-15
// 解析字符串: t, _ := m_time.NewFromString("2023-10-15"); t.Format("2006/01/02") // 输出: 2023/10/15
func (t *Time) Format(layout string) string {
	return t.time.Format(layout)
}

// In 修改时间的时区
// 参数: location *time.Location - 目标时区
// 返回: *Time - 转换到指定时区后的新 Time 实例
// 示例: 
// loc, _ := time.LoadLocation("America/New_York")
// m_time.New().In(loc).Format("2006-01-02 15:04:05") // 输出: 当前时间的纽约时区时间
// loc, _ := time.LoadLocation("Asia/Tokyo")
// m_time.New().In(loc).Format("2006-01-02 15:04:05") // 输出: 当前时间的东京时区时间
func (t *Time) In(location *time.Location) *Time {
	return &Time{time: t.time.In(location)}
}

// Add 添加时间
// 参数: d time.Duration - 要添加的时间长度
// 返回: *Time - 添加时间后的新 Time 实例
// 示例: m_time.New().Add(24 * time.Hour).Format("2006-01-02") // 输出: 当前日期的明天
func (t *Time) Add(d time.Duration) *Time {
	return &Time{time: t.time.Add(d)}
}

// Subtract 减去时间
// 参数: d time.Duration - 要减去的时间长度
// 返回: *Time - 减去时间后的新 Time 实例
// 示例: m_time.New().Subtract(24 * time.Hour).Format("2006-01-02") // 输出: 当前日期的昨天
func (t *Time) Subtract(d time.Duration) *Time {
	return &Time{time: t.time.Add(-d)}
}

// StartOf 获取时间的开始
// 参数: unit string - 时间单位(year, month, day, hour, minute, second)
// 返回: *Time - 对应时间单位开始的新 Time 实例
// 示例: m_time.New().StartOf("day").Format("2006-01-02 15:04:05") // 输出: 当前日期的 00:00:00
func (t *Time) StartOf(unit string) *Time {
	switch unit {
	case "year":
		return &Time{time: time.Date(t.time.Year(), 1, 1, 0, 0, 0, 0, t.time.Location())}
	case "month":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), 1, 0, 0, 0, 0, t.time.Location())}
	case "day":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), 0, 0, 0, 0, t.time.Location())}
	case "hour":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), t.time.Hour(), 0, 0, 0, t.time.Location())}
	case "minute":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), t.time.Hour(), t.time.Minute(), 0, 0, t.time.Location())}
	case "second":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), t.time.Hour(), t.time.Minute(), t.time.Second(), 0, t.time.Location())}
	default:
		return t
	}
}

// EndOf 获取时间的结束
// 参数: unit string - 时间单位(year, month, day, hour, minute, second)
// 返回: *Time - 对应时间单位结束的新 Time 实例
// 示例: m_time.New().EndOf("day").Format("2006-01-02 15:04:05") // 输出: 当前日期的 23:59:59.999999999
func (t *Time) EndOf(unit string) *Time {
	switch unit {
	case "year":
		return &Time{time: time.Date(t.time.Year(), 12, 31, 23, 59, 59, 999999999, t.time.Location())}
	case "month":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.AddDate(0, 1, 0).Add(-time.Nanosecond).Day(), 23, 59, 59, 999999999, t.time.Location())}
	case "day":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), 23, 59, 59, 999999999, t.time.Location())}
	case "hour":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), t.time.Hour(), 59, 59, 999999999, t.time.Location())}
	case "minute":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), t.time.Hour(), t.time.Minute(), 59, 999999999, t.time.Location())}
	case "second":
		return &Time{time: time.Date(t.time.Year(), t.time.Month(), t.time.Day(), t.time.Hour(), t.time.Minute(), t.time.Second(), 999999999, t.time.Location())}
	default:
		return t
	}
}

// Diff 计算时间差
// 参数: t2 *Time - 要比较的另一个 Time 实例, unit string - 时间单位
// 返回: decimal.Decimal - 两个时间之间的差值
// 示例: t1 := m_time.New(); t2 := m_time.New().Add(24 * time.Hour); t2.Diff(t1, "day") // 输出: 1
func (t *Time) Diff(t2 *Time, unit string) decimal.Decimal {
	switch unit {
	case "year":
		return decimal.NewFromInt(int64(t.time.Year() - t2.time.Year()))
	case "month":
		return decimal.NewFromInt(int64((t.time.Year()-t2.time.Year())*12 + int(t.time.Month()-t2.time.Month())))
	case "day":
		return decimal.NewFromInt(int64(t.time.Sub(t2.time).Hours() / 24))
	case "hour":
		return decimal.NewFromInt(int64(t.time.Sub(t2.time).Hours()))
	case "minute":
		return decimal.NewFromInt(int64(t.time.Sub(t2.time).Minutes()))
	case "second":
		return decimal.NewFromInt(int64(t.time.Sub(t2.time).Seconds()))
	case "millisecond":
		return decimal.NewFromInt(t.time.Sub(t2.time).Milliseconds())
	case "microsecond":
		return decimal.NewFromInt(t.time.Sub(t2.time).Microseconds())
	case "nanosecond":
		return decimal.NewFromInt(t.time.Sub(t2.time).Nanoseconds())
	default:
		return decimal.NewFromInt(0)
	}
}

// Unix 获取 Unix 时间戳
// 返回: int64 - Unix 时间戳(秒)
// 示例: m_time.New().Unix() // 输出: 1697347526
func (t *Time) Unix() int64 {
	return t.time.Unix()
}

// UnixMilli 获取 Unix 毫秒时间戳
// 返回: int64 - Unix 时间戳(毫秒)
// 示例: m_time.New().UnixMilli() // 输出: 1697347526123
func (t *Time) UnixMilli() int64 {
	return t.time.UnixMilli()
}

// UnixMicro 获取 Unix 微秒时间戳
// 返回: int64 - Unix 时间戳(微秒)
// 示例: m_time.New().UnixMicro() // 输出: 1697347526123456
func (t *Time) UnixMicro() int64 {
	return t.time.UnixMicro()
}

// UnixNano 获取 Unix 纳秒时间戳
// 返回: int64 - Unix 时间戳(纳秒)
// 示例: m_time.New().UnixNano() // 输出: 1697347526123456789
func (t *Time) UnixNano() int64 {
	return t.time.UnixNano()
}

// DaysInMonth 获取月份的天数
// 返回: int - 月份的天数
// 示例: m_time.New().DaysInMonth() // 输出: 31(取决于当前月份)
func (t *Time) DaysInMonth() int {
	return t.EndOf("month").time.Day()
}

// Time 获取 time.Time
// 返回: time.Time - 底层的 time.Time 实例
// 示例: m_time.New().Time() // 输出: 等同于 time.Now()
func (t *Time) Time() time.Time {
	return t.time
}

// String 实现 Stringer 接口
// 返回: string - 时间的字符串表示
// 示例: m_time.New().String() // 输出: 2023-10-15 13:45:26.123456789 +0800 CST
func (t *Time) String() string {
	return t.time.String()
}
