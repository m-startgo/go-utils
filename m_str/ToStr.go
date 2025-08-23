package m_str

import "fmt"

func ToStr(p any) string {
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
