package m_time

// 封装时间工具函数库，API设计参考dayjs风格
// 使用github.com/araddon/dateparse库进行日期字符串解析

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// 常用时间格式常量
const (
	FormatDateTime  = "2006-01-02 15:04:05"
	FormatDate      = "2006-01-02"
	FormatTime      = "15:04:05"
	FormatYearMonth = "2006-01"
	FormatMonthDay  = "01-02"
	FormatISO8601   = "2006-01-02T15:04:05Z07:00"
	FormatRFC3339   = "2006-01-02T15:04:05Z07:00"
)

// normalizeUnit 归一化时间单位
// 将各种时间单位表示形式统一转换为标准形式
// 参数: u string - 时间单位字符串，可以是多种形式（如'y', 'year', 'years', 'M', 'month'等）
// 返回: string - 归一化后的标准时间单位
// 示例: normalizeUnit("y") // 返回: "year"
//
//	normalizeUnit("M") // 返回: "month"
//	normalizeUnit("days") // 返回: "day"
func normalizeUnit(u string) string {
	// 保留原始输入的大小写判断，以便支持 Dayjs 风格的 'M' 表示 month
	s := strings.TrimSpace(u)
	if s == "M" {
		return "month"
	}
	normalized := strings.ToLower(s)
	switch normalized {
	case "y", "yr", "yrs", "years", "year":
		return "year"
	case "mon", "months", "month":
		return "month"
	case "d", "day", "days":
		return "day"
	case "h", "hour", "hours":
		return "hour"
	case "m", "min", "minute", "minutes":
		return "minute"
	case "s", "sec", "second", "seconds":
		return "second"
	case "ms", "millisecond", "milliseconds":
		return "millisecond"
	case "us", "µs", "microsecond", "microseconds":
		return "microsecond"
	case "ns", "nanosecond", "nanoseconds":
		return "nanosecond"
	default:
		return normalized
	}
}

// Format 格式化时间
// 参数: layout string - 时间格式模板
// 返回: string - 格式化后的时间字符串
// 示例: m_time.New().Format("2006-01-02 15:04:05")
func (t *Time) Format(layout string) string {
	if t == nil {
		return ""
	}
	return t.tm.Format(layout)
}

// In 修改时间的时区
// 参数: location *time.Location - 目标时区
// 返回: *Time - 转换到指定时区后的新 Time 实例
// 示例: loc, _ := time.LoadLocation("America/New_York"); m_time.New().In(loc).Format("2006-01-02 15:04:05")
func (t *Time) In(location *time.Location) *Time {
	if t == nil {
		return nil
	}
	return &Time{tm: t.tm.In(location)}
}

// Add 添加时间
// 参数: d time.Duration - 要添加的时间长度
// 返回: *Time - 添加时间后的新 Time 实例
// 示例: m_time.New().Add(24 * time.Hour).Format("2006-01-02")
func (t *Time) Add(d time.Duration) *Time {
	if t == nil {
		return nil
	}
	return &Time{tm: t.tm.Add(d)}
}

// Subtract 减去时间
// 参数: d time.Duration - 要减去的时间长度
// 返回: *Time - 减去时间后的新 Time 实例
// 示例: m_time.New().Subtract(24 * time.Hour).Format("2006-01-02")
func (t *Time) Subtract(d time.Duration) *Time {
	if t == nil {
		return nil
	}
	return &Time{tm: t.tm.Add(-d)}
}

// StartOf 获取时间的开始
// 参数: unit string - 时间单位(year, month, day, hour, minute, second)
// 返回: *Time - 对应时间单位开始的新 Time 实例
// 示例: m_time.New().StartOf("day").Format("2006-01-02 15:04:05")
func (t *Time) StartOf(unit string) *Time {
	if t == nil {
		return nil
	}
	switch normalizeUnit(unit) {
	case "year":
		return &Time{tm: time.Date(t.tm.Year(), 1, 1, 0, 0, 0, 0, t.tm.Location())}
	case "month":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), 1, 0, 0, 0, 0, t.tm.Location())}
	case "day":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), 0, 0, 0, 0, t.tm.Location())}
	case "hour":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), t.tm.Hour(), 0, 0, 0, t.tm.Location())}
	case "minute":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), t.tm.Hour(), t.tm.Minute(), 0, 0, t.tm.Location())}
	case "second":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), t.tm.Hour(), t.tm.Minute(), t.tm.Second(), 0, t.tm.Location())}
	default:
		return t
	}
}

