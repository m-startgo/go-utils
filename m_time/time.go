package m_time

/*
我想封装一些关于时间方面的工具函数库

它的使用方式和API以及常用功能最好是可以和 dayjs 对齐

你可以使用 github.com/araddon/dateparse 库来实现日期字符串解析相关的操作，这样可以帮我们节省大量的时间和精力

*/

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/shopspring/decimal"
)

// Time 结构体，用于封装 time.Time
type Time struct {
	time time.Time
}

// New 创建一个新的 Time 实例
func New() *Time {
	return &Time{time: time.Now()}
}

// NewFromTime 从 time.Time 创建一个新的 Time 实例
func NewFromTime(t time.Time) *Time {
	return &Time{time: t}
}

// NewFromString 从字符串创建一个新的 Time 实例
func NewFromString(s string) (*Time, error) {
	tp, err := dateparse.ParseAny(s)
	if err != nil {
		return nil, err
	}
	return &Time{time: tp}, nil
}

// Format 格式化时间
func (t *Time) Format(layout string) string {
	return t.time.Format(layout)
}

// Add 添加时间
func (t *Time) Add(d time.Duration) *Time {
	return &Time{time: t.time.Add(d)}
}

// Subtract 减去时间
func (t *Time) Subtract(d time.Duration) *Time {
	return &Time{time: t.time.Add(-d)}
}

// StartOf 获取时间的开始
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
func (t *Time) Unix() int64 {
	return t.time.Unix()
}

// UnixMilli 获取 Unix 毫秒时间戳
func (t *Time) UnixMilli() int64 {
	return t.time.UnixMilli()
}

// UnixMicro 获取 Unix 微秒时间戳
func (t *Time) UnixMicro() int64 {
	return t.time.UnixMicro()
}

// UnixNano 获取 Unix 纳秒时间戳
func (t *Time) UnixNano() int64 {
	return t.time.UnixNano()
}

// DaysInMonth 获取月份的天数
func (t *Time) DaysInMonth() int {
	return t.EndOf("month").time.Day()
}

// Time 获取 time.Time
func (t *Time) Time() time.Time {
	return t.time
}

// String 实现 Stringer 接口
func (t *Time) String() string {
	return t.time.String()
}
