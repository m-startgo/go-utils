package m_time

/*

我想封装一个 m_time 包，提供一些时间处理的工具函数

并且使用 araddon/dateparse 库来解析日期字符串。
请帮我写出这个包的代码。

代码一定要精简，尽量依赖go标准库去实现。

我不太喜欢 go 语言这样的   "2006-01-02 15:04:05.000 -07:00"
我更加喜欢 YYYY-MM-DD HH:mm:ss.SSSSSS ±HH:MM 这样的方式，你可以封装处理一下


代码一定要精简，基于官方标准库。运行速度要快。
Api 可以参考 dayjs  一些深度的API无需支持，只要涉及到常用的 一些方法 就行。


*/

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// 默认的 token 风格格式
const DefaultToken = "YYYY-MM-DD HH:mm:ss.SSSSSS ±HH:MM"

// Time 是对 time.Time 的轻量封装，便于链式调用
type Time struct {
	t time.Time
}

// Parse 使用 araddon/dateparse 尝试解析任意常见日期字符串
func Parse(s string) (Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Time{}, fmt.Errorf("empty time string")
	}
	tt, err := dateparse.ParseAny(s)
	if err != nil {
		return Time{}, err
	}
	return Time{t: tt}, nil
}

// MustParse 解析失败时 panic，方便测试或内联使用
func MustParse(s string) Time {
	t, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}

// Now 返回当前时间的封装
func Now() Time { return Time{t: time.Now()} }

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
	return Time{t: time.Date(y, m, d, 23, 59, 59, int(time.Microsecond*999999), loc)}
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
