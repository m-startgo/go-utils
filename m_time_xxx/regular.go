package m_time

import "time"

// Now 返回当前时间的封装
func Now() Time {
	return Time{t: time.Now()}
}

// NowUnixMilli 返回当前时间的毫秒时间戳（13 位）
func NowUnixMilli() int64 { return time.Now().UnixNano() / 1e6 }

// FormatDefault 返回默认的无参数格式化，格式为 "YYYY-MM-DDTHH:mm:ss"
// 例如: 2020-01-02T15:04:05
func (t Time) FormatDefault() string {
	return t.Format("YYYY-MM-DDTHH:mm:ss")
}

// NowDefaultString 直接返回当前时间的默认格式字符串 "YYYY-MM-DDTHH:mm:ss"
// 例如: 2020-01-02T15:04:05
func NowDefaultString() string {
	return Now().FormatDefault()
}

// UnixMilli 返回以毫秒为单位的时间戳（13 位）
func (t Time) UnixMilli() int64 {
	return t.t.UnixNano() / 1e6
}
