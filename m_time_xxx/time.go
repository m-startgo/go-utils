package m_time

// m_time 提供几个常用的时间工具：解析（支持常见字符串与时间戳）、格式化
// 及一些便捷方法（开始/结束时间、加减天数/小时等）。默认尽量基于标准库实现。

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/araddon/dateparse"
)

// 默认的 token 风格格式
const DefaultToken = "YYYY-MM-DD HH:mm:ss.SSSSSS ±HH:MM"

// Time 是对 time.Time 的轻量封装，便于链式调用
type Time struct {
	t time.Time
}

// DefaultLocation 控制数字时间戳解析后使用的时区（默认 UTC）。
// 如果想使用本地时区，调用 SetDefaultLocation(time.Local)。
var defaultLoc atomic.Value

func init() {
	defaultLoc.Store(time.UTC)
}

// DefaultLocation 返回当前默认时区（用于数值时间戳解析），默认 UTC。
func DefaultLocation() *time.Location {
	v := defaultLoc.Load()
	if loc, ok := v.(*time.Location); ok {
		return loc
	}
	return time.UTC
}

// SetDefaultLocation 设置数字时间戳解析后的默认时区，传入 nil 恢复为 UTC。
func SetDefaultLocation(loc *time.Location) {
	if loc == nil {
		defaultLoc.Store(time.UTC)
		return
	}
	defaultLoc.Store(loc)
}

// Parse 接受多种输入并尝试解析为时间：
// 支持：
// - string（任意可被 dateparse 识别的格式）
// - 数字字符串（表示秒/毫秒/微秒/纳秒 时间戳）
// - 整数类型/无符号整数/浮点数（表示秒或带小数的秒）
// 示例：
//
//	Parse("2020-01-02 15:04:05")
//	Parse(1609459200000)           // 13 位毫秒时间戳
//	Parse("1609459200000")       // 数字字符串（毫秒）
//	Parse(1609459200.123)         // 带小数的秒
//
// ParseString 解析字符串格式的时间，优先将纯数字字符串视为时间戳
func ParseString(s string) (Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Time{}, fmt.Errorf("empty time string")
	}
	if isNumericString(s) {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return Time{t: unixFromInt64(i)}, nil
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return Time{t: timeFromFloatSeconds(f)}, nil
		}
	}
	tt, err := dateparse.ParseAny(s)
	if err != nil {
		return Time{}, err
	}
	return Time{t: tt}, nil
}

// ParseInt64 按整数时间戳解析（自动判断秒/毫秒/微秒/纳秒）
func ParseInt64(n int64) (Time, error) { return Time{t: unixFromInt64(n)}, nil }

// ParseFloat64 按带小数的秒解析为时间
func ParseFloat64(f float64) (Time, error) { return Time{t: timeFromFloatSeconds(f)}, nil }

// Parse 接受任意类型输入并路由到更具体的解析函数（兼容旧 API）
func Parse(v any) (Time, error) {
	if v == nil {
		return Time{}, fmt.Errorf("nil value")
	}
	switch x := v.(type) {
	case string:
		return ParseString(x)
	case int:
		return ParseInt64(int64(x))
	case int8:
		return ParseInt64(int64(x))
	case int16:
		return ParseInt64(int64(x))
	case int32:
		return ParseInt64(int64(x))
	case int64:
		return ParseInt64(x)
	case uint:
		return ParseInt64(int64(x))
	case uint8:
		return ParseInt64(int64(x))
	case uint16:
		return ParseInt64(int64(x))
	case uint32:
		return ParseInt64(int64(x))
	case uint64:
		if x <= math.MaxInt64 {
			return ParseInt64(int64(x))
		}
		return Time{}, fmt.Errorf("uint64 value too large")
	case float32:
		return ParseFloat64(float64(x))
	case float64:
		return ParseFloat64(x)
	default:
		s := fmt.Sprintf("%v", v)
		if isNumericString(s) {
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return ParseInt64(i)
			}
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return ParseFloat64(f)
			}
		}
		return Time{}, fmt.Errorf("unsupported parse type: %T", v)
	}
}

// MustParse 解析失败时返回零值（不再 panic），保留旧名称以兼容调用者。
// 推荐在需要检查错误的场景使用 Parse。
func MustParse(v any) Time {
	t, _ := Parse(v)
	return t
}

