package m_str

import "fmt"

/*
任意类型转字符串

a := []rune("mo7欢迎你")
a := []byte("mo7欢迎你")
a := 10.97
a := os.PathSeparator
str := m_str.ToStr(a)
*/
func ToStr(p any) string {
	// fmt.Println("type: ", reflect.TypeOf(p))
	returnStr := ""
	switch p := p.(type) {
	case []int32:
		returnStr = string(p)
	case []uint8:
		returnStr = string(p)
	case int32:
		returnStr = string(p)
	case uint8:
		returnStr = string(p)
	default:
		returnStr = fmt.Sprintf("%+v", p)
	}

	return returnStr
}
