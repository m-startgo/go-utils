package m_time

import (
	"time"

	"github.com/araddon/dateparse"
)

// Time 结构体，用于封装 time.Time
type Time struct {
	// 将字段名从 `time` 改为 `tm`，避免与包名冲突并提高可读性
	tm time.Time
}

// New 创建一个新的 Time 实例
func New() *Time {
	return &Time{tm: time.Now()}
}

// NewFromTime 从 time.Time 创建一个新的 Time 实例
func NewFromTime(t time.Time) *Time {
	return &Time{tm: t}
}

// NewFromString 从字符串创建一个新的 Time 实例
func NewFromString(s string) (*Time, error) {
	tp, err := dateparse.ParseAny(s)
	if err != nil {
		return nil, err
	}
	return &Time{tm: tp}, nil
}
