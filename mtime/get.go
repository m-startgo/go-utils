package mtime

import (
	"time"
)

// Now 返回当前时间的封装
func Now() MTime {
	return MTime{t: time.Now()}
}

// NowUnixMilli 返回当前时间的毫秒时间戳（13 位）
func NowUnixMilli() int64 { return time.Now().UnixNano() / 1e6 }

// NowDefaultString 直接返回当前时间的默认格式字符串 "YYYY-MM-DDTHH:mm:ss"
// 例如: 2020-01-02T15:04:05
func NowDefaultString() string {
	return Now().FormatDefault()
}

// FormatDefault 返回默认的无参数格式化，格式为 "YYYY-MM-DDTHH:mm:ss"
// 例如: 2020-01-02T15:04:05
func (t MTime) FormatDefault() string {
	return t.Format("YYYY-MM-DDTHH:mm:ss")
}

// UnixMilli 返回以毫秒为单位的时间戳（13 位）
func (t MTime) UnixMilli() int64 {
	return t.t.UnixNano() / 1e6
}

// ParseToTimeWithMillisOffset 解析任意支持的输入为 time.Time，并在结果上加上以毫秒为单位的偏移量。
// 参数:
// - v: 支持 Parse 接受的任意类型（string/数字/浮点等）。
// - offsetMillis: 以毫秒为单位的偏移量，可以为负、正或 0。
// 返回值:
// - 解析并加上偏移后的 time.Time。
// 行为说明:
// - 解析失败时为保持向后兼容，返回 Parse(0) 的结果（即 epoch 对应的 time.Time）。
func ParseToTimeWithMillisOffset(v any, offsetMillis int64) time.Time {
	tm, err := Parse(v)
	if err != nil {
		n, _ := Parse(0)
		return n.ToTime()
	}
	ms := tm.UnixMilli() + offsetMillis
	return time.UnixMilli(ms)
}

// FormatDefaultFrom 将任意支持的输入解析并格式化为默认时间字符串 "YYYY-MM-DDTHH:mm:ss"。
// 解析失败时返回 epoch 的默认格式化字符串以保持向后兼容。
func FormatDefaultFrom(v any) string {
	tm, err := Parse(v)
	if err != nil {
		n, _ := Parse(0)
		return n.FormatDefault()
	}
	return tm.FormatDefault()
}
