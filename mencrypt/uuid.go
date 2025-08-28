package mencrypt

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-startgo/go-utils/mstr"
)

// UUID 生成一个 RFC4122 v4 随机 UUID 并返回字符串表示（小写，带连字符）。
//
// 示例：
//
//	id := UUID()
//	// id -> "550e8400-e29b-41d4-a716-446655440000"
//
// 说明：此函数是对第三方库的轻量封装，调用非常廉价且不会返回错误。
func UUID() string {
	return uuid.New().String()
}

// 生成一个可读的 Time ID，基于当前时间戳和随机数。
func TimeID() string {
	// 格式：YYYYMMDDhhmmss.mmm-xxxxxxxx
	// 例如：20250828T150405.123-1a2b3c4d （返回中不包含字母 T，这里使用紧凑数字格式）
	t := time.Now().Format("2006-01-02-15-04-05-.000")
	// 去掉.
	t = strings.ReplaceAll(t, ".", "")

	return t + "-" + mstr.Rand(8)
}
