package m_str

import (
	"os"
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
	for _, v := range s {
		build.WriteString(ToStr(v))
	}
	return build.String()
}

/*
字符串模板

	var config = `

app.name = ${appName}
app.ip = ${appIP}
app.port = ${appPort}
`

	var dev = map[string]string{
		"appName": "my_ap123p",
		"appIP":   "0.0.0.0",
		"appPort": "8080",
	}
	s := m_str.Temp(config, dev)

	fmt.Println("temp", s)
*/
func Temp(tmplStr string, lMap map[string]string) string {
	s := os.Expand(tmplStr, func(k string) string {
		return lMap[k]
	})
	return s
}
