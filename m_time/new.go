package m_time

import (
	"time"

	"github.com/araddon/dateparse"
)

// Time 结构体，用于封装 time.Time
type Time struct {
	tm time.Time // 避免与包名冲突
}

// New 创建一个新的 Time 实例
// 返回: *Time - 包含当前时间的新实例
func New() *Time {
	return &Time{tm: time.Now()}
}

// NewFromTime 从 time.Time 创建一个新的 Time 实例
// 参数: t time.Time - 基础时间
// 返回: *Time - 新的 Time 实例
func NewFromTime(t time.Time) *Time {
	return &Time{tm: t}
}

// NewFromString 从字符串创建一个新的 Time 实例
// 参数: s string - 时间字符串
// 返回: (*Time, error) - 新的 Time 实例和可能的错误
func NewFromString(s string) (*Time, error) {
	tp, err := dateparse.ParseAny(s)
	if err != nil {
		return nil, err
	}
	return &Time{tm: tp}, nil
}
