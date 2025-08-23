package m_str

import "os"

// TplFormat 使用给定的数据替换模板中的占位符
//
// 参数:
//
//	tplStr - 模板字符串  以 ${key} 的形式表示占位符
//	MapData - 用于替换占位符的数据映射
//
// 返回值:
//
//	string - 替换占位符后的字符串

func TplFormat(tplStr string, MapData map[string]string) string {
	s := os.Expand(tplStr, func(k string) string {
		// 检查键是否存在，如果不存在则返回原占位符
		if val, ok := MapData[k]; ok {
			return val
		}
		return "" // 如果键不存在，返回空字符串
	})
	return s
}
