package mstr

import (
	"math/rand"
)

var baseStr = "0123456789abcdefghijklmnopqrstuvwxyz"

// RandStr 生成一个 length 长度的随机字符串。
//
// 示例：
//
//	s := RandStr(8) // 返回类似 "a1b2c3d4" 的字符串
//
// 特殊情况：
//   - 当 length <= 0 时，返回空字符串
//
// 说明：此函数使用 math/rand 包的全局随机源，适用于一般用途（非加密场景）。
func Rand(length int) string {
	if length <= 0 {
		return ""
	}
	l := len(baseStr)
	if l == 0 {
		return ""
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = baseStr[rand.Intn(l)]
	}
	return string(bytes)
}
