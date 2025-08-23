package m_str

import "os"

func TplFormat(tplStr string, MapData map[string]string) string {
	s := os.Expand(tplStr, func(k string) string {
		return MapData[k]
	})
	return s
}
