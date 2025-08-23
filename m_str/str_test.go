package m_str

import (
	"testing"
)

// go test -v -run Test_Join

func Test_Join(t *testing.T) {
	a := []int32{1, 2, 3, 4, 5}
	joinStr := Join("mo7", "欢迎你", a, "张三")
	t.Log("joinStr", joinStr)
}

// go test -v -run Test_ToStr
func Test_ToStr(t *testing.T) {
	// 测试 []rune 转字符串
	a := []rune("mo7欢迎你")
	t.Log(ToStr(a))

	// 测试 []byte 转字符串
	b := []byte("mo7欢迎你")
	t.Log(ToStr(b))

	// 测试 float64 转字符串
	c := 10.97
	t.Log(ToStr(c))

	// 测试 int 转字符串
	d := 100
	t.Log(ToStr(d))

	// 测试 string 转字符串
	e := "hello"
	t.Log(ToStr(e))

	// 测试 bool 转字符串
	f := true
	t.Log(ToStr(f))

	// 测试 map 转字符串
	g := map[string]int{"one": 1, "two": 2}
	t.Log(ToStr(g))

	// 测试 struct 转字符串
	type Person struct {
		Name string
		Age  int
	}
	h := Person{Name: "Alice", Age: 30}
	t.Log(ToStr(h))

	// 测试 nil 转字符串
	var i any = nil
	t.Log(ToStr(i))

	// 测试 数组 转字符串
	j := [3]int{1, 2, 3}
	t.Log(ToStr(j))

	// 测试空值
	var k any
	t.Log(ToStr(k))

	// 测试负数
	l := -42
	t.Log(ToStr(l))

	// 测试浮点数
	m := 3.14159
	t.Log(ToStr(m))

	n := []int32{1, 2, 3, 4, 5}
	t.Log(ToStr(n))

	o := []int32{97}
	t.Log(ToStr(o))
}

// go test -v -run Test_TplFormat
func Test_TplFormat(t *testing.T) {
	tpl := `
app.name = ${appName}
app.ip = ${appIP}
app.port = ${appPort}
`
	data := map[string]string{
		"appName": "my_ap123p",
		"appIP":   "0.0.0.0",
	}

	s := TplFormat(tpl, data)
	t.Log("temp", s)
}