// EndOf 获取时间的结束
// 参数: unit string - 时间单位(year, month, day, hour, minute, second)
// 返回: *Time - 对应时间单位结束的新 Time 实例
// 示例: m_time.New().EndOf("day").Format("2006-01-02 15:04:05")
func (t *Time) EndOf(unit string) *Time {
	if t == nil {
		return nil
	}
	switch normalizeUnit(unit) {
	case "year":
		return &Time{tm: time.Date(t.tm.Year(), 12, 31, 23, 59, 59, 999999999, t.tm.Location())}
	case "month":
		// 先取下月的第一天的 00:00:00，然后减 1 纳秒，得到本月最后一纳秒；直接返回该时刻
		firstOfNext := time.Date(t.tm.Year(), t.tm.Month(), 1, 0, 0, 0, 0, t.tm.Location()).AddDate(0, 1, 0)
		last := firstOfNext.Add(-time.Nanosecond)
		return &Time{tm: last}
	case "day":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), 23, 59, 59, 999999999, t.tm.Location())}
	case "hour":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), t.tm.Hour(), 59, 59, 999999999, t.tm.Location())}
	case "minute":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), t.tm.Hour(), t.tm.Minute(), 59, 999999999, t.tm.Location())}
	case "second":
		return &Time{tm: time.Date(t.tm.Year(), t.tm.Month(), t.tm.Day(), t.tm.Hour(), t.tm.Minute(), t.tm.Second(), 999999999, t.tm.Location())}
	default:
		return t
	}
}

// 为DaysInMonth添加缓存
var daysInMonthCache = make(map[string]int)

// DaysInMonth 获取当前时间所在月份的天数
// 示例: m_time.New().DaysInMonth() // 如10月返回31
func (t *Time) DaysInMonth() int {
	if t == nil {
		return 0
	}
	// 构造缓存键
	cacheKey := fmt.Sprintf("%d-%d", t.tm.Year(), t.tm.Month())
	if days, ok := daysInMonthCache[cacheKey]; ok {
		return days
	}
	// 计算当月天数
	firstOfNext := time.Date(t.tm.Year(), t.tm.Month(), 1, 0, 0, 0, 0, t.tm.Location()).AddDate(0, 1, 0)
	last := firstOfNext.Add(-time.Nanosecond)
	days := last.Day()
	// 存入缓存
	daysInMonthCache[cacheKey] = days
	return days
}

// Diff 计算时间差
// 对 year/month 做按整月计算并返回小数（月为单位再除以12 得到年）
func (t *Time) Diff(t2 *Time, unit string) decimal.Decimal {
	if t == nil || t2 == nil {
		return decimal.NewFromInt(0)
	}
	switch normalizeUnit(unit) {
	case "year":
		// 以月为单位计算更精确的年（包含月份的小数部分）
		months := int64((t.tm.Year()-t2.tm.Year())*12 + int(t.tm.Month()-t2.tm.Month()))
		// 天的差转换为月的补充（以 t2 当月天数为基准）
		days := t.tm.Sub(t2.tm).Hours() / 24.0
		approxMonths := float64(months) + days/float64(maxInt(1, t2.DaysInMonth()))
		return decimal.NewFromFloat(approxMonths / 12.0)
	case "month":
		months := int64((t.tm.Year()-t2.tm.Year())*12 + int(t.tm.Month()-t2.tm.Month()))
		days := t.tm.Sub(t2.tm).Hours() / 24.0
		approxMonths := float64(months) + days/float64(maxInt(1, t2.DaysInMonth()))
		return decimal.NewFromFloat(approxMonths)
	case "day":
		return decimal.NewFromFloat(t.tm.Sub(t2.tm).Hours() / 24.0)
	case "hour":
		return decimal.NewFromFloat(t.tm.Sub(t2.tm).Hours())
	case "minute":
		return decimal.NewFromFloat(t.tm.Sub(t2.tm).Minutes())
	case "second":
		return decimal.NewFromFloat(t.tm.Sub(t2.tm).Seconds())
	case "millisecond":
		return decimal.NewFromInt(t.tm.Sub(t2.tm).Milliseconds())
	case "microsecond":
		return decimal.NewFromInt(t.tm.Sub(t2.tm).Microseconds())
	case "nanosecond":
		return decimal.NewFromInt(t.tm.Sub(t2.tm).Nanoseconds())
	default:
		return decimal.NewFromInt(0)
	}
}

// Unix 获取 Unix 时间戳
// 返回: int64 - Unix 时间戳(秒)
// 示例: m_time.New().Unix() // 输出: 1697347526
func (t *Time) Unix() int64 {
	if t == nil {
		return 0
	}
	return t.tm.Unix()
}

// UnixMilli 获取 Unix 毫秒时间戳
// 返回: int64 - Unix 时间戳(毫秒)
// 示例: m_time.New().UnixMilli() // 输出: 1697347526123
func (t *Time) UnixMilli() int64 {
	if t == nil {
		return 0
	}
	return t.tm.UnixMilli()
}

