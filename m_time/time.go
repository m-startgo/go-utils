package m_time

import (
	"time"

	"github.com/araddon/dateparse"
)

// m_time 包：时间工具，API 设计参考 dayjs，使用 dateparse 解析日期字符串。

// Time 封装 time.Time
type Time struct{ tm time.Time }

// New 返回当前时间的 *Time
func New() *Time { return &Time{tm: time.Now()} }

// NewFromTime 用已有 time.Time 创建 *Time
func NewFromTime(t time.Time) *Time { return &Time{tm: t} }

// NewFromString 解析字符串生成 *Time（使用 dateparse.ParseAny）
func NewFromString(s string) (*Time, error) {
	tp, err := dateparse.ParseAny(s)
	if err != nil {
		return nil, err
	}
	return &Time{tm: tp}, nil
}
