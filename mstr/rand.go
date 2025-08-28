package mstr

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
)

var baseStr = "0123456789abcdefghijklmnopqrstuvwxyz"

// Rand 生成一个 length 长度的随机字符串。
//
// 示例：
//
//	s := Rand(8) // 返回类似 "a1b2c3d4" 的字符串
//
// 特殊情况：
//   - 当 length <= 0 时，返回空字符串
//
// 说明：此函数使用 crypto/rand 以避免跨进程/重启时伪随机序列重复，适用于需要更高随机性的场景。
func Rand(length int) string {
	if length <= 0 {
		return ""
	}
	l := int64(len(baseStr))
	if l == 0 {
		return ""
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := crand.Int(crand.Reader, big.NewInt(l))
		if err != nil {
			// crypto/rand 读取失败时回退到 math/rand 全局实现
			idx := rand.Intn(int(l))
			bytes[i] = baseStr[idx]
			continue
		}
		bytes[i] = baseStr[n.Int64()]
	}
	return string(bytes)
}
