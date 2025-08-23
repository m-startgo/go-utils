package m_str

import (
	"strings"
)

/*
拼接字符串

var a = []int32{1, 2, 3, 4, 5}
joinStr := m_str.Join("mo7", "欢迎你", a, "张三")

fmt.Println("joinStr", joinStr)
*/
func Join(s ...any) string {
	var build strings.Builder
	// 预分配容量以提高性能
	for _, v := range s {
		build.WriteString(ToStr(v))
	}
	return build.String()
}