// isNumericString 判断字符串是否只包含数字（允许前导 -）
func isNumericString(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	dot := false
	digits := 0
	for _, c := range s {
		if c == '.' {
			if dot {
				return false
			}
			dot = true
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
		digits++
	}
	return digits > 0
}

// unixFromInt64 根据数字长度判断单位并返回 time.Time
func unixFromInt64(n int64) time.Time {
	// 通过数字位数判断： >=18 纳秒, >=16 微秒, >=13 毫秒, else 秒
	abs := n
	if abs < 0 {
		abs = -abs
	}
	d := len(strconv.FormatInt(abs, 10))
	var tt time.Time
	switch {
	case d >= 18:
		// 纳秒
		tt = time.Unix(0, n)
	case d >= 16:
		// 微秒 -> 转为纳秒
		tt = time.Unix(0, n*1000)
	case d >= 13:
		// 毫秒
		sec := n / 1000
		msec := n % 1000
		tt = time.Unix(sec, msec*1e6)
	default:
		// 秒
		tt = time.Unix(n, 0)
	}
	loc := DefaultLocation()
	if loc == time.UTC {
		return tt.UTC()
	}
	return tt.In(loc)
}

// timeFromFloatSeconds 将带小数的秒转换为 time.Time
func timeFromFloatSeconds(f float64) time.Time {
	sec := int64(math.Floor(f))
	frac := f - float64(sec)
	nsec := int64(math.Round(frac * 1e9))
	tt := time.Unix(sec, nsec)
	loc := DefaultLocation()
	if loc == time.UTC {
		return tt.UTC()
	}
	return tt.In(loc)
}

// FromTime 包装一个标准 time.Time
func FromTime(tt time.Time) Time { return Time{t: tt} }

// ToTime 返回底层 time.Time
func (t Time) ToTime() time.Time { return t.t }

// Format 支持两种形式：
// - 传入空字符串或 DefaultToken 则使用默认的 token
// - 传入类似 "YYYY-MM-DD HH:mm:ss" 的 token，会被映射为 go layout
func (t Time) Format(token string) string {
	if token == "" {
		token = DefaultToken
	}
	layout := tokenToLayout(token)
	return t.t.Format(layout)
}

// String 实现 fmt.Stringer，使用默认格式
func (t Time) String() string { return t.Format(DefaultToken) }

// Add 增加 duration
func (t Time) Add(d time.Duration) Time { return Time{t: t.t.Add(d)} }

// AddDays 快捷方法
func (t Time) AddDays(days int) Time { return Time{t: t.t.Add(time.Duration(days) * 24 * time.Hour)} }

// AddHours 快捷方法
func (t Time) AddHours(h int) Time { return Time{t: t.t.Add(time.Duration(h) * time.Hour)} }

// StartOfDay 返回当天零点
func (t Time) StartOfDay() Time {
	y, m, d := t.t.Date()
	loc := t.t.Location()
	return Time{t: time.Date(y, m, d, 0, 0, 0, 0, loc)}
}

// EndOfDay 返回当天 23:59:59.999999
func (t Time) EndOfDay() Time {
	y, m, d := t.t.Date()
	loc := t.t.Location()
	const endOfDayNsec = 999999 * int(time.Microsecond)
	return Time{t: time.Date(y, m, d, 23, 59, 59, endOfDayNsec, loc)}
}

// IsSameDay 判断是否同一天（本地时区）
func (t Time) IsSameDay(o Time) bool {
	y1, m1, d1 := t.t.Date()
	y2, m2, d2 := o.t.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// Diff 返回 t - o 的 duration
func (t Time) Diff(o Time) time.Duration { return t.t.Sub(o.t) }

// Unix 返回秒级时间戳
func (t Time) Unix() int64 { return t.t.Unix() }

// UnixNano 返回纳秒时间戳
func (t Time) UnixNano() int64 { return t.t.UnixNano() }

// UTC 转换为 UTC
func (t Time) UTC() Time { return Time{t: t.t.UTC()} }

// Local 转换为本地时间
func (t Time) Local() Time { return Time{t: t.t.Local()} }

// tokenToLayout 将常见 token 映射为 go time layout，支持部分 dayjs 风格 token
func tokenToLayout(token string) string {
	// 优先处理 SSSSSS
	r := token
	// 保证 ±HH:MM 映射为 -07:00
	r = strings.ReplaceAll(r, "±HH:MM", "-07:00")
	// 有顺序地替换其他 token
	replacements := []struct{ old, new string }{
		{"YYYY", "2006"},
		{"MM", "01"},
		{"DD", "02"},
		{"HH", "15"},
		{"mm", "04"},
		{"ss", "05"},
		{"SSSSSS", "000000"},
		{"SSS", "000"},
	}
	for _, rp := range replacements {
		r = strings.ReplaceAll(r, rp.old, rp.new)
	}
	return r
}

// MillisPer... 常量表示毫秒级别的时间长度，命名更语义化且便于直接使用。
const (
	MillisPerSecond int64 = 1000
	MillisPerMinute int64 = 60 * MillisPerSecond
	MillisPerHour   int64 = 60 * MillisPerMinute
	MillisPerDay    int64 = 24 * MillisPerHour
)
