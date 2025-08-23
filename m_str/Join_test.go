package m_str

import (
	"testing"
)

// go test -v -run Test_Str_Join

func Test_Str_Join(t *testing.T) {
	a := []int32{1, 2, 3, 4, 5}
	joinStr := Join("mo7", "欢迎你", a, "张三")

	t.Log("joinStr", joinStr)
}