// UnixMicro 获取 Unix 微秒时间戳
// 返回: int64 - Unix 时间戳(微秒)
// 示例: m_time.New().UnixMicro() // 输出: 1697347526123456
func (t *Time) UnixMicro() int64 {
	if t == nil {
		return 0
	}
	return t.tm.UnixMicro()
}

// UnixNano 获取 Unix 纳秒时间戳
// 返回: int64 - Unix 时间戳(纳秒)
// 示例: m_time.New().UnixNano() // 输出: 1697347526123456789
func (t *Time) UnixNano() int64 {
	if t == nil {
		return 0
	}
	return t.tm.UnixNano()
}

// Time 获取 time.Time
// 返回: time.Time - 底层的 time.Time 实例
// 示例: m_time.New().Time() // 输出: 等同于 time.Now()
func (t *Time) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.tm
}

// String 实现 Stringer 接口
// 返回: string - 时间的字符串表示
// 示例: m_time.New().String() // 输出: 2023-10-15 13:45:26.123456789 +0800 CST
func (t *Time) String() string {
	if t == nil {
		return "<nil>"
	}
	return t.tm.String()
}

// Before 检查当前时间是否在另一个时间之前
// 参数: t2 *Time - 要比较的另一个时间
// 返回: bool - 如果当前时间在t2之前返回true，否则返回false
// 示例: m_time.New().Before(m_time.New().Add(24 * time.Hour)) // 返回: true
func (t *Time) Before(t2 *Time) bool {
	if t == nil || t2 == nil {
		return false
	}
	return t.tm.Before(t2.tm)
}

// After 检查当前时间是否在另一个时间之后
// 参数: t2 *Time - 要比较的另一个时间
// 返回: bool - 如果当前时间在t2之后返回true，否则返回false
// 示例: m_time.New().After(m_time.New().Subtract(24 * time.Hour)) // 返回: true
func (t *Time) After(t2 *Time) bool {
	if t == nil || t2 == nil {
		return false
	}
	return t.tm.After(t2.tm)
}

// Equal 检查当前时间是否与另一个时间相等
// 参数: t2 *Time - 要比较的另一个时间
// 返回: bool - 如果当前时间与t2相等返回true，否则返回false
// 示例: now := m_time.New(); now.Equal(now) // 返回: true
func (t *Time) Equal(t2 *Time) bool {
	if t == nil || t2 == nil {
		return t == t2
	}
	return t.tm.Equal(t2.tm)
}

// Weekday 获取当前时间是星期几
// 返回: time.Weekday - 星期几（0-6，周日为0）
// 示例: m_time.New().Weekday() // 如周一返回 time.Monday
func (t *Time) Weekday() time.Weekday {
	if t == nil {
		return time.Sunday
	}
	return t.tm.Weekday()
}

// StartOfWeek 获取本周的开始时间（周一0点）
// 返回: *Time - 本周开始时间
// 示例: m_time.New().StartOfWeek().Format("2006-01-02 15:04:05")
func (t *Time) StartOfWeek() *Time {
	if t == nil {
		return nil
	}
	wd := t.tm.Weekday()
	// 调整为周一为一周的开始
	if wd == time.Sunday {
		wd = 6
	} else {
		wd--
	}
	daysToSubtract := time.Duration(wd) * 24 * time.Hour
	return &Time{tm: t.tm.Add(-daysToSubtract).Truncate(24 * time.Hour)}
}

// Parse 按照指定格式解析时间字符串
// 参数: layout string - 时间格式
//
//	s string - 时间字符串
//
// 返回: (*Time, error) - 解析后的Time实例和可能的错误
// 示例: m_time.Parse("2006-01-02", "2023-10-15")
func Parse(layout, s string) (*Time, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return nil, err
	}
	return &Time{tm: t}, nil
}

// ParseInLocation 按照指定格式和时区解析时间字符串
// 参数: layout string - 时间格式
//
//	s string - 时间字符串
//	loc *time.Location - 时区
//
// 返回: (*Time, error) - 解析后的Time实例和可能的错误
// 示例: loc, _ := time.LoadLocation("Asia/Shanghai"); m_time.ParseInLocation("2006-01-02", "2023-10-15", loc)
func ParseInLocation(layout, s string, loc *time.Location) (*Time, error) {
	t, err := time.ParseInLocation(layout, s, loc)
	if err != nil {
		return nil, err
	}
	return &Time{tm: t}, nil
}

// maxInt 返回两个整数中的较大值
// 小工具：避免除以 0
// 参数: a, b int - 要比较的两个整数
// 返回: int - 较大的整数
// 示例: maxInt(3, 5) // 返回: 5
//
//	maxInt(-1, -5) // 返回: -1
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
