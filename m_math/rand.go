package m_math

import (
	"math/rand"
	"time"
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandIntN 返回 [0, n) 的随机整数；当 n <= 0 时返回 0
func RandIntN(n int) int {
	if n <= 0 {
		return 0
	}
	return rnd.Intn(n)
}

// RandInt 返回闭区间 [min, max] 的随机整数；若 min > max 则会交换两者
func RandInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}
	return min + rnd.Intn(max-min+1)
}
